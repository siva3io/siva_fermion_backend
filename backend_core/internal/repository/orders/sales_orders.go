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
type Sales interface {
	Save(data *orders.SalesOrders) error
	FindAll(page *pagination.Paginatevalue) (interface{}, error)
	FindOne(query map[string]interface{}) (interface{}, error)
	Update(query map[string]interface{}, data *orders.SalesOrders) error
	Delete(query map[string]interface{}) error

	SaveOrderLines(orders.SalesOrderLines) error
	UpdateOrderLines(map[string]interface{}, orders.SalesOrderLines) (int64, error)
	DeleteOrderLine(map[string]interface{}) error
	FindOrderLines(map[string]interface{}) (orders.SalesOrderLines, error)

	Search(query string) (interface{}, error)
	GetSalesHistory(productId uint, page *pagination.Paginatevalue) (interface{}, error)
}

type SalesOrders struct {
	db *gorm.DB
}

func NewSalesOrder() *SalesOrders {
	db := db.DbManager()
	return &SalesOrders{db}

}

func (r *SalesOrders) Save(data *orders.SalesOrders) error {
	err := r.db.Model(&orders.SalesOrders{}).Create(data).Error

	if err != nil {

		return err

	}

	return nil
}

func (r *SalesOrders) FindAll(page *pagination.Paginatevalue) (interface{}, error) {

	var data []orders.SalesOrders

	err := r.db.Model(&orders.SalesOrders{}).Scopes(helpers.Paginate(&orders.SalesOrders{}, page, r.db)).Preload(clause.Associations).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *SalesOrders) FindOne(query map[string]interface{}) (interface{}, error) {
	var data orders.SalesOrders

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *SalesOrders) Update(query map[string]interface{}, data *orders.SalesOrders) error {
	res := r.db.Model(&orders.SalesOrders{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *SalesOrders) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&orders.SalesOrders{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *SalesOrders) SaveOrderLines(data orders.SalesOrderLines) error {

	res := r.db.Model(&orders.SalesOrderLines{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *SalesOrders) FindOrderLines(query map[string]interface{}) (orders.SalesOrderLines, error) {
	var result orders.SalesOrderLines
	fmt.Println(query)
	res := r.db.Model(&orders.SalesOrderLines{}).Where(query).First(&result)

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *SalesOrders) UpdateOrderLines(query map[string]interface{}, data orders.SalesOrderLines) (int64, error) {
	res := r.db.Model(&orders.SalesOrderLines{}).Where(query).Updates(&data)

	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}

func (r *SalesOrders) DeleteOrderLine(query map[string]interface{}) error {
	res := r.db.Model(&orders.SalesOrderLines{}).Where(query).Delete(&orders.SalesOrderLines{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *SalesOrders) Search(query string) (interface{}, error) {
	var data []orders.SalesOrders

	fields := []string{"customer_name", "channel_name", "reference_number", "sales_order_number"}

	fields_string, values := helpers.ApplySearch(query, fields)

	err := r.db.Model(&orders.SalesOrders{}).Limit(20).Preload(clause.Associations).Where(fields_string, values...).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *SalesOrders) GetSalesHistory(productId uint, page *pagination.Paginatevalue) (interface{}, error) {
	var data []orders.SalesOrders
	var ids = make([]uint, 0)

	page.Filters = fmt.Sprintf("[[\"product_id\", \"=\", %v]]", productId)
	err := r.db.Model(&orders.SalesOrderLines{}).Select("so_id").Scopes(helpers.Paginate(&orders.SalesOrderLines{}, page, r.db)).Scan(&ids)

	if err.Error != nil {
		return nil, err.Error
	}
	err = r.db.Model(&orders.SalesOrders{}).Where("id IN ?", ids).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}
