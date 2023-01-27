package pick_list

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_tasks"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
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
	PickListNumber  string `json:"pick_list_number,omitempty" query:"pick_list_number"`
	ReferenceNumber string `json:"reference_number,omitempty" query:"reference_number"`
}

// Create and Update Request Payload for Picklist
type (
	PicklistRequest struct {
		PickListNumber           string                                `json:"pick_list_number"`
		AutoCreatePicklistNumber bool                                  `json:"auto_create_picklist_number"`
		ReferenceNumber          string                                `json:"reference_number"`
		SourceDocumentTypeID     uint                                  `json:"source_document_type_id"`
		SourceDocumentType       model_core.Lookupcode                 `json:"source_document_type"`
		AssigneeToID             uint                                  `json:"assignee_to_id"`
		AssigneeTo               mdm.Partner                           `json:"assignee_to"`
		SelectCustomerId         uint                                  `json:"select_customer_id"`
		PartnerID                uint                                  `json:"partner_id"`
		WarehouseID              uint                                  `json:"warehouse_id"`
		Warehouse                shared_pricing_and_location.Locations `json:"warehouse"`
		StatusID                 uint                                  `json:"status_id"`
		Status                   model_core.Lookupcode                 `json:"status"`
		StatusHistory            []map[string]interface{}              `json:"status_history"`
		PicklistLines            []PickListLines                       `json:"picklist_lines"`
		InternalNotes            string                                `json:"internal_notes"`
		ExternalNotes            string                                `json:"external_notes"`
		AttachmentFiles          []map[string]interface{}              `json:"attachment_files"`
		StartDateTime            time.Time                             `json:"start_date_time"`
		EndDateTime              time.Time                             `json:"end_date_time"`
		TotalItemsToPick         uint                                  `json:"items_to_pick"`
		TotalPickedItems         uint                                  `json:"total_picked_items"`
		PriceListID              uint                                  `json:"price_list_id"`
		SourceDocuments          map[string]interface{}                `json:"source_documents"`
		model_core.Model
	}

	PickListLines struct {
		ProductID         uint                     `json:"product_id"`
		ProductVariantID  uint                     `json:"product_variant_id"`
		SalesDocumentID   []map[string]interface{} `json:"sales_document_id"`
		PartnerID         uint                     `json:"partner_id"`
		QuantityOrdered   uint                     `json:"quantity_ordered"`
		QuantityToPick    uint                     `json:"quantity_to_pick"`
		QuantityPicked    uint                     `json:"quantity_picked"`
		RemainingQuantity uint                     `json:"remaining_quantity"`
		CustomerName      string                   `json:"customer_name"`
		model_core.Model
	}
)

// Create Picklist response
type (
	PicklistCreate struct {
		Created_id int
	}

	PicklistCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PicklistCreate
		}
	} // @ name PicklistCreateResponse
)

// Bulk Create Picklist response
type (
	BulkPicklistCreate struct {
		Created_id int
	}

	BulkPicklistCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkPicklistCreate
		}
	} // @ name BulkPicklistCreateResponse
)

// Update Picklist response
type (
	PicklistUpdate struct {
		Updated_id int
	}

	PicklistUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PicklistUpdate
		}
	} // @ name PicklistUpdateResponse
)

// Delete Picklist response
type (
	PicklistDelete struct {
		Deleted_id int
	}
	PicklistDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PicklistDelete
		} // @ name PicklistDeleteResponse
	}
)

// Delete Picklist Lines response
type (
	PicklistLinesDelete struct {
		Deleted_id int
	}
	PicklistLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PicklistLinesDelete
		} // @ name PicklistLinesDeleteResponse
	}
)

// Get Picklist response
type (
	PicklistGet struct {
		inventory_tasks.PickList
	}
	PicklistGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data PicklistGet
		}
	} // @ name PicklistGetResponse
)

