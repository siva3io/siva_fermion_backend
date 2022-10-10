package sales_orders

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/creditnote"
	"fermion/backend_core/controllers/accounting/sales_invoice"
	"fermion/backend_core/controllers/returns/sales_returns"

	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
	"fermion/backend_core/controllers/orders/delivery_orders"
	"fermion/backend_core/controllers/orders/purchase_orders"

	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
	user_repo "fermion/backend_core/internal/repository"
	access_repo "fermion/backend_core/internal/repository/access"
	orders_repo "fermion/backend_core/internal/repository/orders"

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
	CreateSalesOrder(data *orders.SalesOrders, access_template_id string, token_id string) error
	ListSalesOrders(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error)
	ViewSalesOrder(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error)
	UpdateSalesOrder(query map[string]interface{}, data *orders.SalesOrders, access_template_id string, token_id string) error
	DeleteSalesOrder(query map[string]interface{}, access_template_id string, token_id string) error
	DeleteSalesOrderLines(query map[string]interface{}, access_template_id string, token_id string) error
	SearchSalesOrders(query string, access_template_id string, token_id string) (interface{}, error)
	ChannelSalesOrderUpsert(data []orders.SalesOrders, TokenUserId string) (interface{}, error)
	ChannelSalesOrderLinesUpsert(data []orders.SalesOrderLines, salesOrderId uint) (interface{}, error)
	ProcessSalesOrderCalculation(data *orders.SalesOrders) *orders.SalesOrders
	GetSalesHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)

	GetSalesOrderTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	salesOrderRepository      orders_repo.Sales
	basicInventoryService     inv_service.Service
	accessModuleRepository    access_repo.Module
	acccessTemplateRepository access_repo.Template
	userRepository            user_repo.User

	deliveryOrdersService delivery_orders.Service
	purchaseOrdersService purchase_orders.Service
	salesInvoiceService   sales_invoice.Service
	salesReturnsService   sales_returns.Service
	creditNoteService     creditnote.Service
}

func NewService() *service {

	return &service{
		salesOrderRepository:      orders_repo.NewSalesOrder(),
		basicInventoryService:     inv_service.NewService(),
		accessModuleRepository:    access_repo.NewModule(),
		acccessTemplateRepository: access_repo.NewTemplate(),
		userRepository:            user_repo.NewUser(),
		deliveryOrdersService:     delivery_orders.NewService(),
		purchaseOrdersService:     purchase_orders.NewService(),
		salesInvoiceService:       sales_invoice.NewService(),
		salesReturnsService:       sales_returns.NewService(),
		creditNoteService:         creditnote.NewService(),
	}
}

func (s *service) CreateSalesOrder(data *orders.SalesOrders, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create sales order at data level")
	}

	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	temp := data.CreatedByID
	user_id := strconv.Itoa(int(*temp))
	data.SalesOrderNumber = helpers.GenerateSequence("SO", user_id, "sales_orders")
	data = s.ProcessSalesOrderCalculation(data)
	defaultStatus, err := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	err = s.salesOrderRepository.Save(data)
	data_inv := map[string]interface{}{
		"id":            data.ID,
		"order_lines":   data.SalesOrderLines,
		"order_type":    "sales_order",
		"is_update_inv": true,
		"is_credit":     false,
	}

	//updating the inventory Committed Stock
	go s.basicInventoryService.UpdateInventory(data_inv, access_template_id, token_id)

	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return nil

}

func (s *service) ListSalesOrders(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales order at data level")
	}

	data, err := s.salesOrderRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil

}

func (s *service) ViewSalesOrder(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales order at data level")
	}

	data, err := s.salesOrderRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil

}

func (s *service) UpdateSalesOrder(query map[string]interface{}, data *orders.SalesOrders, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update sales order at data level")
	}

	var find_query = map[string]interface{}{
		"id": query["id"],
	}

	found_data, er := s.salesOrderRepository.FindOne(find_query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	old_data := found_data.(orders.SalesOrders)

	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result
		final_status, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "DELIVERED")
		if data.StatusId == final_status {
			data_inv := map[string]interface{}{
				"id":            data.ID,
				"order_lines":   data.SalesOrderLines,
				"order_type":    "sales_order",
				"is_update_inv": false,
				"is_credit":     false,
			}

			//updating the inventory Committed Stock
			go s.basicInventoryService.UpdateInventory(data_inv, access_template_id, token_id)
		}
	}

	data = s.ProcessSalesOrderCalculation(data)

	err := s.salesOrderRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, order_line := range data.SalesOrderLines {
		update_query := map[string]interface{}{
			"product_id": order_line.ProductId,
			"so_id":      uint(query["id"].(int)),
		}
		count, err1 := s.salesOrderRepository.UpdateOrderLines(update_query, order_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			order_line.SoId = uint(query["id"].(int))
			e := s.salesOrderRepository.SaveOrderLines(order_line)
			fmt.Println(e)
			if e != nil {
				return e
			}
		}
	}
	return nil

}

