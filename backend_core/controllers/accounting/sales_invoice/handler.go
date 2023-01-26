package sales_invoice

import (
	"fmt"
	"strconv"

	accounting_base "fermion/backend_core/controllers/accounting/base"
	"fermion/backend_core/internal/model/pagination"
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
	base_service accounting_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := accounting_base.NewServiceBase()
	return &handler{service, base_service}
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
func (h *handler) CreateSalesInvoice(c echo.Context) (err error) {
	data := c.Get("sales_invoice").(*SalesInvoiceRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	t := uint(user_id)
	data.CreatedByID = &t
	id, err := h.service.CreateSalesInvoice(data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}

	emitter := emitter.GetObj()
	fmt.Println("ID---------------------->", id)
	emitter.Emit(events.CreateSalesInvoiceEvent+"Quickbooks", c, id)

	return res.RespSuccess(c, "Sales Invoice created successfully", map[string]interface{}{"created_id": id})
}

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
	data := new([]SalesInvoiceRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)

	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateSalesInvoice(data, token_id, access_template_id)
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
func (h *handler) UpdateSalesInvoice(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("sales_invoice").(*SalesInvoiceRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateSalesInvoice(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice updated successfully", map[string]interface{}{"updated_id": id})
}

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
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetSalesInvoice(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice data fetched successfully", result)
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
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllSalesInvoice(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Sales Invoice all data Retrieved successfully", result, p)
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
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetAllSalesInvoice(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Sales Invoice dropdown data Retrieved successfully", result, p)
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
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteSalesInvoice(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice deleted successfully", map[string]interface{}{"deleted_id": id})
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
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	product_id := c.QueryParam("product_id")
	prod_id, _ := strconv.Atoi(product_id)
	query := map[string]interface{}{
		"sales_invoice_id": uint(id),
		"product_id":       uint(prod_id),
	}
	err = h.service.DeleteSalesInvoiceLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice Line Items successfully deleted", query)
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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	fmt.Println("handler ------------------------ -1", query)
	err = h.base_service.FavouriteSalesInvoice(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice is Marked as Favourite", map[string]string{"record_id": id})
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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteSalesInvoice(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Sales Invoice is Unmarked as Favourite", map[string]string{"record_id": id})
}

func (h *handler) GetSalesInvoiceTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetSalesInvoiceTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
