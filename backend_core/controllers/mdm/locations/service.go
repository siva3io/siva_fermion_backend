package locations

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/repository"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
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
	//------------------------Locations--------------------------
	CreateLocation(metaData core.MetaData, data *shared_pricing_and_location.Locations) error
	UpdateLocation(metaData core.MetaData, data *shared_pricing_and_location.Locations) error
	DeleteLocation(metaData core.MetaData) error
	GetLocation(metaData core.MetaData) (interface{}, error)
	GetLocationList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)

	//------------------------Virtual Location--------------------------
	CreateVirtualLocation(metaData core.MetaData, data *shared_pricing_and_location.VirtualLocation) error
	UpdateVirtualLocation(metaData core.MetaData, data *shared_pricing_and_location.VirtualLocation) error
	DeleteVirtualLocation(metaData core.MetaData) error
	GetVirtualLocation(metaData core.MetaData) (interface{}, error)

	//------------------------Retail Location--------------------------
	CreateRetailLocation(metaData core.MetaData, data *shared_pricing_and_location.Retail) error
	UpdateRetailLocation(metaData core.MetaData, data *shared_pricing_and_location.Retail) error
	DeleteRetailLocation(metaData core.MetaData) error
	GetRetailLocation(metaData core.MetaData) (interface{}, error)

	//------------------------Warehouse Location--------------------------
	CreateWarehouseLocation(metaData core.MetaData, data *shared_pricing_and_location.LocalWarehouse) error
	UpdateWarehouseLocation(metaData core.MetaData, data *shared_pricing_and_location.LocalWarehouse) error
	DeleteWarehouseLocation(metaData core.MetaData) error
	GetWarehouseLocation(metaData core.MetaData) (interface{}, error)

	CreateStorageLocation(metaData core.MetaData, data *shared_pricing_and_location.StorageLocation) error
	GetStorageLocationList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	GetStorageLocation(metaData core.MetaData) (interface{}, error)
	GetStorageQuantityList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)

	//------------------------Office or Factory Location--------------------------
	CreateOfficeLocation(metaData core.MetaData, data *shared_pricing_and_location.Office) error
	UpdateOfficeLocation(metaData core.MetaData, data *shared_pricing_and_location.Office) error
	DeleteOfficeLocation(metaData core.MetaData) error
	GetOfficeLocation(metaData core.MetaData) (interface{}, error)
}

type service struct {
	locationsRepository mdm_repo.Locations
	coreRepository      repository.Core
	pricingRepository   mdm_repo.Pricing
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	locationsRepository := mdm_repo.NewLocations()
	coreRepository := repository.NewCore()
	pricingRepository := mdm_repo.NewPricing()
	newServiceObj = &service{locationsRepository, coreRepository, pricingRepository}
	return newServiceObj
}

