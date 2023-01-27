package pricing

import (
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	mdm_base "fermion/backend_core/controllers/mdm/base"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
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
	pricing_repo mdm_repo.Pricing
}

var PricingHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if PricingHandler != nil {
		return PricingHandler
	}
	PricingHandler = &handler{
		NewService(),
		mdm_base.NewServiceBase(),
		mdm_repo.NewPricing(),
	}
	return PricingHandler
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
func (h *handler) CreatePricingEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("pricing"),
	}
	eda.Produce(eda.CREATE_PRICING, request_payload)
	return res.RespSuccess(c, "Pricing Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreatePricing(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	pricing := new(shared_pricing_and_location.Pricing)
	helpers.JsonMarshaller(data, pricing)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	edaMetaData.Query = map[string]interface{}{
		"price_list_name": pricing.PriceListName,
		"company_id":      edaMetaData.CompanyId,
	}
	_, err := h.service.FindPricing(edaMetaData)
	if err == nil {
		responseMessage.ErrorMessage = fmt.Errorf("price list already exists")
		// eda.Produce(eda.CREATE_PRICING_ACK, responseMessage)
		return
	}

	if pricing.PriceListId == 1 {
		salesPriceList := new(shared_pricing_and_location.SalesPriceList)
		helpers.JsonMarshaller(data, salesPriceList)
		err := h.service.CreateSalesPricing(edaMetaData, salesPriceList)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_PRICING_ACK, responseMessage)
			return
		}
		pricing.SalesPriceListId = &salesPriceList.ID

	} else if pricing.PriceListId == 2 {
		purchasePriceList := new(shared_pricing_and_location.PurchasePriceList)
		helpers.JsonMarshaller(data, purchasePriceList)
		err := h.service.CreatePurchasePricing(edaMetaData, purchasePriceList)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_PRICING_ACK, responseMessage)
			return
		}
		pricing.PurchasePriceListId = &purchasePriceList.ID

	} else if pricing.PriceListId == 3 {
		transferPriceList := new(shared_pricing_and_location.TransferPriceList)
		helpers.JsonMarshaller(data, transferPriceList)
		err := h.service.CreateTransferPricing(edaMetaData, transferPriceList)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_PRICING_ACK, responseMessage)
			return
		}
		pricing.TransferPriceListId = &transferPriceList.ID
	} else {
		responseMessage.ErrorMessage = errors.New("oops! invalid price_list_id")
		// eda.Produce(eda.CREATE_PRICING_ACK, responseMessage)
		return
	}
	err = h.service.CreatePricing(edaMetaData, pricing)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_PRICING_ACK, responseMessage)
		return
	}

	// updating cache
	UpdatePricingInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{
		"created_id": pricing.ID,
	}
	// eda.Produce(eda.CREATE_PRICING_ACK, responseMessage)
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
func (h *handler) UpdatePricingEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	pricingId := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         pricingId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("pricing"),
	}

	eda.Produce(eda.UPDATE_PRICING, request_payload)
	return res.RespSuccess(c, "Update Pricing Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdatePricing(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	err := h.service.UpdatePricing(edaMetaData, data)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_PRICING_ACK, responseMessage)
		return
	}

	//update cache
	UpdatePricingInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{"updated_id": edaMetaData.Query["id"]}
	// eda.Produce(eda.UPDATE_PRICING_ACK, responseMessage)
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
	metaData := c.Get("MetaData").(core.MetaData)
	pricingId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": pricingId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.FindPricing(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response PricingDTO
	helpers.JsonMarshaller(result, &response)
	if response.Price_list_id == 1 {
		response.PurchasePriceList = nil
		response.TransferPriceList = nil
	} else if response.Price_list_id == 2 {
		response.SalesPriceList = nil
		response.TransferPriceList = nil
	} else {
		response.SalesPriceList = nil
		response.PurchasePriceList = nil
	}

	return res.RespSuccess(c, "data Retrieved successfully", response)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []PricingDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetPricingFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.ListPricing(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "data Retrieved successfully", response, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []PricingDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetPricingFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.ListPricing(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "data Retrieved successfully", response, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	pricingId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         pricingId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeletePricing(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// updating cache
	UpdatePricingInCache(metaData)

	return res.RespSuccess(c, "Delete Pricing successfully done", map[string]interface{}{"deleted_id": pricingId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	pricingId := c.Param("id")
	productVariantId := c.QueryParam("product_id")
	metaData.Query = map[string]interface{}{
		"id":         pricingId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	metaData.AdditionalFields = map[string]interface{}{
		"product_variant_id": productVariantId,
	}

	err = h.service.DeleteLineItems(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "line items successfully Deleted", map[string]interface{}{"pricing_id": pricingId, "product_id": productVariantId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	pricingId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         pricingId,
		"company_id": metaData.CompanyId,
	}
	_, err = h.service.FindPricing(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	query := map[string]interface{}{
		"ID":      pricingId,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.FavouritePricings(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Pricing is Marked as Favourite", map[string]string{"record_id": pricingId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	pricingId := c.Param("id")

	query := map[string]interface{}{
		"ID":      pricingId,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouritePricings(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Pricing is Unmarked as Favourite", map[string]string{"record_id": pricingId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	tokenUserId := int(metaData.TokenUserId)
	data, er := h.base_service.GetArrPricing(tokenUserId)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavPricingList(data, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Pricing list Retrieved Successfully", result, p)
}

// -------------------------Channel API's ----------------------------------------------------
func (h *handler) ChannelSalesPriceUpsertEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	var arrayData []interface{}
	var objectData interface{}
	var data interface{}

	c.Bind(&data)
	err = helpers.JsonMarshaller(data, &arrayData)
	if err != nil {
		helpers.JsonMarshaller(data, &objectData)
		arrayData = append(arrayData, objectData)
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      arrayData,
	}
	eda.Produce(eda.UPSERT_SALES_PRICE_LIST, request_payload)
	return res.RespSuccess(c, "Multiple Records Upsert Inprogress", edaMetaData.RequestId)
}
func (h *handler) ChannelSalesPriceUpsert(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	var data []interface{}
	helpers.JsonMarshaller(request["data"], &data)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	response, err := h.service.ChannelSalesPriceUpsert(edaMetaData, data)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPSERT_SALES_PRICE_LIST_ACK, responseMessage)
		return
	}
	responseMessage.Response = response
	// eda.Produce(eda.UPSERT_SALES_PRICE_LIST_ACK, responseMessage)
}
