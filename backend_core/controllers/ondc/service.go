package ondc

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fermion/backend_core/controllers/cores"
	"fermion/backend_core/controllers/mdm/products"
	"fermion/backend_core/controllers/orders/sales_orders"
	"fermion/backend_core/pkg/util/helpers"

	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/repository"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	order_repo "fermion/backend_core/internal/repository/orders"

	// ipaas_core "fermion/backend_core/ipaas_core/app"
	"github.com/google/uuid"
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
	// Search(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{}
	// Cancel(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{}
	// Select(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{}
	// Confirm(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{}
	// Track(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{}
	// Update(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{}
	// Status(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{}

	BuildOndcPayload(company cores.Company, action string) map[string]interface{}

	OnSearch(data BPPRequest) []products.VariantResponseDTO
	OnCancel(data BPPRequest) *sales_orders.SalesOrdersDTO
	OnSelect(data BPPRequest) products.VariantResponseDTO
	OnInit(data BPPRequest) map[string]interface{}
	OnConfirm(data BPPRequest) *sales_orders.SalesOrdersDTO
	OnTrack(data BPPRequest) map[string]interface{}
	OnUpdate(data BPPRequest) *sales_orders.SalesOrdersDTO
	OnStatus(data BPPRequest) *sales_orders.SalesOrdersDTO
}

type service struct {
	userRepository       repository.User
	coreRepository       repository.Core
	productRepository    mdm_repo.Products
	salesOrderRepository order_repo.SalesOrders
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	userRepository := repository.NewUser()
	coreRepository := repository.NewCore()
	productRepository := mdm_repo.NewProducts()
	salesOrderRepository := *order_repo.NewSalesOrder()

	newServiceObj = &service{userRepository, coreRepository, productRepository, salesOrderRepository}
	return newServiceObj
}

// ============================================ BAP Worflow APIs ======================================================== //

// func (s *service) Search(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{} {
// 	request_payload := map[string]interface{}{
// 		"ondc_details": ondc_details,
// 		"payload":      data,
// 	}
// 	// "is_store_pickup":true,
// 	// "is_home_delivery":true,
// 	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
// 	return request_payload
// }

// func (s *service) Cancel(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{} {
// 	request_payload := map[string]interface{}{
// 		"ondc_details": ondc_details,
// 		"payload":      data,
// 	}
// 	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
// 	return request_payload
// }

// func (s *service) Select(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{} {
// 	request_payload := map[string]interface{}{
// 		"ondc_details": ondc_details,
// 		"payload":      data,
// 	}
// 	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
// 	return request_payload
// }

// func (s *service) Confirm(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{} {
// 	request_payload := map[string]interface{}{
// 		"ondc_details": ondc_details,
// 		"payload":      data,
// 	}
// 	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
// 	return request_payload
// }

// func (s *service) Track(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{} {
// 	request_payload := map[string]interface{}{
// 		"ondc_details": ondc_details,
// 		"payload":      data,
// 	}
// 	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
// 	return request_payload
// }

// func (s *service) Update(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{} {
// 	request_payload := map[string]interface{}{
// 		"ondc_details": ondc_details,
// 		"payload":      data,
// 	}
// 	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
// 	return request_payload
// }

// func (s *service) Status(ondc_details map[string]interface{}, data map[string]interface{}) map[string]interface{} {
// 	request_payload := map[string]interface{}{
// 		"ondc_details": ondc_details,
// 		"payload":      data,
// 	}
// 	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
// 	return request_payload
// }

func (s *service) BuildOndcPayload(company cores.Company, action string) map[string]interface{} {
	ondc_details := map[string]interface{}{
		"domain":         company.CompanyDetails.Domain.LookupCode,
		"country":        "IND",
		"city":           company.CompanyDetails.StdCode.LookupCode,
		"action":         action,
		"transaction_id": uuid.New(),
		"message_id":     uuid.New(),
		"timestamp":      time.Now(),
		fmt.Sprintf("%s_id", strings.ToLower(company.CompanyType.LookupCode)):  company.OndcDetails.SubscriberId,
		fmt.Sprintf("%s_url", strings.ToLower(company.CompanyType.LookupCode)): company.OndcDetails.SubscriberURL,
	}
	return ondc_details
}

// ============================================ BPP Worflow APIs ======================================================== //

func (s *service) OnSearch(data BPPRequest) []products.VariantResponseDTO {
	intent := data.Message["intent"]
	var search_type string
	query := make(map[string]interface{}, 0)

	if price, ok := intent.(map[string]interface{})["item"].(map[string]interface{})["price"]; ok {
		query["minimum_value"] = price.(map[string]interface{})["minimum_value"]
		query["maximum_value"] = price.(map[string]interface{})["maximum_value"]
		search_type = "price"

	} else if name, ok := intent.(map[string]interface{})["item"].(map[string]interface{})["descriptor"].(map[string]interface{})["name"]; ok {
		query["name"] = name
		search_type = "product_name"

	} else if sku_id, ok := intent.(map[string]interface{})["item"].(map[string]interface{})["descriptor"].(map[string]interface{})["code"]; ok {
		query["sku_id"] = sku_id
		search_type = "sku_id"

	}

	result, err := s.productRepository.GetAllProducts(query, search_type)
	if err != nil {
		fmt.Println(err)
	}
	var response []products.VariantResponseDTO
	dto, _ := json.Marshal(result)
	_ = json.Unmarshal(dto, &response)

	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
	return response
}

