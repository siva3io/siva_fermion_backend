package shipping_orders

import (
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/controllers/omnichannel/channel"
	shipping_base "fermion/backend_core/controllers/shipping/base"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"

	// "fermion/backend_core/internal/repository/shipping"
	"fermion/backend_core/pkg/util/emitter"
	"fermion/backend_core/pkg/util/events"
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
	base_service shipping_base.ServiceBase
}

var ShippingOrdersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ShippingOrdersHandler != nil {
		return ShippingOrdersHandler
	}
	service := NewService()
	base_service := shipping_base.NewServiceBase()
	ShippingOrdersHandler = &handler{service, base_service}
	return ShippingOrdersHandler
}

// CreateShippingOrder godoc
// @Summary CreateShippingOrder
// @Description Create a ShippingOrder
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   ShippingOrderCreateRequest body ShippingOrderRequest true "Create a ShippingOrder"
// @Success 200 {object} ShippingOrderCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/create [post]
func (h *handler) CreateShippingOrderEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("shipping_orders"),
	}
	eda.Produce(eda.CREATE_SHIPPING_ORDER, request_payload, false)
	return res.RespSuccess(c, "Shipping_order Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})

}

func (h *handler) CreateShippingOrder(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(shipping.ShippingOrder)
	helpers.JsonMarshaller(data, createPayload)
	// t:=&final_data
	err := h.service.CreateShippingOrder(edaMetaData, createPayload)
	id := createPayload.ID
	edaMetaData.Query = map[string]interface{}{
		"id": id,
	}
	response_data_interface, err := h.service.GetShippingOrder(edaMetaData)
	response_data := response_data_interface.(shipping.ShippingOrder)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SHIPPING_ORDER_ACK, responseMessage)
		return
	}
	// t:=*createPayload
	fmt.Println("create payload data : ------>", *createPayload)
	// helpers.PrettyPrint("create payload", *createPayload)
	host := edaMetaData.Host
	scheme := edaMetaData.Scheme
	shipping_partner := response_data.ShippingPartner.PartnerName
	service_type := response_data.ShippingMode.DisplayName
	fmt.Println("host:"+host, "scheme:"+scheme, id, "shipping partner"+shipping_partner)
	fmt.Println("service type:" + service_type)

	emitter := emitter.GetObj()
	if service_type == "Land" {
		fmt.Println("Entered land " + events.CreateShipmentEvent + shipping_partner + service_type)
		emitter.Emit(events.CreateShipmentEvent+shipping_partner+service_type, host, scheme, id)
	}
	if service_type == "Air" {
		fmt.Println("Entered air " + events.CreateShipmentEvent + shipping_partner + service_type)
		emitter.Emit(events.CreateShipmentEvent+shipping_partner+service_type, host, scheme, id)
	}
	// eda.Produce(eda.CREATE_SHIPPING_ORDER_ACK, responseMessage)
	// cache implementation
	UpdateShippingOrderInCache(edaMetaData)

}

// GetAllShippingOrders godoc
// @Summary GetAllShippingOrders
// @Description Get All ShippingOrders
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllShippingOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders [get]
func (h *handler) GetAllShippingOrders(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []ShippingOrderResponse
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetShippingOrderFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllShippingOrders(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "ShippingOrders Fetched Successfully", response, p)
}

// GetAllShippingOrdersDropDown godoc
// @Summary GetAllShippingOrdersDropDown
// @Description Get All ShippingOrdersDropDown
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param per_page 				query int 	 false "per_page"
// @Param page_no           	query int 	 false "page_no"
// @Param sort					query string false "sort"
// @param filters 			    query string false "Filters"
// @Success 200 {object} GetAllShippingOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/dropdown [get]
func (h *handler) GetAllShippingOrdersDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []ShippingOrderResponse
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))

	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetShippingOrderFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllShippingOrders(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "ShippingOrders Fetched Successfully", result, p)
}

// GetShippingOrder godoc
// @Summary GetShippingOrder
// @Description Get a ShippingOrder
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id  path int true "Get ShippingOrder"
// @Success 200 {object} ShippingOrderGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id} [get]
func (h *handler) GetShippingOrder(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	shippingOrderId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": shippingOrderId,
		// "company_id": metaData.CompanyId,
	}

	result, err := h.service.GetShippingOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response ShippingOrderResponse
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "ShippingOrder Fetched Successfully", result)
}

// DeleteShippingOrder godoc
// @Summary DeleteShippingOrder
// @Description Delete a ShippingOrder
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Success 200 {object} DeleteShippingOrderResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id}/delete [delete]
func (h *handler) DeleteShippingOrder(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	shippingOrderId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         shippingOrderId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteShippingOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateShippingOrderInCache(metaData)

	return res.RespSuccess(c, "ShippingOrder Deleted Successfully", map[string]string{"Deleted shipping order id": shippingOrderId})
}

// UpdateShippingOrder godoc
// @Summary UpdateShippingOrder
// @Description Update a ShippingOrder
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param ShippingOrderUpdateRequest body ShippingOrderRequest true "Update a ShippingOrder"
// @Success 200 {object} UpdateShippingOrderResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id}/update [post]
func (h *handler) UpdateShippingOrderEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	shippingOrderId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         shippingOrderId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("shipping_orders"),
	}

	eda.Produce(eda.UPDATE_SHIPPING_ORDER, request_payload)
	return res.RespSuccess(c, "Basic Shipping_order Update Inprogress", map[string]interface{}{"request_id": shippingOrderId})

}

