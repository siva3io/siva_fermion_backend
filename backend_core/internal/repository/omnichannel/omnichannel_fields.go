package omnichannel

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/omnichannel"
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
type OmnichannelField interface {
	Save(data *omnichannel.OmnichannelField) (uint, error)
	FindAll(page *pagination.Paginatevalue) (interface{}, error)
	FindOne(query map[string]interface{}) (interface{}, error)
	Update(query map[string]interface{}, data *omnichannel.OmnichannelField) error
	Delete(query map[string]interface{}) error

	UpsertOmnichannelFieldsData(data *[]omnichannel.OmnichannelFieldData) error
	GetOmnichannelFieldsData(query map[string]interface{}) (interface{}, error)
	FindAllFieldData(page *pagination.Paginatevalue) (interface{}, error)
}

type omnichannelField struct {
	db *gorm.DB
}

var omnichannelFieldRepository *omnichannelField //singleton object

// singleton function
func NewOmnichannelFieldRepo() *omnichannelField {
	if omnichannelFieldRepository != nil {
		return omnichannelFieldRepository
	}
	db := db.DbManager()
	omnichannelFieldRepository = &omnichannelField{db}
	return omnichannelFieldRepository

}

func (r *omnichannelField) Save(data *omnichannel.OmnichannelField) (uint, error) {
	err := r.db.Model(&omnichannel.OmnichannelField{}).Create(data).Error
	if err != nil {
		return 0, err

	}
	return data.ID, err
}

func (r *omnichannelField) FindAll(page *pagination.Paginatevalue) (interface{}, error) {
	var data []omnichannel.OmnichannelField
	err := r.db.Model(&omnichannel.OmnichannelField{}).Scopes(helpers.Paginate(&omnichannel.OmnichannelField{}, page, r.db)).Preload(clause.Associations).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *omnichannelField) FindOne(query map[string]interface{}) (interface{}, error) {

	var data omnichannel.OmnichannelField
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *omnichannelField) Update(query map[string]interface{}, data *omnichannel.OmnichannelField) error {
	res := r.db.Model(&omnichannel.OmnichannelField{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *omnichannelField) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&omnichannel.OmnichannelField{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *omnichannelField) UpsertOmnichannelFieldsData(data *[]omnichannel.OmnichannelFieldData) error {

	err := r.db.Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "field_app_id"}, {Name: "omnichannel_field_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"data"}),
	}).Create(data).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *omnichannelField) GetOmnichannelFieldsData(query map[string]interface{}) (interface{}, error) {
	var data []omnichannel.OmnichannelFieldData

	err := r.db.Debug().Model(&omnichannel.OmnichannelFieldData{}).Preload(clause.Associations + "." + clause.Associations).Joins("INNER JOIN omnichannel_fields ON omnichannel_fields.id = omnichannel_field_data.omnichannel_field_id").Where(query).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *omnichannelField) FindAllFieldData(page *pagination.Paginatevalue) (interface{}, error) {
	var data []omnichannel.OmnichannelFieldData
	err := r.db.Debug().Model(&omnichannel.OmnichannelFieldData{}).Scopes(helpers.Paginate(&omnichannel.OmnichannelFieldData{}, page, r.db)).Preload(clause.Associations).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}
