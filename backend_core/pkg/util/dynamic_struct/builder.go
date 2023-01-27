package dynamic_struct

import "reflect"

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

type Builder interface {
	AddField(fieldNameValue string, fieldTypeValue interface{}, fieldTagValue string) Builder
	RemoveField(fieldName string) Builder
	HasField(fieldName string) bool
	GetField(fieldName string) FieldConfigurationSettings
	Build() DynamicStruct
}

// ================== builder ===============================
type builder struct {
	Fields []*FieldConfiguration
}
type FieldConfiguration struct {
	FieldName string
	FieldType interface{}
	FieldTag  string
}
type FieldConfigurationSettings interface {
	SetNewFieldType(typ interface{}) FieldConfigurationSettings
	SetNewFieldTag(tag string) FieldConfigurationSettings
}

// ================== dynamic struct formatter ===============================
type DynamicStruct interface {
	NewStruct() interface{}
	NewArrayOfStruct() interface{}
}
type dynamicStruct struct {
	definition reflect.Type
}

func NewDynamicStructBuilder() Builder {
	return &builder{
		Fields: []*FieldConfiguration{},
	}
}

// =====================Create Dynamic strict ====================================================
func (b *builder) AddField(fieldNameValue string, fieldTypeValue interface{}, fieldTagValue string) Builder {
	b.Fields = append(b.Fields, &FieldConfiguration{
		FieldName: fieldNameValue,
		FieldType: fieldTypeValue,
		FieldTag:  fieldTagValue,
	})

	return b
}
func (b *builder) RemoveField(fieldName string) Builder {
	for i := range b.Fields {
		if b.Fields[i].FieldName == fieldName {
			b.Fields = append(b.Fields[:i], b.Fields[i+1:]...)
			break
		}
	}
	return b
}
func (b *builder) HasField(fieldName string) bool {
	for i := range b.Fields {
		if b.Fields[i].FieldName == fieldName {
			return true
		}
	}
	return false
}
func (b *builder) GetField(fieldName string) FieldConfigurationSettings {
	for i := range b.Fields {
		if b.Fields[i].FieldName == fieldName {
			return b.Fields[i]
		}
	}
	return nil
}
func (b *builder) Build() DynamicStruct {
	var structFields []reflect.StructField

	for _, field := range b.Fields {
		structFields = append(structFields, reflect.StructField{
			Name: field.FieldName,
			Type: reflect.TypeOf(field.FieldType),
			Tag:  reflect.StructTag(field.FieldTag),
		})
	}

	return &dynamicStruct{
		definition: reflect.StructOf(structFields),
	}
}

// ======================= Fiels Configuration Settings ==============================================
func (f *FieldConfiguration) SetNewFieldType(fieldNewType interface{}) FieldConfigurationSettings {
	f.FieldType = fieldNewType
	return f
}
func (f *FieldConfiguration) SetNewFieldTag(fieldNewTag string) FieldConfigurationSettings {
	f.FieldTag = fieldNewTag
	return f
}

// ================create object or array  dynamic struct =======================================
func (ds *dynamicStruct) NewStruct() interface{} {
	return reflect.New(ds.definition).Interface()
}
func (ds *dynamicStruct) NewArrayOfStruct() interface{} {
	return reflect.New(reflect.SliceOf(ds.definition)).Interface()
}
