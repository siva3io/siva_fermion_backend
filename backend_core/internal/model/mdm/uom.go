package mdm

import model_core "fermion/backend_core/internal/model/core"

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type Uom struct {
	model_core.Model
	ItemTypeId       uint                  `json:"item_type_id"`
	ItemType         model_core.Lookupcode `json:"item_type" gorm:"foreignkey:ItemTypeId"`
	Code             string                `json:"code" gorm:"unique"`
	Name             string                `json:"name" gorm:""`
	Description      string                `json:"description" gorm:""`
	UomClassId       uint                  `json:"uom_class_id"`
	UomClassCode     UomClass              `json:"uom_class_code" gorm:"foreignkey:UomClassId;references:ID"`
	UomClassName     string                `json:"uom_class_name"`
	BaseUom          string                `json:"base_uom" gorm:""`
	ConversionTypeId uint                  `json:"conversion_type_id"`
	ConversionType   model_core.Lookupcode `json:"conversion_type" gorm:"foreignkey:ConversionTypeId"`
	ConversionFactor float64               `json:"conversion_factor" gorm:""`
}

type UomClass struct {
	model_core.Model
	Code    string `json:"code" gorm:"unique"`
	Name    string `json:"name" gorm:""`
	BaseUom string `json:"base_uom" gorm:""`
	//Uom         *Uom   `json:"uom" gorm:"foreignkey:ID;references:ID"`
	Description string `json:"description" gorm:""`
}
