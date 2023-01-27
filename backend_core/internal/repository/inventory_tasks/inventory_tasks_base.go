package inventory_tasks

import (
	"errors"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/inventory_tasks"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

	"github.com/lib/pq"
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
type InventoryTasksBase interface {
	CreateRecord(data inventory_tasks.InventoryTasksFav) error
	Favourite(query map[string]interface{}) error
	UnFavourite(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (inventory_tasks.InventoryTasksFav, error)
	GetCycleCountFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.CycleCount, error)
	GetPickListFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error)
}

type inventoryTasksBase struct {
	db *gorm.DB
}

var inventoryTasksBaseRepository *inventoryTasksBase //singleton object

// singleton function
func NewInventoryTasksBase() *inventoryTasksBase {
	if inventoryTasksBaseRepository != nil {
		return inventoryTasksBaseRepository
	}
	db := db.DbManager()
	inventoryTasksBaseRepository = &inventoryTasksBase{db}
	return inventoryTasksBaseRepository

}

func (r *inventoryTasksBase) CreateRecord(data inventory_tasks.InventoryTasksFav) error {
	err := r.db.Model(&inventory_tasks.InventoryTasksFav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *inventoryTasksBase) Favourite(query map[string]interface{}) error {
	err := r.db.Model(&inventory_tasks.InventoryTasksFav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *inventoryTasksBase) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&inventory_tasks.InventoryTasksFav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *inventoryTasksBase) FindOne(query map[string]interface{}) (inventory_tasks.InventoryTasksFav, error) {
	var data inventory_tasks.InventoryTasksFav
	err := r.db.Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	return data, nil
}

func (r *inventoryTasksBase) GetCycleCountFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.CycleCount, error) {
	var result []inventory_tasks.CycleCount
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&inventory_tasks.CycleCount{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&inventory_tasks.CycleCount{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil

}
func (r *inventoryTasksBase) GetPickListFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error) {
	var result []inventory_tasks.PickList
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&inventory_tasks.PickList{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&inventory_tasks.PickList{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil

}
