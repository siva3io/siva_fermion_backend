package scrap_orders

import (
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/sales_invoice"
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
	service      Service
	base_service orders_base.ServiceBase
}

var ScrapOrdersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ScrapOrdersHandler != nil {
		return ScrapOrdersHandler
	}
	service := NewService()
	base_service := orders_base.NewServiceBase()
	ScrapOrdersHandler = &handler{service, base_service}
	return ScrapOrdersHandler
}

// CreateScrapOrder godoc
// @Summary      do a CreateScrapOrder
// @Description  Create a ScrapOrder
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param CreateScrapOrderRequest body  ScrapOrders true "create a ScrapOrder"
// @Success      200  {object}  ScrapOrderCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/create [post]
func (h *handler) CreateScrapOrderEvent(c echo.Context) (err error) {
	data := new(ScrapOrders)
	edaMetaData := c.Get("MetaData").(core.MetaData)

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("scrap_orders"),
	}
	seq_scrap := helpers.GenerateSequence("SCRAP", fmt.Sprint(edaMetaData.TokenUserId), "scrap_orders")
	seq_ref := helpers.GenerateSequence("REF", fmt.Sprint(edaMetaData.TokenUserId), "scrap_orders")
	if data.AutoCreateScrapNumber {
		data.Scrap_order_number = seq_scrap
	}
	if data.AutoGenerateReferenceNumber {
		data.Reference_id = seq_ref
	}
	eda.Produce(eda.CREATE_SCRAP_ORDERS, requestPayload)
	return res.RespSuccess(c, "Scrap Order Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateScrapOrder(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(ScrapOrders)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.CreateScrapOrder(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SCRAP_ORDERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateScrapOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": requestPayload.ID,
	}
	// eda.Produce(eda.CREATE_SCRAP_ORDERS_ACK, *responseMessage)

}

// UpdateScrap godoc
// @Summary      do a UpdateScrap
// @Description  Update a ScrapOrder
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id                      path  int               true  "id"
// @Param UpdateScrapOrderRequest body  ScrapOrders true "update a ScrapOrder"
// @Success      200  {object}  ScrapOrderUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/{id}/update [post]
func (h *handler) UpdateScrapOrderEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("scrap_orders"),
	}

	eda.Produce(eda.UPDATE_SCRAP_ORDERS, requestPayload)
	return res.RespSuccess(c, "Scrap Order Update Inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateScrapOrder(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	updatePayload := new(ScrapOrders)
	helpers.JsonMarshaller(request["data"], updatePayload)

	err := h.service.UpdateScrapOrder(edaMetaData, updatePayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_SCRAP_ORDERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateScrapOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": updatePayload.ID,
	}
	// eda.Produce(eda.UPDATE_SCRAP_ORDERS_ACK, *responseMessage)

}

// FindAllScrapOrdersDropDown godoc
// @Summary      Find All ScrapOrders and filter it by search query
// @Description  Find All ScrapOrders and filter it by search query
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters 	query    string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  ScrapOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/dropdown [get]
func (h *handler) FindAllScrapOrdersDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetScrapOrderFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.AllScrapOrders(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ScrapOrdersResponseDTO
	helpers.JsonMarshaller(data, &dtoData)
	return res.RespSuccessInfo(c, "list of scrap orders fetched", dtoData, page)
}

// FindAllScrapOrders godoc
// @Summary      Find All ScrapOrders and filter it by search query
// @Description  Find All ScrapOrders and filter it by search query
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters 	query    string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  ScrapOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders [get]
func (h *handler) FindAllScrapOrders(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetScrapOrderFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.AllScrapOrders(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ScrapOrdersResponseDTO
	helpers.JsonMarshaller(data, &dtoData)
	return res.RespSuccessInfo(c, "list of scrap orders fetched", dtoData, page)
}

// ScrapOrdersByid godoc
// @Summary      Find a ScrapOrder By id
// @Description  Find a ScrapOrder By id
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id     path  int  true  "id"
// @Success      200  {object}  ScrapOrderGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/{id} [get]
func (h *handler) ScrapOrdersByid(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": c.Param("id"),
		// "company_id": metaData.CompanyId,
	}
	data, err := h.service.FindScrapOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "scrap orders fetched", data)
}

