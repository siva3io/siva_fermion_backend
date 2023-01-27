package scrap_orders

import (
	core_dto "fermion/backend_core/controllers/cores"
	locations "fermion/backend_core/controllers/mdm/locations"
	products "fermion/backend_core/controllers/mdm/products"
	mdm_UOM "fermion/backend_core/controllers/mdm/uom"
	"fermion/backend_core/controllers/shipping/shipping_orders"
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
type SearchQuery struct {
	Scrap_order_number string `json:"scrap_order_number,omitempty" query:"scrap_order_number"`
	Scrapingstatus     string `json:"scraping_status,omitempty" query:"status"`
	Grn_status         bool   `json:"grn_status,omitempty" query:"grn_status"`
}

// Create and update ScrapOrder RequestPayload

type (
	ScrapOrders struct {
		ID                          uint                                 `json:"id"`
		CreatedByID                 *uint                                `json:"created_by"`
		UpdatedByID                 *uint                                `json:"updated_by"`
		DeletedByID                 *uint                                `json:"deleted_by"`
		CompanyId                   uint                                 `json:"company_id"`
		Scrap_order_number          string                               `json:"scrap_order_no"`
		Schedule_scrap_date         string                               `json:"schedule_scrap_date"`
		AutoCreateScrapNumber       bool                                 `json:"auto_create_scrap_number"`
		AutoGenerateReferenceNumber bool                                 `json:"auto_generate_reference_number"`
		Scrap_reason_id             uint                                 `json:"scrap_reason_id"`
		Scrap_source_location_id    uint                                 `json:"scrap_source_location_id"`
		Scrap_location_id           uint                                 `json:"scrap_location_id"`
		Reference_id                string                               `json:"reference_id"`
		No_of_items                 int                                  `json:"no_of_items"`
		Total_quantity              int                                  `json:"total_quantity,omitempty"`
		Order_lines                 []ScrapOrderLines                    `json:"order_lines"`
		ScrapLocationDetails        map[string]interface{}               `json:"scrap_location_details"`
		PickupDateAndTime           SO_PickupDateAndTime                 `json:"pickup_date_time"`
		Scraping_status             string                               `json:"scraping_status"`
		Grn_status                  *bool                                `json:"grn_status"`
		GrnStatusId                 *uint                                `json:"grn_status_id"`
		GrnStatus                   model_core.Lookupcode                `json:"Grn_status"`
		SourceDocuments             map[string]interface{}               `json:"source_documents"`
		SourceDocumentTypeId        *uint                                `json:"source_document_type_id"`
		Status_id                   *uint                                `json:"status_id"`
		Status                      model_core.Lookupcode                `json:"status"`
		Status_history              map[string]interface{}               `json:"status_history"`
		ShippingOrderId             *uint                                `json:"shipping_order_id" `
		Shipping_details            shipping_orders.ShippingOrderRequest `json:"shipping_details"`
		IsShipping                  *bool                                `json:"is_shipping"`
		ExpectedShippingDate        string                               `json:"expected_shipping_date"`
		Shipping_mode_id            *uint                                `json:"shipping_mode_id"`
		Shipping_Mode               model_core.Lookupcode                `json:"shipping_mode"`
		//ShippingOrder               core_dto.LookupCodesDTO `json:"shipping_order"`
	}

	SO_PickupDateAndTime struct {
		PickupDate     string `json:"pickup_date"`
		PickupFromTime string `json:"pickup_from_time"`
		PickupToTime   string `json:"pickup_to_time"`
	}

	ScrapOrderLines struct {
		ID                  uint    `json:"id"`
		Scrap_id            uint    `json:"scrap_id"`
		Product_id          uint    `json:"product_id"`
		Lot_number          int     `json:"lot_number"`
		InventoryId         uint64  `json:"inventory_id"`
		Uom_id              uint    `json:"uom_id"`
		Scrap_item_quantity int     `json:"scrap_item_quantity"`
		Price               float64 `json:"price"`
	}
)

// Create Scrap Order response
type (
	// SOCreate struct {
	// 	Created_id int
	// }
	ScrapOrderCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name ScrapOrderCreateResponse
)

// Update Scrap Order response
type (
	// SOUpdate struct {
	// 	Updated_id int
	// }
	ScrapOrderUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name ScrapOrderUpdateResponse
)

// Get ScrapOrder response
type (
	SOGet struct {
		orders.ScrapOrders
	}
	ScrapOrderGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SOGet
		}
	} //@ name ScrapOrderGetResponse
)

