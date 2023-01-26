package inventory_orders

import (
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/orders"

	"github.com/lib/pq"
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
type ASN struct {
	model_core.Model
	AsnNumber                  string                                `json:"asn_number" gorm:"type:varchar(50)"`
	ReferenceNumber            string                                `json:"reference_number" gorm:"type:varchar(50)"`
	WarehouseID                uint                                  `json:"warehouse_id"`
	Warehouse                  shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseID; references:ID"`
	TotalQuantity              uint                                  `json:"total_quantity" gorm:"type:integer"`
	StatusID                   uint                                  `json:"status_id" gorm:"type:integer"`
	Status                     model_core.Lookupcode                 `json:"status" gorm:"foreignKey:StatusID"`
	StatusHistory              datatypes.JSON                        `json:"status_history" gorm:"type:json"`
	GrnIDs                     pq.Int64Array                         `json:"grn_ids" gorm:"type:int[]"`
	ShippingModeID             uint                                  `json:"shipping_mode_id" gorm:"type:integer"`
	ShippingMode               model_core.Lookupcode                 `json:"shipping_mode" gorm:"foreignKey:ShippingModeID"`
	LinkPurchaseOrderID        uint                                  `json:"link_purchase_order_id"`
	LinkPurchaseOrder          orders.PurchaseOrders                 `json:"link_po" gorm:"foreignkey:LinkPurchaseOrderID; references:ID"`
	ScrapOrderID               uint                                  `json:"scrap_order_id" gorm:"type:integer"`
	SourceDocumentTypeID       uint                                  `json:"source_document_type_id" gorm:"type:integer"`
	SourceDocumentType         model_core.Lookupcode                 `json:"source_document_type" gorm:"foreignKey:SourceDocumentTypeID"`
	SourceDocuments            datatypes.JSON                        `json:"source_document" gorm:"type:json; default:'[]'; not null"`
	AsnOrderLines              []AsnLines                            `json:"asn_order_lines" gorm:"foreignKey:Asn_id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DispatchLocationDetails    datatypes.JSON                        `json:"dispatch_location_details" gorm:"type:json; default:'[]'; not null"`
	DestinationLocationDetails datatypes.JSON                        `json:"destination_location_details" gorm:"type:json; default:'[]'; not null"`
	ShippingDetails            datatypes.JSON                        `json:"shipping_details" gorm:"type:json; default:'[]'; not null"`
	EunimartWalletAmount       float32                               `json:"eunimart_wallet_amount" gorm:"type:double precision"`
	SchedulePickupDate         string                                `json:"schedule_pickup_date" gorm:"type:date"`
	SchedulePickupFromTime     string                                `json:"schedule_pickup_from_time" gorm:"type:varchar(50)"`
	SchedulePickupToTime       string                                `json:"schedule_pickup_to_time" gorm:"type:varchar(50)"`
	PartnerID                  uint                                  `json:"partner_id"`
	Partner                    mdm.Partner                           `json:"partner" gorm:"foreignKey:PartnerID; references:ID"`
}

type AsnLines struct {
	model_core.Model
	Asn_id              uint                  `json:"-"`
	ProductID           uint                  `json:"product_id"`
	Product             mdm.ProductTemplate   `json:"product" gorm:"foreignKey:ProductID; references:ID"`
	ProductVariantID    uint                  `json:"product_variant_id"`
	ProductVariant      mdm.ProductVariant    `json:"product_variant" gorm:"foreignKey:ProductVariantID; references:ID"`
	PackageTypeID       uint                  `json:"package_type_id"`
	PackageType         model_core.Lookupcode `json:"package_type" gorm:"foreignKey:PackageTypeID"`
	UnitPerBox          uint                  `json:"unit_per_box" gorm:"type:integer"`
	UomID               uint                  `json:"uom_id" gorm:"type:integer"`
	UOM                 mdm.Uom               `json:"uom" gorm:"foreignKey:UomID; references:ID"`
	PackageLength       float32               `json:"package_length" gorm:"type:double precision"`
	PackageWidth        float32               `json:"package_width" gorm:"type:double precision"`
	PackageHeight       float32               `json:"package_height" gorm:"type:double precision"`
	PackageWeight       float32               `json:"package_weight" gorm:"type:double precision"`
	NumberOfBoxes       uint                  `json:"no_of_boxes" gorm:"type:integer"`
	CrossDockingRequest datatypes.JSON        `json:"cross_docking_req" gorm:"type:json"`
}
