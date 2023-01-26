package cores

import (
	access_template "fermion/backend_core/controllers/access/module"
	model_core "fermion/backend_core/internal/model/core"

	"github.com/lib/pq"
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
type LookupCodesDTO struct {
	Id         uint   `json:"id"`
	LookupCode string `json:"lookup_code"`
	Name       string `json:"display_name"`
}
type LookupCodeDTO struct {
	Id         uint   `json:"id"`
	LookupCode string `json:"lookup_code"`
	Name       string `json:"display_name"`
	LookupType string `json:"lookup_type"`
}
type LookupTypesDTO struct {
	Id                  uint   `json:"id"`
	LookupType          string `json:"lookup_type"`
	DisplayName         string `json:"display_name"`
	NumberOfLookUpCodes uint   `json:"number_of_lookup_codes"`
}

type CreateLookupTypeDTO struct {
	LookupType  string `json:"lookup_type"`
	DisplayName string `json:"display_name"`
}

type LookupCode struct {
	LookupTypeId uint   `json:"lookup_type_id"`
	LookupCode   string `json:"lookup_code"`
	DisplayName  string `json:"display_name"`
}

type CreateUpdateLookupCodesDTO struct {
	LookupCodes []LookupCode `json:"lookup_codes"`
}

type CreateLookupTypeResponseDTO struct {
	LookupTypeId uint `json:"lookup_type_id"`
}

type CreateUpdateLookupCodesResponseDTO struct {
	LookupCodes []struct {
		ID           uint   `json:"id"`
		LookupTypeId uint   `json:"lookup_type_id"`
		LookupCode   string `json:"lookup_code"`
	} `json:"lookup_codes"`
}

type UpdateLookupTypeDTO struct {
	DisplayName string `json:"display_name"`
}

type UpdateLookupCodeDTO struct {
	DisplayName string `json:"display_name"`
}

type SearchAppsDTO struct {
	Name string `query:"name"`
	Code string `query:"code"`
}

type AppsListDTO struct {
	Icon          map[string]interface{} `json:"icon"`
	Name          string                 `json:"name"`
	TotalInstalls int64                  `json:"total_installs"`
	//CategoryId    uint                   `json:"category_id"` // need to implement
	CategoryName string `json:"category_name"`
	Ratings      string `json:"ratings"` // need to implement
	SourceLink   string `json:"source_link"`
	State        string `json:"state"` // installed or not installed
	model_core.Model
}

type InstalledAppsDTO struct {
	Name              string         `json:"name"`
	Code              string         `json:"code"`
	Icon              datatypes.JSON `json:"icon"`
	Type              string         `json:"type"`
	CurrentVersion    string         `json:"current_version"`
	VersionCompatible datatypes.JSON `json:"version_compatible"`
	AccessToken       string         `json:"access_token"`
	AppDetails        datatypes.JSON `json:"app_details"`
	AppAccessTemplate datatypes.JSON `json:"app_access_template"`
	IsCore            bool           `json:"is_core"`
	DeveloperIds      datatypes.JSON `json:"developer_ids"`
	Tutorials         datatypes.JSON `json:"tutorials"`
	Gamification      datatypes.JSON `json:"gamification"`
	Schedulers        datatypes.JSON `json:"schedulers"`
	ConcurrencyList   datatypes.JSON `json:"concurrency_list"`
	CategoryId        *uint          `json:"category_id"`
	CategoryName      string         `json:"category_name"`
	model_core.Model
}

type CurrencyDTO struct {
	Id                  uint           `json:"id"`
	Name                string         `json:"name"`
	CurrencySymbol      string         `json:"currency_symbol"`
	CurrencyCode        string         `json:"currency_code"`
	IsBaseCurrency      bool           `json:"is_base_currency"`
	ExchangeRate        float64        `json:"exchange_rate"`
	ExchangeRateHistory datatypes.JSON `json:"exchange_rate_history"`
	AutoUpdateEr        datatypes.JSON `json:"auto_update_er"`
}
type GetAppStoreDTO struct {
	Name           string                   `json:"name"`
	Code           string                   `json:"code"`
	Icon           map[string]interface{}   `json:"icon"`
	Description    string                   `json:"description"`
	VersionHistory string                   `json:"version_history"`
	ImageOptions   []map[string]interface{} `json:"image_options"`
	Version        string                   `json:"version"`
	ReleaseDate    datatypes.Date           `json:"release_date"`
	TotalInstalls  int64                    `json:"total_installs"`
	Publisher      map[string]interface{}   `json:"publisher"`
	Installation   map[string]interface{}   `json:"installation"`
	FAQs           []map[string]interface{} `json:"faqs"`
	Reviews        []map[string]interface{} `json:"reviews"`
	Support        map[string]interface{}   `json:"support"`
	CurrencyId     *uint                    `json:"currency_id"`
	Currency       model_core.Currency      `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	Price          float64                  `json:"price" gorm:"type:double precision"`
	License        string                   `json:"license" gorm:"type:varchar(100)"`
	Website        string                   `json:"website" gorm:"type:varchar(100)"`
	SourceLink     string                   `json:"source_link"`
	State          string                   `json:"state"` // installed or not installed
	model_core.Model
}

type IRModelFieldsDTO struct {
	Id                     uint   `json:"id"`
	TableId                uint   `json:"table_id"`
	TableName              string `json:"table_name"`
	ColumnName             string `json:"column_name"`
	DataType               string `json:"data_type"`
	CharacterMaximumLength string `json:"character_maximum_length"`
	IsNullable             string `json:"is_nullable"`
	ColumnDefault          string `json:"column_default"`
}
type (
	ListMetaDataDTO struct {
		Data []model_core.IRModelFields
	}
	ViewMetaDataDTO struct {
		Data model_core.IRModelFields
	}
)
type CountryDTO struct {
	Id                         uint           `json:"id"`
	Name                       string         `json:"name"`
	DefaultAccountingPrinciple datatypes.JSON `json:"default_accounting_principle"`
	CountryCode                string         `json:"country_code"`
	CountryCode2               string         `json:"country_code2"`
	StateIds                   pq.Int64Array  `json:"state_ids"`
	TimezoneIds                pq.Int64Array  `json:"time_zone_ids"`
	IsDst                      bool           `json:"is_dst"`
	CurrencyId                 uint           `json:"currency_id"`
	Currency                   CurrencyDTO
}
type StateDTO struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	StateCode string `json:"state_code"`
	CountryId uint   `json:"country_id"`
	Country   CountryDTO
}

type Company struct {
	ID                       int64          `json:"id" `
	Name                     string         `json:"name"`
	Addresses                datatypes.JSON `json:"addresses"`
	Phone                    string         `json:"phone"`
	Email                    string         `json:"email"`
	CompanyDetails           datatypes.JSON `json:"company_details"`
	IsEnterpise              bool           `json:"is_enterpise"`
	ParentId                 uint           `json:"parent_id"`
	ChildIds                 pq.Int64Array  `json:"child_ids"`
	Type                     int            `json:"type"`
	CompanyDefaults          datatypes.JSON `json:"company_defaults"`
	NotificationSettingsId   *uint          `json:"notification_settings_id"`
	NotificationTemplatesIds pq.Int32Array  `json:"notification_template_id"`
	IsEnabled                bool           `json:"is_enabled"`
	IsActive                 bool           `json:"is_active"`
	CreatedByID              *uint          `json:"created_by"`
	UpdatedByID              *uint          `json:"updated_by"`
	OrganizationDetails      datatypes.JSON `json:"organization_details"`
	KYCDocuments             []KYCDocuments `json:"KYC_documents"`
	CloudUpload              bool           `json:"cloud_upload"`
	FilePreferenceID         *uint          `json:"file_preference_id"`
	FilePreference           LookupCodesDTO `json:"file_preference"`
	InvoiceGenerationID      *uint          `json:"invoice_generation_id"`
	InvoiceGeneration        LookupCodesDTO `json:"invoice_generation"`
	BusinessTypeID           *uint          `json:"business_type_id"`
	BusinessType             LookupCodesDTO `json:"business_type"`
}
type KYCDocuments struct {
	CountryId uint                   `json:"country_id"`
	Documents map[string]interface{} `json:"documents"`
}

type UserResponseDTO struct {
	ID               uint           `json:"id"`
	Name             string         `json:"name"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Username         string         `json:"username"`
	Email            string         `json:"email"`
	WorkEmail        string         `json:"work_email"`
	MobileNumber     string         `json:"mobile_number"`
	LoginType        int32          `json:"login_type"`
	Auth             datatypes.JSON `json:"auth"`
	FaConf           datatypes.JSON `json:"2fa_conf"`
	DeviceIds        datatypes.JSON `json:"device_ids"`
	Preferences      datatypes.JSON `json:"preferences"`
	AccessIds        datatypes.JSON `json:"access_ids"`
	CompanyId        *uint          `json:"company_id"`
	TeamHead         string         `json:"team_head"`
	ExternalDetails  datatypes.JSON `json:"external_details"`
	Company          Company
	Profile          datatypes.JSON                           `json:"profile"`
	AddressDetails   datatypes.JSON                           `json:"address_details"`
	AccessTemplateId []access_template.AccessTemplateResponse `json:"access_template_id"`

	//PltPointIds   *PlatformPoints `json:"plt_point_ids" gorm:"foreignkey:PltPointId"`
	// CreatedBy []CoreUsers `gorm:"many2many:created_by"`
	IsEnabled bool `json:"is_enabled" gorm:"default:true"`
	IsActive  bool `json:"is_active" gorm:"default:true"`
}
