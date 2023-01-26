package inventory_adjustments

import (
	"errors"
	"fmt"
	"strconv"

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
	service      Service
	base_service inventory_orders_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := inventory_orders_base.NewServiceBase()
	return &handler{service, base_service}
}

// FavouriteInventoryAdjustmentView godoc
// @Summary Get all favourite InventoryAdjustmentList
// @Description Get all favourite InventoryAdjustmentList
// @Tags Inventory Adjustments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} InvAdjGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/inventory_adjustments/favourite_list [get]
func (h *handler) FavouriteInvAdjView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrInvAdj(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavInvAdjList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Inventory adjustment list Retrieved Successfully", result, p)
}

// CreateInventoryAdjustment godoc
// @Summary      do a CreateInventoryAdjustment
// @Description  Create a Inventory Adjustment
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   InvAdjCreateRequest body InventoryAdjustmentsRequest true "create a InvAdj"
// @Success      200  {object}  InvAdjCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/create [post]
func (h *handler) CreateInvAdj(c echo.Context) (err error) {
	data := c.Get("inventory_adjustments").(*InventoryAdjustmentsRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	t := uint(user_id)
	data.CreatedByID = &t
	id, err := h.service.CreateInvAdj(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)

	}
	return res.RespSuccess(c, "Inventory Adjustment created successfully", map[string]interface{}{"created_id": id})
}

// BulkCreateInventoryAdjustments godoc
// @Summary 	do a BulkCreateInventoryAdjustments
// @Description do a BulkCreateInventoryAdjustments
// @Tags 		Inventory Adjustments
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkInvAdjCreateRequest body InventoryAdjustmentsRequest true "Create a Bulk InvAdj"
// @Success 	200 {object} BulkInvAdjCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/inventory_adjustments/bulk_create [post]
func (h *handler) BulkCreateInvAdj(c echo.Context) (err error) {
	data := new([]InventoryAdjustmentsRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateInvAdj(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Inventory Adjustments multiple records created successfully", data)
}

// UpdateInventoryAdjustment godoc
// @Summary      do a UpdateInventoryAdjustment
// @Description  Update a Inventory Adjustment
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   InvAdjUpdateRequest body InventoryAdjustmentsRequest true "Update a InvAdj"
// @Success      200  {object}  InvAdjUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/{id}/edit [put]
func (h *handler) UpdateInvAdj(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("inventory_adjustments").(*InventoryAdjustmentsRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateInvAdj(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "InventoryAdjustment updated successfully", map[string]interface{}{"updated_id": id})
}

// GetInventoryAdjustment godoc
// @Summary 	Find a GetInventoryAdjustment by id
// @Description Find a GetInventoryAdjustment by id
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  InvAdjGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/{id} [get]
func (h *handler) GetInvAdj(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetInvAdj(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Inventory Adjustment fetched successfully", result)
}

// GetAllInventoryAdjustments godoc
// @Summary      get all GetAllInventoryAdjustments with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllInventoryAdjustments with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		query	  string false "filters"
// @Param        per_page 		query     int    false "per_page"
// @Param  	     page_no        query     int    false "page_no"
// @Param        sort		    query     string false "sort"
// @Success      200  {object}  InvAdjGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments [get]
func (h *handler) GetAllInvAdj(c echo.Context) (err error) {
	//per_page, page_no, sort
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllInvAdj(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "InvAdj Retrieved Successfully", result, p)
}

// DeleteInventoryAdjustment godoc
// @Summary      delete a DeleteInventoryAdjustment
// @Description  delete a DeleteInventoryAdjustment
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  InvAdjDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/{id}/delete [delete]
func (h *handler) DeleteInvAdj(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteInvAdj(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Inventory Adjustment deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteInventoryAdjustmentLines godoc
// @Summary      delete a DeleteInventoryAdjustmentLines
// @Description  delete a DeleteInventoryAdjustmentLines
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        inv_adj_id     path      int    true  "inv_adj_id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  InvAdjLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/{id}/delete_products [delete]
func (h *handler) DeleteInvAdjLines(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	product_id := c.QueryParam("product_id")
	prod_id, _ := strconv.Atoi(product_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	query := map[string]interface{}{
		"inv_adj_id": uint(id),
		"product_id": uint(prod_id),
	}
	err = h.service.DeleteInvAdjLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Inventory Adjustment Line Items successfully deleted", query)
}

// SendMailInvAdj godoc
// @Summary      Send Email to specific recevier for Id
// @Description  Send Email to specific recevier for Id
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Param        receiver_email query     string true "receiver_email"
// @Success      200  {object}  SendMailInvAdjResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/{id}/sendemail [get]
func (h *handler) SendMailInvAdj(c echo.Context) (err error) {

	q := new(SendMailInvAdj)
	c.Bind(q)
	// err = h.service.SendMailInvAdj(q)
	// if err != nil {
	// 	return res.RespError(c, err)
	// }

	return res.RespSuccess(c, "Email Sent Successfully", q)
}

// DownloadPdfInvAdj godoc
// @Summary      Download PDF
// @Description  Download PDF
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Success      200  {object}  DownloadPdfInvAdjResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/{id}/printpdf [get]
func (h *handler) DownloadPdfInvAdj(c echo.Context) (err error) {
	return res.RespSuccess(c, "Pdf Generated Successfully", "https://drive.google.com/file/d/1O46B3lgliUGakL4weBjrovCHphTzR5Np/view?usp=sharing")
}

// FavouriteInvAdj godoc
// @Summary FavouriteInvAdj
// @Description FavouriteInvAdj
// @Tags Inventory Adjustments
// @Accept  json
// @Produce  json
// @param id path string true "Inventory Adjustment ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/inventory_adjustments/{id}/favourite [post]
func (h *handler) FavouriteInvAdj(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	fmt.Println("handler ------------------------ -1", query)
	err = h.base_service.FavouriteInvAdj(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Inventory Adjustment is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteInvAdj godoc
// @Summary UnFavouriteInvAdj
// @Description UnFavouriteInvAdj
// @Tags Inventory Adjustments
// @Accept  json
// @Produce  json
// @param id path string true "Inventory Adjustment ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/inventory_adjustments/{id}/unfavourite [post]
func (h *handler) UnFavouriteInvAdj(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteInvAdj(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Inventory Adjustment is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetAllInvAdjDropDown godoc
// @Summary      get all GetAllInvAdjDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllInvAdjDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         Inventory Adjustments
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		query	  string false "filters"
// @Param        per_page 		query     int    false "per_page"
// @Param  	     page_no        query     int    false "page_no"
// @Param        sort		    query     string false "sort"
// @Success      200  {object}  InvAdjGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/inventory_adjustments/dropdown [get]
func (h *handler) GetAllInvAdjDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllInvAdj(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "InvAdj Retrieved Successfully", result, p)
}
