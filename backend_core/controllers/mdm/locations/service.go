package locations

import (
	"encoding/json"
	"fmt"

	app_core "fermion/backend_core/controllers/cores"
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
	CreateLocation(data shared_pricing_and_location.Locations, token_id string, access_template_id string) (interface{}, error)
	UpdateLocation(query map[string]interface{}, data shared_pricing_and_location.Locations, token_id string, access_template_id string) error
	DeleteLocation(query map[string]interface{}, token_id string, access_template_id string) error
	GetLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetLocationList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
	SearchLocation(query string, token_id string, access_template_id string) (interface{}, error)

	//------------------------Virtual Location--------------------------
	CreateVirtualLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error)
	UpdateVirtualLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error
	DeleteVirtualLocation(query map[string]interface{}, token_id string, access_template_id string) error
	GetVirtualLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)

	//------------------------Retail Location--------------------------
	CreateRetailLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error)
	UpdateRetailLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error
	DeleteRetailLocation(query map[string]interface{}, token_id string, access_template_id string) error
	GetRetailLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)

	//------------------------Warehouse Location--------------------------
	CreateWarehouseLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error)
	UpdateWarehouseLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error
	DeleteWarehouseLocation(query map[string]interface{}, token_id string, access_template_id string) error
	GetWarehouseLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)

	GetStorageLocationList(query map[string]interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string) ([]StorageLocationDTO, error)
	GetStorageLocation(query map[string]interface{}, token_id string, access_template_id string) (StorageLocationDTO, error)
	CreateStorageLocation(data StorageLocationDTO, token_id string, access_template_id string) (interface{}, error)
	// GetAuthKeys(query map[string]interface{}) (interface{}, error)
	GetStorageQuantityList(query map[string]interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string) ([]StorageQuantityDTO, error)

	//------------------------Office or Factory Location--------------------------
	CreateOfficeLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error)
	UpdateOfficeLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error
	DeleteOfficeLocation(query map[string]interface{}, token_id string, access_template_id string) error
	GetOfficeLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
}

type service struct {
	locationsRepository mdm_repo.Locations
	coreRepository      repository.Core
	pricingRepository   mdm_repo.Pricing
}

func NewService() *service {
	locationsRepository := mdm_repo.NewLocations()
	coreRepository := repository.NewCore()
	pricingRepository := mdm_repo.NewPricing()
	return &service{locationsRepository, coreRepository, pricingRepository}
}

