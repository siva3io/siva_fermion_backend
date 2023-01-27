package sales_orders

import (
	"time"

	app_core "fermion/backend_core/controllers/cores"
	"fermion/backend_core/controllers/mdm/products"
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
	ListSalesOrdersDTO struct {
		SalesOrderNumber     string                `json:"sales_order_number"`
		ReferenceNumber      string                `json:"reference_number"`
		SoDate               string                `json:"so_date"`
		CustomerName         string                `json:"customer_name"`
		ChannelName          string                `json:"channel_name"`
		TotalQuantity        int64                 `json:"total_quantity"`
		PaymentTypeId        uint                  `json:"payment_type_id"`
		PaymentType          app_core.LookupCode   `json:"payment_type"`
		StatusId             uint                  `json:"status_id"`
		Status               app_core.LookupCode   `json:"status"`
		CancellationReasonId *uint                 `json:"cancellation_reason_id"`
		CancellationReason   model_core.Lookupcode `json:"cancellation_reason"`
		InvoicedId           uint                  `json:"invoiced_id"`
		Invoiced             app_core.LookupCode   `json:"invoiced"`
		PaymentReceivedId    uint                  `json:"payment_received_id"`
		PaymentReceived      app_core.LookupCode   `json:"payment_received"`
		Amount               float32               `json:"amount"`
		PaymentTermId        uint                  `json:"payment_term_id"`
		PaymentTerms         app_core.LookupCode   `json:"payment_terms"`
		ExpectedShippingDate string                `json:"expected_shipping_date"`
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
	ID                            uint      `json:"id"`
	CreatedByID                   *uint     `json:"created_by"`
	UpdatedByID                   *uint     `json:"updated_by"`
	DeletedByID                   *uint     `json:"deleted_by"`
	UpdatedDate                   time.Time `json:"updated_date"`
	CreatedDate                   time.Time `json:"created_date"`
	CreatedBy                     *app_core.UserResponseDTO
	CompanyId                     uint                     `json:"company_id"`
	SalesOrderNumber              string                   `json:"sales_order_number"`
	ReferenceNumber               string                   `json:"reference_number"`
	SoDate                        string                   `json:"so_date,omitempty"`
	CurrencyId                    *uint                    `json:"currency_id"`
	CustomerName                  string                   `json:"customer_name"`
	TotalQuantity                 int64                    `json:"total_quantity"`
	InvoicedId                    *uint                    `json:"invoiced_id" gorm:"type:int"`
	Invoiced                      app_core.LookupCodesDTO  `json:"invoiced"`
	PaymentReceivedId             *uint                    `json:"payment_received_id" gorm:"type:int"`
	PaymentReceived               app_core.LookupCodesDTO  `json:"payment_received"`
	Item                          SalesOrderlines          `json:"item"`
	SalesOrderLines               []SalesOrderlines        `json:"sales_order_lines"`
	AdditionalInformation         AdditionalInformation    `json:"additional_information" gorm:"embedded"`
	SoPaymentDetails              SoPaymentDetails         `json:"so_payment_details" gorm:"embedded"`
	CustomerShippingAddress       map[string]interface{}   `json:"customer_shipping_address,omitempty"`
	CustomerBillingAddress        map[string]interface{}   `json:"customer_billing_address,omitempty"`
	ChannelName                   string                   `json:"channel_name"`
	PaymentTypeId                 uint                     `json:"payment_type_id,omitempty"`
	PaymentType                   app_core.LookupCodesDTO  `json:"payment_type"`
	StatusId                      *uint                    `json:"status_id"`
	Status                        app_core.LookupCodesDTO  `json:"status"`
	StatusHistory                 []map[string]interface{} `json:"status_history"`
	CancellationReasonId          *uint                    `json:"cancellation_reason_id"`
	CancellationReason            app_core.LookupCodesDTO  `json:"cancellation_reason"`
	Amount                        float32                  `json:"amount"`
	PaymentTermsId                *uint                    `json:"payment_terms_id"`
	PaymentTerms                  app_core.LookupCodesDTO  `json:"payment_terms"`
	PaymentDueDate                string                   `json:"payment_due_date,omitempty"`
	ExpectedShippingDate          string                   `json:"expected_shipping_date,omitempty"`
	VendorDetails                 map[string]interface{}   `json:"vendor_details"`
	OndcContext                   map[string]interface{}   `json:"ondc_context,omitempty"`
	SellerNpTypeId                *uint                    `json:"seller_np_type_id"`
	SellerNpType                  app_core.LookupCodesDTO  `json:"seller_np_type"`
	SellerPinCode                 string                   `json:"seller_pin_code"`
	SellerCity                    string                   `json:"seller_city"`
	OrderCategoryId               *uint                    `json:"order_category_id"`
	OrderCategory                 app_core.LookupCodesDTO  `json:"order_category"`
	DeliveryTypeId                *uint                    `json:"delivery_type_id"`
	DeliveryType                  app_core.LookupCodesDTO  `json:"delivery_type"`
	DeliveryDate                  time.Time                `json:"delivery_date"`
	LogisticsNetworkOrderId       string                   `json:"logistics_network_order_id"`
	LogisticsNetworkTransactionId string                   `json:"logistics_network_transaction_id"`
	LogisticsSellerNpName         string                   `json:"logistics_seller_np_name"`
	CancelledByID                 *uint                    `json:"cancelled_by"`
	CancelledBy                   app_core.UserResponseDTO
	CancelledDate                 time.Time               `json:"cancelled_date"`
	CancelById                    *uint                   `json:"cancel_by_id"`
	CancelBy                      app_core.LookupCodesDTO `json:"cancel_by"`
	ShippedDate                   time.Time               `josn:"shipped_date"`
	ReadyTOShip                   time.Time               `json:"ready_to_ship"`
}

type OndcOrdersResponseDTO struct {
	ID                            uint                     `json:"id"`
	CreatedByID                   *uint                    `json:"created_by"`
	UpdatedByID                   *uint                    `json:"updated_by"`
	DeletedByID                   *uint                    `json:"deleted_by"`
	UpdatedDate                   time.Time                `json:"updated_date"`
	CreatedDate                   time.Time                `json:"created_date"`
	CompanyId                     uint                     `json:"company_id"`
	SalesOrderNumber              string                   `json:"sales_order_number"`
	ReferenceNumber               string                   `json:"reference_number"`
	SoDate                        string                   `json:"so_date,omitempty"`
	CurrencyId                    *uint                    `json:"currency_id"`
	CustomerName                  string                   `json:"customer_name"`
	TotalQuantity                 int64                    `json:"total_quantity"`
	InvoicedId                    *uint                    `json:"invoiced_id" gorm:"type:int"`
	PaymentReceivedId             *uint                    `json:"payment_received_id" gorm:"type:int"`
	SalesOrderLines               []SalesOrderlines        `json:"sales_order_lines"`
	AdditionalInformation         AdditionalInformation    `json:"additional_information" gorm:"embedded"`
	SoPaymentDetails              SoPaymentDetails         `json:"so_payment_details" gorm:"embedded"`
	CustomerShippingAddress       map[string]interface{}   `json:"customer_shipping_address,omitempty"`
	CustomerBillingAddress        map[string]interface{}   `json:"customer_billing_address,omitempty"`
	ChannelName                   string                   `json:"channel_name"`
	PaymentTypeId                 uint                     `json:"payment_type_id,omitempty"`
	StatusId                      *uint                    `json:"status_id"`
	Status                        app_core.LookupCodesDTO  `json:"status"`
	StatusHistory                 []map[string]interface{} `json:"status_history"`
	CancellationReasonId          *uint                    `json:"cancellation_reason_id"`
	Amount                        float32                  `json:"amount"`
	PaymentTermsId                *uint                    `json:"payment_terms_id"`
	PaymentDueDate                string                   `json:"payment_due_date,omitempty"`
	ExpectedShippingDate          string                   `json:"expected_shipping_date,omitempty"`
	VendorDetails                 map[string]interface{}   `json:"vendor_details"`
	OndcContext                   map[string]interface{}   `json:"ondc_context,omitempty"`
	DeliveryDate                  time.Time                `json:"delivery_date"`
	LogisticsNetworkOrderId       string                   `json:"logistics_network_order_id"`
	LogisticsNetworkTransactionId string                   `json:"logistics_network_transaction_id"`
	CancelledByID                 *uint                    `json:"cancelled_by"`
	CancelledDate                 time.Time                `json:"cancelled_date"`
}
type SalesOrderlines struct {
	ID                uint                         `json:"id"`
	CreatedByID       *uint                        `json:"created_by"`
	UpdatedByID       *uint                        `json:"updated_by"`
	DeletedByID       *uint                        `json:"deleted_by"`
	CompanyId         *uint                        `json:"company_id"`
	SoId              uint                         `json:"so_id"`
	ProductId         *uint                        `json:"product_id"`
	Product           *products.VariantResponseDTO `json:"product_details"`
	WarehouseId       *uint                        `json:"warehouse_id"`
	ProductTemplateId *uint                        `json:"product_template_id"`
	InventoryId       *uint                        `json:"inventory_id"`
	UomId             *uint                        `json:"uom_id"`
	SerialNumber      string                       `json:"serial_number"`
	Quantity          int64                        `json:"quantity"`
	Price             float32                      `json:"price"`
	Discount          float32                      `json:"discount"`
	Tax               float32                      `json:"tax"`
	Amount            float32                      `json:"amount"`
	Description       string                       `json:"description"`
	Location          string                       `json:"location"`
	PaymentTermsId    *uint                        `json:"payment_terms_id"`
	PaymentTerms      app_core.LookupCodesDTO      `json:"payment_terms"`
	SKUName           string                       `json:"sku_name"`
	SKUCode           string                       `json:"sku_code"`
	StatusId          *uint                        `json:"status_id"`
}
type AdditionalInformation struct {
	Notes              string                   `json:"notes"`
	TermsAndConditions string                   `json:"terms_and_conditions"`
	Attachments        []map[string]interface{} `json:"attachments,omitempty"`
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

type IdNameDTO struct {
	Name string `json:"name"`
	Id   uint   `json:"id"`
}

type CustomerBillingorShippingAddress struct {
	City                string `json:"city"`
	State               string `json:"state"`
	Country             string `json:"country"`
	PinCode             string `json:"pin_code"`
	Landmark            string `json:"land_mark"`
	LocationName        string `json:"location_name"`
	AddressLine1        string `json:"address_line_1"`
	AddressLine2        string `json:"address_line_2"`
	AddressLine3        string `json:"address_line_3"`
	ContactPersonNumber string `json:"contact_person_number"`
	ContactPersonName   string `json:"contact_person_name"`
	GstInNumber         string `json:"gst_in_number"`
}

type csvDto struct {
	Buyer_NP_Name                    string `json:"Buyer_NP_Name"`
	Seller_NP_Name                   string `json:"Seller_NP_Name"`
	Order_Create_Date                string `json:"Order_Create_Date_&_Time"`
	Network_Order_Id                 string `json:"Network_Order_Id"`
	Network_Transaction_Id           string `json:"Network_Transaction_Id"`
	Seller_NP_Order_Item_Id          string `json:"Seller_NP_Order_Item_Id"`
	Seller_NP_Type                   string `json:"Seller_NP_Type"`
	Order_Status                     string `json:"Order_Status"`
	Name_of_Seller                   string `json:"Name_of_Seller"`
	Seller_Pincode                   string `json:"Seller_Pincode"`
	SKU_Name                         string `json:"SKU_Name"`
	SKU_Code                         string `json:"SKU_Code"`
	Order_Category                   string `json:"Order_Category"`
	Shipped_At_Date_                 string `json:"Shipped_At_Date_&_Time"`
	Delivered_At_Date_               string `json:"Delivered_At_Date_&_Time"`
	Delivery_Type                    string `json:"Delivery_Type"`
	Logistics_Network_Order_Id       string `json:"Logistics_Network_Order_Id"`
	Logistics_Network_Transaction_Id string `json:"Logistics_Network_Transaction_Id"`
	Delivery_City                    string `json:"Delivery_City"`
	Delivery_Pincode                 string `json:"Delivery_Pincode"`
	Cancelled_At_Date                string `json:"Cancelled_At_Date_&_Time"`
	Cancelled_By                     string `json:"Cancelled_By"`
	Cancellation_Reason              string `json:"Cancellation_Reason"`
	Cancellation_Remark              string `json:"Cancellation_Remark"`
	Total_Order_Value                string `json:"Total_Order_Value"`
}
