package sales_orders

import (
	app_core "fermion/backend_core/controllers/cores"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
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
	ListSalesOrdersDTO struct {
		SalesOrderNumber     string              `json:"sales_order_number"`
		ReferenceNumber      string              `json:"reference_number"`
		SoDate               string              `json:"So_date"`
		CustomerName         string              `json:"customer_name"`
		ChannelName          string              `json:"channel_name"`
		PaymentTypeId        uint                `json:"payment_type_id"`
		PaymentType          app_core.LookupCode `json:"payment_type"`
		StatusId             uint                `json:"status_id"`
		Status               app_core.LookupCode `json:"status"`
		InvoicedId           uint                `json:"invoiced_id"`
		Invoiced             app_core.LookupCode `json:"invoiced"`
		PaymentReceivedId    uint                `json:"payment_received_id"`
		PaymentReceived      app_core.LookupCode `json:"payment_received"`
		PaymentAmount        float32             `json:"payment_amount"`
		PaymentTermId        uint                `json:"payment_term_id"`
		PaymentTerms         app_core.LookupCode `json:"payment_terms"`
		ExpectedShippingDate string              `json:"expected_shipping_date"`
		model_core.Model
	}
	GetAllSalesOrdersResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ListSalesOrdersDTO
		}
	} //@ name GetAllSalesOrdersResponse
)

type SalesOrdersDTO struct {
	model_core.Model
	ReferenceNumber         string                 `json:"reference_number"`
	SoDate                  string                 `json:"so_date"`
	CurrencyId              uint                   `json:"currency_id"`
	CustomerName            string                 `json:"customer_name"`
	CustomerShippingAddress map[string]interface{} `json:"customer_shipping_address"`
	CustomerBillingAddress  map[string]interface{} `json:"customer_billing_address"`
	ChannelName             string                 `json:"channel_name"`
	PaymentTypeId           uint                   `json:"payment_type_id"`
	StatusId                uint                   `json:"status_id"`
	StatusHistory           datatypes.JSON         `json:"status_history"`
	InvoicedId              uint                   `json:"invoiced_id"`
	PaymentReceivedId       uint                   `json:"payment_received_id"`
	PaymentAmount           float32                `json:"payment_amount"`
	PaymentTermsId          uint                   `json:"payment_terms_id"`
	PaymentDueDate          string                 `json:"payment_due_date"`
	ExpectedShippingDate    string                 `json:"expected_shipping_date"`
	VendorDetails           map[string]interface{} `json:"vendor_details"`
	AdditionalInformation   AdditionalInformation  `json:"additional_information"`
	SalesOrderlines         []SalesOrderlines      `json:"sales_order_lines"`
	SoPaymentDetails        SoPaymentDetails       `json:"so_payment_details"`
}

type SalesOrderlines struct {
	model_core.Model
	SoId              uint                  `json:"so_id"`
	ProductId         uint                  `json:"product_id"`
	WarehouseId       uint                  `json:"warehouse_id"`
	SerialNumber      string                `json:"serial_number"`
	ProductTemplateId uint                  `json:"product_template_id"`
	InventoryId       uint                  `json:"inventory_id"`
	UomId             uint                  `json:"uom_id"`
	Quantity          int                   `json:"quantity"`
	Price             float32               `json:"price"`
	Discount          float32               `json:"discount"`
	Tax               float32               `json:"tax"`
	Amount            float32               `json:"amount"`
	Description       string                `json:"description"`
	Location          string                `json:"location"`
	PaymentTermsId    uint                  `json:"payment_terms_id"`
	PaymentTerms      model_core.Lookupcode `json:"payment_terms"`
}
type AdditionalInformation struct {
	Notes              string                 `json:"notes"`
	TermsAndConditions string                 `json:"terms_and_conditions"`
	Attachments        map[string]interface{} `json:"attachments"`
}
type SoPaymentDetails struct {
	AvailableCustomerCredits float32 `json:"available_customer_credits"`
	UseCreditsForPayment     bool    `json:"use_credits_for_payment"`
	SubTotal                 float32 `json:"sub_total"`
	Tax                      float32 `json:"tax"`
	ShippingCharges          float32 `json:"shipping_charges"`
	AdjustmentAmount         float32 `json:"adjustment_amount"`
	TotalAmount              float32 `json:"total_amount"`
}

// Create Sales Orders response
type (
	SalesOrdersCreate struct {
		Created_id int
	}
	SalesOrdersCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesOrdersCreate
		}
	} //@ name SalesOrdersCreateResponse
)

// Update Sales Orders response
type (
	SalesOrdersUpdate struct {
		Updated_id int
	}
	SalesOrdersUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesOrdersUpdate
		}
	} //@ name SalesOrdersUpdateResponse
)

// Get Sales Orders response
type (
	SalesOrdersGet struct {
		orders.SalesOrders
	}
	SalesOrdersGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesOrdersGet
		}
	} //@ name SalesOrdersGetResponse
)

// Delete Sales Orders response
type (
	SalesOrdersDelete struct {
		Deleted_id int
	}
	SalesOrdersDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesOrdersDelete
		}
	} //@ name SalesOrdersDeleteResponse
)

// Delete Sales Orderlines response
type (
	SalesOrderLinesDelete struct {
		Deleted_id int
		Product_id int
	}
	SalesOrderLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesOrderLinesDelete
		}
	} //@ name SalesOrderLinesDeleteResponse
)
