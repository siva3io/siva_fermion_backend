package core

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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
var Tables = []interface{}{
	&Company{},
	&CoreUsers{},
	&Tags{},
	&State{},
	&Country{},
	&Currency{},
	&L10n{},
	&InstalledApps{},
	&AppCategory{},
	&AppStore{},
	&Attachments{},
	&Lookuptype{},
	&Lookupcode{},
	&Access{},
	&AppsEdit{},
	&EunimartBaseSettings{},
	&PlatformPoints{},
	&NotificationTemplates{},
	&NotificationSettings{},
	&Localization{},
	&DBSchema{},
	&ExternalMapper{},
	&ViewSchema{},
	&ChannelLookupCodes{},
}

type Model struct {
	ID          uint  `json:"id" gorm:"primarykey"`
	IsEnabled   *bool `json:"is_enabled" gorm:"default:true"`
	IsActive    *bool `json:"is_active" gorm:"default:true"`
	CreatedByID *uint `json:"created_by" gorm:"column:created_by"`
	CreatedBy   *CoreUsers
	UpdatedDate time.Time `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedByID *uint     `json:"updated_by" gorm:"column:updated_by"`
	UpdatedBy   *CoreUsers
	DeletedByID *uint `json:"deleted_by" gorm:"column:deleted_by"`
	DeletedBy   *CoreUsers
	CreatedDate time.Time `json:"created_date" gorm:"autoCreateTime"`
	CompanyId   uint      `json:"company_id" gorm:"column:company_id"`
	AppId       *uint     `json:"app_id" gorm:"column:app_id;default:null"`
	App         *InstalledApps
	// Company     *Company
	DeletedAt gorm.DeletedAt `json:"-"`
}

type L10n struct {
	Name         string         `gorm:"type:varchar(50)" json:"name"`
	Params       datatypes.JSON `gorm:"type:json" json:"params"`
	LangCode     string         `gorm:"type:varchar" json:"lang_code"`
	LookupCodeId uint           `json:"lookup_code_id" gorm:"column:lookup_code_id"`
	LookupCode   Lookupcode
	Model
}

// user level table which has installed apps
type InstalledApps struct {
	Name              string         `gorm:"type:varchar(50); not null;" json:"name"`
	Code              string         `gorm:"type:varchar(50); unique; not null;" json:"code"`
	Icon              datatypes.JSON `json:"icon" gorm:"type:json"`
	Type              string         `gorm:"type:varchar(50)" json:"type"`
	CurrentVersion    string         `gorm:"type:varchar(50)" json:"current_version"`
	VersionCompatible datatypes.JSON `gorm:"type:json" json:"version_compatible"`
	AccessToken       string         `gorm:"type:text" json:"access_token"`
	AppDetails        datatypes.JSON `gorm:"type:json" json:"app_details"`
	AppAccessTemplate datatypes.JSON `gorm:"type:json" json:"app_access_template"`
	IsCore            bool           `gorm:"type:boolean" json:"is_core"`
	DeveloperIds      datatypes.JSON `gorm:"type:json" json:"developer_ids"`
	Tutorials         datatypes.JSON `gorm:"type:json" json:"tutorials"`
	Gamification      datatypes.JSON `gorm:"type:json" json:"gamification"`
	Schedulers        datatypes.JSON `gorm:"type:json" json:"schedulers"`
	ConcurrencyList   datatypes.JSON `gorm:"type:json" json:"concurrency_list"`
	CategoryId        *uint          `json:"category_id"`
	CategoryName      string         `json:"category_name"`
	ID                uint           `json:"id" gorm:"primarykey"`
	IsEnabled         bool           `json:"is_enabled" gorm:"default:true"`
	IsActive          bool           `json:"is_active" gorm:"default:true"`
	UpdatedDate       time.Time      `json:"updated_date" gorm:"autoUpdateTime"`
	CreatedDate       time.Time      `json:"created_date" gorm:"autoCreateTime"`
	CreatedByID       *uint          `json:"created_by" gorm:"column:created_by"`
	UpdatedByID       *uint          `json:"updated_by" gorm:"column:updated_by"`
	AppServices       datatypes.JSON `gorm:"type:json;default:'[]'" json:"app_services"`
	CreatedBy         *CoreUsers
	UpdatedBy         *CoreUsers
}

// global table which has all the apps available to install
type AppStore struct {
	Name           string         `json:"name" gorm:"type:varchar(100); not null;"`
	Code           string         `json:"code" gorm:"type:varchar(100); unique"`
	Icon           datatypes.JSON `json:"icon" gorm:"type:json"`
	Description    string         `json:"description" gorm:"text"`
	VersionHistory string         `json:"version_history" gorm:"text"`
	ImageOptions   datatypes.JSON `json:"image_options"`
	Version        string         `json:"version" gorm:"type:varchar(100)"`
	ReleaseDate    datatypes.Date `json:"release_date" gorm:"date"`
	TotalInstalls  int64          `json:"total_installs"`
	Publisher      datatypes.JSON `json:"publisher" gorm:"type:json"`
	Installation   datatypes.JSON `json:"installation" gorm:"type:json"`
	FAQs           datatypes.JSON `json:"faqs" gorm:"type:json"`
	Support        datatypes.JSON `json:"support" gorm:"type:json"`
	CurrencyId     *uint          `json:"currency_id"`
	Currency       Currency       `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	Price          float64        `json:"price" gorm:"type:double precision"`
	License        string         `json:"license" gorm:"type:varchar(100)"`
	Website        string         `json:"website" gorm:"type:varchar(100)"`
	SourceLink     string         `json:"source_link" gorm:"type:varchar(100)"`
	Reviews        datatypes.JSON `json:"reviews" gorm:"type:json"`
	CategoryId     *uint          `json:"category_id"`
	CategoryName   string         `json:"category_name"`
	Category       *AppCategory   `gorm:"foreignkey:CategoryId; references:ID"`
	AppServices    datatypes.JSON `gorm:"type:json;default:'[]'" json:"app_services"`

	//Ratings
	//Language
	//Tags
	Model
}

