package sales_orders

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"fermion/backend_core/controllers/accounting/creditnote"
	"fermion/backend_core/controllers/accounting/sales_invoice"
	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/controllers/omnichannel/channel"
	"fermion/backend_core/controllers/returns/sales_returns"
	invoice_model "fermion/backend_core/internal/model/accounting"

	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
	"fermion/backend_core/controllers/orders/delivery_orders"
	"fermion/backend_core/controllers/orders/purchase_orders"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
	user_repo "fermion/backend_core/internal/repository"
	access_repo "fermion/backend_core/internal/repository/access"
	orders_repo "fermion/backend_core/internal/repository/orders"

	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
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
	CreateSalesOrder(metaData core.MetaData, data *orders.SalesOrders) error
	ListSalesOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ViewSalesOrder(metaData core.MetaData) (interface{}, error)
	UpdateSalesOrder(metaData core.MetaData, data *orders.SalesOrders) error
	DeleteSalesOrder(metaData core.MetaData) error
	DeleteSalesOrderLines(metaData core.MetaData) error
	ChannelSalesOrderUpsert(metaData core.MetaData, data []orders.SalesOrders) (interface{}, error)
	GetSalesHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	TrackSalesOrder(query map[string]interface{}) []map[string]interface{}
	GetSalesOrderTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)

	ChannelSalesOrderLinesUpsert(data []orders.SalesOrderLines, salesOrderId uint) (interface{}, error)
	ProcessSalesOrderCalculation(data *orders.SalesOrders) *orders.SalesOrders

	MapWithSalesInvoice(requestPayload *orders.SalesOrders, edaMetaData core.MetaData, flag bool) *invoice_model.SalesInvoice
	MapWithSalesReturns(requestPayload *orders.SalesOrders, edaMetaData core.MetaData, flag bool) *returns.SalesReturns
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
	channelService        channel.Service
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
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
		channelService:            channel.NewService(),
	}
	return newServiceObj
}

func (s *service) CreateSalesOrder(metaData core.MetaData, data *orders.SalesOrders) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "CREATE", "SALES_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create sales order at data level")
	}

	if data.ChannelName != "ONDC" {
		data = s.ProcessSalesOrderCalculation(data)
	}

	// data.StatusHistory, _ = helpers.UpdateStatusHistory(data.StatusHistory, *data.StatusId)
	data.SalesOrderNumber = helpers.GenerateSequence("SO", fmt.Sprint(metaData.TokenUserId), "sales_orders")

	if data.StatusId == nil || *data.StatusId == 0 {
		defaultStatus, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "DRAFT")
		data.StatusId = &defaultStatus
	}
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, *data.StatusId)
	data.StatusHistory = result
	sales_order_id, err := s.salesOrderRepository.Save(data)
	if err != nil {
		return err
	}

	dataInv := map[string]interface{}{
		"id":            data.ID,
		"order_lines":   data.SalesOrderLines,
		"order_type":    "sales_order",
		"is_update_inv": true,
		"is_credit":     false,
	}
	//updating the inventory Committed Stock
	go s.basicInventoryService.UpdateTransactionInventory(metaData, dataInv)
	// if data.ChannelName == "ONDC" {
	//ondc on_select call
	// }
	Status, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "NEW")
	*data.StatusId = Status
	if *data.StatusId == Status {
		salesInvoice := s.MapWithSalesInvoice(data, metaData, false)
		metaData.Query = map[string]interface{}{
			"sales_order_id": sales_order_id,
		}
		go s.salesInvoiceService.CreateSalesInvoice(metaData, salesInvoice)
	}

	return nil
}

func (s *service) ListSalesOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "SALES_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales order at data level")
	}

	data, err := s.salesOrderRepository.FindAll(metaData.Query, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) ViewSalesOrder(metaData core.MetaData) (interface{}, error) {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "SALES_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales order at data level")
	}

	data, err := s.salesOrderRepository.FindOne(metaData.Query)
	if err != nil {

		return nil, err
	}
	return data, nil
}

