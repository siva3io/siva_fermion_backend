package basic_inventory

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"fermion/backend_core/controllers/mdm/locations"
	"fermion/backend_core/controllers/mdm/products"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	access_checker "fermion/backend_core/pkg/util/access"
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
	CreateCentrailizedInventory(metaData core.MetaData, data *mdm.CentralizedBasicInventory) error
	CreateDecentralizedInventory(metaData core.MetaData, data *mdm.DecentralizedBasicInventory) error

	UpdateCentrailizedInventory(metaData core.MetaData, data *mdm.CentralizedBasicInventory) error
	UpdateDecentralizedInventory(metaData core.MetaData, data *mdm.DecentralizedBasicInventory) error

	DeleteCentrailizedInventory(metaData core.MetaData) error
	DeleteDecentralizedInventory(metaData core.MetaData) error

	GetCentrailizedInventory(metaData core.MetaData) (interface{}, error)
	GetDecentralizedInventory(metaData core.MetaData) (interface{}, error)

	GetCentrailizedInventoryList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	GetDecentralizedInventoryList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)

	UpsertInventoryTemplate(metaData core.MetaData, data []interface{}) (interface{}, error)

	InventoryTransactionCreate(metaData core.MetaData, data *mdm.CentralizedInventoryTransactions) error
	UpdateTransactionInventory(metaData core.MetaData, data map[string]interface{}) error
}

type service struct {
	basicInventory  mdm_repo.BasicInventory
	productService  products.Service
	locationService locations.Service
	mu              sync.Mutex
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}

	basicInventory := mdm_repo.NewBasicInventory()
	productService := products.NewService()
	locationService := locations.NewService()
	newServiceObj = &service{
		basicInventory:  basicInventory,
		productService:  productService,
		locationService: locationService,
	}
	return newServiceObj
}

func (s *service) CreateCentrailizedInventory(metaData core.MetaData, data *mdm.CentralizedBasicInventory) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create inventory at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.basicInventory.CentralizedInventorySave(data)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) CreateDecentralizedInventory(metaData core.MetaData, data *mdm.DecentralizedBasicInventory) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create inventory at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.basicInventory.DecentralizedInventorySave(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateCentrailizedInventory(metaData core.MetaData, data *mdm.CentralizedBasicInventory) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update inventory at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.basicInventory.UpdateCentralizedInventory(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateDecentralizedInventory(metaData core.MetaData, data *mdm.DecentralizedBasicInventory) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update inventory at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.basicInventory.UpdateDecentralizedInventory(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteCentrailizedInventory(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete inventory at data level")
	}
	err := s.basicInventory.DeleteCentralizedInventory(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteDecentralizedInventory(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete inventory at data level")
	}

	err := s.basicInventory.DeleteDecentralizedInventory(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetCentrailizedInventory(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view inventory at data level")
	}

	result, err := s.basicInventory.FindOneCentralizedInventory(metaData.Query)
	if err != nil {
		return nil, err
	}

	//----------Fetch the product information------------------------

	query := map[string]interface{}{"id": result.ProductVariantId}
	response, _ := s.productService.GetVariantView(query, fmt.Sprintf("%v", metaData.TokenUserId), accessTemplateId)
	result.ProductDetails, _ = json.Marshal(response)

	//-------------Fetch the Physical location information---------------------
	if result.PhysicalLocationId != 0 {
		metaData.Query = map[string]interface{}{"id": result.PhysicalLocationId}
		response, _ = s.locationService.GetLocation(metaData)
		result.PhysicalLocation, _ = json.Marshal(response)
	}

	return result, nil
}
func (s *service) GetDecentralizedInventory(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)

	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view inventory at data level")
	}
	result, err := s.basicInventory.FindOneDecentralizedInventory(metaData.Query)
	if err != nil {
		return result, err
	}

	//----------Fetch the product information------------------------
	query := map[string]interface{}{"id": result.ProductVariantId}
	response, _ := s.productService.GetVariantView(query, fmt.Sprintf("%v", metaData.TokenUserId), accessTemplateId)
	result.ProductDetails, _ = json.Marshal(response)

	//-------------Fetch the location information---------------------
	if result.PhysicalLocationId != 0 {
		metaData.Query = map[string]interface{}{"id": result.PhysicalLocationId}
		response, _ = s.locationService.GetLocation(metaData)
		result.PhysicalLocation, _ = json.Marshal(response)
	}

	return result, nil
}
func (s *service) GetCentrailizedInventoryList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view inventory at data level")
	}
	results, err := s.basicInventory.FindAllCentralizedInventory(metaData.Query, p)
	if err != nil {
		return results, err
	}

	for index, result := range results {
		//----------Fetch the product information------------------------
		query := map[string]interface{}{"id": result.ProductVariantId}
		response, _ := s.productService.GetVariantView(query, fmt.Sprintf("%v", metaData.TokenUserId), accessTemplateId)
		result.ProductDetails, _ = json.Marshal(response)

		//-------------Fetch the location information---------------------
		if result.PhysicalLocationId != 0 {
			metaData.Query = map[string]interface{}{"id": result.PhysicalLocationId}
			response, _ = s.locationService.GetLocation(metaData)
			result.PhysicalLocation, _ = json.Marshal(response)
		}

		results[index] = result
	}

	return results, nil
}
func (s *service) GetDecentralizedInventoryList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list inventory at data level")
	}
	results, err := s.basicInventory.FindAllDecentralizedInventory(metaData.Query, p)
	if err != nil {
		return results, err
	}

	for index, result := range results {
		//----------Fetch the product information------------------------
		query := map[string]interface{}{"id": result.ProductVariantId}
		response, _ := s.productService.GetVariantView(query, fmt.Sprintf("%v", metaData.TokenUserId), accessTemplateId)
		result.ProductDetails, _ = json.Marshal(response)

		//-------------Fetch the location information---------------------
		if result.PhysicalLocationId != 0 {
			metaData.Query = map[string]interface{}{"id": result.PhysicalLocationId}
			response, _ = s.locationService.GetLocation(metaData)
			result.PhysicalLocation, _ = json.Marshal(response)
		}

		results[index] = result
	}

	return results, nil
}

