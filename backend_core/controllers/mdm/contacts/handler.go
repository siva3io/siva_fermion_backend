package contacts

import (
	"errors"
	"fmt"
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

var ContactsHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ContactsHandler != nil {
		return ContactsHandler
	}
	service := NewService()
	base_service := mdm_base.NewServiceBase()
	ContactsHandler = &handler{service, base_service}
	return ContactsHandler
}

// CreateContact godoc
// @Summary CreateContact
// @Description CreateContact
// @Tags Contacts
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body PartnerRequestDTO true "Contact Request Body"
// @Success 200 {object} SwaggerCreateOrUpdatePartnerResponseDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/create [post]
func (h *handler) CreateContactEvent(c echo.Context) error {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("partner"),
	}

	eda.Produce(eda.CREATE_CONTACTS, request_payload)

	return res.RespSuccess(c, "Contact Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateContact(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})

	//================= create contact ======================
	createPayload := new(mdm.Partner)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.Create(edaMetaData, createPayload)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_CONTACTS_ACK, responseMessage)
		return
	}

	//================= upsert sub_user =====================
	isAllowedLogin := data["is_allowed_login"]
	if isAllowedLogin != nil && isAllowedLogin == true {
		var createSubUserPayload PartnerRequestDTO
		helpers.JsonMarshaller(data, &createSubUserPayload)
		userId, err := h.service.CreateOrUpdateSubUser(edaMetaData, createSubUserPayload)
		if err != nil {
			fmt.Println("error in create subuser")
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_CONTACTS_ACK, responseMessage)
			return
		}
		edaMetaData.Query = map[string]interface{}{
			"id": createPayload.ID,
		}

		updatePayload := new(mdm.Partner)
		updatePayload.UserId = userId
		err = h.service.Update(edaMetaData, updatePayload)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.UPDATE_CONTACTS_ACK, responseMessage)
			return
		}
	}

	//================= upsert vendor =====================
	isVendorCreate := false
	var contactProperties []map[string]interface{}
	helpers.JsonMarshaller(createPayload.Properties, &contactProperties)
	for _, object := range contactProperties {
		LookupCode, _ := helpers.GetLookupcodeName(object["id"])
		if LookupCode == "VENDOR" {
			isVendorCreate = true
			break
		}
	}
	if isVendorCreate {
		vendor := new(mdm.Vendors)
		if createPayload.FirstName != "" || createPayload.LastName != "" {
			vendor.Name = createPayload.FirstName + " " + createPayload.LastName
		} else {
			if createPayload.CompanyName != "" {
				vendor.Name = createPayload.CompanyName
			}
		}
		vendor.ContactId = &createPayload.ID
		vendor.PrimaryContactId = createPayload.ParentId

		vendorPayload := []interface{}{}

		vendorPayload = append(vendorPayload, vendor)

		request_payload := map[string]interface{}{
			"meta_data": edaMetaData,
			"data":      vendorPayload,
		}
		eda.Produce(eda.UPSERT_VENDOR, request_payload)
	}

	//===============cache implementation=====================
	UpdateContactInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}

	// eda.Produce(eda.CREATE_CONTACTS_ACK, responseMessage)
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
// @param RequestBody body PartnerRequestDTO true "Contact  Request Body"
// @Success 200 {object} SwaggerCreateOrUpdatePartnerResponseDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/contacts/{id}/update [post]
func (h *handler) UpdateContactEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	contactId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         contactId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("partner"),
	}

	eda.Produce(eda.UPDATE_CONTACTS, request_payload)
	return res.RespSuccess(c, "Contact Updation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateContact(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})

	updatePayload := new(mdm.Partner)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.Update(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_CONTACTS_ACK, responseMessage)
		return
	}
	//================= upsert sub_user =====================
	isAllowedLogin := data["is_allowed_login"]
	if isAllowedLogin != nil && isAllowedLogin == true {
		contactDetails, err := h.service.FindOne(edaMetaData)
		if err != nil {
			fmt.Println("error in create subuser")
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_CONTACTS_ACK, responseMessage)
			return
		}

		var createSubUserPayload PartnerRequestDTO
		helpers.JsonMarshaller(data, &createSubUserPayload)
		helpers.JsonMarshaller(contactDetails, &createSubUserPayload)

		userId, err := h.service.CreateOrUpdateSubUser(edaMetaData, createSubUserPayload)
		if err != nil {
			fmt.Println("error in create subuser")
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.CREATE_CONTACTS_ACK, responseMessage)
			return
		}

		updatePayload := new(mdm.Partner)
		updatePayload.UserId = userId
		err = h.service.Update(edaMetaData, updatePayload)
		if err != nil {
			responseMessage.ErrorMessage = err
			// eda.Produce(eda.UPDATE_CONTACTS_ACK, responseMessage)
			return
		}
	}

	//================= upsert vendor =====================
	isVendorCreate := false

	contactResponse, _ := h.service.FindOne(edaMetaData)
	var contactDetails PartnerResponseDTO
	helpers.JsonMarshaller(contactResponse, &contactDetails)
	for _, object := range contactDetails.Properties {
		LookupCode, _ := helpers.GetLookupcodeName(object.Id)
		if LookupCode == "VENDOR" {
			isVendorCreate = true
			break
		}
	}
	if isVendorCreate {
		vendor := new(mdm.Vendors)
		if contactDetails.FirstName != "" || contactDetails.LastName != "" {
			vendor.Name = contactDetails.FirstName + " " + contactDetails.LastName
		} else {
			if contactDetails.CompanyName != "" {
				vendor.Name = contactDetails.CompanyName
			}
		}
		vendor.ContactId = &contactDetails.Id
		vendor.PrimaryContactId = contactDetails.ParentId
		vendorPayload := []interface{}{}
		vendorPayload = append(vendorPayload, vendor)
		request_payload := map[string]interface{}{
			"meta_data": edaMetaData,
			"data":      vendorPayload,
		}
		eda.Produce(eda.UPSERT_VENDOR, request_payload)
	}

	//===============cache implementation===================================
	UpdateContactInCache(edaMetaData)

	responseMessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}

	// eda.Produce(eda.UPDATE_CONTACTS_ACK, responseMessage)
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
	metaData := c.Get("MetaData").(core.MetaData)
	contactId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         contactId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}

	err = h.service.Delete(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	//===============cache implementation===================================
	UpdateContactInCache(metaData)

	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": contactId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	contactId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id": contactId,
		// "company_id": metaData.CompanyId,
	}

	result, err := h.service.FindOne(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	var response PartnerResponseDTO
	helpers.JsonMarshaller(result, &response)

	for _, address := range response.AddressDetails {
		if address.MarkDefaultAddress {
			response.DefaultAddress = address
			break
		}
	}
	return res.RespSuccess(c, "Contact Details Retrieved successfully", response)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []PartnerResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetContactFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "Contact list Retrieved Successfully", response, p)
	}

	result, err := h.service.FindAll(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)

	for index, ContactDetails := range response {
		for _, address := range ContactDetails.AddressDetails {
			if address.MarkDefaultAddress {
				response[index].DefaultAddress = address
				break
			}
		}
	}
	return res.RespSuccessInfo(c, "Contact list Retrieved Successfully", response, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []PartnerResponseDTO

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetContactFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "Contact list Retrieved Successfully", response, p)
	}

	result, err := h.service.FindAll(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)

	for index, ContactDetails := range response {
		for _, address := range ContactDetails.AddressDetails {
			if address.MarkDefaultAddress {
				response[index].DefaultAddress = address
				break
			}
		}
	}
	return res.RespSuccessInfo(c, "Contact list Retrieved Successfully", response, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	contactId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         contactId,
		"company_id": metaData.CompanyId,
	}

	result, err := h.service.FindOne(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Contact Details Retrieved successfully", result)
}

