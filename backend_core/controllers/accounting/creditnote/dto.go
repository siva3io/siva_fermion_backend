package creditnote

import (
	"time"

	"fermion/backend_core/controllers/accounting/purchase_invoice"
	cores "fermion/backend_core/controllers/cores"
	contacts "fermion/backend_core/controllers/mdm/contacts"
	locations "fermion/backend_core/controllers/mdm/locations"
	"fermion/backend_core/controllers/mdm/products"
	mdm_UOM "fermion/backend_core/controllers/mdm/uom"

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
type CreditNoteDTO struct {
	ID                   uint                   `json:"id" `
	IsEnabled            *bool                  `json:"is_enabled"`
	IsActive             *bool                  `json:"is_active"`
	CreatedDate          time.Time              `json:"created_date" `
	UpdatedDate          time.Time              `json:"updated_date" `
	AppId                *uint                  `json:"app_id" `
	CompanyId            uint                   `json:"company_id"`
	CreatedByID          *uint                  `json:"created_by"`
	UpdatedByID          *uint                  `json:"updated_by"`
	CustomerId           *uint                  `json:"customer_id"`
	PurchaseInvoiceId    *uint                  `json:"purchase_invoice_id"`
	CurrencyId           *uint                  `json:"currency_id"`
	GenerateCreditNoteId bool                   `json:"generate_credit_note_id"`
	CreditNoteID         string                 `json:"credit_note_id"`
	GenerateReferenceId  bool                   `json:"generate_reference_id"`
	StatusId             uint                   `json:"status_id"`
	ReferenceId          string                 `json:"reference_id"`
	ReasonId             *uint                  `json:"reason_id"`
	BillingAddressId     map[string]interface{} `json:"billing_address_id"`
	ShippingAddressId    map[string]interface{} `json:"shipping_address_id"`
	CreditNoteLineItems  []CreditNoteLineItems  `json:"credit_note_line_items" gorm:"foreignkey:CreditNoteId;references:ID"`
	InternalNotes        string                 `json:"internal_notes"`
	ExternalNotes        string                 `json:"external_notes"`
	TermsAndConditions   string                 `json:"terms_and_conditions"`
	Attachments          datatypes.JSON         `json:"attachments"`
	SubTotal             float64                `json:"sub_total"`
	Tax                  float64                `json:"tax"`
	ShippingCharges      float64                `json:"shipping_charges"`
	Adjustments          float64                `json:"adjustments"`
	CustomerCredits      float32                `json:"customer_credits"`
	TotalAmount          float64                `json:"total_amount"`
	SourceDocuments      map[string]interface{} `json:"source_documents"`
	SourceDocumentTypeId *uint                  `json:"source_document_type_id"`
}

type CreditNoteLineItems struct {
	ID                uint    `json:"id"`
	CreditNoteId      uint    `json:"credit_note_id"`
	ProductVariantId  uint    `json:"product_variant_id"`
	ProductTemplateId uint    `json:"product_template_id"`
	Quantity          int64   `json:"quantity"`
	UomId             uint    `json:"uom_id"`
	Price             float64 `json:"price"`
	Discount          float64 `json:"discount"`
	Tax               float64 `json:"tax"`
	Amount            float64 `json:"amount"`
	CompanyId         uint    `json:"company_id"`
}
type CreditNoteResponseDTO struct {
	AppId             *uint                                   `json:"app_id" `
	CompanyId         uint                                    `json:"company_id"`
	CustomerId        *uint                                   `json:"customer_id"`
	Customer          contacts.PartnerResponseDTO             `json:"customer"`
	PurchaseInvoiceId *uint                                   `json:"purchase_invoice_id"`
	PurchaseInvoice   purchase_invoice.ListPurchaseInvoiceDTO `json:"purchase_invoice"`
	CurrencyId        *uint                                   `json:"currency_id"`
	Currency          cores.CurrencyDTO                       `json:"currency" `
	CreditNoteID      string                                  `json:"credit_note_id"`
	locations.ModelDto
	ReferenceId          string                           `json:"reference_id"`
	ReasonId             *uint                            `json:"reason_id"`
	Reason               cores.LookupCodeDTO              `json:"reason"`
	StatusId             uint                             `json:"status_id"`
	Status               cores.LookupCodeDTO              `json:"status"`
	BillingAddressId     map[string]interface{}           `json:"billing_address_id"`
	ShippingAddressId    map[string]interface{}           `json:"shipping_address_id"`
	CreditNoteLineItems  []CreditNoteLineItemsResponseDTO `json:"credit_note_line_items"`
	InternalNotes        string                           `json:"internal_notes" `
	ExternalNotes        string                           `json:"external_notes" `
	TermsAndConditions   string                           `json:"terms_and_conditions" `
	Attachments          datatypes.JSON                   `json:"attachments"`
	SubTotal             float64                          `json:"sub_total"`
	Tax                  float64                          `json:"tax"`
	ShippingCharges      float64                          `json:"shipping_charges" `
	Adjustments          float64                          `json:"adjustments"`
	CustomerCredits      float32                          `json:"customer_credits" `
	TotalAmount          float64                          `json:"total_amount"`
	SourceDocuments      map[string]interface{}           `json:"source_documents"`
	SourceDocumentTypeId *uint                            `json:"source_document_type_id"`
	SourceDocumentType   cores.LookupCodeDTO              `json:"source_document"`
}
type CreditNoteLineItemsResponseDTO struct {
	CreditNoteId      uint                        `json:"credit_note_id"`
	ProductVariantId  uint                        `json:"product_variant_id"`
	ProductVariant    products.VariantResponseDTO `json:"product_variant" `
	ProductTemplateId uint                        `json:"product_template_id"`
	ProductTemplate   products.TemplateReponseDTO `json:"product_template" `
	Quantity          int64                       `json:"quantity"`
	UomId             uint                        `json:"uom_id"`
	Uom               mdm_UOM.UomResponseDTO      `json:"uom" `
	Discount          float64                     `json:"discount"`
	Tax               float64                     `json:"tax"`
	Amount            float64                     `json:"amount"`
	CompanyId         uint                        `json:"company_id"`
	Price             float64                     `json:"price"`
}
