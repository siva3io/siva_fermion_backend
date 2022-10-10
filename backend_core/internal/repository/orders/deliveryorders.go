package orders

import (
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
type DeliveryOrders interface {
	CreateDeliveryOrder(data *orders.DeliveryOrders) (uint, error)
	BulkCreateDeliveryOrder(data *[]orders.DeliveryOrders) error
	UpdateDeliveryOrder(id uint, data *orders.DeliveryOrders) error
	FindDeliveryOrder(id uint) (orders.DeliveryOrders, error)
	AllDeliveryOrders(p *pagination.Paginatevalue) ([]orders.DeliveryOrders, error)
	DeleteDeliveryOrder(id uint, user_id uint) error
	CreateDeliveryOrderLines(data orders.DeliveryOrderLines) error
	UpdateDeliveryOrderLines(query interface{}, data orders.DeliveryOrderLines) (int64, error)
	FindDeliveryOrdersLines(query interface{}) (orders.DeliveryOrderLines, error)
	DeleteDeliveryOrderLines(query interface{}) error
}

type deliveryOrders struct {
	db *gorm.DB
}

func NewDo() *deliveryOrders {
	db := db.DbManager()
	return &deliveryOrders{db}
}

func (r *deliveryOrders) CreateDeliveryOrder(data *orders.DeliveryOrders) (uint, error) {

	res := r.db.Model(&orders.DeliveryOrders{}).Create(&data)
	return data.ID, res.Error
}
func (r *deliveryOrders) BulkCreateDeliveryOrder(data *[]orders.DeliveryOrders) error {
	res := r.db.Model(&orders.DeliveryOrders{}).Create(data)
	return res.Error
}

func (r *deliveryOrders) UpdateDeliveryOrder(id uint, data *orders.DeliveryOrders) error {
	res := r.db.Model(&orders.DeliveryOrders{}).Where("id", id).Updates(data)
	return res.Error
}

func (r *deliveryOrders) FindDeliveryOrder(id uint) (orders.DeliveryOrders, error) {
	var result orders.DeliveryOrders
	res := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&orders.DeliveryOrders{}).Preload("DeliveryOrderLines.Product").Preload("DeliveryOrderLines.Product_template").Preload(clause.Associations).Where("id", id).First(&result)
	return result, res.Error
}

func (r *deliveryOrders) AllDeliveryOrders(p *pagination.Paginatevalue) ([]orders.DeliveryOrders, error) {
	var data []orders.DeliveryOrders
	res := r.db.Preload(clause.Associations).Model(&orders.DeliveryOrders{}).Scopes(helpers.Paginate(&orders.DeliveryOrders{}, p, r.db)).Where("is_active = true").Find(&data)
	return data, res.Error
}

func (r *deliveryOrders) DeleteDeliveryOrder(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&orders.DeliveryOrders{}).Where("id", id).Updates(data)
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
