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
	ipaas_models "fermion/backend_core/ipaas_core/model"
	"fermion/backend_core/pkg/util/helpers"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

//-------------error handled functions------------------------------

// func ReadFile(path string) (map[string]interface{}, error) {
// 	data := map[string]interface{}{}
// 	fileDataByte, ioErr := ioutil.ReadFile(path)
// 	if ioErr != nil {
// 		fmt.Println("--------> Read File error <----------------", ioErr)
// 		return data, ioErr
// 	}

//		if err := json.Unmarshal([]byte(fileDataByte), &data); err != nil {
//			fmt.Println("--------> Error while parsing Readfile <----------------" + path)
//			return nil, err
//		}
//		// fmt.Println("ReadFile --->>>>", data)
//		return data, nil
//	}
func GetValueFromSessionVariablesKey(key string, featureSessionVariables []model.KeyValuePair) (interface{}, error) {
	var result interface{}
	for _, sessionVariable := range featureSessionVariables {
		if sessionVariable.Key == key {
			result = sessionVariable.Value
			// fmt.Println("sessionVariable value --->>>>", result)
			return result, nil
		}
	}
	// fmt.Println("GetValueFromSessionVariablesKey --->>>>", result)
	errorString := fmt.Sprintf("%v - key not found", key)
	return result, errors.New(errorString)
}
func ParseObjectsFromConfigPayload(headerArr []interface{}) []model.KeyValuePair {
	var keyValuePairArrObj = make([]model.KeyValuePair, 0)
	for _, item := range headerArr {
		var keyValuePairObj model.KeyValuePair
		keyValuePairObj.Key = item.(map[string]interface{})["key"].(string)
		keyValuePairObj.Type = item.(map[string]interface{})["type"].(string)
		keyValuePairObj.Value = item.(map[string]interface{})["value"]
		propsInterface := item.(map[string]interface{})["props"]
		if propsInterface != nil {
			keyValuePairObj.Props = ParsePropsFromConfig(propsInterface.([]interface{}))
		}
		keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
	}
	// fmt.Println("ParseObjectsFromConfigPayload --->>>>", keyValuePairArrObj)
	return keyValuePairArrObj
}
func ParseKeyValuePair(keyValuePair []model.KeyValuePair, featureSessionVariables []model.KeyValuePair) []model.KeyValuePair {

	for ind, val := range keyValuePair {
		if val.Type == "variable" {
			keyValuePair[ind].Value, _ = GetValueFromSessionVariablesKey(val.Value.(string), featureSessionVariables)
			continue
		}
		if val.Type == "dynamic" {
			result, _ := CallFunctionFromProps(val.Props, featureSessionVariables)
			fmt.Println("--------> Function return Value from Props<----------------", result)
			keyValuePair[ind].Value = result
		}
	}
	// fmt.Println("ParseKeyValuePair --->>>>", keyValuePair)
	return keyValuePair
}
func FormatEndpoint(format string, sessionVariables []model.KeyValuePair) (string, error) {

	regex := regexp.MustCompile(`\{(.*?)\}`)
	keys := regex.FindAllStringSubmatch(format, -1)

	for _, key := range keys {
		result, err := GetValueFromSessionVariablesKey(key[1], sessionVariables)
		if err != nil {
			fmt.Println(err)
			fmt.Println("--------> error in format endpoint url <----------------")
			return format, err
		}
		value := fmt.Sprintf("%v", result)
		format = strings.Replace(format, key[0], value, -1)
	}
	// fmt.Println("FormatEndpoint --->>>>", format)
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
	// fmt.Println("GetUrlQueryParams --->>>>", format)
	return response, nil
}
func ConvertKeyValuePairToInterface(keyValuePair []model.KeyValuePair) interface{} {
	var object = make(map[string]interface{})

	for _, obj := range keyValuePair {
		object[obj.Key] = obj.Value
	}
	// fmt.Println("ConvertKeyValuePairToInterface --->>>>", object)
	return object
}
func MakeAPIRequest(requestMethod string, url string, headers []model.KeyValuePair, requestBody interface{}, featureSessionVariables []model.KeyValuePair, optionalParams ...interface{}) (model.EndpointResponse, error) {

	var req *http.Request
	var err error
	var responseCompletionNorms []interface{}
	RequestBodyUrlEncodedFlag := false

	//-------------------request api headers ---------------------------------------------------------------------
	// fmt.Println(headers)
	headers = ParseKeyValuePair(headers, featureSessionVariables)

	for index := range headers {
		if headers[index].Key == "Content-Type" {
			if headers[index].Value == "application/x-www-form-urlencoded" {
				RequestBodyUrlEncodedFlag = true
				break
			}
		}
	}

	//-------------------optional parameters - [oauth, ...] -----------------------------------------------------------------------
	if len(optionalParams) != 0 {
		endpointTaskFile, ok := optionalParams[0].(map[string]interface{})
		if !ok {
			fmt.Println("--------> error in MakeAPIRequest - optionalParams <----------------")
			return model.EndpointResponse{}, errors.New("error in oauth parameters type conversion")
		}

		if endpointTaskFile["response"] != nil {
			endpointTaskFileResponse, ok := endpointTaskFile["response"].(map[string]interface{})
			if !ok {
				fmt.Println("--------> error in MakeAPIRequest - endpointTaskFile[response] <----------------")
				return model.EndpointResponse{}, errors.New("endpointTaskFile[\"response\"]")
			}
			if endpointTaskFileResponse["response_completion_norms"] != nil {
				responseCompletionNorms, ok = endpointTaskFileResponse["response_completion_norms"].([]interface{})
				if !ok {
					fmt.Println("--------> error in MakeAPIRequest - endpointTaskFile[response][response_completion_norms] <----------------")
					return model.EndpointResponse{}, errors.New("endpointTaskFile[\"response\"][\"response_completion_norms\"]")
				}
			}
		}

		payload, ok := endpointTaskFile["payload"].(map[string]interface{})
		if !ok {
			fmt.Println("--------> error in MakeAPIRequest - endpointTaskFile[payload] <----------------")
			return model.EndpointResponse{}, errors.New("endpointTaskFile[\"payload\"]")
		}

		if payload["signature"] != nil {
			signatureArr, ok := payload["signature"].([]interface{})
			if !ok {
				fmt.Println("--------> error in MakeAPIRequest - signature optionalParams <----------------")
				return model.EndpointResponse{}, errors.New("error in signature optionalParams")
			}
			signatures := ParseObjectsFromConfigPayload(signatureArr)

			for _, signature := range signatures {

				requestData := map[string]interface{}{
					"method":  requestMethod,
					"url":     url,
					"payload": requestBody,
					"headers": headers,
				}

				for index, props := range signature.Props {
					signature.Props[index].Params = append(props.Params, model.KeyValuePair{
						Key:   "requestData",
						Value: requestData,
						Type:  "static",
					})
				}

				authSignature, err := CallFunctionFromProps(signature.Props, featureSessionVariables)

				if err != nil {
					fmt.Println("--------> error in Signature - optionalParams <----------------")
					return model.EndpointResponse{}, errors.New("error in signature functions")
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

	fmt.Println("\n url --------------->", url)
	// helpers.PrettyPrint("requestBody", requestBody)

	//------------------create new api request -------------------------------------------------------------------
	if RequestBodyUrlEncodedFlag {
		urlEncodedpayload, err := helpers.ReturnURLEncodeString(requestBody)
		if err != nil {
			fmt.Println("--------> warning in MakeAPIRequest - url encoded payload <----------------")
		}
		req, err = http.NewRequest(requestMethod, url, strings.NewReader(urlEncodedpayload))
		if err != nil {
			fmt.Println("--------> error in MakeAPIRequest - http.NewRequest <----------------")
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
					fmt.Println("--------> warning in MakeAPIRequest - marshal request body <----------------")
				}
				jsonBody = marshaldata
			} else {
				jsonBody = []byte("")
			}
		}

		// fmt.Println("\n\n body   -------------------------------------", string(jsonBody))

		req, err = http.NewRequest(requestMethod, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("--------> error in MakeAPIRequest - http.NewRequest <----------------")
			return model.EndpointResponse{}, err
		}
	}

	fmt.Println("--------> request successfully ready to call <----------------")

	for index := 0; index < len(headers); index++ {
		req.Header.Set(headers[index].Key, fmt.Sprintf("%v", headers[index].Value))
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("--------> Error while parsing MakeGETRequest <----------------")
		return model.EndpointResponse{}, err
	}
	response, err := ConvertResponseToInterface(res, responseCompletionNorms, featureSessionVariables)
	if err != nil {
		fmt.Println("--------> error in MakeAPIRequest - ConvertResponseToInterface <----------------")
		return response, err
	}

	// helpers.PrettyPrint("response", response)
	return response, err
}
func ConvertResponseToInterface(response *http.Response, responseCompletionNorms []interface{}, featureSessionVariables []model.KeyValuePair) (model.EndpointResponse, error) {
	var responseConverted model.EndpointResponse
	responseConverted.Status = response.Status
	responseConverted.StatusCode = response.StatusCode
	responseConverted.Header = response.Header
	responseConverted.Request = response.Request
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("--------> Error while parsing convertResponseToInterface <----------------")
		return responseConverted, err
	}

	if len(responseCompletionNorms) > 0 {
		responseCompletionNorms := ParseObjectsFromConfigPayload(responseCompletionNorms)
		for _, val := range responseCompletionNorms {
			if val.Key == "raw_response" {
				convertedResponseBody, err := CallFunctionFromProps(val.Props, featureSessionVariables, body)
				if err != nil {
					fmt.Println("--------> Error while parsing rawResponseCompletionNorms <----------------")
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

	// jsonData, _ := json.MarshalIndent(bodyInterface, "", "\t")
	// fmt.Println("\n\n-------------------------------------JSON RESPONSE ----------", string(jsonData))

	// fmt.Println("ConvertResponseToInterface --->>>>", responseConverted)
	return responseConverted, nil
}
func ConvertArrayInterfaceToArrayStructWithJsonKeyParse(responseSessionArr []interface{}, response model.EndpointResponse, featureSessionVariables []model.KeyValuePair) ([]model.KeyValuePair, error) {
	sessionVariables := ParseObjectsFromConfigPayload(responseSessionArr)
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
				fmt.Println("--------> Function return Value from Props<----------------", result)
				keyValuePairObj.Value = result
			}
			keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
			continue
		}

		//----------------type static -----------------------------------------------------------------------
		keyValuePairObj.Value = sessionVariables[index].Value
		keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
	}
	// fmt.Println("ConvertArrayInterfaceToArrayStructWithJsonKeyParse --->>>>", keyValuePairArrObj)
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
			// errorString := fmt.Sprintf("%v - not found in the json_path from object", field)
			fmt.Printf("%v - not found in the json_path from object", field)
			return nil, nil
		}
		if fieldIndex == len(segs)-1 {
			return object[field], nil
		}

		result, ok := object[field].(map[string]interface{})
		if !ok {
			errorString := fmt.Sprintf("%v - cannot proceed json_path further with this object", field)
			return nil, errors.New(errorString)
		}
		object = result
	}
	// fmt.Println("ParseJsonPathFromObject --->>>>", object)
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

	// fmt.Println("AppendOrUpdateKeyValuePair --->>>>", featureSessionVariables)
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
	// fmt.Println("AppendOrUpdateTaskResponseObject --->>>>", taskResponseArray)
	return taskResponseArray
}
func UpdateCurrentBatchRequirements(featureSessionVariables []model.KeyValuePair, currentBatchNumber int, perBatch int, batchRequestDetails []model.KeyValuePair) ([]model.KeyValuePair, error) {
	data, err := GetValueFromSessionVariablesKey("total_batch", featureSessionVariables)
	if err != nil {
		return featureSessionVariables, err
	}
	currentBatchDetails := ArraySlice(data.(([]interface{})), currentBatchNumber, currentBatchNumber+perBatch)
	featureSessionVariables = AppendOrUpdateKeyValuePair(featureSessionVariables, []model.KeyValuePair{
		{Key: "current_batch", Value: currentBatchDetails},
	})
	batchRequestDetails = ParseKeyValuePair(batchRequestDetails, featureSessionVariables)
	featureSessionVariables = AppendOrUpdateKeyValuePair(featureSessionVariables, batchRequestDetails)

	// fmt.Println("UpdateCurrentBatchRequirements --->>>>", featureSessionVariables)
	return featureSessionVariables, nil
}
func ArraySlice(data []interface{}, start int, end int) []interface{} {
	arrayLength := len(data)
	if start >= 0 && start <= arrayLength && end >= 0 && end <= arrayLength && end >= start {
		return data[start:end]
	}
	if start >= 0 && start < arrayLength {
		if end >= arrayLength {
			return data[start:]
		}
	}
	return []interface{}{}
}

