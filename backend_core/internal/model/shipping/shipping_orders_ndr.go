package shipping

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
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
type NDR struct {
	model_core.Model
	ShippingOrderId         *uint          `json:"shipping_order_id"` // Reference,awb,shipping_id,pickup address,destination address,customer details,product details,payment,channel
	ShippingOrder           *ShippingOrder `json:"shipping_order" gorm:"foreignKey:ShippingOrderId;references:ID"`
	Amount                  float64        `gorm:"type:double precision" json:"amount"`
	NDRLines                []NDRLines     `json:"ndr_lines" gorm:"foreignKey:Non_delivery_request_id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DeliveryAttemptLeft     uint           `json:"delivery_attempt_left" gorm:"type:integer"`
	LastDeliveryAttemptDate time.Time      `gorm:"type:time" json:"last_delivery_attempt_date"`
	FailureReason           string         `json:"failure_reason" gorm:"type:varchar"`
}

type NDRLines struct {
	model_core.Model
	Non_delivery_request_id uint      `json:"-"`
	NdrsId                  string    `gorm:"type:varchar" json:"ndrs_id"`
	DeliveryAttemptMade     uint      `json:"delivery_attempt_made" gorm:"type:integer"`
	FailureReason           string    `json:"failure_reason" gorm:"type:varchar"`
	NdrsRaisedDate          time.Time `gorm:"type:time" json:"ndrs_raised_date"`
	LastDeliveryAttemptDate time.Time `gorm:"type:time" json:"last_delivery_attempt_date"`
}
