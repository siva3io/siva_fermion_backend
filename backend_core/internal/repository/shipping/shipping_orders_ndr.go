package shipping

import (
	"errors"
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
type NDR interface {
	CreateNDR(data *shipping.NDR) error
	BulkCreateNDR(data *[]shipping.NDR) error
	GetAllNDR(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.NDR, error)
	GetNDR(query map[string]interface{}) (shipping.NDR, error)
	UpdateNDR(query map[string]interface{}, data *shipping.NDR) error
	DeleteNDR(query map[string]interface{}) error

	CreateNDRLines(data shipping.NDRLines) error
	GetNDRLines(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.NDRLines, error)
	UpdateNDRLines(query map[string]interface{}, data *shipping.NDRLines) (int64, error)
	DeleteNDRLines(query map[string]interface{}) error
}

type ndr struct {
	db *gorm.DB
}

var ndrRepository *ndr //singleton object

// singleton function
func NewNDR() *ndr {
	if ndrRepository != nil {
		return ndrRepository
	}
	db := db.DbManager()
	ndrRepository = &ndr{db}
	return ndrRepository
}

func (r *ndr) CreateNDR(data *shipping.NDR) error {
	err := r.db.Model(&shipping.NDR{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ndr) BulkCreateNDR(data *[]shipping.NDR) error {
	res := r.db.Model(&shipping.NDR{}).Create(&data)
	return res.Error
}

func (r *ndr) GetNDR(query map[string]interface{}) (shipping.NDR, error) {
	var data shipping.NDR
	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Model(&shipping.NDR{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *ndr) GetAllNDR(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.NDR, error) {
	var data []shipping.NDR

	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations + "." + clause.Associations).Model(&shipping.NDR{}).Scopes(helpers.Paginate(&shipping.NDR{}, p, r.db)).Where("is_active = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *ndr) UpdateNDR(query map[string]interface{}, data *shipping.NDR) error {
	err := r.db.Model(&shipping.NDR{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *ndr) UpdateNDRLines(query map[string]interface{}, data *shipping.NDRLines) (int64, error) {
	result := r.db.Model(&shipping.NDRLines{}).Where(query).Updates(&data)
	return result.RowsAffected, result.Error
}

func (r *ndr) CreateNDRLines(data shipping.NDRLines) error {
	result := r.db.Model(&shipping.NDRLines{}).Create(&data)
	return result.Error
}

func (r *ndr) DeleteNDR(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&shipping.NDR{}).Where(query).Updates(data)
	return res.Error
}

func (r *ndr) GetNDRLines(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.NDRLines, error) {
	var data []shipping.NDRLines
	result := r.db.Model(&shipping.NDRLines{}).Where(query).Find(&data)
	return data, result.Error
}

func (r *ndr) DeleteNDRLines(query map[string]interface{}) error {
	var data shipping.NDRLines
	result := r.db.Model(&shipping.NDRLines{}).Where(query).Delete(&data)
	return result.Error
}
