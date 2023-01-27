package sales_returns

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/creditnote"
	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
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
	CreateSalesReturn(metaData core.MetaData, data *returns.SalesReturns) error
	ListSalesReturns(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ViewSalesReturn(metaData core.MetaData) (interface{}, error)
	UpdateSalesReturn(metaData core.MetaData, data *returns.SalesReturns) error
	DeleteSalesReturn(metaData core.MetaData) error
	DeleteSalesReturnLines(metaData core.MetaData) error
	SearchSalesReturns(query string, access_template_id string, token_id string) (interface{}, error)
	GetSalesReturnsHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	GetSalesReturnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ProcessSalesReturnCalculation(data *returns.SalesReturns) *returns.SalesReturns
	TrackSalesReturn(order_id string) map[string]interface{}
}

type service struct {
	salesReturnRepository    returns_repo.SalesReturn
	basicInventoryService    inv_service.Service
	creditNoteService        creditnote.Service
	SalesOrdersRepository    orders_repo.Sales
	DeliveryOrdersRepository orders_repo.DeliveryOrders
	SalesInvoiceRepository   accounting_repo.SalesInvoice
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	salesReturnRepository := returns_repo.NewSalesReturn()
	basicInventoryService := inv_service.NewService()
	newServiceObj = &service{salesReturnRepository, basicInventoryService,
		creditnote.NewService(), orders_repo.NewSalesOrder(), orders_repo.NewDo(), accounting_repo.NewSalesInvoice()}
	return newServiceObj
}

func (s *service) CreateSalesReturn(metaData core.MetaData, data *returns.SalesReturns) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "CREATE", "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create sales return at data level")
	}
	data = s.ProcessSalesReturnCalculation(data)
	var srData *returns.SalesReturns
	helpers.JsonMarshaller(data, &srData)
	//srData.StatusHistory, _ = helpers.UpdateStatusHistory(srData.StatusHistory, srData.StatusId)
	srData.SalesReturnNumber = helpers.GenerateSequence("SR", fmt.Sprint(metaData.TokenUserId), "sales_returns")
	defaultStatus, _ := helpers.GetLookupcodeId("RETURN_STATUS", "DRAFT")
	srData.StatusId = defaultStatus
	srData.CompanyId = metaData.CompanyId
	srData.CreatedByID = &metaData.TokenUserId
	result, _ := helpers.UpdateStatusHistory(srData.StatusHistory, srData.StatusId)
	srData.StatusHistory = result

	err := s.salesReturnRepository.Save(srData)
	if err != nil {
		return err
	}

	dataInv := map[string]interface{}{
		"id":            srData.ID,
		"order_lines":   srData.SalesReturnLines,
		"order_type":    "sales_return",
		"is_update_inv": true,
		"is_credit":     false,
	}
	//updating the inventory Committed Stock
	go s.basicInventoryService.UpdateTransactionInventory(metaData, dataInv)

	return nil
}

