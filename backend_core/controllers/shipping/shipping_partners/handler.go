package shipping_partners

import (
	"fmt"
	"strconv"

	shipping_base "fermion/backend_core/controllers/shipping/base"
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

func NewHandler() *handler {
	service := NewService()
	base_service := shipping_base.NewServiceBase()
	return &handler{service, base_service}
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
func (h *handler) CreateUserShippingPartnerRegistration(c echo.Context) (err error) {
	data := c.Get("shipping_partners").(*UserShippingPartnerRegistrationRequest)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	userID, _ := strconv.Atoi(token_id)
	t := uint(userID)
	data.CreatedByID = &t
	created_id, err := h.service.CreateUserShippingPartnerRegistration(data, uint(userID), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Shippings partner created", map[string]interface{}{"created_id": created_id})
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
func (h *handler) UpdateUserShippingPartnerRegistration(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data := c.Get("shipping_partners").(*UserShippingPartnerRegistrationRequest)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateUserShippingPartnerRegistration(uint(id), *data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Shippings partner updated", map[string]interface{}{"updated_id": id})
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
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.UserShippingPartnerRegistrationbyid(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
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
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AllUserShippingPartnerRegistration(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "shipping data list data Retrieved successfully", result, p)
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
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AllUserShippingPartnerRegistration(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "shipping data dropdown list data Retrieved successfully", result, p)
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
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteUserShippingPartnerRegistration(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "deleted successfully", map[string]int{"deleted_id": id})
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
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.UserShippingPartnerRegistrationbyid(uint(ID), token_id, access_template_id)
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
