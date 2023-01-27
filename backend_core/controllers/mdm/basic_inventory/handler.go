package basic_inventory

import (
	// "fmt"

	"strconv"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
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
	service Service
}

var BasicInventoryHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if BasicInventoryHandler != nil {
		return BasicInventoryHandler
	}
	service := NewService()
	BasicInventoryHandler = &handler{service}
	return BasicInventoryHandler
}

// Create Centralized Inventory godoc
// @Summary Create Centralized Inventory
// @Description Create Centralized Inventory
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body CentralizedBasicInventoryDTO true "Centralized Inventory Request Body"
// @Success 200 {object} CentralizedBasicInventoryDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/centralized/create [post]
func (h *handler) CreateCentrailizedInventoryEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("basic_inventory_centralized"),
	}
	eda.Produce(eda.CREATE_CENTRALIZED_INVENTORY, request_payload)
	return res.RespSuccess(c, "Basic centralized Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateCentrailizedInventory(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(mdm.CentralizedBasicInventory)
	helpers.JsonMarshaller(data, createPayload)
	err := h.service.CreateCentrailizedInventory(edaMetaData, createPayload)

	var responseMessage eda.ConsumerResponse

	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_CENTRALIZED_INVENTORY_ACK, responseMessage)
		return
	}

	// cache implementation
	UpdateCentralizedInventoryInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_CENTRALIZED_INVENTORY_ACK, responseMessage)
}

// CreateDecentralizedInventory godoc
// @Summary CreateDecentralizedInventory
// @Description CreateDecentralizedInventory
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body DecentralizedBasicInventoryDTO true "Decentralized Inventory Request Body"
// @Success 200 {object} DecentralizedBasicInventoryDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/decentralized/create [post]
func (h *handler) CreateDecentralizedInventoryEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("basic_inventory_decentralized"),
	}

	eda.Produce(eda.CREATE_INVENTORY, request_payload)
	return res.RespSuccess(c, "Basic Decentralized Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateDecentralizedInventory(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(mdm.DecentralizedBasicInventory)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.CreateDecentralizedInventory(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_INVENTORY_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateDecentralizedInventoryInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_INVENTORY_ACK, responseMessage)
}

// UpdateCentralizedInventory godoc
// @Summary UpdateCentralizedInventory
// @Description UpdateCentralizedInventory
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @param id path string true "ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body CentralizedBasicInventoryDTO true "Centralized Inventory Request Body"
// @Success 200 {object} CentralizedBasicInventoryDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/centralized/{id}/update [post]
func (h *handler) UpdateCentrailizedInventoryEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)

	var id = c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("basic_inventory_centralized"),
	}

	eda.Produce(eda.UPDATE_CENTRALIZED_INVENTORY, request_payload)
	return res.RespSuccess(c, "Basic centralized Updation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateCentrailizedInventory(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	updatePayload := new(mdm.CentralizedBasicInventory)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.UpdateCentrailizedInventory(edaMetaData, updatePayload)

	var responseMessage eda.ConsumerResponse

	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_CENTRALIZED_INVENTORY_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateCentralizedInventoryInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_CENTRALIZED_INVENTORY_ACK, responseMessage)
}

// UpdateDecentralizedInventory godoc
// @Summary UpdateDecentralizedInventory
// @Description UpdateDecentralizedInventory
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @param id path string true "ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body DecentralizedBasicInventoryDTO true "decentralized Inventory Request Body"
// @Success 200 {object} DecentralizedBasicInventoryDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/decentralized/{id}/update [post]
func (h *handler) UpdateDecentralizedInventoryEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         id,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("basic_inventory_decentralized"),
	}

	eda.Produce(eda.UPDATE_INVENTORY, request_payload)
	return res.RespSuccess(c, "Basic Decentralized Updation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateDecentralizedInventory(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(mdm.DecentralizedBasicInventory)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.UpdateDecentralizedInventory(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse

	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_INVENTORY_ACK, responseMessage)
		return
	}
	// cache implementation
	UpdateDecentralizedInventoryInCache(edaMetaData)
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_INVENTORY_ACK, responseMessage)

}

