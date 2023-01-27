package accounting

import (
	"errors"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/accounting"

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
type AccountingBase interface {
	CreateRecord(data accounting.AccountingFav) error
	Favourite(query map[string]interface{}) error
	UnFavourite(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (accounting.AccountingFav, error)
}

type accountingBase struct {
	db *gorm.DB
}

var accountingBaseRepository *accountingBase //singleton object

// singleton function
func NewAccountingBase() *accountingBase {
	if accountingBaseRepository != nil {
		return accountingBaseRepository
	}
	db := db.DbManager()
	accountingBaseRepository = &accountingBase{db}
	return accountingBaseRepository

}

func (r *accountingBase) CreateRecord(data accounting.AccountingFav) error {
	err := r.db.Model(&accounting.AccountingFav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *accountingBase) Favourite(query map[string]interface{}) error {
	err := r.db.Model(&accounting.AccountingFav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *accountingBase) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&accounting.AccountingFav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *accountingBase) FindOne(query map[string]interface{}) (accounting.AccountingFav, error) {
	var data accounting.AccountingFav
	err := r.db.Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	return data, nil
}
