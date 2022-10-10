package scrap_orders

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
func (a ScrapOrders) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Scrap_order_number,
			validation.When(!a.AutoCreateScrapNumber, validation.Required),
		),
		validation.Field(&a.Schedule_scrap_date, validation.Required),
		validation.Field(&a.Scrap_reason_id, validation.Required),
	)
}

func (o ScrapOrderLines) Validate() error {
	return validation.ValidateStruct(&o,

		validation.Field(&o.Product_id, validation.Required),
	)
}
func ScrapOrdersCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ScrapOrders)
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

		c.Set("scrap_orders", data)
		return next(c)
	}
}

func ScrapOrdersUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ScrapOrders)
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

		c.Set("scrap_orders", data)
		return next(c)
	}
}
