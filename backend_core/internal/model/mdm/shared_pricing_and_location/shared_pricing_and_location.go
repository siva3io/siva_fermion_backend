package shared_pricing_and_location

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

//====================================================Pricings===========================================================================

type Pricing struct {
	model_core.Model
	PriceListName     string                `json:"price_list_name" gorm:"type:varchar(100)"`
	CurrencyId        *uint                 `json:"currency_id"`
	PriceListCurrency *model_core.Currency  `json:"currency" gorm:"foreignkey:CurrencyId;references:ID"`
	StartDate         datatypes.Date        `json:"start_date" gorm:"type:date"`
	EndDate           datatypes.Date        `json:"end_date" gorm:"type:date"`
	Description       string                `json:"description" gorm:"type:text"`
	PriceListRule     string                `json:"price_list_rule" gorm:"type:text"`
	StatusId          *uint                 `json:"status_id" gorm:"type:integer"`
	Status            model_core.Lookupcode `json:"status" gorm:"foreignKey:StatusId; references:ID"`

	// ========== **Price_list_id** field holds the value of which table we need to do the operation ===============================================
	// ========== 1.sales   2.purchase   3.transfer ================================================================================================
	PriceListId         uint              `json:"price_list_id"`
	SalesPriceListId    *uint             `json:"sales_price_list_id,omitempty"`
	SalesPriceList      SalesPriceList    `json:"sales_price_list,omitempty" gorm:"foreignkey:SalesPriceListId; references:ID"`
	PurchasePriceListId *uint             `json:"purchase_price_list_id,omitempty"`
	PurchasePriceList   PurchasePriceList `json:"purchase_price_list,omitempty" gorm:"foreignkey:PurchasePriceListId; references:ID"`
	TransferPriceListId *uint             `json:"transfer_price_list_id,omitempty"`
	TransferPriceList   TransferPriceList `json:"transfer_price_list,omitempty" gorm:"foreignkey:TransferPriceListId; references:ID"`
}