// ------------------------Locations----------------------------------------------------
func (s *service) CreateLocation(data shared_pricing_and_location.Locations, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for create location at data level")
	}
	err := s.locationsRepository.SaveLocation(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (s *service) UpdateLocation(query map[string]interface{}, data shared_pricing_and_location.Locations, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}
	err := s.locationsRepository.UpdateLocation(query, &data)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) DeleteLocation(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	location, er := s.locationsRepository.FindOneLocation(q)
	if er != nil {
		return er
	}

	//--------Fetch Location Type from LookupCode-----------------------------------
	lookupcode_query := map[string]interface{}{"id": location.LocationTypeId}
	lookupcode, err := s.coreRepository.GetLookupCodes(lookupcode_query, new(pagination.Paginatevalue))
	if err != nil || len(lookupcode) == 0 {
		return err
	}

	//------------Related LocationType------------------
	DeleteLocationMap := map[string]interface{}{
		"VIRTUAL_LOCATION": s.DeleteVirtualLocation,
		"LOCAL_WAREHOUSE":  s.DeleteWarehouseLocation,
		"RETAIL":           s.DeleteRetailLocation,
		"OFFICE":           s.DeleteOfficeLocation,
	}

	//---------------Delete the  Related Location Type--------------------
	location_query := map[string]interface{}{"id": location.RelatedLocationID}
	err = DeleteLocationMap[lookupcode[0].LookupCode].(func(map[string]interface{}, string, string) error)(location_query, token_id, access_template_id)
	if err != nil {
		return err
	}

	//------------------ Delete the Location-----------------------------
	err = s.locationsRepository.DeleteLocation(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}
	var location LocationResponseDTO

	//--------Fetch the location details from the Location--------
	result, err := s.locationsRepository.FindOneLocation(query)
	if err != nil {
		return nil, err
	}
	var locationType app_core.LookupCodesDTO
	jsondata, _ := json.Marshal(result.LocationType)
	err = json.Unmarshal(jsondata, &locationType)
	if err != nil {
		return nil, err
	}
	// fmt.Println("-->>>>", locationType)
	location.LocationType = locationType

	//------Convert the location information result into a LocationResponseDTO---------
	jsondata, _ = json.Marshal(result)
	err = json.Unmarshal(jsondata, &location)
	if err != nil {
		return nil, err
	}

	//--------Fetch Location Type from LookupCode---------
	lookupcode_query := map[string]interface{}{"id": result.LocationTypeId}
	lookupcode, err := s.coreRepository.GetLookupCodes(lookupcode_query, new(pagination.Paginatevalue))
	if err != nil {
		return nil, err
	}

	//------------Related LocationType------------------
	GetRelatedLocationMap := map[string]interface{}{
		"VIRTUAL_LOCATION": s.GetVirtualLocation,
		"LOCAL_WAREHOUSE":  s.GetWarehouseLocation,
		"RETAIL":           s.GetRetailLocation,
		"OFFICE":           s.GetOfficeLocation,
	}

	//---------------Get the  Related Location Type details and append to LocationResponseDTO--------------------
	location_query := map[string]interface{}{"id": result.RelatedLocationID}
	RelatedLocationResponse, err := GetRelatedLocationMap[lookupcode[0].LookupCode].(func(map[string]interface{}, string, string) (interface{}, error))(location_query, token_id, access_template_id)
	if err != nil {
		return nil, err
	}

	jsondata, err = json.Marshal(RelatedLocationResponse)
	if err != nil {
		return nil, err
	}
	location.LocationDetails = jsondata

	return location, nil
}
func (s *service) GetLocationList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list location at data level")
	}
	var location_list []LocationListResponseDTO

	//--------Fetch the location details from the Location--------
	result, err := s.locationsRepository.FindAllLocation(query, p)
	if err != nil {
		return nil, err
	}

	//------Convert the location information result into a LocationListResponseDTO---------
	jsondata, _ := json.Marshal(result)
	err = json.Unmarshal(jsondata, &location_list)
	if err != nil {
		return nil, err
	}

	for index := range result {
		var locationType app_core.LookupCodesDTO
		jsondata, _ := json.Marshal(result[index].LocationType)
		err = json.Unmarshal(jsondata, &locationType)
		if err != nil {
			return nil, err
		}
		// fmt.Println("-->>>>", locationType)
		location_list[index].LocationType = locationType
	}

	return location_list, nil
}
func (s *service) SearchLocation(query string, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list location at data level")
	}
	//--------Fetch the available search location details from the Location--------
	result, err := s.locationsRepository.SearchLocation(query)
	var locationData []IdNameDTO
	if err != nil {
		return nil, err
	}

	//------Convert the location information result into a IdNameDTO(id, name)---------
	for _, v := range result {
		var data IdNameDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		locationData = append(locationData, data)
	}

	return locationData, nil
}

