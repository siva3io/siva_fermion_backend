package uom

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

	g.POST("/create", h.CreateUomEvent, cmiddleware.Authorization, UomCreateValidate)
	g.POST("/:id/update", h.UpdateUomEvent, cmiddleware.Authorization, UomUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteUom, cmiddleware.Authorization)
	g.GET("/:id", h.GetUom, cmiddleware.Authorization)
	g.GET("", h.GetUomList, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetUomListDropdown, cmiddleware.Authorization)

	g.POST("/class/create", h.CreateUomClassEvent, cmiddleware.Authorization, UomClassCreateValidate)
	g.POST("/class/:id/update", h.UpdateUomClassEvent, cmiddleware.Authorization, UomClassUpdateValidate)
	g.DELETE("/class/:id/delete", h.DeleteUomClass, cmiddleware.Authorization)
	g.GET("/class/:id", h.GetUomClass, cmiddleware.Authorization)
	g.GET("/class", h.GetUomClassList, cmiddleware.Authorization)
	g.GET("/class/dropdown", h.GetUomClassListDropdown, cmiddleware.Authorization)
}
