package delivery_orders

import (
	core_dto "fermion/backend_core/controllers/cores"
	//app_core "fermion/backend_core/controllers/cores"
	contacts_dto "fermion/backend_core/controllers/mdm/contacts"
	locations_dto "fermion/backend_core/controllers/mdm/locations"
	products_dto "fermion/backend_core/controllers/mdm/products"
	uom_dto "fermion/backend_core/controllers/mdm/uom"
	shipping_dto "fermion/backend_core/controllers/shipping/shipping_orders"
	model_core "fermion/backend_core/internal/model/core"
	shared_pricing_and_location "fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
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
type SearchQuery struct {
	Reference_id          string `json:"reference_id,omitempty" query:"reference_id"`
	Delivery_order_number string `json:"delivery_order_number,omitempty" query:"delivery_order_number"`
	CustomerName          string `json:"customer_name,omitempty" query:"customer_name"`
	Tracking_number       string `json:"tracking_number,omitempty" query:"tracking_number"`
	Channel_name          string `json:"channel_name,omitempty" query:"channel_name"`
	Status                string `json:"status,omitempty" query:"status"`
	Payment_type          string `json:"payment_type,omitempty" query:"payment_type"`
}

// Create and Update Request Payload for Delivery Orders

type (
	DeliveryOrderRequest struct {
		model_core.Model
		DeliveryOrdersDetails  DeliveryOrdersDetails             `json:"delivery_order_details"`
		DeliveryAddressDetails map[string]interface{}            `json:"delivery_address_details"`
		BillingAddressDetails  map[string]interface{}            `json:"billing_address_details"`
		DeliveryOrderLines     []DeliveryOrderLines              `json:"delivery_order_lines"`
		Shipping_mode_id       uint                              `json:"shipping_mode_id"`
		PickupDateAndTime      DO_PickupDateAndTime              `json:"pickup_date_and_time"`
		Additionalinformation  Additionalinformation             `json:"additional_information"`
		Payment_details        Payment_details                   `json:"payment_details"`
		IsCreditUsed           *bool                             `json:"is_credit_used"`
		Status_id              uint                              `json:"status_id"`
		Status_history         map[string]interface{}            `json:"status_history"`
		Tracking_number        string                            `json:"tracking_number"`
		CustomerName           string                            `json:"customer_name"`
		ChannelName            string                            `json:"channel_name"`
		Contact_id             *uint                             `json:"contact_id"`
		App_id                 *uint                             `json:"app_id"`
		Warehouse_id           *int64                            `json:"warehouse_id"`
		Payment_type_id        *uint                             `json:"payment_type_id"`
		Shipping_Carrier_id    *uint                             `json:"shipping_carrier_id"`
		ShippingDetails        shipping_dto.ShippingOrderRequest `json:"shipping_details"`
		Shipping_order_id      *uint                             `json:"shipping_order_id"`
		IsShipping             *bool                             `json:"is_shipping"`
		TotalQuantity          int64                             `json:"total_quantity"`
	}

	DeliveryOrdersDetails struct {
		Reference_id                string                 `json:"reference_id"`
		Delivery_order_number       string                 `json:"delivery_order_number"`
		AutoCreateDoNumber          bool                   `json:"auto_create_do_number"`
		AutoGenerateReferenceNumber bool                   `json:"auto_generate_reference_number"`
		PriceListID                 *uint                  `json:"price_list_id"`
		Do_currency                 *uint                  `json:"do_currency"`
		Payment_due_date            string                 `json:"payment_due_date"`
		Payment_term_id             *uint                  `json:"payment_term_id"`
		Payment_Term                model_core.Lookupcode  `json:"payment_term"`
		Delivery_order_date         string                 `json:"delivery_order_date"`
		ExpectedShippingDate        string                 `json:"expected_shipping_date"`
		SourceDocuments             map[string]interface{} `json:"source_documents"`
		Source_document_id          *uint                  `json:"source_document_id"`
	}

	DeliveryOrderLines struct {
		model_core.Model
		DO_id               uint
		Product_id          *uint                 `json:"product_id"`
		Product_template_id *uint                 `json:"product_template_id"`
		Warehouse           string                `json:"warehouse"`
		Serial_no           string                `json:"serial_no"`
		Available_credit    uint                  `json:"available_credit"`
		Uom_id              *uint                 `json:"uom_id"`
		Rate                float64               `json:"rate"`
		Product_qty         int64                 `json:"quantity"`
		Discount            float64               `json:"discount_value"`
		Tax                 float64               `json:"tax"`
		Amount              float64               `json:"amount"`
		Description         string                `json:"description"`
		Location            string                `json:"location"`
		InventoryId         uint                  `json:"inventory_id"`
		Price               uint                  `json:"price"`
		PaymentTermsId      uint                  `json:"payment_terms_id"`
		PaymentTerms        model_core.Lookupcode `json:"payment_terms"`
	}

	Additionalinformation struct {
		Notes              string                 `json:"notes"`
		Terms_and_condtion string                 `json:"terms_and_condtion"`
		Attachments        map[string]interface{} `json:"attachments"`
	}

	Payment_details struct {
		Currency_id       *uint   `json:"currency_id"`
		Sub_total         float64 `json:"sub_total"`
		Payment_tax       float64 `json:"tax"`
		Shipping_charge   float64 `json:"shipping_charge"`
		Vender_credits    float64 `json:"vender_credits"`
		Adjustment_amount float64 `json:"adjustment_amount"`
		Total_amount      float64 `json:"total_amount"`
	}

	DO_PickupDateAndTime struct {
		PickupDate     string `json:"pickup_date"`
		PickupFromTime string `json:"pickup_from_time"`
		PickupToTime   string `json:"pickup_to_time"`
	}
)

// Create delivery order response
type (
	// DOCreate struct {
	// 	Created_id int
	// }
	DeliveryOrderCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name DeliveryOrderCreateResponse
)

// Update delivery order response
type (
	// DOUpdate struct {
	// 	Updated_id int
	// }
	DeliveryOrderUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name DeliveryOrderUpdateResponse
)

// Get  Delivery Order response
type (
	DOGet struct {
		orders.DeliveryOrders
	}
	DeliveryOrderGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DOGet
		}
	} //@ name DeliveryOrderGetResponse
)

