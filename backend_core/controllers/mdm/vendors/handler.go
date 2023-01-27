package vendors

import (
	"errors"
	"strconv"

	"fermion/backend_core/controllers/eda"
	mdm_base "fermion/backend_core/controllers/mdm/base"

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
	service      Service
	base_service mdm_base.ServiceBase
}

var VendorsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if VendorsHandler != nil {
		return VendorsHandler
	}
	service := NewService()
	base_service := mdm_base.NewServiceBase()
	VendorsHandler = &handler{service, base_service}
	return VendorsHandler
}

//---------------------------------------------Vendors-----------------------------------------

// CreateVendors godoc
// @Summary CreateVendors
// @Description CreateVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorsRequestDTO true "Vendors Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/create [post]
func (h *handler) CreateVendorEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("vendors"),
	}
	eda.Produce(eda.CREATE_VENDOR, request_payload)
	return res.RespSuccess(c, "Vendor Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}
func (h *handler) CreateVendor(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(mdm.Vendors)
	helpers.JsonMarshaller(data, createPayload)
	err := h.service.CreateVendors(edaMetaData, createPayload)
	//response data
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_VENDOR_ACK, responseMessage)
		return
	}

	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_VENDOR_ACK, responseMessage)
}

// UpdateVendors godoc
// @Summary UpdateVendors
// @Description UpdateVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorsRequestDTO true "Vendors Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/update [post]
func (h *handler) UpdateVendorEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	vendortId := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         vendortId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("vendors"),
	}
	eda.Produce(eda.UPDATE_VENDOR, request_payload)
	return res.RespSuccess(c, "Vendor update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateVendor(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(mdm.Vendors)
	helpers.JsonMarshaller(data, updatePayload)
	err := h.service.UpdateVendors(edaMetaData, updatePayload)
	//response data
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_VENDOR_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_VENDOR_ACK, responseMessage)
}

// func (h *handler) UpdateVendors(c echo.Context) (err error) {
// 	var id = c.Param("id")
// 	var query = make(map[string]interface{}, 0)
// 	ID, _ := strconv.Atoi(id)
// 	query["id"] = ID
// 	data := c.Get("vendors").(*mdm.Vendors)
// 	token_id := c.Get("TokenUserID").(string)
// 	access_template_id := c.Get("AccessTemplateId").(string)
// 	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
// 	err = h.service.UpdateVendors(query, data, token_id, access_template_id)
// 	if err != nil {
// 		return res.RespError(c, err)
// 	}
// 	return res.RespSuccess(c, "Details updated succesfully", map[string]interface{}{"updated_id": query["id"], "Details": data})
// }

// DeleteVendors godoc
// @Summary      delete a Vendors
// @Description  delete a Vendors
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/vendors/{id}/delete [delete]
func (h *handler) DeleteVendors(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	vendortId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         vendortId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteVendors(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": vendortId})
}

// GetVendors godoc
// @Summary Get Details of specific vendor based on ID
// @Description Get Details of specific vendor based on ID
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success 200 {object} VendorDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id} [get]
func (h *handler) GetVendors(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	vendorId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         vendorId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetVendors(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor data retrieved successfully", map[string]interface{}{"Vendor_ID": vendorId, "Details": result})
}

// GetVendorList godoc
// @Summary Get all vendors
// @Description Get all vendors
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors [get]
func (h *handler) GetVendorsList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var response []mdm.Vendors

	result, err := h.service.GetVendorsList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "List Retrieved successfully", response, p)
}

// GetVendorListDropdown godoc
// @Summary Get all vendors
// @Description Get all vendors
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/dropdown [get]
func (h *handler) GetVendorsListDropdown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var response []mdm.Vendors

	result, err := h.service.GetVendorsList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "List Retrieved successfully", response, p)
}

// FavouriteVendors godoc
// @Summary FavouriteVendors
// @Description FavouriteVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "Vendor ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/favourite [post]
func (h *handler) FavouriteVendors(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	metaData := c.Get("MetaData").(core.MetaData)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.GetVendors(metaData)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteVendors(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteVendors godoc
// @Summary UnFavouriteVendors
// @Description UnFavouriteVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "Vendor ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/unfavourite [post]
func (h *handler) UnFavouriteVendors(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnfavouriteVendors(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetFavouriteVendors godoc
// @Summary Get all favourite vendors
// @Description Get all favourite vendors
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/favourite_list [get]
func (h *handler) FavouriteVendorsView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrVendor(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavVendorList(data, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Vendors list Retrieved Successfully", result, p)
}

//---------------------------------------------vendor pricelist------------------------------------------

// CreateVendorPriceList godoc
// @Summary CreateVendorPriceList
// @Description CreateVendorPriceList
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorPriceListsDTO true "VendorPricelist Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist/create [post]
func (h *handler) CreateVendorPriceList(c echo.Context) error {
	data := c.Get("vendors_price_list").(*mdm.VendorPriceLists)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err := h.service.CreateVendorPriceList(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor price list created successfully", map[string]interface{}{"created_id": data.ID, "Details": data})
}

// UpdateVendorPricelist godoc
// @Summary UpdateVendorPricelist
// @Description UpdateVendorPricelist
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorPriceListsDTO true "VendorPricelist Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/update [post]
func (h *handler) UpdateVendorPriceLists(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("vendors_price_list").(*mdm.VendorPriceLists)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateVendorPriceLists(query, data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Details updated succesfully", map[string]interface{}{"updated_id": query["id"], "Details": data})
}

// GetVendorPiceList godoc
// @Summary Get Details of specific vendorpricelist based on ID
// @Description Get Details of specific vendorpricelist based on ID
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success 200 {object} VendorPriceListsDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist/{id} [get]
func (h *handler) GetVendorPriceList(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	query["id"] = ID
	result, err := h.service.GetVendorPriceList(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor Pricelist data retrieved successfully", map[string]interface{}{"Vendor_Pricelist_ID": query["id"], "Details": result})
}

// GetVendorPriceLists godoc
// @Summary Get all vendors price list
// @Description Get all vendors price list
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist [get]
func (h *handler) GetVendorPriceLists(c echo.Context) (err error) {
	//q := c.QueryParam("vendor_id")
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	//id,_ := strconv.Atoi(q)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetVendorPriceLists(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "List Retrieved successfully", result, p)
}

// GetVendorPriceListsDropdown godoc
// @Summary Get all vendors price list
// @Description Get all vendors price list
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist/dropdown [get]
func (h *handler) GetVendorPriceListsDropdown(c echo.Context) (err error) {
	//q := c.QueryParam("vendor_id")
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	//id,_ := strconv.Atoi(q)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetVendorPriceLists(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "List Retrieved successfully", result, p)
}

// DeleteVendorPricelist godoc
// @Summary      delete a VendorPricelist
// @Description  delete a VendorPricelist
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/vendors/pricelist/{id}/delete [delete]
func (h *handler) DeleteVendorPricelist(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteVendorPricelist(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

func (h *handler) UpsertVendorEvent(c echo.Context) (err error) {
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
	request_payload["meta_data"] = edaMetaData
	eda.Produce(eda.UPSERT_VENDOR, request_payload)

	return res.RespSuccess(c, "Vendor Upsert Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}
func (h *handler) UpsertVendor(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	var data []interface{}
	helpers.JsonMarshaller(request["data"], &data)

	response, err := h.service.Upsert(edaMetaData, data)
	helpers.PrettyPrint("vendors upsert response", response)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPSERT_VENDORS_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"response": response,
	}

	// eda.Produce(eda.UPSERT_VENDORS_ACK, responseMessage)
}
