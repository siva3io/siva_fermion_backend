package returns

import (
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/orders"

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
type PurchaseReturns struct {
	model_core.Model
	PurchaseReturnNumber  string                       `json:"purchase_return_number" gorm:"type:text"`
	ReferenceNumber       string                       `json:"reference_number"  gorm:"type:text"`
	CurrencyId            *uint                        `json:"currency_id" gorm:"type:int"`
	Currency              *model_core.Currency         `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	PrDate                datatypes.Date               `json:"pr_date" gorm:"type:date"`
	ExpectedDeliveryDate  datatypes.Date               `json:"expected_delivery_date" gorm:"type:date"`
	SourceDocuments       datatypes.JSON               `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	Source_document_id    *uint                        `json:"source_document_id" gorm:"type:integer"`
	Source_document       model_core.Lookupcode        `json:"source_document" gorm:"foreignKey:Source_document_id; references:ID"`
	VendorDetails         datatypes.JSON               `json:"vendor_details" gorm:"type:JSON"`
	Amount                float32                      `json:"amount" gorm:"type:float"`
	PaymentTermsId        *uint                        `json:"payment_terms_id" gorm:"type:int"`
	PaymentTerms          model_core.Lookupcode        `json:"payment_terms" gorm:"foreignkey:PaymentTermsId; references:ID"`
	PaymentDueDate        datatypes.Date               `json:"payment_due_date" gorm:"type:date"`
	DebitNoteIssuedId     *uint                        `json:"debit_note_issued_id" gorm:"type:int"`
	DebitNoteIssued       model_core.Lookupcode        `json:"debit_note_issued" gorm:"foreignkey:DebitNoteIssuedId; references:ID"`
	StatusId              uint                         `json:"status_id" gorm:"type:int"`
	Status                model_core.Lookupcode        `json:"status" gorm:"foreignkey:StatusId; references:ID"`
	StatusHistory         datatypes.JSON               `json:"status_history" gorm:"type:JSON; default:'[]';not null"`
	PurchaseReturnLines   []PurchaseReturnLines        `json:"purchase_return_lines" gorm:"foreignkey:PrId; references:ID"`
	AdditionalInformation orders.AdditionalInformation `json:"additional_information" gorm:"embedded"`
	PickupDateAndTime     orders.PickupDateAndTime     `json:"pickup_date_and_time" gorm:"embedded"`
	PrPaymentDetails      PrPaymentDetails             `json:"pr_payment_details" gorm:"embedded"`
	TotalQuantity         int64                        `json:"total_quantity"`
	// AsnId uint                 `json:"asn_id" gorm:"type:int"`
	// Asn   inventory_orders.ASN `json:"asn" gorm:"foreignkey:AsnId; references:ID"`
	// GrnId uint                 `json:"grn_id" gorm:"type:int"`
	// Grn   inventory_orders.GRN `json:"grn" gorm:"foreignkey:GrnId; references:ID"`
}

type PurchaseReturnLines struct {
	model_core.Model
	PrId              uint                                   `json:"pr_id" gorm:"type:int"`
	ProductId         *uint                                  `json:"product_id" gorm:"type:int"`
	Product           *mdm.ProductVariant                    `json:"product_details" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId *uint                                  `json:"product_template_id" gorm:"type:INT"`
	ProductTemplate   *mdm.ProductTemplate                   `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	UomId             *uint                                  `json:"uom_id" gorm:"type:INT"`
	Uom               *mdm.Uom                               `json:"uom" gorm:"foreignkey:UomId; references:ID"`
	SerialNumber      string                                 `json:"serial_number"  gorm:"type:text"`
	InventoryId       uint64                                 `json:"inventory_id"`
	QuantityPurchased int                                    `json:"quantity_purchased" gorm:"type:int"`
	QuantityReturned  int64                                  `json:"quantity_returned" gorm:"type:int"`
	Rate              int                                    `json:"rate" gorm:"type:int"`
	LocationID        *uint                                  `json:"location_id" gorm:"type:int"`
	ReturnLocation    *shared_pricing_and_location.Locations `json:"return_location" gorm:"foreignkey:LocationID; references:ID"`
	Discount          float32                                `json:"discount" gorm:"type:float"`
	DiscountFormat    string                                 `json:"discount_format" gorm:"type:text"`
	Tax               float32                                `json:"tax" gorm:"type:float"`
	TaxFormat         string                                 `json:"tax_format" gorm:"type:text"`
	Amount            float32                                `json:"amount" gorm:"type:float"`
}

type PrPaymentDetails struct {
	VendorCredits        float32 `json:"vendor_Credits" gorm:"type:float"`
	UseCreditsForPayment bool    `json:"use_credits_for_payment" gorm:"type:boolean"`
	SubTotal             float32 `json:"sub_total" gorm:"type:float"`
	Tax                  float32 `json:"tax" gorm:"type:float"`
	ShippingCharges      float32 `json:"shipping_charges" gorm:"type:float"`
	Adjustments          float32 `json:"adjustments" gorm:"type:float"`
	TotalAmount          float32 `json:"total_amount" gorm:"type:float"`
}
