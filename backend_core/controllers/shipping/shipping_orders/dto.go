package shipping_orders

import (
	"time"

	core_dto "fermion/backend_core/controllers/cores"
	partner_dto "fermion/backend_core/controllers/mdm/contacts"
	locations "fermion/backend_core/controllers/mdm/locations"
	products_dto "fermion/backend_core/controllers/mdm/products"
	marketplace_dto "fermion/backend_core/controllers/omnichannel/marketplace"
	shipping_dto "fermion/backend_core/controllers/shipping/shipping_partners"
	"fermion/backend_core/internal/model/shipping"
	"fermion/backend_core/pkg/util/response"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type SearchKeys struct {
	Status              string `json:"status,omitempty" query:"status"`
	OriginCountry       string `json:"origin_country,omitempty" query:"origin_country"`
	OriginZipcode       int32  `json:"origin_zipcode,omitempty" query:"origin_zipcode"`
	DestinationZipcode  int32  `json:"destination_zipcode,omitempty" query:"destination_zipcode"`
	ShippingPreferences string `json:"shipping_preferences,omitempty" query:"shipping_preferences"`
}

type (
	CreateShippingOrder struct {
		Created_ShippingOrder_ID int
	}
	ShippingOrderCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CreateShippingOrder
		}
	} // @name ShippingOrderCreateResponse
)

type ShippingOrderRequest struct {
	ShippingNumber        string                    `json:"shipping_number"`
	ReferenceNumber       string                    `json:"reference_number"`
	ChannelId             uint                      `json:"channel_id"`
	SenderAddress         SenderAddress             `json:"sender_address"`
	ReceiverAddress       ReceiverOrBillingAddresss `json:"receiver_address"`
	BillingAddress        ReceiverOrBillingAddresss `json:"billing_address"`
	ShippingDate          time.Time                 `json:"shipping_date"`
	PartnerId             uint                      `json:"partner_id"`
	PackageDetails        PackageDetails            `json:"package_details"`
	ShippingCost          float64                   `json:"shipping_cost"`
	ShippingLabelId       map[string]interface{}    `json:"shipping_label_id"`
	ShippingManifestId    map[string]interface{}    `json:"shipping_manifest_id"`
	ShippingInvoiceId     map[string]interface{}    `json:"shipping_invoice_id"`
	OrderValue            float64                   `json:"order_value"`
	ShippingPaymentTypeID uint                      `json:"shipping_payment_type_id"`
	ExpectedDeliveryDate  time.Time                 `json:"expected_delivery_date"`
	PickupAttempted       int                       `json:"pickup_attempted"`
	IsMarketplaceOrder    bool                      `json:"is_marketplace_order"`
	OrderId               uint                      `json:"order_id"`
	OrderDate             time.Time                 `json:"order_date"`
	EunimartWalletAmount  float64                   `json:"eunimart_wallet_amount"`
	SetPickupDate         time.Time                 `json:"set_pickup_date"`
	SetPickupFromTime     time.Time                 `json:"set_pickup_from_time"`
	SetPickupToTime       time.Time                 `json:"set_pickup_to_time"`
	ShippingPartnerId     uint                      `json:"shipping_partner_id"`
	AwbNumber             string                    `json:"awb_number"`
	IsCod                 bool                      `json:"is_cod"`
	ShippingTypeId        uint                      `json:"shipping_type_id"`
	ShippingModeId        uint                      `json:"shipping_mode_id"`
	ShippingOrderLines    []ShippingOrderLines      `json:"shipping_order_lines"`
	Quantity              int32                     `json:"quantity"`
	Billingdetails        BillingDetails            `json:"billing_details"`
	DestinationCountryId  uint                      `json:"destination_country_id"`
	DestinationZipcode    int32                     `json:"destination_zipcode"`
	OriginCountryId       uint                      `json:"origin_country_id"`
	OriginZipcode         int32                     `json:"origin_zipcode"`
	PackageDirectionsId   uint                      `json:"package_direction_id"`
	CODStatus             string                    `json:"cod_status"`
	CODDueAmount          float64                   `json:"cod_due_amount"`
	ShippingStatusId      uint                      `json:"shipping_status_id"`
	IsShippingAddress     *bool                     `json:"is_shipping_address"`
	SupplierId            *uint                     `json:"supplier_id"`
	ID                    uint                      `json:"id"`
	CreatedByID           *uint                     `json:"created_by"`
	UpdatedByID           *uint                     `json:"updated_by"`
}

type ShippingOrderLines struct {
	ProductId         uint    `json:"product_id"`
	ProductTemplateId uint    `json:"product_template_id"`
	ItemQuantity      int32   `json:"item_quantity"`
	UnitPrice         float64 `json:"unit_price"`
	TaxPrice          float32 `json:"tax_price"`
	GSTIN             string  `json:"gstin"`
	Discount          float64 `json:"discount"`
}