// ----------------functions need to handle errors-----------------------------------------------------
func ParsePropsFromConfig(propsArr []interface{}) []model.Props {
	var returnPropsArr = make([]model.Props, 0)
	for _, item := range propsArr {
		var propsObj model.Props
		propsObj.Name = item.(map[string]interface{})["name"].(string)
		propsObj.Params = ParseObjectsFromConfigPayload(item.(map[string]interface{})["params"].([]interface{}))
		returnPropsArr = append(returnPropsArr, propsObj)
	}
	return returnPropsArr
}

func ConvertStringToInterface(input string) map[string]interface{} {
	var bodyInterface map[string]interface{}
	json.Unmarshal([]byte(input), &bodyInterface)
	return bodyInterface
}
func ConvertArrayToKeyValuePair(arrString []map[string]interface{}) []model.KeyValuePair {
	var keyValuePairArrObj []model.KeyValuePair
	for i := 0; i < len(arrString); i++ {
		var keyValuePairObj model.KeyValuePair
		keyValuePairObj.Key = arrString[i]["key"].(string)
		keyValuePairObj.Value = arrString[i]["value"].(string)
		keyValuePairObj.Type = arrString[i]["type"].(string)
		keyValuePairArrObj[i] = keyValuePairObj
	}
	return keyValuePairArrObj
}

