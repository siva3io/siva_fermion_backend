package module

import (
	"encoding/json"
	"fmt"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	access_repo "fermion/backend_core/internal/repository/access"
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

type Services interface {
	RegisterModule(data *model_core.AccessModuleAction) error
	UpdateModuleDetails(id uint, data model_core.AccessModuleAction, p *pagination.Paginatevalue) error
	FindAllModules(p *pagination.Paginatevalue) ([]AccessModuleAction, error)
	FindAllCoreModules(p *pagination.Paginatevalue) ([]CoreAppModule, error)
	FindAllAccessCoreModules(p *pagination.Paginatevalue) ([]CoreAppModule, error)
	GetModuleDetail(id uint, p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error)
	DeleteModuleDetail(id uint, user_id uint, p *pagination.Paginatevalue) error
}

type service struct {
	moduleRepository access_repo.Module
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{access_repo.NewModule()}
	return newServiceObj
}
func (s *service) RegisterModule(data *model_core.AccessModuleAction) error {

	err := s.moduleRepository.CreateModule(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateModuleDetails(id uint, data model_core.AccessModuleAction, p *pagination.Paginatevalue) error {

	_, er := s.moduleRepository.FindByidModule(id, p)
	if er != nil {
		return er
	}

	count, err := s.moduleRepository.UpdateModule(id, data)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("no data gets updated")
	}

	return nil
}

func (s *service) FindAllModules(p *pagination.Paginatevalue) ([]AccessModuleAction, error) {
	var resultdto []AccessModuleAction
	result, _ := s.moduleRepository.FindAllModule(p)
	value, _ := json.Marshal(result)
	err := json.Unmarshal(value, &resultdto)
	if err != nil {
		return nil, err
	}

	return resultdto, nil
}

func (s *service) GetModuleDetail(id uint, p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error) {

	cout, err := s.moduleRepository.FindByidModule(id, p)
	if err != nil {
		return cout, err
	}

	return cout, nil
}

func (s *service) DeleteModuleDetail(id uint, user_id uint, p *pagination.Paginatevalue) error {

	_, er := s.moduleRepository.FindByidModule(id, p)
	if er != nil {
		return er
	}

	err := s.moduleRepository.DeleteModule(id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindAllAccessCoreModules(p *pagination.Paginatevalue) ([]CoreAppModule, error) {
	var resultdto []model_core.AccessModuleAction
	result, _ := s.moduleRepository.FindAllAccessModule(p)
	value, _ := json.Marshal(result)
	err := json.Unmarshal(value, &resultdto)
	if err != nil {
		return nil, err
	}
	filter := "[[\"module_code\",\"in\",["
	for _, temp_value := range result {
		var view_actions []map[string]interface{}
		viewactions_value, _ := json.Marshal(temp_value.ViewActions)
		viewaction_err := json.Unmarshal(viewactions_value, &view_actions)
		if viewaction_err != nil {
			return nil, err
		}
		for _, action := range view_actions {
			if action["ctrl_flag"].(float64) == 1 && action["lookup_code"].(string) == "LIST" {
				filter = filter + "\"" + temp_value.DisplayName + "\","
			}
		}
	}
	sort := "[[\"id\",\"asc\"],[\"item_sequence\",\"asc\"]]"
	filter = filter + "\"USER_PROFILE\",\"ORGANISATION\",\"BUSINESS_SETTINGS\",\"SETTINGS\""
	filter = filter + "]]]"
	page := pagination.Paginatevalue{
		Filters:  filter,
		Sort:     sort,
		Per_page: p.Per_page,
	}
	var core_app_resultdto []CoreAppModule

	core_app_result, _ := s.moduleRepository.FindAllCoreModule(&page)

	core_app_value, _ := json.Marshal(core_app_result)
	final_err := json.Unmarshal(core_app_value, &core_app_resultdto)
	if final_err != nil {
		return nil, final_err
	}

	return core_app_resultdto, nil
}

func (s *service) FindAllCoreModules(p *pagination.Paginatevalue) ([]CoreAppModule, error) {
	var resultdto []CoreAppModule
	result, _ := s.moduleRepository.FindAllCoreModule(p)
	value, _ := json.Marshal(result)
	err := json.Unmarshal(value, &resultdto)
	if err != nil {
		return nil, err
	}

	return resultdto, nil
}
