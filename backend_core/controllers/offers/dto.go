package offers

import (
	"time"

	"fermion/backend_core/controllers/cores"
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
type ListOffersDTO struct {
	ID                   *uint                 `json:"id"`
	PromotionalDetails   PromotionalDetails    `json:"promotional_details"`
	OfferProductDetails  []OfferProductDetails `json:"offer_product_details"`
	OtherDetails         OtherDetails          `json:"other_details"`
	DiscountTypeID       *uint                 `json:"discount_type_id"`
	DiscountType         cores.LookupCodesDTO  `json:"discount_type"`
	DiscountValue        float64               `json:"discount_value"`
	PromotionalStatus    bool                  `json:"promotional_status"`
	IsAppliedForAll      bool                  `json:"is_applied_for_all"`
	IsSameDiscountForAll bool                  `json:"is_same_discount_for_all"`
}
type OffersDTO struct {
	ID                   *uint                 `json:"id"`
	IsEnabled            bool                  `json:"is_enabled"`
	IsActive             bool                  `json:"is_active"`
	CreatedByID          *uint                 `json:"created_by"`
	UpdatedByID          *uint                 `json:"updated_by"`
	CompanyId            *uint                 `json:"company_id"`
	PromotionalDetails   PromotionalDetails    `json:"promotional_details"`
	OfferProductDetails  []OfferProductDetails `json:"offer_product_details"`
	OtherDetails         OtherDetails          `json:"other_details"`
	DiscountTypeID       uint                  `json:"discount_type_id"`
	DiscountType         cores.LookupCodesDTO  `json:"discount_type"`
	DiscountValue        float64               `json:"discount_value"`
	PromotionalStatus    bool                  `json:"promotional_status"`
	IsAppliedForAll      bool                  `json:"is_applied_for_all"`
	IsSameDiscountForAll bool                  `json:"is_same_discount_for_all"`
}
type PromotionalDetails struct {
	Name             string    `json:"name"`
	StartDateAndTime time.Time `json:"start_date_and_time"`
	EndDateAndTime   time.Time `json:"end_date_and_time"`
}

type OfferProductDetails struct {
	ID                *uint                    `json:"id"`
	OfferId           *uint                    `json:"offer_id"`
	ProductId         *uint                    `json:"product_id"`
	Product           offerResponseVarientDTO  `json:"product_details"`
	ProductTemplateId *uint                    `json:"product_template_id"`
	ProductTemplate   offerResponseTemplateDTO `json:"product_template"`
	MRP               float64                  `json:"mrp"`
	SaleRate          float64                  `json:"sale_rate"`
	DiscountTypeID    *uint                    `json:"discount_type_id"`
	DiscountType      cores.LookupCodesDTO     `json:"discount_type"`
	DiscountedPrice   float64                  `json:"discounted_price"`
}
type OtherDetails struct {
	TermsAndConditionsID *uint                `json:"terms_and_conditions_id"`
	TermsAndConditions   cores.LookupCodesDTO `json:"terms_and_conditions"`
	Descriptions         string               `json:"descriptions"`
	TermsAndCondition    string               `json:"terms_and_condition"`
}

type offerResponseVarientDTO struct {
	SkuId             string `json:"sku_id,omitempty"`
	ProductName       string `json:"product_name"`
	ID                *uint  `json:"id"`
	ProductTemplateId *uint  `json:"product_template_id"`
}

type offerResponseTemplateDTO struct {
	Name string `json:"product_name"`
}
