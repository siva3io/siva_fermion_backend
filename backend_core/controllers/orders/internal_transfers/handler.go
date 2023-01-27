package internal_transfers

import (
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"

	concurrency_management "fermion/backend_core/controllers/concurrency_management"
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
	service            Service
	base_service       orders_base.ServiceBase
	concurrencyService concurrency_management.Service
}

var InternalTransfersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if InternalTransfersHandler != nil {
		return InternalTransfersHandler
	}
	InternalTransfersHandler = &handler{
		NewService(),
		orders_base.NewServiceBase(),
		concurrency_management.NewService(),
	}
	return InternalTransfersHandler
}

// CreateInternalTransfers godoc
// @Summary Create InternalTransfers
// @Description Create InternalTransfers
// @Tags InternalTransfers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body InternalTransfersDTO true "Internal Transfers  Request Body"
// @Success 200 {object} InternalTransfersCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/create [post]
func (h *handler) CreateInternalTransfersEvent(c echo.Context) error {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("internal_transfers"),
	}
	eda.Produce(eda.CREATE_INTERNAL_TRANSFERS, requestPayload)
	return res.RespSuccess(c, "IST Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateInternalTransfers(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(orders.InternalTransfers)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.CreateInternalTransfers(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_INTERNAL_TRANSFERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateInternalTransferInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": requestPayload.ID,
	}
	// eda.Produce(eda.CREATE_INTERNAL_TRANSFERS_ACK, *responseMessage)
}

// GetListInternalTransfers godoc
// @Summary Get all InternalTransfers list
// @Description Get all InternalTransfers list
// @Tags InternalTransfers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllInternalTransfersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers [get]
func (h *handler) ListInternalTransfers(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetInternalTransferFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.ListInternalTransfers(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ListInternalTransfersDTO
	helpers.JsonMarshaller(data, &dtoData)
	return res.RespSuccessInfo(c, "list of internal transfers fetched", dtoData, page)
}

// ListInternalTransfersDropDown godoc
// @Summary Get all InternalTransfers DropDown
// @Description Get all InternalTransfers DropDown
// @Tags InternalTransfers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllInternalTransfersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/dropdown [get]
func (h *handler) ListInternalTransfersDropDown(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	helpers.AddMandatoryFilters(page, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	if *page == pagination.BasePaginatevalue {
		cacheResponse, *page = GetInternalTransferFromCache(fmt.Sprint(metaData.TokenUserId))
		if cacheResponse != nil {
			return res.RespSuccessInfo(c, "data retrieved successfully", cacheResponse, page)
		}
	}

	data, err := h.service.ListInternalTransfers(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	var dtoData []ListInternalTransfersDTO
	helpers.JsonMarshaller(data, &dtoData)
	return res.RespSuccessInfo(c, "list of internal transfers fetched", dtoData, page)
}

// ViewInternalTransfers godoc
// @Summary View InternalTransfers
// @Summary View InternalTransfers
// @Description View InternalTransfers
// @Tags InternalTransfers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ist_id"
// @Success 200 {object} InternalTransfersGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/{id} [get]
func (h *handler) ViewInternalTransfers(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": c.Param("id"),
		// "company_id": metaData.CompanyId,
	}
	data, err := h.service.ViewInternalTransfers(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "internal transfers fetched", data)
}

// UpdateInternalTransfers godoc
// @Summary Update InternalTransfers
// @Description Update InternalTransfers
// @Tags InternalTransfers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ist_id"
// @param RequestBody body InternalTransfersDTO true "Internal Transfers  Request Body"
// @Success 200 {object} InternalTransfersUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/{id}/update [post]
func (h *handler) UpdateInternalTransfersEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)
	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	requestPayload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("internal_transfers"),
	}

	//--------------------------concurrency_management_start-----------------------------------------
	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(uint(id), "status_id", "internal_transfers")
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(uint(id), PreviousStatusId, "status_id", "internal_transfers")
	//--------------------------concurrency_management_end-----------------------------------------

	eda.Produce(eda.UPDATE_INTERNAL_TRANSFERS, requestPayload)
	return res.RespSuccess(c, "IST Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateInternalTransfers(request map[string]interface{}) {

	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	requestPayload := new(orders.InternalTransfers)
	helpers.JsonMarshaller(request["data"], requestPayload)

	err := h.service.UpdateInternalTransfers(edaMetaData, requestPayload)
	responseMessage := new(eda.ConsumerResponse)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_INTERNAL_TRANSFERS_ACK, *responseMessage)
		return
	}
	// cache implementation
	UpdateInternalTransferInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": requestPayload.ID,
	}
	// eda.Produce(eda.UPDATE_INTERNAL_TRANSFERS_ACK, *responseMessage)
}

