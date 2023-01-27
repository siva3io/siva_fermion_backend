package asn

import (
	"errors"
	"fmt"
	"strconv"

	concurrency_management "fermion/backend_core/controllers/concurrency_management"
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
	service            Service
	base_service       inventory_orders_base.ServiceBase
	concurrencyService concurrency_management.Service
}

var AsnHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if AsnHandler != nil {
		return AsnHandler
	}
	service := NewService()
	base_service := inventory_orders_base.NewServiceBase()
	concurrencyService := concurrency_management.NewService()
	AsnHandler = &handler{service, base_service, concurrencyService}
	return AsnHandler
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
func (h *handler) CreateAsnEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("asn"),
	}
	eda.Produce(eda.CREATE_ASN, request_payload)
	return res.RespSuccess(c, "Asn creation inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateAsn(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	createPayload := new(inventory_orders.ASN)
	helpers.JsonMarshaller(data, createPayload)
	createPayload.CreatedByID = &edaMetaData.TokenUserId
	err := h.service.CreateAsn(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_ASN_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateAsnInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_ASN_ACK, responseMessage)
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
	data := new([]inventory_orders.ASN)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	var metaData core.MetaData
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateAsn(metaData, data)
	if err != nil {
		return res.RespError(c, err)
	}
	UpdateAsnInCache(metaData)
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
func (h *handler) UpdateAsnEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	asnId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         asnId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("asn"),
	}
	eda.Produce(eda.UPDATE_ASN, request_payload)
	return res.RespSuccess(c, "Asn Update inprogress", edaMetaData.RequestId)
	//--------------------------concurrency_management_start-----------------------------------------
	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus((*helpers.ConvertStringToUint(asnId)), "status_id", "asns")
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(*helpers.ConvertStringToUint(asnId), PreviousStatusId, "status_id", "asn")
	//--------------------------concurrency_management_end-----------------------------------------
	return res.RespSuccess(c, "Asn update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateAsn(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(inventory_orders.ASN)
	helpers.JsonMarshaller(data, updatePayLoad)

	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdateAsn(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_ASN_ACK, responsemessage)
		return
	}
	// cache implementation
	UpdateAsnInCache(edaMetaData)
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_ASN_ACK, responsemessage)
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
	metaData := c.Get("MetaData").(core.MetaData)
	asnId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": asnId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetAsn(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// var response AsnGet
	// helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Credit Note Details Retrieved successfully", result)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []AsnGetAllResponse
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetAsnFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllAsn(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	// helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Asn list Retrieved Successfully", result, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	asnId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         asnId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteAsn(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// cache implementation
	UpdateAsnInCache(metaData)
	return res.RespSuccess(c, "Record deleted successfully", asnId)
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
	product_id := c.QueryParam("product_id")
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	prod_id, _ := strconv.Atoi(product_id)
	query := map[string]interface{}{
		"asn_id":     ID,
		"product_id": uint(prod_id),
	}
	var metaData core.MetaData
	err = h.service.DeleteAsnLines(metaData)
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
	metaData := c.Get("MetaData").(core.MetaData)
	asnId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         asnId,
		"company_id": metaData.CompanyId,
	}
	_, err = h.service.GetAsn(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	query := map[string]interface{}{
		"ID":      asnId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.FavouriteAsn(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Asn Ais Marked as Favourite", map[string]string{"id": asnId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	asnId := c.Param("id")

	query := map[string]interface{}{
		"ID":      asnId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.UnFavouriteAsn(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "asn is Marked as UnFavourite", map[string]string{"credinote_id": asnId})
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
	var metaData core.MetaData
	metaData.AdditionalFields = map[string]interface{}{
		"asn":      id,
		"tab_name": tab,
	}
	data, err := h.service.GetAsnTab(metaData, page)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetAsnFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, p)
	}
	//var data []AsnGetAllResponse
	response, err := h.service.GetAllAsn(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	//helpers.JsonMarshaller(response, &data)
	return res.RespSuccessInfo(c, "dropdown list of asn fetched", response, p)
}
