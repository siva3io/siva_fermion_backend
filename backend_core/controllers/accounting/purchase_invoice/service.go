package purchase_invoice

import (
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/accounting"
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
	CreatePurchaseInvoice(data *accounting.PurchaseInvoice, token_id string, access_template_id string) error
	ListPurchaseInvoice(page *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
	ViewPurchaseInvoice(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	UpdatePurchaseInvoice(query map[string]interface{}, data *accounting.PurchaseInvoice, token_id string, access_template_id string) error
	DeletePurchaseInvoice(query map[string]interface{}, token_id string, access_template_id string) error
	DeletePurchaseInvoiceLines(query map[string]interface{}, token_id string, access_template_id string) error
	SearchPurchaseInvoice(query string, token_id string, access_template_id string) (interface{}, error)
	GetPurchaseInvoiceTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	purchaseInvoiceRepository accounting_repo.Purchase
	debitNoteRepositary       accounting_repo.DebitNotes
	purchaseOrdersRepositary  orders_repo.Purchase
}

func NewService() *service {
	purchaseInvoiceRepository := accounting_repo.NewPurchaseInvoice()
	return &service{purchaseInvoiceRepository, accounting_repo.NewDebitNote(), orders_repo.NewPurchaseOrder()}
}

func (s *service) CreatePurchaseInvoice(data *accounting.PurchaseInvoice, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create purchase invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create purchase invoice at data level")
	}
	temp := data.CreatedByID
	user_id := strconv.Itoa(int(*temp))
	data.PurchaseInvoiceNumber = helpers.GenerateSequence("PI", user_id, "purchase_invoices")

	result, _ := helpers.UpdateStatusHistory(data.StatusHistory, data.StatusId)
	data.StatusHistory = result
	defaultStatus, err := helpers.GetLookupcodeId("PURCHASE_INVOICE_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusId = defaultStatus
	err = s.purchaseInvoiceRepository.Save(data)

	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return nil
}

func (s *service) ListPurchaseInvoice(page *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase invoice at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase invoice at data level")
	}
	data, err := s.purchaseInvoiceRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) ViewPurchaseInvoice(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at data level")
	}
	data, err := s.purchaseInvoiceRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) UpdatePurchaseInvoice(query map[string]interface{}, data *accounting.PurchaseInvoice, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update purchase invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update purchase invoice at data level")
	}
	var find_query = map[string]interface{}{
		"id": query["id"],
	}

	found_data, er := s.purchaseInvoiceRepository.FindOne(find_query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	old_data := found_data.(accounting.PurchaseInvoice)

	old_status := old_data.StatusId
	new_status := data.StatusId

	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusId)
		data.StatusHistory = result
	}

	err := s.purchaseInvoiceRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, invoice_line := range data.PurchaseInvoiceLines {
		update_query := map[string]interface{}{
			"product_id":          invoice_line.ProductId,
			"purchase_invoice_id": uint(query["id"].(int)),
		}
		count, err1 := s.purchaseInvoiceRepository.UpdateInvoiceLines(update_query, invoice_line)
		if err1 != nil {
			return err1
		} else if count == 0 {
			invoice_line.PurchaseInvoiceId = uint(query["id"].(int))
			e := s.purchaseInvoiceRepository.SaveInvoiceLines(invoice_line)
			fmt.Println(e)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) DeletePurchaseInvoice(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase invoice at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.purchaseInvoiceRepository.FindOne(q)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.purchaseInvoiceRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	} else {
		query := map[string]interface{}{"id": query["id"]}
		er := s.purchaseInvoiceRepository.DeleteInvoiceLine(query)
		if er != nil {
			return res.BuildError(res.ErrDataNotFound, er)
		}
	}
	return nil
}

func (s *service) DeletePurchaseInvoiceLines(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete purchase invoice at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete purchase invoice at data level")
	}
	_, er := s.purchaseInvoiceRepository.FindInvoiceLines(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.purchaseInvoiceRepository.DeleteInvoiceLine(query)

	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
}

func (s *service) SearchPurchaseInvoice(query string, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list purchase invoice at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list purchase invoice at data level")
	}
	data, err := s.purchaseInvoiceRepository.Search(query)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}

func (s *service) GetPurchaseInvoiceTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PURCHASE_INVOICE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view purchase invoice at data level")
	}

	debitNoteId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "debit_note" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("DEBIT_NOTE_SOURCE_DOCUMENT_TYPES", "PURCHASE_INVOICE")

		debitNotePage := page
		fmt.Println(sourceDocumentId, debitNoteId)
		debitNotePage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, debitNoteId)
		query := make(map[string]interface{}, 0)
		data, err := s.debitNoteRepositary.FindAllDebitNote(query, debitNotePage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "purchase_orders" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES", "PURCHASE_INVOICE")

		purchaseOrderPage := page
		purchaseOrderPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, debitNoteId)
		data, err := s.purchaseOrdersRepositary.FindAll(purchaseOrderPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
