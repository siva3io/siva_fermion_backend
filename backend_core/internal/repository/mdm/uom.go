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
type Uom interface {
	UomSave(data *mdm.Uom) error
	UomClassSave(data *mdm.UomClass) error

	UpdateUom(query map[string]interface{}, data *mdm.Uom) error
	UpdateUomClass(query map[string]interface{}, data *mdm.UomClass) error

	DeleteUom(query map[string]interface{}) error
	DeleteUomClass(query map[string]interface{}) error

	FindOneUom(query map[string]interface{}) (mdm.Uom, error)
	FindOneUomClass(query map[string]interface{}) (mdm.UomClass, error)

	FindAllUom(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.Uom, error)
	FindAllUomClass(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.UomClass, error)
}

type uom struct {
	Db *gorm.DB
}

var uomRepository *uom //singleton object

// singleton function
func NewUom() *uom {
	if uomRepository != nil {
		return uomRepository
	}
	db := db.DbManager()
	uomRepository = &uom{db}
	return uomRepository
}

func (r *uom) UomSave(data *mdm.Uom) error {

	err := r.Db.Model(&mdm.Uom{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *uom) UomClassSave(data *mdm.UomClass) error {

	err := r.Db.Model(&mdm.UomClass{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *uom) UpdateUom(query map[string]interface{}, data *mdm.Uom) error {

	err := r.Db.Model(&mdm.Uom{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *uom) UpdateUomClass(query map[string]interface{}, data *mdm.UomClass) error {

	err := r.Db.Model(&mdm.UomClass{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *uom) DeleteUom(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.Db.Model(&mdm.Uom{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *uom) DeleteUomClass(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}

	//=========delete class related uom's=============================
	uomDeleteQuery := map[string]interface{}{
		"uom_class_id": query["id"],
		"user_id":      query["user_id"],
		"company_id":   query["company_id"],
	}
	r.DeleteUom(uomDeleteQuery)

	delete(query, "user_id")
	//=========== delete uom class ==================================
	err := r.Db.Model(&mdm.UomClass{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *uom) FindOneUom(query map[string]interface{}) (mdm.Uom, error) {
	var data mdm.Uom
	err := r.Db.Preload(clause.Associations).Model(&mdm.Uom{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *uom) FindOneUomClass(query map[string]interface{}) (mdm.UomClass, error) {
	var data mdm.UomClass
	err := r.Db.Preload(clause.Associations).Model(&mdm.UomClass{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *uom) FindAllUom(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.Uom, error) {
	var data []mdm.Uom
	err := r.Db.Preload(clause.Associations).Model(&mdm.Uom{}).Scopes(helpers.Paginate(&mdm.Uom{}, p, r.Db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *uom) FindAllUomClass(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.UomClass, error) {
	var data []mdm.UomClass
	err := r.Db.Preload(clause.Associations).Model(&mdm.UomClass{}).Scopes(helpers.Paginate(&mdm.UomClass{}, p, r.Db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
