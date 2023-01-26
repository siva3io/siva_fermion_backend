package shipping_partners

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
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
type ShippingPartnerRequest struct {
	PartnerName           string                 `json:"partner_name"`
	ShippingPartnerTypeId uint                   `json:"shipping_partner_type_id"`
	ProfileOptions        map[string]interface{} `json:"profile_options"`
	AuthOptions           map[string]interface{} `json:"auth_options"`
	SubscriptionOptions   map[string]interface{} `json:"subscription_options"`
	SpDefaultPreferences  map[string]interface{} `json:"sp_default_preferences"`
	SpTimePreferences     map[string]interface{} `json:"sp_time_preferences"`
	Id                    uint                   `json:"id"`
	AppId                 *uint                  `json:"app_id"`
} // @name ShippingPartnerRequest

type RateCalcRequest struct {
	PackageDetails map[string]interface{} `json:"package_details"`
} // @name RateCalcRequest

type UserShippingPartnerRegistrationRequest struct {
	model_core.Model
	ShippingPartnerId        uint                     `json:"shipping_partner_id"`
	UserId                   uint                     `json:"user_id"`
	AccountDetails           AccountDetails           `json:"account_details"`
	PickupDetails            PickupDetails            `json:"pickup_details"`
	PackingDetails           PackingDetails           `json:"packing_details"`
	AddAddressLocations      []AddLocations           `json:"add_locations"`
	CommericalInvoiceDetails CommericalInvoiceDetails `json:"commerical_invoice_details"`
	UploadDoc                UploadDoc                `json:"upload_doc"`
	TestAccount              bool                     `json:"test_account"`
	StatusId                 uint                     `json:"status_id"`
	StatusHistory            []map[string]interface{} `json:"status_history"`
} // @name UserShippingPartnerRegistrationRequest

type AddLocations struct {
	LocationId uint      `json:"location_id"`
	PickUpTime time.Time `json:"pick_up_time"`
	EndTime    time.Time `json:"end_time"`
}
type AccountDetails struct {
	Account_number       int64  `json:"account_number"`
	Web_service_key      string `json:"web_service_key"`
	Web_service_password string `json:"web_service_password"`
	Meter_number         string `json:"meter_number"`
	Carrier_weight_id    uint   `json:"carrier_weight_id"`
	Carrier_currency_id  uint   `json:"carrier_currency_id"`
	Show_service_rates   bool   `json:"show_service_rates"`
	Payment_type_id      uint   `json:"payment_type_id"`
}

type PickupDetails struct {
	Shipper_details_id uint      `json:"shipper_details_id"`
	From_time          time.Time `json:"from_time"`
	To_time            time.Time `json:"to_time"`
}

type PackingDetails struct {
	Use_volumetric_weight bool    `json:"use_volumetric_weight"`
	Packing_method_id     uint    `json:"packing_method_id"`
	Maximum_weight        float64 `json:"maximum_weight"`
}

type CommericalInvoiceDetails struct {
	IsCommericalInvoiceDetails bool    `json:"is_commercial_invoice_details"`
	Include_shipping_charges   bool    `json:"include_shipping_charges"`
	Include_insurance_charges  bool    `json:"include_insurance_charges"`
	Customer_declaration       string  `json:"customer_declaration"`
	Maximum_weight_allowed     float64 `json:"maximum_weight_allowed"`
	Terms_of_sales_id          uint    `json:"terms_of_sales_id"`
	Uom_id                     uint    `json:"uom_id"`
	Invoice_currency_id        uint    `json:"invoice_currency_id"`
	ProductId                  uint    `json:"product_id"`
	ProductTemplateId          uint    `json:"product_template_id"`
}

type UploadDoc struct {
	IsUploadDoc         bool                   `json:"is_upload_doc"`
	Customer_signature  map[string]interface{} `json:"customer_signature"`
	Company_letter_head map[string]interface{} `json:"company_letter_head"`
	Additional_document map[string]interface{} `json:"additional_document"`
}

type (
	CreateUserShippingPartnerRegistration struct {
		Created_UserShippingPartnerRegistration_ID int
	}
	UserShippingPartnerRegistrationCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CreateUserShippingPartnerRegistration
		}
	} // @name UserShippingPartnerRegistrationCreateResponse
)

type (
	UserShippingPartnerRegistrationUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@name UserShippingPartnerRegistrationUpdateResponse
)

type (
	UserShippingPartnerRegistrationGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@name UserShippingPartnerRegistrationGetResponse
)

type (
	UserShippingPartnerRegistrationGetAllResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []UserShippingPartnerRegistrationRequest
		}
	} //@name UserShippingPartnerRegistrationGetAllResponse
)

type (
	UserShippingPartnerRegistrationDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@name UserShippingPartnerRegistrationDeleteResponse
)

type (
	GetAllRateCalculator struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@name GetAllRateCalculator
)

type (
	ShippingPartnerGetAllResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ShippingPartnerRequest
		}
	} //@name ShippingPartnerGetAllResponse
)
