package shipping_orders

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
	g.POST("/create", h.CreateShippingOrderEvent, cmiddleware.Authorization, ShippingOrdersCreateValidate)
	g.GET("", h.GetAllShippingOrders, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetAllShippingOrdersDropDown, cmiddleware.Authorization)
	g.GET("/:id", h.GetShippingOrder, cmiddleware.Authorization)
	g.POST("/:id/update", h.UpdateShippingOrderEvent, cmiddleware.Authorization, ShippingOrdersUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteShippingOrder, cmiddleware.Authorization)
	g.DELETE("/order_line/:id/delete", h.DeleteShippingOrderLine, cmiddleware.Authorization)
	g.POST("/bulk_create", h.CreateBulkShippingOrder, cmiddleware.Authorization)
	g.POST("/:id/send_email", h.SendEmailShippingOrder, cmiddleware.Authorization)
	g.POST("/:id/download_pdf", h.DownloadShippingOrderPDF, cmiddleware.Authorization)
	g.POST("/:id/generate_pdf", h.GenerateShippingOrderPDF, cmiddleware.Authorization)
	// g.POST("/:id/print_label", h.PrintShippingOrderLabel, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteShippingOrder, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteShippingOrder, cmiddleware.Authorization)
	g.GET("/:id/filter_module/:tab", h.GetShippingOrderTab, cmiddleware.Authorization)

	//================================ADMIN===========================================
	g.GET("/admin/list", h.GetAllShippingOrdersAdmin, cmiddleware.Authorization)
	g.GET("/admin/:id", h.GetShippingOrderAdmin, cmiddleware.Authorization)
}
