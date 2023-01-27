package cycle_count

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
	Name         string `json:"name,omitempty" query:"name"`
	Status       string `json:"status,omitempty" query:"status"`
	CountMethod  string `json:"count_method,omitempty" query:"count_method"`
	ScheduleDate string `json:"schedule_date,omitempty" query:"schedule_date"`
}

// Create and Update Request Payload for Cycle Count
type (
	CycleCountRequest struct {
		CycleCountNumber          string                   `json:"cycle_count_number"`
		CycleCountDate            time.Time                `json:"cycle_count_date"`
		WarehouseID               uint                     `json:"warehouse_id"`
		PartnerID                 uint                     `json:"partner_id"`
		StatusID                  uint                     `json:"status_id"`
		StatusHistory             []map[string]interface{} `json:"status_history"`
		ItemsCount                uint                     `json:"items_count"`
		CountMethodID             uint                     `json:"count_method_id"`
		OrderLines                []CycleCountLines        `json:"order_lines"`
		ScheduleDate              time.Time                `json:"schedule_date"`
		ThresholdLimit            float64                  `json:"threshold_limit"`
		AutoScheduleCycleCount    bool                     `json:"auto_schedule_cycle_count"`
		FrequencyID               uint                     `json:"frequency_id"`
		CreateInventoryAdjustment bool                     `json:"create_inventory_adjustment"`
		model_core.Model
	}

	CycleCountLines struct {
		ProductID                   uint                     `json:"product_id"`
		ProductVariantID            uint                     `json:"product_variant_id"`
		BinLocationID               []map[string]interface{} `json:"bin_location_id"`
		InventoryCount              uint                     `json:"inventory_count"`
		LocationDetails             []map[string]interface{} `json:"location_details"`
		LocationSpaceTypeID         uint                     `json:"location_space_type_id"`
		LocationInputTypeID         uint                     `json:"location_input_type_id"`
		RowNumber                   uint                     `json:"row_number"`
		RackNumber                  uint                     `json:"rack_number"`
		ShelfNumber                 uint                     `json:"shelf_number"`
		BinNumber                   uint                     `json:"bin_number"`
		TotalBinCount               uint                     `json:"total_bin_count"`
		RecordedQuantity            uint                     `json:"recorded_quantity"`
		CurrentCount                uint                     `json:"current_count"`
		DiscrepancyCount            uint                     `json:"discrepancy_count"`
		DiscrepancyLevel            float32                  `json:"discrepancy_level"`
		CycleCountDiscrepancyReason string                   `json:"cycle_count_discrepancy_reason"`
		model_core.Model
	}
)

// Create Cycle Count response
type (
	CycleCountCreate struct {
		Created_id int
	}

	CycleCountCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CycleCountCreate
		}
	} // @ name CycleCountCreateResponse
)

// Bulk Create Cycle Count response
type (
	BulkCycleCountCreate struct {
		Created_id int
	}

	BulkCycleCountCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkCycleCountCreate
		}
	} // @ name BulkCycleCountCreateResponse
)

// Update Cycle Count response
type (
	CycleCountUpdate struct {
		Updated_id int
	}

	CycleCountUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CycleCountUpdate
		}
	} // @ name CycleCountUpdateResponse
)

// Delete Cycle Count response
type (
	CycleCountDelete struct {
		Deleted_id int
	}
	CycleCountDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CycleCountDelete
		} // @ name CycleCountDeleteResponse
	}
)

// Delete Cycle Count OrderLines response
type (
	CycleCountLinesDelete struct {
		Deleted_id int
	}
	CycleCountLinesDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CycleCountLinesDelete
		} // @ name CycleCountLinesDeleteResponse
	}
)

// Get Cycle Count response
type (
	CycleCountGet struct {
		inventory_tasks.CycleCount
	}
	CycleCountGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CycleCountGet
		}
	} // @ name CycleCountGetResponse
)

