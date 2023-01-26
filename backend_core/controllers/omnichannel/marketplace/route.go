package marketplace

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

	// REGISTER MARKETPLACE
	g.GET("", h.FetchAllMarketplace, cmiddleware.Authorization)
	g.GET("/dropdown", h.FetchAllMarketplaceDropDown, cmiddleware.Authorization)
	g.GET("/:id", h.FindMarketplaceByID, cmiddleware.Authorization)
	g.POST("/register", h.RegisterMarket, cmiddleware.Authorization, MarketPlaceCreateValidate)
	g.POST("/:id/edit", h.UpdateMarket, cmiddleware.Authorization, MarketPlaceUpdateValidate)
	g.DELETE("/:id/delete_marketplace", h.DeleteMarket, cmiddleware.Authorization)

	// MARKETPLACE DETAILS
	g.POST("/create", h.SaveMarketDetails, cmiddleware.Authorization, MarketPlaceDetailsCreateValidate)
	g.POST("/:id/update", h.UpdateMarketDetails, cmiddleware.Authorization, MarketPlaceDetailsUpdateValidate)
	g.POST("/update", h.UpdateMarketDetailsByQuery, cmiddleware.Authorization, MarketPlaceDetailsUpdateValidate)
	g.GET("/view/:id", h.FindMarketDetails, cmiddleware.Authorization)
	g.GET("/list", h.ListMarketDetails, cmiddleware.Authorization)
	g.GET("/list/dropdown", h.ListMarketDetailsDropDown, cmiddleware.Authorization)
	g.DELETE("/:id/delete", h.DeleteMarketDetails, cmiddleware.Authorization)
	g.GET("/:marketplace_code/get_auth_keys", h.GetAuthKeys, cmiddleware.Authorization)

	//AVAILABLE MARKETPLACE
	g.GET("/available", h.AvailableMarketPlaces, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteMarketPlace, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteMarketPlace, cmiddleware.Authorization)

	// MARKETPLACE
	g.POST("/create_marketplace", h.SaveMarketplace, cmiddleware.Authorization)
	g.POST("/:id/update_marketplace", h.UpdateMarketplace, cmiddleware.Authorization)
	g.GET("/:id/view_marketplace", h.FindMarketplace, cmiddleware.Authorization)

}
