package sales_invoice

import (
	"time"

	"fermion/backend_core/internal/model/accounting"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/pkg/util/response"

	"github.com/lib/pq"
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
type SearchQuery struct {
	ReferenceNumber    string `json:"reference_number,omitempty" query:"reference_number"`
	SalesInvoiceNumber string `json:"sales_invoice_number,omitempty" query:"sales_invoice_number"`
	CustomerName       string `json:"customer_name,omitempty" query:"customer_name"`
	Status             string `json:"status,omitempty" query:"status"`
	Payment_type       string `json:"payment_type,omitempty" query:"payment_type"`
}

// Create and Update Request Payload for Sales Invoice
type (
	SalesInvoiceRequest struct {
		SalesInvoiceNumber          string                   `json:"sales_invoice_number"`
		AutoGenerateInvoiceNumber   bool                     `json:"auto_generate_invoice_number"`
		ReferenceNumber             string                   `json:"reference_number"`
		AutoGenerateReferenceNumber bool                     `json:"auto_generate_reference_number"`
		SalesInvoiceDate            string                   `json:"sales_invoice_date"`
		ExpectedShipmentDate        string                   `json:"expected_shipment_date"`
		CurrencyID                  uint                     `json:"currency_id"`
		SalesOrderIds               []int64                  `json:"sales_order_ids"`
		LinkSalesOrders             []map[string]interface{} `json:"link_sales_orders"`
		CustomerID                  uint                     `json:"customer_id"`
		Customer                    mdm.Partner              `json:"customer"`
		ChannelID                   *uint                    `json:"channel_id"`
		Channel                     *omnichannel.Channel     `json:"channel"`
		PaymentTypeID               uint                     `json:"payment_type_id"`
		PaymentType                 model_core.Lookupcode    `json:"payment_type"`
		CompanyId                   uint                     `json:"company_id"`
		StatusID                    uint                     `json:"status_id"`
		Status                      model_core.Lookupcode    `json:"status"`
		StatusHistory               []map[string]interface{} `json:"status_history"`
		Invoiced                    bool                     `json:"is_invoiced"`
		PaymentReceived             bool                     `json:"is_payment_received"`
		BalanceDueAmount            float64                  `json:"balance_due_amount"`
		AvailableCustomerCredits    float64                  `json:"available_customer_credits"`
		PaymentDueDate              time.Time                `json:"payment_due_date"`
		SalesInvoiceLines           []SalesInvoiceLines      `json:"sales_invoice_lines"`
		InternalNotes               string                   `json:"internal_notes"`
		ExternalNotes               string                   `json:"external_notes"`
		TermsAndConditions          string                   `json:"terms_and_conditions"`
		AttachmentFiles             []map[string]interface{} `json:"attachment_files"`
		UseCreditsForPayment        bool                     `json:"use_credits_for_payment"`
		SubTotalAmount              float64                  `json:"sub_total_amount"`
		ShippingAmount              float64                  `json:"shipping_amount"`
		TaxAmount                   float64                  `json:"tax_amount"`
		IgstAmt                     float64                  `json:"igst_amt" gorm:"type:double precision"`
		CgstAmt                     float64                  `json:"cgst_amt" gorm:"type:double precision"`
		SgstAmt                     float64                  `json:"sgst_amt" gorm:"type:double precision"`
		Adjustments                 float64                  `json:"adjustments"`
		CustomerCreditsAmount       float64                  `json:"customer_credits_amount"`
		TotalAmount                 float64                  `json:"total_amount"`
		BillingAddress              []map[string]interface{} `json:"billing_address"`
		DeliveryAddress             []map[string]interface{} `json:"delivery_address"`
		ShippingAddress             []map[string]interface{} `json:"shipping_address"`
		SourceDocuments             map[string]interface{}   `json:"source_documents"`
		SourceDocumentTypeId        *uint                    `json:"source_document_type_id"`
		PaymentTermsID              uint                     `json:"payment_terms_id"`
		PaymentTerms                model_core.Lookupcode    `json:"payment_terms"`

		OrderId   *uint  `json:"order_id"`
		OrderDate string `json:"order_date"`
		Quantity  int32  `json:"quantity"`
		model_core.Model
	}

	SalesInvoiceLines struct {
		ProductID        uint                  `json:"product_id"`
		ProductVariantID uint                  `json:"product_variant_id"`
		Description      string                `json:"description"`
		WarehouseID      uint                  `json:"warehouse_id"`
		InventoryID      uint                  `json:"inventory_id"`
		UomID            uint                  `json:"uom_id"`
		Quantity         uint                  `json:"quantity"`
		Discount         float64               `json:"discount"`
		DiscountTypeID   uint                  `json:"discount_type_id"`
		Tax              float64               `json:"tax"`
		TaxTypeID        uint                  `json:"tax_type_id"`
		PaymentTermsID   uint                  `json:"payment_terms_id"`
		PaymentTerms     model_core.Lookupcode `json:"payment_terms"`
		TotalAmount      float64               `json:"total_amount"`
		Price            float64               `json:"price"`
		model_core.Model
	}
)

// Create Sales Invoice response
type (
	SalesInvoiceCreate struct {
		Created_id int
	}

	SalesInvoiceCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesInvoiceCreate
		}
	} // @ name SalesInvoiceCreateResponse
)

