package core

import (
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
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
var CoreUsersTables = []interface{}{
	&CoreUsers{},
}

type Hours int64

type CoreUsers struct {
	ID               uint              `json:"id" gorm:"primarykey"`
	Name             string            `json:"name"`
	FirstName        string            `json:"first_name"`
	LastName         string            `json:"last_name"`
	Username         string            `json:"username"`
	Email            string            `json:"email"`
	WorkEmail        string            `json:"work_email"`
	MobileNumber     string            `json:"mobile_number"`
	Password         string            `json:"password"`
	LoginType        int32             `json:"login_type"`
	Auth             datatypes.JSON    `gorm:"type:json;default:'{}';not null" json:"auth"`
	FaConf           datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"2fa_conf"`
	DeviceIds        datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"device_ids"`
	Preferences      datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"preferences"`
	Tutorials        datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"tutorials"`
	Gamification     datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"gamification"`
	TotalPoints      int               `json:"total_points"`
	MenuHierarchy    datatypes.JSON    `gorm:"type:json;default:'{}';not null" json:"menu_hierarchy"`
	Constraints      datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"constraints"`
	Schedulers       datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"schedulers"`
	QueueServices    datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"queue_services"`
	AccessIds        datatypes.JSON    `gorm:"type:json;default:'[]';not null" json:"access_ids"`
	CompanyId        *uint             `json:"company_id" gorm:"column:company_id"`
	Company          Company           `json:"company" gorm:""`
	Profile          datatypes.JSON    `gorm:"type:json;default:'{}';not null" json:"profile"`
	PltPointId       pq.Int32Array     `json:"plt_point_id" gorm:"type:int[]"`
	AccessTemplateId []*AccessTemplate `json:"access_template_id" gorm:"many2many:user_template_mapper;ForeignKey:id;References:id"`
	ExternalDetails  datatypes.JSON    `json:"external_details" gorm:"type:json;default:'{}';not null"`
	//PltPointIds   *PlatformPoints `json:"plt_point_ids" gorm:"foreignkey:PltPointId"`
	// CreatedBy []CoreUsers `gorm:"many2many:created_by"`
	IsEnabled   bool      `json:"is_enabled" gorm:"default:true"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	UpdatedDate time.Time `json:"updated_date" gorm:"autoUpdateTime"`
	CreatedDate time.Time `json:"created_date" gorm:"autoCreateTime"`
	CreatedByID *uint     `json:"created_by" gorm:"column:created_by"`
	UpdatedByID *uint     `json:"updated_by" gorm:"column:updated_by"`
	CreatedBy   *CoreUsers
	UpdatedBy   *CoreUsers
	UserTypes   datatypes.JSON `gorm:"type:json;default:'{}';not null" json:"user_types"`
	RoleId      *uint          `json:"role_id" gorm:"type:int"`
}

type Auth struct {
	OtpToken string `json:"otp_token"`
}

func (u *CoreUsers) HashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(bytes)
}

func (u *CoreUsers) GenerateToken(access_template_id uint, duration Hours) (string, error) {
	var (
		jwtKey = os.Getenv("JWT_KEY")
	)
	if access_template_id == 0 {
		access_template_id = 2
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":                 u.ID,
		"Username":           u.Username,
		"first_name":         u.FirstName,
		"last_name":          u.LastName,
		"company_id":         u.CompanyId,
		"access_template_id": access_template_id,
		"user_types":         u.UserTypes,
		"role_id":            u.RoleId,
		"exp":                time.Now().Add(time.Hour * time.Duration(duration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, err
}
