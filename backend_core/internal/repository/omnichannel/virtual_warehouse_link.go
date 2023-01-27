package omnichannel

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
type VirtualWarehouseDetails interface {
	CreateVirtualWarehouseDetail(data *omnichannel.User_Virtual_Warehouse_Link) error
	UpdateVirtualWarehouseDetail(id uint, data omnichannel.User_Virtual_Warehouse_Link) error
	FindVirtualWarehouseDetailByStoreName(id uint) (omnichannel.User_Virtual_Warehouse_Link, error)
	FindAllVirtualWarehouseDetail(p *pagination.Paginatevalue) ([]omnichannel.User_Virtual_Warehouse_Link, error)
	DeleteVirtualWarehouseDetailByStoreName(id uint, user_id uint) error

	GetAuthKey(query map[string]interface{}) (omnichannel.User_Virtual_Warehouse_Link, error)
}

type virtual_warehouse_details struct {
	Db *gorm.DB
}

var virtualWarehouseRepository *virtual_warehouse_details //singleton object

// singleton function
func NewVirtualWarehouseDetails() *virtual_warehouse_details {
	if virtualWarehouseRepository != nil {
		return virtualWarehouseRepository
	}
	db := db.DbManager()
	virtualWarehouseRepository = &virtual_warehouse_details{db}
	return virtualWarehouseRepository
}

func (r *virtual_warehouse_details) CreateVirtualWarehouseDetail(data *omnichannel.User_Virtual_Warehouse_Link) error {
	res := r.Db.Model(&omnichannel.User_Virtual_Warehouse_Link{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *virtual_warehouse_details) UpdateVirtualWarehouseDetail(id uint, data omnichannel.User_Virtual_Warehouse_Link) error {

	res := r.Db.Model(&omnichannel.User_Virtual_Warehouse_Link{}).Where("id", id).Updates(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *virtual_warehouse_details) FindVirtualWarehouseDetailByStoreName(id uint) (omnichannel.User_Virtual_Warehouse_Link, error) {

	var result omnichannel.User_Virtual_Warehouse_Link
	res := r.Db.Preload(clause.Associations+"."+clause.Associations).Model(&omnichannel.User_Virtual_Warehouse_Link{}).Where("id", id).First(&result)
	if res.Error != nil {
		return result, res.Error
	}

	return result, res.Error
}

func (r *virtual_warehouse_details) FindAllVirtualWarehouseDetail(p *pagination.Paginatevalue) ([]omnichannel.User_Virtual_Warehouse_Link, error) {
	var result []omnichannel.User_Virtual_Warehouse_Link
	res := r.Db.Preload(clause.Associations).Model(&omnichannel.User_Virtual_Warehouse_Link{}).Preload(clause.Associations).Scopes(helpers.Paginate(&omnichannel.User_Virtual_Warehouse_Link{}, p, r.Db)).Where("is_active = true").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (r *virtual_warehouse_details) DeleteVirtualWarehouseDetailByStoreName(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.Db.Model(&omnichannel.User_Virtual_Warehouse_Link{}).Where("id", id).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *virtual_warehouse_details) GetAuthKey(query map[string]interface{}) (omnichannel.User_Virtual_Warehouse_Link, error) {
	var data omnichannel.User_Virtual_Warehouse_Link
	err := r.Db.Model(&omnichannel.User_Virtual_Warehouse_Link{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
