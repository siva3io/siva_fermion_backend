package pricing

import (
	app_core "fermion/backend_core/controllers/cores"
	location "fermion/backend_core/controllers/mdm/locations"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
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
type PricingDTO struct {
	Id              uint        `json:"id"`
	Price_list_name string      `json:"price_list_name"`
	Currency_id     *uint       `json:"currency_id"`
	Currency        CurrencyDTO `json:"currency"`
	Price_list_id   uint        `json:"price_list_id"` // this field holds the value of which table we need to do the operation
	// 1.sales   2.purchase   3.transfer
	StartDate           string                `json:"start_date"`
	EndDate             string                `json:"end_date"`
	Description         string                `json:"description"`
	Price_list_rule     string                `json:"price_list_rule"`
	StatusId            *uint                 `json:"status_id" `
	Status              app_core.LookupCode   `json:"status" `
	SalesPriceListId    *uint                 `json:"sales_price_list_id"`
	SalesPriceList      *SalesPriceListDTO    `json:"sales_price_list"`
	PurchasePriceListId *uint                 `json:"purchase_price_list_id"`
	PurchasePriceList   *PurchasePriceListDTO `json:"purchase_price_list"`
	TransferPriceListId *uint                 `json:"transfer_price_list_id"`
	TransferPriceList   *TransferPriceListDTO `json:"transfer_price_list"`
}
type CurrencyDTO struct {
	Name           string `json:"name"`
	CurrencySymbol string `json:"currency_symbol"`
	CurrencyCode   string `json:"currency_code"`
	CurrencyId     string `json:"currency_id"`
	Exchangerate   uint   `json:"exchange_rate"`
}

type PurchasePriceListDTO struct {
	UpdatedByID              *uint                                                `json:"updated_by"`
	CreatedByID              *uint                                                `json:"created_by"`
	Price_list_id            uint                                                 `json:"price_list_id"`
	Vendor_name_id           *uint                                                `json:"vendor_name_id"`
	PurchaseLineItems        []PurchaseLineItemsDTO                               `json:"purchase_line_items"`
	PurchaseListOtherdetails shared_pricing_and_location.PurchaseListOtherdetails `json:"other_details"`
}
type PurchaseLineItemsDTO struct {
	Product_id           uint                   `json:"product_id"`
	MinimumOrderQuantity int                    `json:"minimum_order_quantity"`
	QuantityValueType    app_core.LookupCode    `json:"quantity_value_type"`
	QuantityValue        map[string]interface{} `json:"quantity_value"`
	Price                float64                `json:"price"`
	Price_quantity       float64                `json:"price_quantity"`
	Vendor_rate          int                    `json:"vendor_rate"`
	SalesPeriod          string                 `json:"sales_period"`
	CreditPeriod         string                 `json:"credit_period" `
	ExpectedDeliveryTime string                 `json:"expected_delivery_time" `
	LeadTime             string                 `json:"lead_time"`
}
type SalesPriceListDTO struct {
	CreatedByID            *uint               `json:"created_by"`
	UpdatedByID            *uint               `json:"updated_by"`
	Enter_Manually         *bool               `json:"enter_manually" `
	Add_channel_of_sale_id *uint               `json:"add_channel_of_sale_id" `
	Add_channel_of_sale    app_core.LookupCode `json:"add_channel_of_sale"`
	Currency_id            uint                `json:"currency_id"`
	Price_list_id          uint                `json:"price_list_id"`
	Item_rate_rule         int                 `json:"item_rate_rule"`
	Shipping_cost          float64             `json:"shipping_cost"`
	CustomerName           string              `json:"customer_name"`
	Type                   string              `json:"type"`
	Percentage             float64             `json:"percentage"`
	SalesLineItems         []SalesLineItemsDTO `json:"sales_line_items"`
}
type SalesLineItemsDTO struct {
	SPL_id     uint `json:"spl_id"`
	Product_id uint `json:"product_id"`
	//Product                ProductVariant   `json:"product_details" `
	CategoryCommission     int                    `json:"category_commission"`
	Uom_id                 *uint                  `json:"uom_id"`
	UOM                    UomDTO                 `json:"uom_details"`
	Quantity_value_type_id *uint                  `json:"quantity_value_type_id"`
	QuantityValueType      app_core.LookupCode    `json:"quantity_value_type"`
	QuantityValue          map[string]interface{} `json:"quantity_value"`
	Mrp                    int                    `json:"mrp"`
	SaleRate               int                    `json:"sale_rate"`
	Duties                 float64                `json:"duties"`
	PricingOptions         string                 `json:"pricing_options"`
	Pricing_options_id     *uint                  `json:"pricing_options_id"`
	Price                  int                    `json:"price"`
}
type UomDTO struct {
	Code             string                `json:"code"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	UomClassId       uint                  `json:"uom_class_id"`
	UomClassName     string                `json:"uom_class_name"`
	ConversionTypeId uint                  `json:"conversion_type_id"`
	ConversionType   model_core.Lookupcode `json:"conversion_type"`
	ConversionFactor float64               `json:"conversion_factor"`
}
type TransferPriceListDTO struct {
	CreatedByID              *uint                                       `json:"created_by"`
	UpdatedByID              *uint                                       `json:"updated_by"`
	Currency_id              uint                                        `json:"currency_id"`
	Price_list_id            uint                                        `json:"price_list_id"`
	ContractDetails          shared_pricing_and_location.ContractDetails `json:"contract_details"`
	TransferLineItems        []TransferLineItemsDTO                      `json:"transfer_line_items"`
	TransferListOtherDetails TransferListOtherDetailsDTO                 `json:"transfer_list_other_details"`
}
type TransferLineItemsDTO struct {
	TPL_id         uint    `json:"tpl_id"`
	Product_id     uint    `json:"product_id"`
	Price          float64 `json:"price"`
	Price_quantity float64 `json:"price_quantity"`
	Product_rate   float64 `json:"product_rate"`
}
type TransferListOtherDetailsDTO struct {
	FromAddressLocationId *uint                            `json:"from_address_location_id"`
	ToAddressLocationId   *uint                            `json:"to_address_location_id"`
	LocationFromAddress   location.LocationListResponseDTO `json:"location_from_address" `
	LocationToAddress     location.LocationListResponseDTO `json:"location_to_address" `
	SalesPeriod           string                           `json:"sales_period"`
	CreditPeriod          string                           `json:"credit_period" `
	ExpectedDeliveryTime  string                           `json:"expected_delivery_time" `
	LeadTime              string                           `json:"lead_time"`
	Shipping_cost         int                              `json:"shipping_cost" `
}

// Response for Create Price List
type (
	PricingCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name PricingCreateResponse
)

// Response for Price List Listing
type (
	PricingListResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []shared_pricing_and_location.Pricing
		}
	} //@ name PricingListResponse
)

// Response for Price List By Id
type (
	PricingListById struct {
		Pricing  shared_pricing_and_location.Pricing
		Sales    shared_pricing_and_location.SalesPriceList
		Purchase shared_pricing_and_location.PurchasePriceList
		Transfer shared_pricing_and_location.TransferPriceList
	}
	PricingListByIdResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PricingListById
		}
	} //@ name PricingListByIdResponse
)

// Response for Update Price List
type (
	PricingUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name PricingUpdateResponse
)

// Response for Delete Price List
type (
	PricingDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name PricingDeleteResponse
)

// Response for Delete Line Items
type (
	PricingDeleteLineResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name PricingDeleteLineResponse
)

// Response for Price List Listing
type (
	PurchasePricingListResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []shared_pricing_and_location.PurchasePriceList
		}
	} //@ name PurchasePricingListResponse
)
