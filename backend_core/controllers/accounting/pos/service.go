package pos

import (
	"fmt"

	accounting_repo "fermion/backend_core/internal/repository/accounting"
	"fermion/backend_core/pkg/util/helpers"

	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/pagination"
	access_checker "fermion/backend_core/pkg/util/access"
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
	CreatePos(data *accounting.UserPosLink) error
	UpdatePos(query map[string]interface{}, data *accounting.UserPosLink) error
	GetAuthKeys(query map[string]interface{}) (interface{}, error)
	GetPosLinkList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
}

type service struct {
	posRepository accounting_repo.Pos
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{accounting_repo.NewPosRepo()}
	return newServiceObj
}

func (s *service) CreatePos(data *accounting.UserPosLink) error {

	err := s.posRepository.CreatePos(data)
	return err
}

func (s *service) UpdatePos(query map[string]interface{}, data *accounting.UserPosLink) error {

	_, er := s.posRepository.FindOne(query)
	if er != nil {
		return er
	}
	err := s.posRepository.UpdatePos(query, data)
	return err
}

func (s *service) GetAuthKeys(query map[string]interface{}) (interface{}, error) {
	result, err := s.posRepository.FindOne(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetPosLinkList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "POS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list contacts at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list contacts at data level")
	}
	result, err := s.posRepository.FindAllPosLink(query, p)
	if err != nil {
		return result, err
	}

	return result, nil
}
