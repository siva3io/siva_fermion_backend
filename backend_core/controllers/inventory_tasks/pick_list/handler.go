package pick_list

import (
	"errors"
	"fmt"
	"strconv"

	inventory_tasks_base "fermion/backend_core/controllers/inventory_tasks/base"
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
	base_service inventory_tasks_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := inventory_tasks_base.NewServiceBase()
	return &handler{service, base_service}
}

// FavouritePickListView godoc
// @Summary Get all favourite PickList
// @Description Get all favourite PickList
// @Tags PickList
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} PicklistGetAllResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pick_list/favourite_list [get]
func (h *handler) FavouritePickListView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrPicList(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavPicList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Pick list Retrieved Successfully", result, p)
}

// CreatePicklist godoc
// @Summary      do a CreatePicklist
// @Description  Create a Picklist
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   PicklistCreateRequest body PicklistRequest true "create a Picklist"
// @Success      200  {object}  PicklistCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/create [post]
func (h *handler) CreatePickList(c echo.Context) (err error) {
	data := c.Get("pick_list").(*PicklistRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	t := uint(user_id)
	data.CreatedByID = &t
	id, err := h.service.CreatePickList(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Picklist created successfully", map[string]interface{}{"created_id": id})
}

// BulkCreatePicklist godoc
// @Summary 	do a BulkCreatePicklist
// @Description do a BulkCreatePicklist
// @Tags 		PickList
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkPicklistCreateRequest body PicklistRequest true "Create a Bulk Picklist"
// @Success 	200 {object} BulkPicklistCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/pick_list/bulk_create [post]
func (h *handler) BulkCreatePickList(c echo.Context) (err error) {
	data := new([]PicklistRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreatePickList(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Picklist multiple records created successfully", data)
}

// UpdatePicklist godoc
// @Summary      do a UpdatePicklist
// @Description  Update a Picklist
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   PicklistUpdateRequest body PicklistRequest true "Update a Picklist"
// @Success      200  {object}  PicklistUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/edit [put]
func (h *handler) UpdatePickList(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("pick_list").(*PicklistRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdatePickList(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Picklist updated successfully", map[string]interface{}{"updated_id": id})
}

// GetPicklist godoc
// @Summary      Find a Picklist by id
// @Description  Find a Picklist by id
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  PicklistGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id} [get]
func (h *handler) GetPickList(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetPickList(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Picklist fetched successfully", result)
}

// GetAllPicklist godoc
// @Summary      get all GetAllPicklist with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllPicklist with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		   query	 string false "filters"
// @Param        per_page 		   query     int    false "per_page"
// @Param  	     page_no           query     int    false "page_no"
// @Param        sort			   query     string false "sort"
// @Success      200  {object}  PicklistGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list [get]
func (h *handler) GetAllPickList(c echo.Context) (err error) {
	//per_page, page_no, sort
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllPickList(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Picklist data retrieved Successfully", result, p)
}

// DeletePicklist godoc
// @Summary      delete a DeletePicklist
// @Description  delete a DeletePicklist
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  PicklistDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/delete [delete]
func (h *handler) DeletePickList(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeletePickList(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Picklist deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeletePicklistLines godoc
// @Summary      delete a DeletePicklistLines
// @Description  delete a DeletePicklistLines
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        pick_list_id   path      int    true  "pick_list_id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  PicklistLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/delete_products [delete]
func (h *handler) DeletePickListLines(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	product_id := c.QueryParam("product_id")
	prod_id, _ := strconv.Atoi(product_id)
	query := map[string]interface{}{
		"pick_list_id": uint(id),
		"product_id":   uint(prod_id),
	}
	err = h.service.DeletePickListLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Picklist Line Items successfully deleted", query)
}

// SendMailPickList godoc
// @Summary      Send Email to specific recevier for Id
// @Description  Send Email to specific recevier for Id
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Param        receiver_email query     string true "receiver_email"
// @Success      200  {object}  SendMailPickListResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/sendemail [get]
func (h *handler) SendMailPickList(c echo.Context) (err error) {

	q := new(SendMailPickList)
	c.Bind(q)
	// err = h.service.SendMailPickList(q)
	// if err != nil {
	// 	return res.RespError(c, err)
	// }

	return res.RespSuccess(c, "Email Sent Successfully", q)
}

// DownloadPdfPickList godoc
// @Summary      Download PDF
// @Description  Download PDF
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Success      200  {object}  DownloadPdfPickListResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/printpdf [get]
func (h *handler) DownloadPdfPickList(c echo.Context) (err error) {
	return res.RespSuccess(c, "Pdf Generated Successfully", "https://drive.google.com/file/d/1O46B3lgliUGakL4weBjrovCHphTzR5Np/view?usp=sharing")
}

// FavoritePickList godoc
// @Summary FavoritePickList
// @Description FavoritePickList
// @Tags PickList
// @Accept  json
// @Produce  json
// @param id path string true "PickList ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pick_list/{id}/favourite [post]
func (h *handler) FavouritePickList(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.FavouritePickList(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Pick List is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouritePickList godoc
// @Summary UnFavoritePickList
// @Description UnFavoritePickList
// @Tags PickList
// @Accept  json
// @Produce  json
// @param id path string true "PickList ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pick_list/{id}/unfavourite [post]
func (h *handler) UnFavouritePickList(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouritePickList(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Pick List is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetAllPickListDropDown godoc
// @Summary      get all GetAllPickListDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllPickListDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		   query	 string false "filters"
// @Param        per_page 		   query     int    false "per_page"
// @Param  	     page_no           query     int    false "page_no"
// @Param        sort			   query     string false "sort"
// @Success      200  {object}  PicklistGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/dropdown [get]
func (h *handler) GetAllPickListDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllPickList(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Picklist data retrieved Successfully", result, p)
}
