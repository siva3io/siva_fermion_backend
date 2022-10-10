package mdm

import (
	"errors"
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/mdm"
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
type Contacts interface {
	SaveContact(data *mdm.Partner) error
	UpdateContact(query map[string]interface{}, data *mdm.Partner) error
	DeleteContact(query map[string]interface{}) error
	FindOneContact(query map[string]interface{}) (mdm.Partner, error)
	FindAllContact(query interface{}, p *pagination.Paginatevalue) ([]mdm.Partner, error)
	SearchContact(query string) ([]mdm.Partner, error)
}

type contacts struct {
	db *gorm.DB
}

func NewContacts() *contacts {
	db := db.DbManager()
	return &contacts{db}

}
func (r *contacts) SaveContact(data *mdm.Partner) error {

	err := r.db.Model(&mdm.Partner{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *contacts) UpdateContact(query map[string]interface{}, data *mdm.Partner) error {

	err := r.db.Model(&mdm.Partner{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *contacts) DeleteContact(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)

	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")

	response, _ := r.FindOneContact(query)

	data["primary_email"] = response.PrimaryEmail + fmt.Sprintf("_%v", time.Now().UnixMilli())
	err := r.db.Model(&mdm.Partner{}).Where(query).Updates(data).Error //As it is a soft delete, we use updates inorder to update deleted_by too
	if err != nil {
		return err
	}
	return nil
}
func (r *contacts) FindOneContact(query map[string]interface{}) (mdm.Partner, error) {
	var data mdm.Partner
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Model(&mdm.Partner{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	data.Properties = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(data.Properties)
	return data, nil
}
func (r *contacts) FindAllContact(query interface{}, p *pagination.Paginatevalue) ([]mdm.Partner, error) {
	var data []mdm.Partner
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Model(&mdm.Partner{}).Scopes(helpers.Paginate(&mdm.Partner{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *contacts) SearchContact(query string) ([]mdm.Partner, error) {
	var data []mdm.Partner
	err := r.db.Preload(clause.Associations).Find(&data, "name ILIKE ? OR primary_email ILIKE ?", "%"+query+"%", "%"+query+"%").Error
	if err != nil {
		return data, err
	}
	return data, nil
}
