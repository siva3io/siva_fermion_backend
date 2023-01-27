package mdm

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/core"
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
type Notifications interface {
	CreateNotification(data *core.Notifications) error
	UpdateNotification(data *core.Notifications, query map[string]interface{}) error
	DeleteNotification(query map[string]interface{}) error
	FindNotification(query map[string]interface{}) (core.Notifications, error)
	GetAllNotification(p *pagination.Paginatevalue) ([]core.Notifications, error)
}
type notifications struct {
	db *gorm.DB
}

func NewNotifications() *notifications {
	db := db.DbManager()
	return &notifications{db}
}

// -------------------------Product Template---------------------------------------------------------------------------------------------------
func (r *notifications) CreateNotification(data *core.Notifications) error {
	err := r.db.Model(&core.Notifications{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *notifications) UpdateNotification(data *core.Notifications, query map[string]interface{}) error {
	err := r.db.Model(&core.Notifications{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *notifications) DeleteNotification(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&core.Notifications{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *notifications) FindNotification(query map[string]interface{}) (core.Notifications, error) {
	var result core.Notifications
	err := r.db.Model(&core.Notifications{}).Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations).Where(query).First(&result).Error
	// result.StockTreatmentIds = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(result.StockTreatmentIds)
	// result.ProductProcurementTreatmentIds = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(result.ProductProcurementTreatmentIds)
	return result, err
}
func (r *notifications) GetAllNotification(p *pagination.Paginatevalue) ([]core.Notifications, error) {
	var result []core.Notifications
	err := r.db.Model(&core.Notifications{}).Preload(clause.Associations + "." + clause.Associations).Scopes(helpers.Paginate(&core.Notifications{}, p, r.db)).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
