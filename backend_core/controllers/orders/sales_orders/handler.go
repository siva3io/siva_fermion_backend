package sales_orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	/*"io/ioutil"
	"net/http"
	"strings"*/

	"strconv"

	concurrency_management "fermion/backend_core/controllers/concurrency_management"
	"fermion/backend_core/controllers/cores"
	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/controllers/mdm/locations"
	"fermion/backend_core/controllers/mdm/products"
	"fermion/backend_core/controllers/omnichannel/channel"
	orders_base "fermion/backend_core/controllers/orders/base"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"

	shipping_model "fermion/backend_core/internal/model/shipping"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"fermion/backend_core/controllers/accounting/sales_invoice"
	shipping_order "fermion/backend_core/controllers/shipping/shipping_orders"
	ipaas_utils "fermion/backend_core/ipaas_core/utils"

	"github.com/jszwec/csvutil"
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
	service              Service
	base_service         orders_base.ServiceBase
	concurrencyService   concurrency_management.Service
	ShippingOrderService shipping_order.Service
	SalesInvoiceService  sales_invoice.Service
	ChannelService       channel.Service
	cores_service        cores.Service
	locations_service    locations.Service
	Product_service      products.Service
}

var SalesOrdersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if SalesOrdersHandler != nil {
		return SalesOrdersHandler
	}
	SalesOrdersHandler = &handler{
		NewService(),
		orders_base.NewServiceBase(),
		concurrency_management.NewService(),
		shipping_order.NewService(),
		sales_invoice.NewService(),
		channel.NewService(),
		cores.NewService(),
		locations.NewService(),
		products.NewService(),
	}
	return SalesOrdersHandler
}

