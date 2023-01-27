package purchase_orders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"fermion/backend_core/controllers/accounting/sales_invoice"
	"fermion/backend_core/controllers/eda"
	orders_base "fermion/backend_core/controllers/orders/base"
	"fermion/backend_core/internal/model/core"
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

var PurchaseOrdersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if PurchaseOrdersHandler != nil {
		return PurchaseOrdersHandler
	}
	service := NewService()
	base_service := orders_base.NewServiceBase()
	PurchaseOrdersHandler = &handler{service, base_service}
	return PurchaseOrdersHandler
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
func (h *handler) CreatePurchaseOrdersEvent(c echo.Context) error {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("purchase_orders"),
	}
	eda.Produce(eda.CREATE_PURCHASE_ORDERS, requestPayload)
	return res.RespSuccess(c, "Purchase Order Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreatePurchaseOrders(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(PurchaseOrdersDTO)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.CreatePurchaseOrder(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_PURCHASE_ORDERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdatePurchaseOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": requestPayload.ID,
	}
	// eda.Produce(eda.CREATE_PURCHASE_ORDERS_ACK, *responseMessage)
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

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetPurchaseOrderFromCache(fmt.Sprint(metaData.TokenUserId))
	}

	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
	}

	dataInterface, err := h.service.ListPurchaseOrders(metaData, page)
	data := dataInterface.([]orders.PurchaseOrders)
	if err != nil {
		return res.RespError(c, err)
	}

	var dtoData []ListPurchaseOrdersDTO
	helpers.JsonMarshaller(data, &dtoData)
	for ind, val := range data {
		dtoData[ind].Amount = val.PoPaymentDetails.TotalAmount
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

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetPurchaseOrderFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.ListPurchaseOrders(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ListPurchaseOrdersDTO
	helpers.JsonMarshaller(data, &dtoData)

	return res.RespSuccessInfo(c, "dropdown list of purchase orders fetched", dtoData, page)
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

	metaData := c.Get("MetaData").(core.MetaData)
	poId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": poId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.ViewPurchaseOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response PurchaseOrdersDTO
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Purchase Orders Details retrieved succesfully", response)
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
func (h *handler) UpdatePurchaseOrdersEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("purchase_orders"),
	}

	eda.Produce(eda.UPDATE_PURCHASE_ORDERS, requestPayload)
	return res.RespSuccess(c, "Purchase Orders Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})

}

func (h *handler) UpdatePurchaseOrders(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(PurchaseOrdersDTO)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.UpdatePurchaseOrder(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_PURCHASE_ORDERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdatePurchaseOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_PURCHASE_ORDERS_ACK, *responseMessage)

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

	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeletePurchaseOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdatePurchaseOrderInCache(metaData)
	return res.RespSuccess(c, "purchase order deleted successfully", map[string]interface{}{"deleted_id": id})
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
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"po_id":      id,
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeletePurchaseOrderLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", metaData.Query)
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
	authToken := c.Request().Header.Get("Authorization")
	companyid := c.Request().Header.Get("CompanyId")

	id := c.Param("id")
	req, err := http.NewRequest("GET", "https://dev-api.eunimart.com/api/v1/purchase_orders/"+id, nil)
	req.Header.Add("Authorization", authToken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res.RespErr(c, err)

	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res.RespErr(c, err)

	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return res.RespErr(c, err)

	}

	data, ok := jsonData["data"].(map[string]interface{})
	if !ok {
		return res.RespErr(c, err)

	}
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, data["date_and_time"].(string))
	userdata, ok := data["CreatedBy"].(map[string]interface{})
	if err != nil {
		return res.RespErr(c, err)

	}

	amount, ok := data["po_payment_details"].(map[string]interface{})
	fmt.Println(companyid)
	companyname := "Eunimart Ltd"

	helpers.SendMail("Purchase_Orders", data["purchase_order_number"].(string), companyname, userdata["username"].(string), "<h3>Dear "+userdata["first_name"].(string)+" "+userdata["last_name"].(string)+"</h3><p>Thank you for choosing our business.</p> <p>Kindly find attached below: Sales Orders.</p><p>Here's a quick overview of the Sales Order: </p><p>-----------------------------------------------------------------------------</p><p>Sales Order : <b>"+data["purchase_order_number"].(string)+" </b></p><p>-----------------------------------------------------------------------------</p><p>Order Date  : "+t.Format("2006-01-02")+"</p><p>Amount  :"+strconv.FormatFloat(amount["total_amount"].(float64), 'f', 2, 64)+" </p><p>-----------------------------------------------------------------------------</p><p>Payment Invoice shall be sent separately.</p><p>Assuring you of our best services at all times.</p><p>Regards,</p>", "https://dev-api-files.s3.amazonaws.com/purchase_orders/"+companyname+"/"+companyid+"/test.pdf")
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	_, err = h.service.ViewPurchaseOrder(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	err = h.base_service.FavouritePurchaseOrders(metaData.Query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseOrder is Marked as Favourite", map[string]string{"id": id})
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
	metaData := c.Get("MetaData").(core.MetaData)

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouritePurchaseOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseOrder is Unmarked as Favourite", map[string]string{"id": id})
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

	metaData := c.Get("MetaData").(core.MetaData)

	productIdString := c.Param("product_id")
	productId, _ := strconv.Atoi(productIdString)

	metaData.Query = map[string]interface{}{
		"product_id": productId,
	}
	data, err := h.service.GetPurchaseHistory(metaData, page)
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
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id":  c.Param("id"),
		"tab": c.Param("tab"),
	}

	data, err := h.service.GetPurchaseOrderTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "success", data, page)
}
