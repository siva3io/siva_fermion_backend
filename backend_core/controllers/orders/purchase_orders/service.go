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
	CreatePurchaseOrder(data *PurchaseOrdersDTO, access_template_id string, token_id string) error
	ListPurchaseOrders(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error)
	ViewPurchaseOrder(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error)
	UpdatePurchaseOrder(query map[string]interface{}, data *orders.PurchaseOrders, access_template_id string, token_id string) error
	DeletePurchaseOrder(query map[string]interface{}, access_template_id string, token_id string) error
	DeletePurchaseOrderLines(query map[string]interface{}, access_template_id string, token_id string) error
	SearchPurchaseOrders(query string, access_template_id string, token_id string) (interface{}, error)
	ProcessPurchaseOrderCalculation(data *orders.PurchaseOrders) *orders.PurchaseOrders
	GetPurchaseHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)

	GetPurchaseOrderTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	purchaseOrderRepository orders_repo.Purchase
	basicInventoryService   inv_service.Service

	purchaseReturnsService purchase_returns.Service
	purchaseInvoiceService purchase_invoice.Service
	asnService             asn.Service
	debitNoteService       debitnote.Service
}

func NewService() *service {
	return &service{
		orders_repo.NewPurchaseOrder(),
		inv_service.NewService(),
		purchase_returns.NewService(),
		purchase_invoice.NewService(),
		asn.NewService(),
		debitnote.NewService(),
	}
}

func (s *service) CreatePurchaseOrder(data *PurchaseOrdersDTO, access_template_id string, token_id string) error {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create purchase order at data level")
	}
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	temp := data.CreatedByID
	user_id := strconv.Itoa(int(*temp))
	data.PurchaseOrderNumber = helpers.GenerateSequence("PO", user_id, "purchase_orders")
	defaultStatus, err := helpers.GetLookupcodeId("PURCHASE_ORDER_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	var datadto *orders.PurchaseOrders
	dto, err := json.Marshal(data)
	json.Unmarshal(dto, &datadto)
	fmt.Println("data", datadto.DeliveryTo)

	datadto = s.ProcessPurchaseOrderCalculation(datadto)

	err = s.purchaseOrderRepository.Save(datadto)

	data_inv := map[string]interface{}{
		"id":            data.ID,
		"order_lines":   data.PurchaseOrderlines,
		"order_type":    "purchase_order",
		"is_update_inv": false,
		"is_credit":     true,
	}

	//updating the inventory Expected Stock
	go s.basicInventoryService.UpdateInventory(data_inv, access_template_id, token_id)

	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return nil
	//}
	//return nil
}

func (s *service) ListPurchaseOrders(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase order at data level")
	}

	data, err := s.purchaseOrderRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
	//}
	//return nil, nil
}

func (s *service) ViewPurchaseOrder(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase order at data level")
	}

	data, err := s.purchaseOrderRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
	//}
	//return nil, nil
}

func (s *service) UpdatePurchaseOrder(query map[string]interface{}, data *orders.PurchaseOrders, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update purchase order at data level")
	}

	var find_query = map[string]interface{}{
		"id": query["id"],
	}

	found_data, er := s.purchaseOrderRepository.FindOne(find_query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	old_data := found_data.(orders.PurchaseOrders)

	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result

		final_status, _ := helpers.GetLookupcodeId("PURCHASE_ORDER_STATUS", "RECEIVED")
		if data.StatusId == final_status {
			data_inv := map[string]interface{}{
				"id":            query["id"],
				"order_lines":   data.PurchaseOrderlines,
				"order_type":    "purchase_order",
				"is_update_inv": true,
				"is_credit":     true,
			}

			//updating the inventory Expected Stock
			go s.basicInventoryService.UpdateInventory(data_inv, access_template_id, token_id)
		}

	}

	data = s.ProcessPurchaseOrderCalculation(data)

	err := s.purchaseOrderRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	for _, order_line := range data.PurchaseOrderlines {
		update_query := map[string]interface{}{
			"product_id": order_line.ProductId,
			"po_id":      uint(query["id"].(int)),
		}
		count, err1 := s.purchaseOrderRepository.UpdateOrderLines(update_query, order_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			order_line.PoId = uint(query["id"].(int))
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

func (s *service) DeletePurchaseOrder(query map[string]interface{}, access_template_id string, token_id string) error {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase order at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.purchaseOrderRepository.FindOne(q)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.purchaseOrderRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	} else {
		query := map[string]interface{}{"id": query["id"]}
		er := s.purchaseOrderRepository.DeleteOrderLine(query)
		if er != nil {
			return res.BuildError(res.ErrDataNotFound, er)
		}
	}
	return nil
	//}
	//return nil
}

func (s *service) DeletePurchaseOrderLines(query map[string]interface{}, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase order at data level")
	}
	_, er := s.purchaseOrderRepository.FindOrderLines(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.purchaseOrderRepository.DeleteOrderLine(query)

	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
	//}
	//return nil
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

func (s *service) ProcessPurchaseOrderCalculation(data *orders.PurchaseOrders) *orders.PurchaseOrders {
	var subTotal, totalTax float32

	for index, orderLines := range data.PurchaseOrderlines {
		amountWithoutDiscount := (orderLines.Price * float32(orderLines.Quantity))
		amountWithoutTax := amountWithoutDiscount - orderLines.Discount
		tax := (amountWithoutTax * orderLines.Tax) / 100
		amount := amountWithoutTax + tax

		data.PurchaseOrderlines[index].Amount = amount
		subTotal += amountWithoutTax
		totalTax += tax
	}

	data.PoPaymentDetails.SubTotal = subTotal
	data.PoPaymentDetails.Tax = totalTax
	data.PoPaymentDetails.TotalAmount = subTotal + totalTax + data.PoPaymentDetails.ShippingCharges

	return data
}

func (s *service) GetPurchaseHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PURCHASE_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase order at data level")
	}
	data, err := s.purchaseOrderRepository.GetPurchaseHistory(productId, page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) GetPurchaseOrderTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales order at data level")
	}

	purchaseOrderId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "purchase_returns" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_RETURNS_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		purchaseReturnPage := page
		purchaseReturnPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseOrderId)
		data, err := s.purchaseReturnsService.ListPurchaseReturns(purchaseReturnPage, access_template_id, token_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "purchase_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_INVOICE_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		purchaseInvoicePage := page
		purchaseInvoicePage.Filters = fmt.Sprintf("[[\"link_source_document_type\",\"=\",%v],[\"link_source_document\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseOrderId)
		data, err := s.purchaseInvoiceService.ListPurchaseInvoice(purchaseInvoicePage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "asn" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("ASN_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		asnPage := page
		asnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseOrderId)
		data, err := s.asnService.GetAllAsn(asnPage, token_id, access_template_id, "LIST")
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
		purchaseReturnsData, err := s.purchaseReturnsService.ListPurchaseReturns(purchaseReturnPage, access_template_id, token_id, "LIST")

		if err != nil {
			return nil, err
		}

		for _, purchaseReturn := range purchaseReturnsData.([]returns.PurchaseReturns) {

			sourceDocumentId, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_RETURNS")

			debitNotePage := new(pagination.Paginatevalue)
			debitNotePage.Per_page = 1000
			debitNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseReturn.ID)

			var query = make(map[string]interface{}, 0)
			debitNoteData, err := s.debitNoteService.GetDebitNoteList(query, debitNotePage, token_id, access_template_id, "LIST")

			if err != nil {
				return nil, err
			}
			data = append(data, debitNoteData...)
		}

		return data, nil
	}

	return nil, nil
}
