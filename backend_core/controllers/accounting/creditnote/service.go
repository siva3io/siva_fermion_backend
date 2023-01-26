package creditnote

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/sales_invoice"
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
	CreateCreditNote(data *accounting.CreditNote, token_id string, access_template_id string) error
	UpdateCreditNote(query map[string]interface{}, data *CreditNoteDTO, token_id string, access_template_id string) error
	DeleteCreditNote(query map[string]interface{}, token_id string, access_template_id string) error
	GetCreditNote(query map[string]interface{}, token_id string, access_template_id string) (CreditNoteResponseDTO, error)
	GetCreditNoteList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]CreditNoteResponseDTO, error)
	GetCreditNoteTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)

	//SearchCreditNote(query string) ([]PartnersObjDTO, error)
}

type service struct {
	credit_note_repository accounting_repo.CreditNotes
	salesInvoiceService    sales_invoice.Service
}

func NewService() *service {
	credit_note_repository := accounting_repo.NewCreditNote()
	return &service{credit_note_repository, sales_invoice.NewService()}
}

func (s *service) CreateCreditNote(data *accounting.CreditNote, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "CREDIT_NOTE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create credit note at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create credit note at data level")
	}
	//data.CreditNoteID = strings.ToLower(data.PrimaryEmail)
	query := map[string]interface{}{
		"credit_note_id": data.CreditNoteID,
	}
	defaultStatus, err := helpers.GetLookupcodeId("CREDIT_NOTE_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	_, err = s.credit_note_repository.FindOneCreditNote(query)
	if err == nil {
		return res.BuildError(res.ErrDuplicate, errors.New("oops! Record already Exists"))
	} else {
		err := s.credit_note_repository.SaveCreditNote(data)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	}
	return nil
}

func (s *service) UpdateCreditNote(query map[string]interface{}, data *CreditNoteDTO, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "CREDIT_NOTE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update credit note at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update credit note at data level")
	}
	_, err := s.credit_note_repository.FindOneCreditNote(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	var res accounting.CreditNote
	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &res)
	if err != nil {
		return err
	}
	err = s.credit_note_repository.UpdateCreditNote(query, &res)
	if err != nil {
		return err
	}
	for _, order_line := range res.CreditNoteLineItems {
		query := map[string]interface{}{
			"credit_note_id":     int(query["id"].(int)),
			"product_variant_id": uint(order_line.ProductVariantId),
		}
		count, er := s.credit_note_repository.UpdateCreditLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.CreditNoteId = uint(query["credit_note_id"].(int))
			e := s.credit_note_repository.SaveCreditLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeleteCreditNote(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "CREDIT_NOTE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete credit note at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete credit note at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.credit_note_repository.FindOneCreditNote(q)
	if er != nil {
		return er
	}
	err := s.credit_note_repository.DeleteCreditNote(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetCreditNote(query map[string]interface{}, token_id string, access_template_id string) (CreditNoteResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "CREDIT_NOTE", *token_user_id)
	if !access_module_flag {
		return CreditNoteResponseDTO{}, fmt.Errorf("you dont have access for view credit note at view level")
	}
	if data_access == nil {
		return CreditNoteResponseDTO{}, fmt.Errorf("you dont have access for view credit note at data level")
	}

	var res CreditNoteResponseDTO
	result, _ := s.credit_note_repository.FindOneCreditNote(query)

	sdata, _ := json.Marshal(result)
	json.Unmarshal(sdata, &res)
	return res, nil
}
func (s *service) GetCreditNoteList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]CreditNoteResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "CREDIT_NOTE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list credit note at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list credit note at data level")
	}
	var resultdto []CreditNoteResponseDTO
	result, _ := s.credit_note_repository.FindAllCreditNote(query, p)
	mdata, _ := json.Marshal(result)
	json.Unmarshal(mdata, &resultdto)
	return resultdto, nil
}

func (s *service) GetCreditNoteTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "CREDIT_NOTE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list credit note at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list credit note at data level")
	}

	creditNoteId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "sales_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_INVOICE_SOURCE_DOCUMENT_TYPES", "CREDIT_NOTE")

		salesInvoicePage := page
		salesInvoicePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, creditNoteId)
		data, err := s.salesInvoiceService.GetAllSalesInvoice(salesInvoicePage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

// func (s *service) SearchCreditNote(query string) ([]PartnersObjDTO, error) {

// 	result, err := s.credit_note_repository.SearchCreditNote(query)
// 	var credit_notedetails []PartnersObjDTO
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
// 		credit_notedetails = append(credit_notedetails, data)
// 	}
// 	return credit_notedetails, nil
// }
