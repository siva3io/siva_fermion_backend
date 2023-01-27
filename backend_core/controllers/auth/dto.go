package auth

import (
	access_template "fermion/backend_core/controllers/access/module"
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
// -------------Basic Authentication DTO --------------------------------
type RegisterDTO struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type UserProfileUpdateDTO struct {
	ID                uint   `json:"id" gorm:"primarykey"`
	Name              string `json:"name"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	WorkEmail         string `json:"work_email"`
	EmailNotification bool   `json:"email_notification"`
	SMSNotification   bool   `json:"sms_notification"`
	Reminder          bool   `json:"reminder"`
	MobileNumber      string `json:"mobile_number"`
	Password          string `json:"password"`
	ProfileDpImage    string `json:"profile_image"`
	TimeFormatId      int    `json:"time_format_id"`
	DateFormatId      int    `json:"date_format_id"`
	TimeZoneId        int    `json:"time_zone_id"`

	TwoStepVerificaion bool `json:"two_step_verification"`
	// Company           CompanyDTO
}

//	type CompanyDTO struct {
//		Name           string         `json:"name"`
//		Addresses      datatypes.JSON `json:"addresses"`
//		Phone          string         `json:"phone"`
//		Email          string         `json:"email"`
//		CompanyDetails datatypes.JSON `json:"company_details"`
//		IsEnterpise    bool           `json:"is_enterpise"`
//	}

// register response data
type RegisterRespData struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ------------swagger documentation---------------------
// register response doc
type RegisterResponseDOC struct {
	Meta response.MetaResponse `json:"meta"`
	Data RegisterRespData      `json:"data"`
}

// login response data
type LoginResponseDTO struct {
	Token string `json:"token"`
}

// for response login
type LoginResponseDOC struct {
	Meta response.MetaResponse `json:"meta"`
	Data LoginResponseDTO      `json:"data"`
}

// ---------------------OTP---------------------------------
type UserLoginDTO struct {
	Login string `json:"login"`
}

type UserLoginResponse struct {
	ID      uint   `json:"id"`
	Status  string `json:"status"`
	NewUser bool   `json:"new_user"`
	OtpSend bool   `json:"otp_sent"`
}

type VerifyOtpDTO struct {
	ID               uint   `json:"id"`
	Otp              string `json:"otp"`
	AccessTemplateId uint   `json:"access_template_id"`
}
type VerifyOtpResponse struct {
	Status           bool            `json:"status"`
	Message          string          `json:"message"`
	Token            string          `json:"token,omitempty"`
	AccessTemplateId uint            `json:"access_template_id"`
	User             UserResponseDTO `json:"user,omitempty"`
}

type UserObject struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	MobileNumber string         `json:"mobile_number"`
	TotalPoints  int            `json:"total_points"`
	LoginType    int32          `json:"login_type"`
	IsEnabled    bool           `json:"is_enabled"`
	IsActive     bool           `json:"is_active"`
	UserTypes    datatypes.JSON `json:"user_types,omitempty"`
}

type UpdateProfileDTO struct {
	ID             uint           `json:"id"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	WorkEmail      string         `json:"work_email"`
	MobileNumber   string         `json:"mobile_number"`
	Password       string         `json:"password"`
	LoginType      int32          `json:"login_type"`
	Profile        datatypes.JSON `json:"profile,omitempty"`
	CompanyName    string         `json:"company_name"`
	CompanyEmail   string         `json:"company_email"`
	FaConf         datatypes.JSON `json:"2fa_conf,omitempty"`
	DeviceIds      datatypes.JSON `json:"device_ids,omitempty"`
	Preferences    datatypes.JSON `json:"preferences,omitempty"`
	AccessIds      datatypes.JSON `json:"access_ids,omitempty"`
	UserTypes      datatypes.JSON `json:"user_types,omitempty"`
	AddressDetails datatypes.JSON `json:"address_details,omitempty"`
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
	UserTypes        datatypes.JSON `json:"user_types"`
	LoginType        int32          `json:"login_type"`
	Auth             datatypes.JSON `json:"auth"`
	FaConf           datatypes.JSON `json:"2fa_conf"`
	DeviceIds        datatypes.JSON `json:"device_ids"`
	Preferences      datatypes.JSON `json:"preferences"`
	AccessIds        []uint         `json:"access_ids"`
	CompanyId        *uint          `json:"company_id"`
	TeamHead         string         `json:"team_head"`
	ExternalDetails  datatypes.JSON `json:"external_details"`
	Company          CompanyDTO
	Profile          datatypes.JSON                           `json:"profile"`
	AddressDetails   datatypes.JSON                           `json:"address_details"`
	AccessTemplateId []access_template.AccessTemplateResponse `json:"access_template_id"`

	//PltPointIds   *PlatformPoints `json:"plt_point_ids" gorm:"foreignkey:PltPointId"`
	// CreatedBy []CoreUsers `gorm:"many2many:created_by"`
	IsEnabled bool   `json:"is_enabled" gorm:"default:true"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	Token     string `json:"token"`
}
type CompanyDTO struct {
	Name                   string         `json:"name"`
	Addresses              datatypes.JSON `json:"addresses"`
	Phone                  string         `json:"phone"`
	Email                  string         `json:"email"`
	CompanyDetails         datatypes.JSON `json:"company_details"`
	IsEnterpise            bool           `json:"is_enterpise"`
	ParentId               uint           `json:"parent_id"`
	ChildIds               pq.Int64Array  `json:"child_ids"`
	Type                   int            `json:"type"`
	CompanyDefaults        datatypes.JSON `json:"company_defaults,omitempty"`
	NotificationSettingsId *uint          `json:"notification_settings_id"`
	// NotificationSettings     NotificationSettings
	NotificationTemplatesIds pq.Int32Array `json:"notification_template_id"`
	IsEnabled                bool          `json:"is_enabled"`
}
type UserVisitDTO struct {
	Name         string `json:"name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
	CompanyId    *uint  `json:"company_id"`
}
