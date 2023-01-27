package scheduler_init

import (
	"flag"
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

func Init() {
	go InitSchedulerJob()
}

var cr = cron.New()
var specParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

func InitSchedulerJob() error {

	schedulerFrequencyInMinute := 5
	count := 1

	clearDB := flag.Lookup("clearDB").Value.(flag.Getter).Get().(bool)
	migrateDB := flag.Lookup("migrateDB").Value.(flag.Getter).Get().(bool)
	seedMD := flag.Lookup("seedMD").Value.(flag.Getter).Get().(bool)

	if clearDB && migrateDB && seedMD {
		err := InitJobsAtReboot()
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

		schedulerLogs, err := repo.NewScheduler().FindAllSchedulerLog(nil, query)

		if err != nil {
			panic(err)
		}

		err = repo.NewScheduler().UpdateSchedulerLog(query, model.SchedulerLog{State: "in_progress"})

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

					nextLog, err = repo.NewScheduler().SaveSchedulerLog(nextLog)

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

				err = repo.NewScheduler().UpdateSchedulerLog(schedulerLogUpdateQuery, schedulerLog)

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
func InitJobsAtReboot() error {

	var query map[string]interface{}

	fmt.Println("Rebooting Scheduler...!")

	rebootSchedulerJobs, err := repo.NewScheduler().FindAll(nil, query)
	if err != nil {
		return err
	}

	lenRebootSchedulerJobs := len(rebootSchedulerJobs)

	if lenRebootSchedulerJobs > 0 {
		rebootSchedulerLogs := make([]model.SchedulerLog, lenRebootSchedulerJobs)
		for index, rebootSchedulerJob := range rebootSchedulerJobs {

			if rebootSchedulerJob.Frequency == "@once" || !*rebootSchedulerJob.State {
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
		_, err = repo.NewScheduler().SaveSchedulerLogs(rebootSchedulerLogs)

		if err != nil {
			return err
		}
	}
	return nil
}
