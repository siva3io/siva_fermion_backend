package cores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"

	"fermion/backend_core/pkg/util/cache"
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

var CoreHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if CoreHandler != nil {
		return CoreHandler
	}
	service := NewService()
	CoreHandler = &handler{service}
	return CoreHandler
}

// Lookup Types godoc
// @Summary Lookup Types
// @Description get Lookup Types
// @Tags Core
// @Accept  json
// @Produce  json
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_types [get]
func (h *handler) GetLookupTypes(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetLookupTypes(query, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// GetAllLookupTypes godoc
// @Summary Get all Lookup Types
// @Description Get all Lookup Types
// @Tags Core
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} LookupTypesDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_types/list [get]
func (h *handler) GetAllLookupTypes(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetLookupTypesList(query, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Lookup Types list Retrieved Successfully", result, p)
}

// Lookup Code godoc
// @Summary get Lookup Code
// @Description get Lookup Code
// @Tags Core
// @Accept  json
// @Produce  json
// @Param type   path      string  true  "type"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_codes/{type} [get]
func (h *handler) GetLookupCodes(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)

	var lookup_type = c.Param("type")
	lookup_type = strings.ToUpper(lookup_type)

	var query = make(map[string]interface{}, 0)
	query["lookup_type"] = lookup_type

	result, err := h.service.GetLookupCodes(query, p)
	if err != nil {
		return res.RespError(c, err)
	}
	msg := fmt.Sprintf("%v - Lookup Codes Retrieved successfully", lookup_type)
	return res.RespSuccessInfo(c, msg, result, p)
}

// GetAllLookupCodes godoc
// @Summary Get all Lookup Codes
// @Description Get all Lookup Codes
// @Tags Core
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} LookupCodeDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_codes/list [get]
func (h *handler) GetAllLookupCodes(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetLookupCodesList(query, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Lookup Codes list Retrieved Successfully", result, p)
}

// Countries godoc
// @Summary get Countries
// @Description get Countries
// @Tags Core
// @Accept  json
// @Produce  json
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/countries [get]
func (h *handler) GetCountries(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetCountries(query, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// States godoc
// @Summary get States
// @Description get States
// @Tags Core
// @Accept  json
// @Produce  json
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/states/{id} [get]
func (h *handler) GetStates(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	id := c.Param("id")
	var query = make(map[string]interface{}, 0)
	query["country_id"] = id
	fmt.Println(query)
	result, err := h.service.GetStates(query, p)

	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// Currencies godoc
// @Summary get Currencies
// @Description get Currencies
// @Tags Core
// @Accept  json
// @Produce  json
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/currencies [get]
func (h *handler) GetCurrencies(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetCurrencies(query, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// Search Lookupcodes godoc
// @Summary Get all filtered  Lookupcodes
// @Description Get all filtered Lookupcodes
// @Tags Core
// @Accept json
// @Produce json
// @param q query string true "Search query"
// @Success 200 {array} LookupCodesDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_codes/search [get]
func (h *handler) SearchLookupCodes(c echo.Context) (err error) {
	q := c.QueryParam("q")
	result, err := h.service.SearchLookupCodes(q)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "OK", result)
}

// Lookup Types godoc
// @Summary Lookup Types
// @Description create Lookup Types
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /api/v1/core/lookup_type/create [post]
func (h *handler) CreateLookupType(c echo.Context) (err error) {
	dto := c.Get("lookup_types").(*CreateLookupTypeDTO)
	var data *model_core.Lookuptype
	value, _ := json.Marshal(dto)
	json.Unmarshal(value, &data)
	err = h.service.CreateLookupType(data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "LookupType Created Successfully", map[string]interface{}{"created_id": data.ID})
}

// Lookup Codes godoc
// @Summary Lookup Codes
// @Description create Lookup Codes
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_codes/create [post]
func (h *handler) CreateLookupCodes(c echo.Context) (err error) {
	dto := c.Get("lookup_codes").(*CreateUpdateLookupCodesDTO)
	var data []model_core.Lookupcode
	value, _ := json.Marshal(dto.LookupCodes)
	json.Unmarshal(value, &data)
	err = h.service.CreateLookupCodes(&data)
	if err != nil {
		return res.RespError(c, err)
	}
	var result CreateUpdateLookupCodesResponseDTO
	response, _ := json.Marshal(data)
	json.Unmarshal(response, &result.LookupCodes)

	return res.RespSuccess(c, "LookupType Created Successfully", result)
}

// UpdateLookupType godoc
// @Summary UpdateLookupType
// @Description Update Lookup Type
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_type/{id}/update [post]
func (h *handler) UpdateLookupType(c echo.Context) (err error) {
	dto := c.Get("lookup_types").(*UpdateLookupTypeDTO)
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	query := map[string]interface{}{
		"id": id,
	}
	var data *model_core.Lookuptype
	value, _ := json.Marshal(dto)
	json.Unmarshal(value, &data)

	err = h.service.UpdateLookupType(data, query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "LookupType Updated Successfully", map[string]interface{}{"updated_id": id})
}

// UpdateLookupCode godoc
// @Summary UpdateLookupCode
// @Description Update Lookup Code
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_code/{id}/update [post]
func (h *handler) UpdateLookupCode(c echo.Context) (err error) {
	dto := c.Get("lookup_codes").(*UpdateLookupCodeDTO)
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	query := map[string]interface{}{
		"id": id,
	}
	var data *model_core.Lookupcode
	value, _ := json.Marshal(dto)
	json.Unmarshal(value, &data)

	err = h.service.UpdateLookupCode(data, query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "LookupCode Updated Successfully", map[string]interface{}{"updated_id": id})
}

// DeleteLookupType godoc
// @Summary DeleteLookupType
// @Description Delete Lookup Type
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_type/{id}/delete [delete]
func (h *handler) DeleteLookupType(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	err = h.service.DeleteLookupType(uint(id))
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "LookupType Deleted Successfully", map[string]interface{}{"deleted_id": id})
}

// DeleteLookupCode godoc
// @Summary DeleteLookupCode
// @Description Delete Lookup Code
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/lookup_code/{id}/delete [delete]
func (h *handler) DeleteLookupCode(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	err = h.service.DeleteLookupCode(uint(id))
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "LookupCode Deleted Successfully", map[string]interface{}{"deleted_id": id})
}

// ListStoreApps godoc
// @Summary ListStoreApps
// @Description List Apps
// @Tags Core
// @Accept  json
// @Produce  json
// @Param per_page 			query int 	 false "per_page"
// @Param page_no           query int 	 false "page_no"
// @Param sort				query string false "sort"
// @param filters 			query string false "Filters"
// @Success 200 {object} AppsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/store/apps [get]
func (h *handler) ListStoreApps(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	var response []AppsListDTO
	var data []model_core.AppStore
	var app_type = c.QueryParam("type")

	if app_type == "installed" {
		data, err = h.service.ListStoreAppsInstalled(page)
		value, _ := json.Marshal(data)
		json.Unmarshal(value, &response)
		for ind := range response {
			response[ind].State = "Installed"
		}
		if err != nil {
			return res.RespError(c, err)
		}
		return res.RespSuccessInfo(c, "List retrieved successfully", response, page)
	}

	data, err = h.service.ListStoreApps(app_type, page)
	value, _ := json.Marshal(data)
	json.Unmarshal(value, &response)
	for ind := range response {
		response[ind].State = "Not Installed"
	}
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "List retrieved successfully", response, page)
}

// GetApp godoc
// @Summary GetApp
// @Description Get App
// @Tags Core
// @Accept  json
// @Produce  json
// @Param   id  path int true "Get App"
// @Success 200 {object} GetAppStoreDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/store/apps/{id} [get]
func (h *handler) GetApp(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)

	query := map[string]interface{}{
		"id": uint(id),
	}

	var data model_core.AppStore
	var response_data GetAppStoreDTO
	data, err = h.service.GetApp(query)
	if err != nil {
		return res.RespError(c, err)
	}
	value, _ := json.Marshal(data)
	json.Unmarshal(value, &response_data)
	check_code := h.service.CheckState(response_data.Code)
	response_data.State = "Not Installed"
	if check_code {
		response_data.State = "Installed"
	}
	return res.RespSuccess(c, "App fetched successfully", response_data)
}

// SearchApps godoc
// @Summary SearchApps
// @Description Search Apps
// @Tags Core
// @Accept  json
// @Produce  json
// @param name query string false "name"
// @param code query string false "code"
// @Success 200 {object} GetAppStoreDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/store/search [get]
func (h *handler) SearchApps(c echo.Context) (err error) {
	var data SearchAppsDTO
	c.Bind(&data)
	result, err := h.service.SearchApps(data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "App fetched", result)
}

func (h *handler) UpdateLocalCacheData(c echo.Context) (err error) {
	data := make(map[string]interface{})
	c.Bind(&data)
	fmt.Println("Payload received")
	for key, value := range data {
		fmt.Printf("Completed inserting %s credentials into local cache \n", key)
		cache.SetCacheVariable(key, value)
	}
	return res.RespSuccessInfo(c, " Local cache updated successfully", nil, nil)
}

func (h *handler) GetLocalCacheData(c echo.Context) (err error) {
	key := c.Param("key")
	fmt.Println(key)
	return res.RespSuccessInfo(c, " Local cache Retrieved successfully", cache.GetCacheVariable(key), nil)
}

func (h *handler) UpdateRedisCacheData(c echo.Context) (err error) {
	data := make(map[string]interface{})
	c.Bind(&data)
	fmt.Println("Payload received")
	for key, value := range data {
		fmt.Printf("Completed inserting %s credentials into cache \n", key)
		cache.SetCacheKey(key, value)
	}
	return res.RespSuccessInfo(c, " Cache updated successfully", nil, nil)
}

func (h *handler) GetRedisCacheData(c echo.Context) (err error) {
	key := c.Param("key")
	fmt.Println(key)
	return res.RespSuccessInfo(c, " Cache Retrieved successfully", cache.GetCacheKey(key), nil)
}

func (h *handler) GetChannelLookupCodes(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.AllChannelLookupCodes(query, p)
	if err != nil {
		return res.RespError(c, err)
	}

	return res.RespSuccessInfo(c, "AllChannelLookupCodes data retrieved successfully", result, p)
}

// ListMetaData godoc
// @Summary ListMetaData
// @Description List MetaData
// @Tags Core
// @Accept  json
// @Produce  json
// @Success 200 {object} ListMetaDataDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/meta_data [get]
func (h *handler) ListMetaData(c echo.Context) (err error) {
	response, err := h.service.ListMetaData()
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "success", response)
}

