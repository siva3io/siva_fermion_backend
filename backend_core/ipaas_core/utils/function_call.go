package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fermion/backend_core/ipaas_core/model"
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
type Functions struct{}

// -----------------------------function call-------------------------------------------------------------------------------
func CallFunctionFromProps(propsArr []model.Props, parameters ...interface{}) (interface{}, error) {
	var valueToReturn []reflect.Value
	var err error
	var ReturnValue interface{}

	featureSessionVariables, ok := parameters[0].([]model.KeyValuePair)
	if !ok {
		return ReturnValue, errors.New("first parameter must be feature session variables")
	}

	fmt.Println("--------> Function call from Props <----------------")

	for index := 0; index < len(propsArr); index++ {
		var paramsAfterEvaluation = make([]interface{}, 0)

		for _, item := range propsArr[index].Params {
			sessionValue := item.Value
			if item.Type == "variable" {
				sessionValue, _ = GetValueFromSessionVariablesKey(item.Value.(string), featureSessionVariables)
				// if err != nil {
				// 	return ReturnValue, err
				// }
			}
			paramsAfterEvaluation = append(paramsAfterEvaluation, sessionValue)
		}
		parameters[0] = featureSessionVariables
		paramsAfterEvaluation = append(paramsAfterEvaluation, parameters...)

		//-----------------------call function by name --------------------------------------------------
		valueToReturn, err = CallFuncByName(&Functions{}, propsArr[index].Name, paramsAfterEvaluation...)
		if err != nil {
			return ReturnValue, err
		}

		// fmt.Println("---->>>>", valueToReturn)
		if len(valueToReturn) > 0 {
			functionReturnValue := valueToReturn[0]

			//-----------------temporary error handle-----------------------------
			if len(valueToReturn) > 1 {
				FunctionError := valueToReturn[1]
				errorValue, ok := FunctionError.Interface().(error)
				if ok {
					return ReturnValue, errorValue
				}
			}

			//-------------------data type conversion ------------------------------------------------
			ReturnValue, err = ReturnTheVariableType(functionReturnValue)
			if err != nil {
				return ReturnValue, err
			}

			featureSessionVariables = AppendOrUpdateKeyValuePair(featureSessionVariables, []model.KeyValuePair{
				{Key: propsArr[index].Name, Value: ReturnValue},
			})
		}
	}
	return ReturnValue, nil
}
func CallFuncByName(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {

	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)

	if !m.IsValid() {
		errorString := fmt.Sprintf("method not found \"%s\"", funcName)
		return make([]reflect.Value, 0), errors.New(errorString)
	}
	in := make([]reflect.Value, len(params))
	fmt.Printf("--------> function name<----------------%v(", funcName)
	for index, param := range params {
		if index != 0 {
			fmt.Printf(",")
		}
		fmt.Printf("%v", reflect.TypeOf(param))
		if param == nil {
			in[index] = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
			continue
		}
		in[index] = reflect.ValueOf(param)
	}
	fmt.Printf(")\n")
	out = m.Call(in)
	return
}
func ReturnTheVariableType(value reflect.Value) (interface{}, error) {
	var dataType interface{}
	var ok bool

	dataType, ok = value.Interface().(string)
	if ok {
		return dataType, nil
	}
	dataType, ok = value.Interface().(bool)
	if ok {
		return dataType, nil
	}
	dataType, ok = value.Interface().(int)
	if ok {
		return dataType, nil
	}
	dataType, ok = value.Interface().(float64)
	if ok {
		return dataType, nil
	}
	dataType, ok = value.Interface().(float32)
	if ok {
		return dataType, nil
	}
	dataType, ok = value.Interface().(map[string]interface{})
	if ok {
		return dataType, nil
	}
	dataType, ok = value.Interface().([]map[string]interface{})
	if ok {
		return dataType, nil
	}

	return value.Interface(), nil
}

