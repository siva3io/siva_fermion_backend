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
type Webstore interface {
	CreateWebstore(data *omnichannel.User_Webstore_Link) error
	UpdateWebstore(id uint, data omnichannel.User_Webstore_Link) (int64, error)
	FindAllWebstore(p *pagination.Paginatevalue) ([]omnichannel.User_Webstore_Link, error)
	FindByidWebstore(id uint) (omnichannel.User_Webstore_Link, error)
	DeleteWebstore(id uint, user_id uint) error
	AvailableWebstores(p *pagination.Paginatevalue) ([]omnichannel.Webstore, error)
	GetAuthKey(query map[string]interface{}) (omnichannel.User_Webstore_Link, error)

	// <-------------------------------webstore------------------------------>
	CreateWebstoreData(data *omnichannel.Webstore) error
	FindWebstore(id uint) (omnichannel.Webstore, error)
	UpdateWebstoreData(id uint, data omnichannel.Webstore) error
	DeleteWebstoreData(id uint, user_data uint) error
}

type webstore struct {
	db *gorm.DB
}

var webstoreRepository *webstore //singleton object

// singleton function
func NewWebstore() *webstore {
	if webstoreRepository != nil {
		return webstoreRepository
	}
	db := db.DbManager()
	webstoreRepository = &webstore{db}
	return webstoreRepository
}

func (r *webstore) CreateWebstore(data *omnichannel.User_Webstore_Link) error {

	res := r.db.Model(&omnichannel.User_Webstore_Link{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
func (r *webstore) UpdateWebstore(id uint, data omnichannel.User_Webstore_Link) (int64, error) {

	res := r.db.Model(&omnichannel.User_Webstore_Link{}).Where("id", id).Updates(&data)
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}
func (r *webstore) FindAllWebstore(p *pagination.Paginatevalue) ([]omnichannel.User_Webstore_Link, error) {

	var result []omnichannel.User_Webstore_Link
	res := r.db.Preload(clause.Associations).Model(&omnichannel.User_Webstore_Link{}).Preload(clause.Associations).Scopes(helpers.Paginate(&omnichannel.User_Webstore_Link{}, p, r.db)).Where("is_active = true").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
func (r *webstore) FindByidWebstore(id uint) (omnichannel.User_Webstore_Link, error) {

	var market_user omnichannel.User_Webstore_Link

	res := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&omnichannel.User_Webstore_Link{}).Preload(clause.Associations).Where("id", id).First(&market_user)
	if res.Error != nil {
		return market_user, res.Error
	}
	return market_user, nil
}
func (r *webstore) DeleteWebstore(id uint, user_id uint) error {

	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}

	result := r.db.Model(&omnichannel.User_Webstore_Link{}).Where("id", id).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *webstore) AvailableWebstores(p *pagination.Paginatevalue) ([]omnichannel.Webstore, error) {

	var result []omnichannel.Webstore
	res := r.db.Model(&omnichannel.Webstore{}).Preload(clause.Associations).Scopes(helpers.Paginate(&omnichannel.Webstore{}, p, r.db)).Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil

}

func (r *webstore) GetAuthKey(query map[string]interface{}) (omnichannel.User_Webstore_Link, error) {
	var data omnichannel.User_Webstore_Link
	err := r.db.Model(&omnichannel.User_Webstore_Link{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *webstore) CreateWebstoreData(data *omnichannel.Webstore) error {
	res := r.db.Model(&omnichannel.Webstore{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *webstore) FindWebstore(id uint) (omnichannel.Webstore, error) {

	var result omnichannel.Webstore
	res := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&omnichannel.Webstore{}).Where("id", id).First(&result)
	if res.Error != nil {
		return result, res.Error
	}
	return result, res.Error
}

func (r *webstore) UpdateWebstoreData(id uint, data omnichannel.Webstore) error {
	res := r.db.Model(&omnichannel.Webstore{}).Where("id", id).Updates(&data).Error
	if res != nil {
		return errors.New("unable to update")
	}

	return nil
}

func (r *webstore) DeleteWebstoreData(id uint, user_id uint) error {
	var data omnichannel.Webstore
	res := r.db.Model(&omnichannel.Webstore{}).Where("id", id).Update("deleted_by", user_id)
	if res.Error != nil {
		return res.Error
	}

	res = r.db.Model(&omnichannel.Webstore{}).Where("id", id).Delete(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