// DeleteInternalTransfers godoc
// @Summary Delete InternalTransfers
// @Description Delete InternalTransfers
// @Tags InternalTransfers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ist_id"
// @Success 200 {object} InternalTransfersDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/{id}/delete [delete]
func (h *handler) DeleteInternalTransfers(c echo.Context) error {
	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err := h.service.DeleteInternalTransfers(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateInternalTransferInCache(metaData)

	return res.RespSuccess(c, "internal transfer order deleted successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteInternalTransfersLines godoc
// @Summary Delete InternalTransfers Orderlines
// @Description Delete InternalTransfers Orderlines
// @Tags InternalTransfers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ist_id"
// @param product_id query string true "product_id"
// @Success 200 {object} InternalTransferOrderLinesDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/order_lines/{id}/delete [delete]
func (h *handler) DeleteInternalTransfersLines(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"ist_id":     id,
		"product_id": c.QueryParam("product_id"),
	}
	err := h.service.DeleteIstLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", metaData.Query)
}

// DownloadInternalTransfers godoc
// @Summary Download InternalTransfers
// @Description Download InternalTransfers
// @Tags InternalTransfers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ist_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/{id}/downloadPdf [post]
func (h *handler) DownloadInternalTransfers(c echo.Context) error {

	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}

	return res.RespSuccess(c, "file downloaded successfully", resp_data)
}

// FavouriteInternalTransfers godoc
// @Summary FavouriteInternalTransfers
// @Description FavouriteInternalTransfers
// @Tags InternalTransfers
// @Accept  json
// @Produce  json
// @param id path string true "InternalTransfers ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/{id}/favourite [post]
func (h *handler) FavouriteInternalTransfers(c echo.Context) error {

	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.Query = map[string]interface{}{
		"id": id,
	}
	_, err := h.service.ViewInternalTransfers(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.FavouriteInternalTransfers(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "InternalTransfers is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouriteInternalTransfers godoc
// @Summary UnFavouriteInternalTransfers
// @Description UnFavouriteInternalTransfers
// @Tags InternalTransfers
// @Accept  json
// @Produce  json
// @param id path string true "InternalTransfers ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/{id}/unfavourite [post]
func (h *handler) UnFavouriteInternalTransfers(c echo.Context) error {
	id := c.Param("id")
	metaData := c.Get("MetaData").(core.MetaData)

	query := map[string]interface{}{
		"id":      id,
		"user_id": metaData.TokenUserId,
	}
	err := h.base_service.UnFavouriteInternalTransfers(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "InternalTransfers is Unmarked as Favourite", map[string]string{"id": id})
}

// GetInternalTransferTab godoc
// @Summary GetInternalTransferTab
// @Summary GetInternalTransferTab
// @Description GetInternalTransferTab
// @Tags InternalTransfers
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
// @Router /api/v1/internal_transfers/{id}/filter_module/{tab} [get]
func (h *handler) GetInternalTransfersTab(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	page := new(pagination.Paginatevalue)
	err := c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	metaData.AdditionalFields = map[string]interface{}{
		"id":  c.Param("id"),
		"tab": c.Param("tab"),
	}

	data, err := h.service.GetInternalTransfersTab(metaData, page)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