// -----------------------------ipaas props function--------------------------------------------------------------------------------
func (f *Functions) TokenBearerType(token interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {
	tokenString, ok := token.(string)

	if !ok {
		return "", errors.New("token must be a string")
	}
	if tokenString == "" {
		return "", errors.New("token must not be empty")
	}
	output := "Bearer " + tokenString

	//  fmt.Println("TokenBearerType --->>>>", output)
	return output, nil
}
func (f *Functions) Paginate(initialValue, pageValue, stepValue interface{}, featureSessionVariables []model.KeyValuePair) (interface{}, error) {

	if pageValue == nil || pageValue == "" {
		return initialValue, nil
	}
	currentPage, _ := strconv.Atoi(pageValue.(string))
	stepValueInt, _ := strconv.Atoi(stepValue.(string))

	return strconv.Itoa(currentPage + stepValueInt), nil
}
func (f *Functions) ConditionNorms(key string, value interface{}, optionalParams ...interface{}) (bool, error) {

	lenOptParams := len(optionalParams)
	var prevResponse bool
	operator := "OR"
	response := optionalParams[lenOptParams-2].(model.EndpointResponse)

	if lenOptParams > 2 {
		prevResponse = optionalParams[0].(bool)
	}

	if lenOptParams > 3 {
		operator = optionalParams[1].(string)
	}

	var responseValue interface{}
	if response.Body == nil {
		errorString := fmt.Sprintf("%v is not available in response", key)
		return false, errors.New(errorString)
	}
	responseValue, err := ParseJsonPathFromObject(response.Body, key)
	if err != nil {
		return false, err
	}
	fmt.Printf("OPERATOR --------->%v , Prev Resp -----------> %v , conditionValue --------> %v ,actualValue --------> %v", operator, prevResponse, fmt.Sprintf("%v", value), fmt.Sprintf("%v", responseValue))
	if responseValue == nil {
		return true, nil
	}
	if operator == "OR" {
		return prevResponse || fmt.Sprintf("%v", responseValue) == fmt.Sprintf("%v", value), nil
	}
	if operator == "AND" {
		return prevResponse && fmt.Sprintf("%v", responseValue) == fmt.Sprintf("%v", value), nil
	}
	return true, nil
}
func (f *Functions) GetLength(array []interface{}, featureSessionVariables []model.KeyValuePair) int {
	return len(array)
}
func (f *Functions) ParseJsonPathFromObject(data map[string]interface{}, keypath string, featureSessionVariables []model.KeyValuePair) (interface{}, error) {

	if data == nil {
		return nil, nil
	}
	res, err := ParseJsonPathFromObject(data, keypath)
	if err != nil {
		return res, err
	}
	return res, nil
}
func (f *Functions) ReturnIndexObjectFromArray(array interface{}, index string, featureSessionVariables []model.KeyValuePair) (interface{}, error) {
	var output map[string]interface{}
	var inputArray []interface{}

	madata, err := json.Marshal(array)
	if err != nil {
		return output, err
	}
	err = json.Unmarshal(madata, &inputArray)
	if err != nil {
		return output, err
	}

	i, err := strconv.Atoi(index)
	if err != nil {
		return output, err
	}

	if i < 0 || i >= len(inputArray) {
		return output, nil
	}
	return inputArray[i], nil
}
func (f *Functions) CheckStatus(optionalParams ...interface{}) bool {

	response := optionalParams[0].(model.EndpointResponse)

	return response.StatusCode < 200 || response.StatusCode >= 300
}
func (f *Functions) FloattoString(input interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {
	floatValue, ok := input.(float64)
	if !ok {
		return "", errors.New("invalid input type")
	}
	return fmt.Sprintf("%v", floatValue), nil
}
func (f *Functions) StringtoFloat(input interface{}, featureSessionVariables []model.KeyValuePair) (float32, error) {
	if input == nil {
		return 0, nil
	}
	stringValue, ok := input.(string)
	if !ok {
		return 0, errors.New("invalid input type")
	}
	if stringValue == "" {
		return 0, nil
	}
	output, err := strconv.ParseFloat(stringValue, 32)
	fmt.Println("-=========", output)
	if err != nil {
		panic(err)
		// return 0, err
	}
	return float32(output), nil
}
func (f *Functions) FloattoInt(input interface{}, featureSessionVariables []model.KeyValuePair) (int, error) {
	if input == nil {
		return 0, errors.New("invalid input type")
	}
	floatValue, ok := input.(float64)
	if !ok {
		return 0, errors.New("invalid conversion input type")
	}
	return int(floatValue), nil
}
func (f *Functions) AppendToArray(input interface{}, featureSessionVariables []model.KeyValuePair) ([]interface{}, error) {
	array := []interface{}{}
	if input == nil {
		return array, errors.New("input must be non-nil")
	}
	array = append(array, input)
	return array, nil
}
func (f *Functions) GetValueFromSessionVariables(value interface{}, featureSessionVariables []model.KeyValuePair) interface{} {
	//-------------get the value from the session variables and return it simple------------------------
	return value
}
func (f *Functions) ParseKeyValuePairThenReturnJSONString(input interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {
	array, ok := input.([]interface{})
	if !ok {
		return "", errors.New("cannot convert to []interface{})")
	}

	keyValuePairs := ParseObjectsFromConfigPayload(array)
	keyValuePairs = ParseKeyValuePair(keyValuePairs, featureSessionVariables)
	output := ConvertKeyValuePairToInterface(keyValuePairs)

	jsonByteArray, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	return string(jsonByteArray), nil
}

// ---------------------------mapper api call prop functions-----------------------------------------------------------------
func (f *Functions) ExecuteGetApiCall(taskObject map[string]interface{}, featureSessionVariables []model.KeyValuePair) (interface{}, error) {

	taskResponseArray := []model.TaskEndpointResponse{}
	taskResponseObject := model.TaskEndpointResponse{}
	singleTaskObject := map[string]interface{}{}

	//---------------------base path----------------------------------------------
	// path, err := GetValueFromSessionVariablesKey("basePath", featureSessionVariables)
	// if err != nil {
	// 	return nil, err
	// }
	// basePath := path.(string)

	//---------------------read file----------------------------------------------
	endpointTaskFile, err := ReadFeature("apis", "", "eunimart.mapper_api_call")
	if err != nil {
		return nil, err
	}

	//---------------------endpoint------------------------------------------------
	endpointTaskFile["endpoint"] = taskObject["endpoint"]

	//-----------------query params------------------------------------------------
	query_params := []model.KeyValuePair{}
	if taskObject["query_params"] != nil {
		jsonQueryParamsArray, ok := taskObject["query_params"].([]interface{})
		if !ok {
			return nil, errors.New("invalid query params conversion([]interface{})")
		}

		if len(jsonQueryParamsArray) > 0 {
			for _, josnQueryParam := range jsonQueryParamsArray {
				param, ok := josnQueryParam.(map[string]interface{})
				if !ok {
					return nil, errors.New("invalid query params conversion(map[string]interface{})")
				}
				if param["key"] == "filters" {
					res, err := f.FormatFilter(param["value"], featureSessionVariables)
					if err != nil {
						return nil, err
					}
					query_params = append(query_params, model.KeyValuePair{Key: param["key"].(string), Value: res, Type: "static"})
					continue
				}
				if param["key"] == "sort" {
					res, err := f.FormatSort(param["value"], featureSessionVariables)
					if err != nil {
						return nil, err
					}
					query_params = append(query_params, model.KeyValuePair{Key: param["key"].(string), Value: res, Type: "static"})
					continue
				}
				query_params = append(query_params, model.KeyValuePair{Key: param["key"].(string), Value: param["value"], Type: param["type"].(string)})
			}
		}
	}

	marshaldata, err := json.Marshal(query_params)
	if err != nil {
		return nil, err
	}
	queryParamsInput := []interface{}{}
	err = json.Unmarshal(marshaldata, &queryParamsInput)
	if err != nil {
		return nil, err
	}
	endpointTaskFile["payload"].(map[string]interface{})["query_params"] = queryParamsInput

	taskResponseObject = f.ExecuteGetRequest(taskResponseArray, taskResponseObject, singleTaskObject, endpointTaskFile, featureSessionVariables)

	return taskResponseObject.Body, nil
}
func (f *Functions) FormatFilter(input interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {

	type Filter struct {
		Column   string      `json:"column"`
		Operator string      `json:"operator"`
		Value    interface{} `json:"value"`
		Type     string      `json:"type"`
	}

	var filters []Filter

	if input == nil {
		return "", errors.New("Filters function input must not be nil")
	}
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, &filters)
	if err != nil {
		return "", err
	}

	output := "["
	for index, filter := range filters {
		if index != 0 {
			output = output + ","
		}
		loop_output := ""
		loop_output = loop_output + "["
		loop_output = loop_output + "\"" + filter.Column + "\","
		loop_output = loop_output + "\"" + filter.Operator + "\","
		if filter.Type == "variable" {
			interfaceValue, err := GetValueFromSessionVariablesKey(filter.Value.(string), featureSessionVariables)
			if err != nil {
				return "", err
			}
			stringValue := fmt.Sprintf("%v", interfaceValue)
			loop_output = loop_output + "\"" + stringValue + "\"]"
			output = output + loop_output
			continue
		}
		stringValue := fmt.Sprintf("%v", filter.Value)
		loop_output = loop_output + "\"" + stringValue + "\"]"
		output = output + loop_output
	}
	output = output + "]"
	return output, nil
}
func (f *Functions) FormatSort(input interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {
	type Sort struct {
		Column  string `json:"column"`
		OrderBy string `json:"order_by"`
	}

	var sorts []Sort

	if input == nil {
		return "", errors.New("sorting function input must not be nil")
	}
	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, &sorts)
	if err != nil {
		return "", err
	}

	output := "["
	for index, sort := range sorts {
		if index != 0 {
			output = output + ","
		}
		loop_output := ""
		loop_output = loop_output + "["
		loop_output = loop_output + "\"" + sort.Column + "\","
		loop_output = loop_output + "\"" + sort.OrderBy + "\"]"
		output = output + loop_output
	}
	output = output + "]"

	return output, nil
}

// TODO : redanunt function same as "single request"
func (f *Functions) ExecuteGetRequest(taskResponseArray []model.TaskEndpointResponse, taskResponseObject model.TaskEndpointResponse, singleTaskObject map[string]interface{}, endpointTaskFile map[string]interface{}, featureSessionVariables []model.KeyValuePair) model.TaskEndpointResponse {
	var requestBody interface{}
	var queryParams string
	var urlWithParams string
	headers := []model.KeyValuePair{}
	var requestMethod string
	Errors := map[string]interface{}{}

	//---------------------payload details-[ header-[], body-[],query_params-[], oauth-[] ]--------------------------------------
	payload, ok := endpointTaskFile["payload"].(map[string]interface{})
	if !ok {
		fmt.Println("--------> error in single_request payload <----------------")
		Errors["message"] = "error in single_request payload"
		taskResponseObject.Errors = Errors
		taskResponseObject.Completed = false
		return taskResponseObject
	}

	//-----------------------------end point url---------------------------------------------------------------------------------
	getURL := endpointTaskFile["scheme"].(string) + "://" + endpointTaskFile["endpoint"].(string)
	getURL, err := FormatEndpoint(getURL, featureSessionVariables)
	if err != nil {
		Errors["message"] = "error from format endpoint function"
		Errors["format_endpoint_error"] = fmt.Sprintf("%v", err)
		taskResponseObject.Errors = Errors
		taskResponseObject.Completed = false
		return taskResponseObject
	}

	//-----------------------------query params----------------------------------------------------------------------------------
	if payload["query_params"] != nil {
		jsonQueryParams, ok := payload["query_params"].([]interface{})
		if !ok {
			fmt.Println("--------> error in single_request query_params <----------------")
			Errors["message"] = "error in single_request query_params"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		if len(jsonQueryParams) > 0 {
			params := ParseObjectsFromConfigPayload(jsonQueryParams)
			params = ParseKeyValuePair(params, featureSessionVariables)
			queryParams, err = GetUrlQueryParams(params, featureSessionVariables)
			if err != nil {
				fmt.Println("--------> error in api call function <----------------")
				Errors["message"] = "error in api call function"
				Errors["api_error"] = fmt.Sprintf("%v", err)
				taskResponseObject.Errors = Errors
				taskResponseObject.Completed = false
				return taskResponseObject
			}
		}
	}
	urlWithParams = getURL + queryParams
	fmt.Println("--------> api endpoint <----------------" + urlWithParams)

	//-----------------------------request body----------------------------------------------------------------------------------
	if payload["body"] != nil {
		jsonPayloadBody, ok := payload["body"].([]interface{})
		if !ok {
			fmt.Println("--------> error in single_request payload[body] <----------------")
			Errors["message"] = "error in single_request payload[body]"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		if len(jsonPayloadBody) > 0 {
			bodyKeyValuePair := ParseObjectsFromConfigPayload(jsonPayloadBody)
			bodyKeyValuePair = ParseKeyValuePair(bodyKeyValuePair, featureSessionVariables)
			requestBody = ConvertKeyValuePairToInterface(bodyKeyValuePair)
		}
	}

	//-----------------------------headers---------------------------------------------------------------------------------------
	if payload["headers"] != nil {
		jsonPayloadHeaders, ok := payload["headers"].([]interface{})
		if !ok {
			fmt.Println("--------> error in single_request payload[headers] <----------------")
			Errors["message"] = "error in single_request payload[headers]"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		headers = ParseObjectsFromConfigPayload(jsonPayloadHeaders)
	}
	fmt.Println("--------> passing headers <----------------" + taskResponseObject.Name)

	//------------------------------request method-----------------------------------------------------------------------------
	if endpointTaskFile["request_configuration"] != nil {
		jsonRequestConfiguration, ok := endpointTaskFile["request_configuration"].(map[string]interface{})
		if !ok {
			fmt.Println("--------> error in single_request request_configuration <----------------")
			Errors["message"] = "error in single_request request_configuration"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		requestMethod, ok = jsonRequestConfiguration["method"].(string)
		if !ok {
			fmt.Println("--------> error in single_request request_configuration[method] <----------------")
			Errors["message"] = "error in single_request request_configuration[method]"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
	}

	//-------------call the MakeAPIRequest function and get the response of the API call--------------------------------------------
	response, err := MakeAPIRequest(requestMethod, urlWithParams, headers, requestBody, featureSessionVariables, endpointTaskFile)
	if err != nil {
		fmt.Println("--------> error in api call function <----------------")
		Errors["message"] = "error in api call function"
		Errors["api_error"] = fmt.Sprintf("%v", err)
		taskResponseObject.Errors = Errors
		taskResponseObject.Completed = false
		return taskResponseObject
	}

	//-----------------customize the API response------------------------------------------------
	if endpointTaskFile["response"] != nil {
		jsonResponse, ok := endpointTaskFile["response"].(map[string]interface{})
		if !ok {
			fmt.Println("--------> error in single_request response <----------------")
			Errors["message"] = "error in single_request response"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		//------------------------1.session variables--------------------------------------------------------------------------
		if jsonResponse["session_variables"] != nil {
			jsonResponseSessionVariables, ok := jsonResponse["session_variables"].([]interface{})
			if !ok {
				fmt.Println("--------> error in single_request response[session_variables] <----------------")
				Errors["message"] = "error in single_request response[session_variables]"
				taskResponseObject.Errors = Errors
				taskResponseObject.Completed = false
				return taskResponseObject
			}

			taskSessionVariables, err := ConvertArrayInterfaceToArrayStructWithJsonKeyParse(jsonResponseSessionVariables, response, featureSessionVariables)
			if err != nil {
				fmt.Println("--------> error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function <----------------")
				Errors["message"] = "error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function"
				Errors["api_error"] = fmt.Sprintf("%v", err)
				taskResponseObject.Errors = Errors
				taskResponseObject.Completed = false
				return taskResponseObject
			}
			taskSessionVariables = AppendOrUpdateKeyValuePair(taskSessionVariables, featureSessionVariables)
			taskResponseObject.SessionVariables = taskSessionVariables
		}
	}

	//------------------task response ------------------------------------------------------------------------------------------
	taskResponseObject.EndpointResponse = response
	taskResponseObject.Completed = true

	return taskResponseObject
}

func (f *Functions) ArrayToSingleValue(arr_obj interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {
	fmt.Println(arr_obj)
	arr := arr_obj.([]interface{})
	fmt.Println(arr[0].(string))
	// if len(arr) == 1{
	// 	value := arr[0].(string)
	// 	return value,nil
	// }
	return arr[0].(string), nil
}

func (f *Functions) GetDateTime(isUtc bool, format string, delay float64, featureSessionVariables []model.KeyValuePair) string {

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

func (f *Functions) GetPathFromUrl(url string, featureSessionVariables []model.KeyValuePair) string {

	urlList := strings.Split(url, "?")

	fmt.Println("\n\n -------->", urlList)

	return strings.Split(urlList[0], "//")[1]
}

func (f *Functions) GetParamsFromUrl(urlString string, keyName string, featureSessionVariables []model.KeyValuePair) string {

	urlList := strings.Split(urlString, "?")

	paramsList := strings.Split(urlList[1], "&")

	queryParams := make(map[string]interface{}, 0)

	for _, param := range paramsList {
		x := strings.Split(param, "=")
		queryParams[x[0]], _ = url.QueryUnescape(x[1])
	}

	fmt.Println("\n queryParams ------->", queryParams)

	return queryParams[keyName].(string)
}

func (f *Functions) TsvToJson(data string, featureSessionVariables []model.KeyValuePair, optionalParams ...interface{}) (map[string]interface{}, error) {

	var jsonArrayResponse []interface{}
	var jsonResponse = make(map[string]interface{}, 0)

	responseList := strings.Split(data, "\n")
	headersList := strings.Split(responseList[0], "\t")

	for index := 1; index < len(responseList)-1; index++ {
		valuesList := strings.Split(responseList[index], "\t")
		var appendObj = make(map[string]interface{}, 0)
		for innerIndex, key := range headersList {
			appendObj[key] = valuesList[innerIndex]
		}
		jsonArrayResponse = append(jsonArrayResponse, appendObj)
	}
	response, _ := json.MarshalIndent(jsonArrayResponse, "", "\t")
	fmt.Println("\nTsvToJson converted response -------------->", string(response))

	jsonResponse["data"] = jsonArrayResponse
	return jsonResponse, nil

}

func (f *Functions) XmlToInterface(data string, featureSessionVariables []model.KeyValuePair, optionalParams ...interface{}) (map[string]interface{}, error) {
	var jsonData = make(map[string]interface{}, 1)

	// TODO :- to be converted from xml to interface , now to handle the error, its been used .
	jsonData["data"] = data

	return jsonData, nil
}

func (f *Functions) AesDecrypt(initialVector string, key string, featureSessionVariables []model.KeyValuePair, crypt []byte) (string, error) {

	decodedKey, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		return "", err
	}
	fmt.Println("decoded key -------->", key, decodedKey)

	decodedIv, err := base64.StdEncoding.DecodeString(initialVector)
	if err != nil {
		return "", err
	}
	fmt.Println("decoded key -------->", initialVector, decodedIv)

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		fmt.Println("err", err)
		return "", err
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(decodedIv))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)
	decrypted = f.PKCS5Trimming(decrypted)

	fmt.Println("decryted data----------------->", string(decrypted))

	return string(decrypted), nil
}

func (f *Functions) PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func (f *Functions) NdJsonToJson(featureSessionVariables []model.KeyValuePair, ndJson []byte) (map[string]interface{}, error) {

	fmt.Println("recieved ndjson ------------>", string(ndJson))

	var jsonResponse = make(map[string]interface{}, 0)

	arrayResponse := strings.Split(string(ndJson), "\n")

	arrayResponse = arrayResponse[0 : len(arrayResponse)-1]

	var jsonArray = make([]interface{}, 0)

	for _, value := range arrayResponse {
		var jsonObj interface{}
		json.Unmarshal([]byte(value), &jsonObj)
		jsonArray = append(jsonArray, jsonObj)
	}

	fmt.Println("len of arr --->", len(arrayResponse))

	jsonResponse["data"] = jsonArray

	return jsonResponse, nil

}

func (f *Functions) UnixMilliToDateConversion(s interface{}, featureSessionVariables []model.KeyValuePair) (time.Time, error) {
	s_string := s.(string)
	open_index := strings.Index(s_string, "(")
	value := s_string[open_index+1 : open_index+11]
	int_value, _ := strconv.Atoi(value)
	data := time.Unix(int64(int_value), 0)
	return data, nil
}

func (f *Functions) CurrentTimeStamp(isMilli bool, featureSessionVariables []model.KeyValuePair) string {
	currentTimeInUnix := time.Now().UnixNano()
	if isMilli {
		currentTimeInUnix = currentTimeInUnix / int64(time.Millisecond)
	}
	return strconv.FormatInt(currentTimeInUnix, 10)
}

func (f *Functions) FormatString(format string, params ...interface{}) (string, error) {

	featureSessionVariables := params[len(params)-1].([]model.KeyValuePair)

	for index, param := range params {

		featureSessionVariables = append(featureSessionVariables, model.KeyValuePair{
			Key:   "param" + fmt.Sprint(index+1),
			Value: param,
			Type:  "static",
		})
	}
	// fmt.Println(param)
	regex := regexp.MustCompile(`\{{(.*?)\}}`)
	keys := regex.FindAllStringSubmatch(format, -1)

	for _, key := range keys {
		result, err := GetValueFromSessionVariablesKey(key[1], featureSessionVariables)
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

func (f *Functions) ExponentialToString(number float64, featureSessionVariables []model.KeyValuePair) string {
	return strconv.FormatFloat(number, 'f', 0, 64)
}

func (f *Functions) SleepExecution(sleepTimeInSecs float64, featureSessionVariables []model.KeyValuePair) {
	sleepTime := time.Duration(int(sleepTimeInSecs))

	time.Sleep(sleepTime * time.Second)
}

func (f *Functions) GetOneOfManyValues(params ...interface{}) interface{} {

	paramsLen := len(params)

	params = params[0 : paramsLen-1]

	for _, value := range params {
		if value != nil {
			return value
		}
	}

	return nil
}

func (f *Functions) InterfaceToXml(key string, data interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {

	marshalData, _ := json.Marshal(data)
	json.Unmarshal(marshalData, &data)

	xmlData, err := InterfaceToXml(data, key)
	fmt.Println("----------------xmlData------------", xmlData)
	if err != nil {
		return xmlData, err
	}
	return xmlData, nil
}

func (f *Functions) AesEncrypt(iv string, key string, data string, optionalParams ...interface{}) ([]byte, error) {

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	fmt.Println("decoded key -------->", key, decodedKey)

	decodedIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	fmt.Println("decoded key -------->", iv, decodedIv)

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(decodedIv))

	content := f.PKCS5Padding([]byte(data), block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	fmt.Println("Crypted ------------->", crypted)
	return crypted, nil
}

func (f *Functions) PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (f *Functions) FormatDateTime(input, inputFormat, outputFormat string, featureSessionVariables []model.KeyValuePair) (string, error) {

	dateTime, err := time.Parse(inputFormat, input)
	if err != nil {
		return "", nil
	}

	output := dateTime.Format(outputFormat)

	return output, nil
}

func (f *Functions) SingleValueToArray(value interface{}, featureSessionVariables []model.KeyValuePair) []interface{} {

	var arrayResponse = make([]interface{}, 0)

	arrayResponse = append(arrayResponse, value)

	return arrayResponse
}
func (f *Functions) ReturnMatchedObjectFromArray(array interface{}, keyofObject interface{}, valueofKey interface{}, featureSessionVariables []model.KeyValuePair) (interface{}, error) {
	inputArray, ok := array.([]map[interface{}]interface{})
	if !ok {
		return nil, errors.New("cannot convert to []map[string]interface{}")
	}
	for _, object := range inputArray {
		if object[keyofObject] != nil {
			if object[keyofObject] == valueofKey {
				return object, nil
			}
		}
	}
	return nil, nil
}
