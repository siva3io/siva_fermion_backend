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
type InventoryAdjustments interface {
	CreateInvAdj(data *inventory_orders.InventoryAdjustments, userID string) (uint, error)
	BulkCreateInvAdj(data *[]inventory_orders.InventoryAdjustments, userID string) error
	UpdateInvAdj(id uint, data *inventory_orders.InventoryAdjustments) error
	GetInvAdj(id uint) (inventory_orders.InventoryAdjustments, error)
	GetAllInvAdj(p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error)
	DeleteInvAdj(id uint, user_id uint) error

	CreateInvAdjLines(data inventory_orders.InventoryAdjustmentLines) error
	UpdateInvAdjLines(query interface{}, data inventory_orders.InventoryAdjustmentLines) (int64, error)
	GetInvAdjLines(query interface{}) ([]inventory_orders.InventoryAdjustmentLines, error)
	DeleteInvAdjLines(query interface{}) error
}

type invadj struct {
	db *gorm.DB
}

func NewInvAdj() *invadj {
	db := db.DbManager()
	return &invadj{db}
}

func (r *invadj) CreateInvAdj(data *inventory_orders.InventoryAdjustments, userID string) (uint, error) {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'ADJUSTMENT_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return 0, err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	result := r.db.Model(&inventory_orders.InventoryAdjustments{}).Create(&data)
	return data.ID, result.Error
}

func (r *invadj) BulkCreateInvAdj(data *[]inventory_orders.InventoryAdjustments, userID string) error {
	for _, value := range *data {
		var scode uint
		err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'ADJUSTMENT_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
		if err != nil {
			return err
		}
		value.StatusID = scode
		res, _ := helpers.UpdateStatusHistory(value.StatusHistory, value.StatusID)
		value.StatusHistory = res
		result := r.db.Model(&inventory_orders.InventoryAdjustments{}).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *invadj) UpdateInvAdj(id uint, data *inventory_orders.InventoryAdjustments) error {
	result := r.db.Model(&inventory_orders.InventoryAdjustments{}).Where("id", id).Updates(&data)
	return result.Error
}

func (r *invadj) GetInvAdj(id uint) (inventory_orders.InventoryAdjustments, error) {
	var data inventory_orders.InventoryAdjustments
	result := r.db.Model(&inventory_orders.InventoryAdjustments{}).Preload(clause.Associations).Where("id", id).First(&data)
	return data, result.Error
}

func (r *invadj) GetAllInvAdj(p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error) {
	var data []inventory_orders.InventoryAdjustments
	res := r.db.Model(&inventory_orders.InventoryAdjustments{}).Scopes(helpers.Paginate(&inventory_orders.InventoryAdjustments{}, p, r.db)).Where("is_active = true").Preload("InventoryAdjustmentLines.Product").Preload("InventoryAdjustmentLines.ProductVariant").Preload(clause.Associations).Find(&data)
	return data, res.Error
}

func (r *invadj) DeleteInvAdj(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	result := r.db.Model(&inventory_orders.InventoryAdjustments{}).Where("id", id).Updates(data)
	return result.Error
}

func (r *invadj) CreateInvAdjLines(data inventory_orders.InventoryAdjustmentLines) error {
	result := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Create(&data)
	return result.Error
}

func (r *invadj) UpdateInvAdjLines(query interface{}, data inventory_orders.InventoryAdjustmentLines) (int64, error) {
	result := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Where(query).Updates(&data)
	return result.RowsAffected, result.Error
}

func (r *invadj) GetInvAdjLines(query interface{}) ([]inventory_orders.InventoryAdjustmentLines, error) {
	var data []inventory_orders.InventoryAdjustmentLines
	result := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Where(query).Preload(clause.Associations).Find(&data)
	return data, result.Error
}

func (r *invadj) DeleteInvAdjLines(query interface{}) error {
	result := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Where(query).Delete(&inventory_orders.InventoryAdjustmentLines{})
	return result.Error
}
