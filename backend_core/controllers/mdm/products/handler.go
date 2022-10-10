package products

import (
	"encoding/json"

	"errors"
	"fmt"

	// "log"
	// "os"
	"strconv"

	// eda "fermion/backend_core/controllers/eda_core"
	// consume "fermion/backend_core/controllers/consume"
	mdm_base "fermion/backend_core/controllers/mdm/base"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"

	// "fermion/backend_core/pkg/util/emitter"
	// "fermion/backend_core/pkg/util/events"
	"fermion/backend_core/pkg/util/emitter"
	"fermion/backend_core/pkg/util/events"
	"fermion/backend_core/pkg/util/helpers"

	res "fermion/backend_core/pkg/util/response"

	// "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	// "github.com/segmentio/kafka-go"
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

type handler struct {
	service      Service
	base_service mdm_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := mdm_base.NewServiceBase()
	return &handler{service, base_service}
}

func (h *handler) Init() {
	// go CreateProductConsumer(context.Background())
	// go InitConsumers(h)
	// InitConsumers(h)
}

//-----------------------------Product Brand-------------------------------------------------------------------------------------------

// CreateProductBrand godoc
// @Summary     CreateProductBrand
// @Description  Create ProductBrand
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   BrandRequestAndResponseDTO body BrandRequestAndResponseDTO true "create a Product Brand"
// @Success      200  {object}  BrandRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/brand/create [post]
func (h *handler) CreateProductBrand(c echo.Context) error {
	data := c.Get("product_brand").(*BrandRequestAndResponseDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	var brand mdm.ProductBrand
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &brand)

	err := h.service.CreateBrand(&brand, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Brand Created", data)
}

// UpdateProductBrand godoc
// @Summary      UpdateProductBrandDetails
// @Description  Update a ProductBrandDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   BrandRequestAndResponseDTO body BrandRequestAndResponseDTO true "update a Product Brand"
// @Success      200  {object}  BrandRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/brand/{id}/update [post]
func (h *handler) UpdateProductBrand(c echo.Context) error {
	data := c.Get("product_brand").(*BrandRequestAndResponseDTO)
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	var brand mdm.ProductBrand
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &brand)

	err = h.service.UpdateBrand(&brand, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Brand Updated", map[string]interface{}{"product_brand_id": id})
}

// DeleteProductBrand godoc
// @Summary      DeleteProductBrandDetails
// @Description  DeleteProductBrandDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/brand/{id}/delete [delete]
func (h *handler) DeleteProductBrand(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteBrand(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Brand Deleted", map[string]interface{}{"deleted_id": id})
}

// GetProductBrand godoc
// @Summary     GetProductBrand
// @Description  ProductBrand List
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  BrandRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/brand [get]
func (h *handler) GetAllBrands(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBrand(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Brands retrived successfully", resp, p)
}

// GetProductBrandDropdown godoc
// @Summary     GetProductBrand
// @Description  ProductBrand List
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  BrandRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/brand/dropdown [get]
func (h *handler) GetAllBrandsDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBrand(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Brands retrived successfully", resp, p)
}

func (h *handler) SearchBrand(c echo.Context) (err error) {

	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchBrand(q, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "OK", result)
}

//------------------Product Category----------------------------------------------------------------------------------------------------------

// CreateProductcategory godoc
// @Summary     CreateProductCategory
// @Description  Create ProductCategory
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   CategoryAndSubcategoryRequestAndResponseDTO body CategoryAndSubcategoryRequestAndResponseDTO true "create a Product category"
// @Success      200  {object}  CategoryAndSubcategoryRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/category/create [post]
func (h *handler) CreateProductCategory(c echo.Context) error {
	data := c.Get("product_category").(*CategoryAndSubcategoryRequestAndResponseDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	var category mdm.ProductCategory
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &category)

	err := h.service.CreateCategory(&category, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Category Created", data)
}

// UpdateProductCategory godoc
// @Summary      UpdateProductCategoryDetails
// @Description  Update a ProductCategoryDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   CategoryAndSubcategoryRequestAndResponseDTO body CategoryAndSubcategoryRequestAndResponseDTO true "update a Product category"
// @Success      200  {object}  CategoryAndSubcategoryRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/category/{id}/update [post]
func (h *handler) UpdateProductCategory(c echo.Context) error {
	data := c.Get("product_category").(*CategoryAndSubcategoryRequestAndResponseDTO)
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	var category mdm.ProductCategory
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &category)

	err = h.service.UpdateCategory(&category, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Category Updated", map[string]interface{}{"Category_id": id})
}

// DeleteProductCategory godoc
// @Summary      DeleteProductCategoryDetails
// @Description  DeleteProductCategoryDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/category/{id}/delete [delete]
func (h *handler) DeleteProductCategory(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteCategory(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Category Deleted", map[string]interface{}{"deleted_id": id})
}

// GetAllProductCategory godoc
// @Summary     ListProductCategory
// @Description  List of  ProductCategories
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  CategoryAndSubcategoryRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/category [get]
func (h *handler) GetAllProductCategory(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllCategory(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Categories retrived successfully", resp, p)
}

// GetAllProductCategoryDropdown godoc
// @Summary     ListProductCategory
// @Description  List of  ProductCategories
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  CategoryAndSubcategoryRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/category/dropdown [get]
func (h *handler) GetAllProductCategoryDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllCategory(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Categories retrived successfully", resp, p)
}

func (h *handler) GetAllSubCategory(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllSubCategory(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Sub-Categories retrived successfully", resp)
}

func (h *handler) GetAllSubCategoryDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllSubCategory(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Sub-Categories retrived successfully", resp)
}

func (h *handler) SearchCategory(c echo.Context) (err error) {
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	q := c.QueryParam("q")
	result, err := h.service.SearchCategory(q, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "OK", result)
}

func (h *handler) SearchCategoryDropdown(c echo.Context) (err error) {
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	q := c.QueryParam("q")
	result, err := h.service.SearchCategory(q, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "OK", result)
}

//------------------Product Base Attributes-------------------------------------------------------------------------------

// @Summary     CreateProductBaseAttributes
// @Description  Create ProductBaseAttributes
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   ProductBaseAttributesValuesRequestAndResponseDTO body ProductBaseAttributesValuesRequestAndResponseDTO true "create a Product BaseAttributes"
// @Success      200  {object}  ProductBaseAttributesValuesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes/create [post]
func (h *handler) CreateProductBaseAttributes(c echo.Context) error {
	data := c.Get("product_base_attributes").(*ProductBaseAttributesRequestAndResponseDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	var baseAttribute mdm.ProductBaseAttributes
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &baseAttribute)

	err := h.service.CreateBaseAttribute(&baseAttribute, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product BaseAttributes Created", data)
}

// UpdateProductBaseAttributes godoc
// @Summary      UpdateProductBaseAttributesDetails
// @Description  Update a ProductBaseAttributesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   ProductBaseAttributesRequestAndResponseDTO body ProductBaseAttributesRequestAndResponseDTO true "update a Product BaseAttributes"
// @Success      200  {object}  ProductBaseAttributesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes/{id}/update [post]
func (h *handler) UpdateProductBaseAttributes(c echo.Context) error {
	data := c.Get("product_base_attributes").(*mdm.ProductBaseAttributes)
	var id = c.Param("id")
	// host := c.Request().Host
	// scheme := c.Scheme()
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	var baseAttribute mdm.ProductBaseAttributes
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &baseAttribute)

	err = h.service.UpdateBaseAttribute(&baseAttribute, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product BaseAttributes Updated", map[string]interface{}{"BaseAttributes_id": id})
}

// DeleteProductBaseAttributes godoc
// @Summary      DeleteProductBaseAttributesDetails
// @Description  DeleteProductBaseAttributesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes/{id}/delete [delete]
func (h *handler) DeleteProductBaseAttributes(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteBaseAttribute(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product BaseAttributes Deleted", map[string]interface{}{"deleted_id": id})
}

// ListofProductBaseAttributes godoc
// @Summary     ListProductBaseAttributes
// @Description  List of  ProductBaseAttributes
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductBaseAttributesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes [get]
func (h *handler) GetAllProductBaseAttributes(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBaseAttribute(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "BaseAttributes retrived successfully", resp, p)
}

// ListofProductBaseAttributesDropdown godoc
// @Summary     ListProductBaseAttributes
// @Description  List of  ProductBaseAttributes
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductBaseAttributesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes/dropdown [get]
func (h *handler) GetAllProductBaseAttributesDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBaseAttribute(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "BaseAttributes retrived successfully", resp, p)
}

//------------------Product Base Attributes Values-----------------------------------------------------------------------------------------

// CreateProductBaseAttributesValues godoc
// @Summary     CreateProductBaseAttributesValues
// @Description  Create ProductBaseAttributesValues
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   ProductBaseAttributesValuesRequestAndResponseDTO body ProductBaseAttributesValuesRequestAndResponseDTO true "create a Product BaseAttributesValues"
// @Success      200  {object}  ProductBaseAttributesValuesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes_values/create [post]
func (h *handler) CreateProductBaseAttributesValues(c echo.Context) error {
	data := c.Get("product_base_attributes_values").(*ProductBaseAttributesValuesRequestAndResponseDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	var baseAttributeVAlue mdm.ProductBaseAttributesValues
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &baseAttributeVAlue)

	err := h.service.CreateBaseAttributeValue(&baseAttributeVAlue, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product BaseAttributesValues Created", map[string]interface{}{"BaseAttributesValues_Details": data})
}

// UpdateProductBaseAttributesValues godoc
// @Summary      UpdateProductBaseAttributesValuesDetails
// @Description  Update a ProductBaseAttributesValuesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   ProductBaseAttributesValuesRequestAndResponseDTO body ProductBaseAttributesValuesRequestAndResponseDTO true "update a Product BaseAttributes"
// @Success      200  {object}  ProductBaseAttributesValuesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes_values/{id}/update [post]
func (h *handler) UpdateProductBaseAttributesValues(c echo.Context) error {
	data := c.Get("product_base_attributes_values").(*mdm.ProductBaseAttributesValues)
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	var baseAttributeVAlue mdm.ProductBaseAttributesValues
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &baseAttributeVAlue)

	err = h.service.UpdateBaseAttributeValue(&baseAttributeVAlue, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Base AttributesValues Updated", map[string]interface{}{"BaseAttributesValues_id": id})
}

// DeleteProductBaseAttributesValues godoc
// @Summary      DeleteProductBaseAttributesValuesDetails
// @Description  DeleteProductBaseAttributesValuesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes_values/{id}/delete [delete]
func (h *handler) DeleteProductBaseAttributesValues(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteBaseAttributeValue(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product BaseAttributes values Deleted", map[string]interface{}{"deleted_id": id})
}

// ListofProductBaseAttributesValues godoc
// @Summary     ListProductBaseAttributesValues
// @Description  List of  ProductBaseAttributesValues
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductBaseAttributesValuesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes_values [get]
func (h *handler) GetAllProductBaseAttributesValues(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBaseAttributeValue(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Base Attributes Values retrived successfully", resp, p)
}

// ListofProductBaseAttributesValuesDropdown godoc
// @Summary     ListProductBaseAttributesValues
// @Description  List of  ProductBaseAttributesValues
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductBaseAttributesValuesRequestAndResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/base_attributes_values/dropdown [get]
func (h *handler) GetAllProductBaseAttributesValuesDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBaseAttributeValue(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Base Attributes Values retrived successfully", resp, p)
}

// ------------------Product Selected Attributes-----------------------------------------------------------------------------------------
// CreateProductSelectedAttributes godoc
// @Summary     CreateProductSelectedAttributes
// @Description  Create ProductSelectedAttributes
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   ProductSelectedAttributesRequestAndREsponseDTO body ProductSelectedAttributesRequestAndREsponseDTO true "create a Product Selected Attributes"
// @Success      200  {object}  ProductSelectedAttributesRequestAndREsponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes/create [post]
func (h *handler) CreateProductSelectedAttributes(c echo.Context) error {
	data := c.Get("product_selected_attributes").(*ProductSelectedAttributesRequestAndREsponseDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	var selectedBaseAttributeVAlue mdm.ProductSelectedAttributes
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &selectedBaseAttributeVAlue)

	err := h.service.CreateSelectedAttribute(&selectedBaseAttributeVAlue, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product Selected Attributes Created", map[string]interface{}{"SelectedAttributes_Details": data})
}

// UpdateProductSelectedAttributes godoc
// @Summary      UpdateProductSelectedAttributesDetails
// @Description  Update a ProductSelectedAttributesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   ProductSelectedAttributesRequestAndREsponseDTO body ProductSelectedAttributesRequestAndREsponseDTO true "update a Product Selected Attributes"
// @Success      200  {object}  ProductSelectedAttributesRequestAndREsponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes/{id}/update [post]
func (h *handler) UpdateProductSelectedAttributes(c echo.Context) error {
	data := c.Get("product_selected_attributes").(*ProductSelectedAttributesRequestAndREsponseDTO)
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	var selectedBaseAttributeVAlue mdm.ProductSelectedAttributes
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &selectedBaseAttributeVAlue)

	err = h.service.UpdateSelectedAttribute(&selectedBaseAttributeVAlue, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product Selected Attributes Updated", map[string]interface{}{"SelectedAttributes_id": id})
}

// DeleteProductSelectedAttributes godoc
// @Summary      DeleteProductSelectedAttributesDetails
// @Description  DeleteProductSelectedAttributesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes/{id}/delete [delete]
func (h *handler) DeleteProductSelectedAttributes(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteSelectedAttribute(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product Selected Attributes Deleted", map[string]interface{}{"deleted_id": id})
}

// GetAllProductSelectedAttributes godoc
// @Summary     ListProductSelectedAttributes
// @Description  List of  ProductSelectedAttributes
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductSelectedAttributesRequestAndREsponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes [get]
func (h *handler) GetAllProductSelectedAttributes(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllSelectedAttribute(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "SelectedAttributes retrived successfully", resp, p)
}

// GetAllProductSelectedAttributesDropdown godoc
// @Summary     ListProductSelectedAttributes
// @Description  List of  ProductSelectedAttributes
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductSelectedAttributesRequestAndREsponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes/dropdown [get]
func (h *handler) GetAllProductSelectedAttributesDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllSelectedAttribute(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "SelectedAttributes retrived successfully", resp, p)
}

//------------------Product Selected Attributes Values-----------------------------------------------------------------------------------------

// CreateProductSelectedAttributesValues godoc
// @Summary     CreateProductSelectedAttributesValues
// @Description  Create ProductSelectedAttributesValues
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   ProductSelectedAttributesValuesRequestAndResponseeDTO body ProductSelectedAttributesValuesRequestAndResponseeDTO true "create a Product SelectedAttributesValues"
// @Success      200  {object}  ProductSelectedAttributesValuesRequestAndResponseeDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes_values/create [post]
func (h *handler) CreateProductSelectedAttributesValues(c echo.Context) error {
	data := c.Get("product_selected_attributes_values").(*ProductSelectedAttributesValuesRequestAndResponseeDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	var selectedBaseAttributeVAlue mdm.ProductSelectedAttributesValues
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &selectedBaseAttributeVAlue)

	err := h.service.CreateSelectedAttributeValue(&selectedBaseAttributeVAlue, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product SelectedAttributesValues Created", data)
}

// UpdateProductSelectedAttributesValues godoc
// @Summary      UpdateProductSelectedAttributesValuesDetails
// @Description  Update a ProductSelectedAttributesValuesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   ProductSelectedAttributesValuesRequestAndResponseeDTO body ProductSelectedAttributesValuesRequestAndResponseeDTO true "update a Product SelectedAttributesValues"
// @Success      200  {object}  ProductSelectedAttributesValuesRequestAndResponseeDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes_values/{id}/update [post]
func (h *handler) UpdateProductSelectedAttributesValues(c echo.Context) error {
	data := c.Get("product_selected_attributes_values").(*ProductSelectedAttributesValuesRequestAndResponseeDTO)
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	var selectedBaseAttributeVAlue mdm.ProductSelectedAttributesValues
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &selectedBaseAttributeVAlue)

	err = h.service.UpdateSelectedAttributeValue(&selectedBaseAttributeVAlue, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "ProductSelected Attributes Values Updated", map[string]interface{}{"SelectedAttributesValues_id": id})
}

// DeleteProductBaseAttributesValues godoc
// @Summary      DeleteProductSelectedAttributesValuesDetails
// @Description  DeleteProductSelectedAttributesValuesDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes_values/{id}/delete [delete]
func (h *handler) DeleteProductSelectedAttributesValues(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteSelectedAttributeValue(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product Selected Attributes values Deleted", map[string]interface{}{"deleted_id": id})
}

// GetAllProductSelectedAttributesValues godoc
// @Summary     ListProductSelectedAttributesValues
// @Description  List of  ProductSelectedAttributesValues
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductSelectedAttributesValuesRequestAndResponseeDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes_values [get]
func (h *handler) GetAllProductSelectedAttributesValues(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllSelectedAttributeValue(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Selected Attributes Values retrived successfully", resp, p)
}

// GetAllProductSelectedAttributesValuesDropdown godoc
// @Summary     ListProductSelectedAttributesValues
// @Description  List of  ProductSelectedAttributesValues
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  ProductSelectedAttributesValuesRequestAndResponseeDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/selected_attributes_values/dropdown [get]
func (h *handler) GetAllProductSelectedAttributesValuesDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllSelectedAttributeValue(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Selected Attributes Values retrived successfully", resp, p)
}

//------------------Product Bundles------------------------------------------------------------------------------------------------------------

// CreateBundleDetails godoc
// @Summary      do a CreateBundleDetails
// @Description  Create a BundleDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   CreateBundlePayload body CreateBundlePayload true "create a Bundle"
// @Success      200  {object}  CreateBundlePayload
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/bundle/create [post]
func (h *handler) CreateBundles(c echo.Context) error {
	data := c.Get("product_bundle").(*CreateBundlePayload)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	Bundle, err := h.service.CreateBundle(*data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Bundles Created", Bundle)
}

// UpdateBundleDetails godoc
// @Summary      UpdateBundleDetails
// @Description  Update a Bundle Details
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   CreateBundlePayload body CreateBundlePayload true "update a Bundle Details"
// @Success      200  {object}  CreateBundlePayload
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/bundle/{id}/update [post]
func (h *handler) UpdateBundle(c echo.Context) error {
	data := c.Get("product_bundle").(*CreateBundlePayload)
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateBundle(*data, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Bundle Updated", map[string]interface{}{"bundle_id": id})
}

// DeleteBundleDetails godoc
// @Summary      DeleteBundleDetails
// @Description  DeleteBundleDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/bundle/{id}/delete [delete]
func (h *handler) DeleteBundle(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteBundle(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Bundle Deleted", map[string]interface{}{"deleted_id": id})
}

// GetBundleView godoc
// @Summary      ViewBundleDetails
// @Description  View a BundleDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  BudleResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/bundle/{id} [get]
func (h *handler) GetOneBundle(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	query["id"] = val
	result, err := h.service.GetOneBundle(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Bundle Fetched", result)
}

// GetAllBundles godoc
// @Summary     List of Bundles
// @Description  List of  Bundles
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  BudleResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/bundle [get]
func (h *handler) GetAllBundles(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBundle(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Bundles data retrived successfully", resp, p)
}

// GetAllBundlesDropdown godoc
// @Summary     List of Bundles
// @Description  List of  Bundles
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  BudleResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/bundle/dropdown [get]
func (h *handler) GetAllBundlesDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllBundle(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Bundles data retrived successfully", resp, p)
}

//------------------Product BUlk---------------------------------------------------------------------------------------------

// CreateBulkProducts godoc
// @Summary      do a CreateBulkProduct
// @Description  Create Bulk Product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   CreateProductTemplatePayload body CreateProductTemplatePayload true "create a Product"
// @Success      200  {array}  CreateProductTemplatePayload
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/bulk_create [post]
func (h *handler) CreateBulkProducts(c echo.Context) (err error) {
	var data = new([]mdm.ProductTemplate)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateProduct(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "Multiple records created successfully", data)
}

//-----------------------Product Related API"s -----------------------------------------------------------------------------

// Skucodecheck godoc
// @Summary Check  Productvariant and template Skucode
// @Description Check Productvariant and Template Skucode values
// @Tags Products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param sku query string true "sku"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/sku_code [get]
func (h *handler) Skucodecheck(c echo.Context) (err error) {
	sku := c.QueryParam("sku")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	er := h.service.Skucodecheck(sku, token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "SKU code is available", er)
	}
	return res.RespErr(c, errors.New("sku is already taken!!!! Try something else"))
}

// VariantGeneration godoc
// @Summary  Productvariant VariantGeneration
// @Description  Productvariant VariantGeneration
// @Tags Products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/variant/variant_generation [post]
func (h *handler) VariantGeneration(c echo.Context) error {
	type AttributesValues struct {
		Values [][]map[string]interface{} `json:"attribute_values"`
	}
	// token_id := c.Get("TokenUserID").(string)
	// access_template_id := c.Get("AccessTemplateId").(string)
	var values AttributesValues
	if err := c.Bind(&values); err != nil {
		return res.RespErr(c, err)
	}

	response, err := helpers.VariantGeneration(values.Values)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Variant Generated Successfully", response)
}

// Qrcode godoc
// @Summary print  Productvariant Qrcode
// @Description print Productvariant Qrcode
// @Tags Products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/variant/{id}/qrcode [post]
func (h *handler) Qrcode(c echo.Context) error {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	response := h.service.Qrcode(uint(id), token_id, access_template_id)
	fmt.Println(response)
	return res.RespSuccess(c, "Product Variant Generated Successfully", "https://drive.google.com/file/d/1ODtoo6yofXbHlAKSBfJnMfxvyF2d2nWu/view?usp=sharing")
}

// DownloadProductPDF godoc
// @Summary Download Productvariant
// @Description Download Productvariant
// @Tags Products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/variant/{id}/download_pdf [post]
func (h *handler) DownloadProductPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

//-----------------------------Favourite and UnFavouriteProducts--------------------------------------------------------------------

// FavouriteProducts godoc
// @Summary FavouriteProducts
// @Description FavouriteProducts
// @Tags Products
// @Accept  json
// @Produce  json
// @param id path string true "ProductVariant ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/{id}/favourite [post]
func (h *handler) FavouriteProducts(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	q := map[string]interface{}{
		"ID": ID,
	}
	_, er := h.service.GetVariantView(q, token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteProducts(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteProducts godoc
// @Summary UnFavouriteProducts
// @Description UnFavouriteProducts
// @Tags Products
// @Accept  json
// @Produce  json
// @param id path string true "Product ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/{id}/unfavourite [post]
func (h *handler) UnFavouriteProducts(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	// access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteProducts(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product is Unmarked as Favourite", map[string]string{"record_id": id})
}

// FavouriteProductsView godoc
// @Summary Get all favourite products
// @Description Get all favourite products
// @Tags Products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} VariantResponseDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/favourite_list [get]
func (h *handler) FavouriteProductsView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	// access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrProduct(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	//var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.base_service.GetFavProductList(data, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Products list Retrieved Successfully", result, p)
}

//---------------------Product Template--------------------------------------------------------------------------------------

// CreateProductDetails godoc
// @Summary      do a CreateProductDetails
// @Description  Create a ProductDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   CreateProductTemplatePayload body CreateProductTemplatePayload true "create a Product"
// @Success      200  {object}  TemplateReponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/create [post]
func (h *handler) CreateProductDetails(c echo.Context) error {
	data := c.Get("product_template").(*CreateProductTemplatePayload)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	user_id, _ := strconv.Atoi(token_id)
	t := uint(user_id)
	// data.CreatedByID = &t
	query := map[string]interface{}{
		"id": t,
	}
	host := c.Request().Host
	scheme := c.Scheme()

	// data.CreatedByID = helpers.ConvertStringToUint(s)

	//=============create Template =======================
	user_token := c.Request().Header.Get("Authorization")
	Template, err := h.service.CreateTemplate(*data, token_id, access_template_id, query, host, scheme, user_token)
	if err != nil {
		return res.RespErr(c, err)
	}

	//============catalogue sync event=====================
	template_id := Template.(mdm.ProductTemplate).ID
	emitter := emitter.GetObj()
	emitter.Emit(events.ProductCatalogueSyncEvent+"template", c, template_id)

	return res.RespSuccess(c, "Product Created", Template)
}

// UpdateProduct godoc
// @Summary      do a UpdateProductDetails
// @Description  Update a ProductDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   CreateProductTemplatePayload body CreateProductTemplatePayload true "update a Product"
// @Success      200  {object}  TemplateReponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/{id}/update [post]
func (h *handler) UpdateProduct(c echo.Context) error {
	data := c.Get("product_template").(*CreateProductTemplatePayload)

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	user_id, _ := strconv.Atoi(token_id)
	t := uint(user_id)
	data.CreatedByID = &t
	query2 := map[string]interface{}{
		"id": t,
	}
	host := c.Request().Host
	scheme := c.Scheme()

	user_token := c.Request().Header.Get("Authorization")
	err = h.service.UpdateTemplate(*data, query, token_id, access_template_id, query2, host, scheme, user_token)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Updated", map[string]interface{}{"product_id": id})
}

// DeleteProduct godoc
// @Summary      do a DeleteProductDetails
// @Description  Delete a ProductDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/{id}/delete [delete]
func (h *handler) DeleteProduct(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteTemplate(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Deleted", map[string]interface{}{"deleted_id": id})
}

// GetProductView godoc
// @Summary      ViewProductDetails
// @Description  View a ProductDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  TemplateReponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/{id} [get]
func (h *handler) GetProductView(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	query := map[string]interface{}{
		"id": id,
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetTemplateView(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Fetched", result)
}

// GetAllProducts godoc
// @Summary      ListProducts with MISC OPTIONS of pagination, search,Filter and sort
// @Description  List of All Products  with MISC OPTIONS of pagination, search,Filter and sort
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  TemplateReponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products [get]
func (h *handler) GetAllProducts(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllTemplate(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data retrived successfully", resp, p)
}

// GetAllProductsDropdown godoc
// @Summary      ListProducts with MISC OPTIONS of pagination, search,Filter and sort
// @Description  List of All Products  with MISC OPTIONS of pagination, search,Filter and sort
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  TemplateReponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/dropdown [get]
func (h *handler) GetAllProductsDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllTemplate(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data retrived successfully", resp, p)
}

// ------------------Product Variant----------------------------------------------------------------------------------------
// CreateProductVariant godoc
// @Summary     CreateProductVariantDetails
// @Description  Create ProductVariantDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   CreateProductVariantDTO body CreateProductVariantDTO true "create a Productvariant"
// @Success      200  {object}  VariantResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/variant/create [post]
func (h *handler) CreateProductVariant(c echo.Context) error {
	data := c.Get("product_variant").(*CreateProductVariantDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)

	//===============create new variant product =================
	err := h.service.CreateVariant(*data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Variant Created", data)
}

// UpdateProductVariant godoc
// @Summary      UpdateProductVariantDetails
// @Description  Update a ProductVariantDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Param   CreateProductVariantDTO body CreateProductVariantDTO true "update a Productvariant"
// @Success      200  {object}  VariantResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/variant/{id}/update [post]
func (h *handler) UpdateProductVariant(c echo.Context) error {
	data := c.Get("product_variant").(*CreateProductVariantDTO)
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespErr(c, err)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)

	//============update variant=================
	err = h.service.UpdateVariant(*data, query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Variant Updated", map[string]interface{}{"product_variant_id": id})
}

// DeleteProductVariant godoc
// @Summary      DeleteProductVariantDetails
// @Description  DeleteProductVariantDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/variant/{id}/delete [delete]
func (h *handler) DeleteProductVariant(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	val, err := strconv.Atoi(id)
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	query["id"] = val
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id

	//==============delete variant=================
	err = h.service.DeleteVariant(query, token_id, access_template_id)
	if err != nil && err.Error() == "record not found" {
		return res.RespErr(c, err)
	}
	if err != nil {
		return res.RespError(c, &res.ErrServerError)
	}
	return res.RespSuccess(c, "Product Variant Deleted", map[string]interface{}{"deleted_id": id})
}

// GetProductVariantView godoc
// @Summary      ViewProductDetails
// @Description  View a ProductDetails
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id path int true "id"
// @Success      200  {object}  VariantResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/variant/{id} [get]
func (h *handler) GetProductVariantView(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	query["id"] = id
	result, err := h.service.GetVariantView(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Product Fetched", result)
}

// GetAllProductVariants godoc
// @Summary      ListProductVariants with MISC OPTIONS of pagination, search,Filter and sort
// @Description  List of All ProductVariants  with MISC OPTIONS of pagination, search,Filter and sort
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  VariantResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/variant [get]
func (h *handler) GetAllProductVariants(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllVariant(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "data retrived successfully", resp, p)
}

// GetAllProductVariantsDropdown godoc
// @Summary      ListProductVariants with MISC OPTIONS of pagination, search,Filter and sort
// @Description  List of All ProductVariants  with MISC OPTIONS of pagination, search,Filter and sort
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param       filters query string    false "filters"
// @param        per_page query int    false "per_page"
// @param        page_no  query int    false "page_no"
// @param        sort     query string false "sort"
// @Success      200  {array}  VariantResponseDTO
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/products/variant/dropdown [get]
func (h *handler) GetAllProductVariantsDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.GetAllVariant(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "data retrived successfully", resp, p)
}

// ---------------------Product Upset Api's ---------------------------------------------------------------------------------------------------------
func (h *handler) ChannelProductUpsert(c echo.Context) (err error) {
	var arrayData []interface{}
	var objectData interface{}
	var data interface{}

	c.Bind(&data)

	jsonData, _ := json.Marshal(data)

	err = json.Unmarshal(jsonData, &arrayData)
	if err != nil {
		_ = json.Unmarshal(jsonData, &objectData)
		arrayData = append(arrayData, objectData)
	}

	TokenUserID := c.Get("TokenUserID").(string)
	// access_template_id := c.Get("AccessTemplateId").(string)
	msg, err := h.service.UpsertProductTemplate(arrayData, TokenUserID)
	if err != nil {
		return res.RespErr(c, err)
	}
	// fmt.Println("upsert api response-->>>",msg)

	return res.RespSuccess(c, "Multiple records executed successfully", msg)
}
func (h *handler) ChannelProductVariantUpsert(c echo.Context) (err error) {
	var arrayData []interface{}
	var objectData interface{}

	var data interface{}

	c.Bind(&data)

	jsonData, _ := json.Marshal(data)

	err = json.Unmarshal(jsonData, &arrayData)
	if err != nil {
		_ = json.Unmarshal(jsonData, &objectData)
		arrayData = append(arrayData, objectData)
	}

	TokenUserID := c.Get("TokenUserID").(string)
	// access_template_id := c.Get("AccessTemplateId").(string)
	msg, err := h.service.UpsertProductVariant(arrayData, TokenUserID)
	if err != nil {
		return res.RespErr(c, err)
	}
	// fmt.Println("upsert api response-->>>", msg)
	return res.RespSuccess(c, "Multiple records executed successfully", msg)
}
func (h *handler) GetAllChannelProducts(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)

	query_params := new(QueryParamsDTO)

	c.Bind(query_params)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	mdata, _ := json.Marshal(query_params)
	json.Unmarshal(mdata, &p)

	resp, err := h.service.GetAllChannelProducts(p, *query_params, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Channel Products data Retrieved successfully", resp, p)
}

func (h *handler) GetAllChannelProductsDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)

	query_params := new(QueryParamsDTO)

	c.Bind(query_params)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	mdata, _ := json.Marshal(query_params)
	json.Unmarshal(mdata, &p)

	resp, err := h.service.GetAllChannelProducts(p, *query_params, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Channel Products data Retrieved successfully", resp, p)
}

// GetProductVariantTabs godoc
// @Summary GetProductVariantTabs
// @Summary GetProductVariantTabs
// @Description GetProductVariantTabs
// @Tags Products
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @param tab path string true "tab"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/products/variant/{id}/filter_module/{tab} [get]
func (h *handler) GetProductVariantTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	product_variant_id := c.Param("id")
	tab_name := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetProductVariantTab(product_variant_id, tab_name, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
