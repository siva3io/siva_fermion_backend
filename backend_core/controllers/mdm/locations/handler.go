package locations

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	mdm_base "fermion/backend_core/controllers/mdm/base"
	"fermion/backend_core/internal/repository"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
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
	service        Service
	base_service   mdm_base.ServiceBase
	coreRepository repository.Core
}

var LocationsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if LocationsHandler != nil {
		return LocationsHandler
	}
	service := NewService()
	coreRepository := repository.NewCore()
	base_service := mdm_base.NewServiceBase()
	LocationsHandler = &handler{service, base_service, coreRepository}
	return LocationsHandler
}

// CreateLocation godoc
// @Summary CreateLocation
// @Description CreateLocation
// @Tags Locations
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body LocationRequestDTO true "Location  Request Body"
// @Success 200 {object} LocationViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/create [post]
func (h *handler) CreateLocationEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("locations"),
	}
	eda.Produce(eda.CREATE_LOCATION, request_payload)
	// h.CreateLocation(request_payload)
	return res.RespSuccess(c, "Location Creation In Progress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateLocation(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	var locationDto LocationRequestDTO
	helpers.JsonMarshaller(request["data"], &locationDto)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	lookupCode, _ := helpers.GetLookupcodeName(locationDto.LocationTypeId)

	if lookupCode == "OFFICE" {
		officeLocation := new(shared_pricing_and_location.Office)
		helpers.JsonMarshaller(locationDto.LocationDetails, officeLocation)
		err := h.service.CreateOfficeLocation(edaMetaData, officeLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_LOCATION_ACK, responseMessage)
			return
		}
		locationDto.RelatedLocationID = officeLocation.ID
		helpers.JsonMarshaller(officeLocation.LocationInchargerDetails, &locationDto)

	} else if lookupCode == "VIRTUAL_LOCATION" {
		virtualLocation := new(shared_pricing_and_location.VirtualLocation)
		helpers.JsonMarshaller(locationDto.LocationDetails, virtualLocation)
		err := h.service.CreateVirtualLocation(edaMetaData, virtualLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_LOCATION_ACK, responseMessage)
			return
		}
		locationDto.RelatedLocationID = virtualLocation.ID
		helpers.JsonMarshaller(virtualLocation.LocationInchargerDetails, &locationDto)

	} else if lookupCode == "RETAIL" {
		retialLocation := new(shared_pricing_and_location.Retail)
		helpers.JsonMarshaller(locationDto.LocationDetails, retialLocation)

		err := h.service.CreateRetailLocation(edaMetaData, retialLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_LOCATION_ACK, responseMessage)
			return
		}
		locationDto.RelatedLocationID = retialLocation.ID
		helpers.JsonMarshaller(retialLocation.LocationInchargerDetails, &locationDto)

	} else if lookupCode == "LOCAL_WAREHOUSE" {
		localWarehouseLocation := new(shared_pricing_and_location.LocalWarehouse)
		helpers.JsonMarshaller(locationDto.LocationDetails, localWarehouseLocation)

		err := h.service.CreateWarehouseLocation(edaMetaData, localWarehouseLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_LOCATION_ACK, responseMessage)
			return
		}
		locationDto.RelatedLocationID = localWarehouseLocation.ID

	} else if lookupCode == "PICKUP_LOCATION" {
		fmt.Println("PICKUP_LOCATION")
	} else {
		responseMessage.ErrorMessage = fmt.Errorf("oops! invalid location type")
		// eda.Produce(eda.CREATE_LOCATION_ACK, responseMessage)
		return
	}

	location := new(shared_pricing_and_location.Locations)
	helpers.JsonMarshaller(locationDto, location)

	err := h.service.CreateLocation(edaMetaData, location)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_LOCATION_ACK, responseMessage)
		return
	}

	// updating cache
	UpdateLocationInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{"created_id": location.ID}
	// eda.Produce(eda.CREATE_LOCATION_ACK, responseMessage)
}

