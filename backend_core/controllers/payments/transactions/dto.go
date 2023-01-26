package transactions

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
type TransactionsDto struct {
	ID                   uint
	UserId               uint    `json:"user_id" gorm:""` //core user id
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	Status               string  `json:"status"`       //transaction status
	PaymentType          string  `json:"payment_type"` //credit, Debit
	Name                 string  `json:"name"`         //wallet recharge, wallet debit
	Description          string  `json:"description"`
	ErrorMessage         string  `json:"error_message"`
	Gateway              string  `json:"gateway"`
	PaymentStatus        string  `json:"payment_status"` //paid, pending, failed
	GatewayTransactionId string  `json:"gateway_transaction_id"`
	TransactionType      string  `json:"transaction_type"` //Wallet, Gateway
}

type TransactionResponseWalletDto struct {
	ID              uint
	UserId          uint    `json:"user_id" gorm:""` //core user id
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`       //transaction status
	PaymentType     string  `json:"payment_type"` //credit, Debit
	Name            string  `json:"name"`         //wallet recharge, wallet debit
	Description     string  `json:"description"`
	ErrorMessage    string  `json:"error_message"`
	TransactionType string  `json:"transaction_type"` //Wallet, Gateway
}

type TransactionResponseGatewayDto struct {
	ID                   uint
	UserId               uint    `json:"user_id" gorm:""` //core user id
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	Status               string  `json:"status"`       //transaction status
	PaymentType          string  `json:"payment_type"` //credit, Debit
	Name                 string  `json:"name"`         //Wallet recharge
	Description          string  `json:"description"`
	ErrorMessage         string  `json:"error_message"`
	Gateway              string  `json:"gateway"`
	PaymentStatus        string  `json:"payment_status"` //paid, pending, failed
	GatewayTransactionId string  `json:"gateway_transaction_id"`
	TransactionType      string  `json:"transaction_type"` //Wallet, Gateway
}
