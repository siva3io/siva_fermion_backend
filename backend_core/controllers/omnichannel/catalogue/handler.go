package catalogue

import (
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/omnichannel"
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

var CatalogueHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if CatalogueHandler != nil {
		return CatalogueHandler
	}
	service := NewService()
	CatalogueHandler = &handler{service}
	return CatalogueHandler
}

//----------------------------------CatalogueTemplate--------------------------------------------------

// UpsertCatalogueTemplate godoc
// @Summary      do a UpsertCatalogueTemplate
// @Description  Upsert a CatalogueTemplate
// @Tags         Catalogue
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/catalogue_template/upsert [post]
func (h *handler) UpsertCatalogueTemplate(c echo.Context) (err error) {
	var data omnichannel.CatalogueTemplate
	c.Bind(&data)
	s := c.Get("TokenUserID").(string)
	token_int := helpers.ConvertStringToUint(s)
	data.CreatedByID = token_int
	data.UpdatedByID = token_int
	err = h.service.UpsertCatalogueTemplate(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "CatalogueTemplate created successfully", map[string]interface{}{"created_id": data.ID})
}

//-------------------------------Catalogue--------------------------------------------------------

// Getcatalogue godoc
// @Summary View Catalogue
// @Summary View Catalogue
// @Description View Catalogue
// @Tags Catalogue
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "External Category ID"
// @param id path string true "External MarketPlace ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/catalogue [get]
func (h *handler) GetCatalogue(c echo.Context) (err error) {
	ext_category_id, _ := strconv.Atoi(c.QueryParam("ext_category_id"))
	ext_channel_id, _ := strconv.Atoi(c.QueryParam("ext_channel_id"))
	variant_id, _ := strconv.Atoi(c.QueryParam("variant_id"))
	company_id, _ := strconv.Atoi(c.QueryParam("company_id"))
	query := map[string]interface{}{
		"category_id": ext_category_id,
		"channel_id":  ext_channel_id,
		"variant_id":  variant_id,
		"company_id":  company_id,
	}
	result, err := h.service.GetCatalogue(query)
	fmt.Println(result)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Catalogue Fetched Successfully", result)
}

//-------------------------------CatalogueTemplateData--------------------------------------------------

// UpsertCatalogueTemplateData godoc
// @Summary      do a UpsertCatalogueTemplateData
// @Description  Upsert a CatalogueTemplateData
// @Tags         Catalogue
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200  {object}  res.SuccessResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/catalogue_template_data/upsert [post]
func (h *handler) UpsertCatalogueTemplateData(c echo.Context) (err error) {
	var data omnichannel.CatalogueTemplateData
	c.Bind(&data)
	s := c.Get("TokenUserID").(string)
	token_int := helpers.ConvertStringToUint(s)
	data.CreatedByID = token_int
	data.UpdatedByID = token_int
	err = h.service.UpsertCatalogueTemplateData(&data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "CatalogueTemplateData created/updated successfully", map[string]interface{}{"created_id": data.ID})
}
