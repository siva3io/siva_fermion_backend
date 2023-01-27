package scheduler

import (
	"fmt"
	"strconv"

	utils "fermion/backend_core/controllers/scheduler/utils"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/labstack/echo/v4"
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
type handler struct {
	service Service
}

var SchedulerHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if SchedulerHandler != nil {
		return SchedulerHandler
	}
	service := NewService()
	SchedulerHandler = &handler{service}
	return SchedulerHandler
}

func (h *handler) AddSchedulerJob(c echo.Context) (err error) {

	var schedulerJobDto SchedulerJobDto
	if err := c.Bind(&schedulerJobDto); err != nil {
		return res.RespErr(c, err)
	}

	fmt.Println(schedulerJobDto)

	if err := h.service.AddSchedulerJob(schedulerJobDto); err != nil {
		return res.RespErr(c, err)
	}

	// logs, err := h.service.GetQueuedSchedulerLogs()

	return res.RespSuccess(c, "success", nil)
}
func (h *handler) ListSchedulerJob(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}
	result, err := h.service.GetAllSchedulerJob(page)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data retrieved successfully", result, page)

}

func (h *handler) ViewSchedulerJob(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = map[string]interface{}{
		"id": ID,
	}
	result, err := h.service.GetOneSchedulerJob(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Scheduler Job details retrieved successfully", result)
}
func (h *handler) DeleteSchedulerJob(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query["user_id"] = user_id
	err = h.service.DeleteSchedulerJob(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Scheduler Job deleted successfully", map[string]string{"deleted_id": id})

}
func (h *handler) UpdateSchedulerJob(c echo.Context) (err error) {

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = uint(ID)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query["user_id"] = user_id
	var data SchedulerJobDto
	err = c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = h.service.UpdateSchedulerJob(query, &data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Scheduler Job updated successfully", map[string]interface{}{"updated_id": id})
}
func (h *handler) ListSchedulerLogs(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}
	result, err := h.service.ListSchedulerLogs(page)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "data retrieved successfully", result, page)
}

func (h *handler) ListSchedulerFunctions(c echo.Context) (err error) {

	data := helpers.GetStructMethods(utils.CronFunctions{})

	return res.RespSuccess(c, "data retrieved successfully", data)
}
