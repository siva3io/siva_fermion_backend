package pick_list

import (
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	inventory_tasks_base "fermion/backend_core/controllers/inventory_tasks/base"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_tasks"
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
	base_service inventory_tasks_base.ServiceBase
}

var PickListHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if PickListHandler != nil {
		return PickListHandler
	}
	service := NewService()
	base_service := inventory_tasks_base.NewServiceBase()
	PickListHandler = &handler{service, base_service}
	return PickListHandler
}

// FavouritePickListView godoc
// @Summary Get all favourite PickList
// @Description Get all favourite PickList
// @Tags PickList
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} PicklistGetAllResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pick_list/favourite_list [get]
func (h *handler) FavouritePickListView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrPicList(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavPicList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Pick list Retrieved Successfully", result, p)
}

// CreatePicklist godoc
// @Summary      do a CreatePicklist
// @Description  Create a Picklist
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   PicklistCreateRequest body PicklistRequest true "create a Picklist"
// @Success      200  {object}  PicklistCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/create [post]
func (h *handler) CreatePicklistEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("pick_list"),
	}
	eda.Produce(eda.CREATE_PICK_LIST, request_payload)
	return res.RespSuccess(c, "Pick List Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreatePicklist(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(inventory_tasks.PickList)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.CreatePickList(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_PICK_LIST_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_PICK_LIST_ACK, responseMessage)
	// cache implementation
	UpdatePickListInCache(edaMetaData)
}

// func (h *handler) CreatePickList(c echo.Context) (err error) {
// 	data := c.Get("pick_list").(*PicklistRequest)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	user_id, _ := strconv.Atoi(token_id)
// 	t := uint(user_id)
// 	data.CreatedByID = &t
// 	id, err := h.service.CreatePickList(data, token_id, access_template_id)
// 	if err != nil {
// 		return res.RespErr(c, err)
// 	}

// 	// cache implementation
// 	UpdatePickListInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Picklist created successfully", map[string]interface{}{"created_id": id})
// }

// BulkCreatePicklist godoc
// @Summary 	do a BulkCreatePicklist
// @Description do a BulkCreatePicklist
// @Tags 		PickList
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkPicklistCreateRequest body PicklistRequest true "Create a Bulk Picklist"
// @Success 	200 {object} BulkPicklistCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/pick_list/bulk_create [post]
func (h *handler) BulkCreatePickList(c echo.Context) (err error) {
	data := new([]PicklistRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	var metaData core.MetaData
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreatePickList(metaData, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Picklist multiple records created successfully", data)
}

// UpdatePicklist godoc
// @Summary      do a UpdatePicklist
// @Description  Update a Picklist
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   PicklistUpdateRequest body PicklistRequest true "Update a Picklist"
// @Success      200  {object}  PicklistUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/edit [put]
func (h *handler) UpdatePickListEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	pickListId := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         pickListId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("pick_list"),
	}
	eda.Produce(eda.UPDATE_PICK_LIST, request_payload)
	return res.RespSuccess(c, "Pick List Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdatePickList(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(inventory_tasks.PickList)
	helpers.JsonMarshaller(data, updatePayload)
	err := h.service.UpdatePickList(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_PICK_LIST_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_PICK_LIST_ACK, responseMessage)

	UpdatePickListInCache(edaMetaData)
}

// func (h *handler) UpdatePickList(c echo.Context) (err error) {
// 	ID := c.Param("id")
// 	id, _ := strconv.Atoi(ID)
// 	data := c.Get("pick_list").(*PicklistRequest)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
// 	err = h.service.UpdatePickList(uint(id), data, token_id, access_template_id)
// 	if err != nil {
// 		return res.RespError(c, err)
// 	}

// 	// cache implementation
// 	UpdatePickListInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Picklist updated successfully", map[string]interface{}{"updated_id": id})
// }

// GetPicklist godoc
// @Summary      Find a Picklist by id
// @Description  Find a Picklist by id
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  PicklistGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id} [get]
func (h *handler) GetPickList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	pickListId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": pickListId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetPickList(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response PicklistRequest
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Picklist fetched successfully", response)
}

