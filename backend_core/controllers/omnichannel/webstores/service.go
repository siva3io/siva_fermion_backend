package webstores

import (
	// "context"

	"fmt"

	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"

	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
	omnichannel_repo "fermion/backend_core/internal/repository/omnichannel"
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
	RegisterWebstores(data *omnichannel.User_Webstore_Link, token_id string, access_template_id string) error
	UpdateWebstoreDetails(id uint, data omnichannel.User_Webstore_Link, token_id string, access_template_id string) error
	FindAllWebstore(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]omnichannel.User_Webstore_Link, error)
	GetWebstoreDetail(id uint, token_id string, access_template_id string) (interface{}, error)
	DeleteWebstoreDetail(id uint, user_id uint, token_id string, access_template_id string) error
	AvailableWebstores(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]omnichannel.Webstore, error)
	GetAuthKeys(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
}
type service struct {
	webstoreRepository omnichannel_repo.Webstore
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{omnichannel_repo.NewWebstore()}
	return newServiceObj

}

func (s *service) RegisterWebstores(data *omnichannel.User_Webstore_Link, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "WEBSTORES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create webstores at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create webstores at data level")
	}
	err := s.webstoreRepository.CreateWebstore(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateWebstoreDetails(id uint, data omnichannel.User_Webstore_Link, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "WEBSTORES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update webstores at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update webstores at data level")
	}
	_, er := s.webstoreRepository.FindByidWebstore(id)
	if er != nil {
		return er
	}

	count, err := s.webstoreRepository.UpdateWebstore(id, data)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("no data gets updated")
	}

	return nil
}

func (s *service) FindAllWebstore(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]omnichannel.User_Webstore_Link, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "WEBSTORES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list webstores at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list webstores at data level")
	}
	result, err := s.webstoreRepository.FindAllWebstore(p)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetWebstoreDetail(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "WEBSTORES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view webstores at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view webstores at data level")
	}
	cout, err := s.webstoreRepository.FindByidWebstore(id)
	if err != nil {
		return cout, err
	}

	return cout, nil
}

func (s *service) DeleteWebstoreDetail(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "WEBSTORES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete webstores at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete webstores at data level")
	}
	_, er := s.webstoreRepository.FindByidWebstore(id)
	if er != nil {
		return er
	}

	err := s.webstoreRepository.DeleteWebstore(id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AvailableWebstores(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]omnichannel.Webstore, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "WEBSTORES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list webstores at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list webstores at data level")
	}
	result, err := s.webstoreRepository.AvailableWebstores(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *service) GetAuthKeys(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	result, err := s.webstoreRepository.GetAuthKey(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
