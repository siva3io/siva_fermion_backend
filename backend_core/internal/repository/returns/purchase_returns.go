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
type PurchaseReturn interface {
	Save(data *returns.PurchaseReturns) error
	FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.PurchaseReturns, error)
	FindOne(query map[string]interface{}) (returns.PurchaseReturns, error)
	Update(query map[string]interface{}, data *returns.PurchaseReturns) error
	Delete(query map[string]interface{}) error

	SaveReturnLines(returns.PurchaseReturnLines) error
	UpdateReturnLines(map[string]interface{}, returns.PurchaseReturnLines) (int64, error)
	DeleteReturnLine(map[string]interface{}) error
	FindReturnLines(map[string]interface{}) (returns.PurchaseReturnLines, error)

	Search(query string) (interface{}, error)
	GetPurchaseReturnsHistory(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.PurchaseReturns, error)
}

type PurchaseReturns struct {
	db *gorm.DB
}

var PurchaseReturnsRepository *PurchaseReturns //singleton object

// singleton function
func NewPurchaseReturn() *PurchaseReturns {
	if PurchaseReturnsRepository != nil {
		return PurchaseReturnsRepository
	}
	db := db.DbManager()
	PurchaseReturnsRepository = &PurchaseReturns{db}
	return PurchaseReturnsRepository

}

func (r *PurchaseReturns) Save(data *returns.PurchaseReturns) error {
	err := r.db.Model(&returns.PurchaseReturns{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PurchaseReturns) FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.PurchaseReturns, error) {

	var data []returns.PurchaseReturns
	err := r.db.Model(&returns.PurchaseReturns{}).Preload(clause.Associations).Scopes(helpers.Paginate(&returns.PurchaseReturns{}, page, r.db)).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *PurchaseReturns) FindOne(query map[string]interface{}) (returns.PurchaseReturns, error) {
	var data returns.PurchaseReturns

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}

func (r *PurchaseReturns) Update(query map[string]interface{}, data *returns.PurchaseReturns) error {
	res := r.db.Model(&returns.PurchaseReturns{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseReturns) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(uint),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&returns.PurchaseReturns{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseReturns) SaveReturnLines(data returns.PurchaseReturnLines) error {

	res := r.db.Model(&returns.PurchaseReturnLines{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseReturns) FindReturnLines(query map[string]interface{}) (returns.PurchaseReturnLines, error) {
	var result returns.PurchaseReturnLines
	fmt.Println(query)
	res := r.db.Model(&returns.PurchaseReturnLines{}).Where(query).First(&result)

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *PurchaseReturns) UpdateReturnLines(query map[string]interface{}, data returns.PurchaseReturnLines) (int64, error) {
	res := r.db.Model(&returns.PurchaseReturnLines{}).Where(query).Updates(&data)

	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}

func (r *PurchaseReturns) DeleteReturnLine(query map[string]interface{}) error {
	res := r.db.Model(&returns.PurchaseReturnLines{}).Where(query).Delete(&returns.PurchaseReturnLines{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *PurchaseReturns) Search(query string) (interface{}, error) {
	var data []returns.PurchaseReturns

	fields := []string{"reference_number", "number"}

	fields_string, values := helpers.ApplySearch(query, fields)

	err := r.db.Model(&returns.PurchaseReturns{}).Limit(20).Preload(clause.Associations).Where(fields_string, values...).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *PurchaseReturns) GetPurchaseReturnsHistory(query map[string]interface{}, page *pagination.Paginatevalue) ([]returns.PurchaseReturns, error) {
	var data []returns.PurchaseReturns
	var ids = make([]uint, 0)

	productId := query["product_id"]

	page.Filters = fmt.Sprintf("[[\"product_id\", \"=\", %v]]", productId)
	err := r.db.Model(&returns.PurchaseReturns{}).Select("pr_id").Scopes(helpers.Paginate(&returns.PurchaseReturns{}, page, r.db)).Scan(&ids)

	if err.Error != nil {
		return nil, err.Error
	}
	err = r.db.Model(&returns.PurchaseReturns{}).Where("id IN ?", ids).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}