// Get all Picklist response
type (
	PicklistGetAll struct {
		PickListNumber           string                                `json:"pick_list_number"`
		AutoCreatePicklistNumber bool                                  `json:"auto_create_picklist_number"`
		ReferenceNumber          string                                `json:"reference_number"`
		SourceDocumentTypeID     uint                                  `json:"source_document_type_id"`
		SourceDocumentType       model_core.Lookupcode                 `json:"source_doc_type"`
		SourceDocumentIds        pq.Int64Array                         `json:"source_document_ids"`
		PartnerID                uint                                  `json:"partner_id"`
		Partner                  mdm.Partner                           `json:"partner"`
		WarehouseID              uint                                  `json:"warehouse_id"`
		Warehouse                shared_pricing_and_location.Locations `json:"warehouse"`
		StatusID                 uint                                  `json:"status_id"`
		Status                   model_core.Lookupcode                 `json:"status"`
		StatusHistory            map[string]interface{}                `json:"status_history"`
		PicklistLines            []PickListOrderLines                  `json:"picklist_lines"`
		InternalNotes            string                                `json:"internal_notes"`
		ExternalNotes            string                                `json:"external_notes"`
		AttachmentFiles          map[string]interface{}                `json:"attachment_files"`
		StartDateTime            time.Time                             `json:"start_date_time"`
		EndDateTime              time.Time                             `json:"end_date_time"`
		TotalItemsToPick         uint                                  `json:"items_to_pick"`
		TotalPickedItems         uint                                  `json:"total_picked_items"`
		PriceListID              uint                                  `json:"price_list_id"`
		PriceList                shared_pricing_and_location.Pricing   `json:"price_list"`
	}

	PickListOrderLines struct {
		PickList_Id       uint                   `json:"-"`
		ProductID         uint                   `json:"product_id"`
		Product           mdm.ProductTemplate    `json:"product"`
		ProductVariantID  uint                   `json:"product_variant_id"`
		ProductVariant    mdm.ProductVariant     `json:"product_variant"`
		SalesDocumentID   map[string]interface{} `json:"sales_document_id"`
		PartnerID         uint                   `json:"partner_id"`
		Partner           mdm.Partner            `json:"partner"`
		QuantityOrdered   uint                   `json:"quantity_ordered"`
		QuantityToPick    uint                   `json:"quantity_to_pick"`
		QuantityPicked    uint                   `json:"quantity_picked"`
		RemainingQuantity uint                   `json:"remaining_quantity"`
	}

	PicklistPaginatedResponse struct {
		Total_pages   uint `json:"total_pages"`
		Per_page      uint `json:"per_page"`
		Current_page  uint `json:"current_page"`
		Next_page     uint `json:"next_page"`
		Previous_page uint `json:"previous_page"`
		Total_rows    uint `json:"total_rows"`
	}
	PicklistGetAllResponse struct {
		Body struct {
			Meta       response.MetaResponse
			Data       []PicklistGetAll
			Pagination PicklistPaginatedResponse
		}
	} // @ name PicklistGetAllResponse
)

// Search Picklist response
type (
	PicklistSearch struct {
		PickListNumber           string                                `json:"pick_list_number"`
		AutoCreatePicklistNumber bool                                  `json:"auto_create_picklist_number"`
		ReferenceNumber          string                                `json:"reference_number"`
		SourceDocumentTypeID     uint                                  `json:"source_document_type_id"`
		SourceDocumentType       model_core.Lookupcode                 `json:"source_doc_type"`
		SourceDocumentIds        pq.Int64Array                         `json:"source_document_ids"`
		PartnerID                uint                                  `json:"partner_id"`
		Partner                  mdm.Partner                           `json:"partner"`
		WarehouseID              uint                                  `json:"warehouse_id"`
		Warehouse                shared_pricing_and_location.Locations `json:"warehouse"`
		StatusID                 uint                                  `json:"status_id"`
		Status                   model_core.Lookupcode                 `json:"status"`
		StatusHistory            map[string]interface{}                `json:"status_history"`
		PicklistLines            []PickListOrderLines                  `json:"picklist_lines"`
		InternalNotes            string                                `json:"internal_notes"`
		ExternalNotes            string                                `json:"external_notes"`
		AttachmentFiles          map[string]interface{}                `json:"attachment_files"`
		StartDateTime            time.Time                             `json:"start_date_time"`
		EndDateTime              time.Time                             `json:"end_date_time"`
		TotalItemsToPick         uint                                  `json:"items_to_pick"`
		TotalPickedItems         uint                                  `json:"total_picked_items"`
		PriceListID              uint                                  `json:"price_list_id"`
		PriceList                shared_pricing_and_location.Pricing   `json:"price_list"`
	}

	PicklistSearchResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []PicklistSearch
		}
	} // @ name PicklistSearchResponse
)

// Send Mail PickList response
type (
	SendMailPickList struct {
		ID            string `query:"id" json:"id"`
		ReceiverEmail string `query:"receiver_email" json:"receiver_email"`
	}

	SendMailPickListResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SendMailPickList
		}
	}
)

// Download Pdf Cycle Count response
type (
	DownloadPdfPickList struct {
		ID string `query:"id"`
	}

	DownloadPdfPickListResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DownloadPdfPickList
		}
	}
)
