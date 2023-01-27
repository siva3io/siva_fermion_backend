package accounting

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
type PurchaseInvoice struct {
	model_core.Model
	PurchaseInvoiceNumber  string                       `json:"purchase_invoice_number" gorm:"type:text"`
	PurchaseInvoiceDate    datatypes.Date               `json:"purchase_invoice_date" gorm:"type:date"`
	ReferenceNumber        string                       `json:"reference_number" gorm:"type:text"`
	VendorDetails          datatypes.JSON               `json:"vendor_details" gorm:"type:JSON; default:'[]'"`
	StatusId               uint                         `json:"status_id" gorm:"type:int"`
	Status                 model_core.Lookupcode        `json:"status" gorm:"foreignkey:StatusId; references:ID"`
	StatusHistory          datatypes.JSON               `json:"status_history" gorm:"type:JSON; default:'[]'"`
	PaidId                 uint                         `json:"paid_id" gorm:"type:int"`
	Paid                   model_core.Lookupcode        `json:"paid" gorm:"foreignkey:PaidId; references:ID"`
	PaymentTermsId         uint                         `json:"payment_terms_id" gorm:"type:int"`
	PaymentTerms           model_core.Lookupcode        `json:"payment_terms" gorm:"foreignkey:PaymentTermsId; references:ID"`
	PaymentDueDate         datatypes.Date               `json:"payment_due_date" gorm:"type:date"`
	PaymentAmount          float32                      `json:"payment_amount" gorm:"type:float"`
	BalanceDue             float32                      `json:"balance_due" gorm:"type:float"`
	DueDate                datatypes.Date               `json:"due_date" gorm:"type:date"`
	ExpectedDeliveryDate   datatypes.Date               `json:"expected_delivery_date" gorm:"type:date"`
	LinkSourceDocumentType *uint                        `json:"link_source_document_type" gorm:"type:integer"`
	Source_document        model_core.Lookupcode        `json:"link_source_document_id" gorm:"foreignKey:LinkSourceDocumentType; references:ID"`
	LinkSourceDocument     datatypes.JSON               `json:"link_source_document" gorm:"type:json"`
	CurrencyId             *uint                        `json:"currency_id" gorm:"type:int"`
	Currency               *model_core.Currency         `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	DeliveryAddress        datatypes.JSON               `json:"delivery_address" gorm:"type:JSON; default:'[]'"`
	AdditionalInformation  orders.AdditionalInformation `json:"additional_information" gorm:"embedded"`
	PaymentDetails         PiPaymentDetails             `json:"payment_details" gorm:"embedded"`
	PurchaseInvoiceLines   []PurchaseInvoiceLines       `json:"purchase_invoice_lines" gorm:"foreignkey:PurchaseInvoiceId; references:ID"`
	PoIds                  datatypes.JSON               `json:"po_ids"  gorm:"type:JSON; default:'[]'"`
}

type PurchaseInvoiceLines struct {
	model_core.Model
	PurchaseInvoiceId uint                                   `json:"purchase_invoice_id" gorm:"type:int"`
	ProductId         uint                                   `json:"product_id" gorm:"type:int"`
	Product           *mdm.ProductVariant                    `json:"product_details" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId uint                                   `json:"product_template_id" gorm:"type:int"`
	ProductTemplate   *mdm.ProductTemplate                   `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	WarehouseId       uint                                   `json:"warehouse_id" gorm:"type:int"`
	Warehouse         *shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseId; references:ID"`
	InventoryId       uint                                   `json:"inventory_id" gorm:"type:int"`
	Inventory         *mdm.CentralizedBasicInventory         `json:"inventory" gorm:"foreignkey:InventoryId; references:ID"`
	UomId             uint                                   `json:"uom_id" gorm:"type:int"`
	Uom               mdm.Uom                                `json:"uom" gorm:"foreignkey:UomId; references:ID"`
	PaymentTermsId    uint                                   `json:"payment_terms_id" gorm:"type:int"`
	PaymentTerms      model_core.Lookupcode                  `json:"payment_terms" gorm:"foreignkey:PaymentTermsId; references:ID"`
	SerialNumber      string                                 `json:"serial_number"  gorm:"type:text"`
	Quantity          int                                    `json:"quantity" gorm:"type:int"`
	Price             float32                                `json:"price" gorm:"type:float"`
	Discount          float32                                `json:"discount" gorm:"type:float"`
	Tax               float32                                `json:"tax" gorm:"type:float"`
	Amount            float32                                `json:"amount" gorm:"type:float"`
}
type PiPaymentDetails struct {
	AvailableVendorCredits float32 `json:"available_vendor_credits" gorm:"type:float"`
	UseCreditsForPayment   bool    `json:"use_credits_for_payment" gorm:"boolean"`
	SubTotal               float32 `json:"sub_total" gorm:"type:float"`
	Tax                    float32 `json:"tax" gorm:"type:double precision"`
	ShippingCharges        float32 `json:"shipping_charges" gorm:"type:float"`
	AdjustmentAmount       float32 `json:"adjustment_amount" gorm:"type:float"`
	TotalAmount            float32 `json:"total_amount" gorm:"type:float"`
}
