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
type BasicInventory interface {
	CentralizedInventorySave(data *mdm.CentralizedBasicInventory) error
	DecentralizedInventorySave(data *mdm.DecentralizedBasicInventory) error

	UpdateCentralizedInventory(query map[string]interface{}, data *mdm.CentralizedBasicInventory) error
	UpdateDecentralizedInventory(query map[string]interface{}, data *mdm.DecentralizedBasicInventory) error

	DeleteCentralizedInventory(query map[string]interface{}) error
	DeleteDecentralizedInventory(query map[string]interface{}) error

	FindOneCentralizedInventory(query map[string]interface{}) (mdm.CentralizedBasicInventory, error)
	FindOneDecentralizedInventory(query map[string]interface{}) (mdm.DecentralizedBasicInventory, error)

	FindAllCentralizedInventory(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.CentralizedBasicInventory, error)
	FindAllDecentralizedInventory(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.DecentralizedBasicInventory, error)

	CentralizedInventoryTransactionSave(data *mdm.CentralizedInventoryTransactions) error
}

type basicInventory struct {
	db *gorm.DB
}

var basicInventoryRepository *basicInventory //singleton object

// singleton function
func NewBasicInventory() *basicInventory {
	if basicInventoryRepository != nil {
		return basicInventoryRepository
	}
	db := db.DbManager()
	basicInventoryRepository = &basicInventory{db}
	return basicInventoryRepository

}

func (r *basicInventory) CentralizedInventorySave(data *mdm.CentralizedBasicInventory) error {
	err := r.db.Model(&mdm.CentralizedBasicInventory{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *basicInventory) DecentralizedInventorySave(data *mdm.DecentralizedBasicInventory) error {
	err := r.db.Model(&mdm.DecentralizedBasicInventory{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *basicInventory) UpdateCentralizedInventory(query map[string]interface{}, data *mdm.CentralizedBasicInventory) error {
	AvaliableQuantityPayload := map[string]interface{}{
		"available_stock": 0,
		"stock_expected":  0,
		"commited_stock":  0,
	}
	if data.AvailableStock == 0 || data.StockExpected == 0 || data.CommitedStock == 0 {
		err := r.db.Model(&mdm.CentralizedBasicInventory{}).Where(query).Updates(AvaliableQuantityPayload)
		if err.RowsAffected == 0 {
			return errors.New("oops! record not found")
		}
		if err.Error != nil {
			return err.Error
		}
		return nil
	}
	err := r.db.Model(&mdm.CentralizedBasicInventory{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *basicInventory) UpdateDecentralizedInventory(query map[string]interface{}, data *mdm.DecentralizedBasicInventory) error {

	err := r.db.Model(&mdm.DecentralizedBasicInventory{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *basicInventory) DeleteCentralizedInventory(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.CentralizedBasicInventory{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *basicInventory) DeleteDecentralizedInventory(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.DecentralizedBasicInventory{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil

}
func (r *basicInventory) FindOneCentralizedInventory(query map[string]interface{}) (mdm.CentralizedBasicInventory, error) {
	var data mdm.CentralizedBasicInventory
	err := r.db.Model(&mdm.CentralizedBasicInventory{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *basicInventory) FindOneDecentralizedInventory(query map[string]interface{}) (mdm.DecentralizedBasicInventory, error) {
	var data mdm.DecentralizedBasicInventory
	err := r.db.Model(&mdm.DecentralizedBasicInventory{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *basicInventory) FindAllCentralizedInventory(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.CentralizedBasicInventory, error) {
	var data []mdm.CentralizedBasicInventory
	err := r.db.Preload(clause.Associations).Model(&mdm.CentralizedBasicInventory{}).Scopes(helpers.Paginate(&mdm.CentralizedBasicInventory{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *basicInventory) FindAllDecentralizedInventory(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.DecentralizedBasicInventory, error) {
	var data []mdm.DecentralizedBasicInventory
	err := r.db.Preload(clause.Associations).Model(&mdm.DecentralizedBasicInventory{}).Scopes(helpers.Paginate(&mdm.DecentralizedBasicInventory{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *basicInventory) CentralizedInventoryTransactionSave(data *mdm.CentralizedInventoryTransactions) error {
	err := r.db.Model(&mdm.CentralizedInventoryTransactions{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
