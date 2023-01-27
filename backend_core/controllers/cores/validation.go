package cores

import (
	// "errors"

	"regexp"

	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (data Company) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(
			&data.AadhaarNumber,
			validation.Required,
			validation.Match((regexp.MustCompile("^[0-9]{4}[ -]?[0-9]{4}[ -]?[0-9]{4}$"))),
		),
		validation.Field(
			&data.GSTIN,
			validation.Required,
			validation.Match((regexp.MustCompile("^[0-9]{2}[A-Z]{5}[0-9]{4}"+"[A-Z]{1}[1-9A-Z]{1}"+"Z[0-9A-Z]{1}$"))),
		),

		validation.Field(
			&data.PANNumber,
			validation.Required,
			validation.Match(regexp.MustCompile("^[A-Z]{5}[0-9]{4}[A-Z]{1}$")),
		),
	)
}

func CompanyUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(Company)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data,
			validation.Field(
				&data.Email,
				is.Email,
			),
			validation.Field(
				&data.KYCDocuments,
				validation.Required,
			),
		)
		if err != nil {
			validation_err := helpers.ValidationFieldStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}
		c.Set("company", data)
		return next(c)
	}
}

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
