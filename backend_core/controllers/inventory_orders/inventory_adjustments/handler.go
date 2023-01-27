package inventory_adjustments

import (
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	inventory_orders_base "fermion/backend_core/controllers/inventory_orders/base"
	"fermion/backend_core/internal/model/core"
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

var InventoryAdjustmentsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if InventoryAdjustmentsHandler != nil {
		return InventoryAdjustmentsHandler
	}
	service := NewService()
	base_service := inventory_orders_base.NewServiceBase()
	InventoryAdjustmentsHandler = &handler{service, base_service}
	return InventoryAdjustmentsHandler
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
func (h *handler) CreateInvAdjEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("inventory_adjustments"),
	}
	eda.Produce(eda.CREATE_INVENTORY_ADJUSTMENTS, request_payload)
	return res.RespSuccess(c, "Inventory adjustments creation inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateInvAdj(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	createPayload := new(inventory_orders.InventoryAdjustments)
	helpers.JsonMarshaller(data, createPayload)
	err := h.service.CreateInvAdj(edaMetaData, createPayload)

	var responseMessage eda.ConsumerResponse

	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_INVENTORY_ADJUSTMENTS_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateInvAdjInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_INVENTORY_ADJUSTMENTS_ACK, responseMessage)
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
	data := new([]inventory_orders.InventoryAdjustments)
	c.Bind(&data)
	var metaData core.MetaData
	// access_template_id := c.Get("AccessTemplateId").(string)
	// t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = &metaData.TokenUserId
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateInvAdj(metaData, data)
	if err != nil {
		return res.RespError(c, err)
	}
	UpdateInvAdjInCache(metaData)
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
func (h *handler) UpdateInvAdjEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	invAdjId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         invAdjId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("inventory_adjustments"),
	}
	eda.Produce(eda.UPDATE_INVENTORY_ADJUSTMENTS, request_payload)
	return res.RespSuccess(c, "Inventory Adjustments Update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateInvAdj(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(inventory_orders.InventoryAdjustments)
	helpers.JsonMarshaller(data, updatePayLoad)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdateInvAdj(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_INVENTORY_ADJUSTMENTS_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateInvAdjInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_INVENTORY_ADJUSTMENTS_ACK, responseMessage)
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
	metaData := c.Get("MetaData").(core.MetaData)
	invAdjId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": invAdjId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetInvAdj(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response InvAdjGet
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Inventory Adjustments Details Retrieved successfully", response)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []InvAdjGetAll
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetInvAdjFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllInvAdj(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Inventory Adjustments list Retrieved Successfully", response, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	invAdjId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         invAdjId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteInvAdj(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// cache implementation
	UpdateInvAdjInCache(metaData)
	return res.RespSuccess(c, "Record deleted successfully", invAdjId)
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
	metaData := c.Get("MetaData").(core.MetaData)
	invAdjId := c.Param("id")
	product_id := c.QueryParam("product_id")
	prod_id, _ := strconv.Atoi(product_id)
	metaData.Query = map[string]interface{}{
		"id":         invAdjId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
		"inv_adj_id": invAdjId,
		"product_id": uint(prod_id),
	}
	err = h.service.DeleteInvAdjLines(metaData)
	if err != nil {
		return res.RespError(c, err)
	}
	UpdateInvAdjInCache(metaData)
	return res.RespSuccess(c, "Inventory Adjustment Line Items successfully deleted", metaData.Query)
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
	metaData := c.Get("MetaData").(core.MetaData)
	invAdjId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": invAdjId,
	}
	_, err = h.service.GetInvAdj(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	query := map[string]interface{}{
		"ID":      invAdjId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.FavouriteInvAdj(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Inventory Adjustments is Marked as Favourite", map[string]string{"credinote_id": invAdjId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	invAdjId := c.Param("id")

	query := map[string]interface{}{
		"ID":      invAdjId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.UnFavouriteInvAdj(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Credit Note is Marked as UnFavourite", map[string]string{"invadj_id": invAdjId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []InvAdjGetAll
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetInvAdjFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}
	result, err := h.service.GetAllInvAdj(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "inventory adjustments dropdown list Retrieved Successfully", response, p)
}
