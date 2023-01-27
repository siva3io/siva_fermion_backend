package mdm

import (
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
type Vendors struct {
	model_core.Model
	Name             string         `json:"name" gorm:""`
	ContactId        *uint          `json:"contact_id" gorm:""`
	Contact          *Partner       `json:"contact" gorm:"foreignkey:ContactId; references:ID"`
	VendorDetails    datatypes.JSON `json:"vendor_details" gorm:"type:json"`
	PrimaryContactId *uint          `json:"primary_contact_id" gorm:""`
	PrimaryContact   *Partner       `json:"primary_contact" gorm:"foreignkey:PrimaryContactId; references:ID"`
	VendorDocs       datatypes.JSON `json:"vendor_docs" gorm:"type:json"`
}

// TODO : need to remove vendor_price_lists
type VendorPriceLists struct {
	model_core.Model
	Name             string              `json:"name" gorm:"unique;not null"`
	VendorID         uint                `json:"vendor_id" gorm:""`
	Vendors          Vendors             `json:"vendors" gorm:"foreignkey:VendorID; references:ID"`
	ProductVariantID uint                `json:"product_variant_id" gorm:""`
	ProductVariant   ProductVariant      `json:"product_variant" gorm:"foreignkey:ProductVariantID; references:ID"`
	MinQty           datatypes.JSON      `json:"min_qty" gorm:"type:json"`
	SupplyLeadTime   datatypes.JSON      `json:"supply_lead_time" gorm:"type:json"`
	Price            float32             `json:"price" gorm:""`
	CurrencyId       uint                `json:"currency_id"`
	Currency         model_core.Currency `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	PriceListId      uint                `json:"price_list_id" gorm:""`
	PriceList        datatypes.JSON      `json:"price_list" gorm:"type:json; default:'{}'"`
	// PriceList        Pricing               `json:"price_list" gorm:"foreignkey:PriceListId;references:ID"`
	ValidFrom     datatypes.Date        `json:"valid_from" gorm:"type:date"`
	ValidTo       datatypes.Date        `json:"valid_to" gorm:"type:date"`
	CreditPeriod  datatypes.JSON        `json:"credit_period" gorm:"type:json"`
	PaymentTypeId uint                  `json:"payment_type_id"`
	PaymentType   model_core.Lookupcode `json:"payment_type" gorm:"foreignkey:PaymentTypeId; references:ID"`
	//CreditPeriodUomId   uint             `json:"credit_period_uom_id"`
	//CreditPeriodUom     Uom              `json:"credit_period_uom" gorm:"foreignkey:CreditPeriodUomId; references:ID"`
	//MinQtyUom           Uom              `json:"min_qty_uom" gorm:"foreignkey:MinQtyUomId; references:ID"`
	//SupplyLeadTimeUomId uint             `json:"min_qty_uom_id"`
	//MinQtyUomId         uint             `json:"min_qty_uom_id"`
	//SupplyLeadTimeUom   Uom              `json:"min_qty_uom" gorm:"foreignkey:MinQtyUomId; references:ID"`
}
