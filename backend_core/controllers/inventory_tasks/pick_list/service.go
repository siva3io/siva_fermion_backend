package pick_list

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
	CreatePickList(data *PicklistRequest, token_id string, access_template_id string) (uint, error)
	BulkCreatePickList(data *[]PicklistRequest, token_id string, access_template_id string) error
	UpdatePickList(id uint, data *PicklistRequest, token_id string, access_template_id string) error
	GetPickList(id uint, token_id string, access_template_id string) (interface{}, error)
	GetAllPickList(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_tasks.PickList, error)
	DeletePickList(id uint, user_id uint, token_id string, access_template_id string) error
	DeletePickListLines(query interface{}, token_id string, access_template_id string) error

	SendMailPickList(q *SendMailPickList) error
}

type service struct {
	picklistRepository inventory_tasks_repo.Picklist
}

func NewService() *service {
	PicklistRepository := inventory_tasks_repo.NewPicklist()
	return &service{PicklistRepository}
}

func (s *service) CreatePickList(data *PicklistRequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PICK_LIST", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create pick list at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create pick list at data level")
	}
	var PickListData inventory_tasks.PickList
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &PickListData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	if data.AutoCreatePicklistNumber {
		PickListData.PickListNumber = helpers.GenerateSequence("PICKLIST", token_id, "pick_lists")
	}
	defaultStatus, err := helpers.GetLookupcodeId("PICK_LIST_STATUS", "IN_PROGRESS")
	if err != nil {
		return 0, err
	}
	data.StatusID = defaultStatus

	id, err := s.picklistRepository.CreatePickList(&PickListData, token_id)
	return id, err
}

func (s *service) BulkCreatePickList(data *[]PicklistRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PICK_LIST", *token_user_id)
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
			value.PickListNumber = "PICKLIST/" + token_id + "/000" + strconv.Itoa(count)
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

	err = s.picklistRepository.BulkCreatePickList(&PickListData, token_id)
	return err
}

func (s *service) UpdatePickList(id uint, data *PicklistRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PICK_LIST", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update pick list at data level")
	}
	var PickListData inventory_tasks.PickList
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &PickListData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	old_data, er := s.picklistRepository.GetPickList(id)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := PickListData.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, PickListData.StatusID)
		PickListData.StatusHistory = result
	}
	err = s.picklistRepository.UpdatePickList(id, &PickListData)
	for _, order_line := range PickListData.PicklistLines {
		query := map[string]interface{}{
			"pick_list_id": uint(id),
			"product_id":   uint(order_line.ProductID),
		}
		count, er := s.picklistRepository.UpdatePickListLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.PickList_Id = id
			e := s.picklistRepository.CreatePickListLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *service) GetPickList(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PICK_LIST", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view pick list at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view pick list at data level")
	}
	result, er := s.picklistRepository.GetPickList(id)
	if er != nil {
		return result, er
	}
	query := map[string]interface{}{
		"pick_list_id": id,
	}
	result_order_lines, err := s.picklistRepository.GetPickListLines(query)
	result.PicklistLines = result_order_lines
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) GetAllPickList(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_tasks.PickList, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PICK_LIST", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pick list at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pick list at data level")
	}
	result, err := s.picklistRepository.GetAllPickList(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeletePickList(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PICK_LIST", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pick list at data level")
	}
	_, er := s.picklistRepository.GetPickList(id)
	if er != nil {
		return er
	}
	err := s.picklistRepository.DeletePickList(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"pick_list_id": id}
	err1 := s.picklistRepository.DeletePickListLines(query)
	return err1
}

func (s *service) DeletePickListLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PICK_LIST", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pick list at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pick list at data level")
	}
	data, er := s.picklistRepository.GetPickListLines(query)
	if er != nil {
		return er
	}
	if len(data) <= 0 {
		return er
	}
	err := s.picklistRepository.DeletePickListLines(query)
	return err
}

func (s *service) SendMailPickList(q *SendMailPickList) error {
	id, _ := strconv.Atoi(q.ID)
	result, er := s.picklistRepository.GetPickList(uint(id))
	if er != nil {
		return er
	}
	err := helpers.SendEmail("Pick List", q.ReceiverEmail, "pkg/util/response/static/inventory_tasks/picklist_template.html", result)

	return err
}
