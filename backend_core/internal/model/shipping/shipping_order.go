package shipping

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/orders"

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
	ShippingOrder struct {
		model_core.Model
		ShippingNumber        string                `json:"shipping_number" gorm:"varchar(100)"`
		ReferenceNumber       string                `json:"reference_number" gorm:"type:varchar(100)"`
		ChannelId             *uint                 `json:"channel_id"`
		Channel               *omnichannel.Channel  `json:"channel" gorm:"foreignKey:ChannelId;references:ID"`
		BillingAddress        datatypes.JSON        `json:"billing_address" gorm:"type:json"`
		ReceiverAddress       datatypes.JSON        `json:"receiver_address" gorm:"type:json"`
		PackageDetails        datatypes.JSON        `json:"package_details" gorm:"type:json"`
		EstimatedWeight       float64               `json:"estimated_weight"`
		EstimatedShippingCost float64               `json:"estimated_shipping_cost"`
		ActualShippingCost    float64               `json:"actual_shipping_cost"`
		ShippingCost          float64               `json:"shipping_cost" gorm:"type:double precision"`
		ShippingStatusId      *uint                 `json:"shipping_status_id"`
		ShippingStatus        model_core.Lookupcode `json:"shipping_status" gorm:"foreignKey:ShippingStatusId;references:ID;"`
		ShippingPaymentTypeID *uint                 `json:"shipping_payment_type_id"`
		ShippingPaymentType   model_core.Lookupcode `json:"shipping_payment_type" gorm:"foreignKey:ShippingPaymentTypeID;references:ID"`
		PaymentStatus         string                `json:"payment_status"`
		Quantity              int32                 `json:"quantity" gorm:"type:integer"`

		AwbNumber         string          `json:"awb_number" gorm:"type:varchar(100)"`
		ShippingPartnerId *uint           `json:"shipping_partner_id"`
		ShippingPartner   ShippingPartner `json:"shipping_partner" gorm:"foreignKey:ShippingPartnerId;references:ID"`

		OrderId   *uint               `json:"order_id"`
		Order     *orders.SalesOrders `json:"order" gorm:"foreignKey:OrderId;references:ID"`
		OrderDate time.Time           `json:"order_date" gorm:"type:time"`

		StatusHistory      datatypes.JSON `json:"status_history" gorm:"type:json"`
		ShippingLabelId    datatypes.JSON `json:"shipping_label_id" gorm:"type:json"`
		ShippingManifestId datatypes.JSON `json:"shipping_manifest_id" gorm:"type:json"`
		ShippingInvoiceId  datatypes.JSON `json:"shipping_invoice_id" gorm:"type:json"`

		SenderAddress             datatypes.JSON        `json:"sender_address" gorm:"type:json"`
		IsShippingAddress         *bool                 `json:"is_shipping_address" gorm:"type:boolean"`
		ShippingDate              time.Time             `json:"shipping_date" gorm:"type:time"`
		PartnerId                 *uint                 `json:"partner_id"`
		Partner                   mdm.Partner           `json:"partner" gorm:"foreignKey:PartnerId;references:ID"`
		OrderValue                float64               `json:"order_value" gorm:"double precision"`
		ExpectedDeliveryDate      time.Time             `json:"expected_delivery_date" gorm:"type:time"`
		PickupAttempted           int                   `json:"pickup_attempted" gorm:"type:integer"`
		IsMarketplaceOrder        *bool                 `json:"is_marketplace_order" gorm:"type:boolean"`
		EunimartWalletAmount      float64               `json:"eunimart_wallet_amount" gorm:"type:double precision"`
		SetPickupDate             time.Time             `json:"set_pickup_date" gorm:"type:time"`
		SetPickupFromTime         time.Time             `json:"set_pickup_from_time" gorm:"type:time"`
		SetPickupToTime           time.Time             `json:"set_pickup_to_time" gorm:"type:time"`
		IsCod                     *bool                 `json:"is_cod" gorm:"type:boolean"`
		ShippingTypeId            *uint                 `json:"shipping_type_id"`
		ShippingType              model_core.Lookupcode `json:"shipping_type" gorm:"foreignKey:ShippingTypeId;references:ID"`
		ShippingModeId            *uint                 `json:"shipping_mode_id"`
		ShippingMode              model_core.Lookupcode `json:"shipping_mode" gorm:"foreignKey:ShippingModeId;references:ID"`
		ShippingOrderLines        []ShippingOrderLines  `json:"shipping_order_lines" gorm:"foreignKey:ShippingOrderId; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Billingdetails            datatypes.JSON        `json:"billing_details" gorm:"type:json"`
		DestinationCountryId      *uint                 `json:"destination_country_id" gorm:"column:destination_country_id"`
		DestinationCountry        model_core.Country    `json:"destination_country" gorm:"foreignKey:DestinationCountryId;references:ID"`
		DestinationZipcode        int32                 `json:"destination_zipcode" gorm:"type:integer"`
		OriginCountryId           *uint                 `json:"origin_country_id" gorm:"column:origin_country_id"`
		OriginCountry             model_core.Country    `json:"origin_country" gorm:"foreignKey:OriginCountryId;references:ID"`
		OriginZipcode             int32                 `json:"origin_zipcode" gorm:"type:integer"`
		PackageDirectionsId       *uint                 `json:"package_direction_id"`
		PackageDirections         model_core.Lookupcode `json:"package_direction" gorm:"foreignKey:PackageDirectionsId;references:ID"`
		CODStatus                 string                `json:"cod_status" gorm:"type:varchar(100)"`
		CODDueAmount              float64               `json:"cod_due_amount" gorm:"type:double precision"`
		CODAmountRecived          float64               `json:"cod_amount_received" gorm:"type:double precision"`
		CODDateAndTimeOfReceiving time.Time             `json:"cod_date_and_time_of_receiving" gorm:"type:time"`
		SupplierId                *uint                 `json:"supplier_id"`
		RatingValue               uint                  `json:"rating_value"`
		OndcDetails               datatypes.JSON        `json:"ondc_details" gorm:"type:json"`
	}

	ShippingOrderLines struct {
		model_core.Model
		ShippingOrderId   uint                `json:"-"`
		ProductId         *uint               `json:"product_id"`
		ProductVariant    mdm.ProductVariant  `json:"product_variant" gorm:"foreignKey:ProductId;references:ID"`
		ProductTemplateId *uint               `json:"product_template_id"`
		ProductTemplate   mdm.ProductTemplate `json:"product_template" gorm:"foreignKey:ProductTemplateId;references:ID"`
		ItemQuantity      int32               `json:"item_quantity" gorm:"type:integer"`
		UnitPrice         float64             `json:"unit_price" gorm:"type:double precision"`
		TaxPrice          float32             `json:"tax_price" gorm:"double precision"`
		GSTIN             string              `json:"gstin" gorm:"type:varchar(100)"`
		Discount          float64             `json:"discount" gorm:"double precision"`
	}

	RateCalculator struct {
		model_core.Model
		ShippingPartnerId     uint                `json:"Shipping_partner_id"`
		ShippingPartner       ShippingPartner     `json:"shipping_partner" gorm:"foreignKey:ShippingPartnerId;references:ID"`
		CommissionPreferences datatypes.JSON      `json:"commission_preferences" gorm:"type:json"`
		CurrencyId            uint                `json:"currency_id"`
		Currency              model_core.Currency `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	}

	TrackingDetails struct {
		ShippingOrderId uint             `json:"shipping_order_id"`
		ShippingOrder   ShippingOrder    `json:"shipping_order" gorm:"foreignKey:ShippingOrderId;references:ID"`
		TrackingStatus  []TrackingStatus `json:"tracking_status" gorm:"foreignKey:Tracking_status_id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}

	TrackingStatus struct {
		Status   string    `json:"status" gorm:"type:varchar(100)"`
		Location string    `json:"location" gorm:"type:varchar(100)"`
		DateTime time.Time `json:"date_time" gorm:"type:time"`
	}
)
