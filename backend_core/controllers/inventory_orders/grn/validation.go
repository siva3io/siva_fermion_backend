package grn

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
func (grn GRNRequest) Validate() error {
	return validation.ValidateStruct(&grn,
		validation.Field(
			&grn.GRNNumber,
			validation.When(!grn.AutoGenerateGrnNumber, validation.Required),
		),
		validation.Field(&grn.WarehouseID, validation.Required),
		validation.Field(&grn.GRNOrderLines, validation.Required),
	)
}

func (grn_lines GRNOrderLines) Validate() error {
	return validation.ValidateStruct(&grn_lines,
		validation.Field(&grn_lines.ProductID, validation.Required),
		// validation.Field(&grn_lines.ProductTemplateId, validation.Required),
		validation.Field(&grn_lines.UOMId, validation.Required),
		validation.Field(&grn_lines.OrderedUnits, validation.Required),
		validation.Field(&grn_lines.ReceivedUnits, validation.Required),
		validation.Field(&grn_lines.PendingUnits, validation.Required),
		validation.Field(&grn_lines.QualityCheck, validation.Required),
		validation.Field(&grn_lines.PutawayStatus, validation.Required),
		validation.Field(&grn_lines.PendingUnits, validation.Required),
	)
}

func GrnCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(GRNRequest)
	return func(c echo.Context) error {
		er := c.Bind(data)
		if er != nil {
			validation_err := helpers.BindErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.Validate(data)
		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		c.Set("grn", data)
		return next(c)
	}
}

func GrnUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(GRNRequest)
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

		c.Set("grn", data)
		return next(c)
	}
}
