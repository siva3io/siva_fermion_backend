package shipping_orders_wd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
	shipping_repo "fermion/backend_core/internal/repository/shipping"
	access_checker "fermion/backend_core/pkg/util/access"
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
	UpdateWD(metaData core.MetaData, data *shipping.WD) error
	DeleteWD(metaData core.MetaData) error
	CreateWD(metaData core.MetaData, data *shipping.WD) error
	BulkCreateWD(metaData core.MetaData, data *[]WDRequest) error
	GetAllWD(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	GetWD(metaData core.MetaData) (interface{}, error)
}

type service struct {
	wdRepository shipping_repo.WD
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	wdRepository := shipping_repo.NewWD()
	newServiceObj = &service{wdRepository}
	return newServiceObj
}

func (s *service) CreateWD(metaData core.MetaData, data *shipping.WD) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "WD", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for createwd at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for createwd at data level")
	}
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	err := s.wdRepository.CreateWD(data)
	if err != nil {
		return err
	}
	return err
}

func (s *service) BulkCreateWD(metaData core.MetaData, data *[]WDRequest) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "WD", metaData.TokenUserId)
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

func (s *service) GetWD(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "WD", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for viewwd at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for viewwd at data level")
	}
	response, err := s.wdRepository.GetWD(metaData.Query)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *service) GetAllWD(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "WD", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for listwd at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for listwd at data level")
	}
	result, err := s.wdRepository.GetAllWD(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateWD(metaData core.MetaData, data *shipping.WD) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "WD", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for updatewd at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for updatewd at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.wdRepository.UpdateWD(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteWD(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "WD", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for deletewd at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for deletewd at data level")
	}

	err := s.wdRepository.DeleteWD(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
