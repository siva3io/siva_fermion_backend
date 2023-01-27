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
type PurchaseOrders struct {
	model_core.Model
	PurchaseOrderNumber   string                              `json:"purchase_order_number" gorm:"type:text"`
	ReferenceNumber       string                              `json:"reference_number" gorm:"type:text"`
	CurrencyId            *uint                               `json:"currency_id" gorm:"type:int"`
	Currency              *model_core.Currency                `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	DateAndTime           datatypes.Date                      `json:"date_and_time" gorm:"datetime"`
	ExpectedDeliveryDate  datatypes.Date                      `json:"expected_delivery_date" gorm:"type:date"`
	VendorDetails         datatypes.JSON                      `json:"vendor_details" gorm:"type:JSON"`
	DeliveryToId          *uint                               `json:"delivery_to_id" gorm:"type:INT"`
	DeliveryTo            model_core.Lookupcode               `json:"delivery_to" gorm:"foreignkey:DeliveryToId; references:ID"`
	OrganizationDetails   datatypes.JSON                      `json:"organization_details" gorm:"type:JSON"`
	StatusId              uint                                `json:"status_id" gorm:"type:int"`
	Status                model_core.Lookupcode               `json:"status" gorm:"foreignkey:StatusId; references:ID"`
	StatusHistory         datatypes.JSON                      `json:"status_history" gorm:"type:JSON; default:'[]'"`
	DeliveryAddress       datatypes.JSON                      `json:"delivery_address" gorm:"type:JSON"`
	BillingAddress        datatypes.JSON                      `json:"billing_address" gorm:"type:JSON"`
	TotalQuantity         int64                               `json:"total_quantity"`
	PaymentTermsId        *uint                               `json:"payment_terms_id" gorm:"type:int"`
	PaymentTerms          model_core.Lookupcode               `json:"payment_terms" gorm:"foreignkey:PaymentTermsId; references:ID"`
	PaymentDueDate        datatypes.Date                      `json:"payment_due_date" gorm:"type:date"`
	BilledId              *uint                               `json:"billed_id" gorm:"type:INT"`
	Billed                model_core.Lookupcode               `json:"billed" gorm:"foreignkey:BilledId; references:ID"`
	PaidId                *uint                               `json:"paid_id" gorm:"type:INT"`
	Paid                  model_core.Lookupcode               `json:"paid" gorm:"foreignkey:PaidId; references:ID"`
	PriceListID           *uint                               `json:"price_list_id" gorm:"type:INT"`
	PriceList             shared_pricing_and_location.Pricing `json:"priceList" gorm:"foreignkey:PriceListID; references:ID"`
	PurchaseOrderlines    []PurchaseOrderLines                `json:"purchase_order_lines" gorm:"foreignkey:PoId; references:ID"`
	AdditionalInformation AdditionalInformation               `json:"additional_information" gorm:"embedded"`
	PoPaymentDetails      PoPaymentDetails                    `json:"po_payment_details" gorm:"embedded"`
	SourceDocuments       datatypes.JSON                      `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	SourceDocumentTypeId  *uint                               `json:"source_document_type_id" gorm:"type:integer"`
	SourceDocumentType    model_core.Lookupcode               `json:"source_document" gorm:"foreignKey:SourceDocumentTypeId; references:ID"`
}
type PurchaseOrderLines struct {
	model_core.Model
	PoId                uint                                   `json:"po_id" gorm:"type:int"`
	ProductId           *uint                                  `json:"product_id" gorm:"type:int"`
	Product             *mdm.ProductVariant                    `json:"product_details" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId   *uint                                  `json:"product_template_id" gorm:"type:INT"`
	ProductTemplate     *mdm.ProductTemplate                   `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	WarehouseId         *uint                                  `json:"warehouse_id" gorm:"type:INT"`
	Warehouse           *shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseId; references:ID"`
	InventoryId         *uint                                  `json:"inventory_id" gorm:"type:INT"`
	Inventory           *mdm.CentralizedBasicInventory         `json:"inventory" gorm:"foreignkey:InventoryId; references:ID"`
	UomId               *uint                                  `json:"uom_id" gorm:"type:INT"`
	Uom                 *mdm.Uom                               `json:"uom" gorm:"foreignkey:UomId; references:ID"`
	SerialNumber        string                                 `json:"serial_number"  gorm:"type:text"`
	Quantity            int64                                  `json:"quantity" gorm:"type:int"`
	Price               float32                                `json:"price" gorm:"type:float"`
	Discount            float32                                `json:"discount" gorm:"type:float"`
	SalesPeriod         string                                 `json:"sales_period" gorm:"type:string"`
	CreditPeriod        string                                 `json:"credit_period" gorm:"type:string"`
	LeadTime            string                                 `json:"lead_time" gorm:"type:string"`
	ExpDeliveryLeadTime string                                 `json:"exp_delivery_lead_time" gorm:"type:string"`
	Tax                 float32                                `json:"tax" gorm:"type:float"`
	Amount              float32                                `json:"amount" gorm:"type:float"`
}
type PoPaymentDetails struct {
	AvailableVendorCredits float32 `json:"available_vendor_credits" gorm:"type:float"`
	UseCreditsForPayment   bool    `json:"use_credits_for_payment" gorm:"boolean"`
	SubTotal               float32 `json:"sub_total" gorm:"type:float"`
	Tax                    float32 `json:"tax" gorm:"type:double precision"`
	ShippingCharges        float32 `json:"shipping_charges" gorm:"type:float"`
	AdjustmentAmount       float32 `json:"adjustment_amount" gorm:"type:float"`
	TotalAmount            float32 `json:"total_amount" gorm:"type:float"`
}
