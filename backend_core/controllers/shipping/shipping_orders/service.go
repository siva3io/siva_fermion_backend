package shipping_orders

import (
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/controllers/shipping/shipping_orders_ndr"
	"fermion/backend_core/controllers/shipping/shipping_orders_rto"
	"fermion/backend_core/controllers/shipping/shipping_orders_wd"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/shipping"
	shipping_repo "fermion/backend_core/internal/repository/shipping"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
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
	CreateShippingOrder(metaData core.MetaData, data *shipping.ShippingOrder) error
	GetAllShippingOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	GetShippingOrder(metaData core.MetaData) (interface{}, error)
	DeleteShippingOrder(metaData core.MetaData) error
	UpdateShippingOrder(metaData core.MetaData, data *shipping.ShippingOrder) error
	DeleteShippingOrderLine(metaData core.MetaData) error
	CreateBulkShippingOrder(metaData core.MetaData, data *[]shipping.ShippingOrder) error
	GetShippingOrderTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error)
}

type service struct {
	shippingrepository shipping_repo.ShippingOrder
	rtoService         shipping_orders_rto.Service
	ndrService         shipping_orders_ndr.Service
	wdService          shipping_orders_wd.Service
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	shippingrepository := shipping_repo.NewShippingOrder()
	newServiceObj = &service{shippingrepository,
		shipping_orders_rto.NewService(),
		shipping_orders_ndr.NewService(),
		shipping_orders_wd.NewService()}
	return newServiceObj
}

func (s *service) CreateShippingOrder(metaData core.MetaData, data *shipping.ShippingOrder) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "SHIPPING_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create uom at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	data.ShippingNumber = helpers.GenerateSequence("SHIP", fmt.Sprint(*data.CreatedByID), "shipping_orders")
	data.ReferenceNumber = helpers.GenerateSequence("REF", fmt.Sprint(*data.CreatedByID), "shipping_orders")
	data.AwbNumber = helpers.GenerateSequence("AWB", fmt.Sprint(*data.CreatedByID), "shipping_orders")

	err := s.shippingrepository.Create(data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) GetAllShippingOrders(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "SHIPPING_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list shipping order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list shipping order at data level")
	}
	result, _ := s.shippingrepository.List(metaData.Query, page)
	return result, nil
}

func (s *service) GetShippingOrder(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "SHIPPING_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return ShippingOrderResponse{}, fmt.Errorf("you dont have access for view shipping order at view level")
	}
	if data_access == nil {
		return ShippingOrderResponse{}, fmt.Errorf("you dont have access for view shipping order at data level")
	}
	result, err := s.shippingrepository.Get(metaData.Query)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteShippingOrder(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "SHIPPING_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete shipping order at data level")
	}
	err := s.shippingrepository.Delete(metaData.Query)
	if err != nil {
		return err
	}
	err1 := s.shippingrepository.DeleteOrderLine(metaData.Query)
	if err1 != nil {
		return err
	}
	return nil
}

func (s *service) DeleteShippingOrderLine(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "SHIPPING_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete shipping order at data level")
	}
	err := s.shippingrepository.DeleteOrderLine(metaData.Query)
	return err
}

