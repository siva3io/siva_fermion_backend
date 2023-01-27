package shipping_orders_rto

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
func (r RTORequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ShippingOrderId, validation.Required),
		validation.Field(&r.BookingDate, validation.Required),
		validation.Field(&r.ReturnDate, validation.Required),
		validation.Field(&r.RTOCost, validation.Required),
		validation.Field(&r.EstimatedDeliveryDate, validation.Required),
		validation.Field(&r.Quantity, validation.Required),
		validation.Field(&r.ShippingPartnerId, validation.Required),
		validation.Field(&r.RTOStatusId, validation.Required),
	)
}

func ShippingOrdersRTOCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(RTORequest)
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

		c.Set("shipping_orders_rto", data)
		return next(c)
	}
}

func ShippingOrdersRTOUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(RTORequest)
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

		c.Set("shipping_orders_rto", data)
		return next(c)
	}
}
