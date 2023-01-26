package shipping_orders_rto

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

// CreateRTO godoc
// @Summary CreateRTO
// @Description Create a RTO
// @Tags RTO
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   RTOCreateRequest body RTORequest true "Create a RTO"
// @Success 200 {object} RTOCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_rto/create [post]
func (h *handler) CreateRTO(c echo.Context) (err error) {
	data := new(RTORequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	id, err := h.service.CreateRTO(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Created", map[string]interface{}{"Created_RTO_Id": id})
}

// BulkCreateRTO godoc
// @Summary 	BulkCreateRTO
// @Description Create a Multiple RTO
// @Tags 		RTO
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkRTOCreateRequest body RTORequest true "Create a Bulk RTO"
// @Success 	200 {object} BulkRTOCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/shipping_orders_rto/bulk_create [post]
func (h *handler) BulkCreateRTO(c echo.Context) (err error) {
	data := new([]RTORequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateRTO(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Created", map[string]bool{"Created": true})
}

// GetRTO godoc
// @Summary GetRTO
// @Description Get a RTO
// @Tags RTO
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id  path int true "id"
// @Success 200 {object} RTOGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_rto/{id} [get]
func (h *handler) GetRTO(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetRTO(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// GetAllRTO godoc
// @Summary GetAllRTO
// @Description Get All RTO
// @Tags RTO
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllRTOResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_rto [get]
func (h *handler) GetAllRTO(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllRTO(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "RTO list Retrieved Successfully", result, p)
}

// GetAllRTODropDown godoc
// @Summary GetAllRTODropDown
// @Description Get DropDown RTO
// @Tags RTO
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllRTOResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_rto/dropdown [get]
func (h *handler) GetAllRTODropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllRTO(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "RTO DropDown Retrieved Successfully", result, p)
}

// UpdateRTO godoc
// @Summary UpdateRTO
// @Description Update the RTO
// @Tags RTO
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param RTOUpdateRequest body RTORequest true "Update the RTO"
// @Success 200 {object} UpdateRTOResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_rto/{id}/update [post]
func (h *handler) UpdateRTO(c echo.Context) (err error) {
	var data RTORequest
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateRTO(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Updated Successfully", map[string]int{"Updated id": id})
}

// DeleteRTO godoc
// @Summary DeleteRTO
// @Description Delete a RTO
// @Tags RTO
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Success 200 {object} DeleteRTOResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_rto/{id}/delete [delete]
func (h *handler) DeleteRTO(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteRTO(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Deleted Successfully", map[string]string{"Deleted_RTO_ID": get_id})
}
