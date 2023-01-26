package customers

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

func (h *handler) CreateCustomer(c echo.Context) error {

	input_data := c.Get("customers").(*payments.Customers)
	s := c.Get("TokenUserID").(string)
	input_data.CreatedByID = helpers.ConvertStringToUint(s)
	input_data.UserId = *input_data.CreatedByID
	err := h.service.CreateCustomer(input_data)
	if err != nil {
		return res.SuccessWithError(c, err.Error(), nil)
	}
	return res.RespSuccess(c, "customer Created", map[string]interface{}{"customer_id": input_data.ID})
}

func (h *handler) ListCustomers(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	resp, err := h.service.ListCustomers(page)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "list of customers fetched", resp, page)
}

func (h *handler) ViewCustomer(c echo.Context) error {

	var query = make(map[string]interface{}, 0)

	ID := c.Param("id")
	fmt.Println("ID", ID)
	id, _ := strconv.Atoi(ID)
	query["id"] = id
	fmt.Println("query", query)
	resp, err := h.service.ViewCustomer(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "customer fetched", resp)
}

func (h *handler) UpdateCustomer(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("customers").(*payments.Customers)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateCustomer(query, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "customer updated succesfully", nil)
}

func (h *handler) DeleteCustomer(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{})
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query["created_by"] = user_id
	err = h.service.DeleteCustomer(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "customer deleted successfully", map[string]string{"customer_id": id})
}

func (h *handler) GetCustomer(c echo.Context) error {
	query_params_map := c.QueryParams()
	query_params := map[string]interface{}{}
	for key, value := range query_params_map {
		query_params[key] = value[0]
	}
	// fmt.Println(query_params)
	resp, err := h.service.GetCustomer(query_params)
	if err != nil {
		return res.SuccessWithError(c, err.Error(), nil)
	}
	return res.RespSuccess(c, "customer fetched", resp)
}
