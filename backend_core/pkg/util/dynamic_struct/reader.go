package dynamic_struct

import (
	"fmt"
	"reflect"
	"time"
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

type DynamicStructReader interface {
	HasField(name string) bool
	GetField(name string) Field
	GetAllFields() []Field
	GetValue() interface{}
}
type dynamicStructReader struct {
	fields map[string]fieldStruct
	value  interface{}
}
type Field interface {
	Name() string
	//===========integer converstion types =================
	Int() int
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64

	PointerInt() *int
	PointerInt8() *int8
	PointerInt16() *int16
	PointerInt32() *int32
	PointerInt64() *int64

	//===========Uint converstion types =================
	Uint() uint
	Uint8() uint8
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64

	PointerUint() *uint
	PointerUint8() *uint8
	PointerUint16() *uint16
	PointerUint32() *uint32
	PointerUint64() *uint64

	//===========float converstion types =================
	Float32() float32
	Float64() float64

	PointerFloat32() *float32
	PointerFloat64() *float64

	//===========string converstion types =================
	String() string
	PointerString() *string

	//===========boolean converstion types =================
	Bool() bool
	PointerBool() *bool

	//===========time converstion types =================
	Time() time.Time
	PointerTime() *time.Time

	MapStringOfInterface() map[string]interface{}
	Interface() interface{}
}

type fieldStruct struct {
	field reflect.StructField
	value reflect.Value
}

func NewDynamicStructReader(value interface{}) DynamicStructReader {
	fields := map[string]fieldStruct{}

	valueOf := reflect.Indirect(reflect.ValueOf(value))
	typeOf := valueOf.Type()
	fmt.Println()
	fmt.Println("============== Struct Model ================================================================================================================")
	fmt.Println(typeOf)
	fmt.Println("============================================================================================================================================")
	fmt.Println("============== Struct Values ===============================================================================================================")
	fmt.Println(valueOf)
	fmt.Println("============================================================================================================================================")
	fmt.Println()
	if typeOf.Kind() == reflect.Struct {
		for i := 0; i < valueOf.NumField(); i++ {
			field := typeOf.Field(i)
			fields[field.Name] = fieldStruct{
				field: field,
				value: valueOf.Field(i),
			}
		}
	} else {
		fmt.Println("warning :- input is not a  dynamic struct")
	}

	return dynamicStructReader{
		fields: fields,
		value:  value,
	}
}

func (r dynamicStructReader) HasField(name string) bool {
	_, ok := r.fields[name]
	return ok
}

func (r dynamicStructReader) GetField(name string) Field {
	if !r.HasField(name) {
		return nil
	}
	return r.fields[name]
}

func (r dynamicStructReader) GetAllFields() []Field {
	var fields []Field

	for _, field := range r.fields {
		fields = append(fields, field)
	}

	return fields
}

func (r dynamicStructReader) GetValue() interface{} {
	return r.value
}

func (f fieldStruct) Name() string {
	return f.field.Name
}

func (f fieldStruct) PointerInt() *int {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int()
	return &value
}

func (f fieldStruct) Int() int {
	return int(reflect.Indirect(f.value).Int())
}

func (f fieldStruct) PointerInt8() *int8 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int8()
	return &value
}

func (f fieldStruct) Int8() int8 {
	return int8(reflect.Indirect(f.value).Int())
}

func (f fieldStruct) PointerInt16() *int16 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int16()
	return &value
}

func (f fieldStruct) Int16() int16 {
	return int16(reflect.Indirect(f.value).Int())
}

func (f fieldStruct) PointerInt32() *int32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int32()
	return &value
}

func (f fieldStruct) Int32() int32 {
	return int32(reflect.Indirect(f.value).Int())
}

func (f fieldStruct) PointerInt64() *int64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int64()
	return &value
}

func (f fieldStruct) Int64() int64 {
	return reflect.Indirect(f.value).Int()
}

func (f fieldStruct) PointerUint() *uint {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint()
	return &value
}

func (f fieldStruct) Uint() uint {
	return uint(reflect.Indirect(f.value).Uint())
}

func (f fieldStruct) PointerUint8() *uint8 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint8()
	return &value
}

func (f fieldStruct) Uint8() uint8 {
	return uint8(reflect.Indirect(f.value).Uint())
}

func (f fieldStruct) PointerUint16() *uint16 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint16()
	return &value
}

func (f fieldStruct) Uint16() uint16 {
	return uint16(reflect.Indirect(f.value).Uint())
}

func (f fieldStruct) PointerUint32() *uint32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint32()
	return &value
}

func (f fieldStruct) Uint32() uint32 {
	return uint32(reflect.Indirect(f.value).Uint())
}

func (f fieldStruct) PointerUint64() *uint64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint64()
	return &value
}

func (f fieldStruct) Uint64() uint64 {
	return reflect.Indirect(f.value).Uint()
}

func (f fieldStruct) PointerFloat32() *float32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Float32()
	return &value
}

func (f fieldStruct) Float32() float32 {
	return float32(reflect.Indirect(f.value).Float())
}

func (f fieldStruct) PointerFloat64() *float64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Float64()
	return &value
}

func (f fieldStruct) Float64() float64 {
	return reflect.Indirect(f.value).Float()
}

func (f fieldStruct) PointerString() *string {
	if f.value.IsNil() {
		return nil
	}
	value := f.String()
	return &value
}

func (f fieldStruct) String() string {
	return reflect.Indirect(f.value).String()
}

func (f fieldStruct) PointerBool() *bool {
	if f.value.IsNil() {
		return nil
	}
	value := f.Bool()
	return &value
}

func (f fieldStruct) Bool() bool {
	return reflect.Indirect(f.value).Bool()
}

func (f fieldStruct) PointerTime() *time.Time {
	if f.value.IsNil() {
		return nil
	}
	value := f.Time()
	return &value
}

func (f fieldStruct) Time() time.Time {
	value, ok := reflect.Indirect(f.value).Interface().(time.Time)
	if !ok {
		panic(fmt.Sprintf(`field "%s" is not instance of time.Time`, f.field.Name))
	}

	return value
}

func (f fieldStruct) MapStringOfInterface() map[string]interface{} {
	return f.value.Interface().(map[string]interface{})
}

func (f fieldStruct) Interface() interface{} {
	return f.value.Interface()
}
