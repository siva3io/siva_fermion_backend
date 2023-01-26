package auth

import (
	// "fermion/backend_core/pkg/util/helpers"

	"errors"

	res "fermion/backend_core/pkg/util/response"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

func (r RegisterDTO) Validate() error {

	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, is.Email),
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.ConfirmPassword, validation.Required),
	)
}

func (r LoginDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

func RegisterValidator(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		var data = new(RegisterDTO)

		if err := c.Bind(data); err != nil {
			return res.RespErr(c, err)
		}

		if err := data.Validate(); err != nil {
			return res.RespErr(c, err)
		}

		if data.Password != data.ConfirmPassword {
			return res.RespErr(c, errors.New("password mismatch"))
		}

		c.Set("register_data", data)

		return next(c)
	}
}

func LoginValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data = new(LoginDTO)

		if err := c.Bind(data); err != nil {
			return res.RespErr(c, err)
		}

		if err := data.Validate(); err != nil {
			return res.RespErr(c, err)
		}

		c.Set("login_data", data)

		return next(c)
	}
}

func UpdateProfilerValidator(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		var data = new(UpdateProfileDTO)

		if err := c.Bind(data); err != nil {
			return res.RespErr(c, err)
		}
		c.Set("profile_data", data)

		return next(c)
	}
}
