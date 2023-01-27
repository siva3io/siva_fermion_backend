package orders

import (
	"errors"

	"fermion/backend_core/db"

	"fermion/backend_core/internal/model/orders"

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
type OrdersBase interface {
	Create(data orders.UserOrdersFav) error
	FindOne(query map[string]interface{}) (orders.UserOrdersFav, error)
	Favourite(query map[string]interface{}) error
	UnFavourite(query map[string]interface{}) error
}

type orders_base struct {
	db *gorm.DB
}

var ordersBaseRepository *orders_base //singleton object

// singleton function
func NewOrdersBase() *orders_base {
	if ordersBaseRepository != nil {
		return ordersBaseRepository
	}
	db := db.DbManager()
	ordersBaseRepository = &orders_base{db}
	return ordersBaseRepository

}

func (r *orders_base) Create(data orders.UserOrdersFav) error {
	err := r.db.Model(&orders.UserOrdersFav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *orders_base) Favourite(query map[string]interface{}) error {
	err := r.db.Model(&orders.UserOrdersFav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *orders_base) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&orders.UserOrdersFav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *orders_base) FindOne(query map[string]interface{}) (orders.UserOrdersFav, error) {
	var data orders.UserOrdersFav

	err := r.db.Preload(clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}