func (s *service) UpdateSalesOrder(metaData core.MetaData, data *orders.SalesOrders) error {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "UPDATE", "SALES_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update sales order at data level")
	}

	foundData, err := s.salesOrderRepository.FindOne(metaData.Query)
	if err != nil {
		return err
	}

	if data.StatusId != nil {
		if *data.StatusId != *foundData.StatusId && *data.StatusId != 0 {

			finalStatus, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "DELIVERED")
			cancelStatus, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "CANCELLED")
			shippedStatus, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "SHIPPED")
			readyToShipStatus, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "READY_TO_SHIP")

			zone := os.Getenv("DB_TZ")
			loc, _ := time.LoadLocation(zone)

			if *data.StatusId == shippedStatus {
				data.ShippedDate = time.Now().In(loc)
			}

			if *data.StatusId == readyToShipStatus {
				data.ReadyTOShip = time.Now().In(loc)
			}
			if *data.StatusId == cancelStatus {
				data.CancelledDate = time.Now().In(loc)
			}

			if *data.StatusId == finalStatus {
				data.DeliveryDate = time.Now().In(loc)
				dataInv := map[string]interface{}{
					"id":            data.ID,
					"order_lines":   data.SalesOrderLines,
					"order_type":    "sales_order",
					"is_update_inv": false,
					"is_credit":     false,
				}
				//updating the inventory Committed Stock
				go s.basicInventoryService.UpdateTransactionInventory(metaData, dataInv)
			}
			result, _ := helpers.UpdateStatusHistory(data.StatusHistory, *data.StatusId)
			data.StatusHistory = result
		}
	}

	if foundData.ChannelName != "ONDC" {
		data = s.ProcessSalesOrderCalculation(data)
	}
	data.UpdatedByID = &metaData.TokenUserId

	err = s.salesOrderRepository.Update(metaData.Query, data)
	if err != nil {
		return err
	}
	for _, orderLine := range data.SalesOrderLines {
		updateQuery := map[string]interface{}{
			"product_id": orderLine.ProductId,
			"so_id":      foundData.ID,
		}
		err = s.salesOrderRepository.UpdateOrderLines(updateQuery, &orderLine)
		if err != nil {
			orderLine.SoId = foundData.ID
			err := s.salesOrderRepository.SaveOrderLines(&orderLine)
			if err != nil {
				return err
			}
		}
	}

	var dontProduce interface{}
	if metaData.AdditionalFields != nil {
		dontProduce = metaData.AdditionalFields["don't produce"]
	}

	if foundData.ChannelName == "ONDC" && dontProduce == nil {
		err = s.UpdateInOndc(*data, foundData)
		if err != nil {
			return err
		}
	}

	Status, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "NEW")
	if data.StatusId != nil && *data.StatusId == Status {
		salesInvoice := s.MapWithSalesInvoice(data, metaData, true)
		go s.salesInvoiceService.CreateSalesInvoice(metaData, salesInvoice)
	}
	SoStatus, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "RETURNED")
	data.StatusId = &SoStatus
	if *data.StatusId == SoStatus {
		salesReturns := s.MapWithSalesReturns(data, metaData, false)
		go s.salesReturnsService.CreateSalesReturn(metaData, salesReturns)
	}
	return nil
}