func (s *service) DeleteSalesOrder(query map[string]interface{}, access_template_id string, si string) error {
	user_id := helpers.ConvertStringToUint(si)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SALES_ORDERS", *user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales order at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.salesOrderRepository.FindOne(q)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.salesOrderRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	} else {
		query := map[string]interface{}{"id": query["id"]}
		er := s.salesOrderRepository.DeleteOrderLine(query)
		if er != nil {
			return res.BuildError(res.ErrDataNotFound, er)
		}
	}
	return nil
}

func (s *service) DeleteSalesOrderLines(query map[string]interface{}, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales order at data level")
	}
	_, er := s.salesOrderRepository.FindOrderLines(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.salesOrderRepository.DeleteOrderLine(query)

	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil

}

func (s *service) SearchSalesOrders(query string, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales order at data level")
	}
	data, err := s.salesOrderRepository.Search(query)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil

}

func (s *service) ChannelSalesOrderUpsert(data []orders.SalesOrders, TokenUserId string) (interface{}, error) {

	//access_module_flag, _ := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_ORDERS")

	//if access_module_flag {
	var success []interface{}
	var failures []interface{}

	for index, salesOrder := range data {

		if salesOrder.ReferenceNumber == "" {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "reference_number should not be empty"})
			continue
		}

		salesOrder.UpdatedByID = helpers.ConvertStringToUint(TokenUserId)

		query := map[string]interface{}{"reference_number": salesOrder.ReferenceNumber}

		data, _ := s.salesOrderRepository.FindOne(query)

		if data != nil {

			fetchedSalesOrder := data.(orders.SalesOrders)

			old_status := fetchedSalesOrder.StatusId
			new_status := salesOrder.StatusId

			if new_status != old_status && new_status != 0 {
				result, _ := helpers.UpdateStatusHistory(fetchedSalesOrder.StatusHistory, new_status)
				salesOrder.StatusHistory = result
			}

			err := s.salesOrderRepository.Update(query, &salesOrder)
			if err != nil {
				failures = append(failures, map[string]interface{}{"reference_number": salesOrder.ReferenceNumber, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			result, err := s.ChannelSalesOrderLinesUpsert(salesOrder.SalesOrderLines, fetchedSalesOrder.ID)
			if err != nil {
				failures = append(failures, map[string]interface{}{"reference_number": salesOrder.ReferenceNumber, "status": false, "serial_number": index + 1, "msg": "sales order line error", "sales_order_lines": err})
				continue
			}
			success = append(success, map[string]interface{}{"reference_number": salesOrder.ReferenceNumber, "status": true, "serial_number": index + 1, "msg": "updated", "sales_order_lines": result})
			continue
		}

		salesOrder.SalesOrderNumber = helpers.GenerateSequence("SO", TokenUserId, "sales_orders")
		salesOrder.CreatedByID = helpers.ConvertStringToUint(TokenUserId)
		salesOrder.UpdatedByID = nil
		statusHistory, _ := helpers.UpdateStatusHistory(salesOrder.StatusHistory, salesOrder.StatusId)
		salesOrder.StatusHistory = statusHistory

		err := s.salesOrderRepository.Save(&salesOrder)
		if err != nil {
			failures = append(failures, map[string]interface{}{"reference_number": salesOrder.ReferenceNumber, "status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		success = append(success, map[string]interface{}{"reference_number": salesOrder.ReferenceNumber, "status": true, "serial_number": index + 1, "msg": "created"})
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
	//}
	//return nil, nil
}

func (s *service) ChannelSalesOrderLinesUpsert(data []orders.SalesOrderLines, salesOrderId uint) (interface{}, error) {
	var success []interface{}
	var failures []interface{}

	for index, orderLine := range data {
		update_query := map[string]interface{}{
			"product_id": orderLine.ProductId,
			"so_id":      salesOrderId,
		}
		count, err := s.salesOrderRepository.UpdateOrderLines(update_query, orderLine)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		if count == 0 {
			orderLine.SoId = salesOrderId
			err := s.salesOrderRepository.SaveOrderLines(orderLine)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			success = append(success, map[string]interface{}{"status": true, "serial_number": index + 1, "msg": "created"})
			continue
		}
		success = append(success, map[string]interface{}{"status": true, "serial_number": index + 1, "msg": "updated"})
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
}

func (s *service) ProcessSalesOrderCalculation(data *orders.SalesOrders) *orders.SalesOrders {
	var subTotal, totalTax, shippingCharges float32

	for index, orderLines := range data.SalesOrderLines {
		amountWithoutDiscount := (orderLines.Price * float32(orderLines.Quantity))
		amountWithoutTax := amountWithoutDiscount - orderLines.Discount
		tax := (amountWithoutTax * orderLines.Tax) / 100
		amount := amountWithoutTax + tax

		data.SalesOrderLines[index].Amount = amount
		subTotal += amountWithoutTax
		totalTax += tax
	}

	data.SoPaymentDetails.SubTotal = subTotal
	data.SoPaymentDetails.Tax = totalTax
	var vendorDetails map[string]interface{}
	err := json.Unmarshal(data.VendorDetails, &vendorDetails)
	if err == nil {
		if vendorDeliveryCharges, ok := vendorDetails["vendor_delivery_charges"].(float64); ok {
			shippingCharges = float32(vendorDeliveryCharges)
		}
	}
	data.SoPaymentDetails.ShippingCharges = shippingCharges
	data.SoPaymentDetails.TotalAmount = subTotal + totalTax + shippingCharges

	return data
}

func (s *service) GetSalesHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales order at data level")
	}

	data, err := s.salesOrderRepository.GetSalesHistory(productId, page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil

}

func (s *service) GetSalesOrderTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales order at data level")
	}

	salesOrderId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "delivery_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		deliveryOrderPage := page
		deliveryOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		data, err := s.deliveryOrdersService.AllDeliveryOrders(deliveryOrderPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "purchase_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		purchaseOrderPage := page
		purchaseOrderPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		data, err := s.purchaseOrdersService.ListPurchaseOrders(purchaseOrderPage, access_template_id, token_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "sales_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_INVOICE_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		salesInvoicePage := page
		salesInvoicePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		data, err := s.salesInvoiceService.GetAllSalesInvoice(salesInvoicePage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "sales_returns" {

		data := make([]returns.SalesReturns, 0)

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		deliveryOrderPage := new(pagination.Paginatevalue)
		deliveryOrderPage.Per_page = 1000
		deliveryOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		deliveryOrderData, err := s.deliveryOrdersService.AllDeliveryOrders(deliveryOrderPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}

		for _, deliveryOrder := range deliveryOrderData {

			sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDERS")

			salesReturnsPage := new(pagination.Paginatevalue)
			salesReturnsPage.Per_page = 1000
			salesReturnsPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrder.ID)
			salesReturnsDataInterface, err := s.salesReturnsService.ListSalesReturns(salesReturnsPage, access_template_id, token_id, "LIST")
			if err != nil {
				return nil, err
			}
			salesReturnsData := salesReturnsDataInterface.([]returns.SalesReturns)
			data = append(data, salesReturnsData...)
		}
		return data, nil
	}

	if tab == "credit_note" {

		data := make([]creditnote.CreditNoteResponseDTO, 0)
		dto, _ := json.Marshal(data)
		err := json.Unmarshal(dto, &data)

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		deliveryOrderPage := new(pagination.Paginatevalue)
		deliveryOrderPage.Per_page = 1000
		deliveryOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		deliveryOrderData, err := s.deliveryOrdersService.AllDeliveryOrders(deliveryOrderPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}

		for _, deliveryOrder := range deliveryOrderData {

			sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDERS")

			salesReturnsPage := new(pagination.Paginatevalue)
			salesReturnsPage.Per_page = 1000
			salesReturnsPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrder.ID)
			salesReturnsDataInterface, err := s.salesReturnsService.ListSalesReturns(salesReturnsPage, access_template_id, token_id, "LIST")
			if err != nil {
				return nil, err
			}
			salesReturnsData := salesReturnsDataInterface.([]returns.SalesReturns)

			for _, salesReturn := range salesReturnsData {

				sourceDocumentId, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_RETURS")

				creditNotePage := new(pagination.Paginatevalue)
				creditNotePage.Per_page = 1000
				creditNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesReturn.ID)

				var query = make(map[string]interface{}, 0)

				creditNoteData, err := s.creditNoteService.GetCreditNoteList(query, creditNotePage, token_id, access_template_id, "LIST")
				if err != nil {
					return nil, err
				}
				data = append(data, creditNoteData...)
			}
		}
		return data, nil
	}
	return nil, nil
}