func ConvertObjectToKeyValuePair(objString map[string]interface{}) []model.KeyValuePair {
	var keyValuePairArrObj = make([]model.KeyValuePair, 0)
	fmt.Println("--------> getKeysFromInterface starting <----------------")
	data := GetKeysFromInterface(objString)
	fmt.Println("--------> getKeysFromInterface ending <----------------")
	// if err := json.Unmarshal(arrString, &data); err != nil {
	// 	panic(err)
	// }
	for i := 0; i < len(data); i++ {
		var keyValuePairObj model.KeyValuePair
		keyValuePairObj.Key = data[i].(string)
		keyValuePairObj.Value = objString[data[i].(string)]
		keyValuePairArrObj = append(keyValuePairArrObj, keyValuePairObj)
	}
	return keyValuePairArrObj
}
func GetKeysFromInterface(m map[string]interface{}) []interface{} {
	keys := make([]interface{}, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func ContainString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func MakeKeyValuePair(key string, value interface{}, v_type string, props []model.Props) model.KeyValuePair {
	var data model.KeyValuePair
	data.Key = key
	data.Value = value
	data.Type = v_type
	data.Props = props
	return data
}
func GetBodyData(c echo.Context) interface{} {
	var body interface{}
	json.NewDecoder(c.Request().Body).Decode(&body)
	reqBytes, _ := json.Marshal(body)
	json.Unmarshal(reqBytes, &body)
	return body
}

func InterfaceToXml(data interface{}, key string) (string, error) {

	var xmlString string
	specialCases := strings.Split(key, "@@")
	leftKey := key
	rightKey := key
	if len(specialCases) > 1 {
		leftKey = fmt.Sprintf("%v %v", strings.TrimSpace(specialCases[0]), strings.TrimSpace(specialCases[1]))
		rightKey = strings.TrimSpace(specialCases[0])
	}

	switch data.(type) {
	case map[string]interface{}:
		obj := data.(map[string]interface{})
		sequences := GetSequenceFromMap(obj)
		for _, sequence := range sequences {
			seqString := sequence.(string)
			formatedString, _ := InterfaceToXml(obj[seqString], seqString)
			xmlString += formatedString
		}
		if key == "" {
			return xmlString, nil
		}
		return fmt.Sprintf("<%v>%v</%v>", leftKey, xmlString, rightKey), nil

	case []interface{}:
		arrayData := data.([]interface{})
		for _, value := range arrayData {
			formatedString, _ := InterfaceToXml(value, key)
			xmlString += formatedString
		}
		return xmlString, nil

	default:
		return fmt.Sprintf("<%v>%v</%v>", leftKey, data, rightKey), nil
	}
}

func GetSequenceFromMap(obj map[string]interface{}) []interface{} {

	sequences, ok := obj["sequence"].([]interface{})
	if ok {
		return sequences
	}

	keys := make([]interface{}, len(obj))
	index := 0
	for key := range obj {
		keys[index] = key
		index++
	}

	return keys
}
func ContainsValueInArray(arr []int64, v int64) bool {
	for _, value := range arr {
		if value == v {
			return true
		}
	}
	return false
}

func BearerAppend(token string) string {
	var properties []ipaas_models.Props
	property := ipaas_models.Props{
		Name: "TokenBearerType",
		Params: []ipaas_models.KeyValuePair{
			{Key: "token", Value: token}},
	}
	properties = append(properties, property)
	var kvp []ipaas_models.KeyValuePair
	bearer_token, _ := CallFunctionFromProps(properties, kvp)
	return bearer_token.(string)
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
			return nil, errors.New("invalid parameters")
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

	err = helpers.JsonMarshaller(dataBsonM, &data)
	if err != nil {
		return nil, err
	}

	if data["encrypted"] == nil {
		return nil, errors.New("failed to decrypt")
	}

	decryptedString, err := helpers.AESGCMDecrypt(data["encrypted"].(string))
	if err != nil {
		return nil, err
	}

	var decryptedData map[string]interface{}
	err = json.Unmarshal([]byte(decryptedString), &decryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil

	// TODO :- need to implement read from file
	// featureData, err := ReadFile(path)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return featureData, nil
	// }

}
