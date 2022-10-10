package shipping_partners

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
	CreateUserShippingPartnerRegistration(data *UserShippingPartnerRegistrationRequest, userID uint, token_id string, access_template_id string) (uint, error)
	UpdateUserShippingPartnerRegistration(id uint, data UserShippingPartnerRegistrationRequest, token_id string, access_template_id string) error
	UserShippingPartnerRegistrationbyid(id uint, token_id string, access_template_id string) (interface{}, error)
	AllUserShippingPartnerRegistration(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.UserShippingPartnerRegistration, error)
	DeleteUserShippingPartnerRegistration(id uint, user_id uint, token_id string, access_template_id string) error

	ShippingPartnerEstimateCosts(data *shipping.RateCalculator, token_id string, access_template_id string) ([]shipping.RateCalculator, error)

	AllShippingPartner(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]shipping.ShippingPartner, error)
	GetShippingPartnerAuthByName(name string, token_id string, access_template_id string) (interface{}, error)
	UpdateShippingPartnerByName(query map[string]interface{}, token_id string, access_template_id string) (int64, error)
	GetShippingPartnerById(id int, token_id string, access_template_id string) (interface{}, error)
}

type service struct {
	shippingPartnerRepository shipping_repo.ShippingPartner
}

func NewService() *service {
	return &service{shipping_repo.NewShipping()}
}

func (s *service) CreateUserShippingPartnerRegistration(data *UserShippingPartnerRegistrationRequest, userID uint, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create logistic partner at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create logistic partner at data level")
	}
	var SPRData shipping.UserShippingPartnerRegistration
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &SPRData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	result, _ := helpers.UpdateStatusHistory(SPRData.StatusHistory, SPRData.StatusId)
	SPRData.StatusHistory = result

	id, err := s.shippingPartnerRepository.Create(&SPRData, userID)
	return id, err
}

func (s *service) UpdateUserShippingPartnerRegistration(id uint, data UserShippingPartnerRegistrationRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update logistic partner at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update logistic partner at data level")
	}
	var SPRData shipping.UserShippingPartnerRegistration
	dto, err := json.Marshal(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &SPRData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	old_data, er := s.shippingPartnerRepository.FindOne(id)
	if er != nil {
		return er
	}
	old_status := old_data.StatusId
	new_status := SPRData.StatusId
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		SPRData.StatusHistory = result
	}
	_, err = s.shippingPartnerRepository.FindOne(id)
	if err != nil {
		return err
	}
	err1 := s.shippingPartnerRepository.Update(id, SPRData)
	if err1 != nil {
		return err1
	}
	return nil
}

func (s *service) UserShippingPartnerRegistrationbyid(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view logistic partner at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view logistic partner at data level")
	}
	result, err := s.shippingPartnerRepository.FindOne(id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) AllUserShippingPartnerRegistration(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]shipping.UserShippingPartnerRegistration, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "LOGISTICS_PARTNERS", *token_user_id)
	fmt.Println(access_action)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list logistic partner at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list logistic partner at data level")
	}
	res, err := s.shippingPartnerRepository.FindAll(p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *service) DeleteUserShippingPartnerRegistration(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete logistic partner at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete logistic partner at data level")
	}
	_, err := s.shippingPartnerRepository.FindOne(id)
	if err != nil {
		return err
	}
	er := s.shippingPartnerRepository.Delete(id, user_id)
	if er != nil {
		return er
	}
	return nil
}

func (s *service) ShippingPartnerEstimateCosts(data *shipping.RateCalculator, token_id string, access_template_id string) ([]shipping.RateCalculator, error) {
	result, err := s.shippingPartnerRepository.ShippingPartnerEstimateCosts(data)
	return result, err
}

func (s service) AllShippingPartner(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]shipping.ShippingPartner, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list logistic partner at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list logistic partner at data level")
	}
	res, err := s.shippingPartnerRepository.FindAllShippingpartner(p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s service) GetShippingPartnerAuthByName(name string, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view logistic partner at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view logistic partner at data level")
	}
	res, err := s.shippingPartnerRepository.FindOneShippingpartnerByName(name)
	if err != nil {
		return res, err
	}
	return res.AuthOptions, nil
}

func (s service) UpdateShippingPartnerByName(query map[string]interface{}, token_id string, access_template_id string) (int64, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for update logistic partner at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for update logistic partner at data level")
	}
	name := query["partner_name"].(string)
	_, err := s.shippingPartnerRepository.FindOneShippingpartnerByName(name)
	if err != nil {
		return 0, err
	}
	resp, err := s.shippingPartnerRepository.UpdateShippingPartnerByName(query)
	if resp != 0 {
		return resp, nil
	}
	return resp, err
}

func (s service) GetShippingPartnerById(id int, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOGISTICS_PARTNERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view logistic partner at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view logistic partner at data level")
	}
	res, err := s.shippingPartnerRepository.FindOneShippingpartnerById(id)
	if err != nil {
		return res, err
	}
	return res, nil
}