// ------------------------Virtual Location----------------------------------------------------
func (s *service) CreateVirtualLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for create location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Virtual shared_pricing_and_location.VirtualLocation

	//------convert DTO to Location objcet-------------------------
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return nil, err
	}

	//--------convert DTO to Virtual Location objcet-----------------------
	dto, err = json.Marshal(data.LocationDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Virtual)
	if err != nil {
		return nil, err
	}

	//-----------Save the Virtual Location --------------------------------
	err = s.locationsRepository.SaveVirtualLocation(&Virtual)
	if err != nil {
		return nil, err
	}

	//----------Append the Virtual Location InchargeDetails to the Location fields------------
	dto, err = json.Marshal(Virtual.LocationInchargerDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return nil, err
	}

	//---------------Save Location Details with Related Location ID---------
	Location.RelatedLocationID = Virtual.ID
	response, err := s.CreateLocation(Location, token_id, access_template_id)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) UpdateVirtualLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Virtual shared_pricing_and_location.VirtualLocation

	//------convert DTO to Location objcet-------------------------
	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return err
	}

	//-----Fetch the details of the location----------------
	result, err := s.locationsRepository.FindOneLocation(query)
	if err != nil {
		return err
	}

	//--------convert DTO to Virtual Location objcet-----------------------
	err = json.Unmarshal([]byte(data.LocationDetails), &Virtual)
	if err != nil {
		return err
	}

	//----Update Query for Location Type and Update Virtual Location--------------
	location_query := map[string]interface{}{"id": result.RelatedLocationID}
	err = s.locationsRepository.UpdateVirtualLocation(location_query, &Virtual)
	if err != nil {
		return err
	}

	//----------Append the Virtual Location InchargeDetails to the Location fields------------
	dto, err = json.Marshal(Virtual.LocationInchargerDetails)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return err
	}

	//---------------Update  Location Details with Related Location ID---------
	err = s.UpdateLocation(query, Location, token_id, access_template_id)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) DeleteVirtualLocation(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteVirtualLocation(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetVirtualLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}
	var virtual_location VirtualLocationDetailsResponseDTO

	//--------Fetch the virtual location information ---------------
	result, err := s.locationsRepository.FindOneVirtualLocation(query)
	if err != nil {
		return nil, err
	}

	//------Convert the virtual location information into a VirtualLocationDetailsResponseDTO---------------
	jsondata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsondata, &virtual_location)
	if err != nil {
		return nil, err
	}

	return virtual_location, nil
}

// ------------------------Retail Location----------------------------------------------------
func (s *service) CreateRetailLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for create location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Retail shared_pricing_and_location.Retail

	//------convert DTO to Location objcet-------------------------
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return nil, err
	}

	//--------convert DTO to Retail Location objcet-----------------------
	dto, err = json.Marshal(data.LocationDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Retail)
	if err != nil {
		return nil, err
	}

	//-----------Save the Retail Location --------------------------------
	err = s.locationsRepository.SaveRetailLocation(&Retail)
	if err != nil {
		return nil, err
	}

	//---Append the Retail Location InchargeDetails to the Location fields------------
	dto, err = json.Marshal(Retail.LocationInchargerDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return nil, err
	}

	//---------------Save Location Details with Related Location ID---------
	Location.RelatedLocationID = Retail.ID
	response, err := s.CreateLocation(Location, token_id, access_template_id)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) UpdateRetailLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Retail shared_pricing_and_location.Retail

	//------convert DTO to Location objcet-------------------------
	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return err
	}

	//-----Fetch the details of the location----------------
	result, err := s.locationsRepository.FindOneLocation(query)
	if err != nil {
		return err
	}

	//--------convert DTO to Retail Location objcet-----------------------
	err = json.Unmarshal([]byte(data.LocationDetails), &Retail)
	if err != nil {
		return err
	}

	//----Update Query for Location Type and Update Retail Location--------------
	location_query := map[string]interface{}{"id": result.RelatedLocationID}
	err = s.locationsRepository.UpdateRetailLocation(location_query, &Retail)
	if err != nil {
		return err
	}

	//----------Append the Retail Location InchargeDetails to the Location fields------------
	dto, err = json.Marshal(Retail.LocationInchargerDetails)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return err
	}

	//---------------Update Location Details with Related Location ID---------
	err = s.UpdateLocation(query, Location, token_id, access_template_id)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) DeleteRetailLocation(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteRetailLocation(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetRetailLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}
	var retail RetailResponseDTO

	//--------Fetch the retail location information ---------------
	result, err := s.locationsRepository.FindOneRetailLocation(query)
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

	//------Convert the retail location information into a RetailResponseDTO---------------
	jsondata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsondata, &retail)
	if err != nil {
		return nil, err
	}

	return retail, nil
}