// ---------------------channels-----------------------------------------------------------------------------------
func (s *service) UpsertInventoryTemplate(metaData core.MetaData, data []interface{}) (interface{}, error) {
	var success []interface{}
	var failures []interface{}
	for index, payload := range data {
		var DeCentralizedBasicInventory mdm.DecentralizedBasicInventory
		jsonPayloadByteArray, err := json.Marshal(payload)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		err = json.Unmarshal(jsonPayloadByteArray, &DeCentralizedBasicInventory)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		DeCentralizedBasicInventory.UpdatedByID = &metaData.TokenUserId

		InventoryQuery := map[string]interface{}{"product_variant_id": DeCentralizedBasicInventory.ProductVariantId, "channel_code": DeCentralizedBasicInventory.ChannelCode}
		res, _ := s.basicInventory.FindOneDecentralizedInventory(InventoryQuery)
		if res.ProductVariantId != 0 {
			err = s.basicInventory.UpdateDecentralizedInventory(InventoryQuery, &DeCentralizedBasicInventory)
			if err != nil {
				failures = append(failures, map[string]interface{}{"id": DeCentralizedBasicInventory.ProductVariantId, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			success = append(success, map[string]interface{}{"id": DeCentralizedBasicInventory.ProductVariantId, "status": true, "serial_number": index + 1, "msg": "updated"})
			continue
		}
		DeCentralizedBasicInventory.CreatedByID = &metaData.TokenUserId
		DeCentralizedBasicInventory.CompanyId = metaData.CompanyId
		err = s.basicInventory.DecentralizedInventorySave(&DeCentralizedBasicInventory)
		if err != nil {
			failures = append(failures, map[string]interface{}{"id": DeCentralizedBasicInventory.ProductVariantId, "status": false, "serial_number": index + 1, "msg": err})
			continue
		}

		success = append(success, map[string]interface{}{"id": DeCentralizedBasicInventory.ProductVariantId, "status": true, "serial_number": index + 1, "msg": "created"})
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
}

func (s *service) InventoryTransactionCreate(metaData core.MetaData, data *mdm.CentralizedInventoryTransactions) error {
	err := s.basicInventory.CentralizedInventoryTransactionSave(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateTransactionInventory(metaData core.MetaData, data map[string]interface{}) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "INVENTORY_DECENTRALIZED", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update inventory at data level")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	inv_model := new(mdm.CentralizedBasicInventory)

	var list []map[string]interface{}
	var line_item = make(map[string]interface{}, 0)

	if data["order_type"].(string) == "sales_return" {
		sales_return_arr := data["order_lines"].([]returns.SalesReturnLines)
		for _, sales_return := range sales_return_arr {
			line_item["inventory_id"] = sales_return.InventoryId
			line_item["quantity"] = sales_return.QuantityReturned
			list = append(list, line_item)
		}
	}
	if data["order_type"].(string) == "sales_order" {
		sales_order_arr := data["order_lines"].([]orders.SalesOrderLines)
		for _, sales_order := range sales_order_arr {
			line_item["inventory_id"] = sales_order.InventoryId
			line_item["quantity"] = sales_order.Quantity
			list = append(list, line_item)
		}
	}
	if data["order_type"].(string) == "purchase_return" {
		purchase_return_arr := data["order_lines"].([]returns.PurchaseReturnLines)
		for _, purchase_return := range purchase_return_arr {
			line_item["inventory_id"] = purchase_return.InventoryId
			line_item["quantity"] = purchase_return.QuantityReturned
			list = append(list, line_item)
		}
	}
	if data["order_type"].(string) == "purchase_order" {
		purchase_order_arr := data["order_lines"].([]orders.PurchaseOrderLines)
		for _, purchase_order := range purchase_order_arr {
			line_item["inventory_id"] = purchase_order.InventoryId
			line_item["quantity"] = purchase_order.Quantity
			list = append(list, line_item)
		}
	}
	if data["order_type"].(string) == "scrap_order" {
		scrap_order_arr := data["order_lines"].([]orders.ScrapOrderLines)
		for _, scrap_order := range scrap_order_arr {
			line_item["inventory_id"] = scrap_order.InventoryId
			line_item["quantity"] = scrap_order.Scrap_item_quantity
			list = append(list, line_item)
		}
	}

	for _, value := range list {

		metaData.Query = map[string]interface{}{
			"id": value["inventory_id"],
		}

		inv_data, _ := s.GetCentrailizedInventory(metaData)
		byteStream, _ := json.Marshal(inv_data)
		_ = json.Unmarshal(byteStream, inv_model)

		centralizedInventoryTransactionPayload := map[string]interface{}{
			"order_id":             data["id"],
			"order_type":           data["order_type"],
			"description":          fmt.Sprintf("%v with ID: %v", data["order_type"], data["id"]),
			"product_variant_id":   inv_model.ProductVariantId,
			"physical_location_id": inv_model.PhysicalLocationId,
			"inventory_id":         inv_model.ID,
			"opening_inventory":    inv_model.AvailableStock,
		}
		// fmt.Println("centralizedInventoryTransactionPayload", centralizedInventoryTransactionPayload)
		if data["is_update_inv"] == true {
			if data["is_credit"] == true {
				inv_model.AvailableStock = inv_model.AvailableStock + value["quantity"].(int64)
				inv_model.StockExpected = inv_model.StockExpected - value["quantity"].(int64)
				centralizedInventoryTransactionPayload["deposited_stock"] = value["quantity"].(int64)

			} else {
				inv_model.AvailableStock = inv_model.AvailableStock - value["quantity"].(int64)
				inv_model.CommitedStock = inv_model.CommitedStock + value["quantity"].(int64)
				centralizedInventoryTransactionPayload["withdrawn_stock"] = value["quantity"].(int64)
			}
			centralizedInventoryTransactionPayload["closing_inventory"] = inv_model.AvailableStock

			invTransactionPayload := new(mdm.CentralizedInventoryTransactions)

			byteStream, _ = json.Marshal(centralizedInventoryTransactionPayload)
			_ = json.Unmarshal(byteStream, invTransactionPayload)

			err := s.InventoryTransactionCreate(metaData, invTransactionPayload)
			if err != nil {
				fmt.Println("error", err.Error())
				return res.BuildError(res.ErrUnprocessableEntity, err)
			}

		} else {
			if data["is_credit"] == true {
				inv_model.StockExpected = inv_model.StockExpected + value["quantity"].(int64)
			} else {
				inv_model.CommitedStock = inv_model.CommitedStock - value["quantity"].(int64)
			}
		}

		err := s.UpdateCentrailizedInventory(metaData, inv_model)
		if err != nil {
			fmt.Println("error", err.Error())
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	}
	return nil
}
