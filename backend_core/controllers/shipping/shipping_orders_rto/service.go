package shipping_orders_rto

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
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type Service interface {
	CreateRTO(data *RTORequest, token_id string, access_template_id string) (uint, error)
	BulkCreateRTO(data *[]RTORequest, token_id string, access_template_id string) error
	GetAllRTO(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.RTO, error)
	GetRTO(id uint, token_id string, access_template_id string) (interface{}, error)
	UpdateRTO(id uint, data RTORequest, token_id string, access_template_id string) error
	DeleteRTO(id uint, user_id uint, token_id string, access_template_id string) error
}

type service struct {
	rtoRepository shipping_repo.RTO
}

func NewService() *service {
	rtoRepository := shipping_repo.NewRTO()
	return &service{rtoRepository}
}

func (s *service) CreateRTO(data *RTORequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "RTO", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create rto at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create rto at data level")
	}
	var RTOData shipping.RTO
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &RTOData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	id, err := s.rtoRepository.CreateRTO(&RTOData)
	return id, err
}

func (s *service) BulkCreateRTO(data *[]RTORequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "RTO", *token_user_id)
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

func (s *service) GetRTO(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "RTO", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view rto at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view rto at data level")
	}
	result, err := s.rtoRepository.GetRTO(id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetAllRTO(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.RTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "RTO", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list rto at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list rto at data level")
	}
	result, err := s.rtoRepository.GetAllRTO(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateRTO(id uint, data RTORequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "RTO", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update rto at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update rto at data level")
	}
	var RTOData shipping.RTO
	dto, err := json.Marshal(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &RTOData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	_, errmsg := s.rtoRepository.GetRTO(id)
	if errmsg != nil {
		return errmsg
	}
	err = s.rtoRepository.UpdateRTO(id, RTOData)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteRTO(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "RTO", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete rto at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete rto at data level")
	}
	_, errmsg := s.rtoRepository.GetRTO(id)
	if errmsg != nil {
		return errmsg
	}
	err := s.rtoRepository.DeleteRTO(id, user_id)
	if err != nil {
		return err
	}
	return nil
}
