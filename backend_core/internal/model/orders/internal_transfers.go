package orders

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
type InternalTransfers struct {
	model_core.Model
	IstNumber             string                                 `json:"ist_number" gorm:"type:TEXT"`
	ReferenceNumber       string                                 `json:"reference_number" gorm:"type:TEXT"`
	ScheduledDeliveryDate datatypes.Date                         `json:"scheduled_delivery_date" gorm:"type:DATE"`
	ReasonId              *uint                                  `json:"reason_id" gorm:"type:INT"`
	Reason                model_core.Lookupcode                  `json:"reason" gorm:"foreignkey:ReasonId; references:ID"`
	IstDate               datatypes.Date                         `json:"ist_date" gorm:"type:date"`
	ReceiptRoutingId      *uint                                  `json:"receipt_routing_id" gorm:"type:int"`
	ReceiptRouting        model_core.Lookupcode                  `json:"receipt_routing" gorm:"foreignkey:ReceiptRoutingId; references:ID"`
	SourceLocation        datatypes.JSON                         `json:"source_location_details" gorm:"type:json; default:'[]'; not null"`
	DestinationLocation   datatypes.JSON                         `json:"destination_location_details" gorm:"type:json; default:'[]'; not null"`
	SourceLocationId      *uint                                  `json:"source_location_id" gorm:"type:INT"`
	SourceWarehouse       *shared_pricing_and_location.Locations `json:"source_warehouse" gorm:"foreignkey:SourceLocationId; references:ID"`
	DestinationLocationId *uint                                  `json:"destination_location_id" gorm:"type:INT"`
	DestinationWarehouse  *shared_pricing_and_location.Locations `json:"destination_warehouse" gorm:"foreignkey:DestinationLocationId; references:ID"`
	ShippingDetails       datatypes.JSON                         `json:"shipping_details" gorm:"type:JSON;"`
	StatusId              *uint                                  `json:"status_id" gorm:"type:INT"`
	Status                model_core.Lookupcode                  `json:"status" gorm:"foreignkey:StatusId; references:ID"`
	NoOfItems             int                                    `json:"no_of_items" gorm:"type:INT"`
	TotalQuantity         int                                    `json:"total_quantity" gorm:"type:INT"`
	ShippingModeId        *uint                                  `json:"shipping_mode_id" gorm:"type:INT"`
	ShippingMode          model_core.Lookupcode                  `json:"shipping_mode" gorm:"foreignkey:ShippingModeId; references:ID"`
	StatusHistory         datatypes.JSON                         `json:"status_history" gorm:"type:JSON; default:'[]'"`
	InternalTransferLines []InternalTransferLines                `json:"internal_transfer_lines" gorm:"foreignkey:IstId; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`
	PickupDateAndTime     PickupDateAndTime                      `json:"pickup_date_and_time" gorm:"embedded"`
	SourceDocuments       datatypes.JSON                         `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	SourceDocumentTypeId  *uint                                  `json:"source_document_type_id" gorm:"type:integer"`
	SourceDocumentType    model_core.Lookupcode                  `json:"source_document" gorm:"foreignKey:SourceDocumentTypeId; references:ID"`
}

type InternalTransferLines struct {
	model_core.Model
	IstId             uint                           `json:"ist_id" gorm:"type:INT"`
	ProductId         *uint                          `json:"product_id" gorm:"type:INT"`
	Product           *mdm.ProductVariant            `json:"product_details" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId *uint                          `json:"product_template_id" gorm:"type:INT"`
	ProductTemplate   *mdm.ProductTemplate           `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	InventoryId       *uint                          `json:"inventory_id" gorm:"type:INT"`
	Inventory         *mdm.CentralizedBasicInventory `json:"inventory" gorm:"foreignkey:InventoryId; references:ID"`
	UomId             *uint                          `json:"uom_id" gorm:"type:INT"`
	Uom               *mdm.Uom                       `json:"uom" gorm:"foreignkey:UomId; references:ID"`
	SerialNumber      string                         `json:"serial_number"  gorm:"type:text"`
	SourceStock       int                            `json:"source_stock" gorm:"type:INT"`
	DestinationStock  int                            `json:"destination_stock" gorm:"type:INT"`
	IsScrap           bool                           `json:"is_scrap" gorm:"type:BOOLEAN"`
	TransferQuantity  int                            `json:"transfer_quantity" gorm:"type:INT"`
}

type PickupDateAndTime struct {
	PickupDate     datatypes.Date `json:"pickup_date" gorm:"type:date"`
	PickupFromTime datatypes.Time `json:"pickup_from_time" gorm:"type:time"`
	PickupToTime   datatypes.Time `json:"pickup_to_time" gorm:"type:time"`
}
