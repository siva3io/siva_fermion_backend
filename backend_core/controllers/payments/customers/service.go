package customers

import (
	"encoding/json"
	"errors"
	"fmt"

	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/model/payments"
	payments_repo "fermion/backend_core/internal/repository/payments"
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
	CreateCustomer(data *payments.Customers) error
	ListCustomers(page *pagination.Paginatevalue) (interface{}, error)
	ViewCustomer(query map[string]interface{}) (interface{}, error)
	UpdateCustomer(query map[string]interface{}, data *payments.Customers) error
	DeleteCustomer(query map[string]interface{}) error
	GetCustomer(query map[string]interface{}) (interface{}, error)
}

type service struct {
	CustomerRepository payments_repo.Customer
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	CustomerRepository := payments_repo.NewCustomer()
	newServiceObj = &service{CustomerRepository}
	return newServiceObj
}

func (s *service) CreateCustomer(data *payments.Customers) error {

	err := s.CustomerRepository.Save(data)
	if err != nil {
		fmt.Println("error", err.Error())
		return errors.New(err.Error())
	}

	return nil
}

func (s *service) ListCustomers(page *pagination.Paginatevalue) (interface{}, error) {

	data, err := s.CustomerRepository.FindAll(page)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	var dto []CustomersResponseDto
	jsonData, err := json.Marshal(data)
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return dto, nil
}

func (s *service) ViewCustomer(query map[string]interface{}) (interface{}, error) {

	data, err := s.CustomerRepository.FindOne(query)
	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	var dto CustomersResponseDto
	jsonData, err := json.Marshal(data)
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return dto, nil
}

func (s *service) UpdateCustomer(query map[string]interface{}, data *payments.Customers) error {

	_, er := s.CustomerRepository.FindOne(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.CustomerRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}

func (s *service) DeleteCustomer(query map[string]interface{}) error {

	_, er := s.CustomerRepository.FindOne(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.CustomerRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}

func (s *service) GetCustomer(query map[string]interface{}) (interface{}, error) {

	data, err := s.CustomerRepository.FindOne(query)
	if err != nil && err.Error() == "record not found" {
		err := errors.New("Record not found")
		return nil, err
	}
	if err != nil {
		return nil, errors.New("Invalid parameters or payload")
	}
	var dto CustomersResponseDto
	jsonData, err := json.Marshal(data)
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		return nil, errors.New("Invalid parameters or payload")
	}
	return dto, nil
}
