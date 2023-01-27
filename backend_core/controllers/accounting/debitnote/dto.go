package debitnote

import (
	purchaseinvoice_dto "fermion/backend_core/controllers/accounting/purchase_invoice"
	core_dto "fermion/backend_core/controllers/cores"
	partner_dto "fermion/backend_core/controllers/mdm/contacts"
	locations "fermion/backend_core/controllers/mdm/locations"
	products_dto "fermion/backend_core/controllers/mdm/products"
	uom_dto "fermion/backend_core/controllers/mdm/uom"
	vendors_dto "fermion/backend_core/controllers/mdm/vendors"

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
type DebitNoteDTO struct {
	CreatedByID            *uint                  `json:"created_by" gorm:"column:created_by"`
	VendorId               uint                   `json:"vendor_id"`
	PurchaseInvoiceId      *uint                  `json:"purchase_invoice_id"`
	CurrencyId             *uint                  `json:"currency_id"`
	GenerateDebitNoteId    bool                   `json:"generate_debit_note_id"`
	DebitNoteID            string                 `json:"debit_note_id"`
	CompanyId              uint                   `json:"company_id"`
	GenerateReferenceId    bool                   `json:"generate_reference_id"`
	ReferenceId            string                 `json:"reference_id"`
	ReasonId               *uint                  `json:"reason_id"`
	BillingAddressId       *uint                  `json:"billing_address_id"`
	DeliveryAddressId      *uint                  `json:"delivery_address_id"`
	StatusId               uint                   `json:"status_id"`
	AvailableVendorCredits float64                `json:"available_vendor_credits"`
	DebitNoteLineItems     []DebitNoteLineItems   `json:"debit_note_line_items"`
	InternalNotes          string                 `json:"internal_notes"`
	ExternalNotes          string                 `json:"external_notes"`
	TermsAndConditions     string                 `json:"terms_and_conditions"`
	Attachments            datatypes.JSON         `json:"attachments"`
	SubTotal               float64                `json:"sub_total"`
	Tax                    float64                `json:"tax"`
	ShippingCharges        float64                `json:"shipping_charges"`
	Adjustments            float64                `json:"adjustments"`
	CustomerCredits        float32                `json:"customer_credits"`
	TotalAmount            float64                `json:"total_amount"`
	SourceDocuments        map[string]interface{} `json:"source_documents"`
	SourceDocumentTypeId   *uint                  `json:"source_document_type_id"`
	ID                     uint                   `json:"id"`
	UpdatedByID            *uint                  `json:"updated_by"`
}

type DebitNoteLineItems struct {
	ProductVariantId  uint    `json:"product_variant_id"`
	ProductTemplateId uint    `json:"product_template_id"`
	Quantity          int64   `json:"quantity"`
	UomId             uint    `json:"uom_id"`
	Price             float64 `json:"price"`
	Discount          float64 `json:"discount"`
	Tax               float64 `json:"tax"`
	Amount            float64 `json:"amount"`
	ID                uint    `json:"id"`
	DebitNoteId       uint    `json:"debit_note_id"`
}
type DebitNoteResponseDTO struct {
	locations.ModelDto
	VendorId               *uint                                      `json:"vendor_id"`
	Vendor                 vendors_dto.VendorsRequestDTO              `json:"vendor"`
	PurchaseInvoiceId      *uint                                      `json:"purchase_invoice_id"`
	PurchaseInvoice        purchaseinvoice_dto.ListPurchaseInvoiceDTO `json:"purchase_invoice"`
	CurrencyId             *uint                                      `json:"currency_id"`
	Currency               core_dto.CurrencyDTO                       `json:"currency"`
	DebitNoteID            string                                     `json:"debit_note_id"`
	ReferenceId            string                                     `json:"reference_id"`
	ReasonId               *uint                                      `json:"reason_id"`
	Reason                 core_dto.LookupCodesDTO                    `json:"reason"`
	StatusId               uint                                       `json:"status_id"`
	CompanyId              uint                                       `json:"company_id"`
	Status                 core_dto.LookupCodesDTO                    `json:"status"`
	BillingAddressId       *uint                                      `json:"billing_address_id"`
	BillingAddress         partner_dto.PartnerResponseDTO             `json:"billing_address"`
	DeliveryAddressId      *uint                                      `json:"delivery_address_id"`
	DeliveryAddress        partner_dto.PartnerResponseDTO             `json:"delivery_address"`
	AvailableVendorCredits float64                                    `json:"available_vendor_credits"`
	DebitNoteLineItems     []DebitNoteLineItem                        `json:"debit_note_line_items"`
	InternalNotes          string                                     `json:"internal_notes"`
	ExternalNotes          string                                     `json:"external_notes"`
	TermsAndConditions     string                                     `json:"terms_and_conditions"`
	Attachments            map[string]interface{}                     `json:"attachments"`
	SubTotal               float64                                    `json:"sub_total"`
	Tax                    float64                                    `json:"tax"`
	ShippingCharges        float64                                    `json:"shipping_charges"`
	Adjustments            float64                                    `json:"adjustments"`
	CustomerCredits        float32                                    `json:"customer_credits"`
	TotalAmount            float64                                    `json:"total_amount"`
	SourceDocuments        map[string]interface{}                     `json:"source_documents"`
	SourceDocumentTypeId   *uint                                      `json:"source_document_type_id"`
	SourceDocumentType     core_dto.LookupCodesDTO                    `json:"source_document"`
	ID                     uint                                       `json:"id"`
	CreatedByID            *uint                                      `json:"created_by"`
	UpdatedByID            *uint                                      `json:"updated_by"`
}
type DebitNoteLineItem struct {
	locations.ModelDto
	DebitNoteId       uint                            `json:"debit_note_id"`
	ProductVariantId  uint                            `json:"product_variant_id"`
	ProductVariant    products_dto.VariantResponseDTO `json:"product_variant"`
	ProductTemplateId uint                            `json:"product_template_id"`
	ProductTemplate   products_dto.VariantResponseDTO `json:"product_template"`
	Quantity          int64                           `json:"quantity"`
	UomId             uint                            `json:"uom_id"`
	Uom               uom_dto.UomResponseDTO          `json:"uom"`
	Price             float64                         `json:"price"`
	Discount          float64                         `json:"discount"`
	Tax               float64                         `json:"tax"`
	Amount            float64                         `json:"amount"`
}
