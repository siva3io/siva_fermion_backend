package grn

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
	g.POST("/create", h.CreateGRNEvent, cmiddleware.Authorization, GrnCreateValidate)
	g.POST("/bulk_create", h.BulkCreateGRN, cmiddleware.Authorization)
	g.GET("/:id", h.GetGRN, cmiddleware.Authorization)
	g.GET("", h.GetAllGRN, cmiddleware.Authorization)
	g.POST("/:id/update", h.UpdateGRNEvent, cmiddleware.Authorization, GrnUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteGRN, cmiddleware.Authorization)
	g.DELETE("/order_line/:id/delete", h.DeleteGRNOrderLine, cmiddleware.Authorization)
	g.POST("/:id/send_email", h.SendEmailGRN, cmiddleware.Authorization)
	g.POST("/:id/download_pdf", h.DownloadGRNPDF, cmiddleware.Authorization)
	g.POST("/:id/generate_pdf", h.GenerateGRNPDF, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteGrn, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteGrn, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteGrnView, cmiddleware.Authorization)
	g.GET("/:id/filter_module/:tab", h.GetGrnTab, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetAllGRNDropDown, cmiddleware.Authorization)
}