func (s *service) UpdateShippingOrder(metaData core.MetaData, data *shipping.ShippingOrder) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "SHIPPING_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update shipping order at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	// var shipping_data shipping.ShippingOrder
	// dto, err := json.Marshal(*data)
	// if err != nil {
	// 	return res.BuildError(res.ErrUnprocessableEntity, err)
	// }
	// err = json.Unmarshal(dto, &shipping_data)
	// if err != nil {
	// 	return res.BuildError(res.ErrUnprocessableEntity, err)
	// }
	// old_data, er := s.shippingrepository.Get(metaData.Query)
	// if er != nil {
	// 	return er
	// }
	// old_status := old_data.ShippingStatusId
	// new_status := shipping_data.ShippingStatusId
	// if new_status != old_status && new_status != 0 {
	// 	result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.ShippingStatusId)
	// 	shipping_data.StatusHistory = result
	// }

	// if data.BillingAddress == (ReceiverOrBillingAddresss{}) {
	// 	shipping_data.BillingAddress = nil
	// }

	// if data.SenderAddress == (SenderAddress{}) {
	// 	shipping_data.SenderAddress = nil
	// }

	// if data.ReceiverAddress == (ReceiverOrBillingAddresss{}) {
	// 	shipping_data.ReceiverAddress = nil
	// }

	// if data.PackageDetails == (PackageDetails{}) {
	// 	shipping_data.PackageDetails = nil
	// }

	// if data.Billingdetails == (BillingDetails{}) {
	// 	shipping_data.Billingdetails = nil
	// }

	// if data.ShippingManifestId == nil {
	// 	shipping_data.ShippingManifestId = nil
	// }

	// if data.ShippingInvoiceId == nil {
	// 	shipping_data.ShippingInvoiceId = nil
	// }

	status_id := *data.ShippingStatusId
	service_err := s.shippingrepository.Update(metaData.Query, data)
	if service_err != nil {
		return service_err
	}
	name, err := helpers.GetLookupcodeName(status_id)
	if err != nil {
		fmt.Println("ERROR---------", err)
	}

	salesOrderStatusMap := map[string]interface{}{
		"SHIPPED":   "SHIPPED",
		"CANCELLED": "CANCELLED",
		"DELIVERED": "DELIVERED",
	}

	if salesOrderStatusMap[name] != nil {
		status_id, err = helpers.GetLookupcodeId("SALES_ORDER_STATUS", salesOrderStatusMap[name].(string))
		if err != nil {
			fmt.Println("ERROR---------", err)
		}

		shippingOrderResponse, err := s.shippingrepository.Get(metaData.Query)
		if err != nil {
			fmt.Println("ERROR---------", err)
		}

		salesOrderId := *shippingOrderResponse.OrderId
		metaData.Query = map[string]interface{}{
			"id": salesOrderId,
		}
		status_data := map[string]interface{}{
			"status_id": status_id,
		}

		request_payload := map[string]interface{}{
			"meta_data": metaData,
			"data":      status_data,
		}

		eda.Produce(eda.UPDATE_SALES_ORDER, request_payload)

	}
	// for _, product := range shipping_data.ShippingOrderLines {
	// 	query := map[string]interface{}{
	// 		"shipping_order_id": id,
	// 		"product_id":        product.ProductId,
	// 	}
	// 	count, err1 := s.shippingrepository.UpdateProductList(query, product)
	// 	if err1 != nil {
	// 		return err1
	// 	} else if count == 0 {
	// 		product.ShippingOrderId = id
	// 		err2 := s.shippingrepository.AddProduct(product)
	// 		if err2 != nil {
	// 			return err2
	// 		}
	// 	}
	// }
	return nil
}

func (s *service) CreateBulkShippingOrder(metaData core.MetaData, data *[]shipping.ShippingOrder) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "SHIPPING_ORDERS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create shipping order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create shipping order at data level")
	}
	service_err := s.shippingrepository.CreateBulkOrder(data)
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
		var metaData core.MetaData
		data, err := s.rtoService.GetAllRTO(metaData, rtoPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "ndr" {
		var metaData core.MetaData
		p := new(pagination.Paginatevalue)
		ndrPage := page
		ndrPage.Filters = fmt.Sprintf("[[\"shipping_order_id\",\"=\",%v]]", shippingOrderId)
		data, err := s.ndrService.GetAllNDR(metaData, p)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "wd" {
		var metaData core.MetaData
		p := new(pagination.Paginatevalue)
		wdPage := page
		wdPage.Filters = fmt.Sprintf("[[\"shipping_order_id\",\"=\",%v]]", shippingOrderId)
		data, err := s.wdService.GetAllWD(metaData, p)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
