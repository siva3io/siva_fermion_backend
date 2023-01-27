package sales_returns

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	returns_base "fermion/backend_core/controllers/returns/base"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
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
	base_service returns_base.ServiceBase
}

var SalesReturnsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if SalesReturnsHandler != nil {
		return SalesReturnsHandler
	}
	service := NewService()
	base_service := returns_base.NewServiceBase()
	SalesReturnsHandler = &handler{service, base_service}
	return SalesReturnsHandler
}

// CreateSalesReturns godoc
// @Summary Create SalesReturns
// @Description Create SalesReturns
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body SalesReturnsDTO true "Sales Returns Request Body"
// @Success 200 {object} SalesReturnsCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/create [post]
func (h *handler) CreateSalesReturnsEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("sales_returns"),
	}
	eda.Produce(eda.CREATE_SALES_RETURNS, requestPayload)
	return res.RespSuccess(c, "Sales Return Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateSalesReturns(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(returns.SalesReturns)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.CreateSalesReturn(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SALES_RETURNS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateSalesReturnsInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": requestPayload.ID,
	}
	// eda.Produce(eda.CREATE_SALES_RETURNS_ACK, *responseMessage)
}

// func (h *handler) CreateSalesReturns(c echo.Context) error {

// 	input_data := c.Get("sales_returns").(*returns.SalesReturns)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	input_data.CreatedByID = helpers.ConvertStringToUint(token_id)
// 	err := h.service.CreateSalesReturn(input_data, access_template_id, token_id)

// 	if err != nil {
// 		return res.RespError(c, err)
// 	}

// 	// cache implementation
// 	UpdateSalesReturnsInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Sales Return Created", map[string]interface{}{"sales_return_id": input_data.ID})
// }

// GetListSalesReturns godoc
// @Summary Get all SalesReturns list
// @Description Get all SalesReturns list
// @Tags SalesReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllSalesReturnsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns [get]
func (h *handler) ListSalesReturns(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetSalesReturnsFromCache(fmt.Sprint(metaData.TokenUserId))
	}

	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
	}

	dataInterface, err := h.service.ListSalesReturns(metaData, page)
	data := dataInterface.([]returns.SalesReturns)
	if err != nil {
		return res.RespError(c, err)
	}

	var dtoData []ListSalesReturnsDTO
	helpers.JsonMarshaller(data, &dtoData)
	for ind, val := range data {
		dtoData[ind].Amount = val.SrPaymentDetails.TotalAmount
	}
	return res.RespSuccessInfo(c, "list of sales returns fetched", data, page)
}

// ListSalesReturnsDropDown godoc
// @Summary Get all SalesReturns dropdown
// @Description Get all SalesReturns dropdown
// @Tags SalesReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllSalesReturnsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/dropdown [get]
func (h *handler) ListSalesReturnsDropDown(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetSalesReturnsFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.ListSalesReturns(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ListSalesReturnsDTO
	helpers.JsonMarshaller(data, &dtoData)

	return res.RespSuccessInfo(c, "dropdown list of sales returns fetched", dtoData, page)
}

// ViewSalesReturns godoc
// @Summary View SalesReturns
// @Description View SalesReturns
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "sr_id"
// @Success 200 {object} SalesReturnsGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id} [get]
func (h *handler) ViewSalesReturns(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	srId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": srId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.ViewSalesReturn(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response returns.SalesReturns
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Sales Returns Details retrieved succesfully", response)
}

// UpdateSalesReturns godoc
// @Summary Update SalesReturns
// @Description Update SalesReturns
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "sr_id"
// @param RequestBody body SalesReturnsDTO true "Sales Returns Request Body"
// @Success 200 {object} SalesReturnsUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id}/update [post]
func (h *handler) UpdateSalesReturnsEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("sales_returns"),
	}

	eda.Produce(eda.UPDATE_SALES_RETURNS, requestPayload)
	return res.RespSuccess(c, "Sales Returns Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateSalesReturns(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(returns.SalesReturns)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.UpdateSalesReturn(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_SALES_RETURNS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateSalesReturnsInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_SALES_RETURNS_ACK, *responseMessage)
}

// func (h *handler) UpdateSalesReturns(c echo.Context) (err error) {
// 	var id = c.Param("id")

// 	var query = make(map[string]interface{}, 0)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	ID, _ := strconv.Atoi(id)

// 	query["id"] = ID

