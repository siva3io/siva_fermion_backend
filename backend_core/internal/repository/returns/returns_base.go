package returns

import (
	"errors"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/returns"

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
type ReturnsBase interface {
	CreateRecord(data returns.UserReturnsFav) error
	Favourite(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (returns.UserReturnsFav, error)
	UnFavourite(query map[string]interface{}) error
}

type returnsBase struct {
	db *gorm.DB
}

var returnsBaseRepository *returnsBase //singleton object

// singleton function
func NewReturnsBase() *returnsBase {
	if returnsBaseRepository != nil {
		return returnsBaseRepository
	}
	db := db.DbManager()
	returnsBaseRepository = &returnsBase{db}
	return returnsBaseRepository

}

func (r *returnsBase) CreateRecord(data returns.UserReturnsFav) error {
	err := r.db.Model(&returns.UserReturnsFav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *returnsBase) Favourite(query map[string]interface{}) error {
	err := r.db.Model(&returns.UserReturnsFav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *returnsBase) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&returns.UserReturnsFav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *returnsBase) FindOne(query map[string]interface{}) (returns.UserReturnsFav, error) {
	var data returns.UserReturnsFav

	err := r.db.Preload(clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}
