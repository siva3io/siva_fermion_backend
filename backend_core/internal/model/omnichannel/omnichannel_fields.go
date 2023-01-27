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

type OmnichannelField struct {
	model_core.Model
	ChannelTypeId     *uint                  `json:"channel_type_id"`
	ChannelType       model_core.AppCategory `json:"channel_type" gorm:"foreignKey:ChannelTypeId;references:ID"`
	ChannelFunctionId *uint                  `json:"channel_function_id"`
	ChannelFunction   model_core.Lookupcode  `json:"channel_function" gorm:"foreignKey:ChannelFunctionId;references:ID"`
	SectionSequence   int                    `json:"section_sequence"`
	FieldSequence     int                    `json:"field_sequence"`
	SectionName       string                 `json:"section_name"`
	Name              string                 `json:"name"`
	Placeholder       string                 `json:"placeholder"`
	Type              string                 `json:"type"`
	DataType          string                 `json:"data_type"`
	DataSource        string                 `json:"data_source"`
	DataSourceParams  datatypes.JSON         `json:"data_source_params"`
	AllowedValues     datatypes.JSON         `json:"allowed_values"`
	IsMandatory       *bool                  `json:"is_mandatory"`
	DisplayName       string                 `json:"display_name"`
	IsDisabled        *bool                  `json:"is_disabled" gorm:"default:false"`
	IsHidden          *bool                  `json:"is_hidden" gorm:"default:false"`
	IsEditable        *bool                  `json:"is_editable" gorm:"default:true"`
}

type OmnichannelFieldData struct {
	model_core.Model
	FieldAppId         *uint `json:"field_app_id" gorm:"uniqueIndex:idx_field_app_id_omnichannel_field_id"`
	FieldApp           *model_core.InstalledApps
	OmnichannelFieldId *uint            `json:"omnichannel_field_id" gorm:"uniqueIndex:idx_field_app_id_omnichannel_field_id"`
	OmnichannelField   OmnichannelField `json:"omnichannel_field" gorm:"foreignKey:OmnichannelFieldId;references:ID"`
	Data               datatypes.JSON   `json:"data"`
	DataType           string           `json:"data_type"`
}
