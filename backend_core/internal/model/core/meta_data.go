package core

import "gorm.io/gorm"

var MetaTables = []interface{}{
	&IRModel{},
	&IRModelFields{},
	&IRModelConstraints{},
}

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
type (
	IRModel struct {
		SchemaName  string               `json:"schema_name"`
		TableName   string               `json:"table_name" gorm:"not null"`
		Fields      []IRModelFields      `json:"fields" gorm:"foreignkey:TableId;references:ID"`
		Constraints []IRModelConstraints `json:"constraints" gorm:"foreignkey:TableId;references:ID"`
		gorm.Model
	}

	IRModelFields struct {
		TableId                uint   `json:"table_id"`
		TableName              string `json:"table_name" gorm:"not null"`
		ColumnName             string `json:"column_name" gorm:"not null"`
		DataType               string `json:"data_type"`
		CharacterMaximumLength string `json:"character_maximum_length"`
		IsNullable             string `json:"is_nullable"`
		ColumnDefault          string `json:"column_default"`
		gorm.Model
	}

	IRModelConstraints struct {
		TableId        uint   `json:"table_id"`
		ConstraintName string `json:"constraint_name" gorm:"not null"`
		TableName      string `json:"table_name" gorm:"not null"`
		ConstraintType string `json:"constraint_type"`
		gorm.Model
	}
)
