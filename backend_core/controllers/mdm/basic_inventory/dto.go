package basic_inventory

import (
	"fermion/backend_core/pkg/util/response"

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
type CentralizedSearchObjDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type DecentralizedSearchObjDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type CentralizedListDTO struct {
	Meta response.MetaResponse
	Data []CentralizedBasicInventoryDTO
}
type DeentralizedListDTO struct {
	Meta response.MetaResponse
	Data []DecentralizedBasicInventoryDTO
}
type CentralizedViewDTO struct {
	Meta response.SuccessResponse
	Data CentralizedBasicInventoryDTO
}
type DeentralizedViewDTO struct {
	Meta response.SuccessResponse
	Data DecentralizedBasicInventoryDTO
}

// -------------------------DTO for Struct Model --------------------------------
type CentralizedBasicInventoryDTO struct {
	ChannelCode        string `json:"channel_code"`
	ProductVariantId   uint   `json:"product_variant_id"`
	PhysicalLocationId uint   `json:"physical_location_id"`
	Quantity           int64  `json:"quantity"`
	BinTag             string `json:"bin_tag"`
	InventoryCount     int64  `json:"inventory_count"`
	AvailableStock     int64  `json:"available_stock"`
	NoOfUnits          int64  `json:"no_of_units"`
	IsEnabled          *bool  `json:"is_enabled" gorm:"default:true"`
	IsActive           *bool  `json:"is_active" gorm:"default:true"`
}
type DecentralizedBasicInventoryDTO struct {
	CompanyId          uint   `json:"company_id"`
	ChannelCode        string `json:"channel_code"`
	ProductVariantId   uint   `json:"product_variant_id"`
	PhysicalLocationId uint   `json:"physical_location_id"`
	Quantity           int64  `json:"quantity"`
	OnHandQuantity     int64  `json:"on_hand_quantity"`
	PlannedInQuantity  int64  `json:"planned_in_quantity"`
	PlannedOutQuantity int64  `json:"planned_out_quantity"`
	BinTag             string `json:"bin_tag"`
	AvailableStock     int64  `json:"available_stock"`
	StockExpected      int64  `json:"stock_expected"`
	CommitedStock      int64  `json:"commited_stock"`
	NoOfUnits          int64  `json:"no_of_units"`
	CreatedByID        *uint  `json:"created_by"`
	UpdatedByID        *uint  `json:"updated_by"`
	IsEnabled          *bool  `json:"is_enabled" gorm:"default:true"`
	IsActive           *bool  `json:"is_active" gorm:"default:true"`
}
type CentralizedBasicInventoryResponseDTO struct {
	Id                 uint           `json:"id"`
	ChannelCode        string         `json:"channel_code"`
	ProductVariantId   uint           `json:"product_variant_id"`
	ProductDetails     datatypes.JSON `json:"product_details"`
	PhysicalLocationId uint           `json:"physical_location_id"`
	PhysicalLocation   datatypes.JSON `json:"physical_location"`
	Quantity           int64          `json:"quantity"`
	BinTag             string         `json:"bin_tag"`
	InventoryCount     int64          `json:"inventory_count"`
	AvailableStock     int64          `json:"available_stock"`
	StockExpected      int64          `json:"stock_expected"`
	CommitedStock      int64          `json:"commited_stock"`
	NoOfUnits          int64          `json:"no_of_units"`
	CreatedByID        *uint          `json:"created_by"`
	UpdatedByID        *uint          `json:"updated_by"`
	DeletedByID        *uint          `json:"deleted_by"`
	IsEnabled          *bool          `json:"is_enabled" gorm:"default:true"`
	IsActive           *bool          `json:"is_active" gorm:"default:true"`
}
type DecentralizedBasicInventoryResponseDTO struct {
	Id                 uint           `json:"id"`
	ChannelCode        string         `json:"channel_code"`
	ProductVariantId   uint           `json:"product_variant_id"`
	ProductDetails     datatypes.JSON `json:"product_details"`
	PhysicalLocationId uint           `json:"physical_location_id"`
	PhysicalLocation   datatypes.JSON `json:"physical_location"`
	Quantity           int64          `json:"quantity"`
	OnHandQuantity     int64          `json:"on_hand_quantity"`
	PlannedInQuantity  int64          `json:"planned_in_quantity"`
	PlannedOutQuantity int64          `json:"planned_out_quantity"`
	BinTag             string         `json:"bin_tag"`
	AvailableStock     int64          `json:"available_stock"`
	StockExpected      int64          `json:"stock_expected"`
	CommitedStock      int64          `json:"commited_stock"`
	NoOfUnits          int64          `json:"no_of_units"`
	CreatedByID        *uint          `json:"created_by"`
	UpdatedByID        *uint          `json:"updated_by"`
	DeletedByID        *uint          `json:"deleted_by"`
	IsEnabled          *bool          `json:"is_enabled" gorm:"default:true"`
	IsActive           *bool          `json:"is_active" gorm:"default:true"`
}
