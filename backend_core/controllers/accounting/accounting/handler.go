package accounting

import (
	"strconv"

	"fermion/backend_core/internal/model/accounting"
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

var AccountingHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if AccountingHandler != nil {
		return AccountingHandler
	}
	service := NewService()
	AccountingHandler = &handler{service}
	return AccountingHandler
}

func (h *handler) CreateAccounting(c echo.Context) (err error) {

	var data accounting.UserAccountingLink
	err = c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.CreateAccounting(&data)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateUserAccountingLinkInCache(token_id, access_template_id)

	return res.RespSuccess(c, "account created successfully", map[string]interface{}{"created_id": data.ID})
}

func (h *handler) UpdateAccounting(c echo.Context) (err error) {

	var data accounting.UserAccountingLink
	err = c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)

	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	query := map[string]interface{}{
		"id":         id,
		"created_by": user_id,
	}

	err = h.service.UpdateAccounting(query, &data)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateUserAccountingLinkInCache(token_id, access_template_id)

	return res.RespSuccess(c, "account created successfully", map[string]interface{}{"updated_id": id})
}

func (h *handler) GetAuthKeys(c echo.Context) (err error) {

	accounting_code := c.Param("accounting_code")
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"accounting_code": accounting_code,
		"created_by":      user_id,
	}
	result, err := h.service.GetAuthKeys(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}
