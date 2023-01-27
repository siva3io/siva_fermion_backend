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
type PaymentTerm interface {
	CreatePaymentTerm(data *accounting.PaymentTerms) error
	UpdatePaymentTerm(query map[string]interface{}, data *accounting.PaymentTerms) error
	DeletePaymentTerm(query map[string]interface{}) error
	FindOnePaymentTerm(query map[string]interface{}) (accounting.PaymentTerms, error)
	FindAllPaymentTerm(query map[string]interface{}, p *pagination.Paginatevalue) ([]accounting.PaymentTerms, error)
}
type Payment_term struct {
	db *gorm.DB
}

var PaymentTermRepository *Payment_term //singleton object

// singleton function
func NewPaymentTerm() *Payment_term {
	if PaymentTermRepository != nil {
		return PaymentTermRepository
	}
	db := db.DbManager()
	PaymentTermRepository = &Payment_term{db}
	return PaymentTermRepository
}

func (r *Payment_term) CreatePaymentTerm(data *accounting.PaymentTerms) error {

	err := r.db.Model(&accounting.PaymentTerms{}).Create(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {

		return err.Error
	}
	return nil
}
func (r *Payment_term) UpdatePaymentTerm(query map[string]interface{}, data *accounting.PaymentTerms) error {

	err := r.db.Model(&accounting.PaymentTerms{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {

		return err.Error
	}
	return nil
}
func (r *Payment_term) DeletePaymentTerm(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(uint),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&accounting.PaymentTerms{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {

		return err.Error
	}
	return nil
}

func (r *Payment_term) FindOnePaymentTerm(query map[string]interface{}) (accounting.PaymentTerms, error) {
	var data accounting.PaymentTerms
	err := r.db.Preload(clause.Associations).Model(&accounting.PaymentTerms{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *Payment_term) FindAllPaymentTerm(query map[string]interface{}, p *pagination.Paginatevalue) ([]accounting.PaymentTerms, error) {
	var data []accounting.PaymentTerms
	err := r.db.Preload(clause.Associations).Model(&accounting.PaymentTerms{}).Scopes(helpers.Paginate(&accounting.PaymentTerms{}, p, r.db)).Where(query).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}
