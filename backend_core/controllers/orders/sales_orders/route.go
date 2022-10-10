package sales_orders

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
	g.GET("", h.ListSalesOrders, cmiddleware.Authorization)
	g.GET("/dropdown", h.ListSalesOrdersDropDown, cmiddleware.Authorization)
	g.GET("/:id", h.ViewSalesOrders, cmiddleware.Authorization)
	g.POST("/create", h.CreateSalesOrders, cmiddleware.Authorization, SalesOrdersCreateValidate)
	g.POST("/:id/update", h.UpdateSalesOrders, cmiddleware.Authorization, SalesOrdersUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteSalesOrders, cmiddleware.Authorization)
	g.DELETE("/order_lines/:id/delete", h.DeleteSalesOrderLines, cmiddleware.Authorization)
	g.GET("/search", h.SearchSalesOrders, cmiddleware.Authorization)

	g.POST("/:id/downloadPdf", h.DownloadSalesOrders, cmiddleware.Authorization)
	g.POST("/:id/sendEmail", h.EmailSalesOrders, cmiddleware.Authorization)
	g.POST("/:id/generatePdf", h.GenerateSalesOrdersPDF, cmiddleware.Authorization)

	g.POST("/:id/favourite", h.FavouriteSalesOrders, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteSalesOrders, cmiddleware.Authorization)

	//-------------Channel API's----------------------------------------------------
	g.POST("/channels/upsert", h.ChannelSalesOrderUpsert, cmiddleware.Authorization)

	g.GET("/:product_id/history", h.GetSalesHistory, cmiddleware.Authorization)

	g.GET("/:id/filter_module/:tab", h.GetSalesOrderTab, cmiddleware.Authorization)
}
