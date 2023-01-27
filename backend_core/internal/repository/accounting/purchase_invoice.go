package accounting

import (
	"errors"
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
	Save(data *accounting.PurchaseInvoice) error
	Update(query map[string]interface{}, data *accounting.PurchaseInvoice) error
	Delete(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (accounting.PurchaseInvoice, error)
	FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]accounting.PurchaseInvoice, error)

	SaveInvoiceLines(data accounting.PurchaseInvoiceLines) error
	UpdateInvoiceLines(query map[string]interface{}, data accounting.PurchaseInvoiceLines) (int64, error)
	DeleteInvoiceLine(query map[string]interface{}) error
	FindInvoiceLines(query map[string]interface{}) (accounting.PurchaseInvoiceLines, error)
}

type PurchaseInvoice struct {
	db *gorm.DB
}

var PurchaseInvoiceRepository *PurchaseInvoice //singleton object

// singleton function
func NewPurchaseInvoice() *PurchaseInvoice {
	if PurchaseInvoiceRepository != nil {
		return PurchaseInvoiceRepository
	}
	db := db.DbManager()
	PurchaseInvoiceRepository = &PurchaseInvoice{db}
	return PurchaseInvoiceRepository

}

func (r *PurchaseInvoice) Save(data *accounting.PurchaseInvoice) error {
	err := r.db.Model(&accounting.PurchaseInvoice{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PurchaseInvoice) FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]accounting.PurchaseInvoice, error) {
	var data []accounting.PurchaseInvoice
	err := r.db.Model(&accounting.PurchaseInvoice{}).Scopes(helpers.Paginate(&accounting.PurchaseInvoice{}, page, r.db)).Preload(clause.Associations).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *PurchaseInvoice) FindOne(query map[string]interface{}) (accounting.PurchaseInvoice, error) {
	var data accounting.PurchaseInvoice
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops!record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *PurchaseInvoice) Update(query map[string]interface{}, data *accounting.PurchaseInvoice) error {
	res := r.db.Model(&accounting.PurchaseInvoice{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *PurchaseInvoice) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&accounting.PurchaseInvoice{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
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
	res := r.db.Model(&accounting.PurchaseInvoiceLines{}).Where(query).First(&result)
	if res.RowsAffected == 0 {
		return result, errors.New("oops! record not found")
	}
	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *PurchaseInvoice) UpdateInvoiceLines(query map[string]interface{}, data accounting.PurchaseInvoiceLines) (int64, error) {
	res := r.db.Model(&accounting.PurchaseInvoiceLines{}).Where(query).Updates(&data)
	if res.RowsAffected == 0 {
		return 1, errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}
	return res.RowsAffected, nil
}

func (r *PurchaseInvoice) DeleteInvoiceLine(query map[string]interface{}) error {
	res := r.db.Model(&accounting.PurchaseInvoiceLines{}).Where(query).Delete(&accounting.PurchaseInvoiceLines{})
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}

	return nil
}
