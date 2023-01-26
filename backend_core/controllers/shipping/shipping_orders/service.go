package shipping_orders

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/shipping/shipping_orders_ndr"
	"fermion/backend_core/controllers/shipping/shipping_orders_rto"
	"fermion/backend_core/controllers/shipping/shipping_orders_wd"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
	shipping_repo "fermion/backend_core/internal/repository/shipping"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"
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
type Service interface {
	CreateShippingOrder(data *ShippingOrderRequest, token_id string, access_template_id string) (*shipping.ShippingOrder, error)
	GetAllShippingOrders(page *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ShippingOrderResponse, error)
	GetShippingOrder(id uint, token_id string, access_template_id string) (ShippingOrderResponse, error)
	DeleteShippingOrder(id uint, user_id uint, token_id string, access_template_id string) error
	UpdateShippingOrder(id uint, data *ShippingOrderRequest, token_id string, access_template_id string) error
	DeleteShippingOrderLine(query interface{}, token_id string, access_template_id string) error
	CreateBulkShippingOrder(data *[]shipping.ShippingOrder, token_id string, access_template_id string) error
	GetShippingOrderTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error)
}

type service struct {
	shippingrepository shipping_repo.ShippingOrder
	rtoService         shipping_orders_rto.Service
	ndrService         shipping_orders_ndr.Service
	wdService          shipping_orders_wd.Service
}

func NewService() *service {
	shippingrepository := shipping_repo.NewShippingOrder()
	return &service{shippingrepository,
		shipping_orders_rto.NewService(),
		shipping_orders_ndr.NewService(),
		shipping_orders_wd.NewService()}
}

func (s *service) CreateShippingOrder(data *ShippingOrderRequest, token_id string, access_template_id string) (*shipping.ShippingOrder, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for create shipping order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for create shipping order at data level")
	}
	var shipping_data shipping.ShippingOrder
	dto, err := json.Marshal(*data)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &shipping_data)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	shipping_data.ShippingNumber = helpers.GenerateSequence("SHIP", fmt.Sprint(*shipping_data.CreatedByID), "shipping_orders")
	shipping_data.ReferenceNumber = helpers.GenerateSequence("REF", fmt.Sprint(*shipping_data.CreatedByID), "shipping_orders")
	response, service_err := s.shippingrepository.Create(&shipping_data)
	return response, service_err
}

func (s *service) GetAllShippingOrders(page *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ShippingOrderResponse, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list shipping order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list shipping order at data level")
	}
	result, _ := s.shippingrepository.List(page)
	var response []ShippingOrderResponse
	for _, data := range result {
		var deliverydata ShippingOrderResponse
		value, _ := json.Marshal(data)
		json.Unmarshal(value, &deliverydata)
		response = append(response, deliverydata)
	}
	return response, nil
}

func (s *service) GetShippingOrder(id uint, token_id string, access_template_id string) (ShippingOrderResponse, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return ShippingOrderResponse{}, fmt.Errorf("you dont have access for view shipping order at view level")
	}
	if data_access == nil {
		return ShippingOrderResponse{}, fmt.Errorf("you dont have access for view shipping order at data level")
	}
	result, err := s.shippingrepository.Get(id)
	var responses ShippingOrderResponse
	value, _ := json.Marshal(result)
	json.Unmarshal(value, &responses)
	if err != nil {
		return responses, err
	}
	return responses, nil
}

func (s *service) DeleteShippingOrder(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete shipping order at data level")
	}
	err := s.shippingrepository.Delete(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"shipping_order_id": id}
	err1 := s.shippingrepository.DeleteOrderLine(query)
	if err1 != nil {
		return err
	}
	return nil
}

func (s *service) DeleteShippingOrderLine(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete shipping order at data level")
	}
	err := s.shippingrepository.DeleteOrderLine(query)
	return err
}

func (s *service) UpdateShippingOrder(id uint, data *ShippingOrderRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update shipping order at data level")
	}
	var shipping_data shipping.ShippingOrder
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &shipping_data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	old_data, er := s.shippingrepository.Get(id)
	if er != nil {
		return er
	}
	old_status := old_data.ShippingStatusId
	new_status := shipping_data.ShippingStatusId
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.ShippingStatusId)
		shipping_data.StatusHistory = result
	}

	if data.BillingAddress == (ReceiverOrBillingAddresss{}) {
		shipping_data.BillingAddress = nil
	}

	if data.SenderAddress == (SenderAddress{}) {
		shipping_data.SenderAddress = nil
	}

	if data.ReceiverAddress == (ReceiverOrBillingAddresss{}) {
		shipping_data.ReceiverAddress = nil
	}

	if data.PackageDetails == (PackageDetails{}) {
		shipping_data.PackageDetails = nil
	}

	if data.Billingdetails == (BillingDetails{}) {
		shipping_data.Billingdetails = nil
	}

	if data.ShippingManifestId == nil {
		shipping_data.ShippingManifestId = nil
	}

	if data.ShippingInvoiceId == nil {
		shipping_data.ShippingInvoiceId = nil
	}

	service_err := s.shippingrepository.Update(id, &shipping_data)
	if service_err != nil {
		return service_err
	}
	for _, product := range shipping_data.ShippingOrderLines {
		query := map[string]interface{}{
			"shipping_order_id": id,
			"product_id":        product.ProductId,
		}
		count, err1 := s.shippingrepository.UpdateProductList(query, product)
		if err1 != nil {
			return err1
		} else if count == 0 {
			product.ShippingOrderId = id
			err2 := s.shippingrepository.AddProduct(product)
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func (s *service) CreateBulkShippingOrder(data *[]shipping.ShippingOrder, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create shipping order at data level")
	}
	service_err := s.shippingrepository.CreateBulkOrder(data, token_id)
	return service_err
}
func (s *service) GetShippingOrderTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "SHIPPING_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list shipping order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list shipping order at data level")
	}

	shippingOrderId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "rto" {

		rtoPage := page
		rtoPage.Filters = fmt.Sprintf("[[\"shipping_order_id\",\"=\",%v]]", shippingOrderId)
		data, err := s.rtoService.GetAllRTO(rtoPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "ndr" {

		ndrPage := page
		ndrPage.Filters = fmt.Sprintf("[[\"shipping_order_id\",\"=\",%v]]", shippingOrderId)
		data, err := s.ndrService.GetAllNDR(ndrPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "wd" {

		wdPage := page
		wdPage.Filters = fmt.Sprintf("[[\"shipping_order_id\",\"=\",%v]]", shippingOrderId)
		data, err := s.wdService.GetAllWD(wdPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
