package vendors

import (
	"errors"
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

//---------------------------------------------Vendors-----------------------------------------

// CreateVendors godoc
// @Summary CreateVendors
// @Description CreateVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorsRequestDTO true "Vendors Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/create [post]
func (h *handler) CreateVendors(c echo.Context) error {
	data := c.Get("vendors").(*mdm.Vendors)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err := h.service.CreateVendors(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor created successfully", map[string]interface{}{"created_id": data.ID, "Details": data})
}

// UpdateVendors godoc
// @Summary UpdateVendors
// @Description UpdateVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorsRequestDTO true "Vendors Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/update [post]
func (h *handler) UpdateVendors(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("vendors").(*mdm.Vendors)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateVendors(query, data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Details updated succesfully", map[string]interface{}{"updated_id": query["id"], "Details": data})
}

// DeleteVendors godoc
// @Summary      delete a Vendors
// @Description  delete a Vendors
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/vendors/{id}/delete [delete]
func (h *handler) DeleteVendors(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteVendors(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}

// GetVendors godoc
// @Summary Get Details of specific vendor based on ID
// @Description Get Details of specific vendor based on ID
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success 200 {object} VendorDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id} [get]
func (h *handler) GetVendors(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetVendors(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor data retrieved successfully", map[string]interface{}{"Vendor_ID": query["id"], "Details": result})
}

// GetVendorList godoc
// @Summary Get all vendors
// @Description Get all vendors
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors [get]
func (h *handler) GetVendorsList(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetVendorsList(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "List Retrieved successfully", result, p)
}

// GetVendorListDropdown godoc
// @Summary Get all vendors
// @Description Get all vendors
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/dropdown [get]
func (h *handler) GetVendorsListDropdown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetVendorsList(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "List Retrieved successfully", result, p)
}

// Search Vendors godoc
// @Summary Search Vendors
// @Description Search Vendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} VendorsObjDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/search [get]
func (h *handler) SearchVendors(c echo.Context) (err error) {
	q := c.QueryParam("q")
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.SearchVendors(q, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "OK", result)
}

// FavouriteVendors godoc
// @Summary FavouriteVendors
// @Description FavouriteVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "Vendor ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/favourite [post]
func (h *handler) FavouriteVendors(c echo.Context) (err error) {
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
	_, er := h.service.GetVendors(q, token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteVendors(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor is Marked as Favourite", map[string]string{"record_id": id})
}

// UnFavouriteVendors godoc
// @Summary UnFavouriteVendors
// @Description UnFavouriteVendors
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "Vendor ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/unfavourite [post]
func (h *handler) UnFavouriteVendors(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnfavouriteVendors(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor is Unmarked as Favourite", map[string]string{"record_id": id})
}

// GetFavouriteVendors godoc
// @Summary Get all favourite vendors
// @Description Get all favourite vendors
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/favourite_list [get]
func (h *handler) FavouriteVendorsView(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	data, er := h.base_service.GetArrVendor(user_id)
	if er != nil {
		return errors.New("favourite list not found for the specified user")
	}
	result, err := h.base_service.GetFavVendorList(data, p)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Favourite Vendors list Retrieved Successfully", result, p)
}

//---------------------------------------------vendor pricelist------------------------------------------

// CreateVendorPriceList godoc
// @Summary CreateVendorPriceList
// @Description CreateVendorPriceList
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorPriceListsDTO true "VendorPricelist Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist/create [post]
func (h *handler) CreateVendorPriceList(c echo.Context) error {
	data := c.Get("vendors_price_list").(*mdm.VendorPriceLists)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err := h.service.CreateVendorPriceList(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor price list created successfully", map[string]interface{}{"created_id": data.ID, "Details": data})
}

// UpdateVendorPricelist godoc
// @Summary UpdateVendorPricelist
// @Description UpdateVendorPricelist
// @Tags Vendors
// @Accept  json
// @Produce  json
// @param id path string true "ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body VendorPriceListsDTO true "VendorPricelist Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/{id}/update [post]
func (h *handler) UpdateVendorPriceLists(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("vendors_price_list").(*mdm.VendorPriceLists)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateVendorPriceLists(query, data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Details updated succesfully", map[string]interface{}{"updated_id": query["id"], "Details": data})
}

// GetVendorPiceList godoc
// @Summary Get Details of specific vendorpricelist based on ID
// @Description Get Details of specific vendorpricelist based on ID
// @Tags Vendors
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success 200 {object} VendorPriceListsDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist/{id} [get]
func (h *handler) GetVendorPriceList(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	query["id"] = ID
	result, err := h.service.GetVendorPriceList(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Vendor Pricelist data retrieved successfully", map[string]interface{}{"Vendor_Pricelist_ID": query["id"], "Details": result})
}

// GetVendorPriceLists godoc
// @Summary Get all vendors price list
// @Description Get all vendors price list
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist [get]
func (h *handler) GetVendorPriceLists(c echo.Context) (err error) {
	//q := c.QueryParam("vendor_id")
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	//id,_ := strconv.Atoi(q)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetVendorPriceLists(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "List Retrieved successfully", result, p)
}

// GetVendorPriceListsDropdown godoc
// @Summary Get all vendors price list
// @Description Get all vendors price list
// @Tags Vendors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} VendorsListDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/vendors/pricelist/dropdown [get]
func (h *handler) GetVendorPriceListsDropdown(c echo.Context) (err error) {
	//q := c.QueryParam("vendor_id")
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	//id,_ := strconv.Atoi(q)
	// id := c.Param("id")
	// ID, _ := strconv.Atoi(id)
	// //println("--------->", ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var query = make(map[string]interface{}, 0)
	// query["id"] = ID
	result, err := h.service.GetVendorPriceLists(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "List Retrieved successfully", result, p)
}

// DeleteVendorPricelist godoc
// @Summary      delete a VendorPricelist
// @Description  delete a VendorPricelist
// @Tags         Vendors
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        id             path      int    true  "id"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/vendors/pricelist/{id}/delete [delete]
func (h *handler) DeleteVendorPricelist(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteVendorPricelist(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", map[string]string{"deleted_id": id})
}