// DeleteScrap godoc
// @Summary      Delete a ScrapOrder
// @Description  Delete a ScrapOrder
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id     path  int  true  "id"
// @Success      200  {object}  ScrapOrderDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/{id}/delete [delete]
func (h *handler) DeleteScrapOrder(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteScrapOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateScrapOrderInCache(metaData)

	return res.RespSuccess(c, "scrap order deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteScrapOrders godoc
// @Summary      Delete a ScrapOrderLine and provide the product_id as queryParam
// @Description  Delete a ScrapOrderLine and provide the product_id as queryParam
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"

// @Param        line_id             path      int    true  "line_id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  ScrapOrderLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse

// @Router     /api/v1/scrap_orders/{line_id}/delete_product [delete]
func (h *handler) DeleteScrapOrderLines(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"scrap_id":   id,
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeleteScrapOrderLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", metaData.Query)
}

// BulkCreateScrapOrder godoc
// @Summary 	do a BulkCreateScrapOrder
// @Description do a BulkCreateScrapOrder
// @Tags 		ScrapOrders
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkScrapOrderCreateRequest body ScrapOrders true "Create a Bulk ScrapOrder"
// @Success 	200 {object} BulkScrapOrderCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/scrap_orders/bulkcreate [post]
func (h *handler) BulkCreateScrapOrder(c echo.Context) (err error) {
	data := new([]ScrapOrders)
	c.Bind(data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateScrapOrder(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrders order added successfully", map[string]interface{}{"created_data": data})
}

// SendEmailScrapOrders godoc
// @Summary Email ScrapOrders
// @Description Email ScrapOrders
// @Tags ScrapOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/send_email [post]
func (h *handler) SendEmailScrapOrders(c echo.Context) error {
	return res.RespSuccess(c, "Email sent successfully", nil)
}

// DownloadScrapOrdersPDF godoc
// @Summary Download ScrapOrders
// @Description Download ScrapOrders
// @Tags ScrapOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/download_pdf [post]
func (h *handler) DownloadScrapOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

// GenerateScrapOrdersPDF godoc
// @Summary Generate ScrapOrders PDF
// @Description Generate ScrapOrders PDF
// @Tags ScrapOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/generate_pdf [post]
func (h *handler) GenerateScrapOrdersPDF(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	purchaseOrderId := c.Param("id")
	heading := c.QueryParam("heading")
	helpers.PrettyPrint("=====================>", heading)
	helpers.PrettyPrint("=====================>", metaData)
	helpers.PrettyPrint("=====================>", purchaseOrderId)
	// company_id := c.Get("CompanyId").(string)
	x := sales_invoice.NewHandler()
	er := x.GetSalesInvoicePdf(c)
	if er != nil {
		return res.RespSuccess(c, "Generated Successfully", "")
	}
	return c.Attachment("test.pdf", "test")
}

// FavouriteScrapOrders godoc
// @Summary FavouriteScrapOrders
// @Description FavouriteScrapOrders
// @Tags ScrapOrders
// @Accept  json
// @Produce  json
// @param id path string true "ScrapOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/favourite [post]
func (h *handler) FavouriteScrapOrder(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	_, err = h.service.FindScrapOrder(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.FavouriteScrapOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrder is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouriteScrapOrders godoc
// @Summary UnFavouriteScrapOrders
// @Description UnFavouriteScrapOrders
// @Tags ScrapOrders
// @Accept  json
// @Produce  json
// @param id path string true "ScrapOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/unfavourite [post]
func (h *handler) UnFavouriteScrapOrder(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouriteScrapOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrder is Unmarked as Favourite", map[string]string{"id": id})
}

// GetScrapOrderTab godoc
// @Summary GetScrapOrderTab
// @Summary GetScrapOrderTab
// @Description GetScrapOrderTab
// @Tags ScrapOrders
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
// @Router /api/v1/scrap_orders/{id}/filter_module/{tab} [get]
func (h *handler) GetScrapOrderTab(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	page := new(pagination.Paginatevalue)
	err := c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	metaData.AdditionalFields = map[string]interface{}{
		"id":  c.Param("id"),
		"tab": c.Param("tab"),
	}

	data, err := h.service.GetScrapOrderTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
