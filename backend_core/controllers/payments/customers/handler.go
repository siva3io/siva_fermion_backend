package customers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/payments"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/google/uuid"
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

var CustomersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if CustomersHandler != nil {
		return CustomersHandler
	}
	service := NewService()
	CustomersHandler = &handler{service}
	return CustomersHandler
}

func (h *handler) CreateCustomerEvent(c echo.Context) error {
	input_data := c.Get("customers").(*payments.Customers)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	input_data.CreatedByID = helpers.ConvertStringToUint(token_id)
	input_data.UserId = *input_data.CreatedByID
	request_payload := make(map[string]interface{}, 0)
	request_payload["data"] = input_data
	meta_data := make(map[string]interface{}, 0)
	meta_data["user_id"] = user_id
	//unique_id
	request_id := uuid.New().String()
	meta_data["request_id"] = request_id
	request_payload["meta_data"] = meta_data
	eda.Produce(eda.CREATE_CUSTOMER, request_payload)
	return res.RespSuccess(c, "Customer Creation Inprogress", map[string]interface{}{"request_id": request_id})

}

func (h *handler) CreateCustomer(request map[string]interface{}) {

	data := request["data"].(map[string]interface{})
	var final_data payments.Customers
	marshal_data, _ := json.Marshal(data)
	json.Unmarshal(marshal_data, &final_data)

	err := h.service.CreateCustomer(&final_data)
	response_message := new(eda.ConsumerResponse)
	if err != nil {
		response_message.ErrorMessage = err
		// eda.Produce(eda.CREATE_CUSTOMER_ACK,*response_message)
		return
	}
	response_message.Response = map[string]interface{}{
		"updated_id": final_data.ID,
	}
	// eda.Produce(eda.CREATE_CUSTOMER_ACK,*response_message)

}

func (h *handler) ListCustomers(c echo.Context) error {

	page := new(pagination.Paginatevalue)
	c.Bind(page)
	resp, err := h.service.ListCustomers(page)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccessInfo(c, "list of customers fetched", resp, page)
}

func (h *handler) ViewCustomer(c echo.Context) error {

	var query = make(map[string]interface{}, 0)

	ID := c.Param("id")
	fmt.Println("ID", ID)
	id, _ := strconv.Atoi(ID)
	query["id"] = id
	fmt.Println("query", query)
	resp, err := h.service.ViewCustomer(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "customer fetched", resp)
}

func (h *handler) UpdateCustomerEvent(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{}, 0)
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	data := c.Get("customers").(*payments.Customers)
	s := c.Get("TokenUserID").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(s)
	request_payload := make(map[string]interface{}, 0)
	request_payload["data"] = &data
	meta_data := make(map[string]interface{}, 0)
	meta_data["query"] = query
	//unique_id
	request_id := uuid.New().String()
	meta_data["request_id"] = request_id
	request_payload["meta_data"] = meta_data
	eda.Produce(eda.UPDATE_CUSTOMER, request_payload)
	return res.RespSuccess(c, "Customer Update Inprogress", map[string]interface{}{"request_id": request_id})

}

func (h *handler) UpdateCustomer(request map[string]interface{}) {
	data := request["data"].(map[string]interface{})
	var final_data payments.Customers
	marshal_data, _ := json.Marshal(data)
	json.Unmarshal(marshal_data, &final_data)
	meta_data := request["meta_data"].(map[string]interface{})
	query := meta_data["query"].(map[string]interface{})
	//--------------------------concurrency_management_end-----------------------------------------
	err := h.service.UpdateCustomer(query, &final_data)
	response_message := new(eda.ConsumerResponse)
	if err != nil {
		response_message.ErrorMessage = err
		// eda.Produce(eda.UPDATE_CUSTOMER_ACK,*response_message)
		return
	}
	response_message.Response = map[string]interface{}{
		"updated_id": final_data.ID,
	}
	// eda.Produce(eda.UPDATE_CUSTOMER_ACK,*response_message)

}
func (h *handler) DeleteCustomer(c echo.Context) (err error) {
	var id = c.Param("id")
	var query = make(map[string]interface{})
	ID, _ := strconv.Atoi(id)
	query["id"] = ID
	s := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(s)
	query["created_by"] = user_id
	err = h.service.DeleteCustomer(query)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "customer deleted successfully", map[string]string{"customer_id": id})
}

func (h *handler) GetCustomer(c echo.Context) error {
	query_params_map := c.QueryParams()
	query_params := map[string]interface{}{}
	for key, value := range query_params_map {
		query_params[key] = value[0]
	}
	// fmt.Println(query_params)
	resp, err := h.service.GetCustomer(query_params)
	if err != nil {
		return res.SuccessWithError(c, err.Error(), nil)
	}
	return res.RespSuccess(c, "customer fetched", resp)
}
