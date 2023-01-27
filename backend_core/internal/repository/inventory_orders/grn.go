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
type GRN interface {
	CreateGRN(data *inventory_orders.GRN) error
	UpdateGRN(query map[string]interface{}, data *inventory_orders.GRN) error
	GetAllGRN(query map[string]interface{}, page *pagination.Paginatevalue) ([]inventory_orders.GRN, error)
	GetGRN(query map[string]interface{}) (inventory_orders.GRN, error)
	DeleteGRN(query map[string]interface{}) error
	BulkCreateGRN(data *[]inventory_orders.GRN, userID string) error

	CreateOrderLines(data inventory_orders.GRNOrderLines) error
	UpdateOrderLines(query map[string]interface{}, data inventory_orders.GRNOrderLines) (int64, error)
	DeleteOrderLines(query map[string]interface{}) error
}

type grn struct {
	db *gorm.DB
}

var grnRepository *grn //singleton object

// singleton function
func NewGRN() *grn {
	if grnRepository != nil {
		return grnRepository
	}
	db := db.DbManager()
	grnRepository = &grn{db}
	return grnRepository
}

func (r *grn) CreateGRN(data *inventory_orders.GRN) error {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'GRN_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return err
	}
	data.StatusId = scode
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	res := r.db.Model(&inventory_orders.GRN{}).Create(data).Error
	if res != nil {
		return err
	}
	return nil
}

func (r *grn) BulkCreateGRN(data *[]inventory_orders.GRN, userID string) error {
	for _, value := range *data {
		var scode uint
		err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'GRN_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
		if err != nil {
			return err
		}
		value.StatusId = scode
		res, _ := helpers.UpdateStatusHistory(value.StatusHistory, value.StatusId)
		value.StatusHistory = res
		result := r.db.Model(&inventory_orders.GRN{}).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *grn) CreateOrderLines(data inventory_orders.GRNOrderLines) error {
	res := r.db.Model(&inventory_orders.GRNOrderLines{}).Create(&data)
	if res != nil {
		return res.Error
	}
	return nil
}

func (r *grn) GetGRN(query map[string]interface{}) (inventory_orders.GRN, error) {
	var data inventory_orders.GRN
	result := r.db.Model(&inventory_orders.GRN{}).Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations).Preload("GRNOrderLines.Product").Preload("GRNOrderLines.ProductTemplate").Preload("GRNOrderLines.UOM").Preload(clause.Associations).Where(query).First(&data)
	if result.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

func (r *grn) GetAllGRN(query map[string]interface{}, page *pagination.Paginatevalue) ([]inventory_orders.GRN, error) {
	var data []inventory_orders.GRN
	result := r.db.Scopes(helpers.Paginate(&inventory_orders.GRN{}, page, r.db)).Preload("GRNOrderLines.Product").Preload("GRNOrderLines.ProductTemplate").Preload("GRNOrderLines.UOM").Preload(clause.Associations).Where("is_active = true").Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (r *grn) UpdateGRN(query map[string]interface{}, data *inventory_orders.GRN) error {
	res := r.db.Model(&inventory_orders.GRN{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *grn) UpdateOrderLines(query map[string]interface{}, data inventory_orders.GRNOrderLines) (int64, error) {
	res := r.db.Model(&inventory_orders.GRNOrderLines{}).Where(query).Updates(&data)
	if res.RowsAffected == 0 {
		return res.RowsAffected, errors.New("oops! record not found")
	}
	return res.RowsAffected, nil

}

func (r *grn) DeleteGRN(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&inventory_orders.GRN{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *grn) DeleteOrderLines(query map[string]interface{}) error {
	var data inventory_orders.GRNOrderLines
	res := r.db.Model(&inventory_orders.GRNOrderLines{}).Where(query).Delete(&data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	return res.Error
}
