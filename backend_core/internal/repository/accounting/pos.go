package accounting

import (
	"errors"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/accounting"

	"gorm.io/gorm"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type Pos interface {
	CreatePos(data *accounting.UserPosLink) error
	UpdatePos(query map[string]interface{}, data *accounting.UserPosLink) error
	FindOne(query map[string]interface{}) (accounting.UserPosLink, error)
}

type PosRepo struct {
	Db *gorm.DB
}

func NewPosRepo() *PosRepo {
	db := db.DbManager()
	return &PosRepo{db}
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
