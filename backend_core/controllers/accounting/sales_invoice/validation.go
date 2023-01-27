package sales_invoice

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
func (c SalesInvoiceRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.SalesInvoiceNumber,
			validation.When(!c.AutoGenerateInvoiceNumber, validation.Required),
		),
		validation.Field(
			&c.ReferenceNumber,
			validation.When(!c.AutoGenerateReferenceNumber, validation.Required),
		),
		validation.Field(&c.CustomerID, validation.Required),
		validation.Field(&c.CurrencyID, validation.Required),
		validation.Field(&c.SalesInvoiceDate, validation.Required),
		validation.Field(&c.BillingAddress, validation.Required),
		validation.Field(&c.DeliveryAddress, validation.Required),
		validation.Field(&c.SalesInvoiceLines, validation.Required),
		validation.Field(&c.AvailableCustomerCredits, validation.Required),
	)
}

func (d SalesInvoiceLines) Validate() error {
	return validation.ValidateStruct(&d,
		//validation.Field(&d.ProductID, validation.Required),
		validation.Field(&d.ProductVariantID, validation.Required),
		validation.Field(&d.Discount, validation.Required),
		validation.Field(&d.Tax, validation.Required),
		validation.Field(&d.Quantity, validation.Required),
		validation.Field(&d.DiscountTypeID, validation.Required),
		validation.Field(&d.TaxTypeID, validation.Required),
		validation.Field(&d.TotalAmount, validation.Required),
		validation.Field(&d.Price, validation.Required),
	)
}

func SalesInvoiceCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(SalesInvoiceRequest)
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

		c.Set("sales_invoice", data)
		return next(c)
	}
}

func SalesInvoiceUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(SalesInvoiceRequest)
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

		c.Set("sales_invoice", data)
		return next(c)
	}
}
