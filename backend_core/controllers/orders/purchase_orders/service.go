package purchase_orders

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/debitnote"
	"fermion/backend_core/controllers/accounting/purchase_invoice"
	"fermion/backend_core/controllers/inventory_orders/asn"
	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
	"fermion/backend_core/controllers/returns/purchase_returns"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
	orders_repo "fermion/backend_core/internal/repository/orders"

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
	CreatePurchaseOrder(metaData core.MetaData, data *PurchaseOrdersDTO) error
	ListPurchaseOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ViewPurchaseOrder(metaData core.MetaData) (interface{}, error)
	UpdatePurchaseOrder(metaData core.MetaData, data *PurchaseOrdersDTO) error
	DeletePurchaseOrder(metaData core.MetaData) error
	DeletePurchaseOrderLines(metaData core.MetaData) error
	SearchPurchaseOrders(query string, access_template_id string, token_id string) (interface{}, error)
	ProcessPurchaseOrderCalculation(data *PurchaseOrdersDTO) *PurchaseOrdersDTO
	GetPurchaseHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)

	GetPurchaseOrderTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	purchaseOrderRepository orders_repo.Purchase
	basicInventoryService   inv_service.Service

	purchaseReturnsService purchase_returns.Service
	purchaseInvoiceService purchase_invoice.Service
	asnService             asn.Service
	debitNoteService       debitnote.Service
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		orders_repo.NewPurchaseOrder(),
		inv_service.NewService(),
		purchase_returns.NewService(),
		purchase_invoice.NewService(),
		asn.NewService(),
		debitnote.NewService(),
	}
	return newServiceObj
}

func (s *service) CreatePurchaseOrder(metaData core.MetaData, data *PurchaseOrdersDTO) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "CREATE", "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create purchase order at data level")
	}
	data = s.ProcessPurchaseOrderCalculation(data)
	var poData *orders.PurchaseOrders
	helpers.JsonMarshaller(data, &poData)
	result, _ := helpers.UpdateStatusHistory(poData.StatusHistory, poData.StatusId)
	poData.StatusHistory = result
	poData.PurchaseOrderNumber = helpers.GenerateSequence("PO", fmt.Sprint(metaData.TokenUserId), "purchase_orders")
	defaultStatus, _ := helpers.GetLookupcodeId("PURCHASE_ORDER_STATUS", "DRAFT")
	poData.StatusId = defaultStatus
	poData.CompanyId = metaData.CompanyId
	poData.CreatedByID = &metaData.TokenUserId

	err := s.purchaseOrderRepository.Save(poData)
	if err != nil {
		return err
	}
	data_inv := map[string]interface{}{
		"id":            poData.ID,
		"order_lines":   poData.PurchaseOrderlines,
		"order_type":    "purchase_order",
		"is_update_inv": false,
		"is_credit":     true,
	}

	//updating the inventory Expected Stock
	go s.basicInventoryService.UpdateTransactionInventory(metaData, data_inv)

	return nil
}

func (s *service) ListPurchaseOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase order at data level")
	}

	data, err := s.purchaseOrderRepository.FindAll(metaData.Query, page)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) ViewPurchaseOrder(metaData core.MetaData) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase order at data level")
	}

	data, err := s.purchaseOrderRepository.FindOne(metaData.Query)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) UpdatePurchaseOrder(metaData core.MetaData, data *PurchaseOrdersDTO) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "UPDATE", "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update purchase order at data level")
	}
	data = s.ProcessPurchaseOrderCalculation(data)
	var PoData orders.PurchaseOrders
	helpers.JsonMarshaller(data, &PoData)
	old_data, err := s.purchaseOrderRepository.FindOne(metaData.Query)
	if err != nil {
		return err
	}

	old_status := old_data.StatusId
	new_status := PoData.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, PoData.StatusId)
		PoData.StatusHistory = result

		final_status, _ := helpers.GetLookupcodeId("PURCHASE_ORDER_STATUS", "RECEIVED")
		if data.StatusId == final_status {
			data_inv := map[string]interface{}{
				"id":            data.ID,
				"order_lines":   data.PurchaseOrderlines,
				"order_type":    "purchase_order",
				"is_update_inv": true,
				"is_credit":     true,
			}

			//updating the inventory Expected Stock
			go s.basicInventoryService.UpdateTransactionInventory(metaData, data_inv)
		}

	}
	data.UpdatedByID = &metaData.TokenUserId
	err = s.purchaseOrderRepository.Update(metaData.Query, &PoData)
	if err != nil {
		return err
	}

	for _, order_line := range PoData.PurchaseOrderlines {
		update_query := map[string]interface{}{
			"product_id": order_line.ProductId,
			"po_id":      metaData.Query["id"],
		}
		count, err1 := s.purchaseOrderRepository.UpdateOrderLines(update_query, order_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			order_line.PoId = uint(metaData.Query["id"].(float64))
			e := s.purchaseOrderRepository.SaveOrderLines(order_line)
			fmt.Println(e)
			if e != nil {
				return e
			}
		}
	}
	return nil
	//}
	//return nil
}

