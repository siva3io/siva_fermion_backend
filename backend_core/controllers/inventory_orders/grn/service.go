package grn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/orders/scrap_orders"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_orders"
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
	CreateGRN(metaData core.MetaData, data *inventory_orders.GRN) error
	UpdateGRN(metaData core.MetaData, data *inventory_orders.GRN) error
	GetGRN(metaData core.MetaData) (interface{}, error)
	GetAllGRN(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	DeleteGRN(metaData core.MetaData) error
	DeleteGRNOrderLines(metaData core.MetaData) error
	BulkCreateGRN(data *[]inventory_orders.GRN, token_id string, access_template_id string) error

	GetGrnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	grnrepository     inventory_orders_repo.GRN
	scrapOrderService scrap_orders.Service
	orderRepository   *order_repo.PurchaseOrders
	asnRepository     inventory_orders_repo.Asn
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	grnrepository := inventory_orders_repo.NewGRN()
	asnRepository := inventory_orders_repo.NewAsn()
	orderRepository := order_repo.NewPurchaseOrder()
	newServiceObj = &service{grnrepository,
		scrap_orders.NewService(), orderRepository, asnRepository}
	return newServiceObj
}

func (s *service) CreateGRN(metaData core.MetaData, data *inventory_orders.GRN) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "GRN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create grn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create grn at data level")
	}
	var grn_data GRNRequest
	if grn_data.AutoGenerateGrnNumber {
		grn_data.GRNNumber = helpers.GenerateSequence("GRN", fmt.Sprint(metaData.TokenUserId), "grns")
	}
	if grn_data.AutoGenerateReferenceNumber {
		grn_data.ReferenceNumber = helpers.GenerateSequence("REF", fmt.Sprint(metaData.TokenUserId), "grns")
	}
	defaultStatus, err := helpers.GetLookupcodeId("ASN_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	data.StatusId = defaultStatus
	err = s.grnrepository.CreateGRN(data)
	if err != nil {
		return err
	}
	return nil
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

func (s *service) GetGRN(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "GRN", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view grn at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view grn at data level")
	}
	result, err := s.grnrepository.GetGRN(metaData.Query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetAllGRN(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "GRN", metaData.TokenUserId)

	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list grn at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list grn at data level")
	}
	result, err := s.grnrepository.GetAllGRN(metaData.Query, page)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateGRN(metaData core.MetaData, data *inventory_orders.GRN) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "GRN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update grn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update grn at data level")
	}

	var find_query = map[string]interface{}{
		"id": metaData.Query["id"],
	}
	old_data, er := s.grnrepository.GetGRN(find_query)
	if er != nil {
		return er
	}
	data.UpdatedByID = &metaData.TokenUserId
	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result
	}
	service_err := s.grnrepository.UpdateGRN(metaData.Query, data)
	if service_err != nil {
		return service_err
	}
	for _, ord_lines := range data.GRNOrderLines {
		query := map[string]interface{}{
			"grn_id":     old_data.ID,
			"product_id": ord_lines.ProductID,
		}
		ord_lines.UpdatedByID = &metaData.TokenUserId
		count, err1 := s.grnrepository.UpdateOrderLines(query, ord_lines)
		if err1 != nil {
			return err1
		} else if count == 0 {
			ord_lines.GRN_ID = old_data.ID
			ord_lines.CreatedByID = &metaData.TokenUserId
			ord_lines.CompanyId = metaData.CompanyId
			ord_lines.UpdatedByID = nil
			err2 := s.grnrepository.CreateOrderLines(ord_lines)
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func (s *service) DeleteGRN(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "GRN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete grn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete grn at data level")
	}
	err := s.grnrepository.DeleteGRN(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{"grn_id": metaData.Query["id"]}
		er := s.grnrepository.DeleteOrderLines(query)
		if er != nil {
			return er
		}
	}

	return nil
}

func (s *service) DeleteGRNOrderLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "GRN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete grn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete grn at data level")
	}
	err := s.grnrepository.DeleteOrderLines(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetGrnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "GRN", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list grn at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list grn at data level")
	}
	tab := metaData.AdditionalFields["tab"].(string)
	scrapOrderId := metaData.AdditionalFields["id"]

	if tab == "scrap_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SCRAP_ORDERS_SOURCE_DOCUMENT_TYPES", "GRN")

		scrapOrderPage := page
		scrapOrderPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, scrapOrderId)
		data, err := s.scrapOrderService.AllScrapOrders(metaData, scrapOrderPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "asn" {
		grnPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("GRN_SOURCE_DOCUMENT_TYPES", "ASN")
		grnPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", scrapOrderId, source_document_type_id)
		grn_order_interface, err := s.GetAllGRN(metaData, grnPagination)
		if err != nil {
			return nil, err
		}
		grn_order := grn_order_interface.([]inventory_orders.GRN)
		if len(grn_order) == 0 {
			return nil, errors.New("no source document")
		}

		var grnOrderSourceDoc map[string]interface{}
		grnOrderSourceDocJson := grn_order[0].SourceDocuments
		dto, err := json.Marshal(grnOrderSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &grnOrderSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		grnOrderSourceDocId := grnOrderSourceDoc["id"]
		if grnOrderSourceDocId == nil {
			return nil, errors.New("no source document")
		}

		asnPage := page
		asnPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", grnOrderSourceDoc["id"])
		data, err := s.asnRepository.GetAllAsn(metaData.Query, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	return nil, nil
}
