package shipping_orders_rto

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
	CreateRTO(metaData core.MetaData, data *shipping.RTO) error
	BulkCreateRTO(metaData core.MetaData, data *[]RTORequest) error
	GetAllRTO(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	GetRTO(metaData core.MetaData) (interface{}, error)
	UpdateRTO(metaData core.MetaData, data *shipping.RTO) error
	DeleteRTO(metaData core.MetaData) error
}

type service struct {
	rtoRepository shipping_repo.RTO
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	rtoRepository := shipping_repo.NewRTO()
	newServiceObj = &service{rtoRepository}
	return newServiceObj
}

func (s *service) CreateRTO(metaData core.MetaData, data *shipping.RTO) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "RTO", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create rto at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create rto at data level")
	}
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	err := s.rtoRepository.CreateRTO(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) BulkCreateRTO(metaData core.MetaData, data *[]RTORequest) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "RTO", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create rto at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create rto at data level")
	}
	var RTOData []shipping.RTO
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &RTOData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = s.rtoRepository.BulkCreateRTO(&RTOData)
	return err
}

func (s *service) GetRTO(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "RTO", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view rto at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view rto at data level")
	}
	response, err := s.rtoRepository.GetRTO(metaData.Query)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *service) GetAllRTO(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "RTO", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list rto at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list rto at data level")
	}
	result, err := s.rtoRepository.GetAllRTO(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) UpdateRTO(metaData core.MetaData, data *shipping.RTO) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "RTO", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update rto at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update rto at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.rtoRepository.UpdateRTO(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteRTO(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "RTO", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete rto at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete rto at data level")
	}
	err := s.rtoRepository.DeleteRTO(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
