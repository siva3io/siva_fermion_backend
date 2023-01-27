package accounting

import (
	"fmt"

	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/pagination"
	accounting_repo "fermion/backend_core/internal/repository/accounting"
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
	CreateAccounting(data *accounting.UserAccountingLink) error
	UpdateAccounting(query map[string]interface{}, data *accounting.UserAccountingLink) error
	GetAuthKeys(query map[string]interface{}) (interface{}, error)
	GetUserAccountLinkList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
}

type service struct {
	accountingRepository accounting_repo.Accounting
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{accounting_repo.NewAccountingRepo()}
	return newServiceObj
}

func (s *service) CreateAccounting(data *accounting.UserAccountingLink) error {

	err := s.accountingRepository.CreateAccounting(data)
	return err
}

func (s *service) UpdateAccounting(query map[string]interface{}, data *accounting.UserAccountingLink) error {

	_, er := s.accountingRepository.FindOne(query)
	if er != nil {
		return er
	}
	err := s.accountingRepository.UpdateAccounting(query, data)
	return err
}

func (s *service) GetAuthKeys(query map[string]interface{}) (interface{}, error) {
	result, err := s.accountingRepository.FindOne(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetUserAccountLinkList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "ACCOUNTING", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list contacts at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list contacts at data level")
	}
	result, err := s.accountingRepository.FindAllAccountingink(query, p)
	if err != nil {
		return result, err
	}

	return result, nil
}