// CreateSalesOrders godoc
// @Summary Create SalesOrders
// @Description Create SalesOrders
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body SalesOrdersDTO true "Sales Orders Request Body"
// @Success 200 {object} SalesOrdersCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/create [post]
func (h *handler) CreateSalesOrdersEvent(c echo.Context) error {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("sales_orders"),
	}
	eda.Produce(eda.CREATE_SALES_ORDER, requestPayload)
	return res.RespSuccess(c, "Sales Order Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateSalesOrders(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(orders.SalesOrders)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.CreateSalesOrder(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SALES_ORDER_ACK, *responseMessage)
		return

	}
	// cache implementation
	UpdateSalesOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": requestPayload.ID,
	}
	// eda.Produce(eda.CREATE_SALES_ORDER_ACK, *responseMessage)

	//================create Shipment =========================
	if len(requestPayload.SalesOrderLines) == 0 {
		fmt.Println("---------------No Order Lines----------------------")
		return
	}
	shippingOrder := new(shipping_model.ShippingOrder)
	shippingOrder.OrderId = &requestPayload.ID
	channel_name := requestPayload.ChannelName
	var channel_data channel.ChannelDTO
	fmt.Println("chanel name------------", channel_name)
	if channel_name != "" {
		query := map[string]interface{}{
			"name": channel_name,
		}
		channel_data, _ = h.ChannelService.GetChannelService(query)
	} else {
		channel_data.ID = 1
	}
	shippingOrder.ChannelId = &channel_data.ID
	fmt.Println("channel id-----------------------", channel_data.ID)
	var billingAddress CustomerBillingorShippingAddress
	json.Unmarshal(requestPayload.CustomerBillingAddress, &billingAddress)

	var receivingAddress CustomerBillingorShippingAddress
	json.Unmarshal(requestPayload.CustomerShippingAddress, &receivingAddress)

	shippingOrder.BillingAddress, _ = json.Marshal(map[string]interface{}{
		"sender_name":   billingAddress.ContactPersonName,
		"mobile_number": billingAddress.ContactPersonNumber,
		"address_line1": billingAddress.AddressLine1,
		"address_line2": billingAddress.AddressLine2,
		"address_line3": billingAddress.AddressLine3,
		"zipcode":       billingAddress.PinCode,
		"city":          billingAddress.City,
		"state":         billingAddress.State,
		"country":       billingAddress.Country,
		"landmark":      billingAddress.Landmark,
	})
	shippingOrder.ReceiverAddress, _ = json.Marshal(map[string]interface{}{
		"receiver_name": receivingAddress.ContactPersonName,
		"mobile_number": receivingAddress.ContactPersonNumber,
		"address_line1": receivingAddress.AddressLine1,
		"address_line2": receivingAddress.AddressLine2,
		"address_line3": receivingAddress.AddressLine3,
		"zipcode":       receivingAddress.PinCode,
		"city":          receivingAddress.City,
		"state":         receivingAddress.State,
		"country":       receivingAddress.Country,
		"landmark":      receivingAddress.Landmark,
	})
	fmt.Println("company id-------", requestPayload.CompanyId)
	company_details, err := h.cores_service.GetCompanyPreferences(map[string]interface{}{
		"id": requestPayload.CompanyId,
	})

	locationPagnation := new(pagination.Paginatevalue)
	lookup_id, err := helpers.GetLookupcodeId("LOCATION_TYPE", "OFFICE")
	if err != nil {
		return
	}
	locationPagnation.Filters = fmt.Sprintf("[[\"company_id\",\"=\",%v],[\"location_type_id\",\"=\",%v]]", edaMetaData.CompanyId, lookup_id)
	location_resp, err := h.locations_service.GetLocationList(edaMetaData, locationPagnation)

	if err != nil {
		fmt.Println("-----------------LOCATION_ERROR------", err.Error())
		return
	}
	var location_data = make([]map[string]interface{}, 0)
	helpers.JsonMarshaller(location_resp, &location_data)

	length := len(location_data)
	if length > 0 {
		address_data := location_data[0]["address"].(map[string]interface{})

		shippingOrder.SenderAddress, _ = json.Marshal(map[string]interface{}{
			"sender_name":   company_details.Name,
			"email":         company_details.Email,
			"mobile_number": company_details.Phone,
			"address_line1": address_data["address_line_1"],
			"address_line2": address_data["address_line_2"],
			"address_line3": address_data["address_line_3"],
			"state":         address_data["state"].(map[string]interface{})["name"],
			"country":       address_data["country"].(map[string]interface{})["name"],
			"landmark":      address_data["land_mark"],
			"zipcode":       address_data["pin_code"],
			"city":          address_data["city"],
		})
	}
	pincode, _ := strconv.Atoi(receivingAddress.PinCode)
	shippingOrder.DestinationZipcode = int32(pincode)
	pincode, _ = strconv.Atoi(billingAddress.PinCode)
	shippingOrder.OriginZipcode = int32(pincode)
	shippingPartnerId := uint(8)
	shippingOrder.ShippingPartnerId = &shippingPartnerId

	shippingOrder.ShippingPaymentTypeID = requestPayload.PaymentTypeId
	for _, sales_order_line := range requestPayload.SalesOrderLines {
		ShippingOrderLine := new(shipping_model.ShippingOrderLines)
		ShippingOrderLine.ProductId = sales_order_line.ProductId
		ShippingOrderLine.ProductTemplateId = sales_order_line.ProductTemplateId
		ShippingOrderLine.ItemQuantity = int32(sales_order_line.Quantity)
		ShippingOrderLine.UnitPrice = float64(sales_order_line.Price)
		shippingOrder.ShippingOrderLines = append(shippingOrder.ShippingOrderLines, *ShippingOrderLine)

	}
	shippingOrder.Quantity = int32(requestPayload.TotalQuantity)
	shippingOrder.OrderDate = requestPayload.CreatedDate

	shippingOrder.SetPickupDate = requestPayload.ReadyTOShip
	productQuery := map[string]interface{}{
		"id": requestPayload.SalesOrderLines[0].ProductTemplateId,
	}
	product_details, _ := h.Product_service.GetTemplateView(productQuery, fmt.Sprint(edaMetaData.TokenUserId), fmt.Sprint(edaMetaData.AccessTemplateId))
	var template products.TemplateReponseDTO
	helpers.JsonMarshaller(product_details, &template)
	shippingOrder.PackageDetails = template.PackageDimensions
	shippingOrder.OrderValue = float64(requestPayload.SoPaymentDetails.TotalAmount)
	shippingOrder.OrderDate = requestPayload.CreatedDate
	lookup_id, _ = helpers.GetLookupcodeId("SHIPPING_TYPE", "DOMESTIC")
	shippingOrder.ShippingTypeId = &lookup_id
	h.ShippingOrderService.CreateShippingOrder(edaMetaData, shippingOrder)

}

