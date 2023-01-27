package module

import (
	"strconv"

	model_core "fermion/backend_core/internal/model/core"
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
	service Services
	// base_service module_base.ServiceBase
}

var AccessModuleHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if AccessModuleHandler != nil {
		return AccessModuleHandler
	}
	service := NewService()
	AccessModuleHandler = &handler{service}
	return AccessModuleHandler
}

func (h *handler) RegisterModule(c echo.Context) (err error) {

	//data := c.Get("template").(*access.AccessTemplate)
	var data *model_core.AccessModuleAction
	c.Bind(&data)
	s := c.Get("TokenUserID").(string)
	data.CreatedByID = helpers.ConvertStringToUint(s)
	err = h.service.RegisterModule(data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Module registered successfully", map[string]interface{}{"created_id": data.ID})
}

func (h *handler) UpdateModule(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	//data := c.Get("template").(*access.AccessTemplate)
	var data model_core.AccessModuleAction
	c.Bind(&data)
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateModuleDetails(uint(id), data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Module details updates successfully", map[string]interface{}{"updated_id": id})
}

func (h *handler) FetchAllModule(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	result, err := h.service.FindAllModules(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

func (h *handler) GetModule(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	result, err := h.service.GetModuleDetail(uint(id), p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

func (h *handler) DeleteModule(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	err = h.service.DeleteModuleDetail(uint(id), uint(user_id), p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "deleted Module details.", map[string]int{"deleted_id": id})
}

// GetAllCoreModule godoc
// @Summary Get all CoreModule
// @Description Get all CoreModule
// @Tags Access_Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CoreAppModule
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/module/core_module [get]
func (h *handler) FetchAllCoreModule(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	result, err := h.service.FindAllCoreModules(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// GetAllAccessCoreModule godoc
// @Summary Get all CoreModule
// @Description Get all CoreModule
// @Tags Access_Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CoreAppModule
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/module/side__module [get]
func (h *handler) FetchAllAccessCoreModule(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	access_template_id := c.Get("AccessTemplateId")
	p.Filters = "[[\"access_template_id\",\"=\",\"" + access_template_id.(string) + "\"]]"
	result, err := h.service.FindAllAccessCoreModules(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}
