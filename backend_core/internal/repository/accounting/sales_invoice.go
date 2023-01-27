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
	CreateSalesInvoice(data *accounting.SalesInvoice) (uint, error)
	BulkCreateSalesInvoice(data *[]accounting.SalesInvoice) error
	UpdateSalesInvoice(query map[string]interface{}, data *accounting.SalesInvoice) error
	GetSalesInvoice(query map[string]interface{}) (accounting.SalesInvoice, error)
	GetAllSalesInvoice(query map[string]interface{}, p *pagination.Paginatevalue) ([]accounting.SalesInvoice, error)
	DeleteSalesInvoice(query map[string]interface{}) error

	CreateSalesInvoiceLines(data accounting.SalesInvoiceLines) error
	UpdateSalesInvoiceLines(query map[string]interface{}, data accounting.SalesInvoiceLines) (int64, error)
	GetSalesInvoiceLines(query map[string]interface{}) ([]accounting.SalesInvoiceLines, error)
	DeleteSalesInvoiceLines(query map[string]interface{}) error
}

type salesInvoice struct {
	db *gorm.DB
}

var salesInvoiceRepository *salesInvoice //singleton object

// singleton function
func NewSalesInvoice() *salesInvoice {
	if salesInvoiceRepository != nil {
		return salesInvoiceRepository
	}
	db := db.DbManager()
	salesInvoiceRepository = &salesInvoice{db}
	return salesInvoiceRepository
}

func (r *salesInvoice) CreateSalesInvoice(data *accounting.SalesInvoice) (uint, error) {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'SALES_INVOICE_STATUS' AND lookupcodes.lookup_code = 'NEW'").First(&scode).Error
	if err != nil {
		return 0, err
	}
	data.StatusID = &scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, *data.StatusID)
	data.StatusHistory = res
	err = r.db.Model(&accounting.SalesInvoice{}).Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (r *salesInvoice) BulkCreateSalesInvoice(data *[]accounting.SalesInvoice) error {
	for _, value := range *data {
		var scode uint
		err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'SALES_INVOICE_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
		if err != nil {
			return err
		}
		value.StatusID = &scode
		res, _ := helpers.UpdateStatusHistory(value.StatusHistory, *value.StatusID)
		value.StatusHistory = res
		result := r.db.Model(&accounting.SalesInvoice{}).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *salesInvoice) UpdateSalesInvoice(query map[string]interface{}, data *accounting.SalesInvoice) error {
	err := r.db.Model(&accounting.SalesInvoice{}).Where(query).Updates(&data)
	if err.RowsAffected == 0 {

		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *salesInvoice) GetSalesInvoice(query map[string]interface{}) (accounting.SalesInvoice, error) {
	var data accounting.SalesInvoice
	fmt.Println("-----------------------------response-----------", query)
	// fmt.Println("======",query)
	err := r.db.Model(&accounting.SalesInvoice{}).Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations).Where(query).First(&data)
	// fmt.Println("error",err)
	// helpers.PrettyPrint("data",data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *salesInvoice) GetAllSalesInvoice(query map[string]interface{}, p *pagination.Paginatevalue) ([]accounting.SalesInvoice, error) {
	var data []accounting.SalesInvoice
	err := r.db.Model(&accounting.SalesInvoice{}).Scopes(helpers.Paginate(&accounting.SalesInvoice{}, p, r.db)).Where(query).Preload("SalesInvoiceLines.Product").Preload("SalesInvoiceLines.ProductVariant").Preload("SalesInvoiceLines.Warehouse").Preload("SalesInvoiceLines.Inventory").Preload("SalesInvoiceLines.UOM").Preload("SalesInvoiceLines.DiscountType").Preload("SalesInvoiceLines.TaxType").Preload("SalesInvoiceLines.PaymentTerms").Preload(clause.Associations).Find(&data)
	if err != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *salesInvoice) DeleteSalesInvoice(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&accounting.SalesInvoice{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *salesInvoice) CreateSalesInvoiceLines(data accounting.SalesInvoiceLines) error {
	result := r.db.Model(&accounting.SalesInvoiceLines{}).Create(&data)
	if result != nil {
		return result.Error
	}
	return nil
}

func (r *salesInvoice) UpdateSalesInvoiceLines(query map[string]interface{}, data accounting.SalesInvoiceLines) (int64, error) {
	result := r.db.Model(&accounting.SalesInvoiceLines{}).Where(query).Updates(&data)
	if result.Error != nil {
		return result.RowsAffected, result.Error
	}
	return result.RowsAffected, nil
}

func (r *salesInvoice) GetSalesInvoiceLines(query map[string]interface{}) ([]accounting.SalesInvoiceLines, error) {
	var data []accounting.SalesInvoiceLines
	result := r.db.Model(&accounting.SalesInvoiceLines{}).Where(query).Preload(clause.Associations).Find(&data)
	if result.RowsAffected == 0 {
		return data, errors.New("oops!record not found")
	}
	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}

func (r *salesInvoice) DeleteSalesInvoiceLines(query map[string]interface{}) error {
	res := r.db.Model(&accounting.SalesInvoiceLines{}).Where(query).Delete(&accounting.SalesInvoiceLines{})
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}

	return nil
}
