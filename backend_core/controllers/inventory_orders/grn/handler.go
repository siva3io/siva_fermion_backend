package grn

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

var GrnHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if GrnHandler != nil {
		return GrnHandler
	}
	service := NewService()
	base_service := inventory_orders_base.NewServiceBase()
	GrnHandler = &handler{service, base_service}
	return GrnHandler
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
func (h *handler) CreateGRNEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("grn"),
	}
	eda.Produce(eda.CREATE_GRN, request_payload)
	return res.RespSuccess(c, "Grn creation inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateGRN(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	createPayload := new(inventory_orders.GRN)
	helpers.JsonMarshaller(data, createPayload)
	createPayload.CreatedByID = &edaMetaData.TokenUserId
	err := h.service.CreateGRN(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_GRN_ACK, responseMessage)
		return
	}

	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_GRN_ACK, responseMessage)
	// cache implementation
	UpdateGrnInCache(edaMetaData)
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
	var metaData core.MetaData
	UpdateGrnInCache(metaData)
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
	metaData := c.Get("MetaData").(core.MetaData)
	grnId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": grnId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetGRN(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// var response GRNGetResponse
	// helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Grn Details Retrieved successfully", result)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []GetAllGRN
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetGrnFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllGRN(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	//helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Crn list Retrieved Successfully", result, p)
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
func (h *handler) UpdateGRNEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	grnId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         grnId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("grn"),
	}
	eda.Produce(eda.UPDATE_GRN, request_payload)
	return res.RespSuccess(c, "Credit Note Update inprogress", edaMetaData.RequestId)

}

func (h *handler) UpdateGRN(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(inventory_orders.GRN)
	helpers.JsonMarshaller(data, updatePayLoad)

	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdateGRN(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_ASN_ACK, responsemessage)
		return
	}
	// cache implementation
	UpdateGrnInCache(edaMetaData)
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_GRN_ACK, responsemessage)

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
	metaData := c.Get("MetaData").(core.MetaData)
	grnId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         grnId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteGRN(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// cache implementation
	UpdateGrnInCache(metaData)
	return res.RespSuccess(c, "Record deleted successfully", grnId)
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
	metaData := c.Get("MetaData").(core.MetaData)
	grnId := c.Param("id")
	query := map[string]interface{}{
		"grn_id":     grnId,
		"product_id": c.QueryParam("product_id"),
		"id":         grnId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteGRNOrderLines(metaData)
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
	metaData := c.Get("MetaData").(core.MetaData)
	grnId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": grnId,
	}
	_, err = h.service.GetGRN(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	query := map[string]interface{}{
		"ID":         grnId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}

	err = h.base_service.FavouriteGrn(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Credit Note is Marked as Favourite", map[string]string{"grn_id": grnId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	grnId := c.Param("id")

	query := map[string]interface{}{
		"ID":      grnId,
		"user_id": metaData.TokenUserId,

		"company_id": metaData.CompanyId,
	}

	err = h.base_service.UnFavouriteGrn(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Credit Note is Marked as UnFavourite", map[string]string{"grn_id": grnId})
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

	metaData := c.Get("MetaData").(core.MetaData)
	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}
	id := c.Param("id")
	tab := c.Param("tab")
	metaData.AdditionalFields = map[string]interface{}{
		"grn_id":   id,
		"tab_name": tab,
	}
	data, err := h.service.GetGrnTab(metaData, page)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []GetAllGRNResponse
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetGrnFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		//helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}
	result, err := h.service.GetAllGRN(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	//helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Grn dropdown list Retrieved Successfully", result, p)
}
