package utils

import (
	"encoding/json"
	"strconv"
	"strings"

	ipaas_models "fermion/backend_core/ipaas_core/model"
	// ipaas_utils "fermion/backend_core/ipaas_core/utils"
)

// Upload file to Colud provider assigned to user
//
//	input format:
//
// [
//
//	{
//		"data": "image base64 data"
//	},
//	{
//		"data": "image base64 data"
//	}
//
// ]
//
//	output format:
//
// [
//
//	{
//		"link": "url",
//		"file_name": "name"
//	},
//	{
//		"link": "url",
//		"file_name": "name"
//	}
//
// ]

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
type ServiceProvider string
type ModeOfOperation string

const (
	AWS          ServiceProvider = "AWS"
	AZURE        ServiceProvider = "AZURE"
	RawData      ServiceProvider = "RAWDATA"
	GCP          ServiceProvider = "GCP"
	DROP_BOX     ServiceProvider = "DROP_BOX"
	GOOGLE_DRIVE ServiceProvider = "GOOGLE_DRIVE"
)

const (
	FILE    ModeOfOperation = "File"
	BLOB    ModeOfOperation = "Blob"
	RAWDATA ModeOfOperation = "RawData"
	LINK    ModeOfOperation = "Link"
	BASE64  ModeOfOperation = "Base64"
	BYTES   ModeOfOperation = "Bytes"
)

func UploadFile(imageOptionsArray interface{}, fileName string, uniqueId string, userId uint, scheme string, host string, cloudProvider ServiceProvider, moduleName string, subModuleName string, extenstion string, mode ModeOfOperation, user_token string) ([]map[string]interface{}, error) {
	var imageLinkArray = make([]map[string]interface{}, 0)
	var headers []ipaas_models.KeyValuePair
	var req = make(map[string]interface{}, 0)
	bearer_token := user_token
	switch cloudProvider {
	case AWS:

		headers = append(headers, MakeKeyValuePair("Content-Type", "application/json", "static", nil))
		headers = append(headers, MakeKeyValuePair("Authorization", bearer_token, "static", nil))
		headers = append(headers, MakeKeyValuePair("Content-Length", "999999", "static", nil))
		headers = append(headers, MakeKeyValuePair("Accept", "*/*", "static", nil))

		var arr []map[string]interface{}
		dto, _ := json.Marshal(imageOptionsArray)
		json.Unmarshal(dto, &arr)

		for i, j := range arr {
			if j["data"] == nil || j["data"] == "" {
				imageLinkArray = append(imageLinkArray, j)
				continue
			}
			temp_fileName := strings.ReplaceAll(fileName, " ", "_")
			temp_fileName = temp_fileName + "_" + strconv.Itoa(i) + extenstion
			req["Bucket_name"] = "dev-api-files"
			req["File_raw_data"] = j["data"]
			req["Mode"] = "base64"
			req["file_name"] = temp_fileName
			req["User_data"] = map[string]interface{}{
				"module_name":     moduleName,
				"sub_module_name": subModuleName,
				"unique_id":       uniqueId,
				"user_id":         userId,
			}
			url := scheme + "://" + host + "/integrations/aws" + "/file_upload"
			// response, err := MakeAPIRequest("POST", url, headers, req, nil)
			response, _ := MakeAPIRequest("POST", url, headers, req, nil)

			// fmt.Println("Cloud Link After update:")
			// fmt.Println(response, err)
			imageLinkArray = append(imageLinkArray, map[string]interface{}{"link": response.Body["data"].(map[string]interface{})["url"], "file_name": temp_fileName})
		}
		return imageLinkArray, nil
	default:
		return imageLinkArray, nil

	}

}

func SendEmail(sender string, receiver string, subject string, body string, file string, scheme string, host string, cloudProvider ServiceProvider, user_token string) {

	var headers []ipaas_models.KeyValuePair
	var req = make(map[string]interface{}, 0)
	switch cloudProvider {
	case AWS:
		headers = append(headers, MakeKeyValuePair("Content-Type", "application/json", "static", nil))
		headers = append(headers, MakeKeyValuePair("Authorization", user_token, "static", nil))
		headers = append(headers, MakeKeyValuePair("Content-Length", "1000", "static", nil))
		headers = append(headers, MakeKeyValuePair("Accept", "*/*", "static", nil))
		req["sender"] = sender
		req["receiver"] = receiver
		req["subject"] = subject
		req["body"] = body
		if file == "" {
			req["file"] = nil
		} else {
			req["file"] = file
		}

		url := scheme + "://" + host + "/integrations/aws" + "/send_email"
		MakeAPIRequest("POST", url, headers, req, nil)

	default:

	}

}
