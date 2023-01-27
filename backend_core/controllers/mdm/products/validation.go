package products

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
//--------------------------------Product Brand -------------------------------------------------------------
func ProductBrandCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(BrandRequestAndResponseDTO)
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

		c.Set("product_brand", data)
		return next(c)
	}
}
func ProductBrandUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(BrandRequestAndResponseDTO)
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

		c.Set("product_brand", data)
		return next(c)
	}
}

// --------------------------------Product Category -------------------------------------------------------------
func ProductCategoryCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CategoryAndSubcategoryRequestAndResponseDTO)
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

		c.Set("product_category", data)
		return next(c)
	}
}
func ProductCategoryUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CategoryAndSubcategoryRequestAndResponseDTO)
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

		c.Set("product_category", data)
		return next(c)
	}
}

// --------------------------------Product Base Attribute -------------------------------------------------------------
func ProductBaseAttributesCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductBaseAttributesRequestAndResponseDTO)
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

		c.Set("product_base_attributes", data)
		return next(c)
	}
}
func ProductBaseAttributesUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductBaseAttributesRequestAndResponseDTO)
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

		c.Set("product_base_attributes", data)
		return next(c)
	}
}

// --------------------------------Product Base Attribute Value -------------------------------------------------------------
func ProductBaseAttributesValuesCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductBaseAttributesValuesRequestAndResponseDTO)
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

		c.Set("product_base_attributes_values", data)
		return next(c)
	}
}
func ProductBaseAttributesValuesUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductBaseAttributesValuesRequestAndResponseDTO)
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

		c.Set("product_base_attributes_values", data)
		return next(c)
	}
}

// --------------------------------Product Selected Base Attribute -------------------------------------------------------------
func ProductSelectedAttributesCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductSelectedAttributesRequestAndREsponseDTO)
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

		c.Set("product_selected_attributes", data)
		return next(c)
	}
}
func ProductSelectedAttributesUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductSelectedAttributesRequestAndREsponseDTO)
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

		c.Set("product_selected_attributes", data)
		return next(c)
	}
}

// --------------------------------Product Selected Base Attribute Value-------------------------------------------------------------
func ProductSelectedAttributesValuesCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductSelectedAttributesValuesRequestAndResponseeDTO)
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

		c.Set("product_selected_attributes_values", data)
		return next(c)
	}
}
func ProductSelectedAttributesValuesUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(ProductSelectedAttributesValuesRequestAndResponseeDTO)
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

		c.Set("product_selected_attributes_values", data)
		return next(c)
	}
}

// --------------------------------Product Bundles-----------------------------------------------------------------------------------
func ProductBundleCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateBundlePayload)
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

		c.Set("product_bundle", data)
		return next(c)
	}
}
func ProductBundleUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateBundlePayload)
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

		c.Set("product_bundle", data)
		return next(c)
	}
}

// --------------------------------Product Template-----------------------------------------------------------------------------------
func (c CreateProductTemplatePayload) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.SKUCode,
			validation.Required,
		),
		validation.Field(
			&c.Name,
			validation.Required,
		),
		// validation.Field(
		// 	&c. ProductProcurementTreatmentIds,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.PrimaryCategoryID,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.SecondaryCategoryID,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.Description,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.Location.Pincode,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.Location.City,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.Location.StateId,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.Location.AddressLine1,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.Location.AddressLine2,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.Location.AddressLine3,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.FoodItemDetails.IngredientsInfo,
		// 	validation.Required,
		// ),
		// validation.Field(
		// 	&c.FoodItemDetails.NutritionalInfo,
		// 	validation.Required,
		// ),
	)
}

func ProductTemplateCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateProductTemplatePayload)
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

		c.Set("product_template", data)
		return next(c)
	}
}
func ProductTemplateUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateProductTemplatePayload)
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

		c.Set("product_template", data)
		return next(c)
	}
}

// --------------------------------Product Variant-----------------------------------------------------------------------------------
func ProductVariantCreateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateProductVariantDTO)
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

		c.Set("product_variant", data)
		return next(c)
	}
}
func ProductVariantUpdateValidate(next echo.HandlerFunc) echo.HandlerFunc {

	var data = new(CreateProductVariantDTO)
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

		c.Set("product_variant", data)
		return next(c)
	}
}