// GetContactTab godoc
// @Summary GetContactTab
// @Summary GetContactTab
// @Description GetContactTab
// @Tags Contacts
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
// @Router /api/v1/contacts/{id}/filter_module/{tab} [get]
func (h *handler) GetContactTab(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	p := new(pagination.Paginatevalue)
	err = c.Bind(p)
	if err != nil {
		return res.RespErr(c, err)
	}

	contactId := c.Param("id")
	tabName := c.Param("tab")

	metaData.AdditionalFields = map[string]interface{}{
		"contact_id": contactId,
		"tab_name":   tabName,
	}

	data, err := h.service.GetRelatedTabs(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, p)
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
	metaData := c.Get("MetaData").(core.MetaData)
	contactId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         contactId,
		"company_id": metaData.CompanyId,
	}

	_, err = h.service.FindOne(metaData)
	if err != nil {
		return res.RespSuccess(c, "Specified record not found", err)
	}

	query := map[string]interface{}{
		"ID":      contactId,
		"user_id": metaData.TokenUserId,
	}
	err = h.base_service.FavouriteContacts(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Contact is Marked as Favourite", map[string]string{"contact_id": contactId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	contactId := c.Param("id")

	query := map[string]interface{}{
		"ID":      contactId,
		"user_id": metaData.TokenUserId,
	}

	err = h.base_service.UnFavouriteContacts(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Contact is Unmarked as Favourite", map[string]string{"contact_id": contactId})
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
	metaData := c.Get("MetaData").(core.MetaData)
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	tokenUserId := int(metaData.TokenUserId)
	data, err := h.base_service.GetArrContact(tokenUserId)
	if err != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavContactList(data, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Contact list Retrieved Successfully", result, p)
}

func (h *handler) UpsertContactEvent(c echo.Context) (err error) {
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
	eda.Produce(eda.UPSERT_CONTACT, request_payload)

	return res.RespSuccess(c, "Contact Upsert Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}
func (h *handler) UpsertContact(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	var data []interface{}
	helpers.JsonMarshaller(request["data"], &data)

	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData

	response, err := h.service.Upsert(edaMetaData, data)
	helpers.PrettyPrint("contacts upsert response", response)

	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPSERT_CONTACTS_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"response": response,
	}

	// eda.Produce(eda.UPSERT_CONTACTS_ACK, responseMessage)
}
