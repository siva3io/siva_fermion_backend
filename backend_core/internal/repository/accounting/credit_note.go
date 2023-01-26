package accounting

import (
	"errors"
	"fmt"
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
type CreditNotes interface {
	SaveCreditNote(data *accounting.CreditNote) error
	UpdateCreditNote(query map[string]interface{}, data *accounting.CreditNote) error
	DeleteCreditNote(query map[string]interface{}) error
	FindOneCreditNote(query map[string]interface{}) (accounting.CreditNote, error)
	FindAllCreditNote(query interface{}, p *pagination.Paginatevalue) ([]accounting.CreditNote, error)
	//SearchCreditNote(query string) ([]accounting.CreditNote, error)

	SaveCreditLines(accounting.CreditNoteLineItems) error
	UpdateCreditLines(map[string]interface{}, accounting.CreditNoteLineItems) (int64, error)
	DeleteCreditLine(map[string]interface{}) error
	FindCreditLines(map[string]interface{}) (accounting.CreditNoteLineItems, error)
}
type credit_note struct {
	db *gorm.DB
}

func NewCreditNote() *credit_note {
	db := db.DbManager()
	return &credit_note{db}
}

func (r *credit_note) SaveCreditNote(data *accounting.CreditNote) error {

	// err := tx.Model(&accounting.CreditNote{}).Create(data).Error
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err := r.db.Model(&accounting.CreditNote{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *credit_note) UpdateCreditNote(query map[string]interface{}, data *accounting.CreditNote) error {

	err := r.db.Model(&accounting.CreditNote{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *credit_note) DeleteCreditNote(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&accounting.CreditNote{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *credit_note) FindOneCreditNote(query map[string]interface{}) (accounting.CreditNote, error) {
	var data accounting.CreditNote
	err := r.db.Preload(clause.Associations).Model(&accounting.CreditNote{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *credit_note) FindAllCreditNote(query interface{}, p *pagination.Paginatevalue) ([]accounting.CreditNote, error) {
	var data []accounting.CreditNote
	err := r.db.Preload(clause.Associations).Model(&accounting.CreditNote{}).Scopes(helpers.Paginate(&accounting.CreditNote{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// func (r *credit_note) SearchCreditNote(query string) ([]accounting.CreditNote, error) {
// 	var data []accounting.CreditNote
// 	err := r.db.Preload(clause.Associations).Find(&data, "name ILIKE ? OR primary_email ILIKE ?", "%"+query+"%", "%"+query+"%").Error
// 	if err != nil {
// 		return data, err
// 	}
// 	return data, nil
// }

func (r *credit_note) SaveCreditLines(data accounting.CreditNoteLineItems) error {

	res := r.db.Model(&accounting.CreditNoteLineItems{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *credit_note) FindCreditLines(query map[string]interface{}) (accounting.CreditNoteLineItems, error) {
	var result accounting.CreditNoteLineItems
	fmt.Println(query)
	res := r.db.Model(&accounting.CreditNoteLineItems{}).Where(query).First(&result)

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *credit_note) UpdateCreditLines(query map[string]interface{}, data accounting.CreditNoteLineItems) (int64, error) {
	res := r.db.Model(&accounting.CreditNoteLineItems{}).Where(query).Updates(&data)
	fmt.Println(data)
	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}

func (r *credit_note) DeleteCreditLine(query map[string]interface{}) error {
	res := r.db.Model(&accounting.CreditNoteLineItems{}).Where(query).Delete(&accounting.CreditNoteLineItems{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}