// GetListSalesOrders godoc
// @Summary Get all SalesOrders list
// @Description Get all SalesOrders list
// @Tags SalesOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllSalesOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders [get]
func (h *handler) ListSalesOrders(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetSalesOrderFromCache(fmt.Sprint(metaData.TokenUserId))
	}

	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
	}

	dataInterface, err := h.service.ListSalesOrders(metaData, page)
	data := dataInterface.([]orders.SalesOrders)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ListSalesOrdersDTO
	helpers.JsonMarshaller(data, &dtoData)
	for ind, val := range data {
		dtoData[ind].Amount = val.SoPaymentDetails.TotalAmount
	}
	return res.RespSuccessInfo(c, "list of sales orders fetched", data, page)
}

// GetListSalesOrdersDropDown godoc
// @Summary Get all SalesOrders DropDown
// @Description Get all SalesOrders DropDown
// @Tags SalesOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllSalesOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/dropdown [get]
func (h *handler) ListSalesOrdersDropDown(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetSalesOrderFromCache(fmt.Sprint(metaData.TokenUserId))
	}
	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
	}

	dataInterface, err := h.service.ListSalesOrders(metaData, page)
	data := dataInterface.([]orders.SalesOrders)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ListSalesOrdersDTO
	helpers.JsonMarshaller(data, &dtoData)
	for ind, val := range data {
		dtoData[ind].Amount = val.SoPaymentDetails.TotalAmount
	}

	return res.RespSuccessInfo(c, "dropdown list of sales orders fetched", data, page)
}

// ViewSalesOrders godoc
// @Summary View SalesOrders
// @Summary View SalesOrders
// @Description View SalesOrders
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @Success 200 {object} SalesOrdersGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id} [get]
func (h *handler) ViewSalesOrders(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": c.Param("id"),
		// "company_id": metaData.CompanyId,
	}
	data, err := h.service.ViewSalesOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "sales order fetched", data)
}

// UpdateSalesOrders godoc
// @Summary Update SalesOrders
// @Description Update SalesOrders
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @param RequestBody body SalesOrdersDTO true "Sales Orders Request Body"
// @Success 200 {object} SalesOrdersUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id}/update [post]
func (h *handler) UpdateSalesOrdersEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("sales_orders"),
	}
	//--------------------------concurrency_management_start-----------------------------------------
	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(uint(id), "status_id", "sales_orders")
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(uint(id), PreviousStatusId, "status_id", "sales_orders")
	//--------------------------concurrency_management_end-----------------------------------------

	eda.Produce(eda.UPDATE_SALES_ORDER, requestPayload)
	return res.RespSuccess(c, "Sales orders Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateSalesOrders(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(orders.SalesOrders)
	helpers.JsonMarshaller(request["data"], requestPayload)

	if request["data"].(map[string]interface{})["status_id"] != nil {
		statusId := uint(request["data"].(map[string]interface{})["status_id"].(float64))
		requestPayload.StatusId = &statusId
	}
	err := h.service.UpdateSalesOrder(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_SALES_ORDER_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateSalesOrderInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": requestPayload.ID,
	}
	// eda.Produce(eda.UPDATE_SALES_ORDER_ACK, *responseMessage)

}