func (s *service) DeleteSalesOrder(metaData core.MetaData) error {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "SALES_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales order at data level")
	}

	err := s.salesOrderRepository.Delete(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{
			"so_id": metaData.Query["id"],
		}
		err = s.salesOrderRepository.DeleteOrderLine(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteSalesOrderLines(metaData core.MetaData) error {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "SALES_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales order at data level")
	}

	err := s.salesOrderRepository.DeleteOrderLine(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ChannelSalesOrderUpsert(metaData core.MetaData, data []orders.SalesOrders) (interface{}, error) {

	var success []interface{}
	var failures []interface{}

	for index, salesOrder := range data {

		if salesOrder.ReferenceNumber == "" {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "reference_number should not be empty"})
			continue
		}
		query := map[string]interface{}{
			"reference_number": salesOrder.ReferenceNumber,
		}
		fetchedSalesOrder, err := s.salesOrderRepository.FindOne(query)
		if err == nil {
			// fetchedSalesOrder := data.(orders.SalesOrders)
			if salesOrder.StatusId != nil {
				if *salesOrder.StatusId != *fetchedSalesOrder.StatusId && *salesOrder.StatusId != 0 {
					salesOrder.StatusHistory, _ = helpers.UpdateStatusHistory(fetchedSalesOrder.StatusHistory, *salesOrder.StatusId)
				}
			}
			salesOrder.UpdatedByID = &metaData.TokenUserId
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

		salesOrder.SalesOrderNumber = helpers.GenerateSequence("SO", fmt.Sprint(metaData.TokenUserId), "sales_orders")
		salesOrder.CreatedByID = &metaData.TokenUserId
		salesOrder.CompanyId = metaData.CompanyId
		salesOrder.StatusHistory, _ = helpers.UpdateStatusHistory(salesOrder.StatusHistory, *salesOrder.StatusId)

		_, err = s.salesOrderRepository.Save(&salesOrder)
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
}

func (s *service) ChannelSalesOrderLinesUpsert(data []orders.SalesOrderLines, salesOrderId uint) (interface{}, error) {
	var success []interface{}
	var failures []interface{}

	for index, orderLine := range data {
		updateQuery := map[string]interface{}{
			"product_id": orderLine.ProductId,
			"so_id":      salesOrderId,
		}
		err := s.salesOrderRepository.UpdateOrderLines(updateQuery, &orderLine)
		if err != nil {
			orderLine.SoId = salesOrderId
			err := s.salesOrderRepository.SaveOrderLines(&orderLine)
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
	var totalQuantity int64 = 0

	for index, orderLines := range data.SalesOrderLines {
		amountWithoutDiscount := (orderLines.Price * float32(orderLines.Quantity))
		amountWithoutTax := amountWithoutDiscount - orderLines.Discount
		tax := (amountWithoutTax * orderLines.Tax) / 100
		amount := amountWithoutTax + tax

		data.SalesOrderLines[index].Amount = amount
		subTotal += amountWithoutTax
		totalTax += tax
		totalQuantity = totalQuantity + orderLines.Quantity
	}
	data.TotalQuantity = totalQuantity
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

func (s *service) GetSalesHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "LIST", "SALES_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales order at data level")
	}

	data, err := s.salesOrderRepository.GetSalesHistory(metaData.Query, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) GetSalesOrderTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "SALES_ORDERS", metaData.TokenUserId)
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

	salesOrderId, _ := strconv.Atoi(id)

	if tab == "delivery_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		deliveryOrderPage := page
		deliveryOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		data, err := s.deliveryOrdersService.AllDeliveryOrders(metaData, deliveryOrderPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "purchase_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		purchaseOrderPage := page
		purchaseOrderPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		data, err := s.purchaseOrdersService.ListPurchaseOrders(metaData, purchaseOrderPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "sales_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_INVOICE_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		salesInvoicePage := page
		salesInvoicePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		var metaData core.MetaData
		data, err := s.salesInvoiceService.GetAllSalesInvoice(metaData, salesInvoicePage)
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
		deliveryOrderData, err := s.deliveryOrdersService.AllDeliveryOrders(metaData, deliveryOrderPage)
		if err != nil {
			return nil, err
		}

		for _, deliveryOrder := range deliveryOrderData {

			sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDERS")

			salesReturnsPage := new(pagination.Paginatevalue)
			salesReturnsPage.Per_page = 1000
			salesReturnsPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrder.ID)
			salesReturnsDataInterface, err := s.salesReturnsService.ListSalesReturns(metaData, salesReturnsPage)
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
		if err != nil {
			return nil, err
		}

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")

		deliveryOrderPage := new(pagination.Paginatevalue)
		deliveryOrderPage.Per_page = 1000
		deliveryOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesOrderId)
		deliveryOrderData, err := s.deliveryOrdersService.AllDeliveryOrders(metaData, deliveryOrderPage)
		if err != nil {
			return nil, err
		}

		for _, deliveryOrder := range deliveryOrderData {

			sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_RETURNS_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDERS")

			salesReturnsPage := new(pagination.Paginatevalue)
			salesReturnsPage.Per_page = 1000
			salesReturnsPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrder.ID)
			salesReturnsDataInterface, err := s.salesReturnsService.ListSalesReturns(metaData, salesReturnsPage)
			if err != nil {
				return nil, err
			}
			salesReturnsData := salesReturnsDataInterface.([]returns.SalesReturns)

			for _, salesReturn := range salesReturnsData {

				sourceDocumentId, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_RETURS")

				creditNotePage := new(pagination.Paginatevalue)
				creditNotePage.Per_page = 1000
				creditNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesReturn.ID)

				var metaData core.MetaData
				creditNoteDataInterface, err := s.creditNoteService.GetCreditNoteList(metaData, creditNotePage)
				if err != nil {
					return nil, err
				}

				creditNoteData := creditNoteDataInterface.([]creditnote.CreditNoteResponseDTO)
				data = append(data, creditNoteData...)
			}
		}
		return data, nil
	}
	return nil, nil
}

