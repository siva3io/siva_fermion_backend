package marketplace

import (
	"fmt"

	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
	omnichannel_repo "fermion/backend_core/internal/repository/omnichannel"
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
	//---------------------Registration-------------------------------
	RegisterMarketplace(data *omnichannel.User_Marketplace_Registration, token_id string, access_template_id string) error
	UpdatemarketplaceDetails(id uint, data omnichannel.User_Marketplace_Registration, token_id string, access_template_id string) error
	FindAllMarketplace(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
	FindByidMarketplace(id uint, token_id string, access_template_id string) (interface{}, error)
	DeleteMarketplace(id uint, user_id uint, token_id string, access_template_id string) error

	//-------------------------Link-------------------------------------
	SaveMarketplaceDetails(data *omnichannel.User_Marketplace_Link, token_id string, access_template_id string) error
	UpdateDetails(id uint, data omnichannel.User_Marketplace_Link, token_id string, access_template_id string) error
	UpdateMarketDetailsByQuery(query map[string]interface{}, data omnichannel.User_Marketplace_Link, token_id string, access_template_id string) error
	FindMarketplaceDetails(id uint, token_id string, access_template_id string) (interface{}, error)
	ListMarketplaceDetails(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]omnichannel.User_Marketplace_Link, error)
	DeleteMarketplaceDetails(id uint, user_id uint, token_id string, access_template_id string) error
	GetAuthKeys(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)

	//------------------------Marketplace-----------------------------------------
	AvailableMarketPlaces(p *pagination.Paginatevalue, token_id string, access_template_id string) (interface{}, error)

	//---------------------Marketplace-------------------------------------------
	SaveMarketplace(data *omnichannel.Marketplace, token_id string, access_template_id string) error
	UpdateMarketplace(id uint, data omnichannel.Marketplace, token_id string, access_template_id string) error
	FindMarketplace(id uint, token_id string, access_template_id string) (interface{}, error)
}

type service struct {
	marketRepository        omnichannel_repo.Marketplace
	marketdetailsrepository omnichannel_repo.MarketplaceDetails
	marketplaceRepository   omnichannel_repo.OmnichannelBase
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{omnichannel_repo.NewMarketplace(), omnichannel_repo.NewMarketPlaceDetails(), omnichannel_repo.NewmarketPlaceBase()}
	return newServiceObj
}

// Register marketplace Services

func (s *service) RegisterMarketplace(data *omnichannel.User_Marketplace_Registration, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create marketplace at data level")
	}

	err := s.marketRepository.CreateMarketplace(data)
	return err
}

func (s *service) UpdatemarketplaceDetails(id uint, data omnichannel.User_Marketplace_Registration, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update marketplace at data level")
	}
	// To check marketplace is present or not
	_, er := s.marketRepository.FindByIDMarketplace(id)
	if er != nil {
		return er
	}
	err := s.marketRepository.Updatemarketplace(id, &data)
	return err
}

func (s *service) FindAllMarketplace(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list marketplace at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list marketplace at data level")
	}
	result, err := s.marketRepository.FindAllmarketplace(p)
	if err != nil {
		return result, nil
	}
	return result, nil
}

func (s *service) FindByidMarketplace(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view marketplace at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view marketplace at data level")
	}
	result, err := s.marketRepository.FindByIDMarketplace(id)
	return result, err
}

func (s *service) DeleteMarketplace(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete marketplace at data level")
	}
	// To check marketplace is present or not
	_, er := s.marketRepository.FindByIDMarketplace(id)
	if er != nil {
		return er
	}

	err := s.marketRepository.DeleteMarketplace(id, user_id)
	return err
}

// Marketplace details Services

func (s *service) SaveMarketplaceDetails(data *omnichannel.User_Marketplace_Link, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create marketplace at data level")
	}
	err := s.marketdetailsrepository.CreateMarketplaceDetail(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindMarketplaceDetails(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view marketplace at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view marketplace at data level")
	}
	query := map[string]interface{}{
		"id": id,
	}
	result, err := s.marketdetailsrepository.FindByStoreName(query)
	if err != nil {
		return result, err
	}

	return result, nil
}
func (s *service) UpdateDetails(id uint, data omnichannel.User_Marketplace_Link, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update marketplace at data level")
	}
	query := map[string]interface{}{
		"id": id,
	}
	_, er := s.marketdetailsrepository.FindByStoreName(query)
	if er != nil {
		return er
	}
	err := s.marketdetailsrepository.UpdateMarketplaceDetail(query, data)
	return err
}

func (s *service) UpdateMarketDetailsByQuery(query map[string]interface{}, data omnichannel.User_Marketplace_Link, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update marketplace at data level")
	}
	_, er := s.marketdetailsrepository.FindByStoreName(query)
	if er != nil {
		return er
	}
	err := s.marketdetailsrepository.UpdateMarketplaceDetail(query, data)
	return err
}

func (s *service) ListMarketplaceDetails(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]omnichannel.User_Marketplace_Link, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list marketplace at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list marketplace at data level")
	}
	result, err := s.marketdetailsrepository.FindAllMarketplaceDetail(p)
	return result, err

	// return data, nil
}

func (s *service) DeleteMarketplaceDetails(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete marketplace at data level")
	}
	err := s.marketdetailsrepository.DeleteByStoreName(id, user_id)
	return err
}

func (s *service) AvailableMarketPlaces(p *pagination.Paginatevalue, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list marketplace at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list marketplace at data level")
	}
	result, err := s.marketRepository.AvailableMarketPlaces(p)
	if err != nil {
		return result, nil
	}
	return result, nil
}

func (s *service) GetAuthKeys(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	result, err := s.marketdetailsrepository.GetAuthKey(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) SaveMarketplace(data *omnichannel.Marketplace, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create marketplace at data level")
	}
	err := s.marketplaceRepository.CreateMarketplace(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateMarketplace(id uint, data omnichannel.Marketplace, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update marketplace at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update marketplace at data level")
	}

	_, er := s.marketplaceRepository.FindMarketplace(id)
	if er != nil {
		return er
	}
	err := s.marketplaceRepository.UpdateMarketplace(id, data)
	return err
}

func (s *service) FindMarketplace(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "MARKETPLACES", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view marketplace at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view marketplace at data level")
	}
	result, err := s.marketplaceRepository.FindMarketplace(id)
	if err != nil {
		return result, err
	}

	return result, nil
}
