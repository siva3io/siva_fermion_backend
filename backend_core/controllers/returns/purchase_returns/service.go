package purchase_returns

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/debitnote"
	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
	"fermion/backend_core/internal/repository/orders"
	returns_repo "fermion/backend_core/internal/repository/returns"
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
	CreatePurchaseReturn(metaData core.MetaData, data *PurchaseReturnsDTO) error
	ListPurchaseReturns(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ViewPurchaseReturn(metaData core.MetaData) (interface{}, error)
	UpdatePurchaseReturn(metaData core.MetaData, data *PurchaseReturnsDTO) error
	DeletePurchaseReturn(metaData core.MetaData) error
	DeletePurchaseReturnLines(metaData core.MetaData) error
	SearchPurchaseReturns(query string, access_template_id string, token_id string) (interface{}, error)
	GetPurchaseReturnsHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	GetPurchaseReturnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ProcessPurchaseReturnCalculation(data *PurchaseReturnsDTO) *PurchaseReturnsDTO
}

type service struct {
	purchaseReturnRepository returns_repo.PurchaseReturn
	basicInventoryService    inv_service.Service
	debitNoteService         debitnote.Service
	purchaseOrderRepositary  orders.Purchase
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	purchaseReturnRepository := returns_repo.NewPurchaseReturn()
	basicInventoryService := inv_service.NewService()
	purchaseOrderRepositary := orders.NewPurchaseOrder()
	newServiceObj = &service{purchaseReturnRepository, basicInventoryService, debitnote.NewService(), purchaseOrderRepositary}
	return newServiceObj
}

func (s *service) CreatePurchaseReturn(metaData core.MetaData, data *PurchaseReturnsDTO) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "CREATE", "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create purchase return at data level")
	}

	data = s.ProcessPurchaseReturnCalculation(data)
	var prData *returns.PurchaseReturns
	helpers.JsonMarshaller(data, &prData)
	prData.StatusHistory, _ = helpers.UpdateStatusHistory(prData.StatusHistory, prData.StatusId)
	prData.PurchaseReturnNumber = helpers.GenerateSequence("PR", fmt.Sprint(metaData.TokenUserId), "purchase_returns")
	defaultStatus, _ := helpers.GetLookupcodeId("RETURN_STATUS", "DRAFT")
	prData.StatusId = defaultStatus
	prData.CompanyId = metaData.CompanyId
	prData.CreatedByID = &metaData.TokenUserId

	err := s.purchaseReturnRepository.Save(prData)
	if err != nil {
		return err
	}

	dataInv := map[string]interface{}{
		"id":            prData.ID,
		"order_lines":   prData.PurchaseReturnLines,
		"order_type":    "purchase_return",
		"is_update_inv": true,
		"is_credit":     false,
	}
	//updating the inventory Committed Stock
	go s.basicInventoryService.UpdateTransactionInventory(metaData, dataInv)
	return nil
}

