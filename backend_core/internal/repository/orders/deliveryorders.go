package orders

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/orders"
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
type DeliveryOrders interface {
	CreateDeliveryOrder(data *orders.DeliveryOrders) (uint, error)
	BulkCreateDeliveryOrder(data *[]orders.DeliveryOrders) error
	UpdateDeliveryOrder(query map[string]interface{}, data *orders.DeliveryOrders) error
	FindDeliveryOrder(query map[string]interface{}) (orders.DeliveryOrders, error)
	AllDeliveryOrders(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.DeliveryOrders, error)
	DeleteDeliveryOrder(query map[string]interface{}) error
	CreateDeliveryOrderLines(data orders.DeliveryOrderLines) error
	UpdateDeliveryOrderLines(query interface{}, data orders.DeliveryOrderLines) (int64, error)
	FindDeliveryOrdersLines(query interface{}) (orders.DeliveryOrderLines, error)
	DeleteDeliveryOrderLines(query interface{}) error
}

type deliveryOrders struct {
	db *gorm.DB
}

var deliveryOrdersRepository *deliveryOrders //singleton object

// singleton function
func NewDo() *deliveryOrders {
	if deliveryOrdersRepository != nil {
		return deliveryOrdersRepository
	}
	db := db.DbManager()
	deliveryOrdersRepository = &deliveryOrders{db}
	return deliveryOrdersRepository
}

func (r *deliveryOrders) CreateDeliveryOrder(data *orders.DeliveryOrders) (uint, error) {

	res := r.db.Model(&orders.DeliveryOrders{}).Create(&data)
	return data.ID, res.Error
}
func (r *deliveryOrders) BulkCreateDeliveryOrder(data *[]orders.DeliveryOrders) error {
	res := r.db.Model(&orders.DeliveryOrders{}).Create(data)
	return res.Error
}

func (r *deliveryOrders) UpdateDeliveryOrder(query map[string]interface{}, data *orders.DeliveryOrders) error {
	res := r.db.Model(&orders.DeliveryOrders{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *deliveryOrders) FindDeliveryOrder(query map[string]interface{}) (orders.DeliveryOrders, error) {
	var data orders.DeliveryOrders

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}

func (r *deliveryOrders) AllDeliveryOrders(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.DeliveryOrders, error) {
	var data []orders.DeliveryOrders

	err := r.db.Model(&orders.DeliveryOrders{}).Preload(clause.Associations).Scopes(helpers.Paginate(&orders.DeliveryOrders{}, page, r.db)).Where(query).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *deliveryOrders) DeleteDeliveryOrder(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(uint),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&orders.DeliveryOrders{}).Where(query).Updates(data)
	return res.Error
}

func (r *deliveryOrders) CreateDeliveryOrderLines(data orders.DeliveryOrderLines) error {
	res := r.db.Model(&orders.DeliveryOrderLines{}).Create(&data)
	return res.Error
}

func (r *deliveryOrders) UpdateDeliveryOrderLines(query interface{}, data orders.DeliveryOrderLines) (int64, error) {
	res := r.db.Model(&orders.DeliveryOrderLines{}).Where(query).Updates(&data)
	return res.RowsAffected, res.Error
}

func (r *deliveryOrders) FindDeliveryOrdersLines(query interface{}) (orders.DeliveryOrderLines, error) {
	var result orders.DeliveryOrderLines
	res := r.db.Model(&orders.DeliveryOrderLines{}).Where(query).First(&result)
	return result, res.Error
}

func (r *deliveryOrders) DeleteDeliveryOrderLines(query interface{}) error {
	res := r.db.Model(&orders.DeliveryOrderLines{}).Where(query).Delete(&orders.DeliveryOrderLines{})
	return res.Error
}
