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
type SalesOrders struct {
	model_core.Model
	SalesOrderNumber        string                `json:"sales_order_number" gorm:"type:text"`
	ReferenceNumber         string                `json:"reference_number" gorm:"type:text"`
	CurrencyId              *uint                 `json:"currency_id" gorm:"type:int"`
	Currency                model_core.Currency   `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	SoDate                  datatypes.Date        `json:"so_date" gorm:"type:date"`
	CustomerName            string                `json:"customer_name" gorm:"type:text"`
	CustomerShippingAddress datatypes.JSON        `json:"customer_shipping_address" gorm:"type:JSON"`
	CustomerBillingAddress  datatypes.JSON        `json:"customer_billing_address" gorm:"type:JSON"`
	PaymentTermsId          *uint                 `json:"payment_terms_id" gorm:"type:int"`
	PaymentTerms            model_core.Lookupcode `json:"payment_terms" gorm:"foreignkey:PaymentTermsId; references:ID"`
	PaymentDueDate          datatypes.Date        `json:"payment_due_date" gorm:"type:date"`
	VendorDetails           datatypes.JSON        `json:"vendor_details" gorm:"type:JSON"`
	ChannelName             string                `json:"channel_name" gorm:"type:text"`
	PaymentTypeId           *uint                 `json:"payment_type_id" gorm:"type:int"`
	PaymentType             model_core.Lookupcode `json:"payment_type" gorm:"foreignkey:PaymentTypeId; references:ID"`
	StatusId                uint                  `json:"status_id" gorm:"type:int"`
	Status                  model_core.Lookupcode `json:"status" gorm:"foreignkey:StatusId; references:ID"`
	StatusHistory           datatypes.JSON        `json:"status_history" gorm:"type:JSON; default:'[]';not null"`
	InvoicedId              *uint                 `json:"invoiced_id" gorm:"type:int"`
	Invoiced                model_core.Lookupcode `json:"invoiced" gorm:"foreignkey:InvoicedId; references:ID"`
	PaymentReceivedId       *uint                 `json:"payment_received_id" gorm:"type:int"`
	PaymentReceived         model_core.Lookupcode `json:"payment_received" gorm:"foreignkey:PaymentReceivedId; references:ID"`
	ExpectedShippingDate    datatypes.Date        `json:"expected_shipping_date" gorm:"type:date"`
	SalesOrderLines         []SalesOrderLines     `json:"sales_order_lines" gorm:"foreignkey:SoId; references:ID"`
	AdditionalInformation   AdditionalInformation `json:"additional_information" gorm:"embedded"`
	SoPaymentDetails        SoPaymentDetails      `json:"so_payment_details" gorm:"embedded"`
}

type SalesOrderLines struct {
	model_core.Model
	SoId              uint                                  `json:"so_id" gorm:"type:int"`
	ProductId         *uint                                 `json:"product_id" gorm:"type:int"`
	Product           mdm.ProductVariant                    `json:"product_details" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId *uint                                 `json:"product_template_id" gorm:"type:INT"`
	ProductTemplate   mdm.ProductTemplate                   `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	WarehouseId       *uint                                 `json:"warehouse_id" gorm:"type:INT"`
	Warehouse         shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseId; references:ID"`
	InventoryId       *uint                                 `json:"inventory_id" gorm:"type:INT"`
	Inventory         mdm.CentralizedBasicInventory         `json:"inventory" gorm:"foreignkey:InventoryId; references:ID"`
	UomId             *uint                                 `json:"uom_id" gorm:"type:INT"`
	Uom               mdm.Uom                               `json:"uom" gorm:"foreignkey:UomId; references:ID"`
	SerialNumber      string                                `json:"serial_number"  gorm:"type:text"`
	Quantity          int64                                 `json:"quantity" gorm:"type:int"`
	Price             float32                               `json:"price" gorm:"type:float"`
	Discount          float32                               `json:"discount" gorm:"type:float"`
	Amount            float32                               `json:"amount" gorm:"type:float"`
	Tax               float32                               `json:"tax" gorm:"type:float"`
	PaymentTermId     *uint                                 `json:"payment_term_id" gorm:"type:int"`
	PaymentTerms      model_core.Lookupcode                 `json:"payment_terms" gorm:"foreignkey:PaymentTermId; references:ID"`
	ExternalDetails   datatypes.JSON                        `json:"external_details" gorm:"column:external_details;type:json;default:'[]';not null"`
	Description       string                                `json:"description" gorm:"type:varchar(100)"`
	Location          string                                `json:"location" gorm:"type:varchar(100)"`
}

type SoPaymentDetails struct {
	AvailableCustomerCredits float32 `json:"available_customer_credits" gorm:"type:float"`
	UseCreditsForPayment     bool    `json:"use_credits_for_payment" gorm:"boolean"`
	SubTotal                 float32 `json:"sub_total" gorm:"type:float"`
	Tax                      float32 `json:"tax" gorm:"type:double precision"`
	ShippingCharges          float32 `json:"shipping_charges" gorm:"type:float"`
	AdjustmentAmount         float32 `json:"adjustment_amount" gorm:"type:float"`
	TotalAmount              float32 `json:"total_amount" gorm:"type:float"`
}

type AdditionalInformation struct {
	Notes              string         `json:"notes" gorm:"type:text"`
	TermsAndConditions string         `json:"terms_and_conditions" gorm:"type:text"`
	Attachments        datatypes.JSON `json:"attachments" gorm:"type:JSON; default:'[]';not null"`
}
