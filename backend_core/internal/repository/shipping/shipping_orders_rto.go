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
type RTO interface {
	CreateRTO(data *shipping.RTO) error
	BulkCreateRTO(data *[]shipping.RTO) error
	GetAllRTO(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.RTO, error)
	GetRTO(query map[string]interface{}) (shipping.RTO, error)
	UpdateRTO(query map[string]interface{}, data *shipping.RTO) error
	DeleteRTO(query map[string]interface{}) error
}

type rto struct {
	db *gorm.DB
}

var rtoRepository *rto //singleton object

// singleton function
func NewRTO() *rto {
	if rtoRepository != nil {
		return rtoRepository
	}
	db := db.DbManager()
	rtoRepository = &rto{db}
	return rtoRepository
}

func (r *rto) CreateRTO(data *shipping.RTO) error {
	err := r.db.Model(&shipping.RTO{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *rto) BulkCreateRTO(data *[]shipping.RTO) error {
	res := r.db.Model(&shipping.RTO{}).Create(&data)
	return res.Error
}

func (r *rto) GetRTO(query map[string]interface{}) (shipping.RTO, error) {
	var data shipping.RTO
	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Model(&shipping.RTO{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *rto) GetAllRTO(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.RTO, error) {
	var data []shipping.RTO
	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations + "." + clause.Associations).Model(&shipping.RTO{}).Scopes(helpers.Paginate(&shipping.RTO{}, p, r.db)).Where("is_active = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *rto) UpdateRTO(query map[string]interface{}, data *shipping.RTO) error {
	err := r.db.Model(&shipping.RTO{}).Where(query).Updates(&data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return err.Error
}

func (r *rto) DeleteRTO(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shipping.RTO{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
