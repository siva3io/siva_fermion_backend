package template

import (
	// "context"

	"encoding/json"
	"fmt"

	access_module "fermion/backend_core/controllers/access/module"
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
	FindAllTemplate(p *pagination.Paginatevalue) ([]access_module.AccessModuleAction, error)
	RegisterTemplates(data *model_core.AccessTemplate, module_data []model_core.AccessModuleAction) error
	UpdateTemplateDetails(id uint, data model_core.AccessTemplate, module_data []model_core.AccessModuleAction) error
	GetTemplateDetail(id uint, p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error)
	GetAllTemplateList(p *pagination.Paginatevalue) ([]AccessTemplate, error)
	DeleteTemplateDetail(id uint, user_id uint) error
}
type service struct {
	templateRepository access_repo.Template
	moduleRepository   access_repo.Module
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{access_repo.NewTemplate(), access_repo.NewModule()}
	return newServiceObj

}

func (s *service) RegisterTemplates(data *model_core.AccessTemplate, module_data []model_core.AccessModuleAction) error {

	access_template_id, err := s.templateRepository.CreateTemplate(data)

	if err != nil {
		return err
	}
	for _, module := range module_data {

		module.AccessTemplateId = access_template_id

		err = s.moduleRepository.CreateModule(&module)
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *service) UpdateTemplateDetails(id uint, data model_core.AccessTemplate, module_data []model_core.AccessModuleAction) error {
	_, er := s.templateRepository.FindByidTemplate(id)

	if er != nil {
		return er
	}

	count, err := s.templateRepository.UpdateTemplate(id, data)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("no data gets updated")
	}

	for _, module := range module_data {

		if module.ID == 0 {
			module.AccessTemplateId = id
			err = s.moduleRepository.CreateModule(&module)
			if err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}

		_, err = s.moduleRepository.UpdateModule(module.ID, module)
	}
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *service) GetAllTemplateList(p *pagination.Paginatevalue) ([]AccessTemplate, error) {
	var resultdto []AccessTemplate
	result, _ := s.templateRepository.FindAllTemplate(p)
	value, _ := json.Marshal(result)
	err := json.Unmarshal(value, &resultdto)
	if err != nil {
		return nil, err
	}
	return resultdto, nil
}

func (s *service) FindAllTemplate(p *pagination.Paginatevalue) ([]access_module.AccessModuleAction, error) {
	var resultdto []access_module.AccessModuleAction
	result, _ := s.moduleRepository.FindAllModule(p)
	value, _ := json.Marshal(result)
	err := json.Unmarshal(value, &resultdto)
	if err != nil {
		return nil, err
	}

	return resultdto, nil
}

func (s *service) GetTemplateDetail(id uint, p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error) {

	cout, err := s.moduleRepository.FindByidModule(id, p)
	if err != nil {
		return cout, err
	}

	return cout, nil
}

func (s *service) DeleteTemplateDetail(id uint, user_id uint) error {

	_, er := s.templateRepository.FindByidTemplate(id)
	if er != nil {
		return er
	}

	err := s.templateRepository.DeleteTemplate(id, user_id)
	if err != nil {
		return err
	}
	return nil
}
