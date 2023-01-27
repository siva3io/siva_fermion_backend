package shipping_orders_ndr

import (
	"encoding/json"
	"fmt"
	"strconv"

	ds "fermion/backend_core/pkg/util/dynamic_struct"

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
	CreateNDR(metaData core.MetaData, data *shipping.NDR) error
	BulkCreateNDR(metaData core.MetaData, data *[]NDRRequest) error
	GetAllNDR(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	GetNDR(metaData core.MetaData) (interface{}, error)
	UpdateNDR(metaData core.MetaData, data *shipping.NDR) error
	DeleteNDR(metaData core.MetaData) error
	DeleteNDRLines(metaData core.MetaData) error
}

type service struct {
	ndrRepository shipping_repo.NDR
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	ndrRepository := shipping_repo.NewNDR()
	newServiceObj = &service{ndrRepository}
	return newServiceObj
}

func (s *service) CreateNDR(metaData core.MetaData, data *shipping.NDR) error {
	metaDataReader := ds.NewDynamicStructReader(metaData)

	accessTemplateId := strconv.FormatUint(uint64(metaDataReader.GetField("AccessTemplateId").Uint()), 10)
	tokenUserId := metaDataReader.GetField("TokenUserId").Uint()
	companyId := metaDataReader.GetField("CompanyId").Uint()
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "NDR", tokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create contacts at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create contacts at data level")
	}
	data.CreatedByID = &tokenUserId
	data.CompanyId = companyId
	err := s.ndrRepository.CreateNDR(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) BulkCreateNDR(metaData core.MetaData, data *[]NDRRequest) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "NDR", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create ndr at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create ndr at data level")
	}
	var NDRData []shipping.NDR
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &NDRData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = s.ndrRepository.BulkCreateNDR(&NDRData)
	return err
}

func (s *service) GetNDR(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "NDR", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view ndr at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view ndr at data level")
	}
	res, er := s.ndrRepository.GetNDR(metaData.Query)
	if er != nil {
		return res, er
	}

	return res, nil
}

func (s *service) GetAllNDR(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "NDR", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list ndr at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list ndr at data level")
	}
	result, err := s.ndrRepository.GetAllNDR(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateNDR(metaData core.MetaData, data *shipping.NDR) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "NDR", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update ndr at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update ndr at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	var NDRData shipping.NDR
	dto, err := json.Marshal(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &NDRData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	res := s.ndrRepository.UpdateNDR(metaData.Query, data)
	for _, non_delivery_request := range NDRData.NDRLines {
		query := map[string]interface{}{
			"non_delivery_request_id": metaData.Query["id"],
		}
		count, er := s.ndrRepository.UpdateNDRLines(query, &non_delivery_request)
		if er != nil {
			return er
		} else if count == 0 {
			non_delivery_request.Non_delivery_request_id = uint(metaData.Query["id"].(float64))
			e := s.ndrRepository.CreateNDRLines(non_delivery_request)
			if e != nil {
				return e
			}
		}
	}
	return res
}

func (s *service) DeleteNDR(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "NDR", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ndr at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ndr at data level")
	}

	err := s.ndrRepository.DeleteNDR(metaData.Query)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteNDRLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "NDR", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ndr at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ndr at data level")
	}

	err := s.ndrRepository.DeleteNDRLines(metaData.Query)
	return err
}