// DeleteSalesOrders godoc
// @Summary Delete SalesOrders
// @Description Delete SalesOrders
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @Success 200 {object} SalesOrdersDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id}/delete [delete]
func (h *handler) DeleteSalesOrders(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteSalesOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateSalesOrderInCache(metaData)
	return res.RespSuccess(c, "sales order deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteSalesOrderLines godoc
// @Summary Delete SalesOrderLines
// @Description Delete SalesOrderLines
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @param product_id query string true "product_id"
// @Success 200 {object} SalesOrderLinesDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/order_lines/{id}/delete [delete]
func (h *handler) DeleteSalesOrderLines(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"so_id":      id,
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeleteSalesOrderLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", metaData.Query)
}

// EmailSalesOrders godoc
// @Summary Email SalesOrders
// @Description Email SalesOrders
// @Tags SalesOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id}/sendEmail [post]
func (h *handler) EmailSalesOrders(c echo.Context) error {

	authToken := c.Request().Header.Get("Authorization")
	companyid := c.Request().Header.Get("CompanyId")

	id := c.Param("id")
	req, err := http.NewRequest("GET", "https://dev-api.eunimart.com/api/v1/sales_orders/"+id, nil)
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
	t, err := time.Parse(layout, data["so_date"].(string))
	useremail, ok := data["CreatedBy"].(map[string]interface{})
	if err != nil {
		return res.RespErr(c, err)

	}

	amount, ok := data["so_payment_details"].(map[string]interface{})

	companyname := "Eunimart Ltd"

	helpers.SendMail("Sales_orders", data["sales_order_number"].(string), companyname, useremail["username"].(string), "<h3>Dear "+data["customer_name"].(string)+"</h3><p>Thank you for choosing our business.</p> <p>Kindly find attached below: Sales Orders.</p><p>Here's a quick overview of the Sales Order: </p><p>-----------------------------------------------------------------------------</p><p>Sales Order : <b>"+data["sales_order_number"].(string)+" </b></p><p>-----------------------------------------------------------------------------</p><p>Order Date  : "+t.Format("2006-01-02")+"</p><p>Amount  :"+strconv.FormatFloat(amount["total_amount"].(float64), 'f', 2, 64)+" </p><p>-----------------------------------------------------------------------------</p><p>Payment Invoice shall be sent separately.</p><p>Assuring you of our best services at all times.</p><p>Regards,</p>", "https://dev-api-files.s3.amazonaws.com/sales_orders/"+companyname+"/"+companyid+"/test.pdf")
	return res.RespSuccess(c, "email sent successfully", nil)

}

// DownloadSalesOrders godoc
// @Summary Download SalesOrders
// @Description Download SalesOrders
// @Tags SalesOrdersb
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id}/downloadPdf [post]
func (h *handler) DownloadSalesOrders(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file downloaded successfully", resp_data)
}

// GenerateSalesOrdersPDF godoccore_service
// @Summary Generate PDF SalesOrders
// @Description Generate PDF SalesOrders
// @Tags SalesOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "so_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id}/generatePdf [post]
func (h *handler) GenerateSalesOrdersPDF(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	salesOrderId := c.Param("id")
	heading := c.QueryParam("heading")
	helpers.PrettyPrint("=====================>", heading)
	helpers.PrettyPrint("=====================>", metaData)
	helpers.PrettyPrint("=====================>", salesOrderId)
	// company_id := c.Get("CompanyId").(string)
	x := sales_invoice.NewHandler()
	er := x.GetSalesInvoicePdf(c)
	if er != nil {
		return res.RespSuccess(c, "Generated Successfully", "")
	}
	// metaData.Query = map[string]interface{}{
	// 	"id": salesOrderId,
	// }

	// data, err := h.service.ViewSalesOrder(metaData)
	// if err != nil {
	// 	return res.RespErr(c, err)
	// }
	// var result orders.SalesOrders
	// jsonData, _ := json.Marshal(data)
	// json.Unmarshal(jsonData, &result)
	// billing_address := make(map[string]interface{}, 0)
	// shipping_address := make(map[string]interface{}, 0)
	// var products = make([]map[string]interface{}, 0)
	// byte_billing_json, _ := json.Marshal(result.CustomerBillingAddress)
	// byte_shipping_json, _ := json.Marshal(result.CustomerShippingAddress)
	// products_json, _ := json.Marshal(&result.SalesOrderLines)
	// json.Unmarshal(products_json, &products)

	// json.Unmarshal(byte_billing_json, &billing_address)
	// json.Unmarshal(byte_shipping_json, &shipping_address)

	// // length := len(billing_address)
	// // if length <= 0 {
	// // 	fmt.Println("------------BILLING ADDRESS EMPTY---------")
	// // 	return res.RespErr(c, errors.New("BILLING ADDRESS IS EMPTY"))
	// // }
	// // length = len(shipping_address)
	// // if length <= 0 {
	// // 	fmt.Println("------------SHIPPING ADDRESS EMPTY---------")
	// // 	return res.RespErr(c, errors.New("SHIPPING ADDRESS IS EMPTY"))
	// // }

	// exp_shipment_date, _ := result.ExpectedShippingDate.Value()
	// var shipment_data string
	// byte_time, _ := json.Marshal(exp_shipment_date)
	// json.Unmarshal(byte_time, &shipment_data)
	// company_details, err := h.cores_service.GetCompanyPreferences(map[string]interface{}{
	// 	"id": company_id,
	// })
	// if err != nil {
	// 	return res.RespErr(c, err)
	// }
	// r := NewRequestPdf("")

	// templatePath := "backend_core/views/sales_orders.html"
	// outputPath := "test.pdf"

	// // lookup_id, err := helpers.GetLookupcodeId("LOCATION_TYPE", "OFFICE")

	// // if err != nil {
	// // 	return res.RespErr(c, err)
	// // }
	// /*locationPagnation := new(pagination.Paginatevalue)
	// locationPagnation.Filters = "[[\"company_id\",\"=\",\"1\"],[\"location_type_id\",\"=\",\"68\"]]"
	// locationPagnation.Filters = fmt.Sprintf("[[\"company_id\",\"=\",%v],[\"location_type_id\",\"=\",%v]]", company_id, lookup_id)
	// location_resp, err := h.location_service.GetLocationList(metaData, locationPagnation)
	// if err != nil {
	// 	fmt.Println("-----------------LOCATION_ERROR------", err.Error())
	// 	return res.RespErr(c, err)
	// }
	// var location_data = make([]map[string]interface{}, 0)
	// json_data, _ := json.Marshal(location_resp)
	// json.Unmarshal(json_data, &location_data)
	// length := len(location_data)
	// if length <= 0 {
	// 	fmt.Println("------------LOCATION EMPTY---------")
	// 	return res.RespErr(c, errors.New("LOCATION IS EMPTY"))
	// }

	// address_data := location_data[0]["address"].(map[string]interface{})*/

	// template_data := map[string]interface{}{
	// 	"bill_address_1":     billing_address["data"].(map[string]interface{})["address_line_1"],
	// 	"bill_address_2":     billing_address["data"].(map[string]interface{})["address_line2_2"],
	// 	"bill_address_3":     billing_address["data"].(map[string]interface{})["address_line3_3"],
	// 	"bill_address_code":  billing_address["data"].(map[string]interface{})["pin_code"],
	// 	"bill_address_phone": billing_address["data"].(map[string]interface{})["contact_person_number"],
	// 	"bill_land_mark":     billing_address["data"].(map[string]interface{})["land_mark"],
	// 	"bill_location_name": billing_address["data"].(map[string]interface{})["location_name"],

	// 	"ship_address_1":     shipping_address["address_line_1"],
	// 	"ship_address_2":     shipping_address["address_line_2"],
	// 	"ship_address_3":     shipping_address["address_line_3"],
	// 	"ship_address_code":  shipping_address["pin_code"],
	// 	"ship_address_phone": shipping_address["contact_person_number"],
	// 	"ship_land_mark":     shipping_address["land_mark"],
	// 	"ship_location_name": shipping_address["location_name"],

	// 	"sales_order_number": result.SalesOrderNumber,
	// 	"sales_order_date":   shipment_data[:10],
	// 	"order_number":       result.ReferenceNumber,
	// 	"product_details":    products,
	// 	"sub_amt":            result.SoPaymentDetails.SubTotal,
	// 	"tot_amt":            result.SoPaymentDetails.TotalAmount,
	// 	"tot_tax":            result.SoPaymentDetails.Tax,
	// 	"payment":            result.PaymentType.DisplayName,
	// 	"company_name":       company_details.Name,
	// 	"company_email":      company_details.Email,
	// 	"company_website":    company_details.Website,
	// 	/*"company_adrs_1":       address_data["address_line_1"],
	// 	"company_adrs_2":       address_data["address_line_2"],
	// 	"company_state":        address_data["state"].(map[string]interface{})["name"],
	// 	"company_country":      address_data["country"].(map[string]interface{})["name"],*/
	// }
	// if err := r.ParseTemplate(templatePath, template_data); err == nil {
	// 	ok, _ := r.GeneratePDF(outputPath)
	// 	fmt.Println("pdf generated:", ok)
	// } else {
	// 	fmt.Println(err)
	// }
	// c.Response().Header().Set(
	// 	"Content-Disposition",
	// 	"attachment; filename="+outputPath,
	// )

	// //Upload Files

	// path := "sales_orders/" + company_details.Name + "/" + salesOrderId + "/"
	// err = helpers.UploadFile("./test.pdf", path)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// s3path := "https://dev-api-files.s3.amazonaws.com/" + path + "test.pdf"
	// defer os.Remove("./test.pdf")

	return c.Attachment("test.pdf", "test")
}

// FavouriteSalesOrders godoc
// @Summary FavouriteSalesOrders
// @Description FavouriteSalesOrders
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @param id path string true "SalesOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id}/favourite [post]
func (h *handler) FavouriteSalesOrders(c echo.Context) error {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	_, err := h.service.ViewSalesOrder(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	err = h.base_service.FavouriteSalesOrders(metaData.Query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "SalesOrders is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouriteSalesOrders godoc
// @Summary UnFavouriteSalesOrders
// @Description UnFavouriteSalesOrders
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @param id path string true "SalesOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{id}/unfavourite [post]
func (h *handler) UnFavouriteSalesOrders(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouriteSalesOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "SalesOrders is Unmarked as Favourite", map[string]string{"id": id})
}

// ----------------------Channel API's --------------------------------------------------------
func (h *handler) ChannelSalesOrderUpsertEvent(c echo.Context) (err error) {

	var arrayData []orders.SalesOrders
	var objectData orders.SalesOrders
	var data interface{}

	c.Bind(&data)

	jsonData, _ := json.Marshal(data)

	err = json.Unmarshal(jsonData, &arrayData)
	if err != nil {
		_ = json.Unmarshal(jsonData, &objectData)
		arrayData = append(arrayData, objectData)
	}

	edaMetaData := c.Get("MetaData").(core.MetaData)
	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      arrayData,
	}

	eda.Produce(eda.UPSERT_SALES_ORDER, requestPayload)
	return res.RespSuccess(c, "Multiple Records Upsert Inprogress", "Upsert Inprogress")
}

func (h *handler) ChannelSalesOrderUpsert(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	var requestPayload []orders.SalesOrders
	helpers.JsonMarshaller(request["data"], requestPayload)

	_, err := h.service.ChannelSalesOrderUpsert(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPSERT_SALES_ORDER_ACK, *responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_message": "channel sales order upserted successfully",
	}
	// eda.Produce(eda.UPSERT_SALES_ORDER_ACK, *responseMessage)
}

// GetSalesHistory godoc
// @Summary GetSalesHistory
// @Summary GetSalesHistory
// @Description GetSalesHistory
// @Tags SalesOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param product_id path string true "product_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/{product_id}/history [get]
func (h *handler) GetSalesHistory(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)

	productIdString := c.Param("product_id")
	productId, _ := strconv.Atoi(productIdString)

	metaData.Query = map[string]interface{}{
		"product_id": productId,
	}
	data, err := h.service.GetSalesHistory(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

// GetSalesOrderTab godoc
// @Summary GetSalesOrderTab
// @Summary GetSalesOrderTab
// @Description GetSalesOrderTab
// @Tags SalesOrders
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
// @Router /api/v1/sales_orders/{id}/filter_module/{tab} [get]
func (h *handler) GetSalesOrderTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id":  c.Param("id"),
		"tab": c.Param("tab"),
	}

	data, err := h.service.GetSalesOrderTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "success", data, page)
}

