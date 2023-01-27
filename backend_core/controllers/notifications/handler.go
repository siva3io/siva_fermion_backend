package notifications

import (
	"fmt"
	"strconv"

	mdm_base "fermion/backend_core/controllers/mdm/base"
	location_service "fermion/backend_core/controllers/mdm/locations"

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
	service          Service
	base_service     mdm_base.ServiceBase
	location_service location_service.Service
}

func NewHandler() *handler {
	service := NewService()
	base_service := mdm_base.NewServiceBase()
	location_service := location_service.NewService()
	return &handler{service, base_service, location_service}
}

// GetAllProducts godoc
// @Summary      ListProducts with MISC OPTIONS of pagination, search,Filter and sort
// @Description  List of All Products  with MISC OPTIONS of pagination, search,Filter and sort
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  TemplateReponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products [get]
func (h *handler) GetAllNotifications(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	// access_template_id := c.Get("AccessTemplateId").(string)

	var response interface{}
	if *p == pagination.BasePaginatevalue {
		response, *p = GetNotificationFromCache(token_id)
		fmt.Println("response", response)
	}

	if response != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	resp, err := h.service.GetAllNotifications(p)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "data retrieved successfully", resp, p)
}

func (h *handler) MarkAsReadNotifications(c echo.Context) (err error) {
	var id = c.Param("id")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	var data = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)

	query["id"] = ID
	data["updated_by"] = helpers.ConvertStringToUint(token_id)
	data["is_read"] = true

	err = h.service.UpdateNotification(data, query)
	if err != nil {
		return res.RespError(c, err)
	}

	// cache implementation
	UpdateNotificationInCache(token_id, access_template_id)

	return res.RespSuccess(c, "Notification updated succesfully", nil)
}

func (h *handler) MarkAllAsReadNotifications(c echo.Context) (err error) {
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	var data = make(map[string]interface{}, 0)

	user_id := helpers.ConvertStringToUint(token_id)
	query["created_by"] = user_id
	data["updated_by"] = user_id
	data["is_read"] = true

	err = h.service.UpdateNotification(data, query)
	if err != nil {
		return res.RespError(c, err)
	}

	// cache implementation
	UpdateNotificationInCache(token_id, access_template_id)

	return res.RespSuccess(c, "Notification updated succesfully", nil)
}
