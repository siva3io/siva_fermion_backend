package purchase_invoice

import (
	"fmt"
	"strconv"

	accounting_base "fermion/backend_core/controllers/accounting/base"
	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"

	// "fermion/backend_core/pkg/util/emitter"
	// "fermion/backend_core/pkg/util/events"
	"fermion/backend_core/controllers/eda"
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
	base_service accounting_base.ServiceBase
}

var PurchaseInvoiceHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if PurchaseInvoiceHandler != nil {
		return PurchaseInvoiceHandler
	}
	service := NewService()
	base_service := accounting_base.NewServiceBase()
	PurchaseInvoiceHandler = &handler{service, base_service}
	return PurchaseInvoiceHandler
}

// CreatePurchaseInvoice godoc
// @Summary Create PurchaseInvoice
// @Description Create PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept  json
// @Produce  json
// @param RequestBody body PurchaseInvoiceDTO true "Purchase Invoice Request Body"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} PurchaseInvoiceCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/create [post]
func (h *handler) CreatePurchaseInvoiceEvent(c echo.Context) error {

	edaMetaData := c.Get("MetaData").(core.MetaData)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("purchase_invoice"),
	}
	eda.Produce(eda.CREATE_PURCHASE_INVOICE, request_payload)
	return res.RespSuccess(c, "Purchase Invoice Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreatePurchaseInvoice(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	createPayLoad := new(accounting.PurchaseInvoice)
	helpers.JsonMarshaller(data, createPayLoad)
	createPayLoad.CreatedByID = &edaMetaData.TokenUserId
	err := h.service.CreatePurchaseInvoice(edaMetaData, createPayLoad)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_PURCHASE_INVOICE_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayLoad.ID,
	}
	// eda.Produce(eda.CREATE_PURCHASE_INVOICE_ACK, responseMessage)
	fmt.Println("-------cur_id----------", createPayLoad.CurrencyId)

	// emitter := emitter.GetObj()
	// fmt.Println("ID---------------------->", createPayLoad.ID)
	// emitter.Emit(events.CreatePurchaseInvoiceEvent+"Quickbooks", createPayLoad.ID)
	// cache implementation
	UpdatePurchaseInvoiceInCache(edaMetaData)

}

// GetListPurchaseInvoice godoc
// @Summary Get all PurchaseInvoice list
// @Description Get all PurchaseInvoice list
// @Tags PurchaseInvoice
// @Accept json
// @Produce json
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {array} GetAllPurchaseInvoiceResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice [get]
func (h *handler) ListPurchaseInvoice(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var data []ListPurchaseInvoiceDTO
	var cacheResponse interface{}
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetPurchaseInvoiceFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &data)
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, p)
	}
	resp, err := h.service.ListPurchaseInvoice(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	helpers.JsonMarshaller(resp, &data)
	return res.RespSuccessInfo(c, "list of purchase invoice fetched", data, p)
}

// ListPurchaseInvoiceDropDown godoc
// @Summary Get all PurchaseInvoice dropdown list
// @Description Get all PurchaseInvoice dropdown list
// @Tags PurchaseInvoice
// @Accept json
// @Produce json
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {array} GetAllPurchaseInvoiceResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/dropdown [get]
func (h *handler) ListPurchaseInvoiceDropDown(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var data []ListPurchaseInvoiceDTO
	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetPurchaseInvoiceFromCache(tokenUserId)
	}
	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &data)
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, p)
	}

	response, err := h.service.ListPurchaseInvoice(metaData, p)
	if err != nil {
		return res.RespError(c, err)
	}
	helpers.JsonMarshaller(response, &data)
	return res.RespSuccessInfo(c, "dropdown list of purchase invoice fetched", data, p)
}

// ViewPurchaseInvoice godoc
// @Summary View PurchaseInvoice
// @Summary View PurchaseInvoice
// @Description View PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept  json
// @Produce  json
// @param id path string true "purchase_invoice_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} PurchaseInvoiceGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id} [get]
func (h *handler) ViewPurchaseInvoice(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	purchaseInoviceId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         purchaseInoviceId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.ViewPurchaseInvoice(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// var response PurchaseInvoiceDTO
	// helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Credit Note Details Retrieved successfully", result)
}

// UpdatePurchaseInvoice godoc
// @Summary Update PurchaseInvoice
// @Description Update PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept  json
// @Produce  json
// @param id path string true "purchase_invoice_id"
// @param RequestBody body PurchaseInvoiceDTO true "Purchase Invoice Request Body"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} PurchaseInvoiceUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/update [post]
func (h *handler) UpdatePurchaseInvoiceEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	purchaseInoviceId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         purchaseInoviceId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("purchase_invoice"),
	}
	eda.Produce(eda.UPDATE_PURCHASE_INVOICE, request_payload)
	return res.RespSuccess(c, "Purchase Invoice Update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdatePurchaseInvoice(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(accounting.PurchaseInvoice)
	helpers.JsonMarshaller(data, updatePayLoad)
	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdatePurchaseInvoice(edaMetaData, *updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_PURCHASE_INVOICE_ACK, responsemessage)
		return
	}
	// cache implementation
	UpdatePurchaseInvoiceInCache(edaMetaData)
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_PURCHASE_INVOICE_ACK, responsemessage)

}

