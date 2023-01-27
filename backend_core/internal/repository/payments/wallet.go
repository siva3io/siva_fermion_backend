package payments

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/payments"
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
type Wallet interface {
	Save(data *payments.Wallets) error
	FindAll(page *pagination.Paginatevalue) (interface{}, error)
	FindOne(query map[string]interface{}) (interface{}, error)
	Update(query map[string]interface{}, data *payments.Wallets) error
	Delete(query map[string]interface{}) error
}

type Wallets struct {
	db *gorm.DB
}

var WalletsRepository *Wallets //singleton object

// singleton function
func NewWallet() *Wallets {
	if WalletsRepository != nil {
		return WalletsRepository
	}
	db := db.DbManager()
	WalletsRepository = &Wallets{db}
	return WalletsRepository

}

func (r *Wallets) Save(data *payments.Wallets) error {
	err := r.db.Model(&payments.Wallets{}).Create(data).Error

	if err != nil {

		return err

	}

	return nil
}

func (r *Wallets) FindAll(page *pagination.Paginatevalue) (interface{}, error) {
	var data []payments.Wallets

	err := r.db.Model(&payments.Wallets{}).Scopes(helpers.Paginate(&payments.Wallets{}, page, r.db)).Preload(clause.Associations).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *Wallets) FindOne(query map[string]interface{}) (interface{}, error) {
	var data payments.Wallets

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *Wallets) Update(query map[string]interface{}, data *payments.Wallets) error {
	res := r.db.Model(&payments.Wallets{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *Wallets) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["created_by"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&payments.Wallets{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}