// UpdateLocation godoc
// @Summary UpdateLocation
// @Description UpdateLocation
// @Tags Locations
// @Accept  json
// @Produce  json
// @param id path string true "Location ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body LocationsDTO true "Location  Request Body"
// @Success 200 {object} LocationViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/{id}/update [post]
func (h *handler) UpdateLocationEvent(c echo.Context) (err error) {

	edaMetaData := c.Get("MetaData").(core.MetaData)
	locationId := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         locationId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("locations"),
	}
	eda.Produce(eda.UPDATE_LOCATION, request_payload)
	// h.UpdateLocation(request_payload)
	return res.RespSuccess(c, "Location Details updation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateLocation(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"]
	updateQuery := edaMetaData.Query

	var locationDto LocationRequestDTO
	helpers.JsonMarshaller(data, &locationDto)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	locationDetails, err := h.service.GetLocation(edaMetaData)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)
		return
	}
	var locationResponse shared_pricing_and_location.Locations
	helpers.JsonMarshaller(locationDetails, &locationResponse)

	lookupCode, _ := helpers.GetLookupcodeName(locationResponse.LocationTypeId)

	if lookupCode == "OFFICE" {
		edaMetaData.Query["iCreateWarehouseLocationd"] = locationResponse.RelatedLocationID

		officeLocation := new(shared_pricing_and_location.Office)
		helpers.JsonMarshaller(locationDto.LocationDetails, officeLocation)

		err := h.service.UpdateOfficeLocation(edaMetaData, officeLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)
			return
		}
		helpers.JsonMarshaller(officeLocation.LocationInchargerDetails, &locationDto)

	} else if lookupCode == "VIRTUAL_LOCATION" {
		edaMetaData.Query["id"] = locationResponse.RelatedLocationID

		virtualLocation := new(shared_pricing_and_location.VirtualLocation)
		helpers.JsonMarshaller(locationDto.LocationDetails, virtualLocation)
		err := h.service.UpdateVirtualLocation(edaMetaData, virtualLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)
			return
		}
		helpers.JsonMarshaller(virtualLocation.LocationInchargerDetails, &locationDto)

	} else if lookupCode == "RETAIL" {
		edaMetaData.Query["id"] = locationResponse.RelatedLocationID

		retialLocation := new(shared_pricing_and_location.Retail)
		helpers.JsonMarshaller(locationDto.LocationDetails, retialLocation)

		err := h.service.UpdateRetailLocation(edaMetaData, retialLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)
			return
		}
		helpers.JsonMarshaller(retialLocation.LocationInchargerDetails, &locationDto)

	} else if lookupCode == "LOCAL_WAREHOUSE" {
		edaMetaData.Query["id"] = locationResponse.RelatedLocationID

		localWarehouseLocation := new(shared_pricing_and_location.LocalWarehouse)
		helpers.JsonMarshaller(locationDto.LocationDetails, localWarehouseLocation)

		err := h.service.UpdateWarehouseLocation(edaMetaData, localWarehouseLocation)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)
			return
		}
	} else {
		responseMessage.ErrorMessage = fmt.Errorf("oops! invalid location type")
		// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)
		return
	}

	location := new(shared_pricing_and_location.Locations)
	helpers.JsonMarshaller(locationDto, location)

	location.LocationTypeId = 0

	edaMetaData.Query["id"] = locationResponse.ID
	err = h.service.UpdateLocation(edaMetaData, location)
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)
		return
	}

	// updating cache
	UpdateLocationInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{"updated_id": updateQuery["id"]}
	// eda.Produce(eda.UPDATE_LOCATION_ACK, responseMessage)

}

// DeleteLocation godoc
// @Summary DeleteLocation
// @Description DeleteLocation
// @Tags Locations
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Location ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/{id}/delete [delete]
func (h *handler) DeleteLocation(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	locationId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         locationId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteLocation(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// updating cache
	UpdateLocationInCache(metaData)

	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": locationId})
}

// LocationView godoc
// @Summary View LocationView
// @Summary View LocationView
// @Description View LocationView
// @Tags Locations
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Location ID"
// @Success 200 {object} LocationViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/{id} [get]
func (h *handler) GetLocation(c echo.Context) error {

	metaData := c.Get("MetaData").(core.MetaData)
	pricingId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": pricingId,
		// "company_id": metaData.CompanyId,
	}
	location, err := h.service.GetLocation(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response LocationResponseDTO
	helpers.JsonMarshaller(location, &response)

	lookupCode := response.LocationType.LookupCode
	//------------Related LocationType------------------
	LocationMap := map[string]interface{}{
		"VIRTUAL_LOCATION": h.service.GetVirtualLocation,
		"LOCAL_WAREHOUSE":  h.service.GetWarehouseLocation,
		"RETAIL":           h.service.GetRetailLocation,
		"OFFICE":           h.service.GetOfficeLocation,
	}

	metaData.Query = map[string]interface{}{
		"id": response.RelatedLocationID,
	}
	relatedLocationResponse, _ := LocationMap[lookupCode].(func(core.MetaData) (interface{}, error))(metaData)

	response.LocationTypeDetails = response.LocationType
	response.LocationDetails, _ = json.Marshal(relatedLocationResponse)

	return res.RespSuccess(c, "Location Details Retrieved successfully", response)
}

// GetLocationList godoc
// @Summary Get all Location  list
// @Description Get all Location  list
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} LocationListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations [get]
func (h *handler) GetLocationList(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	location_type_id, _ := helpers.GetLookupcodeId("LOCATION_TYPE", "VIRTUAL_LOCATION")
	// fmt.Println("=======>>>", location_type_id)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	helpers.AddMandatoryFilters(p, "location_type_id", "!=", location_type_id)

	var cacheResponse interface{}
	var response []LocationResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetLocationFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetLocationList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	for i, dto := range response {
		response[i].LocationTypeDetails = dto.LocationType
	}

	return res.RespSuccessInfo(c, " Location List Retrieved successfully", response, p)
}

