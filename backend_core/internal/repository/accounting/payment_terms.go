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
type PaymentTerm interface {
	CreatePaymentTerm(data *accounting.PaymentTerms) error
	UpdatePaymentTerm(query map[string]interface{}, data *accounting.PaymentTerms) error
	DeletePaymentTerm(query map[string]interface{}) error
	FindOnePaymentTerm(query map[string]interface{}) (accounting.PaymentTerms, error)
	FindAllPaymentTerm(query interface{}, p *pagination.Paginatevalue) ([]accounting.PaymentTerms, error)
}
type Payment_term struct {
	db *gorm.DB
}

func NewPaymentTerm() *Payment_term {
	db := db.DbManager()
	return &Payment_term{db}
}

func (r *Payment_term) CreatePaymentTerm(data *accounting.PaymentTerms) error {

	err := r.db.Model(&accounting.PaymentTerms{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Payment_term) UpdatePaymentTerm(query map[string]interface{}, data *accounting.PaymentTerms) error {

	err := r.db.Model(&accounting.PaymentTerms{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Payment_term) DeletePaymentTerm(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&accounting.PaymentTerms{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *Payment_term) FindOnePaymentTerm(query map[string]interface{}) (accounting.PaymentTerms, error) {
	var data accounting.PaymentTerms
	err := r.db.Preload(clause.Associations).Model(&accounting.PaymentTerms{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *Payment_term) FindAllPaymentTerm(query interface{}, p *pagination.Paginatevalue) ([]accounting.PaymentTerms, error) {
	var data []accounting.PaymentTerms
	err := r.db.Preload(clause.Associations).Model(&accounting.PaymentTerms{}).Scopes(helpers.Paginate(&accounting.PaymentTerms{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
