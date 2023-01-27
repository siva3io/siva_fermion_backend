package sales_returns

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
	g.GET("", h.ListSalesReturns, cmiddleware.Authorization)
	g.GET("/dropdown", h.ListSalesReturnsDropDown, cmiddleware.Authorization)
	g.GET("/:id", h.ViewSalesReturns, cmiddleware.Authorization)
	g.POST("/create", h.CreateSalesReturnsEvent, cmiddleware.Authorization, SalesReturnsCreateValidate)
	g.POST("/:id/update", h.UpdateSalesReturnsEvent, cmiddleware.Authorization, SalesReturnsUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteSalesReturns, cmiddleware.Authorization)
	g.DELETE("/return_lines/:id/delete", h.DeleteSalesReturnLines, cmiddleware.Authorization)
	g.GET("/search", h.SearchSalesReturns, cmiddleware.Authorization)
	g.GET("/track", h.TrackSalesReturn)

	g.POST("/:id/downloadPdf", h.DownloadSalesReturns, cmiddleware.Authorization)
	g.POST("/:id/sendEmail", h.EmailSalesReturns, cmiddleware.Authorization)
	g.POST("/:id/generatePdf", h.GenerateSalesReturnsPDF, cmiddleware.Authorization)

	g.POST("/:id/favourite", h.FavouriteSalesReturns, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteSalesReturns, cmiddleware.Authorization)

	g.GET("/:product_id/history", h.GetSalesReturnsHistory, cmiddleware.Authorization)
	g.GET("/:id/filter_module/:tab", h.GetSalesReturnTab, cmiddleware.Authorization)
}
