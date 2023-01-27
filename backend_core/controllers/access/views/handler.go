package views

import (
	"strconv"

	// view_base "fermion/backend_core/controllers/access/base"
	model_core "fermion/backend_core/internal/model/core"
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
	service Services
}

var AccessViewHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if AccessViewHandler != nil {
		return AccessViewHandler
	}
	service := NewService()
	AccessViewHandler = &handler{service}
	return AccessViewHandler
}

//----------------------------------Views--------------------------------------------------

func (h *handler) RegisterView(c echo.Context) (err error) {

	var data *model_core.ViewsLevelAccessItems

	c.Bind(&data)
	s := c.Get("TokenUserID").(string)
	data.CreatedByID = helpers.ConvertStringToUint(s)
	err = h.service.RegisterViews(data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "View registered successfully", map[string]interface{}{"created_id": data.ID})
}

func (h *handler) UpdateView(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	var data model_core.ViewsLevelAccessItems
	c.Bind(&data)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateViewDetails(uint(id), data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "View details updates successfully", map[string]interface{}{"updated_id": id})
}

func (h *handler) FetchAllViews(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	result, err := h.service.FindAllView(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

func (h *handler) GetView(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	result, err := h.service.GetViewDetail(uint(id))
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

func (h *handler) DeleteView(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	err = h.service.DeleteViewDetail(uint(id), uint(user_id))
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "deleted view details.", map[string]int{"deleted_id": id})
}
