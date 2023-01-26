package grn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	inventory_orders_base "fermion/backend_core/controllers/inventory_orders/base"
	"fermion/backend_core/internal/model/inventory_orders"
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
	service      Service
	base_service inventory_orders_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := inventory_orders_base.NewServiceBase()
	return &handler{service, base_service}
}

// FavouriteGrnView godoc
// @Summary Get all favourite GrnList
// @Description Get all favourite GrnList
// @Tags GRN
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} GRNGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/favourite_list [get]
func (h *handler) FavouriteGrnView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrGrn(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetGrnFavList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Grn list Retrieved Successfully", result, p)
}

// CreateGRN godoc
// @Summary CreateGRN
// @Description Create a GRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   GRNCreateRequest body GRNRequest true "Create a GRN"
// @Success 200 {object} GRNCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/create [post]
func (h *handler) CreateGRN(c echo.Context) (err error) {
	data := c.Get("grn").(*GRNRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	t := uint(user_id)
	data.CreatedByID = &t
	id, err := h.service.CreateGRN(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "GRN Successfully Created", map[string]interface{}{"Created_Id": id})
}

// BulkCreateGRN godoc
// @Summary 	BulkCreateGRN
// @Description Create a Multiple GRN
// @Tags 		GRN
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkGRNCreateRequest body GRNRequest true "Create a Bulk GRN"
// @Success 	200 {object} BulkGRNCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/grn/bulk_create [post]
func (h *handler) BulkCreateGRN(c echo.Context) (err error) {
	data := new([]inventory_orders.GRN)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateGRN(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Created", map[string]bool{"Created": true})
}

// GetGRN godoc
// @Summary GetGRN
// @Description Get a GRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id  path int true "Get GRN"
// @Success 200 {object} GRNGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id} [get]
func (h *handler) GetGRN(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetGRN(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// GetAllGRN godoc
// @Summary GetAllGRN
// @Description Get All GRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 			query int 	 false "per_page"
// @Param page_no           query int 	 false "page_no"
// @Param sort				query string false "sort"
// @param filters 			query string false "Filters"
// @Success 200 {object} GetAllGRNResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn [get]
func (h *handler) GetAllGRN(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllGRN(page, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "GRN Successfully Retrieved", result, page)
}

// UpdateGRN godoc
// @Summary UpdateGRN
// @Description Update a GRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param GRNUpdateRequest body GRNRequest true "Update a GRN"
// @Success 200 {object} UpdateGRNResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/update [Post]
func (h *handler) UpdateGRN(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	data := c.Get("grn").(*GRNRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateGRN(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Updated Successfully", map[string]int{"Updated id": id})
}

// SearchGRN godoc
// @Summary Search SearchGRN
// @Description Search SearchGRN
// @Tags GRN
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   query 	query 	string 	true "query"
// @Success 200 {array} GetAllGRNResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/search [get]
func (h *handler) SearchGRN(c echo.Context) error {
	query := c.QueryParam("query")
	var data []inventory_orders.GRN
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.SearchGRN(query, token_id, access_template_id)
	value, _ := json.Marshal(resp)
	_ = json.Unmarshal(value, &data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Fetched", data)
}

// DeleteGRN godoc
// @Summary DeleteGRN
// @Description Delete a GRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Success 200 {object} DeleteGRNResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/delete [delete]
func (h *handler) DeleteGRN(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteGRN(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Deleted Successfully", map[string]string{"Deleted grn": get_id})
}

// DeleteGRN godoc
// @Summary DeleteGRNOrderLine
// @Description Delete a GRN Order Line
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param product_id query string true "product_id"
// @Success 200 {object} DeleteGRNOrderLineResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/order_line/{id}/delete [delete]
func (h *handler) DeleteGRNOrderLine(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	query := map[string]interface{}{
		"grn_id":     uint(id),
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeleteGRNOrderLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Line Item Deleted Successfully", query)
}

// SendEmailGRN godoc
// @Summary Email GRN
// @Description Email GRN
// @Tags GRN
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/send_email [post]
func (h *handler) SendEmailGRN(c echo.Context) error {
	return res.RespSuccess(c, "Email sent successfully", nil)
}

// DownloadGRNPDF godoc
// @Summary Download GRN
// @Description Download GRN
// @Tags GRN
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/download_pdf [post]
func (h *handler) DownloadGRNPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

// GenerateGRNPDF godoc
// @Summary Generate GRN PDF
// @Description Generate GRN PDF
// @Tags GRN
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/generate_pdf [post]
func (h *handler) GenerateGRNPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF Generated", resp_data)
}

// FavouriteGRN godoc
// @Summary FavouriteGRN
// @Description FavouriteGRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path string true "GRN ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/favourite [post]
func (h *handler) FavouriteGrn(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.FavouriteGrn(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "GRN is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteGRN godoc
// @Summary UnFavouriteGRN
// @Description UnFavouriteGRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path string true "GRN ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/unfavourite [post]
func (h *handler) UnFavouriteGrn(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteGrn(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "GRN is marked as Unfavourite", map[string]string{"record_id": id})
}

// GetGrnTab godoc
// @Summary GetGrnTab
// @Summary GetGrnTab
// @Description GetGrnTab
// @Tags Grn
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @param tab path string true "tab"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/{id}/filter_module/{tab} [get]
func (h *handler) GetGrnTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetGrnTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

// GetAllGRNDropDown godoc
// @Summary GetAllGRNDropDown
// @Description Get All GRN
// @Tags GRN
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 			query int 	 false "per_page"
// @Param page_no           query int 	 false "page_no"
// @Param sort				query string false "sort"
// @param filters 			query string false "Filters"
// @Success 200 {object} GetAllGRNResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/grn/dropdown [get]
func (h *handler) GetAllGRNDropDown(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllGRN(page, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "GRN Successfully Retrieved", result, page)
}
