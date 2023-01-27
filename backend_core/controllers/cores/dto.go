package cores

import (
	"time"

	access_template "fermion/backend_core/controllers/access/module"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/pkg/util/response"

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
	SourceCode string `json:"source_code,omitempty"`
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
	ID                       *uint           `json:"id" `
	Name                     string          `json:"name"`
	Phone                    string          `json:"phone"`
	Email                    string          `json:"email"`
	Website                  string          `json:"website"`
	LanguagePreferred        string          `json:"language_preferred"`
	TimeSlotsPreferred       string          `json:"time_slots_preferred"`
	WorkingHoursStartTime    time.Time       `json:"working_hours_start_time"`
	WorkingHoursEndTime      time.Time       `json:"working_hours_end_time"`
	CompanyDetails           *CompanyDetails `json:"company_details,omitempty"`
	OndcDetailsId            *uint           `json:"ondc_details_id"`
	OndcDetails              OndcDetails     `json:"ondc_details"`
	IsEnterpise              bool            `json:"is_enterpise"`
	ParentId                 uint            `json:"parent_id"`
	ChildIds                 pq.Int64Array   `json:"child_ids"`
	Type                     *uint           `json:"type,omitempty"`
	CompanyType              LookupCodesDTO  `json:"company_type"`
	CompanyDefaults          datatypes.JSON  `json:"company_defaults,omitempty"`
	NotificationSettingsId   *uint           `json:"notification_settings_id"`
	NotificationTemplatesIds pq.Int32Array   `json:"notification_template_id"`
	IsEnabled                bool            `json:"is_enabled"`
	IsActive                 bool            `json:"is_active"`
	CreatedByID              *uint           `json:"created_by"`
	UpdatedByID              *uint           `json:"updated_by"`
	KYCDocuments             KycDocuments    `json:"KYC_documents"`
	CloudUpload              bool            `json:"cloud_upload"`
	BankDetails              BankDetails     `json:"bank_details,omitempty"`
	FilePreferenceID         *uint           `json:"file_preference_id"`
	FilePreference           LookupCodesDTO  `json:"file_preference"`
	InvoiceGenerationID      *uint           `json:"invoice_generation_id"`
	InvoiceGeneration        LookupCodesDTO  `json:"invoice_generation"`
	BusinessTypeID           *uint           `json:"business_type_id"`
	BusinessType             LookupCodesDTO  `json:"business_type"`
	ShippingPreferenceID     *uint           `json:"shipping_preference_id"`
	ShippingPreference       LookupCodesDTO  `json:"shipping_preference"`
	DeliveryTypeID           *uint           `json:"delivery_type_id"`
	DeliveryType             LookupCodesDTO  `json:"delivery_type"`
	DeliveryPreferencesID    *uint           `json:"delivery_preferences_id"`
	DeliveryPreferences      LookupCodesDTO  `json:"delivery_preferences"`
	OTPPreferenceID          *uint           `json:"otp_preference_id"`
	OTPPreferences           LookupCodesDTO  `json:"otp_preferences"`
	NotificationPreferenceID *uint           `json:"notification_preference_id"`
	NotificationPreferences  LookupCodesDTO  `json:"notification_preferences"`
	AadhaarNumber            string          `json:"aadhaar_number,omitempty"`
	GSTIN                    string          `json:"gstin,omitempty"`
	PANNumber                string          `json:"pan_number,omitempty"`
	RatingAverage            float64         `json:"rating_average"`
	RatingCount              uint            `json:"rating_count"`
}

type KycDocuments struct {
	AdharCard datatypes.JSON `json:"aadhar_card,omitempty"`
	Gst       datatypes.JSON `json:"gst,omitempty"`
	Iec       datatypes.JSON `json:"iec,omitempty"`
	PanCard   datatypes.JSON `json:"pan_card,omitempty"`
	PassPort  datatypes.JSON `json:"passport,omitempty"`
	VoterId   datatypes.JSON `json:"voter_id,omitempty"`
}

