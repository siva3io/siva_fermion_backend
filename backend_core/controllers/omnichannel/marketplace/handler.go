package marketplace

import (
	"fmt"
	"strconv"

	marketplace_base "fermion/backend_core/controllers/omnichannel/base"
	"fermion/backend_core/internal/model/omnichannel"
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
	base_service marketplace_base.ServiceBase
}

var MarketPlaceHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if MarketPlaceHandler != nil {
		return MarketPlaceHandler
	}
	service := NewService()
	base_service := marketplace_base.NewServiceBase()
	MarketPlaceHandler = &handler{service, base_service}
	return MarketPlaceHandler
}

// REGISTER A MARKETPLACE
// CreateMarketplaces godoc
// @Summary      do a CreateMarketplace
// @Description  Create a Marketplace
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   MarketplaceCreateRequest body Marketplace_Request true "create a Marketplace"
// @Success      200  {object}  MarketplaceCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/register [post]
func (h *handler) RegisterMarket(c echo.Context) (err error) {

	data := c.Get("marketplace").(*omnichannel.User_Marketplace_Registration)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.RegisterMarketplace(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "marketplace registered successfully", map[string]interface{}{"created_id": data.ID})
}

// UpdateMarketplaces godoc
// @Summary      do a UpdateMarketplaces
// @Description  Update a Marketplace
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   MarketplaceUpdateRequest body Marketplace_Request true "Update a Marketplace"
// @Success      200  {object}  MarketplaceUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/{id}/edit [post]
func (h *handler) UpdateMarket(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("marketplace").(*omnichannel.User_Marketplace_Registration)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdatemarketplaceDetails(uint(id), *data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "marketplace details updates successfully", map[string]interface{}{"updated_id": id})
}

// FindMarketplaces godoc
// @Summary      Find a Marketplace by id
// @Description  Find a Marketplace by id
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  MarketplaceGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/{id} [get]
func (h *handler) FindMarketplaceByID(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data, err := h.service.FindByidMarketplace(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "data Retrieved successfully", data)
}

