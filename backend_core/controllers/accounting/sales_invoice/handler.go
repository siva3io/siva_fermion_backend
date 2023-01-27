package sales_invoice

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"os"
	"strconv"
	"time"

	accounting_base "fermion/backend_core/controllers/accounting/base"
	"fermion/backend_core/controllers/cores"
	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/controllers/mdm/locations"
	"fermion/backend_core/internal/model/accounting"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/labstack/echo/v4"
)

const (
	AWS_S3_REGION = "ap-southeast-1"
	AWS_S3_BUCKET = "dev-api-files"
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
	service           Service
	base_service      accounting_base.ServiceBase
	cores_service     cores.Service
	locations_service locations.Service
}

var SalesInvoiceHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if SalesInvoiceHandler != nil {
		return SalesInvoiceHandler
	}
	service := NewService()
	base_service := accounting_base.NewServiceBase()
	cores_service := cores.NewService()
	locations_service := locations.NewService()
	SalesInvoiceHandler = &handler{service, base_service, cores_service, locations_service}
	return SalesInvoiceHandler
}

// CreateSalesInvoice godoc
// @Summary      Create Sales Invoice
// @Description  Create Sales Invoice
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   SalesInvoiceCreateRequest body SalesInvoiceRequest true "create Sales Invoice"
// @Success      200  {object}  SalesInvoiceCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/create [post]
func (h *handler) CreateSalesInvoiceEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("sales_invoice"),
	}
	eda.Produce(eda.CREATE_SALES_INVOICE, request_payload)
	return res.RespSuccess(c, "Sales Invoice Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateSalesInvoice(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	createPayLoad := new(accounting.SalesInvoice)
	helpers.JsonMarshaller(data, createPayLoad)
	createPayLoad.CreatedByID = &edaMetaData.TokenUserId
	err := h.service.CreateSalesInvoice(edaMetaData, createPayLoad)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SALES_INVOICE_ACK, responseMessage)

		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayLoad.ID,
	}
	// eda.Produce(eda.CREATE_SALES_INVOICE_ACK, responseMessage)
	// emitter := emitter.GetObj()
	// fmt.Println("ID---------------------->", createPayLoad.ID)
	// emitter.Emit(events.CreateSalesInvoiceEvent+"Quickbooks", createPayLoad.ID)

	// cache implementation
	UpdateSalesInvoiceInCache(edaMetaData)
}

// func (h *handler) CreateSalesInvoice(c echo.Context) (err error) {
// 	data := c.Get("sales_invoice").(*SalesInvoiceRequest)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	user_id, _ := strconv.Atoi(token_id)
// 	t := uint(user_id)
// 	data.CreatedByID = &t
// 	id, err := h.service.CreateSalesInvoice(data, token_id, access_template_id)
// 	if err != nil {
// 		return res.RespError(c, err)
// 	}

// 	emitter := emitter.GetObj()
// 	fmt.Println("ID---------------------->", id)
// 	emitter.Emit(events.CreateSalesInvoiceEvent+"Quickbooks", c, id)

// 	// cache implementation
// 	UpdateSalesInvoiceInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Sales Invoice created successfully", map[string]interface{}{"created_id": id})
// }

// BulkCreateSalesInvoice godoc
// @Summary 	Bulk Create Sales Invoice
// @Description Bulk Create Sales Invoice
// @Tags 		Sales Invoice
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkSalesInvoiceCreateRequest body SalesInvoiceRequest true "Create a Bulk Sales Invoice"
// @Success 	200 {object} BulkSalesInvoiceCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/sales_invoice/bulk_create [post]
func (h *handler) BulkCreateSalesInvoice(c echo.Context) (err error) {
	data := new([]accounting.SalesInvoice)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	t := helpers.ConvertStringToUint(token_id)
	var metaData core.MetaData
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateSalesInvoice(metaData, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice multiple records created successfully", data)
}

// UpdateSalesInvoice godoc
// @Summary      Update Sales Invoice
// @Description  Update Sales Invoice
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   SalesInvoiceUpdateRequest body SalesInvoiceRequest true "Update Sales Invoice"
// @Success      200  {object}  SalesInvoiceUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/{id}/edit [put]
func (h *handler) UpdateSalesInvoiceEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	salesInvoiceId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         salesInvoiceId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("sales_invoice"),
	}
	eda.Produce(eda.UPDATE_SALES_INVOICE, request_payload)
	return res.RespSuccess(c, "Sales Invoice Update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateSalesInvoice(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(accounting.SalesInvoice)
	helpers.JsonMarshaller(data, updatePayLoad)
	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdateSalesInvoice(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err

		// eda.Produce(eda.UPDATE_SALES_INVOICE_ACK, responsemessage)
		return
	}
	// cache implementation
	UpdateSalesInvoiceInCache(edaMetaData)
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_SALES_INVOICE_ACK, responsemessage)
}

