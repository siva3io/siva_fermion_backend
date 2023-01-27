package creditnote

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/sales_invoice"
	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	accounting_repo "fermion/backend_core/internal/repository/accounting"
	orders_repo "fermion/backend_core/internal/repository/orders"
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
	CreateCreditNote(metaData core.MetaData, data *accounting.CreditNote) error
	UpdateCreditNote(metaData core.MetaData, data *accounting.CreditNote) error
	GetCreditNote(metaData core.MetaData) (interface{}, error)
	GetCreditNoteList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeleteCreditNote(metaData core.MetaData) error
	GetCreditNoteTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)

	//SearchCreditNote(query string) ([]PartnersObjDTO, error)
}

type service struct {
	creditNoteRepository  accounting_repo.CreditNotes
	salesInvoiceService   sales_invoice.Service
	salesReturnRepository returns_repo.SalesReturn
	salesorderRepository  orders_repo.Sales
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	//credit_note_repository := accounting_repo.NewCreditNote()
	newServiceObj = &service{
		accounting_repo.NewCreditNote(),
		sales_invoice.NewService(),
		returns_repo.NewSalesReturn(),
		orders_repo.NewSalesOrder()}
	return newServiceObj
}

// ============================ CRUD services =================================================================

func (s *service) CreateCreditNote(metaData core.MetaData, data *accounting.CreditNote) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "CREDIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create credit note at view level")
	}

	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create credit note at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	var dtoData CreditNoteDTO
	seqCredit := helpers.GenerateSequence("Cred", fmt.Sprint(metaData.TokenUserId), "credit_notes")
	seqRef := helpers.GenerateSequence("Ref", fmt.Sprint(metaData.TokenUserId), "credit_notes")
	if dtoData.GenerateCreditNoteId {
		dtoData.CreditNoteID = seqCredit
	}
	if dtoData.GenerateReferenceId {
		dtoData.ReferenceId = seqRef
	}
	defaultStatus, err := helpers.GetLookupcodeId("CREDIT_NOTE_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	err = s.creditNoteRepository.SaveCreditNote(data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return nil
}

func (s *service) UpdateCreditNote(metaData core.MetaData, data *accounting.CreditNote) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "CREDIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update credit note at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update credit note at data level")
	}
	data.UpdatedByID = &metaData.TokenUserId
	err := s.creditNoteRepository.UpdateCreditNote(metaData.Query, data)
	if err != nil {
		return err
	}
	for _, orderLine := range data.CreditNoteLineItems {

		query := map[string]interface{}{
			"credit_note_id":     helpers.ConvertStringToUint(metaData.Query["id"].(string)),
			"product_variant_id": uint(orderLine.ProductVariantId),
		}
		orderLine.UpdatedByID = &metaData.TokenUserId
		count, er := s.creditNoteRepository.UpdateCreditLines(query, orderLine)
		if er != nil {
			return er
		} else if count == 0 {
			orderLine.CreditNoteId = uint(metaData.Query["id"].(int))
			e := s.creditNoteRepository.SaveCreditLines(orderLine)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeleteCreditNote(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "CREDIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete credit note at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete credit note at data level")
	}
	err := s.creditNoteRepository.DeleteCreditNote(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{"credit_note_id": metaData.Query["id"]}
		er := s.creditNoteRepository.DeleteCreditLine(query)
		if er != nil {
			return er
		}
	}
	return nil
}

func (s *service) GetCreditNote(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "CREDIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view credit note at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view credit note at data level")
	}
	result, err := s.creditNoteRepository.FindOneCreditNote(metaData.Query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetCreditNoteList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "CREDIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list credit note at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list credit note at data level")
	}
	result, err := s.creditNoteRepository.FindAllCreditNote(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil

}

func (s *service) GetCreditNoteTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "CREDIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list credit note at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list credit note at data level")
	}
	tab := metaData.AdditionalFields["tab"]
	creditNoteId := metaData.AdditionalFields["id"]
	if tab == "sales_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("SALES_INVOICE_SOURCE_DOCUMENT_TYPES", "CREDIT_NOTE")

		salesInvoicePage := page
		salesInvoicePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, creditNoteId)
		data, err := s.salesInvoiceService.GetAllSalesInvoice(metaData, salesInvoicePage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	if tab == "sales_returns" {
		creditNotePagination := new(pagination.Paginatevalue)

		source_document_type_id, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_RETURNS")
		creditNotePagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", creditNoteId, source_document_type_id)
		sales_returns_interface, err := s.GetCreditNoteList(metaData, creditNotePagination)
		if err != nil {
			return nil, err
		}
		sales_returns := sales_returns_interface.([]accounting.CreditNote)
		if len(sales_returns) == 0 {
			return nil, errors.New("no source document")
		}

		var CreditNoteSourceDoc map[string]interface{}
		CreditNoteSourceDocJson := sales_returns[0].SourceDocuments
		dto, err := json.Marshal(CreditNoteSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &CreditNoteSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		CreditNoteSourceDocId := CreditNoteSourceDoc["id"]
		if CreditNoteSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		salesReturnsPage := page
		salesReturnsPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", CreditNoteSourceDoc["id"])
		data, err := s.salesReturnRepository.FindAll(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "sales_orders" {
		creditNotePagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("CREDIT_NOTE_SOURCE_DOCUMENT_TYPES", "SALES_ORDERS")
		creditNotePagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", creditNoteId, source_document_type_id)
		sales_orders_interface, err := s.GetCreditNoteList(metaData, creditNotePagination)
		if err != nil {
			return nil, err
		}
		sales_orders := sales_orders_interface.([]accounting.CreditNote)
		if len(sales_orders) == 0 {
			return nil, errors.New("no source document")
		}

		var CreditNoteSourceDoc map[string]interface{}
		CreditNoteSourceDocJson := sales_orders[0].SourceDocuments
		dto, err := json.Marshal(CreditNoteSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &CreditNoteSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		CreditNoteSourceDocId := CreditNoteSourceDoc["id"]
		if CreditNoteSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		salesOrdersPage := page
		salesOrdersPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", CreditNoteSourceDoc["id"])
		data, err := s.salesorderRepository.FindAll(nil, page)
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
