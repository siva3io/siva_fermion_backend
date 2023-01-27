package currency_exchange

import (
	"fmt"
	"strconv"

	accounting "fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	repo "fermion/backend_core/internal/repository/accounting"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/response"
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
	CreateCurrency(metaData core.MetaData, data *accounting.CurrencyExchange) error
	UpdateExchangePair(metaData core.MetaData, data *accounting.CurrencyExchange) error
	GetExchangePair(metaData core.MetaData, P *pagination.Paginatevalue) (interface{}, error)
	GetExchangePairData(metaData core.MetaData) (interface{}, error)
	DeleteExchangePair(metaData core.MetaData) error
}

type service struct {
	currencyExchangeRepository repo.CurrencyExchange
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	repositroy := repo.NewCurrency()
	newServiceObj = &service{repositroy}
	return newServiceObj
}

func (s *service) CreateCurrency(metaData core.MetaData, data *accounting.CurrencyExchange) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "CURRENCY_EXCHANGE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create currency exchange at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create currency exchange at data level")
	}
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	err := s.currencyExchangeRepository.CreateExchangePair(data)
	if err != nil {
		return response.BuildError(response.ErrUnprocessableEntity, err)
	}
	return nil
}
func (s *service) UpdateExchangePair(metaData core.MetaData, data *accounting.CurrencyExchange) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "CURRENCY_EXCHANGE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update currency exchange at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update currency exchange at data level")
	}
	_, err := s.currencyExchangeRepository.FindOne(metaData.Query)
	if err != nil {
		return err
	}
	err = s.currencyExchangeRepository.UpdateExchangePair(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetExchangePair(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "CURRENCY_EXCHANGE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list currency exchange at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list currency exchange at data level")
	}
	result, err := s.currencyExchangeRepository.FindAll(metaData.Query, p)
	if err != nil {
		return nil, err

	}
	return result, nil
}
func (s *service) GetExchangePairData(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "CURRENCY_EXCHANGE", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view currency exchange at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view currency exchange at data level")
	}
	result, err := s.currencyExchangeRepository.FindOne(metaData.Query)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (s *service) DeleteExchangePair(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "CURRENCY_EXCHANGE", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete currency exchange at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete currency exchange at data level")
	}
	err := s.currencyExchangeRepository.DeleteExchangePair(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