// func (h *handler) UpdateSalesInvoice(c echo.Context) (err error) {
// 	ID := c.Param("id")
// 	id, _ := strconv.Atoi(ID)
// 	data := c.Get("sales_invoice").(*SalesInvoiceRequest)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
// 	err = h.service.UpdateSalesInvoice(uint(id), data, token_id, access_template_id)
// 	if err != nil {
// 		return res.RespError(c, err)
// 	}

// 	// cache implementation
// 	UpdateSalesInvoiceInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Sales Invoice updated successfully", map[string]interface{}{"updated_id": id})
// }

// GetSalesInvoice godoc
// @Summary 	Get Sales Invoice by id
// @Description Get Sales Invoice by id
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  SalesInvoiceGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/{id} [get]
func (h *handler) GetSalesInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	purchaseInoviceId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         purchaseInoviceId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetSalesInvoice(metaData)
	fmt.Println(result)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Sales InvoiceDetails Retrieved successfully", result)
}

// GetAllSalesInvoice godoc
// @Summary      Get All Sales Invoice with MiscOptions of Pagination, Search, Filter and Sort
// @Description  Get All Sales Invoice with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        filters		query	  string false "filters"
// @Param        per_page 		query     int    false "per_page"
// @Param  	     page_no        query     int    false "page_no"
// @Param        sort			query     string false "sort"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  SalesInvoiceGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice [get]
func (h *handler) GetAllSalesInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var data []SalesInvoiceRequest
	var cacheResponse interface{}
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetSalesInvoiceFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, p)
	}
	resp, err := h.service.GetAllSalesInvoice(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	helpers.JsonMarshaller(resp, &data)
	return res.RespSuccessInfo(c, "list of sales invoice fetched", data, p)
}

// GetAllSalesInvoiceDropDown godoc
// @Summary      Get All Sales Invoice dropdown with MiscOptions of Pagination, Search, Filter and Sort
// @Description  Get All Sales Invoice dropdown with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        filters		query	  string false "filters"
// @Param        per_page 		query     int    false "per_page"
// @Param  	     page_no        query     int    false "page_no"
// @Param        sort			query     string false "sort"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  SalesInvoiceGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/dropdown [get]
func (h *handler) GetAllSalesInvoiceDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetSalesInvoiceFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, p)
	}
	var data []SalesInvoiceGetAll
	response, err := h.service.GetAllSalesInvoice(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	helpers.JsonMarshaller(response, &data)
	return res.RespSuccessInfo(c, "dropdown list of sales invoice fetched", data, p)
}

// DeleteSalesInvoice godoc
// @Summary      Delete Sales Invoice
// @Description  Delete Sales Invoice
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  SalesInvoiceDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/{id}/delete [delete]
func (h *handler) DeleteSalesInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	salesInvoiceId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         salesInvoiceId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteSalesInvoice(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// cache implementation
	UpdateSalesInvoiceInCache(metaData)

	return res.RespSuccess(c, "Sales Invoice deleted successfully", salesInvoiceId)
}

// DeleteSalesInvoiceLines godoc
// @Summary      Delete Sales Invoice Lines
// @Description  Delete Sales Invoice Lines
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        sales_invoice_id         path      int    true  "sales_invoice_id"
// @Param        product_id     query     string true "product_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  SalesInvoiceLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/{id}/delete_products [delete]
func (h *handler) DeleteSalesInvoiceLines(c echo.Context) (err error) {
	var product_id = c.QueryParam("product_id")
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = map[string]interface{}{
		"sales_invoice_id": ID,
		"product_id":       product_id,
	}
	var metaData core.MetaData
	err = h.service.DeleteSalesInvoiceLines(metaData)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "purchase invoice line deleted successfully", query)
}

// SendMailSalesInvoice godoc
// @Summary      Send Email to specific recevier for Id
// @Description  Send Email to specific recevier for Id
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        id         	query     int    true "id"
// @Param        receiver_email query     string true "receiver_email"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  SendMailSalesInvoiceResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/{id}/sendemail [get]
func (h *handler) SendMailSalesInvoice(c echo.Context) (err error) {

	q := new(SendMailSalesInvoice)
	c.Bind(q)
	// err = h.service.SendMailSalesInvoice(q)
	// if err != nil {
	// 	return res.RespError(c, err)
	// }

	return res.RespSuccess(c, "Email Sent Successfully", q)
}