// Get all Cycle Count response
type (
	CycleCountGetAll struct {
		model_core.Model
		CycleCountNumber          string                                `json:"cycle_count_number"`
		WarehouseID               uint                                  `json:"warehouse_id"`
		Warehouse                 shared_pricing_and_location.Locations `json:"warehouse"`
		PartnerID                 uint                                  `json:"partner_id"`
		Partner                   mdm.Partner                           `json:"partner"`
		StatusID                  uint                                  `json:"status_id"`
		Status                    model_core.Lookupcode                 `json:"status"`
		StatusHistory             map[string]interface{}                `json:"status_history"`
		ItemsCount                uint                                  `json:"items_count"`
		CountMethodID             uint                                  `json:"count_method_id"`
		CountMethod               model_core.Lookupcode                 `json:"count_method"`
		OrderLines                []CycleCountLines                     `json:"order_lines"`
		ScheduleDate              time.Time                             `json:"schedule_date"`
		ThresholdLimit            float64                               `json:"threshold_limit"`
		AutoScheduleCycleCount    bool                                  `json:"auto_schedule_cycle_count"`
		FrequencyID               uint                                  `json:"frequency_id"`
		Frequency                 model_core.Lookupcode                 `json:"frequency"`
		CreateInventoryAdjustment bool                                  `json:"create_inventory_adjustment"`
	}

	CycleCountOrderLines struct {
		Cycle_count_id              uint                   `json:"-"`
		ProductID                   uint                   `json:"product_id"`
		Product                     mdm.ProductTemplate    `json:"product"`
		ProductVariantID            uint                   `json:"product_variant_id"`
		ProductVariant              mdm.ProductVariant     `json:"product_variant"`
		BinLocationID               map[string]interface{} `json:"bin_location_id"`
		InventoryCount              uint                   `json:"inventory_count"`
		LocationDetails             map[string]interface{} `json:"location_details"`
		LocationSpaceTypeID         uint                   `json:"location_space_type_id"`
		LocationSpaceType           model_core.Lookupcode  `json:"location_space_type"`
		LocationInputTypeID         uint                   `json:"location_input_type_id"`
		LocationInputType           model_core.Lookupcode  `json:"location_input_type"`
		RowNumber                   uint                   `json:"row_number"`
		RackNumber                  uint                   `json:"rack_number"`
		ShelfNumber                 uint                   `json:"shelf_number"`
		BinNumber                   uint                   `json:"bin_number"`
		TotalBinCount               uint                   `json:"total_bin_count"`
		RecordedQuantity            uint                   `json:"recorded_quantity"`
		CurrentCount                uint                   `json:"current_count"`
		DiscrepancyCount            uint                   `json:"discrepancy_count"`
		DiscrepancyLevel            float32                `json:"discrepancy_level"`
		CycleCountDiscrepancyReason string                 `json:"cycle_count_discrepancy_reason"`
	}

	CycleCountPaginatedResponse struct {
		Total_pages   uint `json:"total_pages"`
		Per_page      uint `json:"per_page"`
		Current_page  uint `json:"current_page"`
		Next_page     uint `json:"next_page"`
		Previous_page uint `json:"previous_page"`
		Total_rows    uint `json:"total_rows"`
	}
	CycleCountGetAllResponse struct {
		Body struct {
			Meta       response.MetaResponse
			Data       []CycleCountGetAll
			Pagination CycleCountPaginatedResponse
		}
	} // @ name CycleCountGetAllResponse
)

type GetProductsListDTO struct {
	Data pq.Int64Array `json:"data"`
}

// Search Cycle Count response
type (
	CycleCountSearch struct {
		CycleCountNumber          string                                `json:"cycle_count_number"`
		WarehouseID               uint                                  `json:"warehouse_id"`
		Warehouse                 shared_pricing_and_location.Locations `json:"warehouse"`
		PartnerID                 uint                                  `json:"partner_id"`
		Partner                   mdm.Partner                           `json:"partner"`
		StatusID                  uint                                  `json:"status_id"`
		Status                    model_core.Lookupcode                 `json:"status"`
		StatusHistory             map[string]interface{}                `json:"status_history"`
		ItemsCount                uint                                  `json:"items_count"`
		CountMethodID             uint                                  `json:"count_method_id"`
		CountMethod               model_core.Lookupcode                 `json:"count_method"`
		OrderLines                []CycleCountLines                     `json:"order_lines"`
		ScheduleDate              time.Time                             `json:"schedule_date"`
		ThresholdLimit            float64                               `json:"threshold_limit"`
		AutoScheduleCycleCount    bool                                  `json:"auto_schedule_cycle_count"`
		FrequencyID               uint                                  `json:"frequency_id"`
		Frequency                 model_core.Lookupcode                 `json:"frequency"`
		CreateInventoryAdjustment bool                                  `json:"create_inventory_adjustment"`
	}

	CycleCountSearchResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []CycleCountSearch
		}
	} // @ name CycleCountSearchResponse
)

// Send Mail CycleCount response
type (
	SendMailCycleCount struct {
		ID            string `query:"id" json:"id"`
		ReceiverEmail string `query:"receiver_email" json:"receiver_email"`
	}

	SendMailCycleCountResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data SendMailCycleCount
		}
	}
)

// Download Pdf CycleCount response
type (
	DownloadPdfCycleCount struct {
		ID string `query:"id"`
	}

	DownloadPdfCycleCountResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DownloadPdfCycleCount
		}
	}
)
