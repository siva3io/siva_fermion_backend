package delivery_orders

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/inventory_orders/asn"
	"fermion/backend_core/controllers/orders/internal_transfers"
	"fermion/backend_core/controllers/orders/purchase_orders"
	"fermion/backend_core/controllers/returns/purchase_returns"
	shipping_order "fermion/backend_core/controllers/shipping/shipping_orders"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	orders_repo "fermion/backend_core/internal/repository/orders"

	"gorm.io/datatypes"

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
	CreateDeliveryOrder(data *DeliveryOrderRequest, token_id string, access_template_id string) (uint, error)
	BulkCreateDeliveryOrder(data *[]DeliveryOrderRequest, token_id string, access_template_id string) error
	UpdateDeliveryOrder(id uint, data *DeliveryOrderRequest, token_id string, access_template_id string) error
	FindDeliveryOrder(id uint, token_id string, access_template_id string) (interface{}, error)
	AllDeliveryOrders(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]DeliveryOrderResponse, error)
	DeleteDeliveryOrder(id uint, user_id uint, token_id string, access_template_id string) error
	DeleteDeliveryOrderLines(query interface{}, token_id string, access_template_id string) error
	ProcessDeliveryOrderCalculation(data *DeliveryOrderRequest) *DeliveryOrderRequest
	GetDeliveryOrderTab(id string, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	doRepository           orders_repo.DeliveryOrders
	shipping_order_service shipping_order.Service
	ISTService             internal_transfers.Service
	purchaseOrdersService  purchase_orders.Service
	ASNSerivce             asn.Service
	purchaseReturnsService purchase_returns.Service
}

func NewService() *service {
	return &service{orders_repo.NewDo(), shipping_order.NewService(), internal_transfers.NewService(), purchase_orders.NewService(), asn.NewService(), purchase_returns.NewService()}
}

func (s *service) CreateDeliveryOrder(data *DeliveryOrderRequest, token_id string, access_template_id string) (uint, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create delivery order at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create delivery order at data level")
	}
	if data.IsShipping != nil && *data.IsShipping {
		data.ShippingDetails.CreatedByID = data.CreatedByID
		res, err := s.shipping_order_service.CreateShippingOrder(&data.ShippingDetails, token_id, access_template_id)
		if err != nil {
			return 0, err
		}
		data.Shipping_order_id = &res.ID
	}
	data = s.ProcessDeliveryOrderCalculation(data)
	defaultStatus, err := helpers.GetLookupcodeId("DELIVERY_ORDER_STATUS", "DRAFT")
	if err != nil {
		return 0, err
	}
	data.Status_id = defaultStatus
	var DoData orders.DeliveryOrders
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(dto, &DoData)
	if err != nil {
		return 0, err
	}

	// var fetchData orders.DeliveryOrders
	if data.DeliveryOrdersDetails.AutoCreateDoNumber {
		DoData.DeliveryOrdersDetails.Delivery_order_number = helpers.GenerateSequence("DO", token_id, "delivery_orders")
	}
	if data.DeliveryOrdersDetails.AutoGenerateReferenceNumber {
		DoData.DeliveryOrdersDetails.Reference_id = helpers.GenerateSequence("REF", token_id, "delivery_orders")
	}

	result, _ := helpers.UpdateStatusHistory(DoData.Status_history, DoData.Status_id)
	DoData.Status_history = result

	id, err := s.doRepository.CreateDeliveryOrder(&DoData)
	return id, err
	//}
	//return 0, nil
}

func (s *service) UpdateDeliveryOrder(id uint, data *DeliveryOrderRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update delivery order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update delivery order at data level")
	}

	data = s.ProcessDeliveryOrderCalculation(data)

	var DoData orders.DeliveryOrders
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &DoData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	old_data, er := s.doRepository.FindDeliveryOrder(id)
	if er != nil {
		return er
	}

	old_status := old_data.Status_id
	new_status := DoData.Status_id

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.Status_history, DoData.Status_id)
		DoData.Status_history = result
	}

	err = s.doRepository.UpdateDeliveryOrder(id, &DoData)
	if err != nil {
		return err
	}
	for _, order_line := range DoData.DeliveryOrderLines {
		query := map[string]interface{}{"do_id": id, "product_id": order_line.Product_id}
		count, er := s.doRepository.UpdateDeliveryOrderLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.DO_id = id
			e := s.doRepository.CreateDeliveryOrderLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
	//}
	//return nil
}

func (s *service) FindDeliveryOrder(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view delivery order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view delivery order at data level")

	}
	result, er := s.doRepository.FindDeliveryOrder(id)
	if er != nil {
		return result, er
	}
	if result.Shipping_order_id != nil && *result.Shipping_order_id != 0 {
		shipping_order_response, _ := s.shipping_order_service.GetShippingOrder(*result.Shipping_order_id, token_id, access_template_id)
		var shipping_order_data datatypes.JSON
		dto, _ := json.Marshal(shipping_order_response)
		json.Unmarshal(dto, &shipping_order_data)
	}
	var response DeliveryOrderResponse
	mdata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(mdata, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
	//}
	//return orders.DeliveryOrders{}, nil
}

