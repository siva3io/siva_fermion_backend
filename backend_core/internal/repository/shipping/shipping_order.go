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
type ShippingOrder interface {
	Create(data *shipping.ShippingOrder) (*shipping.ShippingOrder, error)
	AddProduct(data shipping.ShippingOrderLines) error
	List(page *pagination.Paginatevalue) ([]shipping.ShippingOrder, error)
	Get(id uint) (shipping.ShippingOrder, error)
	Delete(id uint, user_id uint) error
	Update(id uint, data *shipping.ShippingOrder) error
	UpdateProductList(query interface{}, data shipping.ShippingOrderLines) (int64, error)
	DeleteOrderLine(query interface{}) error
	CreateBulkOrder(data *[]shipping.ShippingOrder, userID string) error
}

type shippingorder struct {
	db *gorm.DB
}

func NewShippingOrder() *shippingorder {
	db := db.DbManager()
	return &shippingorder{db}
}

func (r *shippingorder) Create(data *shipping.ShippingOrder) (*shipping.ShippingOrder, error) {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'SHIPPING_ORDER_STATUS' AND lookupcodes.lookup_code = 'PICKUP_SCHEDULED'").First(&scode).Error
	if err != nil {
		return nil, err
	}
	data.ShippingStatusId = scode
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.ShippingStatusId)
	data.StatusHistory = result
	res := r.db.Model(&shipping.ShippingOrder{}).Create(data)
	response, _ := r.Get(data.ID)
	return &response, res.Error
}

func (r *shippingorder) AddProduct(data shipping.ShippingOrderLines) error {
	res := r.db.Model(&shipping.ShippingOrderLines{}).Create(&data)
	return res.Error
}

func (r *shippingorder) List(page *pagination.Paginatevalue) ([]shipping.ShippingOrder, error) {
	var data []shipping.ShippingOrder
	result := r.db.Scopes(helpers.Paginate(&shipping.ShippingOrder{}, page, r.db)).Preload("ShippingOrderLines.ProductVariant").Preload("ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Where("is_active = true").Find(&data)
	return data, result.Error
}

func (r *shippingorder) Get(id uint) (shipping.ShippingOrder, error) {
	var data shipping.ShippingOrder
	result := r.db.Model(&shipping.ShippingOrder{}).Preload("ShippingOrderLines.ProductVariant").Preload("ShippingOrderLines.ProductTemplate").Preload(clause.Associations).Where("id", id).First(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

func (r *shippingorder) Delete(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&shipping.ShippingOrder{}).Where("id", id).Updates(data)
	return res.Error
}

func (r *shippingorder) Update(id uint, data *shipping.ShippingOrder) error {
	res := r.db.Model(&shipping.ShippingOrder{}).Where("id", id).Updates(data)
	return res.Error
}

func (r *shippingorder) UpdateProductList(query interface{}, data shipping.ShippingOrderLines) (int64, error) {
	res := r.db.Model(&shipping.ShippingOrderLines{}).Where(query).Updates(&data)
	return res.RowsAffected, res.Error
}

func (r *shippingorder) DeleteOrderLine(query interface{}) error {
	var data shipping.ShippingOrderLines
	res := r.db.Model(&shipping.ShippingOrderLines{}).Where(query).Delete(&data)
	return res.Error
}

func (r *shippingorder) CreateBulkOrder(data *[]shipping.ShippingOrder, userID string) error {
	res := r.db.Model(&shipping.ShippingOrder{}).Create(&data)
	return res.Error
}
