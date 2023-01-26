package asn

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/shipping"
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
	ReferenceNumber string `json:"reference_number,omitempty" query:"reference_number"`
	Asn_number      string `json:"asn_number,omitempty" query:"asn_number"`
	CustomerName    string `json:"customer_name,omitempty" query:"customer_name"`
	Status          string `json:"status,omitempty" query:"status"`
	Payment_type    string `json:"payment_type,omitempty" query:"payment_type"`
}

// Create and Update Request Payload for Asn
type (
	AsnRequest struct {
		AsnNumber                   string                   `json:"asn_number"`
		AutoCreateAsnNumber         bool                     `json:"auto_create_asn_number"`
		ReferenceNumber             string                   `json:"reference_number"`
		AutoGenerateReferenceNumber bool                     `json:"auto_generate_reference_number"`
		WarehouseID                 uint                     `json:"warehouse_id"`
		TotalQuantity               uint                     `json:"total_quantity"`
		StatusID                    uint                     `json:"status_id"`
		StatusHistory               []map[string]interface{} `json:"status_history"`
		GrnIDs                      []int64                  `json:"grn_ids"`
		ShippingModeID              uint                     `json:"shipping_mode_id"`
		LinkPurchaseOrderID         uint                     `json:"link_purchase_order_id"`
		ScrapOrderID                uint                     `json:"scrap_order_id"`
		SourceDocumentTypeID        uint                     `json:"source_document_type_id"`
		SourceDocuments             map[string]interface{}   `json:"source_document"`
		AsnOrderLines               []AsnLines               `json:"asn_order_lines"`
		DispatchLocationDetails     map[string]interface{}   `json:"dispatch_location_details"`
		DestinationLocationDetails  map[string]interface{}   `json:"destination_location_details"`
		ShippingDetails             map[string]interface{}   `json:"shipping_details"`
		EunimartWalletAmount        float32                  `json:"eunimart_wallet_amount"`
		SchedulePickupDate          string                   `json:"schedule_pickup_date"`
		SchedulePickupFromTime      string                   `json:"schedule_pickup_from_time"`
		SchedulePickupToTime        string                   `json:"schedule_pickup_to_time"`
		PartnerID                   uint                     `json:"partner_id"`
		model_core.Model
	}

	AsnLines struct {
		ProductID           uint                     `json:"product_id"`
		ProductVariantID    uint                     `json:"product_variant_id"`
		PackageTypeID       uint                     `json:"package_type_id"`
		UnitPerBox          uint                     `json:"unit_per_box"`
		UomID               uint                     `json:"uom_id"`
		PackageLength       float32                  `json:"package_length"`
		PackageWidth        float32                  `json:"package_width"`
		PackageHeight       float32                  `json:"package_height"`
		PackageWeight       float32                  `json:"package_weight"`
		NumberOfBoxes       uint                     `json:"no_of_boxes"`
		CrossDockingRequest []map[string]interface{} `json:"cross_docking_req"`
		model_core.Model
	}
)

// Create Asn response
type (
	AsnCreate struct {
		Created_id int
	}

	AsnCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data AsnCreate
		}
	} // @ name AsnCreateResponse
)

// Bulk Create Asn response
type (
	BulkAsnCreate struct {
		Created_id int
	}

	BulkAsnCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkAsnCreate
		}
	} // @ name BulkAsnCreateResponse
)

// Update Asn response
type (
	AsnUpdate struct {
		Updated_id int
	}

	AsnUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data AsnUpdate
		}
	} // @ name AsnUpdateResponse
)

// Delete Asn response
type (
	AsnDelete struct {
		Deleted_id int
	}
	AsnDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data AsnDelete
		} // @ name AsnDeleteResponse
	}
)

// Delete Asn Lines response
type (
	AsnLinesDelete struct {
		Deleted_id int
	}
	AsnLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data AsnLinesDelete
		} // @ name AsnLinesDeleteResponse
	}
)

// Get Asn response
type (
	AsnGet struct {
		inventory_orders.ASN
	}
	AsnGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data AsnGet
		}
	} // @ name AsnGetResponse
)

