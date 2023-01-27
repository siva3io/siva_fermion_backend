package cycle_count

import (
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	inventory_tasks_base "fermion/backend_core/controllers/inventory_tasks/base"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_tasks"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	mdmrepo "fermion/backend_core/internal/repository/mdm"
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
	service             Service
	base_service        inventory_tasks_base.ServiceBase
	locationsRepository mdmrepo.Locations
}

var CycleCountHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if CycleCountHandler != nil {
		return CycleCountHandler
	}
	service := NewService()
	base_service := inventory_tasks_base.NewServiceBase()
	locationsRepository := mdmrepo.NewLocations()
	CycleCountHandler = &handler{service, base_service, locationsRepository}
	return CycleCountHandler
}

// GetCountCapacity godoc
// @Summary 	do a GetCountCapacity
// @Description do a GetCountCapacity
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   variant_id     query   int    	true "variant_id"
// @Param   storage_location_id     query   int    	true "storage_location_id"
// @Success 	200 {object} res.SuccessResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/get_count [get]
func (h *handler) GetCount(c echo.Context) (err error) {
	variant_id, _ := strconv.Atoi(c.QueryParam("variant_id"))
	storage_location_id, _ := strconv.Atoi(c.QueryParam("storage_location_id"))
	data, _ := h.locationsRepository.GetStorageLocation(map[string]interface{}{"id": storage_location_id})
	//query := map[string]interface{}{"id": data.Type.DisplayName}
	var count int64 = CountCapacity(variant_id, data, int64(0))
	return res.RespSuccess(c, "Count of the product in the specified location calculated", count)
	// data.Type.LookupCode != "BIN"
	//lookup,err := h.coreRepository.GetLookupCodes(query, new(pagination.Paginatevalue))
}

func CountCapacity(variant_id int, childStorageLocation shared_pricing_and_location.StorageLocation, count int64) (result int64) {
	if childStorageLocation.Type.LookupCode != "BIN" {
		for _, storage_location := range childStorageLocation.ChildStorageLocations {
			count = count + CountCapacity(variant_id, *storage_location, 0)
		}
	} else {
		for _, storage_quantity := range childStorageLocation.StorageQuantity {
			if variant_id == int(*storage_quantity.ProductVariantId) {
				count = count + storage_quantity.StorageCount
				fmt.Println(storage_quantity.ID)
			}
		}
	}
	return count
}

// GetProductlist godoc
// @Summary 	do a GetProductlist
// @Description do a GetProductlist
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param 		Authorization header string true "Authorization"
// @Param		GetProductList body GetProductsListDTO true "Create a CycleCount"
// @Success 	200 {object} res.SuccessResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/product_list [post]
func (h *handler) GetProductList(c echo.Context) (err error) {
	var request_data interface{}
	c.Bind(&request_data)
	var data []shared_pricing_and_location.StorageLocation
	data_ids := request_data.(map[string]interface{})["data"].([]interface{})
	for _, id := range data_ids {
		loc_data, _ := h.locationsRepository.GetStorageLocation(map[string]interface{}{"id": id.(float64)})
		data = append(data, loc_data)
	}
	result := make(map[int64]interface{}, 0)
	var prod_ids []int64
	for _, storage_location_data := range data {
		for _, storage_quantity_data := range storage_location_data.StorageQuantity {
			variant_id := int64(*storage_quantity_data.ProductVariantId)
			if helpers.Contains(prod_ids, variant_id) {
				result_kvp := result[variant_id].(map[string]interface{})
				storage_loc_data, _ := h.locationsRepository.GetStorageLocation(map[string]interface{}{"id": storage_quantity_data.StorageLocationId})
				result_kvp["bins"] = append(result_kvp["bins"].([]shared_pricing_and_location.StorageLocation), storage_loc_data)
				result_kvp["count"] = result_kvp["count"].(int64) + storage_quantity_data.StorageCount
				result[variant_id] = result_kvp
				fmt.Println("successs")
			} else {
				result_index := make(map[string]interface{}, 0)
				var bins_data []shared_pricing_and_location.StorageLocation
				storage_loc, _ := h.locationsRepository.GetStorageLocation(map[string]interface{}{"id": storage_quantity_data.StorageLocationId})
				bins_data = append(bins_data, storage_loc)
				result_index["bins"] = bins_data
				result_index["count"] = storage_quantity_data.StorageCount
				result[variant_id] = result_index
				fmt.Println("else condtn")
				prod_ids = append(prod_ids, variant_id)
			}
		}
	}
	return res.RespSuccess(c, "The product list in the given list of bins is retrieved", result)
}

