package omnichannel

import (
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/pkg/util/helpers"

	"fermion/backend_core/internal/model/pagination"

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
type VirtualWarehouse interface {
	CreateVirtualWarehouse(data *omnichannel.User_Virtual_Warehouse_Registration) error
	UpdateVirtualWarehouse(id uint, data *omnichannel.User_Virtual_Warehouse_Registration) error
	FindAllVirtualWarehouses(p *pagination.Paginatevalue) ([]omnichannel.User_Virtual_Warehouse_Registration, error)
	FindByIDVirtualWarehouse(id uint) (interface{}, error)
	DeleteVirtualWarehouse(id uint, user_id uint) error
	AvailableVirtualWarehouses(p *pagination.Paginatevalue) ([]omnichannel.VirtualWarehouse, error)
}

type virtual_warehouse struct {
	db *gorm.DB
}

var virtualWarehouseRegisterRepository *virtual_warehouse //singleton object

// singleton function
func NewVirtualWarehouse() *virtual_warehouse {
	if virtualWarehouseRegisterRepository != nil {
		return virtualWarehouseRegisterRepository
	}
	db := db.DbManager()
	virtualWarehouseRegisterRepository = &virtual_warehouse{db}
	return virtualWarehouseRegisterRepository
}

func (r *virtual_warehouse) CreateVirtualWarehouse(data *omnichannel.User_Virtual_Warehouse_Registration) error {

	res := r.db.Model(&omnichannel.User_Virtual_Warehouse_Registration{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *virtual_warehouse) UpdateVirtualWarehouse(id uint, data *omnichannel.User_Virtual_Warehouse_Registration) error {

	res := r.db.Model(&omnichannel.User_Virtual_Warehouse_Registration{}).Where("id", id).Updates(&data)
	return res.Error

}

func (r *virtual_warehouse) FindAllVirtualWarehouses(p *pagination.Paginatevalue) ([]omnichannel.User_Virtual_Warehouse_Registration, error) {

	var result []omnichannel.User_Virtual_Warehouse_Registration
	res := r.db.Preload(clause.Associations).Model(&omnichannel.User_Virtual_Warehouse_Registration{}).Preload(clause.Associations).Scopes(helpers.Paginate(&omnichannel.User_Virtual_Warehouse_Registration{}, p, r.db)).Where("is_active = true").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (r *virtual_warehouse) FindByIDVirtualWarehouse(id uint) (interface{}, error) {
	var result omnichannel.User_Virtual_Warehouse_Registration
	res := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&omnichannel.User_Virtual_Warehouse_Registration{}).Where("id", id).Preload(clause.Associations).Find(&result)
	if res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected == 0 {
		return nil, fmt.Errorf("no Data or record Found")
	}
	return result, nil
}

func (r *virtual_warehouse) DeleteVirtualWarehouse(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&omnichannel.User_Virtual_Warehouse_Registration{}).Where("id", id).Updates(data)
	r.db.Where("id", id).Delete(&omnichannel.Bank_info{})
	r.db.Where("id", id).Delete(&omnichannel.Kyc_info{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *virtual_warehouse) AvailableVirtualWarehouses(p *pagination.Paginatevalue) ([]omnichannel.VirtualWarehouse, error) {

	var result []omnichannel.VirtualWarehouse
	res := r.db.Model(&omnichannel.VirtualWarehouse{}).Preload(clause.Associations).Scopes(helpers.Paginate(&omnichannel.VirtualWarehouse{}, p, r.db)).Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
