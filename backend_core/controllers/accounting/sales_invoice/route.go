package sales_invoice

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
	g.GET("", h.GetAllSalesInvoice, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetAllSalesInvoiceDropDown, cmiddleware.Authorization)
	g.GET("/:id", h.GetSalesInvoice, cmiddleware.Authorization)
	g.POST("/create", h.CreateSalesInvoice, cmiddleware.Authorization, SalesInvoiceCreateValidate)
	g.POST("/bulk_create", h.BulkCreateSalesInvoice, cmiddleware.Authorization)
	g.PUT("/:id/edit", h.UpdateSalesInvoice, cmiddleware.Authorization, SalesInvoiceUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteSalesInvoice, cmiddleware.Authorization)
	g.DELETE("/:id/delete_products", h.DeleteSalesInvoiceLines, cmiddleware.Authorization)

	g.GET("/:id/sendemail", h.SendMailSalesInvoice, cmiddleware.Authorization)
	g.GET("/:id/printpdf", h.DownloadPdfSalesInvoice, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteSalesInvoice, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteSalesInvoice, cmiddleware.Authorization)
	g.GET("/:id/filter_module/:tab", h.GetSalesInvoiceTab, cmiddleware.Authorization)

}