func (s *service) TrackSalesOrder(query map[string]interface{}) []map[string]interface{} {
	var sales_order_res SalesOrdersDTO
	res, _ := s.salesOrderRepository.FindOne(query)
	dto, _ := json.Marshal(res)
	json.Unmarshal(dto, &sales_order_res)
	return sales_order_res.StatusHistory
}

func (s *service) UpdateInOndc(newOrder orders.SalesOrders, oldOrder orders.SalesOrders) (err error) {

	orderStatusMap := map[string]string{
		"CANCELLED":     "Cancelled",
		"RETURNED":      "RTO-Initiated",
		"DELIVERED":     "Completed",
		"SHIPPED":       "Order-picked-up",
		"READY_TO_SHIP": "Packed",
		"PROCESSING":    "In-progress",
		"NEW":           "Accepted",
		"DRAFT":         "Created",
	}

	cancelledStatusId, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "CANCELLED")
	partiallyCancelledStatusId, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "PARTIALLY_CANCELLED")

	if (newOrder.StatusId == nil || *newOrder.StatusId == *oldOrder.StatusId) && len(newOrder.SalesOrderLines) == 0 {
		fmt.Println("Not Produced to ONDC 1", *newOrder.StatusId, len(newOrder.SalesOrderLines))
		return nil
	}
	payload := map[string]interface{}{}

	var ondcContext map[string]interface{}
	err = json.Unmarshal(oldOrder.OndcContext, &ondcContext)
	if err != nil {
		return err
	}

	if *newOrder.StatusId != cancelledStatusId && *newOrder.StatusId != partiallyCancelledStatusId {
		orderState, err := helpers.GetLookupcodeName(*newOrder.StatusId)
		if err != nil {
			return err
		}
		ondcOrderState := orderStatusMap[orderState]
		if ondcOrderState == "" {
			return nil
		}
		ondcContext["message"].(map[string]interface{})["order"].(map[string]interface{})["state"] = ondcOrderState
		oldOrder.OndcContext, _ = json.Marshal(ondcContext)

		producePayload := map[string]interface{}{
			"meta_data": core.MetaData{
				AdditionalFields: map[string]interface{}{
					"type": "get_order_status",
				},
			},
			"data":          oldOrder,
			"order_request": ondcContext,
		}

		eda.Produce(eda.SEND_SALES_ORDER, producePayload, false)
		return nil
	}

	payload["context"] = ondcContext["context"]

	order := map[string]interface{}{
		"id": oldOrder.ReferenceNumber,
	}

	if newOrder.StatusId != nil && *newOrder.StatusId != *oldOrder.StatusId {

		orderState, err := helpers.GetLookupcodeName(*newOrder.StatusId)

		fmt.Println("orderState", orderState, *newOrder.StatusId)

		if err != nil {
			fmt.Println("Not Produced to ONDC 1", *newOrder.StatusId, len(newOrder.SalesOrderLines))
			return nil
		}

		if ondcOrderState := orderStatusMap[orderState]; ondcOrderState != "" {
			order["state"] = ondcOrderState
		}

	}

	if (newOrder.StatusId != nil && *newOrder.StatusId != cancelledStatusId) && len(newOrder.SalesOrderLines) > 0 {
		items := []map[string]interface{}{}
		for _, newItem := range newOrder.SalesOrderLines {
			for _, oldItem := range oldOrder.SalesOrderLines {
				if newItem.ProductId != nil && oldItem.ProductId != nil {
					if *newItem.ProductId == *oldItem.ProductId {
						item := map[string]interface{}{
							"id": oldItem.Product.SkuId,
						}
						if newItem.Quantity > 0 && newItem.Quantity != oldItem.Quantity {
							item["quantity"] = map[string]interface{}{
								"count": newItem.Quantity,
							}
						}
						if newItem.StatusId != nil && (*newItem.StatusId == cancelledStatusId || *newItem.StatusId == partiallyCancelledStatusId) {
							item["tags"] = map[string]interface{}{
								"update_type": "cancel",
							}
						}
						if item["quantity"] != nil || item["tags"] != nil {
							items = append(items, item)
						}
					}
				}
			}
		}
		if len(items) > 0 {
			order["items"] = items
		}
	}

	if order["items"] == nil && order["state"] == nil {
		fmt.Println("Not Produced to ONDC 2")
		return nil
	}

	if ondcContext != nil {
		order["quote"] = ondcContext["message"].(map[string]interface{})["order"].(map[string]interface{})["quote"]
	}

	payload["message"] = map[string]interface{}{
		"order": order,
	}

	producePayload := map[string]interface{}{
		"meta_data": core.MetaData{},
		"data":      payload,
	}

	// fetchedOrder, err := s.salesOrderRepository.FindOne(map[string]interface{}{"id": oldOrder.ID})
	// if err != nil {
	// 	return err
	// }
	// var orderDto SalesOrdersDTO
	// err = helpers.JsonMarshaller(fetchedOrder, &orderDto)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("--------Order Produced to ONDC-------->", oldOrder.ReferenceNumber, eda.ONDC_UPDATE_SALES_ORDER_STATUS)
	eda.Produce(eda.ONDC_UPDATE_SALES_ORDER_STATUS, producePayload, false)

	// helpers.PrettyPrint("payload", producePayload)

	return nil
}

