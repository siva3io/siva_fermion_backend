package debitnote

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/pagination"

	accounting_base "fermion/backend_core/controllers/accounting/base"
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
	service      Service
	base_service accounting_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := accounting_base.NewServiceBase()
	return &handler{service, base_service}
}

// CreateDebitNote godoc
// @Summary CreateDebitNote
// @Description CreateDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body DebitNoteDTO true "Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/create [post]
func (h *handler) CreateDebitNote(c echo.Context) (err error) {
	data := new(accounting.DebitNote)
	dto_data := c.Get("debit_note").(*DebitNoteDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	dto_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	deb_seq := helpers.GenerateSequence("Acc", token_id, "debit_notes")
	ref_seq := helpers.GenerateSequence("Ref", token_id, "debit_notes")
	if dto_data.GenerateDebitNoteId {
		dto_data.DebitNoteID = deb_seq
	}
	if dto_data.GenerateReferenceId {
		dto_data.ReferenceId = ref_seq
	}

	byte_data, _ := json.Marshal(dto_data)
	json.Unmarshal(byte_data, &data)
	err = h.service.CreateDebitNote(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Debit note created successfully", map[string]interface{}{"created_id": data.ID, "Details": data})
}

// UpdateDebitNote godoc
// @Summary UpdateDebitNote
// @Description UpdateDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @param id path string true "Debit Note ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body DebitNoteDTO true "DebitNote  Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/update [post]
func (h *handler) UpdateDebitNote(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("debit_note").(*DebitNoteDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateDebitNote(query, data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Debit Note updated succesfully", data)
}

// DebitNoteView godoc
// @Summary View DebitNote
// @Summary View DebitNote
// @Description View DebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "DebitNote ID"
// @Success 200 {object} DebitNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id} [get]
func (h *handler) GetDebitNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetDebitNote(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Debit Note Details Retrieved successfully", result)
}

// GetAllDebitNote godoc
// @Summary Get all Debit note
// @Description Get all Debit note
// @Tags DebitNote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} DebitNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote [get]
func (h *handler) GetAllDebitNote(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetDebitNoteList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Debit Note list Retrieved Successfully", result, p)
}

// GetAllDebitNoteDropDown godoc
// @Summary Get Debit note dropdown list
// @Description Get Debit note dropdown list
// @Tags DebitNote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} DebitNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/dropdown [get]
func (h *handler) GetAllDebitNoteDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetDebitNoteList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Debit Note dropdown list Retrieved Successfully", result, p)
}

// DeleteDebitNote godoc
// @Summary DeleteDebitNote
// @Description DeleteDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "DebitNote ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/delete [delete]
func (h *handler) DeleteDebitNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteDebitNote(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// FavouriteDebitNote godoc
// @Summary FavouriteDebitNote
// @Description FavouriteDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @param id path string true "DebitNote ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/favourite [post]
func (h *handler) FavouriteDebitNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//var token = Token
	//userId, _ := helpers.ExtractToken(token)
	//var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id) //to be extracted from token payload
	fmt.Println(user_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	q := map[string]interface{}{
		"id": ID,
	}
	_, er := h.service.GetDebitNote(q, token_id, access_template_id)
	if er != nil {
		//fmt.Println(er.Error())
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteDebitNote(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Debit Note is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteDebitNote godoc
// @Summary UnFavouriteDebitNote
// @Description UnFavouriteDebitNote
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @param id path string true "Contact ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/unfavourite [post]
func (h *handler) UnFavouriteDebitNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//var token = Token
	//userId, _ := helpers.ExtractToken(token)
	//var query = make(map[string]interface{}, 0)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteDebitNote(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Debit Note is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetDebitNoteTab godoc
// @Summary GetDebitNoteTab
// @Summary GetDebitNoteTab
// @Description GetDebitNoteTab
// @Tags DebitNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @param tab path string true "tab"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/debitnote/{id}/filter_module/{tab} [get]
func (h *handler) GetDebitNoteTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetDebitNoteTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
