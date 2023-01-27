package inventory_orders

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/inventory_orders"
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
type Asn interface {
	CreateAsn(data *inventory_orders.ASN) error
	BulkCreateAsn(data *[]inventory_orders.ASN) error
	UpdateAsn(query map[string]interface{}, data *inventory_orders.ASN) error
	GetAsn(query map[string]interface{}) (inventory_orders.ASN, error)
	GetAllAsn(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_orders.ASN, error)
	DeleteAsn(query map[string]interface{}) error
	CreateAsnLines(data inventory_orders.AsnLines) error
	UpdateAsnLines(query interface{}, data inventory_orders.AsnLines) (int64, error)
	GetAsnLines(query map[string]interface{}) (inventory_orders.AsnLines, error)
	DeleteAsnLines(query interface{}) error
}

type asn struct {
	db *gorm.DB
}

var asnRepository *asn //singleton object

// singleton function
func NewAsn() *asn {
	if asnRepository != nil {
		return asnRepository
	}
	db := db.DbManager()
	asnRepository = &asn{db}
	return asnRepository
}

func (r *asn) CreateAsn(data *inventory_orders.ASN) error {

	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'ASN_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	result := r.db.Model(&inventory_orders.ASN{}).Create(data).Error
	if result != nil {
		return err
	}
	return nil
}

func (r *asn) BulkCreateAsn(data *[]inventory_orders.ASN) error {

	for _, value := range *data {
		var scode uint
		err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'ASN_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
		if err != nil {
			return err
		}
		value.StatusID = scode
		res, _ := helpers.UpdateStatusHistory(value.StatusHistory, value.StatusID)
		value.StatusHistory = res
		result := r.db.Model(&inventory_orders.ASN{}).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *asn) UpdateAsn(query map[string]interface{}, data *inventory_orders.ASN) error {
	result := r.db.Model(&inventory_orders.ASN{}).Where(query).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if result != nil {
		return result.Error
	}
	return nil
}

func (r *asn) GetAsn(query map[string]interface{}) (inventory_orders.ASN, error) {
	var data inventory_orders.ASN
	result := r.db.Model(&inventory_orders.ASN{}).Preload(clause.Associations).Where(query).First(&data)
	if result.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if result != nil {
		return data, result.Error
	}
	return data, nil
}

func (r *asn) GetAllAsn(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_orders.ASN, error) {
	var data []inventory_orders.ASN
	res := r.db.Model(&inventory_orders.ASN{}).Scopes(helpers.Paginate(&inventory_orders.ASN{}, p, r.db)).Where(query).Preload("AsnOrderLines.Product").Preload("AsnOrderLines.ProductVariant").Preload("AsnOrderLines.PackageType").Preload("AsnOrderLines.UOM").Preload(clause.Associations).Find(&data)
	if res != nil {
		return data, res.Error
	}
	return data, nil

}

func (r *asn) DeleteAsn(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	result := r.db.Model(&inventory_orders.ASN{}).Where(query).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if result != nil {
		return result.Error
	}
	return nil
}

func (r *asn) CreateAsnLines(data inventory_orders.AsnLines) error {
	result := r.db.Model(&inventory_orders.AsnLines{}).Create(&data)
	if result != nil {
		return result.Error
	}
	return nil
}

func (r *asn) UpdateAsnLines(query interface{}, data inventory_orders.AsnLines) (int64, error) {
	result := r.db.Model(&inventory_orders.AsnLines{}).Where(query).Updates(data)
	if result != nil {
		return result.RowsAffected, result.Error
	}
	return result.RowsAffected, nil
}

func (r *asn) GetAsnLines(query map[string]interface{}) (inventory_orders.AsnLines, error) {
	var data inventory_orders.AsnLines
	result := r.db.Model(&inventory_orders.AsnLines{}).Where(query).Preload(clause.Associations).Find(&data)
	if result.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if result != nil {
		return data, result.Error
	}
	return data, nil

}

func (r *asn) DeleteAsnLines(query interface{}) error {
	result := r.db.Model(&inventory_orders.AsnLines{}).Where(query).Delete(&inventory_orders.AsnLines{})
	if result.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if result != nil {
		return result.Error
	}
	return nil
}
