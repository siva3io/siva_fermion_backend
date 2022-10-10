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
type GRN interface {
	UpdateGRN(id uint, data *inventory_orders.GRN) error
	DeleteGRN(id uint, user_id uint) error
	UpdateOrderLines(query interface{}, data inventory_orders.GRNOrderLines) (int64, error)
	CreateGRN(data *inventory_orders.GRN, userID string) (uint, error)
	BulkCreateGRN(data *[]inventory_orders.GRN, userID string) error
	CreateOrderLines(data inventory_orders.GRNOrderLines) error
	GetAllGRN(page *pagination.Paginatevalue) ([]inventory_orders.GRN, error)
	GetGRN(id uint) (inventory_orders.GRN, error)
	DeleteOrderLines(query interface{}) error
	SearchGRN(key string) (interface{}, error)
}

type grn struct {
	db *gorm.DB
}

func NewGRN() *grn {
	db := db.DbManager()
	return &grn{db}
}

func (r *grn) CreateGRN(data *inventory_orders.GRN, userID string) (uint, error) {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'GRN_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return 0, err
	}
	data.StatusId = scode
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	res := r.db.Model(&inventory_orders.GRN{}).Create(&data)
	return data.ID, res.Error
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
	return res.Error
}

func (r *grn) GetGRN(id uint) (inventory_orders.GRN, error) {
	var data inventory_orders.GRN
	result := r.db.Model(&inventory_orders.GRN{}).Preload(clause.Associations+"."+clause.Associations+"."+clause.Associations+"."+clause.Associations).Preload("GRNOrderLines.Product").Preload("GRNOrderLines.ProductTemplate").Preload("GRNOrderLines.UOM").Preload(clause.Associations).Where("id", id).First(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

func (r *grn) GetAllGRN(page *pagination.Paginatevalue) ([]inventory_orders.GRN, error) {
	var data []inventory_orders.GRN
	result := r.db.Scopes(helpers.Paginate(&inventory_orders.GRN{}, page, r.db)).Preload("GRNOrderLines.Product").Preload("GRNOrderLines.ProductTemplate").Preload("GRNOrderLines.UOM").Preload(clause.Associations).Where("is_active = true").Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (r *grn) UpdateGRN(id uint, data *inventory_orders.GRN) error {
	res := r.db.Model(&inventory_orders.GRN{}).Where("id", id).Updates(&data)
	return res.Error
}

func (r *grn) UpdateOrderLines(query interface{}, data inventory_orders.GRNOrderLines) (int64, error) {
	res := r.db.Model(&inventory_orders.GRNOrderLines{}).Where(query).Updates(&data)
	return res.RowsAffected, res.Error
}

func (r *grn) SearchGRN(key string) (interface{}, error) {
	var datas []inventory_orders.GRN
	fields := []string{"reference_number", "grn_number", "putaway_location"}
	fields_string, values := helpers.ApplySearch(key, fields)
	res := r.db.Model(&inventory_orders.GRN{}).Limit(50).Preload(clause.Associations).Where(fields_string, values...).Find(&datas)
	return datas, res.Error
}

func (r *grn) DeleteGRN(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&inventory_orders.GRN{}).Where("id", id).Updates(data)
	return res.Error
}

func (r *grn) DeleteOrderLines(query interface{}) error {
	var data inventory_orders.GRNOrderLines
	res := r.db.Model(&inventory_orders.GRNOrderLines{}).Where(query).Delete(&data)
	return res.Error
}
