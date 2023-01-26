package shipping_orders

import (
	"fmt"
	"strconv"

	shipping_base "fermion/backend_core/controllers/shipping/base"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"

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

func NewHandler() *handler {
	service := NewService()
	base_service := shipping_base.NewServiceBase()
	return &handler{service, base_service}
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
func (h *handler) CreateShippingOrder(c echo.Context) (err error) {
	data := c.Get("shipping_orders").(*ShippingOrderRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	response, err := h.service.CreateShippingOrder(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	id := response.ID
	shipping_partner := response.ShippingPartner.PartnerName
	service_type := response.ShippingMode.DisplayName
	emitter := emitter.GetObj()
	if service_type == "Land" {
		emitter.Emit(events.CreateShipmentEvent+shipping_partner+service_type, c, id)
	}
	if service_type == "Air" {
		emitter.Emit(events.CreateShipmentEvent+shipping_partner+service_type, c, id)
	}
	return res.RespSuccess(c, "ShippingOrder Created Successfully", map[string]interface{}{"Created ShippingOrder Id": response})
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
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllShippingOrders(page, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "ShippingOrders Fetched Successfully", result, page)
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
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllShippingOrders(page, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "dropdown of ShippingOrders Fetched Successfully", result, page)
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
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetShippingOrder(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
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
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteShippingOrder(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "ShippingOrder Deleted Successfully", map[string]string{"Deleted shipping order id": get_id})
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
func (h *handler) UpdateShippingOrder(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	data := c.Get("shipping_orders").(*ShippingOrderRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateShippingOrder(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	response, _ := h.service.GetShippingOrder(uint(id), token_id, access_template_id)
	status := response.ShippingStatus.Name
	partner := response.ShippingPartner.PartnerName
	shipping_id := response.AwbNumber
	if status == "Delivery unsuccessful" {
		emitter := emitter.GetObj()
		emitter.Emit(events.CancelShipmentEvent+partner, c, shipping_id)
	}
	return res.RespSuccess(c, "ShippingOrder Updated Successfully", map[string]int{"Updated id": id})
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
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	get_prod_id := c.QueryParam("product_id")
	product_id, _ := strconv.Atoi(get_prod_id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	query := map[string]interface{}{
		"shipping_order_id": uint(id),
		"product_id":        uint(product_id),
	}
	err = h.service.DeleteShippingOrderLine(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "ShippingOrder LineItem Deleted Successfully", query)
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
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.CreateBulkShippingOrder(data, token_id, access_template_id)
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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.GetShippingOrder(uint(ID), token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteShippingOrder(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "ShippingOrder is Marked as Favourite", map[string]string{"record_id": id})
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
