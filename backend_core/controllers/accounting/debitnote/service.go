package debitnote

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/purchase_invoice"
	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/pagination"
	accounting_repo "fermion/backend_core/internal/repository/accounting"
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
	CreateDebitNote(data *accounting.DebitNote, token_id string, access_template_id string) error
	UpdateDebitNote(query map[string]interface{}, data *DebitNoteDTO, token_id string, access_template_id string) error
	DeleteDebitNote(query map[string]interface{}, token_id string, access_template_id string) error
	GetDebitNote(query map[string]interface{}, token_id string, access_template_id string) (DebitNoteResponseDTO, error)
	GetDebitNoteList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]DebitNoteResponseDTO, error)
	GetDebitNoteTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
	//SearchDebitNote(query string) ([]PartnersObjDTO, error)
}

type service struct {
	debit_note_repository  accounting_repo.DebitNotes
	purchaseInvoiceService purchase_invoice.Service
}

func NewService() *service {
	debit_note_repository := accounting_repo.NewDebitNote()
	return &service{debit_note_repository, purchase_invoice.NewService()}
}

func (s *service) CreateDebitNote(data *accounting.DebitNote, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "DEBIT_NOTE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create debit note at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create debit note at data level")
	}
	//data.DebitNoteID = strings.ToLower(data.PrimaryEmail)
	query := map[string]interface{}{
		"debit_note_id": data.DebitNoteID,
	}
	defaultStatus, err := helpers.GetLookupcodeId("DEBIT_NOTE_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	_, err = s.debit_note_repository.FindOneDebitNote(query)
	if err == nil {
		return res.BuildError(res.ErrDuplicate, errors.New("oops! Record already Exists"))
	} else {
		err := s.debit_note_repository.SaveDebitNote(data)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	}
	return nil
}

func (s *service) UpdateDebitNote(query map[string]interface{}, data *DebitNoteDTO, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "DEBIT_NOTE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update debit note at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update debit note at data level")
	}
	_, err := s.debit_note_repository.FindOneDebitNote(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	var d accounting.DebitNote
	value, _ := json.Marshal(data)
	json.Unmarshal(value, &d)
	err = s.debit_note_repository.UpdateDebitNote(query, &d)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, order_line := range d.DebitNoteLineItems {
		query := map[string]interface{}{
			"debit_note_id":      int(query["id"].(int)),
			"product_variant_id": order_line.ProductVariantId,
		}
		count, er := s.debit_note_repository.UpdateDebitLines(query, order_line)
		if er != nil && order_line.ID != uint(0) {
			return er
		} else if count == 0 {
			order_line.DebitNoteId = query["id"].(uint)
			e := s.debit_note_repository.SaveDebitLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeleteDebitNote(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "DEBIT_NOTE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete debit note at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete debit note at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.debit_note_repository.FindOneDebitNote(q)
	if er != nil {
		return er
	}
	err := s.debit_note_repository.DeleteDebitNote(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetDebitNote(query map[string]interface{}, token_id string, access_template_id string) (DebitNoteResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "DEBIT_NOTE", *token_user_id)
	if !access_module_flag {
		return DebitNoteResponseDTO{}, fmt.Errorf("you dont have access for view debit note at view level")
	}
	if data_access == nil {
		return DebitNoteResponseDTO{}, fmt.Errorf("you dont have access for view debit note at data level")
	}
	result, err := s.debit_note_repository.FindOneDebitNote(query)
	var responses DebitNoteResponseDTO
	value, _ := json.Marshal(result)
	json.Unmarshal(value, &responses)
	if err != nil {
		return responses, err
	}
	return responses, nil
}

func (s *service) GetDebitNoteList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]DebitNoteResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "DEBIT_NOTE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list debit note at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list debit note at data level")
	}
	result, _ := s.debit_note_repository.FindAllDebitNote(query, p)
	var response []DebitNoteResponseDTO
	for _, data := range result {
		var debitnotedata DebitNoteResponseDTO
		value, _ := json.Marshal(data)
		json.Unmarshal(value, &debitnotedata)
		response = append(response, debitnotedata)
	}
	return response, nil
}

func (s *service) GetDebitNoteTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "DEBIT_NOTE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list debit note at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list debit note at data level")
	}

	debitNoteId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "purchase_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_INVOICE_SOURCE_DOCUMENT_TYPES", "DEBIT_NOTE")

		purchaseInvoicePage := page
		purchaseInvoicePage.Filters = fmt.Sprintf("[[\"link_source_document_type\",\"=\",%v],[\"link_source_document\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, debitNoteId)
		data, err := s.purchaseInvoiceService.ListPurchaseInvoice(purchaseInvoicePage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

// func (s *service) SearchDebitNote(query string) ([]PartnersObjDTO, error) {

// 	result, err := s.debit_note_repository.SearchDebitNote(query)
// 	var debit_notedetails []PartnersObjDTO
// 	if err != nil {
// 		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
// 	}
// 	for _, v := range result {
// 		var data PartnersObjDTO
// 		value, _ := json.Marshal(v)
// 		err := json.Unmarshal(value, &data)
// 		if err != nil {
// 			return nil, err
// 		}
// 		debit_notedetails = append(debit_notedetails, data)
// 	}
// 	return debit_notedetails, nil
// }
