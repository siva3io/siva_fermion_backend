package omnichannel

import (
	"errors"
	"fmt"

	"fermion/backend_core/internal/model/omnichannel"

	"fermion/backend_core/db"

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
type OmnichannelBase interface {
	CreateRecord(data omnichannel.OmnichannelFav) error
	Favourite(query map[string]interface{}) error
	UnFavourite(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (omnichannel.OmnichannelFav, error)

	CreateMarketplace(data *omnichannel.Marketplace) error
	UpdateMarketplace(id uint, data omnichannel.Marketplace) error
	FindMarketplace(id uint) (omnichannel.Marketplace, error)
	DeleteMarketplace(id uint, user_id uint) error
}

type omnichannelBase struct {
	db *gorm.DB
}

var omnichannelBaseRepository *omnichannelBase //singleton object

// singleton function
func NewmarketPlaceBase() *omnichannelBase {
	if omnichannelBaseRepository != nil {
		return omnichannelBaseRepository
	}
	db := db.DbManager()
	omnichannelBaseRepository = &omnichannelBase{db}
	return omnichannelBaseRepository

}

func (r *omnichannelBase) CreateRecord(data omnichannel.OmnichannelFav) error {
	fmt.Println(data)
	err := r.db.Model(&omnichannel.OmnichannelFav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *omnichannelBase) Favourite(query map[string]interface{}) error {
	//fmt.Println(data)
	err := r.db.Model(&omnichannel.OmnichannelFav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *omnichannelBase) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&omnichannel.OmnichannelFav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *omnichannelBase) FindOne(query map[string]interface{}) (omnichannel.OmnichannelFav, error) {
	var data omnichannel.OmnichannelFav
	err := r.db.Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	return data, nil
}

func (r *omnichannelBase) CreateMarketplace(data *omnichannel.Marketplace) error {
	res := r.db.Model(&omnichannel.Marketplace{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *omnichannelBase) FindMarketplace(id uint) (omnichannel.Marketplace, error) {

	var result omnichannel.Marketplace
	res := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&omnichannel.Marketplace{}).Where("id", id).First(&result)
	if res.Error != nil {
		return result, res.Error
	}

	return result, res.Error
}

func (r *omnichannelBase) UpdateMarketplace(id uint, data omnichannel.Marketplace) error {

	res := r.db.Model(&omnichannel.Marketplace{}).Where("id", id).Updates(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
func (r *omnichannelBase) DeleteMarketplace(id uint, user_id uint) error {
	var data omnichannel.Marketplace
	res := r.db.Model(&omnichannel.Marketplace{}).Where("id", id).Update("deleted_by", user_id)
	if res.Error != nil {
		return res.Error
	}
	res = r.db.Model(&omnichannel.Marketplace{}).Where("id", id).Delete(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
