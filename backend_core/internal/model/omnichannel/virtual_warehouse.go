package omnichannel

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
type VirtualWarehouse struct {
	model_core.Model
	Code                    string         `json:"virtual_warehouse_code" gorm:"unique"`
	Name                    string         `json:"name" gorm:"varchar(50)"`
	Image_options           datatypes.JSON `json:"image_options" gorm:"json"`
	Location_served         datatypes.JSON `json:"location_served" gorm:"json"`
	VirtualWarehouseDetails datatypes.JSON `json:"virtual_warehouse_details" gorm:"json"`
}

type User_Virtual_Warehouse_Link struct {
	model_core.Model
	VirtualWarehouseCode           string                         `json:"virtual_warehouse_code" gorm:""`
	VirtualWarehouse               VirtualWarehouse               `json:"virtual_warehouse" gorm:"foreignKey:VirtualWarehouseCode; references:Code"`
	Auth                           datatypes.JSON                 `gorm:"type:json;default:'[]';not null" json:"auth"`
	Channel_details                Channel_details                `gorm:"embedded"`
	Fullfilment_details            Fullfilment_details            `gorm:"embedded"`
	Inventory_details              Inventory_details              `gorm:"embedded"`
	Order_details                  Order_details                  `gorm:"embedded"`
	Advanced_channel_configuration Advanced_channel_configuration `gorm:"embedded"`
	Payment_mapping                Payment_mapping                `gorm:"embedded"`
	Pre_sync_configuration         Pre_sync_configuration         `gorm:"embedded"`
	Sync_configuration             Sync_configuration             `gorm:"embedded"`
	Sync_fetch_data                Sync_fetch_data                `gorm:"embedded"`
	Inventory_Automation_Details   Inventory_Automation_Details   `gorm:"embedded"`
}

// register for virtual warehouse
type User_Virtual_Warehouse_Registration struct {
	model_core.Model
	VirtualWarehouseCode string           `json:"virtual_warehouse_code" gorm:"type:text"`
	VirtualWarehouse     VirtualWarehouse `json:"virtual_warehouse" gorm:"foreignKey:VirtualWarehouseCode; references:Code"`
	Full_name            string           `json:"full_name" gorm:"varchar"`
	Email                string           `json:"email" gorm:"varchar"`
	Mobile_no            int64            `json:"mobile_no" `
	Date_of_birth        string           `json:"date_of_birth" gorm:"date"`
	Bank_info            Bank_info        `gorm:"embedded"`
	KYC_info             Kyc_info         `gorm:"embedded"`
	Status               string           `json:"status"`
	Status_history       string           `json:"status_history" gorm:"status"`
}
