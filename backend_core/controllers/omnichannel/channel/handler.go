package channel

import (
	"errors"
	"strconv"

	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/repository"
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
	service        Service
	coreRepository repository.Core
}

var ChannelsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ChannelsHandler != nil {
		return ChannelsHandler
	}
	service := NewService()
	coreRepository := repository.NewCore()
	ChannelsHandler = &handler{service, coreRepository}
	return ChannelsHandler
}

// CreateChannel godoc
// @Summary      Do Create Channel
// @Description  Create a Channel
// @Tags         Channel
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router       /api/v1/channels/create [post]
func (h *handler) CreateChannel(c echo.Context) error {
	s := c.Get("TokenUserID").(string)
	var data *ChannelDTO
	err := c.Bind(&data)
	data.CreatedByID = helpers.ConvertStringToUint(s)
	if err != nil {
		return res.RespError(c, err)
	}
	var resp interface{}
	query := map[string]interface{}{"id": data.ChannelTypeId}
	lookupcode, err := h.coreRepository.GetLookupCodes(query, new(pagination.Paginatevalue))
	external_id := data.ExternalId
	external_value := *external_id
	if (err != nil || len(lookupcode) == 0) && external_value != uint(0) {
		return res.RespErr(c, errors.New("invalid Channel type"))
	}
	if len(lookupcode) > 0 {
		if lookupcode[0].LookupCode == "MARKETPLACE" {
			resp, err = h.service.CreateMarketplace(*data)
		} else if lookupcode[0].LookupCode == "WEBSTORE" {
			resp, err = h.service.CreateWebstore(*data)
		} else if lookupcode[0].LookupCode == "CORE" {
			resp, err = h.service.CreateChannel(*data)
		}
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Channel created successfully", map[string]interface{}{"created_id": resp.(*omnichannel.Channel).ID})

}

// GetChannel godoc
// @Summary      Do Channel fetching
// @Description  fetch a Channel
// @Tags         Channel
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router       /api/v1/channels/:id [get]
func (h *handler) GetChannel(c echo.Context) error {

	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	channels, err := h.service.GetChannelService(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	query = map[string]interface{}{"id": channels.ChannelTypeId}
	lookupcode, err := h.coreRepository.GetLookupCodes(query, new(pagination.Paginatevalue))
	external_id := channels.ExternalId
	external_value := *external_id
	if external_value != uint(0) {
		channel, err := h.service.GetChannel(channels, lookupcode[0].LookupCode)

		if err != nil {
			return res.RespErr(c, err)
		}

		return res.RespSuccess(c, "Channel Details Retrieved successfully", channel)
	}
	return res.RespSuccess(c, "Channel Details Retrieved successfully", channels)
}

// UpdateChannel godoc
// @Summary      Do Update Channel
// @Description  Updates a Channel
// @Tags         Channel
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router       /api/v1/channels/:id/update [post]
func (h *handler) UpdateChannel(c echo.Context) (err error) {

	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	var data ChannelDTO

	err = c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	channels, err := h.service.GetChannelService(query)
	external_id := channels.ExternalId
	external_value := *external_id
	var result ChannelDTO
	if external_value != uint(0) {
		query = map[string]interface{}{"id": channels.ChannelTypeId}
		lookupcode, _ := h.coreRepository.GetLookupCodes(query, new(pagination.Paginatevalue))
		result, _ = h.service.UpdateChannel(data, channels, lookupcode[0].LookupCode)
	}
	result, _ = h.service.UpdateChannel(data, channels, "")
	return res.RespSuccess(c, "Channel Updated successfully", map[string]interface{}{"updated_id": result.UpdatedByID})

}

// DeleteChannel godoc
// @Summary      Do Delete Channel
// @Description  Deletes a Channel
// @Tags         Channel
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router       /api/v1/channels/:id/delete [delete]
func (h *handler) DeleteChannel(c echo.Context) error {

	var id = c.Param("id")
	s := c.Get("TokenUserID").(string)
	var data ChannelDTO
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	channels, err := h.service.GetChannelService(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	external_id := channels.ExternalId
	external_value := *external_id
	user_id := helpers.ConvertStringToUint(s)
	if external_value != uint(0) {
		query = map[string]interface{}{"id": channels.ChannelTypeId}
		lookupcode, err := h.coreRepository.GetLookupCodes(query, new(pagination.Paginatevalue))
		data.DeletedByID = helpers.ConvertStringToUint(s)
		_, err = h.service.DeleteChannel(channels, *user_id, lookupcode[0].LookupCode)
		if err != nil {
			return res.RespErr(c, err)
		}
	}
	_, err = h.service.DeleteChannel(channels, *user_id, "")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Channel Deleted successfully", map[string]interface{}{"id": ID})

}

// GetAllChannels godoc
// @Summary      Do fetch Channel
// @Description  Fetch all Channels
// @Tags         Channel
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router       /api/v1/channels [get]
func (h *handler) GetAllChannels(c echo.Context) error {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	query := make(map[string]interface{}, 0)
	channel_types, err := h.coreRepository.GetLookupTypes(map[string]interface{}{
		"lookup_type": "CHANNEL_TYPE",
	}, new(pagination.Paginatevalue))
	for _, channel := range channel_types[0].Lookupcodes {
		if channel.LookupCode == "MARKETPLACE" {
			query["marketplace"] = channel.ID
		}
		if channel.LookupCode == "WEBSTORE" {
			query["webstore"] = channel.ID
		}
	}

	result, err := h.service.FindAllChannels(query, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data retrieved successfully", result, p)
}
