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
	DeliveryOrders struct {
		DeliveryOrdersDetails  DeliveryOrdersDetails                 `json:"delivery_order_details" gorm:"embedded"`
		DeliveryAddressDetails datatypes.JSON                        `json:"delivery_address_details" gorm:"type:json; default:'[]'; not null"`
		BillingAddressDetails  datatypes.JSON                        `json:"billing_address_details" gorm:"type:json; default:'[]'; not null"`
		DeliveryOrderLines     []DeliveryOrderLines                  `json:"delivery_order_lines,omitempty" gorm:"foreignkey:DO_id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Shipping_mode_id       uint                                  `json:"shipping_mode_id" gorm:"type:int"`
		Shipping_Mode          model_core.Lookupcode                 `json:"shipping_mode" gorm:"foreignKey:Shipping_mode_id; references:ID"`
		PickupDateAndTime      DO_PickupDateAndTime                  `json:"pickup_date_and_time" gorm:"embedded"`
		Additionalinformation  Additionalinformation                 `json:"additional_information" gorm:"embedded"`
		Payment_details        Payment_details                       `json:"payment_details" gorm:"embedded"`
		IsCreditUsed           *bool                                 `json:"is_credit_used" gorm:"type:boolean"`
		Status_id              uint                                  `json:"status_id" gorm:"type:int"`
		Status                 model_core.Lookupcode                 `json:"status" gorm:"foreignKey:Status_id; references:ID"`
		Status_history         datatypes.JSON                        `json:"status_history" gorm:"type:json; default:'[]'; not null"`
		ChannelName            string                                `json:"channel_name" gorm:"type:text"`
		Tracking_number        string                                `json:"tracking_number" gorm:"type:varchar(100)"`
		CustomerName           string                                `json:"customer_name" gorm:"type:varchar(100)"`
		Contact_id             *uint                                 `json:"contact_id" gorm:"integer"`
		Contact_details        mdm.Partner                           `json:"contact_details" gorm:"foreignKey:Contact_id; references:ID"`
		App_id                 uint                                  `json:"app_id" gorm:"type:int"`
		Warehouse_id           *int64                                `json:"warehouse_id" gorm:"int"`
		Warehouse              shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignKey:Warehouse_id; references:ID"`
		Payment_type_id        *uint                                 `json:"payment_type_id" gorm:"type:int"`
		Payment_Type           model_core.Lookupcode                 `json:"payment_type" gorm:"foreignKey:Payment_type_id; references:ID"`
		Shipping_Carrier_id    *uint                                 `json:"shipping_carrier_id" gorm:"type:int"`
		Shipping_Carrier       model_core.Lookupcode                 `json:"shipping_carrier" gorm:"foreignKey:Shipping_Carrier_id; references:ID"`
		Shipping_order_id      *uint                                 `json:"shipping_order_id" gorm:"type:integer"`
		Shipping_order         model_core.Lookupcode                 `json:"shipping_order" gorm:"foreignKey:Shipping_order_id; references:ID"`
		IsShipping             *bool                                 `json:"is_shipping" gorm:"default:true"`
		TotalQuantity          int64                                 `json:"total_quantity"`
		model_core.Model
	}

	DeliveryOrdersDetails struct {
		Reference_id          string                              `json:"reference_id" gorm:"type:varchar(50)"`
		Delivery_order_number string                              `json:"delivery_order_number" gorm:"type:varchar(100)"`
		PriceListID           *uint                               `json:"price_list_id" gorm:"type:INT"`
		PriceList             shared_pricing_and_location.Pricing `json:"pricing_details" gorm:"foreignkey:PriceListID; references:ID"`
		Do_currency           *uint                               `json:"do_currency"`
		Currency_details      *model_core.Currency                `json:"currency_details" gorm:"foreignKey:Do_currency; references:ID"`
		Payment_due_date      string                              `json:"payment_due_date" gorm:"type:date"`
		Payment_term_id       *uint                               `json:"payment_term_id" gorm:"type:integer"`
		Payment_Term          model_core.Lookupcode               `json:"payment_term" gorm:"foreignKey:Payment_term_id; references:ID"`
		Delivery_order_date   string                              `json:"delivery_order_date" gorm:"type:date"`
		ExpectedShippingDate  string                              `json:"expected_shipping_date" gorm:"type:date"`
		SourceDocuments       datatypes.JSON                      `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
		Source_document_id    *uint                               `json:"source_document_id" gorm:"type:integer"`
		Source_document       model_core.Lookupcode               `json:"source_document" gorm:"foreignKey:Source_document_id; references:ID"`
	}

	DeliveryOrderLines struct {
		DO_id               uint
		Product_id          *uint                 `json:"product_id"`
		Product             mdm.ProductVariant    `json:"product_details" gorm:"foreignkey:product_id; references:ID"`
		Product_template_id *uint                 `json:"product_template_id"`
		Product_template    mdm.ProductTemplate   `json:"product_template" gorm:"foreignkey:Product_template_id; references:ID"`
		Warehouse           string                `json:"warehouse" gorm:"type:varchar(100)"`
		Serial_no           string                `json:"serial_no" gorm:"type:varchar(100)"`
		Available_credit    uint                  `json:"available_credit" gorm:"type:int"`
		Uom_id              *uint                 `json:"uom_id"`
		UOM                 mdm.Uom               `json:"uom_details" gorm:"foreignkey:Uom_id; references:ID"`
		Rate                float64               `json:"rate" gorm:"type:double precision"`
		Product_qty         int64                 `json:"quantity" gorm:"type:int"`
		Discount            float64               `json:"discount_value" gorm:"type:double precision"`
		Tax                 float64               `json:"tax" gorm:"type:double precision"`
		Amount              float64               `json:"amount" gorm:"type:double precision"`
		Description         string                `json:"description" gorm:"type:varchar(100)"`
		Location            string                `json:"location" gorm:"type:varchar(100)"`
		InventoryId         uint                  `json:"inventory_id" gorm:"type:int"`
		Price               uint                  `json:"price"`
		PaymentTermsId      *uint                 `json:"payment_terms_id" gorm:"type:int"`
		PaymentTerms        model_core.Lookupcode `json:"payment_terms" gorm:"foreignkey:PaymentTermsId; references:ID"`
		model_core.Model
	}

	Additionalinformation struct {
		Notes              string         `json:"notes" gorm:"type:text"`
		Terms_and_condtion string         `json:"terms_and_condtion" gorm:"type:text"`
		Attachments        datatypes.JSON `json:"attachments" gorm:"type:json; default:'[]'; not null"`
	}

	Payment_details struct {
		Currency_id       *uint                `json:"currency_id"`
		Currency          *model_core.Currency `json:"currency_info" gorm:"foreignKey:Currency_id; references:ID"`
		Sub_total         float64              `gorm:"type:double precision" json:"sub_total"`
		Payment_tax       float64              `gorm:"type:double precision" json:"tax"`
		Shipping_charge   float64              `gorm:"type:double precision" json:"shipping_charge"`
		Vender_credits    float64              `gorm:"type:double precision" json:"vender_credits"`
		Adjustment_amount float64              `gorm:"type:double precision" json:"adjustment_amount"`
		Total_amount      float64              `gorm:"type:double precision" json:"total_amount"`
	}

	DO_PickupDateAndTime struct {
		PickupDate     string `json:"pickup_date" gorm:"type:date"`
		PickupFromTime string `json:"pickup_from_time" gorm:"type:varchar(50)"`
		PickupToTime   string `json:"pickup_to_time" gorm:"type:varchar(50)"`
	}
)
