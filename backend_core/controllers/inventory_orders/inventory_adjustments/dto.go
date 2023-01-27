package inventory_adjustments

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
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
	InvAdjNumber    string `json:"inv_adj_number,omitempty" query:"inv_adj_number"`
	AdjustmentType  string `json:"adjustment_type,omitempty" query:"adjustment_type"`
}

// Create and Update Request Payload for Inventory Adjustment
type (
	InventoryAdjustmentsRequest struct {
		ReasonID                    uint                       `json:"reason_id"`
		AdjustmentTypeID            uint                       `json:"adjustment_type_id"`
		ReferenceNumber             string                     `json:"reference_number"`
		AdjustmentDate              string                     `json:"adjustment_date"`
		WarehouseID                 uint                       `json:"warehouse_id"`
		PartnerID                   uint                       `json:"partner_id"`
		StatusID                    uint                       `json:"status_id"`
		StatusHistory               []map[string]interface{}   `json:"status_history"`
		InventoryAdjustmentLines    []InventoryAdjustmentLines `json:"inventory_adjustment_lines"`
		TotalChangeInInventory      uint                       `json:"total_change_in_inventory"`
		TotalChangeInInventoryCount uint                       `json:"total_change_in_inventory_count"`
		InternalNotes               string                     `json:"internal_notes"`
		ExternalNotes               string                     `json:"external_notes"`
		FileAttachId                []map[string]interface{}   `json:"file_attach_id"`
		ShippingOrderID             uint                       `json:"shipping_order_id"`
		ShippingAddress             []map[string]interface{}   `json:"shipping_address"`
		PartnerAddress              []map[string]interface{}   `json:"partner_address"`
		WarehouseAddress            []map[string]interface{}   `json:"warehouse_address"`
		SourceDocuments             map[string]interface{}     `json:"source_documents"`
		SourceDocumentTypeId        *uint                      `json:"source_document_type_id"`
		model_core.Model
	}

	InventoryAdjustmentLines struct {
		ProductID        uint    `json:"product_id"`
		ProductVariantID uint    `json:"product_variant_id"`
		Description      string  `json:"description"`
		StockInHand      uint    `json:"stock_in_hand"`
		AdjustedQuantity uint    `json:"adjusted_quantity"`
		AdjustedPrice    float32 `json:"adjusted_price"`
		BalanceQuantity  uint    `json:"balance_quantity"`
		UnitPrice        uint    `json:"unit_price"`
		model_core.Model
	}
)

// Create Inventory Adjustment response
type (
	InvAdjCreate struct {
		Created_id int
	}

	InvAdjCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data InvAdjCreate
		}
	} // @ name InvAdjCreateResponse
)

// Bulk Create Inventory Adjustments response
type (
	BulkInvAdjCreate struct {
		Created_id int
	}

	BulkInvAdjCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkInvAdjCreate
		}
	} // @ name BulkInvAdjCreateResponse
)

// Update Inventory Adjustment response
type (
	InvAdjUpdate struct {
		Updated_id int
	}

	InvAdjUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data InvAdjUpdate
		}
	} // @ name InvAdjUpdateResponse
)

// Delete Inventory Adjustment response
type (
	InvAdjDelete struct {
		Deleted_id int
	}
	InvAdjDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data InvAdjDelete
		} // @ name InvAdjDeleteResponse
	}
)

// Delete Inventory Adjustment Lines response
type (
	InvAdjLinesDelete struct {
		Deleted_id int
	}
	InvAdjLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data InvAdjLinesDelete
		} // @ name InvAdjLinesDeleteResponse
	}
)

// Get Inventory Adjustment response
type (
	InvAdjGet struct {
		inventory_orders.InventoryAdjustments
	}
	InvAdjGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data InvAdjGet
		}
	} // @ name InvAdjGetResponse
)

