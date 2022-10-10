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
type RTO interface {
	CreateRTO(data *shipping.RTO) (uint, error)
	BulkCreateRTO(data *[]shipping.RTO) error
	GetAllRTO(p *pagination.Paginatevalue) ([]shipping.RTO, error)
	GetRTO(id uint) (shipping.RTO, error)
	UpdateRTO(id uint, data shipping.RTO) error
	DeleteRTO(id uint, user_id uint) error
}

type rto struct {
	db *gorm.DB
}

func NewRTO() *rto {
	db := db.DbManager()
	return &rto{db}
}

func (r *rto) CreateRTO(data *shipping.RTO) (uint, error) {
	res := r.db.Model(&shipping.RTO{}).Create(&data)
	return data.ID, res.Error
}

func (r *rto) BulkCreateRTO(data *[]shipping.RTO) error {
	res := r.db.Model(&shipping.RTO{}).Create(&data)
	return res.Error
}

func (r *rto) GetRTO(id uint) (shipping.RTO, error) {
	var data shipping.RTO
	result := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Model(&shipping.RTO{}).Where("id", id).First(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

func (r *rto) GetAllRTO(p *pagination.Paginatevalue) ([]shipping.RTO, error) {
	var data []shipping.RTO
	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Model(&shipping.RTO{}).Scopes(helpers.Paginate(&shipping.RTO{}, p, r.db)).Where("is_active = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *rto) UpdateRTO(id uint, data shipping.RTO) error {
	res := r.db.Model(&shipping.RTO{}).Where("id", id).Updates(&data)
	return res.Error
}

func (r *rto) DeleteRTO(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&shipping.RTO{}).Where("id", id).Updates(data)
	return res.Error
}
