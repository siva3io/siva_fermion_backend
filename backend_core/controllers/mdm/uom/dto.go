package uom

import (
	app_core "fermion/backend_core/controllers/cores"
	"fermion/backend_core/pkg/util/response"
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
type UomListDTO struct {
	Meta response.MetaResponse
	Data []UomResponseDTO
}
type UomViewDTO struct {
	Meta response.SuccessResponse
	Data UomResponseDTO
}
type UomClassListDTO struct {
	Meta response.MetaResponse
	Data []UomClassResponseDTO
}
type UomClassViewtDTO struct {
	Meta response.SuccessResponse
	Data UomClassResponseDTO
}

//-------------------DTO for struct Model----------------------------

type IdNameDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type UomRequestDTO struct {
	ItemTypeId       *uint   `json:"item_type_id"`
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	UomClassId       uint    `json:"uom_class_id"`
	UomClassName     string  `json:"uom_class_name"`
	BaseUom          string  `json:"base_uom"`
	ConversionTypeId *uint   `json:"conversion_type_id"`
	ConversionFactor float64 `json:"conversion_factor"`
	IsEnabled        *bool   `json:"is_enabled" gorm:"default:true"`
	IsActive         *bool   `json:"is_active" gorm:"default:true"`
}
type UomClassRequestDTO struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	BaseUom     string `json:"base_uom"`
	Description string `json:"description"`
	IsEnabled   *bool  `json:"is_enabled" gorm:"default:true"`
	IsActive    *bool  `json:"is_active" gorm:"default:true"`
}
type UomResponseDTO struct {
	Id               uint                    `json:"id"`
	ItemType         app_core.LookupCodesDTO `json:"item_type"`
	Code             string                  `json:"code"`
	Name             string                  `json:"name"`
	Description      string                  `json:"description"`
	UomClassCode     UomClassResponseDTO     `json:"uom_class_code"`
	UomClassName     string                  `json:"uom_class_name"`
	BaseUom          string                  `json:"base_uom"`
	ConversionType   app_core.LookupCodesDTO `json:"conversion_type"`
	ConversionFactor float64                 `json:"conversion_factor"`
	IsEnabled        *bool                   `json:"is_enabled" gorm:"default:true"`
	IsActive         *bool                   `json:"is_active" gorm:"default:true"`
}

type UomClassResponseDTO struct {
	Id          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	BaseUom     string `json:"base_uom"`
	Description string `json:"description"`
	IsEnabled   *bool  `json:"is_enabled" gorm:"default:true"`
	IsActive    *bool  `json:"is_active" gorm:"default:true"`
}
