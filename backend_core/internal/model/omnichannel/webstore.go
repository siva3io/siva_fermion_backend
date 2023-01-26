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
type (
	Webstore struct {
		model_core.Model
		Code            string         `json:"webstore_code" gorm:"unique"`
		Name            string         `json:"name" gorm:"varchar(50)"`
		Image_options   datatypes.JSON `json:"image_options" gorm:"json"`
		Location_served datatypes.JSON `json:"location_served" gorm:"json"`
		WebstoreDetails datatypes.JSON `json:"webstore_details" gorm:"json"`
	}
)
type User_Webstore_Link struct {
	model_core.Model
	WebstoreCode                   string                         `json:"webstore_code" gorm:""`
	Webstore                       Webstore                       `json:"webstore" gorm:"foreignKey:WebstoreCode; references:Code"`
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