// ====================== Sales Price List =========================================================================================================
type SalesPriceList struct {
	model_core.Model
	CustomerName       string                `json:"customer_name" gorm:"type:varchar(100)"`
	EnterManually      *bool                 `json:"enter_manually" gorm:"type:boolean"`
	AddChannelOfSaleId *uint                 `json:"add_channel_of_sale_id" gorm:"type:integer"`
	AddChannelOfSale   model_core.Lookupcode `json:"add_channel_of_sale" gorm:"foreignKey:AddChannelOfSaleId; references:ID"`
	SelectType         *bool                 `json:"select_type" gorm:"type:boolean"`
	Percentage         float64               `json:"percentage" gorm:"type:double precision"`
	SalesLineItems     []SalesLineItems      `json:"sales_line_items,omitempty" gorm:"foreignkey:SalesPriceListId; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ShippingCost       float64               `json:"shipping_cost" gorm:"type:double precision"`
}
type SalesLineItems struct {
	model_core.Model
	SalesPriceListId    uint                  `json:"spl_id"`
	ProductVariantId    *uint                 `json:"product_id"`
	Product             mdm.ProductVariant    `json:"product_details" gorm:"foreignkey:ProductVariantId; references:ID"`
	CategoryCommission  int                   `json:"category_commission" gorm:"type:int"`
	UomId               *uint                 `json:"uom_id"`
	Uom                 mdm.Uom               `json:"uom_details" gorm:"foreignkey:UomId; references:ID"`
	QuantityValue       datatypes.JSON        `json:"quantity_value" gorm:"type:json; default:'[]'; not null"`
	QuantityValueTypeId *uint                 `json:"quantity_value_type_id" gorm:"type:integer"`
	QuantityValueType   model_core.Lookupcode `json:"quantity_value_type" gorm:"foreignKey:QuantityValueTypeId; references:ID"`
	Mrp                 float64               `json:"mrp"`
	SaleRate            float64               `json:"sale_rate"`
	Duties              float64               `json:"duties"` //product_taxes
	PricingOptionsId    *uint                 `json:"pricing_options_id" gorm:"type:integer"`
	PricingOptions      model_core.Lookupcode `json:"pricing_options" gorm:"foreignKey:PricingOptionsId; references:ID"`
	CostPrice           float64               `json:"price"`
}

// ====================== Purchase Price List =====================================================================================================
type PurchasePriceList struct {
	model_core.Model
	VendorId                 *uint                    `json:"vendor_name_id"`
	VendorDetails            *mdm.Vendors             `json:"vendor_name" gorm:"foreignkey:VendorId; references:ID"`
	PurchaseLineItems        []PurchaseLineItems      `json:"purchase_line_items,omitempty" gorm:"foreignkey:PurchasePriceListId; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PurchaseListOtherdetails PurchaseListOtherdetails `json:"other_details,omitempty" gorm:"embedded"`
}
type PurchaseLineItems struct {
	model_core.Model
	PurchasePriceListId  uint                  `json:"ppl_id"`
	ProductVariantId     *uint                 `json:"product_id"`
	Product              mdm.ProductVariant    `json:"product_details" gorm:"foreignkey:ProductVariantId; references:ID"`
	MinimumOrderQuantity int                   `json:"minimum_order_quantity" gorm:"type:int"`
	QuantityValue        datatypes.JSON        `json:"quantity_value" gorm:"type:json; default:'[]'; not null"`
	QuantityValueTypeId  *uint                 `json:"quantity_value_type_id" gorm:"type:integer"`
	QuantityValueType    model_core.Lookupcode `json:"quantity_value_type" gorm:"foreignKey:QuantityValueTypeId; references:ID"`
	CostPrice            float64               `json:"price" gorm:"type:double precision"`
	Mrp                  float64               `json:"mrp"`
	PriceQuantity        float64               `json:"price_quantity" gorm:"type:double precision"`
	VendorRate           int                   `json:"vendor_rate" gorm:"type:int"`
	SalesPeriod          string                `json:"sales_period" gorm:"type:varchar(100)"`
	CreditPeriod         string                `json:"credit_period" gorm:"type:varchar(100)"`
	ExpectedDeliveryTime string                `json:"expected_delivery_time" gorm:"type:varchar(100)"`
	LeadTime             string                `json:"lead_time" gorm:"type:varchar(100)"`
}
type PurchaseListOtherdetails struct {
	ShippingCost int `json:"shipping_cost" gorm:"type:int"`
}

