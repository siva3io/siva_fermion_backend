package omnichannel_fields

import (
	"strconv"

	"fermion/backend_core/internal/model/pagination"
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

var OmnichannelFieldsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if OmnichannelFieldsHandler != nil {
		return OmnichannelFieldsHandler
	}
	service := NewService()
	OmnichannelFieldsHandler = &handler{service}
	return OmnichannelFieldsHandler
}

// CreateOmnichannelFields godoc
// @Summary Create OmnichannelFields
// @Description Create OmnichannelFields
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body OmnichannelFieldRequestDto true "OmnichannelField Request Body"
// @Success 200 {object} OmnichannelFieldResponseDto
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/create [post]
func (h *handler) CreateOmnichannelFields(c echo.Context) error {

	input_data := c.Get("omnichannel_fields").(*OmnichannelFieldRequestDto)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	id, err := h.service.CreateOmnichannelField(input_data, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Omnichannel Field Created", map[string]interface{}{"omnichannel_field_id": id})
}

// GetListOmnichannelFields godoc
// @Summary Get all OmnichannelFields list
// @Description Get all OmnichannelFields list
// @Tags OmnichannelFields
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllOmnichannelFieldsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields [get]
func (h *handler) ListOmnichannelFields(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data, err := h.service.ListOmnichannelFields(page, access_template_id, token_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "list of omnichannel fields fetched", data, page)
}

// ViewOmnichannelFields godoc
// @Summary View OmnichannelFields
// @Summary View OmnichannelFields
// @Description View OmnichannelFields
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @Success 200 {object} OmnichannelFieldResponseDto
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/{id} [get]
func (h *handler) ViewOmnichannelFields(c echo.Context) error {

	var query = make(map[string]interface{}, 0)
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	query["id"] = int(id)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.ViewOmnichannelField(query, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "omnichannel fields fetched", resp)
}

// UpdateOmnichannelFields godoc
// @Summary Update OmnichannelFields
// @Description Update OmnichannelFields
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @param RequestBody body OmnichannelFieldRequestDto true "OmnichannelField Request Body"
// @Success 200 {object} OmnichannelFieldResponseDto
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/{id}/update [post]
func (h *handler) UpdateOmnichannelFields(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("omnichannel_fields").(*OmnichannelFieldRequestDto)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.UpdateOmnichannelField(query, data, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "omnichannel fields updated succesfully", map[string]interface{}{"omnichannel_field_id": ID})
}

// DeleteOmnichannelFields godoc
// @Summary Delete OmnichannelFields
// @Description Delete OmnichannelFields
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "field_id"
// @Success 200 {object} OmnichannelFieldDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/{id}/delete [delete]
func (h *handler) DeleteOmnichannelFields(c echo.Context) (err error) {
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{})
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.DeleteOmnichannelField(query, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "omnichannel fields deleted successfully", map[string]int{"deleted_id": ID})
}

// ViewAppFields godoc
// @Summary View AppFields
// @Summary View AppFields
// @Description View AppFields
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   app_id 					query 	int 	false "app_id"
// @Param   channel_type_id			query   int 	false "channel_type_id"
// @Param   channel_function_id 	query   int    	false "channel_function_id"
// @Success 200 {object} ViewAppFieldsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/app_fields [get]
func (h *handler) ViewAppFields(c echo.Context) error {

	query := c.Get("app_field_query").(*ViewAppFieldsQueryDto)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	accessToken := c.Request().Header.Get("Authorization")

	resp, err := h.service.ViewAppFields(*query, accessToken, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "omnichannel app fields fetched", resp)
}

// UpsertOmnichannelFieldsData godoc
// @Summary Upsert OmnichannelFieldsData
// @Description Upsert OmnichannelFieldsData
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body OmnichannelFieldDataRequestDto true "OmnichannelFieldDataRequestDto Request Body"
// @Success 200 {object} OmnichannelFieldResponseDto
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/upsert_data [post]
func (h *handler) UpsertOmnichannelFieldsData(c echo.Context) error {

	register, err := strconv.ParseBool(c.QueryParam("register"))
	if err != nil {
		register = true
	}

	input_data := c.Get("omnichannel_fields_data").(*OmnichannelFieldDataRequestDto)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	response, err := h.service.UpsertOmnichannelFieldsData(input_data, access_template_id, token_id, register)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Omnichannel Fields Data Upserted", response)
}

// ViewAppFieldsData godoc
// @Summary View AppFieldsData
// @Summary View AppFieldsData
// @Description View AppFieldsData
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   app_id 					query 	int 	false "app_id"
// @Param   channel_function_id 	query   int    	false "channel_function_id"
// @Success 200 {object} ViewAppFieldsDataResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/app_data [get]
func (h *handler) ViewAppFieldsData(c echo.Context) error {

	query := c.Get("app_field_data_query").(*ViewAppFieldsDataQueryDto)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	resp, err := h.service.ViewAppFieldsData(*query, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "omnichannel app fields data fetched", resp)
}

// GetAppSyncSettings godoc
// @Summary Get AppSyncSettings
// @Description Get AppSyncSettings
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} GetAppSyncSettingsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/{app_code}/sync_settings [get]
func (h *handler) GetAppSyncSettings(c echo.Context) error {

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	appCode := c.Param("app_code")
	data, err := h.service.GetAppSyncSettings(appCode, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", data)
}

// UpsertAppSyncSettings godoc
// @Summary Upsert AppSyncSettings
// @Description Upsert AppSyncSettings
// @Tags OmnichannelFields
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body OmnichannelSyncSettingsRequestDto true "OmnichannelSyncSettingsRequestDto Request Body"
// @Success 200 {object} UpsertAppSyncSettingsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/omnichannel_fields/{app_code}/upsert_sync_settings [post]
func (h *handler) UpsertAppSyncSettings(c echo.Context) error {

	input_data := c.Get("omnichannel_sync_settings").(*[]OmnichannelSyncSettingsRequestDto)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	appCode := c.Param("app_code")
	err := h.service.UpsertAppSyncSettings(appCode, *input_data, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Omnichannel App Sync Settings Data Upserted", nil)
}

func (h *handler) GetAppDataFilter(c echo.Context) error {

	query := c.Get("get_app_data_filter_query").(*GetAppDataFilterQueryDto)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	resp, err := h.service.GetAppDataFilter(*query, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "omnichannel app fields data fetched", resp)
}
