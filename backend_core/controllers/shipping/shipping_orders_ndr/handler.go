package shipping_orders_ndr

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

// CreateNDR godoc
// @Summary CreateNDR
// @Description Create a NDR
// @Tags NDR
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   NDRCreateRequest body NDRRequest true "Create a NDR"
// @Success 200 {object} NDRCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_ndr/create [post]
func (h *handler) CreateNDR(c echo.Context) (err error) {
	data := new(NDRRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	id, err := h.service.CreateNDR(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Created", map[string]interface{}{"Created_NDR_Id": id})
}

// BulkCreateNDR godoc
// @Summary 	BulkCreateNDR
// @Description Create a Multiple NDR
// @Tags 		NDR
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkNDRCreateRequest body NDRRequest true "Create a Bulk NDR"
// @Success 	200 {object} BulkNDRCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/shipping_orders_ndr/bulk_create [post]
func (h *handler) BulkCreateNDR(c echo.Context) (err error) {
	data := new([]NDRRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateNDR(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Created", map[string]bool{"Created": true})
}

// GetNDR godoc
// @Summary GetNDR
// @Description Get a NDR
// @Tags NDR
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id  path int true "id"
// @Success 200 {object} NDRGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_ndr/{id} [get]
func (h *handler) GetNDR(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetNDR(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// GetAllNDR godoc
// @Summary Get all NDR and filter it by search query
// @Description Get All NDR
// @Tags NDR
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllNDRResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_ndr [get]
func (h *handler) GetAllNDR(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllNDR(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "NDR list Retrieved successfully", result, p)
}

// GetAllNDRDropDown godoc
// @Summary Get all NDR and filter it by search query
// @Description Get NDR dropdown
// @Tags NDR
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllNDRResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_ndr/dropdown [get]
func (h *handler) GetAllNDRDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllNDR(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "DropDown NDR list Retrieved successfully", result, p)
}

// UpdateNDR godoc
// @Summary UpdateNDR
// @Description Update the NDR
// @Tags NDR
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param NDRUpdateRequest body NDRRequest true "Update the NDR"
// @Success 200 {object} UpdateNDRResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_ndr/{id}/update [post]
func (h *handler) UpdateNDR(c echo.Context) (err error) {
	var data NDRRequest
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateNDR(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Updated Successfully", map[string]int{"Updated id": id})
}

// DeleteNDR godoc
// @Summary DeleteNDR
// @Description Delete a NDR
// @Tags NDR
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Success 200 {object} DeleteNDRResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_ndr/{id}/delete [delete]
func (h *handler) DeleteNDR(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteNDR(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Deleted Successfully", map[string]string{"Deleted_NDR_ID": get_id})
}

// DeleteNDRLine godoc
// @Summary DeleteNDRLine
// @Description Delete a NDR Line
// @Tags NDR
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param ndrs_id query string true "ndrs_id"
// @Success 200 {object} DeleteNDRLineResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders_ndr/ndr_lines/{id}/delete [delete]
func (h *handler) DeleteNDRLines(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	query := map[string]interface{}{
		"non_delivery_request_id": uint(id),
		"ndrs_id":                 c.QueryParam("ndrs_id"),
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	fmt.Println("!!!", query)
	err = h.service.DeleteNDRLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "NDR Line successfully deleted", query)
}
