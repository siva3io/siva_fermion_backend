package inventory_orders

import (
	"errors"
	"fmt"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/inventory_orders"
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
type InventoryOrdersBase interface {
	CreateRecord(data inventory_orders.InventoryOrdersFav) error
	Favourite(query map[string]interface{}) error
	UnFavourite(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (inventory_orders.InventoryOrdersFav, error)
	GetAsnFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.ASN, error)
	GetGrnFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.GRN, error)
	GetInvAdjFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error)
}

type inventoryOrdersBase struct {
	db *gorm.DB
}

var inventoryOrdersBaseRepository *inventoryOrdersBase //singleton object

// singleton function
func NewInventoryOrdersBase() *inventoryOrdersBase {
	if inventoryOrdersBaseRepository != nil {
		return inventoryOrdersBaseRepository
	}
	db := db.DbManager()
	inventoryOrdersBaseRepository = &inventoryOrdersBase{db}
	return inventoryOrdersBaseRepository

}
func (r *inventoryOrdersBase) GetInvAdjFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error) {
	var result []inventory_orders.InventoryAdjustments
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&inventory_orders.InventoryAdjustments{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&inventory_orders.InventoryAdjustments{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *inventoryOrdersBase) CreateRecord(data inventory_orders.InventoryOrdersFav) error {
	err := r.db.Model(&inventory_orders.InventoryOrdersFav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *inventoryOrdersBase) Favourite(query map[string]interface{}) error {
	err := r.db.Model(&inventory_orders.InventoryOrdersFav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *inventoryOrdersBase) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&inventory_orders.InventoryOrdersFav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *inventoryOrdersBase) FindOne(query map[string]interface{}) (inventory_orders.InventoryOrdersFav, error) {
	var data inventory_orders.InventoryOrdersFav
	err := r.db.Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	fmt.Println("Findone:", data)
	return data, nil
}

func (r *inventoryOrdersBase) GetAsnFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.ASN, error) {
	var result []inventory_orders.ASN
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&inventory_orders.ASN{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&inventory_orders.ASN{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *inventoryOrdersBase) GetGrnFavourites(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.GRN, error) {
	var result []inventory_orders.GRN
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&inventory_orders.GRN{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&inventory_orders.GRN{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
