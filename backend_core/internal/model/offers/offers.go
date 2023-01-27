package offers

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
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
type Offers struct {
	model_core.Model
	PromotionalDetails   PromotionalDetails    `json:"promotional_details" gorm:"embedded"`
	OfferProductDetails  []OfferProductDetails `json:"offer_product_details" gorm:"foreignkey:OfferId; references:ID"`
	OtherDetails         OtherDetails          `json:"other_details" gorm:"embedded"`
	DiscountTypeID       *uint                 `json:"discount_type_id" gorm:"type:integer"`
	DiscountType         model_core.Lookupcode `json:"discount_type" gorm:"foreignKey:DiscountTypeID; references:ID"`
	DiscountValue        float64               `json:"discount_value" gorm:"float"`
	PromotionalStatus    bool                  `json:"promotional_status" gorm:"boolean"`
	IsAppliedForAll      bool                  `json:"is_applied_for_all" gorm:"boolean"`
	IsSameDiscountForAll bool                  `json:"is_same_discount_for_all" gorm:"boolean"`
}
type PromotionalDetails struct {
	Name             string    `json:"name" gorm:"type:text"`
	StartDateAndTime time.Time `json:"start_date_and_time" gorm:"type:date"`
	EndDateAndTime   time.Time `json:"end_date_and_time" gorm:"type:date"`
}
type OfferProductDetails struct {
	model_core.Model
	OfferId           uint                  `json:"offer_id" gorm:"type:int"`
	ProductId         *uint                 `json:"product_id" gorm:"type:int"`
	Product           *mdm.ProductVariant   `json:"product_details" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId *uint                 `json:"product_template_id" gorm:"type:INT"`
	ProductTemplate   *mdm.ProductTemplate  `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	MRP               float64               `json:"mrp" gorm:"type:float"`
	SaleRate          float64               `json:"sale_rate" gorm:"type:float"`
	DiscountTypeID    *uint                 `json:"discount_type_id" gorm:"type:integer"`
	DiscountType      model_core.Lookupcode `json:"discount_type" gorm:"foreignKey:DiscountTypeID; references:ID"`
	DiscountedPrice   float64               `json:"discounted_price" gorm:"type:float"`
}
type OtherDetails struct {
	TermsAndConditionsID *uint                 `json:"terms_and_conditions_id" gorm:"type:integer"`
	TermsAndConditions   model_core.Lookupcode `json:"terms_and_conditions" gorm:"foreignKey:TermsAndConditionsID; references:ID"`
	Descriptions         string                `json:"descriptions" gorm:"type:text"`
	TermsAndCondition    string                `json:"terms_and_condition" gorm:"type:text"`
}
