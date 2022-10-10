package basic_inventory

import (
	// "fmt"

	"encoding/json"
	"fmt"
	"strconv"

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
	service Service
}

func NewHandler() *handler {
	service := NewService()
	return &handler{service}
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
func (h *handler) CreateCentrailizedInventory(c echo.Context) (err error) {

	data := c.Get("basic_inventory_centralized").(*CentralizedBasicInventoryDTO)

	var request_data mdm.CentralizedBasicInventory

	marshaldata, err := json.Marshal(*data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(marshaldata, &request_data)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	request_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.CreateCentrailizedInventory(&request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Basic Centralized Inventory created successfully", request_data)
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
func (h *handler) CreateDecentralizedInventory(c echo.Context) (err error) {

	data := c.Get("basic_inventory_decentralized").(*DecentralizedBasicInventoryDTO)
	var request_data mdm.DecentralizedBasicInventory

	fmt.Println("------------------------------")
	marshaldata, err := json.Marshal(*data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(marshaldata, &request_data)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	request_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.CreateDecentralizedInventory(&request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Basic Decentralized created successfully", request_data)
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
func (h *handler) UpdateCentrailizedInventory(c echo.Context) (err error) {

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID

	data := c.Get("basic_inventory_centralized").(*CentralizedBasicInventoryDTO)

	var request_data mdm.CentralizedBasicInventory

	marshaldata, err := json.Marshal(*data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(marshaldata, &request_data)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	request_data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateCentrailizedInventory(query, &request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Centralized Inventory Details updated succesfully", request_data)
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
func (h *handler) UpdateDecentralizedInventory(c echo.Context) (err error) {

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID

	data := c.Get("basic_inventory_decentralized").(*DecentralizedBasicInventoryDTO)
	var request_data mdm.DecentralizedBasicInventory

	marshaldata, err := json.Marshal(*data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(marshaldata, &request_data)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	request_data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateDecentralizedInventory(query, &request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Decentralized Inventory Details updated succesfully", request_data)
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

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteCentrailizedInventory(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

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

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id

	err = h.service.DeleteDecentralizedInventory(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

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

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetCentrailizedInventory(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Centralized Details Retrieved succesfully", result)
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

	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetDecentralizedInventory(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Decentralized Details Retrieved succesfully", result)
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

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)

	result, err := h.service.GetCentrailizedInventoryList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, " Centralized List Retrieved successfully", result, p)
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

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)

	result, err := h.service.GetCentrailizedInventoryList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, " Centralized List Retrieved successfully", result, p)
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

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetDecentralizedInventoryList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, " Decentralized List Retrieved successfully", result, p)
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

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetDecentralizedInventoryList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, " Decentralized List Retrieved successfully", result, p)
}

// SearchCentralizedInventory godoc
// @Summary SearchCentralizedInventory
// @Description SearchCentralizedInventory
// @Tags Basic_Inventory
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param q query string true "Search query"
// @Success 200 {array} CentralizedSearchObjDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/centralized/search [get]
func (h *handler) SearchCentrailizedInventory(c echo.Context) (err error) {
	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchCentrailizedInventory(q, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "OK", result)
}

// SearchDecentralizedInventory godoc
// @Summary SearchDecentralizedInventory
// @Description SearchDecentralizedInventory
// @Tags Basic_Inventory
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param q query string true "Search query"
// @Success 200 {array} DecentralizedSearchObjDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/basic_inventory/decentralized/search [get]
func (h *handler) SearchDecentralizedInventory(c echo.Context) (err error) {
	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchDecentralizedInventory(q, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "OK", result)
}

// -----------Channel Upsert Api--------------------------------------
func (h *handler) DeCentralisedInventoryUpsert(c echo.Context) (err error) {
	var arrayData []interface{}
	var objectData interface{}
	var data interface{}

	c.Bind(&data)
	token_id := c.Get("TokenUserID").(string)
	jsonData, _ := json.Marshal(data)

	err = json.Unmarshal(jsonData, &arrayData)
	if err != nil {
		_ = json.Unmarshal(jsonData, &objectData)
		arrayData = append(arrayData, objectData)
	}

	msg, err := h.service.UpsertInventoryTemplate(arrayData, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Multiple records excuted successfully", msg)
}

func (h *handler) InventoryTransactionCreate(c echo.Context) (err error) {

	var data = new(mdm.CentralizedInventoryTransactions)

	err = c.Bind(data)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.InventoryTransactionCreate(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Centralized Inventory Transaction created successfully", data)
}
