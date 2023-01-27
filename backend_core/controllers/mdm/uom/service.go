package uom

import (
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
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
	CreateUom(metaData core.MetaData, data *mdm.Uom) error
	CreateUomClass(metaData core.MetaData, data *mdm.UomClass) error

	UpdateUom(metaData core.MetaData, data *mdm.Uom) error
	UpdateUomClass(metaData core.MetaData, data *mdm.UomClass) error

	DeleteUom(metaData core.MetaData) error
	DeleteUomClass(metaData core.MetaData) error

	GetUom(metaData core.MetaData) (interface{}, error)
	GetUomClass(metaData core.MetaData) (interface{}, error)

	GetUomList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	GetUomClassList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
}
type service struct {
	uom mdm_repo.Uom
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		mdm_repo.NewUom(),
	}
	return newServiceObj
}

func (s *service) CreateUom(metaData core.MetaData, data *mdm.Uom) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create uom at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.uom.UomSave(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) CreateUomClass(metaData core.MetaData, data *mdm.UomClass) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create uom at data level")
	}
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId

	err := s.uom.UomClassSave(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateUom(metaData core.MetaData, data *mdm.Uom) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update uom at data level")
	}

	data.UpdatedByID = &metaData.TokenUserId
	err := s.uom.UpdateUom(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateUomClass(metaData core.MetaData, data *mdm.UomClass) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update uom at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.uom.UpdateUomClass(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteUom(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete uom at data level")
	}

	err := s.uom.DeleteUom(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteUomClass(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete uom at data level")
	}

	err := s.uom.DeleteUomClass(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetUom(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view uom at data level")
	}
	response, err := s.uom.FindOneUom(metaData.Query)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s *service) GetUomClass(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view uom at data level")
	}

	response, err := s.uom.FindOneUomClass(metaData.Query)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s *service) GetUomList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "UOM", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list uom at data level")
	}

	result, err := s.uom.FindAllUom(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (s *service) GetUomClassList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "UOM", metaData.TokenUserId)

	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list uom at data level")
	}
	result, err := s.uom.FindAllUomClass(metaData.Query, p)
	if err != nil {
		return result, err
	}

	return result, nil
}