type SenderAddress struct {
	SenderName     string  `json:"name"`
	MobileNumber   string  `json:"mobile"`
	PickupNickname string  `json:"nickname"`
	Email          string  `json:"email"`
	AddressLine1   string  `json:"addressline1"`
	AddressLine2   string  `json:"addressline2"`
	AddressLine3   string  `json:"addressline3"`
	Zipcode        int32   `json:"pincode"`
	City           string  `json:"city"`
	State          string  `json:"state"`
	Country        string  `json:"country"`
	Landmark       string  `json:"landmark"`
	Latitude       float32 `json:"latitude"`
	Longitude      float32 `json:"longitude"`
}

type ReceiverOrBillingAddresss struct {
	ReceiverName string  `json:"name"`
	MobileNumber string  `json:"mobile"`
	Email        string  `json:"email"`
	AddressLine1 string  `json:"addressline1"`
	AddressLine2 string  `json:"addressline2"`
	AddressLine3 string  `json:"addressline3"`
	Zipcode      int32   `json:"pincode"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Country      string  `json:"country"`
	Landmark     string  `json:"landmark"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
}

type PackageDetails struct {
	PackageHeight    int `json:"package_height"`
	PackageLength    int `json:"package_length"`
	PackageWeight    int `json:"package_weight"`
	PackageWidth     int `json:"package_width"`
	VolumetricWeight int `json:"volumetric_weight"`
	NoOfItems        int `json:"no_of_items"`
}

type BillingDetails struct {
	OrderId       string  `json:"order_id"`
	InvoiceNumber string  `json:"invoice_number"`
	Currency      string  `json:"currency"`
	TaxAmount     float64 `json:"tax_amount"`
	InvoiceAmount float64 `json:"invoice_amount"`
	GSTIN         string  `json:"gstin"`
	IECnumber     string  `json:"iec_number"`
	HSNcode       string  `json:"hsn_code"`
}

type (
	GetShippingOrder struct {
		shipping.ShippingOrder
	}
	ShippingOrderGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data GetShippingOrder
		}
	} // @name ShippingOrderGetResponse
)

type (
	GetAllShippingOrders struct {
		shipping.ShippingOrder
	}
	GetAllShippingOrdersResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []GetAllShippingOrders
		}
	} // @name GetAllShippingOrdersResponse
)

type (
	UpdateShippingOrder struct {
		Updated_ShippingOrder_Id int
	}
	UpdateShippingOrderResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data UpdateShippingOrder
		}
	} // @name UpdateShippingOrderResponse
)

type (
	DeleteShippingOrder struct {
		Deleted_ShippingOrder_Id int
	}
	DeleteShippingOrderResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteShippingOrder
		}
	} // @name DeleteShippingOrderResponse
)

type (
	DeleteShippingOrderLine struct {
		Deleted_GRN_Line_Item string
		Grn_Id                int
	}
	DeleteShippingOrderLineResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteShippingOrderLine
		}
	} // @name DeleteShippingOrderLineResponse
)

type (
	BulkShippingOrderCreate struct {
		Created bool
	}

	BulkShippingOrderCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkShippingOrderCreate
		}
	} // @ name BulkShippingOrderCreateResponse
)

