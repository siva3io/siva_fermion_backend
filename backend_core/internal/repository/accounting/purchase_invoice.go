package accounting

import (
	"errors"
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/accounting"
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
type Purchase interface {
	Save(data *accounting.PurchaseInvoice) error
	FindAll(page *pagination.Paginatevalue) (interface{}, error)
	FindOne(query map[string]interface{}) (interface{}, error)
	Update(query map[string]interface{}, data *accounting.PurchaseInvoice) error
	Delete(query map[string]interface{}) error

	SaveInvoiceLines(accounting.PurchaseInvoiceLines) error
	UpdateInvoiceLines(map[string]interface{}, accounting.PurchaseInvoiceLines) (int64, error)
	DeleteInvoiceLine(map[string]interface{}) error
	FindInvoiceLines(map[string]interface{}) (accounting.PurchaseInvoiceLines, error)

	Search(query string) (interface{}, error)
}

type PurchaseInvoice struct {
	db *gorm.DB
}

func NewPurchaseInvoice() *PurchaseInvoice {
	db := db.DbManager()
	return &PurchaseInvoice{db}

}

func (r *PurchaseInvoice) Save(data *accounting.PurchaseInvoice) error {
	err := r.db.Model(&accounting.PurchaseInvoice{}).Create(data).Error

	if err != nil {

		return err

	}

	return nil
}

func (r *PurchaseInvoice) FindAll(page *pagination.Paginatevalue) (interface{}, error) {
	var data []accounting.PurchaseInvoice

	err := r.db.Model(&accounting.PurchaseInvoice{}).Scopes(helpers.Paginate(&accounting.PurchaseInvoice{}, page, r.db)).Preload(clause.Associations).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *PurchaseInvoice) FindOne(query map[string]interface{}) (interface{}, error) {
	var data accounting.PurchaseInvoice

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *PurchaseInvoice) Update(query map[string]interface{}, data *accounting.PurchaseInvoice) error {
	res := r.db.Model(&accounting.PurchaseInvoice{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseInvoice) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&accounting.PurchaseInvoice{}).Where(query).Delete(data)
	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseInvoice) SaveInvoiceLines(data accounting.PurchaseInvoiceLines) error {

	res := r.db.Model(&accounting.PurchaseInvoiceLines{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *PurchaseInvoice) FindInvoiceLines(query map[string]interface{}) (accounting.PurchaseInvoiceLines, error) {
	var result accounting.PurchaseInvoiceLines
	fmt.Println(query)
	res := r.db.Model(&accounting.PurchaseInvoiceLines{}).Where(query).First(&result)

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *PurchaseInvoice) UpdateInvoiceLines(query map[string]interface{}, data accounting.PurchaseInvoiceLines) (int64, error) {
	res := r.db.Model(&accounting.PurchaseInvoiceLines{}).Where(query).Updates(&data)

	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}

func (r *PurchaseInvoice) DeleteInvoiceLine(query map[string]interface{}) error {
	res := r.db.Model(&accounting.PurchaseInvoiceLines{}).Where(query).Delete(&accounting.PurchaseInvoiceLines{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *PurchaseInvoice) Search(query string) (interface{}, error) {
	var data []accounting.PurchaseInvoice

	fields := []string{"reference_number", "number"}

	fields_string, values := helpers.ApplySearch(query, fields)

	err := r.db.Model(&accounting.PurchaseInvoice{}).Limit(20).Preload(clause.Associations).Where(fields_string, values...).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}