func (h *handler) TrackSalesOrder(c echo.Context) (err error) {
	queryParams := c.Request().URL.Query()

	query := make(map[string]interface{}, 0)

	for key, value := range queryParams {
		query[key] = value
	}

	status_history := h.service.TrackSalesOrder(query)

	return res.RespSuccess(c, "sales order history successfully", status_history)
}

func (h *handler) ListSalesOrdersAdmin(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	dataInterface, err := h.service.ListSalesOrders(metaData, page)
	data := dataInterface.([]orders.SalesOrders)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []SalesOrdersDTO
	helpers.JsonMarshaller(data, &dtoData)
	for ind, val := range data {
		dtoData[ind].Amount = val.SoPaymentDetails.TotalAmount
		dtoData[ind].Item = dtoData[ind].SalesOrderLines[0]
	}
	return res.RespSuccessInfo(c, "list of sales orders fetched", dtoData, page)
}

func (h *handler) ViewSalesOrderAdmin(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": c.Param("id"),
	}
	data, err := h.service.ViewSalesOrder(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "sales order fetched", data)
}

func (h *handler) GenerateSalesOrdersCSV(c echo.Context) error {
	host := c.Request().Host
	scheme := c.Scheme()

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	page.Per_page = -1

	dataInterface, err := h.service.ListSalesOrders(metaData, page)
	data := dataInterface.([]orders.SalesOrders)
	if err != nil {
		return res.RespErr(c, err)
	}
	var dtoData []SalesOrdersDTO
	helpers.JsonMarshaller(data, &dtoData)
	for ind, val := range data {
		dtoData[ind].Amount = val.SoPaymentDetails.TotalAmount
	}

	mapper, err := helpers.ReadFile("./backend_core/controllers/orders/sales_orders/csv_mapper_template.json")
	if err != nil {
		return res.RespErr(c, err)
	}
	authToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTEsIlVzZXJuYW1lIjoic3VwZXJfdXNlckBldW5pbWFydC5jb20iLCJhY2Nlc3NfdGVtcGxhdGVfaWQiOjEsImNvbXBhbnlfaWQiOjEsImV4cCI6MTcwNzU0NTI0NCwiZmlyc3RfbmFtZSI6IlN1cGVyIiwibGFzdF9uYW1lIjoiVXNlciIsInJvbGVfaWQiOm51bGwsInVzZXJfdHlwZXMiOlt7ImlkIjo2MzQsIm5hbWUiOiJTRUxMRVIifV19.yFfpIOAKL2--8uutOdVvgGnO9Lyruz72bkM0cpFT8fU"

	response, err := helpers.MakeRequest(helpers.Request{
		Method: "POST",
		Scheme: "http",
		Host:   "localhost:3031",
		Path:   "ipaas/boson_convertor",
		Header: map[string]string{
			"Authorization": authToken,
			"Content-Type":  "application/json",
		},
		Body: map[string]interface{}{
			"data": map[string]interface{}{
				"input_data":      dtoData,
				"mapper_template": mapper,
				"feature_session_variables": []map[string]interface{}{
					{
						"key":   "host",
						"value": "localhost:3031",
					},
					{
						"key":   "Authorization",
						"value": authToken,
					},
				},
			},
		},
	})
	if err != nil {
		return res.RespErr(c, err)
	}

	var mappedResponse interface{}
	var errorResponse interface{}
	if data := response.(map[string]interface{})["data"]; data != nil {
		mappedResponse = data.(map[string]interface{})["mapped_response"]
		errorResponse = data.(map[string]interface{})["error_message"]
	}

	if mappedResponse == nil {
		errorResponseString, _ := json.Marshal(errorResponse)
		return res.RespErr(c, errors.New(string(errorResponseString)))
	}

	var csvDto []csvDto
	err = helpers.JsonMarshaller(response.(map[string]interface{})["data"].(map[string]interface{})["mapped_response"], &csvDto)
	if err != nil {
		return res.RespErr(c, err)
	}

	dta, err := csvutil.Marshal(csvDto)
	if err != nil {
		return res.RespErr(c, err)
	}

	module := "s3"
	subModule := "files"
	token := c.Request().Header.Get("Authorization")

	var request interface{}
	err = helpers.JsonMarshaller(map[string]interface{}{"data": dta}, &request)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dt []interface{}
	dt = append(dt, request)
	unique := helpers.GenerateRandID(1, 2000, "tech")
	fileName := fmt.Sprintf("orders_%s", unique)
	response, err = ipaas_utils.UploadFile(dt, fileName, unique, 12, scheme, host, ipaas_utils.AWS, module, subModule, ".csv", ipaas_utils.BYTES, token)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Link of the downloaded data", response)
}
