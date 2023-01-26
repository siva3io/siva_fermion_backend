package pricing

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	mdm_base "fermion/backend_core/controllers/mdm/base"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

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
type handler struct {
	service      Service
	base_service mdm_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := mdm_base.NewServiceBase()
	return &handler{service, base_service}
}

// CreatePricing godoc
// @Summary      Do Create Price List
// @Description  Create a Price List
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   PurchasePriceListCreateRequest     body PurchasePriceListDTO    false "create a Purchase Price List"
// @Param   TransferPriceListCreateRequest     body  TransferPriceListDTO       false "create a Transfer Price List"
// @Param   SalesPriceListMarkupCreateRequest  body SalesPriceListDTO false "create a Sales Price List for Markup and Markdown By age"
// @Param   SalesPriceListAllItemCreateRequest body SalesPriceListDTO false "create a Sales Price List for Enter All Items"
// @Success      200  {object}  PricingCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pricing/create [post]
func (h *handler) CreatePricing(c echo.Context) (err error) {

	var reqPayload map[string]interface{}
	json.NewDecoder(c.Request().Body).Decode(&reqPayload)
	reqBytes, _ := json.Marshal(reqPayload)

	pricing := c.Get("pricing").(*shared_pricing_and_location.Pricing)
	json.Unmarshal(reqBytes, pricing)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	pricing.CreatedByID = t
	if pricing.Price_list_id == 1 {
		sales := c.Get("sales").(*shared_pricing_and_location.SalesPriceList)
		json.Unmarshal(reqBytes, sales)
		sales.CreatedByID = t
		err = h.service.CreateSalesPricing(pricing, sales, token_id, access_template_id)
		if err != nil {
			return res.RespErr(c, err)
		}
		return res.RespSuccess(c, "created successfully", map[string]interface{}{"created_sales_id": pricing.ID})
	} else if pricing.Price_list_id == 2 {
		purchase := c.Get("purchase").(*shared_pricing_and_location.PurchasePriceList)
		json.Unmarshal(reqBytes, purchase)
		purchase.CreatedByID = t
		err = h.service.CreatePurchasePricing(pricing, purchase, token_id, access_template_id)
		if err != nil {
			return res.RespErr(c, err)
		}
		return res.RespSuccess(c, "created successfully", map[string]interface{}{"created_purchase_id": pricing.ID})
	} else if pricing.Price_list_id == 3 {
		transfer := c.Get("transfer").(*shared_pricing_and_location.TransferPriceList)
		json.Unmarshal(reqBytes, transfer)
		transfer.CreatedByID = t
		err = h.service.CreateTransferPricing(pricing, transfer, token_id, access_template_id)
		if err != nil {
			return res.RespErr(c, err)
		}
		return res.RespSuccess(c, "created successfully", map[string]interface{}{"created_transfer_id": pricing.ID})
	} else {
		return res.RespErr(c, fmt.Errorf("invalid price list id"))
	}

	// return res.RespSuccess(c, "created successfully", map[string]interface{}{"created_id": pricing.ID})
}

