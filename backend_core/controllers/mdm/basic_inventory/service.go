package basic_inventory

import (
	"encoding/json"
	"fmt"
	"sync"

	"fermion/backend_core/controllers/mdm/locations"
	"fermion/backend_core/controllers/mdm/products"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
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
	CreateCentrailizedInventory(data *mdm.CentralizedBasicInventory, token_id string, access_template_id string) error
	CreateDecentralizedInventory(data *mdm.DecentralizedBasicInventory, token_id string, access_template_id string) error

	UpdateCentrailizedInventory(query map[string]interface{}, data *mdm.CentralizedBasicInventory, token_id string, access_template_id string) error
	UpdateDecentralizedInventory(query map[string]interface{}, data *mdm.DecentralizedBasicInventory, token_id string, access_template_id string) error

	DeleteCentrailizedInventory(query map[string]interface{}, token_id string, access_template_id string) error
	DeleteDecentralizedInventory(query map[string]interface{}, token_id string, access_template_id string) error

	GetCentrailizedInventory(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetDecentralizedInventory(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)

	GetCentrailizedInventoryList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
	GetDecentralizedInventoryList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_actions string) (interface{}, error)

	SearchCentrailizedInventory(query string, token_id string, access_template_id string) ([]CentralizedSearchObjDTO, error)
	SearchDecentralizedInventory(query string, token_id string, access_template_id string) ([]DecentralizedSearchObjDTO, error)

	UpsertInventoryTemplate(data []interface{}, TokenUserId string) (interface{}, error)

	InventoryTransactionCreate(data *mdm.CentralizedInventoryTransactions, token_id string, access_template_id string) error

	UpdateInventory(data map[string]interface{}, token_id string, access_template_id string) error
}

type service struct {
	basicInventory  mdm_repo.BasicInventory
	productService  products.Service
	locationService locations.Service
	mu              sync.Mutex
}

func NewService() *service {
	basicInventory := mdm_repo.NewBasicInventory()
	productService := products.NewService()
	locationService := locations.NewService()
	return &service{
		basicInventory:  basicInventory,
		productService:  productService,
		locationService: locationService,
	}
}

