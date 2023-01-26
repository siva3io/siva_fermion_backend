package orders

import (
	"errors"
	"fmt"
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
type Purchase interface {
	Save(data *orders.PurchaseOrders) error
	FindAll(page *pagination.Paginatevalue) (interface{}, error)
	FindOne(query map[string]interface{}) (interface{}, error)
	Update(query map[string]interface{}, data *orders.PurchaseOrders) error
	Delete(query map[string]interface{}) error

	SaveOrderLines(orders.PurchaseOrderLines) error
	UpdateOrderLines(map[string]interface{}, orders.PurchaseOrderLines) (int64, error)
	DeleteOrderLine(map[string]interface{}) error
	FindOrderLines(map[string]interface{}) (orders.PurchaseOrderLines, error)

	Search(query string) (interface{}, error)
	GetPurchaseHistory(productId uint, page *pagination.Paginatevalue) (interface{}, error)
}

type PurchaseOrders struct {
	db *gorm.DB
}

func NewPurchaseOrder() *PurchaseOrders {
	db := db.DbManager()
	return &PurchaseOrders{db}

}

func (r *PurchaseOrders) Save(data *orders.PurchaseOrders) error {
	err := r.db.Model(&orders.PurchaseOrders{}).Create(data).Error

	if err != nil {

		return err

	}

	return nil
}

func (r *PurchaseOrders) FindAll(page *pagination.Paginatevalue) (interface{}, error) {
	var data []orders.PurchaseOrders

	err := r.db.Model(&orders.PurchaseOrders{}).Scopes(helpers.Paginate(&orders.PurchaseOrders{}, page, r.db)).Preload(clause.Associations).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *PurchaseOrders) FindOne(query map[string]interface{}) (interface{}, error) {
	var data orders.PurchaseOrders

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *PurchaseOrders) Update(query map[string]interface{}, data *orders.PurchaseOrders) error {
	res := r.db.Model(&orders.PurchaseOrders{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseOrders) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&orders.PurchaseOrders{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseOrders) SaveOrderLines(data orders.PurchaseOrderLines) error {

	res := r.db.Model(&orders.PurchaseOrderLines{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseOrders) FindOrderLines(query map[string]interface{}) (orders.PurchaseOrderLines, error) {
	var result orders.PurchaseOrderLines
	fmt.Println(query)
	res := r.db.Model(&orders.PurchaseOrderLines{}).Where(query).First(&result)

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *PurchaseOrders) UpdateOrderLines(query map[string]interface{}, data orders.PurchaseOrderLines) (int64, error) {
	res := r.db.Model(&orders.PurchaseOrderLines{}).Where(query).Updates(&data)

	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}

func (r *PurchaseOrders) DeleteOrderLine(query map[string]interface{}) error {
	res := r.db.Model(&orders.PurchaseOrderLines{}).Where(query).Delete(&orders.PurchaseOrderLines{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *PurchaseOrders) Search(query string) (interface{}, error) {
	var data []orders.PurchaseOrders

	fields := []string{"reference_number", "number"}

	fields_string, values := helpers.ApplySearch(query, fields)

	err := r.db.Model(&orders.PurchaseOrders{}).Limit(20).Preload(clause.Associations).Where(fields_string, values...).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *PurchaseOrders) GetPurchaseHistory(productId uint, page *pagination.Paginatevalue) (interface{}, error) {
	var data []orders.PurchaseOrders
	var ids = make([]uint, 0)

	page.Filters = fmt.Sprintf("[[\"product_id\", \"=\", %v]]", productId)
	err := r.db.Model(&orders.PurchaseOrderLines{}).Select("po_id").Scopes(helpers.Paginate(&orders.PurchaseOrderLines{}, page, r.db)).Scan(&ids)

	if err.Error != nil {
		return nil, err.Error
	}
	err = r.db.Model(&orders.PurchaseOrders{}).Where("id IN ?", ids).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}
