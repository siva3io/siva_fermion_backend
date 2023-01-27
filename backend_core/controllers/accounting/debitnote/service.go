package debitnote

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/accounting/purchase_invoice"
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
	CreateDebitNote(metaData core.MetaData, data *accounting.DebitNote) error
	UpdateDebitNote(metaData core.MetaData, data *accounting.DebitNote) error
	DeleteDebitNote(metaData core.MetaData) error
	GetDebitNote(metaData core.MetaData) (interface{}, error)
	GetDebitNoteList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	GetDebitNoteTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	//SearchDebitNote(query string) ([]PartnersObjDTO, error)
}

type service struct {
	debitNoteRepository       accounting_repo.DebitNotes
	purchaseInvoiceService    purchase_invoice.Service
	purchaseReturnsRepository returns_repo.PurchaseReturn
	purchaseOrdersRepository  orders_repo.Purchase
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	debitNoteRepository := accounting_repo.NewDebitNote()
	newServiceObj = &service{
		debitNoteRepository,
		purchase_invoice.NewService(),
		returns_repo.NewPurchaseReturn(),
		orders_repo.NewPurchaseOrder()}
	return newServiceObj
}

func (s *service) CreateDebitNote(metaData core.MetaData, data *accounting.DebitNote) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "DEBIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create debit note at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create debit note at data level")
	}
	var dtoData DebitNoteDTO
	seqDebit := helpers.GenerateSequence("debt", fmt.Sprint(metaData.TokenUserId), "debit_notes")
	seqRef := helpers.GenerateSequence("Ref", fmt.Sprint(metaData.TokenUserId), "debit_notes")
	if dtoData.GenerateDebitNoteId {
		dtoData.DebitNoteID = seqDebit
	}
	if dtoData.GenerateReferenceId {
		dtoData.ReferenceId = seqRef
	}
	defaultStatus, err := helpers.GetLookupcodeId("DEBIT_NOTE_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	searchQuery := map[string]interface{}{
		"company_id": data.CompanyId,
	}
	data.StatusId = defaultStatus
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	_, err = s.debitNoteRepository.FindOneDebitNote(searchQuery)
	if err == nil {
		return res.BuildError(res.ErrDuplicate, errors.New("oops! Record already Exists"))
	} else {
		err := s.debitNoteRepository.SaveDebitNote(data)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	}
	return nil
}

func (s *service) UpdateDebitNote(metaData core.MetaData, data *accounting.DebitNote) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "DEBIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update debit note at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update debit note at data level")
	}
	_, err := s.debitNoteRepository.FindOneDebitNote(metaData.Query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	data.UpdatedByID = &metaData.TokenUserId
	err = s.debitNoteRepository.UpdateDebitNote(metaData.Query, data)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, order_line := range data.DebitNoteLineItems {
		query := map[string]interface{}{
			"debit_note_id":      helpers.ConvertStringToUint(metaData.Query["id"].(string)),
			"product_variant_id": order_line.ProductVariantId,
		}
		count, er := s.debitNoteRepository.UpdateDebitLines(query, order_line)
		if er != nil && order_line.ID != uint(0) {
			return er
		} else if count == 0 {
			order_line.DebitNoteId = query["id"].(uint)
			e := s.debitNoteRepository.SaveDebitLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeleteDebitNote(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "DEBIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete debit note at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete debit note at data level")
	}
	err := s.debitNoteRepository.DeleteDebitNote(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{"debit_note_id": metaData.Query["id"]}
		er := s.debitNoteRepository.DeleteDebitLine(query)
		if er != nil {
			return err
		}
	}
	return nil
}
func (s *service) GetDebitNote(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "DEBIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return DebitNoteResponseDTO{}, fmt.Errorf("you dont have access for view debit note at view level")
	}
	if dataAccess == nil {
		return DebitNoteResponseDTO{}, fmt.Errorf("you dont have access for view debit note at data level")
	}
	result, err := s.debitNoteRepository.FindOneDebitNote(metaData.Query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetDebitNoteList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "DEBIT_NOTE", metaData.AccessTemplateId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list debit note at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list debit note at data level")
	}
	result, err := s.debitNoteRepository.FindAllDebitNote(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetDebitNoteTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "DEBIT_NOTE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list debit note at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list debit note at data level")
	}
	tab := metaData.AdditionalFields["tab"].(string)
	debitNoteId := metaData.AdditionalFields["id"]
	if tab == "purchase_invoice" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_INVOICE_SOURCE_DOCUMENT_TYPES", "DEBIT_NOTE")
		var metaData core.MetaData
		purchaseInvoicePage := page
		purchaseInvoicePage.Filters = fmt.Sprintf("[[\"link_source_document_type\",\"=\",%v],[\"link_source_document\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, debitNoteId)
		data, err := s.purchaseInvoiceService.ListPurchaseInvoice(metaData, purchaseInvoicePage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "purchase_returns" {
		debitPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_RETURNS")
		debitPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", debitNoteId, source_document_type_id)
		debit_order_interface, err := s.GetDebitNoteList(metaData, debitPagination)
		if err != nil {
			return nil, err
		}
		debit_order := debit_order_interface.([]accounting.DebitNote)
		if len(debit_order) == 0 {
			return nil, errors.New("no source document")
		}

		var debitOrderSourceDoc map[string]interface{}
		debitOrderSourceDocJson := debit_order[0].SourceDocuments
		dto, err := json.Marshal(debitOrderSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &debitOrderSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		debitOrderSourceDocId := debitOrderSourceDoc["id"]
		if debitOrderSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		purchaseReturnsPage := page
		purchaseReturnsPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", debitOrderSourceDoc["id"])
		data, err := s.purchaseReturnsRepository.FindAll(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "purchase_orders" {
		debitPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_ORDERS")
		debitPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", debitNoteId, source_document_type_id)
		debit_order_interface, err := s.GetDebitNoteList(metaData, debitPagination)
		if err != nil {
			return nil, err
		}
		debit_order := debit_order_interface.([]accounting.DebitNote)
		if len(debit_order) == 0 {
			return nil, errors.New("no source document")
		}

		var debitOrderSourceDoc map[string]interface{}
		debitOrderSourceDocJson := debit_order[0].SourceDocuments
		dto, err := json.Marshal(debitOrderSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &debitOrderSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		debitOrderSourceDocId := debitOrderSourceDoc["id"]
		if debitOrderSourceDocId == nil {
			return nil, errors.New("no source document")
		}
		purchaseOrdersPage := page
		purchaseOrdersPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", debitOrderSourceDoc["id"])
		data, err := s.purchaseOrdersRepository.FindAll(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

// func (s *service) SearchDebitNote(query string) ([]PartnersObjDTO, error) {

// 	result, err := s.debitNoteRepository.SearchDebitNote(query)
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