// Bulk Create Sales Invoice response
type (
	BulkSalesInvoiceCreate struct {
		Created_id int
	}

	BulkSalesInvoiceCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkSalesInvoiceCreate
		}
	} // @ name BulkSalesInvoiceCreateResponse
)

// Update Sales Invoice response
type (
	SalesInvoiceUpdate struct {
		Updated_id int
	}

	SalesInvoiceUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesInvoiceUpdate
		}
	} // @ name SalesInvoiceUpdateResponse
)

// Delete Sales Invoice response
type (
	SalesInvoiceDelete struct {
		Deleted_id int
	}
	SalesInvoiceDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesInvoiceDelete
		} // @ name SalesInvoiceDeleteResponse
	}
)

// Delete Sales Invoice Lines response
type (
	SalesInvoiceLinesDelete struct {
		Deleted_id int
	}
	SalesInvoiceLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesInvoiceLinesDelete
		} // @ name SalesInvoiceLinesDeleteResponse
	}
)

// Get Sales Invoice response
type (
	SalesInvoiceGet struct {
		accounting.SalesInvoice
	}
	SalesInvoiceGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SalesInvoiceGet
		}
	} // @ name SalesInvoiceGetResponse
)

// Get all Sales Invoice response
type (
	SalesInvoiceGetAll struct {
		SalesInvoiceNumber          string                   `json:"sales_invoice_number"`
		AutoGenerateInvoiceNumber   bool                     `json:"auto_generate_invoice_number"`
		ReferenceNumber             string                   `json:"reference_number"`
		AutoGenerateReferenceNumber bool                     `json:"auto_generate_reference_number"`
		ExpectedShipmentDate        time.Time                `json:"expected_shipment_date"`
		CurrencyID                  uint                     `json:"currency_id"`
		Currency                    model_core.Currency      `json:"currency"`
		SalesOrderIds               pq.Int64Array            `json:"sales_order_ids"`
		LinkSalesOrders             map[string]interface{}   `json:"link_sales_orders"`
		CustomerID                  uint                     `json:"customer_id"`
		Customer                    mdm.Partner              `json:"customer"`
		ChannelID                   *uint                    `json:"channel_id"`
		Channel                     *omnichannel.Marketplace `json:"channel"`
		PaymentTypeID               uint                     `json:"payment_type_id"`
		PaymentType                 model_core.Lookupcode    `json:"payment_type"`
		StatusID                    uint                     `json:"status_id"`
		Status                      model_core.Lookupcode    `json:"status"`
		StatusHistory               map[string]interface{}   `json:"status_history"`
		Invoiced                    bool                     `json:"is_invoiced"`
		PaymentReceived             bool                     `json:"is_payment_received"`
		BalanceDueAmount            float64                  `json:"balance_due_amount"`
		PaymentTermsID              uint                     `json:"payment_terms_id"`
		PaymentTerms                model_core.Lookupcode    `json:"payment_terms"`
		AvailableCustomerCredits    float64                  `json:"available_customer_credits"`
		PaymentDueDate              time.Time                `json:"payment_due_date"`
		SalesInvoiceLines           []SalesInvoiceOrderLines `json:"sales_invoice_lines"`
		InternalNotes               string                   `json:"internal_notes"`
		ExternalNotes               string                   `json:"external_notes"`
		TermsAndConditions          string                   `json:"terms_and_conditions"`
		AttachmentFiles             map[string]interface{}   `json:"attachment_files"`
		UseCreditsForPayment        bool                     `json:"use_credits_for_payment"`
		SubTotalAmount              float64                  `json:"sub_total_amount"`
		ShippingAmount              float64                  `json:"shipping_amount"`
		TaxAmount                   float64                  `json:"tax_amount"`
		Adjustments                 float64                  `json:"adjustments"`
		CompanyId                   uint                     `json:"company_id"`
		CustomerCreditsAmount       float64                  `json:"customer_credits_amount"`
		TotalAmount                 float64                  `json:"total_amount"`
		BillingAddress              map[string]interface{}   `json:"billing_address"`
		DeliveryAddress             map[string]interface{}   `json:"delivery_address"`
		ShippingAddress             map[string]interface{}   `json:"shipping_address"`
		SourceDocuments             map[string]interface{}   `json:"source_documents"`
		SourceDocumentTypeId        *uint                    `json:"source_document_type_id"`
	}

	SalesInvoiceOrderLines struct {
		ProductID        uint                                  `json:"product_id"`
		Product          mdm.ProductTemplate                   `json:"products"`
		ProductVariantID uint                                  `json:"product_variant_id"`
		ProductVariant   mdm.ProductVariant                    `json:"product_variant"`
		Description      string                                `json:"description"`
		WarehouseID      uint                                  `json:"warehouse_id"`
		Warehouse        shared_pricing_and_location.Locations `json:"warehouse"`
		InventoryID      uint                                  `json:"inventory_id"`
		Inventory        mdm.CentralizedBasicInventory         `json:"inventory"`
		UomID            uint                                  `json:"uom_id"`
		UOM              mdm.Uom                               `json:"uom"`
		Quantity         uint                                  `json:"quantity"`
		Discount         float64                               `json:"discount"`
		DiscountTypeID   uint                                  `json:"discount_type_id"`
		DiscountType     model_core.Lookupcode                 `json:"discount_type"`
		Tax              float64                               `json:"tax"`
		TaxTypeID        uint                                  `json:"tax_type_id"`
		TaxType          model_core.Lookupcode                 `json:"tax_type"`
		PaymentTermsID   uint                                  `json:"payment_terms_id"`
		PaymentTerms     model_core.Lookupcode                 `json:"payment_terms"`
		TotalAmount      float64                               `json:"total_amount"`
	}

	SalesInvoicePaginatedResponse struct {
		Total_pages   uint `json:"total_pages"`
		Per_page      uint `json:"per_page"`
		Current_page  uint `json:"current_page"`
		Next_page     uint `json:"next_page"`
		Previous_page uint `json:"previous_page"`
		Total_rows    uint `json:"total_rows"`
	}
	SalesInvoiceGetAllResponse struct {
		Body struct {
			Meta       response.MetaResponse
			Data       []SalesInvoiceGetAll
			Pagination SalesInvoicePaginatedResponse
		}
	} // @ name SalesInvoiceGetAllResponse
)