func (h *handler) UpdateShippingOrder(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	updatePayload := new(shipping.ShippingOrder)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.UpdateShippingOrder(edaMetaData, updatePayload)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_UOM_CLASS_ACK, responseMessage)
		return
	}
	host := edaMetaData.Host
	scheme := edaMetaData.Scheme

	h.service.GetShippingOrder(edaMetaData)
	status := updatePayload.ShippingStatus.DisplayName
	partner := updatePayload.ShippingPartner.PartnerName
	shipping_id := updatePayload.AwbNumber
	if status == "Delivery unsuccessful" {
		emitter := emitter.GetObj()
		emitter.Emit(events.CancelShipmentEvent+partner, host, scheme, shipping_id)
	}

	// cache implementation
	UpdateShippingOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_UOM_CLASS_ACK, responseMessage)

}

// DeleteShippingOrderLine godoc
// @Summary DeleteShippingOrderLine
// @Description Delete a ShippingOrderLine
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path int true "id"
// @Param product_id query string true "Product_id"
// @Success 200 {object} DeleteShippingOrderLineResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/order_line/{id}/delete [delete]
func (h *handler) DeleteShippingOrderLine(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	shippingOrderId := c.Param("id")
	productId := c.Param("product_id")
	metaData.Query = map[string]interface{}{
		"id":         shippingOrderId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
		"product_id": productId,
	}

	err = h.service.DeleteShippingOrderLine(metaData)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "ShippingOrder LineItem Deleted Successfully", shippingOrderId)
}

// BulkCreateShippingOrder godoc
// @Summary 	BulkCreateShippingOrder
// @Description Create a Multiple ShippingOrder
// @Tags 		ShippingOrder
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkShippingOrderCreateRequest body ShippingOrderRequest true "Create a Bulk ShippingOrder"
// @Success 	200 {object} BulkShippingOrderCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/shipping_orders/bulk_create [post]
func (h *handler) CreateBulkShippingOrder(c echo.Context) (err error) {
	data := new([]shipping.ShippingOrder)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	var metaData core.MetaData
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.CreateBulkShippingOrder(metaData, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, " Bulk ShippingOrders Created Successfully", map[string]bool{"All ShippingOrders Created Successfully": true})
}

// SendEmailShippingOrder godoc
// @Summary Email ShippingOrder
// @Description Email ShippingOrder
// @Tags ShippingOrder
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ShippingOrder_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id}/send_email [post]
func (h *handler) SendEmailShippingOrder(c echo.Context) error {
	return res.RespSuccess(c, "Email sent successfully", nil)
}

// DownloadShippingOrderPDF godoc
// @Summary Download ShippingOrder
// @Description Download ShippingOrder
// @Tags ShippingOrder
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ShippingOrder_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id}/download_pdf [post]
func (h *handler) DownloadShippingOrderPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

// GenerateShippingOrderPDF godoc
// @Summary Generate ShippingOrder PDF
// @Description Generate ShippingOrder PDF
// @Tags ShippingOrder
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ShippingOrder_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id}/generate_pdf [post]
func (h *handler) GenerateShippingOrderPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF Generated", resp_data)
}

// FavouriteShippingOrder godoc
// @Summary FavouriteShippingOrder
// @Description FavouriteShippingOrder
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path string true "ShippingOrder ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id}/favourite [post]
func (h *handler) FavouriteShippingOrder(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	shippingOrderId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         shippingOrderId,
		"company_id": metaData.CompanyId,
	}
	_, er := h.service.GetShippingOrder(metaData)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteShippingOrder(metaData.Query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "ShippingOrder is Marked as Favourite", map[string]string{"record_id": shippingOrderId})
}

// UnFavouriteShippingOrder godoc
// @Summary UnFavouriteShippingOrder
// @Description UnFavouriteShippingOrder
// @Tags ShippingOrder
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path string true "ShippingOrder ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_orders/{id}/unfavourite [post]
func (h *handler) UnFavouriteShippingOrder(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteShippingOrder(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "ShippingOrder is marked as Unfavourite", map[string]string{"record_id": id})
}

// GetShippingOrderTab godoc
// @Summary GetShippingOrderTab
// @Summary GetShippingOrderTab
// @Description GetShippingOrderTab
// @Tags ShippingOrder
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
// @Router /api/v1/shipping_orders/{id}/filter_module/{tab} [get]
func (h *handler) GetShippingOrderTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetShippingOrderTab(id, tab, page, access_template_id, token_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

func (h *handler) GetAllShippingOrdersAdmin(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	channelName := c.QueryParam("type")
	if channelName == "ONDC" {
		query := map[string]interface{}{
			"name": channelName,
		}
		res, _ := channel.NewService().GetChannelService(query)
		fmt.Println("===========>>>>", res.ID)
		if res.ID != 0 {
			helpers.AddMandatoryFilters(p, "channel_id", "=", res.ID)
		}
	}

	var response []ShippingOrderResponseDTOForBPPAdmin

	result, err := h.service.GetAllShippingOrders(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	for i := 0; i < len(response); i++ {
		for j := 0; j < len(response[i].ShippingOrderLines); j++ {
			response[i].ShippingOrderLine.ProductName = response[i].ShippingOrderLines[j].ProductVariant.ProductName
			response[i].ShippingOrderLine.SkuId = response[i].ShippingOrderLines[j].ProductVariant.SkuId
		}
	}
	return res.RespSuccessInfo(c, "ShippingOrders Fetched Successfully", response, p)
}

func (h *handler) GetShippingOrderAdmin(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	shippingOrderId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": shippingOrderId,
	}

	result, err := h.service.GetShippingOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response ShippingOrderResponse
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "ShippingOrder Fetched Successfully", result)
}
