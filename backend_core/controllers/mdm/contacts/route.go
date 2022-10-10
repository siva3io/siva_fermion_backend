package contacts

import (
	"github.com/labstack/echo/v4"

	cmiddleware "fermion/backend_core/middleware"
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
	g.POST("/create", h.CreateContact, cmiddleware.Authorization, ContactsCreateValidate)
	g.POST("/:id/update", h.UpdateContact, cmiddleware.Authorization, ContactsUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteContacts, cmiddleware.Authorization)
	g.GET("/:id", h.ContactView, cmiddleware.Authorization)
	g.GET("/:id/related", h.GetRelatedContacts, cmiddleware.Authorization)
	g.GET("", h.GetContacts, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetContactsDropdown, cmiddleware.Authorization)
	g.GET("/search", h.SearchContacts, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteContacts, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteContacts, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteContactsView, cmiddleware.Authorization)
}
