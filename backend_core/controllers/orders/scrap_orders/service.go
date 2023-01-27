package scrap_orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	// "fermion/backend_core/controllers/inventory_orders/grn"
	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
	shipping_order "fermion/backend_core/controllers/shipping/shipping_orders"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
	inventory_order_repo "fermion/backend_core/internal/repository/inventory_orders"
	orders_repo "fermion/backend_core/internal/repository/orders"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"gorm.io/datatypes"
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
	CreateScrapOrder(metaData core.MetaData, data *ScrapOrders) error
	BulkCreateScrapOrder(data *[]ScrapOrders, token_id string, access_template_id string) error
	UpdateScrapOrder(metaData core.MetaData, data *ScrapOrders) error
	AllScrapOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	FindScrapOrder(metaData core.MetaData) (interface{}, error)
	DeleteScrapOrder(metaData core.MetaData) error
	DeleteScrapOrderLines(metaData core.MetaData) error
	GetScrapOrderTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	scrapRepository        orders_repo.ScrapOrders
	basicInventoryService  inv_service.Service
	shipping_order_service shipping_order.Service
	grn_repositary         inventory_order_repo.GRN
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	shipping_order_service := shipping_order.NewService()
	grn_service := inventory_order_repo.NewGRN()
	newServiceObj = &service{orders_repo.NewScrap(), inv_service.NewService(), shipping_order_service, grn_service}
	return newServiceObj
}

func (s *service) CreateScrapOrder(metaData core.MetaData, data *ScrapOrders) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "CREATE", "SCRAP_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create scrap order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create scrap order at data level")
	}
	createPayload := new(shipping.ShippingOrder)
	helpers.JsonMarshaller(data, createPayload)
	if data.IsShipping != nil && *data.IsShipping {
		data.Shipping_details.CreatedByID = data.CreatedByID
		err := s.shipping_order_service.CreateShippingOrder(metaData, createPayload)
		if err != nil {
			return err
		}
		// data.ShippingOrderId = &res.ID
	}

	var ScrapData *orders.ScrapOrders
	helpers.JsonMarshaller(data, &ScrapData)
	defaultStatus, err := helpers.GetLookupcodeId("SCRAP_ORDER_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	ScrapData.Status_id = &defaultStatus
	// var fetchData orders.DeliveryOrders
	if data.AutoCreateScrapNumber {
		ScrapData.Scrap_order_number = helpers.GenerateSequence("SCRAP", fmt.Sprint(metaData.TokenUserId), "scrap_orders")
	}
	if data.AutoGenerateReferenceNumber {
		ScrapData.Reference_id = helpers.GenerateSequence("REF", fmt.Sprint(metaData.TokenUserId), "scrap_orders")
	}
	ScrapData.CompanyId = metaData.CompanyId
	ScrapData.CreatedByID = &metaData.TokenUserId
	result, _ := helpers.UpdateStatusHistory(ScrapData.Status_history, *ScrapData.Status_id)
	ScrapData.Status_history = result
	_, err = s.scrapRepository.SaveScrapOrder(ScrapData)
	if err != nil {
		return err
	}
	data_inv := map[string]interface{}{
		"id":            ScrapData.ID,
		"order_lines":   ScrapData.Order_lines,
		"order_type":    "scrap_order",
		"is_update_inv": true,
		"is_credit":     false,
	}

	//updating the inventory & Committed Stock
	go s.basicInventoryService.UpdateTransactionInventory(metaData, data_inv)

	return nil
}