func (s *service) ListSalesReturns(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}
	data, err := s.salesReturnRepository.FindAll(metaData.Query, page)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) ViewSalesReturn(metaData core.MetaData) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales return at data level")
	}

	data, err := s.salesReturnRepository.FindOne(metaData.Query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) UpdateSalesReturn(metaData core.MetaData, data *returns.SalesReturns) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "UPDATE", "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update sales return at data level")
	}

	data = s.ProcessSalesReturnCalculation(data)
	var srData returns.SalesReturns
	helpers.JsonMarshaller(data, &srData)
	old_data, err := s.salesReturnRepository.FindOne(metaData.Query)
	if err != nil {
		return err
	}

	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, srData.StatusId)
		srData.StatusHistory = result
		final_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "FULLY_RETURNED")
		partial_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "PARTIALLY_RETURNED")
		if data.StatusId == final_status || data.StatusId == partial_status {
			dataInv := map[string]interface{}{
				"id":            data.ID,
				"order_lines":   data.SalesReturnLines,
				"order_type":    "sales_return",
				"is_update_inv": false,
				"is_credit":     false,
			}

			//updating the Committed Stock
			go s.basicInventoryService.UpdateTransactionInventory(metaData, dataInv)
		}
	}
	data.UpdatedByID = &metaData.TokenUserId
	err = s.salesReturnRepository.Update(metaData.Query, &srData)
	if err != nil {
		return err
	}
	for _, return_line := range srData.SalesReturnLines {
		update_query := map[string]interface{}{
			"product_id": return_line.ProductId,
			"sr_id":      metaData.Query["id"],
		}
		count, err1 := s.salesReturnRepository.UpdateReturnLines(update_query, return_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			return_line.SrId = uint(metaData.Query["id"].(float64))
			e := s.salesReturnRepository.SaveReturnLines(return_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeleteSalesReturn(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales return at data level")
	}
	err := s.salesReturnRepository.Delete(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{
			"sr_id": metaData.Query["id"],
		}
		err = s.salesReturnRepository.DeleteReturnLine(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteSalesReturnLines(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales return at data level")
	}
	err := s.salesReturnRepository.DeleteReturnLine(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SearchSalesReturns(query string, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}
	data, err := s.salesReturnRepository.Search(query)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) GetSalesReturnsHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "LIST", "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}
	data, err := s.salesReturnRepository.GetSalesReturnsHistory(metaData.Query, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (s *service) GetSalesReturnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "SALES_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}
	id := metaData.Query["id"].(string)
	tab := metaData.Query["tab"].(string)
	// token_id := fmt.Sprint(metaData.TokenUserId)
	// access_template_id := fmt.Sprint(metaData.AccessTemplateId)
	salesReturnId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "credit_note" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_RETURNS")

		creditNotePage := page
		creditNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesReturnId)
		// var query = make(map[string]interface{}, 0)

		creditNoteData, err := s.creditNoteService.GetCreditNoteList(metaData, creditNotePage)
		if err != nil {
			return nil, err
		}
		return creditNoteData, nil
	}
	if tab == "sales_orders" {
		salesReturnsPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")
		salesReturnsPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", salesReturnId, source_document_type_id)
		sales_returns_interface, err := s.ListSalesReturns(metaData, salesReturnsPagination)
		if err != nil {
			return nil, err
		}
		sales_returns := sales_returns_interface.([]returns.SalesReturns)
		if len(sales_returns) == 0 {
			return nil, errors.New("no source document")
		}

		var salesReturnsSourceDoc map[string]interface{}
		salesReturnsSourceDocJson := sales_returns[0].SourceDocuments
		dto, err := json.Marshal(salesReturnsSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &salesReturnsSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		salesReturnsSourceDocId := salesReturnsSourceDoc["id"]
		if salesReturnsSourceDocId == nil {
			return nil, errors.New("no source document")
		}

		salesOrdersPage := page
		salesOrdersPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", salesReturnsSourceDoc["id"])
		data, err := s.SalesOrdersRepository.FindAll(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "sales_invoice" {
		salesReturnsPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "SALES_INVOICE")
		salesReturnsPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", salesReturnId, source_document_type_id)
		sales_returns_interface, err := s.ListSalesReturns(core.MetaData{}, salesReturnsPagination)
		if err != nil {
			return nil, err
		}
		sales_returns := sales_returns_interface.([]returns.SalesReturns)
		if len(sales_returns) == 0 {
			return nil, errors.New("no source document")
		}

		var salesReturnsSourceDoc map[string]interface{}
		salesReturnsSourceDocJson := sales_returns[0].SourceDocuments
		dto, err := json.Marshal(salesReturnsSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &salesReturnsSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		salesReturnsSourceDocId := salesReturnsSourceDoc["id"]
		if salesReturnsSourceDocId == nil {
			return nil, errors.New("no source document")
		}

		salesInvoicePage := page
		salesInvoicePage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", salesReturnsSourceDoc["id"])
		data, err := s.SalesInvoiceRepository.GetAllSalesInvoice(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "delivery_orders" {
		salesReturnsPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDERS")
		salesReturnsPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", salesReturnId, source_document_type_id)
		sales_returns_interface, err := s.ListSalesReturns(core.MetaData{}, salesReturnsPagination)
		if err != nil {
			return nil, err
		}
		sales_returns := sales_returns_interface.([]returns.SalesReturns)
		if len(sales_returns) == 0 {
			return nil, errors.New("no source document")
		}

		var salesReturnsSourceDoc map[string]interface{}
		salesReturnsSourceDocJson := sales_returns[0].SourceDocuments
		dto, err := json.Marshal(salesReturnsSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &salesReturnsSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		salesReturnsSourceDocId := salesReturnsSourceDoc["id"]
		if salesReturnsSourceDocId == nil {
			return nil, errors.New("no source document")
		}

		deliveryOrdersPage := page
		deliveryOrdersPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", salesReturnsSourceDoc["id"])
		data, err := s.DeliveryOrdersRepository.AllDeliveryOrders(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
func (s *service) ProcessSalesReturnCalculation(data *returns.SalesReturns) *returns.SalesReturns {
	var totalQuantity int64 = 0

	for _, returnLines := range data.SalesReturnLines {
		totalQuantity = totalQuantity + returnLines.QuantityReturned
	}
	data.TotalQuantity = totalQuantity
	return data
}

func (s *service) TrackSalesReturn(order_id string) map[string]interface{} {
	var sales_order_res SalesReturnsDTO
	res, _ := s.salesReturnRepository.FindOne(map[string]interface{}{"sales_return_number": order_id})
	dto, _ := json.Marshal(res)
	json.Unmarshal(dto, &sales_order_res)
	return sales_order_res.StatusHistory
}
