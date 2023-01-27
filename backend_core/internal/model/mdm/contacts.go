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
type Partner struct {
	model_core.Model
	ContactTypeId         *uint                 `json:"contact_type_id" gorm:""`
	ContactType           model_core.Lookupcode `json:"contact_type" gorm:"foreignkey:ContactTypeId;references:ID"`
	FirstName             string                `json:"first_name" gorm:"type:text"`
	LastName              string                `json:"last_name" gorm:"type:text"`
	IsAllowedLogin        *bool                 `json:"is_allowed_login" gorm:"default:false"`
	CompanyName           string                `json:"company_name" gorm:""`
	Properties            datatypes.JSON        `json:"properties,omitempty" gorm:"type:json"`
	PrimaryEmail          string                `json:"primary_email" gorm:"not null;unique"`
	PrimaryPhone          string                `json:"primary_phone"`
	ParentId              *uint                 `json:"parent_id"`
	Parent                *Partner              `json:"parent"`
	ChildIds              []Partner             `json:"child_contacts" gorm:"foreignkey:ParentId"`
	ImageOptions          datatypes.JSON        `json:"image_options,omitempty" gorm:"type:json; default:'{}'"`
	AddressDetails        datatypes.JSON        `json:"address_details,omitempty" gorm:"type:json; default:'[]'"`
	BillingDetails        datatypes.JSON        `json:"billing_details,omitempty" gorm:"type:json; default:'{}'"`
	ProfileInfo           datatypes.JSON        `json:"profile_info,omitempty" gorm:"type:json; default:'{}'"`
	AdditionalInformation datatypes.JSON        `json:"additional_information,omitempty" gorm:"type:json; default:'{}'"`
	UserId                uint                  `json:"user_id" gorm:""`
}
