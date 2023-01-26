package grn

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/orders/scrap_orders"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	inventory_orders_repo "fermion/backend_core/internal/repository/inventory_orders"
	order_repo "fermion/backend_core/internal/repository/orders"
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
	UpdateGRN(id uint, data *GRNRequest, token_id string, access_template_id string) error
	DeleteGRN(id uint, user_id uint, token_id string, access_template_id string) error
	CreateGRN(data *GRNRequest, token_id string, access_template_id string) (uint, error)
	GetAllGRN(page *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_orders.GRN, error)
	GetGRN(id uint, token_id string, access_template_id string) (interface{}, error)
	SearchGRN(key string, token_id string, access_template_id string) (interface{}, error)
	DeleteGRNOrderLines(query interface{}, token_id string, access_template_id string) error
	BulkCreateGRN(data *[]inventory_orders.GRN, token_id string, access_template_id string) error

	GetGrnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	grnrepository     inventory_orders_repo.GRN
	scrapOrderService scrap_orders.Service
	orderRepository   *order_repo.PurchaseOrders
	asnRepository     inventory_orders_repo.Asn
}

func NewService() *service {
	grnrepository := inventory_orders_repo.NewGRN()
	asnRepository := inventory_orders_repo.NewAsn()
	orderRepository := order_repo.NewPurchaseOrder()
	return &service{grnrepository,
		scrap_orders.NewService(), orderRepository, asnRepository}
}

func (s *service) CreateGRN(data *GRNRequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "GRN", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create grn at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create grn at data level")
	}
	var grn_data inventory_orders.GRN
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &grn_data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	if data.AutoGenerateGrnNumber {
		grn_data.GRNNumber = helpers.GenerateSequence("GRN", token_id, "grns")
	}
	if data.AutoGenerateReferenceNumber {
		grn_data.ReferenceNumber = helpers.GenerateSequence("REF", token_id, "grns")
	}
	defaultStatus, err := helpers.GetLookupcodeId("ASN_STATUS", "DRAFT")
	if err != nil {
		return 0, err
	}
	data.StatusId = defaultStatus
	id, service_err := s.grnrepository.CreateGRN(&grn_data, token_id)
	fmt.Println(grn_data)
	return id, service_err
}

func (s *service) BulkCreateGRN(data *[]inventory_orders.GRN, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "GRN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create grn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create grn at data level")
	}
	service_err := s.grnrepository.BulkCreateGRN(data, token_id)
	return service_err
}

func (s *service) GetGRN(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "GRN", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view grn at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view grn at data level")
	}
	result, err := s.grnrepository.GetGRN(id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetAllGRN(page *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_orders.GRN, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "GRN", *token_user_id)

	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list grn at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list grn at data level")
	}
	result, err := s.grnrepository.GetAllGRN(page)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateGRN(id uint, data *GRNRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "GRN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update grn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update grn at data level")
	}
	var grn_data inventory_orders.GRN
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &grn_data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	old_data, er := s.grnrepository.GetGRN(id)
	if er != nil {
		return er
	}
	old_status := old_data.StatusId
	new_status := grn_data.StatusId
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		grn_data.StatusHistory = result
	}
	service_err := s.grnrepository.UpdateGRN(id, &grn_data)
	if service_err != nil {
		return service_err
	}
	for _, ord_lines := range grn_data.GRNOrderLines {
		query := map[string]interface{}{
			"grn_id":     id,
			"product_id": ord_lines.ProductID,
		}
		count, err1 := s.grnrepository.UpdateOrderLines(query, ord_lines)
		if err1 != nil {
			return err1
		} else if count == 0 {
			ord_lines.GRN_ID = id
			err2 := s.grnrepository.CreateOrderLines(ord_lines)
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func (s *service) SearchGRN(key string, token_id string, access_template_id string) (interface{}, error) {
	result, err := s.grnrepository.SearchGRN(key)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *service) DeleteGRN(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "GRN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete grn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete grn at data level")
	}
	err := s.grnrepository.DeleteGRN(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"grn_id": id}
	err1 := s.grnrepository.DeleteOrderLines(query)
	if err1 != nil {
		return err
	}
	return nil
}

func (s *service) DeleteGRNOrderLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "GRN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete grn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete grn at data level")
	}
	err := s.grnrepository.DeleteOrderLines(query)
	return err
}
func (s *service) GetGrnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "GRN", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list grn at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list grn at data level")
	}

	scrapOrderId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "scrap_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SCRAP_ORDERS_SOURCE_DOCUMENT_TYPES", "GRN")

		scrapOrderPage := page
		scrapOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"link_source_document\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, scrapOrderId)
		data, err := s.scrapOrderService.AllScrapOrders(scrapOrderPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "asn" {
		var data []inventory_orders.ASN
		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES", "GRN")
		purchaseOrderPage := page
		purchaseOrderPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, scrapOrderId)
		PurchaseOrdersInterface, err := s.orderRepository.FindAll(purchaseOrderPage)
		PurchaseOrdersData := PurchaseOrdersInterface.([]orders.PurchaseOrders)
		if err != nil {
			return nil, err
		}
		for _, POValue := range PurchaseOrdersData {

			sourceDocumentId, _ := helpers.GetLookupcodeId("ASN_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

			AsnPage := new(pagination.Paginatevalue)
			AsnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, POValue.ID)
			data, err = s.asnRepository.GetAllAsn(AsnPage)
			if err != nil {
				return nil, err
			}
		}
		return data, nil
	}

	return nil, nil
}
