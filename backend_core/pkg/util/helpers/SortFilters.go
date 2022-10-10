package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
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
func ApplyFilter(filters string) ([]string, []interface{}, error) {
	var filterArray [][]interface{}
	var fields []string
	var values []interface{}
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
		//this "in" filter query converts each individual filter value to string type
		if filter[1] == "IN" || filter[1] == "in" {
			valStr := fmt.Sprint(filter[2])
			temp_filter_value := strings.Split(valStr[1:len(valStr)-1], " ")
			fields = append(fields, filter[0].(string))
			values = append(values, temp_filter_value)
			continue
		}

		value := fmt.Sprintf("%v", filter[2])

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