// Get all Inventory Adjustments response
type (
	InvAdjGetAll struct {
		model_core.Model
		ReasonID                    uint                                  `json:"reason_id"`
		Reason                      model_core.Lookupcode                 `json:"reason"`
		AdjustmentTypeID            uint                                  `json:"adjustment_type_id"`
		AdjustmentType              model_core.Lookupcode                 `json:"adjustment_type"`
		ReferenceNumber             string                                `json:"reference_number"`
		AdjustmentDate              time.Time                             `json:"adjustment_date"`
		WarehouseID                 uint                                  `json:"warehouse_id"`
		Warehouse                   shared_pricing_and_location.Locations `json:"warehouse"`
		PartnerID                   uint                                  `json:"partner_id"`
		Partner                     mdm.Partner                           `json:"partner"`
		StatusID                    uint                                  `json:"status_id"`
		Status                      model_core.Lookupcode                 `json:"status"`
		StatusHistory               map[string]interface{}                `json:"status_history"`
		InventoryAdjustmentLines    []InventoryAdjustmentOrderLines       `json:"inventory_adjustment_lines"`
		TotalChangeInInventory      uint                                  `json:"total_change_in_inventory"`
		TotalChangeInInventoryCount uint                                  `json:"total_change_in_inventory_count"`
		InternalNotes               string                                `json:"internal_notes"`
		ExternalNotes               string                                `json:"external_notes"`
		FileAttachId                map[string]interface{}                `json:"file_attach_id"`
		ShippingOrderID             uint                                  `json:"shipping_order_id"`
		ShippingOrder               shipping.ShippingOrder                `json:"shipping_order"`
		ShippingAddress             map[string]interface{}                `json:"shipping_address"`
		PartnerAddress              map[string]interface{}                `json:"partner_address"`
		WarehouseAddress            map[string]interface{}                `json:"warehouse_address"`
	}

	InventoryAdjustmentOrderLines struct {
		ProductID        uint                `json:"product_id"`
		Product          mdm.ProductTemplate `json:"product"`
		ProductVariantID uint                `json:"product_variant_id"`
		ProductVariant   mdm.ProductVariant  `json:"product_variant"`
		Description      string              `json:"description"`
		StockInHand      uint                `json:"stock_in_hand"`
		AdjustedQuantity uint                `json:"adjusted_quantity"`
		AdjustedPrice    float32             `json:"adjusted_price"`
		BalanceQuantity  uint                `json:"balance_quantity"`
	}

	InvAdjPaginatedResponse struct {
		Total_pages   uint `json:"total_pages"`
		Per_page      uint `json:"per_page"`
		Current_page  uint `json:"current_page"`
		Next_page     uint `json:"next_page"`
		Previous_page uint `json:"previous_page"`
		Total_rows    uint `json:"total_rows"`
	}
	InvAdjGetAllResponse struct {
		Body struct {
			Meta       response.MetaResponse
			Data       []InvAdjGetAll
			Pagination InvAdjPaginatedResponse
		}
	} // @ name InvAdjGetAllResponse
)

// Search Inventory Adjustments response
type (
	InvAdjSearch struct {
		ReasonID                    uint                                  `json:"reason_id"`
		Reason                      model_core.Lookupcode                 `json:"reason"`
		AdjustmentTypeID            uint                                  `json:"adjustment_type_id"`
		AdjustmentType              model_core.Lookupcode                 `json:"adjustment_type"`
		ReferenceNumber             string                                `json:"reference_number"`
		AdjustmentDate              time.Time                             `json:"adjustment_date"`
		WarehouseID                 uint                                  `json:"warehouse_id"`
		Warehouse                   shared_pricing_and_location.Locations `json:"warehouse"`
		PartnerID                   uint                                  `json:"partner_id"`
		Partner                     mdm.Partner                           `json:"partner"`
		StatusID                    uint                                  `json:"status_id"`
		Status                      model_core.Lookupcode                 `json:"status"`
		StatusHistory               map[string]interface{}                `json:"status_history"`
		InventoryAdjustmentLines    []InventoryAdjustmentOrderLines       `json:"inventory_adjustment_lines"`
		TotalChangeInInventory      uint                                  `json:"total_change_in_inventory"`
		TotalChangeInInventoryCount uint                                  `json:"total_change_in_inventory_count"`
		InternalNotes               string                                `json:"internal_notes"`
		ExternalNotes               string                                `json:"external_notes"`
		FileAttachId                map[string]interface{}                `json:"file_attach_id"`
		ShippingOrderID             uint                                  `json:"shipping_order_id"`
		ShippingOrder               shipping.ShippingOrder                `json:"shipping_order"`
		ShippingAddress             map[string]interface{}                `json:"shipping_address"`
		PartnerAddress              map[string]interface{}                `json:"partner_address"`
		WarehouseAddress            map[string]interface{}                `json:"warehouse_address"`
	}

	InvAdjSearchResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []InvAdjSearch
		}
	} // @ name InvAdjSearchResponse
)

// Send Mail Inventory Adjustment response
type (
	SendMailInvAdj struct {
		ID            string `query:"id" json:"id"`
		ReceiverEmail string `query:"receiver_email" json:"receiver_email"`
	}

	SendMailInvAdjResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SendMailInvAdj
		}
	}
)

// Download Pdf Inventory Adjustment response
type (
	DownloadPdfInvAdj struct {
		ID string `query:"id"`
	}

	DownloadPdfInvAdjResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DownloadPdfInvAdj
		}
	}
)
