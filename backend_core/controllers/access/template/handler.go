package template

import (
	"encoding/json"
	"fmt"
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
}

var AccessTemplateHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if AccessTemplateHandler != nil {
		return AccessTemplateHandler
	}
	service := NewService()
	AccessTemplateHandler = &handler{service}
	return AccessTemplateHandler
}

// CreateTemplate godoc
// @Summary CreateTemplate
// @Description CreateTemplate
// @Tags Access_Management
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body AccessTemplate true "Template Request Body"
// @Success 200 {object} AccessTemplate
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/template/create [post]
func (h *handler) RegisterTemplate(c echo.Context) (err error) {

	// data := c.Get("templates").(*access.AccessTemplate)

	var data map[string]interface{}
	var access_template model_core.AccessTemplate
	var module_data []model_core.AccessModuleAction

	c.Bind(&data)
	s := c.Get("TokenUserID").(string)
	access_template.CreatedByID = helpers.ConvertStringToUint(s)
	dto, _ := json.Marshal(data)
	_ = json.Unmarshal(dto, &access_template)
	dto, _ = json.Marshal(data["module_data"])
	_ = json.Unmarshal(dto, &module_data)
	err = h.service.RegisterTemplates(&access_template, module_data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Template registered successfully", map[string]interface{}{"created_id": access_template.ID})
}

// UpdateTemplate godoc
// @Summary UpdateTemplate
// @Description UpdateTemplate
// @Tags Access_Management
// @Accept  json
// @Produce  json
// @param id path string true "Template ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body AccessTemplate true "Template  Request Body"
// @Success 200 {object} AccessTemplate
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/template/{id}/update [post]
func (h *handler) UpdateTemplate(c echo.Context) (err error) {

	var data = make(map[string]interface{}, 0)
	var access_template model_core.AccessTemplate
	var module_data []model_core.AccessModuleAction
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	c.Bind(&data)
	//s := c.Get("TokenUserID").(string)

	dto, _ := json.Marshal(data)
	_ = json.Unmarshal(dto, &access_template)
	dto, _ = json.Marshal(data["module_data"])
	_ = json.Unmarshal(dto, &module_data)
	//access_template.CreatedByID = helpers.ConvertStringToUint(s)

	err = h.service.UpdateTemplateDetails(uint(id), access_template, module_data)

	//var data model_core.AccessTemplate
	//c.Bind(&data)
	//s := c.Get("TokenUserID").(string)
	//data.UpdatedByID = helpers.ConvertStringToUint(s)
	//err = h.service.UpdateTemplateDetails(uint(id), data)
	if err != nil {
		fmt.Println(err)
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Template details updates successfully", map[string]interface{}{"updated_id": id})
}

// GetTemplatesList godoc
// @Summary Get all Templates
// @Description Get all Templates
// @Tags Access_Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} AccessTemplate
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/template/list [get]
func (h *handler) FetchAllTemplatesList(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	result, err := h.service.GetAllTemplateList(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// GetTemplates godoc
// @Summary Get all Templates
// @Description Get all Templates
// @Tags Access_Management
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} AccessTemplate
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/template [get]
func (h *handler) FetchAllTemplates(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	result, err := h.service.FindAllTemplate(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

func (h *handler) GetTemplate(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	result, err := h.service.GetTemplateDetail(uint(id), p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// DeleteTemplate godoc
// @Summary DeleteTemplate
// @Description DeleteTemplate
// @Tags Access_Management
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Template ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/template/{id}/delete [delete]
func (h *handler) DeleteTemplate(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	err = h.service.DeleteTemplateDetail(uint(id), uint(user_id))
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "deleted template details.", map[string]int{"deleted_id": id})
}
