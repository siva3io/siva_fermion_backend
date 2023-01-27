package accounting

import (
	"errors"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/mdm"
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
type Pos interface {
	CreatePos(data *accounting.UserPosLink) error
	UpdatePos(query map[string]interface{}, data *accounting.UserPosLink) error
	FindOne(query map[string]interface{}) (accounting.UserPosLink, error)
	FindAllPosLink(query interface{}, p *pagination.Paginatevalue) ([]accounting.UserPosLink, error)
}

type PosRepo struct {
	Db *gorm.DB
}

var PosRepository *PosRepo //singleton object

// singleton function
func NewPosRepo() *PosRepo {
	if PosRepository != nil {
		return PosRepository
	}
	db := db.DbManager()
	PosRepository = &PosRepo{db}
	return PosRepository
}

func (r *PosRepo) CreatePos(data *accounting.UserPosLink) error {
	res := r.Db.Model(&accounting.UserPosLink{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *PosRepo) UpdatePos(query map[string]interface{}, data *accounting.UserPosLink) error {

	res := r.Db.Model(&accounting.UserPosLink{}).Where(query).Updates(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *PosRepo) FindOne(query map[string]interface{}) (accounting.UserPosLink, error) {
	var data accounting.UserPosLink
	err := r.Db.Model(&accounting.UserPosLink{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *PosRepo) FindAllPosLink(query interface{}, p *pagination.Paginatevalue) ([]accounting.UserPosLink, error) {
	var data []accounting.UserPosLink
	err := r.Db.Preload(clause.Associations).Model(&accounting.UserPosLink{}).Scopes(helpers.Paginate(&mdm.Partner{}, p, r.Db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
