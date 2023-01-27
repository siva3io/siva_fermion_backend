package returns

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/shipping"

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
type SalesReturns struct {
	model_core.Model
	SalesReturnNumber      string                       `json:"sales_return_number" gorm:"type:text"`
	ReferenceNumber        string                       `json:"reference_number" gorm:"type:text"`
	CustomerName           string                       `json:"customer_name" gorm:"type:text"`
	ChannelName            string                       `json:"channel_name" gorm:"type:text"`
	ReasonId               *uint                        `json:"reason_id" gorm:"type:int"`
	Reason                 model_core.Lookupcode        `json:"reason" gorm:"foreignkey:ReasonId; references:ID"`
	ShippingModeId         *uint                        `json:"shipping_mode_id" gorm:"type:int"`
	ShippingMode           model_core.Lookupcode        `json:"shipping_mode" gorm:"foreignkey:ShippingModeId; references:ID"`
	Amount                 float32                      `json:"amount" gorm:"type:float"`
	ShippingCarrierId      *uint                        `json:"shipping_carrier_id" gorm:"type:int"`
	ShippingCarrier        *shipping.ShippingPartner    `json:"shipping_carrier" gorm:"foreignkey:ShippingCarrierId; references:ID"`
	CurrencyId             *uint                        `json:"currency_id" gorm:"type:int"`
	Currency               *model_core.Currency         `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	ExpectedDeliveyDate    datatypes.Date               `json:"expected_delivery_date" gorm:"type:date"`
	CustomerPickupAddress  datatypes.JSON               `json:"customer_pickup_address" gorm:"type:JSON"`
	CustomerBillingAddress datatypes.JSON               `json:"customer_billing_address" gorm:"type:JSON"`
	ShippingDetails        datatypes.JSON               `json:"shipping_details" gorm:"type:JSON;"`
	SrDate                 datatypes.Date               `json:"sr_date" gorm:"type:date"`
	StatusId               uint                         `json:"status_id" gorm:"type:INT"`
	Status                 model_core.Lookupcode        `json:"status" gorm:"foreignkey:StatusId; references:ID"`
	StatusHistory          datatypes.JSON               `json:"status_history" gorm:"type:JSON; default:'[]';not null"`
	CreditIssuedId         *uint                        `json:"credit_issued_id" gorm:"type:int"`
	CreditIssued           model_core.Lookupcode        `json:"credit_issued" gorm:"foreignkey:CreditIssuedId; references:ID"`
	SalesReturnLines       []SalesReturnLines           `json:"sales_return_lines" gorm:"foreignkey:SrId; references:ID"`
	SrPaymentDetails       SrPaymentDetails             `json:"sr_payment_details" gorm:"embedded"`
	PickupDateAndTime      orders.PickupDateAndTime     `json:"pickup_date_and_time" gorm:"embedded"`
	AdditionalInformation  orders.AdditionalInformation `json:"additional_information" gorm:"embedded"`
	SourceDocuments        datatypes.JSON               `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	SourceDocumentTypeId   *uint                        `json:"source_document_type_id" gorm:"type:integer"`
	SourceDocumentsType    model_core.Lookupcode        `json:"source_document_type" gorm:"foreignKey:SourceDocumentTypeId; references:ID"`
	TotalQuantity          int64                        `json:"total_quantity"`
	OrderId                *uint                        `json:"order_id"`
	Order                  *orders.SalesOrders          `json:"order" gorm:"foreignKey:OrderId;references:ID"`
	OrderDate              time.Time                    `json:"order_date" gorm:"type:time"`
}

type SalesReturnLines struct {
	model_core.Model
	SrId              uint                                   `json:"sr_id" gorm:"type:int"`
	ProductId         *uint                                  `json:"product_id" gorm:"type:int"`
	Product           *mdm.ProductVariant                    `json:"product_details" gorm:"foreignkey:ProductId; references:ID"`
	ProductTemplateId *uint                                  `json:"product_template_id" gorm:"type:INT"`
	ProductTemplate   *mdm.ProductTemplate                   `json:"product_template" gorm:"foreignkey:ProductTemplateId; references:ID"`
	UomId             *uint                                  `json:"uom_id" gorm:"type:INT"`
	Uom               *mdm.Uom                               `json:"uom" gorm:"foreignkey:UomId; references:ID"`
	SerialNumber      string                                 `json:"serial_number"  gorm:"type:text"`
	InventoryId       uint64                                 `json:"inventory_id"`
	QuantitySold      int                                    `json:"quantity_sold" gorm:"type:int"`
	QuantityReturned  int64                                  `json:"quantity_returned" gorm:"type:int"`
	ReturnTypeId      *uint                                  `json:"return_type_id" gorm:"type:int"`
	ReturnType        model_core.Lookupcode                  `json:"return_type" gorm:"foreignkey:ReturnTypeId; references:ID"`
	Rate              int                                    `json:"rate" gorm:"type:int"`
	ReturnLocationID  *uint                                  `json:"return_location_id" gorm:"type:int"`
	ReturnLocation    *shared_pricing_and_location.Locations `json:"return_location" gorm:"foreignkey:ReturnLocationID; references:ID"`
	Discount          float32                                `json:"discount" gorm:"type:float"`
	Tax               float32                                `json:"tax" gorm:"type:float"`
	Amount            float32                                `json:"amount" gorm:"type:float"`
}
type SrPaymentDetails struct {
	CustomerCredits      float32 `json:"customer_credits" gorm:"type:float"`
	UseCreditsForPayment bool    `json:"use_credits_for_payment" gorm:"type:boolean"`
	SubTotal             float32 `json:"sub_total" gorm:"type:float"`
	Tax                  float32 `json:"tax" gorm:"type:float"`
	ShippingCharges      float32 `json:"shipping_charges" gorm:"type:float"`
	Adjustments          float32 `json:"adjustments" gorm:"type:float"`
	TotalAmount          float32 `json:"total_amount" gorm:"type:float"`
}