// ViewMetaDta godoc
// @Summary ViewMetaDta
// @Description View MetaDta
// @Tags Core
// @Accept  json
// @Produce  json
// @Param   model_name  path string true "model_name"
// @Success 200 {object} ViewMetaDataDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/meta_data/{model_name} [get]
func (h *handler) ViewMetaData(c echo.Context) (err error) {
	model_name := c.Param("model_name")
	response, err := h.service.ViewMetaData(model_name)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "success", response)
}

// ListInstalledApps godoc
// @Summary ListInstalledApps
// @Description List Installed Apps
// @Tags Core
// @Accept  json
// @Produce  json
// @Param per_page 			query int 	 false "per_page"
// @Param page_no           query int 	 false "page_no"
// @Param sort				query string false "sort"
// @param filters 			query string false "Filters"
// @param registered 		query string false "registered"
// @Success 200 {object} InstalledAppsDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/apps/installed [get]
func (h *handler) ListInstalledApps(c echo.Context) (err error) {
	queryParam := make(map[string]interface{})
	var registered bool
	page := new(pagination.Paginatevalue)
	err = c.Bind(&queryParam)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.ParseQueryToPaginate(queryParam, page)
	if queryParam["registered"] != nil {
		registered, _ = strconv.ParseBool(queryParam["registered"].(string))
	}
	response, err := h.service.ListInstalledApps(page, registered)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "success", response, page)
}

