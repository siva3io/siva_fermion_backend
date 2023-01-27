package shipping_orders_ndr

import (
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
func (rd NDRRequest) Validate() error {
	return validation.ValidateStruct(&rd) // validation.Field(&rd.ReceiverName, validation.Required),
	// validation.Field(&rd.Email, validation.Required),
	// validation.Field(&rd.MobileNumber, validation.Required),
	// validation.Field(&rd.AddressLine1, validation.Required),
	// validation.Field(&rd.Zipcode, validation.Required),
	// validation.Field(&rd.City, validation.Required),
	// validation.Field(&rd.State, validation.Required),
	// validation.Field(&rd.Country, validation.Required),
	// validation.Field(&rd.PickupNickname, validation.Required),

}

func ShippingOrdersNDRCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(NDRRequest)
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

		c.Set("shipping_orders_ndr", data)
		return next(c)
	}
}

func ShippingOrdersUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(NDRRequest)
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

		c.Set("shipping_orders_ndr", data)
		return next(c)
	}
}
