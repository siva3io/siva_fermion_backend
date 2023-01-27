package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"fermion/backend_core/db"
	"fermion/backend_core/ipaas_core/model"
	pkg_helpers "fermion/backend_core/pkg/util/helpers"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetValueFromSessionVariablesKey(key string, featureSessionVariables []model.KeyValuePair) interface{} {
	var result interface{}
	for _, sessionVariable := range featureSessionVariables {
		if sessionVariable.Key == key {
			result = sessionVariable.Value
			return result
		}
	}
	//  fmt.Println("GetValueFromSessionVariablesKey ===>>>>", result)
	return result
}
func ParseObjectsFromConfigPayload(headerArr []interface{}) ([]model.KeyValuePair, error) {
	var keyValuePairArrObj []model.KeyValuePair
	err := pkg_helpers.JsonMarshaller(headerArr, &keyValuePairArrObj)
	if err != nil {
		return keyValuePairArrObj, err
	}
	//  fmt.Println("ParseObjectsFromConfigPayload ===>>>>", keyValuePairArrObj)
	return keyValuePairArrObj, nil
}
func ParseKeyValuePair(keyValuePair []model.KeyValuePair, featureSessionVariables []model.KeyValuePair) []model.KeyValuePair {
	for ind, val := range keyValuePair {
		if val.Type == "variable" {
			keyValuePair[ind].Value = GetValueFromSessionVariablesKey(val.Value.(string), featureSessionVariables)
			continue
		}
		if val.Type == "dynamic" {
			result, _ := CallFunctionFromProps(val.Props, featureSessionVariables)
			fmt.Println("=============> Function return Value from Props <============= ", result)
			keyValuePair[ind].Value = result
		}
	}
	//  fmt.Println("ParseKeyValuePair ===>>>>", keyValuePair)
	return keyValuePair
}
func FormatEndpoint(format string, sessionVariables []model.KeyValuePair) (string, error) {

	regex := regexp.MustCompile(`\{(.*?)\}`)
	keys := regex.FindAllStringSubmatch(format, -1)

	for _, key := range keys {
		result := GetValueFromSessionVariablesKey(key[1], sessionVariables)
		value := fmt.Sprintf("%v", result)
		format = strings.Replace(format, key[0], value, -1)
	}
	//  fmt.Println("FormatEndpoint ===>>>>", format)
	return format, nil
}
func GetUrlQueryParams(queryParams []model.KeyValuePair, featureSessionVariables []model.KeyValuePair) (string, error) {

	var response string
	for index := 0; index < len(queryParams); index++ {
		if queryParams[index].Value == nil {
			continue
		}
		response = response + queryParams[index].Key + "=" + url.QueryEscape(fmt.Sprintf("%v", queryParams[index].Value)) + "&"
	}
	if len(response) > 0 {
		response = "?" + response[0:len(response)-1]
	}
	//  fmt.Println("GetUrlQueryParams ===>>>>", response)
	return response, nil
}
func ConvertKeyValuePairToInterface(keyValuePair []model.KeyValuePair) interface{} {
	var object = make(map[string]interface{})

	for _, obj := range keyValuePair {
		object[obj.Key] = obj.Value
	}
	//  fmt.Println("ConvertKeyValuePairToInterface ===>>>>", object)
	return object
}
func MakeAPIRequest(requestMethod string, url string, headers []model.KeyValuePair, requestBody interface{}, featureSessionVariables []model.KeyValuePair, optionalParams ...interface{}) (model.EndpointResponse, error) {

	// fmt.Println("=============> Function call from Props <=============")
	var req *http.Request
	var responseCompletionNorms []model.KeyValuePair
	RequestBodyUrlEncodedFlag := false

	var err error

	//==================================request api headers=======================================================================================================================================================================
	headers = ParseKeyValuePair(headers, featureSessionVariables)

	//============================check content-type is application/x-www-form-urlencoded or not==================================================================================================================================
	for index := range headers {
		if headers[index].Key == "Content-Type" {
			if headers[index].Value == "application/x-www-form-urlencoded" {
				RequestBodyUrlEncodedFlag = true
				break
			}
		}
	}

	//============================optional parameters - [signature ...]===========================================================================================================================================================
	if len(optionalParams) != 0 {
		endpointTaskFile, ok := optionalParams[0].(model.APIFile)
		if !ok {
			fmt.Println("=============> error in MakeAPIRequest - optionalParams <=============")
			return model.EndpointResponse{}, errors.New("oops! issue in endpoint_task_file conversion from MakeAPIRequest")
		}
		responseCompletionNorms = endpointTaskFile.Response.ResponseCompletionNorms

		payload := endpointTaskFile.Payload

		if len(payload.Signature) > 0 {
			for _, signature := range payload.Signature {

				requestData := map[string]interface{}{
					"method":  requestMethod,
					"url":     url,
					"payload": requestBody,
					"headers": headers,
				}
				newSignatureProps := make([]model.Props, 0)
				newSignatureProps = append(newSignatureProps, signature.Props...)

				for index, props := range signature.Props {
					newSignatureProps[index].Params = append(props.Params, model.KeyValuePair{
						Key:   "requestData",
						Value: requestData,
						Type:  "static",
					})
				}

				authSignature, err := CallFunctionFromProps(newSignatureProps, featureSessionVariables)
				fmt.Println("authSignature", authSignature)

				if err != nil {
					fmt.Println("=============> error in Signature - optionalParams <=============")
					fmt.Println("--------> error in Signature - optionalParams <----------------")
					return model.EndpointResponse{}, errors.New("oops! issue in signature functions from MakeAPIRequest")
				}
				if signature.Value.(string) == "headers" {
					headers = append(headers, model.KeyValuePair{
						Key:   signature.Key,
						Value: authSignature.(string),
						Type:  "",
						Props: []model.Props{},
					})
				}
				if signature.Value.(string) == "query_params" {
					url += "&" + signature.Key + "=" + authSignature.(string)
				}
			}
		}
	}

	//------------------create new api request -------------------------------------------------------------------
	if RequestBodyUrlEncodedFlag {
		urlEncodedpayload, err := pkg_helpers.ReturnURLEncodeString(requestBody)
		if err != nil {
			fmt.Println("=============> warning in MakeAPIRequest - url encoded payload <=============")
		}
		req, err = http.NewRequest(requestMethod, url, strings.NewReader(urlEncodedpayload))
		if err != nil {
			fmt.Println("=============> error in MakeAPIRequest - http.NewRequest <=============")
			return model.EndpointResponse{}, err
		}
	}
	if !RequestBodyUrlEncodedFlag {
		//----------------request body--------------------------------------------------------------------------------
		jsonBody, ok := requestBody.([]byte)
		if !ok {
			if requestBody != nil {
				marshaldata, err := json.Marshal(requestBody)
				if err != nil {
					fmt.Println("=============> warning in MakeAPIRequest - marshal request body <=============")
				}
				jsonBody = marshaldata
			} else {
				jsonBody = []byte("")
			}
		}

		// fmt.Println("\n\n body   -------------------------------------", string(jsonBody))

		req, err = http.NewRequest(requestMethod, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("=============> error in MakeAPIRequest - http.NewRequest <=============")
			return model.EndpointResponse{}, err
		}
	}

	for index := 0; index < len(headers); index++ {
		req.Header.Set(headers[index].Key, fmt.Sprintf("%v", headers[index].Value))
	}

	fmt.Println("=============> api_endpoint <============= ", url)
	fmt.Println("=============> api_request_type <============= ", requestMethod)
	// pkg_helpers.PrettyPrint("api_request_body",requestBody)
	// pkg_helpers.PrettyPrint("api_request_headers",headers)

	fmt.Println("=============> request successfully ready to call <=============")
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("=============> Error while parsing MakeGETRequest <=============")
		return model.EndpointResponse{}, err
	}
	response, err := ConvertResponseToInterface(res, responseCompletionNorms, featureSessionVariables)
	if err != nil {
		fmt.Println("=============> error in MakeAPIRequest - ConvertResponseToInterface <=============")
		return response, err
	}

	fmt.Println("=============> api_response_status <============= ", response.Status)
	// pkg_helpers.PrettyPrint("api_response",response.Body)
	return response, err
}
func ConvertResponseToInterface(response *http.Response, responseCompletionNorms []model.KeyValuePair, featureSessionVariables []model.KeyValuePair) (model.EndpointResponse, error) {
	var responseConverted model.EndpointResponse
	responseConverted.Status = response.Status
	responseConverted.StatusCode = response.StatusCode
	responseConverted.Header = response.Header
	responseConverted.Request = response.Request
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("=============> Error while parsing convertResponseToInterface <============= ")
		return responseConverted, err
	}

	if len(responseCompletionNorms) > 0 {
		for _, val := range responseCompletionNorms {
			if val.Key == "raw_response" {
				newProps := make([]model.Props, 0)
				newProps = append(newProps, val.Props...)
				for index, props := range val.Props {
					newProps[index].Params = append(props.Params, model.KeyValuePair{
						Key:   "body",
						Value: body,
						Type:  "static",
					})
				}
				convertedResponseBody, err := CallFunctionFromProps(newProps, featureSessionVariables)
				if err != nil {
					fmt.Println("=============> Error while parsing rawResponseCompletionNorms <============= ")
					return responseConverted, err
				}
				responseConverted.Body = convertedResponseBody.(map[string]interface{})
				return responseConverted, nil
			}
		}
	}

	bodyInterface := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyInterface)
	if err != nil {
		var bodyArrayInterface []map[string]interface{}
		json.Unmarshal(body, &bodyArrayInterface)
		bodyInterface["data"] = bodyArrayInterface
	}
	responseConverted.Body = bodyInterface

	//  fmt.Println("ConvertResponseToInterface ===>>>>", responseConverted)
	return responseConverted, nil
}
func ConvertArrayInterfaceToArrayStructWithJsonKeyParse(responseSessionArray interface{}, response model.EndpointResponse, featureSessionVariables []model.KeyValuePair) ([]model.KeyValuePair, error) {
	sessionVariables, ok := responseSessionArray.([]model.KeyValuePair)
	if !ok {
		return nil, errors.New("cannot convert response_sessiona_rray to []keyValuePair")
	}
	var keyValuePairArrObj = make([]model.KeyValuePair, 0)

	for index := 0; index < len(sessionVariables); index++ {
		var keyValuePairObj model.KeyValuePair
		keyValuePairObj.Key = sessionVariables[index].Key
		keyValuePairObj.Type = sessionVariables[index].Type

		//----------------type variable -----------------------------------------------------------------------
		if sessionVariables[index].Type == "variable" {
			result, err := ParseJsonPathFromObject(response.Body, sessionVariables[index].Value.(string))
			if err != nil {
				return keyValuePairArrObj, err
			}
			keyValuePairObj.Value = result
			keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
			continue

			//----------------type dynamic -----------------------------------------------------------------------
		} else if sessionVariables[index].Type == "dynamic" {
			sessionValue, err := ParseJsonPathFromObject(response.Body, sessionVariables[index].Value.(string))
			if err != nil {
				return keyValuePairArrObj, err
			}
			if sessionValue != nil {
				featureSessionVariablesWithResponse := append(featureSessionVariables, model.KeyValuePair{
					Key:   "response",
					Value: sessionValue,
					Type:  "static",
				})
				result, _ := CallFunctionFromProps(sessionVariables[index].Props, featureSessionVariablesWithResponse)
				fmt.Println("=============> Function return Value from Props <============= ", result)
				keyValuePairObj.Value = result
			}
			keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
			continue
		}

		//----------------type static -----------------------------------------------------------------------
		keyValuePairObj.Value = sessionVariables[index].Value
		keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
	}
	//  fmt.Println("ConvertArrayInterfaceToArrayStructWithJsonKeyParse ===>>>>", keyValuePairArrObj)
	return keyValuePairArrObj, nil
}
func ParseJsonPathFromObject(data map[string]interface{}, keypath string) (interface{}, error) {

	if keypath == "" {
		return data, nil
	}
	var segs []string = strings.Split(keypath, ".")
	object := data

	for fieldIndex, field := range segs {
		if object[field] == nil {
			fmt.Printf("oops! %v - not found in the json_path from object\n", field)
			return nil, nil
		}
		if fieldIndex == len(segs)-1 {
			return object[field], nil
		}
		result, ok := object[field].(map[string]interface{})
		if !ok {
			errorString := fmt.Sprintf("oops! %v - cannot proceed json_path further with this object\n", field)
			return nil, errors.New(errorString)
		}
		object = result
	}
	//  fmt.Println("ParseJsonPathFromObject ===>>>>", object)
	return object, nil
}
func AppendOrUpdateKeyValuePair(featureSessionVariables []model.KeyValuePair, sessionVariables []model.KeyValuePair) []model.KeyValuePair {
	FlagForUpdated := false
	for i := 0; i < len(sessionVariables); i++ {
		for j := 0; j < len(featureSessionVariables); j++ {
			if sessionVariables[i].Key == featureSessionVariables[j].Key {
				featureSessionVariables[j].Value = sessionVariables[i].Value
				featureSessionVariables[j].Props = sessionVariables[i].Props
				featureSessionVariables[j].Type = sessionVariables[i].Type
				FlagForUpdated = true
				break
			}
		}
		if !FlagForUpdated {
			featureSessionVariables = append(featureSessionVariables, sessionVariables[i])
		}
		FlagForUpdated = false
	}

	//  fmt.Println("AppendOrUpdateKeyValuePair ===>>>>", featureSessionVariables)
	return featureSessionVariables
}
func AppendOrUpdateTaskResponseObject(taskResponseArray []model.TaskEndpointResponse, currentTaskResponseObject model.TaskEndpointResponse) []model.TaskEndpointResponse {

	for index, taskResponseObject := range taskResponseArray {
		if currentTaskResponseObject.Name == taskResponseObject.Name {
			taskResponseArray[index].Header = currentTaskResponseObject.Header
			taskResponseArray[index].EndpointResponse = currentTaskResponseObject.EndpointResponse
			taskResponseArray[index].MappedResponse = currentTaskResponseObject.MappedResponse
			taskResponseArray[index].Completed = currentTaskResponseObject.Completed
			taskResponseArray[index].SessionVariables = currentTaskResponseObject.SessionVariables

			return taskResponseArray
		}
	}
	taskResponseArray = append(taskResponseArray, currentTaskResponseObject)

	//  fmt.Println("AppendOrUpdateTaskResponseObject ===>>>>", taskResponseArray)
	return taskResponseArray
}
func UpdateBatchRequirements(featureSessionVariables []model.KeyValuePair, currentBatchNumber int, perBatch int, batchRequestDetails []model.KeyValuePair) ([]model.KeyValuePair, error) {
	data := GetValueFromSessionVariablesKey("total_batch", featureSessionVariables)
	currentBatchDetails := ArraySlice(data.(([]interface{})), currentBatchNumber, currentBatchNumber+perBatch)
	featureSessionVariables = AppendOrUpdateKeyValuePair(featureSessionVariables, []model.KeyValuePair{
		{Key: "current_batch", Value: currentBatchDetails},
	})
	batchRequestDetails = ParseKeyValuePair(batchRequestDetails, featureSessionVariables)
	featureSessionVariables = AppendOrUpdateKeyValuePair(featureSessionVariables, batchRequestDetails)

	//  fmt.Println("UpdateBatchRequirements ===>>>>", featureSessionVariables)
	return featureSessionVariables, nil
}
func ArraySlice(data []interface{}, start int, end int) []interface{} {
	arrayLength := len(data)
	if start >= 0 && start <= arrayLength && end >= 0 && end <= arrayLength && end >= start {

		//  fmt.Println("ArraySlice ===>>>>", data[start:end])
		return data[start:end]
	}
	if start >= 0 && start < arrayLength {
		if end >= arrayLength {

			//  fmt.Println("ArraySlice ===>>>>", data[start:])
			return data[start:]
		}
	}
	return []interface{}{}
}
func ParsePropsFromConfig(propsArr []interface{}) ([]model.Props, error) {
	var returnPropsArr []model.Props
	err := pkg_helpers.JsonMarshaller(propsArr, &returnPropsArr)
	if err != nil {
		return returnPropsArr, err
	}

	//  fmt.Println("ArraySlice ===>>>>", data[start:])
	return returnPropsArr, nil
}
func ConvertStringToInterface(input string) (map[string]interface{}, error) {
	var bodyInterface map[string]interface{}
	err := pkg_helpers.JsonMarshaller(input, &bodyInterface)
	if err != nil {
		return bodyInterface, err
	}

	//  fmt.Println("ConvertStringToInterface ===>>>>", bodyInterface)
	return bodyInterface, nil
}
func ConvertArrayToKeyValuePair(arrString []map[string]interface{}) ([]model.KeyValuePair, error) {
	var keyValuePairArrObj []model.KeyValuePair
	err := pkg_helpers.JsonMarshaller(arrString, &keyValuePairArrObj)
	if err != nil {
		return keyValuePairArrObj, err
	}

	//  fmt.Println("ConvertArrayToKeyValuePair ===>>>>", keyValuePairArrObj)
	return keyValuePairArrObj, nil
}
func ConvertObjectToKeyValuePair(objString map[string]interface{}) []model.KeyValuePair {
	var keyValuePairArrObj = make([]model.KeyValuePair, 0)

	for key := range objString {
		var keyValuePairObj model.KeyValuePair
		keyValuePairObj.Key = key
		keyValuePairObj.Value = objString[key]
		keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
	}

	//  fmt.Println("ConvertObjectToKeyValuePair ===>>>>", keyValuePairArrObj)
	return keyValuePairArrObj
}
func MakeKeyValuePair(key string, value interface{}, v_type string, props []model.Props) model.KeyValuePair {
	var data model.KeyValuePair
	data.Key = key
	data.Value = value
	data.Type = v_type
	data.Props = props

	//  fmt.Println("MakeKeyValuePair ===>>>>", data)
	return data
}
func GetBodyData(c echo.Context) interface{} {
	var body interface{}
	json.NewDecoder(c.Request().Body).Decode(&body)

	//  fmt.Println("GetBodyData ===>>>>", body)
	return body
}
func GetSequenceFromMap(obj map[string]interface{}) []interface{} {

	sequences, ok := obj["sequence"].([]interface{})
	if ok {
		//  fmt.Println("GetSequenceFromMap ===>>>>", sequences)
		return sequences
	}

	keys := make([]interface{}, len(obj))
	index := 0
	for key := range obj {
		keys[index] = key
		index++
	}

	//  fmt.Println("GetSequenceFromMap ===>>>>", keys)
	return keys
}
func ReadFeature(collectionName, id string, optionalParams ...interface{}) (map[string]interface{}, error) {
	var Collections = map[string]string{
		"features": "ipaas_features",
		"mappers":  "ipaas_mappers",
		"apis":     "ipaas_apis",
	}
	collectionName = Collections[collectionName]
	var query primitive.M
	if len(optionalParams) > 0 {
		name, ok := optionalParams[0].(string)
		if !ok {
			return nil, errors.New("oops! invalid parameters for ReadFeature")
		}
		query = bson.M{"name": name}
	} else {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		query = bson.M{"_id": objectId}
	}
	var dataBsonM bson.M
	var data map[string]interface{}
	db := db.NoSqlDbManager()
	collection := db.Collection(collectionName)
	err := collection.FindOne(context.Background(), query).Decode(&dataBsonM)
	if err != nil {
		return nil, err
	}
	err = pkg_helpers.JsonMarshaller(dataBsonM, &data)
	if err != nil {
		return nil, err
	}
	if data["encrypted"] == nil {
		return nil, errors.New("failed to decrypt")
	}
	decryptedString, err := pkg_helpers.AESGCMDecrypt(data["encrypted"].(string))
	if err != nil {
		return nil, err
	}
	var decryptedData map[string]interface{}
	err = json.Unmarshal([]byte(decryptedString), &decryptedData)
	if err != nil {
		return nil, err
	}

	//  fmt.Println("ReadFeature ===>>>>", decryptedData)
	return decryptedData, nil
}
