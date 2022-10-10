package accounting

import (
	"errors"
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/accounting"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
type CurrencyExchange interface {
	CreateExchangePair(data *accounting.CurrencyExchange) error
	UpdateExchangePair(query map[string]interface{}, data *accounting.CurrencyExchange) error
	FindOne(query map[string]interface{}) (accounting.CurrencyExchange, error)
	FindAll() ([]accounting.CurrencyExchange, error)
	DeleteExchangePair(query map[string]interface{}) error
}

type currency_exchange struct {
	db *gorm.DB
}

func NewCurrency() *currency_exchange {
	db := db.DbManager()
	return &currency_exchange{db}
}

func (r *currency_exchange) CreateExchangePair(data *accounting.CurrencyExchange) error {
	err := r.db.Model(&accounting.CurrencyExchange{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *currency_exchange) UpdateExchangePair(query map[string]interface{}, data *accounting.CurrencyExchange) error {
	err := r.db.Preload(clause.Associations).Model(&accounting.CurrencyExchange{}).Where(query).Updates(data).Error
	fmt.Println("--------------------------", err)
	if err != nil {
		return err
	}
	return nil
}

func (r *currency_exchange) FindAll() ([]accounting.CurrencyExchange, error) {
	// var data []accounting.BasicPairInfo
	var data []accounting.CurrencyExchange
	err := r.db.Preload(clause.Associations).Model(&accounting.CurrencyExchange{}).Find(&data).Error
	fmt.Println("//////////////////////////////////////", data)
	if err != nil {
		return nil, err
	}
	fmt.Println("||||||||||||||||||||||||||||||||||||||||||||||", data)
	return data, nil
}

func (r *currency_exchange) FindOne(query map[string]interface{}) (accounting.CurrencyExchange, error) {
	var data accounting.CurrencyExchange
	err := r.db.Preload(clause.Associations).Model(&accounting.CurrencyExchange{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
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
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&accounting.CurrencyExchange{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
