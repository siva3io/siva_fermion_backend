package sales_orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	concurrency_management "fermion/backend_core/controllers/concurrency_management"
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
	service            Service
	base_service       orders_base.ServiceBase
	concurrencyService concurrency_management.Service
}

func NewHandler() *handler {
	service := NewService()
	base_service := orders_base.NewServiceBase()
	concurrencyService := concurrency_management.NewService()
	return &handler{service, base_service, concurrencyService}
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
func (h *handler) CreateSalesOrders(c echo.Context) error {

	input_data := c.Get("sales_orders").(*orders.SalesOrders)
	token_id := c.Get("TokenUserID").(string)
	input_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	access_template_id := c.Get("AccessTemplateId").(string)
	err := h.service.CreateSalesOrder(input_data, access_template_id, token_id)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "Sales Order Created", map[string]interface{}{"sales_order_id": input_data.ID})
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

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	var data []ListSalesOrdersDTO
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.ListSalesOrders(page, access_template_id, token_id, "LIST")

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
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

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	var data []ListSalesOrdersDTO
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.ListSalesOrders(page, access_template_id, token_id, "DROPDOWN_LIST")

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
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

	var query = make(map[string]interface{}, 0)

	ID := c.Param("id")

	id, _ := strconv.Atoi(ID)

	query["id"] = int(id)

	fmt.Println("query", query)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.ViewSalesOrder(query, access_template_id, token_id)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "sales order fetched", resp)
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
func (h *handler) UpdateSalesOrders(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID

	data := c.Get("sales_orders").(*orders.SalesOrders)

	//--------------------------concurrency_management_start-----------------------------------------
	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(uint(ID), "status_id", "sales_orders")
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(uint(ID), PreviousStatusId, "status_id", "sales_orders")
	//--------------------------concurrency_management_end-----------------------------------------

	token_id := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.UpdateSalesOrder(query, data, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "sales order updated succesfully", nil)
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
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)

	var query = make(map[string]interface{})

	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id

	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.DeleteSalesOrder(query, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "sales order deleted successfully", map[string]int{"deleted_id": ID})
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
	var product_id = c.QueryParam("product_id")
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = map[string]interface{}{
		"so_id":      ID,
		"product_id": product_id,
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.DeleteSalesOrderLines(query, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", query)
}

// SearchSalesOrders godoc
// @Summary Search SalesOrders
// @Description Search SalesOrders
// @Tags SalesOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   query 	query 	string 	true "query"
// @Success 200 {array} GetAllSalesOrdersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_orders/search [get]
func (h *handler) SearchSalesOrders(c echo.Context) error {

	query := c.QueryParam("query")
	var data []ListSalesOrdersDTO
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.SearchSalesOrders(query, access_template_id, token_id)

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "OK", data)
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
	return res.RespSuccess(c, "email sent successfully", nil)
}

// DownloadSalesOrders godoc
// @Summary Download SalesOrders
// @Description Download SalesOrders
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
// @Router /api/v1/sales_orders/{id}/downloadPdf [post]
func (h *handler) DownloadSalesOrders(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file downloaded successfully", resp_data)
}

// GenerateSalesOrdersPDF godoc
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
func (h *handler) GenerateSalesOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF generated", resp_data)
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
func (h *handler) FavouriteSalesOrders(c echo.Context) (err error) {
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
	_, er := h.service.ViewSalesOrder(q, access_template_id, token_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteSalesOrders(query)
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
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteSalesOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "SalesOrders is Unmarked as Favourite", map[string]string{"id": id})
}

// ----------------------Channel API's --------------------------------------------------------
func (h *handler) ChannelSalesOrderUpsert(c echo.Context) (err error) {

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

	TokenUserID := c.Get("TokenUserID").(string)
	msg, err := h.service.ChannelSalesOrderUpsert(arrayData, TokenUserID)
	if err != nil {
		return res.RespErr(c, err)
	}
	fmt.Println(msg)

	return res.RespSuccess(c, "Multiple records excuted successfully", msg)
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

	productIdString := c.Param("product_id")

	productId, _ := strconv.Atoi(productIdString)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data, err := h.service.GetSalesHistory(uint(productId), page, access_template_id, token_id)

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
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetSalesOrderTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
