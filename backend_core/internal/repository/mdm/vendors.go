package mdm

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/mdm"
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
type Vendors interface {
	//============= vendors CRUD=====================================
	Create(data *mdm.Vendors) error
	Update(query map[string]interface{}, data *mdm.Vendors) error
	Delete(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (mdm.Vendors, error)
	FindAll(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.Vendors, error)

	//TODO : Pricelist need to remove
	//============= vendors pricelist CRUD=====================================
	CreatePriceList(data *mdm.VendorPriceLists) error
	UpdatePriceList(query map[string]interface{}, data *mdm.VendorPriceLists) error
	DeletePriceList(query map[string]interface{}) error
	FindOnePriceList(query map[string]interface{}) (mdm.VendorPriceLists, error)
	FindAllPriceList(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.VendorPriceLists, error)
}

type vendors struct {
	Db *gorm.DB
}

var vendorsRepository *vendors //singleton object

// singleton function
func NewVendors() *vendors {
	if vendorsRepository != nil {
		return vendorsRepository
	}
	db := db.DbManager()
	vendorsRepository = &vendors{db}
	return vendorsRepository
}

func (r *vendors) Create(data *mdm.Vendors) error {
	err := r.Db.Model(&mdm.Vendors{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *vendors) Update(query map[string]interface{}, data *mdm.Vendors) error {
	err := r.Db.Model(&mdm.Vendors{}).Where(query).Updates(&data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *vendors) Delete(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.Db.Model(&mdm.Vendors{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *vendors) FindOne(query map[string]interface{}) (mdm.Vendors, error) {
	var data mdm.Vendors
	err := r.Db.Preload(clause.Associations).Model(&mdm.Vendors{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *vendors) FindAll(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.Vendors, error) {
	var data []mdm.Vendors
	err := r.Db.Preload(clause.Associations + "." + clause.Associations).Model(&mdm.Vendors{}).Scopes(helpers.Paginate(&mdm.Vendors{}, p, r.Db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *vendors) CreatePriceList(data *mdm.VendorPriceLists) error {
	err := r.Db.Model(&mdm.VendorPriceLists{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *vendors) UpdatePriceList(query map[string]interface{}, data *mdm.VendorPriceLists) error {
	res := r.Db.Model(&mdm.VendorPriceLists{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *vendors) DeletePriceList(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.Db.Model(&mdm.VendorPriceLists{}).Where(query).Updates(data).Error
	return err
}
func (r *vendors) FindAllPriceList(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.VendorPriceLists, error) {
	var data []mdm.VendorPriceLists
	err := r.Db.Preload(clause.Associations + "." + clause.Associations).Model(&mdm.VendorPriceLists{}).Scopes(helpers.Paginate(&mdm.VendorPriceLists{}, p, r.Db)).Where(query).Find(&data).Error

	if err != nil {
		return data, nil
	}
	return data, nil
}
func (r *vendors) FindOnePriceList(query map[string]interface{}) (mdm.VendorPriceLists, error) {
	var data mdm.VendorPriceLists
	err := r.Db.Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
