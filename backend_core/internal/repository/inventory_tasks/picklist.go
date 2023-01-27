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
type Picklist interface {
	CreatePickList(data *inventory_tasks.PickList) error
	BulkCreatePickList(data *[]inventory_tasks.PickList) error
	UpdatePickList(query map[string]interface{}, data *inventory_tasks.PickList) error
	GetPickList(query map[string]interface{}) (inventory_tasks.PickList, error)
	GetAllPickList(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error)
	DeletePickList(query map[string]interface{}) error

	CreatePickListLines(data *inventory_tasks.PickListLines) error
	UpdatePickListLines(query map[string]interface{}, data *inventory_tasks.PickListLines) (int64, error)
	GetPickListLines(query interface{}) ([]inventory_tasks.PickListLines, error)
	DeletePickListLines(query interface{}) error
}

type picklist struct {
	db *gorm.DB
}

var pickListRepository *picklist //singleton object

// singleton function
func NewPicklist() *picklist {
	if pickListRepository != nil {
		return pickListRepository
	}
	db := db.DbManager()
	pickListRepository = &picklist{db}
	return pickListRepository
}

func (r *picklist) CreatePickList(data *inventory_tasks.PickList) error {
	var scode uint
	err := r.db.Raw("SELECT lookupcodes.id FROM lookuptypes,lookupcodes WHERE lookuptypes.id = lookupcodes.lookup_type_id AND lookuptypes.lookup_type = 'PICK_LIST_STATUS' AND lookupcodes.lookup_code = 'IN_PROGRESS'").First(&scode).Error
	if err != nil {
		return err
	}
	data.StatusID = scode
	res, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusID)
	data.StatusHistory = res
	er := r.db.Model(&inventory_tasks.PickList{}).Create(&data).Error
	if er != nil {
		return er
	}
	return nil
}

func (r *picklist) BulkCreatePickList(data *[]inventory_tasks.PickList) error {
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

func (r *picklist) UpdatePickList(query map[string]interface{}, data *inventory_tasks.PickList) error {
	err := r.db.Model(&inventory_tasks.PickList{}).Where(query).Updates(&data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *picklist) GetPickList(query map[string]interface{}) (inventory_tasks.PickList, error) {
	var data inventory_tasks.PickList
	err := r.db.Model(&inventory_tasks.PickList{}).Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *picklist) GetAllPickList(query map[string]interface{}, p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error) {
	var data []inventory_tasks.PickList
	res := r.db.Model(&inventory_tasks.PickList{}).Scopes(helpers.Paginate(&inventory_tasks.PickList{}, p, r.db)).Where("is_active = true").Preload("PicklistLines.Product").Preload("PicklistLines.ProductVariant").Preload("PicklistLines.Partner").Preload(clause.Associations).Find(&data)
	return data, res.Error
}

func (r *picklist) DeletePickList(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&inventory_tasks.PickList{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *picklist) CreatePickListLines(data *inventory_tasks.PickListLines) error {
	result := r.db.Model(&inventory_tasks.PickListLines{}).Create(&data)
	return result.Error
}

func (r *picklist) UpdatePickListLines(query map[string]interface{}, data *inventory_tasks.PickListLines) (int64, error) {
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