// ------------------------Warehouse Location----------------------------------------------------
func (s *service) CreateWarehouseLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for create location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Warehouse shared_pricing_and_location.LocalWarehouse

	//------convert DTO to Location objcet-------------------------
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return nil, err
	}

	//--------convert DTO to Warehouse Location objcet-----------------------
	dto, err = json.Marshal(data.LocationDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Warehouse)
	if err != nil {
		return nil, err
	}

	//-----------Save the Warehouse Location --------------------------------
	err = s.locationsRepository.SaveWarehouseLocation(&Warehouse)
	if err != nil {
		return nil, err
	}

	//---Append the Warehouse Location ContactDetails to the Location fields------------
	var warehouse_contact_details []interface{}
	dto, err = json.Marshal(Warehouse.ContactDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &warehouse_contact_details)
	if err != nil {
		return nil, err
	}
	if len(warehouse_contact_details) > 0 {
		contactDetails := warehouse_contact_details[0].(map[string]interface{})
		if contactDetails["email"].(string) != "" {
			Location.Email = contactDetails["email"].(string)
		}
		if contactDetails["mobile_number"] != nil && contactDetails["mobile_number"].(string) != "" {
			Location.MobileNumber = contactDetails["mobile_number"].(string)
		}
		// dto, err = json.Marshal(warehouse_contact_details[0])
		// if err != nil {
		// 	return nil, err
		// }
		// err = json.Unmarshal(dto, &Location)
		// if err != nil {
		// 	return nil, err
		// }
	}

	//---------------Save Location Details with Related Location ID---------
	Location.RelatedLocationID = Warehouse.ID
	err = s.locationsRepository.SaveLocation(&Location)
	if err != nil {
		return nil, err
	}

	return Location, nil
}
func (s *service) UpdateWarehouseLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Warehouse shared_pricing_and_location.LocalWarehouse

	//------convert DTO to Location objcet-------------
	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return err
	}

	//-----Fetch the details of the location----------------
	result, err := s.locationsRepository.FindOneLocation(query)
	if err != nil {
		return err
	}

	//--------convert DTO to Warehouse Location objcet-----------------------
	err = json.Unmarshal([]byte(data.LocationDetails), &Warehouse)
	if err != nil {
		return err
	}

	//----Update Query for Location Type and Update Warehouse Location--------------
	location_query := map[string]interface{}{"id": result.RelatedLocationID}
	err = s.locationsRepository.UpdateWarehouseLocation(location_query, &Warehouse)
	if err != nil {
		return err
	}

	//----------Append the Warehouse Location Contactdetails to the Location fields------------
	var warehouse_contact_details []interface{}
	dto, err = json.Marshal(Warehouse.ContactDetails)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &warehouse_contact_details)
	if err != nil {
		return err
	}
	if len(warehouse_contact_details) != 0 {
		dto, err = json.Marshal(warehouse_contact_details[0])
		if err != nil {
			return err
		}
		err = json.Unmarshal(dto, &Location)
		if err != nil {
			return err
		}
	}

	//---------------Update Location Details with Related Location ID---------
	err = s.UpdateLocation(query, Location, token_id, access_template_id)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) DeleteWarehouseLocation(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteWarehouseLocation(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetWarehouseLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}
	var local_warehouse LocalWarehouseDetailsResponseDTO

	//--------Fetch the warehouse location information ---------------
	result, err := s.locationsRepository.FindOneWarehouseLocation(query)
	if err != nil {
		return nil, err
	}

	//------Convert the warehouse location information into a LocalWarehouseDetailsResponseDTO---------------
	jsondata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsondata, &local_warehouse)
	if err != nil {
		return nil, err
	}

	return local_warehouse, nil
}

//------------------------------------------Storage Location----------------------------------------------------

func (s *service) CreateStorageLocation(data StorageLocationDTO, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for create location at data level")
	}
	var response_data shared_pricing_and_location.StorageLocation
	marshal_data, _ := json.Marshal(data)
	json.Unmarshal(marshal_data, &response_data)
	result := s.locationsRepository.SaveStorageLocation(&response_data)
	return result, nil
}

