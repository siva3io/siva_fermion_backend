package purchase_orders

import (
	"encoding/json"
	"fmt"
	"strconv"

	orders_base "fermion/backend_core/controllers/orders/base"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
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
	base_service orders_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := orders_base.NewServiceBase()
	return &handler{service, base_service}
}

// CreatePurchaseOrders godoc
// @Summary Create PurchaseOrders
// @Description Create PurchaseOrders
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body PurchaseOrdersDTO true "Purchase Orders Request Body"
// @Success 200 {object} PurchaseOrdersCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/create [post]
func (h *handler) CreatePurchaseOrders(c echo.Context) (err error) {

	data := c.Get("purchase_orders").(*PurchaseOrdersDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	seq_purchase := helpers.GenerateSequence("purchase", token_id, "purchase_orders")
	seq_ref := helpers.GenerateSequence("Ref", token_id, "purchase_orders")

	if data.GeneratePurchaseOrderNumber {
		data.PurchaseOrderNumber = seq_purchase
	}
	if data.GenerateReferenceId {
		data.ReferenceNumber = seq_ref
	}

	if err != nil {
		return res.RespError(c, err)
	}
	err = h.service.CreatePurchaseOrder(data, access_template_id, token_id)
	return res.RespSuccess(c, "Purchase Order Created", map[string]interface{}{"purchase_order_id": data.CreatedByID})
}

// GetListPurchaseOrders godoc
// @Summary Get all PurchaseOrders list
// @Description Get all PurchaseOrders list
// @Tags PurchaseOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllPurchaseOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders [get]
func (h *handler) ListPurchaseOrders(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListPurchaseOrdersDTO

	resp, err := h.service.ListPurchaseOrders(page, access_template_id, token_id, "LIST")

	data1 := resp.([]orders.PurchaseOrders)

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	for ind, val := range data1 {
		data[ind].Amount = val.PoPaymentDetails.TotalAmount
	}
	return res.RespSuccessInfo(c, "list of purchase orders fetched", data, page)
}

// ListPurchaseOrdersDropDown godoc
// @Summary Get all PurchaseOrders dropdown
// @Description Get all PurchaseOrders dropdown
// @Tags PurchaseOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllPurchaseOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/dropdown [get]
func (h *handler) ListPurchaseOrdersDropDown(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListPurchaseOrdersDTO

	resp, err := h.service.ListPurchaseOrders(page, access_template_id, token_id, "DROPDOWN_LIST")

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "dropdown of purchase orders fetched", data, page)
}

// ViewPurchaseOrders godoc
// @Summary View PurchaseOrders
// @Summary View PurchaseOrders
// @Description View PurchaseOrders
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "po_id"
// @Success 200 {object} PurchaseOrdersGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id} [get]
func (h *handler) ViewPurchaseOrders(c echo.Context) error {

	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	ID := c.Param("id")

	id, _ := strconv.Atoi(ID)

	query["id"] = int(id)

	fmt.Println("query", query)

	resp, err := h.service.ViewPurchaseOrder(query, access_template_id, token_id)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "purchase order fetched", resp)
}

// UpdatePurchaseOrders godoc
// @Summary Update PurchaseOrders
// @Description Update PurchaseOrders
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "po_id"
// @param RequestBody body PurchaseOrdersDTO true "Purchase Orders Request Body"
// @Success 200 {object} PurchaseOrdersUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id}/update [post]
func (h *handler) UpdatePurchaseOrders(c echo.Context) (err error) {
	var id = c.Param("id")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)

	ID, _ := strconv.Atoi(id)

	query["id"] = ID

	data := c.Get("purchase_orders").(*orders.PurchaseOrders)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdatePurchaseOrder(query, data, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "purchase order updated succesfully", nil)
}

