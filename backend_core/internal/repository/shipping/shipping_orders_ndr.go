package shipping

import (
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
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
type NDR interface {
	CreateNDR(data *shipping.NDR) (uint, error)
	BulkCreateNDR(data *[]shipping.NDR) error
	GetAllNDR(p *pagination.Paginatevalue) ([]shipping.NDR, error)
	GetNDR(id uint) (shipping.NDR, error)
	UpdateNDR(id uint, data shipping.NDR) error
	DeleteNDR(id uint, user_id uint) error

	CreateNDRLines(data shipping.NDRLines) error
	GetNDRLines(query interface{}) ([]shipping.NDRLines, error)
	UpdateNDRLines(query interface{}, data shipping.NDRLines) (int64, error)
	DeleteNDRLines(query interface{}) error
}

type ndr struct {
	db *gorm.DB
}

func NewNDR() *ndr {
	db := db.DbManager()
	return &ndr{db}
}

func (r *ndr) CreateNDR(data *shipping.NDR) (uint, error) {
	res := r.db.Model(&shipping.NDR{}).Create(&data)
	return data.ID, res.Error
}

func (r *ndr) BulkCreateNDR(data *[]shipping.NDR) error {
	res := r.db.Model(&shipping.NDR{}).Create(&data)
	return res.Error
}

func (r *ndr) GetNDR(id uint) (shipping.NDR, error) {
	var data shipping.NDR
	result := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Model(&shipping.NDR{}).Where("id", id).First(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

func (r *ndr) GetAllNDR(p *pagination.Paginatevalue) ([]shipping.NDR, error) {
	var data []shipping.NDR

	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Model(&shipping.NDR{}).Scopes(helpers.Paginate(&shipping.NDR{}, p, r.db)).Where("is_active = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *ndr) UpdateNDR(id uint, data shipping.NDR) error {
	res := r.db.Model(&shipping.NDR{}).Where("id", id).Updates(&data)
	return res.Error
}

func (r *ndr) UpdateNDRLines(query interface{}, data shipping.NDRLines) (int64, error) {
	result := r.db.Model(&shipping.NDRLines{}).Where(query).Updates(&data)
	return result.RowsAffected, result.Error
}

func (r *ndr) CreateNDRLines(data shipping.NDRLines) error {
	result := r.db.Model(&shipping.NDRLines{}).Create(&data)
	return result.Error
}

func (r *ndr) DeleteNDR(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&shipping.NDR{}).Where("id", id).Updates(data)
	return res.Error
}

func (r *ndr) GetNDRLines(query interface{}) ([]shipping.NDRLines, error) {
	var data []shipping.NDRLines
	result := r.db.Model(&shipping.NDRLines{}).Where(query).Find(&data)
	return data, result.Error
}

func (r *ndr) DeleteNDRLines(query interface{}) error {
	var data shipping.NDRLines
	result := r.db.Model(&shipping.NDRLines{}).Where(query).Delete(&data)
	return result.Error
}
