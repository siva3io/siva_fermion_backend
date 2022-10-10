package uom

import (
	"encoding/json"
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

// CreateUom godoc
// @Summary CreateUom
// @Description CreateUom
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomRequestDTO true "Uom Request Body"
// @Success 200 {object} UomRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/create [post]
func (h *handler) CreateUom(c echo.Context) (err error) {
	uom := c.Get("uom").(*UomRequestDTO)
	var request_data mdm.Uom

	marshaldata, err := json.Marshal(*uom)
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
	err = h.service.CreateUom(&request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "UOM created successfully", request_data)
}

// CreateUomClass godoc
// @Summary CreateUomClass
// @Description CreateUomClass
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomClassRequestDTO true "Uom Class Request Body"
// @Success 200 {object} UomClassRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/create [post]
func (h *handler) CreateUomClass(c echo.Context) (err error) {
	uomClass := c.Get("uom_class").(*UomClassRequestDTO)
	var request_data mdm.UomClass

	marshaldata, err := json.Marshal(*uomClass)
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
	err = h.service.CreateUomClass(&request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "UomClass created successfully", request_data)
}

// UpdateUom godoc
// @Summary UpdateUom
// @Description UpdateUom
// @Tags UOM
// @Accept  json
// @Produce  json
// @param id path string true "UOM ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomRequestDTO true "Uom Request Body"
// @Success 200 {object} UomRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/{id}/update [post]
func (h *handler) UpdateUom(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("uom").(*UomRequestDTO)
	var request_data mdm.Uom

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
	err = h.service.UpdateUom(query, &request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "UOM Details updated succesfully", request_data)
}

// UpdateUomClass godoc
// @Summary UpdateUomClass
// @Description UpdateUomClass
// @Tags UOM
// @Accept  json
// @Produce  json
// @param id path string true "UOM class ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body UomClassRequestDTO true "Uom class Request Body"
// @Success 200 {object} UomClassRequestDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/{id}/update [post]
func (h *handler) UpdateUomClass(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("uom_class").(*UomClassRequestDTO)
	var request_data mdm.UomClass
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	marshaldata, err := json.Marshal(*data)
	if err != nil {
		return res.RespErr(c, err)
	}
	err = json.Unmarshal(marshaldata, &request_data)
	if err != nil {
		return res.RespErr(c, err)
	}
	s := c.Get("TokenUserID").(string)
	request_data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateUomClass(query, &request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "UOMClass Details updated succesfully", request_data)
}

// DeleteUom godoc
// @Summary DeleteUom
// @Description DeleteUom
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/{id}/delete [delete]
func (h *handler) DeleteUom(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteUom(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// DeleteUomClass godoc
// @Summary DeleteUomClass
// @Description DeleteUomClass
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM class ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/{id}/delete [delete]
func (h *handler) DeleteUomClass(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteUomClass(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// UomView godoc
// @Summary View UomView
// @Summary View UomView
// @Description View UomView
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM ID"
// @Success 200 {object} UomViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/{id} [get]
func (h *handler) GetUom(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetUom(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Uom DEtails retried succesfully", result)
}

// UomClassView godoc
// @Summary View UomClassView
// @Summary View UomClassView
// @Description View UomClassView
// @Tags UOM
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UOM class ID"
// @Success 200 {object} UomClassViewtDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/{id} [get]
func (h *handler) GetUomClass(c echo.Context) error {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetUomClass(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "UomClass Details retried succesfully", result)
}

// UomList godoc
// @Summary UomList
// @Description UomList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom [get]
func (h *handler) GetUomList(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetUomList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// UomListDropdown godoc
// @Summary UomList
// @Description UomList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/dropdown [get]
func (h *handler) GetUomListDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetUomList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// UomClassList godoc
// @Summary UomClassList
// @Description UomClassList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomClassListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class [get]
func (h *handler) GetUomClassList(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetUomClassList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// UomClassListDropdown godoc
// @Summary UomClassList
// @Description UomClassList
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} UomClassListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/dropdown [get]
func (h *handler) GetUomClassListDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetUomClassList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", result, p)
}

// UomSearch godoc
// @Summary UomSearch
// @Description UomSearch
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param q query string true "Search query"
// @Success 200 {array} IdNameDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/search [get]
func (h *handler) SearchUom(c echo.Context) (err error) {
	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchUom(q, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "OK", result)
}

// UomClassSearch godoc
// @Summary UomClassSearch
// @Description UomClassSearch
// @Tags UOM
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param q query string true "Search query"
// @Success 200 {array} IdNameDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/uom/class/search [get]
func (h *handler) SearchUomClass(c echo.Context) (err error) {
	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchUomClass(q, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "OK", result)
}
