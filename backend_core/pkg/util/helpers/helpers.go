package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
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
func Contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func BasicAuth(username string, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetDateTime(isUtc bool, format string, delay float64) string {

	dateTime := time.Now()

	if isUtc {
		dateTime = dateTime.UTC()
	}

	dateTime = dateTime.Add(time.Hour*time.Duration(0) +
		time.Minute*time.Duration(delay) +
		time.Second*time.Duration(0))

	formatedDateTime := dateTime.Format(format)

	return formatedDateTime
}

func Mod(input1 int, input2 int) int {
	return input1 % input2
}

func ReturnURLEncodeString(input interface{}) (string, error) {
	data := url.Values{}
	var inputconversion map[string]interface{}

	jsondata, err := json.Marshal(input)
	if err != nil {
		return "", errors.New("invalid conversion(map[string]interface{})")
	}
	err = json.Unmarshal(jsondata, &inputconversion)
	if err != nil {
		return "", errors.New("invalid conversion(map[string]interface{})")
	}

	for i := range inputconversion {
		data.Set(i, fmt.Sprintf("%v", inputconversion[i]))
	}
	output := data.Encode()

	fmt.Println("ReturnURLEncodeString --->>>>", output)
	return output, nil
}

func PrettyPrint(key string, data interface{}) {
	byteData, err := json.MarshalIndent(data, "", "\t")
	fmt.Println("--------------------------" + key + "---------------------------------")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(byteData))
	fmt.Println("----------------------------------------------------------------------")
}

func JsonMarshaller(input interface{}, output interface{}) error {

	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, output)
	if err != nil {
		return err
	}
	return nil
}
