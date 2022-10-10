package transactions

import (
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/payments"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/labstack/echo/v4"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type handler struct {
	service Service
}

func NewHandler() *handler {
	service := NewService()
	return &handler{service}
}

func (h *handler) CreateTransactions(c echo.Context) error {

	input_data := c.Get("transactions").(*payments.Transactions)
	s := c.Get("TokenUserID").(string)
	input_data.CreatedByID = helpers.ConvertStringToUint(s)
	err := h.service.CreateTransaction(input_data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, " Created", map[string]interface{}{"transaction_id": input_data.ID})
}

func (h *handler) ListWalletTransactions(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	resp, err := h.service.ListWalletTransactions(page)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "list of transactions fetched", resp, page)
}

func (h *handler) ListGatewayTransactions(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	resp, err := h.service.ListGatewayTransactions(page)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "list of transactions fetched", resp, page)
}

func (h *handler) ViewTransactions(c echo.Context) error {

	var query = make(map[string]interface{}, 0)

	ID := c.Param("id")
	fmt.Println("ID", ID)

	// id, _ := strconv.Atoi(ID)

	query["id"] = ID

	fmt.Println("query", query)

	resp, err := h.service.ViewTransaction(query)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "transaction fetched", resp)
}

func (h *handler) UpdateTransactions(c echo.Context) (err error) {
	var id = c.Param("id")

	var query = make(map[string]interface{}, 0)
	query["id"] = id

	fmt.Println(query)

	data := c.Get("transactions").(*payments.Transactions)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateTransaction(query, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "transaction updated succesfully", nil)
}

func (h *handler) DeleteTransactions(c echo.Context) (err error) {
	var id = c.Param("id")

	var query = make(map[string]interface{})

	query["id"] = id
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query["created_by"] = user_id
	err = h.service.DeleteTransaction(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "transaction deleted successfully", map[string]string{"transaction_id": id})
}

func (h *handler) GetTransaction(c echo.Context) error {
	query_params_map := c.QueryParams()
	query_params := map[string]interface{}{}
	for key, value := range query_params_map {
		query_params[key] = value[0]
	}
	// fmt.Println(query_params)
	resp, err := h.service.GetTransaction(query_params)
	if err != nil {
		return res.SuccessWithError(c, err.Error(), nil)
	}
	return res.RespSuccess(c, "customer fetched", resp)
}
