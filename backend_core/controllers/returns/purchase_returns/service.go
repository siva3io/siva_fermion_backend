package purchase_returns

import (
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/debitnote"
	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
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
	CreatePurchaseReturn(data *returns.PurchaseReturns, access_template_id string, token_id string) error
	ListPurchaseReturns(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error)
	ViewPurchaseReturn(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error)
	UpdatePurchaseReturn(query map[string]interface{}, data *returns.PurchaseReturns, access_template_id string, token_id string) error
	DeletePurchaseReturn(query map[string]interface{}, access_template_id string, token_id string) error
	DeletePurchaseReturnLines(query map[string]interface{}, access_template_id string, token_id string) error
	SearchPurchaseReturns(query string, access_template_id string, token_id string) (interface{}, error)
	GetPurchaseReturnsHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
	GetPurchaseReturnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	purchaseReturnRepository returns_repo.PurchaseReturn
	basicInventoryService    inv_service.Service
	debitNoteService         debitnote.Service
	purchaseOrderRepositary  orders.Purchase
}

func NewService() *service {
	purchaseReturnRepository := returns_repo.NewPurchaseReturn()
	basicInventoryService := inv_service.NewService()
	purchaseOrderRepositary := orders.NewPurchaseOrder()
	return &service{purchaseReturnRepository, basicInventoryService, debitnote.NewService(), purchaseOrderRepositary}
}

func (s *service) CreatePurchaseReturn(data *returns.PurchaseReturns, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create purchase return at data level")
	}
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	temp := data.CreatedByID
	user_id := strconv.Itoa(int(*temp))
	data.PurchaseReturnNumber = helpers.GenerateSequence("PR", user_id, "purchase_returns")

	err := s.purchaseReturnRepository.Save(data)

	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	data_inv := map[string]interface{}{
		"id":            data.ID,
		"order_lines":   data.PurchaseReturnLines,
		"order_type":    "purchase_return",
		"is_update_inv": true,
		"is_credit":     false,
	}

	//updating the inventory & Committed Stock
	go s.basicInventoryService.UpdateInventory(data_inv, token_id, access_template_id)

	return nil
}

func (s *service) ListPurchaseReturns(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase return at data level")
	}

	data, err := s.purchaseReturnRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) ViewPurchaseReturn(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase return at data level")
	}
	data, err := s.purchaseReturnRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) UpdatePurchaseReturn(query map[string]interface{}, data *returns.PurchaseReturns, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update purchase return at data level")
	}
	var find_query = map[string]interface{}{
		"id": query["id"],
	}

	found_data, er := s.purchaseReturnRepository.FindOne(find_query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	old_data := found_data.(returns.PurchaseReturns)

	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result
		final_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "FULLY_RETURNED")
		partial_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "PARTIALLY_RETURNED")
		if data.StatusId == final_status || data.StatusId == partial_status {
			data_inv := map[string]interface{}{
				"id":            data.ID,
				"order_lines":   data.PurchaseReturnLines,
				"order_type":    "purchase_return",
				"is_update_inv": false,
				"is_credit":     false,
			}

			//updating the Committed Stock
			go s.basicInventoryService.UpdateInventory(data_inv, token_id, access_template_id)
		}
	}

	err := s.purchaseReturnRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, return_line := range data.PurchaseReturnLines {
		update_query := map[string]interface{}{
			"product_id": return_line.ProductId,
			"pr_id":      uint(query["id"].(int)),
		}
		count, err1 := s.purchaseReturnRepository.UpdateReturnLines(update_query, return_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			return_line.PrId = uint(query["id"].(int))
			e := s.purchaseReturnRepository.SaveReturnLines(return_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeletePurchaseReturn(query map[string]interface{}, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase return at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.purchaseReturnRepository.FindOne(q)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.purchaseReturnRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	} else {
		query := map[string]interface{}{"id": query["id"]}
		er := s.purchaseReturnRepository.DeleteReturnLine(query)
		if er != nil {
			return res.BuildError(res.ErrDataNotFound, er)
		}
	}
	return nil
}

func (s *service) DeletePurchaseReturnLines(query map[string]interface{}, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase return at data level")
	}
	_, er := s.purchaseReturnRepository.FindReturnLines(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.purchaseReturnRepository.DeleteReturnLine(query)

	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
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

func (s *service) GetPurchaseReturnsHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PURCHASE_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase return at data level")
	}
	data, err := s.purchaseReturnRepository.GetPurchaseReturnsHistory(productId, page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}
func (s *service) GetPurchaseReturnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}

	purchaseReturnId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "debit_note" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_RETURNS")

		debitNotePage := page
		debitNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseReturnId)
		var query = make(map[string]interface{}, 0)

		debitNoteData, err := s.debitNoteService.GetDebitNoteList(query, debitNotePage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return debitNoteData, nil
	}

	if tab == "purchase_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_RETURNS_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")

		debitNotePage := page
		debitNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, purchaseReturnId)
		debitNoteData, err := s.purchaseOrderRepositary.FindAll(page)
		if err != nil {
			return nil, err
		}
		return debitNoteData, nil
	}
	return nil, nil
}