// DownloadPdfSalesInvoice godoc
// @Summary      Download PDF
// @Description  Download PDF
// @Tags         Sales Invoice
// @Accept       json
// @Produce      json
// @Param        id         	query     int    true "id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  DownloadPdfSalesInvoiceResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/sales_invoice/{id}/printpdf [get]
func (h *handler) DownloadPdfSalesInvoice(c echo.Context) (err error) {
	return res.RespSuccess(c, "Pdf Generated Successfully", "https://drive.google.com/file/d/1O46B3lgliUGakL4weBjrovCHphTzR5Np/view?usp=sharing")
}

// FavouriteSalesInvoice godoc
// @Summary Favourite Sales Invoice
// @Description Favourite Sales Invoice
// @Tags Sales Invoice
// @Accept  json
// @Produce  json
// @param id path string true "Sales Invoice ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_invoice/{id}/favourite [post]
func (h *handler) FavouriteSalesInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	salesInvoiceId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": salesInvoiceId,
	}
	_, err = h.service.GetSalesInvoice(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	query := map[string]interface{}{
		"ID":      salesInvoiceId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.FavouriteSalesInvoice(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice is Marked as Favourite", map[string]string{"id": salesInvoiceId})
}

// UnFavouriteSalesInvoice godoc
// @Summary UnFavourite Sales Invoice
// @Description UnFavourite Sales Invoice
// @Tags Sales Invoice
// @Accept  json
// @Produce  json
// @param id path string true "Sales Invoice ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_invoice/{id}/unfavourite [post]
func (h *handler) UnFavouriteSalesInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	salesInvoiceId := c.Param("id")

	query := map[string]interface{}{
		"ID":      salesInvoiceId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.UnFavouriteSalesInvoice(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseInvoice is Marked as UnFavourite", map[string]string{"credinote_id": salesInvoiceId})
}

func (h *handler) GetSalesInvoiceTab(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")
	metaData.AdditionalFields = map[string]interface{}{
		"purchase_invoice": id,
		"tab_name":         tab,
	}

	data, err := h.service.GetSalesInvoiceTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

func (h *handler) GetSalesInvoicePdf(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	salesInvoiceId := c.Param("id")
	company_id := c.Get("CompanyId").(string)
	// fmt.Println("==========cid", company_id)
	query := map[string]interface{}{
		"id": salesInvoiceId,
	}

	heading := c.QueryParam("heading")

	switch heading {
	case "sales_orders":
		heading = "SalesOrders"
	case "purchase_orders":
		heading = "PurchaseOrders"
	case "delivery_orders":
		heading = "DeliveryOrders"
	case "scrap_orders":
		heading = "ScrapOrders"
	default:
		heading = "Invoice"
	}

	result, err := h.service.GetSalesInvoicePdf(query)
	if err != nil {
		return res.RespErr(c, err)
	}

	billing_address := make([]map[string]interface{}, 0)
	shipping_address := make([]map[string]interface{}, 0)
	products := make([]map[string]interface{}, 0)
	json.Unmarshal(result.BillingAddress, &billing_address)
	json.Unmarshal(result.ShippingAddress, &shipping_address)
	helpers.JsonMarshaller(result.SalesInvoiceLines, &products)

	length := len(billing_address)
	if length < 0 {
		fmt.Println("------------BILLING ADDRESS EMPTY---------")
		return res.RespErr(c, errors.New("BILLING ADDRESS IS EMPTY"))
	}
	length = len(shipping_address)
	if length < 0 {
		fmt.Println("------------SHIPPING ADDRESS EMPTY---------")
		return res.RespErr(c, errors.New("SHIPPING ADDRESS IS EMPTY"))
	}
	exp_shipment_date, _ := result.ExpectedShipmentDate.Value()
	var shipment_data string
	byte_time, _ := json.Marshal(exp_shipment_date)
	json.Unmarshal(byte_time, &shipment_data)
	company_details, err := h.cores_service.GetCompanyPreferences(map[string]interface{}{
		"id": company_id,
	})
	if err != nil {
		return res.RespErr(c, err)
	}
	r := NewRequestPdf("")

	templatePath := "backend_core/views/index.html"
	outputPath := "test.pdf"
	lookup_id, err := helpers.GetLookupcodeId("LOCATION_TYPE", "OFFICE")
	if err != nil {
		return res.RespErr(c, err)
	}
	locationPagnation := new(pagination.Paginatevalue)
	locationPagnation.Filters = fmt.Sprintf("[[\"company_id\",\"=\",%v],[\"location_type_id\",\"=\",%v]]", company_id, lookup_id)
	location_resp, err := h.locations_service.GetLocationList(metaData, locationPagnation)
	if err != nil {
		fmt.Println("-----------------LOCATION_ERROR------", err.Error())
		return res.RespErr(c, err)
	}
	var location_data = make([]map[string]interface{}, 0)
	json_data, _ := json.Marshal(location_resp)
	json.Unmarshal(json_data, &location_data)
	length = len(location_data)
	if length < 0 {
		fmt.Println("------------LOCATION EMPTY---------")
		return res.RespErr(c, errors.New("LOCATION IS EMPTY"))
	}
	address_data := location_data[0]["address"].(map[string]interface{})

	// Company_Id := result.CompanyId
	// fmt.Println("company id",Company_Id)
	query1 := map[string]interface{}{
		"id": company_id,
	}
	gst_no, err := h.cores_service.GetCompanyPreferences(query1)
	gst := gst_no.GSTIN
	template_data := map[string]interface{}{
		"heading":              heading,
		"gst":                  gst,
		"bill_address_1":       billing_address[0]["address_line1"],
		"bill_address_2":       billing_address[0]["address_line2"],
		"bill_address_3":       billing_address[0]["address_line3"],
		"bill_address_code":    billing_address[0]["zipcode"],
		"bill_address_state":   billing_address[0]["state"],
		"bill_address_country": billing_address[0]["country"],
		"bill_address_phone":   billing_address[0]["mobile_number"],
		"ship_address_1":       shipping_address[0]["address_line1"],
		"ship_address_2":       shipping_address[0]["address_line2"],
		"ship_address_3":       shipping_address[0]["address_line3"],
		"ship_address_code":    shipping_address[0]["zipcode"],
		"ship_address_state":   shipping_address[0]["state"],
		"ship_address_country": shipping_address[0]["country"],
		"ship_address_phone":   shipping_address[0]["mobile_number"],
		"invoice_number":       result.SalesInvoiceNumber,
		"ship_date":            result.SalesInvoiceDate.Format(time.RFC3339Nano)[:10],
		"invoice_date":         shipment_data[:10],
		"order_number":         result.ReferenceNumber,
		"product_details":      products,
		"sub_amt":              result.SubTotalAmount,
		"tot_amt":              result.TotalAmount,
		"tot_tax":              result.TaxAmount,
		"igst":                 result.IgstAmt,
		"cgst":                 result.CgstAmt,
		"sgst":                 result.SgstAmt,
		"payment":              result.PaymentType.DisplayName,
		"company_name":         company_details.Name,
		"company_email":        company_details.Email,
		"company_website":      company_details.Website,
		"company_adrs_1":       address_data["address_line_1"],
		"company_adrs_2":       address_data["address_line_2"],
		"company_state":        address_data["state"].(map[string]interface{})["name"],
		"company_country":      address_data["country"].(map[string]interface{})["name"],
	}
	if err := r.ParseTemplate(templatePath, template_data); err == nil {
		ok, _ := r.GeneratePDF(outputPath)
		fmt.Println("pdf generated:", ok)
	} else {
		fmt.Println(err)
	}
	c.Response().Header().Set(
		"Content-Disposition",
		"attachment; filename="+outputPath,
	)

	// Upload Files
	var s3path = ""
	var path = ""
	switch heading {
	case "SalesOrders":
		path = "sales_orders/" + company_details.Name + "/" + salesInvoiceId + "/"
	case "PurchaseOrders":
		path = "purchase_orders/" + company_details.Name + "/" + salesInvoiceId + "/"
	case "DeliveryOrders":
		path = "delivery_orders/" + company_details.Name + "/" + salesInvoiceId + "/"
	case "ScrapOrders":
		path = "scrap_orders/" + company_details.Name + "/" + salesInvoiceId + "/"
	default:
		path = "sales_invoice/" + company_details.Name + "/" + salesInvoiceId + "/"
	}

	err = helpers.UploadFile("./test.pdf", path)
	if err != nil {
		log.Fatal(err)
	}

	s3path = "https://dev-api-files.s3.amazonaws.com/" + path + "test.pdf"
	defer os.Remove("./test.pdf")

	return c.Attachment("test.pdf", s3path)

}
