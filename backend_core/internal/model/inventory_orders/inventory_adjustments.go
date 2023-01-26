package inventory_orders

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/shipping"

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
type InventoryAdjustments struct {
	ReasonID                    uint                                  `json:"reason_id"`
	Reason                      model_core.Lookupcode                 `json:"reason" gorm:"foreignkey:ReasonID; references:ID"`
	AdjustmentTypeID            uint                                  `json:"adjustment_type_id"`
	AdjustmentType              model_core.Lookupcode                 `json:"adjustment_type" gorm:"foreignkey:AdjustmentTypeID; references:ID"`
	ReferenceNumber             string                                `json:"reference_number" gorm:"type:varchar(50)"`
	AdjustmentDate              time.Time                             `json:"adjustment_date" gorm:"type:time"`
	WarehouseID                 uint                                  `json:"warehouse_id"`
	Warehouse                   shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseID; references:ID"`
	PartnerID                   uint                                  `json:"partner_id"`
	Partner                     mdm.Partner                           `json:"partner" gorm:"foreignkey:PartnerID; references:ID"`
	StatusID                    uint                                  `json:"status_id" gorm:"type:integer"`
	Status                      model_core.Lookupcode                 `json:"status" gorm:"foreignKey:StatusID"`
	StatusHistory               datatypes.JSON                        `json:"status_history" gorm:"type:json"`
	InventoryAdjustmentLines    []InventoryAdjustmentLines            `json:"inventory_adjustment_lines" gorm:"foreignKey:Inv_Adj_Id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TotalChangeInInventory      uint                                  `json:"total_change_in_inventory" gorm:"type:integer"`
	TotalChangeInInventoryCount uint                                  `json:"total_change_in_inventory_count" gorm:"type:integer"`
	InternalNotes               string                                `json:"internal_notes" gorm:"type:varchar(250)"`
	ExternalNotes               string                                `json:"external_notes" gorm:"type:varchar(250)"`
	FileAttachId                datatypes.JSON                        `json:"file_attach_id" gorm:"type:json"`
	ShippingOrderID             uint                                  `json:"shipping_order_id"`
	ShippingOrder               shipping.ShippingOrder                `json:"shipping_order" gorm:"foreignkey:ShippingOrderID; references:ID"`
	ShippingAddress             datatypes.JSON                        `json:"shipping_address" gorm:"type:json"`
	PartnerAddress              datatypes.JSON                        `json:"partner_address" gorm:"type:json"`
	WarehouseAddress            datatypes.JSON                        `json:"warehouse_address" gorm:"type:json"`
	SourceDocuments             datatypes.JSON                        `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	SourceDocumentTypeId        *uint                                 `json:"source_document_type_id" gorm:"type:integer"`
	SourceDocument              model_core.Lookupcode                 `json:"source_document" gorm:"foreignKey:SourceDocumentTypeId; references:ID"`
	model_core.Model
}

type InventoryAdjustmentLines struct {
	Inv_Adj_Id       uint                `json:"-"`
	ProductID        uint                `json:"product_id"`
	Product          mdm.ProductTemplate `json:"product" gorm:"foreignkey:ProductID; references:ID"`
	ProductVariantID uint                `json:"product_variant_id"`
	ProductVariant   mdm.ProductVariant  `json:"product_variant" gorm:"foreignKey:ProductVariantID; references:ID"`
	Description      string              `json:"description" gorm:"type:varchar(250)"`
	StockInHand      uint                `json:"stock_in_hand" gorm:"type:integer"`
	AdjustedQuantity uint                `json:"adjusted_quantity" gorm:"type:integer"`
	AdjustedPrice    float32             `json:"adjusted_price" gorm:"type:double precision"`
	BalanceQuantity  uint                `json:"balance_quantity" gorm:"type:integer"`
	UnitPrice        uint                `json:"unit_price" gorm:"type:integer"`
	model_core.Model
}