func (s *service) DeletePurchaseOrder(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase order at data level")
	}

	err := s.purchaseOrderRepository.Delete(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{
			"po_id": metaData.Query["id"],
		}
		err = s.purchaseOrderRepository.DeleteOrderLine(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeletePurchaseOrderLines(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase order at data level")
	}
	err := s.purchaseOrderRepository.DeleteOrderLine(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SearchPurchaseOrders(query string, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase order at data level")
	}

	data, err := s.purchaseOrderRepository.Search(query)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
	//}
	//return nil, nil
}

func (s *service) ProcessPurchaseOrderCalculation(data *PurchaseOrdersDTO) *PurchaseOrdersDTO {

	var subTotal, totalTax float32
	var totalQuantity int64 = 0

	for index, orderLines := range data.PurchaseOrderlines {
		amountWithoutDiscount := (orderLines.Price * float32(orderLines.Quantity))
		amountWithoutTax := amountWithoutDiscount - orderLines.Discount
		tax := (amountWithoutTax * orderLines.Tax) / 100
		amount := amountWithoutTax + tax
		totalQuantity = totalQuantity + orderLines.Quantity

		data.PurchaseOrderlines[index].Amount = amount
		subTotal += amountWithoutTax
		totalTax += tax
	}
	data.TotalQuantity = totalQuantity
	data.PoPaymentDetails.SubTotal = subTotal
	data.PoPaymentDetails.Tax = totalTax
	data.PoPaymentDetails.TotalAmount = subTotal + totalTax + data.PoPaymentDetails.ShippingCharges

	return data
}

func (s *service) GetPurchaseHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "LIST", "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase order at data level")
	}
	data, err := s.purchaseOrderRepository.GetPurchaseHistory(metaData.Query, page)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) GetPurchaseOrderTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "LIST", "PURCHASE_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales order at data level")
	}
	id := metaData.Query["id"].(string)
	tab := metaData.Query["tab"].(string)
	// token_id := fmt.Sprint(metaData.TokenUserId)
	// access_template_id := fmt.Sprint(metaData.AccessTemplateId)
	purchaseOrderId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "purchase_returns" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_RETURNS_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		purchaseReturnPage := page
		purchaseReturnPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseOrderId)
		data, err := s.purchaseReturnsService.ListPurchaseReturns(core.MetaData{}, purchaseReturnPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "purchase_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_INVOICE_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		purchaseInvoicePage := page
		purchaseInvoicePage.Filters = fmt.Sprintf("[[\"link_source_document_type\",\"=\",%v],[\"link_source_document\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseOrderId)
		data, err := s.purchaseInvoiceService.ListPurchaseInvoice(metaData, purchaseInvoicePage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "asn" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("ASN_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		asnPage := page
		asnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseOrderId)
		data, err := s.asnService.GetAllAsn(metaData, asnPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "debit_note" {

		data := make([]debitnote.DebitNoteResponseDTO, 0)
		dto, _ := json.Marshal(data)
		json.Unmarshal(dto, &data)

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_RETURNS_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		purchaseReturnPage := new(pagination.Paginatevalue)
		purchaseReturnPage.Per_page = 1000
		purchaseReturnPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseOrderId)
		purchaseReturnsData, err := s.purchaseReturnsService.ListPurchaseReturns(metaData, purchaseReturnPage)

		if err != nil {
			return nil, err
		}

		for _, purchaseReturn := range purchaseReturnsData.([]returns.PurchaseReturns) {

			sourceDocumentId, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_RETURNS")

			debitNotePage := new(pagination.Paginatevalue)
			debitNotePage.Per_page = 1000
			debitNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseReturn.ID)

			var metaData core.MetaData
			debitNoteDataInterface, err := s.debitNoteService.GetDebitNoteList(metaData, debitNotePage)

			if err != nil {
				return nil, err
			}
			debitNoteData := debitNoteDataInterface.([]debitnote.DebitNoteResponseDTO)
			data = append(data, debitNoteData...)
		}

		return data, nil
	}

	return nil, nil
}
