package delivery_orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	concurrency_management "fermion/backend_core/controllers/concurrency_management"
	orders_base "fermion/backend_core/controllers/orders/base"
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
	concurrencyService concurrency_management.Service
	base_service       orders_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	concurrencyService := concurrency_management.NewService()
	base_service := orders_base.NewServiceBase()
	return &handler{service, concurrencyService, base_service}
}

// CreateDeliveryOrders godoc
// @Summary      do a CreateDeliveryOrder
// @Description  Create a DeliveryOrder
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   DeliveryOrderCreateRequest body DeliveryOrderRequest true "create a DeliveryOrder"
// @Success      200  {object}  DeliveryOrderCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/create [post]
func (h *handler) CreateDeliveryOrders(c echo.Context) (err error) {
	data := c.Get("delivery_orders").(*DeliveryOrderRequest)
	// lookupcode_query := map[string]interface{}{"id": location.(LocationResponseDTO).LocationType.Id}
	// lookupcode, err := h.coreRepository.GetLookupCodes(lookupcode_query, new(pagination.Paginatevalue))
	// if err != nil || len(lookupcode) == 0 {
	// 	return res.RespError(c, err)
	// }

	//--------------------------concurrency_management_start-----------------------------------------
	var source_documents map[string]interface{}
	jsondata, err := json.Marshal(data.DeliveryOrdersDetails.SourceDocuments)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(jsondata, &source_documents)
	if err != nil {
		return res.RespErr(c, err)
	}

	if source_documents["id"] == nil || source_documents["id"].(float64) <= 0 {
		return res.RespErr(c, errors.New("in source documents id must be present"))
	}

	SourceDocumentType, err := helpers.GetLookupcodeName(int(*data.DeliveryOrdersDetails.Source_document_id))
	if err != nil {
		return res.RespError(c, err)
	}

	recordID := uint(source_documents["id"].(float64))
	modelName := helpers.SourceDocumentTypeWithModelName[SourceDocumentType]
	columnName := helpers.SourceDocumentTypeWithColumnName[SourceDocumentType]

	response, PreviousStatusId, err := h.concurrencyService.CheckConcurrencyStatus(recordID, columnName, modelName)
	if response.Block {
		return res.RespErr(c, errors.New(response.Message))
	}
	defer h.concurrencyService.ReleaseConcurrencyLock(recordID, PreviousStatusId, columnName, modelName)

	//--------------------------concurrency_management_end-----------------------------------------

	token_id := c.Get("TokenUserID").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	access_template_id := c.Get("AccessTemplateId").(string)
	seq_delivery := helpers.GenerateSequence("scrap", token_id, "scrap_orders")
	seq_ref := helpers.GenerateSequence("Ref", token_id, "scrap_orders")

	if data.DeliveryOrdersDetails.AutoCreateDoNumber {
		data.DeliveryOrdersDetails.Delivery_order_number = seq_delivery
	}
	if data.DeliveryOrdersDetails.AutoGenerateReferenceNumber {
		data.DeliveryOrdersDetails.Reference_id = seq_ref
	}
	id, err := h.service.CreateDeliveryOrder(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DO created successfully", map[string]interface{}{"created_id": id})
}

// UpdateDeliveryOrders godoc
// @Summary      do a UpdateDeliveryOrders
// @Description  Update a DeliveryOrder
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   DeliveryOrderUpdateRequest body DeliveryOrderRequest true "Update a DeliveryOrder"
// @Success      200  {object}  DeliveryOrderUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id}/update [post]
func (h *handler) UpdateDeliveryOrders(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("delivery_orders").(*DeliveryOrderRequest)
	token_id := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.UpdateDeliveryOrder(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DO updated successfully", map[string]interface{}{"updated_id": id})
}

// FindDeliveryOrders godoc
// @Summary      Find a DeliveryOrder by id
// @Description  Find a DeliveryOrder by id
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  DeliveryOrderGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id} [get]
func (h *handler) FindDeliveryOrders(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.FindDeliveryOrder(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DO data retrieved successfully", result)
}

// AllDeliveryOrders godoc
// @Summary      get AllDeliveryOrders and filter it by search query
// @Description  a get AllDeliveryOrders and filter it by search query
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters           query     string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  DeliveryOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders [get]
func (h *handler) AllDeliveryOrders(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AllDeliveryOrders(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "DO data retrieved successfully", result, p)
}

// AllDeliveryOrdersDropDown godoc
// @Summary      get AllDeliveryOrdersDropDown and filter it by search query
// @Description  a get AllDeliveryOrdersDropDown and filter it by search query
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters           query     string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  DeliveryOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/dropdown [get]
func (h *handler) AllDeliveryOrdersDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AllDeliveryOrders(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "DO dropdown retrieved successfully", result, p)
}

// DeleteDeliveryOrders godoc
// @Summary      delete a DeleteDeliveryOrders
// @Description  delete a DeleteDeliveryOrders
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  DeliveryOrderDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id}/delete [delete]
func (h *handler) DeleteDeliveryOrders(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.DeleteDeliveryOrder(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DO deleted successfully", map[string]int{"deleted_id": id})
}

// DeleteDeliveryOrderLines godoc
// @Summary      delete a DeleteDeliveryOrderLines
// @Description  delete a DeleteDeliveryOrderLines
// @Tags         DeliveryOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  DeliveryOrderLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/delivery_orders/{id}/delete_products [delete]
func (h *handler) DeleteDeliveryOrderLines(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	product_id := c.QueryParam("product_id")
	prod_id, _ := strconv.Atoi(product_id)
	query := map[string]interface{}{
		"do_id":      uint(id),
		"product_id": uint(prod_id),
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.DeleteDeliveryOrderLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "line items successfully Deleted", query)
}

// BulkCreateDeliveryOrders godoc
// @Summary 	do a BulkCreateDeliveryOrder
// @Description do a BulkCreateDeliveryOrder
// @Tags 		DeliveryOrders
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkDeliveryOrderCreateRequest body DeliveryOrderRequest true "Create a Bulk DeliveryOrder"
// @Success 	200 {object} BulkDeliveryOrderCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/delivery_orders/bulkcreate [post]
func (h *handler) BulkCreateDeliveryOrders(c echo.Context) (err error) {
	data := new([]DeliveryOrderRequest)
	c.Bind(data)
	s := c.Get("TokenUserID").(string)
	t := helpers.ConvertStringToUint(s)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.BulkCreateDeliveryOrder(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DO created successfully", map[string]interface{}{"created_data": data})
}

// SendEmailDeliveryOrders godoc
// @Summary Email DeliveryOrders
// @Description Email DeliveryOrders
// @Tags DeliveryOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/send_email [post]
func (h *handler) SendEmailDeliveryOrders(c echo.Context) error {
	return res.RespSuccess(c, "Email sent successfully", nil)
}

// DownloadDeliveryOrdersPDF godoc
// @Summary Download DeliveryOrders
// @Description Download DeliveryOrders
// @Tags DeliveryOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/download_pdf [post]
func (h *handler) DownloadDeliveryOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

// GenerateDeliveryOrdersPDF godoc
// @Summary Generate DeliveryOrders PDF
// @Description Generate DeliveryOrders PDF
// @Tags DeliveryOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/generate_pdf [post]
func (h *handler) GenerateDeliveryOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF Generated", resp_data)
}

// FavouriteDeliveryOrders godoc
// @Summary FavouriteDeliveryOrders
// @Description FavouriteDeliveryOrders
// @Tags DeliveryOrders
// @Accept  json
// @Produce  json
// @param id path string true "DeliveryOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/favourite [post]
func (h *handler) FavouriteDeliveryOrders(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}

	_, er := h.service.FindDeliveryOrder(uint(ID), token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteDeliveryOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DeliveryOrders is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouriteDeliveryOrders godoc
// @Summary UnFavouriteDeliveryOrders
// @Description UnFavouriteDeliveryOrders
// @Tags DeliveryOrders
// @Accept  json
// @Produce  json
// @param id path string true "DeliveryOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/delivery_orders/{id}/unfavourite [post]
func (h *handler) UnFavouriteDeliveryOrders(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteDeliveryOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "DeliveryOrders is Unmarked as Favourite", map[string]string{"id": id})
}

//-----------------------------------------------------

func (h *handler) GetDeliveryOrderTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetDeliveryOrderTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
