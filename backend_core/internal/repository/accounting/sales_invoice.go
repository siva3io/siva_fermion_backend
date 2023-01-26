package accounting

import (
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
type SalesInvoice interface {
	CreateSalesInvoice(data *accounting.SalesInvoice, userID string) (uint, error)
	BulkCreateSalesInvoice(data *[]accounting.SalesInvoice, userID string) error
	UpdateSalesInvoice(id uint, data *accounting.SalesInvoice) error
	GetSalesInvoice(id uint) (accounting.SalesInvoice, error)
	GetAllSalesInvoice(p *pagination.Paginatevalue) ([]accounting.SalesInvoice, error)
	DeleteSalesInvoice(id uint, user_id uint) error

	CreateSalesInvoiceLines(data accounting.SalesInvoiceLines) error
	UpdateSalesInvoiceLines(query interface{}, data accounting.SalesInvoiceLines) (int64, error)
	GetSalesInvoiceLines(query interface{}) ([]accounting.SalesInvoiceLines, error)
	DeleteSalesInvoiceLines(query interface{}) error
}

type salesInvoice struct {
	db *gorm.DB
}

func NewSalesInvoice() *salesInvoice {
	db := db.DbManager()
	return &salesInvoice{db}
}

func (r *salesInvoice) CreateSalesInvoice(data *accounting.SalesInvoice, userID string) (uint, error) {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'SALES_INVOICE_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return 0, err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	result := r.db.Model(&accounting.SalesInvoice{}).Create(&data)
	return data.ID, result.Error
}

func (r *salesInvoice) BulkCreateSalesInvoice(data *[]accounting.SalesInvoice, userID string) error {
	for _, value := range *data {
		var scode uint
		err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'SALES_INVOICE_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
		if err != nil {
			return err
		}
		value.StatusID = scode
		res, _ := helpers.UpdateStatusHistory(value.StatusHistory, value.StatusID)
		value.StatusHistory = res
		result := r.db.Model(&accounting.SalesInvoice{}).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *salesInvoice) UpdateSalesInvoice(id uint, data *accounting.SalesInvoice) error {
	result := r.db.Model(&accounting.SalesInvoice{}).Where("id", id).Updates(&data)
	return result.Error
}

func (r *salesInvoice) GetSalesInvoice(id uint) (accounting.SalesInvoice, error) {
	var data accounting.SalesInvoice
	result := r.db.Model(&accounting.SalesInvoice{}).Preload(clause.Associations).Where("id", id).First(&data)
	return data, result.Error
}

func (r *salesInvoice) GetAllSalesInvoice(p *pagination.Paginatevalue) ([]accounting.SalesInvoice, error) {
	var data []accounting.SalesInvoice
	res := r.db.Model(&accounting.SalesInvoice{}).Scopes(helpers.Paginate(&accounting.SalesInvoice{}, p, r.db)).Where("is_active = true").Preload("SalesInvoiceLines.Product").Preload("SalesInvoiceLines.ProductVariant").Preload("SalesInvoiceLines.Warehouse").Preload("SalesInvoiceLines.Inventory").Preload("SalesInvoiceLines.UOM").Preload("SalesInvoiceLines.DiscountType").Preload("SalesInvoiceLines.TaxType").Preload("SalesInvoiceLines.PaymentTerms").Preload(clause.Associations).Find(&data)
	return data, res.Error
}

func (r *salesInvoice) DeleteSalesInvoice(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	result := r.db.Model(&accounting.SalesInvoice{}).Where("id", id).Updates(data)
	return result.Error
}

func (r *salesInvoice) CreateSalesInvoiceLines(data accounting.SalesInvoiceLines) error {
	result := r.db.Model(&accounting.SalesInvoiceLines{}).Create(&data)
	return result.Error
}

func (r *salesInvoice) UpdateSalesInvoiceLines(query interface{}, data accounting.SalesInvoiceLines) (int64, error) {
	result := r.db.Model(&accounting.SalesInvoiceLines{}).Where(query).Updates(&data)
	return result.RowsAffected, result.Error
}

func (r *salesInvoice) GetSalesInvoiceLines(query interface{}) ([]accounting.SalesInvoiceLines, error) {
	var data []accounting.SalesInvoiceLines
	result := r.db.Model(&accounting.SalesInvoiceLines{}).Where(query).Preload(clause.Associations).Find(&data)

	return data, result.Error
}

func (r *salesInvoice) DeleteSalesInvoiceLines(query interface{}) error {
	result := r.db.Model(&accounting.SalesInvoiceLines{}).Where(query).Delete(&accounting.SalesInvoiceLines{})
	return result.Error
}
