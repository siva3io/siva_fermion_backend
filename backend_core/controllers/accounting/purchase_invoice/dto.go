package purchase_invoice

import (
	"fermion/backend_core/internal/model/accounting"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/pkg/util/response"
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
	ListPurchaseInvoiceDTO struct {
		model_core.Model
		PurchaseInvoiceNumber string                 `json:"purchase_invoice_number"`
		PurchaseInvoiceDate   string                 `json:"purchase_invoice_date"`
		ReferenceNumber       string                 `json:"reference_number"`
		VendorDetails         map[string]interface{} `json:"vendor_details"`
		StatusId              uint                   `json:"status_id"`
		Status                model_core.Lookupcode  `json:"status"`
		PaidId                uint                   `json:"paid_id"`
		Paid                  model_core.Lookupcode  `json:"paid"`
		PaymentTermsId        uint                   `json:"payment_terms_id"`
		PaymentTerms          model_core.Lookupcode  `json:"payment_terms"`
		PaymentAmount         float32                `json:"payment_amount"`
		CompanyId             uint                   `json:"company_id"`
		BalanceDue            float32                `json:"balance_due"`
		DueDate               string                 `json:"due_date"`
		ExpectedDeliveryDate  string                 `json:"expected_delivery_date"`
	}
	GetAllPurchaseInvoiceResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ListPurchaseInvoiceDTO
		}
	} //@ name GetAllPurchaseInvoiceResponse
)
type PurchaseInvoiceDTO struct {
	model_core.Model
	PurchaseInvoiceNumber  string                       `json:"purchase_invoice_number"`
	PurchaseInvoiceDate    string                       `json:"purchase_invoice_date"`
	ReferenceNumber        string                       `json:"reference_number"`
	VendorDetails          map[string]interface{}       `json:"vendor_details"`
	StatusId               uint                         `json:"status_id"`
	StatusHistory          map[string]interface{}       `json:"status_history"`
	PaidId                 uint                         `json:"paid_id"`
	PaymentTermsId         uint                         `json:"payment_terms_id"`
	PaymentDueDate         string                       `json:"payment_due_date"`
	PaymentAmount          float32                      `json:"payment_amount"`
	BalanceDue             float32                      `json:"balance_due"`
	DueDate                string                       `json:"due_date"`
	CompanyId              uint                         `json:"company_id"`
	ExpectedDeliveryDate   string                       `json:"expected_delivery_date"`
	LinkSourceDocumentType *uint                        `json:"link_source_document_type"`
	LinkSourceDocument     map[string]interface{}       `json:"link_source_document"`
	CurrencyId             uint                         `json:"currency_id"`
	DeliveryAddress        map[string]interface{}       `json:"delivery_address"`
	AdditionalInformation  orders.AdditionalInformation `json:"additional_information"`
	PaymentDetails         PiPaymentDetails             `json:"payment_details"`
	PurchaseInvoiceLines   []PurchaseInvoiceLines       `json:"purchase_invoice_lines"`
	PoIds                  []map[string]interface{}     `json:"po_ids"`
}

type PurchaseInvoiceLines struct {
	model_core.Model
	PurchaseInvoiceId uint    `json:"purchase_invoice_id"`
	ProductId         uint    `json:"product_id"`
	ProductTemplateId uint    `json:"product_template_id"`
	WarehouseId       uint    `json:"warehouse_id"`
	InventoryId       uint    `json:"inventory_id"`
	UomId             uint    `json:"uom_id"`
	PaymentTermsId    uint    `json:"payment_terms_id"`
	SerialNumber      string  `json:"serial_number"`
	Quantity          int     `json:"quantity"`
	Price             float32 `json:"price"`
	Discount          float32 `json:"discount"`
	Tax               float32 `json:"tax"`
	Amount            float32 `json:"amount"`
}
type PiPaymentDetails struct {
	AvailableVendorCredits float32 `json:"available_vendor_credits"`
	UseCreditsForPayment   bool    `json:"use_credits_for_payment"`
	SubTotal               float32 `json:"sub_total"`
	Tax                    float32 `json:"tax"`
	ShippingCharges        float32 `json:"shipping_charges"`
	AdjustmentAmount       float32 `json:"adjustment_amount"`
	TotalAmount            float32 `json:"total_amount"`
}

// Create Purchase Invoice response
type (
	PurchaseInvoiceCreate struct {
		Created_id int
	}
	PurchaseInvoiceCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseInvoiceCreate
		}
	} //@ name PurchaseInvoiceCreateResponse
)

// Update Purchase Invoice response
type (
	PurchaseInvoiceUpdate struct {
		Updated_id int
	}
	PurchaseInvoiceUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseInvoiceUpdate
		}
	} //@ name PurchaseInvoiceUpdateResponse
)

// Get Purchase Invoice response
type (
	PurchaseInvoiceGet struct {
		accounting.PurchaseInvoice
	}
	PurchaseInvoiceGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseInvoiceGet
		}
	} //@ name PurchaseInvoiceGetResponse
)

// Delete Purchase Invoice response
type (
	PurchaseInvoiceDelete struct {
		Deleted_id int
	}
	PurchaseInvoiceDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseInvoiceDelete
		}
	} //@ name PurchaseInvoiceDeleteResponse
)

// Delete Purchase Invoice lines response
type (
	PurchaseInvoiceLinesDelete struct {
		Deleted_id int
		Product_id int
	}
	PurchaseInvoiceLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseInvoiceLinesDelete
		}
	} //@ name PurchaseInvoiceLinesDeleteResponse
)
