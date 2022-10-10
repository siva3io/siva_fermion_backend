package uom

import (
	"encoding/json"
	"fmt"

	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
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
	CreateUom(data *mdm.Uom, token_id string, access_template_id string) error
	CreateUomClass(data *mdm.UomClass, token_id string, access_template_id string) error

	UpdateUom(query map[string]interface{}, data *mdm.Uom, token_id string, access_template_id string) error
	UpdateUomClass(query map[string]interface{}, data *mdm.UomClass, token_id string, access_template_id string) error

	DeleteUom(query map[string]interface{}, token_id string, access_template_id string) error
	DeleteUomClass(query map[string]interface{}, token_id string, access_template_id string) error

	GetUom(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetUomClass(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)

	GetUomList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
	GetUomClassList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)

	SearchUom(query string, token_id string, access_template_id string) ([]IdNameDTO, error)
	SearchUomClass(query string, token_id string, access_template_id string) ([]IdNameDTO, error)
}

type service struct {
	uom mdm_repo.Uom
}

func NewService() *service {
	uom := mdm_repo.NewUom()
	return &service{uom}
}

func (s *service) CreateUom(data *mdm.Uom, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "UOM", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create uom at data level")
	}
	err := s.uom.UomSave(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) CreateUomClass(data *mdm.UomClass, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "UOM", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create uom at data level")
	}
	err := s.uom.UomClassSave(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateUom(query map[string]interface{}, data *mdm.Uom, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "UOM", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update uom at data level")
	}
	_, err := s.uom.FindOneUom(query)
	if err != nil {
		return err
	}
	err = s.uom.UpdateUom(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateUomClass(query map[string]interface{}, data *mdm.UomClass, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "UOM", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update uom at data level")
	}
	_, err := s.uom.FindOneUomClass(query)
	if err != nil {
		return err
	}
	err = s.uom.UpdateUomClass(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteUom(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "UOM", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete uom at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.uom.FindOneUom(q)
	if er != nil {
		return er
	}
	err := s.uom.DeleteUom(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteUomClass(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "UOM", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete uom at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.uom.FindOneUomClass(q)
	if er != nil {
		return er
	}
	err := s.uom.DeleteUomClass(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetUom(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "UOM", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view uom at data level")
	}
	result, err := s.uom.FindOneUom(query)
	if err != nil {
		return result, err
	}
	var response UomResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) GetUomClass(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "UOM", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view uom at data level")
	}
	result, err := s.uom.FindOneUomClass(query)
	if err != nil {
		return result, err
	}
	var response UomClassResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) GetUomList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "UOM", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list uom at data level")
	}
	result, err := s.uom.FindAllUom(query, p)
	if err != nil {
		return result, err
	}
	var response []UomResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) GetUomClassList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "UOM", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list uom at data level")
	}
	result, err := s.uom.FindAllUomClass(query, p)
	if err != nil {
		return result, err
	}
	var response []UomClassResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) SearchUom(query string, token_id string, access_template_id string) ([]IdNameDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "UOM", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list uom at data level")
	}
	result, err := s.uom.SearchUom(query)
	var uomDTO []IdNameDTO
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		var data IdNameDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		uomDTO = append(uomDTO, data)
	}
	return uomDTO, nil
}
func (s *service) SearchUomClass(query string, token_id string, access_template_id string) ([]IdNameDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "UOM", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list uom at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list uom at data level")
	}
	result, err := s.uom.SearchUom(query)
	var uomClassDTO []IdNameDTO
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		var data IdNameDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		uomClassDTO = append(uomClassDTO, data)
	}
	return uomClassDTO, nil
}
