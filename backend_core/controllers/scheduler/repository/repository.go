package scheduler

import (
	"errors"
	"os"
	"time"

	model "fermion/backend_core/controllers/scheduler/model"
	"fermion/backend_core/db"
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
type Scheduler interface {
	Save(data model.SchedulerJob) (model.SchedulerJob, error)
	FindAll(page *pagination.Paginatevalue, query map[string]interface{}) ([]model.SchedulerJob, error)
	SaveSchedulerLog(data model.SchedulerLog) (model.SchedulerLog, error)
	SaveSchedulerLogs(data []model.SchedulerLog) ([]model.SchedulerLog, error)
	FindAllSchedulerLog(page *pagination.Paginatevalue, query map[string]interface{}) ([]model.SchedulerLog, error)

	FindOneSchedulerJob(query map[string]interface{}) (model.SchedulerJob, error)
	UpdateSchedulerLog(query map[string]interface{}, data model.SchedulerLog) error
	DeleteSchedulerJob(query map[string]interface{}) error
	DeleteSchedulerLog(query map[string]interface{}) error
	UpdateSchedulerJob(query map[string]interface{}, data model.SchedulerJob) error
}

type scheduler struct {
	db *gorm.DB
}

func NewScheduler() *scheduler {
	db := db.DbManager()
	return &scheduler{db}

}

func (r *scheduler) Save(data model.SchedulerJob) (model.SchedulerJob, error) {

	err := r.db.Model(&model.SchedulerJob{}).Create(&data).Error

	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *scheduler) FindAll(page *pagination.Paginatevalue, query map[string]interface{}) ([]model.SchedulerJob, error) {
	var data []model.SchedulerJob

	if page == nil || page.Per_page == -1 {
		err := r.db.Model(&model.SchedulerJob{}).Where(query).Find(&data)
		if err.Error != nil {
			return nil, err.Error
		}
		return data, nil
	}

	err := r.db.Model(&model.SchedulerJob{}).Scopes(helpers.Paginate(&model.SchedulerJob{}, page, r.db)).Where(query).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *scheduler) SaveSchedulerLog(data model.SchedulerLog) (model.SchedulerLog, error) {
	err := r.db.Model(&model.SchedulerLog{}).Create(&data).Error

	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *scheduler) SaveSchedulerLogs(data []model.SchedulerLog) ([]model.SchedulerLog, error) {
	err := r.db.Model(&model.SchedulerLog{}).Create(&data).Error

	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *scheduler) FindAllSchedulerLog(page *pagination.Paginatevalue, query map[string]interface{}) ([]model.SchedulerLog, error) {
	var data []model.SchedulerLog

	if page == nil || page.Per_page == -1 {
		if query["state"] != nil && query["start_time"] != nil {
			err := r.db.Model(&model.SchedulerLog{}).Preload(clause.Associations).Where("state = ? AND start_time < ?", query["state"], query["start_time"]).Find(&data)
			if err.Error != nil {
				return nil, err.Error
			}
			return data, nil
		}

		err := r.db.Model(&model.SchedulerLog{}).Preload(clause.Associations).Where(query).Find(&data)

		if err.Error != nil {
			return nil, err.Error
		}
		return data, nil
	}
	err := r.db.Model(&model.SchedulerLog{}).Scopes(helpers.Paginate(&model.SchedulerLog{}, page, r.db)).Preload(clause.Associations).Where(query).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *scheduler) FindOneSchedulerJob(query map[string]interface{}) (model.SchedulerJob, error) {
	var data model.SchedulerJob

	err := r.db.Model(&model.SchedulerJob{}).Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {

		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *scheduler) UpdateSchedulerLog(query map[string]interface{}, data model.SchedulerLog) error {

	if query["state"] != nil && query["start_time"] != nil {
		res := r.db.Model(&model.SchedulerLog{}).Where("state = ? AND start_time < ?", query["state"], query["start_time"]).Updates(&data)

		if res.Error != nil {
			return res.Error
		}
		return nil
	}

	res := r.db.Model(&model.SchedulerLog{}).Where(query).Updates(&data)

	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *scheduler) DeleteSchedulerJob(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&model.SchedulerJob{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *scheduler) DeleteSchedulerLog(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)

	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")

	res := r.db.Model(&model.SchedulerLog{}).Where(query).Updates(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *scheduler) UpdateSchedulerJob(query map[string]interface{}, data model.SchedulerJob) error {
	err := r.db.Model(&model.SchedulerJob{}).Where(query).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}
