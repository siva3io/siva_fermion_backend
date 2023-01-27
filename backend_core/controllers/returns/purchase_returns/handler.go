package purchase_returns

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	concurrency_management "fermion/backend_core/controllers/concurrency_management"
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
	service            Service
	base_service       returns_base.ServiceBase
	concurrencyService concurrency_management.Service
}

var PurchaseReturnsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if PurchaseReturnsHandler != nil {
		return PurchaseReturnsHandler
	}
	service := NewService()
	base_service := returns_base.NewServiceBase()
	concurrencyService := concurrency_management.NewService()
	PurchaseReturnsHandler = &handler{service, base_service, concurrencyService}
	return PurchaseReturnsHandler
}

// CreatePurchaseReturns godoc
// @Summary Create PurchaseReturns
// @Description Create PurchaseReturns
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body PurchaseReturnsDTO true "Purchase Returns Request Body"
// @Success 200 {object} PurchaseReturnsCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/create [post]
func (h *handler) CreatePurchaseReturnsEvent(c echo.Context) error {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("purchase_returns"),
	}
	eda.Produce(eda.CREATE_PURCHASE_RETURNS, requestPayload)
	return res.RespSuccess(c, "Purchase Return Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreatePurchaseReturns(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(PurchaseReturnsDTO)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.CreatePurchaseReturn(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_PURCHASE_RETURNS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdatePurchaseReturnInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": requestPayload.ID,
	}
	// eda.Produce(eda.CREATE_PURCHASE_RETURNS_ACK, *responseMessage)
}

// func (h *handler) CreatePurchaseReturns(c echo.Context) error {

// 	input_data := c.Get("purchase_returns").(*returns.PurchaseReturns)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	input_data.CreatedByID = helpers.ConvertStringToUint(token_id)
// 	err := h.service.CreatePurchaseReturn(input_data, access_template_id, token_id)

// 	if err != nil {
// 		return res.RespError(c, err)
// 	}

// 	// cache implementation
// 	UpdatePurchaseReturnInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Purchase Return Created", map[string]interface{}{"purchase_return_id": input_data.ID})
// }

// ListPurchaseReturnsDropdown godoc
// @Summary Get all PurchaseReturns list
// @Description Get all PurchaseReturns list
// @Tags PurchaseReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllPurchaseReturnsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns [get]
func (h *handler) ListPurchaseReturns(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetPurchaseReturnFromCache(fmt.Sprint(metaData.TokenUserId))
	}

	if cacheResponse != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
	}

	dataInterface, err := h.service.ListPurchaseReturns(metaData, page)
	data := dataInterface.([]returns.PurchaseReturns)
	if err != nil {
		return res.RespError(c, err)
	}

	var dtoData []ListPurchaseReturnsDTO
	helpers.JsonMarshaller(data, &dtoData)
	for ind, val := range data {
		dtoData[ind].Amount = val.PrPaymentDetails.TotalAmount
	}
	return res.RespSuccessInfo(c, "list of purchase returns fetched", data, page)
}

// ViewPurchaseReturns godoc
// @Summary View PurchaseReturns
// @Description View PurchaseReturns
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "pr_id"
// @Success 200 {object} PurchaseReturnsGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id} [get]
func (h *handler) ViewPurchaseReturns(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	prId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": prId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.ViewPurchaseReturn(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response returns.PurchaseReturns
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Purchase Returns Details retrieved succesfully", response)
}

// UpdatePurchaseReturns godoc
// @Summary Update PurchaseReturns
// @Description Update PurchaseReturns
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "pr_id"
// @param RequestBody body PurchaseReturnsDTO true "Purchase Returns Request Body"
// @Success 200 {object} PurchaseReturnsUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id}/update [post]
func (h *handler) UpdatePurchaseReturnsEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("purchase_returns"),
	}
	//--------------------------concurrency_management_start-----------------------------------------
	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(uint(id), "status_id", "purchase_returns")
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(uint(id), PreviousStatusId, "status_id", "purchase_returns")
	//--------------------------concurrency_management_end-----------------------------------------

	eda.Produce(eda.UPDATE_PURCHASE_RETURNS, requestPayload)
	return res.RespSuccess(c, "Purchase Returns Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdatePurchaseReturns(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(PurchaseReturnsDTO)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.UpdatePurchaseReturn(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_PURCHASE_RETURNS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdatePurchaseReturnInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_PURCHASE_RETURNS_ACK, *responseMessage)
}

// func (h *handler) UpdatePurchaseReturns(c echo.Context) (err error) {
// 	var id = c.Param("id")
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	var query = make(map[string]interface{}, 0)

// 	ID, _ := strconv.Atoi(id)

// 	query["id"] = ID
// 	data := c.Get("purchase_returns").(*returns.PurchaseReturns)

// 	//--------------------------concurrency_management_start-----------------------------------------
// 	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(uint(ID), "status_id", "purchase_returns")
// 	if response.Block {
// 		return res.RespErr(c, errors.New(response.Message))
// 	}
// 	defer h.concurrencyService.ReleaseConcurrencyLock(uint(ID), PreviousStatusId, "status_id", "purchase_returns")
// 	//--------------------------concurrency_management_end-----------------------------------------