// UpdatePriceList godoc
// @Summary      do  Update Price List
// @Description  Update a Price List
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   PurchasePriceListCreateRequest     body    PurchasePriceListDTO    false "create a Purchase Price List"
// @Param   TransferPriceListCreateRequest     body TransferPriceListDTO     false "create a Transfer Price List"
// @Param   SalesPriceListMarkupCreateRequest  body SalesPriceListDTO   false "create a Sales Price List for Markup and Markdown By age"
// @Param   SalesPriceListAllItemCreateRequest body  SalesPriceListDTO false "create a Sales Price List for Enter All Items"
// @Success      200  {object}  PricingUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pricing/{id}/edit [post]
func (h *handler) UpdatePricing(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var reqPayload map[string]interface{}
	json.NewDecoder(c.Request().Body).Decode(&reqPayload)
	reqBytes, _ := json.Marshal(reqPayload)

	err = h.service.UpdatePricing(uint(id), reqBytes, c, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Update Pricing successfully done", map[string]interface{}{"updated_pricing_id": id})
}

// FindPriceList godoc
// @Summary      Find a PriceList by id
// @Description  Find a PriceList by id
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  PricingListByIdResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/{id}/pricing [get]
func (h *handler) FindPricing(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.PricingByid(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "data Retrieved successfully", result)
}

// PricingList godoc
// @Summary      get all Price List and filter it by search query
// @Description  a get all Price List and filter it by search query
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters				  query     string false "filters"
// @Param        per_page 			 query     int false "per_page"
// @Param        page_no              query     int false "page_no"
// @Param        sort				  query     string false "sort"
// @Success      200  {object}  PricingListResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pricing [get]
func (h *handler) PricingList(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.ListPricing(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// PricingListDropdown godoc
// @Summary      get all Price List and filter it by search query
// @Description  a get all Price List and filter it by search query
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters				  query     string false "filters"
// @Param        per_page 			 query     int false "per_page"
// @Param        page_no              query     int false "page_no"
// @Param        sort				  query     string false "sort"
// @Success      200  {object}  PricingListResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pricing/dropdown [get]
func (h *handler) PricingListDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.ListPricing(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// DeletePricing godoc
// @Summary      delete a Price List
// @Description  delete a Price List
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  PricingDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pricing/{id}/delete [delete]
func (h *handler) DeletePricing(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeletePricing(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Delete Pricing successfully done", map[string]interface{}{"deleted_pricing_id": id})
}

// DeletePricingLineItems godoc
// @Summary      delete a DeletePricingLineItems
// @Description  delete a DeletePricingLineItems
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  PricingDeleteLineResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pricing/{id}/delete_line_items [delete]
func (h *handler) DeleteLineItems(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	prod_id := c.QueryParam("product_id")
	product_id, _ := strconv.Atoi(prod_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.DeleteLineItems(uint(id), uint(product_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "line items successfully Deleted", map[string]interface{}{"pricing_id": id, "product_id": product_id})
}

// PricingList godoc
// @Summary      get all Purchase Price List and filter it by search query
// @Description  a get all Purchase Price List and filter it by search query
// @Tags         Pricing
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters				  query     string false "filters"
// @Param        per_page 			 query     int false "per_page"
// @Param        page_no              query     int false "page_no"
// @Param        sort				  query     string false "sort"
// @Success      200  {object}  PricingListResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pricing [get]
func (h *handler) PurchasePricingList(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.ListPurchasePricing(p, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// FavouritePricings godoc
// @Summary FavouritePricings
// @Description FavouritePricings
// @Tags Pricing
// @Accept  json
// @Produce  json
// @param id path string true "Pricing ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pricing/{id}/favourite [post]
func (h *handler) FavouritePricings(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.PricingByid(uint(ID), token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouritePricings(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Pricing is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouritePricings godoc
// @Summary UnFavouritePricings
// @Description UnFavouritePricings
// @Tags Pricing
// @Accept  json
// @Produce  json
// @param id path string true "Pricing ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pricing/{id}/unfavourite [post]
func (h *handler) UnFavouritePricings(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouritePricings(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Pricing is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetFavouritePricing godoc
// @Summary Get all favourite pricing
// @Description Get all favourite pricing
// @Tags Pricing
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pricing/favourite_list [get]
func (h *handler) FavouritePricingView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	data, er := h.base_service.GetArrPricing(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavPricingList(data, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Pricing list Retrieved Successfully", result, p)
}

// -----------sales Price line Items-------------------------
func (h *handler) SalesPriceLineItems(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.ListSalesPiceLineItems(p, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// -------------------------Channel API's ----------------------------------------------------
func (h *handler) ChannelSalesPriceUpsert(c echo.Context) (err error) {
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
	msg, err := h.service.ChannelSalesPriceUpsert(arrayData, TokenUserID)
	if err != nil {
		return res.RespErr(c, err)
	}
	fmt.Println("upsert api response-->>>", msg)

	return res.RespSuccess(c, "Multiple records executed successfully", msg)
}
