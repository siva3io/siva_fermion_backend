package ipaas_core

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
	g.POST("/shipping/rate_calculator", h.RateCalculator, cmiddleware.Authorization)
	g.GET("/:module/:task", h.GenericHandler, cmiddleware.Authorization)

	g.GET("/integrations/:collection/:id", h.GetFeature, cmiddleware.Authorization)
	g.POST("/integrations/:collection/create", h.CreateFeature, cmiddleware.Authorization)
	g.POST("/integrations/:collection/:id/update", h.UpdateFeature, cmiddleware.Authorization)
	g.DELETE("/integrations/:collection/:id/delete", h.DeleteFeature, cmiddleware.Authorization)
	g.POST("/boson_convertor", h.BosonConvertor, cmiddleware.Authorization)
	g.GET("/integrations/app_features/:app_code", h.GetAppFeatures, cmiddleware.Authorization)
}
