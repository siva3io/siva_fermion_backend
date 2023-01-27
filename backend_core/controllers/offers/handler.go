package offers

import (
	"encoding/json"
	"strconv"

	"fermion/backend_core/internal/model/offers"
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
	service Service
}

var OfferHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if OfferHandler != nil {
		return OfferHandler
	}
	service := NewService()
	OfferHandler = &handler{service: service}
	return OfferHandler
}

// Createoffers godoc
// @Summary Create Offers
// @Description Create Offers
// @Tags Offers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body OffersDTO true "Offers Request Body"
// @Success 200 {object} OffersCreateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/offers/create [post]
func (h *handler) CreateOffers(c echo.Context) error {

	var data offers.Offers
	err := c.Bind(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	token_id := c.Get("TokenUserID").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.CreateOffers(&data, token_id, "1")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Offers created successfully", map[string]interface{}{"created_id": data.ID})
}

// GetListOffers godoc
// @Summary Get all Offers list
// @Description Get all Offers list
// @Tags Offers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "filters"
// @Param   sort		query   string 	false "sort"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {array} GetAllOffersResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/offers [get]
func (h *handler) ListOffers(c echo.Context) (err error) {
	page := new(pagination.Paginatevalue)
	c.Bind(page)
	token_id := c.Get("TokenUserID").(string)

	response, err := h.service.ListOffers(page, token_id)
	if err != nil {
		return res.RespError(c, err)
	}
	var offerResponse []OffersDTO
	dt, err := json.Marshal(response)
	json.Unmarshal(dt, &offerResponse)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "success", offerResponse, page)
}

// ViewOffers godoc
// @Summary View Offers
// @Summary View Offers
// @Description View Offers
// @Tags Offers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "offer_id"
// @Success 200 {object} OffersGetResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/offers/{id} [get]
func (h *handler) ViewOffers(c echo.Context) error {

	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	query["id"] = int(id)
	response, err := h.service.ViewOffers(query, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	var offerResponse OffersDTO
	dt, err := json.Marshal(response)
	json.Unmarshal(dt, &offerResponse)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "success", offerResponse)
}

// UpdateOffers godoc
// @Summary Update Offers
// @Description Update Offers
// @Tags Offers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "offer_id"
// @param RequestBody body OffersDTO true "Offers Request Body"
// @Success 200 {object} OffersUpdateResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/offers/{id}/update [post]
func (h *handler) UpdateOffers(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = uint(ID)
	token_id := c.Get("TokenUserID").(string)
	company_id := c.Get("TokenUserID").(string)
	var data offers.Offers
	err = c.Bind(&data)
	if err != nil {
		return err
	}
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateOffers(query, &data, token_id, company_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Offers updated successfully", map[string]interface{}{"updated_id": id})
}

// DeleteOffers godoc
// @Summary Delete Offers
// @Description Delete Offers
// @Tags Offers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "offer_id"
// @Success 200 {object} OffersDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/offers/{id}/delete [delete]
func (h *handler) DeleteOffers(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeleteOffers(query, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Offer deleted successfully", map[string]string{"deleted_id": id})
}

// DeleteOfferLines godoc
// @Summary Delete OffersLines
// @Description Delete OfferLines
// @Tags Offers
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "offer_id"
// @param product_id query string true "product_id"
// @Success 200 {object} OfferLinesDeleteResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/offers/offer_lines/{id}/delete [delete]
func (h *handler) DeleteOffersLines(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	query = map[string]interface{}{
		"offer_id":   id,
		"product_id": c.QueryParam("product_id"),
	}
	err = h.service.DeleteOffersLines(query, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "offer line deleted successfully", map[string]string{"deleted_id": id})
}
