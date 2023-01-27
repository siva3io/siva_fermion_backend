package helpers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	Pagination "fermion/backend_core/internal/model/pagination"

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
func Paginate(value interface{}, pagination *Pagination.Paginatevalue, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	fields, values, err := ApplyFilter(pagination.Filters)
	if err != nil {
		panic("filter error")
	}

	sortfields, err := ApplySort(pagination.Sort)
	if err != nil {
		fmt.Println(err)
		panic("Sort error")
	}

	per_page := pagination.GetLimit()

	joinedFilter := strings.Join(fields, " AND ")
	if pagination.FilterType == "OR" || pagination.FilterType == "or" {
		joinedFilter = strings.Join(fields, " OR ")
	}

	if per_page == -1 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Model(value).Order(sortfields).Where(joinedFilter, values...)
		}
	}

	pagination.Filters = ""
	pagination.Sort = ""
	pagination.FilterType = ""

	db.Model(value).Where(joinedFilter, values...).Count(&pagination.TotalRows)

	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(per_page)))

	return func(db *gorm.DB) *gorm.DB {
		return db.Model(value).Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(sortfields).Where(joinedFilter, values...)
	}
}

func ApplyFilter(filters string) ([]string, []interface{}, error) {
	var filterArray [][]interface{}
	var fields []string
	var values []interface{}
	var value interface{}

	if filters == "" {
		return fields, values, nil
	}
	err := json.Unmarshal([]byte(filters), &filterArray)

	if err != nil {
		return fields, values, err
	}

	for _, filter := range filterArray {
		field := fmt.Sprintf("%v %v ?", filter[0], filter[1])
		if filter[1] == "ILIKE" || filter[1] == "ilike" {
			filter[2] = "%" + filter[2].(string) + "%"
		}
		value = fmt.Sprintf("%v", filter[2])

		//this "in" filter query converts each individual filter value to string type
		if filter[1] == "IN" || filter[1] == "in" {
			valStr := fmt.Sprint(filter[2])
			value = strings.Split(valStr[1:len(valStr)-1], " ")
		}

		fields = append(fields, field)
		values = append(values, value)
	}
	return fields, values, nil
}

func ApplySort(sort string) (string, error) {

	if sort == "" {
		return "Id desc", nil
	}
	var sortArray [][]interface{}
	err := json.Unmarshal([]byte(sort), &sortArray)
	if err != nil {
		return "", err
	}
	var fields []string
	for _, sort := range sortArray {
		field := fmt.Sprintf("%v %v", sort[0], sort[1])
		fields = append(fields, field)
	}
	return strings.Join(fields, ","), err
}

func ApplySearch(query string, fields []string) (string, []interface{}) {

	var field_array []string
	var values []interface{}

	for _, field := range fields {

		temp1 := fmt.Sprintf("%v ILIKE ?", field)

		temp2 := "%" + query + "%"

		field_array = append(field_array, temp1)
		values = append(values, temp2)
	}

	field_string := strings.Join(field_array, " OR ")

	return field_string, values
}

func ParseQueryToPaginate(query map[string]interface{}, page *Pagination.Paginatevalue) {
	JsonMarshaller(query, page)
	if query["per_page"] != nil {
		page.Per_page, _ = strconv.Atoi(query["per_page"].(string))
	}
	if query["page_no"] != nil {
		page.Page_no, _ = strconv.Atoi(query["page_no"].(string))
	}
}

func AddMandatoryFilters(p *Pagination.Paginatevalue, fieldName interface{}, operator interface{}, fieldValue interface{}) {
	if fieldName.(string) == "company_id" && fieldValue == uint(1) {
		return
	}
	mandatoryFilter := fmt.Sprintf("[[\"%v\",\"%v\",%v]]", fieldName, operator, fieldValue)
	if p.Filters == "" {
		p.Filters = mandatoryFilter
		return
	}
	p.Filters = p.Filters[:len(p.Filters)-1] + "," + mandatoryFilter[1:]
}
