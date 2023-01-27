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
type ScrapOrders interface {
	SaveScrapOrder(data *orders.ScrapOrders) (uint, error)
	BulkSaveScrapOrder(data *[]orders.ScrapOrders) error
	UpdateScrapOrder(query map[string]interface{}, data *orders.ScrapOrders) error
	FindAllScrapOrders(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.ScrapOrders, error)
	FindScrapOrderById(query map[string]interface{}) (orders.ScrapOrders, error)
	DeleteScrapOrder(query map[string]interface{}) error

	CreateScrapOrderLines(data orders.ScrapOrderLines) error
	UpdateScrapOrderLines(query interface{}, data orders.ScrapOrderLines) (int64, error)
	DeleteScrapOrderLines(query interface{}) error
	FindScrapOrderLines(query interface{}) (orders.ScrapOrderLines, error)
}

type scrapOrders struct {
	db *gorm.DB
}

var scrapOrdersRepository *scrapOrders //singleton object

// singleton function
func NewScrap() *scrapOrders {
	if scrapOrdersRepository != nil {
		return scrapOrdersRepository
	}
	db := db.DbManager()
	scrapOrdersRepository = &scrapOrders{db}
	return scrapOrdersRepository
}

func (r *scrapOrders) SaveScrapOrder(data *orders.ScrapOrders) (uint, error) {

	res := r.db.Model(&orders.ScrapOrders{}).Create(&data)
	return data.ID, res.Error
}
func (r *scrapOrders) BulkSaveScrapOrder(data *[]orders.ScrapOrders) error {
	fmt.Println("repo")
	res := r.db.Model(&orders.ScrapOrders{}).Create(data)
	return res.Error
}

func (r *scrapOrders) UpdateScrapOrder(query map[string]interface{}, data *orders.ScrapOrders) error {
	err := r.db.Model(&orders.ScrapOrders{}).Where(query).Updates(data)

	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {

		return err.Error
	}
	return nil
}
func (r *scrapOrders) FindAllScrapOrders(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.ScrapOrders, error) {
	var data []orders.ScrapOrders

	err := r.db.Model(&orders.ScrapOrders{}).Preload(clause.Associations).Scopes(helpers.Paginate(&orders.ScrapOrders{}, page, r.db)).Where(query).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *scrapOrders) FindScrapOrderById(query map[string]interface{}) (orders.ScrapOrders, error) {
	var data orders.ScrapOrders

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}

func (r *scrapOrders) DeleteScrapOrder(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(uint),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&orders.ScrapOrders{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {

		return err.Error
	}
	return nil
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
