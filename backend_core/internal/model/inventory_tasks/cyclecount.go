package inventory_tasks

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"

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
type CycleCount struct {
	CycleCountNumber          string                                `json:"cycle_count_number" gorm:"type:varchar(50);unique"`
	CycleCountDate            time.Time                             `json:"cycle_count_date" gorm:"type:time"`
	WarehouseID               uint                                  `json:"warehouse_id" gorm:"type:integer"`
	Warehouse                 shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseID; references:ID"`
	PartnerID                 uint                                  `json:"partner_id" gorm:"type:integer"`
	Partner                   mdm.Partner                           `json:"partner" gorm:"foreignKey:PartnerID; references:ID"`
	StatusID                  uint                                  `json:"status_id" gorm:"type:integer"`
	Status                    model_core.Lookupcode                 `json:"status" gorm:"foreignKey:StatusID; references:ID"`
	StatusHistory             datatypes.JSON                        `json:"status_history" gorm:"type:json"`
	ItemsCount                uint                                  `json:"items_count" gorm:"type:integer"`
	CountMethodID             uint                                  `json:"count_method_id" gorm:"type:integer"`
	CountMethod               model_core.Lookupcode                 `json:"count_method" gorm:"foreignKey:CountMethodID; references:ID"`
	OrderLines                []CycleCountLines                     `json:"order_lines" gorm:"foreignKey:Cycle_count_id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScheduleDate              time.Time                             `json:"schedule_date" gorm:"type:time"`
	ThresholdLimit            float64                               `json:"threshold_limit" gorm:"type:double precision"`
	AutoScheduleCycleCount    bool                                  `json:"auto_schedule_cycle_count" gorm:"type: boolean"`
	FrequencyID               uint                                  `json:"frequency_id" gorm:"type:integer"`
	Frequency                 model_core.Lookupcode                 `json:"frequency" gorm:"foreignKey:FrequencyID; references:ID"`
	CreateInventoryAdjustment bool                                  `json:"create_inventory_adjustment" gorm:"type:boolean"`
	model_core.Model
}

type CycleCountLines struct {
	Cycle_count_id              uint                  `json:"-"`
	ProductID                   *uint                 `json:"product_id" gorm:"type:integer"`
	Product                     *mdm.ProductTemplate  `json:"product" gorm:"foreignKey:ProductID"`
	ProductVariantID            *uint                 `json:"product_variant_id"`
	ProductVariant              *mdm.ProductVariant   `json:"product_variant" gorm:"foreignKey:ProductVariantID; references:ID"`
	BinLocationID               datatypes.JSON        `json:"bin_location_id" gorm:"type:json"`
	InventoryCount              uint                  `json:"inventory_count" gorm:"type:integer"`
	LocationDetails             datatypes.JSON        `json:"location_details" gorm:"type:json"`
	LocationSpaceTypeID         uint                  `json:"location_space_type_id" gorm:"type:integer"`
	LocationSpaceType           model_core.Lookupcode `json:"location_space_type" gorm:"foreignKey:LocationSpaceTypeID"`
	LocationInputTypeID         uint                  `json:"location_input_type_id" gorm:"type:integer"`
	LocationInputType           model_core.Lookupcode `json:"location_input_type" gorm:"foreignKey:LocationInputTypeID"`
	RowNumber                   uint                  `json:"row_number" gorm:"type:integer"`
	RackNumber                  uint                  `json:"rack_number" gorm:"type:integer"`
	ShelfNumber                 uint                  `json:"shelf_number" gorm:"type:integer"`
	BinNumber                   uint                  `json:"bin_number" gorm:"type:integer"`
	TotalBinCount               uint                  `json:"total_bin_count" gorm:"type:integer"`
	RecordedQuantity            uint                  `json:"recorded_quantity" gorm:"type:integer"`
	CurrentCount                uint                  `json:"current_count" gorm:"type:integer"`
	DiscrepancyCount            uint                  `json:"discrepancy_count" gorm:"type:integer"`
	DiscrepancyLevel            float32               `json:"discrepancy_level" gorm:"type:double precision"`
	CycleCountDiscrepancyReason string                `json:"cycle_count_discrepancy_reason" gorm:"type:varchar(250)"`
	model_core.Model
}
