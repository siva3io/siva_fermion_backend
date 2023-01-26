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
type WD interface {
	UpdateWD(id uint, data shipping.WD) error
	DeleteWD(id uint, user_id uint) error
	CreateWD(data *shipping.WD) (uint, error)
	BulkCreateWD(data *[]shipping.WD) error
	GetAllWD(p *pagination.Paginatevalue) ([]shipping.WD, error)
	GetWD(id uint) (shipping.WD, error)
}

type wd struct {
	db *gorm.DB
}

func NewWD() *wd {
	db := db.DbManager()
	return &wd{db}
}

func (r *wd) CreateWD(data *shipping.WD) (uint, error) {
	res := r.db.Model(&shipping.WD{}).Create(&data)
	return data.ID, res.Error
}
func (r *wd) BulkCreateWD(data *[]shipping.WD) error {
	res := r.db.Model(&shipping.WD{}).Create(&data)
	return res.Error
}

func (r *wd) GetWD(id uint) (shipping.WD, error) {
	var data shipping.WD
	result := r.db.Model(&shipping.WD{}).Where("id", id).Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).First(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

func (r *wd) GetAllWD(p *pagination.Paginatevalue) ([]shipping.WD, error) {
	var data []shipping.WD
	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Model(&shipping.WD{}).Scopes(helpers.Paginate(&shipping.WD{}, p, r.db)).Where("is_active = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *wd) UpdateWD(id uint, data shipping.WD) error {
	res := r.db.Model(&shipping.WD{}).Where("id", id).Updates(&data)
	return res.Error
}

func (r *wd) DeleteWD(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&shipping.WD{}).Where("id", id).Updates(data)
	return res.Error
}
