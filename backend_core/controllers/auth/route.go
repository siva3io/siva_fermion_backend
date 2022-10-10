package auth

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

	//--------------basic Auth------------------
	g.POST("/register", h.Register, RegisterValidator)
	g.POST("/login", h.Login, LoginValidator)

	//-------------OTP---------------
	g.POST("/user_login", h.UserLogin)
	g.POST("/verify_otp", h.VerifyOTP)
	g.POST("/:id/update", h.AssignTemplate)
	g.POST("/update_profile", h.UpdateProfile, UpdateProfilerValidator)
	g.GET("/:id/get_user", h.GetUserById)
	g.GET("/me", h.GetUserProfile, cmiddleware.Authorization)
}