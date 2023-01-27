package scheduler

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"

	"gorm.io/datatypes"
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
var SchedulerTables = []interface{}{
	&SchedulerJob{},
	&SchedulerLog{},
}

type SchedulerJob struct {
	model_core.Model
	Name         string         `json:"name"` //name
	Type         string         `json:"type"`
	Function     string         `json:"function"`
	FunctionType string         `json:"function_type"` //api / function
	Params       datatypes.JSON `json:"params"`
	State        *bool          `json:"state"` //active / inactive
	StartTime    time.Time      `json:"start_time"`
	Frequency    string         `json:"frequency"`
	SourceType   string         `json:"source_type"`
	SourceId     string         `json:"source_id"`
	// FrequencyType string
}

type SchedulerLog struct {
	model_core.Model
	SchedulerJobId uint           `json:"scheduler_job_id"`
	SchedulerJob   SchedulerJob   `json:"scheduler_job" gorm:"foreignkey:SchedulerJobId; references:ID"`
	State          string         `json:"state"` //success / failed / queued
	Notes          datatypes.JSON `json:"notes"`
	StartTime      time.Time      `json:"start_time"`
	EndTime        time.Time      `json:"end_time"`
}
