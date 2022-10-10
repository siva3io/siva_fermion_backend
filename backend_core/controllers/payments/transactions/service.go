package transactions

import (
	"encoding/json"
	"errors"
	"fmt"

	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/payments"
	payments_repo "fermion/backend_core/internal/repository/payments"

	//access_checker "fermion/backend_core/pkg/util/access"

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
	CreateTransaction(data *payments.Transactions) error
	ListWalletTransactions(page *pagination.Paginatevalue) (interface{}, error)
	ListGatewayTransactions(page *pagination.Paginatevalue) (interface{}, error)
	ViewTransaction(query map[string]interface{}) (interface{}, error)
	UpdateTransaction(query map[string]interface{}, data *payments.Transactions) error
	DeleteTransaction(query map[string]interface{}) error
	GetTransaction(query map[string]interface{}) (interface{}, error)
}

type service struct {
	TransactionRepository payments_repo.Transaction
}

func NewService() *service {
	TransactionRepository := payments_repo.NewTransaction()
	return &service{TransactionRepository}
}

func (s *service) CreateTransaction(data *payments.Transactions) error {

	err := s.TransactionRepository.Save(data)
	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return nil
}

func (s *service) ListWalletTransactions(page *pagination.Paginatevalue) (interface{}, error) {

	data, err := s.TransactionRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	var dto []TransactionResponseWalletDto
	jsonData, err := json.Marshal(data)

	err = json.Unmarshal(jsonData, &dto)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return dto, nil
}

func (s *service) ListGatewayTransactions(page *pagination.Paginatevalue) (interface{}, error) {

	data, err := s.TransactionRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	var dto []TransactionResponseGatewayDto
	jsonData, err := json.Marshal(data)

	err = json.Unmarshal(jsonData, &dto)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return dto, nil
}

func (s *service) ViewTransaction(query map[string]interface{}) (interface{}, error) {

	data, err := s.TransactionRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	var dto TransactionsDto
	jsonData, err := json.Marshal(data)

	err = json.Unmarshal(jsonData, &dto)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return dto, nil
}

func (s *service) UpdateTransaction(query map[string]interface{}, data *payments.Transactions) error {

	_, er := s.TransactionRepository.FindOne(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	err := s.TransactionRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
}

func (s *service) DeleteTransaction(query map[string]interface{}) error {

	_, er := s.TransactionRepository.FindOne(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.TransactionRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}

func (s *service) GetTransaction(query map[string]interface{}) (interface{}, error) {

	data, err := s.TransactionRepository.FindOne(query)
	if err != nil && err.Error() == "record not found" {
		err := errors.New("Record not found")
		return nil, err
	}
	if err != nil {
		return nil, errors.New("Invalid parameters or payload")
	}
	var dto TransactionResponseGatewayDto
	jsonData, err := json.Marshal(data)
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		return nil, errors.New("Invalid parameters or payload")
	}
	return dto, nil
}
