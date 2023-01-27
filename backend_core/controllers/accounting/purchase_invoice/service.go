package purchase_invoice

import (
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	accounting_repo "fermion/backend_core/internal/repository/accounting"
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
	CreatePurchaseInvoice(metaData core.MetaData, data *accounting.PurchaseInvoice) error
	UpdatePurchaseInvoice(metaData core.MetaData, data accounting.PurchaseInvoice) error
	DeletePurchaseInvoice(metaData core.MetaData) error
	ListPurchaseInvoice(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
	ViewPurchaseInvoice(metaData core.MetaData) (interface{}, error)
	DeletePurchaseInvoiceLines(metaData core.MetaData) error
	GetPurchaseInvoiceTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	purchaseInvoiceRepository accounting_repo.Purchase
	debitNoteRepositary       accounting_repo.DebitNotes
	purchaseOrdersRepositary  orders_repo.Purchase
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	//purchaseInvoiceRepository := accounting_repo.NewPurchaseInvoice()
	newServiceObj = &service{
		accounting_repo.NewPurchaseInvoice(),
		accounting_repo.NewDebitNote(),
		orders_repo.NewPurchaseOrder()}
	return newServiceObj
}

func (s *service) CreatePurchaseInvoice(metaData core.MetaData, data *accounting.PurchaseInvoice) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PURCHASE_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create purchase invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create purchase invoice at data level")
	}
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	data.PurchaseInvoiceNumber = helpers.GenerateSequence("PI", fmt.Sprint(metaData.TokenUserId), "purchase_invoices")
	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	defaultStatus, err := helpers.GetLookupcodeId("PURCHASE_INVOICE_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	err = s.purchaseInvoiceRepository.Save(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdatePurchaseInvoice(metaData core.MetaData, data accounting.PurchaseInvoice) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "PURCHASE_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update purchase invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update purchase invoice at data level")
	}

	var find_query = map[string]interface{}{
		"id": metaData.Query["id"],
	}
	data.UpdatedByID = &metaData.TokenUserId
	found_data, er := s.purchaseInvoiceRepository.FindOne(find_query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	//found_data:=found_data_model.(*accounting.PurchaseInvoice)
	old_data := found_data //.(accounting.PurchaseInvoice)
	old_status := old_data.StatusId
	new_status := data.StatusId
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result
	}
	err := s.purchaseInvoiceRepository.Update(metaData.Query, &data)
	if er != nil {
		return err
	}
	for _, invoice_line := range data.PurchaseInvoiceLines {
		update_query := map[string]interface{}{
			"product_id":          invoice_line.ProductId,
			"purchase_invoice_id": found_data.ID,
		}
		invoice_line.UpdatedByID = &metaData.TokenUserId
		count, err1 := s.purchaseInvoiceRepository.UpdateInvoiceLines(update_query, invoice_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			invoice_line.PurchaseInvoiceId = found_data.ID
			invoice_line.CreatedByID = &metaData.TokenUserId
			invoice_line.CompanyId = metaData.CompanyId
			invoice_line.UpdatedByID = nil
			e := s.purchaseInvoiceRepository.SaveInvoiceLines(invoice_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
func (s *service) ListPurchaseInvoice(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "PURCHASE_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list purchase invoice at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase invoice at data level")
	}
	data, err := s.purchaseInvoiceRepository.FindAll(metaData.Query, page)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return data, nil
}

func (s *service) ViewPurchaseInvoice(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "PURCHASE_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at data level")
	}
	data, err := s.purchaseInvoiceRepository.FindOne(metaData.Query)
	if err != nil {
		return data, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return data, nil
}

func (s *service) DeletePurchaseInvoice(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "PURCHASE_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete purchase invoice at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete purchase invoice at data level")
	}
	err := s.purchaseInvoiceRepository.Delete(metaData.Query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	} else {
		query := map[string]interface{}{"purchase_invoice_id": metaData.Query["id"]}
		er := s.purchaseInvoiceRepository.DeleteInvoiceLine(query)
		if er != nil {
			return res.BuildError(res.ErrDataNotFound, er)
		}
	}
	return nil
}

func (s *service) DeletePurchaseInvoiceLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "PURCHASE_INVOICE", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase invoice at data level")
	}
	err := s.purchaseInvoiceRepository.DeleteInvoiceLine(metaData.Query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}

func (s *service) GetPurchaseInvoiceTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "PURCHASE_INVOICE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at data level")
	}
	tab := metaData.AdditionalFields["tab"].(string)
	debitNoteId := metaData.AdditionalFields["id"]
	if tab == "debit_note" {
		sourceDocumentId, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_INVOICE")
		debitNotePage := page
		fmt.Println(sourceDocumentId, debitNoteId)
		debitNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, debitNoteId)
		data, err := s.debitNoteRepositary.FindAllDebitNote(metaData, debitNotePage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "purchase_orders" {
		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES", "PURCHASE_INVOICE")
		purchaseOrderPage := page
		purchaseOrderPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, debitNoteId)
		data, err := s.purchaseOrdersRepositary.FindAll(nil, purchaseOrderPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
