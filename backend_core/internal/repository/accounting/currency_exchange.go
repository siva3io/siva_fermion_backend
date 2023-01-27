package accounting

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/accounting"
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
type CurrencyExchange interface {
	CreateExchangePair(data *accounting.CurrencyExchange) error
	UpdateExchangePair(query map[string]interface{}, data *accounting.CurrencyExchange) error
	FindOne(query map[string]interface{}) (accounting.CurrencyExchange, error)
	FindAll(query map[string]interface{}, p *pagination.Paginatevalue) ([]accounting.CurrencyExchange, error)
	DeleteExchangePair(query map[string]interface{}) error
}

type currency_exchange struct {
	db *gorm.DB
}

var currencyExchangeRepository *currency_exchange //singleton object

// singleton function
func NewCurrency() *currency_exchange {
	if currencyExchangeRepository != nil {
		return currencyExchangeRepository
	}
	db := db.DbManager()
	currencyExchangeRepository = &currency_exchange{db}
	return currencyExchangeRepository
}

func (r *currency_exchange) CreateExchangePair(data *accounting.CurrencyExchange) error {
	err := r.db.Model(&accounting.CurrencyExchange{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *currency_exchange) UpdateExchangePair(query map[string]interface{}, data *accounting.CurrencyExchange) error {
	err := r.db.Preload(clause.Associations).Model(&accounting.CurrencyExchange{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *currency_exchange) FindAll(query map[string]interface{}, p *pagination.Paginatevalue) ([]accounting.CurrencyExchange, error) {
	var data []accounting.CurrencyExchange
	err := r.db.Preload(clause.Associations).Model(&accounting.CurrencyExchange{}).Scopes(helpers.Paginate(&accounting.CurrencyExchange{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *currency_exchange) FindOne(query map[string]interface{}) (accounting.CurrencyExchange, error) {
	var data accounting.CurrencyExchange
	err := r.db.Preload(clause.Associations).Model(&accounting.CurrencyExchange{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil

}
func (r *currency_exchange) DeleteExchangePair(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&accounting.CurrencyExchange{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}
