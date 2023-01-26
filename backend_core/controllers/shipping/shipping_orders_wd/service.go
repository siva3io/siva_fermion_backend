package shipping_orders_wd

import (
	"encoding/json"
	"fmt"

	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
	shipping_repo "fermion/backend_core/internal/repository/shipping"
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
	UpdateWD(id uint, data WDRequest, token_id string, access_template_id string) error
	DeleteWD(id uint, user_id uint, token_id string, access_template_id string) error
	CreateWD(data *WDRequest, token_id string, access_template_id string) (uint, error)
	BulkCreateWD(data *[]WDRequest, token_id string, access_template_id string) error
	GetAllWD(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.WD, error)
	GetWD(id uint, token_id string, access_template_id string) (interface{}, error)
}

type service struct {
	wdRepository shipping_repo.WD
}

func NewService() *service {
	wdRepository := shipping_repo.NewWD()
	return &service{wdRepository}
}

func (s *service) CreateWD(data *WDRequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "WD", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for createwd at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for createwd at data level")
	}
	var WDData shipping.WD
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &WDData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	ID, err1 := s.wdRepository.CreateWD(&WDData)
	return ID, err1
}

func (s *service) BulkCreateWD(data *[]WDRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "WD", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for createwd at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for createwd at data level")
	}
	var WDData []shipping.WD
	dto, err := json.Marshal(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &WDData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = s.wdRepository.BulkCreateWD(&WDData)
	return err
}

func (s *service) GetWD(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "WD", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for viewwd at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for viewwd at data level")
	}
	result, err := s.wdRepository.GetWD(id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetAllWD(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.WD, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "WD", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for listwd at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for listwd at data level")
	}
	result, err := s.wdRepository.GetAllWD(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateWD(id uint, data WDRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "WD", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for updatewd at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for updatewd at data level")
	}
	var WDData shipping.WD
	dto, err := json.Marshal(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &WDData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	_, errmsg := s.wdRepository.GetWD(id)
	if errmsg != nil {
		return errmsg
	}
	err = s.wdRepository.UpdateWD(id, WDData)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteWD(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "WD", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for deletewd at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for deletewd at data level")
	}
	_, errmsg := s.wdRepository.GetWD(id)
	if errmsg != nil {
		return errmsg
	}
	err := s.wdRepository.DeleteWD(id, user_id)
	if err != nil {
		return err
	}
	return nil
}
