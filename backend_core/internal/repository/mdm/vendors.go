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
type Vendors interface {
	SaveVendors(data *mdm.Vendors) error
	SaveVendorPriceLists(data *mdm.VendorPriceLists) error
	UpdateVendors(query map[string]interface{}, data *mdm.Vendors) error
	UpdateVendorPriceLists(query map[string]interface{}, data *mdm.VendorPriceLists) (int64, error)
	DeleteVendors(query map[string]interface{}) error
	FindOneVendors(query map[string]interface{}) (mdm.Vendors, error)
	FindAll(query interface{}, p *pagination.Paginatevalue) ([]mdm.Vendors, error)
	Search(q string) ([]mdm.Vendors, error)

	FindById(id string) (mdm.Vendors, error)
	FindByName(name string) (*mdm.Vendors, error)
	FindOneVendorPriceList(query map[string]interface{}) (mdm.VendorPriceLists, error)
	GetVendorPriceList(query interface{}, p *pagination.Paginatevalue) ([]mdm.VendorPriceLists, error)
	DeleteVendorPriceList(query map[string]interface{}) error
}

type vendors struct {
	db *gorm.DB
}

func NewVendors() *vendors {
	db := db.DbManager()
	return &vendors{db}
}

func (r *vendors) SaveVendors(data *mdm.Vendors) error {
	err := r.db.Model(&mdm.Vendors{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *vendors) SaveVendorPriceLists(data *mdm.VendorPriceLists) error {
	err := r.db.Model(&mdm.VendorPriceLists{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *vendors) UpdateVendors(query map[string]interface{}, data *mdm.Vendors) error {
	res := r.db.Model(&mdm.Vendors{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *vendors) UpdateVendorPriceLists(query map[string]interface{}, data *mdm.VendorPriceLists) (int64, error) {
	res := r.db.Model(&mdm.VendorPriceLists{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}
	return res.RowsAffected, nil
}
func (r *vendors) DeleteVendors(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&mdm.Vendors{}).Where(query).Updates(data)
	return res.Error
}
func (r *vendors) FindOneVendors(query map[string]interface{}) (mdm.Vendors, error) {
	var data mdm.Vendors
	err := r.db.Model(&mdm.Vendors{}).Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *vendors) FindAll(query interface{}, p *pagination.Paginatevalue) ([]mdm.Vendors, error) {
	var data []mdm.Vendors
	err := r.db.Preload(clause.Associations).Model(&mdm.Vendors{}).Scopes(helpers.Paginate(&mdm.Vendors{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *vendors) FindById(id string) (mdm.Vendors, error) {
	var data mdm.Vendors
	err := r.db.Preload(clause.Associations).Where("ID = ?", id).First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}
func (r *vendors) FindByName(name string) (*mdm.Vendors, error) {
	var data *mdm.Vendors
	err := r.db.Preload(clause.Associations).Where("name = ?", name).First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *vendors) Search(q string) ([]mdm.Vendors, error) {
	var data []mdm.Vendors
	err := r.db.Preload(clause.Associations).Find(&data, "name ILIKE ?", "%"+q+"%").Error

	//err := r.db.Where("name ILIKE ? OR primary_email ILIKE ?", "%"+q+"%", "%"+q+"%").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *vendors) GetVendorPriceList(query interface{}, p *pagination.Paginatevalue) ([]mdm.VendorPriceLists, error) {
	var data []mdm.VendorPriceLists
	err := r.db.Model(&mdm.VendorPriceLists{}).Scopes(helpers.Paginate(&mdm.VendorPriceLists{}, p, r.db)).Where(query).Find(&data).Error

	if err != nil {
		return data, nil
	}
	return data, nil
}

func (r *vendors) FindOneVendorPriceList(query map[string]interface{}) (mdm.VendorPriceLists, error) {
	var data mdm.VendorPriceLists
	err := r.db.Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *vendors) DeleteVendorPriceList(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&mdm.VendorPriceLists{}).Where(query).Updates(data)
	return res.Error
}
