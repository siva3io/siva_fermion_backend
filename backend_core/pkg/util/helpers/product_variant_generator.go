package helpers

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
func VariantGeneration(array [][]map[string]interface{}) ([]interface{}, error) {

	/*
		{
			{{"id":1,"name":"S"},{"id":2,"name":"M"},{"id":3,"name":"L"}},
			{{"id":4,"name":"Red"},{"id":5,"name":"Blue"}},
			{{"id":6,"name":"Cotton"},{"id":7,"name":"Fabric"}},
			{{"id":8,"name":"Plain"},{"id":9,"name":"Checked"}}
		}
	*/

	length := len(array)
	var combinations []interface{}
	indices := make([]int, length)
	for i := 0; i < length; i++ {
		indices[i] = 0
	}
	for {
		var combination []interface{}
		for i := 0; i < length; i++ {
			j := indices[i]
			combination = append(combination, array[i][j])
		}
		combinations = append(combinations, map[string]interface{}{"combination": combination})
		next := length - 1

		for next >= 0 && (indices[next]+1) >= len(array[next]) {
			next = next - 1
		}
		if next < 0 {
			break
		}
		indices[next]++

		for i := next + 1; i < length; i++ {
			indices[i] = 0
		}
	}
	return combinations, nil
}
