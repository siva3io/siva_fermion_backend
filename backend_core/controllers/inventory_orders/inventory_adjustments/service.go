package inventory_adjustments

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/pagination"
	inventory_orders_repo "fermion/backend_core/internal/repository/inventory_orders"
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
	CreateInvAdj(data *InventoryAdjustmentsRequest, token_id string, access_template_id string) (uint, error)
	BulkCreateInvAdj(data *[]InventoryAdjustmentsRequest, token_id string, access_template_id string) error
	UpdateInvAdj(id uint, data *InventoryAdjustmentsRequest, token_id string, access_template_id string) error
	GetInvAdj(id uint, token_id string, access_template_id string) (interface{}, error)
	GetAllInvAdj(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_orders.InventoryAdjustments, error)
	DeleteInvAdj(id uint, user_id uint, token_id string, access_template_id string) error
	DeleteInvAdjLines(query interface{}, token_id string, access_template_id string) error

	SendMailInvAdj(q *SendMailInvAdj) error
}

type service struct {
	invadjRepository inventory_orders_repo.InventoryAdjustments
}

func NewService() *service {
	InvAdjRepository := inventory_orders_repo.NewInvAdj()
	return &service{InvAdjRepository}
}

func (s *service) CreateInvAdj(data *InventoryAdjustmentsRequest, token_id string, access_template_id string) (uint, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "INVENTORY_ADJUSTMENT", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create inventory_adjustment at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create inventory_adjustment at data level")
	}
	var InvAdjData inventory_orders.InventoryAdjustments
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &InvAdjData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	id, err := s.invadjRepository.CreateInvAdj(&InvAdjData, token_id)
	return id, err
}

func (s *service) BulkCreateInvAdj(data *[]InventoryAdjustmentsRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "INVENTORY_ADJUSTMENT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create inventory_adjustment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create inventory_adjustment at data level")
	}
	var InvAdjData []inventory_orders.InventoryAdjustments
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &InvAdjData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	err = s.invadjRepository.BulkCreateInvAdj(&InvAdjData, token_id)
	return err
}

func (s *service) UpdateInvAdj(id uint, data *InventoryAdjustmentsRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "INVENTORY_ADJUSTMENT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update inventory_adjustment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update inventory_adjustment at data level")
	}
	var InvAdjData inventory_orders.InventoryAdjustments
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &InvAdjData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	old_data, er := s.invadjRepository.GetInvAdj(id)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := InvAdjData.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, InvAdjData.StatusID)
		InvAdjData.StatusHistory = result
	}
	err = s.invadjRepository.UpdateInvAdj(id, &InvAdjData)
	if err != nil {
		return err
	}
	for _, order_line := range InvAdjData.InventoryAdjustmentLines {
		query := map[string]interface{}{
			"inv_adj_id": uint(id),
			"product_id": uint(order_line.ProductID),
		}
		count, er := s.invadjRepository.UpdateInvAdjLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.Inv_Adj_Id = id
			e := s.invadjRepository.CreateInvAdjLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) GetInvAdj(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "INVENTORY_ADJUSTMENT", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view inventory_adjustment at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view inventory_adjustment at data level")
	}
	result, er := s.invadjRepository.GetInvAdj(id)
	if er != nil {
		return result, er
	}
	query := map[string]interface{}{
		"inv_adj_id": id,
	}
	result_order_lines, err := s.invadjRepository.GetInvAdjLines(query)
	result.InventoryAdjustmentLines = result_order_lines
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) GetAllInvAdj(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_orders.InventoryAdjustments, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "INVENTORY_ADJUSTMENT", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list inventory_adjustment at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list inventory_adjustment at data level")
	}
	result, err := s.invadjRepository.GetAllInvAdj(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteInvAdj(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "INVENTORY_ADJUSTMENT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at data level")
	}
	_, er := s.invadjRepository.GetInvAdj(id)
	if er != nil {
		return er
	}
	err := s.invadjRepository.DeleteInvAdj(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"inv_adj_id": id}
	err1 := s.invadjRepository.DeleteInvAdjLines(query)
	return err1
}

func (s *service) DeleteInvAdjLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "INVENTORY_ADJUSTMENT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete inventory_adjustment at data level")
	}
	data, er := s.invadjRepository.GetInvAdjLines(query)
	if er != nil {
		return er
	}
	if len(data) <= 0 {
		return er
	}
	err := s.invadjRepository.DeleteInvAdjLines(query)
	return err
}

func (s *service) SendMailInvAdj(q *SendMailInvAdj) error {
	id, _ := strconv.Atoi(q.ID)
	result, er := s.invadjRepository.GetInvAdj(uint(id))
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Inventory Adjustment", q.ReceiverEmail, "pkg/util/response/static/inventory_orders/inv_adj_template.html", result)

	return err
}
