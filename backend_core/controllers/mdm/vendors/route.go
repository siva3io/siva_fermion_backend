package vendors

import (
	"github.com/labstack/echo/v4"

	cmiddleware "fermion/backend_core/middleware"
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
	g.POST("/upsert", h.UpsertVendorEvent, cmiddleware.Authorization)
	g.POST("/create", h.CreateVendorEvent, cmiddleware.Authorization, VendorsCreateValidate)
	g.GET("/:id", h.GetVendors, cmiddleware.Authorization)
	g.GET("", h.GetVendorsList, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetVendorsListDropdown, cmiddleware.Authorization)
	g.DELETE("/:id/delete", h.DeleteVendors, cmiddleware.Authorization)
	g.POST("/:id/update", h.UpdateVendorEvent, cmiddleware.Authorization, VendorsUpdateValidate)
	g.POST("/:id/favourite", h.FavouriteVendors, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteVendors, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteVendorsView, cmiddleware.Authorization)

	g.GET("/pricelist", h.GetVendorPriceLists, cmiddleware.Authorization)
	g.GET("/pricelist/dropdown", h.GetVendorPriceListsDropdown, cmiddleware.Authorization)
	g.POST("/pricelist/create", h.CreateVendorPriceList, cmiddleware.Authorization, VendorsPriceListCreateValidate)
	g.POST("/pricelist/:id/update", h.UpdateVendorPriceLists, cmiddleware.Authorization, VendorsPriceListUpdateValidate)
	g.GET("/pricelist/:id", h.GetVendorPriceList, cmiddleware.Authorization)
	g.DELETE("/pricelist/:id/delete", h.DeleteVendorPricelist, cmiddleware.Authorization)
}
