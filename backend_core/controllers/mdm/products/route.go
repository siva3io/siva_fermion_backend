package products

import (
	cmiddleware "fermion/backend_core/middleware"

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
func (h *handler) Route(g *echo.Group) {

	//---------------------Product Brands----------------------------------------------------------------------------------------
	g.POST("/brand/create", h.CreateProductBrand, cmiddleware.Authorization, ProductBrandCreateValidate)
	g.POST("/brand/:id/update", h.UpdateProductBrand, cmiddleware.Authorization, ProductBrandUpdateValidate)
	g.DELETE("/brand/:id/delete", h.DeleteProductBrand, cmiddleware.Authorization)
	g.GET("/brand", h.GetAllBrands, cmiddleware.Authorization)
	g.GET("/brand/dropdown", h.GetAllBrandsDropdown, cmiddleware.Authorization)
	g.GET("/brand/search", h.SearchBrand, cmiddleware.Authorization)

	//---------------------Product category---------------------------------------------------------------------------------------------
	g.POST("/category/create", h.CreateProductCategory, cmiddleware.Authorization, ProductCategoryCreateValidate)
	g.POST("/category/:id/update", h.UpdateProductCategory, cmiddleware.Authorization, ProductCategoryUpdateValidate)
	g.DELETE("/category/:id/delete", h.DeleteProductCategory, cmiddleware.Authorization)
	g.GET("/category", h.GetAllProductCategory, cmiddleware.Authorization)
	g.GET("/category/dropdown", h.GetAllProductCategoryDropdown, cmiddleware.Authorization)
	g.GET("/category/subcategory", h.GetAllSubCategory, cmiddleware.Authorization)
	g.GET("/category/subcategory/dropdown", h.GetAllSubCategoryDropdown, cmiddleware.Authorization)
	g.GET("/category/search", h.SearchCategory, cmiddleware.Authorization)
	g.GET("/category/search/dropdown", h.SearchCategoryDropdown, cmiddleware.Authorization)

	//---------------------Product Base Attributes--------------------------------------------------------------------------------------------
	g.POST("/base_attributes/create", h.CreateProductBaseAttributes, cmiddleware.Authorization, ProductBaseAttributesCreateValidate)
	g.POST("/base_attributes/:id/update", h.UpdateProductBaseAttributes, cmiddleware.Authorization, ProductBaseAttributesUpdateValidate)
	g.DELETE("/base_attributes/:id/delete", h.DeleteProductBaseAttributes, cmiddleware.Authorization)
	g.GET("/base_attributes", h.GetAllProductBaseAttributes, cmiddleware.Authorization)
	g.GET("/base_attributes/dropdown", h.GetAllProductBaseAttributesDropdown, cmiddleware.Authorization)

	//---------------------Product Base Attribute Values-----------------------------------------------------------------------------------------
	g.POST("/base_attributes_values/create", h.CreateProductBaseAttributesValues, cmiddleware.Authorization, ProductBaseAttributesValuesCreateValidate)
	g.POST("/base_attributes_values/:id/update", h.UpdateProductBaseAttributesValues, cmiddleware.Authorization, ProductBaseAttributesValuesUpdateValidate)
	g.DELETE("/base_attributes_values/:id/delete", h.DeleteProductBaseAttributesValues, cmiddleware.Authorization)
	g.GET("/base_attributes_values", h.GetAllProductBaseAttributesValues, cmiddleware.Authorization)
	g.GET("/base_attributes_values/dropdown", h.GetAllProductBaseAttributesValuesDropdown, cmiddleware.Authorization)

	//---------------------Product Selected Attributes--------------------------------------------------------------------------------------------
	g.POST("/selected_attributes/create", h.CreateProductSelectedAttributes, cmiddleware.Authorization, ProductSelectedAttributesCreateValidate)
	g.POST("/selected_attributes/:id/update", h.UpdateProductSelectedAttributes, cmiddleware.Authorization, ProductSelectedAttributesUpdateValidate)
	g.DELETE("/selected_attributes/:id/delete", h.DeleteProductSelectedAttributes, cmiddleware.Authorization)
	g.GET("/selected_attributes", h.GetAllProductSelectedAttributes, cmiddleware.Authorization)
	g.GET("/selected_attributes/dropdown", h.GetAllProductSelectedAttributesDropdown, cmiddleware.Authorization)

	//---------------------Product Selected Attributes Values--------------------------------------------------------------------------------------------
	g.POST("/selected_attributes_values/create", h.CreateProductSelectedAttributesValues, cmiddleware.Authorization, ProductSelectedAttributesValuesCreateValidate)
	g.POST("/selected_attributes_values/:id/update", h.UpdateProductSelectedAttributesValues, cmiddleware.Authorization, ProductSelectedAttributesValuesUpdateValidate)
	g.DELETE("/selected_attributes_values/:id/delete", h.DeleteProductSelectedAttributesValues, cmiddleware.Authorization)
	g.GET("/selected_attributes_values", h.GetAllProductSelectedAttributesValues, cmiddleware.Authorization)
	g.GET("/selected_attributes_values/dropdown", h.GetAllProductSelectedAttributesValuesDropdown, cmiddleware.Authorization)

	//---------------------Product Bundles-----------------------------------------------------------------------------------------------------------------------
	g.POST("/bundle/create", h.CreateBundles, cmiddleware.Authorization, ProductBundleCreateValidate)
	g.POST("/bundle/:id/update", h.UpdateBundle, cmiddleware.Authorization, ProductBundleUpdateValidate)
	g.DELETE("/bundle/:id/delete", h.DeleteBundle, cmiddleware.Authorization)
	g.GET("/bundle/:id", h.GetOneBundle, cmiddleware.Authorization)
	g.GET("/bundle", h.GetAllBundles, cmiddleware.Authorization)
	g.GET("/bundle/dropdown", h.GetAllBundlesDropdown, cmiddleware.Authorization)

	//---------------------Product Template------------------
	g.POST("/create", h.CreateProductDetails, cmiddleware.Authorization, ProductTemplateCreateValidate)
	// g.POST("/create", h.CreateProductDetailsEvent, cmiddleware.Authorization, ProductTemplateCreateValidate)

	g.POST("/:id/update", h.UpdateProduct, cmiddleware.Authorization, ProductTemplateUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteProduct, cmiddleware.Authorization)
	g.GET("/:id", h.GetProductView, cmiddleware.Authorization)
	g.GET("", h.GetAllProducts, cmiddleware.Authorization)
	g.GET("/dropdown", h.GetAllProductsDropdown, cmiddleware.Authorization)
	g.GET("/sku_code", h.Skucodecheck, cmiddleware.Authorization)

	//---------------------Archive or Unarchive------------------
	// g.POST("/archive/:id/update", h.UpdateProductArchiveStatus)

	//---------------------Product Variant------------------
	g.POST("/variant/create", h.CreateProductVariant, cmiddleware.Authorization, ProductVariantCreateValidate)
	g.POST("/variant/:id/update", h.UpdateProductVariant, cmiddleware.Authorization, ProductVariantUpdateValidate)
	g.DELETE("/variant/:id/delete", h.DeleteProductVariant, cmiddleware.Authorization)
	g.GET("/variant/:id", h.GetProductVariantView, cmiddleware.Authorization)
	g.GET("/variant", h.GetAllProductVariants, cmiddleware.Authorization)
	g.GET("/variant/dropdown", h.GetAllProductVariantsDropdown, cmiddleware.Authorization)
	g.POST("/variant/variant_generation", h.VariantGeneration, cmiddleware.Authorization)
	g.GET("/variant/:id/qrcode", h.Qrcode, cmiddleware.Authorization)
	g.GET("/variant/:id/download", h.DownloadProductPDF, cmiddleware.Authorization)

	//---------------------Favourite Product Variant---------------------------
	g.POST("/:id/favourite", h.FavouriteProducts, cmiddleware.Authorization)
	g.POST("/:id/unfavourite", h.UnFavouriteProducts, cmiddleware.Authorization)
	g.GET("/favourite_list", h.FavouriteProductsView, cmiddleware.Authorization)

	//---------------------Bulk------------------
	g.POST("/bulk_create", h.CreateBulkProducts, cmiddleware.Authorization)

	//--------------------Channel Api's ------------------------------------------------------------------
	g.POST("/channels/upsert", h.ChannelProductUpsert, cmiddleware.Authorization)
	g.POST("/channels/variant/upsert", h.ChannelProductVariantUpsert, cmiddleware.Authorization)
	g.GET("/channels", h.GetAllChannelProducts, cmiddleware.Authorization)
	g.GET("/channels/dropdown", h.GetAllChannelProductsDropdown, cmiddleware.Authorization)

	//-------------------Filter Tabs-----------------------------------------------------------------
	g.GET("/variant/:id/filter_module/:tab", h.GetProductVariantTab, cmiddleware.Authorization)

}
