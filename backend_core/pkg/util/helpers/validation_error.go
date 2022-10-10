package helpers

import (
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
func ValidationErrorStructure(validation_err error) []map[string]interface{} {

	errs := strings.Split(validation_err.Error(), ";")

	fmt.Println(errs)

	var valErrMsg []map[string]interface{}

	for _, v := range errs {
		value := strings.Split(v, ":")
		errMsg := map[string]interface{}{
			"field":   value[0],
			"message": value[1],
		}
		valErrMsg = append([]map[string]interface{}{errMsg}, valErrMsg...)
	}

	return valErrMsg
}

func BindErrorStructure(bind_err error) map[string]interface{} {

	errs := strings.Split(bind_err.Error(), ",")

	fmt.Println(errs)

	bindErrMsg := map[string]interface{}{
		"field":   strings.Split(errs[3], "=")[1],
		"message": strings.Split(errs[1], "=")[1] + " type : " + strings.Split(errs[1], "=")[2],
	}

	return bindErrMsg
}

func JsonMarshalErrorStructure(marshal_err error) map[string]interface{} {

	errs := strings.Split(marshal_err.Error(), " ")

	fmt.Println(errs)

	marshalErrMsg := map[string]interface{}{
		"field":   errs[8],
		"message": marshal_err.Error(),
	}

	return marshalErrMsg
}
