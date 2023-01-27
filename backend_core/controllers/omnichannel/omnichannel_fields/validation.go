package omnichannel_fields

import (
	"errors"

	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	validation "github.com/go-ozzo/ozzo-validation"
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
func OmnichannelFieldCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(OmnichannelFieldRequestDto)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		c.Set("omnichannel_fields", data)
		return next(c)
	}
}

func OmnichannelFieldUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(OmnichannelFieldRequestDto)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		c.Set("omnichannel_fields", data)
		return next(c)
	}
}

func ViewAppFieldsQueryValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ViewAppFieldsQueryDto)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		err = validation.Validate(data)
		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespError(c, res.BuildError(res.ErrValidation, errors.New("invalid payload")))
			}
		}

		c.Set("app_field_query", data)
		return next(c)
	}
}

func (dto ViewAppFieldsQueryDto) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(
			&dto.AppId,
			validation.Required,
		),
		validation.Field(
			&dto.ChannelFunctionId,
			validation.Required,
		),
		validation.Field(
			&dto.ChannelTypeId,
			validation.Required,
		),
	)
}

func ViewAppFieldsDataQueryValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ViewAppFieldsDataQueryDto)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		err = validation.Validate(data)
		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespError(c, res.BuildError(res.ErrValidation, errors.New("invalid payload")))
			}
		}

		c.Set("app_field_data_query", data)
		return next(c)
	}
}

func (dto ViewAppFieldsDataQueryDto) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(
			&dto.FieldAppId,
			validation.Required,
		),
		validation.Field(
			&dto.ChannelFunctionId,
			validation.Required,
		),
	)
}

func OmnichannelFieldsDataUpsertValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(OmnichannelFieldDataRequestDto)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		c.Set("omnichannel_fields_data", data)
		return next(c)
	}
}

func OmnichannelAppSyncSettingsUpsertValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new([]OmnichannelSyncSettingsRequestDto)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		for _, obj := range *data {
			err := validation.ValidateStruct(&obj)

			if err != nil {
				validation_err := helpers.ValidationErrorStructure(err)
				if validation_err != nil {
					return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
				}
			}
		}

		c.Set("omnichannel_sync_settings", data)
		return next(c)
	}
}

func GetAppDataFilterQueryValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(GetAppDataFilterQueryDto)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		err = validation.Validate(data)
		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespError(c, res.BuildError(res.ErrValidation, errors.New("invalid payload")))
			}
		}

		c.Set("get_app_data_filter_query", data)
		return next(c)
	}
}

func (dto GetAppDataFilterQueryDto) Validate() error {
	return validation.ValidateStruct(&dto,
		validation.Field(
			&dto.AppId,
			validation.Required,
		),
		validation.Field(
			&dto.Fields,
			validation.Required,
		),
	)
}
