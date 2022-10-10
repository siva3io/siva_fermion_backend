package asn

import (
	"errors"

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

func (c AsnRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.AsnNumber,
			validation.When(!c.AutoCreateAsnNumber, validation.Required),
		),
		validation.Field(
			&c.ReferenceNumber,
			validation.When(!c.AutoGenerateReferenceNumber, validation.Required),
		),
		validation.Field(
			&c.WarehouseID,
			validation.Required,
		),
		validation.Field(
			&c.ShippingModeID,
			validation.Required,
		),
		validation.Field(
			&c.LinkPurchaseOrderID,
			validation.Required,
		),
		validation.Field(
			&c.SourceDocumentTypeID,
			validation.Required,
		),
		validation.Field(
			&c.AsnOrderLines,
			validation.Required,
		),
		validation.Field(
			&c.DestinationLocationDetails,
			validation.Required,
		),
		validation.Field(
			&c.DispatchLocationDetails,
			validation.Required,
		),
	)
}

func (d AsnLines) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(
			&d.ProductID,
			validation.Required,
		),
		validation.Field(
			&d.ProductVariantID,
			validation.Required,
		),
		validation.Field(
			&d.PackageTypeID,
			validation.Required,
		),
	)
}

func AsnCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(AsnRequest)
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
				return res.RespError(c, res.BuildError(res.ErrValidation, errors.New("invalid payload")))
			}
		}

		c.Set("asn", data)
		return next(c)
	}
}

func AsnUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(AsnRequest)
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

		c.Set("asn", data)
		return next(c)
	}
}
