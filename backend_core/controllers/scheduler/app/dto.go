package scheduler

import (
	"time"

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
type SchedulerJobDto struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Function     string         `json:"function"`
	FunctionType string         `json:"function_type"`
	Params       datatypes.JSON `json:"params"`
	State        *bool          `json:"state"`
	StartTime    time.Time      `json:"start_time"`
	Frequency    string         `json:"frequency"`
	SourceType   string         `json:"source_type"`
	SourceId     string         `json:"source_id"`
}
type SchedulerJobResponseDTO struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Function     string         `json:"function"`
	FunctionType string         `json:"function_type"`
	Params       datatypes.JSON `json:"params"`
	State        *bool          `json:"state"`
	StartTime    time.Time      `json:"start_time"`
	Frequency    string         `json:"frequency"`
	SourceType   string         `json:"source_type"`
	SourceId     string         `json:"source_id"`
}
type SchedulerLogResponseDTO struct {
	ID             uint                    `json:"id"`
	SchedulerJobId uint                    `json:"scheduler_job_id"`
	SchedulerJob   SchedulerJobResponseDTO `json:"scheduler_job"`
	State          string                  `json:"state"` //success / failed / queued
	Notes          datatypes.JSON          `json:"notes"`
	StartTime      time.Time               `json:"start_time"`
	EndTime        time.Time               `json:"end_time"`
}
