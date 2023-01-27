package access

import (
	"os"
	"time"

	"fermion/backend_core/db"
	model_core "fermion/backend_core/internal/model/core"
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
type Module interface {
	CreateModule(data *model_core.AccessModuleAction) error
	UpdateModule(id uint, data model_core.AccessModuleAction) (int64, error)
	FindAllModule(p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error)
	FindAllAccessModule(p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error)
	FindByidModule(id uint, p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error)
	DeleteModule(id uint, user_id uint) error
	FindAllCoreModule(p *pagination.Paginatevalue) ([]model_core.CoreAppModule, error)
}

type module struct {
	db *gorm.DB
}

var ModuleRepository *module //singleton object

// singleton function
func NewModule() *module {
	if ModuleRepository != nil {
		return ModuleRepository
	}
	db := db.DbManager()
	ModuleRepository = &module{db}
	return ModuleRepository
}

func (r *module) CreateModule(data *model_core.AccessModuleAction) error {

	res := r.db.Model(&model_core.AccessModuleAction{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *module) UpdateModule(id uint, data model_core.AccessModuleAction) (int64, error) {

	res := r.db.Model(&model_core.AccessModuleAction{}).Where("id", id).Updates(&data)
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}

func (r *module) FindAllModule(p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error) {

	var result []model_core.AccessModuleAction
	res := r.db.Model(&model_core.AccessModuleAction{}).Preload(clause.Associations + "." + clause.Associations).Scopes(helpers.Paginate(&model_core.AccessModuleAction{}, p, r.db)).Where("is_active = true").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
func (r *module) FindAllAccessModule(p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error) {

	var result []model_core.AccessModuleAction
	res := r.db.Model(&model_core.AccessModuleAction{}).Scopes(helpers.Paginate(&model_core.AccessModuleAction{}, p, r.db)).Where("is_active = true AND view_actions @> ?", "[{ \"ctrl_flag\": 1 }]").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
func (r *module) FindByidModule(id uint, p *pagination.Paginatevalue) ([]model_core.AccessModuleAction, error) {

	var market_user []model_core.AccessModuleAction

	res := r.db.Preload(clause.Associations).Model(&model_core.AccessModuleAction{}).Scopes(helpers.Paginate(&model_core.AccessModuleAction{}, p, r.db)).Where("access_template_id", id).Find(&market_user)
	if res.Error != nil {
		return market_user, res.Error
	}
	return market_user, nil
}

func (r *module) DeleteModule(id uint, user_id uint) error {

	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}

	result := r.db.Model(&model_core.AccessModuleAction{}).Where("id", id).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *module) FindAllCoreModule(p *pagination.Paginatevalue) ([]model_core.CoreAppModule, error) {

	var result []model_core.CoreAppModule
	p.Sort = "[[\"id\",\"asc\"],[\"item_sequence\",\"asc\"]]"
	res := r.db.Model(&model_core.CoreAppModule{}).Preload(clause.Associations).Scopes(helpers.Paginate(&model_core.CoreAppModule{}, p, r.db)).Where("is_active = true").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
