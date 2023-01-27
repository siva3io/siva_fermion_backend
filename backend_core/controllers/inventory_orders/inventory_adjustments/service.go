package inventory_adjustments

import (
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/pagination"
	inventory_orders_repo "fermion/backend_core/internal/repository/inventory_orders"
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
	CreateInvAdj(metaData core.MetaData, data *inventory_orders.InventoryAdjustments) error
	BulkCreateInvAdj(metaData core.MetaData, data *[]inventory_orders.InventoryAdjustments) error
	UpdateInvAdj(metaData core.MetaData, data *inventory_orders.InventoryAdjustments) error
	GetInvAdj(metaData core.MetaData) (interface{}, error)
	GetAllInvAdj(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeleteInvAdj(metaData core.MetaData) error
	DeleteInvAdjLines(metaData core.MetaData) error

	SendMailInvAdj(metaData core.MetaData, q *SendMailInvAdj) error
}

type service struct {
	invadjRepository inventory_orders_repo.InventoryAdjustments
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	InvAdjRepository := inventory_orders_repo.NewInvAdj()
	newServiceObj = &service{InvAdjRepository}
	return newServiceObj
}

func (s *service) CreateInvAdj(metaData core.MetaData, data *inventory_orders.InventoryAdjustments) error {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "INVENTORY_ADJUSTMENT", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create inventory_adjustment at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create inventory_adjustment at data level")
	}
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	err := s.invadjRepository.CreateInvAdj(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) BulkCreateInvAdj(metaData core.MetaData, data *[]inventory_orders.InventoryAdjustments) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "INVENTORY_ADJUSTMENT", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create inventory_adjustment at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create inventory_adjustment at data level")
	}
	err := s.invadjRepository.BulkCreateInvAdj(metaData, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateInvAdj(metaData core.MetaData, data *inventory_orders.InventoryAdjustments) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "INVENTORY_ADJUSTMENT", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update inventory_adjustment at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update inventory_adjustment at data level")
	}
	var find_query = map[string]interface{}{
		"id": metaData.Query["id"],
	}

	old_data, er := s.invadjRepository.GetInvAdj(find_query)
	if er != nil {
		return er
	}
	data.UpdatedByID = &metaData.TokenUserId
	old_status := old_data.StatusID
	new_status := data.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusID)
		data.StatusHistory = result
	}
	err := s.invadjRepository.UpdateInvAdj(metaData.Query, data)
	if err != nil {
		return err
	}
	for _, order_line := range data.InventoryAdjustmentLines {
		query := map[string]interface{}{
			"inv_adj_id": old_data.ID,
			"product_id": order_line.ProductID,
		}
		count, er := s.invadjRepository.UpdateInvAdjLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.Inv_Adj_Id = metaData.TokenUserId
			order_line.CreatedByID = &metaData.TokenUserId
			order_line.UpdatedByID = nil
			order_line.CompanyId = metaData.CompanyId
			e := s.invadjRepository.CreateInvAdjLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) GetInvAdj(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "INVENTORY_ADJUSTMENT", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view inventory_adjustment at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view inventory_adjustment at data level")
	}
	data, err := s.invadjRepository.GetInvAdj(metaData.Query)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *service) GetAllInvAdj(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "INVENTORY_ADJUSTMENT", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list inventory_adjustment at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list inventory_adjustment at data level")
	}
	result, err := s.invadjRepository.GetAllInvAdj(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteInvAdj(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "INVENTORY_ADJUSTMENT", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at data level")
	}
	err := s.invadjRepository.DeleteInvAdj(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{"inv_adj_id": metaData.Query["id"]}
		err1 := s.invadjRepository.DeleteInvAdjLines(query)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func (s *service) DeleteInvAdjLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "INVENTORY_ADJUSTMENT", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at data level")
	}
	err := s.invadjRepository.DeleteInvAdjLines(metaData.Query)
	return err
}

func (s *service) SendMailInvAdj(metaData core.MetaData, q *SendMailInvAdj) error {
	query := map[string]interface{}{
		"id1": q.ID,
	}
	result, er := s.invadjRepository.GetInvAdj(query)
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Inventory Adjustment", q.ReceiverEmail, "pkg/util/response/static/inventory_orders/inv_adj_template.html", result)

	return err
}