// InstallApp godoc
// @Summary InstallApp
// @Description Install App
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id  path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/store/apps/{id}/install [post]
func (h *handler) InstallApp(c echo.Context) (err error) {
	app_code := c.Param("app_code")
	app_code = strings.ToUpper(app_code)
	user_id := c.Get("TokenUserID").(string)
	err = h.service.InstallApp(app_code, user_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "success", map[string]interface{}{"message": "App Successfully Installed"})
}

// UninstallApp godoc
// @Summary UninstallApp
// @Description Uninstall App
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   id  path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/store/apps/{id}/uninstall [post]
func (h *handler) UninstallApp(c echo.Context) (err error) {
	app_code := c.Param("app_code")
	app_code = strings.ToUpper(app_code)
	err = h.service.UninstallApp(app_code)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "success", map[string]interface{}{"message": "App Successfully UnInstalled"})
}

func (h *handler) UpdateInstalledApp(c echo.Context) error {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	query := map[string]interface{}{
		"id": id,
	}
	var data model_core.InstalledApps
	err := c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = h.service.UpdateInstalledApp(&data, query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "successfully Update App  details", data)

}

func (h *handler) RenewSubscription(c echo.Context) (err error) {
	app_code := c.QueryParam("code")
	// company_id:= c.QueryParam("company_id")
	user_id := c.Get("TokenUserID").(string)
	fmt.Println("userid------------>", user_id)
	query := map[string]interface{}{
		"app_code": app_code,
		// "company_id": company_id,

	}
	result := h.service.RenewSubscription(query, user_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Subscription is updated successfully!", result)
}

func (h *handler) GetCompanyPreferences(c echo.Context) (err error) {
	var id = c.Param("id")
	company_id, _ := strconv.Atoi(id)
	var query = make(map[string]interface{})
	query["id"] = company_id

	result, err := h.service.GetCompanyPreferences(query)
	if err != nil {
		fmt.Println(err)
	}

	return res.RespSuccess(c, "success", result)
}

func (h *handler) ListCompanyPreferences(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	var data []ChannelStatusDTO
	response, err := h.service.ListCompanies(page)
	value, _ := json.Marshal(response)

	_ = json.Unmarshal(value, &data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "success", response, page)
}

func (h *handler) CreateAppStore(c echo.Context) error {
	var data model_core.AppStore
	err := c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = h.service.CreateAppStore(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "successfully created app", data)
}

func (h *handler) UpdateAppStore(c echo.Context) error {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	query := map[string]interface{}{
		"id": id,
	}
	var data model_core.AppStore
	err := c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = h.service.UpdateAppStore(&data, query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "successfully Update App store details", data)

}

func (h *handler) ListCompanyUsers(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var id = c.Param("id")
	company_id, _ := strconv.Atoi(id)
	var query = make(map[string]interface{})
	query["company_id"] = company_id

	result, err := h.service.ListCompanyUsers(query, p)
	if err != nil {
		fmt.Println(err)
	}
	return res.RespSuccessInfo(c, "Users list Retrieved Successfully", result, p)

}

func (h *handler) UpdateCompany(c echo.Context) (err error) {
	var id = c.Param("id")
	company_id, _ := strconv.Atoi(id)

	data := make(map[string]interface{})
	c.Bind(&data)

	var query = make(map[string]interface{})
	query["id"] = company_id

	var company_details Company
	data["id"] = company_id
	dto, _ := json.Marshal(data)
	err = json.Unmarshal(dto, &company_details)

	if err != nil {
		fmt.Println(err)
		return res.RespError(c, err)
	}
	host := c.Request().Host
	scheme := c.Scheme()

	s := c.Get("TokenUserID").(string)
	company_details.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateCompany(query, company_details, host, scheme)
	if err != nil {
		fmt.Println(err)
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "Company Details updated successfully", map[string]interface{}{"updated_id": company_id})
	// return res.RespSuccess(c, "Company Details updated successfully", map[string]interface{}{"updated_id": company_id, "ondc": ondc})
}

func (h *handler) OndcRegistration(c echo.Context) (err error) {
	var id = c.Param("id")
	company_id, _ := strconv.Atoi(id)

	data := make(map[string]interface{})
	c.Bind(&data)

	var query = make(map[string]interface{})
	query["id"] = company_id

	var company_details Company
	data["id"] = company_id
	dto, _ := json.Marshal(data)
	err = json.Unmarshal(dto, &company_details)

	if err != nil {
		fmt.Println(err)
		return res.RespError(c, err)
	}
	ondc, err := h.service.OndcRegistration(query, company_details)
	if err != nil {
		fmt.Println(err)
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "Company Details updated successfully", map[string]interface{}{"updated_id": company_id, "ondc": ondc})
}

// ListChannelStatus godoc
// @Summary ListChannelStatus
// @Description get ChannelStatus
// @Tags Core
// @Accept  json
// @Produce  json
// @Param per_page 			query int 	 false "per_page"
// @Param page_no           query int 	 false "page_no"
// @Param sort				query string false "sort"
// @param filters 			query string false "Filters"
// @param Authorization header string true "Authorization"
// @Success 200 {object} GetAllChannelStatusResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/channel_status [get]
func (h *handler) ListChannelStatus(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	var data []ChannelStatusDTO
	response, err := h.service.ListChannelStatus(page)
	value, _ := json.Marshal(response)

	_ = json.Unmarshal(value, &data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "success", response, page)
}

// ViewChannelStatus godoc
// @Summary ViewChannelStatus
// @Description ViewChannelStatus
// @Tags Core
// @Accept  json
// @Produce  json
// @param Authorization header string true "Authorization"
// @Param   channel_status_id  path string true "id"
// @Success 200 {object} ChannelStatusGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/channel_status/{id} [get]
func (h *handler) ViewChannelStatus(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	result, err := h.service.ViewChannelStatus(uint(id))
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// CreateChannelStatus godoc
// @Summary Create ChannelStatus
// @Description Create ChannelStatus
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body ChannelStatusDTO true "Channel status Request Body"
// @Success 200 {object} ChannelStatusCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/channel_status/create [post]
func (h *handler) CreateChannelStatus(c echo.Context) error {
	var data model_core.ChannelLookupCodes
	err := c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = h.service.CreateChannelStatus(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "channel status created successfully", data)
}

// UpdateChannelStatus godoc
// @Summary Update ChannelStatus
// @Description Update ChannelStatus
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "channel_status_id"
// @param RequestBody body ChannelStatusDTO true "Channel status Request Body"
// @Success 200 {object} ChannelStatusUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/channel_status/{id}/update [post]
func (h *handler) UpdateChannelStatus(c echo.Context) error {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	query := map[string]interface{}{
		"id": id,
	}
	var data model_core.ChannelLookupCodes
	err := c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = h.service.UpdateChannelStatus(&data, query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Channel Status Updated Successfully", data)

}

// DeleteChannelStatus godoc
// @Summary Delete ChannelStatus
// @Description Delete ChannelStatus
// @Tags Core
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "channel_status_id"
// @Success 200 {object} ChannelStatusDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/core/channel_status/{id}/delete [delete]
func (h *handler) DeleteChannelStatus(c echo.Context) (err error) {
	get_id := c.Param("id")
	id, _ := strconv.Atoi(get_id)
	err = h.service.DeleteChannelStatus(uint(id))
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Channel Status Deleted Successfully", map[string]interface{}{"deleted_id": id})
}

func (h *handler) UpdateOndcDetails(c echo.Context) (err error) {
	company_id := c.Get("CompanyId").(string)

	data := make(map[string]interface{})
	c.Bind(&data)

	query := map[string]interface{}{
		"id": company_id,
	}
	result, err := h.service.GetCompanyPreferences(query)
	if err != nil {
		fmt.Println(err)
	}
	OndcID := result.OndcDetailsId
	query = map[string]interface{}{
		"id": OndcID,
	}
	err = h.service.UpdateOndcDetails(query, data)
	if err != nil {
		fmt.Println(err)
		return res.RespError(c, err)
	}

	return res.RespSuccess(c, "Ondc Details updated successfully", map[string]interface{}{"updated_id": company_id})
}

func (h *handler) AadhaarValidation(c echo.Context) (err error) {
	data := make(map[string]interface{})
	c.Bind(&data)
	regex, _ := regexp.MatchString("^[0-9]{4}[ -]?[0-9]{4}[ -]?[0-9]{4}$", data["aadhar_number"].(string))

	return c.JSON(http.StatusOK, map[string]interface{}{"success": regex})
}

func (h *handler) RequestSolution(c echo.Context) (err error) {
	var data model_core.CustomSolution
	c.Bind(&data)
	CompanyId := c.Get("CompanyId").(string)
	data.CompanyId = helpers.ConvertStringToUint(CompanyId)
	_ = h.service.RequestSolution(data)
	return c.JSON(http.StatusOK, map[string]interface{}{"success": data})
}

func (h *handler) ListCompaniesAdmin(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	result, err := h.service.ListCompanies(page)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved Successfully", result, page)

}
