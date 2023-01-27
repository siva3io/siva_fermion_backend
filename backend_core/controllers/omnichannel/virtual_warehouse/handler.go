package virtual_warehouse

import (
	"fmt"
	"strconv"

	virtual_warehouse_base "fermion/backend_core/controllers/omnichannel/base"
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
	base_service virtual_warehouse_base.ServiceBase
}

var VirtualWarehouseHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if VirtualWarehouseHandler != nil {
		return VirtualWarehouseHandler
	}
	service := NewService()
	base_service := virtual_warehouse_base.NewServiceBase()
	VirtualWarehouseHandler = &handler{service, base_service}
	return VirtualWarehouseHandler
}

// REGISTER VIRTUAL WAREHOUSE
// CreateVirtualWarehouse godoc
// @Summary      do a CreateVirtualWarehouse
// @Description  Create a VirtualWarehouse
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Param        VirtualWarehouseCreateRequest body VirtualWarehouse_Request true "create a VirtualWarehouse"
// @Success      200  {object}  VirtualWarehouseCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/register [post]
func (h *handler) RegisterVirtualWarehouse(c echo.Context) (err error) {
	data := c.Get("virtual_warehouse").(*omnichannel.User_Virtual_Warehouse_Registration)
	s := c.Get("TokenUserID").(string)
	data.CreatedByID = helpers.ConvertStringToUint(s)
	err = h.service.RegisterVirtualWarehouse(data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Virtual Warehouse Registered Successfully", map[string]interface{}{"created_id": data.ID})
}

// UpdateVirtualWarehouse godoc
// @Summary      do a UpdateVirtualWarehouse
// @Description  Update a VirtualWarehouse
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param        VirtualWarehouseUpdateRequest body VirtualWarehouse_Request true "Update a VirtualWarehouse"
// @Success      200  {object}  VirtualWarehouseUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/{id}/edit [post]
func (h *handler) UpdateVirtualWarehouse(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("virtual_warehouse").(*omnichannel.User_Virtual_Warehouse_Registration)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateVirtualWarehouse(uint(id), *data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Virtual Warehouse is updated successfully", map[string]interface{}{"updated_id": id})
}

// FindVirtualWarehouse godoc
// @Summary      Find a VirtualWarehouse by id
// @Description  Find a VirtualWarehouse by id
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  VirtualWarehouseGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/{id} [get]
func (h *handler) FindVirtualWarehouseByID(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data, err := h.service.FindByidVirtualWarehouse(uint(id))
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Virtual Warehouse data Retrieved successfully", data)
}

// AllVirtualWarehouses godoc
// @Summary      get AllVirtualWarehouses and filter it by search query
// @Description  a get  AllVirtualWarehouses and filter it by search query
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Param		 filters	  query	string false "filters"
// @Param   	 per_page 	  query int    false "per_page"
// @Param  		 page_no      query int    false "page_no"
// @Param   	 sort		  query string false "sort"
// @Success      200  {object}  VirtualWarehouseGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse [get]
func (h *handler) FetchAllVirtualWarehouses(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	result, err := h.service.FindAllVirtualWarehouses(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "All Virtual Warehouse data Retrieved successfully", result, p)
}

// DeleteVirtualWarehouse godoc
// @Summary      delete a DeleteVirtualWarehouse
// @Description  delete a DeleteVirtualWarehouse
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  VirtualWarehouseLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/{id}/delete [delete]
func (h *handler) DeleteVirtualWarehouse(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	err = h.service.DeleteVirtualWarehouse(uint(id), uint(user_id))
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Virtual Warehouse deleted successfully", map[string]interface{}{"deleted_id": id})
}

// LINK VIRTUAL WAREHOUSE DETAILS
// CreateVirtualWarehouseDetail godoc
// @Summary      do a LinkVirtualWarehouseDetail
// @Description  Link a VirtualWarehouseDetail
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Param        VirtualWarehouse_Details body Virtual_Warehouse_Details true "link a VirtualWarehouse"
// @Success      200  {object}  VirtualWarehouseDetailCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/create [post]
func (h *handler) SaveVirtualWarehouseDetails(c echo.Context) (err error) {

	data := c.Get("virtual_warehouse_details").(*omnichannel.User_Virtual_Warehouse_Link)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.SaveVirtualWarehouseDetails(data)
	// @Summary      Find a VirtualWarehouse by id
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateVirtualWarehouseInCache(token_id, access_template_id)

	return res.RespSuccess(c, "Virtual Warehouse Details created successfully", map[string]interface{}{"created_id": data.ID})
}

// UpdateVirtualWarehouseDetail godoc
// @Summary      do a UpdateVirtualWarehouseDetail
// @Description  Update a VirtualWarehouseDetail
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Param        VirtualWarehouse_details body Virtual_Warehouse_Details true "Update a VirtualWarehouse"
// @Success      200  {object}  VirtualWarehouseDetailUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/{id}/update [post]
func (h *handler) UpdateVirtualWarehouseDetails(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("virtual_warehouse_details").(*omnichannel.User_Virtual_Warehouse_Link)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateVirtualWarehouseDetails(uint(id), *data)
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateVirtualWarehouseInCache(token_id, access_template_id)

	return res.RespSuccess(c, "Update Virtual Warehouse Details successfully done", map[string]interface{}{"updated_id": id})
}

// FindVirtualWarehouseDetails godoc
// @Summary      Find a VirtualWarehouseDetails by id
// @Description  Find a VirtualWarehouseDetails by id
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  VirtualWarehouseDetailGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/{id}/view [get]
func (h *handler) FindVirtualWarehouseDetails(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	result, err := h.service.FindVirtualWarehouseDetails(uint(id))
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Data Retrieved successdully", result)
}

// AllVirtualWarehouseDetails godoc
// @Summary      get all AllVirtualWarehouseDetails and filter it by search query
// @Description  a get all AllVirtualWarehouseDetails and filter it by search query
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Success      200  {object}  VirtualWarehouseDetailGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/list [get]
func (h *handler) ListVirtualWarehouseDetails(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)

	token_id := c.Get("TokenUserID").(string)

	var response interface{}
	if *p == pagination.BasePaginatevalue {
		response, *p = GetVirtualWarehouseFromCache(token_id)
	}

	if response != nil {
		return res.RespSuccessInfo(c, "data retrieved successfully", response, p)
	}

	result, err := h.service.ListVirtualWarehouseDetails(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Data Retrieved successfully", result, p)
}

// DeleteVirtualWarehouseDetails godoc
// @Summary      delete a DeleteVirtualWarehouseDetails
// @Description  delete a DeleteVirtualWarehouseDetails
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id   path      int  true  "id"
// @Success      200  {object}  VirtualWarehouseDetailDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/{id}/remove [delete]
func (h *handler) DeleteVirtualWarehouseDetails(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteVirtualWarehouseDetails(uint(id), uint(user_id))
	if err != nil {
		return res.RespErr(c, err)
	}

	// cache implementation
	UpdateVirtualWarehouseInCache(token_id, access_template_id)

	return res.RespSuccess(c, "Virtual Warehouse Details deleted successfully", map[string]int{"deleted_id": id})
}

// AvailableVirtualWarehouses godoc
// @Summary      get AvailableVirtualWarehouses and filter it by search query
// @Description  get AvailableVirtualWarehouses and filter it by search query
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @param        Authorization header string true "Authorization"
// @Success      200  {object}  AvailableVirtualWarehousesResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/available [get]
func (h *handler) AvailableVirtualWarehouses(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	result, err := h.service.AvailableVirtualWarehouses(p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Available Virtual Warehouse data Retrieved successfully", result, p)
}

// Get AuthKeys godoc
// @Summary      Get AuthKeys
// @Description  Get AuthKeys
// @Tags         VirtualWarehouse
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        virtual_warehouse_code   path      string  true  "virtual_warehouse_code"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/virtual_warehouse/{virtual_warehouse_code}/get_auth_keys [get]
func (h *handler) GetAuthKeys(c echo.Context) (err error) {

	virtual_warehouse := c.Param("virtual_warehouse_code")
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"virtual_warehouse_code": virtual_warehouse,
		"created_by":             user_id,
	}
	result, err := h.service.GetAuthKeys(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", result)
}

// FavouriteVirtualWarehouse godoc
// @Summary     FavouriteVirtualWarehouse
// @Description FavouriteVirtualWarehouse
// @Tags        VirtualWarehouse
// @Accept      json
// @Produce     json
// @param id path string true "VirtualWarehouse ID"
// @Security    ApiKeyAuth
// @param       Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/virtual_warehouse/{id}/favourite [post]
func (h *handler) FavouriteVirtualWarehouse(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.FindByidVirtualWarehouse(uint(ID))
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	fmt.Println("handler ------------------------ -1", query)
	err = h.base_service.FavouriteVirtualWarehouseFav(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Virtual Warehouse is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteVirtualWarehouse godoc
// @Summary     UnFavouriteVirtualWarehouse
// @Description UnFavouriteVirtualWarehouse
// @Tags        VirtualWarehouse
// @Accept      json
// @Produce     json
// @param id path string true "VirtualWarehouse ID"
// @Security    ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/virtual_warehouse/{id}/unfavourite [post]
func (h *handler) UnFavouriteVirtualWarehouse(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteVirtualWarehouseFav(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Virtual Warehouse is Unmarked as Favourite", map[string]string{"record_id": id})
}
