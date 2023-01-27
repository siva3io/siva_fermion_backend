package repository

import (
	"fmt"

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
type User interface {

	//----------------------core users-------------------------------
	Save(data *model_core.CoreUsers) error
	UpdateUser(query map[string]interface{}, data *model_core.CoreUsers) error
	FindUser(query map[string]interface{}) (model_core.CoreUsers, error)
	FindAll() ([]model_core.CoreUsers, error)
	FindLogin(query map[string]interface{}) (model_core.CoreUsers, error)
	UpdateAssignTemplate(id uint, data model_core.CoreUsers) (int64, error)
	UpdateCompany(query map[string]interface{}, data model_core.Company) (model_core.Company, error)
	GetCompany(query map[string]interface{}) (model_core.Company, error)
	FindAllCompanyUsers(query map[string]interface{}, page *pagination.Paginatevalue) ([]model_core.CoreUsers, error)
	UpdateUserDetails(id uint, data model_core.CoreUsers) (int64, error)
	FindAllCoreUser(query map[string]interface{}, p *pagination.Paginatevalue) ([]model_core.CoreUsers, error)

	FindAllCompanies(page *pagination.Paginatevalue) ([]model_core.Company, error)
}

type user struct {
	db *gorm.DB
}

func NewUser() *user {
	db := db.DbManager()
	return &user{db}
}

//------------------------core user----------------------------------------------------------------

func (r *user) Save(data *model_core.CoreUsers) error {
	err := r.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *user) UpdateUser(query map[string]interface{}, data *model_core.CoreUsers) error {
	err := r.db.Model(&model_core.CoreUsers{}).Where(query).Updates(data).Error
	return err
}
func (r *user) FindUser(query map[string]interface{}) (model_core.CoreUsers, error) {
	var data model_core.CoreUsers
	fmt.Println(query)
	err := r.db.Model(&model_core.CoreUsers{}).Preload(clause.Associations).Where(query).First(&data)
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *user) FindAll() ([]model_core.CoreUsers, error) {
	var datas []model_core.CoreUsers
	err := r.db.Find(&datas).Error
	if err != nil {
		return datas, err
	}
	return datas, nil
}
func (r *user) FindLogin(query map[string]interface{}) (model_core.CoreUsers, error) {
	var data model_core.CoreUsers
	// fmt.Println(query)
	err := r.db.Model(&model_core.CoreUsers{}).Where("email=? OR mobile_number=? OR username=?", query["login"], query["login"], query["login"]).First(&data)
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *user) UpdateAssignTemplate(id uint, data model_core.CoreUsers) (int64, error) {
	if data.ID == id {
		res := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id", id).Updates(&data)
		if res.Error != nil {
			return res.RowsAffected, res.Error
		}
		return res.RowsAffected, nil
	}
	return 0, fmt.Errorf("id not matched")

}

func (r *user) UpdateCompany(query map[string]interface{}, data model_core.Company) (model_core.Company, error) {
	err := r.db.Model(&model_core.Company{}).Where(query).Updates(&data).Error
	return data, err
}

func (r *user) GetCompany(query map[string]interface{}) (model_core.Company, error) {
	var data model_core.Company
	err := r.db.Model(&model_core.Company{}).Preload(clause.Associations).Where(query).First(&data).Error

	return data, err
}

func (r *user) FindAllCompanyUsers(query map[string]interface{}, p *pagination.Paginatevalue) ([]model_core.CoreUsers, error) {
	var data []model_core.CoreUsers
	err := r.db.Preload(clause.Associations).Model(&model_core.CoreUsers{}).Scopes(helpers.Paginate(&model_core.CoreUsers{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *user) UpdateUserDetails(id uint, data model_core.CoreUsers) (int64, error) {

	res := r.db.Model(&model_core.CoreUsers{}).Where("id", id).Updates(&data)
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}

	return res.RowsAffected, nil
}

func (r *user) FindAllCoreUser(query map[string]interface{}, p *pagination.Paginatevalue) ([]model_core.CoreUsers, error) {

	var result []model_core.CoreUsers
	res := r.db.Model(&model_core.CoreUsers{}).Scopes(helpers.Paginate(&model_core.CoreUsers{}, p, r.db)).Where("is_active = true")
	final_res := res.Find(&result)
	fmt.Println(final_res.Error)

	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (r *user) FindAllCompanies(page *pagination.Paginatevalue) ([]model_core.Company, error) {
	var data []model_core.Company
	err := r.db.Preload(clause.Associations).Model(&model_core.Company{}).Scopes(helpers.Paginate(&model_core.Company{}, page, r.db)).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
