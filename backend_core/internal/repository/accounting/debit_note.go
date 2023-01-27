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
type DebitNotes interface {
	SaveDebitNote(data *accounting.DebitNote) error
	UpdateDebitNote(query map[string]interface{}, data *accounting.DebitNote) error
	DeleteDebitNote(query map[string]interface{}) error
	FindOneDebitNote(query map[string]interface{}) (accounting.DebitNote, error)
	FindAllDebitNote(query interface{}, p *pagination.Paginatevalue) ([]accounting.DebitNote, error)
	//SearchDebitNote(query string) ([]accounting.DebitNote, error)
	SaveDebitLines(data accounting.DebitNoteLineItems) error
	UpdateDebitLines(query map[string]interface{}, data accounting.DebitNoteLineItems) (int64, error)
	DeleteDebitLine(query map[string]interface{}) error
	FindDebitLines(query map[string]interface{}) (accounting.DebitNoteLineItems, error)
}
type debit_note struct {
	db *gorm.DB
}

var debitNoteRepository *debit_note //singleton object

// singleton function
func NewDebitNote() *debit_note {
	if debitNoteRepository != nil {
		return debitNoteRepository
	}
	db := db.DbManager()
	debitNoteRepository = &debit_note{db}
	return debitNoteRepository
}

func (r *debit_note) SaveDebitNote(data *accounting.DebitNote) error {
	err := r.db.Model(&accounting.DebitNote{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *debit_note) UpdateDebitNote(query map[string]interface{}, data *accounting.DebitNote) error {

	err := r.db.Model(&accounting.DebitNote{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil
}
func (r *debit_note) DeleteDebitNote(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&accounting.DebitNote{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *debit_note) FindOneDebitNote(query map[string]interface{}) (accounting.DebitNote, error) {
	var data accounting.DebitNote
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Model(&accounting.DebitNote{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *debit_note) FindAllDebitNote(query interface{}, p *pagination.Paginatevalue) ([]accounting.DebitNote, error) {
	var data []accounting.DebitNote
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Model(&accounting.DebitNote{}).Scopes(helpers.Paginate(&accounting.DebitNote{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// func (r *debit_note) SearchDebitNote(query string) ([]accounting.DebitNote, error) {
// 	var data []accounting.DebitNote
// 	err := r.db.Preload(clause.Associations).Find(&data, "name ILIKE ? OR primary_email ILIKE ?", "%"+query+"%", "%"+query+"%").Error
// 	if err != nil {
// 		return data, err
// 	}
// 	return data, nil
// }

func (r *debit_note) SaveDebitLines(data accounting.DebitNoteLineItems) error {
	res := r.db.Model(&accounting.DebitNoteLineItems{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *debit_note) FindDebitLines(query map[string]interface{}) (accounting.DebitNoteLineItems, error) {
	var result accounting.DebitNoteLineItems
	fmt.Println(query)
	res := r.db.Model(&accounting.DebitNoteLineItems{}).Where(query).First(&result)
	if res.RowsAffected == 0 {
		return result, errors.New("oops! record not found")
	}
	if res.Error != nil {
		return result, res.Error
	}
	return result, nil
}

func (r *debit_note) UpdateDebitLines(query map[string]interface{}, data accounting.DebitNoteLineItems) (int64, error) {
	res := r.db.Model(&accounting.DebitNoteLineItems{}).Where(query).Updates(&data)
	if res.Error != nil {

		return res.RowsAffected, res.Error

	}
	return res.RowsAffected, nil
}

func (r *debit_note) DeleteDebitLine(query map[string]interface{}) error {
	res := r.db.Model(&accounting.DebitNoteLineItems{}).Where(query).Delete(&accounting.DebitNoteLineItems{})
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}