func (s *service) ListPurchaseReturns(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase return at data level")
	}

	data, err := s.purchaseReturnRepository.FindAll(metaData.Query, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) ViewPurchaseReturn(metaData core.MetaData) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase return at data level")
	}
	data, err := s.purchaseReturnRepository.FindOne(metaData.Query)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *service) UpdatePurchaseReturn(metaData core.MetaData, data *PurchaseReturnsDTO) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "UPDATE", "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update purchase return at data level")
	}

	data = s.ProcessPurchaseReturnCalculation(data)
	var PrData returns.PurchaseReturns
	helpers.JsonMarshaller(data, &PrData)
	old_data, err := s.purchaseReturnRepository.FindOne(metaData.Query)
	if err != nil {
		return err
	}

	old_status := old_data.StatusId
	new_status := PrData.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, PrData.StatusId)
		PrData.StatusHistory = result
		final_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "FULLY_RETURNED")
		partial_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "PARTIALLY_RETURNED")
		if PrData.StatusId == final_status || PrData.StatusId == partial_status {
			dataInv := map[string]interface{}{
				"id":            data.ID,
				"order_lines":   data.PurchaseReturnLines,
				"order_type":    "purchase_return",
				"is_update_inv": false,
				"is_credit":     false,
			}

			//updating the Committed Stock
			go s.basicInventoryService.UpdateTransactionInventory(metaData, dataInv)
		}
	}
	data.UpdatedByID = &metaData.TokenUserId
	err = s.purchaseReturnRepository.Update(metaData.Query, &PrData)
	if err != nil {
		return err
	}
	for _, return_line := range PrData.PurchaseReturnLines {
		update_query := map[string]interface{}{
			"product_id": return_line.ProductId,
			"pr_id":      metaData.Query["id"],
		}
		count, err1 := s.purchaseReturnRepository.UpdateReturnLines(update_query, return_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			return_line.PrId = uint(metaData.Query["id"].(float64))
			e := s.purchaseReturnRepository.SaveReturnLines(return_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeletePurchaseReturn(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase return at data level")
	}
	err := s.purchaseReturnRepository.Delete(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{
			"pr_id": metaData.Query["id"],
		}
		err = s.purchaseReturnRepository.DeleteReturnLine(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeletePurchaseReturnLines(metaData core.MetaData) error {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase return at data level")
	}
	err := s.purchaseReturnRepository.DeleteReturnLine(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SearchPurchaseReturns(query string, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase return at data level")
	}
	data, err := s.purchaseReturnRepository.Search(query)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) GetPurchaseReturnsHistory(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "LIST", "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase return at data level")
	}
	data, err := s.purchaseReturnRepository.GetPurchaseReturnsHistory(metaData.Query, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (s *service) GetPurchaseReturnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "PURCHASE_RETURNS", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purhase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase return at data level")
	}

	id := metaData.Query["id"].(string)
	tab := metaData.Query["tab"].(string)
	// token_id := fmt.Sprint(metaData.TokenUserId)
	// access_template_id := fmt.Sprint(metaData.AccessTemplateId)
	purchaseReturnId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "debit_note" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_RETURNS")

		debitNotePage := page
		debitNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseReturnId)
		// var query = make(map[string]interface{}, 0)

		debitNoteData, err := s.debitNoteService.GetDebitNoteList(metaData, debitNotePage)
		if err != nil {
			return nil, err
		}
		return debitNoteData, nil
	}

	if tab == "purchase_orders" {
		purchaseReturnPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("PURCHASE_RETURNS_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")
		purchaseReturnPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_id\",\"=\",%v]]", purchaseReturnId, source_document_type_id)
		purchaseReturn_order_interface, err := s.ListPurchaseReturns(core.MetaData{}, purchaseReturnPagination)
		if err != nil {
			return nil, err
		}

		purchaseReturn_order := purchaseReturn_order_interface.([]returns.PurchaseReturns)
		if len(purchaseReturn_order) == 0 {
			return nil, errors.New("no source document")
		}

		var purchaseReturnOrderSourceDoc map[string]interface{}
		purchaseReturnOrderSourceDocJson := purchaseReturn_order[0].SourceDocuments
		dto, err := json.Marshal(purchaseReturnOrderSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &purchaseReturnOrderSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		purchaseReturnOrderSourceDocId := purchaseReturnOrderSourceDoc["id"]
		if purchaseReturnOrderSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		purchaseOrderPage := page
		purchaseOrderPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", purchaseReturnOrderSourceDoc["id"])
		data, err := s.purchaseOrderRepositary.FindAll(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
func (s *service) ProcessPurchaseReturnCalculation(data *PurchaseReturnsDTO) *PurchaseReturnsDTO {
	var totalQuantity int64 = 0

	for _, returnLines := range data.PurchaseReturnLines {
		totalQuantity = totalQuantity + returnLines.QuantityReturned
	}
	data.TotalQuantity = totalQuantity
	return data
}
