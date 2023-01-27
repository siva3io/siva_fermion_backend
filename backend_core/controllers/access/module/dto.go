package module

import (
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

type CoreAppModule struct {
	ID             uint   `json:"id"`
	PortNumber     *uint  `json:"port_number"`
	RoutePath      string `json:"route_path"`
	ImageOption    string `json:"image_option"`
	ModuleCode     string `json:"module_code"`
	ModuleName     string `json:"module_name"`
	ParentModuleId *uint  `json:"parent_module_id"`
	ParentModule   *CoreAppModule
	ItemSequence   *uint `json:"item_sequence"`
	ExternalId     *uint `json:"external_id"`
}
type AccessModuleAction struct {
	AccessTemplateId uint `json:"access_template_id"`
	AccessTemplate   *AccessTemplateResponse
	CoreAppModuleId  uint `json:"core_app_module_id"`
	CoreAppModule    *CoreAppModule
	ViewActions      datatypes.JSON `json:"view_actions_json"`
	DataActions      datatypes.JSON `json:"data_actions_json"`
	ModuleCtrlFlag   *bool          `json:"module_ctrl_flag"`
	DisplayName      string         `json:"display_name"`
}

type AccessModuleActionSideMenuResponse struct {
	AccessTemplateId uint   `json:"access_template_id"`
	DisplayName      string `json:"display_name"`
}

type AccessTemplateResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
