package orders

import (
	model_core "fermion/backend_core/internal/model/core"

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
type UserOrdersFav struct {
	model_core.Model
	UserID              uint                 `json:"user_id"`
	User                model_core.CoreUsers `json:"user" gorm:"foreignkey:UserID;references:ID"`
	DeliveryOrderIds    pq.Int64Array        `json:"delivery_order_ids" gorm:"type:int[]"`
	InternalTransferIds pq.Int64Array        `json:"internal_transfer_ids" gorm:"type:int[]"`
	PurchaseOrderIds    pq.Int64Array        `json:"purchase_order_ids" gorm:"type:int[]"`
	SalesOrderIds       pq.Int64Array        `json:"sales_order_ids" gorm:"type:int[]"`
	ScrapOrderIds       pq.Int64Array        `json:"scrap_order_ids" gorm:"type:int[]"`
}
