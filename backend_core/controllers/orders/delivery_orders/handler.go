package delivery_orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/sales_invoice"
	concurrency_management "fermion/backend_core/controllers/concurrency_management"
	"fermion/backend_core/controllers/eda"
	orders_base "fermion/backend_core/controllers/orders/base"
	"fermion/backend_core/internal/model/core"
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
	service            Service
	concurrencyService concurrency_management.Service
	base_service       orders_base.ServiceBase
}

var DeliveryOrdersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if DeliveryOrdersHandler != nil {
		return DeliveryOrdersHandler
	}
	service := NewService()
	concurrencyService := concurrency_management.NewService()
	base_service := orders_base.NewServiceBase()
	DeliveryOrdersHandler = &handler{service, concurrencyService, base_service}
	return DeliveryOrdersHandler
}

// CreateDeliveryOrders godoc
// @Summary      do a CreateDeliveryOrder
// @Description  Create a DeliveryOrder
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   DeliveryOrderCreateRequest body DeliveryOrderRequest true "create a DeliveryOrder"
// @Success      200  {object}  DeliveryOrderCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/create [post]
func (h *handler) CreateDeliveryOrdersEvent(c echo.Context) error {
	data := new(DeliveryOrderRequest)
	edaMetaData := c.Get("MetaData").(core.MetaData)

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("delivery_orders"),
	}
	//--------------------------concurrency_management_start-----------------------------------------
	var source_documents map[string]interface{}
	jsondata, err := json.Marshal(data.DeliveryOrdersDetails.SourceDocuments)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(jsondata, &source_documents)
	if err != nil {
		return res.RespErr(c, err)
	}

	if source_documents["id"] != nil && source_documents["id"].(float64) > 0 {
		SourceDocumentType, err := helpers.GetLookupcodeName(int(*data.DeliveryOrdersDetails.Source_document_id))
		if err != nil {
			return res.RespError(c, err)
		}
		recordID := uint(source_documents["id"].(float64))
		modelName := helpers.SourceDocumentTypeWithModelName[SourceDocumentType]
		columnName := helpers.SourceDocumentTypeWithColumnName[SourceDocumentType]

		response, PreviousStatusId, _ := h.concurrencyService.CheckConcurrencyStatus(recordID, columnName, modelName)
		if response.Block {
			return res.RespErr(c, errors.New(response.Message))
		}
		defer h.concurrencyService.ReleaseConcurrencyLock(recordID, PreviousStatusId, columnName, modelName)
	}
	//--------------------------concurrency_management_end-----------------------------------------
	seq_delivery := helpers.GenerateSequence("DO", fmt.Sprint(edaMetaData.TokenUserId), "delivery_orders")
	seq_ref := helpers.GenerateSequence("REF", fmt.Sprint(edaMetaData.TokenUserId), "delivery_orders")

	if data.DeliveryOrdersDetails.AutoCreateDoNumber {
		data.DeliveryOrdersDetails.Delivery_order_number = seq_delivery
	}
	if data.DeliveryOrdersDetails.AutoGenerateReferenceNumber {
		data.DeliveryOrdersDetails.Reference_id = seq_ref
	}
	eda.Produce(eda.CREATE_DELIVERY_ORDERS, requestPayload)
	return res.RespSuccess(c, "Delivery Order Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateDeliveryOrders(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(DeliveryOrderRequest)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.CreateDeliveryOrder(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_DELIVERY_ORDERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateDeliverOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": requestPayload.ID,
	}
	// eda.Produce(eda.CREATE_DELIVERY_ORDERS_ACK, *responseMessage)

}

// UpdateDeliveryOrders godoc
// @Summary      do a UpdateDeliveryOrders
// @Description  Update a DeliveryOrder
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   DeliveryOrderUpdateRequest body DeliveryOrderRequest true "Update a DeliveryOrder"
// @Success      200  {object}  DeliveryOrderUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id}/update [post]
func (h *handler) UpdateDeliveryOrdersEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("delivery_orders"),
	}

	eda.Produce(eda.UPDATE_DELIVERY_ORDERS, requestPayload)
	return res.RespSuccess(c, "Delivery Order Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateDeliveryOrders(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(DeliveryOrderRequest)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.UpdateDeliveryOrder(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_DELIVERY_ORDERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateDeliverOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_DELIVERY_ORDERS_ACK, *responseMessage)
}

// FindDeliveryOrders godoc
// @Summary      Find a DeliveryOrder by id
// @Description  Find a DeliveryOrder by id
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  DeliveryOrderGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id} [get]
func (h *handler) FindDeliveryOrders(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	doId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": doId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.FindDeliveryOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response DeliveryOrderRequest
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Delivery Order Details retrieved succesfully", response)
}

