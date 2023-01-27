package shipping_orders_wd

import (
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
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

var ShippingOrdersWdHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ShippingOrdersWdHandler != nil {
		return ShippingOrdersWdHandler
	}
	service := NewService()
	ShippingOrdersWdHandler = &handler{service}
	return ShippingOrdersWdHandler
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
func (h *handler) CreateWDEvent(c echo.Context) (err error) {
	data := new(WDRequest)
	c.Bind(&data)
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      data,
	}
	eda.Produce(eda.CREATE_SHIPPING_ORDER_WD, request_payload)
	return res.RespSuccess(c, "Shipping_order wd Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateWD(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(shipping.WD)
	helpers.JsonMarshaller(data, createPayload)
	err := h.service.CreateWD(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SHIPPING_ORDER_WD_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// cache implementation
	UpdateWdInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	eda.Produce(eda.CREATE_SHIPPING_ORDER_WD_ACK, responseMessage)

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
	var metaData core.MetaData
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateWD(metaData, data)
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
	metaData := c.Get("MetaData").(core.MetaData)
	wdId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": wdId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetWD(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response GetWD
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "success", response)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []GetAllWD

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))

	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetWdFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)

	}

	result, err := h.service.GetAllWD(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "WD list Retrieved Successfully", response, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []GetAllWD

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))

	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetWdFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)

	}

	result, err := h.service.GetAllWD(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "WD list Retrieved Successfully", response, p)
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
func (h *handler) UpdateWDEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	wdId := c.Param("id")
	data := new(WDRequest)
	c.Bind(&data)
	edaMetaData.Query = map[string]interface{}{
		"id":         wdId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      data,
	}

	eda.Produce(eda.UPDATE_SHIPPING_ORDER_WD, request_payload)
	return res.RespSuccess(c, "Shipping_order wd Update Inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateWD(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(shipping.WD)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.UpdateWD(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_SHIPPING_ORDER_WD_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateWdInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_SHIPPING_ORDER_WD_ACK, responseMessage)

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
	metaData := c.Get("MetaData").(core.MetaData)
	wdId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         wdId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteWD(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateWdInCache(metaData)

	return res.RespSuccess(c, "Deleted Successfully", wdId)
}