type (
	ShippingOrderResponse struct {
		locations.ModelDto
		ShippingNumber            string                              `json:"shipping_number"`
		ReferenceNumber           string                              `json:"reference_number"`
		ChannelId                 uint                                `json:"channel_id"`
		Channel                   marketplace_dto.Marketplace         `json:"channel"`
		SenderAddress             map[string]interface{}              `json:"sender_address"`
		ReceiverAddress           map[string]interface{}              `json:"receiver_address"`
		BillingAddress            map[string]interface{}              `json:"billing_address"`
		IsShippingAddress         *bool                               `json:"is_shipping_address"`
		ShippingDate              time.Time                           `json:"shipping_date"`
		PartnerId                 uint                                `json:"partner_id"`
		Partner                   partner_dto.PartnerResponseDTO      `json:"partner"`
		PackageDetails            map[string]interface{}              `json:"package_details"`
		ShippingCost              float64                             `json:"shipping_cost"`
		ShippingStatusId          uint                                `json:"shipping_status_id"`
		ShippingStatus            core_dto.LookupCodesDTO             `json:"shipping_status"`
		StatusHistory             map[string]interface{}              `json:"status_history"`
		ShippingLabelId           map[string]interface{}              `json:"shipping_label_id"`
		ShippingManifestId        map[string]interface{}              `json:"shipping_manifest_id"`
		ShippingInvoiceId         map[string]interface{}              `json:"shipping_invoice_id"`
		OrderValue                float64                             `json:"order_value"`
		ShippingPaymentTypeID     uint                                `json:"shipping_payment_type_id"`
		ShippingPaymentType       core_dto.LookupCodesDTO             `json:"shipping_payment_type"`
		ExpectedDeliveryDate      time.Time                           `json:"expected_delivery_date"`
		PickupAttempted           int                                 `json:"pickup_attempted"`
		IsMarketplaceOrder        *bool                               `json:"is_marketplace_order"`
		OrderId                   uint                                `json:"order_id"`
		Order                     ListSalesOrdersDTO                  `json:"order"`
		OrderDate                 time.Time                           `json:"order_date"`
		EunimartWalletAmount      float64                             `json:"eunimart_wallet_amount"`
		SetPickupDate             time.Time                           `json:"set_pickup_date"`
		SetPickupFromTime         time.Time                           `json:"set_pickup_from_time"`
		SetPickupToTime           time.Time                           `json:"set_pickup_to_time"`
		ShippingPartnerId         uint                                `json:"shipping_partner_id"`
		ShippingPartner           shipping_dto.ShippingPartnerRequest `json:"shipping_partner"`
		AwbNumber                 string                              `json:"awb_number"`
		IsCod                     *bool                               `json:"is_cod"`
		ShippingTypeId            uint                                `json:"shipping_type_id"`
		ShippingType              core_dto.LookupCodesDTO             `json:"shipping_type"`
		ShippingModeId            uint                                `json:"shipping_mode_id"`
		ShippingMode              core_dto.LookupCodesDTO             `json:"shipping_mode"`
		ShippingOrderLines        []ShippingOrderLine                 `json:"shipping_order_lines"`
		Quantity                  int32                               `json:"quantity"`
		Billingdetails            map[string]interface{}              `json:"billing_details"`
		DestinationCountryId      uint                                `json:"destination_country_id"`
		DestinationCountry        core_dto.CountryDTO                 `json:"destination_country"`
		DestinationZipcode        int32                               `json:"destination_zipcode"`
		OriginCountryId           uint                                `json:"origin_country_id"`
		OriginCountry             core_dto.CountryDTO                 `json:"origin_country"`
		OriginZipcode             int32                               `json:"origin_zipcode"`
		PackageDirectionsId       uint                                `json:"package_direction_id"`
		PackageDirections         core_dto.LookupCodesDTO             `json:"package_direction"`
		CODStatus                 string                              `json:"cod_status"`
		CODDueAmount              float64                             `json:"cod_due_amount"`
		CODAmountRecived          float64                             `json:"cod_amount_received"`
		CODDateAndTimeOfReceiving time.Time                           `json:"cod_date_and_time_of_receiving"`
		SupplierId                *uint                               `json:"supplier_id"`
	}

	ShippingOrderLine struct {
		locations.ModelDto
		ShippingOrderId   uint                            `json:"-"`
		ProductId         uint                            `json:"product_id"`
		ProductVariant    products_dto.VariantResponseDTO `json:"product_variant"`
		ProductTemplateId uint                            `json:"product_template_id"`
		ProductTemplate   products_dto.TemplateReponseDTO `json:"product_template"`
		ItemQuantity      int32                           `json:"item_quantity"`
		UnitPrice         float64                         `json:"unit_price"`
		TaxPrice          float32                         `json:"tax_price"`
		GSTIN             string                          `json:"gstin"`
		Discount          float64                         `json:"discount"`
	}

	RateCalculator struct {
		locations.ModelDto
		ShippingPartnerId     uint                                `json:"Shipping_partner_id"`
		ShippingPartner       shipping_dto.ShippingPartnerRequest `json:"shipping_partner"`
		CommissionPreferences map[string]interface{}              `json:"commission_preferences"`
		CurrencyId            uint                                `json:"currency_id"`
		Currency              core_dto.CurrencyDTO                `json:"currency"`
	}

	TrackingDetails struct {
		ShippingOrderId uint                  `json:"shipping_order_id"`
		ShippingOrder   ShippingOrderResponse `json:"shipping_order"`
		TrackingStatus  []TrackingStatus      `json:"tracking_status"`
	}

	TrackingStatus struct {
		Status   string    `json:"status"`
		Location string    `json:"location"`
		DateTime time.Time `json:"date_time"`
	}
)

type ListSalesOrdersDTO struct {
	SalesOrderNumber     string              `json:"sales_order_number"`
	ReferenceNumber      string              `json:"reference_number"`
	SoDate               string              `json:"So_date"`
	CustomerName         string              `json:"customer_name"`
	ChannelName          string              `json:"channel_name"`
	PaymentTypeId        uint                `json:"payment_type_id"`
	PaymentType          core_dto.LookupCode `json:"payment_type"`
	StatusId             uint                `json:"status_id"`
	Status               core_dto.LookupCode `json:"status"`
	InvoicedId           uint                `json:"invoiced_id"`
	Invoiced             core_dto.LookupCode `json:"invoiced"`
	PaymentReceivedId    uint                `json:"payment_received_id"`
	PaymentReceived      core_dto.LookupCode `json:"payment_received"`
	PaymentAmount        float32             `json:"payment_amount"`
	PaymentTermId        uint                `json:"payment_term_id"`
	PaymentTerms         core_dto.LookupCode `json:"payment_terms"`
	ExpectedShippingDate string              `json:"expected_shipping_date"`
	core_dto.LookupCodesDTO
}