// GetLocationListDropdown godoc
// @Summary Get all Location  list
// @Description Get all Location  list
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} LocationListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/dropdown [get]
func (h *handler) GetLocationListDropdown(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []LocationResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetLocationFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.GetLocationList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccessInfo(c, " Location List Retrieved successfully", response, p)
}

// FavouriteLocations godoc
// @Summary FavouriteLocations
// @Description FavouriteLocations
// @Tags Locations
// @Accept  json
// @Produce  json
// @param id path string true "Location ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/{id}/favourite [post]
func (h *handler) FavouriteLocations(c echo.Context) (err error) {

	metaData := c.Get("MetaData").(core.MetaData)
	locationId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         locationId,
		"company_id": metaData.CompanyId,
	}

	_, er := h.service.GetLocation(metaData)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	query := map[string]interface{}{
		"ID":      locationId,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.FavouriteLocations(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Location is Marked as Favourite", map[string]string{"record_id": locationId})
}

// UnFavouriteLocations godoc
// @Summary UnFavouriteLocations
// @Description UnFavouriteLocations
// @Tags Locations
// @Accept  json
// @Produce  json
// @param id path string true "Location ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/{id}/unfavourite [post]
func (h *handler) UnFavouriteLocations(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	locationId := c.Param("id")

	query := map[string]interface{}{
		"ID":      locationId,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.UnFavouriteLocations(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Location is Unmarked as Favourite", map[string]string{"record_id": locationId})
}

// GetFavouriteLocations godoc
// @Summary Get all favourite locations
// @Description Get all favourite locations
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} LocationListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/favourite_list [get]
func (h *handler) FavouriteLocationsView(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	tokenUserId := int(metaData.TokenUserId)
	data, er := h.base_service.GetArrLocation(tokenUserId)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavLocationList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Location list Retrieved Successfully", result, p)
}

//-----------------------------------Storage Location---------------------------------------------------------

// CreateStorageLocation godoc
// @Summary CreateStorageLocation
// @Description CreateStorageLocation
// @Tags Locations
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body StorageLocationDTO true "Storage Location  Request Body"
// @Success 200 {object} StorageLocationDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/storage_location/create [post]
func (h *handler) CreateStorageLocation(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)

	var request_data StorageLocationDTO
	c.Bind(&request_data)

	metaData.Query = map[string]interface{}{
		"id": *request_data.ParentStorageLocationId,
	}
	parentData, _ := h.service.GetStorageLocation(metaData)
	var parentdataResponse StorageLocationDTO
	helpers.JsonMarshaller(parentData, &parentdataResponse)

	if request_data.LocalWarehouseId == nil {
		request_data.LocalWarehouseId = parentdataResponse.LocalWarehouseId
	}
	if request_data.ParentStorageLocationId != nil {
		if parentdataResponse.ZoneId != nil {
			request_data.ZoneId = parentdataResponse.ZoneId
			if parentdataResponse.RowId != nil {
				request_data.RowId = parentdataResponse.RowId
				if parentdataResponse.RackId != nil {
					request_data.RackId = parentdataResponse.RackId
					request_data.ShelfId = request_data.ParentStorageLocationId
				} else {
					request_data.RackId = request_data.ParentStorageLocationId
				}
			} else {
				request_data.RowId = request_data.ParentStorageLocationId
			}
		} else {
			request_data.ZoneId = request_data.ParentStorageLocationId
		}
	}

	createPayload := new(shared_pricing_and_location.StorageLocation)
	helpers.JsonMarshaller(request_data, createPayload)
	err = h.service.CreateStorageLocation(metaData, createPayload)
	return res.RespSuccess(c, "Storage Location created successfully", map[string]interface{}{"created_id": createPayload.ID})
}

// StorageLocationView godoc
// @Summary View StorageLocationView
// @Summary View StorageLocationView
// @Description View StorageLocationView
// @Tags Locations
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Storage Location ID"
// @Success 200 {object} StorageLocationDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/storage_location/{id} [get]
func (h *handler) GetStorageLocation(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	storageLocationId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         storageLocationId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetStorageLocation(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response StorageLocationDTO
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "Storage Location Details Retrieved successfully", response)
}

// GetStorageLocationList godoc
// @Summary Get all Storage Location  list
// @Description Get all Storage Location  list
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} StorageLocationListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/storage_location [get]
func (h *handler) GetStorageLocationList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	result, err := h.service.GetStorageLocationList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response []StorageLocationDTO
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccessInfo(c, "Storage Location List Retrieved successfully", response, p)
}

// GetStorageQuantityList godoc
// @Summary Get all Storage Quantity  list
// @Description Get all Storage Quantity  list
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} StorageQuantityListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/storage_quantity [get]
func (h *handler) GetStorageQuantityList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	result, err := h.service.GetStorageQuantityList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response []StorageQuantityDTO
	helpers.JsonMarshaller(result, &response)

	return res.RespSuccessInfo(c, "Storage Quantity List Retrieved successfully", response, p)
}