// Get all ScrapOrders response
type (
	ScrapOrdersResponseDTO struct {
		locations.ModelDto
		AppId                    *uint                         `json:"app_id" `
		CompanyId                uint                          `json:"company_id"`
		Scrap_order_number       string                        `json:"scrap_order_no"`
		Schedule_scrap_date      string                        `json:"schedule_scrap_date"`
		Scrap_reason_id          uint                          `json:"scrap_reason_id"`
		Scrap_reason             core_dto.LookupCodeDTO        `json:"scrap_reason"`
		Scrap_source_location    locations.LocationResponseDTO `json:"scrap_source_location"`
		Scrap_source_location_id uint                          `json:"scrap_source_location_id"`
		Scrap_location_id        uint                          `json:"scrap_location_id"`
		Scrap_location           locations.LocationResponseDTO `json:"scrap_location"`
		Reference_id             string                        `json:"reference_id"`
		No_of_items              int                           `json:"no_of_items"`
		Total_quantity           int                           `json:"total_quantity,omitempty"`
		Order_lines              []ScrapOrderLinesResponseDTO  `json:"order_lines"`
		ScrapLocationDetails     map[string]interface{}        `json:"scrap_location_details"`
		PickupDateAndTime        SO_PickupDateAndTime          `json:"pickup_date_time"`
		Scraping_status          string                        `json:"scraping_status"`
		Grn_status               *bool                         `json:"grn_status"`
		SourceDocumentType       core_dto.LookupCodeDTO        `json:"source_document"`
		SourceDocuments          map[string]interface{}        `json:"source_documents"`
		SourceDocumentTypeId     *uint                         `json:"source_document_type_id"`
		Status_id                *uint                         `json:"status_id"`
		Status                   core_dto.LookupCodeDTO        `json:"status"`
		Status_history           map[string]interface{}        `json:"status_history"`
		ShippingOrderId          *uint                         `json:"shipping_order_id" `
		IsShipping               *bool                         `json:"is_shipping"`
		ExpectedShippingDate     string                        `json:"expected_shipping_date"`
		Shipping_mode_id         *uint                         `json:"shipping_mode_id"`
		Shipping_Mode            core_dto.LookupCodeDTO        `json:"shipping_mode"`
		GrnStatusId              *uint                         `json:"grn_status_id"`
		GrnStatus                core_dto.LookupCodeDTO        `json:"Grn_status"`
		//ShippingOrder        core_dto.LookupCodesDTO               `json:"shipping_order"`
	}
	ScrapOrderLinesResponseDTO struct {
		ID                  uint                        `json:"id"`
		Scrap_id            uint                        `json:"scrap_id"`
		Product_id          uint                        `json:"product_id"`
		Product             products.VariantResponseDTO `json:"product_Details"`
		Lot_number          int                         `json:"lot_number"`
		InventoryId         uint64                      `json:"inventory_id"`
		Uom_id              uint                        `json:"uom_id"`
		UOM                 mdm_UOM.UomResponseDTO      `json:"uom"`
		Scrap_item_quantity int                         `json:"scrap_item_quantity"`
		Price               float64                     `json:"price"`
	}

	ScrapOrderGetAllResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ScrapOrdersResponseDTO
		}
	} //@ name ScrapOrderGetAllResponse
)

// Delete Scrap Order response
type (
	// SODelete struct {
	// 	Deleted_id int
	// }
	ScrapOrderDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name ScrapOrderDeleteResponse
)

// Delete Scrap Orderlines response
type (
	// SOLDelete struct {
	// 	Deleted_id int
	// 	Product_id int
	// }
	ScrapOrderLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name ScrapOrderLinesDeleteResponse
)

// Bulk Create Scrap order response
type (
	BulkScrapOrderCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} // @ name BulkScrapOrderCreateResponse
)
