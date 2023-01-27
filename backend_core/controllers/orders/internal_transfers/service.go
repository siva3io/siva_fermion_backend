package internal_transfers

import (
	"encoding/json"
	"errors"
	"fmt"

	"fermion/backend_core/controllers/inventory_orders/asn"
	"fermion/backend_core/controllers/inventory_orders/grn"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	orders_repo "fermion/backend_core/internal/repository/orders"
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
	CreateInternalTransfers(metaData core.MetaData, data *orders.InternalTransfers) error
	ListInternalTransfers(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ViewInternalTransfers(metaData core.MetaData) (interface{}, error)
	UpdateInternalTransfers(metaData core.MetaData, data *orders.InternalTransfers) error
	DeleteInternalTransfers(metaData core.MetaData) error
	DeleteIstLines(metaData core.MetaData) error

	GetInternalTransfersTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	istRepository      orders_repo.InternalTransfers
	grnService         grn.Service
	asnService         asn.Service
	doRepository       orders_repo.DeliveryOrders
	purchaseRepositary orders_repo.Purchase
	salesRepositary    orders_repo.Sales
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		istRepository:      orders_repo.NewInternalTransfer(),
		grnService:         grn.NewService(),
		asnService:         asn.NewService(),
		doRepository:       orders_repo.NewDo(),
		purchaseRepositary: orders_repo.NewPurchaseOrder(),
		salesRepositary:    orders_repo.NewSalesOrder(),
	}
	return newServiceObj
}

func (s *service) CreateInternalTransfers(metaData core.MetaData, data *orders.InternalTransfers) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "CREATE", "IST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create ist order at data level")
	}

	data.StatusHistory, _ = helpers.UpdateStatusHistory(data.StatusHistory, *data.StatusId)
	data.IstNumber = helpers.GenerateSequence("IST", fmt.Sprint(metaData.TokenUserId), "internal_transfers")
	defaultStatus, _ := helpers.GetLookupcodeId("IST_STATUS", "DRAFT")
	data.StatusId = &defaultStatus
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	data = s.GetTotalQuantity(data)
	err := s.istRepository.Save(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ListInternalTransfers(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "IST", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list ist order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list ist order at data level")
	}

	data, err := s.istRepository.FindAll(metaData.Query, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) ViewInternalTransfers(metaData core.MetaData) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "IST", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view ist order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view ist order at data level")
	}

	data, err := s.istRepository.FindOne(metaData.Query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) UpdateInternalTransfers(metaData core.MetaData, data *orders.InternalTransfers) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "UPDATE", "IST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update ist order at data level")
	}

	foundData, err := s.istRepository.FindOne(metaData.Query)
	if err != nil {
		return err
	}

	if data.StatusId != nil {
		if *data.StatusId != *foundData.StatusId && *data.StatusId != 0 {
			data.StatusHistory, _ = helpers.UpdateStatusHistory(foundData.StatusHistory, *data.StatusId)
		}
	}
	if len(data.InternalTransferLines) > 0 {
		data = s.GetTotalQuantity(data)
	}
	data.UpdatedByID = &metaData.TokenUserId
	err = s.istRepository.Update(metaData.Query, data)
	if err != nil {
		return err
	}
	for _, orderLine := range data.InternalTransferLines {
		updateQuery := map[string]interface{}{
			"product_id": orderLine.ProductId,
			"ist_id":     metaData.Query["id"],
		}
		err = s.istRepository.UpdateOrderLines(updateQuery, &orderLine)
		if err != nil {
			orderLine.IstId = uint(metaData.Query["id"].(float64))
			err = s.istRepository.SaveOrderLines(&orderLine)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) DeleteInternalTransfers(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "IST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ist order at data level")
	}

	err := s.istRepository.Delete(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{
			"ist_id": metaData.Query["id"],
		}
		err = s.istRepository.DeleteOrderLine(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteIstLines(metaData core.MetaData) error {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "IST", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete ist order at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete ist order at data level")
	}
	err := s.istRepository.DeleteOrderLine(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetInternalTransfersTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "IST", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view ist order at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view ist order at data level")
	}

	id := metaData.AdditionalFields["id"].(string)
	tab := metaData.AdditionalFields["tab"].(string)

	if tab == "grn" {

		var query = make(map[string]interface{}, 0)
		query["id"] = id
		query["source_document_type_id"], _ = helpers.GetLookupcodeId("IST_SOURCE_DOCUMENT_TYPES", "GRN")

		internal_transfer_order, err := s.ViewInternalTransfers(metaData)
		if err != nil {
			return nil, errors.New("no source document")
		}
		var internalTransfersSourceDoc map[string]interface{}
		internalTransfersSourceDocJson := internal_transfer_order.(orders.InternalTransfers).SourceDocuments
		dto, err := json.Marshal(internalTransfersSourceDocJson)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(dto, &internalTransfersSourceDoc)
		if err != nil {
			return 0, err
		}

		internalTransfersSourceDocId := internalTransfersSourceDoc["id"]
		if internalTransfersSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		grnPage := page
		grnPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", internalTransfersSourceDoc["id"])
		data, err := s.grnService.GetAllGRN(metaData, grnPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "asn" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("ASN_SOURCE_DOCUMENT_TYPES", "IST")

		asnPage := page
		asnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, id)
		data, err := s.asnService.GetAllAsn(metaData, asnPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "delivery_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES", "IST")

		deliveryOrderPage := page
		deliveryOrderPage.Filters = fmt.Sprintf("[[\"source_document_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, id)
		data, err := s.doRepository.AllDeliveryOrders(nil, deliveryOrderPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "purchase_orders" {
		istPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("IST_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")
		istPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", id, source_document_type_id)
		ist_order_interface, err := s.ListInternalTransfers(metaData, page)
		if err != nil {
			return nil, err
		}
		ist_order := ist_order_interface.([]orders.InternalTransfers)
		if len(ist_order) == 0 {
			return nil, errors.New("no source document")
		}

		var istOrderSourceDoc map[string]interface{}
		istOrderSourceDocJson := ist_order[0]
		dto, err := json.Marshal(istOrderSourceDocJson)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(dto, &istOrderSourceDoc)
		if err != nil {
			return 0, err
		}

		istOrderSourceDocId := istOrderSourceDoc["id"]
		if istOrderSourceDocId == nil {
			return nil, errors.New("no source document")
		}

		purchaseOrdersPage := page
		purchaseOrdersPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", istOrderSourceDoc["id"])
		data, err := s.purchaseRepositary.FindAll(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "sales_orders" {
		istPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("IST_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")
		istPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", id, source_document_type_id)
		ist_order_interface, err := s.ListInternalTransfers(metaData, page)
		if err != nil {
			return nil, err
		}
		ist_order := ist_order_interface.([]orders.InternalTransfers)
		if len(ist_order) == 0 {
			return nil, errors.New("no source document")
		}

		var istOrderSourceDoc map[string]interface{}
		istOrderSourceDocJson := ist_order[0]
		dto, err := json.Marshal(istOrderSourceDocJson)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(dto, &istOrderSourceDoc)
		if err != nil {
			return 0, err
		}

		istOrderSourceDocId := istOrderSourceDoc["id"]
		if istOrderSourceDocId == nil {
			return nil, errors.New("no source document")
		}

		salesOrdersPage := page
		salesOrdersPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", istOrderSourceDoc["id"])
		data, err := s.salesRepositary.FindAll(nil, page)
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
