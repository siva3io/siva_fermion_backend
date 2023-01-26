package accounting

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
type DebitNote struct {
	model_core.Model
	VendorId               *uint                 `json:"vendor_id"`
	Vendor                 mdm.Vendors           `json:"vendor" gorm:"foreignkey:VendorId;references:ID"`
	PurchaseInvoiceId      *uint                 `json:"purchase_invoice_id"`
	PurchaseInvoice        PurchaseInvoice       `json:"purchase_invoice" gorm:"foreignkey:PurchaseInvoiceId;references:ID"`
	CurrencyId             *uint                 `json:"currency_id"`
	Currency               model_core.Currency   `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	DebitNoteID            string                `json:"debit_note_id"`
	ReferenceId            string                `json:"reference_id"`
	ReasonId               *uint                 `json:"reason_id"`
	Reason                 model_core.Lookupcode `json:"reason" gorm:"foreignkey:ReasonId; references:ID"`
	StatusId               uint                  `json:"status_id" gorm:"type:int"`
	Status                 model_core.Lookupcode `json:"status" gorm:"foreignkey:StatusId; references:ID"`
	BillingAddressId       *uint                 `json:"billing_address_id"`
	BillingAddress         mdm.Partner           `json:"billing_address" gorm:"foreignkey:BillingAddressId; references:ID"`
	DeliveryAddressId      *uint                 `json:"delivery_address_id"`
	DeliveryAddress        mdm.Partner           `json:"delivery_address" gorm:"foreignkey:DeliveryAddressId;references:ID"`
	AvailableVendorCredits float64               `json:"available_vendor_credits"`
	DebitNoteLineItems     []DebitNoteLineItems  `json:"debit_note_line_items" gorm:"foreignkey:DebitNoteId;references:ID"`
	InternalNotes          string                `json:"internal_notes" gorm:""`
	ExternalNotes          string                `json:"external_notes" gorm:""`
	TermsAndConditions     string                `json:"terms_and_conditions" gorm:""`
	Attachments            datatypes.JSON        `json:"attachments"`
	SubTotal               float64               `json:"sub_total" gorm:""`
	Tax                    float64               `json:"tax" gorm:""`
	ShippingCharges        float64               `json:"shipping_charges" gorm:""`
	Adjustments            float64               `json:"adjustments" gorm:""`
	CustomerCredits        float32               `json:"customer_credits" gorm:""`
	TotalAmount            float64               `json:"total_amount" gorm:""`

	SourceDocuments      datatypes.JSON        `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	SourceDocumentTypeId *uint                 `json:"source_document_type_id" gorm:"type:integer"`
	SourceDocumentType   model_core.Lookupcode `json:"source_document" gorm:"foreignKey:SourceDocumentTypeId; references:ID"`
}

type DebitNoteLineItems struct {
	model_core.Model
	DebitNoteId       uint                `json:"debit_note_id"`
	ProductVariantId  uint                `json:"product_variant_id"`
	ProductVariant    mdm.ProductVariant  `json:"product_variant" gorm:"foreignkey:ProductVariantId;references:ID"`
	ProductTemplateId uint                `json:"product_template_id"`
	ProductTemplate   mdm.ProductTemplate `json:"product_template" gorm:"foreignkey:ProductTemplateId;references:ID"`
	Quantity          int64               `json:"quantity"`
	UomId             uint                `json:"uom_id"`
	Uom               mdm.Uom             `json:"uom" gorm:"foreignkey:UomId;references:ID"`
	Price             float64             `json:"price"`
	Discount          float64             `json:"discount"`
	Tax               float64             `json:"tax"`
	Amount            float64             `json:"amount"`
}
