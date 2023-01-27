package wallets

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
	CreateWallet(data *payments.Wallets) error
	ListWallets(page *pagination.Paginatevalue) (interface{}, error)
	ViewWallet(query map[string]interface{}) (interface{}, error)
	UpdateWallet(query map[string]interface{}, data *payments.Wallets) error
	DeleteWallet(query map[string]interface{}) error
	GetWallet(query map[string]interface{}) (interface{}, error)
	AddMoney(query map[string]interface{}, data *payments.Wallets) error
	DeductMoney(query map[string]interface{}, data *payments.Wallets) error
}

type service struct {
	WalletsRepository payments_repo.Wallet
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	WalletsRepository := payments_repo.NewWallet()
	newServiceObj = &service{WalletsRepository}
	return newServiceObj
}

func (s *service) CreateWallet(data *payments.Wallets) error {
	walletData, err := s.WalletsRepository.FindOne(map[string]interface{}{"user_id": data.UserId, "currency": data.Currency, "gateway": data.Gateway})
	if walletData != nil {
		err := errors.New("Wallet already exist with given data")
		return res.BuildError(res.ErrDuplicate, err)
	}
	err = s.WalletsRepository.Save(data)
	if err != nil {
		fmt.Println("error", err.Error())
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return nil
}

func (s *service) ListWallets(page *pagination.Paginatevalue) (interface{}, error) {

	data, err := s.WalletsRepository.FindAll(page)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	var dto []WalletResponseDto
	jsonData, err := json.Marshal(data)

	err = json.Unmarshal(jsonData, &dto)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return dto, nil
}

func (s *service) ViewWallet(query map[string]interface{}) (interface{}, error) {

	data, err := s.WalletsRepository.FindOne(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	var dto WalletResponseDto
	jsonData, err := json.Marshal(data)

	err = json.Unmarshal(jsonData, &dto)

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return dto, nil
}

func (s *service) UpdateWallet(query map[string]interface{}, data *payments.Wallets) error {

	_, er := s.WalletsRepository.FindOne(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	err := s.WalletsRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
}

func (s *service) DeleteWallet(query map[string]interface{}) error {

	_, er := s.WalletsRepository.FindOne(query)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}
	err := s.WalletsRepository.Delete(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}

func (s *service) GetWallet(query map[string]interface{}) (interface{}, error) {

	data, err := s.WalletsRepository.FindOne(query)
	if err != nil && err.Error() == "record not found" {
		err := errors.New("Record not found")
		return nil, err
	}
	if err != nil {
		return nil, errors.New("Invalid parameters or payload")
	}
	var dto WalletResponseDto
	jsonData, err := json.Marshal(data)
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		return nil, errors.New("Invalid parameters or payload")
	}
	return dto, nil
}

func (s *service) AddMoney(query map[string]interface{}, data *payments.Wallets) error {

	query["gateway"] = data.Gateway
	walletDataInterface, er := s.WalletsRepository.FindOne(query)
	// fmt.Println(walletDataInterface)
	walletData := walletDataInterface.(payments.Wallets)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	data.TotalUserTopupAmount = walletData.TotalUserTopupAmount + data.UserAmount
	data.UserAmount += walletData.UserAmount

	err := s.WalletsRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
}

func (s *service) DeductMoney(query map[string]interface{}, data *payments.Wallets) error {

	query["gateway"] = data.Gateway
	walletDataInterface, er := s.WalletsRepository.FindOne(query)
	walletData := walletDataInterface.(payments.Wallets)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, er)
	}

	data.UserAmount = walletData.UserAmount - data.UserAmount

	err := s.WalletsRepository.Update(query, data)
	if er != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}

	return nil
}