// Search Sales Invoice response
type (
	SalesInvoiceSearch struct {
		SalesInvoiceNumber          string                   `json:"sales_invoice_number"`
		AutoGenerateInvoiceNumber   bool                     `json:"auto_generate_invoice_number"`
		ReferenceNumber             string                   `json:"reference_number"`
		AutoGenerateReferenceNumber bool                     `json:"auto_generate_reference_number"`
		ExpectedShipmentDate        time.Time                `json:"expected_shipment_date"`
		CurrencyID                  uint                     `json:"currency_id"`
		Currency                    model_core.Currency      `json:"currency"`
		SalesOrderIds               pq.Int64Array            `json:"sales_order_ids"`
		SalesOrders                 orders.SalesOrders       `json:"sales_orders"`
		LinkSalesOrders             map[string]interface{}   `json:"link_sales_orders"`
		CustomerID                  uint                     `json:"customer_id"`
		Customer                    mdm.Partner              `json:"customer"`
		ChannelID                   *uint                    `json:"channel_id"`
		Channel                     *omnichannel.Marketplace `json:"channel"`
		PaymentTypeID               uint                     `json:"payment_type_id"`
		PaymentType                 model_core.Lookupcode    `json:"payment_type"`
		StatusID                    uint                     `json:"status_id"`
		Status                      model_core.Lookupcode    `json:"status"`
		StatusHistory               map[string]interface{}   `json:"status_history"`
		Invoiced                    bool                     `json:"is_invoiced"`
		CompanyId                   uint                     `json:"company_id"`
		PaymentReceived             bool                     `json:"is_payment_received"`
		BalanceDueAmount            float64                  `json:"balance_due_amount"`
		PaymentTermsID              uint                     `json:"payment_terms_id"`
		PaymentTerms                model_core.Lookupcode    `json:"payment_terms"`
		AvailableCustomerCredits    float64                  `json:"available_customer_credits"`
		PaymentDueDate              time.Time                `json:"payment_due_date"`
		SalesInvoiceLines           []SalesInvoiceOrderLines `json:"sales_invoice_lines"`
		InternalNotes               string                   `json:"internal_notes"`
		ExternalNotes               string                   `json:"external_notes"`
		TermsAndConditions          string                   `json:"terms_and_conditions"`
		AttachmentFiles             map[string]interface{}   `json:"attachment_files"`
		UseCreditsForPayment        bool                     `json:"use_credits_for_payment"`
		SubTotalAmount              float64                  `json:"sub_total_amount"`
		ShippingAmount              float64                  `json:"shipping_amount"`
		TaxAmount                   float64                  `json:"tax_amount"`
		Adjustments                 float64                  `json:"adjustments"`
		CustomerCreditsAmount       float64                  `json:"customer_credits_amount"`
		TotalAmount                 float64                  `json:"total_amount"`
		BillingAddress              map[string]interface{}   `json:"billing_address"`
		DeliveryAddress             map[string]interface{}   `json:"delivery_address"`
		ShippingAddress             map[string]interface{}   `json:"shipping_address"`
		SourceDocuments             map[string]interface{}   `json:"source_documents"`
		SourceDocumentTypeId        *uint                    `json:"source_document_type_id"`
	}

	SalesInvoiceSearchResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []SalesInvoiceSearch
		}
	} // @ name SalesInvoiceSearchResponse
)

// Send Mail Sales Invoice response
type (
	SendMailSalesInvoice struct {
		ID            string `query:"id" json:"id"`
		ReceiverEmail string `query:"receiver_email" json:"receiver_email"`
	}

	SendMailSalesInvoiceResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SendMailSalesInvoice
		}
	} // @ name SendMailSalesInvoiceResponse
)

// Download Pdf Sales Invoice response
type (
	DownloadPdfSalesInvoice struct {
		ID string `query:"id"`
	}

	DownloadPdfSalesInvoiceResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DownloadPdfSalesInvoice
		}
	}
)
