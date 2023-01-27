package asn

import (
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

func (c AsnRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AsnNumber, validation.When(!c.AutoCreateAsnNumber, validation.Required)),
		validation.Field(&c.ReferenceNumber, validation.When(!c.AutoGenerateReferenceNumber, validation.Required)),
		validation.Field(&c.WarehouseID, validation.Required),
		validation.Field(&c.ShippingModeID, validation.Required),
		validation.Field(&c.LinkPurchaseOrderID, validation.Required),
		validation.Field(&c.SourceDocumentTypeID, validation.Required),
		validation.Field(&c.AsnOrderLines, validation.Required),
		validation.Field(&c.DestinationLocationDetails, validation.Required),
		validation.Field(&c.DispatchLocationDetails, validation.Required),
	)
}

func (d AsnLines) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ProductVariantID, validation.Required),
		validation.Field(&d.PackageLength, validation.Required),
		validation.Field(&d.PackageWidth, validation.Required),
		validation.Field(&d.PackageHeight, validation.Required),
		validation.Field(&d.PackageWeight, validation.Required),
		validation.Field(&d.PackageTypeID, validation.Required),
		validation.Field(&d.NumberOfBoxes, validation.Required),
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

		err := validation.Validate(data)
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