// 	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
// 	err = h.service.UpdatePurchaseReturn(query, data, access_template_id, token_id)
// 	if err != nil {
// 		return res.RespError(c, err)
// 	}

// 	// cache implementation
// 	UpdatePurchaseReturnInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "purchase return updated succesfully", nil)
// }

// DeletePurchaseReturns godoc
// @Summary Delete PurchaseReturns
// @Description Delete PurchaseReturns
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "pr_id"
// @Success 200 {object} PurchaseReturnsDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id}/delete [delete]
func (h *handler) DeletePurchaseReturns(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeletePurchaseReturn(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdatePurchaseReturnInCache(metaData)
	return res.RespSuccess(c, "purchase return deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeletePurchaseReturnLines godoc
// @Summary Delete PurchaseReturnLines
// @Description Delete PurchaseReturnLines
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "pr_id"
// @param product_id query string true "product_id"
// @Success 200 {object} PurchaseReturnLinesDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/return_lines/{id}/delete [delete]
func (h *handler) DeletePurchaseReturnLines(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"pr_id":      id,
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeletePurchaseReturnLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "return line deleted successfully", metaData.Query)
}

// SearchPurchaseReturns godoc
// @Summary Search PurchaseReturns
// @Description Search PurchaseReturns
// @Tags PurchaseReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   query 	query 	string 	true "query"
// @Success 200 {array} GetAllPurchaseReturnsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/search [get]
func (h *handler) SearchPurchaseReturns(c echo.Context) error {

	query := c.QueryParam("query")

	var data []ListPurchaseReturnsDTO
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.SearchPurchaseReturns(query, access_template_id, token_id)

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "OK", data)
}

// EmailPurchaseReturns godoc
// @Summary Email PurchaseReturns
// @Description Email PurchaseReturns
// @Tags PurchaseReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "pr_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id}/sendEmail [post]
func (h *handler) EmailPurchaseReturns(c echo.Context) error {
	return res.RespSuccess(c, "email sent successfully", nil)
}

// DownloadPurchaseReturns godoc
// @Summary Download PurchaseReturns
// @Description Download PurchaseReturns
// @Tags PurchaseReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "pr_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id}/downloadPdf [post]
func (h *handler) DownloadPurchaseReturns(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file downloaded successfully", resp_data)
}

// GeneratePurchaseReturnsPDF godoc
// @Summary Generate PDF PurchaseReturns
// @Description Generate PDF PurchaseReturns
// @Tags PurchaseReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "pr_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id}/generatePdf [post]
func (h *handler) GeneratePurchaseReturnsPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF generated", resp_data)
}

// FavouritePurchaseReturns godoc
// @Summary FavouritePurchaseReturns
// @Description FavouritePurchaseReturns
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @param id path string true "PurchaseReturns ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id}/favourite [post]
func (h *handler) FavouritePurchaseReturns(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	_, err = h.service.ViewPurchaseReturn(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}
	err = h.base_service.FavouritePurchaseReturns(metaData.Query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseReturn is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouritePurchaseReturns godoc
// @Summary UnFavouritePurchaseReturns
// @Description UnFavouritePurchaseReturns
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @param id path string true "PurchaseReturns ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{id}/unfavourite [post]
func (h *handler) UnFavouritePurchaseReturns(c echo.Context) (err error) {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouritePurchaseReturns(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "PurchaseReturn is Unmarked as Favourite", map[string]string{"id": id})
}

// GetPurchaseReturnsHistory godoc
// @Summary GetPurchaseReturnsHistory
// @Summary GetPurchaseReturnsHistory
// @Description GetPurchaseReturnsHistory
// @Tags PurchaseReturns
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param product_id path string true "product_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/{product_id}/history [get]
func (h *handler) GetPurchaseReturnsHistory(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)

	productIdString := c.Param("product_id")
	productId, _ := strconv.Atoi(productIdString)

	metaData.Query = map[string]interface{}{
		"product_id": productId,
	}
	data, err := h.service.GetPurchaseReturnsHistory(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}

// ListPurchaseReturnsDropdown godoc
// @Summary Get all PurchaseReturns list
// @Description Get all PurchaseReturns list
// @Tags PurchaseReturns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllPurchaseReturnsResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/purchase_returns/dropdown [get]
func (h *handler) ListPurchaseReturnsDropdown(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetPurchaseReturnFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.ListPurchaseReturns(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ListPurchaseReturnsDTO
	helpers.JsonMarshaller(data, &dtoData)

	return res.RespSuccessInfo(c, "dropdown list of purchase returns fetched", dtoData, page)
}

// GetPurchaseReturnTab godoc
// @Summary GetPurchaseReturnTab
// @Summary GetPurchaseReturnTab
// @Description GetPurchaseReturnTab
// @Tags PurchaseReturns
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
// @Router /api/v1/purchase_returns/{id}/filter_module/{tab} [get]
func (h *handler) GetPurchaseReturnTab(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id":  c.Param("id"),
		"tab": c.Param("tab"),
	}

	data, err := h.service.GetPurchaseReturnTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "success", data, page)
}
