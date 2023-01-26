package inventory_tasks

import (
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
type Picklist interface {
	CreatePickList(data *inventory_tasks.PickList, userID string) (uint, error)
	BulkCreatePickList(data *[]inventory_tasks.PickList, userID string) error
	UpdatePickList(id uint, data *inventory_tasks.PickList) error
	GetPickList(id uint) (inventory_tasks.PickList, error)
	GetAllPickList(p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error)
	DeletePickList(id uint, user_id uint) error

	CreatePickListLines(data inventory_tasks.PickListLines) error
	UpdatePickListLines(query interface{}, data inventory_tasks.PickListLines) (int64, error)
	GetPickListLines(query interface{}) ([]inventory_tasks.PickListLines, error)
	DeletePickListLines(query interface{}) error
}

type picklist struct {
	db *gorm.DB
}

func NewPicklist() *picklist {
	db := db.DbManager()
	return &picklist{db}
}

func (r *picklist) CreatePickList(data *inventory_tasks.PickList, userID string) (uint, error) {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'PICK_LIST_STATUS' AND lookupcodes.lookup_code = 'IN_PROGRESS'").First(&scode).Error
	if err != nil {
		return 0, err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	result := r.db.Model(&inventory_tasks.PickList{}).Create(&data)
	return data.ID, result.Error
}

func (r *picklist) BulkCreatePickList(data *[]inventory_tasks.PickList, userID string) error {
	for _, value := range *data {
		var scode uint
		err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'PICK_LIST_STATUS' AND lookupcodes.lookup_code = 'IN_PROGRESS'").First(&scode).Error
		if err != nil {
			return err
		}
		value.StatusID = scode
		res, _ := helpers.UpdateStatusHistory(value.StatusHistory, value.StatusID)
		value.StatusHistory = res
		result := r.db.Model(&inventory_tasks.PickList{}).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *picklist) UpdatePickList(id uint, data *inventory_tasks.PickList) error {
	result := r.db.Model(&inventory_tasks.PickList{}).Where("id", id).Updates(&data)
	return result.Error
}

func (r *picklist) GetPickList(id uint) (inventory_tasks.PickList, error) {
	var data inventory_tasks.PickList
	result := r.db.Model(&inventory_tasks.PickList{}).Preload(clause.Associations).Where("id", id).First(&data)
	return data, result.Error
}

func (r *picklist) GetAllPickList(p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error) {
	var data []inventory_tasks.PickList
	res := r.db.Model(&inventory_tasks.PickList{}).Scopes(helpers.Paginate(&inventory_tasks.PickList{}, p, r.db)).Where("is_active = true").Preload("PicklistLines.Product").Preload("PicklistLines.ProductVariant").Preload("PicklistLines.Partner").Preload(clause.Associations).Find(&data)
	return data, res.Error
}

func (r *picklist) DeletePickList(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	result := r.db.Model(&inventory_tasks.PickList{}).Where("id", id).Updates(data)
	return result.Error
}

func (r *picklist) CreatePickListLines(data inventory_tasks.PickListLines) error {
	result := r.db.Model(&inventory_tasks.PickListLines{}).Create(&data)
	return result.Error
}

func (r *picklist) UpdatePickListLines(query interface{}, data inventory_tasks.PickListLines) (int64, error) {
	result := r.db.Model(&inventory_tasks.PickListLines{}).Where(query).Updates(&data)
	return result.RowsAffected, result.Error
}

func (r *picklist) GetPickListLines(query interface{}) ([]inventory_tasks.PickListLines, error) {
	var data []inventory_tasks.PickListLines
	result := r.db.Model(&inventory_tasks.PickListLines{}).Where(query).Preload(clause.Associations).Find(&data)
	return data, result.Error
}

func (r *picklist) DeletePickListLines(query interface{}) error {
	result := r.db.Model(&inventory_tasks.PickListLines{}).Where(query).Delete(&inventory_tasks.PickListLines{})
	return result.Error
}
