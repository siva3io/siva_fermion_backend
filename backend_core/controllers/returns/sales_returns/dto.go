package sales_returns

import (
	"time"

	app_core "fermion/backend_core/controllers/cores"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/returns"
	"fermion/backend_core/internal/model/shipping"
	"fermion/backend_core/pkg/util/response"

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
type (
	ListSalesReturnsDTO struct {
		model_core.Model
		SalesReturnNumber    string                   `json:"sales_return_number"`
		ReferenceNumber      string                   `json:"reference_number"`
		SrDate               string                   `json:"sr_date"`
		CustomerName         string                   `json:"customer_name"`
		ShippingModeId       uint                     `json:"shipping_mode_id"`
		ShippingMode         app_core.LookupCode      `json:"shipping_mode"`
		ShippingDetails      map[string]interface{}   `json:"shipping_details"`
		ChannelName          string                   `json:"channel_name"`
		Amount               float32                  `json:"amount"`
		ReasonId             uint                     `json:"reason_id"`
		Reason               app_core.LookupCode      `json:"reason"`
		ShippingCarrierId    uint                     `json:"shipping_carrier_id"`
		ShippingCarrier      shipping.ShippingPartner `json:"shipping_carrier"`
		StatusId             uint                     `json:"status_id"`
		Status               app_core.LookupCode      `json:"status"`
		CreditIssuedId       uint                     `json:"credit_issued_id"`
		CreditIssued         app_core.LookupCode      `json:"credit_issued"`
		SourceDocuments      map[string]interface{}   `json:"source_documents"`
		SourceDocumentTypeId *uint                    `json:"source_document_type_id"`
		TotalQuantity        int64                    `json:"total_quantity"`
	}
	GetAllSalesReturnsResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ListSalesReturnsDTO
		}
	} //@ name GetAllSalesReturnsResponse
)
type SalesReturnsDTO struct {
	model_core.Model
	SalesReturnNumber      string                       `json:"sales_return_number" gorm:"type:text"`
	TotalQuantity          int64                        `json:"total_quantity"`
	ReferenceNumber        string                       `json:"reference_number"`
	SrDate                 string                       `json:"sr_date"`
	CustomerName           string                       `json:"customer_name"`
	ShippingModeId         *uint                        `json:"shipping_mode_id"`
	ShippingDetails        map[string]interface{}       `json:"shipping_details"`
	ChannelName            string                       `json:"channel_name"`
	Amount                 float32                      `json:"amount"`
	ReasonId               *uint                        `json:"reason_id"`
	ShippingCarrierId      *uint                        `json:"shipping_carrier_id"`
	StatusId               uint                         `json:"status_id"`
	StatusHistory          map[string]interface{}       `json:"status_history"`
	CreditIssuedId         *uint                        `json:"credit_issued_id"`
	CurrencyId             *uint                        `json:"currency_id"`
	ExpectedDeliveyDate    datatypes.Date               `json:"expected_delivery_date"`
	PickupDateAndTime      orders.PickupDateAndTime     `json:"pickup_date_and_time"`
	CustomerPickupAddress  map[string]interface{}       `json:"customer_pickup_address"`
	CustomerBillingAddress map[string]interface{}       `json:"customer_billing_address"`
	AdditionalInformation  orders.AdditionalInformation `json:"additional_information"`
	SalesReturnLines       []SalesReturnLines           `json:"sales_return_lines"`
	SrPaymentDetails       SrPaymentDetails             `json:"sr_payment_details"`
	SourceDocuments        map[string]interface{}       `json:"source_documents"`
	SourceDocumentTypeId   *uint                        `json:"source_document_type_id"`
	OrderId                *uint                        `json:"order_id"`
	Order                  *orders.SalesOrders          `json:"order" gorm:"foreignKey:OrderId;references:ID"`
	OrderDate              time.Time                    `json:"order_date" gorm:"type:time"`
}

type SalesReturnLines struct {
	model_core.Model
	SrId              uint    `json:"sr_id"`
	ProductId         uint    `json:"product_id"`
	ProductTemplateId uint    `json:"product_template_id"`
	InventoryId       uint64  `json:"inventory_id"`
	UomId             uint    `json:"uom_id"`
	QuantitySold      int     `json:"quantity_sold"`
	QuantityReturned  int64   `json:"quantity_returned"`
	ReturnTypeId      uint    `json:"return_type_id"`
	SerialNumber      string  `json:"serial_number"`
	Rate              int     `json:"rate"`
	ReturnLocationID  uint    `json:"return_location_id"`
	Discount          float32 `json:"discount"`
	Tax               float32 `json:"tax"`
	Amount            float32 `json:"amount"`
}
type SrPaymentDetails struct {
	SubTotal        float32 `json:"sub_total"`
	Tax             float32 `json:"tax"`
	ShippingCharges float32 `json:"shipping_charges"`
	CustomerCredits float32 `json:"customer_credits"`
	Adjustments     float32 `json:"adjustments"`
	TotalAmount     float32 `json:"total_amount"`
}

// Create Sales Returns response
type (
	SalesReturnsCreate struct {
		Created_id int
	}
	SalesReturnsCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesReturnsCreate
		}
	} //@ name SalesReturnsCreateResponse
)

// Update Sales Returns response
type (
	SalesReturnsUpdate struct {
		Updated_id int
	}
	SalesReturnsUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesReturnsUpdate
		}
	} //@ name SalesReturnsUpdateResponse
)

// Get Sales Returns response
type (
	SalesReturnsGet struct {
		returns.SalesReturns
	}
	SalesReturnsGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesReturnsGet
		}
	} //@ name SalesReturnsGetResponse
)

// Delete Sales Returns response
type (
	SalesReturnsDelete struct {
		Deleted_id int
	}
	SalesReturnsDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesReturnsDelete
		}
	} //@ name SalesReturnsDeleteResponse
)

// Delete Sales Returnlines response
type (
	SalesReturnLinesDelete struct {
		Deleted_id int
		Product_id int
	}
	SalesReturnLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesReturnLinesDelete
		}
	} //@ name SalesReturnLinesDeleteResponse
)
