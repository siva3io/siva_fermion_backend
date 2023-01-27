package uom

import (
	"strconv"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
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

var UomHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if UomHandler != nil {
		return UomHandler
	}
	service := NewService()
	UomHandler = &handler{service}
	return UomHandler
}

// CreateUom godoc
// @Summary CreateUom
// @Description CreateUom
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomRequestDTO true "Uom Request Body"
// @Success 200 {object} UomRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/create [post]
func (h *handler) CreateUomEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("uom"),
	}

	eda.Produce(eda.CREATE_UOM, request_payload)
	return res.RespSuccess(c, "Uom Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateUom(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(mdm.Uom)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.CreateUom(edaMetaData, createPayload)

	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_UOM_ACK, responseMessage)
		return
	}
	// updating cache
	UpdateUomInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_UOM_ACK, responseMessage)

}

// CreateUomClass godoc
// @Summary CreateUomClass
// @Description CreateUomClass
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomClassRequestDTO true "Uom Class Request Body"
// @Success 200 {object} UomClassRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/create [post]
func (h *handler) CreateUomClassEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("uom_class"),
	}
	eda.Produce(eda.CREATE_UOM_CLASS, request_payload)
	return res.RespSuccess(c, "Uom class Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateUomClass(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})

	createPayload := new(mdm.UomClass)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.CreateUomClass(edaMetaData, createPayload)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_UOM_CLASS_ACK, responseMessage)
		return
	}

	// updating cache
	UpdateUomClassInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_UOM_CLASS_ACK, responseMessage)

}

// UpdateUom godoc
// @Summary UpdateUom
// @Description UpdateUom
// @Tags UOM
// @Accept  json
// @Produce  json
// @param id path string true "UOM ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomRequestDTO true "Uom Request Body"
// @Success 200 {object} UomRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/{id}/update [post]
func (h *handler) UpdateUomEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	uomId := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         uomId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("uom"),
	}

	eda.Produce(eda.UPDATE_UOM, request_payload)
	return res.RespSuccess(c, "Uom Update Inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateUom(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(mdm.Uom)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.UpdateUom(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_UOM_ACK, responseMessage)
		return
	}

	// updating cache
	UpdateUomInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_UOM_ACK, responseMessage)
}

// UpdateUomClass godoc
// @Summary UpdateUomClass
// @Description UpdateUomClass
// @Tags UOM
// @Accept  json
// @Produce  json
// @param id path string true "UOM class ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomClassRequestDTO true "Uom class Request Body"
// @Success 200 {object} UomClassRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/{id}/update [post]
func (h *handler) UpdateUomClassEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	uonClassId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         uonClassId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("uom_class"),
	}
	eda.Produce(eda.UPDATE_UOM_CLASS, request_payload)
	return res.RespSuccess(c, "Uom class Update Inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateUomClass(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	updatePayload := new(mdm.UomClass)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.UpdateUomClass(edaMetaData, updatePayload)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_UOM_CLASS_ACK, responseMessage)
		return
	}

	// updating cache
	UpdateUomClassInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_UOM_CLASS_ACK, responseMessage)

}

// DeleteUom godoc
// @Summary DeleteUom
// @Description DeleteUom
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/{id}/delete [delete]
func (h *handler) DeleteUom(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	uomId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         uomId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteUom(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// updating cache
	UpdateUomInCache(metaData)

	return res.RespSuccess(c, "Record deleted successfully", uomId)
}

// DeleteUomClass godoc
// @Summary DeleteUomClass
// @Description DeleteUomClass
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM class ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/{id}/delete [delete]
func (h *handler) DeleteUomClass(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	uomClassId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         uomClassId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteUomClass(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// updating cache
	UpdateUomClassInCache(metaData)

	return res.RespSuccess(c, "Record deleted successfully", uomClassId)
}

// UomView godoc
// @Summary View UomView
// @Summary View UomView
// @Description View UomView
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM ID"
// @Success 200 {object} UomViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/{id} [get]
func (h *handler) GetUom(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	uomId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": uomId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetUom(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response UomResponseDTO
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Uom Details retried succesfully", response)
}

// UomClassView godoc
// @Summary View UomClassView
// @Summary View UomClassView
// @Description View UomClassView
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM class ID"
// @Success 200 {object} UomClassViewtDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/{id} [get]
func (h *handler) GetUomClass(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	uomClassId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": uomClassId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetUomClass(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response UomClassResponseDTO
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccess(c, "UomClass Details retried succesfully", result)
}

// UomList godoc
// @Summary UomList
// @Description UomList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom [get]
func (h *handler) GetUomList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []UomResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetUomFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetUomList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, " List Retrieved successfully", response, p)
}

// UomListDropdown godoc
// @Summary UomList
// @Description UomList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/dropdown [get]
func (h *handler) GetUomListDropdown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []UomResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetUomFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetUomList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, " List Retrieved successfully", response, p)
}

// UomClassList godoc
// @Summary UomClassList
// @Description UomClassList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomClassListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class [get]
func (h *handler) GetUomClassList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []UomClassResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetUomClassFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetUomClassList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, " List Retrieved successfully", response, p)
}

// UomClassListDropdown godoc
// @Summary UomClassList
// @Description UomClassList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomClassListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/dropdown [get]
func (h *handler) GetUomClassListDropdown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []UomResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetUomClassFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetUomClassList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, " List Retrieved successfully", response, p)
}
