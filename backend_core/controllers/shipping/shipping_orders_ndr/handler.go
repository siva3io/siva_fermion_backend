package shipping_orders_ndr

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

var ShippingOrdersNdrHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ShippingOrdersNdrHandler != nil {
		return ShippingOrdersNdrHandler
	}
	service := NewService()
	ShippingOrdersNdrHandler = &handler{service}
	return ShippingOrdersNdrHandler
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
func (h *handler) CreateNDREvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	data := new(NDRRequest)
	c.Bind(&data)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      data,
	}

	eda.Produce(eda.CREATE_SHIPPING_ORDER_NDR, request_payload)
	return res.RespSuccess(c, "Ndr Creation Inprogress", edaMetaData.RequestId)

}

func (h *handler) CreateNDR(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(shipping.NDR)
	helpers.JsonMarshaller(data, createPayload)
	err := h.service.CreateNDR(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SHIPPING_ORDER_NDR_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{"created_id": createPayload.ID}
	// eda.Produce(eda.CREATE_SHIPPING_ORDER_NDR_ACK, responseMessage)
	// cache implementation
	UpdateNdrInCache(edaMetaData)

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
	var metaData core.MetaData
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateNDR(metaData, data)
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
	metaData := c.Get("MetaData").(core.MetaData)
	ndrId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": ndrId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetNDR(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response NDRRequest
	helpers.JsonMarshaller(result, &response)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []NDRGetResponse

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetNdrFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllNDR(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []NDRRequest

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))

	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetNdrFromCache(tokenUserId)
	}

	if response != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllNDR(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
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
func (h *handler) UpdateNDREvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	ndrId := c.Param("id")
	data := new(NDRRequest)
	c.Bind(&data)
	edaMetaData.Query = map[string]interface{}{
		"id":         ndrId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      data,
	}

	eda.Produce(eda.UPDATE_SHIPPING_ORDER_NDR, request_payload)
	return res.RespSuccess(c, "Shipping_order ndr Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})

}

func (h *handler) UpdateNDR(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	updatePayload := new(shipping.NDR)
	helpers.JsonMarshaller(data, updatePayload)
	err := h.service.UpdateNDR(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_SHIPPING_ORDER_NDR_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{"updated_id": edaMetaData.Query["id"]}
	// eda.Produce(eda.UPDATE_SHIPPING_ORDER_NDR_ACK, responseMessage)
	// cache implementation
	UpdateNdrInCache(edaMetaData)

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
	metaData := c.Get("MetaData").(core.MetaData)
	ndrId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         ndrId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteNDR(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateNdrInCache(metaData)

	return res.RespSuccess(c, "Deleted Successfully", map[string]string{"Deleted_NDR_ID": ndrId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	ndrId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         ndrId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteNDRLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "NDR Line successfully deleted", ndrId)
}
