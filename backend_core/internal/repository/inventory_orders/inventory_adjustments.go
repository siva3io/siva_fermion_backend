package inventory_orders

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/core"
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
	CreateInvAdj(data *inventory_orders.InventoryAdjustments) error
	BulkCreateInvAdj(metaData core.MetaData, data *[]inventory_orders.InventoryAdjustments) error
	UpdateInvAdj(query map[string]interface{}, data *inventory_orders.InventoryAdjustments) error
	GetInvAdj(query map[string]interface{}) (inventory_orders.InventoryAdjustments, error)
	GetAllInvAdj(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error)
	DeleteInvAdj(query map[string]interface{}) error

	CreateInvAdjLines(data inventory_orders.InventoryAdjustmentLines) error
	UpdateInvAdjLines(query map[string]interface{}, data inventory_orders.InventoryAdjustmentLines) (int64, error)
	GetInvAdjLines(query map[string]interface{}) ([]inventory_orders.InventoryAdjustmentLines, error)
	DeleteInvAdjLines(query map[string]interface{}) error
}

type invadj struct {
	db *gorm.DB
}

var invadjRepository *invadj //singleton object

// singleton function
func NewInvAdj() *invadj {
	if invadjRepository != nil {
		return invadjRepository
	}
	db := db.DbManager()
	invadjRepository = &invadj{db}
	return invadjRepository
}

func (r *invadj) CreateInvAdj(data *inventory_orders.InventoryAdjustments) error {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'ADJUSTMENT_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	result := r.db.Model(&inventory_orders.InventoryAdjustments{}).Create(data).Error
	if result != nil {
		return err
	}
	return nil
}

func (r *invadj) BulkCreateInvAdj(metaData core.MetaData, data *[]inventory_orders.InventoryAdjustments) error {
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

func (r *invadj) UpdateInvAdj(query map[string]interface{}, data *inventory_orders.InventoryAdjustments) error {
	err := r.db.Model(&inventory_orders.InventoryAdjustments{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *invadj) GetInvAdj(query map[string]interface{}) (inventory_orders.InventoryAdjustments, error) {
	var data inventory_orders.InventoryAdjustments
	err := r.db.Model(&inventory_orders.InventoryAdjustments{}).Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err != nil {
		return data, err.Error
	}
	return data, nil

}

func (r *invadj) GetAllInvAdj(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error) {
	var data []inventory_orders.InventoryAdjustments
	err := r.db.Model(&inventory_orders.InventoryAdjustments{}).Scopes(helpers.Paginate(&inventory_orders.InventoryAdjustments{}, p, r.db)).Where("is_active = true").Preload("InventoryAdjustmentLines.Product").Preload("InventoryAdjustmentLines.ProductVariant").Preload(clause.Associations).Find(&data)
	if err != nil {
		return data, err.Error
	}
	return data, nil

}

func (r *invadj) DeleteInvAdj(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&inventory_orders.InventoryAdjustments{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil

}

func (r *invadj) CreateInvAdjLines(data inventory_orders.InventoryAdjustmentLines) error {
	res := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Create(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *invadj) UpdateInvAdjLines(query map[string]interface{}, data inventory_orders.InventoryAdjustmentLines) (int64, error) {
	err := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Where(query).Updates(data)
	if err != nil {
		return err.RowsAffected, err.Error
	}
	return err.RowsAffected, nil
}

func (r *invadj) GetInvAdjLines(query map[string]interface{}) ([]inventory_orders.InventoryAdjustmentLines, error) {
	var data []inventory_orders.InventoryAdjustmentLines
	err := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Where(query).Preload(clause.Associations).Find(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *invadj) DeleteInvAdjLines(query map[string]interface{}) error {
	err := r.db.Model(&inventory_orders.InventoryAdjustmentLines{}).Where(query).Delete(&inventory_orders.InventoryAdjustmentLines{})
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil

}
