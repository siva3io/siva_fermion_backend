package shipping_orders_rto

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

var ShippingOrdersRtoHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ShippingOrdersRtoHandler != nil {
		return ShippingOrdersRtoHandler
	}
	service := NewService()
	ShippingOrdersRtoHandler = &handler{service}
	return ShippingOrdersRtoHandler
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
func (h *handler) CreateRTOEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	data := new(RTORequest)
	c.Bind(&data)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      data,
	}
	eda.Produce(eda.CREATE_SHIPPING_ORDER_RTO, request_payload)
	return res.RespSuccess(c, "Shipping_order rto Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateRTO(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(shipping.RTO)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.CreateRTO(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		eda.Produce(eda.CREATE_SHIPPING_ORDER_RTO_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateRtoInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	eda.Produce(eda.CREATE_SHIPPING_ORDER_RTO_ACK, responseMessage)

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
	var metaData core.MetaData
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateRTO(metaData, data)
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
	metaData := c.Get("MetaData").(core.MetaData)
	uomId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": uomId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetRTO(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response GetRTO
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "success", response)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []GetAllRTO
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetRtoFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}
	result, err := h.service.GetAllRTO(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "RTO list Retrieved Successfully", response, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []GetAllRTO
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetRtoFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}
	result, err := h.service.GetAllRTO(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "RTO list Retrieved Successfully", response, p)
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
func (h *handler) UpdateRTOEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	rtoId := c.Param("id")
	data := new(RTORequest)
	c.Bind(&data)
	edaMetaData.Query = map[string]interface{}{
		"id":         rtoId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      data,
	}

	eda.Produce(eda.UPDATE_SHIPPING_ORDER_RTO, request_payload)
	return res.RespSuccess(c, "Shipping_order rtoUpdate Inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateRTO(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(shipping.RTO)
	helpers.JsonMarshaller(data, updatePayload)
	err := h.service.UpdateRTO(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		eda.Produce(eda.UPDATE_SHIPPING_ORDER_RTO_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateRtoInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	eda.Produce(eda.UPDATE_SHIPPING_ORDER_RTO_ACK, responseMessage)

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
	metaData := c.Get("MetaData").(core.MetaData)
	rtoId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         rtoId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteRTO(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateRtoInCache(metaData)

	return res.RespSuccess(c, "Deleted Successfully", rtoId)
}
