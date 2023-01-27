package purchase_returns

import (
	app_core "fermion/backend_core/controllers/cores"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/returns"
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
	ListPurchaseReturnsDTO struct {
		PurchaseReturnNumber string                 `json:"purchase_return_number"`
		ReferenceNumber      string                 `json:"reference_number"`
		PurchaseOrder        orders.PurchaseOrders  `json:"purchase_order"`
		VendorDetails        map[string]interface{} `json:"vendor_details"`
		DebitNoteIssuedId    uint                   `json:"debit_note_issued_id"`
		DebitNoteIssued      app_core.LookupCode    `json:"debit_note_issued"`
		Amount               float32                `json:"amount"`
		PaymentTermsId       uint                   `json:"payment_terms_id"`
		PaymentTerms         app_core.LookupCode    `json:"payment_terms"`
		PrDate               string                 `json:"pr_date"`
		ExpectedDeliveryDate string                 `json:"expected_delivery_date"`
		StatusId             uint                   `json:"status_id"`
		Status               app_core.LookupCode    `json:"status"`
		Source_document_id   *uint                  `json:"source_document_id"`
		Source_document      map[string]interface{} `json:"source_documents"`
		TotalQuantity        int64                  `json:"total_quantity"`
		model_core.Model
	}
	GetAllPurchaseReturnsResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ListPurchaseReturnsDTO
		}
	} //@ name GetAllPurchaseReturnsResponse
)
type PurchaseReturnsDTO struct {
	model_core.Model
	PurchaseReturnNumber  string                       `json:"purchase_return_number" gorm:"type:text"`
	ReferenceNumber       string                       `json:"reference_number"`
	VendorDetails         map[string]interface{}       `json:"vendor_details"`
	DebitNoteIssuedId     uint                         `json:"debit_note_issued_id"`
	Amount                float32                      `json:"amount"`
	PaymentTermsId        uint                         `json:"payment_terms_id"`
	PaymentDueDate        string                       `json:"payment_due_date"`
	PrDate                string                       `json:"pr_date"`
	ExpectedDeliveryDate  string                       `json:"expected_delivery_date"`
	StatusId              uint                         `json:"status_id"`
	StatusHistory         map[string]interface{}       `json:"status_history"`
	CurrencyId            uint                         `json:"currency_id"`
	PickupDateAndTime     orders.PickupDateAndTime     `json:"pickup_date_and_time"`
	AdditionalInformation orders.AdditionalInformation `json:"additional_information"`
	PrPaymentDetails      PrPaymentDetails             `json:"pr_payment_details"`
	Source_document_id    *uint                        `json:"source_document_id"`
	SourceDocuments       map[string]interface{}       `json:"source_documents"`
	PurchaseReturnLines   []PurchaseReturnLines        `json:"purchase_return_lines"`
	TotalQuantity         int64                        `json:"total_quantity"`
}

type PurchaseReturnLines struct {
	model_core.Model
	PrId              uint    `json:"pr_id"`
	ProductId         uint    `json:"product_id"`
	ProductTemplateId uint    `json:"product_template_id"`
	InventoryId       uint64  `json:"inventory_id"`
	UomId             uint    `json:"uom_id"`
	QuantityPurchased int     `json:"quantity_purchased"`
	QuantityReturned  int64   `json:"quantity_returned"`
	SerialNumber      string  `json:"serial_number"`
	Rate              int     `json:"rate"`
	LocationID        uint    `json:"location_id"`
	Discount          float32 `json:"discount"`
	DiscountFormat    string  `json:"discount_format"`
	Tax               float32 `json:"tax"`
	TaxFormat         string  `json:"tax_format"`
	Amount            float32 `json:"amount"`
}
type PrPaymentDetails struct {
	SubTotal             float32 `json:"sub_total"`
	Tax                  float32 `json:"tax"`
	ShippingCharges      float32 `json:"shipping_charges"`
	VendorCredits        float32 `json:"vendor_credits"`
	UseCreditsForPayment bool    `json:"use_credits_for_payment"`
	Adjustments          float32 `json:"adjustments"`
	TotalAmount          float32 `json:"total_amount"`
}

// Create Purchase Returns response
type (
	PurchaseReturnsCreate struct {
		Created_id int
	}
	PurchaseReturnsCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseReturnsCreate
		}
	} //@ name PurchaseReturnsCreateResponse
)

// Update Purchase Returns response
type (
	PurchaseReturnsUpdate struct {
		Updated_id int
	}
	PurchaseReturnsUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseReturnsUpdate
		}
	} //@ name PurchaseReturnsUpdateResponse
)

// Get Purchase Returns response
type (
	PurchaseReturnsGet struct {
		returns.PurchaseReturns
	}
	PurchaseReturnsGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseReturnsGet
		}
	} //@ name PurchaseReturnsGetResponse
)

// Delete Purchase Returns response
type (
	PurchaseReturnsDelete struct {
		Deleted_id int
	}
	PurchaseReturnsDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseReturnsDelete
		}
	} //@ name PurchaseReturnsDeleteResponse
)

// Delete Purchase Returnlines response
type (
	PurchaseReturnLinesDelete struct {
		Deleted_id int
		Product_id int
	}
	PurchaseReturnLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PurchaseReturnLinesDelete
		}
	} //@ name PurchaseReturnLinesDeleteResponse
)
