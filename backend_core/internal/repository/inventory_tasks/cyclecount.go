package inventory_tasks

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/inventory_tasks"
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
type CycleCount interface {
	CreateCycleCount(data *inventory_tasks.CycleCount) error
	BulkCreateCycleCount(data *[]inventory_tasks.CycleCount) error
	UpdateCycleCount(query map[string]interface{}, data *inventory_tasks.CycleCount) error
	GetCycleCount(query map[string]interface{}) (inventory_tasks.CycleCount, error)
	GetAllCycleCount(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_tasks.CycleCount, error)
	DeleteCycleCount(query map[string]interface{}) error

	CreateCycleCountLines(data *inventory_tasks.CycleCountLines) error
	UpdateCycleCountLines(query map[string]interface{}, data *inventory_tasks.CycleCountLines) (int64, error)
	GetCycleCountLines(query map[string]interface{}) ([]inventory_tasks.CycleCountLines, error)
	DeleteCycleCountLines(interface{}) error
}

type cycleCount struct {
	db *gorm.DB
}

var cycleCountRepository *cycleCount //singleton object

// singleton function
func NewCycleCount() *cycleCount {
	if cycleCountRepository != nil {
		return cycleCountRepository
	}
	db := db.DbManager()
	cycleCountRepository = &cycleCount{db}
	return cycleCountRepository
}

func (r *cycleCount) CreateCycleCount(data *inventory_tasks.CycleCount) error {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'CYCLE_COUNT_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
	if err != nil {
		return err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	er := r.db.Model(&inventory_tasks.CycleCount{}).Create(&data).Error
	if er != nil {
		return er
	}
	return nil
}

func (r *cycleCount) BulkCreateCycleCount(data *[]inventory_tasks.CycleCount) error {
	for _, value := range *data {
		var scode uint
		err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'CYCLE_COUNT_STATUS' AND lookupcodes.lookup_code = 'DRAFT'").First(&scode).Error
		if err != nil {
			return err
		}
		value.StatusID = scode
		res, _ := helpers.UpdateStatusHistory(value.StatusHistory, value.StatusID)
		value.StatusHistory = res
		result := r.db.Model(&inventory_tasks.CycleCount{}).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *cycleCount) UpdateCycleCount(query map[string]interface{}, data *inventory_tasks.CycleCount) error {
	result := r.db.Model(&inventory_tasks.CycleCount{}).Where(query).Updates(&data)
	return result.Error
}

func (r *cycleCount) GetCycleCount(query map[string]interface{}) (inventory_tasks.CycleCount, error) {
	var data inventory_tasks.CycleCount
	result := r.db.Model(&inventory_tasks.CycleCount{}).Preload(clause.Associations).Where(query).First(&data)
	return data, result.Error
}

func (r *cycleCount) GetAllCycleCount(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_tasks.CycleCount, error) {
	var data []inventory_tasks.CycleCount
	res := r.db.Model(&inventory_tasks.CycleCount{}).Scopes(helpers.Paginate(&inventory_tasks.CycleCount{}, p, r.db)).Where("is_active = true").Preload("OrderLines.Product").Preload("OrderLines.ProductVariant").Preload("OrderLines.LocationSpaceType").Preload("OrderLines.LocationInputType").Preload(clause.Associations).Find(&data)
	return data, res.Error
}

func (r *cycleCount) DeleteCycleCount(query map[string]interface{}) error {

	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&inventory_tasks.CycleCount{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *cycleCount) CreateCycleCountLines(data *inventory_tasks.CycleCountLines) error {
	result := r.db.Model(&inventory_tasks.CycleCountLines{}).Create(&data)
	return result.Error
}

func (r *cycleCount) UpdateCycleCountLines(query map[string]interface{}, data *inventory_tasks.CycleCountLines) (int64, error) {
	result := r.db.Model(&inventory_tasks.CycleCountLines{}).Where(query).Updates(&data)
	return result.RowsAffected, result.Error
}

func (r *cycleCount) GetCycleCountLines(query map[string]interface{}) ([]inventory_tasks.CycleCountLines, error) {
	var data []inventory_tasks.CycleCountLines
	result := r.db.Model(&inventory_tasks.CycleCountLines{}).Where(query).Preload(clause.Associations).Find(&data)
	return data, result.Error
}

func (r *cycleCount) DeleteCycleCountLines(query interface{}) error {
	result := r.db.Model(&inventory_tasks.CycleCountLines{}).Where(query).Delete(&inventory_tasks.CycleCountLines{})
	return result.Error
}
