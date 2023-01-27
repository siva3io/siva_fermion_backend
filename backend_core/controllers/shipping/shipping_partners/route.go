package shipping_partners

import (
	// mddlwr "euni_plt_go_gorm/backend_core/internal/middleware"

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
	g.POST("/register", h.CreateUserShippingPartnerRegistrationEvent, cmiddleware.Authorization, ShippingPartnersCreateValidate)
	g.GET("", h.GetAllUserShippingPartnerRegistration, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetAllUserShippingPartnerRegistrationDropDown, cmiddleware.Authorization)
	g.GET("/registered/:id", h.GetUserShippingPartnerRegistration, cmiddleware.Authorization)
	g.POST("/:id/update", h.UpdateUserShippingPartnerRegistrationEvent, cmiddleware.Authorization, ShippingPartnersUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteUserShippingPartnerRegistration, cmiddleware.Authorization)
	g.POST("/rate_calc", h.ShippingPartnerEstimateCosts, cmiddleware.Authorization)
	g.POST("/:id/favourite", h.FavouriteUserShippingPartnerRegistration, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteUserShippingPartnerRegistration, cmiddleware.Authorization)
	g.POST("/:id/download_pdf", h.DownloadUserShippingPartnerRegistrationPDF, cmiddleware.Authorization)
	g.POST("/:id/upload_document", h.UploadUserShippingPartnerRegistrationDocument, cmiddleware.Authorization)

	g.GET("/list", h.GetAllShippingPartner, cmiddleware.Authorization)
	g.GET("/:partner_name/get_auth_options", h.GetShippingPartnerAuth, cmiddleware.Authorization)
	g.GET("/:id", h.GetShippingPartner, cmiddleware.Authorization)
	g.POST("/:partner_name/update_partner", h.UpdateShippingPartnerByName, cmiddleware.Authorization)
}
