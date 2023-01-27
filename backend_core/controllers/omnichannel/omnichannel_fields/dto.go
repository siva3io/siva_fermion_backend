package omnichannel_fields

import (
	"fermion/backend_core/controllers/cores"
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

type OmnichannelFieldRequestDto struct {
	AppId             *uint          `json:"app_id"`
	Name              string         `json:"name"`
	Placeholder       string         `json:"placeholder"`
	Type              string         `json:"type"`
	DataType          string         `json:"data_type"`
	DataSource        string         `json:"data_source"`
	DataSourceParams  datatypes.JSON `json:"data_source_params"`
	AllowedValues     datatypes.JSON `json:"allowed_values"`
	IsMandatory       *bool          `json:"is_mandatory"`
	DisplayName       string         `json:"display_name"`
	ChannelFunctionId *uint          `json:"channel_function_id"`
	ChannelTypeId     *uint          `json:"channel_type_id"`
	SectionSequence   int            `json:"section_sequence"`
	FieldSequence     int            `json:"field_sequence"`
	SectionName       string         `json:"section_name"`
	IsDisabled        *bool          `json:"is_disabled"`
	IsHidden          *bool          `json:"is_hidden"`
	IsEditable        *bool          `json:"is_editable"`
}

type OmnichannelFieldResponseDto struct {
	ID                uint                         `json:"id"`
	AppId             *uint                        `json:"app_id"`
	SectionSequence   int                          `json:"section_sequence"`
	FieldSequence     int                          `json:"field_sequence"`
	SectionName       string                       `json:"section_name"`
	ChannelFunctionId *uint                        `json:"channel_function_id"`
	ChannelFunction   cores.LookupCodeDTO          `json:"channel_function"`
	ChannelTypeId     *uint                        `json:"channel_type_id"`
	ChannelType       cores.AppCategoryResponseDTO `json:"channel_type"`
	Name              string                       `json:"name"`
	Placeholder       string                       `json:"placeholder"`
	Type              string                       `json:"type"`
	DataType          string                       `json:"data_type"`
	DataSource        string                       `json:"data_source"`
	DataSourceParams  datatypes.JSON               `json:"data_source_params"`
	AllowedValues     datatypes.JSON               `json:"allowed_values"`
	IsMandatory       *bool                        `json:"is_mandatory"`
	DisplayName       string                       `json:"display_name"`
	IsDisabled        *bool                        `json:"is_disabled"`
	IsEditable        *bool                        `json:"is_editable"`
}

type OmnichannelFieldDeleteDto struct {
	Deleted_id int
	Field_id   int
}

type ViewAppFieldsQueryDto struct {
	AppId             uint `query:"app_id"`
	ChannelFunctionId uint `query:"channel_function_id"`
	ChannelTypeId     uint `query:"channel_type_id"`
}

type ViewAppFieldsDataQueryDto struct {
	FieldAppId        uint `query:"field_app_id"`
	ChannelFunctionId uint `query:"channel_function_id"`
}

type GetAppDataFilterQueryDto struct {
	Fields string `query:"fields"`
	AppId  uint   `query:"app_id"`
}

type OmnichannelFieldDataRequestDto struct {
	AppCode string                 `json:"app_code"`
	Fields  []OmnichannelFieldData `json:"fields"`
}

type OmnichannelFieldData struct {
	OmnichannelFieldId *uint       `json:"omnichannel_field_id"`
	Data               interface{} `json:"data"`
	DataType           string      `json:"data_type"`
}

type OmnichannelFieldDataResponseDto struct {
	FieldAppId         *uint                       `json:"field_app_id"`
	OmnichannelFieldId *uint                       `json:"omnichannel_field_id"`
	OmnichannelField   OmnichannelFieldResponseDto `json:"omnichannel_field"`
	Data               interface{}                 `json:"data"`
	DataType           string                      `json:"data_type"`
}

type OmnichannelSyncSettingsResponseDto struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	JobId       uint   `json:"job_id"`
	State       bool   `json:"state"`
	Frequency   string `json:"frequency"`
}

type OmnichannelSyncSettingsResponseDtoWithDropDownDto struct {
	SyncSettings      []OmnichannelSyncSettingsResponseDto `json:"sync_settings"`
	FrequencyDropDown []cores.LookupCodesDTO               `json:"frequency_dropdown"`
}

type OmnichannelSyncSettingsRequestDto struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	Frequency string `json:"frequency"`
	State     bool   `json:"state"`
}

type AppDataFilterResponseDto struct {
	ID   uint           `json:"id"`
	Name string         `json:"name"`
	Data datatypes.JSON `json:"data"`
}

type (
	ViewAppFieldsResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data OmnichannelFieldResponseDto
		}
	}
	UpsertOmnichannelFieldsDataResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	}
	ViewAppFieldsDataResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data OmnichannelFieldDataResponseDto
		}
	}
	GetAppSyncSettingsResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data OmnichannelSyncSettingsResponseDto
		}
	}
	UpsertAppSyncSettingsResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	}
	OmnichannelFieldDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data OmnichannelFieldResponseDto
		}
	}
	GetAllOmnichannelFieldsResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []OmnichannelFieldResponseDto
		}
	}
)