func (s *service) GetStorageLocationList(query map[string]interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string) ([]StorageLocationDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list location at data level")
	}
	result, err := s.locationsRepository.GetStorageLocationList(query, p)
	if err != nil {
		return nil, err
	}
	var response []StorageLocationDTO
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

func (s *service) GetStorageLocation(query map[string]interface{}, token_id string, access_template_id string) (StorageLocationDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "LOCATIONS", *token_user_id)
	var response StorageLocationDTO
	if !access_module_flag {
		return response, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return response, fmt.Errorf("you dont have access for list location at data level")
	}
	result, err := s.locationsRepository.GetStorageLocation(query)
	if err != nil {
		return response, err
	}
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(marshaldata, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *service) GetStorageQuantityList(query map[string]interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string) ([]StorageQuantityDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list location at data level")
	}
	result, err := s.locationsRepository.GetStorageQuantityList(query, p)
	if err != nil {
		return nil, err
	}
	var response []StorageQuantityDTO
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

// func (s *service) GetAuthKeys(query map[string]interface{}) (interface{}, error) {
// 	result, err := s.locationsRepository.GetAuthKey(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result.Auth, nil
// }

// ------------------------Office or Factory Location----------------------------------------------------
func (s *service) CreateOfficeLocation(data LocationRequestDTO, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for create location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for create location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Office shared_pricing_and_location.Office

	//----convert DTO to Location objcet--------------
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return nil, err
	}

	//-----convert DTO to Office Location objcet--------------
	dto, err = json.Marshal(data.LocationDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Office)
	if err != nil {
		return nil, err
	}

	//-------Save the Office Location ---------------------
	err = s.locationsRepository.SaveOfficeLocation(&Office)
	if err != nil {
		return nil, err
	}

	//---Append the Office Location Location Incharge to the Location fields------
	dto, err = json.Marshal(Office.LocationInchargerDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return nil, err
	}

	//---------Save Location Details with Related Location ID------
	Location.RelatedLocationID = Office.ID
	response, err := s.CreateLocation(Location, token_id, access_template_id)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) UpdateOfficeLocation(query map[string]interface{}, data LocationRequestDTO, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update location at data level")
	}
	var Location shared_pricing_and_location.Locations
	var Office shared_pricing_and_location.Office

	//------convert DTO to Location objcet-----------
	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return err
	}

	//-----Fetch the details of the location-------------------
	result, err := s.locationsRepository.FindOneLocation(query)
	if err != nil {
		return err
	}

	//--------convert DTO to Office Location objcet-----------------------
	err = json.Unmarshal([]byte(data.LocationDetails), &Office)
	if err != nil {
		return err
	}

	//----Update Query for Location Type and Update Office Location--------------
	location_query := map[string]interface{}{"id": result.RelatedLocationID}
	err = s.locationsRepository.UpdateOfficeLocation(location_query, &Office)
	if err != nil {
		return err
	}

	//----------Append the Office Location Incharge Details to the Location fields------------
	dto, err = json.Marshal(Office.LocationInchargerDetails)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Location)
	if err != nil {
		return err
	}

	//---------------Update Location Details with Related Location ID---------
	err = s.UpdateLocation(query, Location, token_id, access_template_id)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) DeleteOfficeLocation(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete location at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete location at data level")
	}
	err := s.locationsRepository.DeleteOfficeLocation(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetOfficeLocation(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "LOCATIONS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view location at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view location at data level")
	}
	var office_or_factory OfficeOrFactoryResponseDTO

	//--------Fetch the Office location information ---------------
	result, err := s.locationsRepository.FindOneOfficeLocation(query)
	if err != nil {
		return nil, err
	}

	//------Convert the Office location information into a OfficeOrFactoryResponseDTO---------------
	jsondata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsondata, &office_or_factory)
	if err != nil {
		return nil, err
	}

	return office_or_factory, nil
}