// AllDeliveryOrders godoc
// @Summary      get AllDeliveryOrders and filter it by search query
// @Description  a get AllDeliveryOrders and filter it by search query
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters           query     string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  DeliveryOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders [get]
func (h *handler) AllDeliveryOrders(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetDeliverOrderFromCache(fmt.Sprint(metaData.TokenUserId))
	}

	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
	}

	data, err := h.service.AllDeliveryOrders(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []DeliveryOrderResponse
	helpers.JsonMarshaller(data, &dtoData)
	return res.RespSuccessInfo(c, "list of delivery orders fetched", data, page)
}

// AllDeliveryOrdersDropDown godoc
// @Summary      get AllDeliveryOrdersDropDown and filter it by search query
// @Description  a get AllDeliveryOrdersDropDown and filter it by search query
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters           query     string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  DeliveryOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/dropdown [get]
func (h *handler) AllDeliveryOrdersDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetDeliverOrderFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.AllDeliveryOrders(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []DeliveryOrderResponse
	helpers.JsonMarshaller(data, &dtoData)

	return res.RespSuccessInfo(c, "dropdown list of delivery orders fetched", dtoData, page)
}

// DeleteDeliveryOrders godoc
// @Summary      delete a DeleteDeliveryOrders
// @Description  delete a DeleteDeliveryOrders
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  DeliveryOrderDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id}/delete [delete]
func (h *handler) DeleteDeliveryOrders(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteDeliveryOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateDeliverOrderInCache(metaData)
	return res.RespSuccess(c, "delivery order deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteDeliveryOrderLines godoc
// @Summary      delete a DeleteDeliveryOrderLines
// @Description  delete a DeleteDeliveryOrderLines
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  DeliveryOrderLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id}/delete_products [delete]
func (h *handler) DeleteDeliveryOrderLines(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"do_id":      id,
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeleteDeliveryOrderLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", metaData.Query)
}

// BulkCreateDeliveryOrders godoc
// @Summary 	do a BulkCreateDeliveryOrder
// @Description do a BulkCreateDeliveryOrder
// @Tags 		DeliveryOrders
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkDeliveryOrderCreateRequest body DeliveryOrderRequest true "Create a Bulk DeliveryOrder"
// @Success 	200 {object} BulkDeliveryOrderCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/delivery_orders/bulkcreate [post]
func (h *handler) BulkCreateDeliveryOrders(c echo.Context) (err error) {
	data := new([]DeliveryOrderRequest)
	c.Bind(data)
	s := c.Get("TokenUserID").(string)
	t := helpers.ConvertStringToUint(s)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.BulkCreateDeliveryOrder(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DO created successfully", map[string]interface{}{"created_data": data})
}

// SendEmailDeliveryOrders godoc
// @Summary Email DeliveryOrders
// @Description Email DeliveryOrders
// @Tags DeliveryOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/send_email [post]
func (h *handler) SendEmailDeliveryOrders(c echo.Context) error {
	return res.RespSuccess(c, "Email sent successfully", nil)
}

// DownloadDeliveryOrdersPDF godoc
// @Summary Download DeliveryOrders
// @Description Download DeliveryOrders
// @Tags DeliveryOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/download_pdf [post]
func (h *handler) DownloadDeliveryOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

// GenerateDeliveryOrdersPDF godoc
// @Summary Generate DeliveryOrders PDF
// @Description Generate DeliveryOrders PDF
// @Tags DeliveryOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/generate_pdf [post]
func (h *handler) GenerateDeliveryOrdersPDF(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	deliveryOrderId := c.Param("id")
	heading := c.QueryParam("heading")
	helpers.PrettyPrint("=====================>", heading)
	helpers.PrettyPrint("=====================>", metaData)
	helpers.PrettyPrint("=====================>", deliveryOrderId)
	// company_id := c.Get("CompanyId").(string)
	x := sales_invoice.NewHandler()
	er := x.GetSalesInvoicePdf(c)
	if er != nil {
		return res.RespSuccess(c, "Generated Successfully", "")
	}
	return c.Attachment("test.pdf", "test")
}

// FavouriteDeliveryOrders godoc
// @Summary FavouriteDeliveryOrders
// @Description FavouriteDeliveryOrders
// @Tags DeliveryOrders
// @Accept  json
// @Produce  json
// @param id path string true "DeliveryOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/favourite [post]
func (h *handler) FavouriteDeliveryOrders(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	_, err = h.service.FindDeliveryOrder(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	err = h.base_service.FavouriteDeliveryOrders(metaData.Query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DeliveryOrder is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouriteDeliveryOrders godoc
// @Summary UnFavouriteDeliveryOrders
// @Description UnFavouriteDeliveryOrders
// @Tags DeliveryOrders
// @Accept  json
// @Produce  json
// @param id path string true "DeliveryOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/unfavourite [post]
func (h *handler) UnFavouriteDeliveryOrders(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouriteDeliveryOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DeliveryOrder is Unmarked as Favourite", map[string]string{"id": id})
}

//-----------------------------------------------------

func (h *handler) GetDeliveryOrderTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id":  c.Param("id"),
		"tab": c.Param("tab"),
	}

	data, err := h.service.GetDeliveryOrderTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "success", data, page)
}
