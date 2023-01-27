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
type (
	ScrapOrders struct {
		Scrap_order_number       string                                `json:"scrap_order_no" gorm:"type:varchar(50)"`
		Schedule_scrap_date      string                                `json:"schedule_scrap_date" gorm:"type:date"`
		Scrap_reason_id          uint                                  `json:"scrap_reason_id" gorm:"type:integer"`
		Scrap_reason             model_core.Lookupcode                 `json:"scrap_reason" gorm:"foreignkey:Scrap_reason_id; refersences:ID"`
		Reference_id             string                                `json:"reference_id" gorm:"type:varchar(50)"`
		Scrap_source_location_id int                                   `json:"scrap_source_location_id" gorm:"type:int"`
		Scrap_source_location    shared_pricing_and_location.Locations `json:"scrap_source_location" gorm:"foreignKey:Scrap_source_location_id; references:ID"`
		Scrap_location_id        int                                   `json:"scrap_location_id" gorm:"type:int"`
		Scrap_location           shared_pricing_and_location.Locations `json:"scrap_location" gorm:"foreignKey:Scrap_location_id; references:ID"`
		No_of_items              int                                   `json:"no_of_items,omitempty" gorm:"type:int"`
		Total_quantity           int                                   `json:"total_quantity,omitempty" gorm:"type:int"`
		Order_lines              []ScrapOrderLines                     `json:"order_lines" gorm:"foreignKey:Scrap_id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		ScrapLocationDetails     datatypes.JSON                        `json:"scrap_location_details" gorm:"type:json; default:'[]'; not null"`
		//DeliveryAddress_id uint                 `json:"deliveryaddress_id" gorm:"type:int"`
		//DeliveryAddress    mdm.Partner          `json:"deliveryaddress" gorm:"foreignkey:DeliveryAddress_id; references:ID"`
		PickupDateAndTime    SO_PickupDateAndTime  `json:"pickup_date_time" gorm:"embedded"`
		Scraping_status      string                `json:"scraping_status" gorm:"type:text"`
		Grn_status           *bool                 `json:"grn_status" gorm:"type:boolean"`
		GrnStatusId          *uint                 `json:"grn_status_id"  gorm:"type:int"`
		GrnStatus            model_core.Lookupcode `json:"Grn_status" gorm:"foreignKey:GrnStatusId; references:ID"`
		SourceDocuments      datatypes.JSON        `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
		SourceDocumentTypeId *uint                 `json:"source_document_type_id" gorm:"type:integer"`
		SourceDocumentType   model_core.Lookupcode `json:"source_document" gorm:"foreignKey:SourceDocumentTypeId; references:ID"`
		Status_id            *uint                 `json:"status_id" gorm:"type:int"`
		Status               model_core.Lookupcode `json:"status" gorm:"foreignKey:Status_id; references:ID"`
		Status_history       datatypes.JSON        `json:"status_history" gorm:"type:json; default:'[]'; not null"`
		ShippingOrderId      *uint                 `json:"shipping_order_id" gorm:"type:integer"`
		IsShipping           *bool                 `json:"is_shipping" gorm:"default:true"`
		ExpectedShippingDate string                `json:"expected_shipping_date"`
		Shipping_mode_id     *uint                 `json:"shipping_mode_id" gorm:"type:int"`
		Shipping_Mode        model_core.Lookupcode `json:"shipping_mode" gorm:"foreignKey:Shipping_mode_id; references:ID"`
		model_core.Model
	}

	// ScrapLocationDetails struct {
	// 	Receiver_name string        `json:"receiver_name" gorm:"type:varchar(100)"`
	// 	Email         string        `json:"email" gorm:"type:varchar(100)"`
	// 	Mobile_number int64         `json:"mobile_number" gorm:"type:bigint"`
	// 	Address_line  string        `json:"address_line" gorm:"varchar(100)"`
	// 	Address_line2 string        `json:"address_line2" gorm:"varchar(100)"`
	// 	Address_line3 string        `json:"address_line3" gorm:"varchar(100)"`
	// 	Landmark      string        `json:"landmark" gorm:"varchar(100)"`
	// 	Zip           string        `json:"zip" gorm:"varchar(60)"`
	// 	City          string        `json:"city" gorm:"varchar(50)"`
	// 	StateId       uint          `json:"state_id"`
	// 	State         model.State   `json:"state" gorm:"foreignkey:StateId; references:ID"`
	// 	CountryId     uint          `json:"country_id"`
	// 	Country       model.Country `json:"country" gorm:"foreignkey:CountryId; references:ID"`
	// }

	SO_PickupDateAndTime struct {
		PickupDate     string `json:"pickup_date" gorm:"type:date"`
		PickupFromTime string `json:"pickup_from_time" gorm:"type:varchar(50)"`
		PickupToTime   string `json:"pickup_to_time" gorm:"type:varchar(50)"`
	}

	ScrapOrderLines struct {
		Scrap_id    uint
		Product_id  uint               `json:"product_id"`
		Product     mdm.ProductVariant `json:"product_Details" gorm:"foreignkey:product_id; references:ID"`
		Lot_number  int                `json:"lot_number" gorm:"type:bigint"`
		InventoryId uint64             `json:"inventory_id"`
		//Uom        string             `json:"uom"`
		Uom_id              uint    `json:"uom_id"`
		UOM                 mdm.Uom `json:"uom" gorm:"foreignkey:Uom_id; references:ID"`
		Scrap_item_quantity int64   `json:"scrap_item_quantity" gorm:"type:bigint"`
		Price               float64 `json:"price" gorm:"type:double precision"`
		model_core.Model
	}
)