// GetAllPicklist godoc
// @Summary      get all GetAllPicklist with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllPicklist with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		   query	 string false "filters"
// @Param        per_page 		   query     int    false "per_page"
// @Param  	     page_no           query     int    false "page_no"
// @Param        sort			   query     string false "sort"
// @Success      200  {object}  PicklistGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list [get]
func (h *handler) GetAllPickList(c echo.Context) (err error) {
	//per_page, page_no, sort
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []PicklistRequest

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetPickListFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllPickList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Picklist data retrieved Successfully", response, p)
}

// DeletePicklist godoc
// @Summary      delete a DeletePicklist
// @Description  delete a DeletePicklist
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  PicklistDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/delete [delete]
func (h *handler) DeletePickList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	priceListId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         priceListId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeletePickList(metaData)
	if err != nil {
		return res.RespError(c, err)
	}

	// cache implementation
	UpdatePickListInCache(metaData)

	return res.RespSuccess(c, "Picklist deleted successfully", map[string]interface{}{"deleted_id": priceListId})
}

// DeletePicklistLines godoc
// @Summary      delete a DeletePicklistLines
// @Description  delete a DeletePicklistLines
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        pick_list_id   path      int    true  "pick_list_id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  PicklistLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/delete_products [delete]
func (h *handler) DeletePickListLines(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	pickListId := c.Param("id")
	product_id := c.QueryParam("product_id")
	metaData.Query = map[string]interface{}{
		"id":         pickListId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
		"product_id": product_id,
	}

	err = h.service.DeletePickListLines(metaData)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Picklist Line Items successfully deleted", pickListId)
}

// SendMailPickList godoc
// @Summary      Send Email to specific recevier for Id
// @Description  Send Email to specific recevier for Id
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Param        receiver_email query     string true "receiver_email"
// @Success      200  {object}  SendMailPickListResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/sendemail [get]
func (h *handler) SendMailPickList(c echo.Context) (err error) {

	q := new(SendMailPickList)
	c.Bind(q)
	// err = h.service.SendMailPickList(q)
	// if err != nil {
	// 	return res.RespError(c, err)
	// }

	return res.RespSuccess(c, "Email Sent Successfully", q)
}

// DownloadPdfPickList godoc
// @Summary      Download PDF
// @Description  Download PDF
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Success      200  {object}  DownloadPdfPickListResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/{id}/printpdf [get]
func (h *handler) DownloadPdfPickList(c echo.Context) (err error) {
	return res.RespSuccess(c, "Pdf Generated Successfully", "https://drive.google.com/file/d/1O46B3lgliUGakL4weBjrovCHphTzR5Np/view?usp=sharing")
}

// FavoritePickList godoc
// @Summary FavoritePickList
// @Description FavoritePickList
// @Tags PickList
// @Accept  json
// @Produce  json
// @param id path string true "PickList ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pick_list/{id}/favourite [post]
func (h *handler) FavouritePickList(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.FavouritePickList(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Pick List is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouritePickList godoc
// @Summary UnFavoritePickList
// @Description UnFavoritePickList
// @Tags PickList
// @Accept  json
// @Produce  json
// @param id path string true "PickList ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/pick_list/{id}/unfavourite [post]
func (h *handler) UnFavouritePickList(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouritePickList(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Pick List is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetAllPickListDropDown godoc
// @Summary      get all GetAllPickListDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Description  get all GetAllPickListDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Tags         PickList
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters		   query	 string false "filters"
// @Param        per_page 		   query     int    false "per_page"
// @Param  	     page_no           query     int    false "page_no"
// @Param        sort			   query     string false "sort"
// @Success      200  {object}  PicklistGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/pick_list/dropdown [get]
func (h *handler) GetAllPickListDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []PicklistRequest

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetPickListFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllPickList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Picklist data retrieved Successfully", response, p)
}
