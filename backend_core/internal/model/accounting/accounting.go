package accounting

import (
	model_core "fermion/backend_core/internal/model/core"

	"gorm.io/datatypes"
)

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
type UserAccountingLink struct {
	model_core.Model
	AccountingCode string         `json:"accounting_code" gorm:""`
	Accounting     Accounting     `json:"accounting" gorm:"foreignKey:AccountingCode; references:Code"`
	Auth           datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"auth"`
}

type (
	Accounting struct {
		model_core.Model
		Code              string         `json:"code" gorm:"unique"`
		Name              string         `json:"name" gorm:"varchar(50)"`
		ImageOptions      datatypes.JSON `json:"image_options" gorm:"json"`
		LocationServed    datatypes.JSON `json:"location_served" gorm:"json"`
		AccountingDetails datatypes.JSON `json:"accounting_details" gorm:"json"`
	}
)