func (s *service) AllDeliveryOrders(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]DeliveryOrderResponse, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list delivery order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list delivery order at data level")
	}
	result, err := s.doRepository.AllDeliveryOrders(p)
	var response []DeliveryOrderResponse
	for _, data := range result {
		var scraporderdata DeliveryOrderResponse
		mdata, _ := json.Marshal(data)
		json.Unmarshal(mdata, &scraporderdata)
		response = append(response, scraporderdata)
	}
	return response, err
	//}
	//return nil, nil
}

func (s *service) DeleteDeliveryOrder(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete delivery order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete delivery order at data level")
	}
	_, er := s.doRepository.FindDeliveryOrder(id)
	if er != nil {
		return er
	}
	err1 := s.doRepository.DeleteDeliveryOrder(id, user_id)
	if err1 != nil {
		return err1
	}
	query := map[string]interface{}{"do_id": id}
	err := s.doRepository.DeleteDeliveryOrderLines(query)
	return err
	//}
	//return nil
}

func (s *service) DeleteDeliveryOrderLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete delivery order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete delivery order at data level")
	}
	_, e := s.doRepository.FindDeliveryOrdersLines(query)
	if e != nil {
		return e
	}

	er := s.doRepository.DeleteDeliveryOrderLines(query)
	return er
	//}
	//return nil
}

func (s *service) BulkCreateDeliveryOrder(data *[]DeliveryOrderRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "DELIVERY_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create delivery order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create delivery order at data level")
	}

	var DoData []orders.DeliveryOrders
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &DoData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = s.doRepository.BulkCreateDeliveryOrder(&DoData)
	return err
	//}
	//return nil

}

func (s *service) ProcessDeliveryOrderCalculation(data *DeliveryOrderRequest) *DeliveryOrderRequest {
	var subTotal, totalTax, shippingCharges float64

	for index, orderLines := range data.DeliveryOrderLines {
		amountWithoutDiscount := (orderLines.Rate * float64(orderLines.Product_qty))
		amountWithoutTax := amountWithoutDiscount - orderLines.Discount
		tax := (amountWithoutTax * orderLines.Tax) / 100
		amount := amountWithoutTax + tax

		data.DeliveryOrderLines[index].Amount = amount
		subTotal += amountWithoutTax
		totalTax += tax
	}

	// if estimatedCost, ok := data.ShippingDetails["estimated_cost"].([]interface{}); ok {
	// 	if estimatedCostObj := estimatedCost[0]; len(estimatedCost) > 0 {
	// 		if estimatedCostObj, ok := estimatedCostObj.(map[string]interface{}); ok {
	// 			if cost, ok := estimatedCostObj["cost"].((float64)); ok {
	// 				shippingCharges = cost
	// 			}
	// 		}
	// 	}
	// }
	data.Payment_details.Sub_total = subTotal
	data.Payment_details.Payment_tax = totalTax
	data.Payment_details.Shipping_charge = data.ShippingDetails.ShippingCost
	data.Payment_details.Total_amount = subTotal + totalTax + shippingCharges

	return data

}

func (s *service) GetDeliveryOrderTab(id string, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

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

	if tab == "ist" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("IST_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDERS")

		deliveryOrdersPage := page
		fmt.Println(sourceDocumentId, deliveryOrderId)
		deliveryOrdersPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		fmt.Println(deliveryOrdersPage)
		data, err := s.ISTService.ListInternalTransfers(deliveryOrdersPage, access_template_id, token_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "asn" {
		var data []inventory_orders.ASN
		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDER")
		deliveryReturnPage := page
		deliveryReturnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		PurchaseOrdersInterface, err := s.purchaseOrdersService.ListPurchaseOrders(deliveryReturnPage, access_template_id, token_id, "LIST")
		PurchaseOrdersData := PurchaseOrdersInterface.([]orders.PurchaseOrders)
		if err != nil {
			return nil, err
		}
		for _, POValue := range PurchaseOrdersData {

			sourceDocumentId, _ := helpers.GetLookupcodeId("ASN_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

			AsnPage := new(pagination.Paginatevalue)
			AsnPage.Per_page = 1000
			AsnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, POValue.ID)
			data, err = s.ASNSerivce.GetAllAsn(AsnPage, token_id, access_template_id, "LIST")
			// fmt.Println(data)
			if err != nil {
				return nil, err
			}
		}
		return data, nil
	}

	if tab == "purchase_returns" {
		var data interface{}
		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES", "DELIVERY_ORDER")
		deliveryReturnPage := page
		deliveryReturnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, deliveryOrderId)
		PurchaseOrdersInterface, err := s.purchaseOrdersService.ListPurchaseOrders(deliveryReturnPage, access_template_id, token_id, "LIST")
		PurchaseOrdersData := PurchaseOrdersInterface.([]orders.PurchaseOrders)
		if err != nil {
			return nil, err
		}
		for _, POValue := range PurchaseOrdersData {

			sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_RETURNS_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

			purchaseReturnsPage := new(pagination.Paginatevalue)
			purchaseReturnsPage.Per_page = 1000
			purchaseReturnsPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, POValue.ID)
			data, err = s.purchaseReturnsService.ListPurchaseReturns(purchaseReturnsPage, access_template_id, token_id, "LIST")
			if err != nil {
				return nil, err
			}
		}
		return data, nil
	}
	return nil, nil
}