// DeletePurchaseInvoice godoc
// @Summary Delete PurchaseInvoice
// @Description Delete PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept  json
// @Produce  json
// @param id path string true "purchase_invoice_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} PurchaseInvoiceDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/delete [delete]
func (h *handler) DeletePurchaseInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	purchaseInvoiceId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         purchaseInvoiceId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeletePurchaseInvoice(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	// cache implementation
	UpdatePurchaseInvoiceInCache(metaData)
	return res.RespSuccess(c, "Record deleted successfully", purchaseInvoiceId)
}

// DeletePurchaseInvoiceLines godoc
// @Summary Delete PurchaseInvoiceLines
// @Description Delete PurchaseInvoiceLines
// @Tags PurchaseInvoice
// @Accept  json
// @Produce  json
// @param id path string true "purchase_invoice_id"
// @param product_id query string true "product_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} PurchaseInvoiceLinesDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/invoice_lines/{id}/delete [delete]
func (h *handler) DeletePurchaseInvoiceLines(c echo.Context) (err error) {
	var product_id = c.QueryParam("product_id")
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = map[string]interface{}{
		"purchase_invoice_id": ID,
		"product_id":          product_id,
	}
	var metaData core.MetaData
	err = h.service.DeletePurchaseInvoiceLines(metaData)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "purchase invoice line deleted successfully", query)
}

// EmailPurchaseInvoice godoc
// @Summary Email PurchaseInvoice
// @Description Email PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept json
// @Produce json
// @param id path string true "purchase_invoice_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/sendEmail [post]
func (h *handler) EmailPurchaseInvoice(c echo.Context) error {
	return res.RespSuccess(c, "email sent successfully", nil)
}

// DownloadPurchaseInvoice godoc
// @Summary Download PurchaseInvoice
// @Description Download PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept json
// @Produce json
// @param id path string true "purchase_invoice_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/downloadPdf [post]
func (h *handler) DownloadPurchaseInvoice(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file downloaded successfully", resp_data)
}

// PrintPurchaseInvoice godoc
// @Summary Print PurchaseInvoice
// @Description print PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept json
// @Produce json
// @param id path string true "purchase_invoice_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/printPdf [post]
func (h *handler) PrintPurchaseInvoice(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file printed successfully", resp_data)
}

// GeneratePurchaseInvoicePDF godoc
// @Summary Generate PDF PurchaseInvoice
// @Description Generate PDF PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept json
// @Produce json
// @param id path string true "purchase_invoice_id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/generatePdf [post]
func (h *handler) GeneratePurchaseInvoicePDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF generated", resp_data)
}

// FavouritePurchaseInvoice godoc
// @Summary FavouritePurchaseInvoice
// @Description FavouritePurchaseInvoice
// @Tags PurchaseInvoice
// @Accept  json
// @Produce  json
// @param id path string true "PurchaseInvoice ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/favourite [post]
func (h *handler) FavouritePurchaseInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	purchaseInoviceId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": purchaseInoviceId,
	}
	_, err = h.service.ViewPurchaseInvoice(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	query := map[string]interface{}{
		"ID":      purchaseInoviceId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.FavouritePurchaseInvoice(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseInvoice is Marked as Favourite", map[string]string{"id": purchaseInoviceId})
}

// UnFavouritePurchaseInvoice godoc
// @Summary UnFavouritePurchaseInvoice
// @Description UnFavouritePurchaseInvoice
// @Tags PurchaseInvoice
// @Accept  json
// @Produce  json
// @param id path string true "PurchaseInvoice ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/{id}/unfavourite [post]
func (h *handler) UnFavouritePurchaseInvoice(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	purchaseInvoiceId := c.Param("id")

	query := map[string]interface{}{
		"ID":      purchaseInvoiceId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.UnFavouritePurchaseInvoice(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseInvoice is Marked as UnFavourite", map[string]string{"credinote_id": purchaseInvoiceId})
}

func (h *handler) GetPurchaseInvoiceTab(c echo.Context) (err error) {

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
	data, err := h.service.GetPurchaseInvoiceTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
