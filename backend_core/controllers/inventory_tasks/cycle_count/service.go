package cycle_count

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/core"
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
	CreateCycleCount(metaData core.MetaData, data *inventory_tasks.CycleCount) error
	BulkCreateCycleCount(metaData core.MetaData, data *[]CycleCountRequest) error
	UpdateCycleCount(metaData core.MetaData, data *inventory_tasks.CycleCount) error
	GetCycleCount(metaData core.MetaData) (interface{}, error)
	GetAllCycleCount(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeleteCycleCount(metaData core.MetaData) error
	DeleteCycleCountLines(metaData core.MetaData) error

	SendMailCycleCount(q *SendMailCycleCount) error
}

type service struct {
	cycleCountRepository inventory_tasks_repo.CycleCount
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	CycleCountRepository := inventory_tasks_repo.NewCycleCount()
	newServiceObj = &service{CycleCountRepository}
	return newServiceObj
}

func (s *service) CreateCycleCount(metaData core.MetaData, data *inventory_tasks.CycleCount) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "CYCLE_COUNT", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create cycle count at data level")
	}
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	defaultStatus, err := helpers.GetLookupcodeId("CYCLE_COUNT_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusID = defaultStatus
	er := s.cycleCountRepository.CreateCycleCount(data)
	if er != nil {
		return er
	}
	return nil
}

func (s *service) BulkCreateCycleCount(metaData core.MetaData, data *[]CycleCountRequest) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "CYCLE_COUNT", metaData.TokenUserId)
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

func (s *service) UpdateCycleCount(metaData core.MetaData, data *inventory_tasks.CycleCount) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "CYCLE_COUNT", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update cycle count at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	old_data, er := s.cycleCountRepository.GetCycleCount(metaData.Query)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := data.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusID)
		data.StatusHistory = result
	}
	err := s.cycleCountRepository.UpdateCycleCount(metaData.Query, data)
	for _, order_line := range data.OrderLines {
		query := map[string]interface{}{
			"product_id":     order_line.ProductID,
			"cycle_count_id": metaData.Query["id"],
		}
		_, er := s.cycleCountRepository.UpdateCycleCountLines(query, &order_line)
		if er != nil {
			return er
		}
		// else if count == 0 {
		// 	order_line.Cycle_count_id = uint(metaData.Query["id"].(float64))
		// 	e := s.cycleCountRepository.CreateCycleCountLines(&order_line)
		// 	if e != nil {
		// 		return e
		// 	}
		// }
	}
	return err
}

func (s *service) GetCycleCount(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "CYCLE_COUNT", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view cycle count at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view cycle count at data level")
	}
	result, er := s.cycleCountRepository.GetCycleCount(metaData.Query)
	if er != nil {
		return nil, er
	}
	// query := map[string]interface{}{
	// 	"cycle_count_id": id,
	// }
	// result_order_lines, err := s.cycleCountRepository.GetCycleCountLines(query)
	// result.OrderLines = result_order_lines
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}

func (s *service) GetAllCycleCount(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "CYCLE_COUNT", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list cycle count at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list cycle count at data level")
	}
	result, err := s.cycleCountRepository.GetAllCycleCount(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteCycleCount(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "CYCLE_COUNT", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete cycle count at data level")
	}

	err := s.cycleCountRepository.DeleteCycleCount(metaData.Query)
	if err != nil {
		return err
	}
	err1 := s.cycleCountRepository.DeleteCycleCountLines(metaData.Query)
	if err != nil {
		return err1
	}
	return nil
}

func (s *service) DeleteCycleCountLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "CYCLE_COUNT", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete cycle count at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete cycle count at data level")
	}
	data, err := s.cycleCountRepository.GetCycleCountLines(metaData.Query)
	if err != nil {
		return err
	}
	if len(data) <= 0 {
		return err
	}
	err = s.cycleCountRepository.DeleteCycleCountLines(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SendMailCycleCount(q *SendMailCycleCount) error {
	var metaData core.MetaData
	result, er := s.cycleCountRepository.GetCycleCount(metaData.Query)
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Cycle Count", q.ReceiverEmail, "pkg/util/response/static/inventory_tasks/cycle_count_template.html", result)

	return err
}
