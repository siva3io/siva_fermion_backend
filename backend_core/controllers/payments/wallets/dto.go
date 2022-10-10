package wallets

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
type WalletResponseDto struct {
	ID                   uint
	UserId               uint    `json:"user_id" gorm:""`
	UserAmount           float64 `json:"user_amount"`
	TotalUserTopupAmount float64 `json:"total_user_topup_amount"`
	Currency             string  `json:"currency"`
	Status               string  `json:"status"`
	Gateway              string  `json:"gateway"`
}