// ======================= Transfer Price List ====================================================================================================
type TransferPriceList struct {
	model_core.Model
	ContractDetails          ContractDetails          `json:"contract_details" gorm:"embedded"`
	TransferLineItems        []TransferLineItems      `json:"transfer_line_items" gorm:"foreignkey:TransferPriceListId; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransferListOtherDetails TransferListOtherDetails `json:"transfer_list_other_details" gorm:"embedded"`
}

type ContractDetails struct {
	SenderName   string `json:"sender_name" gorm:"type:varchar(100)"`
	ReceiverName string `json:"receiver_name" gorm:"type:varchar(100)"`
}

type TransferLineItems struct {
	TransferPriceListId uint               `json:"tpl_id"`
	ProductVariantId    *uint              `json:"product_id"`
	Product             mdm.ProductVariant `json:"product_details" gorm:"foreignkey:ProductVariantId; references:ID"`
	Price               float64            `json:"price" gorm:"type:double precision"`
	PriceQuantity       float64            `json:"price_quantity" gorm:"type:double precision"`
	ProductRate         float64            `json:"product_rate" gorm:"type:double precision"`
	model_core.Model
}

type TransferListOtherDetails struct {
	FromAddressLocationId *uint      `json:"from_address_location_id"`
	LocationFromAddress   *Locations `json:"location_from_address" gorm:"foreignkey:FromAddressLocationId; references:ID"`
	ToAddressLocationId   *uint      `json:"to_address_location_id"`
	LocationToAddress     *Locations `json:"location_to_address" gorm:"foreignkey:ToAddressLocationId; references:ID"`
	SalesPeriod           string     `json:"sales_period" gorm:"type:varchar(100)"`
	CreditPeriod          string     `json:"credit_period" gorm:"type:varchar(100)"`
	ExpectedDeliveryTime  string     `json:"expected_delivery_time" gorm:"type:varchar(100)"`
	LeadTime              string     `json:"lead_time" gorm:"type:varchar(100)"`
	ShippingCost          int        `json:"shipping_cost" gorm:"type:int"`
}

//====================================================Locations===========================================================================

type Locations struct {
	model_core.Model
	Name               string                 `json:"name" gorm:""`
	LocationTypeId     uint                   `json:"location_type_id" gorm:""`
	LocationType       *model_core.Lookupcode `json:"location_type" gorm:"foreignkey:LocationTypeId;references:ID"`
	LocationCode       string                 `json:"location_id" gorm:""`
	ParentLocationId   *uint                  `json:"parent_id"`
	Parent             *Locations             `json:"parent_location" gorm:"foreignkey:ParentLocationId"`
	ChildLocationIds   []*Locations           `json:"child_locations" gorm:"foreignkey:ParentLocationId"`
	Address            datatypes.JSON         `json:"address" gorm:"type:json; default:'[]'"`
	LocationDocs       datatypes.JSON         `json:"location_docs" gorm:"type:json; default:'{}'"`
	Latitude           float32                `json:"latitude" gorm:""`
	Longitude          float32                `json:"longitude" gorm:""`
	ServiceableAreaIds datatypes.JSON         `json:"serviceable_area_ids" gorm:"type:json; default:'[]'"`
	RelatedLocationID  uint                   `json:"related_location_id" gorm:""`
	Email              string                 `json:"email"`
	MobileNumber       string                 `json:"mobile_number"`
	Notes              string                 `json:"notes"`
}

//-------------------------------------Location Types - [Virtual Location, Local Warehouse, Office or Factory , Retail]----------------------------------------------------------------------------------

type VirtualLocation struct {
	model_core.Model
	ExternalDetails          datatypes.JSON `json:"external_details" gorm:"type:json; default:'{}'"`
	PaymentMapping           datatypes.JSON `json:"payment_mapping" gorm:"type:json; default:'[]'"`
	LocationInchargerDetails datatypes.JSON `json:"location_incharge_details" gorm:"type:json; default:'{}'"`
}
type LocalWarehouse struct {
	model_core.Model
	ContactDetails             datatypes.JSON    `json:"contact_details" gorm:"type:json; default:'[]'"`
	IsScrapLocation            bool              `json:"is_scrap_location" gorm:""`
	IsReturnLocation           bool              `json:"is_return_location" gorm:""`
	ShippingPartners           datatypes.JSON    `json:"shipping_partners" gorm:"type:json; default:'{}'"`
	IntegratedChannels         datatypes.JSON    `json:"integrated_channels" gorm:"type:json; default:'{}'"`
	WarehouseStorageManagement datatypes.JSON    `json:"warehouse_storage_management" gorm:"type:json; default:'[]'"`
	Zone                       []StorageLocation `json:"zone" gorm:"foreignkey:LocalWarehouseId;references:ID"`
	// RacksCapacity              uint           `json:"racks_capacity" gorm:""`
	// ShelvesCapacity            uint           `json:"shelves_capacity" gorm:""`
	// BinsCapacity               uint           `json:"bins_capacity" gorm:""`
	// RacksUomId                 *uint          `json:"racks_uom_id" gorm:""`
	// RacksUom                   mdm.Uom        `json:"racks_uom" gorm:"foreignkey:RacksUomId;references:ID"`
	// ShelvesUomId               *uint          `json:"shelves_uom_id" gorm:""`
	// ShelvesUom                 mdm.Uom        `json:"shelves_uom" gorm:"foreignkey:ShelvesUomId;references:ID"`
	// BinsUomId                  *uint          `json:"bins_uom_id" gorm:""`
	// BinsUom                    mdm.Uom        `json:"bins_uom" gorm:"foreignkey:BinsUomId;references:ID"`
	// WarehouseCode              string         `json:"warehouse_code" gorm:""`
	// Auth                       datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"auth"`
	// ShippingPartner            shipping.ShippingPartner `json:"shipping_partner" gorm:"foreignkey:ShippingPartnerID;references:ID"`
}

