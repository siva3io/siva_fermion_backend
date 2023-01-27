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
type ShippingOrder interface {
	Create(data *shipping.ShippingOrder) error
	AddProduct(data shipping.ShippingOrderLines) error
	List(query map[string]interface{}, page *pagination.Paginatevalue) ([]shipping.ShippingOrder, error)
	Get(query map[string]interface{}) (shipping.ShippingOrder, error)
	Delete(query map[string]interface{}) error
	Update(query map[string]interface{}, data *shipping.ShippingOrder) error
	UpdateProductList(query interface{}, data shipping.ShippingOrderLines) (int64, error)
	DeleteOrderLine(query map[string]interface{}) error
	CreateBulkOrder(data *[]shipping.ShippingOrder) error
}

type shippingorder struct {
	db *gorm.DB
}

var shippingOrderRepository *shippingorder //singleton object

// singleton function
func NewShippingOrder() *shippingorder {
	if shippingOrderRepository != nil {
		return shippingOrderRepository
	}
	db := db.DbManager()
	shippingOrderRepository = &shippingorder{db}
	return shippingOrderRepository
}

func (r *shippingorder) Create(data *shipping.ShippingOrder) error {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'SHIPPING_ORDER_STATUS' AND lookupcodes.lookup_code = 'PICKUP_SCHEDULED'").First(&scode).Error
	if err != nil {
		return err
	}
	data.ShippingStatusId = &scode
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, *data.ShippingStatusId)
	data.StatusHistory = result
	err = r.db.Model(&shipping.ShippingOrder{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
	// response, _ := r.Get(data.ID)
	// return &response, res.Error
}

func (r *shippingorder) AddProduct(data shipping.ShippingOrderLines) error {
	res := r.db.Model(&shipping.ShippingOrderLines{}).Create(&data)
	return res.Error
}

func (r *shippingorder) List(query map[string]interface{}, page *pagination.Paginatevalue) ([]shipping.ShippingOrder, error) {
	var data []shipping.ShippingOrder
	err := r.db.Scopes(helpers.Paginate(&shipping.ShippingOrder{}, page, r.db)).Preload("ShippingOrderLines.ProductVariant").Preload("ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Preload("CreatedBy.Company").Where("is_active = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *shippingorder) Get(query map[string]interface{}) (shipping.ShippingOrder, error) {
	var data shipping.ShippingOrder
	err := r.db.Model(&shipping.ShippingOrder{}).Preload("ShippingOrderLines.ProductVariant").Preload("ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *shippingorder) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&shipping.ShippingOrder{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *shippingorder) Update(query map[string]interface{}, data *shipping.ShippingOrder) error {
	err := r.db.Model(&shipping.ShippingOrder{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *shippingorder) UpdateProductList(query interface{}, data shipping.ShippingOrderLines) (int64, error) {
	res := r.db.Model(&shipping.ShippingOrderLines{}).Where(query).Updates(&data)
	return res.RowsAffected, res.Error
}

func (r *shippingorder) DeleteOrderLine(query map[string]interface{}) error {
	var data shipping.ShippingOrderLines
	res := r.db.Model(&shipping.ShippingOrderLines{}).Where(query).Delete(&data)
	return res.Error
}

func (r *shippingorder) CreateBulkOrder(data *[]shipping.ShippingOrder) error {
	res := r.db.Model(&shipping.ShippingOrder{}).Create(&data)
	return res.Error
}
