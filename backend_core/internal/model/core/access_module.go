package core

import "gorm.io/datatypes"

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
type CoreAppModule struct {
	Model
	PortNumber     *uint  `json:"port_number"`
	RoutePath      string `json:"route_path"`
	ImageOption    string `json:"image_option"`
	ModuleCode     string `json:"module_code"`
	ModuleName     string `json:"module_name"`
	ParentModuleId *uint  `json:"parent_module_id" gorm:"column:parent_module"`
	ParentModule   *CoreAppModule
	ItemSequence   *uint `json:"item_sequence"`
	ExternalId     *uint `json:"external_id"`
}
type AccessModuleAction struct {
	Model
	CoreAppModuleId  uint            `json:"core_app_module_id"`
	CoreAppModule    *CoreAppModule  `gorm:"foreignkey:CoreAppModuleId; reference:ID"`
	ViewActions      datatypes.JSON  `json:"view_actions_json" gorm:"type:json; default:'[]'"`
	DataActions      datatypes.JSON  `json:"data_actions_json" gorm:"type:json; default:'[]'"`
	ModuleCtrlFlag   *bool           `json:"module_ctrl_flag" gorm:"default:true"`
	DisplayName      string          `gorm:"type:varchar(50); not null;" json:"display_name"`
	AccessTemplateId uint            `json:"access_template_id"`
	AccessTemplate   *AccessTemplate `gorm:"foreignkey:AccessTemplateId; reference:ID"`
}
