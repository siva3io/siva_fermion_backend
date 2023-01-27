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
	Save(data *orders.SalesOrders) (uint, error)
	SaveOrderLines(data *orders.SalesOrderLines) error

	Update(query map[string]interface{}, data *orders.SalesOrders) error
	UpdateOrderLines(query map[string]interface{}, data *orders.SalesOrderLines) error

	Delete(query map[string]interface{}) error
	DeleteOrderLine(map[string]interface{}) error

	FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.SalesOrders, error)

	FindOne(query map[string]interface{}) (orders.SalesOrders, error)
	FindOrderLines(map[string]interface{}) (orders.SalesOrderLines, error)

	GetSalesHistory(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.SalesOrders, error)
}

type SalesOrders struct {
	db *gorm.DB
}

var SalesOrdersRepository *SalesOrders //singleton object

// singleton function
func NewSalesOrder() *SalesOrders {
	if SalesOrdersRepository != nil {
		return SalesOrdersRepository
	}
	SalesOrdersRepository = &SalesOrders{
		db.DbManager(),
	}
	return SalesOrdersRepository

}

func (r *SalesOrders) Save(data *orders.SalesOrders) (uint, error) {

	err := r.db.Model(&orders.SalesOrders{}).Create(data).Error
	if err != nil {
		return data.ID, err
	}
	return data.ID, nil
}

func (r *SalesOrders) FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.SalesOrders, error) {

	var data []orders.SalesOrders
	err := r.db.Model(&orders.SalesOrders{}).Preload(clause.Associations + "." + clause.Associations).Preload("CreatedBy.Company").Preload("SalesOrderLines.Product.Category").Scopes(helpers.Paginate(&orders.SalesOrders{}, page, r.db)).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *SalesOrders) FindOne(query map[string]interface{}) (orders.SalesOrders, error) {

	var data orders.SalesOrders
	err := r.db.Model(&orders.SalesOrders{}).Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *SalesOrders) Update(query map[string]interface{}, data *orders.SalesOrders) error {

	err := r.db.Model(&orders.SalesOrders{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *SalesOrders) Delete(query map[string]interface{}) error {

	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&orders.SalesOrders{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *SalesOrders) SaveOrderLines(data *orders.SalesOrderLines) error {

	err := r.db.Model(&orders.SalesOrderLines{}).Create(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *SalesOrders) FindOrderLines(query map[string]interface{}) (orders.SalesOrderLines, error) {
	var data orders.SalesOrderLines
	err := r.db.Model(&orders.SalesOrderLines{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *SalesOrders) UpdateOrderLines(query map[string]interface{}, data *orders.SalesOrderLines) error {

	err := r.db.Model(&orders.SalesOrderLines{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *SalesOrders) DeleteOrderLine(query map[string]interface{}) error {

	err := r.db.Model(&orders.SalesOrderLines{}).Where(query).Delete(&orders.SalesOrderLines{})
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *SalesOrders) GetSalesHistory(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.SalesOrders, error) {
	var data []orders.SalesOrders
	var ids = make([]uint, 0)

	productId := query["product_id"]

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