func (s *service) UpdateScrapOrder(metaData core.MetaData, data *ScrapOrders) error {
	//accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "UPDATE", "SCRAP_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update scrap order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update scrap order at data level")
	}
	var ScrapData *orders.ScrapOrders
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &ScrapData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	old_data, er := s.scrapRepository.FindScrapOrderById(metaData.Query)
	if er != nil {
		return er
	}

	old_status := old_data.Status_id
	new_status := data.Status_id

	if new_status != old_status && new_status != nil {
		result, _ := helpers.UpdateStatusHistory(old_data.Status_history, *ScrapData.Status_id)
		ScrapData.Status_history = result
		final_status, _ := helpers.GetLookupcodeId("SCRAP_ORDER_STATUS", "SHIPPED_AND_RECEIVED")
		partial_status1, _ := helpers.GetLookupcodeId("SCRAP_ORDER_STATUS", "SHIPPED_AND_PARTIALLY_RECEIVED")
		partial_status, _ := helpers.GetLookupcodeId("SCRAP_ORDER_STATUS", "PARTIALLY_SHIPPED_AND_PARTIALLY_RECEIVED")
		if *ScrapData.Status_id == final_status || *ScrapData.Status_id == partial_status || *ScrapData.Status_id == partial_status1 {
			data_inv := map[string]interface{}{
				"id":            ScrapData.ID,
				"order_lines":   ScrapData.Order_lines,
				"order_type":    "scrap_order",
				"is_update_inv": false,
				"is_credit":     false,
			}

			//updating the Committed Stock
			go s.basicInventoryService.UpdateTransactionInventory(metaData, data_inv)
		}
	}
	ScrapData.UpdatedByID = &metaData.TokenUserId
	err = s.scrapRepository.UpdateScrapOrder(metaData.Query, ScrapData)

	if err != nil {
		return err
	}

	for _, order_line := range ScrapData.Order_lines {
		query := map[string]interface{}{"scrap_id": metaData.Query["id"], "product_id": order_line.Product_id}
		count, err1 := s.scrapRepository.UpdateScrapOrderLines(query, order_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			order_line.Scrap_id = uint(metaData.Query["id"].(float64))
			e := s.scrapRepository.CreateScrapOrderLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
	//}
	//return nil
}

func (s *service) AllScrapOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "SCRAP_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list scrap order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list scrap order at data level")
	}
	data, err := s.scrapRepository.FindAllScrapOrders(metaData.Query, page)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) FindScrapOrder(metaData core.MetaData) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "SCRAP_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view scrap order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view scrap order at data level")
	}
	result, err := s.scrapRepository.FindScrapOrderById(metaData.Query)

	if err != nil {
		return result, err
	}

	if result.ShippingOrderId != nil && *result.ShippingOrderId != 0 {
		shipping_order_response, _ := s.shipping_order_service.GetShippingOrder(metaData)
		var shipping_order_data datatypes.JSON
		dto, _ := json.Marshal(shipping_order_response)
		json.Unmarshal(dto, &shipping_order_data)
	}

	var response ScrapOrdersResponseDTO
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
	//return orders.ScrapOrders{}, nil
}

func (s *service) DeleteScrapOrder(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "SCRAP_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete scrap order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete scrap order at data level")
	}
	err := s.scrapRepository.DeleteScrapOrder(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{
			"scrap_id": metaData.Query["id"],
		}
		err = s.scrapRepository.DeleteScrapOrderLines(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteScrapOrderLines(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "SCRAP_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase order at data level")
	}
	err := s.scrapRepository.DeleteScrapOrderLines(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) BulkCreateScrapOrder(data *[]ScrapOrders, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SCRAP_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create scrap order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create scrap order at data level")
	}
	var ScrapData *[]orders.ScrapOrders
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &ScrapData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = s.scrapRepository.BulkSaveScrapOrder(ScrapData)
	return err
	//}
	//return nil
}

func (s *service) GetScrapOrderTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "LIST", "SCRAP_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view scrap order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view scrap order at data level")
	}
	id := metaData.Query["id"].(string)
	tab := metaData.Query["tab"].(string)
	// token_id := fmt.Sprint(metaData.TokenUserId)
	//access_template_id := fmt.Sprint(metaData.AccessTemplateId)
	scrapOrderId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "grn" {
		scrapPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("SCRAP_ORDERS_SOURCE_DOCUMENT_TYPES", "GRN")
		scrapPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", scrapOrderId, source_document_type_id)
		scrap_order_response, err := s.AllScrapOrders(metaData, scrapPagination)
		if err != nil {
			return nil, err
		}
		var scrap_order []ScrapOrders
		helpers.JsonMarshaller(scrap_order_response, &scrap_order)

		if len(scrap_order) == 0 {
			return nil, errors.New("noo source document")
		}

		var scrapOrderSourceDoc map[string]interface{}
		scrapOrderSourceDocJson := scrap_order[0].SourceDocuments
		dto, err := json.Marshal(scrapOrderSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &scrapOrderSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		scrapOrderSourceDocId := scrapOrderSourceDoc["id"]
		if scrapOrderSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		grnPage := page
		grnPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", scrapOrderSourceDoc["id"])
		data, err := s.grn_repositary.GetAllGRN(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	return nil, nil
}
