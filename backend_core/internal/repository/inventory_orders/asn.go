package inventory_orders

import (
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
	CreateAsn(data *inventory_orders.ASN, userID string) (uint, error)
	BulkCreateAsn(data *[]inventory_orders.ASN, userID string) error
	UpdateAsn(id uint, data *inventory_orders.ASN) error
	GetAsn(id uint) (inventory_orders.ASN, error)
	GetAllAsn(p *pagination.Paginatevalue) ([]inventory_orders.ASN, error)
	DeleteAsn(id uint, user_id uint) error

	CreateAsnLines(data inventory_orders.AsnLines) error
	UpdateAsnLines(query interface{}, data inventory_orders.AsnLines) (int64, error)
	GetAsnLines(query interface{}) ([]inventory_orders.AsnLines, error)
	DeleteAsnLines(query interface{}) error
}

type asn struct {
	db *gorm.DB
}

func NewAsn() *asn {
	db := db.DbManager()
	return &asn{db}
}

func (r *asn) CreateAsn(data *inventory_orders.ASN, userID string) (uint, error) {

	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'ASN_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return 0, err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	result := r.db.Model(&inventory_orders.ASN{}).Create(&data)
	return data.ID, result.Error
}

func (r *asn) BulkCreateAsn(data *[]inventory_orders.ASN, userID string) error {

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

func (r *asn) UpdateAsn(id uint, data *inventory_orders.ASN) error {
	result := r.db.Model(&inventory_orders.ASN{}).Where("id", id).Updates(&data)
	return result.Error
}

func (r *asn) GetAsn(id uint) (inventory_orders.ASN, error) {
	var data inventory_orders.ASN
	result := r.db.Model(&inventory_orders.ASN{}).Preload(clause.Associations).Where("id", id).First(&data)
	return data, result.Error
}

func (r *asn) GetAllAsn(p *pagination.Paginatevalue) ([]inventory_orders.ASN, error) {
	var data []inventory_orders.ASN
	res := r.db.Model(&inventory_orders.ASN{}).Scopes(helpers.Paginate(&inventory_orders.ASN{}, p, r.db)).Where("is_active = true").Preload("AsnOrderLines.Product").Preload("AsnOrderLines.ProductVariant").Preload("AsnOrderLines.PackageType").Preload("AsnOrderLines.UOM").Preload(clause.Associations).Find(&data)
	return data, res.Error
}

func (r *asn) DeleteAsn(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	result := r.db.Model(&inventory_orders.ASN{}).Where("id", id).Updates(data)
	return result.Error
}

func (r *asn) CreateAsnLines(data inventory_orders.AsnLines) error {
	result := r.db.Model(&inventory_orders.AsnLines{}).Create(&data)
	return result.Error
}

func (r *asn) UpdateAsnLines(query interface{}, data inventory_orders.AsnLines) (int64, error) {
	result := r.db.Model(&inventory_orders.AsnLines{}).Where(query).Updates(&data)
	return result.RowsAffected, result.Error
}

func (r *asn) GetAsnLines(query interface{}) ([]inventory_orders.AsnLines, error) {
	var data []inventory_orders.AsnLines
	result := r.db.Model(&inventory_orders.AsnLines{}).Where(query).Preload(clause.Associations).Find(&data)

	return data, result.Error
}

func (r *asn) DeleteAsnLines(query interface{}) error {
	result := r.db.Model(&inventory_orders.AsnLines{}).Where(query).Delete(&inventory_orders.AsnLines{})
	return result.Error
}