// Get all Delivery Orders response
type (
	DeliveryOrderResponse struct {
		DeliveryOrdersDetails  DeliveryOrdersDetail              `json:"delivery_order_details"`
		DeliveryAddressDetails map[string]interface{}            `json:"delivery_address_details"`
		BillingAddressDetails  map[string]interface{}            `json:"billing_address_details"`
		DeliveryOrderLines     []DeliveryOrderLinesResponse      `json:"delivery_order_lines,omitempty"`
		Shipping_mode_id       uint                              `json:"shipping_mode_id"`
		Shipping_Mode          core_dto.LookupCodesDTO           `json:"shipping_mode"`
		PickupDateAndTime      DO_PickupDateAndTime              `json:"pickup_date_and_time"`
		Additionalinformation  AdditionalInformation             `json:"additional_information"`
		Payment_details        PaymentDetails                    `json:"payment_details"`
		IsCreditUsed           *bool                             `json:"is_credit_used"`
		Status_id              uint                              `json:"status_id"`
		Status                 core_dto.LookupCodesDTO           `json:"status"`
		Status_history         map[string]interface{}            `json:"status_history"`
		ChannelName            string                            `json:"channel_name"`
		Tracking_number        string                            `json:"tracking_number"`
		CustomerName           string                            `json:"customer_name"`
		Contact_id             *uint                             `json:"contact_id"`
		Contact_details        contacts_dto.PartnerResponseDTO   `json:"contact_details"`
		App_id                 uint                              `json:"app_id"`
		Warehouse_id           *int64                            `json:"warehouse_id"`
		Warehouse              locations_dto.LocationResponseDTO `json:"warehouse"`
		Payment_type_id        *uint                             `json:"payment_type_id"`
		Payment_Type           core_dto.LookupCodesDTO           `json:"payment_type"`
		Shipping_Carrier_id    *uint                             `json:"shipping_carrier_id"`
		Shipping_Carrier       core_dto.LookupCodesDTO           `json:"shipping_carrier"`
		Shipping_order_id      *uint                             `json:"shipping_order_id"`
		Shipping_order         core_dto.LookupCodesDTO           `json:"shipping_order"`
		IsShipping             *bool                             `json:"is_shipping"`
		ID                     uint                              `json:"id"`
		CreatedByID            *uint                             `json:"created_by"`
		UpdatedByID            *uint                             `json:"updated_by"`
	}

	DeliveryOrdersDetail struct {
		Reference_id          string                              `json:"reference_id"`
		Delivery_order_number string                              `json:"delivery_order_number"`
		PriceListID           *uint                               `json:"price_list_id"`
		PriceList             shared_pricing_and_location.Pricing `json:"pricing_details"`
		Do_currency           *uint                               `json:"do_currency"`
		Currency_details      core_dto.CurrencyDTO                `json:"currency_details"`
		Payment_due_date      string                              `json:"payment_due_date"`
		Payment_term_id       *uint                               `json:"payment_term_id"`
		Payment_Term          core_dto.LookupCodesDTO             `json:"payment_term"`
		Delivery_order_date   string                              `json:"delivery_order_date"`
		ExpectedShippingDate  string                              `json:"expected_shipping_date"`
		SourceDocuments       map[string]interface{}              `json:"source_documents"`
		Source_document_id    *uint                               `json:"source_document_id"`
		Source_document       core_dto.LookupCodesDTO             `json:"source_document"`
	}

	DeliveryOrderLinesResponse struct {
		DO_id               uint
		Product_id          *uint                           `json:"product_id"`
		Product             products_dto.VariantResponseDTO `json:"product_details"`
		Product_template_id *uint                           `json:"product_template_id"`
		Product_template    products_dto.TemplateReponseDTO `json:"product_template"`
		Warehouse           string                          `json:"warehouse"`
		Serial_no           string                          `json:"serial_no"`
		Available_credit    uint                            `json:"available_credit"`
		Uom_id              *uint                           `json:"uom_id"`
		UOM                 uom_dto.UomResponseDTO          `json:"uom_details"`
		Rate                float64                         `json:"rate"`
		Product_qty         int64                           `json:"quantity"`
		Discount            float64                         `json:"discount_value"`
		Tax                 float64                         `json:"tax"`
		Amount              float64                         `json:"amount"`
		Description         string                          `json:"description"`
		Location            string                          `json:"location"`
		InventoryId         uint                            `json:"inventory_id"`
		Price               uint                            `json:"price"`
		PaymentTermsId      uint                            `json:"payment_terms_id"`
		PaymentTerms        core_dto.LookupCodesDTO         `json:"payment_terms"`
		ID                  uint                            `json:"id"`
		CreatedByID         *uint                           `json:"created_by"`
		UpdatedByID         *uint                           `json:"updated_by"`
	}

	AdditionalInformation struct {
		Notes              string                 `json:"notes"`
		Terms_and_condtion string                 `json:"terms_and_condtion"`
		Attachments        map[string]interface{} `json:"attachments"`
	}

	PaymentDetails struct {
		Currency_id       *uint                `json:"currency_id"`
		Currency          core_dto.CurrencyDTO `json:"currency_info"`
		Sub_total         float64              `json:"sub_total"`
		Payment_tax       float64              `json:"tax"`
		Shipping_charge   float64              `json:"shipping_charge"`
		Vender_credits    float64              `json:"vender_credits"`
		Adjustment_amount float64              `json:"adjustment_amount"`
		Total_amount      float64              `json:"total_amount"`
	}
	DeliveryOrderGetAllResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []DeliveryOrderResponse
		}
	} //@ name DeliveryOrderGetAllResponse
)

// Delete Delivery Order response
type (
	// DODelete struct {
	// 	Deleted_id int
	// }
	DeliveryOrderDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name DeliveryOrderDeleteResponse
)

// Delete Deliver Orderlines response
type (
	// DOLDelete struct {
	// 	Deleted_id int
	// 	Product_id int
	// }
	DeliveryOrderLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name DeliveryOrderLinesDeleteResponse
)

// Bulk Create Delivery orders response
type (
	BulkDeliveryOrderCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} // @ name BulkDeliveryOrderCreateResponse
)
