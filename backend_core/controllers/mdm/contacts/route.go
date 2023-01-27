package contacts

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
	g.POST("/upsert", h.UpsertContactEvent, cmiddleware.Authorization)
	g.POST("/create", h.CreateContactEvent, cmiddleware.Authorization, ContactsCreateValidate)
	g.POST("/:id/update", h.UpdateContactEvent, cmiddleware.Authorization, ContactsUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteContacts, cmiddleware.Authorization)
	g.GET("/:id", h.ContactView, cmiddleware.Authorization)
	g.GET("/:id/related", h.GetRelatedContacts, cmiddleware.Authorization)
	g.GET("", h.GetContacts, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetContactsDropdown, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteContacts, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteContacts, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteContactsView, cmiddleware.Authorization)
	g.GET("/:id/filter_module/:tab", h.GetContactTab, cmiddleware.Authorization)
}
