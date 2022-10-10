package pricing

import (
	"encoding/json"

	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	validation "github.com/go-ozzo/ozzo-validation"
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
func PricingCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(shared_pricing_and_location.Pricing)
	return func(c echo.Context) error {
		var reqPayload map[string]interface{}
		json.NewDecoder(c.Request().Body).Decode(&reqPayload)
		reqBytes, _ := json.Marshal(reqPayload)
		er := json.Unmarshal(reqBytes, data)
		if er != nil {
			validation_err := helpers.JsonMarshalErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		if data.Price_list_id == 1 {
			sales := new(shared_pricing_and_location.SalesPriceList)
			er := json.Unmarshal(reqBytes, sales)
			if er != nil {
				validation_err := helpers.JsonMarshalErrorStructure(er)
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}

			err := validation.ValidateStruct(sales)

			if err != nil {
				validation_err := helpers.ValidationErrorStructure(err)
				if validation_err != nil {
					return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
				}
			}
			c.Set("sales", sales)
		} else if data.Price_list_id == 2 {
			purchase := new(shared_pricing_and_location.PurchasePriceList)
			er := json.Unmarshal(reqBytes, purchase)
			if er != nil {
				validation_err := helpers.JsonMarshalErrorStructure(er)
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}

			err := validation.ValidateStruct(purchase)

			if err != nil {
				validation_err := helpers.ValidationErrorStructure(err)
				if validation_err != nil {
					return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
				}
			}
			c.Set("purchase", purchase)
		} else if data.Price_list_id == 3 {
			transfer := new(shared_pricing_and_location.TransferPriceList)
			er := json.Unmarshal(reqBytes, transfer)
			if er != nil {
				validation_err := helpers.JsonMarshalErrorStructure(er)
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}

			err := validation.ValidateStruct(transfer)

			if err != nil {
				validation_err := helpers.ValidationErrorStructure(err)
				if validation_err != nil {
					return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
				}
			}
			c.Set("transfer", transfer)
		}

		c.Set("pricing", data)
		return next(c)
	}
}

func PricingUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(shared_pricing_and_location.Pricing)
	return func(c echo.Context) error {
		var reqPayload map[string]interface{}
		json.NewDecoder(c.Request().Body).Decode(&reqPayload)
		reqBytes, _ := json.Marshal(reqPayload)
		er := json.Unmarshal(reqBytes, data)
		if er != nil {
			validation_err := helpers.JsonMarshalErrorStructure(er)
			return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
		}

		err := validation.ValidateStruct(data)

		if err != nil {
			validation_err := helpers.ValidationErrorStructure(err)
			if validation_err != nil {
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}
		}

		if data.Price_list_id == 1 {
			sales := new(shared_pricing_and_location.SalesPriceList)
			er := json.Unmarshal(reqBytes, sales)
			if er != nil {
				validation_err := helpers.JsonMarshalErrorStructure(er)
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}

			err := validation.ValidateStruct(sales)

			if err != nil {
				validation_err := helpers.ValidationErrorStructure(err)
				if validation_err != nil {
					return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
				}
			}
			c.Set("sales", sales)
		} else if data.Price_list_id == 2 {
			purchase := new(shared_pricing_and_location.PurchasePriceList)
			er := json.Unmarshal(reqBytes, purchase)
			if er != nil {
				validation_err := helpers.JsonMarshalErrorStructure(er)
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}

			err := validation.ValidateStruct(purchase)

			if err != nil {
				validation_err := helpers.ValidationErrorStructure(err)
				if validation_err != nil {
					return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
				}
			}
			c.Set("purchase", purchase)
		} else if data.Price_list_id == 3 {
			transfer := new(shared_pricing_and_location.TransferPriceList)
			er := json.Unmarshal(reqBytes, transfer)
			if er != nil {
				validation_err := helpers.JsonMarshalErrorStructure(er)
				return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
			}

			err := validation.ValidateStruct(transfer)

			if err != nil {
				validation_err := helpers.ValidationErrorStructure(err)
				if validation_err != nil {
					return res.RespValidationErr(c, "Invalid Fields or Parameter Found", validation_err)
				}
			}
			c.Set("transfer", transfer)
		}

		c.Set("pricing", data)
		return next(c)
	}
}
