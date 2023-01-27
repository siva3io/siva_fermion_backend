package shipping

import (
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"

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
type WD struct {
	model_core.Model
	ShippingOrderId        uint           `json:"shipping_order_id"`
	ShippingOrder          ShippingOrder  `json:"shipping_order" gorm:"foreignKey:ShippingOrderId;references:ID"` // Reference,awb,shipping_id, initial amount==>shipping cost
	CourierPartnerFileId   datatypes.JSON `json:"courier_partner_file_id" gorm:"type:json"`
	TransactionId          string         `json:"transaction_id" gorm:"type:varchar"`
	InitialWeightTypeId    uint           `json:"initial_weight_type_id" gorm:"type:int"`
	InitialWeightType      mdm.Uom        `json:"initial_weight_type" gorm:"foreignKey:InitialWeightTypeId; references:ID"`
	FinalWeightTypeId      uint           `json:"final_weight_type_id" gorm:"type:int"`
	FinalWeightType        mdm.Uom        `json:"final_weight_type" gorm:"foreignKey:FinalWeightTypeId; references:ID"`
	InitialWeightTaken     float64        `gorm:"type:double precision" json:"initial_weight_taken"`
	FinalWeightTaken       float64        `gorm:"type:double precision" json:"final_weight_taken"`
	InitialAmount          float64        `gorm:"type:double precision" json:"initial_amount"`
	FinalAmount            float64        `gorm:"type:double precision" json:"final_amount"`
	DiscrepancyAmount      float64        `gorm:"type:double precision" json:"discrepancy_amount"`
	WeightDiscrepancyProof datatypes.JSON `json:"Weight_discrepancy_proof" gorm:"type:json"`
}
