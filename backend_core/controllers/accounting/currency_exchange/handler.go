package currency_exchange

import (
	//"strconv"

	"fmt"
	"net/http"
	"strconv"

	//"fermion/backend_core/internal/model/accounting"
	db "fermion/backend_core/internal/model/accounting"

	//"fermion/backend_core/pkg/util/helpers"

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

func NewHandler() *handler {
	service := NewService()
	return &handler{service}
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
func (h *handler) CreateExchangePair(c echo.Context) error {
	// data := c.Get("currency_exchange").(*accounting.BasicPairInfo)
	data := db.CurrencyExchange{}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	if er := c.Bind(&data); er != nil {
		return res.RespErr(c, er)
	}
	err := h.service.CreateCurrency(&data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Exchange Pair created successfully", map[string]interface{}{"created_id": data.ID})
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
func (h *handler) UpdateExchangePair(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	ID, _ := strconv.Atoi(id)
	query["id"] = uint(ID)
	var data db.CurrencyExchange
	err = c.Bind(&data)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	err = h.service.UpdateExchangePair(query, &data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Exchange Pair updated successfully", map[string]interface{}{"updated_id": id})

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
	var result []db.CurrencyExchange
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err = h.service.GetExchangePair(token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)

	}
	return c.JSON(http.StatusOK, result)
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
	var result []db.CurrencyExchange
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err = h.service.GetExchangePair(token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)

	}
	return c.JSON(http.StatusOK, result)
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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = map[string]interface{}{
		"id": ID,
	}
	result, err := h.service.GetExchangePairData(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Exchange Pair Details Retrieved successfully", result)

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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteExchangePair(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Exchange Pair Record deleted successfully", map[string]string{"deleted_id": id})
}
