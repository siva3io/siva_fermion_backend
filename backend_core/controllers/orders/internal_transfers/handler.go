package internal_transfers

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

func NewHandler() *handler {
	service := NewService()
	base_service := orders_base.NewServiceBase()
	concurrencyService := concurrency_management.NewService()
	return &handler{service, base_service, concurrencyService}
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
func (h *handler) CreateInternalTransfers(c echo.Context) error {

	input_data := c.Get("internal_transfers").(*orders.InternalTransfers)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	input_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err := h.service.CreateInternalTransfers(input_data, access_template_id, token_id)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "IST Created", map[string]interface{}{"internal_transfer_id": input_data.ID})
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

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListInternalTransfersDTO

	resp, err := h.service.ListInternalTransfers(page, access_template_id, token_id, "LIST")

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "list of internal transfers fetched", data, page)
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

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data []ListInternalTransfersDTO

	resp, err := h.service.ListInternalTransfers(page, access_template_id, token_id, "DROPDOWN_LIST")

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "Dropdown of internal transfers fetched", data, page)
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

	var query = make(map[string]interface{}, 0)

	ID := c.Param("id")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	id, _ := strconv.Atoi(ID)

	query["id"] = int(id)

	fmt.Println("query", query)

	resp, err := h.service.ViewInternalTransfers(query, access_template_id, token_id)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "internal transfers fetched", resp)
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
func (h *handler) UpdateInternalTransfers(c echo.Context) (err error) {
	var id = c.Param("id")

	var query = make(map[string]interface{}, 0)

	ID, _ := strconv.Atoi(id)

	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data := c.Get("internal_transfers").(*orders.InternalTransfers)

	//--------------------------concurrency_management_start-----------------------------------------
	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(uint(ID), "status_id", "internal_transfers")
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(uint(ID), PreviousStatusId, "status_id", "internal_transfers")
	//--------------------------concurrency_management_end-----------------------------------------

	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateInternalTransfers(query, data, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "internal transfers updated succesfully", nil)
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
func (h *handler) DeleteInternalTransfers(c echo.Context) (err error) {
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)

	var query = make(map[string]interface{})
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteInternalTransfers(query, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "internal transfer order deleted successfully", map[string]int{"deleted_id": ID})
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
func (h *handler) DeleteInternalTransfersLines(c echo.Context) (err error) {
	var product_id = c.QueryParam("product_id")
	var id = c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = map[string]interface{}{
		"ist_id":     ID,
		"product_id": product_id,
	}
	err = h.service.DeleteIstLines(query, access_template_id, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", query)
}

// SearchInternalTransfers godoc
// @Summary Search InternalTransfers
// @Description Search InternalTransfers
// @Tags InternalTransfers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   query 	query 	string 	true "query"
// @Success 200 {array} GetAllInternalTransfersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/internal_transfers/search [get]
func (h *handler) SearchInternalTransfers(c echo.Context) error {

	query := c.QueryParam("query")

	var data []ListInternalTransfersDTO
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	resp, err := h.service.SearchInternalTransfers(query, access_template_id, token_id)

	value, _ := json.Marshal(resp)

	_ = json.Unmarshal(value, &data)

	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "OK", data)
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
func (h *handler) FavouriteInternalTransfers(c echo.Context) (err error) {
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
	_, er := h.service.ViewInternalTransfers(q, access_template_id, token_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
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
func (h *handler) UnFavouriteInternalTransfers(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteInternalTransfers(query)
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
func (h *handler) GetInternalTransfersTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetInternalTransfersTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
