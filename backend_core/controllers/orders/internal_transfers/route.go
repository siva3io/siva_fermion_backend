package internal_transfers

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
	g.GET("", h.ListInternalTransfers, cmiddleware.Authorization)
	g.GET("/dropdown", h.ListInternalTransfersDropDown, cmiddleware.Authorization)
	g.GET("/:id", h.ViewInternalTransfers, cmiddleware.Authorization)
	g.POST("/create", h.CreateInternalTransfers, cmiddleware.Authorization, InternalTransfersCreateValidate)
	g.POST("/:id/update", h.UpdateInternalTransfers, cmiddleware.Authorization, InternalTransfersUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteInternalTransfers, cmiddleware.Authorization)
	g.DELETE("/order_lines/:id/delete", h.DeleteInternalTransfersLines, cmiddleware.Authorization)
	g.GET("/search", h.SearchInternalTransfers, cmiddleware.Authorization)

	g.POST("/:id/downloadPdf", h.DownloadInternalTransfers, cmiddleware.Authorization)

	g.POST("/:id/favourite", h.FavouriteInternalTransfers, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteInternalTransfers, cmiddleware.Authorization)
	g.GET("/:id/filter_module/:tab", h.GetInternalTransfersTab, cmiddleware.Authorization)
}
