package orders

import (
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
type ScrapOrders interface {
	SaveScrapOrder(data *orders.ScrapOrders) (uint, error)
	BulkSaveScrapOrder(data *[]orders.ScrapOrders) error
	UpdateScrapOrder(id uint, data *orders.ScrapOrders) error
	FindAllScrapOrders(p *pagination.Paginatevalue) ([]orders.ScrapOrders, error)
	FindScrapOrderById(id uint) (orders.ScrapOrders, error)
	DeleteScrapOrder(id uint, user_id uint) error

	CreateScrapOrderLines(data orders.ScrapOrderLines) error
	UpdateScrapOrderLines(query interface{}, data orders.ScrapOrderLines) (int64, error)
	DeleteScrapOrderLines(query interface{}) error
	FindScrapOrderLines(query interface{}) (orders.ScrapOrderLines, error)
}

type scrapOrders struct {
	db *gorm.DB
}

func NewScrap() *scrapOrders {
	db := db.DbManager()
	return &scrapOrders{db}
}

func (r *scrapOrders) SaveScrapOrder(data *orders.ScrapOrders) (uint, error) {
	fmt.Println("repo")

	res := r.db.Model(&orders.ScrapOrders{}).Create(&data)
	return data.ID, res.Error
}
func (r *scrapOrders) BulkSaveScrapOrder(data *[]orders.ScrapOrders) error {
	fmt.Println("repo")
	res := r.db.Model(&orders.ScrapOrders{}).Create(data)
	return res.Error
}

func (r *scrapOrders) UpdateScrapOrder(id uint, Updatedata *orders.ScrapOrders) error {
	res := r.db.Model(&orders.ScrapOrders{}).Where("id", id).Updates(Updatedata)
	return res.Error
}

func (r *scrapOrders) FindAllScrapOrders(p *pagination.Paginatevalue) ([]orders.ScrapOrders, error) {
	var result []orders.ScrapOrders
	res := r.db.Preload(clause.Associations).Model(&orders.ScrapOrders{}).Scopes(helpers.Paginate(&orders.ScrapOrders{}, p, r.db)).Where("is_active = true").Preload("Order_lines.Product").Preload(clause.Associations).Find(&result)
	return result, res.Error
}

func (r *scrapOrders) FindScrapOrderById(id uint) (orders.ScrapOrders, error) {
	var res orders.ScrapOrders
	err := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&orders.ScrapOrders{}).Preload("Order_lines.Product").Preload(clause.Associations).Where("id", id).First(&res)
	return res, err.Error
}

func (r *scrapOrders) DeleteScrapOrder(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&orders.ScrapOrders{}).Where("id", id).Updates(data)
	return res.Error
}

func (r *scrapOrders) CreateScrapOrderLines(data orders.ScrapOrderLines) error {
	res := r.db.Model(&orders.ScrapOrderLines{}).Create(&data)
	return res.Error
}

func (r *scrapOrders) UpdateScrapOrderLines(query interface{}, data orders.ScrapOrderLines) (int64, error) {
	res := r.db.Model(&orders.ScrapOrderLines{}).Where(query).Updates(&data)
	return res.RowsAffected, res.Error
}

func (r *scrapOrders) FindScrapOrderLines(query interface{}) (orders.ScrapOrderLines, error) {
	var result orders.ScrapOrderLines
	fmt.Println(query)
	res := r.db.Model(&orders.ScrapOrderLines{}).Where(query).First(&result)
	return result, res.Error
}

func (r *scrapOrders) DeleteScrapOrderLines(query interface{}) error {
	res := r.db.Model(&orders.ScrapOrderLines{}).Where(query).Delete(&orders.ScrapOrderLines{})
	return res.Error
}
