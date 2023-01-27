package vendors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/core"
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
	Upsert(metaData core.MetaData, data []interface{}) (interface{}, error)
	CreateVendors(metaData core.MetaData, data *mdm.Vendors) error
	UpdateVendors(metaData core.MetaData, data *mdm.Vendors) error
	DeleteVendors(metaData core.MetaData) error
	GetVendors(metaData core.MetaData) (interface{}, error)
	GetVendorsList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)

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

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	vendorsRepository := mdm_repo.NewVendors()
	pricingRepository := mdm_repo.NewPricing()
	newServiceObj = &service{vendorsRepository, pricingRepository}
	return newServiceObj
}

func (s *service) CreateVendors(metaData core.MetaData, data *mdm.Vendors) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PRODUCTS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create vendor at data level")
	}
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	err := s.vendorsRepository.Create(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return nil
}
func (s *service) UpdateVendors(metaData core.MetaData, data *mdm.Vendors) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "PRODUCTS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update vendor at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.vendorsRepository.Update(metaData.Query, data)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}
func (s *service) DeleteVendors(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "PRODUCTS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete vendor at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete vendor at data level")
	}

	err := s.vendorsRepository.Delete(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetVendors(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "PRODUCTS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view vendor at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view vendor at data level")
	}
	result, err := s.vendorsRepository.FindOne(metaData.Query)
	if err != nil {
		return result, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return result, nil
}
func (s *service) GetVendorsList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "PRODUCTS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list vendor at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list vendor at data level")
	}
	result, err := s.vendorsRepository.FindAll(metaData.Query, p)
	if err != nil {
		return result, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return result, nil
}
func (s *service) Upsert(metaData core.MetaData, data []interface{}) (interface{}, error) {
	tokenUserId := metaData.TokenUserId
	companyId := metaData.CompanyId

	var success []interface{}
	var failures []interface{}

	for index, payload := range data {
		vendor := new(mdm.Vendors)
		helpers.JsonMarshaller(payload, vendor)
		if *vendor.ContactId == 0 {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "contact_id should not be empty"})
			continue
		}
		findQuery := map[string]interface{}{
			"contact_id": vendor.ContactId,
		}
		response, _ := s.vendorsRepository.FindOne(findQuery)
		if response.ContactId != nil && *response.ContactId != 0 {
			vendor.UpdatedByID = &tokenUserId
			vendor.CompanyId = companyId
			err := s.vendorsRepository.Update(findQuery, vendor)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			success = append(success, map[string]interface{}{
				"status":        true,
				"serial_number": index + 1,
				"msg":           "updated",
			})
			continue
		}
		vendor.CreatedByID = &tokenUserId
		vendor.CompanyId = companyId
		err := s.vendorsRepository.Create(vendor)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		success = append(success, map[string]interface{}{
			"status":        true,
			"serial_number": index + 1,
			"msg":           "created",
		})
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
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
	result, err := s.vendorsRepository.FindAllPriceList(query.(map[string]interface{}), p)
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
	err := s.vendorsRepository.CreatePriceList(data)
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
	_, err := s.vendorsRepository.FindOnePriceList(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	err = s.vendorsRepository.UpdatePriceList(query, data)
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
	result, err := s.vendorsRepository.FindOnePriceList(query)
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
	_, er := s.vendorsRepository.FindOnePriceList(q)
	if er != nil {
		return er
	}
	err := s.vendorsRepository.DeletePriceList(query)
	if err != nil {
		return err
	}
	return nil
}