// DeleteCentralizedInventory godoc
// @Summary DeleteCentralizedInventory
// @Description DeleteCentralizedInventory
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/centralized/{id}/delete [delete]
func (h *handler) DeleteCentrailizedInventory(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteCentrailizedInventory(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateCentralizedInventoryInCache(metaData)

	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// DeleteDecentralizedInventory godoc
// @Summary DeleteDecentralizedInventory
// @Description DeleteDecentralizedInventory
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/decentralized/{id}/delete [delete]
func (h *handler) DeleteDecentralizedInventory(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         id,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}

	err = h.service.DeleteDecentralizedInventory(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateDecentralizedInventoryInCache(metaData)

	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// CentralizedInventoryView godoc
// @Summary View CentralizedInventoryView
// @Summary View CentralizedInventoryView
// @Description View CentralizedInventoryView
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ID"
// @Success 200 {object} CentralizedViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/centralized/{id} [get]
func (h *handler) GetCentrailizedInventory(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": id,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetCentrailizedInventory(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response CentralizedBasicInventoryResponseDTO
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccess(c, "Centralized Details Retrieved succesfully", response)
}

// DecentralizedInventoryView godoc
// @Summary View DecentralizedInventoryView
// @Summary View DecentralizedInventoryView
// @Description View DecentralizedInventoryView
// @Tags Basic_Inventory
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "ID"
// @Success 200 {object} DeentralizedViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/decentralized/{id} [get]
func (h *handler) GetDecentralizedInventory(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	id := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": id,
		// "company_id": metaData.CompanyId,
	}
	result, err := h.service.GetDecentralizedInventory(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response DecentralizedBasicInventoryResponseDTO
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccess(c, "Decentralized Details Retrieved succesfully", response)
}

// CentralizedInventoryList godoc
// @Summary CentralizedInventoryList
// @Description CentralizedInventoryList
// @Tags Basic_Inventory
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CentralizedListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/centralized [get]
func (h *handler) GetCentrailizedInventoryList(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []CentralizedBasicInventoryResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetCentralizedInventoryFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetCentrailizedInventoryList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccessInfo(c, " Centralized List Retrieved successfully", response, p)
}

// CentralizedInventoryList godoc
// @Summary CentralizedInventoryList
// @Description CentralizedInventoryList
// @Tags Basic_Inventory
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CentralizedListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/centralized/dropdown [get]
func (h *handler) GetCentrailizedInventoryListDropdown(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []CentralizedBasicInventoryResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetCentralizedInventoryFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetCentrailizedInventoryList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccessInfo(c, " Centralized List Retrieved successfully", response, p)
}

// DecentralizedInventoryList godoc
// @Summary DecentralizedList
// @Description DecentralizedList
// @Tags Basic_Inventory
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} DeentralizedListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/decentralized [get]
func (h *handler) GetDecentralizedInventoryList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []DecentralizedBasicInventoryResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetDecentralizedInventoryFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetDecentralizedInventoryList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccessInfo(c, " Decentralized List Retrieved successfully", response, p)
}

// DecentralizedInventoryListDropdown godoc
// @Summary DecentralizedList
// @Description DecentralizedList
// @Tags Basic_Inventory
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} DeentralizedListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/decentralized/dropdown [get]
func (h *handler) GetDecentralizedInventoryListDropdown(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []DecentralizedBasicInventoryResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetDecentralizedInventoryFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetDecentralizedInventoryList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccessInfo(c, " Decentralized List Retrieved successfully", response, p)
}

// -----------Channel Upsert Api--------------------------------------
func (h *handler) DeCentralisedInventoryUpsertEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	var arrayData []interface{}
	var objectData interface{}
	var data interface{}

	c.Bind(&data)
	err = helpers.JsonMarshaller(data, &arrayData)
	if err != nil {
		helpers.JsonMarshaller(data, &objectData)
		arrayData = append(arrayData, objectData)
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      arrayData,
	}
	eda.Produce(eda.UPSERT_DECENTRALIZED_INVENTORY, request_payload)
	return res.RespSuccess(c, "Multiple Records Upsert Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}
func (h *handler) DeCentralisedInventoryUpsert(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	var data []interface{}
	helpers.JsonMarshaller(request["data"], &data)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	_, err := h.service.UpsertInventoryTemplate(edaMetaData, data)
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPSERT_DECENTRALIZED_INVENTORY_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"upserted_message": "Decentralised inventory upserted successfully",
	}
	// eda.Produce(eda.UPSERT_DECENTRALIZED_INVENTORY_ACK, responseMessage)
}

func (h *handler) InventoryTransactionCreate(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)

	var data = new(mdm.CentralizedInventoryTransactions)

	err = c.Bind(data)
	if err != nil {
		return res.RespErr(c, err)
	}

	err = h.service.InventoryTransactionCreate(metaData, data)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Centralized Inventory Transaction created successfully", map[string]interface{}{"created_id": data.ID})
}
