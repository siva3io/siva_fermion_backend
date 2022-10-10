package sales_returns

import (
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/creditnote"
	inv_service "fermion/backend_core/controllers/mdm/basic_inventory"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/returns"
	returns_repo "fermion/backend_core/internal/repository/returns"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"
)

/*
 Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
 All rights reserved.
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

type Service interface {
	CreateSalesReturn(data *returns.SalesReturns, access_template_id string, token_id string) error
	ListSalesReturns(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error)
	ViewSalesReturn(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error)
	UpdateSalesReturn(query map[string]interface{}, data *returns.SalesReturns, access_template_id string, token_id string) error
	DeleteSalesReturn(query map[string]interface{}, access_template_id string, token_id string) error
	DeleteSalesReturnLines(query map[string]interface{}, access_template_id string, token_id string) error
	SearchSalesReturns(query string, access_template_id string, token_id string) (interface{}, error)
	GetSalesReturnsHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
	GetSalesReturnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	salesReturnRepository returns_repo.SalesReturn
	basicInventoryService inv_service.Service
	creditNoteService     creditnote.Service
}

func NewService() *service {
	salesReturnRepository := returns_repo.NewSalesReturn()
	basicInventoryService := inv_service.NewService()
	return &service{salesReturnRepository, basicInventoryService,
		creditnote.NewService()}
}

func (s *service) CreateSalesReturn(data *returns.SalesReturns, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create sales return at data level")
	}
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	temp := data.CreatedByID
	user_id := strconv.Itoa(int(*temp))
	data.SalesReturnNumber = helpers.GenerateSequence("SR", user_id, "sales_returns")

	err := s.salesReturnRepository.Save(data)
	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	data_inv := map[string]interface{}{
		"id":            data.ID,
		"order_lines":   data.SalesReturnLines,
		"order_type":    "sales_return",
		"is_update_inv": false,
		"is_credit":     true,
	}

	//updating the Expected Stock
	go s.basicInventoryService.UpdateInventory(data_inv, token_id, access_template_id)

	return nil
}

func (s *service) ListSalesReturns(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "SALES_RETURNS", *token_user_id)
	fmt.Println(access_action)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}
	data, err := s.salesReturnRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) ViewSalesReturn(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view sales return at data level")
	}
	data, err := s.salesReturnRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) UpdateSalesReturn(query map[string]interface{}, data *returns.SalesReturns, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update sales return at data level")
	}
	var find_query = map[string]interface{}{
		"id": query["id"],
	}

	found_data, er := s.salesReturnRepository.FindOne(find_query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	old_data := found_data.(returns.SalesReturns)

	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result

		final_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "FULLY_RETURNED")
		partial_status, _ := helpers.GetLookupcodeId("RETURN_STATUS", "PARTIALLY_RETURNED")
		if data.StatusId == final_status || data.StatusId == partial_status {
			data_inv := map[string]interface{}{
				"id":            query["id"],
				"order_lines":   data.SalesReturnLines,
				"order_type":    "sales_return",
				"is_update_inv": true,
				"is_credit":     true,
			}

			//updating the inventory & Expected Stock
			go s.basicInventoryService.UpdateInventory(data_inv, token_id, access_template_id)
		}
	}

	err := s.salesReturnRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, return_line := range data.SalesReturnLines {
		update_query := map[string]interface{}{
			"product_id": return_line.ProductId,
			"sr_id":      uint(query["id"].(int)),
		}
		count, err1 := s.salesReturnRepository.UpdateReturnLines(update_query, return_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			return_line.SrId = uint(query["id"].(int))
			e := s.salesReturnRepository.SaveReturnLines(return_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeleteSalesReturn(query map[string]interface{}, access_template_id string, token_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales return at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.salesReturnRepository.FindOne(q)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.salesReturnRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	} else {
		query := map[string]interface{}{"id": query["id"]}
		er := s.salesReturnRepository.DeleteReturnLine(query)
		if er != nil {
			return res.BuildError(res.ErrDataNotFound, er)
		}
	}
	return nil
}

func (s *service) DeleteSalesReturnLines(query map[string]interface{}, access_template_id string, token_id string) error {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete sales return at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete sales return at data level")
	}
	_, er := s.salesReturnRepository.FindReturnLines(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.salesReturnRepository.DeleteReturnLine(query)

	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
}

func (s *service) SearchSalesReturns(query string, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}
	data, err := s.salesReturnRepository.Search(query)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) GetSalesReturnsHistory(productId uint, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}
	data, err := s.salesReturnRepository.GetSalesReturnsHistory(productId, page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}
func (s *service) GetSalesReturnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "SALES_RETURNS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list sales return at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list sales return at data level")
	}

	salesReturnId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "credit_note" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_RETURNS")

		creditNotePage := page
		creditNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, salesReturnId)
		var query = make(map[string]interface{}, 0)

		creditNoteData, err := s.creditNoteService.GetCreditNoteList(query, creditNotePage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return creditNoteData, nil
	}
	return nil, nil
}
