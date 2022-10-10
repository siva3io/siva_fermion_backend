package sales_invoice

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/accounting"
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
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

type Service interface {
	CreateSalesInvoice(data *SalesInvoiceRequest, token_id string, access_template_id string) (uint, error)
	BulkCreateSalesInvoice(data *[]SalesInvoiceRequest, token_id string, access_template_id string) error
	UpdateSalesInvoice(id uint, data *SalesInvoiceRequest, token_id string, access_template_id string) error
	GetSalesInvoice(id uint, token_id string, access_template_id string) (interface{}, error)
	GetAllSalesInvoice(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]accounting.SalesInvoice, error)
	DeleteSalesInvoice(id uint, user_id uint, token_id string, access_template_id string) error
	DeleteSalesInvoiceLines(query interface{}, token_id string, access_template_id string) error

	SendMailSalesInvoice(q *SendMailSalesInvoice, token_id string, access_template_id string) error
	GetSalesInvoiceTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	salesInvoiceRepository accounting_repo.SalesInvoice
	creditNoteRepo         accounting_repo.CreditNotes
	doRepository           orders_repo.DeliveryOrders
	salesReturnRepository  returns_repo.SalesReturn
}

func NewService() *service {
	SalesInvoiceRepository := accounting_repo.NewSalesInvoice()
	return &service{SalesInvoiceRepository, accounting_repo.NewCreditNote(), orders_repo.NewDo(), returns_repo.NewSalesReturn()}
}

func (s *service) CreateSalesInvoice(data *SalesInvoiceRequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SALES_INVOICE", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create sales invoice at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create sales invoice at data level")
	}
	var SalesInvoiceData accounting.SalesInvoice
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &SalesInvoiceData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	if data.AutoGenerateInvoiceNumber {
		SalesInvoiceData.SalesInvoiceNumber = helpers.GenerateSequence("SALESINVOICE", token_id, "sales_invoices")
	}
	if data.AutoGenerateReferenceNumber {
		SalesInvoiceData.ReferenceNumber = helpers.GenerateSequence("REF", token_id, "sales_invoices")
	}
	defaultStatus, err := helpers.GetLookupcodeId("SALES_INVOICE_STATUS", "DRAFT")
	if err != nil {
		return 0, err
	}
	data.StatusID = defaultStatus

	id, err := s.salesInvoiceRepository.CreateSalesInvoice(&SalesInvoiceData, token_id)
	return id, err
}

func (s *service) BulkCreateSalesInvoice(data *[]SalesInvoiceRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SALES_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create sales invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create sales invoice at data level")
	}
	bulk_data := []map[string]interface{}{}
	for index, value := range *data {
		v := map[string]interface{}{}
		count := helpers.GetCount("SELECT COUNT(*) FROM sales_invoices") + 1 + index
		if value.AutoGenerateInvoiceNumber {
			value.SalesInvoiceNumber = "SALESINVOICE/" + token_id + "/000" + strconv.Itoa(count)
		}
		if value.AutoGenerateReferenceNumber {
			value.ReferenceNumber = "REF/" + token_id + "/000" + strconv.Itoa(count)
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

	err = s.salesInvoiceRepository.BulkCreateSalesInvoice(&SalesInvoiceData, token_id)
	return err
}

func (s *service) UpdateSalesInvoice(id uint, data *SalesInvoiceRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "SALES_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update sales invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update sales invoice at data level")
	}
	var SalesInvoiceData accounting.SalesInvoice
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &SalesInvoiceData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	old_data, er := s.salesInvoiceRepository.GetSalesInvoice(id)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := SalesInvoiceData.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, SalesInvoiceData.StatusID)
		SalesInvoiceData.StatusHistory = result
	}

	err = s.salesInvoiceRepository.UpdateSalesInvoice(id, &SalesInvoiceData)
	for _, order_line := range SalesInvoiceData.SalesInvoiceLines {
		query := map[string]interface{}{
			"sales_invoice_id": uint(id),
			"product_id":       uint(order_line.ProductID),
		}
		count, er := s.salesInvoiceRepository.UpdateSalesInvoiceLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.SalesInvoiceID = id
			e := s.salesInvoiceRepository.CreateSalesInvoiceLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *service) GetSalesInvoice(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_INVOICE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales invoice at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales invoice at data level")
	}
	result, er := s.salesInvoiceRepository.GetSalesInvoice(id)
	if er != nil {
		return result, er
	}
	query := map[string]interface{}{
		"sales_invoice_id": id,
	}
	result_order_lines, err := s.salesInvoiceRepository.GetSalesInvoiceLines(query)
	result.SalesInvoiceLines = result_order_lines
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) GetAllSalesInvoice(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]accounting.SalesInvoice, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "SALES_INVOICE", *token_user_id)
	fmt.Println("access data", access_module_flag, data_access)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales invoice at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales invoice at data level")
	}
	result, err := s.salesInvoiceRepository.GetAllSalesInvoice(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteSalesInvoice(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SALES_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales invoice at data level")
	}
	_, er := s.salesInvoiceRepository.GetSalesInvoice(id)
	if er != nil {
		return er
	}
	err := s.salesInvoiceRepository.DeleteSalesInvoice(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"sales_invoice_id": id}
	err1 := s.salesInvoiceRepository.DeleteSalesInvoiceLines(query)
	return err1
}

func (s *service) DeleteSalesInvoiceLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SALES_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales invoice at data level")
	}
	data, err := s.salesInvoiceRepository.GetSalesInvoiceLines(query)
	if err != nil {
		return err
	}
	if len(data) <= 0 {
		return err
	}
	err = s.salesInvoiceRepository.DeleteSalesInvoiceLines(query)
	return err
}

func (s *service) SendMailSalesInvoice(q *SendMailSalesInvoice, token_id string, access_template_id string) error {
	id, _ := strconv.Atoi(q.ID)
	result, er := s.salesInvoiceRepository.GetSalesInvoice(uint(id))
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Sales Invoice", q.ReceiverEmail, "pkg/util/response/static/accounting/sales_invoice_template.html", result)

	return err
}

func (s *service) GetSalesInvoiceTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales order at data level")
	}

	deliveryOrderId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "credit_note" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_INVOICE")

		deliveryOrdersPage := page
		fmt.Println(sourceDocumentId, deliveryOrderId)
		deliveryOrdersPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		fmt.Println(deliveryOrdersPage)
		query := make(map[string]interface{}, 0)
		data, err := s.creditNoteRepo.FindAllCreditNote(query, deliveryOrdersPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "delivery_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_INVOICE")

		deliveryOrdersPage := page
		deliveryOrdersPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		data, err := s.doRepository.AllDeliveryOrders(deliveryOrdersPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "sales_returns" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "SALES_INVOICE")

		salesReturnsPage := page
		salesReturnsPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		data, err := s.salesReturnRepository.FindAll(salesReturnsPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
