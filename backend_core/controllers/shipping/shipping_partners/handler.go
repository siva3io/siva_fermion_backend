package shipping_partners

import (
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	shipping_base "fermion/backend_core/controllers/shipping/base"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
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
	base_service shipping_base.ServiceBase
}

var ShippingPartnersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if ShippingPartnersHandler != nil {
		return ShippingPartnersHandler
	}
	service := NewService()
	base_service := shipping_base.NewServiceBase()
	ShippingPartnersHandler = &handler{service, base_service}
	return ShippingPartnersHandler
}

// Create UserShippingPartnerRegistration godoc
// @Summary      CreateUserShippingPartnerRegistration
// @Description  Create a UserShippingPartnerRegistration
// @Tags         UserShippingPartnerRegistration
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        UserShippingPartnerRegistrationCreateRequest body UserShippingPartnerRegistrationRequest true "create a Shipping Partner"
// @Success      200  {object}  UserShippingPartnerRegistrationCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/register [post]
func (h *handler) CreateUserShippingPartnerRegistrationEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("shipping_partners"),
	}
	eda.Produce(eda.CREATE_SHIPPING_PARTNER, request_payload)
	return res.RespSuccess(c, "Shipping Partner Creation Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) CreateUserShippingPartnerRegistration(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(shipping.UserShippingPartnerRegistration)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.CreateUserShippingPartnerRegistration(edaMetaData, createPayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_SHIPPING_PARTNER_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{"created_id": createPayload.ID}
	// eda.Produce(eda.CREATE_SHIPPING_PARTNER_ACK, responseMessage)
}

// UpdateUserShippingPartnerRegistration godoc
// @Summary      UpdateUserShippingPartnerRegistration
// @Description  UpdateUserShippingPartnerRegistration
// @Tags         UserShippingPartnerRegistration
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param        UserShippingPartnerRegistrationUpdateRequest body UserShippingPartnerRegistrationRequest true "Update a Shipping Partner"
// @Success      200  {object}  UserShippingPartnerRegistrationUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/{id}/update [post]
func (h *handler) UpdateUserShippingPartnerRegistrationEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	shippingPartnerId := c.Param("id")

	edaMetaData.Query = map[string]interface{}{
		"id":         shippingPartnerId,
		"company_id": edaMetaData.CompanyId,
	}

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("shipping_partners"),
	}
	eda.Produce(eda.UPDATE_SHIPPING_PARTNER, request_payload)
	return res.RespSuccess(c, "Shipping_partner Update Inprogress", map[string]interface{}{"request_id": edaMetaData.RequestId})
}

func (h *handler) UpdateUserShippingPartnerRegistration(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})
	updatePayload := new(shipping.UserShippingPartnerRegistration)
	helpers.JsonMarshaller(data, updatePayload)

	err := h.service.UpdateUserShippingPartnerRegistration(edaMetaData, updatePayload)
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_SHIPPING_PARTNER_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{"updated_id": edaMetaData.Query["id"]}
	// eda.Produce(eda.UPDATE_SHIPPING_PARTNER_ACK, responseMessage)
	// cache implementation
	UpdateUserLogisticPartnerInCache(edaMetaData)

}