// AllMarketplaces godoc
// @Summary      get AllMarketplaces and filter it by search query
// @Description  a get  AllMarketplaces and filter it by search query
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		filters		  query	string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  MarketplaceGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace [get]
func (h *handler) FetchAllMarketplace(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.FindAllMarketplace(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// FetchAllMarketplaceDropDown godoc
// @Summary      get AllMarketplaces dropdown list and filter it by search query
// @Description  get AllMarketplaces dropdown list and filter it by search query
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		filters		  query	string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  MarketplaceGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/dropdown [get]
func (h *handler) FetchAllMarketplaceDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.FindAllMarketplace(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "dropdown data Retrieved successfully", result, p)
}

// DeleteMarketplaces godoc
// @Summary      delete a DeleteMarketplaces
// @Description  delete a DeleteMarketplaces
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  MarketplaceLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/{id}/delete_marketplace [delete]
func (h *handler) DeleteMarket(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteMarketplace(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Marketplace deleted successfully", map[string]int{"deleteD_id": id})
}

// LINK MARKETPLACE DETAILS
// CreateMarketplacesdetails godoc
// @Summary      do a LinkMarketplacedetails
// @Description  Link a Marketplacedetails
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   Marketplace_Details body Marketplace_Details true "link a Marketplace"
// @Success      200  {object}  MarketplacedetailCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/create [post]
func (h *handler) SaveMarketDetails(c echo.Context) (err error) {

	data := c.Get("marketplace_details").(*omnichannel.User_Marketplace_Link)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.SaveMarketplaceDetails(data, token_id, access_template_id)
	// @Summary      Find a Marketplace by id
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateMarketPlaceLinkInCache(token_id, access_template_id)

	return res.RespSuccess(c, "details created successfully", map[string]interface{}{"created_id": data.ID})
}

// UpdateMarketplacesdetails godoc
// @Summary      do a UpdateMarketplacesdetails
// @Description  Update a Marketplacedetails
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param   Marketplace_details body Marketplace_Details true "Update a Marketplace"
// @Success      200  {object}  MarketplacedetailUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/{id}/update [post]
func (h *handler) UpdateMarketDetails(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("marketplace_details").(*omnichannel.User_Marketplace_Link)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateDetails(uint(id), *data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateMarketPlaceLinkInCache(token_id, access_template_id)

	return res.RespSuccess(c, "update details successfully done", map[string]interface{}{"updated_id": id})
}

// UpdateMarketDetailsByQuery godoc
// @Summary      do a UpdateMarketDetailsByQuery
// @Description  Update a Marketplacedetails
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security 	 ApiKeyAuth
// @param 		 Authorization header string true "Authorization"
// @Param		 marketplace_code		  query	string false "marketplace_code"
// @Param   	 Marketplace_details body Marketplace_Details true "Update a Marketplace"
// @Success      200  {object}  MarketplacedetailUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/update [post]
func (h *handler) UpdateMarketDetailsByQuery(c echo.Context) (err error) {

	marketplace_code := c.QueryParam("marketplace_code")
	data := c.Get("marketplace_details").(*omnichannel.User_Marketplace_Link)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	query := map[string]interface{}{
		"market_place_code": marketplace_code,
	}
	err = h.service.UpdateMarketDetailsByQuery(query, *data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "update details successfully done", map[string]interface{}{"updated_id": data.ID})
}

// FindMarketplacesdetails godoc
// @Summary      Find a Marketplacedetails by id
// @Description  Find a Marketplacedetails by id
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  MarketplacedetailGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/view/{id} [get]
func (h *handler) FindMarketDetails(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.FindMarketplaceDetails(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "data Retrieved successdully", result)
}

// AllMarketplacesdetails godoc
// @Summary      get all AllMarketplacesdetails and filter it by search query
// @Description  a get all AllMarketplacesdetails and filter it by search query
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  MarketplacedetailGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/list [get]
func (h *handler) ListMarketDetails(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	var response interface{}
	if *p == pagination.BasePaginatevalue {
		response, *p = GetMarketPlaceLinkFromCache(token_id)
	}

	if response != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.ListMarketplaceDetails(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// ListMarketDetailsDropDown godoc
// @Summary      get AllMarketplacesdetails dropdown list and filter it by search query
// @Description  get AllMarketplacesdetails dropdown list and filter it by search query
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  MarketplacedetailGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/list/dropdown [get]
func (h *handler) ListMarketDetailsDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	var response interface{}
	if *p == pagination.BasePaginatevalue {
		response, *p = GetMarketPlaceLinkFromCache(token_id)
	}

	if response != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.ListMarketplaceDetails(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "dropdown data Retrieved successfully", result, p)
}

// DeleteMarketplacesdetails godoc
// @Summary      delete a DeleteMarketplacesdetails
// @Description  delete a DeleteMarketplacesdetails
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  MarketplacedetailDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/{id}/delete [delete]
func (h *handler) DeleteMarketDetails(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteMarketplaceDetails(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateMarketPlaceLinkInCache(token_id, access_template_id)

	return res.RespSuccess(c, "marketplace details deleted successfully", map[string]int{"deleted_id": id})
}

// AvailableMarketPlaces godoc
// @Summary      get AvailableMarketPlaces and filter it by search query
// @Description  get AvailableMarketPlaces and filter it by search query
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  AvailableMarketPlacesResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/available [get]
func (h *handler) AvailableMarketPlaces(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AvailableMarketPlaces(p, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successfully", result, p)
}

// Get AuthKeys godoc
// @Summary      Get AuthKeys
// @Description  Get AuthKeys
// @Tags         Marketplaces
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        marketplace_code   path      string  true  "marketplace_code"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/marketplace/{marketplace_code}/get_auth_keys [get]
func (h *handler) GetAuthKeys(c echo.Context) (err error) {

	marketplace := c.Param("marketplace_code")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"market_place_code": marketplace,
		"created_by":        user_id,
	}
	result, err := h.service.GetAuthKeys(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// FavouriteMarketPlaces godoc
// @Summary FavouriteMarketPlaces
// @Description FavouriteMarketPlaces
// @Tags Marketplaces
// @Accept  json
// @Produce  json
// @param id path string true "MarketPlace ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/marketplace/{id}/favourite [post]
func (h *handler) FavouriteMarketPlace(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.FindByidMarketplace(uint(ID), token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	fmt.Println("handler ------------------------ -1", query)
	err = h.base_service.FavouriteMarketPlaceFav(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "MarketPlace is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteMarketPlaces godoc
// @Summary UnFavouriteMarketPlaces
// @Description UnFavouriteMarketPlaces
// @Tags Marketplaces
// @Accept  json
// @Produce  json
// @param id path string true "MarketPlace ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/marketplace/{id}/unfavourite [post]
func (h *handler) UnFavouriteMarketPlace(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteMarketPlaceFav(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "MarketPlace is Unmarked as Favourite", map[string]string{"record_id": id})
}

func (h *handler) SaveMarketplace(c echo.Context) (err error) {
	var data = new(omnichannel.Marketplace)
	_ = c.Bind(data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.SaveMarketplace(data, token_id, access_template_id)

	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "details created successfully", map[string]interface{}{"created_id": data.ID})
}

func (h *handler) UpdateMarketplace(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	var data = new(omnichannel.Marketplace)
	_ = c.Bind(data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateMarketplace(uint(id), *data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "update details successfully done", map[string]interface{}{"updated_id": id})
}

func (h *handler) FindMarketplace(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.FindMarketplace(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "data Retrieved successdully", result)
}
