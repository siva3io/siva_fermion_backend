package payment_terms_and_record_payment

import (
	"fmt"

	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	accounting_repo "fermion/backend_core/internal/repository/accounting"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
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
	CreatePaymentTerm(metaData core.MetaData, data *accounting.PaymentTerms) error
	UpdatePaymentTerm(metaData core.MetaData, data *accounting.PaymentTerms) error
	DeletePaymentTerm(metaData core.MetaData) error
	GetPaymentTerm(metaData core.MetaData) (interface{}, error)
	GetPaymentTermlist(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	payment_Repository accounting_repo.PaymentTerm
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	//payment_Repository := accounting_repo.NewPaymentTerm()
	newServiceObj = &service{accounting_repo.NewPaymentTerm()}
	return newServiceObj

}

func (s *service) CreatePaymentTerm(metaData core.MetaData, data *accounting.PaymentTerms) error {

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "CREATE", "PAYMENT_TERMS_RECORD_PAYMENT", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create payment terms & record payment at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create payment terms & record payment at data level")
	}
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	err := s.payment_Repository.CreatePaymentTerm(data)
	if err == nil {
		return err
	}
	return nil
}

func (s *service) UpdatePaymentTerm(metaData core.MetaData, data *accounting.PaymentTerms) error {

	accessmoduleflag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "UPDATE", "PAYMENT_TERMS_RECORD_PAYMENT", metaData.TokenUserId)
	if !accessmoduleflag {
		return fmt.Errorf("you dont have access for update payment terms & record payment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update payment terms & record payment at data level")
	}

	data.UpdatedByID = &metaData.TokenUserId
	err := s.payment_Repository.UpdatePaymentTerm(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeletePaymentTerm(metaData core.MetaData) error {

	accessmoduleflag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "DELETE", "PAYMENT_TERMS_RECORD_PAYMENT", metaData.TokenUserId)
	if !accessmoduleflag {
		return fmt.Errorf("you dont have access for delete payment terms & record payment at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete payment terms & record payment at data level")
	}
	err := s.payment_Repository.DeletePaymentTerm(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetPaymentTerm(metaData core.MetaData) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), "READ", "PAYMENT_TERMS_RECORD_PAYMENT", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view payment terms & record payment at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view payment terms & record payment at data level")
	}
	response, err := s.payment_Repository.FindOnePaymentTerm(metaData.Query)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *service) GetPaymentTermlist(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {

	access_module_flag, data_access := access_checker.ValidateUserAccess(fmt.Sprint(metaData.AccessTemplateId), metaData.ModuleAccessAction, "PAYMENT_TERMS_RECORD_PAYMENT", metaData.TokenUserId)
	//fmt.Println(metaData.ModuleAccessAction)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list payment terms & record payment at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list payment terms & record payment at data level")
	}
	result, err := s.payment_Repository.FindAllPaymentTerm(metaData.Query, p)
	var response []PaymentTermsDTO
	helpers.JsonMarshaller(result, &response)
	return response, err
}