func (s *service) CreateCentrailizedInventory(data *mdm.CentralizedBasicInventory, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create inventory at data level")
	}
	err := s.basicInventory.CentralizedInventorySave(data)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) CreateDecentralizedInventory(data *mdm.DecentralizedBasicInventory, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create inventory at data level")
	}
	err := s.basicInventory.DecentralizedInventorySave(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateCentrailizedInventory(query map[string]interface{}, data *mdm.CentralizedBasicInventory, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update inventory at data level")
	}
	_, err := s.basicInventory.FindOneCentralizedInventory(query)
	if err != nil {
		return err
	}
	err = s.basicInventory.UpdateCentralizedInventory(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateDecentralizedInventory(query map[string]interface{}, data *mdm.DecentralizedBasicInventory, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update inventory at data level")
	}
	_, err := s.basicInventory.FindOneDecentralizedInventory(query)
	if err != nil {
		return err
	}
	err = s.basicInventory.UpdateDecentralizedInventory(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteCentrailizedInventory(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete inventory at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.basicInventory.FindOneCentralizedInventory(q)
	if er != nil {
		return er
	}
	err := s.basicInventory.DeleteCentralizedInventory(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteDecentralizedInventory(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete inventory at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete inventory at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.basicInventory.FindOneDecentralizedInventory(q)
	if er != nil {
		return er
	}

	err := s.basicInventory.DeleteDecentralizedInventory(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetCentrailizedInventory(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view inventory at data level")
	}
	result, err := s.basicInventory.FindOneCentralizedInventory(query)
	if err != nil {
		return nil, err
	}

	//----------Fetch the product information------------------------
	query = map[string]interface{}{"id": result.ProductVariantId}
	response, _ := s.productService.GetVariantView(query, token_id, access_template_id)
	result.ProductDetails, _ = json.Marshal(response)

	//-------------Fetch the Physical location information---------------------
	if result.PhysicalLocationId != 0 {
		query = map[string]interface{}{"id": result.PhysicalLocationId}
		response, _ = s.locationService.GetLocation(query, token_id, access_template_id)
		result.PhysicalLocation, _ = json.Marshal(response)
	}

	var centralized_basic_inventory_response CentralizedBasicInventoryResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &centralized_basic_inventory_response)
	if err != nil {
		return nil, err
	}

	return centralized_basic_inventory_response, nil
}
func (s *service) GetDecentralizedInventory(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view inventory at data level")
	}
	result, err := s.basicInventory.FindOneDecentralizedInventory(query)
	if err != nil {
		return result, err
	}

	//----------Fetch the product information------------------------
	query = map[string]interface{}{"id": result.ProductVariantId}
	response, _ := s.productService.GetVariantView(query, token_id, access_template_id)
	result.ProductDetails, _ = json.Marshal(response)

	//-------------Fetch the location information---------------------
	if result.PhysicalLocationId != 0 {
		query = map[string]interface{}{"id": result.PhysicalLocationId}
		response, _ = s.locationService.GetLocation(query, token_id, access_template_id)
		result.PhysicalLocation, _ = json.Marshal(response)
	}

	var de_centralized_basic_inventory_response DecentralizedBasicInventoryResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &de_centralized_basic_inventory_response)
	if err != nil {
		return nil, err
	}

	return de_centralized_basic_inventory_response, nil
}
func (s *service) GetCentrailizedInventoryList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view inventory at data level")
	}
	results, err := s.basicInventory.FindAllCentralizedInventory(query, p)
	if err != nil {
		return results, err
	}

	for index, result := range results {
		//----------Fetch the product information------------------------
		query := map[string]interface{}{"id": result.ProductVariantId}
		response, _ := s.productService.GetVariantView(query, token_id, access_template_id)
		result.ProductDetails, _ = json.Marshal(response)

		//-------------Fetch the location information---------------------
		if result.PhysicalLocationId != 0 {
			query = map[string]interface{}{"id": result.PhysicalLocationId}
			response, _ = s.locationService.GetLocation(query, token_id, access_template_id)
			result.PhysicalLocation, _ = json.Marshal(response)
		}

		results[index] = result
	}

	var centralized_basic_inventory_response []CentralizedBasicInventoryResponseDTO
	marshaldata, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &centralized_basic_inventory_response)
	if err != nil {
		return nil, err
	}

	return centralized_basic_inventory_response, nil
}
func (s *service) GetDecentralizedInventoryList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list inventory at data level")
	}
	results, err := s.basicInventory.FindAllDecentralizedInventory(query, p)
	if err != nil {
		return results, err
	}

	for index, result := range results {
		//----------Fetch the product information------------------------
		query := map[string]interface{}{"id": result.ProductVariantId}
		response, _ := s.productService.GetVariantView(query, token_id, access_template_id)
		result.ProductDetails, _ = json.Marshal(response)

		//-------------Fetch the location information---------------------
		if result.PhysicalLocationId != 0 {
			query = map[string]interface{}{"id": result.PhysicalLocationId}
			response, _ = s.locationService.GetLocation(query, token_id, access_template_id)
			result.PhysicalLocation, _ = json.Marshal(response)
		}

		results[index] = result
	}

	var de_centralized_basic_inventory_response []DecentralizedBasicInventoryResponseDTO
	marshaldata, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &de_centralized_basic_inventory_response)
	if err != nil {
		return nil, err
	}

	return de_centralized_basic_inventory_response, nil
}
func (s *service) SearchCentrailizedInventory(query string, token_id string, access_template_id string) ([]CentralizedSearchObjDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list inventory at data level")
	}
	result, err := s.basicInventory.SearchCentralizedInventory(query)
	var centralizedData []CentralizedSearchObjDTO
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	for _, v := range result {
		var data CentralizedSearchObjDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		centralizedData = append(centralizedData, data)
	}
	return centralizedData, nil
}
func (s *service) SearchDecentralizedInventory(query string, token_id string, access_template_id string) ([]DecentralizedSearchObjDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "INVENTORY_DECENTRALIZED", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list inventory at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list inventory at data level")
	}
	result, err := s.basicInventory.SearchDecentralizedInventory(query)
	var decentralizeddata []DecentralizedSearchObjDTO
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	for _, v := range result {
		var data DecentralizedSearchObjDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		decentralizeddata = append(decentralizeddata, data)
	}
	return decentralizeddata, nil
}

// ---------------------channels-----------------------------------------------------------------------------------
func (s *service) UpsertInventoryTemplate(data []interface{}, TokenUserId string) (interface{}, error) {
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
		DeCentralizedBasicInventory.UpdatedByID = helpers.ConvertStringToUint(TokenUserId)

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
		DeCentralizedBasicInventory.CreatedByID = helpers.ConvertStringToUint(TokenUserId)
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

func (s *service) InventoryTransactionCreate(data *mdm.CentralizedInventoryTransactions, token_id string, access_template_id string) error {
	err := s.basicInventory.CentralizedInventoryTransactionSave(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateInventory(data map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "INVENTORY_DECENTRALIZED", *token_user_id)
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

		query := map[string]interface{}{
			"id": value["inventory_id"],
		}

		inv_data, _ := s.GetCentrailizedInventory(query, token_id, access_template_id)
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

			err := s.InventoryTransactionCreate(invTransactionPayload, token_id, access_template_id)
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

		err := s.UpdateCentrailizedInventory(query, inv_model, token_id, access_template_id)
		if err != nil {
			fmt.Println("error", err.Error())
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	}
	return nil
}
