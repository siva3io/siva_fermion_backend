package shipping_orders_wd

import (
	"fmt"
	"strconv"

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

func NewHandler() *handler {
	service := NewService()
	return &handler{service}
}

// CreateWD godoc
// @Summary CreateWD
// @Description Create a WD
// @Tags WD
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   WDCreateRequest body WDRequest true "Create a WD"
// @Success 200 {object} WDCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_wd/create [post]
func (h *handler) CreateWD(c echo.Context) (err error) {
	data := new(WDRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	id, err := h.service.CreateWD(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Created", map[string]interface{}{"Created_WD_Id": id})
}

// BulkCreateWD godoc
// @Summary 	BulkCreateWD
// @Description Create a Multiple WD
// @Tags 		WD
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkWDCreateRequest body WDRequest true "Create a Bulk WD"
// @Success 	200 {object} BulkWDCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/shipping_orders_wd/bulk_create [post]
func (h *handler) BulkCreateWD(c echo.Context) (err error) {
	data := new([]WDRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateWD(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Created", map[string]bool{"Created": true})
}

// GetWD godoc
// @Summary GetWD
// @Description Get a WD
// @Tags WD
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Success 200 {object} WDGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_wd/{id} [get]
func (h *handler) GetWD(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetWD(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// GetAllWD godoc
// @Summary GetAllWD
// @Description Get All WD
// @Tags WD
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllWDResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_wd [get]
func (h *handler) GetAllWD(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllWD(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "WD list Retrieved Successfully", result, p)
}

// GetAllWDDropDown godoc
// @Summary GetAllWDDropDown
// @Description Get All WD
// @Tags WD
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllWDResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_wd/dropdown [get]
func (h *handler) GetAllWDDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllWD(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "WD dropdown list Retrieved Successfully", result, p)
}

// UpdateWD godoc
// @Summary UpdateWD
// @Description Update the WD
// @Tags WD
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param WDUpdateRequest body WDRequest true "Update the WD"
// @Success 200 {object} UpdateWDResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_wd/{id}/update [post]
func (h *handler) UpdateWD(c echo.Context) (err error) {
	var data WDRequest
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateWD(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Updated Successfully", map[string]int{"Updated id": id})
}

// DeleteWD godoc
// @Summary DeleteWD
// @Description Delete a WD
// @Tags WD
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Success 200 {object} DeleteWDResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_wd/{id}/delete [delete]
func (h *handler) DeleteWD(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteWD(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Deleted Successfully", map[string]string{"Deleted_WD_ID": get_id})
}
