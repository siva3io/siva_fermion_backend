package creditnote

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

// CreateCreditNote godoc
// @Summary CreateCreditNote
// @Description CreateCreditNote
// @Tags CreditNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body CreditNoteDTO true "Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote/create [post]
func (h *handler) CreateCreditNote(c echo.Context) (err error) {
	data := new(accounting.CreditNote)
	dto_data := c.Get("credit_note").(*CreditNoteDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	dto_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	seq_credit := helpers.GenerateSequence("Cred", token_id, "credit_notes")
	seq_ref := helpers.GenerateSequence("Ref", token_id, "credit_notes")
	if dto_data.GenerateCreditNoteId {
		dto_data.CreditNoteID = seq_credit
	}
	if dto_data.GenerateReferenceId {
		dto_data.ReferenceId = seq_ref
	}
	byte_data, _ := json.Marshal(dto_data)
	json.Unmarshal(byte_data, &data)
	err = h.service.CreateCreditNote(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Credit note created successfully", map[string]interface{}{"created_id": data.ID, "Details": data})
}

// UpdateCreditNote godoc
// @Summary UpdateCreditNote
// @Description UpdateCreditNote
// @Tags CreditNote
// @Accept  json
// @Produce  json
// @param id path string true "Credit Note ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body CreditNoteDTO true "CreditNote  Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote/{id}/update [post]
func (h *handler) UpdateCreditNote(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("credit_note").(*CreditNoteDTO)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateCreditNote(query, data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Credit Note updated succesfully", data)
}

// CreditNoteView godoc
// @Summary View CreditNote
// @Summary View CreditNote
// @Description View CreditNote
// @Tags CreditNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "CreditNote ID"
// @Success 200 {object} CreditNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote/{id} [get]
func (h *handler) GetCreditNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	result, err := h.service.GetCreditNote(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Credit Note Details Retrieved successfully", result)
}

// GetAllCreditNote godoc
// @Summary Get all Credit note
// @Description Get all Credit note
// @Tags CreditNote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CreditNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote [get]
func (h *handler) GetAllCreditNote(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetCreditNoteList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Credit Note dropdown list Retrieved Successfully", result, p)
}

// GetAllCreditNoteDropDown godoc
// @Summary Get all Credit note dropdown
// @Description Get all Credit note dropdown
// @Tags CreditNote
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} CreditNoteDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote/dropdown [get]
func (h *handler) GetAllCreditNoteDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetCreditNoteList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Credit Note dropdown list Retrieved Successfully", result, p)
}

// DeleteCreditNote godoc
// @Summary DeleteCreditNote
// @Description DeleteCreditNote
// @Tags CreditNote
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "CreditNote ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote/{id}/delete [delete]
func (h *handler) DeleteCreditNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteCreditNote(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// FavouriteCreditNote godoc
// @Summary FavouriteCreditNote
// @Description FavouriteCreditNote
// @Tags CreditNote
// @Accept  json
// @Produce  json
// @param id path string true "CreditNote ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote/{id}/favourite [post]
func (h *handler) FavouriteCreditNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	fmt.Println(user_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	q := map[string]interface{}{
		"id": ID,
	}
	access_template_id := c.Get("AccessTemplateId").(string)
	_, er := h.service.GetCreditNote(q, token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteCreditNote(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Credit Note is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteCreditNote godoc
// @Summary UnFavouriteCreditNote
// @Description UnFavouriteCreditNote
// @Tags CreditNote
// @Accept  json
// @Produce  json
// @param id path string true "Contact ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/creditnote/{id}/unfavourite [post]
func (h *handler) UnFavouriteCreditNote(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteCreditNote(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Credit Note is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetCreditNoteTab godoc
// @Summary GetCreditNoteTab
// @Summary GetCreditNoteTab
// @Description GetCreditNoteTab
// @Tags CreditNote
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
// @Router /api/v1/creditnote/{id}/filter_module/{tab} [get]
func (h *handler) GetCreditNoteTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetCreditNoteTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