// 	data := c.Get("sales_returns").(*returns.SalesReturns)
// 	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
// 	err = h.service.UpdateSalesReturn(query, data, access_template_id, token_id)
// 	if err != nil {
// 		return res.RespError(c, err)
// 	}

// 	// cache implementation
// 	UpdateSalesReturnsInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "sales return updated succesfully", nil)
// }

// DeleteSalesReturns godoc
// @Summary Delete SalesReturns
// @Description Delete SalesReturns
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "sr_id"
// @Success 200 {object} SalesReturnsDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id}/delete [delete]
func (h *handler) DeleteSalesReturns(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteSalesReturn(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateSalesReturnsInCache(metaData)
	return res.RespSuccess(c, "sales return deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteSalesReturnLines godoc
// @Summary Delete SalesReturnLines
// @Description Delete SalesReturnLines
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "sr_id"
// @param product_id query string true "product_id"
// @Success 200 {object} SalesReturnLinesDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/return_lines/{id}/delete [delete]
func (h *handler) DeleteSalesReturnLines(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"sr_id":      id,
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeleteSalesReturnLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "return line deleted successfully", metaData.Query)
}

// SearchSalesReturns godoc
// @Summary Search SalesReturns
// @Description Search SalesReturns
// @Tags SalesReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   query 	query 	string 	true "query"
// @Success 200 {array} GetAllSalesReturnsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/search [get]
func (h *handler) SearchSalesReturns(c echo.Context) error {

	query := c.QueryParam("query")

	var data []ListSalesReturnsDTO
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.SearchSalesReturns(query, access_template_id, token_id)

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "OK", data)
}

// EmailSalesReturns godoc
// @Summary Email SalesReturns
// @Description Email SalesReturns
// @Tags SalesReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "sr_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id}/sendEmail [post]
func (h *handler) EmailSalesReturns(c echo.Context) error {
	return res.RespSuccess(c, "email sent successfully", nil)
}

// DownloadSalesReturns godoc
// @Summary Download SalesReturns
// @Description Download SalesReturns
// @Tags SalesReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "sr_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id}/downloadPdf [post]
func (h *handler) DownloadSalesReturns(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file downloaded successfully", resp_data)

}

// GenerateSalesReturnsPDF godoc
// @Summary Generate PDF SalesReturns
// @Description Generate PDF SalesReturns
// @Tags SalesReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "sr_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id}/generatePdf [post]
func (h *handler) GenerateSalesReturnsPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF generated", resp_data)
}

// FavouriteSalesReturns godoc
// @Summary FavouriteSalesReturns
// @Description FavouriteSalesReturns
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @param id path string true "SalesReturns ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id}/favourite [post]
func (h *handler) FavouriteSalesReturns(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	_, err = h.service.ViewSalesReturn(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	err = h.base_service.FavouriteSalesReturns(metaData.Query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "SalesReturn is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouriteSalesReturns godoc
// @Summary UnFavouriteSalesReturns
// @Description UnFavouriteSalesReturns
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @param id path string true "SalesReturns ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{id}/unfavourite [post]
func (h *handler) UnFavouriteSalesReturns(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouriteSalesReturns(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "SalesReturn is Unmarked as Favourite", map[string]string{"id": id})
}

// GetSalesReturnsHistory godoc
// @Summary GetSalesReturnsHistory
// @Summary GetSalesReturnsHistory
// @Description GetSalesReturnsHistory
// @Tags SalesReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param product_id path string true "product_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/sales_returns/{product_id}/history [get]
func (h *handler) GetSalesReturnsHistory(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)

	productIdString := c.Param("product_id")
	productId, _ := strconv.Atoi(productIdString)

	metaData.Query = map[string]interface{}{
		"product_id": productId,
	}
	data, err := h.service.GetSalesReturnsHistory(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

// GetSalesReturnTab godoc
// @Summary GetSalesReturnTab
// @Summary GetSalesReturnTab
// @Description GetSalesReturnTab
// @Tags SalesReturns
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
// @Router /api/v1/sales_returns/{id}/filter_module/{tab} [get]
func (h *handler) GetSalesReturnTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id":  c.Param("id"),
		"tab": c.Param("tab"),
	}

	data, err := h.service.GetSalesReturnTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "success", data, page)
}

func (h *handler) TrackSalesReturn(c echo.Context) (err error) {
	var order_id = c.QueryParam("sales_return_number")

	status_history := h.service.TrackSalesReturn(order_id)

	return res.RespSuccess(c, "sales order history successfully", status_history)
}
