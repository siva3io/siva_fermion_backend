package currency_exchange

import (
	//"strconv"

	"strconv"

	//"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/controllers/eda"
	accounting "fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"

	//"fermion/backend_core/pkg/util/helpers"

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

var CurrencyExchangeHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if CurrencyExchangeHandler != nil {
		return CurrencyExchangeHandler
	}
	service := NewService()
	CurrencyExchangeHandler = &handler{service}
	return CurrencyExchangeHandler
}

// CreateExchangePair godoc
// @Summary CreateExchangePair
// @Description CreateExchangePair
// @Tags CurrencyExchange
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body CurrencyExchangeDTO true "Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/currency_and_exchange/create [post]
func (h *handler) CreateExchangePairEvent(c echo.Context) error {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("currency_exchange"),
	}
	eda.Produce(eda.CREATE_CURRENCY_EXCHANGE, request_payload)
	return res.RespSuccess(c, "Currency Exchange creation inprogress", edaMetaData.RequestId)

}

func (h *handler) CreateExchangePair(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	createPayload := new(accounting.CurrencyExchange)
	helpers.JsonMarshaller(data, createPayload)
	createPayload.CreatedByID = &edaMetaData.TokenUserId
	err := h.service.CreateCurrency(edaMetaData, createPayload)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_CURRENCY_EXCHANGE_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateCurrencyExchangeInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_CURRENCY_EXCHANGE_ACK, responseMessage)
}

// UpdateExchangePair godoc
// @Summary UpdateExchangePair
// @Description UpdateExchangePair
// @Tags CurrencyExchange
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body CurrencyExchangeDTO true "Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/currency_and_exchange/{id}/update [post]
func (h *handler) UpdateExchangePairEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	currencyId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         currencyId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("currency_exchange"),
	}
	eda.Produce(eda.UPDATE_CURRENCY_EXCHANGE, request_payload)
	return res.RespSuccess(c, "Currency Exchange Update inprogress", edaMetaData.RequestId)

}
func (h *handler) UpdateExchangePair(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(accounting.CurrencyExchange)
	helpers.JsonMarshaller(data, updatePayLoad)
	err := h.service.UpdateExchangePair(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_CURRENCY_EXCHANGE_ACK, responsemessage)
		return
	}
	// cache implementation
	UpdateCurrencyExchangeInCache(edaMetaData)
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_CURRENCY_EXCHANGE_ACK, responsemessage)

}

// GetExchangePair godoc
// @Summary Get all CurrencyExchange
// @Description Get all CurrencyExchange
// @Tags CurrencyExchange
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CurrencyExchangeDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/currency_and_exchange [get]
func (h *handler) GetExchangePair(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []CurrencyExchangeDTO
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetCurrencyExchangeFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetExchangePair(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Currency Exchange list Retrieved Successfully", response, p)
}

// GetExchangePairDropDown godoc
// @Summary Get all CurrencyExchangeDropDown
// @Description Get all CurrencyExchangeDropDown
// @Tags CurrencyExchange
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CurrencyExchangeDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/currency_and_exchange/dropdown [get]
func (h *handler) GetExchangePairDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []CurrencyExchangeDTO
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetCurrencyExchangeFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}
	result, err := h.service.GetExchangePair(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Currency Exchange dropdown list Retrieved Successfully", response, p)
}

// GetExchangePairData godoc
// @Summary Get one CurrencyExchange
// @Summary Get one CurrencyExchange
// @Description GetExchangePairData
// @Tags CurrencyExchange
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "CurrencyExchange ID"
// @Success 200 {object} CurrencyExchangeDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/currency_and_exchange/{id} [get]
func (h *handler) GetExchangePairData(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	currencyId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         currencyId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetExchangePairData(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response CurrencyExchangeDTO
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Currency Exchange Details Retrieved successfully", response)

}

// DeleteExchangePair godoc
// @Summary DeleteExchangePair
// @Description DeleteExchangePair
// @Tags CurrencyExchange
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "CurrencyExchange ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/currency_and_exchange/{id}/delete  [delete]
func (h *handler) DeleteExchangePair(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	currencyId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         currencyId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteExchangePair(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// cache implementation
	UpdateCurrencyExchangeInCache(metaData)
	return res.RespSuccess(c, "Record deleted successfully", currencyId)
}
