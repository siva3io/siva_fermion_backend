package mdm

import (
	model_core "fermion/backend_core/internal/model/core"

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
type CentralizedBasicInventory struct {
	model_core.Model
	ChannelCode        string         `json:"channel_code" gorm:"column:channel_code"`
	ProductVariantId   uint           `json:"product_variant_id"`
	ProductDetails     datatypes.JSON `json:"product_details"`
	PhysicalLocationId uint           `json:"physical_location_id"`
	PhysicalLocation   datatypes.JSON `json:"physical_location" gorm:"type:json; default:'{}'"`
	Quantity           int64          `json:"quantity"`
	BinTag             string         `json:"bin_tag"`
	InventoryCount     int64          `json:"inventory_count"`
	AvailableStock     int64          `json:"available_stock"`
	StockExpected      int64          `json:"stock_expected"`
	CommitedStock      int64          `json:"commited_stock"`
	NoOfUnits          int64          `json:"no_of_units"`
}

type DecentralizedBasicInventory struct {
	model_core.Model
	ChannelCode        string         `json:"channel_code" gorm:"column:channel_code"`
	ProductVariantId   uint           `json:"product_variant_id"`
	ProductDetails     datatypes.JSON `json:"product_details"`
	PhysicalLocationId uint           `json:"physical_location_id"`
	PhysicalLocation   datatypes.JSON `json:"physical_location" gorm:"type:json; default:'{}'"`
	Quantity           int64          `json:"quantity" gorm:"default:0"`
	OnHandQuantity     int64          `json:"on_hand_quantity" gorm:"default:0"`
	PlannedInQuantity  int64          `json:"planned_in_quantity" gorm:"default:0"`
	PlannedOutQuantity int64          `json:"planned_out_quantity" gorm:"default:0"`
	BinTag             string         `json:"bin_tag"`
	AvailableStock     int64          `json:"available_stock" gorm:"default:0"`
	StockExpected      int64          `json:"stock_expected" gorm:"default:0"`
	CommitedStock      int64          `json:"commited_stock" gorm:"default:0"`
	NoOfUnits          int64          `json:"no_of_units" gorm:"default:0"`
}

type CentralizedInventoryTransactions struct {
	model_core.Model
	OrderID            uint   `json:"order_id"`
	OrderType          string `json:"order_type"`
	Description        string `json:"description"`
	ProductVariantId   uint   `json:"product_variant_id"`
	PhysicalLocationId uint   `json:"physical_location_id"`
	InventoryId        uint   `json:"inventory_id"`
	OpeningInventory   int64  `json:"opening_inventory"`
	DepositedStock     int64  `json:"deposited_stock"`
	WithdrawnStock     int64  `json:"withdrawn_stock"`
	ClosingInventory   int64  `json:"closing_inventory"`
}
