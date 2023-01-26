package payment_terms_and_record_payment

import (
	"errors"
	"fmt"

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
	CreatePaymentTerm(data *accounting.PaymentTerms, token_id string, access_template_id string) error
	UpdatePaymentTerm(query map[string]interface{}, data *accounting.PaymentTerms, token_id string, access_template_id string) error
	DeletePaymentTerm(query map[string]interface{}, token_id string, access_template_id string) error
	GetPaymentTerm(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetPaymentTermlist(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]accounting.PaymentTerms, error)
}

type service struct {
	payment_Repository accounting_repo.PaymentTerm
}

func NewService() *service {
	payment_Repository := accounting_repo.NewPaymentTerm()
	return &service{payment_Repository}

}

func (s *service) CreatePaymentTerm(data *accounting.PaymentTerms, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PAYMENT_TERMS_RECORD_PAYMENT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create payment terms & record payment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create payment terms & record payment at data level")
	}
	query := map[string]interface{}{
		"payment_term_name": data.PaymentTermName,
	}
	_, err := s.payment_Repository.FindOnePaymentTerm(query)
	if err == nil {
		return res.BuildError(res.ErrDuplicate, errors.New("oops! Record already Exists"))
	} else {
		err := s.payment_Repository.CreatePaymentTerm(data)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	}
	return nil
}

func (s *service) UpdatePaymentTerm(query map[string]interface{}, data *accounting.PaymentTerms, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PAYMENT_TERMS_RECORD_PAYMENT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update payment terms & record payment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update payment terms & record payment at data level")
	}
	_, err := s.payment_Repository.FindOnePaymentTerm(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	err = s.payment_Repository.UpdatePaymentTerm(query, data)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}

func (s *service) DeletePaymentTerm(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PAYMENT_TERMS_RECORD_PAYMENT", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete payment terms & record payment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete payment terms & record payment at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.payment_Repository.FindOnePaymentTerm(q)
	if er != nil {
		return er
	}
	err := s.payment_Repository.DeletePaymentTerm(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetPaymentTerm(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PAYMENT_TERMS_RECORD_PAYMENT", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view payment terms & record payment at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view payment terms & record payment at data level")
	}
	result, err := s.payment_Repository.FindOnePaymentTerm(query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetPaymentTermlist(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]accounting.PaymentTerms, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PAYMENT_TERMS_RECORD_PAYMENT", *token_user_id)
	fmt.Println(access_action)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list payment terms & record payment at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list payment terms & record payment at data level")
	}
	result, err := s.payment_Repository.FindAllPaymentTerm(query, p)
	if err != nil {
		return result, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return result, nil
}
