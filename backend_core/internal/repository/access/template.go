package access

import (
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

	"gorm.io/gorm"
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
type Template interface {
	CreateTemplate(data *model_core.AccessTemplate) (uint, error)
	UpdateTemplate(id uint, data model_core.AccessTemplate) (int64, error)
	FindAllTemplate(p *pagination.Paginatevalue) ([]model_core.AccessTemplate, error)
	FindByidTemplate(id uint) (model_core.AccessTemplate, error)
	DeleteTemplate(id uint, user_id uint) error
}

type template struct {
	db *gorm.DB
}

var TemplateRepository *template //singleton object

//singleton function

func NewTemplate() *template {
	if TemplateRepository != nil {
		return TemplateRepository
	}
	db := db.DbManager()
	TemplateRepository = &template{db}
	return TemplateRepository
}

func (r *template) CreateTemplate(data *model_core.AccessTemplate) (uint, error) {

	res := r.db.Model(&model_core.AccessTemplate{}).Create(&data)
	fmt.Println(data.ID)
	if res.Error != nil {
		return 0, res.Error
	}

	return data.ID, nil
}
func (r *template) UpdateTemplate(id uint, data model_core.AccessTemplate) (int64, error) {

	res := r.db.Model(&model_core.AccessTemplate{}).Where("id", id).Updates(&data)
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}
func (r *template) FindAllTemplate(p *pagination.Paginatevalue) ([]model_core.AccessTemplate, error) {

	var result []model_core.AccessTemplate
	//res := r.db.Preload(clause.Associations).Model(&model_core.AccessTemplate{}).Preload(clause.Associations).Scopes(helpers.Paginate(&model_core.AccessTemplate{}, p, r.db)).Where("is_active = true").Find(&result)
	res := r.db.Model(&model_core.AccessTemplate{}).Scopes(helpers.Paginate(&model_core.AccessTemplate{}, p, r.db)).Where("is_active = true")
	res.Find(&result)
	// fmt.Println(final_res.Error)

	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
func (r *template) FindByidTemplate(id uint) (model_core.AccessTemplate, error) {

	var market_user model_core.AccessTemplate

	res := r.db.Model(&model_core.AccessTemplate{}).Where("id", id).First(&market_user)
	if res.Error != nil {
		return market_user, res.Error
	}
	return market_user, nil
}
func (r *template) DeleteTemplate(id uint, user_id uint) error {

	if id == 1 {
		return fmt.Errorf("you are not able to delete super admin")
	}
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}

	result := r.db.Model(&model_core.AccessTemplate{}).Where("id", id).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
