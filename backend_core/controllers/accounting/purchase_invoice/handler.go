package purchase_invoice

import (
	"encoding/json"
	"fmt"
	"strconv"

	accounting_base "fermion/backend_core/controllers/accounting/base"
	"fermion/backend_core/internal/model/accounting"
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
	base_service accounting_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := accounting_base.NewServiceBase()
	return &handler{service, base_service}
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
func (h *handler) CreatePurchaseInvoice(c echo.Context) error {

	input_data := c.Get("purchase_invoice").(*accounting.PurchaseInvoice)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	input_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err := h.service.CreatePurchaseInvoice(input_data, token_id, access_template_id)

	if err != nil {
		return res.RespError(c, err)
	}

	emitter := emitter.GetObj()
	fmt.Println("ID---------------------->", input_data.ID)
	emitter.Emit(events.CreatePurchaseInvoiceEvent+"Quickbooks", c, input_data.ID)

	return res.RespSuccess(c, "Purchase Invoice Created", map[string]interface{}{"purchase_invoice_id": input_data.ID})
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

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListPurchaseInvoiceDTO

	resp, err := h.service.ListPurchaseInvoice(page, token_id, access_template_id, "LIST")

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "list of purchase invoice fetched", data, page)
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

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListPurchaseInvoiceDTO

	resp, err := h.service.ListPurchaseInvoice(page, token_id, access_template_id, "DROPDOWN_LIST")

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "dropdown list of purchase invoice fetched", data, page)
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

	var query = make(map[string]interface{}, 0)

	ID := c.Param("id")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	id, _ := strconv.Atoi(ID)

	query["id"] = int(id)

	fmt.Println("query", query)

	resp, err := h.service.ViewPurchaseInvoice(query, token_id, access_template_id)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "purchase invoice fetched", resp)
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
func (h *handler) UpdatePurchaseInvoice(c echo.Context) (err error) {
	var id = c.Param("id")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)

	ID, _ := strconv.Atoi(id)

	query["id"] = ID

	data := c.Get("purchase_invoice").(*accounting.PurchaseInvoice)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdatePurchaseInvoice(query, data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "purchase invoice updated succesfully", nil)
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
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)

	var query = make(map[string]interface{})

	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeletePurchaseInvoice(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "purchase invoice deleted successfully", map[string]int{"deleted_id": ID})
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
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	ID, _ := strconv.Atoi(id)
	var query = map[string]interface{}{
		"purchase_invoice_id": ID,
		"product_id":          product_id,
	}
	err = h.service.DeletePurchaseInvoiceLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "purchase invoice line deleted successfully", query)
}

// SearchPurchaseInvoice godoc
// @Summary Search PurchaseInvoice
// @Description Search PurchaseInvoice
// @Tags PurchaseInvoice
// @Accept json
// @Produce json
// @Param   query 	query 	string 	true "query"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {array} GetAllPurchaseInvoiceResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_invoice/search [get]
func (h *handler) SearchPurchaseInvoice(c echo.Context) error {

	query := c.QueryParam("query")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListPurchaseInvoiceDTO

	resp, err := h.service.SearchPurchaseInvoice(query, token_id, access_template_id)

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "OK", data)
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
		"id": ID,
	}
	_, er := h.service.ViewPurchaseInvoice(q, token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouritePurchaseInvoice(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseInvoice is Marked as Favourite", map[string]string{"id": id})
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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouritePurchaseInvoice(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseInvoice is Unmarked as Favourite", map[string]string{"id": id})
}

func (h *handler) GetPurchaseInvoiceTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetPurchaseInvoiceTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