type StorageLocation struct {
	model_core.Model
	ParentStorageLocationId *uint                 `json:"parent_storage_location_id"`
	ParentStorageLocation   *StorageLocation      `json:"parent_storage_location" gorm:"foreignkey:ParentStorageLocationId;references:ID"`
	LocalWarehouseId        *uint                 `json:"local_warehouse_id"`
	LocalWarehouse          *LocalWarehouse       `json:"local_warehouse"`
	ZoneId                  *uint                 `json:"zone_id"`
	Zone                    *StorageLocation      `json:"zone" gorm:"foreignkey:ZoneId;references:ID"`
	RowId                   *uint                 `json:"row_id"`
	Row                     *StorageLocation      `json:"row" gorm:"foreignkey:RowId;references:ID"`
	RackId                  *uint                 `json:"rack_id"`
	Rack                    *StorageLocation      `json:"rack" gorm:"foreignkey:RackId;references:ID"`
	ShelfId                 *uint                 `json:"shelf_id"`
	Shelf                   *StorageLocation      `json:"shelf" gorm:"foreignkey:ShelfId;references:ID"`
	Capacity                int64                 `json:"capacity"`
	UomId                   *uint                 `json:"uom_id"`
	Uom                     mdm.Uom               `json:"uom"`
	Name                    string                `json:"name"`
	TypeId                  uint                  `json:"type_id"`
	Type                    model_core.Lookupcode `json:"type"`
	ChildStorageLocations   []*StorageLocation    `json:"child_storage_locations" gorm:"foreignkey:ParentStorageLocationId"`
	StorageQuantity         []*StorageQuantity    `json:"storage_quantity" gorm:"foreignkey:StorageLocationId;references:ID"`
}

type StorageQuantity struct {
	model_core.Model
	StorageLocationId uint                `json:"storage_location_id"`
	ProductVariantId  *uint               `json:"product_variant_id"`
	ProductVariant    *mdm.ProductVariant `json:"product_variant" gorm:"foreignkey:ProductVariantId;references:ID"`
	StorageCount      int64               `json:"stock_count" gorm:"default:0"`
}

type Office struct {
	model_core.Model
	LocationInchargerDetails datatypes.JSON `json:"location_incharge_details" gorm:"type:json; default:'{}'"`
}
type Retail struct {
	model_core.Model
	Storename                string                   `json:"store_name" gorm:""`
	CurrencyId               *uint                    `json:"currency_id" gorm:""`
	Currency                 model_core.Currency      `json:"currency" gorm:"foreignkey:CurrencyId;references:ID"`
	PriceListId              uint                     `json:"price_list_id" gorm:""`
	PriceList                datatypes.JSON           `json:"price_list" gorm:"type:json; default:'{}'"`
	LinkedFacilityId         *uint                    `json:"linked_facility_id"`
	LinkedFacility           Locations                `json:"linked_facility" gorm:"foreignkey:LinkedFacilityId;references:ID"`
	AllowBackOrder           bool                     `json:"allow_back_order" gorm:""`
	OrderTagIds              []*model_core.Lookupcode `json:"order_tag_ids" gorm:"many2many:retail_order_Tags;foreignkey:ID;references:ID"`
	PriceIncludesTaxId       *uint                    `json:"price_includes_tax_id" gorm:""`
	PriceIncludesTax         model_core.Lookupcode    `json:"price_includes_tax" gorm:"foreignkey:PriceIncludesTaxId;references:ID"`
	EmailNotification        bool                     `json:"email_notification" gorm:""`
	PartialFulfilment        bool                     `json:"partial_fulfilment" gorm:""`
	InventoryTypeId          *uint                    `json:"inventory_type_id" gorm:""`
	InventoryType            model_core.Lookupcode    `json:"inventory_type" gorm:"foreignkey:InventoryTypeId;references:ID"`
	SourceFacilityId         *uint                    `json:"source_facility_id"`
	SourceFacility           Locations                `json:"source_facility" gorm:"foreignkey:SourceFacilityId;references:ID"`
	PaymentMapping           datatypes.JSON           `json:"payment_mapping" gorm:"type:json; default:'[]'"`
	LocationInchargerDetails datatypes.JSON           `json:"location_incharge_details" gorm:"type:json; default:'{}'"`
}
