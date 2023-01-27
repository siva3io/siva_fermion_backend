package virtual_warehouse

import (
	omnichannel_repo "fermion/backend_core/internal/repository/omnichannel"

	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
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
	RegisterVirtualWarehouse(data *omnichannel.User_Virtual_Warehouse_Registration) error
	UpdateVirtualWarehouse(id uint, data omnichannel.User_Virtual_Warehouse_Registration) error
	FindAllVirtualWarehouses(p *pagination.Paginatevalue) (interface{}, error)
	FindByidVirtualWarehouse(id uint) (interface{}, error)
	DeleteVirtualWarehouse(id uint, user_id uint) error

	//-------------------------Link-------------------------------------
	SaveVirtualWarehouseDetails(data *omnichannel.User_Virtual_Warehouse_Link) error
	UpdateVirtualWarehouseDetails(id uint, data omnichannel.User_Virtual_Warehouse_Link) error
	FindVirtualWarehouseDetails(id uint) (omnichannel.User_Virtual_Warehouse_Link, error)
	ListVirtualWarehouseDetails(p *pagination.Paginatevalue) ([]omnichannel.User_Virtual_Warehouse_Link, error)
	DeleteVirtualWarehouseDetails(id uint, user_id uint) error
	GetAuthKeys(query map[string]interface{}) (interface{}, error)

	//------------------------VirtualWarehouse-----------------------------------------
	AvailableVirtualWarehouses(p *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	virtual_warehouse_Repository         omnichannel_repo.VirtualWarehouse
	virtual_warehouse_details_Repository omnichannel_repo.VirtualWarehouseDetails
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{omnichannel_repo.NewVirtualWarehouse(), omnichannel_repo.NewVirtualWarehouseDetails()}
	return newServiceObj
	// return &service{omnichannel_repo.NewVirtualWarehouseDetails()}
}

// // Register VirtualWarehouse Services

func (s *service) RegisterVirtualWarehouse(data *omnichannel.User_Virtual_Warehouse_Registration) error {

	err := s.virtual_warehouse_Repository.CreateVirtualWarehouse(data)
	return err
}

func (s *service) UpdateVirtualWarehouse(id uint, data omnichannel.User_Virtual_Warehouse_Registration) error {

	// To check VirtualWarehouse is present or not
	_, er := s.virtual_warehouse_Repository.FindByIDVirtualWarehouse(id)
	if er != nil {
		return er
	}
	err := s.virtual_warehouse_Repository.UpdateVirtualWarehouse(id, &data)
	return err
}

func (s *service) FindAllVirtualWarehouses(p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.virtual_warehouse_Repository.FindAllVirtualWarehouses(p)
	if err != nil {
		return result, nil
	}
	return result, nil
}

func (s *service) FindByidVirtualWarehouse(id uint) (interface{}, error) {
	result, err := s.virtual_warehouse_Repository.FindByIDVirtualWarehouse(id)
	return result, err
}

func (s *service) DeleteVirtualWarehouse(id uint, user_id uint) error {

	// To check VirtualWarehouse is present or not
	_, er := s.virtual_warehouse_Repository.FindByIDVirtualWarehouse(id)
	if er != nil {
		return er
	}

	err := s.virtual_warehouse_Repository.DeleteVirtualWarehouse(id, user_id)
	return err
}

// Virtual Warehouse details Services

func (s *service) SaveVirtualWarehouseDetails(data *omnichannel.User_Virtual_Warehouse_Link) error {

	err := s.virtual_warehouse_details_Repository.CreateVirtualWarehouseDetail(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindVirtualWarehouseDetails(id uint) (omnichannel.User_Virtual_Warehouse_Link, error) {

	result, err := s.virtual_warehouse_details_Repository.FindVirtualWarehouseDetailByStoreName(id)
	if err != nil {
		return result, err
	}

	return result, nil
}
func (s *service) UpdateVirtualWarehouseDetails(id uint, data omnichannel.User_Virtual_Warehouse_Link) error {

	_, er := s.virtual_warehouse_details_Repository.FindVirtualWarehouseDetailByStoreName(id)
	if er != nil {
		return er
	}
	err := s.virtual_warehouse_details_Repository.UpdateVirtualWarehouseDetail(id, data)
	return err
}

func (s *service) ListVirtualWarehouseDetails(p *pagination.Paginatevalue) ([]omnichannel.User_Virtual_Warehouse_Link, error) {
	result, err := s.virtual_warehouse_details_Repository.FindAllVirtualWarehouseDetail(p)
	return result, err

	// return data, nil
}

func (s *service) DeleteVirtualWarehouseDetails(id uint, user_id uint) error {
	err := s.virtual_warehouse_details_Repository.DeleteVirtualWarehouseDetailByStoreName(id, user_id)
	return err
}

func (s *service) AvailableVirtualWarehouses(p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.virtual_warehouse_Repository.AvailableVirtualWarehouses(p)
	if err != nil {
		return result, nil
	}
	return result, nil
}

func (s *service) GetAuthKeys(query map[string]interface{}) (interface{}, error) {
	result, err := s.virtual_warehouse_details_Repository.GetAuthKey(query)
	if err != nil {
		return nil, err
	}
	return result.Auth, nil
}
