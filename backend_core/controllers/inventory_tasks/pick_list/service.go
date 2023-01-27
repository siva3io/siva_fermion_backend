package pick_list

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
	CreatePickList(metaData core.MetaData, data *inventory_tasks.PickList) error
	BulkCreatePickList(metaData core.MetaData, data *[]PicklistRequest) error
	UpdatePickList(metaData core.MetaData, data *inventory_tasks.PickList) error
	GetPickList(metaData core.MetaData) (interface{}, error)
	GetAllPickList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeletePickList(metaData core.MetaData) error
	DeletePickListLines(metaData core.MetaData) error

	SendMailPickList(q *SendMailPickList) error
}

type service struct {
	picklistRepository inventory_tasks_repo.Picklist
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	PicklistRepository := inventory_tasks_repo.NewPicklist()
	newServiceObj = &service{PicklistRepository}
	return newServiceObj
}

func (s *service) CreatePickList(metaData core.MetaData, data *inventory_tasks.PickList) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PICK_LIST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pick list at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	data.PickListNumber = helpers.GenerateSequence("PICKLIST", fmt.Sprint(metaData.TokenUserId), "pick_lists")
	defaultStatus, er := helpers.GetLookupcodeId("PICK_LIST_STATUS", "IN_PROGRESS")
	if er != nil {
		return er
	}
	data.StatusID = defaultStatus

	err := s.picklistRepository.CreatePickList(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) BulkCreatePickList(metaData core.MetaData, data *[]PicklistRequest) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PICK_LIST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pick list at data level")
	}
	bulk_data := []map[string]interface{}{}
	for index, value := range *data {
		v := map[string]interface{}{}
		count := helpers.GetCount("SELECT COUNT(*) FROM pick_lists") + 1 + index
		if value.AutoCreatePicklistNumber {
			value.PickListNumber = "PICKLIST/" + fmt.Sprint(metaData.TokenUserId) + "/000" + strconv.Itoa(count)
		}

		val, err := json.Marshal(value)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(val, &v)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}

		bulk_data = append(bulk_data, v)
	}
	var PickListData []inventory_tasks.PickList
	dto, err := json.Marshal(bulk_data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &PickListData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	err = s.picklistRepository.BulkCreatePickList(&PickListData)
	return err
}

func (s *service) UpdatePickList(metaData core.MetaData, data *inventory_tasks.PickList) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "PICK_LIST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update pick list at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId

	old_data, er := s.picklistRepository.GetPickList(metaData.Query)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := data.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusID)
		data.StatusHistory = result
	}

	err := s.picklistRepository.UpdatePickList(metaData.Query, data)
	for _, order_line := range data.PicklistLines {
		query := map[string]interface{}{
			"product_id":   order_line.ProductID,
			"pick_list_id": metaData.Query["id"],
		}
		count, er := s.picklistRepository.UpdatePickListLines(query, &order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.PickList_Id = uint(metaData.Query["id"].(float64))
			e := s.picklistRepository.CreatePickListLines(&order_line)
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *service) GetPickList(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "PICK_LIST", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view pick list at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view pick list at data level")
	}
	result, er := s.picklistRepository.GetPickList(metaData.Query)
	if er != nil {
		return result, er
	}

	result_order_lines, err := s.picklistRepository.GetPickListLines(metaData.Query)
	result.PicklistLines = result_order_lines
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) GetAllPickList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "PICK_LIST", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pick list at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pick list at data level")
	}
	result, err := s.picklistRepository.GetAllPickList(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeletePickList(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "PICK_LIST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pick list at data level")
	}

	err := s.picklistRepository.DeletePickList(metaData.Query)
	if err != nil {
		return err
	}
	err1 := s.picklistRepository.DeletePickListLines(metaData.Query)
	return err1
}

func (s *service) DeletePickListLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "PICK_LIST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pick list at data level")
	}
	data, er := s.picklistRepository.GetPickListLines(metaData.Query)
	if er != nil {
		return er
	}
	if len(data) <= 0 {
		return er
	}
	err := s.picklistRepository.DeletePickListLines(metaData.Query)
	return err
}

func (s *service) SendMailPickList(q *SendMailPickList) error {
	// id, _ := strconv.Atoi(q.ID)
	var metaData core.MetaData
	result, er := s.picklistRepository.GetPickList(metaData.Query)
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Pick List", q.ReceiverEmail, "pkg/util/response/static/inventory_tasks/picklist_template.html", result)

	return err
}
