package shipping

import (
	"time"

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
type ShippingPartner struct {
	model_core.Model
	PartnerName           string                `json:"partner_name" gorm:"type:varchar(50); unique; not null;"`
	ShippingPartnerTypeId uint                  `json:"shipping_partner_type_id" gorm:"type:int"`
	ShippingPartnerType   model_core.Lookupcode `json:"shipping_partner_type" gorm:"foreignKey:ShippingPartnerTypeId;references:ID"`
	ProfileOptions        datatypes.JSON        `json:"profile_options" gorm:"type:json"`
	AuthOptions           datatypes.JSON        `json:"auth_options" gorm:"type:json:"`
	SubscriptionOptions   datatypes.JSON        `json:"subscription_options" gorm:"type:json"`
	SpDefaultPreferences  datatypes.JSON        `json:"sp_default_preferences" gorm:"type:json"`
	SpTimePreferences     datatypes.JSON        `json:"sp_time_preferences" gorm:"type:json"`
	SpContactDetails      datatypes.JSON        `json:"sp_contact_details" gorm:"type:json"`
}

type UserShippingPartnerRegistration struct {
	model_core.Model
	ShippingPartnerId        *uint                    `json:"shipping_partner_id"`
	ShippingPartner          ShippingPartner          `json:"shipping_partner" gorm:"foreignKey:ShippingPartnerId;references:ID"`
	UserId                   uint                     `json:"user_id"`
	User                     model_core.CoreUsers     `json:"user" gorm:"foreignKey:UserId;references:ID"`
	AccountDetails           AccountDetails           `gorm:"embedded" json:"account_details"`
	PickupDetails            PickupDetails            `gorm:"embedded" json:"pickup_details"`
	PackingDetails           PackingDetails           `gorm:"embedded" json:"packing_details"`
	AddAddressLocations      datatypes.JSON           `json:"add_locations" gorm:"type:json"`
	CommericalInvoiceDetails CommericalInvoiceDetails `gorm:"embedded" json:"commerical_invoice_details"`
	UploadDoc                UploadDoc                `gorm:"embedded" json:"upload_doc"`
	TestAccount              *bool                    `json:"test_account" gorm:"type:boolean"`
	StatusId                 uint                     `json:"status_id"`
	Status                   model_core.Lookupcode    `json:"status" gorm:"foreignKey:StatusId;references:ID;"`
	StatusHistory            datatypes.JSON           `json:"status_history" gorm:"type:json"`
}

type AccountDetails struct {
	Account_number       int64                 `json:"account_number" gorm:"type:bigint"`
	Web_service_key      string                `json:"web_service_key" gorm:"type:varchar(50)"`
	Web_service_password string                `json:"web_service_password" gorm:"type:varchar(50)"`
	Meter_number         string                `json:"meter_number" gorm:"type:varchar"`
	Carrier_weight_id    *uint                 `json:"carrier_weight_id" gorm:"type:int"`
	Carrier_weight       mdm.Uom               `json:"carrier_weight" gorm:"foreignKey:Carrier_weight_id; references:ID"`
	Carrier_currency_id  *uint                 `json:"carrier_currency_id" gorm:"type:integer"`
	Currency             model_core.Currency   `json:"currency" gorm:"foreignKey:Carrier_currency_id; references:ID"`
	Show_service_rates   *bool                 `json:"show_service_rates" gorm:"type:boolean"`
	Payment_type_id      *uint                 `json:"payment_type_id" gorm:"type:int"`
	Payment_Type         model_core.Lookupcode `json:"payment_type" gorm:"foreignKey:Payment_type_id; references:ID"`
}

type PickupDetails struct {
	Shipper_details_id *uint                                 `json:"shipper_details_id" gorm:"type:integer"`
	Shipper_details    shared_pricing_and_location.Locations `json:"shipper_details" gorm:"foreignKey:Shipper_details_id; references:ID"`
	From_time          time.Time                             `json:"from_time" gorm:"type:time"`
	To_time            time.Time                             `json:"to_time" gorm:"type:time"`
}

type PackingDetails struct {
	Use_volumetric_weight *bool                 `json:"use_volumetric_weight" gorm:"type:boolean"`
	Packing_method_id     *uint                 `json:"packing_method_id" gorm:"type:integer"`
	Packing_method        model_core.Lookupcode `json:"packing_method" gorm:"foreignKey:Packing_method_id; references:ID"`
	Maximum_weight        float64               `json:"maximum_weight" gorm:"type:double precision"`
}

type CommericalInvoiceDetails struct {
	IsCommericalInvoiceDetails *bool                 `json:"is_commercial_invoice_details" gorm:"type:boolean"`
	Include_shipping_charges   *bool                 `json:"include_shipping_charges" gorm:"type:boolean"`
	Include_insurance_charges  *bool                 `json:"include_insurance_charges" gorm:"type:boolean"`
	Customer_declaration       string                `json:"customer_declaration" gorm:"type:varchar"`
	Maximum_weight_allowed     float64               `json:"maximum_weight_allowed" gorm:"type:double precision"`
	Terms_of_sales_id          *uint                 `json:"terms_of_sales_id" gorm:"type:integer"`
	Terms_of_sales             model_core.Lookupcode `json:"terms_of_sales" gorm:"foreignKey:Terms_of_sales_id; references:ID"`
	Uom_id                     *uint                 `json:"uom_id" gorm:"type:integer"`
	UOM                        mdm.Uom               `json:"uom" gorm:"foreignkey:Uom_id; references:ID"`
	Invoice_currency_id        *uint                 `json:"invoice_currency_id" gorm:"type:integer"`
	Currency_info              model_core.Currency   `json:"currency" gorm:"foreignKey:Invoice_currency_id; references:ID"`
	ProductId                  *uint                 `json:"product_id" gorm:"type:int"`
	ProductVariant             mdm.ProductVariant    `json:"product_variant" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId          *uint                 `json:"product_template_id" gorm:"type:integer"`
	ProductTemplate            mdm.ProductTemplate   `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
}

type UploadDoc struct {
	IsUploadDoc         *bool          `json:"is_upload_doc" gorm:"type:boolean"`
	Customer_signature  datatypes.JSON `json:"customer_signature" gorm:"json"`
	Company_letter_head datatypes.JSON `json:"company_letter_head" gorm:"json"`
	Additional_document datatypes.JSON `json:"additional_document" gorm:"json"`
	// Company_signature_id   uint        `json:"company_signature_id"`
	// Attachments            Attachments `json:"-" gorm:"foreignKey:Company_signature_id; references:ID"`
	// Company_letter_head_id uint        `json:"company_letter_head_id"`
	// Company_Attachments    Attachments `json:"-" gorm:"foreignKey:Company_letter_head_id; references:ID"`
	// Additional_document_id uint        `json:"additional_document_id"`
	// Additional_Attachments Attachments `json:"-" gorm:"foreignKey:Additional_document_id; references:ID"`
}
