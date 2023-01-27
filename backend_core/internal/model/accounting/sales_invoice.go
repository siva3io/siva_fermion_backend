package accounting

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/orders"

	"github.com/lib/pq"
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
type SalesInvoice struct {
	model_core.Model
	SalesInvoiceNumber       string                `json:"sales_invoice_number" gorm:"type:varchar(50)"`
	ReferenceNumber          string                `json:"reference_number" gorm:"type:varchar(50)"`
	ExpectedShipmentDate     datatypes.Date        `json:"expected_shipment_date" gorm:"type:date"`
	CurrencyID               *uint                 `json:"currency_id" gorm:"type:integer"`
	SalesInvoiceDate         time.Time             `json:"sales_invoice_date" gorm:"type:time"`
	Currency                 model_core.Currency   `json:"currency" gorm:"foreignKey:CurrencyID; references:ID"`
	SalesOrderIds            pq.Int64Array         `json:"sales_order_ids" gorm:"type:int[]"`
	LinkSalesOrders          datatypes.JSON        `json:"link_sales_orders" gorm:"type:json"`
	CustomerID               *uint                 `json:"customer_id" gorm:"type:integer"`
	Customer                 mdm.Partner           `json:"customer" gorm:"foreignKey:CustomerID; references:ID"`
	ChannelID                *uint                 `json:"channel_id" gorm:"type:integer"`
	Channel                  *omnichannel.Channel  `json:"channel" gorm:"foreignKey:ChannelID; references:ID"`
	PaymentTypeID            *uint                 `json:"payment_type_id" gorm:"type:integer"`
	PaymentType              model_core.Lookupcode `json:"payment_type" gorm:"foreignKey:PaymentTypeID; references:ID"`
	StatusID                 *uint                 `json:"status_id" gorm:"type:integer"`
	Status                   model_core.Lookupcode `json:"status" gorm:"foreignKey:StatusID; references:ID"`
	StatusHistory            datatypes.JSON        `json:"status_history" gorm:"type:json"`
	Invoiced                 bool                  `json:"is_invoiced" gorm:"type:boolean"`
	PaymentReceived          bool                  `json:"is_payment_received" gorm:"type:boolean"`
	BalanceDueAmount         float64               `json:"balance_due_amount" gorm:"type:double precision"`
	PaymentTermsID           *uint                 `json:"payment_terms_id" gorm:"type:integer"`
	PaymentTerms             model_core.Lookupcode `json:"payment_terms" gorm:"foreignKey:PaymentTermsID; references:ID"`
	AvailableCustomerCredits float64               `json:"available_customer_credits" gorm:"type:double precision"`
	PaymentDueDate           time.Time             `json:"payment_due_date" gorm:"type:time"`
	SalesInvoiceLines        []SalesInvoiceLines   `json:"sales_invoice_lines" gorm:"foreignKey:SalesInvoiceID;references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	InternalNotes            string                `json:"internal_notes" gorm:"type:varchar(250)"`
	ExternalNotes            string                `json:"external_notes" gorm:"type:varchar(250)"`
	TermsAndConditions       string                `json:"terms_and_conditions" gorm:"type:varchar(250)"`
	AttachmentFiles          datatypes.JSON        `json:"attachment_files" gorm:"type:json"`
	UseCreditsForPayment     bool                  `json:"use_credits_for_payment" gorm:"type:boolean"`
	SubTotalAmount           float64               `json:"sub_total_amount" gorm:"type:double precision"`
	ShippingAmount           float64               `json:"shipping_amount" gorm:"type:double precision"`
	TaxAmount                float64               `json:"tax_amount" gorm:"type:double precision"`
	IgstAmt                  float64               `json:"igst_amt" gorm:"type:double precision"`
	CgstAmt                  float64               `json:"cgst_amt" gorm:"type:double precision"`
	SgstAmt                  float64               `json:"sgst_amt" gorm:"type:double precision"`
	Adjustments              float64               `json:"adjustments" gorm:"type:double precision"`
	CustomerCreditsAmount    float64               `json:"customer_credits_amount" gorm:"type:double precision"`
	TotalAmount              float64               `json:"total_amount" gorm:"type: double precision"`
	BillingAddress           datatypes.JSON        `json:"billing_address" gorm:"type:json"`
	DeliveryAddress          datatypes.JSON        `json:"delivery_address" gorm:"type:json"`
	ShippingAddress          datatypes.JSON        `json:"shipping_address" gorm:"type:json"`
	SourceDocuments          datatypes.JSON        `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	SourceDocumentTypeId     *uint                 `json:"source_document_type_id" gorm:"type:integer"`
	SourceDocumentsType      model_core.Lookupcode `json:"source_document_type" gorm:"foreignKey:SourceDocumentTypeId; references:ID"`

	OrderId   *uint               `json:"order_id"`
	Order     *orders.SalesOrders `json:"order" gorm:"foreignKey:OrderId;references:ID"`
	OrderDate time.Time           `json:"order_date" gorm:"type:time"`
	Quantity  int32               `json:"quantity" gorm:"type:integer"`
}

type SalesInvoiceLines struct {
	model_core.Model
	SalesInvoiceID   uint                                  `json:"sales_invoice_id"`
	ProductID        *uint                                 `json:"product_id" gorm:"type:integer"`
	Product          mdm.ProductTemplate                   `json:"products" gorm:"foreignKey:ProductID; references:ID"`
	ProductVariantID *uint                                 `json:"product_variant_id" gorm:"type:integer"`
	ProductVariant   mdm.ProductVariant                    `json:"product_variant" gorm:"foreignKey:ProductVariantID; references:ID"`
	Description      string                                `json:"description" gorm:" type:varchar(250)"`
	WarehouseID      *uint                                 `json:"warehouse_id" gorm:"type:integer"`
	Warehouse        shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignKey:WarehouseID; references:ID"`
	InventoryID      *uint                                 `json:"inventory_id" gorm:"type:integer"`
	Inventory        mdm.CentralizedBasicInventory         `json:"inventory" gorm:"foreignKey:InventoryID; references:ID"`
	UOMID            *uint                                 `json:"uom_id" gorm:"type:integer"`
	UOM              mdm.Uom                               `json:"uom" gorm:"foreignKey:UOMID; references:ID"`
	Quantity         int32                                 `json:"quantity" gorm:"type:integer"`
	Discount         float64                               `json:"discount" gorm:"type:double precision"`
	DiscountTypeID   *uint                                 `json:"discount_type_id" gorm:"type:integer"`
	DiscountType     model_core.Lookupcode                 `json:"discount_type" gorm:"foreignKey:DiscountTypeID; references:ID"`
	Tax              float64                               `json:"tax" gorm:"type:double precision"`
	CGSTRate         float64                               `json:"cgst_rate" gorm:"type:double precision"`
	IGSTRate         float64                               `json:"igst_rate" gorm:"type:double precision"`
	SGSTRate         float64                               `json:"sgst_rate" gorm:"type:double precision"`
	Price            float64                               `json:"price" gorm:"type:double precision"`
	TaxTypeID        *uint                                 `json:"tax_type_id" gorm:"type:integer"`
	TaxType          model_core.Lookupcode                 `json:"tax_type" gorm:"foreignKey:TaxTypeID; references:ID"`
	PaymentTermsID   *uint                                 `json:"payment_terms_id" gorm:"type:integer"`
	PaymentTerms     model_core.Lookupcode                 `json:"payment_terms" gorm:"foreignKey:PaymentTermsID; references:ID"`
	TotalAmount      float64                               `json:"total_amount" gorm:"type:double precision"`
}
