package asn

import (
	"errors"
	"fmt"
	"strconv"

	concurrency_management "fermion/backend_core/controllers/concurrency_management"
	inventory_orders_base "fermion/backend_core/controllers/inventory_orders/base"
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
	service            Service
	base_service       inventory_orders_base.ServiceBase
	concurrencyService concurrency_management.Service
}

func NewHandler() *handler {
	service := NewService()
	base_service := inventory_orders_base.NewServiceBase()
	concurrencyService := concurrency_management.NewService()
	return &handler{service, base_service, concurrencyService}
}

// FavouriteAsnView godoc
// @Summary Get all favourite AsnList
// @Description Get all favourite AsnList
// @Tags Asn
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} AsnGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/asn/favourite_list [get]
func (h *handler) FavouriteAsnView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrAsn(user_id)
	if er != nil {
		return errors.New("favourite ASN list not found for the specified user")
	}
	result, err := h.base_service.GetFavAsnList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Asn list Retrieved Successfully", result, p)
}

// CreateAsn godoc
// @Summary      do a CreateAsn
// @Description  Create a Asn
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   AsnCreateRequest body AsnRequest true "create a Asn"
// @Success      200  {object}  AsnCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/create [post]
func (h *handler) CreateAsn(c echo.Context) (err error) {
	data := c.Get("asn").(*AsnRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	t := uint(user_id)
	data.CreatedByID = &t
	id, err := h.service.CreateAsn(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn created successfully", map[string]interface{}{"created_id": id})
}

// BulkCreateAsn godoc
// @Summary 	do a BulkCreateAsn
// @Description do a BulkCreateAsn
// @Tags 		Asn
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkAsnCreateRequest body AsnRequest true "Create a Bulk Asn"
// @Success 	200 {object} BulkAsnCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/asn/bulk_create [post]
func (h *handler) BulkCreateAsn(c echo.Context) (err error) {
	data := new([]AsnRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateAsn(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn multiple records created successfully", data)
}

// UpdateAsn godoc
// @Summary      do a UpdateAsn
// @Description  Update a Asn
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   AsnUpdateRequest body AsnRequest true "Update a Asn"
// @Success      200  {object}  AsnUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/{id}/edit [put]
func (h *handler) UpdateAsn(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("asn").(*AsnRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	//--------------------------concurrency_management_start-----------------------------------------
	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(uint(id), "status_id", "asns")
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(uint(id), PreviousStatusId, "status_id", "asns")
	//--------------------------concurrency_management_end-----------------------------------------
	err = h.service.UpdateAsn(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn updated successfully", map[string]interface{}{"updated_id": id})
}

// GetAsn godoc
// @Summary 	Find a GetAsn by id
// @Description Find a GetAsn by id
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  AsnGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/{id} [get]
func (h *handler) GetAsn(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAsn(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn fetched successfully", result)
}

// GetAllAsn godoc
// @Summary      get all GetAllAsn with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllAsn with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		query	  string false "filters"
// @Param        per_page 		query     int    false "per_page"
// @Param  	     page_no        query     int    false "page_no"
// @Param        sort			query     string false "sort"
// @Success      200  {object}  AsnGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn [get]
func (h *handler) GetAllAsn(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllAsn(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Asn data Retrieved successfully", result, p)
}

// DeleteAsn godoc
// @Summary      delete a DeleteAsn
// @Description  delete a DeleteAsn
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  AsnDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/{id}/delete [delete]
func (h *handler) DeleteAsn(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteAsn(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteAsnLines godoc
// @Summary      delete a DeleteAsnLines
// @Description  delete a DeleteAsnLines
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        asn_id         path      int    true  "asn_id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  AsnLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/{id}/delete_products [delete]
func (h *handler) DeleteAsnLines(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	product_id := c.QueryParam("product_id")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	prod_id, _ := strconv.Atoi(product_id)
	query := map[string]interface{}{
		"asn_id":     uint(id),
		"product_id": uint(prod_id),
	}
	err = h.service.DeleteAsnLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn Line Items successfully deleted", query)
}

// SendMailAsn godoc
// @Summary      Send Email to specific recevier for Id
// @Description  Send Email to specific recevier for Id
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Param        receiver_email query     string true "receiver_email"
// @Success      200  {object}  SendMailAsnResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/{id}/sendemail [get]
func (h *handler) SendMailAsn(c echo.Context) (err error) {

	q := new(SendMailAsn)
	c.Bind(q)
	// err = h.service.SendMailAsn(q)
	// if err != nil {
	// 	return res.RespError(c, err)
	// }

	return res.RespSuccess(c, "Email Sent Successfully", q)
}

// DownloadPdfAsn godoc
// @Summary      Download PDF
// @Description  Download PDF
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Success      200  {object}  DownloadPdfAsnResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/{id}/printpdf [get]
func (h *handler) DownloadPdfAsn(c echo.Context) (err error) {
	return res.RespSuccess(c, "Pdf Generated Successfully", "https://drive.google.com/file/d/1O46B3lgliUGakL4weBjrovCHphTzR5Np/view?usp=sharing")
}

// FavouriteAsn godoc
// @Summary FavouriteAsn
// @Description FavouriteAsn
// @Tags Asn
// @Accept  json
// @Produce  json
// @param id path string true "Asn ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/asn/{id}/favourite [post]
func (h *handler) FavouriteAsn(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	fmt.Println("handler ------------------------ -1", query)
	err = h.base_service.FavouriteAsn(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteAsn godoc
// @Summary UnFavouriteAsn
// @Description UnFavouriteAsn
// @Tags Asn
// @Accept  json
// @Produce  json
// @param id path string true "Asn ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/asn/{id}/unfavourite [post]
func (h *handler) UnFavouriteAsn(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteAsn(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Asn is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetAsnTab godoc
// @Summary GetAsnTab
// @Summary GetAsnTab
// @Description GetAsnTab
// @Tags Asn
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
// @Router /api/v1/asn/{id}/filter_module/{tab} [get]
func (h *handler) GetAsnTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetAsnTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

// GetAllAsnDropdown godoc
// @Summary      get all GetAllAsnDropdown with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllAsnDropdown with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         Asn
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		query	  string false "filters"
// @Param        per_page 		query     int    false "per_page"
// @Param  	     page_no        query     int    false "page_no"
// @Param        sort			query     string false "sort"
// @Success      200  {object}  AsnGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/asn/dropdown [get]
func (h *handler) GetAllAsnDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllAsn(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Asn data Retrieved successfully", result, p)
}
