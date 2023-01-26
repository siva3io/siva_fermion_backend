package shipping

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"

	"gorm.io/datatypes"
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
type RTO struct {
	model_core.Model
	ShippingOrderId       uint                  `json:"shipping_order_id"` // Reference,awb,shipping_id, consignee,payment
	ShippingOrder         ShippingOrder         `json:"shipping_order" gorm:"foreignKey:ShippingOrderId;references:ID;"`
	BookingDate           time.Time             `gorm:"type:time" json:"booking_date"`
	ReturnDate            time.Time             `gorm:"type:time" json:"return_date"`
	RTOCost               float64               `gorm:"type:double precision" json:"rto_cost"`
	ActualDeliveryDate    time.Time             `gorm:"type:time" json:"actual_delivery_date"`
	EstimatedDeliveryDate time.Time             `gorm:"type:time" json:"estimate_delivery_date"`
	ShippingPartnerId     uint                  `json:"shipping_partner_id"`
	ShippingPartner       ShippingPartner       `json:"shipping_partner" gorm:"foreignKey:ShippingPartnerId;references:ID;"`
	Quantity              int32                 `gorm:"type:integer" json:"quantity"`
	TotalAmount           float64               `gorm:"type:double precision" json:"total_amount"`
	RTOStatusId           uint                  `json:"rto_status_id"`
	RTOStatus             model_core.Lookupcode `json:"rto_status" gorm:"foreignKey:RTOStatusId;references:ID;"`
	StatusHistory         datatypes.JSON        `json:"status_history" gorm:"type:json"`
}
