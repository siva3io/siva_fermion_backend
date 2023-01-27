package mdm

import (
	"errors"
	"fmt"
	"os"
	"time"

	db "fermion/backend_core/db"
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
type Contacts interface {
	Create(data *mdm.Partner) error
	Update(query map[string]interface{}, data *mdm.Partner) error
	Delete(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (mdm.Partner, error)
	FindAll(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.Partner, error)
}

type contacts struct {
	Db *gorm.DB
}

var contactsRepository *contacts //singleton object

// singleton function
func NewContacts() *contacts {
	if contactsRepository != nil {
		return contactsRepository
	}
	db := db.DbManager()
	contactsRepository = &contacts{db}
	return contactsRepository

}
func (r *contacts) Create(data *mdm.Partner) error {

	err := r.Db.Model(&mdm.Partner{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *contacts) Update(query map[string]interface{}, data *mdm.Partner) error {

	err := r.Db.Model(&mdm.Partner{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *contacts) Delete(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)

	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")

	contactDetails, _ := r.FindOne(query)
	data["primary_email"] = contactDetails.PrimaryEmail + fmt.Sprintf("_%v", time.Now().UnixMilli())

	err := r.Db.Model(&mdm.Partner{}).Where(query).Updates(data) //As it is a soft delete, we use updates inorder to update deleted_by too
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *contacts) FindOne(query map[string]interface{}) (mdm.Partner, error) {
	var data mdm.Partner

	err := r.Db.Preload(clause.Associations + "." + clause.Associations).Model(&mdm.Partner{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}

	return data, nil
}
func (r *contacts) FindAll(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.Partner, error) {
	var data []mdm.Partner
	err := r.Db.Preload(clause.Associations + "." + clause.Associations).Model(&mdm.Partner{}).Scopes(helpers.Paginate(&mdm.Partner{}, p, r.Db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
