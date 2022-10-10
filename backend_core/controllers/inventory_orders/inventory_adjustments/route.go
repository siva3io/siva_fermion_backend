package inventory_adjustments

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
	g.GET("", h.GetAllInvAdj, cmiddleware.Authorization)
	g.GET("/:id", h.GetInvAdj, cmiddleware.Authorization)
	g.POST("/create", h.CreateInvAdj, cmiddleware.Authorization, InventoryAdjustmentsCreateValidate)
	g.POST("/bulk_create", h.BulkCreateInvAdj, cmiddleware.Authorization)
	g.PUT("/:id/edit", h.UpdateInvAdj, cmiddleware.Authorization, InventoryAdjustmentsUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteInvAdj, cmiddleware.Authorization)
	g.DELETE("/:id/delete_products", h.DeleteInvAdjLines, cmiddleware.Authorization)

	g.GET("/:id/sendemail", h.SendMailInvAdj, cmiddleware.Authorization)
	g.GET("/:id/printpdf", h.DownloadPdfInvAdj, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteInvAdj, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteInvAdj, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteInvAdjView, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetAllInvAdjDropDown, cmiddleware.Authorization)
}
