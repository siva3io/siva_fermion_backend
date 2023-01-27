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
	pkg_helpers "fermion/backend_core/pkg/util/helpers"
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

// ====================================function calls===================================================================================================================================================
func CallFunctionFromProps(propsArr []model.Props, optionalParams ...interface{}) (interface{}, error) {
	var valueToReturn []reflect.Value
	var err error
	var ReturnValue interface{}

	SessionVariables, ok := optionalParams[0].([]model.KeyValuePair)
	if !ok {
		return ReturnValue, errors.New("first parameter must be feature session variables")
	}
	var featureSessionVariables []model.KeyValuePair
	featureSessionVariables = append(featureSessionVariables, SessionVariables...)

	fmt.Println("=============> Function call from Props <=============")

	for index := 0; index < len(propsArr); index++ {
		var paramsAfterEvaluation = make([]interface{}, 0)

		for _, item := range propsArr[index].Params {
			sessionValue := item.Value
			if item.Type == "variable" {
				sessionValue = GetValueFromSessionVariablesKey(item.Value.(string), featureSessionVariables)
			}
			paramsAfterEvaluation = append(paramsAfterEvaluation, sessionValue)
		}
		optionalParams[0] = featureSessionVariables
		paramsAfterEvaluation = append(paramsAfterEvaluation, optionalParams...)

		//==============================call function by name================================================================================================================
		valueToReturn, err = CallFuncByName(&Functions{}, propsArr[index].Name, paramsAfterEvaluation...)
		if err != nil {
			return ReturnValue, err
		}

		if len(valueToReturn) > 0 {
			functionReturnValue := valueToReturn[0]

			//=================================temporary error handle========================================================================================================
			if len(valueToReturn) > 1 {
				FunctionError := valueToReturn[1]
				errorValue, ok := FunctionError.Interface().(error)
				if ok {
					return ReturnValue, errorValue
				}
			}

			//=============================data type conversion==============================================================================================================
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
func CallFuncByName(myClass interface{}, funcName string, optionalParams ...interface{}) (out []reflect.Value, err error) {

	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)

	if !m.IsValid() {
		errorString := fmt.Sprintf("method not found \"%s\"", funcName)
		return make([]reflect.Value, 0), errors.New(errorString)
	}
	in := make([]reflect.Value, len(optionalParams))
	fmt.Printf("=============> function name =============> %v(", funcName)
	for index, param := range optionalParams {
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
func ReturnTheVariableType(value reflect.Value, optionalParams ...interface{}) (interface{}, error) {
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

// ==================================ipaas props functions==============================================================================================================================================
func (f *Functions) TokenBearerType(token interface{}, optionalParams ...interface{}) (string, error) {
	tokenString, ok := token.(string)
	if !ok {
		return "", errors.New("oops! token must be a string")
	}
	output, err := pkg_helpers.TokenBearerType(tokenString)
	if err != nil {
		return output, err
	}

	//  fmt.Println("TokenBearerType ===>>>>", output)
	return output, nil
}
func (f *Functions) Paginate(initialValue, pageValue, stepValue interface{}, optionalParams ...interface{}) (interface{}, error) {

	if pageValue == nil || pageValue == "" {
		return initialValue, nil
	}
	IntegerCurrentPageNumber, err := strconv.Atoi(pageValue.(string))
	if err != nil {
		return nil, err
	}
	IntegerStepValue, _ := strconv.Atoi(stepValue.(string))
	if err != nil {
		return nil, err
	}

	//  fmt.Println("Paginate ===>>>>", strconv.Itoa(IntegerCurrentPageNumber + IntegerStepValue))
	return strconv.Itoa(IntegerCurrentPageNumber + IntegerStepValue), nil
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
		errorString := fmt.Sprintf("oops! %v is not available in response", key)
		return false, errors.New(errorString)
	}
	responseValue, err := ParseJsonPathFromObject(response.Body, key)
	if err != nil {
		return false, err
	}
	// fmt.Printf("OPERATOR ===>>>> %v , Prev Resp ===>>>> %v , conditionValue ===>>>> %v ,actualValue ===>>>> %v", operator, prevResponse, fmt.Sprintf("%v", value), fmt.Sprintf("%v", responseValue))
	if responseValue == nil {
		return true, nil
	}
	if operator == "OR" {
		return prevResponse || fmt.Sprintf("%v", responseValue) == fmt.Sprintf("%v", value), nil
	}
	if operator == "AND" {
		return prevResponse && fmt.Sprintf("%v", responseValue) == fmt.Sprintf("%v", value), nil
	}

	//  fmt.Println("ConditionNorms ===>>>>")
	return true, nil
}
func (f *Functions) ParseJsonPathFromObject(data interface{}, keypath string, optionalParams ...interface{}) (interface{}, error) {

	if data == nil {
		return nil, nil
	}
	object, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("oops! cannot convert to object")
	}
	res, err := ParseJsonPathFromObject(object, keypath)
	if err != nil {
		return res, err
	}

	//  fmt.Println("ParseJsonPathFromObject ===>>>>", res)
	return res, nil
}
func (f *Functions) ReturnIndexObjectFromArray(array interface{}, index string, optionalParams ...interface{}) (interface{}, error) {
	output := map[string]interface{}{}
	var inputArray []map[string]interface{}

	err := pkg_helpers.JsonMarshaller(array, &inputArray)
	if err != nil {
		return nil, err
	}

	i, err := strconv.Atoi(index)
	if err != nil {
		return output, err
	}
	if len(inputArray) == 0 {
		return output, nil
	}

	if i < 0 || i >= len(inputArray) {
		return output, errors.New("oops! index out of range for ReturnIndexObjectFromArray")
	}

	//  fmt.Println("ReturnIndexObjectFromArray ===>>>>", inputArray[i])
	return inputArray[i], nil
}
func (f *Functions) GetLength(array []interface{}, optionalParams ...interface{}) int {

	//  fmt.Println("GetLength ===>>>>", len(array))
	return len(array)
}
func (f *Functions) CheckStatus(optionalParams ...interface{}) bool {

	response := optionalParams[0].(model.EndpointResponse)

	//  fmt.Println("CheckStatus ===>>>>", response.StatusCode < 200 || response.StatusCode >= 300)
	return response.StatusCode < 200 || response.StatusCode >= 300
}
func (f *Functions) ParseKeyValuePairThenReturnJSONString(input interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {
	array, ok := input.([]interface{})
	if !ok {
		return "", errors.New("oops! issue in coversion_type from interface{} to []interface{}")
	}

	keyValuePairs, err := ParseObjectsFromConfigPayload(array)
	if err != nil {
		return "", err
	}
	keyValuePairs = ParseKeyValuePair(keyValuePairs, featureSessionVariables)
	output := ConvertKeyValuePairToInterface(keyValuePairs)

	jsonByteArray, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	//  fmt.Println("ParseKeyValuePairThenReturnJSONString ===>>>>", string(jsonByteArray))
	return string(jsonByteArray), nil
}
func (f *Functions) AppendToArray(input interface{}, optionalParams ...interface{}) ([]interface{}, error) {
	array := []interface{}{}
	if input == nil {
		return array, errors.New("oops! input must be non-nil")
	}
	array = append(array, input)

	//  fmt.Println("AppendToArray ===>>>>", array)
	return array, nil
}
func (f *Functions) GetValueFromSessionVariables(value interface{}, optionalParams ...interface{}) interface{} {
	// whatever the value getting as a input just return it because, already processed in params itself

	//  fmt.Println("GetValueFromSessionVariables ===>>>>", value)
	return value
}

// ==================================ipaas type_conversion or type_casting props functions=====================================================================================================================
func (f *Functions) StringtoFloat(input interface{}, optionalParams ...interface{}) (float32, error) {
	if input == nil {
		return 0, nil
	}
	stringValue, ok := input.(string)
	if !ok {
		return 0, errors.New("oops! invalid input_type for StringtoFloat Function")
	}
	if stringValue == "" {
		return 0, nil
	}
	output, err := strconv.ParseFloat(stringValue, 32)
	if err != nil {
		return 0, err
	}
	//  fmt.Println("StringtoFloat ===>>>>", float32(output))
	return float32(output), nil
}
func (f *Functions) FloattoString(input interface{}, optionalParams ...interface{}) (string, error) {
	if input == nil {
		return "", nil
	}
	floatValue, ok := input.(float64)
	if !ok {
		return "", errors.New("oops! invalid input_type for FloattoString Function")
	}

	//  fmt.Println("StringtoFloat ===>>>>", fmt.Sprintf("%v", floatValue))
	return fmt.Sprintf("%v", floatValue), nil
}
func (f *Functions) FloattoInt(input interface{}, optionalParams ...interface{}) (int, error) {
	if input == nil {
		return 0, nil
	}
	floatValue, ok := input.(float64)
	if !ok {
		return 0, errors.New("oops! invalid input_type for FloattoInt Function")
	}

	//  fmt.Println("FloattoInt ===>>>>", int(floatValue))
	return int(floatValue), nil
}
func (f *Functions) ExponentialToString(number float64, optionalParams ...interface{}) string {

	//  fmt.Println("ExponentialToString ===>>>>", strconv.FormatFloat(number, 'f', 0, 64))
	return strconv.FormatFloat(number, 'f', 0, 64)
}
func (f *Functions) TsvToJson(data string, optionalParams ...interface{}) (map[string]interface{}, error) {

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

	jsonResponse["data"] = jsonArrayResponse

	//  fmt.Println("TsvToJson ===>>>>", jsonResponse)
	return jsonResponse, nil
}
func (f *Functions) InterfaceToXml(key string, data interface{}, optionalParams ...interface{}) (string, error) {

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
			formatedString, _ := f.InterfaceToXml(seqString, obj[seqString])
			xmlString += formatedString
		}
		if key == "" {
			return xmlString, nil
		}
		//  fmt.Println("InterfaceToXml ===>>>>", fmt.Sprintf("<%v>%v</%v>", leftKey, xmlString, rightKey))
		return fmt.Sprintf("<%v>%v</%v>", leftKey, xmlString, rightKey), nil

	case []interface{}:
		arrayData := data.([]interface{})
		for _, value := range arrayData {
			formatedString, _ := f.InterfaceToXml(key, value)
			xmlString += formatedString
		}
		//  fmt.Println("InterfaceToXml ===>>>>", xmlString)
		return xmlString, nil

	default:
		//  fmt.Println("InterfaceToXml ===>>>>",fmt.Sprintf("<%v>%v</%v>", leftKey, data, rightKey))
		return fmt.Sprintf("<%v>%v</%v>", leftKey, data, rightKey), nil
	}
}
func (f *Functions) XmlToInterface(data string, optionalParams ...interface{}) (map[string]interface{}, error) {
	var jsonData = make(map[string]interface{}, 1)
	// TODO :- to be converted from xml to interface , now to handle the error, its been used .
	jsonData["data"] = data

	//  fmt.Println("XmlToInterface ===>>>>", jsonData)
	return jsonData, nil
}
func (f *Functions) NdJsonToJson(ndJson []byte, optionalParams ...interface{}) (map[string]interface{}, error) {
	var jsonResponse = make(map[string]interface{}, 0)
	arrayResponse := strings.Split(string(ndJson), "\n")
	arrayResponse = arrayResponse[0 : len(arrayResponse)-1]

	var jsonArray = make([]interface{}, 0)

	for _, value := range arrayResponse {
		var jsonObj interface{}
		json.Unmarshal([]byte(value), &jsonObj)
		jsonArray = append(jsonArray, jsonObj)
	}

	jsonResponse["data"] = jsonArray

	//  fmt.Println("NdJsonToJson ===>>>>", jsonResponse)
	return jsonResponse, nil
}

// ================================mapper api_call props functions========================================================================================================================================
func (f *Functions) ExecuteGetApiCall(taskObject map[string]interface{}, featureSessionVariables []model.KeyValuePair) (interface{}, error) {

	taskResponseArray := []model.TaskEndpointResponse{}
	taskResponseObject := model.TaskEndpointResponse{}

	var endpointTaskFile model.APIFile

	endpointTaskFile.Scheme = "http"
	endpointTaskFile.Endpoint = taskObject["endpoint"].(string)
	endpointTaskFile.RequestConfiguration.Method = "GET"
	endpointTaskFile.RequestConfiguration.Type = "single"
	endpointTaskFile.Payload.Headers = []model.KeyValuePair{
		{Key: "Content-Type", Value: "application/json", Type: "static"},
		{Key: "Authorization", Value: "Authorization", Type: "variable"}}

	//========================query params======================================================================
	query_params := []model.KeyValuePair{}
	if taskObject["query_params"] != nil {
		jsonQueryParamsArray, ok := taskObject["query_params"].([]interface{})
		if !ok {
			return nil, errors.New("oops! invalid query_params conversion from interface{} to []interface{} in ExecuteGetApiCall")
		}
		if len(jsonQueryParamsArray) > 0 {
			for _, josnQueryParam := range jsonQueryParamsArray {
				param, ok := josnQueryParam.(map[string]interface{})
				if !ok {
					return nil, errors.New("oops! invalid query_params conversion from interface to map[string]interface{} in ExecuteGetApiCall")
				}
				//====================filters=================================================
				if param["key"] == "filters" {
					res, err := f.FormatFilter(param["value"], featureSessionVariables)
					if err != nil {
						return nil, err
					}
					query_params = append(query_params, model.KeyValuePair{Key: param["key"].(string), Value: res, Type: "static"})
					continue
				}
				//====================sorts===================================================
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

	endpointTaskFile.Payload.QueryParams = append(endpointTaskFile.Payload.QueryParams, query_params...)

	//=======================================execute the get_request_method=============================================================================================================================
	response, err := f.ExecuteGetRequest(taskResponseArray, taskResponseObject, endpointTaskFile, featureSessionVariables)
	if err != nil {
		return nil, err
	}

	fmt.Println("ExecuteGetApiCall ===>>>>", response)
	return response, nil
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
		return "", errors.New("oops! input must not be nil in FormatFilter")
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
			interfaceValue := GetValueFromSessionVariablesKey(filter.Value.(string), featureSessionVariables)
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

	//  fmt.Println("FormatFilter ===>>>>", output)
	return output, nil
}
func (f *Functions) FormatSort(input interface{}, featureSessionVariables []model.KeyValuePair) (string, error) {
	type Sort struct {
		Column  string `json:"column"`
		OrderBy string `json:"order_by"`
	}

	var sorts []Sort

	if input == nil {
		return "", errors.New("oops! input must not be nil in FormatSort")
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

	//  fmt.Println("FormatSort ===>>>>", output)
	return output, nil
}
func (f *Functions) ExecuteGetRequest(taskResponseArray []model.TaskEndpointResponse, taskResponseObject model.TaskEndpointResponse, endpointTaskFile model.APIFile, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var requestBody interface{}
	var queryParams string
	var urlWithParams string
	var requestMethod string

	payload := endpointTaskFile.Payload
	getURL := endpointTaskFile.Scheme + "://" + endpointTaskFile.Endpoint
	getURL, err := FormatEndpoint(getURL, featureSessionVariables)
	if err != nil {
		return nil, err
	}

	if len(payload.QueryParams) > 0 {
		params := ParseKeyValuePair(payload.QueryParams, featureSessionVariables)
		queryParams, err = GetUrlQueryParams(params, featureSessionVariables)
		if err != nil {
			return nil, err
		}
	}
	urlWithParams = getURL + queryParams
	headers := payload.Headers
	requestMethod = endpointTaskFile.RequestConfiguration.Method

	//==============call the MakeAPIRequest function and get the response of the API call==============================================================================================
	response, err := MakeAPIRequest(requestMethod, urlWithParams, headers, requestBody, featureSessionVariables, endpointTaskFile)
	if err != nil {
		return nil, err
	}

	fmt.Println("ExecuteGetRequest ===>>>>", response.Body)
	return response.Body, nil
}

// ==================================ipaas time_and_date related props functions=====================================================================================================================
func (f *Functions) GetDateTime(isUtc bool, format string, delay float64, optionalParams ...interface{}) string {
	output := pkg_helpers.GetDateTime(isUtc, format, delay)

	//  fmt.Println("GetDateTime ===>>>>", output)
	return output
}
func (f *Functions) UnixMilliToDateConversion(s interface{}, optionalParams ...interface{}) (time.Time, error) {
	s_string := s.(string)
	open_index := strings.Index(s_string, "(")
	value := s_string[open_index+1 : open_index+11]
	int_value, _ := strconv.Atoi(value)
	data := time.Unix(int64(int_value), 0)

	//  fmt.Println("UnixMilliToDateConversion ===>>>>", data)
	return data, nil
}
func (f *Functions) CurrentTimeStamp(isMilli bool, optionalParams ...interface{}) string {
	currentTimeInUnix := time.Now().UnixNano()
	if isMilli {
		currentTimeInUnix = currentTimeInUnix / int64(time.Millisecond)
	}

	//  fmt.Println("CurrentTimeStamp ===>>>>", strconv.FormatInt(currentTimeInUnix, 10))
	return strconv.FormatInt(currentTimeInUnix, 10)
}
func (f *Functions) SleepExecution(sleepTimeInSecs float64, optionalParams ...interface{}) {
	sleepTime := time.Duration(int(sleepTimeInSecs))

	time.Sleep(sleepTime * time.Second)
}
func (f *Functions) FormatDateTime(input, inputFormat, outputFormat string, optionalParams ...interface{}) (string, error) {

	dateTime, err := time.Parse(inputFormat, input)
	if err != nil {
		return "", err
	}
	output := dateTime.Format(outputFormat)

	//  fmt.Println("FormatDateTime ===>>>>", output)
	return output, nil
}

// ==================================ipaas props functions=====================================================================================================================
func (f *Functions) GetPathFromUrl(url string, optionalParams ...interface{}) string {
	urlList := strings.Split(url, "?")

	//  fmt.Println("GetPathFromUrl ===>>>>", strings.Split(urlList[0], "//")[1])
	return strings.Split(urlList[0], "//")[1]
}
func (f *Functions) GetParamsFromUrl(urlString string, keyName string, optionalParams ...interface{}) string {

	urlList := strings.Split(urlString, "?")
	paramsList := strings.Split(urlList[1], "&")
	queryParams := make(map[string]interface{}, 0)
	for _, param := range paramsList {
		x := strings.Split(param, "=")
		queryParams[x[0]], _ = url.QueryUnescape(x[1])
	}

	//  fmt.Println("GetParamsFromUrl ===>>>>", queryParams[keyName].(string))
	return queryParams[keyName].(string)
}
func (f *Functions) FormatString(format string, optionalParams ...interface{}) (string, error) {

	featureSessionVariables := optionalParams[len(optionalParams)-1].([]model.KeyValuePair)
	for index, param := range optionalParams {

		featureSessionVariables = append(featureSessionVariables, model.KeyValuePair{
			Key:   "param" + fmt.Sprint(index+1),
			Value: param,
			Type:  "static",
		})
	}
	regex := regexp.MustCompile(`\{{(.*?)\}}`)
	keys := regex.FindAllStringSubmatch(format, -1)

	for _, key := range keys {
		result := GetValueFromSessionVariablesKey(key[1], featureSessionVariables)
		value := fmt.Sprintf("%v", result)
		format = strings.Replace(format, key[0], value, -1)
	}
	// fmt.Println("FormatString --->>>>", format)
	return format, nil
}
func (f *Functions) ReturnMatchedObjectFromArray(array interface{}, keyofObject string, valueofKey interface{}, optionalParams ...interface{}) (interface{}, error) {
	inputArray, ok := array.([]map[string]interface{})
	if !ok {
		return nil, errors.New("oops! issue in conversion from interface{} to []map[string]interface{}")
	}
	for _, object := range inputArray {
		if object[keyofObject] != nil {
			if object[keyofObject] == valueofKey {

				//  fmt.Println("ReturnMatchedObjectFromArray ===>>>>", object)
				return object, nil
			}
		}
	}
	return nil, nil
}
func (f *Functions) GetOneOfManyValues(optionalParams ...interface{}) interface{} {

	paramsLen := len(optionalParams)
	optionalParams = optionalParams[0 : paramsLen-1]

	for _, value := range optionalParams {
		if value != nil {

			//  fmt.Println("GetOneOfManyValues ===>>>>", value)
			return value
		}
	}
	return nil
}
func (f *Functions) ArrofObjecttoArrofInterface(input interface{}, key string, optionalParams ...interface{}) (interface{}, error) {
	arrayOfObject := []map[string]interface{}{}
	err := pkg_helpers.JsonMarshaller(input, &arrayOfObject)
	if err != nil {
		return nil, err
	}
	var output []interface{}
	for _, obj := range arrayOfObject {
		if obj[key] != nil {
			output = append(output, obj[key])
		}
	}
	return output, nil
}

// === TODO: Below are redundant functions :-> drupadh needs to check it and remove from either pkg_helpers or here=========
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
func (f *Functions) AesDecrypt(initialVector string, key string, crypt []byte, optionalParams ...interface{}) (string, error) {

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
func (f *Functions) PKCS5Padding(ciphertext []byte, blockSize int, optionalParams ...interface{}) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func (f *Functions) PKCS5Trimming(encrypt []byte, optionalParams ...interface{}) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func (f *Functions) FieldMapping(input interface{}, mappedValue map[string]interface{}, optionalParams ...interface{}) (interface{}, error) {

	if len(optionalParams) > 1 && optionalParams[0] == true {
		for key, value := range mappedValue {
			if value == input {
				return key, nil
			}
		}
		return nil, nil
	}
	return mappedValue[input.(string)], nil
}