type AppCategory struct {
	DisplayName  string         `gorm:"type:varchar(50); not null;" json:"display_name"`
	CategoryCode string         `json:"category_code"`
	Description  datatypes.JSON `json:"description"`
	//AppIds      []AppStore     `json:"app_ids"`
	AppCount int          `json:"app_count"`
	ParentId uint         `json:"parent_id"`
	Parent   *AppCategory `json:"parent"`
	//ChildIds    []AppCategory  `json:"child_ids"`
	Sequence int  `json:"sequence"`
	Visible  bool `json:"visible"`
	gorm.Model
}

type Lookuptype struct {
	LookupType  string       `gorm:"type:varchar(50); unique; not null;" json:"lookup_type"`
	DisplayName string       `gorm:"type:varchar(50); not null;" json:"display_name"`
	Lookupcodes []Lookupcode `gorm:"foreignkey:LookupTypeId; references:ID;"`
	IsEnabled   *bool        `json:"is_enabled" gorm:"default:true"`
	gorm.Model
}

type Lookupcode struct {
	LookupTypeId uint ` gorm:"type:integer; not null; uniqueIndex:idx_lookup_type_id_lookup_code;" json:"lookup_type_id"`
	LookupType   Lookuptype
	LookupCode   string `gorm:"type:varchar(50); not null; uniqueIndex:idx_lookup_type_id_lookup_code;" json:"lookup_code"`
	DisplayName  string `gorm:"type:varchar(50); not null;" json:"display_name"`
	IsEnabled    *bool  `json:"is_enabled" gorm:"default:true"`
	gorm.Model
}