// GetBinsList godoc
// @Summary 	do a GetBinsList
// @Description do a GetBinsList
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 	200 {object} res.SuccessResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/bin_list/{variant_id} [get]
func (h *handler) GetBinsList(c echo.Context) (err error) {
	variant_id_int, _ := strconv.Atoi(c.Param("variant_id"))
	variant_id := uint(variant_id_int)
	var response []shared_pricing_and_location.StorageLocation
	query := make(map[string]interface{}, 0)
	var bin_ids []int64
	p := new(pagination.Paginatevalue)
	//temporary hardcoding
	p.Per_page = 1000
	storage_quantity_data, _ := h.locationsRepository.GetStorageQuantityList(query, p)
	for _, data := range storage_quantity_data {
		data_variant_id := *data.ProductVariantId
		if variant_id == data_variant_id {
			query["id"] = data.ID
			bin_data, _ := h.locationsRepository.GetStorageLocation(query)
			if !helpers.Contains(bin_ids, int64(bin_data.ID)) {
				response = append(response, bin_data)
				bin_ids = append(bin_ids, int64(bin_data.ID))
			}
		}
	}
	return res.RespSuccess(c, "The bins list is retrieved", response)
}

// CreateCycleCount godoc
// @Summary 	do a CreateCycleCount
// @Description do a CreateCycleCount
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		CycleCountCreateRequest body CycleCountRequest true "Create a CycleCount"
// @Success 	200 {object} CycleCountCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/create [post]
func (h *handler) CreateCycleCountEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("cycle_count"),
	}
	eda.Produce(eda.CREATE_CYCLE_COUNT, request_payload)
	return res.RespSuccess(c, "Cycle Count Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateCycleCount(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	createPayload := new(inventory_tasks.CycleCount)
	helpers.JsonMarshaller(data, createPayload)
	err := h.service.CreateCycleCount(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_CYCLE_COUNT_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_CYCLE_COUNT_ACK, responseMessage)
	// cache implementation
	UpdateCycleCountInCache(edaMetaData)
}

// func (h *handler) CreateCycleCount(c echo.Context) (err error) {
// 	data := c.Get("cycle_count").(*CycleCountRequest)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	data.CreatedByID = helpers.ConvertStringToUint(token_id)
// 	id, err := h.service.CreateCycleCount(data, token_id, access_template_id)
// 	if err != nil {
// 		return res.RespErr(c, err)
// 	}

// 	// cache implementation
// 	UpdateCycleCountInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Cycle count created successfully", map[string]interface{}{"created_id": id})
// }

// BulkCreateCycleCount godoc
// @Summary 	do a BulkCreateCycleCount
// @Description do a BulkCreateCycleCount
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkCycleCountCreateRequest body CycleCountRequest true "Create a Bulk CycleCount"
// @Success 	200 {object} BulkCycleCountCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/bulk_create [post]
func (h *handler) BulkCreateCycleCount(c echo.Context) (err error) {
	data := new([]CycleCountRequest)
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	var metaData core.MetaData
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateCycleCount(metaData, data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Cycle count records created successfully", data)
}

// UpdateCycleCount godoc
// @Summary 	do a UpdateCycleCount
// @Description do a UpdateCycleCount
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param 		id 	path 	 int true "id"
// @Param		CycleCountUpdateRequest body CycleCountRequest true "Update a CycleCount"
// @Success 	200 {object} CycleCountUpdateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/{id}/edit [put]
func (h *handler) UpdateCycleCountEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	cycleCountId := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         cycleCountId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("cycle_count"),
	}
	eda.Produce(eda.UPDATE_CYCLE_COUNT, request_payload)
	return res.RespSuccess(c, "Cycle Count Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateCycleCount(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(inventory_tasks.CycleCount)
	helpers.JsonMarshaller(data, updatePayload)
	err := h.service.UpdateCycleCount(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_CYCLE_COUNT_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_CYCLE_COUNT_ACK, responseMessage)

	UpdateCycleCountInCache(edaMetaData)
}

// func (h *handler) UpdateCycleCount(c echo.Context) (err error) {
// 	ID := c.Param("id")
// 	id, _ := strconv.Atoi(ID)
// 	data := c.Get("cycle_count").(*CycleCountRequest)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
// 	err = h.service.UpdateCycleCount(uint(id), data, token_id, access_template_id)
// 	if err != nil {
// 		return res.RespErr(c, err)
// 	}

// 	// cache implementation
// 	UpdateCycleCountInCache(token_id, access_template_id)

// 	return res.RespSuccess(c, "Cycle count updated successfully", map[string]interface{}{"updated_id": id})
// }

// GetCycleCount godoc
// @Summary 	Find a GetCycleCount by id
// @Description Find a GetCycleCount by id
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		id  path 	 int true "id"
// @Success 	200 {object} CycleCountGetResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/{id} [get]
func (h *handler) GetCycleCount(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	cycleCountId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": cycleCountId,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetCycleCount(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response CycleCountRequest
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Cycle count fetched successfully", response)
}

