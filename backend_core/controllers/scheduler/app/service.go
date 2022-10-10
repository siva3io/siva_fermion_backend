package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	model "fermion/backend_core/controllers/scheduler/model"
	repo "fermion/backend_core/controllers/scheduler/repository"
	utils "fermion/backend_core/controllers/scheduler/utils"

	"github.com/robfig/cron/v3"
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
var cr = cron.New()

var specParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type Service interface {
	AddSchedulerJob(data SchedulerJobDto) error
	InitSchedulerJob() error
	GetAllSchedulerJob() ([]SchedulerJobResponseDTO, error)
	GetOneSchedulerJob(query map[string]interface{}) (SchedulerJobResponseDTO, error)
	ListSchedulerLogs(query map[string]interface{}) ([]SchedulerLogResponseDTO, error)
	DeleteSchedulerJob(query map[string]interface{}) error
	UpdateSchedulerJob(query map[string]interface{}, data *SchedulerJobDto) error

	// GetQueuedSchedulerLogs() ([]model.SchedulerLog, error)
	InitJobsAtReboot() error
}

type service struct {
	schedulerRepository repo.Scheduler
}

func NewService() *service {
	schedulerRepository := repo.NewScheduler()
	return &service{schedulerRepository}
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

func (s *service) InitSchedulerJob() error {

	schedulerFrequencyInMinute := 5
	count := 1

	var query map[string]interface{}
	schedulerLogs, err := s.schedulerRepository.FindAllSchedulerLog(query)
	if err != nil {
		panic(err)
	}
	if len(schedulerLogs) == 0 {
		err = s.InitJobsAtReboot()
		if err != nil {
			panic(err)
		}
	}

	entryId, _ := cr.AddFunc(fmt.Sprintf("@every %vm", schedulerFrequencyInMinute), func() {

		var wg sync.WaitGroup

		fmt.Println("-------------------Scheduler started-------------------", count)

		query := map[string]interface{}{
			"state":      "queued",
			"start_time": time.Now(),
		}

		schedulerLogs, err := s.schedulerRepository.FindAllSchedulerLog(query)

		if err != nil {
			panic(err)
		}

		err = s.schedulerRepository.UpdateSchedulerLog(query, model.SchedulerLog{State: "in_progress"})

		if err != nil {
			panic(err)
		}

		wg.Add(len(schedulerLogs))

		for index, schedulerLog := range schedulerLogs {

			go func(index int, schedulerLog model.SchedulerLog) {

				defer wg.Done()

				schedulerJob := schedulerLog.SchedulerJob

				fmt.Printf("------------JobId:%v JobName:%v------------\n", schedulerJob.ID, schedulerJob.Name)

				if schedulerJob.Frequency != "@once" && schedulerJob.Frequency != "@reboot" {

					sched, err := specParser.Parse(schedulerJob.Frequency)

					if err != nil {
						panic(err)
					}

					nextLog := model.SchedulerLog{
						SchedulerJobId: schedulerJob.ID,
						State:          "queued",
						StartTime:      sched.Next(time.Now()),
					}

					nextLog, err = s.schedulerRepository.SaveSchedulerLog(nextLog)

					if err != nil {
						panic(err)
					}
				}

				returnValue, err := utils.CallFuncByName(&utils.CronFunctions{}, schedulerJob.Function, schedulerJob.Params)

				if err != nil {
					panic(err)
				}

				responseLog := returnValue[0].Interface().(model.SchedulerLog)

				schedulerLog.State = responseLog.State
				schedulerLog.Notes = responseLog.Notes
				schedulerLog.EndTime = time.Now()

				schedulerLogUpdateQuery := map[string]interface{}{
					"id": schedulerLog.ID,
				}

				err = s.schedulerRepository.UpdateSchedulerLog(schedulerLogUpdateQuery, schedulerLog)

				if err != nil {
					panic(err)
				}

			}(index, schedulerLog)
		}

		wg.Wait()
		fmt.Println("-------------------Scheduler ended-------------------", count)

		count++
	})
	fmt.Println("Scheduler Initialized! Cron Id: ", entryId)

	cr.Start()

	return nil
}
func (s *service) GetAllSchedulerJob() ([]SchedulerJobResponseDTO, error) {
	var res []SchedulerJobResponseDTO
	var query map[string]interface{}
	result, _ := s.schedulerRepository.FindAll(query)
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
func (s *service) UpdateSchedulerJob(query map[string]interface{}, data *SchedulerJobDto) error {
	var res model.SchedulerJob
	_, _ = s.schedulerRepository.FindOneSchedulerJob(query)

	sdata, _ := json.Marshal(data)
	json.Unmarshal(sdata, &res)
	er := s.schedulerRepository.UpdateSchedulerJob(query, res)
	if er != nil {
		return er
	}
	return nil
}
func (s *service) ListSchedulerLogs(query map[string]interface{}) ([]SchedulerLogResponseDTO, error) {
	var sdata []SchedulerLogResponseDTO
	result, _ := s.schedulerRepository.FindAllSchedulerLog(query)
	kdata, _ := json.Marshal(result)
	err := json.Unmarshal(kdata, &sdata)
	if err != nil {
		return nil, err
	}
	return sdata, nil
}

func (s *service) InitJobsAtReboot() error {

	var query map[string]interface{}

	fmt.Println("Rebooting Scheduler...!")

	rebootSchedulerJobs, err := s.schedulerRepository.FindAll(query)
	if err != nil {
		return err
	}

	lenRebootSchedulerJobs := len(rebootSchedulerJobs)

	if lenRebootSchedulerJobs > 0 {
		rebootSchedulerLogs := make([]model.SchedulerLog, lenRebootSchedulerJobs)
		for index, rebootSchedulerJob := range rebootSchedulerJobs {

			if rebootSchedulerJob.Frequency == "@once" {
				continue
			}

			startTime := time.Now()
			if rebootSchedulerJob.Frequency != "@reboot" {
				sched, err := specParser.Parse(rebootSchedulerJob.Frequency)
				if err != nil {
					return err
				}
				startTime = sched.Next(rebootSchedulerJob.StartTime)
			}

			schedulerLog := model.SchedulerLog{
				SchedulerJobId: rebootSchedulerJob.ID,
				State:          "queued",
				Notes:          []byte{},
				StartTime:      startTime,
			}
			rebootSchedulerLogs[index] = schedulerLog
		}
		_, err = s.schedulerRepository.SaveSchedulerLogs(rebootSchedulerLogs)

		if err != nil {
			return err
		}
	}
	return nil
}
