package contacts

import (
	app_core "fermion/backend_core/controllers/cores"
	"fermion/backend_core/internal/model/core"
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
type ContactListDTO struct {
	Meta response.MetaResponse
	Data []PartnerResponseDTO
}
type ContactViewDTO struct {
	Meta response.SuccessResponse
	Data PartnerResponseDTO
}
type SwaggerCreateOrUpdatePartnerResponseDTO struct {
	RequestId string `json:"request_id"`
}

type Ids struct {
	ID uint `json:"id"`
}
type IdNameDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type ImageOptionsDTO struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Data string `json:"data"`
	Type string `json:"type"`
	Size int    `json:"size"`
}
type BillingDetailsDTO struct {
	BankId        string `json:"bank_id"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	Upi           string `json:"upi"`
	AccountName   string `json:"account_name"`
	IfscCode      string `json:"ifsc_code"`
}
type AdditionalInformation struct {
	AdditionalInformation string `json:"additional_information"`
	Notes                 string `json:"notes"`
	GSTNo                 string `json:"gst_no"`
	EmergencyContact      string `json:"emergency_contact"`
	AdditionalContact     string `json:"additional_contact"`
	WebsiteUrl            string `json:"website_url"`
	GSTDoc                string `json:"gst_doc"`
	DateOfBirth           string `json:"date_of_birth"`
}
type AddressDetailsDTO struct {
	LocationName        string                  `json:"location_name"`
	LocationType        app_core.LookupCodesDTO `json:"location"`
	AddressLine1        string                  `json:"address_line_1"`
	AddressLine2        string                  `json:"address_line_2"`
	AddressLine3        string                  `json:"address_line_3"`
	Landmark            string                  `json:"land_mark"`
	Pincode             string                  `json:"pin_code"`
	Country             app_core.CountryDTO     `json:"country"`
	State               app_core.StateDTO       `json:"state"`
	City                string                  `json:"city"`
	GSTINNumber         string                  `json:"gst_in_number"`
	ContactPersonName   string                  `json:"contact_person_name"`
	ContactPersonNumber string                  `json:"contact_person_number"`
	MarkDefaultAddress  bool                    `json:"mark_default_address"`
	AddressTypeID       uint                    `json:"address_type_id"` //1,2,3
	AddressType         string                  `json:"address_type"`    //shipping, billing, both
}

// ----------------------DTO for struct Model---------------------------
type PartnerRequestDTO struct {
	CreatedByID           *uint                  `json:"created_by_id"`
	UpdatedByID           *uint                  `json:"updated_by_id"`
	DeletedByID           *uint                  `json:"deleted_by_id"`
	ContactTypeId         *uint                  `json:"contact_type_id"`
	FirstName             string                 `json:"first_name"`
	LastName              string                 `json:"last_name"`
	IsAllowedLogin        *bool                  `json:"is_allowed_login"`
	AccessTemplateId      []*core.AccessTemplate `json:"access_template_id"`
	AccessIds             datatypes.JSON         `json:"access_ids"`
	CompanyName           string                 `json:"company_name"`
	Properties            []*Ids                 `json:"properties,omitempty"`
	PrimaryEmail          string                 `json:"primary_email"`
	PrimaryPhone          string                 `json:"primary_phone"`
	ParentId              *uint                  `json:"parent_id"`
	ImageOptions          datatypes.JSON         `json:"image_options,omitempty"`
	AddressDetails        []AddressDetailsDTO    `json:"address_details,omitempty"`
	BillingDetails        datatypes.JSON         `json:"billing_details,omitempty"`
	ProfileInfo           datatypes.JSON         `json:"profile_info,omitempty"`
	ExternalDetails       datatypes.JSON         `json:"external_details"`
	AdditionalInformation datatypes.JSON         `json:"additional_information,omitempty"`
	IsEnabled             *bool                  `json:"is_enabled"`
	IsActive              *bool                  `json:"is_active"`
	RoleId                *uint                  `json:"role_id"`
	UserTypes             datatypes.JSON         `json:"user_types"`
}
type PartnerResponseDTO struct {
	Id             uint                    `json:"id"`
	ContactType    app_core.LookupCodesDTO `json:"contact_type"`
	FirstName      string                  `json:"first_name"`
	LastName       string                  `json:"last_name"`
	IsAllowedLogin *bool                   `json:"is_allowed_login"`
	// AccessTemplateId      uint                      `json:"access_template_id"`
	CompanyName           string                    `json:"company_name"`
	Properties            []app_core.LookupCodesDTO `json:"properties"`
	PrimaryEmail          string                    `json:"primary_email"`
	PrimaryPhone          string                    `json:"primary_phone"`
	ParentId              *uint                     `json:"parent_id"`
	Parent                *PartnerResponseDTO       `json:"parent"`
	ChildIds              []PartnerResponseDTO      `json:"child_contacts"`
	ImageOptions          datatypes.JSON            `json:"image_options"`
	AddressDetails        []AddressDetailsDTO       `json:"address_details"`
	DefaultAddress        AddressDetailsDTO         `json:"default_address"`
	BillingDetails        BillingDetailsDTO         `json:"billing_details"`
	ProfileInfo           IdNameDTO                 `json:"profile_info"`
	AdditionalInformation datatypes.JSON            `json:"additional_information"`
	IsEnabled             *bool                     `json:"is_enabled"`
	IsActive              *bool                     `json:"is_active" `
}
