package cycle_count

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/inventory_tasks"
	"fermion/backend_core/internal/model/pagination"
	inventory_tasks_repo "fermion/backend_core/internal/repository/inventory_tasks"
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
	CreateCycleCount(data *CycleCountRequest, token_id string, access_template_id string) (uint, error)
	BulkCreateCycleCount(data *[]CycleCountRequest, token_id string, access_template_id string) error
	UpdateCycleCount(id uint, data *CycleCountRequest, token_id string, access_template_id string) error
	GetCycleCount(id uint, token_id string, access_template_id string) (interface{}, error)
	GetAllCycleCount(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_tasks.CycleCount, error)
	DeleteCycleCount(id uint, user_id uint, token_id string, access_template_id string) error
	DeleteCycleCountLines(query interface{}, token_id string, access_template_id string) error

	SendMailCycleCount(q *SendMailCycleCount) error
}

type service struct {
	cycleCountRepository inventory_tasks_repo.CycleCount
}

func NewService() *service {
	CycleCountRepository := inventory_tasks_repo.NewCycleCount()
	return &service{CycleCountRepository}
}

func (s *service) CreateCycleCount(data *CycleCountRequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "CYCLE_COUNT", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create cycle count at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create cycle count at data level")
	}
	var CycleCountData inventory_tasks.CycleCount
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &CycleCountData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	defaultStatus, err := helpers.GetLookupcodeId("CYCLE_COUNT_STATUS", "DRAFT")
	if err != nil {
		return 0, err
	}
	data.StatusID = defaultStatus
	id, err := s.cycleCountRepository.CreateCycleCount(&CycleCountData)
	return id, err
}

func (s *service) BulkCreateCycleCount(data *[]CycleCountRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "CYCLE_COUNT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create cycle count at data level")
	}
	var CycleCountData []inventory_tasks.CycleCount
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &CycleCountData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = s.cycleCountRepository.BulkCreateCycleCount(&CycleCountData)
	return err
}

func (s *service) UpdateCycleCount(id uint, data *CycleCountRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "CYCLE_COUNT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update cycle count at data level")
	}
	var CycleCountData inventory_tasks.CycleCount
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &CycleCountData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	old_data, er := s.cycleCountRepository.GetCycleCount(id)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := CycleCountData.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, CycleCountData.StatusID)
		CycleCountData.StatusHistory = result
	}

	err = s.cycleCountRepository.UpdateCycleCount(id, &CycleCountData)
	for _, order_line := range CycleCountData.OrderLines {
		query := map[string]interface{}{
			"cycle_count_id": uint(id),
			"product_id":     uint(order_line.ProductID),
		}
		count, er := s.cycleCountRepository.UpdateCycleCountLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.Cycle_count_id = id
			e := s.cycleCountRepository.CreateCycleCountLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *service) GetCycleCount(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "CYCLE_COUNT", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view cycle count at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view cycle count at data level")
	}
	result, er := s.cycleCountRepository.GetCycleCount(id)
	if er != nil {
		return result, er
	}
	query := map[string]interface{}{
		"cycle_count_id": id,
	}
	result_order_lines, err := s.cycleCountRepository.GetCycleCountLines(query)
	result.OrderLines = result_order_lines
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) GetAllCycleCount(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_tasks.CycleCount, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "CYCLE_COUNT", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list cycle count at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list cycle count at data level")
	}
	result, err := s.cycleCountRepository.GetAllCycleCount(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteCycleCount(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "CYCLE_COUNT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete cycle count at data level")
	}
	_, er := s.cycleCountRepository.GetCycleCount(id)
	if er != nil {
		return er
	}
	err := s.cycleCountRepository.DeleteCycleCount(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"cycle_count_id": id}
	err1 := s.cycleCountRepository.DeleteCycleCountLines(query)
	return err1
}

func (s *service) DeleteCycleCountLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "CYCLE_COUNT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete cycle count at data level")
	}
	data, err := s.cycleCountRepository.GetCycleCountLines(query)
	if err != nil {
		return err
	}
	if len(data) <= 0 {
		return err
	}
	err = s.cycleCountRepository.DeleteCycleCountLines(query)
	return err
}

func (s *service) SendMailCycleCount(q *SendMailCycleCount) error {
	id, _ := strconv.Atoi(q.ID)
	result, er := s.cycleCountRepository.GetCycleCount(uint(id))
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Cycle Count", q.ReceiverEmail, "pkg/util/response/static/inventory_tasks/cycle_count_template.html", result)

	return err
}
