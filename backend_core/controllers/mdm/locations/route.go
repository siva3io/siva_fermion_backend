package locations

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
	g.POST("/create", h.CreateLocationEvent, cmiddleware.Authorization, LocationsCreateValidate)
	g.POST("/:id/update", h.UpdateLocationEvent, cmiddleware.Authorization, LocationsUpdateValidate)
	g.GET("", h.GetLocationList, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetLocationListDropdown, cmiddleware.Authorization)
	g.GET("/:id", h.GetLocation, cmiddleware.Authorization)
	g.DELETE("/:id/delete", h.DeleteLocation, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteLocations, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteLocations, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteLocationsView, cmiddleware.Authorization)
	//storage location
	g.POST("/storage_location/create", h.CreateStorageLocation, cmiddleware.Authorization)
	g.GET("/storage_location", h.GetStorageLocationList, cmiddleware.Authorization)
	g.GET("/storage_location/:id", h.GetStorageLocation, cmiddleware.Authorization)

	//storage quuantity
	g.GET("/storage_quantity", h.GetStorageQuantityList, cmiddleware.Authorization)
	// g.GET("/:warehouse_code/get_auth_keys", h.GetAuthKeys, cmiddleware.Authorization)
}