func (s *service) OnCancel(data BPPRequest) *sales_orders.SalesOrdersDTO {
	var query = map[string]interface{}{
		"sales_order_number": data.Message["order_id"],
	}
	var sales_order_data orders.SalesOrders
	statusId, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "CANCELLED")
	sales_order_data.StatusId = &statusId
	err := s.salesOrderRepository.Update(query, &sales_order_data)
	if err != nil {
		return nil
	}
	res, _ := s.salesOrderRepository.FindOne(query)
	var sales_order_res sales_orders.SalesOrdersDTO
	dto, _ := json.Marshal(res)
	json.Unmarshal(dto, &sales_order_res)
	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
	return &sales_order_res
}

func (s *service) OnSelect(data BPPRequest) products.VariantResponseDTO {
	fmt.Println("data.Message['order']", data.Message["order"])
	items := data.Message["order"].(map[string]interface{})["items"].([]interface{})
	var item_id interface{}
	var result mdm.ProductVariant

	for _, value := range items {
		item_id = value.(map[string]interface{})["id"]
		result, _ = s.productRepository.FindVariant(map[string]interface{}{"id": item_id})
	}
	var response products.VariantResponseDTO
	dto, _ := json.Marshal(result)
	_ = json.Unmarshal(dto, &response)

	return response
}

func (s *service) OnInit(data BPPRequest) map[string]interface{} {
	order := data.Message["order"]
	// items := order.(map[string]interface{})["items"].([]interface{})
	// billingAddress := order.(map[string]interface{})["billing"]
	// fulfillment := order.(map[string]interface{})["fulfillment"]
	// startLocation := fulfillment.(map[string]interface{})["start"]
	// endLocation := fulfillment.(map[string]interface{})["end"]
	// fmt.Println(items, billingAddress, startLocation, endLocation)
	// err := ipaas_core.NewIpaasHandler().RateCalculator()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//-----> Need to write a function for getting RATES for user ----//

	return order.(map[string]interface{})
}

func (s *service) OnConfirm(data BPPRequest) *sales_orders.SalesOrdersDTO {
	order := data.Message["order"]
	items := order.(map[string]interface{})["items"].([]interface{})
	billingAddress := order.(map[string]interface{})["billing"]
	fulfillment := order.(map[string]interface{})["fulfillment"]
	startLocation := fulfillment.(map[string]interface{})["start"]
	endLocation := fulfillment.(map[string]interface{})["end"]
	fmt.Println(items, billingAddress, startLocation, endLocation)

	var result mdm.ProductVariant
	var sales_order_res sales_orders.SalesOrdersDTO

	for _, value := range items {
		item_id := value.(map[string]interface{})["id"]
		result, _ = s.productRepository.FindVariant(map[string]interface{}{"id": item_id})
		// Need to create sales order --->(Implementing mapper file for converting payload into create sales_order payload)
	}
	fmt.Println(result)
	// ...........after completion of mapper need to return created sales order response..............
	res, _ := s.salesOrderRepository.FindOne(map[string]interface{}{"sales_order_number": "SO/1/0002"})
	dto, _ := json.Marshal(res)
	json.Unmarshal(dto, &sales_order_res)
	// ...........after mapper completion need to return created sales order response..............

	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
	return &sales_order_res
}

func (s *service) OnTrack(data BPPRequest) map[string]interface{} {
	res, err := s.salesOrderRepository.FindOne(map[string]interface{}{"sales_order_number": data.Message["order_id"]})
	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{
			"error": "Sales order not found",
		}
	}
	status, _ := helpers.GetLookupcodeName(int(*res.StatusId))

	response := map[string]interface{}{
		"tl_method":    "https/get",
		"status":       status,
		"tracking_url": fmt.Sprintf("https://dev-api.eunimart.com/api/v1/sales_orders/track?sales_order_number=%v", data.Message["order_id"]),
	}

	// Sandbox API call -->[NEED TO BE IMPLEMENTED]

	return response
}

func (s *service) OnStatus(data BPPRequest) *sales_orders.SalesOrdersDTO {
	var sales_order_res sales_orders.SalesOrdersDTO
	res, _ := s.salesOrderRepository.FindOne(map[string]interface{}{"sales_order_number": data.Message["order_id"]})
	dto, _ := json.Marshal(res)
	err := json.Unmarshal(dto, &sales_order_res)
	if err != nil {
		fmt.Println(err)
	}

	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
	return &sales_order_res
}

func (s *service) OnUpdate(data BPPRequest) *sales_orders.SalesOrdersDTO {
	var sales_order_res sales_orders.SalesOrdersDTO
	res, _ := s.salesOrderRepository.FindOne(map[string]interface{}{"sales_order_number": data.Message["order"].(map[string]interface{})["id"]})
	dto, _ := json.Marshal(res)
	err := json.Unmarshal(dto, &sales_order_res)
	if err != nil {
		fmt.Println(err)
	}

	// Sandbox API call -->[NEED TO BE IMPLEMENTED]
	return &sales_order_res
}
