package scheduler

import (
	"time"

	"gorm.io/datatypes"
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
type SchedulerJobDto struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Function     string         `json:"function"`
	FunctionType string         `json:"function_type"`
	Params       datatypes.JSON `json:"params"`
	State        string         `json:"state"`
	StartTime    time.Time      `json:"start_time"`
	Frequency    string         `json:"frequency"`
}
type SchedulerJobResponseDTO struct {
	ID           uint           `json:"id"`
	CreatedDate  uint           `json:"created_date"`
	CreatedBy    string         `json:"created_by"`
	UpdatedBy    string         `json:"updated_by"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Function     string         `json:"function"`
	FunctionType string         `json:"function_type"`
	Params       datatypes.JSON `json:"params"`
	State        string         `json:"state"`
	StartTime    time.Time      `json:"start_time"`
	Frequency    string         `json:"frequency"`
}
type SchedulerLogResponseDTO struct {
	ID             uint           `json:"id"`
	SchedulerJobId uint           `json:"scheduler_job_id"`
	State          string         `json:"state"` //success / failed / queued
	Notes          datatypes.JSON `json:"notes"`
	StartTime      time.Time      `json:"start_time"`
	EndTime        time.Time      `json:"end_time"`
}
