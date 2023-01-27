package returns

import (
	"errors"
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
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
type SalesReturn interface {
	Save(data *returns.SalesReturns) error
	FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.SalesReturns, error)
	FindOne(query map[string]interface{}) (returns.SalesReturns, error)
	Update(query map[string]interface{}, data *returns.SalesReturns) error
	Delete(query map[string]interface{}) error

	SaveReturnLines(returns.SalesReturnLines) error
	UpdateReturnLines(map[string]interface{}, returns.SalesReturnLines) (int64, error)
	DeleteReturnLine(map[string]interface{}) error
	FindReturnLines(map[string]interface{}) (returns.SalesReturnLines, error)

	Search(query string) (interface{}, error)
	GetSalesReturnsHistory(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.SalesReturns, error)
}

type SalesReturns struct {
	db *gorm.DB
}

var SalesReturnsRepository *SalesReturns //singleton object

// singleton function
func NewSalesReturn() *SalesReturns {
	if SalesReturnsRepository != nil {
		return SalesReturnsRepository
	}
	db := db.DbManager()
	SalesReturnsRepository = &SalesReturns{db}
	return SalesReturnsRepository

}

func (r *SalesReturns) Save(data *returns.SalesReturns) error {
	err := r.db.Model(&returns.SalesReturns{}).Create(data).Error

	if err != nil {

		return err

	}

	return nil
}

func (r *SalesReturns) FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.SalesReturns, error) {
	var data []returns.SalesReturns
	err := r.db.Model(&returns.SalesReturns{}).Preload(clause.Associations).Scopes(helpers.Paginate(&returns.SalesReturns{}, page, r.db)).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *SalesReturns) FindOne(query map[string]interface{}) (returns.SalesReturns, error) {
	var data returns.SalesReturns

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}

func (r *SalesReturns) Update(query map[string]interface{}, data *returns.SalesReturns) error {
	res := r.db.Model(&returns.SalesReturns{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *SalesReturns) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(uint),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&returns.SalesReturns{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *SalesReturns) SaveReturnLines(data returns.SalesReturnLines) error {

	res := r.db.Model(&returns.SalesReturnLines{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *SalesReturns) FindReturnLines(query map[string]interface{}) (returns.SalesReturnLines, error) {
	var result returns.SalesReturnLines
	fmt.Println(query)
	res := r.db.Model(&returns.SalesReturnLines{}).Where(query).First(&result)

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *SalesReturns) UpdateReturnLines(query map[string]interface{}, data returns.SalesReturnLines) (int64, error) {
	res := r.db.Model(&returns.SalesReturnLines{}).Where(query).Updates(&data)

	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}

func (r *SalesReturns) DeleteReturnLine(query map[string]interface{}) error {
	res := r.db.Model(&returns.SalesReturnLines{}).Where(query).Delete(&returns.SalesReturnLines{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *SalesReturns) Search(query string) (interface{}, error) {
	var data []returns.SalesReturns

	fields := []string{"customer_name", "channel_name", "number"}

	fields_string, values := helpers.ApplySearch(query, fields)

	err := r.db.Model(&returns.SalesReturns{}).Limit(20).Preload(clause.Associations).Where(fields_string, values...).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *SalesReturns) GetSalesReturnsHistory(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.SalesReturns, error) {
	var data []returns.SalesReturns
	var ids = make([]uint, 0)

	productId := query["product_id"]

	page.Filters = fmt.Sprintf("[[\"product_id\", \"=\", %v]]", productId)
	err := r.db.Model(&returns.SalesReturns{}).Select("sr_id").Scopes(helpers.Paginate(&returns.SalesReturns{}, page, r.db)).Scan(&ids)

	if err.Error != nil {
		return nil, err.Error
	}
	err = r.db.Model(&returns.SalesReturns{}).Where("id IN ?", ids).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}
