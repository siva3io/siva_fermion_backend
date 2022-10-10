package contacts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	mdm_base "fermion/backend_core/controllers/mdm/base"

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
	service      Service
	base_service mdm_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := mdm_base.NewServiceBase()
	return &handler{service, base_service}
}

// CreateContact godoc
// @Summary CreateContact
// @Description CreateContact
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body SwaggerPartnerRequestDTO true "Contact Request Body"
// @Success 200 {object} PartnerResponseDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/create [post]
func (h *handler) CreateContact(c echo.Context) (err error) {
	var request_data mdm.Partner

	partner := c.Get("partner").(*PartnerRequestDTO)
	s := c.Get("TokenUserID").(string)
	partner.CreatedByID = helpers.ConvertStringToUint(s)

	fmt.Println("=====>>>", &partner.CreatedByID)

	marshaldata, err := json.Marshal(*partner)
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
	err = h.service.CreateContact(&request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	isAllowedLogin := partner.IsAllowedLogin

	//------------check IsAllowedLogin true or false-----------------
	if isAllowedLogin != nil && *isAllowedLogin {
		err = h.service.CreateOrUpdateSubUser(*partner)
		if err != nil {
			return res.RespErr(c, errors.New("error in create subuser"))
		}
	}

	return res.RespSuccess(c, "Contact created successfully", request_data)
}

// UpdateContact godoc
// @Summary UpdateContact
// @Description UpdateContact
// @Tags Contacts
// @Accept  json
// @Produce  json
// @param id path string true "Contact ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body SwaggerPartnerRequestDTO true "Contact  Request Body"
// @Success 200 {object} PartnerResponseDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/{id}/update [post]
func (h *handler) UpdateContact(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	partner := c.Get("partner").(*PartnerRequestDTO)
	s := c.Get("TokenUserID").(string)
	partner.UpdatedByID = helpers.ConvertStringToUint(s)
	var request_data mdm.Partner

	marshaldata, err := json.Marshal(*partner)
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
	err = h.service.UpdateContact(query, &request_data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	isAllowedLogin := partner.IsAllowedLogin

	//------------check IsAllowedLogin true or false-----------------
	if isAllowedLogin != nil && *isAllowedLogin {
		err = h.service.CreateOrUpdateSubUser(*partner)
		if err != nil {
			return res.RespErr(c, errors.New("error in create subuser"))
		}
	}
	return res.RespSuccess(c, "Contact updated succesfully", request_data)
}

// DeleteContacts godoc
// @Summary DeleteContacts
// @Description DeleteContacts
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Contact ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/{id}/delete [delete]
func (h *handler) DeleteContacts(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteContact(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// ContactView godoc
// @Summary View Contact
// @Summary View Contact
// @Description View Contact
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Contact ID"
// @Success 200 {object} ContactViewDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/{id} [get]
func (h *handler) ContactView(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetContact(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Contact Details Retrieved successfully", result)
}

// GetRelatedContacts godoc
// @Summary Get related contacts
// @Description Get related contacts
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "Contact ID"
// @Success 200 {object} ContactListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/{id}/related [get]
func (h *handler) GetRelatedContacts(c echo.Context) error {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetContact(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Contact Details Retrieved successfully", result)
}

// GetContacts godoc
// @Summary Get all contacts
// @Description Get all contacts
// @Tags Contacts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} ContactListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts [get]
func (h *handler) GetContacts(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetContactList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Contact list Retrieved Successfully", result, p)
}

// GetContactsDropdown godoc
// @Summary Get all contacts
// @Description Get all contacts
// @Tags Contacts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} ContactListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/dropdown [get]
func (h *handler) GetContactsDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	result, err := h.service.GetContactList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Contact list Retrieved Successfully", result, p)
}

// SearchContacts godoc
// @Summary Get all filtered  contacts
// @Description Get all filtered contacts
// @Tags Contacts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param q query string true "Search query"
// @Success 200 {array} IdNameDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/search [get]
func (h *handler) SearchContacts(c echo.Context) (err error) {
	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchContact(q, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "OK", result)
}

// FavouriteContacts godoc
// @Summary FavouriteContacts
// @Description FavouriteContacts
// @Tags Contacts
// @Accept  json
// @Produce  json
// @param id path string true "Contact ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/{id}/favourite [post]
func (h *handler) FavouriteContacts(c echo.Context) (err error) {
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
	_, er := h.service.GetContact(q, token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteContacts(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Contact is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteContacts godoc
// @Summary UnFavouriteContacts
// @Description UnFavouriteContacts
// @Tags Contacts
// @Accept  json
// @Produce  json
// @param id path string true "Contact ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/{id}/unfavourite [post]
func (h *handler) UnFavouriteContacts(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteContacts(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Contact is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetFavouriteContacts godoc
// @Summary Get all favourite contacts
// @Description Get all favourite contacts
// @Tags Contacts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} ContactListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/favourite_list [get]
func (h *handler) FavouriteContactsView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrContact(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavContactList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Contact list Retrieved Successfully", result, p)
}
