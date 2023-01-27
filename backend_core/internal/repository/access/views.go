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
type View interface {
	CreateView(data *model_core.ViewsLevelAccessItems) error
	UpdateView(id uint, data model_core.ViewsLevelAccessItems) (int64, error)
	FindAllView(p *pagination.Paginatevalue) ([]model_core.ViewsLevelAccessItems, error)
	FindByidView(id uint) (model_core.ViewsLevelAccessItems, error)
	DeleteView(id uint, user_id uint) error
}

type view struct {
	db *gorm.DB
}

var viewRepository *view //singleton object

// singleton function
func NewView() *view {
	if viewRepository != nil {
		return viewRepository
	}
	db := db.DbManager()
	viewRepository = &view{db}
	return &view{db}
}

func (r *view) CreateView(data *model_core.ViewsLevelAccessItems) error {

	res := r.db.Model(&model_core.ViewsLevelAccessItems{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
func (r *view) UpdateView(id uint, data model_core.ViewsLevelAccessItems) (int64, error) {

	res := r.db.Model(&model_core.ViewsLevelAccessItems{}).Where("id", id).Updates(&data)
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}
func (r *view) FindAllView(p *pagination.Paginatevalue) ([]model_core.ViewsLevelAccessItems, error) {

	var result []model_core.ViewsLevelAccessItems
	res := r.db.Preload(clause.Associations).Model(&model_core.ViewsLevelAccessItems{}).Preload(clause.Associations).Scopes(helpers.Paginate(&model_core.ViewsLevelAccessItems{}, p, r.db)).Where("is_active = true").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
func (r *view) FindByidView(id uint) (model_core.ViewsLevelAccessItems, error) {

	var market_user model_core.ViewsLevelAccessItems

	res := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&model_core.ViewsLevelAccessItems{}).Preload(clause.Associations).Where("id", id).First(&market_user)
	if res.Error != nil {
		return market_user, res.Error
	}
	return market_user, nil
}
func (r *view) DeleteView(id uint, user_id uint) error {

	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}

	result := r.db.Model(&model_core.ViewsLevelAccessItems{}).Where("id", id).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
