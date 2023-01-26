package vendors

import (
	"encoding/json"
	"errors"
	"fmt"

	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
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
	//vendors
	CreateVendors(data *mdm.Vendors, token_id string, access_template_id string) error
	UpdateVendors(query map[string]interface{}, data *mdm.Vendors, token_id string, access_template_id string) error
	DeleteVendors(query map[string]interface{}, token_id string, access_template_id string) error
	GetVendors(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetVendorsList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]mdm.Vendors, error)
	SearchVendors(query string, token_id string, access_template_id string) ([]VendorsObjDTO, error)

	//vendor price list
	CreateVendorPriceList(data *mdm.VendorPriceLists, token_id string, access_template_id string) error
	UpdateVendorPriceLists(query map[string]interface{}, data *mdm.VendorPriceLists, token_id string, access_template_id string) error
	GetVendorPriceLists(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]mdm.VendorPriceLists, error)
	GetVendorPriceList(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	DeleteVendorPricelist(query map[string]interface{}, token_id string, access_template_id string) error
}

type service struct {
	vendorsRepository mdm_repo.Vendors
	pricingRepository mdm_repo.Pricing
}

func NewService() *service {
	vendorsRepository := mdm_repo.NewVendors()
	pricingRepository := mdm_repo.NewPricing()
	return &service{vendorsRepository, pricingRepository}
}

func (s *service) CreateVendors(data *mdm.Vendors, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create vendor at data level")
	}
	err := s.vendorsRepository.SaveVendors(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return nil
}
func (s *service) UpdateVendors(query map[string]interface{}, data *mdm.Vendors, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update vendor at data level")
	}
	_, err := s.vendorsRepository.FindOneVendors(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	err = s.vendorsRepository.UpdateVendors(query, data)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}
func (s *service) DeleteVendors(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete vendor at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.vendorsRepository.FindOneVendors(q)
	if er != nil {
		return er
	}
	err := s.vendorsRepository.DeleteVendors(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetVendors(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view vendor at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view vendor at data level")
	}
	result, err := s.vendorsRepository.FindOneVendors(query)
	if err != nil {
		return result, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return result, nil
}
func (s *service) GetVendorsList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]mdm.Vendors, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list vendor at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list vendor at data level")
	}
	result, err := s.vendorsRepository.FindAll(query, p)
	if err != nil {
		return result, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return result, nil
}
func (s *service) SearchVendors(q string, token_id string, access_template_id string) ([]VendorsObjDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list vendor at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list vendor at data level")
	}
	result, err := s.vendorsRepository.Search(q)
	var vendorsDto []VendorsObjDTO
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	for _, v := range result {
		var data VendorsObjDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		vendorsDto = append(vendorsDto, data)
	}
	return vendorsDto, nil
}

//---------------------------vendor price list-----------------------------------------

func (s *service) GetVendorPriceLists(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]mdm.VendorPriceLists, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, res.BuildError(res.ErrDataNotFound, errors.New("you dont have access for list vendor at view level"))
	}
	if data_access == nil {
		return nil, res.BuildError(res.ErrDataNotFound, errors.New("you dont have access for list vendor at data level"))
	}
	result, err := s.vendorsRepository.GetVendorPriceList(query, p)
	if err != nil {
		return result, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	for index, priceList := range result {
		if priceList.PriceListId != 0 {
			query := map[string]interface{}{
				"id": priceList.PriceListId,
			}
			pricingdetails, err := s.pricingRepository.FindPricing(query)
			if err != nil {
				return result, err
			}
			jsondata, _ := json.Marshal(pricingdetails)
			result[index].PriceList = jsondata
		}
	}

	return result, nil
}

func (s *service) CreateVendorPriceList(data *mdm.VendorPriceLists, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create vendor at data level")
	}
	err := s.vendorsRepository.SaveVendorPriceLists(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return nil
}

func (s *service) UpdateVendorPriceLists(query map[string]interface{}, data *mdm.VendorPriceLists, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update vendor at data level")
	}
	_, err := s.vendorsRepository.FindOneVendorPriceList(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	_, err = s.vendorsRepository.UpdateVendorPriceLists(query, data)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}

func (s *service) GetVendorPriceList(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view vendor at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view vendor at data level")
	}
	result, err := s.vendorsRepository.FindOneVendorPriceList(query)
	if err != nil {
		return result, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	if result.PriceListId != 0 {
		query := map[string]interface{}{
			"id": result.PriceListId,
		}
		pricingdetails, err := s.pricingRepository.FindPricing(query)
		if err != nil {
			return result, err
		}
		jsondata, _ := json.Marshal(pricingdetails)
		result.PriceList = jsondata
	}
	return result, nil
}

func (s *service) DeleteVendorPricelist(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete vendor at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.vendorsRepository.FindOneVendorPriceList(q)
	if er != nil {
		return er
	}
	err := s.vendorsRepository.DeleteVendorPriceList(query)
	if err != nil {
		return err
	}
	return nil
}