// GetAllCycleCount godoc
// @Summary 	Get all GetAllCycleCount with MiscOptions of Pagination, Search, Filter and Sort
// @Description Get all GetAllCycleCount with MiscOptions of Pagination, Search, Filter and Sort
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		filters		  query	string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success 	200 {object} CycleCountGetAllResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count [get]
func (h *handler) GetAllCycleCount(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []CycleCountGetAll

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetCycleCountFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllCycleCount(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Cycle Count data Retrieved successfully", response, p)
}

// DeleteCycleCount godoc
// @Summary 	delete a DeleteCycleCount
// @Description delete a DeleteCycleCount
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		id	path	 int true "id"
// @Success 	200 {object} CycleCountDeleteResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/{id}/delete [delete]
func (h *handler) DeleteCycleCount(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	cycleCountId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         cycleCountId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteCycleCount(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateCycleCountInCache(metaData)

	return res.RespSuccess(c, " Cycle Count deleted successfully", map[string]interface{}{"deleted_id": cycleCountId})
}

// DeleteCycleCountLines godoc
// @Summary 	delete a DeleteCycleCountLines
// @Description delete a DeleteCycleCountLines
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		id	path	 int true "id"
// @Param       product_id query string true "product_id"
// @Success 	200 {object} CycleCountLinesDeleteResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/{id}/delete_products [delete]
func (h *handler) DeleteCycleCountLines(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	cycleCountId := c.Param("id")
	productId := c.QueryParam("product_id")
	metaData.Query = map[string]interface{}{
		"id":         cycleCountId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
		"product_id": productId,
	}

	err = h.service.DeleteCycleCountLines(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Cycle Count Line Items successfully deleted", cycleCountId)
}

// SendMailCycleCount godoc
// @Summary      Send Email to specific recevier for Id
// @Description  Send Email to specific recevier for Id
// @Tags         CycleCount
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Param        receiver_email query     string true "receiver_email"
// @Success      200  {object}  SendMailCycleCountResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/cycle_count/{id}/sendemail [get]
func (h *handler) SendMailCycleCount(c echo.Context) (err error) {

	q := new(SendMailCycleCount)
	c.Bind(q)
	// err = h.service.SendMailCycleCount(q)
	// if err != nil {
	// 	return res.RespErr(c, err)
	// }

	return res.RespSuccess(c, "Email Sent Successfully", q)
}

// DownloadPdfCycleCount godoc
// @Summary      Download PDF
// @Description  Download PDF
// @Tags         CycleCount
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id         	query     int    true "id"
// @Success      200  {object}  DownloadPdfCycleCountResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/cycle_count/{id}/printpdf [get]
func (h *handler) DownloadPdfCycleCount(c echo.Context) (err error) {
	return res.RespSuccess(c, "Pdf Generated Successfully", "https://drive.google.com/file/d/1O46B3lgliUGakL4weBjrovCHphTzR5Np/view?usp=sharing")
}

// FavouriteCycleCountView godoc
// @Summary Get all favourite CycleCount
// @Description Get all favourite CycleCount
// @Tags CycleCount
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CycleCountGetAllResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/cycle_count/favourite_list [get]
func (h *handler) FavouriteCycleCountView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrCycleCount(user_id)
	if er != nil {
		return errors.New("favourite CycleCount list not found for the specified user")
	}
	result, err := h.base_service.GetFavCycleCountList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Cycle count list Retrieved Successfully", result, p)
}

// FavoriteCycleCount godoc
// @Summary FavoriteCycleCount
// @Description FavoriteCycleCount
// @Tags CycleCount
// @Accept  json
// @Produce  json
// @param id path string true "CycleCount ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/cycle_count/{id}/favourite [post]
func (h *handler) FavouriteCycleCount(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.FavouriteCycleCount(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Cycle Count is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteCycleCount godoc
// @Summary UnFavoriteCycleCount
// @Description UnFavoriteCycleCount
// @Tags CycleCount
// @Accept  json
// @Produce  json
// @param id path string true "CycleCount ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/cycle_count/{id}/unfavourite [post]
func (h *handler) UnFavouriteCycleCount(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteCycleCount(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Cycle Count is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetAllCycleCountDropDown godoc
// @Summary 	Get all GetAllCycleCountDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Description Get all GetAllCycleCountDropDown with MiscOptions of Pagination, Search, Filter and Sort
// @Tags 		CycleCount
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		filters		  query	string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success 	200 {object} CycleCountGetAllResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/cycle_count/dropdown [get]
func (h *handler) GetAllCycleCountDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	var cacheResponse interface{}
	var response []CycleCountGetAll

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetCycleCountFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetAllCycleCount(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "Cycle Count data Retrieved successfully", response, p)
}