func (s *service) MapWithSalesInvoice(requestPayload *orders.SalesOrders, edaMetaData core.MetaData, flag bool) *invoice_model.SalesInvoice {
	//================create Sales Invoice =========================
	salesInvoice := new(invoice_model.SalesInvoice)
	salesInvoice.OrderId = &requestPayload.ID
	if flag {
		delete(edaMetaData.Query, "company_id")
		so_data, _ := s.ViewSalesOrder(edaMetaData)
		json_data, _ := json.Marshal(so_data)
		json.Unmarshal(json_data, &requestPayload)
	}
	channel_name := requestPayload.ChannelName
	var channel_data channel.ChannelDTO
	if channel_name != "" {
		query := map[string]interface{}{
			"name": channel_name,
		}
		channel_data, _ = s.channelService.GetChannelService(query)
	} else {
		channel_data.ID = 1
	}
	salesInvoice.ChannelID = &channel_data.ID

	var InvoiceBillingAddress CustomerBillingorShippingAddress
	var InvoiceReceivingAddress CustomerBillingorShippingAddress
	json.Unmarshal(requestPayload.CustomerBillingAddress, &InvoiceBillingAddress)
	json.Unmarshal(requestPayload.CustomerShippingAddress, &InvoiceReceivingAddress)
	var billing_data = make([]map[string]interface{}, 0)
	var shipping_data = make([]map[string]interface{}, 0)
	billing_data = append(billing_data, map[string]interface{}{
		"sender_name":   InvoiceBillingAddress.ContactPersonName,
		"mobile_number": InvoiceBillingAddress.ContactPersonNumber,
		"address_line1": InvoiceBillingAddress.AddressLine1,
		"address_line2": InvoiceBillingAddress.AddressLine2,
		"address_line3": InvoiceBillingAddress.AddressLine3,
		"zipcode":       InvoiceBillingAddress.PinCode,
		"city":          InvoiceBillingAddress.City,
		"state":         InvoiceBillingAddress.State,
		"country":       InvoiceBillingAddress.Country,
		"landmark":      InvoiceBillingAddress.Landmark,
	})
	salesInvoice.BillingAddress, _ = json.Marshal(billing_data)
	shipping_data = append(shipping_data, map[string]interface{}{
		"receiver_name": InvoiceReceivingAddress.ContactPersonName,
		"mobile_number": InvoiceReceivingAddress.ContactPersonNumber,
		"address_line1": InvoiceReceivingAddress.AddressLine1,
		"address_line2": InvoiceReceivingAddress.AddressLine2,
		"address_line3": InvoiceReceivingAddress.AddressLine3,
		"zipcode":       InvoiceReceivingAddress.PinCode,
		"city":          InvoiceReceivingAddress.City,
		"state":         InvoiceReceivingAddress.State,
		"country":       InvoiceReceivingAddress.Country,
		"landmark":      InvoiceReceivingAddress.Landmark,
	})
	salesInvoice.ShippingAddress, _ = json.Marshal(shipping_data)

	salesInvoice.PaymentTypeID = requestPayload.PaymentTypeId
	var InvoiceTotalQuantity int
	for _, sales_order_line := range requestPayload.SalesOrderLines {
		InvoiceTotalQuantity += int(sales_order_line.Quantity)
		SalesInvoiceLine := new(invoice_model.SalesInvoiceLines)
		SalesInvoiceLine.ProductVariantID = sales_order_line.ProductId
		SalesInvoiceLine.ProductID = sales_order_line.ProductTemplateId
		SalesInvoiceLine.Quantity = int32(sales_order_line.Quantity)
		SalesInvoiceLine.Tax = float64(sales_order_line.Tax)
		SalesInvoiceLine.Discount = float64(sales_order_line.Discount)
		SalesInvoiceLine.Price = float64(sales_order_line.Price)
		SalesInvoiceLine.TotalAmount = float64(sales_order_line.Amount)
		salesInvoice.SalesInvoiceLines = append(salesInvoice.SalesInvoiceLines, *SalesInvoiceLine)
	}
	salesInvoice.Quantity = int32(InvoiceTotalQuantity)

	salesInvoice.TotalAmount = float64(requestPayload.SoPaymentDetails.TotalAmount)
	salesInvoice.SubTotalAmount = float64(requestPayload.SoPaymentDetails.SubTotal)
	salesInvoice.ExpectedShipmentDate = requestPayload.ExpectedShippingDate
	salesInvoice.OrderDate = requestPayload.CreatedDate
	salesInvoice.InternalNotes = requestPayload.AdditionalInformation.Notes
	salesInvoice.SalesInvoiceNumber = helpers.GenerateSequence("SALESINVOICE", fmt.Sprint(edaMetaData.TokenUserId), "sales_invoices")
	salesInvoice.ReferenceNumber = helpers.GenerateSequence("REF", fmt.Sprint(edaMetaData.TokenUserId), "sales_invoices")

	return salesInvoice

}

