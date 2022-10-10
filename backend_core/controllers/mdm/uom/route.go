package uom

import (
	cmiddleware "fermion/backend_core/middleware"

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
func (h *handler) Route(g *echo.Group) {

	g.POST("/create", h.CreateUom, cmiddleware.Authorization, UomCreateValidate)
	g.POST("/:id/update", h.UpdateUom, cmiddleware.Authorization, UomUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteUom, cmiddleware.Authorization)
	g.GET("/:id", h.GetUom, cmiddleware.Authorization)
	g.GET("", h.GetUomList, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetUomListDropdown, cmiddleware.Authorization)
	g.GET("/search", h.SearchUom, cmiddleware.Authorization)

	g.POST("/class/create", h.CreateUomClass, cmiddleware.Authorization, UomClassCreateValidate)
	g.POST("/class/:id/update", h.UpdateUomClass, cmiddleware.Authorization, UomClassUpdateValidate)
	g.DELETE("/class/:id/delete", h.DeleteUomClass, cmiddleware.Authorization)
	g.GET("/class/:id", h.GetUomClass, cmiddleware.Authorization)
	g.GET("/class", h.GetUomClassList, cmiddleware.Authorization)
	g.GET("/class/dropdown", h.GetUomClassListDropdown, cmiddleware.Authorization)
	g.GET("/class/search", h.SearchUomClass, cmiddleware.Authorization)
}
