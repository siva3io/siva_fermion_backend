package core

import (
	"time"

	"gorm.io/gorm"
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
type AccessTemplate struct {
	// Model
	Name        string `json:"name"`
	ID          uint   `json:"id" gorm:"primarykey"`
	IsEnabled   *bool  `json:"is_enabled" gorm:"default:true"`
	IsActive    *bool  `json:"is_active" gorm:"default:true"`
	CreatedByID *uint  `json:"created_by" gorm:"column:created_by"`
	// CreatedBy   *CoreUsers
	UpdatedDate time.Time `json:"updated_date" gorm:"autoUpdateTime"`
	UpdatedByID *uint     `json:"updated_by" gorm:"column:updated_by"`
	// UpdatedBy   *CoreUsers
	DeletedByID *uint `json:"deleted_by" gorm:"column:deleted_by"`
	// DeletedBy   *CoreUsers
	CreatedDate time.Time `json:"created_date" gorm:"autoCreateTime"`
	CompanyId   uint      `json:"company_id" gorm:"column:company_id"`
	AppId       *uint     `json:"app_id" gorm:"column:app_id;default:null"`
	// App         *InstalledApps
	// Company     *Company
	DeletedAt gorm.DeletedAt `json:"-"`
}