// DeletePurchaseOrders godoc
// @Summary Delete PurchaseOrders
// @Description Delete PurchaseOrders
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "po_id"
// @Success 200 {object} PurchaseOrdersDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id}/delete [delete]
func (h *handler) DeletePurchaseOrders(c echo.Context) (err error) {
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{})

	query["id"] = ID
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeletePurchaseOrder(query, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "purchase order deleted successfully", map[string]int{"deleted_id": ID})
}

// DeletePurchaseOrderLines godoc
// @Summary Delete PurchaseOrderLines
// @Description Delete PurchaseOrderLines
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "po_id"
// @param product_id query string true "product_id"
// @Success 200 {object} PurchaseOrderLinesDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/order_lines/{id}/delete [delete]
func (h *handler) DeletePurchaseOrderLines(c echo.Context) (err error) {
	var product_id = c.QueryParam("product_id")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = map[string]interface{}{
		"po_id":      ID,
		"product_id": product_id,
	}
	err = h.service.DeletePurchaseOrderLines(query, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", query)
}

// SearchPurchaseOrders godoc
// @Summary Search PurchaseOrders
// @Description Search PurchaseOrders
// @Tags PurchaseOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   query 	query 	string 	true "query"
// @Success 200 {array} GetAllPurchaseOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/search [get]
func (h *handler) SearchPurchaseOrders(c echo.Context) error {

	query := c.QueryParam("query")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListPurchaseOrdersDTO

	resp, err := h.service.SearchPurchaseOrders(query, access_template_id, token_id)

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "OK", data)
}

// EmailPurchaseOrders godoc
// @Summary Email PurchaseOrders
// @Description Email PurchaseOrders
// @Tags PurchaseOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "po_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id}/sendEmail [post]
func (h *handler) EmailPurchaseOrders(c echo.Context) error {
	return res.RespSuccess(c, "email sent successfully", nil)
}

// DownloadPurchaseOrders godoc
// @Summary Download PurchaseOrders
// @Description Download PurchaseOrders
// @Tags PurchaseOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "po_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id}/downloadPdf [post]
func (h *handler) DownloadPurchaseOrders(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file downloaded successfully", resp_data)
}

// GeneratePurchaseOrdersPDF godoc
// @Summary Generate PDF PurchaseOrders
// @Description Generate PDF PurchaseOrders
// @Tags PurchaseOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "po_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id}/generatePdf [post]
func (h *handler) GeneratePurchaseOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF generated", resp_data)
}

// FavouritePurchaseOrders godoc
// @Summary FavouritePurchaseOrders
// @Description FavouritePurchaseOrders
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @param id path string true "PurchaseOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id}/favourite [post]
func (h *handler) FavouritePurchaseOrders(c echo.Context) (err error) {
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
	_, er := h.service.ViewPurchaseOrder(q, access_template_id, token_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouritePurchaseOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseOrders is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouritePurchaseOrders godoc
// @Summary UnFavouritePurchaseOrders
// @Description UnFavouritePurchaseOrders
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @param id path string true "PurchaseOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{id}/unfavourite [post]
func (h *handler) UnFavouritePurchaseOrders(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouritePurchaseOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseOrders is Unmarked as Favourite", map[string]string{"id": id})
}

// GetPurchaseHistory godoc
// @Summary GetPurchaseHistory
// @Summary GetPurchaseHistory
// @Description GetPurchaseHistory
// @Tags PurchaseOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param product_id path string true "product_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_orders/{product_id}/history [get]
func (h *handler) GetPurchaseHistory(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	productIdString := c.Param("product_id")

	productId, _ := strconv.Atoi(productIdString)

	data, err := h.service.GetPurchaseHistory(uint(productId), page, access_template_id, token_id)

	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

// GetPurchaseOrderTab godoc
// @Summary GetPurchaseOrderTab
// @Summary GetPurchaseOrderTab
// @Description GetPurchaseOrderTab
// @Tags PurchaseOrders
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
// @Router /api/v1/purchase_orders/{id}/filter_module/{tab} [get]
func (h *handler) GetPurchaseOrderTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetPurchaseOrderTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
