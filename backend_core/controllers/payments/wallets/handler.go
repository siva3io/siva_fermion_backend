package wallets

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

func (h *handler) CreateWallet(c echo.Context) error {

	input_data := c.Get("wallet").(*payments.Wallets)
	s := c.Get("TokenUserID").(string)
	input_data.CreatedByID = helpers.ConvertStringToUint(s)
	input_data.UserId = *input_data.CreatedByID
	err := h.service.CreateWallet(input_data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "wallet Created", map[string]interface{}{"wallet_id": input_data.ID})
}

func (h *handler) ListWallets(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	resp, err := h.service.ListWallets(page)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "list of wallets fetched", resp, page)
}

func (h *handler) ViewWallet(c echo.Context) error {

	var query = make(map[string]interface{}, 0)

	ID := c.Param("id")
	fmt.Println("ID", ID)
	id, _ := strconv.Atoi(ID)
	query["id"] = id
	fmt.Println("query", query)
	resp, err := h.service.ViewWallet(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "wallet fetched", resp)
}

func (h *handler) UpdateWallet(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("wallet").(*payments.Wallets)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateWallet(query, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "wallet updated succesfully", nil)
}

func (h *handler) DeleteWallet(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{})
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query["created_by"] = user_id
	err = h.service.DeleteWallet(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Wallet deleted successfully", map[string]string{"wallet_id": id})
}

func (h *handler) GetWallet(c echo.Context) error {
	query_params_map := c.QueryParams()
	query_params := map[string]interface{}{}
	for key, value := range query_params_map {
		query_params[key] = value[0]
	}
	// fmt.Println(query_params)
	resp, err := h.service.GetWallet(query_params)
	if err != nil {
		return res.SuccessWithError(c, err.Error(), nil)
	}
	return res.RespSuccess(c, "customer fetched", resp)
}

func (h *handler) AddMoney(c echo.Context) (err error) {
	var query = make(map[string]interface{}, 0)
	s := c.Get("TokenUserID").(string)
	data := c.Get("wallet").(*payments.Wallets)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	ID, _ := strconv.Atoi(s)
	query["user_id"] = ID
	err = h.service.AddMoney(query, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "wallet updated succesfully", nil)
}

func (h *handler) DeductMoney(c echo.Context) (err error) {
	var query = make(map[string]interface{}, 0)
	s := c.Get("TokenUserID").(string)
	data := c.Get("wallet").(*payments.Wallets)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	ID, _ := strconv.Atoi(s)
	query["user_id"] = ID
	err = h.service.DeductMoney(query, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "wallet updated succesfully", nil)
}
