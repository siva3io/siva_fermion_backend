package omnichannel_fields

import (
	cmiddleware "fermion/backend_core/middleware"

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

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.ListOmnichannelFields, cmiddleware.Authorization)
	g.GET("/:id", h.ViewOmnichannelFields, cmiddleware.Authorization)
	g.POST("/create", h.CreateOmnichannelFields, cmiddleware.Authorization, OmnichannelFieldCreateValidate)
	g.POST("/:id/update", h.UpdateOmnichannelFields, cmiddleware.Authorization, OmnichannelFieldUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteOmnichannelFields, cmiddleware.Authorization)

	g.GET("/app_fields", h.ViewAppFields, cmiddleware.Authorization, ViewAppFieldsQueryValidate)
	g.GET("/app_data", h.ViewAppFieldsData, cmiddleware.Authorization, ViewAppFieldsDataQueryValidate)
	g.POST("/upsert_data", h.UpsertOmnichannelFieldsData, cmiddleware.Authorization, OmnichannelFieldsDataUpsertValidate)

	g.GET("/:app_code/sync_settings", h.GetAppSyncSettings, cmiddleware.Authorization)
	g.POST("/:app_code/upsert_sync_settings", h.UpsertAppSyncSettings, cmiddleware.Authorization, OmnichannelAppSyncSettingsUpsertValidate)

	g.GET("/app_data/filter", h.GetAppDataFilter, cmiddleware.Authorization, GetAppDataFilterQueryValidate)
}
