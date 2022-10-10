package shipping_orders_ndr

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
	CreateNDR(data *NDRRequest, token_id string, access_template_id string) (uint, error)
	BulkCreateNDR(data *[]NDRRequest, token_id string, access_template_id string) error
	GetAllNDR(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.NDR, error)
	GetNDR(id uint, token_id string, access_template_id string) (interface{}, error)
	UpdateNDR(id uint, data NDRRequest, token_id string, access_template_id string) error
	DeleteNDR(id uint, user_id uint, token_id string, access_template_id string) error
	DeleteNDRLines(query interface{}, token_id string, access_template_id string) error
}

type service struct {
	ndrRepository shipping_repo.NDR
}

func NewService() *service {
	ndrRepository := shipping_repo.NewNDR()
	return &service{ndrRepository}
}

func (s *service) CreateNDR(data *NDRRequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "NDR", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create ndr at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create ndr at data level")
	}
	var NDRData shipping.NDR
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &NDRData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	id, err := s.ndrRepository.CreateNDR(&NDRData)
	return id, err
}

func (s *service) BulkCreateNDR(data *[]NDRRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "NDR", *token_user_id)
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

func (s *service) GetNDR(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "NDR", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view ndr at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view ndr at data level")
	}
	res, er := s.ndrRepository.GetNDR(id)
	if er != nil {
		return res, er
	}
	query := map[string]interface{}{
		"non_delivery_request_id": id,
	}
	result_non_delivery_request, err := s.ndrRepository.GetNDRLines(query)
	res.NDRLines = result_non_delivery_request
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *service) GetAllNDR(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.NDR, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "NDR", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list ndr at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list ndr at data level")
	}
	result, err := s.ndrRepository.GetAllNDR(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateNDR(id uint, data NDRRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "NDR", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update ndr at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update ndr at data level")
	}
	var NDRData shipping.NDR
	dto, err := json.Marshal(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &NDRData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	_, er := s.ndrRepository.GetNDR(id)
	if er != nil {
		return er
	}
	err = s.ndrRepository.UpdateNDR(id, NDRData)
	for _, non_delivery_request := range NDRData.NDRLines {
		query := map[string]interface{}{
			"non_delivery_request_id": uint(id),
			"ndrs_id":                 non_delivery_request.NdrsId,
		}
		fmt.Println("!!!!", query)
		count, er := s.ndrRepository.UpdateNDRLines(query, non_delivery_request)
		if er != nil {
			return er
		} else if count == 0 {
			non_delivery_request.Non_delivery_request_id = id
			e := s.ndrRepository.CreateNDRLines(non_delivery_request)
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *service) DeleteNDR(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "NDR", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ndr at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ndr at data level")
	}
	_, errmsg := s.ndrRepository.GetNDR(id)
	if errmsg != nil {
		return errmsg
	}
	err := s.ndrRepository.DeleteNDR(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"non_delivery_request_id": id}
	err1 := s.ndrRepository.DeleteNDRLines(query)
	return err1
}

func (s *service) DeleteNDRLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "NDR", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ndr at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ndr at data level")
	}
	data, errmsg := s.ndrRepository.GetNDRLines(query)
	if errmsg != nil {
		return errmsg
	}
	if len(data) <= 0 {
		return errmsg
	}
	err := s.ndrRepository.DeleteNDRLines(query)
	return err
}
