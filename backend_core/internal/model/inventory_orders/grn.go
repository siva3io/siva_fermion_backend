package inventory_orders

import (
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"

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
type GRN struct {
	GRNNumber            string                                `json:"grn_number" gorm:"type:varchar(100)"`
	ReferenceNumber      string                                `json:"reference_number" gorm:"type:varchar(100)"`
	SourceDocumentTypeId uint                                  `json:"source_document_type_id"` // asn, ist, purchaseOrder, purchaseReturn
	DocumentType         model_core.Lookupcode                 `json:"document_type" gorm:"foreignKey:SourceDocumentTypeId;references:ID"`
	SourceDocuments      datatypes.JSON                        `json:"source_document" gorm:"type:json; default:'[]'; not null"`
	CreateScrapOrder     bool                                  `json:"create_scrap_order" gorm:"type:boolean"`
	GRNOrderLines        []GRNOrderLines                       `json:"grn_order_lines" gorm:"foreignkey:GRN_ID; references:ID"`
	StatusId             uint                                  `json:"status_id"`
	Status               model_core.Lookupcode                 `json:"status" gorm:"foreignKey:StatusId;references:ID"`
	StatusHistory        datatypes.JSON                        `json:"status_history" gorm:"type:json"`
	WarehouseID          *uint                                 `json:"warehouse_id"`
	Warehouse            shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseID;references:ID"`
	model_core.Model
}

type GRNOrderLines struct {
	GRN_ID             uint                `json:"grn_id" gorm:"type:integer"`
	ProductID          uint                `json:"product_id" gorm:"type:integer"`
	Product            mdm.ProductVariant  `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	ProductTemplateId  uint                `json:"product_template_id" gorm:"type:integer"`
	ProductTemplate    mdm.ProductTemplate `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	LotNumber          string              `json:"lot_number" gorm:"type:varchar(100)"`
	UOMId              uint                `json:"uom_id" gorm:"type:integer"`
	UOM                mdm.Uom             `json:"uom" gorm:"foreignKey:UOMId;references:ID"`
	OrderedUnits       int                 `json:"ordered_units" gorm:"type:integer"`
	ReceivedUnits      int                 `json:"received_units" gorm:"type:integer"`
	PendingUnits       int                 `json:"pending_units" gorm:"type:integer"`
	QualityCheck       bool                `json:"quality_check" gorm:"type:boolean"`
	ShelfLocation      string              `json:"shelf_location" gorm:"type:varchar(100)"`
	RejectedQuantities int                 `json:"rejected_quantities" gorm:"type:integer"`
	ReasonOfRejection  string              `json:"reason_of_rejection" gorm:"type:varchar(100)"`
	PutawayStatus      bool                `json:"putaway_status" gorm:"type:boolean"`
	PutawayLocation    string              `json:"putaway_location" gorm:"type:varchar(100)"`
	model_core.Model
}