type BankDetails struct {
	Bank_name                string         `json:"bank_name"`
	Holder_name              string         `json:"holder_name"`
	Account_number           string         `json:"account_number"`
	Ifsc_code                string         `json:"ifsc_code"`
	Bank_statement           datatypes.JSON `json:"bank_statement,omitempty"`
	State                    string         `json:"state"`
	City                     string         `json:"city"`
	PennyTransferVerfication bool           `json:"penny_transfer_verification"`
	UPIAddress               string         `json:"upi_address"`
	BranchName               string         `json:"branch_name"`
	Cancelled_cheque         datatypes.JSON `json:"cancelled_cheque,omitempty"`
}
type CompanyDetails struct {
	BusinessName               string         `json:"business_name"`
	BusinessAddress            string         `json:"business_address"`
	FinancialYearStartId       *uint          `json:"financial_year_start_id"`
	FinancialYearStart         LookupCodesDTO `json:"financial_year_start" gorm:"foreignkey:FinancialYearStartId; references:ID"`
	FinancialYearEndId         *uint          `json:"financial_year_end_id"`
	FinancialYearEnd           LookupCodesDTO `json:"financial_year_end" gorm:"foreignkey:FinancialYearEndId; references:ID"`
	AuthorisedSignatory        string         `json:"authorised_signatory"`
	AuthorisedSignatoryAddress string         `json:"authorised_signatory_address"`
	StdCodeID                  *uint          `json:"std_code_id"`
	StdCode                    LookupCodesDTO `json:"std_code"`
	StoreName                  string         `json:"store_name"`
	StoreDescription           string         `json:"store_description"`
	ServiceableAreas           datatypes.JSON `json:"serviceable_areas"`
	DomainId                   *uint          `json:"domain_id"`
	Domain                     LookupCodesDTO `json:"domain"`
	EstablishedOn              time.Time      `json:"established_on"`
	StoreTimings               datatypes.JSON `json:"store_timings,omitempty"`
	SellerApps                 []string       `json:"seller_apps,omitempty"`
	EnableEmailNotifications   bool           `json:"enable_email_notifications"`
	EnablePhoneNotifications   bool           `json:"enable_phone_notifications"`
}

type OndcDetails struct {
	ID                      *uint     `json:"id"`
	SubscriberId            string    `json:"subscriber_id,omitempty"`
	SubscriberURL           string    `json:"subscriber_url,omitempty"`
	SigningPublicKey        string    `json:"signing_public_key,omitempty"`
	SigningPrivateKey       string    `json:"signing_private_key,omitempty"`
	EncryptionPrivateKey    string    `json:"encryption_private_key,omitempty"`
	EncryptionPublicKey     string    `json:"encryption_public_key,omitempty"`
	UniqueId                string    `json:"unique_id,omitempty"`
	Type                    *uint     `json:"type"`
	BuyerAppFinderFeeType   string    `json:"buyer_app_finder_fee_type"`
	BuyerAppFinderFeeAmount string    `json:"buyer_app_finder_fee_amount"`
	IsCollector             bool      `gorm:"type:boolean" json:"is_collector"`
	CreatedDate             time.Time `json:"created_date"`
	UpdatedDate             time.Time `json:"updated_date"`
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
	Company          *Company
	Profile          datatypes.JSON                           `json:"profile"`
	AddressDetails   datatypes.JSON                           `json:"address_details"`
	AccessTemplateId []access_template.AccessTemplateResponse `json:"access_template_id"`

	//PltPointIds   *PlatformPoints `json:"plt_point_ids" gorm:"foreignkey:PltPointId"`
	// CreatedBy []CoreUsers `gorm:"many2many:created_by"`
	IsEnabled bool `json:"is_enabled" gorm:"default:true"`
	IsActive  bool `json:"is_active" gorm:"default:true"`
}
type AppCategoryResponseDTO struct {
	ID           uint   `json:"id"`
	DisplayName  string `json:"display_name"`
	CategoryCode string `json:"category_code"`
}
type (
	ChannelStatusDTO struct {
		ChannelCode  string `json:"channel_code" gorm:"type:varchar(20)"`
		InternalId   uint   `json:"internal_id" gorm:"type:integer"`
		LookupCode   string `json:"lookup_code" gorm:"type:varchar(20)"`
		ExternalCode string `json:"external_code" gorm:"type:varchar(20)"`
		ExternalId   int    `json:"external_id" gorm:"type:integer"`
		model_core.Model
	}
	GetAllChannelStatusResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []ChannelStatusDTO
		}
	} //@ name GetAllChannelStatusResponse
)

// Create channel status response
type (
	ChannelStatusCreate struct {
		Created_id int
	}
	ChannelStatusCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data ChannelStatusCreate
		}
	} //@ name ChannelStatusCreateResponse
)

// Update channel status response
type (
	ChannelStatusUpdate struct {
		Updated_id int
	}
	ChannelStatusUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data ChannelStatusUpdate
		}
	} //@ name ChannelStatusUpdateResponse
)

// Get channel status response
type (
	ChannelStatusGet struct {
		model_core.ChannelLookupCodes
	}
	ChannelStatusGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data ChannelStatusGet
		}
	} //@ name ChannelStatusGetResponse
)

// Delete channel status response
type (
	ChannelStatusDelete struct {
		Deleted_id int
	}
	ChannelStatusDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data ChannelStatusDelete
		}
	} //@ name ChannelStatusDeleteResponse
)
