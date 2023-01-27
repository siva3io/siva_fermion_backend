package purchase_orders

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
	ListPurchaseOrdersDTO struct {
		PurchaseOrderNumber  string                 `json:"purchase_order_number"`
		DateAndTime          string                 `json:"date_and_time"`
		ReferenceNumber      string                 `json:"reference_number"`
		VendorDetails        map[string]interface{} `json:"vendor_details"`
		CurrencyId           uint                   `json:"currency_id"`
		DeliveryToId         uint                   `json:"delivery_to_id"`
		OrganizationDetails  datatypes.JSON         `json:"organization_details"`
		StatusId             uint                   `json:"status_id"`
		Status               app_core.LookupCode    `json:"status"`
		BilledId             uint                   `json:"billed_id"`
		Billed               app_core.LookupCode    `json:"billed"`
		PaidId               uint                   `json:"paid_id"`
		Paid                 app_core.LookupCode    `json:"paid"`
		Amount               float32                `json:"amount"`
		TotalQuantity        int64                  `json:"total_quantity"`
		PaymentTermsId       uint                   `json:"payment_terms_id"`
		PaymentTerms         app_core.LookupCode    `json:"payment_terms"`
		ExpectedDeliveryDate string                 `json:"expected_delivery_date"`
		SourceDocuments      datatypes.JSON         `json:"source_documents" `
		SourceDocumentTypeId *uint                  `json:"source_document_type_id"`
		model_core.Model
	}
	GetAllPurchaseOrdersResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ListPurchaseOrdersDTO
		}
	} //@ name GetAllPurchaseOrdersResponse
)

type PurchaseOrdersDTO struct {
	model_core.Model
	PurchaseOrderNumber         string                 `json:"purchase_order_number"`
	ReferenceNumber             string                 `json:"reference_number"`
	CurrencyId                  *uint                  `json:"currency_id"`
	DateAndTime                 string                 `json:"date_and_time" gorm:"datetime"`
	ExpectedDeliveryDate        string                 `json:"expected_delivery_date" gorm:"type:date"`
	OrganizationDetails         map[string]interface{} `json:"organization_details"`
	TotalQuantity               int64                  `json:"total_quantity"`
	PaymentDueDate              string                 `json:"payment_due_date" gorm:"type:date"`
	PriceListID                 *uint                  `json:"price_list_id" gorm:"type:INT"`
	GenerateReferenceId         bool                   `json:"generate_reference_id"`
	StatusId                    uint                   `json:"status_id"`
	DeliveryToId                *uint                  `json:"delivery_to_id"`
	StatusHistory               datatypes.JSON         `json:"status_history"`
	BilledId                    *uint                  `json:"billed_id"`
	PaidId                      *uint                  `json:"paid_id"`
	GeneratePurchaseOrderNumber bool                   `json:"generate_purchase_order_number"`
	Amount                      float32                `json:"amount"`
	PaymentTermsId              *uint                  `json:"payment_terms_id"`
	VendorDetails               map[string]interface{} `json:"vendor_details"`
	DeliveryAddress             map[string]interface{} `json:"delivery_address"`
	BillingAddress              map[string]interface{} `json:"billing_address"`
	AdditionalInformation       AdditionalInformation  `json:"additional_information"`
	PoPaymentDetails            PoPaymentDetails       `json:"po_payment_details"`
	PurchaseOrderlines          []PurchaseOrderLines   `json:"purchase_order_lines"`
	SourceDocuments             datatypes.JSON         `json:"source_documents" `
	SourceDocumentTypeId        *uint                  `json:"source_document_type_id"`
}

type PurchaseOrderLines struct {
	model_core.Model
	PoId                uint    `json:"po_id"`
	ProductId           uint    `json:"product_id"`
	ProductTemplateId   uint    `json:"product_template_id"`
	WarehouseId         uint    `json:"warehouse_id"`
	InventoryId         uint    `json:"inventory_id"`
	UomId               uint    `json:"uom_id"`
	SerialNumber        string  `json:"serial_number"`
	Quantity            int64   `json:"quantity"`
	Price               float32 `json:"price"`
	Discount            float32 `json:"discount"`
	Tax                 float32 `json:"tax"`
	Amount              float32 `json:"amount"`
	SalesPeriod         string  `json:"sales_period"`
	CreditPeriod        string  `json:"credit_period"`
	LeadTime            string  `json:"lead_time"`
	ExpDeliveryLeadTime string  `json:"exp_delivery_lead_time"`
}

type AdditionalInformation struct {
	Notes              string                 `json:"notes"`
	TermsAndConditions string                 `json:"terms_and_conditions"`
	Attachments        map[string]interface{} `json:"attachments"`
}
type PoPaymentDetails struct {
	AvailableVendorCredits float32 `json:"available_vendor_credits"`
	UseCreditsForPayment   bool    `json:"use_credits_for_payment"`
	SubTotal               float32 `json:"sub_total"`
	Tax                    float32 `json:"tax"`
	ShippingCharges        float32 `json:"shipping_charges"`
	AdjustmentAmount       float32 `json:"adjustment_amount"`
	TotalAmount            float32 `json:"total_amount"`
}

// Create Purchase Orders response
type (
	PurchaseOrdersCreate struct {
		Created_id int
	}
	PurchaseOrdersCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseOrdersCreate
		}
	} //@ name PurchaseOrdersCreateResponse
)

// Update Purchase Orders response
type (
	PurchaseOrdersUpdate struct {
		Updated_id int
	}
	PurchaseOrdersUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseOrdersUpdate
		}
	} //@ name PurchaseOrdersUpdateResponse
)

// Get Purchase Orders response
type (
	PurchaseOrdersGet struct {
		orders.PurchaseOrders
	}
	PurchaseOrdersGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseOrdersGet
		}
	} //@ name PurchaseOrdersGetResponse
)

// Delete Purchase orders response
type (
	PurchaseOrdersDelete struct {
		Deleted_id int
	}
	PurchaseOrdersDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseOrdersDelete
		}
	} //@ name PurchaseOrdersDeleteResponse
)

// Delete Purchase Orderlines response
type (
	PurchaseOrderLinesDelete struct {
		Deleted_id int
		Product_id int
	}
	PurchaseOrderLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseOrderLinesDelete
		}
	} //@ name PurchaseOrderLinesDeleteResponse
)
