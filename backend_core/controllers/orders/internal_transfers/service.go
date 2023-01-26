package internal_transfers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/inventory_orders/asn"
	"fermion/backend_core/controllers/inventory_orders/grn"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	orders_repo "fermion/backend_core/internal/repository/orders"
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
	CreateInternalTransfers(data *orders.InternalTransfers, access_template_id string, token_id string) error
	ListInternalTransfers(page *pagination.Paginatevalue, access_template_id string, token_id string, access_access string) (interface{}, error)
	ViewInternalTransfers(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error)
	UpdateInternalTransfers(query map[string]interface{}, data *orders.InternalTransfers, access_template_id string, token_id string) error

	DeleteInternalTransfers(query map[string]interface{}, access_template_id string, token_id string) error
	DeleteIstLines(query map[string]interface{}, access_template_id string, token_id string) error
	SearchInternalTransfers(query string, access_template_id string, token_id string) (interface{}, error)
	GetInternalTransfersTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	istRepository orders_repo.InternalTransfers
	grnService    grn.Service
	asnService    asn.Service
	doRepository  orders_repo.DeliveryOrders
}

func NewService() *service {
	return &service{
		istRepository: orders_repo.NewInternalTransfer(),
		grnService:    grn.NewService(),
		asnService:    asn.NewService(),
		doRepository:  orders_repo.NewDo(),
	}
}

func (s *service) CreateInternalTransfers(data *orders.InternalTransfers, access_template_id string, token_id string) error {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "IST", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create ist order at data level")
	}

	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	temp := data.CreatedByID
	user_id := strconv.Itoa(int(*temp))
	data.IstNumber = helpers.GenerateSequence("IST", user_id, "internal_transfers")
	defaultStatus, err := helpers.GetLookupcodeId("IST_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	data = s.GetTotalQuantity(data)
	err = s.istRepository.Save(data)

	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return nil
	//}
	//return nil
}

func (s *service) ListInternalTransfers(page *pagination.Paginatevalue, access_template_id string, token_id string, access_access string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_access, "IST", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list ist order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list ist order at data level")
	}

	data, err := s.istRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
	//}
	//return nil, nil
}

func (s *service) ViewInternalTransfers(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "IST", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view ist order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view ist order at data level")
	}
	data, err := s.istRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
	//}
	//return nil, nil
}

func (s *service) UpdateInternalTransfers(query map[string]interface{}, data *orders.InternalTransfers, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "IST", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update ist order at data level")
	}

	var find_query = map[string]interface{}{
		"id": query["id"],
	}

	found_data, er := s.istRepository.FindOne(find_query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	old_data := found_data.(orders.InternalTransfers)

	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result
	}
	if len(data.InternalTransferLines) > 0 {
		data = s.GetTotalQuantity(data)
	}
	err := s.istRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, order_line := range data.InternalTransferLines {
		update_query := map[string]interface{}{
			"product_id": order_line.ProductId,
			"ist_id":     uint(query["id"].(int)),
		}
		count, err1 := s.istRepository.UpdateOrderLines(update_query, order_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			order_line.IstId = uint(query["id"].(int))
			e := s.istRepository.SaveOrderLines(order_line)
			fmt.Println(e)
			if e != nil {
				return e
			}
		}
	}
	return nil
	//}
	//return nil
}

func (s *service) DeleteInternalTransfers(query map[string]interface{}, access_template_id string, token_id string) error {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "IST", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ist order at data level")
	}

	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.istRepository.FindOne(q)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.istRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	} else {
		query := map[string]interface{}{"ist_id": query["id"]}
		er := s.istRepository.DeleteOrderLine(query)
		if er != nil {
			return res.BuildError(res.ErrDataNotFound, er)
		}
	}
	return nil
	//}
	//return nil
}

func (s *service) DeleteIstLines(query map[string]interface{}, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "IST", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ist order at data level")
	}
	_, er := s.istRepository.FindOrderLines(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.istRepository.DeleteOrderLine(query)

	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
	//}
	//return nil
}

func (s *service) SearchInternalTransfers(query string, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "IST", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list ist order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list ist order at data level")
	}

	data, err := s.istRepository.Search(query)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
	//}
	//return nil, nil
}

func (s *service) GetInternalTransfersTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "IST", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view ist order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view ist order at data level")
	}

	istId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if tab == "grn" {

		var query = make(map[string]interface{}, 0)
		query["id"] = (uint(istId))
		query["source_document_type_id"], _ = helpers.GetLookupcodeId("IST_SOURCE_DOCUMENT_TYPES", "GRN")

		internal_transfer_order, err := s.ViewInternalTransfers(query, access_template_id, token_id)
		if err != nil {
			return nil, errors.New("no source document")
		}
		var internalTransfersSourceDoc map[string]interface{}
		internalTransfersSourceDocJson := internal_transfer_order.(orders.InternalTransfers).SourceDocuments
		dto, err := json.Marshal(internalTransfersSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &internalTransfersSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		internalTransfersSourceDocId := internalTransfersSourceDoc["id"]
		if internalTransfersSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		grnPage := page
		grnPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", internalTransfersSourceDoc["id"])
		data, err := s.grnService.GetAllGRN(grnPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "asn" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("ASN_SOURCE_DOCUMENT_TYPES", "IST")

		asnPage := page
		asnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, istId)
		data, err := s.asnService.GetAllAsn(asnPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "delivery_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "IST")

		deliveryOrderPage := page
		deliveryOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, istId)
		data, err := s.doRepository.AllDeliveryOrders(deliveryOrderPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	return nil, nil
}

func (s *service) GetTotalQuantity(data *orders.InternalTransfers) *orders.InternalTransfers {

	data.NoOfItems = len(data.InternalTransferLines)
	data.TotalQuantity = 0

	for _, orderLine := range data.InternalTransferLines {
		data.TotalQuantity += orderLine.TransferQuantity
	}
	return data
}
