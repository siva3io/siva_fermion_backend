package cycle_count

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
	g.GET("/get_count", h.GetCount, cmiddleware.Authorization)
	g.GET("/bin_list/:variant_id", h.GetBinsList, cmiddleware.Authorization)
	g.POST("/product_list", h.GetProductList, cmiddleware.Authorization)

	g.GET("", h.GetAllCycleCount, cmiddleware.Authorization)
	g.GET("/:id", h.GetCycleCount, cmiddleware.Authorization)
	g.POST("/create", h.CreateCycleCountEvent, cmiddleware.Authorization, CycleCountCreateValidate)
	g.POST("/bulk_create", h.BulkCreateCycleCount, cmiddleware.Authorization)
	g.PUT("/:id/edit", h.UpdateCycleCountEvent, cmiddleware.Authorization, CycleCountUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteCycleCount, cmiddleware.Authorization)
	g.DELETE("/:id/delete_products", h.DeleteCycleCountLines, cmiddleware.Authorization)

	g.GET("/:id/sendemail", h.SendMailCycleCount, cmiddleware.Authorization)
	g.GET("/:id/printpdf", h.DownloadPdfCycleCount, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteCycleCount, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteCycleCount, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteCycleCountView, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetAllCycleCountDropDown, cmiddleware.Authorization)
}
