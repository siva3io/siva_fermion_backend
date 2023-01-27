package debitnote

import (
	"strconv"

	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"

	accounting_base "fermion/backend_core/controllers/accounting/base"
	"fermion/backend_core/controllers/eda"
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
	base_service accounting_base.ServiceBase
}

var DebitNoteHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if DebitNoteHandler != nil {
		return DebitNoteHandler
	}
	service := NewService()
	base_service := accounting_base.NewServiceBase()
	DebitNoteHandler = &handler{service, base_service}
	return DebitNoteHandler
}

// CreateDebitNote godoc
// @Summary CreateDebitNote
// @Description CreateDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body DebitNoteDTO true "Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/create [post]
func (h *handler) CreateDebitNoteEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("debit_note"),
	}
	eda.Produce(eda.CREATE_DEBITNOTE, request_payload)
	return res.RespSuccess(c, "Debit note creation inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateDebitNote(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	debitNotePayload := new(accounting.DebitNote)
	helpers.JsonMarshaller(data, debitNotePayload)
	err := h.service.CreateDebitNote(edaMetaData, debitNotePayload)
	debitNotePayload.CreatedByID = &edaMetaData.TokenUserId
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_DEBITNOTE_ACK, responseMessage)
		return
	}
	//cache implementation
	UpdateDebitNoteInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": debitNotePayload.ID,
	}
	// eda.Produce(eda.CREATE_DEBITNOTE_ACK, responseMessage)
}

// UpdateDebitNote godoc
// @Summary UpdateDebitNote
// @Description UpdateDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @param id path string true "Debit Note ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body DebitNoteDTO true "DebitNote  Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/update [post]
func (h *handler) UpdateDebitNoteEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	DebitNoteId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         DebitNoteId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("debit_note"),
	}
	eda.Produce(eda.UPDATE_DEBITNOTE, request_payload)
	return res.RespSuccess(c, "Debit Note Update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateDebitNote(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(accounting.DebitNote)
	helpers.JsonMarshaller(data, updatePayLoad)
	err := h.service.UpdateDebitNote(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_DEBITNOTE_ACK, responsemessage)
		return
	}
	// cache implementation
	UpdateDebitNoteInCache(edaMetaData)
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_DEBITNOTE_ACK, responsemessage)

}

// DebitNoteView godoc
// @Summary View DebitNote
// @Summary View DebitNote
// @Description View DebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "DebitNote ID"
// @Success 200 {object} DebitNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id} [get]
func (h *handler) GetDebitNote(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	debitNoteId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         debitNoteId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetDebitNote(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response DebitNoteResponseDTO
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "DebitNote Note Details Retrieved successfully", response)
}

// GetAllDebitNote godoc
// @Summary Get all Debit note
// @Description Get all Debit note
// @Tags DebitNote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} DebitNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote [get]
func (h *handler) GetAllDebitNote(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []DebitNoteDTO
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetDebitNoteFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetDebitNoteList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Debit Note list Retrieved Successfully", response, p)
}

// GetAllDebitNoteDropDown godoc
// @Summary Get Debit note dropdown list
// @Description Get Debit note dropdown list
// @Tags DebitNote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} DebitNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/dropdown [get]
func (h *handler) GetAllDebitNoteDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []DebitNoteResponseDTO
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetDebitNoteFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}
	result, err := h.service.GetDebitNoteList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Debit Note dropdown list Retrieved Successfully", response, p)

}

// DeleteDebitNote godoc
// @Summary DeleteDebitNote
// @Description DeleteDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "DebitNote ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/delete [delete]
func (h *handler) DeleteDebitNote(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	debitNoteId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         debitNoteId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteDebitNote(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// cache implementation
	UpdateDebitNoteInCache(metaData)
	return res.RespSuccess(c, "Record deleted successfully", debitNoteId)
}

// FavouriteDebitNote godoc
// @Summary FavouriteDebitNote
// @Description FavouriteDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @param id path string true "DebitNote ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/favourite [post]
func (h *handler) FavouriteDebitNote(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	debitNoteId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": debitNoteId,
	}
	_, err = h.service.GetDebitNote(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	query := map[string]interface{}{
		"ID":      debitNoteId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.FavouriteDebitNote(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Debit Note is Marked as Favourite", map[string]string{"debitnote_id": debitNoteId})
}

// UnFavouriteDebitNote godoc
// @Summary UnFavouriteDebitNote
// @Description UnFavouriteDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @param id path string true "Contact ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/unfavourite [post]
func (h *handler) UnFavouriteDebitNote(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	debitNoteId := c.Param("id")

	query := map[string]interface{}{
		"ID":      debitNoteId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.UnFavouriteCreditNote(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Debit Note is Marked as UnFavourite", map[string]string{"debitnote_id": debitNoteId})
}

// GetDebitNoteTab godoc
// @Summary GetDebitNoteTab
// @Summary GetDebitNoteTab
// @Description GetDebitNoteTab
// @Tags DebitNote
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
// @Router /api/v1/debitnote/{id}/filter_module/{tab} [get]
func (h *handler) GetDebitNoteTab(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}
	id := c.Param("id")
	tab := c.Param("tab")
	metaData.AdditionalFields = map[string]interface{}{
		"debitnote_id": id,
		"tab_name":     tab,
	}
	data, err := h.service.GetDebitNoteTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
