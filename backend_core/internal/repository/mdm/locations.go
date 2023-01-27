package mdm

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
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
type Locations interface {
	//-------------------------------Locations--------------------------------------------
	SaveLocation(data *shared_pricing_and_location.Locations) error
	UpdateLocation(query map[string]interface{}, data *shared_pricing_and_location.Locations) error
	DeleteLocation(query map[string]interface{}) error
	FindOneLocation(query map[string]interface{}) (shared_pricing_and_location.Locations, error)
	FindAllLocation(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Locations, error)

	//------------------------------Virtual Locations----------------------------------------
	SaveVirtualLocation(data *shared_pricing_and_location.VirtualLocation) error
	UpdateVirtualLocation(query map[string]interface{}, data *shared_pricing_and_location.VirtualLocation) error
	DeleteVirtualLocation(query map[string]interface{}) error
	FindOneVirtualLocation(query map[string]interface{}) (shared_pricing_and_location.VirtualLocation, error)

	//------------------------------Local Warehouse Locations ------------------------------
	SaveWarehouseLocation(data *shared_pricing_and_location.LocalWarehouse) error
	UpdateWarehouseLocation(query map[string]interface{}, data *shared_pricing_and_location.LocalWarehouse) error
	DeleteWarehouseLocation(query map[string]interface{}) error
	FindOneWarehouseLocation(query map[string]interface{}) (shared_pricing_and_location.LocalWarehouse, error)
	//----------------------------Storage Location--------------------------------------
	SaveStorageLocation(data *shared_pricing_and_location.StorageLocation) error
	GetStorageLocationList(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.StorageLocation, error)
	GetStorageLocation(query map[string]interface{}) (shared_pricing_and_location.StorageLocation, error)
	GetStorageQuantityList(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.StorageQuantity, error)

	//------------------------------Retail Locations-------------------------------------------
	SaveRetailLocation(data *shared_pricing_and_location.Retail) error
	UpdateRetailLocation(query map[string]interface{}, data *shared_pricing_and_location.Retail) error
	DeleteRetailLocation(query map[string]interface{}) error
	FindOneRetailLocation(query map[string]interface{}) (shared_pricing_and_location.Retail, error)

	//------------------------------Office Locations ------------------------------
	SaveOfficeLocation(data *shared_pricing_and_location.Office) error
	UpdateOfficeLocation(query map[string]interface{}, data *shared_pricing_and_location.Office) error
	DeleteOfficeLocation(query map[string]interface{}) error
	FindOneOfficeLocation(query map[string]interface{}) (shared_pricing_and_location.Office, error)
}

type locations struct {
	db *gorm.DB
}

var locationsRepository *locations //singleton object

// singleton function
func NewLocations() *locations {
	if locationsRepository != nil {
		return locationsRepository
	}
	db := db.DbManager()
	locationsRepository = &locations{db}
	return locationsRepository
}

// ---------------------------------Locations----------------------------------
func (r *locations) SaveLocation(data *shared_pricing_and_location.Locations) error {
	err := r.db.Model(&shared_pricing_and_location.Locations{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *locations) UpdateLocation(query map[string]interface{}, data *shared_pricing_and_location.Locations) error {
	err := r.db.Model(&shared_pricing_and_location.Locations{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) DeleteLocation(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.Locations{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) FindOneLocation(query map[string]interface{}) (shared_pricing_and_location.Locations, error) {
	var data shared_pricing_and_location.Locations
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Model(&shared_pricing_and_location.Locations{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}
func (r *locations) FindAllLocation(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Locations, error) {
	var data []shared_pricing_and_location.Locations
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.Locations{}).Scopes(helpers.Paginate(&shared_pricing_and_location.Locations{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// ---------------------------------Virtaul Locations--------------------------------
func (r *locations) SaveVirtualLocation(data *shared_pricing_and_location.VirtualLocation) error {
	err := r.db.Model(&shared_pricing_and_location.VirtualLocation{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *locations) UpdateVirtualLocation(query map[string]interface{}, data *shared_pricing_and_location.VirtualLocation) error {
	err := r.db.Model(&shared_pricing_and_location.VirtualLocation{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) DeleteVirtualLocation(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.VirtualLocation{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) FindOneVirtualLocation(query map[string]interface{}) (shared_pricing_and_location.VirtualLocation, error) {
	var data shared_pricing_and_location.VirtualLocation
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.VirtualLocation{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

// ---------------------------------Warehouse Locations--------------------------------
func (r *locations) SaveWarehouseLocation(data *shared_pricing_and_location.LocalWarehouse) error {
	err := r.db.Model(&shared_pricing_and_location.LocalWarehouse{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *locations) UpdateWarehouseLocation(query map[string]interface{}, data *shared_pricing_and_location.LocalWarehouse) error {
	err := r.db.Model(&shared_pricing_and_location.LocalWarehouse{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) DeleteWarehouseLocation(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.LocalWarehouse{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) FindOneWarehouseLocation(query map[string]interface{}) (shared_pricing_and_location.LocalWarehouse, error) {
	var data shared_pricing_and_location.LocalWarehouse
	err := r.db.Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations).Model(&shared_pricing_and_location.LocalWarehouse{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

//-----------------------------------------------Storage Location--------------------------------------------

func (r *locations) SaveStorageLocation(data *shared_pricing_and_location.StorageLocation) error {
	err := r.db.Model(&shared_pricing_and_location.StorageLocation{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *locations) GetStorageLocationList(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.StorageLocation, error) {
	var data []shared_pricing_and_location.StorageLocation
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.StorageLocation{}).Scopes(helpers.Paginate(shared_pricing_and_location.StorageLocation{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *locations) GetStorageLocation(query map[string]interface{}) (shared_pricing_and_location.StorageLocation, error) {
	var data shared_pricing_and_location.StorageLocation
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.StorageLocation{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *locations) GetStorageQuantityList(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.StorageQuantity, error) {
	var data []shared_pricing_and_location.StorageQuantity
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.StorageQuantity{}).Scopes(helpers.Paginate(shared_pricing_and_location.StorageQuantity{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// ---------------------------------Retails Locations--------------------------------
func (r *locations) SaveRetailLocation(data *shared_pricing_and_location.Retail) error {
	err := r.db.Model(&shared_pricing_and_location.Retail{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *locations) UpdateRetailLocation(query map[string]interface{}, data *shared_pricing_and_location.Retail) error {
	err := r.db.Model(&shared_pricing_and_location.Retail{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) DeleteRetailLocation(query map[string]interface{}) error {

	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.Retail{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) FindOneRetailLocation(query map[string]interface{}) (shared_pricing_and_location.Retail, error) {
	var data shared_pricing_and_location.Retail
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.Retail{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

// ---------------------------------Office Locations--------------------------------
func (r *locations) SaveOfficeLocation(data *shared_pricing_and_location.Office) error {
	err := r.db.Model(&shared_pricing_and_location.Office{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *locations) UpdateOfficeLocation(query map[string]interface{}, data *shared_pricing_and_location.Office) error {
	err := r.db.Model(&shared_pricing_and_location.Office{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) DeleteOfficeLocation(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}

	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.Office{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *locations) FindOneOfficeLocation(query map[string]interface{}) (shared_pricing_and_location.Office, error) {
	var data shared_pricing_and_location.Office
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.Office{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
