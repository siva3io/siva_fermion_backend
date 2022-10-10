package locations

import (
	"encoding/json"
	"errors"
	"strconv"

	mdm_base "fermion/backend_core/controllers/mdm/base"
	"fermion/backend_core/internal/repository"

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
	service        Service
	base_service   mdm_base.ServiceBase
	coreRepository repository.Core
}

func NewHandler() *handler {
	service := NewService()
	coreRepository := repository.NewCore()
	base_service := mdm_base.NewServiceBase()
	return &handler{service, base_service, coreRepository}
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
func (h *handler) CreateLocation(c echo.Context) (err error) {
	data := c.Get("locations").(*LocationRequestDTO)
	s := c.Get("TokenUserID").(string)
	data.CreatedByID = helpers.ConvertStringToUint(s)
	//--------Fetch Location Type from LookupCode--------------------
	query := map[string]interface{}{"id": data.LocationTypeId}
	lookupcode, err := h.coreRepository.GetLookupCodes(query, new(pagination.Paginatevalue))
	if err != nil || len(lookupcode) == 0 {
		return res.RespError(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	CreateLocationMap := map[string]interface{}{
		"VIRTUAL_LOCATION": h.service.CreateVirtualLocation,
		"LOCAL_WAREHOUSE":  h.service.CreateWarehouseLocation,
		"RETAIL":           h.service.CreateRetailLocation,
		"OFFICE":           h.service.CreateOfficeLocation,
	}
	if _, ok := CreateLocationMap[lookupcode[0].LookupCode]; !ok {
		return res.RespErr(c, errors.New("invalid location type"))
	}
	resp, err := CreateLocationMap[lookupcode[0].LookupCode].(func(LocationRequestDTO, string, string) (interface{}, error))(*data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Location created successfully", resp)
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
func (h *handler) UpdateLocation(c echo.Context) (err error) {

	var id = c.Param("id")

	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data := c.Get("locations").(*LocationRequestDTO)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	//---------------Fetch the location details from Location--------------
	location, err := h.service.GetLocation(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}

	//--------Fetch Location Type from LookupCode-----------------------------------
	lookupcode_query := map[string]interface{}{"id": location.(LocationResponseDTO).LocationType.Id}
	lookupcode, err := h.coreRepository.GetLookupCodes(lookupcode_query, new(pagination.Paginatevalue))
	if err != nil || len(lookupcode) == 0 {
		return res.RespError(c, err)
	}

	//----------LocationType Functuon Map--------------------------------
	UpdateLocationMap := map[string]interface{}{
		"VIRTUAL_LOCATION": h.service.UpdateVirtualLocation,
		"LOCAL_WAREHOUSE":  h.service.UpdateWarehouseLocation,
		"RETAIL":           h.service.UpdateRetailLocation,
		"OFFICE":           h.service.UpdateOfficeLocation,
	}

	if _, ok := UpdateLocationMap[lookupcode[0].LookupCode]; !ok {
		return res.RespErr(c, errors.New("invalid location type"))
	}

	if data.LocationTypeId != 0 && location.(LocationResponseDTO).LocationType.Id != data.LocationTypeId {
		return res.RespErr(c, errors.New("location type not matching"))
	}

	err = UpdateLocationMap[lookupcode[0].LookupCode].(func(map[string]interface{}, LocationRequestDTO, string, string) error)(query, *data, token_id, access_template_id)

	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Location Details updated succesfully", data)
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

	var id = c.Param("id")

	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteLocation(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
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

	var id = c.Param("id")

	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	location, err := h.service.GetLocation(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Location Details Retrieved successfully", location)
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

	p := new(pagination.Paginatevalue)

	err = c.Bind(p)
	if err != nil {
		return res.RespErr(c, err)
	}

	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetLocationList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, " Location List Retrieved successfully", result, p)
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

	p := new(pagination.Paginatevalue)

	err = c.Bind(p)
	if err != nil {
		return res.RespErr(c, err)
	}

	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetLocationList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, " Location List Retrieved successfully", result, p)
}

// SearchLocation godoc
// @Summary Get all filtered Location
// @Description Get all filtered Location
// @Tags Locations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param q query string true "Search query"
// @Success 200 {array} IdNameDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/locations/search [get]
func (h *handler) SearchLocation(c echo.Context) (err error) {

	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchLocation(q, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "OK", result)
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
	_, er := h.service.GetLocation(q, token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteLocations(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Location is Marked as Favourite", map[string]string{"record_id": id})
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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteLocations(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Location is Unmarked as Favourite", map[string]string{"record_id": id})
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
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrLocation(user_id)
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
	var data shared_pricing_and_location.StorageLocation
	var request_data StorageLocationDTO
	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	marshaldata, err := json.Marshal(data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(marshaldata, &request_data)
	if err != nil {
		return res.RespErr(c, err)
	}
	var query = make(map[string]interface{}, 0)
	query["id"] = request_data.ParentStorageLocationId
	parent_data, _ := h.service.GetStorageLocation(query, token_id, access_template_id)
	if request_data.LocalWarehouseId == nil {
		request_data.LocalWarehouseId = parent_data.LocalWarehouseId
	}
	if request_data.ParentStorageLocationId != nil {
		if parent_data.ZoneId != nil {
			request_data.ZoneId = parent_data.ZoneId
			if parent_data.RowId != nil {
				request_data.RowId = parent_data.RowId
				if parent_data.RackId != nil {
					request_data.RackId = parent_data.RackId
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
	s := c.Get("TokenUserID").(string)
	request_data.CreatedByID = helpers.ConvertStringToUint(s)
	response_data, err := h.service.CreateStorageLocation(request_data, token_id, access_template_id)
	return res.RespSuccess(c, "Storage Location created successfully", response_data)
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
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	result, err := h.service.GetStorageLocation(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Storage Location Details Retrieved successfully", result)
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
	p := new(pagination.Paginatevalue)
	//c.Bind(p)
	err = c.Bind(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)

	result, err := h.service.GetStorageLocationList(query, p, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "Storage Location List Retrieved successfully", result, p)
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
	p := new(pagination.Paginatevalue)
	//c.Bind(p)
	err = c.Bind(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)

	result, err := h.service.GetStorageQuantityList(query, p, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "Storage Quantity List Retrieved successfully", result, p)
}

// Get AuthKeys godoc
// @Summary      Get AuthKeys
// @Description  Get AuthKeys
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        virtual_warehouse   path      string  true  "warehouse_code"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/locations/{warehouse_code}/get_auth_keys [get]

// func (h *handler) GetAuthKeys(c echo.Context) (err error) {

// 	virtual_warehouse := c.Param("warehouse_code")
// 	s := c.Get("TokenUserID").(string)
// 	user_id, _ := strconv.Atoi(s)
// 	query := map[string]interface{}{
// 		"warehouse_code": virtual_warehouse,
// 		"created_by":     user_id,
// 	}
// 	result, err := h.service.GetAuthKeys(query)
// 	if err != nil {
// 		return res.RespErr(c, err)
// 	}
// 	return res.RespSuccess(c, "success", result)
// }