// GetUserShippingPartnerRegistration godoc
// @Summary      FindUserShippingPartnerRegistrationById
// @Description  Find a UserShippingPartnerRegistration by id
// @Tags         UserShippingPartnerRegistration
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  UserShippingPartnerRegistrationGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/registered/{id} [get]
func (h *handler) GetUserShippingPartnerRegistration(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	shippingPartnerId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         shippingPartnerId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.UserShippingPartnerRegistrationbyid(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	var response UserShippingPartnerRegistrationRequest
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccess(c, "data Retrieved successfully", result)
}

// GetAllUserShippingPartnerRegistration godoc
// @Summary     FindAllUserShippingPartnerRegistration
// @Description  FindAllUserShippingPartnerRegistration
// @Tags         UserShippingPartnerRegistration
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   	per_page 	      query int false "per_page"
// @Param  		page_no           query int false "page_no"
// @Param   	sort		      query string false "sort"
// @param		filters 	      query string false "Filters"
// @Success      200  {object}  UserShippingPartnerRegistrationGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners [get]
func (h *handler) GetAllUserShippingPartnerRegistration(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []UserShippingPartnerRegistrationRequest

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetUserLogisticPartnerFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.AllUserShippingPartnerRegistration(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "shipping data list data Retrieved successfully", response, p)
}

// GetAllUserShippingPartnerRegistrationDropDown godoc
// @Summary     GetAllUserShippingPartnerRegistrationDropDown
// @Description  GetAllUserShippingPartnerRegistrationDropDown
// @Tags         UserShippingPartnerRegistrationDropDown
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   	per_page 	      query int false "per_page"
// @Param  		page_no           query int false "page_no"
// @Param   	sort		      query string false "sort"
// @param		filters 	      query string false "Filters"
// @Success      200  {object}  UserShippingPartnerRegistrationGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/dropdown [get]
func (h *handler) GetAllUserShippingPartnerRegistrationDropDown(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "DROPDOWN_LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	var cacheResponse interface{}
	var response []UserShippingPartnerRegistrationRequest

	tokenUserId := strconv.Itoa(int(metaData.TokenUserId))
	if *p == pagination.BasePaginatevalue {
		cacheResponse, *p = GetUserLogisticPartnerFromCache(tokenUserId)
	}

	if cacheResponse != nil {
		helpers.JsonMarshaller(cacheResponse, &response)
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.AllUserShippingPartnerRegistration(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	helpers.JsonMarshaller(result, &response)
	return res.RespSuccessInfo(c, "shipping data list data Retrieved successfully", response, p)
}

// DeleteUserShippingPartnerRegistration godoc
// @Summary      DeleteUserShippingPartnerRegistration
// @Description  Delete a UserShippingPartnerRegistration
// @Tags         UserShippingPartnerRegistration
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  UserShippingPartnerRegistrationDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/{id}/delete [delete]
func (h *handler) DeleteUserShippingPartnerRegistration(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	shippingPartnerId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         shippingPartnerId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
	}
	err = h.service.DeleteUserShippingPartnerRegistration(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateUserLogisticPartnerInCache(metaData)

	return res.RespSuccess(c, "deleted successfully", shippingPartnerId)
}

// FavouriteUserShippingPartnerRegistration godoc
// @Summary FavouriteUserShippingPartnerRegistration
// @Description FavouriteUserShippingPartnerRegistration
// @Tags UserShippingPartnerRegistration
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path string true "UserShippingPartnerRegistration ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_partners/{id}/favourite [post]
func (h *handler) FavouriteUserShippingPartnerRegistration(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	metaData := c.Get("MetaData").(core.MetaData)

	_, er := h.service.UserShippingPartnerRegistrationbyid(metaData)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteUserShippingPartnerRegistration(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Shipping Partner is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteUserShippingPartnerRegistration godoc
// @Summary UnFavouriteUserShippingPartnerRegistration
// @Description UnFavouriteUserShippingPartnerRegistration
// @Tags UserShippingPartnerRegistration
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id path string true "UserShippingPartnerRegistration ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_partners/{id}/unfavourite [post]
func (h *handler) UnFavouriteUserShippingPartnerRegistration(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteUserShippingPartnerRegistration(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Shipping Partner is marked as Unfavourite", map[string]string{"record_id": id})
}

// UserShippingPartnerRegistrationEstimateCosts godoc
// @Summary      UserShippingPartnerRegistrationEstimateCosts
// @Description  UserShippingPartnerRegistrationEstimateCosts
// @Tags         UserShippingPartnerRegistration
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        RateCalcCreateRequest body RateCalcRequest true "get all rate calculator cost"
// @Success      200  {object}  GetAllRateCalculator
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/rate_calc [post]
func (h *handler) ShippingPartnerEstimateCosts(c echo.Context) (err error) {
	data := new(shipping.RateCalculator)
	c.Bind(&data)
	fmt.Println(data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.ShippingPartnerEstimateCosts(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Estimated Cost of Shipping Partners is Fetched Successfully", result)
}

// DownloadUserShippingPartnerRegistrationPDF godoc
// @Summary Download UserShippingPartnerRegistration
// @Description Download UserShippingPartnerRegistration
// @Tags UserShippingPartnerRegistration
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UserShippingPartnerRegistration_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_partners/{id}/download_pdf [post]
func (h *handler) DownloadUserShippingPartnerRegistrationPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

// UploadUserShippingPartnerRegistrationDocument godoc
// @Summary Upload UserShippingPartnerRegistration Document
// @Description Upload UserShippingPartnerRegistration Document
// @Tags UserShippingPartnerRegistration
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "UserShippingPartnerRegistration_id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/shipping_partners/{id}/upload_document [post]
func (h *handler) UploadUserShippingPartnerRegistrationDocument(c echo.Context) error {
	file_name := c.QueryParam("file_name")
	return res.RespSuccess(c, "Document uploaded successfully", map[string]interface{}{"file_name": file_name})
}

// GetAllShippingPartner godoc
// @Summary     FindAllShippingPartner
// @Description  FindAllShippingPartner
// @Tags         ShippingPartner
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param       per_page           query int false "per_page"
// @Param          page_no           query int false "page_no"
// @Param       sort              query string false "sort"
// @param        filters           query string false "Filters"
// @Success      200  {object}  ShippingPartnerGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/list [get]
func (h *handler) GetAllShippingPartner(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AllShippingPartner(p, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// GetShippingPartnerAuthById godoc
// @Summary      GetShippingPartnerAuthById
// @Description  Find a GetShippingPartnerAuth by id
// @Tags         ShippingPartner
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/{partner_name}/get_auth_options [get]
func (h *handler) GetShippingPartnerAuth(c echo.Context) (err error) {
	Name := c.Param("partner_name")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetShippingPartnerAuthByName(Name, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "data retrieved successfully", result)
}

// GetShippingPartnerByName godoc
// @Summary      GetShippingPartnerByName
// @Description  Find a GetShippingPartner by name
// @Tags         ShippingPartner
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        partner_name   path      string  true  "partner_name"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/{partner_name}/update_partner [get]
func (h *handler) UpdateShippingPartnerByName(c echo.Context) (err error) {
	Name := c.Param("partner_name")
	var data shipping.ShippingPartner
	c.Bind(&data)
	//fmt.Println(Name)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	query := map[string]interface{}{
		"partner_name": Name,
		"data":         data,
	}
	response, _ := h.service.UpdateShippingPartnerByName(query, token_id, access_template_id)

	return res.RespSuccess(c, "Shipping Partner Updated Successfully", response)
}

// GetShippingPartnerById godoc
// @Summary      GetShippingPartnerById
// @Description  Find a GetShippingPartner by id
// @Tags         ShippingPartner
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/shipping_partners/{id} [get]
func (h *handler) GetShippingPartner(c echo.Context) (err error) {
	pr := c.Param("id")
	id, _ := strconv.Atoi(pr)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetShippingPartnerById(id, token_id, access_template_id)

	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "data retrieved successfully", result)
}
