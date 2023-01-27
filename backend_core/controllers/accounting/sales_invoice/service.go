package sales_invoice

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	accounting_repo "fermion/backend_core/internal/repository/accounting"
	orders_repo "fermion/backend_core/internal/repository/orders"
	returns_repo "fermion/backend_core/internal/repository/returns"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"
)

/*
 Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Lesser General Public License v3.0 as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Lesser General Public License v3.0 for more details.
 You should have received a copy of the GNU Lesser General Public License v3.0
 along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/

type Service interface {
	CreateSalesInvoice(metaData core.MetaData, data *accounting.SalesInvoice) error
	BulkCreateSalesInvoice(metaData core.MetaData, data *[]accounting.SalesInvoice) error
	UpdateSalesInvoice(metaData core.MetaData, data *accounting.SalesInvoice) error
	GetSalesInvoicePdf(query map[string]interface{}) (accounting.SalesInvoice, error)
	GetSalesInvoice(metaData core.MetaData) (interface{}, error)
	GetAllSalesInvoice(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeleteSalesInvoice(metaData core.MetaData) error
	DeleteSalesInvoiceLines(metaData core.MetaData) error
	ProcessTaxCalculation(*accounting.SalesInvoice)

	SendMailSalesInvoice(metaData core.MetaData, q *SendMailSalesInvoice) error
	GetSalesInvoiceTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	salesInvoiceRepository accounting_repo.SalesInvoice
	creditNoteRepo         accounting_repo.CreditNotes
	doRepository           orders_repo.DeliveryOrders
	salesReturnRepository  returns_repo.SalesReturn
	SalesOrdersRepository  orders_repo.SalesOrders
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	SalesInvoiceRepository := accounting_repo.NewSalesInvoice()
	newServiceObj = &service{
		SalesInvoiceRepository,
		accounting_repo.NewCreditNote(),
		orders_repo.NewDo(),
		returns_repo.NewSalesReturn(),
		*orders_repo.NewSalesOrder(),
	}
	return newServiceObj
}

func (s *service) CreateSalesInvoice(metaData core.MetaData, data *accounting.SalesInvoice) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "SALES_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create sales invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create sales invoice at data level")
	}
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId

	var requestData SalesInvoiceRequest
	if requestData.AutoGenerateInvoiceNumber {
		data.SalesInvoiceNumber = helpers.GenerateSequence("SALESINVOICE", fmt.Sprint(metaData.TokenUserId), "sales_invoices")
	}
	if requestData.AutoGenerateReferenceNumber {
		data.ReferenceNumber = helpers.GenerateSequence("REF", fmt.Sprint(metaData.TokenUserId), "sales_invoices")
	}
	defaultStatus, err := helpers.GetLookupcodeId("SALES_INVOICE_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusID = &defaultStatus

	id, err := s.salesInvoiceRepository.CreateSalesInvoice(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	var so_data orders.SalesOrders
	so_data.InvoiceId = id
	err = s.SalesOrdersRepository.Update(map[string]interface{}{"id": metaData.Query["sales_order_id"]}, &so_data)
	if err != nil {
		return err
	}
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	datas, _ := s.GetSalesInvoice(metaData)

	var si_data accounting.SalesInvoice
	json_data, _ := json.Marshal(datas)
	json.Unmarshal(json_data, &si_data)

	s.ProcessTaxCalculation(&si_data)
	s.UpdateSalesInvoice(metaData, &si_data)

	return nil
}

func (s *service) BulkCreateSalesInvoice(metaData core.MetaData, data *[]accounting.SalesInvoice) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "SALES_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create sales invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create sales invoice at data level")
	}
	bulk_data := []map[string]interface{}{}
	var dtoData SalesInvoiceRequest
	for index, value := range *data {
		v := map[string]interface{}{}
		count := helpers.GetCount("SELECT COUNT(*) FROM sales_invoices") + 1 + index
		if dtoData.AutoGenerateInvoiceNumber {
			value.SalesInvoiceNumber = "SALESINVOICE/" + fmt.Sprint(metaData.TokenUserId) + "/000" + strconv.Itoa(count)
		}
		if dtoData.AutoGenerateReferenceNumber {
			value.ReferenceNumber = "REF/" + fmt.Sprint(metaData.TokenUserId) + "/000" + strconv.Itoa(count)
		}

		val, err := json.Marshal(value)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(val, &v)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}

		bulk_data = append(bulk_data, v)
	}
	var SalesInvoiceData []accounting.SalesInvoice
	dto, err := json.Marshal(bulk_data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &SalesInvoiceData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	err = s.salesInvoiceRepository.BulkCreateSalesInvoice(&SalesInvoiceData)
	return err
}

func (s *service) UpdateSalesInvoice(metaData core.MetaData, data *accounting.SalesInvoice) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "SALES_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update sales invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update sales invoice at data level")
	}
	old_data, er := s.salesInvoiceRepository.GetSalesInvoice(metaData.Query)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := data.StatusID
	if new_status != old_status && new_status != nil {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, *data.StatusID)
		data.StatusHistory = result
	}

	err := s.salesInvoiceRepository.UpdateSalesInvoice(metaData.Query, data)
	if err != nil {
		return err
	}
	for _, order_line := range data.SalesInvoiceLines {
		query := map[string]interface{}{
			"sales_invoice_id": old_data.ID,
			"product_id":       order_line.ProductID,
		}
		order_line.UpdatedByID = &metaData.TokenUserId
		count, er := s.salesInvoiceRepository.UpdateSalesInvoiceLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.SalesInvoiceID = order_line.ID
			order_line.CreatedByID = &metaData.TokenUserId
			order_line.CompanyId = metaData.CompanyId
			e := s.salesInvoiceRepository.CreateSalesInvoiceLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *service) GetSalesInvoice(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "SALES_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view sales invoice at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view sales invoice at data level")
	}
	result, er := s.salesInvoiceRepository.GetSalesInvoice(metaData.Query)
	if er != nil {
		return result, er
	}
	// result_order_lines, err := s.salesInvoiceRepository.GetSalesInvoiceLines(metaData.Query)
	// result.SalesInvoiceLines = result_order_lines
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}

func (s *service) GetAllSalesInvoice(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "SALES_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list sales invoice at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list sales invoice at data level")
	}
	result, err := s.salesInvoiceRepository.GetAllSalesInvoice(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteSalesInvoice(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "SALES_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete sales invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete sales invoice at data level")
	}
	err := s.salesInvoiceRepository.DeleteSalesInvoice(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{"sales_invoice_id": metaData.Query["id"]}
		er := s.salesInvoiceRepository.DeleteSalesInvoiceLines(query)
		if er != nil {
			return err
		}
	}
	return nil

}

func (s *service) DeleteSalesInvoiceLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "SALES_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete sales invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete sales invoice at data level")
	}
	err := s.salesInvoiceRepository.DeleteSalesInvoiceLines(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SendMailSalesInvoice(metaData core.MetaData, q *SendMailSalesInvoice) error {
	result, er := s.salesInvoiceRepository.GetSalesInvoice(metaData.Query)
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Sales Invoice", q.ReceiverEmail, "pkg/util/response/static/accounting/sales_invoice_template.html", result)

	return err
}

func (s *service) GetSalesInvoiceTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "DELIVERY_ORDERS", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view sales order at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view sales order at data level")
	}
	tab := metaData.AdditionalFields["tab"].(string)
	salesInvoiceId := metaData.AdditionalFields["id"]

	deliveryOrderId := salesInvoiceId
	if tab == "credit_note" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_INVOICE")

		deliveryOrdersPage := page
		fmt.Println(sourceDocumentId, deliveryOrderId)
		deliveryOrdersPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		fmt.Println(deliveryOrdersPage)
		var metaData core.MetaData
		data, err := s.creditNoteRepo.FindAllCreditNote(metaData.Query, deliveryOrdersPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "delivery_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_INVOICE")

		deliveryOrdersPage := page
		deliveryOrdersPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		data, err := s.doRepository.AllDeliveryOrders(nil, deliveryOrdersPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "sales_returns" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "SALES_INVOICE")

		salesReturnsPage := page
		salesReturnsPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		data, err := s.salesReturnRepository.FindAll(nil, salesReturnsPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

func (s *service) GetSalesInvoicePdf(query map[string]interface{}) (accounting.SalesInvoice, error) {
	result, err := s.salesInvoiceRepository.GetSalesInvoice(query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) ProcessTaxCalculation(data *accounting.SalesInvoice) {

	billing_address := make([]map[string]interface{}, 0)
	shipping_address := make([]map[string]interface{}, 0)

	json.Unmarshal(data.BillingAddress, &billing_address)
	json.Unmarshal(data.ShippingAddress, &shipping_address)

	bill_add := billing_address[0]["state"]
	ship_add := shipping_address[0]["state"]
	data.SubTotalAmount = 0
	var subTotal, tot_igst, tot_cgst, tot_sgst, igst_rate, cgst_rate, sgst_rate float64

	for index, orderLines := range data.SalesInvoiceLines {
		amtwithouttax := float64(orderLines.Price) * float64(orderLines.Quantity)
		data.SalesInvoiceLines[index].TotalAmount = amtwithouttax
		if bill_add != ship_add {

			igst_rate = data.SalesInvoiceLines[index].Product.HSNCodesData.IGSTRate
			tot_igst += igst_rate / 100 * (amtwithouttax)

		} else {

			cgst_rate = data.SalesInvoiceLines[index].Product.HSNCodesData.CGSTRate
			sgst_rate = data.SalesInvoiceLines[index].Product.HSNCodesData.SGSTRate
			tot_sgst += sgst_rate / 100 * (amtwithouttax)
			tot_cgst += cgst_rate / 100 * (amtwithouttax)

		}

		amtwithouttax = amtwithouttax - (float64(orderLines.Discount) / 100 * (amtwithouttax))
		subTotal += amtwithouttax
		data.SalesInvoiceLines[index].TotalAmount = amtwithouttax
		data.SalesInvoiceLines[index].IGSTRate = igst_rate
		data.SalesInvoiceLines[index].SGSTRate = sgst_rate
		data.SalesInvoiceLines[index].CGSTRate = cgst_rate

	}

	data.SubTotalAmount += subTotal
	data.IgstAmt = tot_igst
	data.CgstAmt = tot_cgst
	data.SgstAmt = tot_sgst
	data.TotalAmount = tot_cgst + tot_igst + tot_sgst + data.SubTotalAmount

}
