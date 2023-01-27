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
type WD interface {
	UpdateWD(query map[string]interface{}, data *shipping.WD) error
	DeleteWD(query map[string]interface{}) error
	CreateWD(data *shipping.WD) error
	BulkCreateWD(data *[]shipping.WD) error
	GetAllWD(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.WD, error)
	GetWD(query map[string]interface{}) (shipping.WD, error)
}

type wd struct {
	db *gorm.DB
}

var wdRepository *wd //singleton object

// singleton function
func NewWD() *wd {
	if wdRepository != nil {
		return wdRepository
	}
	db := db.DbManager()
	wdRepository = &wd{db}
	return wdRepository
}

func (r *wd) CreateWD(data *shipping.WD) error {
	err := r.db.Model(&shipping.WD{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *wd) BulkCreateWD(data *[]shipping.WD) error {
	res := r.db.Model(&shipping.WD{}).Create(&data)
	return res.Error
}

func (r *wd) GetWD(query map[string]interface{}) (shipping.WD, error) {

	var data shipping.WD
	err := r.db.Model(&shipping.WD{}).Where(query).Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *wd) GetAllWD(query map[string]interface{}, p *pagination.Paginatevalue) ([]shipping.WD, error) {
	var data []shipping.WD
	err := r.db.Preload("ShippingOrder.Channel").Preload("ShippingOrder.Partner").Preload("ShippingOrder.Order").Preload("ShippingOrder.ShippingPartner").Preload("ShippingOrder.ShippingOrderLines").Preload("ShippingOrder.ShippingOrderLines.ProductVariant").Preload("ShippingOrder.ShippingOrderLines.ProductTemplate").Preload(clause.Associations + "." + clause.Associations).Model(&shipping.WD{}).Scopes(helpers.Paginate(&shipping.WD{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *wd) UpdateWD(query map[string]interface{}, data *shipping.WD) error {
	err := r.db.Model(&shipping.WD{}).Where(query).Updates(&data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return err.Error
}

func (r *wd) DeleteWD(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shipping.WD{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
