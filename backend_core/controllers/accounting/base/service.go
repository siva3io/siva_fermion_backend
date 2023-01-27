package accounting_base

import (
	"errors"

	accounting_repo "fermion/backend_core/internal/repository/accounting"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/lib/pq"
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
type ServiceBase interface {
	//SalesInvoice
	FavouriteSalesInvoice(map[string]interface{}) error
	UnFavouriteSalesInvoice(map[string]interface{}) error

	//PurchaseInvoice
	FavouritePurchaseInvoice(map[string]interface{}) error
	UnFavouritePurchaseInvoice(map[string]interface{}) error

	//DebitNote
	FavouriteDebitNote(q map[string]interface{}) error
	UnFavouriteDebitNote(q map[string]interface{}) error

	//CreditNote
	FavouriteCreditNote(q map[string]interface{}) error
	UnFavouriteCreditNote(q map[string]interface{}) error
}

type serviceBase struct {
	accountingRepository accounting_repo.AccountingBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	accountingRepository := accounting_repo.NewAccountingBase()
	newServiceObj = &serviceBase{accountingRepository}
	return newServiceObj
}

//-----------------------SalesInvoice------------------------------------

func (s *serviceBase) FavouriteSalesInvoice(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.accountingRepository.FindOne(query)
	samp := record.SalesInvoiceIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.SalesInvoiceIds = pq.Int64Array([]int64{int64(id)})
		err := s.accountingRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.SalesInvoiceIds = append(record.SalesInvoiceIds, int64(id))
		query := map[string]interface{}{
			"user_id":           uint(user_id),
			"sales_invoice_ids": record.SalesInvoiceIds,
		}
		err := s.accountingRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteSalesInvoice(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.accountingRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.SalesInvoiceIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "sales_invoice_ids",
			"query":       "array_remove(sales_invoice_ids,?)",
		}
		err = s.accountingRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

//-----------------------PurchaseInvoice-----------------------------------

func (s *serviceBase) FavouritePurchaseInvoice(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.accountingRepository.FindOne(query)
	samp := record.PurchaseInvoiceIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.PurchaseInvoiceIds = pq.Int64Array([]int64{int64(id)})
		err := s.accountingRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.PurchaseInvoiceIds = append(record.PurchaseInvoiceIds, int64(id))
		query := map[string]interface{}{
			"user_id":              uint(user_id),
			"purchase_invoice_ids": record.PurchaseInvoiceIds,
		}
		err := s.accountingRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouritePurchaseInvoice(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.accountingRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.PurchaseInvoiceIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "purchase_invoice_ids",
			"query":       "array_remove(purchase_invoice_ids,?)",
		}
		err = s.accountingRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

//-----------------------DebitNote------------------------------------------

func (s *serviceBase) FavouriteDebitNote(q map[string]interface{}) error {
	//fmt.Println("start")
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	//fmt.Println(id, user_id)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.accountingRepository.FindOne(query)
	samp := record.DebitNoteIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		//fmt.Println("-----------------------------Create Record---------------------------------------------")
		record.UserID = uint(user_id)
		record.DebitNoteIds = pq.Int64Array([]int64{int64(id)}) //append(record.DebitNoteIds, id)
		//fmt.Println("contact id:----------------------------->", record.DebitNoteIds)
		err := s.accountingRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		//fmt.Println("----------------------------Add id in Record---")
		record.DebitNoteIds = append(record.DebitNoteIds, int64(id))
		//fmt.Println("contact id:----------------------------->", record.DebitNoteIds)
		query := map[string]interface{}{
			"user_id":     uint(user_id),
			"contact_ids": record.DebitNoteIds,
		}
		err := s.accountingRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteDebitNote(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.accountingRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.DebitNoteIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "debit_note_ids",
			"query":       "array_remove(debit_note_ids,?)",
		}
		err = s.accountingRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

//-----------------------CreditNote------------------------------------------

func (s *serviceBase) FavouriteCreditNote(q map[string]interface{}) error {
	//fmt.Println("start")
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	//fmt.Println(id, user_id)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.accountingRepository.FindOne(query)
	samp := record.CreditNoteIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		//fmt.Println("-----------------------------Create Record---------------------------------------------")
		record.UserID = uint(user_id)
		record.CreditNoteIds = pq.Int64Array([]int64{int64(id)}) //append(record.CreditNoteIds, id)
		//fmt.Println("contact id:----------------------------->", record.CreditNoteIds)
		err := s.accountingRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		//fmt.Println("----------------------------Add id in Record---")
		record.CreditNoteIds = append(record.CreditNoteIds, int64(id))
		//fmt.Println("contact id:----------------------------->", record.CreditNoteIds)
		query := map[string]interface{}{
			"user_id":     uint(user_id),
			"contact_ids": record.CreditNoteIds,
		}
		err := s.accountingRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteCreditNote(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.accountingRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.CreditNoteIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "credit_note_ids",
			"query":       "array_remove(credit_note_ids,?)",
		}
		err = s.accountingRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}
