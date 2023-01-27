package webstores

import (
	"fmt"
	"strconv"

	webstore_base "fermion/backend_core/controllers/omnichannel/base"
	"fermion/backend_core/internal/model/omnichannel"
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
	service      Services
	base_service webstore_base.ServiceBase
}

var WebstoresHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if WebstoresHandler != nil {
		return WebstoresHandler
	}
	service := NewService()
	base_service := webstore_base.NewServiceBase()
	WebstoresHandler = &handler{service, base_service}
	return WebstoresHandler
}

//----------------------------------Webstore--------------------------------------------------

// CreateWebstores godoc
// @Summary      do a CreateWebstore
// @Description  Create a Webstore
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   WebstoreCreateRequest body User_Webstore_Link true "create a Webstore"
// @Success      200  {object}  WebstoredetailCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores/create [post]
func (h *handler) RegisterWebstore(c echo.Context) (err error) {

	data := c.Get("webstores").(*omnichannel.User_Webstore_Link)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.RegisterWebstores(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateWebstoreInCache(token_id, access_template_id)

	return res.RespSuccess(c, "webstore registered successfully", map[string]interface{}{"created_id": data.ID})
}

// UpdateWebstores godoc
// @Summary      do a UpdateWebstores
// @Description  Update a Webstore
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   WebstoreUpdateRequest body User_Webstore_Link true "Update a Webstore"
// @Success      200  {object}  WebstoredetailUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores/{id}/update [post]
func (h *handler) UpdateWebstore(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("webstores").(*omnichannel.User_Webstore_Link)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateWebstoreDetails(uint(id), *data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateWebstoreInCache(token_id, access_template_id)

	return res.RespSuccess(c, "webstore details updates successfully", map[string]interface{}{"updated_id": id})
}

// AllWebstores godoc
// @Summary      get AllWebstores and filter it by search query
// @Description  a get AllWebstores and filter it by search query
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		filters		  query	string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  WebstoredetailGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores [get]
func (h *handler) FetchAllWebstores(c echo.Context) (err error) {
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	p := new(pagination.Paginatevalue)
	c.Bind(p)

	var response interface{}
	if *p == pagination.BasePaginatevalue {
		response, *p = GetWebstoreFromCache(token_id)
	}

	if response != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.FindAllWebstore(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data retrieved successfully", result, p)
}

// AllWebstores godoc
// @Summary      get AllWebstores dropdown list and filter it by search query
// @Description  a get AllWebstores dropdown list and filter it by search query
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		filters		  query	string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  WebstoredetailGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores/dropdown [get]
func (h *handler) FetchAllWebstoresDropDown(c echo.Context) (err error) {
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	p := new(pagination.Paginatevalue)
	c.Bind(p)

	var response interface{}
	if *p == pagination.BasePaginatevalue {
		response, *p = GetWebstoreFromCache(token_id)
	}

	if response != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.FindAllWebstore(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "dropdown data retrieved successfully", result, p)
}

// FindWebstores godoc
// @Summary      Find a Webstore by id
// @Description  Find a Webstore by id
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  WebstoredetailGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores/{id} [get]
func (h *handler) GetWebstore(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetWebstoreDetail(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// DeleteWebstores godoc
// @Summary      delete a DeleteWebstores
// @Description  delete a DeleteWebstores
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  WebstoredetailDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores/{id}/delete [delete]
func (h *handler) DeleteWebstore(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteWebstoreDetail(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateWebstoreInCache(token_id, access_template_id)

	return res.RespSuccess(c, "deleted webstore details.", map[string]int{"deleted_id": id})
}

// AvailableWebstores godoc
// @Summary      get AvailableWebstores and filter it by search query
// @Description  get AvailableWebstores and filter it by search query
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		filters		  query	string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  AvailableWebstoresResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores/available [get]
func (h *handler) AvailableWebstores(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AvailableWebstores(p, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data retrieved successfully", result, p)
}

// Get AuthKeys godoc
// @Summary      Get AuthKeys
// @Description  Get AuthKeys
// @Tags         Webstores
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        webstore   path      string  true  "webstore_code"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/webstores/{webstore_code}/get_auth_keys [get]
func (h *handler) GetAuthKeys(c echo.Context) (err error) {

	webstore := c.Param("webstore_code")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"webstore_code": webstore,
		"created_by":    user_id,
	}
	result, err := h.service.GetAuthKeys(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// FavouriteWebstores godoc
// @Summary FavouriteWebstores
// @Description FavouriteWebstores
// @Tags Webstores
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Webstore ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/webstores/{id}/favourite [post]
func (h *handler) FavouriteWebstore(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.GetWebstoreDetail(uint(ID), token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	fmt.Println("handler ------------------------ -1", query)
	err = h.base_service.FavouriteWebstoreFav(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Webstore is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteWebstore godoc
// @Summary UnFavouriteWebstore
// @Description UnFavouriteWebstore
// @Tags Webstores
// @Accept  json
// @Produce  json
// @param id path string true "Webstore ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/webstores/{id}/unfavourite [post]
func (h *handler) UnFavouriteWebstore(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteWebstoreFav(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Webstore is Unmarked as Favourite", map[string]string{"record_id": id})
}