// ------------------------Locations----------------------------------------------------
func (s *service) CreateLocation(metaData core.MetaData, data *shared_pricing_and_location.Locations) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create location at data level")
	}

	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	err := s.locationsRepository.SaveLocation(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateLocation(metaData core.MetaData, data *shared_pricing_and_location.Locations) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.locationsRepository.UpdateLocation(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteLocation(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	q := map[string]interface{}{
		"id":         metaData.Query["id"],
		"company_id": metaData.CompanyId,
	}
	location, er := s.locationsRepository.FindOneLocation(q)
	if er != nil {
		return er
	}

	lookupCode, _ := helpers.GetLookupcodeName(location.LocationTypeId)
	fmt.Println(lookupCode)

	//------------Related LocationType------------------
	DeleteLocationMap := map[string]interface{}{
		"VIRTUAL_LOCATION": s.DeleteVirtualLocation,
		"LOCAL_WAREHOUSE":  s.DeleteWarehouseLocation,
		"RETAIL":           s.DeleteRetailLocation,
		"OFFICE":           s.DeleteOfficeLocation,
	}

	//---------------Delete the  Related Location Type--------------------
	metaData.Query = map[string]interface{}{
		"id":         location.RelatedLocationID,
		"company_id": metaData.CompanyId,
		"user_id":    metaData.TokenUserId,
	}
	err := DeleteLocationMap[lookupCode].(func(core.MetaData) error)(metaData)
	if err != nil {
		return err
	}

	fmt.Println("I'm here")

	//------------------ Delete the Location-----------------------------
	query := map[string]interface{}{
		"id":         location.ID,
		"company_id": metaData.CompanyId,
		"user_id":    metaData.TokenUserId,
	}
	err = s.locationsRepository.DeleteLocation(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetLocation(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}
	result, err := s.locationsRepository.FindOneLocation(metaData.Query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *service) GetLocationList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list location at data level")
	}

	//--------Fetch the location details from the Location--------
	result, err := s.locationsRepository.FindAllLocation(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ------------------------Virtual Location----------------------------------------------------
func (s *service) CreateVirtualLocation(metaData core.MetaData, data *shared_pricing_and_location.VirtualLocation) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create location at data level")
	}

	//-----------Save the Virtual Location --------------------------------

	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	err := s.locationsRepository.SaveVirtualLocation(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateVirtualLocation(metaData core.MetaData, data *shared_pricing_and_location.VirtualLocation) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}

	data.UpdatedByID = &metaData.TokenUserId
	err := s.locationsRepository.UpdateVirtualLocation(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteVirtualLocation(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteVirtualLocation(metaData.Query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetVirtualLocation(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}
	//--------Fetch the virtual location information ---------------
	result, err := s.locationsRepository.FindOneVirtualLocation(metaData.Query)
	if err != nil {
		return nil, err
	}

	var response VirtualLocationDetailsResponseDTO
	helpers.JsonMarshaller(result, &response)

	return response, nil
}

// ------------------------Retail Location----------------------------------------------------
func (s *service) CreateRetailLocation(metaData core.MetaData, data *shared_pricing_and_location.Retail) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create location at data level")
	}

	//-----------Save the Retail Location --------------------------------

	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	err := s.locationsRepository.SaveRetailLocation(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateRetailLocation(metaData core.MetaData, data *shared_pricing_and_location.Retail) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.locationsRepository.UpdateRetailLocation(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteRetailLocation(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteRetailLocation(metaData.Query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetRetailLocation(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}

	//--------Fetch the retail location information ---------------
	result, err := s.locationsRepository.FindOneRetailLocation(metaData.Query)
	if err != nil {
		return nil, err
	}
	if result.PriceListId != 0 {
		query := map[string]interface{}{
			"id": result.PriceListId,
		}
		pricingdetails, err := s.pricingRepository.FindPricing(query)
		if err != nil {
			return nil, err
		}
		jsondata, _ := json.Marshal(pricingdetails)
		result.PriceList = jsondata
	}

	var response RetailResponseDTO
	helpers.JsonMarshaller(result, &response)

	return response, nil
}

// ------------------------Warehouse Location----------------------------------------------------
func (s *service) CreateWarehouseLocation(metaData core.MetaData, data *shared_pricing_and_location.LocalWarehouse) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create location at data level")
	}

	//----------- Save the Warehouse Location --------------------------------

	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	err := s.locationsRepository.SaveWarehouseLocation(data)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) UpdateWarehouseLocation(metaData core.MetaData, data *shared_pricing_and_location.LocalWarehouse) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}

	data.UpdatedByID = &metaData.TokenUserId
	err := s.locationsRepository.UpdateWarehouseLocation(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteWarehouseLocation(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteWarehouseLocation(metaData.Query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetWarehouseLocation(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}

	//--------Fetch the warehouse location information ---------------
	result, err := s.locationsRepository.FindOneWarehouseLocation(metaData.Query)
	if err != nil {
		return nil, err
	}

	var response LocalWarehouseDetailsResponseDTO
	helpers.JsonMarshaller(result, &response)

	return response, nil
}

// ------------------------------------------Storage Location----------------------------------------------------
func (s *service) CreateStorageLocation(metaData core.MetaData, data *shared_pricing_and_location.StorageLocation) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create location at data level")
	}
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId

	err := s.locationsRepository.SaveStorageLocation(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetStorageLocationList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list location at data level")
	}
	result, err := s.locationsRepository.GetStorageLocationList(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *service) GetStorageLocation(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "LOCATIONS", metaData.TokenUserId)
	var response StorageLocationDTO
	if !access_module_flag {
		return response, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return response, fmt.Errorf("you dont have access for list location at data level")
	}
	result, err := s.locationsRepository.GetStorageLocation(metaData.Query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *service) GetStorageQuantityList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list location at data level")
	}
	result, err := s.locationsRepository.GetStorageQuantityList(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ------------------------Office or Factory Location----------------------------------------------------
func (s *service) CreateOfficeLocation(metaData core.MetaData, data *shared_pricing_and_location.Office) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create location at data level")
	}

	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId

	//-------Save the Office Location ---------------------
	err := s.locationsRepository.SaveOfficeLocation(data)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) UpdateOfficeLocation(metaData core.MetaData, data *shared_pricing_and_location.Office) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}

	data.UpdatedByID = &metaData.TokenUserId
	err := s.locationsRepository.UpdateOfficeLocation(metaData.Query, data)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) DeleteOfficeLocation(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteOfficeLocation(metaData.Query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetOfficeLocation(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "LOCATIONS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}

	//--------Fetch the Office location information ---------------
	result, err := s.locationsRepository.FindOneOfficeLocation(metaData.Query)
	if err != nil {
		return nil, err
	}

	var response OfficeOrFactoryResponseDTO
	helpers.JsonMarshaller(result, &response)

	return response, nil
}
