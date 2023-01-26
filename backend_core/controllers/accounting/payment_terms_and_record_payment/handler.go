package payment_terms_and_record_payment

import (
	"fmt"
	"strconv"

	accounting_base "fermion/backend_core/controllers/accounting/base"
	"fermion/backend_core/internal/model/accounting"
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
	base_service accounting_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := accounting_base.NewServiceBase()
	return &handler{service, base_service}
}

// CreatePaymentTerm godoc
// @Summary CreatePaymentTerm
// @Description CreatePaymentTerm
// @Tags PaymentTerms and Record Payments
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body PaymentTermsDTO true "Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/payment_terms/create [post]
func (h *handler) CreatePaymentTerm(c echo.Context) (err error) {

	var data accounting.PaymentTerms
	c.Bind(&data)
	//data := c.Get("payment_terms_and_record_payment").(*accounting.PaymentTerms)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	//seq := helpers.GenerateSequence("Acc", s, "payment_terms")
	//fmt.Println("The tokenuserid is ", seq)
	//data.PaymentTermName = seq
	err = h.service.CreatePaymentTerm(&data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Payment term created successfully", map[string]interface{}{"created_id": data.ID})
}

// UpdatePaymentTerm godoc
// @Summary UpdatePaymentTerm
// @Description UpdatePaymentTerm
// @Tags PaymentTerms and Record Payments
// @Accept  json
// @Produce  json
// @param id path string true "PaymentTerm ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param RequestBody body PaymentTermsDTO true "PaymentTerms  Request Body"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/payment_terms/{id}/update [post]
func (h *handler) UpdatePaymentTerm(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = uint(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	var data accounting.PaymentTerms
	err = c.Bind(&data)
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	err = h.service.UpdatePaymentTerm(query, &data, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Payment Term updated successfully", map[string]interface{}{"updated_id": data.ID})
}

// GetPaymentTerm godoc
// @Summary View PaymentTerms
// @Summary View PaymentTerms
// @Description View PaymentTerms
// @Tags PaymentTerms and Record Payments
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "PaymentTerms ID"
// @Success 200 {object} PaymentTermsDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/payment_terms/{id} [get]
func (h *handler) GetPaymentTerm(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	//println("--------->", ID)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetPaymentTerm(query, token_id, access_template_id)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "Payment Term Details Retrieved successfully", result)
}

// GetAllPaymentTerm godoc
// @Summary Get all Payment Terms
// @Description Get all Payment Terms
// @Tags  PaymentTerms and Record Payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} PaymentTermsDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/payment_terms [get]
func (h *handler) GetAllPaymentTerm(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetPaymentTermlist(query, p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Payment Term list Retrieved Successfully", result, p)
}

// GetAllPaymentTermDropDown godoc
// @Summary Get all Payment Terms dropdown
// @Description Get all Payment Terms dropdown
// @Tags  PaymentTerms and Record Payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param   filters 	query 	string 	false "[[column_name,operator,value]]"
// @Param   sort		query   string 	false "[[column_name,asc or desc]]"
// @Param   per_page 	query   int    	false "per_page"
// @Param   page_no     query   int    	false "page_no"
// @Success 200 {object} PaymentTermsDTO
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/payment_terms/dropdown [get]
func (h *handler) GetAllPaymentTermDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	var query = make(map[string]interface{}, 0)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetPaymentTermlist(query, p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "Payment Term dropdown list Retrieved Successfully", result, p)
}

// DeletePaymentTerm godoc
// @Summary DeletePaymentTerm
// @Description DeletePaymentTerm
// @Tags PaymentTerms and Record Payments
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "PaymentTerms ID"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/payment_terms/{id}/delete [delete]
func (h *handler) DeletePaymentTerm(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query["user_id"] = user_id
	err = h.service.DeletePaymentTerm(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Payment Term Record deleted successfully", map[string]string{"deleted_id": id})
}
