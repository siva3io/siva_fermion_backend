package sales_returns

import (
	"fermion/backend_core/internal/model/returns"
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
func (a SalesReturnsDTO) Validate() error {
	return validation.ValidateStruct(&a,
		//validation.Field(&a.SalesReturnNumber, validation.Required),
		validation.Field(&a.SrDate, validation.Required),
		validation.Field(&a.CustomerPickupAddress, validation.Required),
		validation.Field(&a.CustomerBillingAddress, validation.Required),
		validation.Field(&a.CurrencyId, validation.Required),
		validation.Field(&a.SalesReturnLines, validation.Required),
	)
}
func (a SalesReturnLines) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ProductId, validation.Required),
		validation.Field(&a.UomId, validation.Required),
		validation.Field(&a.QuantitySold, validation.Required),
		validation.Field(&a.QuantityReturned, validation.Required),
		validation.Field(&a.ReturnTypeId, validation.Required),
		validation.Field(&a.Rate, validation.Required),
		validation.Field(&a.Amount, validation.Required),
	)
}
func SalesReturnsCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(SalesReturnsDTO)
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

		c.Set("sales_returns", data)
		return next(c)
	}
}

func SalesReturnsUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(returns.SalesReturns)
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

		c.Set("sales_returns", data)
		return next(c)
	}
}
