package scheduler

import (
	"encoding/json"
	"errors"
	"time"

	model "fermion/backend_core/controllers/scheduler/model"
	repo "fermion/backend_core/controllers/scheduler/repository"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

	"github.com/robfig/cron/v3"
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

var specParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type Service interface {
	AddSchedulerJob(data SchedulerJobDto) error
	GetAllSchedulerJob(page *pagination.Paginatevalue) ([]SchedulerJobResponseDTO, error)
	GetOneSchedulerJob(query map[string]interface{}) (SchedulerJobResponseDTO, error)
	ListSchedulerLogs(page *pagination.Paginatevalue) ([]SchedulerLogResponseDTO, error)
	DeleteSchedulerJob(query map[string]interface{}) error
	UpdateSchedulerJob(query map[string]interface{}, data *SchedulerJobDto) error
}

type service struct {
	schedulerRepository repo.Scheduler
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	schedulerRepository := repo.NewScheduler()
	newServiceObj = &service{schedulerRepository}
	return newServiceObj
}

func (s *service) AddSchedulerJob(data SchedulerJobDto) error {

	var schedulerJob model.SchedulerJob

	dto, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(dto, &schedulerJob); err != nil {
		return err
	}

	if schedulerJob.State == nil {
		state := false
		schedulerJob.State = &state
	}

	if schedulerJob.StartTime.IsZero() {
		schedulerJob.StartTime = time.Now()
	}

	if schedulerJob.StartTime.Before(time.Now().Add(-10 * time.Second)) {
		return errors.New("invalid start time")
	}

	schedulerJob, err = s.schedulerRepository.Save(schedulerJob)

	if err != nil {
		return err
	}

	if schedulerJob.Frequency == "@reboot" {
		return nil
	}
	startTime := schedulerJob.StartTime

	if schedulerJob.Frequency != "@once" {
		sched, err := specParser.Parse(schedulerJob.Frequency)
		if err != nil {
			return err
		}
		startTime = sched.Next(schedulerJob.StartTime)
	}

	if !*schedulerJob.State {
		return nil
	}

	log := model.SchedulerLog{
		SchedulerJobId: schedulerJob.ID,
		State:          "queued",
		StartTime:      startTime,
	}

	log, err = s.schedulerRepository.SaveSchedulerLog(log)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAllSchedulerJob(page *pagination.Paginatevalue) ([]SchedulerJobResponseDTO, error) {
	var res []SchedulerJobResponseDTO
	result, _ := s.schedulerRepository.FindAll(page, nil)
	sdata, _ := json.Marshal(result)
	json.Unmarshal(sdata, &res)
	return res, nil
}

func (s *service) GetOneSchedulerJob(query map[string]interface{}) (SchedulerJobResponseDTO, error) {
	var sdata SchedulerJobResponseDTO
	result, _ := s.schedulerRepository.FindOneSchedulerJob(query)
	kdata, _ := json.Marshal(result)
	json.Unmarshal(kdata, &sdata)
	return sdata, nil
}
func (s *service) DeleteSchedulerJob(query map[string]interface{}) error {

	updatedQuery := map[string]interface{}{
		"id": query["id"],
	}
	_, err := s.schedulerRepository.FindOneSchedulerJob(updatedQuery)
	if err != nil {
		return err
	}
	updatedQuery = map[string]interface{}{
		"id":      query["id"],
		"user_id": query["user_id"],
	}
	err = s.schedulerRepository.DeleteSchedulerJob(updatedQuery)
	if err != nil {
		return err
	}
	updatedQuery = map[string]interface{}{
		"scheduler_job_id": query["id"],
		"state":            "queued",
		"user_id":          query["user_id"],
	}
	err = s.schedulerRepository.DeleteSchedulerLog(updatedQuery)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateSchedulerJob(query map[string]interface{}, dtoData *SchedulerJobDto) error {
	var data model.SchedulerJob
	jobQuery := map[string]interface{}{
		"id": query["id"],
	}
	schedulerJob, err := s.schedulerRepository.FindOneSchedulerJob(jobQuery)
	if err != nil {
		return err
	}
	if (dtoData.Frequency != schedulerJob.Frequency) && (dtoData.Frequency == "@reboot" || dtoData.Frequency == "@once" || schedulerJob.Frequency == "@reboot" || schedulerJob.Frequency == "@once") {
		return errors.New("cannot update frequency of @once or @reboot schedulers")
	}

	err = helpers.JsonMarshaller(*dtoData, &data)
	if err != nil {
		return err
	}
	err = s.schedulerRepository.UpdateSchedulerJob(jobQuery, data)
	if err != nil {
		return err
	}
	frequency := schedulerJob.Frequency
	if data.Frequency != "" {
		frequency = data.Frequency
	}
	if data.State != nil && *data.State != *schedulerJob.State {
		if !*data.State {
			logQuery := map[string]interface{}{
				"scheduler_job_id": query["id"],
				"state":            "queued",
				"user_id":          query["user_id"],
			}
			err = s.schedulerRepository.DeleteSchedulerLog(logQuery)
			if err != nil {
				return err
			}
		} else {
			startTime := schedulerJob.StartTime
			if schedulerJob.Frequency != "@once" {
				sched, err := specParser.Parse(frequency)
				if err != nil {
					return err
				}
				startTime = sched.Next(schedulerJob.StartTime)
			}
			log := model.SchedulerLog{
				SchedulerJobId: schedulerJob.ID,
				State:          "queued",
				StartTime:      startTime,
			}
			_, err = s.schedulerRepository.SaveSchedulerLog(log)
			if err != nil {
				return err
			}
		}
	}
	if data.Frequency != "" && data.Frequency != schedulerJob.Frequency {
		logQuery := map[string]interface{}{
			"scheduler_job_id": query["id"],
			"state":            "queued",
			"user_id":          query["user_id"],
		}
		err = s.schedulerRepository.DeleteSchedulerLog(logQuery)
		if err != nil {
			return err
		}
		startTime := schedulerJob.StartTime
		if schedulerJob.Frequency != "@once" {
			sched, err := specParser.Parse(frequency)
			if err != nil {
				return err
			}
			startTime = sched.Next(schedulerJob.StartTime)
		}
		log := model.SchedulerLog{
			SchedulerJobId: schedulerJob.ID,
			State:          "queued",
			StartTime:      startTime,
		}
		_, err = s.schedulerRepository.SaveSchedulerLog(log)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *service) ListSchedulerLogs(page *pagination.Paginatevalue) ([]SchedulerLogResponseDTO, error) {
	var sdata []SchedulerLogResponseDTO
	result, _ := s.schedulerRepository.FindAllSchedulerLog(page, nil)
	kdata, _ := json.Marshal(result)
	err := json.Unmarshal(kdata, &sdata)
	if err != nil {
		return nil, err
	}
	return sdata, nil
}
