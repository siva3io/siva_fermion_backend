package currency_exchange

import (
	"fmt"

	accounting "fermion/backend_core/internal/model/accounting"
	repo "fermion/backend_core/internal/repository/accounting"
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
	CreateCurrency(data *accounting.CurrencyExchange, token_id string, access_template_id string) error
	UpdateExchangePair(query map[string]interface{}, data *accounting.CurrencyExchange, token_id string, access_template_id string) error
	GetExchangePair(token_id string, access_template_id string, access_action string) ([]accounting.CurrencyExchange, error)
	GetExchangePairData(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	DeleteExchangePair(query map[string]interface{}, token_id string, access_template_id string) error
}

type service struct {
	currencyExchangeRepository repo.CurrencyExchange
}

func NewService() *service {
	repositroy := repo.NewCurrency()
	return &service{repositroy}
}

func (s *service) CreateCurrency(data *accounting.CurrencyExchange, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "CURRENCY_EXCHANGE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create currency exchange at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create currency exchange at data level")
	}
	err := s.currencyExchangeRepository.CreateExchangePair(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateExchangePair(query map[string]interface{}, data *accounting.CurrencyExchange, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "CURRENCY_EXCHANGE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update currency exchange at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update currency exchange at data level")
	}
	_, err := s.currencyExchangeRepository.FindOne(query)
	if err != nil {
		return err
	}
	er := s.currencyExchangeRepository.UpdateExchangePair(query, data)
	if er != nil {
		return er
	}
	return nil
}
func (s *service) GetExchangePair(token_id string, access_template_id string, access_action string) ([]accounting.CurrencyExchange, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "CURRENCY_EXCHANGE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list currency exchange at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list currency exchange at data level")
	}
	var result []accounting.CurrencyExchange
	result, err := s.currencyExchangeRepository.FindAll()
	if err != nil {
		return nil, err

	}
	return result, nil
}
func (s *service) GetExchangePairData(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "CURRENCY_EXCHANGE", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view currency exchange at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view currency exchange at data level")
	}
	result, err := s.currencyExchangeRepository.FindOne(query)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (s *service) DeleteExchangePair(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "CURRENCY_EXCHANGE", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete currency exchange at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete currency exchange at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.currencyExchangeRepository.FindOne(q)
	if er != nil {
		return er
	}
	err := s.currencyExchangeRepository.DeleteExchangePair(query)
	if err != nil {
		return err
	}
	return nil
}
