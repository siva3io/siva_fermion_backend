package cores

import (
	// "errors"

	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
func (data CreateLookupTypeDTO) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.LookupType, validation.Required),
		validation.Field(&data.DisplayName, validation.Required),
	)
}

func LookupTypeCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateLookupTypeDTO)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		if err := data.Validate(); err != nil {
			return res.RespError(c, res.BuildError(res.ErrValidation, err))
		}

		c.Set("lookup_types", data)
		return next(c)
	}
}

func (data CreateUpdateLookupCodesDTO) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.LookupCodes, validation.Required),
	)
}

func (data LookupCode) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.DisplayName, validation.Required),
		validation.Field(&data.LookupCode, validation.Required),
		validation.Field(&data.LookupTypeId, validation.Required),
	)
}

func LookupCodeCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateUpdateLookupCodesDTO)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		if err := data.Validate(); err != nil {
			return res.RespError(c, res.BuildError(res.ErrValidation, err))
		}

		c.Set("lookup_codes", data)
		return next(c)
	}
}

func (data UpdateLookupTypeDTO) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.DisplayName, validation.Required),
	)
}

func LookupTypeUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {
	var data = new(UpdateLookupTypeDTO)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := data.Validate()

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		c.Set("lookup_types", data)
		return next(c)
	}
}

func (data UpdateLookupCodeDTO) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.DisplayName, validation.Required),
	)
}

func LookupCodeUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {
	var data = new(UpdateLookupCodeDTO)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := data.Validate()

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		c.Set("lookup_codes", data)
		return next(c)
	}
}