type ViewSchema struct {
	Name     string         `json:"name" gorm:"type:varchar(50)"`
	AppsList datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"apps_list"`
	ViewList datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"view_list"`
	Model
}
type ExternalMapper struct {
	Name           string         `json:"name" gorm:"type:varchar(50)"`
	AppsList       datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"apps_list"`
	ServiceList    datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"service_list"`
	AuthToken      datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"auth_token"`
	RequestPayload datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"request_payload"`
	Response       datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"response"`
	ResponseMap    datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"response_map"`
	Model
}
type DBSchema struct {
	Name      string         `json:"name" gorm:"type:varchar(50)"`
	AppsList  datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"apps_list"`
	FieldList datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"field_list"`
	Model
}
type Localization struct {
	Name       string         `json:"name" gorm:"type:varchar(50)"`
	CountryIds datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"country_ids"`
	StateIds   datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"state_ids"`
	Details    datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"details"`
	Model
}

type PlatformPoints struct {
	Points       int            `json:"points" gorm:"type:integer"`
	Type         uint           `json:"type" gorm:"type:integer"`
	UserId       uint           `json:"user_id"`
	CoreUsers    CoreUsers      `gorm:"foreignkey:UserId"`
	IsRedeemable bool           `json:"is_redeemable" gorm:"type:boolean"`
	Redeemed     bool           `json:"redeemed" gorm:"type:boolean"`
	ExpiryDate   time.Time      `json:"expiry_date" gorm:"type:time"`
	Expired      bool           `json:"expired" gorm:"type:boolean"`
	Constraints  datatypes.JSON `json:"constraints" gorm:"type:json; default:'[]'; not null"`
	Model
}

type NotificationSettings struct {
	Name                 string         `json:"name" gorm:"type:varchar(50)"`
	EmailSMTP            datatypes.JSON `json:"email_smtp" gorm:"type:json; default:'[]'; not null"`
	EmailPOP             datatypes.JSON `json:"email_pop" gorm:"type:json; default:'[]'; not null"`
	EmailDomain          string         `json:"email_domain" gorm:"type:varchar(50)"`
	SmsSettings          datatypes.JSON `json:"sms_settings" gorm:"type:json; default:'[]'; not null"`
	BrowserNotifications datatypes.JSON `json:"browser_notifications" gorm:"type:json; default:'[]'; not null"`
	FirebaseSettings     datatypes.JSON `json:"firebase_settings" gorm:"type:json; default:'[]'; not null"`
	MobilePushSettings   datatypes.JSON `json:"mobile_push_settings" gorm:"type:json; default:'[]'; not null"`
	PixelSettings        datatypes.JSON `json:"pixel_settings" gorm:"type:json; default:'[]'; not null"`
	TagManagerSettings   datatypes.JSON `json:"tag_manager_settings" gorm:"type:json; default:'[]'; not null"`
	HotjarSettings       datatypes.JSON `json:"hotjar_settings" gorm:"type:json; default:'[]'; not null"`
	AppsFlyerSettings    datatypes.JSON `json:"apps_flyer_settings" gorm:"type:json; default:'[]'; not null"`
	BranchSettings       datatypes.JSON `json:"branch_settings" gorm:"type:json; default:'[]'; not null"`
	ThirdPartySettings   datatypes.JSON `json:"third_party_settings" gorm:"type:json; default:'[]'; not null"`
	AppHeaderSettings    datatypes.JSON `json:"app_header_settings" gorm:"type:json; default:'[]'; not null"`
	ID                   uint           `json:"id" gorm:"primarykey"`
	// Model
}

type NotificationTemplates struct {
	Name           string         `json:"name" gorm:"type:varchar(50)"`
	Type           uint           `json:"type" gorm:"type:integer"`
	AppIds         pq.Int64Array  `json:"app_ids" gorm:"type:int[]"`
	ModelIds       pq.Int64Array  `json:"model_ids" gorm:"type:int[]"`
	ViewIds        pq.Int64Array  `json:"view_ids" gorm:"type:int[]"`
	ViewSectionIds pq.Int64Array  `json:"view_section_ids" gorm:"type:int[]"`
	ControllerIds  pq.Int64Array  `json:"controller_ids" gorm:"type:int[]"`
	TriggerRule    datatypes.JSON `json:"trigger_rule" gorm:"type:json; default:'[]'; not null"`
	Model
}

type EunimartBaseSettings struct {
	Name               string         `gorm:"type:varchar(50)" json:"name"`
	MenuHierarchy      datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"menu_hierarchy"`
	Theme              datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"theme"`
	AppsList           datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"apps_list"`
	ConcurrencyList    datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"concurrency_list"`
	DecimalSettings    datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"decimal_settings"`
	RoundOff           datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"round_off"`
	SchedulerSettings  datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"scheduler_settings"`
	SecuritySettings   datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"security_settings"`
	ControllerSettings datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"controller_settings"`
	DomainSettings     datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"domain_settings"`
	CdnSettings        datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"cdn_settings"`
	SslSettings        datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"ssl_settings"`
	MiscSettings       datatypes.JSON `gorm:"type:json; default:'[]'; not null" json:"misc_settings"`
	Model
}

type Company struct {
	ID                       uint           `json:"id" gorm:"primarykey"`
	Name                     string         `gorm:"type:text;unique; not null;" json:"company_name"`
	Addresses                datatypes.JSON `gorm:"type:json;default:'[]'" json:"addresses"`
	Phone                    string         `gorm:"varchar(50)" json:"phone"`
	Email                    string         `gorm:"type:varchar(50);not null" json:"email"`
	CompanyDetails           datatypes.JSON `gorm:"type:json;default:'{}';not null" json:"company_details"`
	IsEnterpise              bool           `json:"is_enterpise"`
	ParentId                 uint           `gorm:"type:integer" json:"parent_id"`
	ChildIds                 pq.Int64Array  `json:"child_ids" gorm:"type:int[]"`
	Type                     int            `json:"type"`
	CompanyDefaults          datatypes.JSON `gorm:"type:json;default:'{}';not null" json:"company_defaults"`
	PltPointIds              pq.Int32Array  `json:"plt_point_ids" gorm:"type:int[]"`
	TotalPoints              int            `json:"total_points" gorm:"type:integer"`
	Constraints              datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"constraints"`
	Schedulers               datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"schedulers"`
	QueueServices            datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"queue_services"`
	NotificationSettingsId   *uint          `json:"notification_settings_id" gorm:"column:notification_settings_id"`
	NotificationSettings     NotificationSettings
	NotificationTemplatesIds pq.Int32Array  `json:"notification_template_id" gorm:"type:int[]"`
	MenuHierarchy            datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"menu_hierarchy"`
	IsEnabled                bool           `json:"is_enabled" gorm:"default:true"`
	IsActive                 bool           `json:"is_active" gorm:"default:true"`
	CreatedByID              *uint          `json:"created_by" gorm:"column:created_by"`
	CreatedBy                *CoreUsers
	UpdatedDate              time.Time `json:"updated_date" gorm:"autoCreateTime"`
	UpdatedByID              *uint     `json:"updated_by" gorm:"column:updated_by"`
	UpdatedBy                *CoreUsers
	CreatedDate              time.Time      `json:"created_date" gorm:"autoUpdateTime"`
	OrganizationDetails      datatypes.JSON `json:"organization_details" gorm:"type:json;default:'{}'"`
	KycDocuments             datatypes.JSON `json:"kyc_documents" gorm:"type:json;default:'[]'"`
	FilePreferenceID         *uint          `json:"file_preference_id"`
	FilePreference           Lookupcode     `json:"file_preference" gorm:"foreignkey:FilePreferenceID; references:ID"`
	InvoiceGenerationId      *uint          `json:"invoice_generation_id"`
	InvoiceGeneration        Lookupcode     `json:"invoice_generation" gorm:"foreignkey:InvoiceGenerationId; references:ID"`
	BusinessTypeID           *uint          `json:"business_type_id"`
	BusinessType             Lookupcode     `json:"business_type" gorm:"foreignkey:BusinessTypeID; references:ID"`
	// PltPointId      *uint           `json:"plt_point_id"`
}

type AppsEdit struct {
	AppsId      int            `gorm:"type:varchar(100)" json:"apps_id"`
	UserId      int            `gorm:"type:varchar(100)" json:"user_id"`
	EditHistory datatypes.JSON `gorm:"type:json" json:"edit_history"`
	Model
}

type Access struct {
	Name            string         `json:"name"`
	ParentId        int            `gorm:"type:integer" json:"parent_id"`
	ChildIds        datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"child_ids"`
	ParentModuleIds datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"parent_module_ids"`
	Models          datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"models"`
	Views           datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"views"`
	Controllers     datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"controllers"`
	CompanyIds      datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"company_ids"`
	UserIds         datatypes.JSON `gorm:"type:json;default:'[]';not null" json:"user_ids"`
	Model
}

type Country struct {
	Name                       string         `json:"name" gorm:"type:varchar(100); not null;"`
	DefaultAccountingPrinciple datatypes.JSON `json:"default_accounting_principle" gorm:"type:json"`
	CountryCode                string         `json:"country_code" gorm:"type:varchar(50); unique; not null;"`
	CountryCode2               string         `json:"country_code2" gorm:"type:varchar(50); unique; null;"`
	StateIds                   pq.Int64Array  `json:"state_ids" gorm:"type:int[]"`
	TimezoneIds                pq.Int64Array  `json:"time_zone_ids" gorm:"type:int[]"`
	IsDst                      bool           `json:"is_dst" gorm:"type:boolean"`
	CurrencyId                 uint           `json:"currency_id" gorm:"column:currency_id"`
	Currency                   Currency
	Model
}

type State struct {
	Name      string `json:"name" gorm:"type:varchar(100); not null;"`
	StateCode string `json:"state_code" gorm:"type:varchar(50); unique; not null"`
	CountryId uint   `json:"country_id" gorm:"column:country_id; not null"`
	Country   Country
	Model
}

type Currency struct {
	Name                string         `json:"name" gorm:"type:varchar(20); not null;"`
	CurrencySymbol      string         `json:"currency_symbol" gorm:"type:varchar(20); not null;"`
	CurrencyCode        string         `json:"currency_code" gorm:"type:varchar(10); unique; not null;"`
	IsBaseCurrency      bool           `json:"is_base_currency" gorm:"type:boolean"`
	ExchangeRate        float64        `json:"exchange_rate" gorm:"type:float"`
	ExchangeRateHistory datatypes.JSON `json:"exchange_rate_history" gorm:"type:json"`
	AutoUpdateEr        datatypes.JSON `json:"auto_update_er" gorm:"type:json"`
	Model
}
type Tags struct {
	Name    string         `json:"name" gorm:"type:varchar(100)"`
	Details datatypes.JSON `json:"details" gorm:"type:json"`
	Model
}
type Attachments struct {
	Name         string         `json:"name" gorm:"type:varchar(100)"`
	Description  string         `json:"description" gorm:"type:text"`
	Details      datatypes.JSON `json:"details" gorm:"type:json"`
	Url          string         `json:"url" gorm:"type:varchar(50)"`
	FilePath     string         `json:"path" gorm:"type:varchar(100)"`
	TypeId       uint           `json:"lookup_type_id" gorm:"column:lookup_type_id"`
	LookupType   Lookuptype
	ThumbnailUrl string  `json:"thumbnail_url" gorm:"type:varchar(50)"`
	Size         float64 `json:"size" gorm:"type:float"`
	Duration     string  `json:"duration" gorm:"type:varchar(100)"`
	Model
}
type ChannelLookupCodes struct {
	ChannelCode  string `json:"channel_code" gorm:"type:varchar(20)"`
	InternalId   uint   `json:"internal_id" gorm:"type:integer"`
	LookupCode   string `json:"lookup_code" gorm:"type:varchar(20)"`
	ExternalCode string `json:"external_code" gorm:"type:varchar(20)"`
	ExternalId   int    `json:"external_id" gorm:"type:integer"`
	Model
}