// Get all Asn response
type (
	AsnGetAll struct {
		AsnNumber                   string                                `json:"asn_number"`
		AutoCreateAsnNumber         bool                                  `json:"auto_create_asn_number"`
		ReferenceNumber             string                                `json:"reference_number"`
		AutoGenerateReferenceNumber bool                                  `json:"auto_generate_reference_number"`
		WarehouseID                 uint                                  `json:"warehouse_id"`
		Warehouse                   shared_pricing_and_location.Locations `json:"warehouse"`
		TotalQuantity               uint                                  `json:"total_quantity"`
		StatusID                    uint                                  `json:"status_id"`
		Status                      model_core.Lookupcode                 `json:"status"`
		StatusHistory               map[string]interface{}                `json:"status_history"`
		ShippingModeID              uint                                  `json:"shipping_mode_id"`
		ShippingMode                model_core.Lookupcode                 `json:"shipping_mode"`
		LinkPurchaseOrderID         uint                                  `json:"link_purchase_order_id"`
		LinkPurchaseOrder           orders.PurchaseOrders                 `json:"link_po"`
		SourceDocumentTypeID        uint                                  `json:"source_document_type_id"`
		SourceDocumentType          model_core.Lookupcode                 `json:"source_document_type"`
		SourceDocumentID            uint                                  `json:"source_document_id"`
		AsnOrderLines               []AsnOrderLines                       `json:"asn_order_lines"`
		DispatchLocationID          uint                                  `json:"dispatch_location_id"`
		DispatchLocation            shared_pricing_and_location.Locations `json:"dispatch_location"`
		DestinationLocationID       uint                                  `json:"destination_location_id"`
		DestinationLocation         shared_pricing_and_location.Locations `json:"destination_location"`
		ShippingID                  uint                                  `json:"shipping_id"`
		Shipping                    shipping.ShippingOrder                `json:"shipping"`
		EunimartWalletAmount        float32                               `json:"eunimart_wallet_amount"`
		SchedulePickupDateTime      time.Time                             `json:"schedule_pickup_date_time"`
		PartnerID                   uint                                  `json:"partner_id"`
		Partner                     mdm.Partner                           `json:"partner"`
		GrnIDs                      map[string]interface{}                `json:"grn_id"`
		ScrapOrderID                uint                                  `json:"scrap_order_id"`
	}

	AsnOrderLines struct {
		ProductID           uint                   `json:"product_id"`
		Product             mdm.ProductTemplate    `json:"product"`
		ProductVariantID    uint                   `json:"product_variant_id"`
		ProductVariant      mdm.ProductVariant     `json:"product_variant"`
		PackageTypeID       uint                   `json:"package_type_id"`
		PackageType         model_core.Lookupcode  `json:"package_type"`
		UnitPerBox          uint                   `json:"unit_per_box"`
		UomID               uint                   `json:"uom_id"`
		UOM                 mdm.Uom                `json:"uom"`
		PackageLength       float32                `json:"package_length"`
		PackageWidth        float32                `json:"package_width"`
		PackageHeight       float32                `json:"package_height"`
		PackageWeight       float32                `json:"package_weight"`
		NumberOfBoxes       uint                   `json:"no_of_boxes"`
		CrossDockingRequest map[string]interface{} `json:"cross_docking_req"`
	}

	AsnPaginatedResponse struct {
		Total_pages   uint `json:"total_pages"`
		Per_page      uint `json:"per_page"`
		Current_page  uint `json:"current_page"`
		Next_page     uint `json:"next_page"`
		Previous_page uint `json:"previous_page"`
		Total_rows    uint `json:"total_rows"`
	}
	AsnGetAllResponse struct {
		Body struct {
			Meta       response.MetaResponse
			Data       []AsnGetAll
			Pagination AsnPaginatedResponse
		}
	} // @ name AsnGetAllResponse
)

// Search Asn response
type (
	AsnSearch struct {
		AsnNumber                   string                                `json:"asn_number"`
		AutoCreateAsnNumber         bool                                  `json:"auto_create_asn_number"`
		ReferenceNumber             string                                `json:"reference_number"`
		AutoGenerateReferenceNumber bool                                  `json:"auto_generate_reference_number"`
		WarehouseID                 uint                                  `json:"warehouse_id"`
		Warehouse                   shared_pricing_and_location.Locations `json:"warehouse"`
		TotalQuantity               uint                                  `json:"total_quantity"`
		StatusID                    uint                                  `json:"status_id"`
		Status                      model_core.Lookupcode                 `json:"status"`
		StatusHistory               map[string]interface{}                `json:"status_history"`
		ShippingModeID              uint                                  `json:"shipping_mode_id"`
		ShippingMode                model_core.Lookupcode                 `json:"shipping_mode"`
		LinkPurchaseOrderID         uint                                  `json:"link_purchase_order_id"`
		LinkPurchaseOrder           orders.PurchaseOrders                 `json:"link_po"`
		SourceDocumentTypeID        uint                                  `json:"source_document_type_id"`
		SourceDocumentType          model_core.Lookupcode                 `json:"source_document_type"`
		SourceDocumentID            uint                                  `json:"source_document_id"`
		AsnOrderLines               []AsnLines                            `json:"asn_order_lines"`
		DispatchLocationID          uint                                  `json:"dispatch_location_id"`
		DispatchLocation            shared_pricing_and_location.Locations `json:"dispatch_location"`
		DestinationLocationID       uint                                  `json:"destination_location_id"`
		DestinationLocation         shared_pricing_and_location.Locations `json:"destination_location"`
		ShippingID                  uint                                  `json:"shipping_id"`
		Shipping                    shipping.ShippingOrder                `json:"shipping"`
		EunimartWalletAmount        float32                               `json:"eunimart_wallet_amount"`
		SchedulePickupDateTime      time.Time                             `json:"schedule_pickup_date_time"`
		PartnerID                   uint                                  `json:"partner_id"`
		Partner                     mdm.Partner                           `json:"partner"`
		GrnIDs                      map[string]interface{}                `json:"grn_id"`
		ScrapOrderID                uint                                  `json:"scrap_order_id"`
	}

	AsnSearchResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []AsnSearch
		}
	} // @ name AsnSearchResponse
)

// Send Mail Asn response
type (
	SendMailAsn struct {
		ID            string `query:"id" json:"id"`
		ReceiverEmail string `query:"receiver_email" json:"receiver_email"`
	}

	SendMailAsnResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SendMailAsn
		}
	}
)

// Download Pdf Asn response
type (
	DownloadPdfAsn struct {
		ID string `query:"id"`
	}

	DownloadPdfAsnResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DownloadPdfAsn
		}
	}
)