func (s *service) MapWithSalesReturns(requestPayload *orders.SalesOrders, edaMetaData core.MetaData, flag bool) *returns.SalesReturns {
	//================create Sales Returns =========================
	salesReturns := new(returns.SalesReturns)
	//salesReturns.OrderId = &requestPayload.ID
	if flag {
		delete(edaMetaData.Query, "company_id")
		so_data, _ := s.ViewSalesOrder(edaMetaData)
		json_data, _ := json.Marshal(so_data)
		json.Unmarshal(json_data, &requestPayload)
	}
	salesReturns.CustomerName = requestPayload.CustomerName
	salesReturns.ChannelName = requestPayload.ChannelName

	var returnBillingAddress CustomerBillingorShippingAddress
	var returnReceivingAddress CustomerBillingorShippingAddress
	json.Unmarshal(requestPayload.CustomerBillingAddress, &returnBillingAddress)
	json.Unmarshal(requestPayload.CustomerShippingAddress, &returnReceivingAddress)
	var billing_data = make([]map[string]interface{}, 0)
	var shipping_data = make([]map[string]interface{}, 0)
	billing_data = append(billing_data, map[string]interface{}{
		"sender_name":   returnBillingAddress.ContactPersonName,
		"mobile_number": returnBillingAddress.ContactPersonNumber,
		"address_line1": returnBillingAddress.AddressLine1,
		"address_line2": returnBillingAddress.AddressLine2,
		"address_line3": returnBillingAddress.AddressLine3,
		"zipcode":       returnBillingAddress.PinCode,
		"city":          returnBillingAddress.City,
		"state":         returnBillingAddress.State,
		"country":       returnBillingAddress.Country,
		"landmark":      returnBillingAddress.Landmark,
	})
	salesReturns.CustomerBillingAddress, _ = json.Marshal(billing_data)
	shipping_data = append(shipping_data, map[string]interface{}{
		"receiver_name": returnReceivingAddress.ContactPersonName,
		"mobile_number": returnReceivingAddress.ContactPersonNumber,
		"address_line1": returnReceivingAddress.AddressLine1,
		"address_line2": returnReceivingAddress.AddressLine2,
		"address_line3": returnReceivingAddress.AddressLine3,
		"zipcode":       returnReceivingAddress.PinCode,
		"city":          returnReceivingAddress.City,
		"state":         returnReceivingAddress.State,
		"country":       returnReceivingAddress.Country,
		"landmark":      returnReceivingAddress.Landmark,
	})
	salesReturns.CustomerPickupAddress, _ = json.Marshal(shipping_data)

	var returnTotalQuantity int
	for _, sales_order_line := range requestPayload.SalesOrderLines {
		returnTotalQuantity += int(sales_order_line.Quantity)
		salesReturnline := new(returns.SalesReturnLines)
		salesReturnline.ProductId = sales_order_line.ProductId
		salesReturnline.UomId = sales_order_line.UomId
		salesReturnline.ProductTemplateId = sales_order_line.ProductTemplateId
		salesReturnline.SerialNumber = sales_order_line.SerialNumber
		salesReturnline.QuantityReturned = int64(sales_order_line.Quantity)
		salesReturnline.Tax = float32(sales_order_line.Tax)
		salesReturnline.Discount = float32(sales_order_line.Discount)
		salesReturnline.Rate = int(sales_order_line.Price)
		salesReturnline.Amount = float32(sales_order_line.Amount)
		salesReturns.SalesReturnLines = append(salesReturns.SalesReturnLines, *salesReturnline)
	}
	salesReturns.TotalQuantity = int64(returnTotalQuantity)

	salesReturns.SrPaymentDetails.TotalAmount = float32(requestPayload.SoPaymentDetails.TotalAmount)
	salesReturns.SrPaymentDetails.SubTotal = float32(requestPayload.SoPaymentDetails.SubTotal)
	salesReturns.ExpectedDeliveyDate = requestPayload.ExpectedShippingDate
	salesReturns.OrderDate = time.Time(requestPayload.SoDate)
	salesReturns.AdditionalInformation.Notes = requestPayload.AdditionalInformation.Notes
	salesReturns.AdditionalInformation.TermsAndConditions = requestPayload.AdditionalInformation.TermsAndConditions
	salesReturns.AdditionalInformation.Attachments = requestPayload.AdditionalInformation.Attachments
	salesReturns.SalesReturnNumber = helpers.GenerateSequence("SR", fmt.Sprint(edaMetaData.TokenUserId), "sales_returns")
	salesReturns.ReferenceNumber = helpers.GenerateSequence("REF", fmt.Sprint(edaMetaData.TokenUserId), "sales_returns")

	return salesReturns

}
