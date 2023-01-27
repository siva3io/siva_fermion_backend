package shipping

import (
	"errors"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/shipping"

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
type ShippingBase interface {
	CreateRecord(data shipping.UserShippingFavUnfav) error
	Favourite(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (shipping.UserShippingFavUnfav, error)
	UnFavourite(query map[string]interface{}) error
}

type shippingBase struct {
	db *gorm.DB
}

var shippingBaseRepository *shippingBase //singleton object

// singleton function
func NewShippingBase() *shippingBase {
	if shippingBaseRepository != nil {
		return shippingBaseRepository
	}
	db := db.DbManager()
	shippingBaseRepository = &shippingBase{db}
	return shippingBaseRepository

}

func (r *shippingBase) CreateRecord(data shipping.UserShippingFavUnfav) error {
	err := r.db.Model(&shipping.UserShippingFavUnfav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *shippingBase) Favourite(query map[string]interface{}) error {
	err := r.db.Model(&shipping.UserShippingFavUnfav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *shippingBase) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&shipping.UserShippingFavUnfav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *shippingBase) FindOne(query map[string]interface{}) (shipping.UserShippingFavUnfav, error) {
	var data shipping.UserShippingFavUnfav

	err := r.db.Preload(clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	return data, nil
}
