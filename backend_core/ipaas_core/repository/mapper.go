package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fermion/backend_core/ipaas_core/model"
	"fermion/backend_core/ipaas_core/utils"
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

func ExecuteMapperForTask(mapperTasks []model.MapperTask, response map[string]interface{}, mappedResponseCompletionNorms, featureSessionVariables []model.KeyValuePair) ([]model.KeyValuePair, error) {

	var mapperTasksResponse []model.KeyValuePair

	//==============================lopping the number of mappers need to execute  for this response==============================================================================
	for _, task := range mapperTasks {
		var mapperResponse model.KeyValuePair
		var err error

		//=============================read the mapper file and convert to mapper struct format===================================================================================
		fileData, err := utils.ReadFeature("mappers", task.MapperId)
		if err != nil {
			return nil, err
		}
		var mapperConfigTemplateFile model.Mapper
		err = pkg_helpers.JsonMarshaller(fileData, &mapperConfigTemplateFile)
		if err != nil {
			return nil, err
		}

		//==============================get the input from the api response=======================================================================================================
		inputPath, err := utils.FormatEndpoint(task.InputKey, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		jsonInput, err := utils.ParseJsonPathFromObject(response, inputPath)
		if err != nil {
			return nil, err
		}

		//===========================Produce the mapper and consume the response====================================================================================================
		// requestId := uuid.New().String()
		// requestPayload := map[string]interface{}{}
		// requestPayload["meta_data"] = map[string]interface{}{
		// 	"request_id": requestId,
		// }
		// requestPayload["data"] = map[string]interface{}{
		// 	"input_data":                jsonInput,
		// 	"mapper_template":           mapperConfigTemplateFile,
		// 	"feature_session_variables": featureSessionVariables,
		// }
		// eda.Produce(eda.BOSON_MAPPER_CONVERTOR, requestPayload)

		//=================find the payload input_type======================================================
		mapperInputType := "object"
		_, ok := jsonInput.([]interface{})
		if ok {
			mapperInputType = "array"
		}

		var mappedResponse interface{}

		//=================execute object mapper=======================================================
		if mapperInputType == "object" {
			mappedResponse, err = ExecuteObjMapper(jsonInput, mapperConfigTemplateFile, featureSessionVariables)
			if err != nil {
				return nil, err
			}
		}
		//=================execute array mapper=======================================================
		if mapperInputType == "array" {
			var mapperInput []interface{}

			err = pkg_helpers.JsonMarshaller(mapperInputType, &mapperInput)
			if err != nil {
				return nil, err
			}
			mappedResponse, err = ExecuteArrayMapper(mapperInput, mapperConfigTemplateFile, featureSessionVariables)
			if err != nil {
				return nil, err
			}
		}

		// edaResponse := map[string]interface{}{}
		// for {
		// 	output := cache.GetCacheVariable(requestId)
		// 	if output != nil {
		// 		var resp map[string]interface{}
		// 		err := pkg_helpers.JsonMarshaller(output, &resp)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		metaData := resp["meta_data"].(map[string]interface{})
		// 		if metaData["encryption"].(bool) {
		// 			decrypted_payload, _ := pkg_helpers.AESGCMDecrypt(resp["data"].(string))
		// 			var decryptedResponseData interface{}
		// 			err := json.Unmarshal([]byte(decrypted_payload), &decryptedResponseData)
		// 			if err != nil {
		// 				return nil, err
		// 			}
		// 			resp["data"] = decryptedResponseData
		// 		}
		// 		err = pkg_helpers.JsonMarshaller(resp, &edaResponse)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		break
		// 	}
		// }
		// edaResponseData := edaResponse["data"].(map[string]interface{})
		// if edaResponseData["error_message"] != nil {
		// 	jsonerror, _ := json.Marshal(edaResponseData["error_message"])
		// 	return nil, errors.New(string(jsonerror))
		// }
		//==============================end of EDA================================================================================================================================
		mapperResponse = model.KeyValuePair{
			Key:   task.MapperName,
			Value: mappedResponse,
		}

		//====================================mapper completion Norms===========================================================================================================
		mapperResponse.Value, err = MappedResponseCompletionNorms(mappedResponseCompletionNorms, mapperResponse, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mapperTasksResponse = append(mapperTasksResponse, mapperResponse)
		continue
	}
	// pkg_helpers.PrettyPrint("Mapped Response", mapperTasksResponse)
	return mapperTasksResponse, nil
}
func MergeMultipleMapperPayload(required_mapped_output_array []model.RequiredMappedOutput, taskResponseArr []model.TaskEndpointResponse) interface{} {
	var payload interface{}
	for _, required_map := range required_mapped_output_array {
		if required_map.Type == "payload" {
			for _, taskResponse := range taskResponseArr {
				if taskResponse.Name == required_map.TaskName {
					for _, mapper := range taskResponse.MappedResponse {
						fmt.Println(mapper.Key, required_map.MapperName)
						if mapper.Key == required_map.MapperName {
							return mapper.Value
						}
					}
				}
			}
		}
	}
	return payload
}

func ExecuteArrayMapper(responseArrToMap []interface{}, mapperConfigTemplateFile model.Mapper, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	mappedArrData := make([]map[string]interface{}, 0)
	mapperFieldsArrToMap := mapperConfigTemplateFile.Fields
	for responseObjIndex := 0; responseObjIndex < len(responseArrToMap); responseObjIndex++ {
		var mappedObjData = make(map[string]interface{})
		responseObj, ok := responseArrToMap[responseObjIndex].(map[string]interface{})
		if !ok {
			return nil, errors.New("cannot convert input to map[string]interface{}")
		}
		for mapperFieldObjIndex := 0; mapperFieldObjIndex < len(mapperFieldsArrToMap); mapperFieldObjIndex++ {
			mapperFieldObj := mapperFieldsArrToMap[mapperFieldObjIndex]
			if mapperFieldObj.ExecutionType == "nested" {
				continue
			}
			mappedResponseObj, err := checkFieldTypeAndMap(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			if mapperFieldObj.OutputType == "array" {
				MergeInterfacesIntoArray(mappedObjData, mappedResponseObj)
				continue
			}
			MergeInterfaces(mappedObjData, mappedResponseObj)
		}
		mappedArrData = append(mappedArrData, mappedObjData)
	}
	return mappedArrData, nil
}
func ExecuteObjMapper(responseObjToMapper interface{}, mapperConfigTemplateFile model.Mapper, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	responseObj, ok := responseObjToMapper.(map[string]interface{})
	if !ok {
		return nil, errors.New("cannot convert input to map[string]interface{}")
	}
	mapperFieldsArrToMap := mapperConfigTemplateFile.Fields

	var mappedObjData = make(map[string]interface{})

	for mapperFieldObjIndex := 0; mapperFieldObjIndex < len(mapperFieldsArrToMap); mapperFieldObjIndex++ {
		mapperFieldObj := mapperFieldsArrToMap[mapperFieldObjIndex]
		if mapperFieldObj.ExecutionType == "nested" {
			continue
		}
		mappedResponseObj, err := checkFieldTypeAndMap(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		if mapperFieldObj.OutputType == "array" {
			MergeInterfacesIntoArray(mappedObjData, mappedResponseObj)
			continue
		}
		MergeInterfaces(mappedObjData, mappedResponseObj)
	}
	return mappedObjData, nil
}
func checkFieldTypeAndMap(mapperFieldObj model.MapperField, responseObj interface{}, mapperFieldsArrToMap []model.MapperField, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	if mapperFieldObj.InputType == "keyvalue" && mapperFieldObj.OutputType == "keyvalue" {
		mappedResForKey, err := MapKeyValueToKeyValueTypeWithMapper(mapperFieldObj, responseObj, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForKey)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "object" && mapperFieldObj.OutputType == "object" {
		mappedResForObject, err := MapObjectToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj.Output, mappedResForObject)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "array" && mapperFieldObj.OutputType == "array" {
		mappedResForArr, err := MapArrToArrTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj.Output, mappedResForArr)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "array" && mapperFieldObj.OutputType == "object" {
		mappedResForObject, err := MapArrToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj.Output, mappedResForObject)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "object" && mapperFieldObj.OutputType == "keyvalue" {
		mappedResForKeyValue, err := MapObjectToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForKeyValue)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "keyvalue" && mapperFieldObj.OutputType == "object" {
		mappedResForObject, err := MapKeyValueToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForObject)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "keyvalue" && mapperFieldObj.OutputType == "array" {
		mappedResForArray, err := MapKeyValueToArrayTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj.Output, mappedResForArray)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "object" && mapperFieldObj.OutputType == "array" {
		mappedResForArray, err := MapObjectToArrayTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj.Output, mappedResForArray)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj.InputType == "interface" && mapperFieldObj.OutputType == "keyvalue" {
		mappedResForKey, err := MapInterfaceToKeyValueTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForKey)
		if err != nil {
			return nil, err
		}
	}
	return mappedObjData, nil
}

func MapKeyValueToKeyValueTypeWithMapper(mapperFieldObj model.MapperField, responseInterface interface{}, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}
	valueObj := responseObj[mapperFieldObj.Input]
	if valueObj != nil {
		mappedObjData[mapperFieldObj.Output] = responseObj[mapperFieldObj.Input]
		//----helper function ----------
		if len(mapperFieldObj.HelperFunction) > 0 {
			helperFunctionOutput, err := MapperHelperFunction(mapperFieldObj, mappedObjData, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			mappedObjData[mapperFieldObj.Output] = helperFunctionOutput
		}
		return mappedObjData, nil
	}
	defaultValueToMap := mapperFieldObj.DefaultValue
	if defaultValueToMap != false {
		mappedObjData[mapperFieldObj.Output] = defaultValueToMap
	}
	//----helper function ----------
	if len(mapperFieldObj.HelperFunction) > 0 {
		helperFunctionOutput, err := MapperHelperFunction(mapperFieldObj, mappedObjData, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData[mapperFieldObj.Output] = helperFunctionOutput
		return mappedObjData, nil
	}
	return mappedObjData, nil
}
func MapObjectToObjectTypeWithMapper(mapperFieldObj model.MapperField, responseInterface interface{}, mapperFieldsArrToMap []model.MapperField, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}
	valueObj := responseObj[mapperFieldObj.Input]
	if valueObj == nil {
		valueObj = responseObj
	}
	if valueObj != nil {
		// fmt.Println(valueObj)
		mappedToBe := valueObj.(map[string]interface{})
		fieldsToMapWithinObj := mapperFieldObj.Fields
		for fieldToMapIndex := 0; fieldToMapIndex < len(fieldsToMapWithinObj); fieldToMapIndex++ {
			fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[fieldToMapIndex], mapperFieldsArrToMap)
			if err != nil {
				return nil, err
			}
			ObjectToBeMapped, err := checkFieldTypeAndMap(fieldObj, mappedToBe, mapperFieldsArrToMap, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			if fieldObj.OutputType == "array" {
				MergeInterfacesIntoArray(mappedObjData, ObjectToBeMapped)
				continue
			}
			mappedObjData, err = MergeInterfaces(mappedObjData, ObjectToBeMapped)
			if err != nil {
				return nil, err
			}
		}
		return mappedObjData, nil
	}
	return mappedObjData, nil
}
func MapArrToArrTypeWithMapper(mapperFieldObj model.MapperField, responseInterface interface{}, mapperFieldsArrToMap []model.MapperField, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	var mappedArrData = make([]map[string]interface{}, 0)
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}
	if responseObj[mapperFieldObj.Input] == nil {
		return mappedArrData, nil
	}
	responseArr := responseObj[mapperFieldObj.Input].([]interface{})
	for responseObjIndex := 0; responseObjIndex < len(responseArr); responseObjIndex++ {
		responseObj := responseArr[responseObjIndex]
		var mappedObjData = make(map[string]interface{})
		if responseObj != nil {
			objToMapWithinArr := mapperFieldObj.Fields
			for fieldToMapIndex := 0; fieldToMapIndex < len(objToMapWithinArr); fieldToMapIndex++ {
				fieldObj, err := fetchFieldObjFromFieldsArr(objToMapWithinArr[fieldToMapIndex], mapperFieldsArrToMap)
				if err != nil {
					return nil, err
				}
				ObjectToBeMapped, err := checkFieldTypeAndMap(fieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
				if err != nil {
					return nil, err
				}
				mappedObjData, err = MergeInterfaces(mappedObjData, ObjectToBeMapped)
				if err != nil {
					return nil, err
				}
			}
		}
		mappedArrData = append(mappedArrData, mappedObjData)
	}
	return mappedArrData, nil
}
func MapArrToObjectTypeWithMapper(mapperFieldObj model.MapperField, responseInterface interface{}, mapperFieldsArrToMap []model.MapperField, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	keyForInputWithIndex := mapperFieldObj.Input
	keyNameArr := strings.Split(keyForInputWithIndex, "[")
	var mappedObjData = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}

	responseArr, ok := responseObj[keyNameArr[0]].([]interface{})
	if !ok {
		responseArr = []interface{}{}
	}
	responseObjIndex, _ := strconv.Atoi(strings.Split(keyNameArr[1], "]")[0])

	responseObj = map[string]interface{}{}
	if responseObjIndex < len(responseArr) {
		responseObj = responseArr[responseObjIndex].(map[string]interface{})
	}

	if responseObj != nil {
		objToMapWithinArr := mapperFieldObj.Fields
		for fieldToMapIndex := 0; fieldToMapIndex < len(objToMapWithinArr); fieldToMapIndex++ {
			fieldObj, err := fetchFieldObjFromFieldsArr(objToMapWithinArr[fieldToMapIndex], mapperFieldsArrToMap)
			if err != nil {
				return nil, err
			}
			ObjectToBeMapped, err := checkFieldTypeAndMap(fieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			mappedObjData, err = MergeInterfaces(mappedObjData, ObjectToBeMapped)
			if err != nil {
				return nil, err
			}
		}
	}
	return mappedObjData, nil
}
func MapKeyValueToObjectTypeWithMapper(mapperFieldObj model.MapperField, responseInterface interface{}, mapperFieldsArrToMap []model.MapperField, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	var valueObj = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}

	valueObj[mapperFieldObj.Input] = responseObj[mapperFieldObj.Input]
	if len(valueObj) == 0 {
		valueObj = responseObj
	}
	fieldsToMapWithinObj := mapperFieldObj.Fields
	fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[0], mapperFieldsArrToMap)
	if err != nil {
		return nil, err
	}

	mappedObjData[mapperFieldObj.Output], err = checkFieldTypeAndMap(fieldObj, valueObj, mapperFieldsArrToMap, featureSessionVariables)
	if err != nil {
		return nil, err
	}
	return mappedObjData, nil
}
func MapKeyValueToArrayTypeWithMapper(mapperFieldObj model.MapperField, responseInterface interface{}, mapperFieldsArrToMap []model.MapperField, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	var mappedObjData = make([]map[string]interface{}, 0)
	var valueObj = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}

	valueObj[mapperFieldObj.Input] = responseObj[mapperFieldObj.Input]
	if len(valueObj) == 0 {
		valueObj = responseObj
	}
	fieldsToMapWithinObj := mapperFieldObj.Fields
	fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[0], mapperFieldsArrToMap)
	if err != nil {
		return nil, err
	}
	ObjectToBeMapped, err := checkFieldTypeAndMap(fieldObj, valueObj, mapperFieldsArrToMap, featureSessionVariables)
	if err != nil {
		return nil, err
	}
	mappedObjData = append(mappedObjData, ObjectToBeMapped)
	return mappedObjData, nil
}
func MapObjectToArrayTypeWithMapper(mapperFieldObj model.MapperField, responseInterface interface{}, mapperFieldsArrToMap []model.MapperField, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	var mappedArrData = make([]map[string]interface{}, 0)
	var mappedObjData = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}

	valueObj := responseObj[mapperFieldObj.Input]
	if valueObj == nil {
		valueObj = responseObj
	}
	if valueObj != nil {
		fieldsToMapWithinObj := mapperFieldObj.Fields
		for fieldToMapIndex := 0; fieldToMapIndex < len(fieldsToMapWithinObj); fieldToMapIndex++ {
			fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[fieldToMapIndex], mapperFieldsArrToMap)
			if err != nil {
				return nil, err
			}
			ObjectToBeMapped, err := checkFieldTypeAndMap(fieldObj, valueObj.(map[string]interface{}), mapperFieldsArrToMap, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			if fieldObj.OutputType == "array" {
				MergeInterfacesIntoArray(mappedObjData, ObjectToBeMapped)
				continue
			}
			mappedObjData, err = MergeInterfaces(mappedObjData, ObjectToBeMapped)
			if err != nil {
				return nil, err
			}
		}
		mappedArrData = append(mappedArrData, mappedObjData)
		return mappedArrData, nil
	}
	return mappedArrData, nil
}
func MapInterfaceToKeyValueTypeWithMapper(mapperFieldObj model.MapperField, responseObj interface{}, mapperFieldsArrToMap []model.MapperField) (map[string]interface{}, error) {

	var mappedObjData = make(map[string]interface{})
	if responseObj != nil {
		mappedObjData[mapperFieldObj.Output] = responseObj
		return mappedObjData, nil
	}
	defaultValueToMap := mapperFieldObj.DefaultValue
	if defaultValueToMap != false {
		mappedObjData[mapperFieldObj.Output] = defaultValueToMap
	}
	return mappedObjData, nil
}

func MergeInterfaces(interfaceToMerge map[string]interface{}, newInterface map[string]interface{}) (map[string]interface{}, error) {
	for k := range newInterface {
		indexKey := k
		interfaceToMerge[indexKey] = newInterface[indexKey]
	}
	return interfaceToMerge, nil
}
func MergeInterfacesIntoArray(interfaceToMerge map[string]interface{}, newInterface map[string]interface{}) (map[string]interface{}, error) {
	for k := range newInterface {
		indexKey := k
		if interfaceToMerge[indexKey] != nil {
			array, ok := interfaceToMerge[indexKey].([]map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unable to cast interfaceToMerge[%s] to []map[string]interface{}", indexKey)
			}
			newArray, ok := newInterface[indexKey].([]map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unable to cast interfaceToMerge[%s] to []map[string]interface{}", indexKey)
			}
			array = append(array, newArray...)
			interfaceToMerge[indexKey] = array
			continue
		}
		interfaceToMerge[indexKey] = newInterface[indexKey]
	}
	return interfaceToMerge, nil
}
func MergeObjectOrArrIntoInterfaces(interfaceToMerge map[string]interface{}, key string, newInterface interface{}) (map[string]interface{}, error) {
	interfaceToMerge[key] = newInterface
	if key == "" {
		interfaceToMerge = newInterface.(map[string]interface{})
	}
	return interfaceToMerge, nil
}
func fetchFieldObjFromFieldsArr(fieldName string, fieldsArr []model.MapperField) (model.MapperField, error) {
	var fieldObj model.MapperField
	for i := 0; i < len(fieldsArr); i++ {
		fieldObj = fieldsArr[i]
		if fieldObj.Id == fieldName {
			return fieldObj, nil
		}
	}
	return fieldObj, nil
}

// ======================================call mapper helper functiion=====================================================================================================================================
func MapperHelperFunction(mapperFieldObj model.MapperField, mappedObjData map[string]interface{}, featureSessionVariables []model.KeyValuePair) (interface{}, error) {
	helpers := mapperFieldObj.HelperFunction
	sessionVariables := utils.ConvertObjectToKeyValuePair(mappedObjData)
	featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, sessionVariables)
	helperFuncOutput, err := utils.CallFunctionFromProps(helpers, featureSessionVariables)
	if err != nil {
		return nil, err
	}
	fmt.Println("=============> Function return Value from Props <============= ", helperFuncOutput)
	return helperFuncOutput, nil
}
func MappedResponseCompletionNorms(completionNorms []model.KeyValuePair, MappedResponse model.KeyValuePair, featureSessionVariables []model.KeyValuePair) (interface{}, error) {
	for _, norms := range completionNorms {
		if norms.Key == "mapped_response" && norms.Value == MappedResponse.Key {
			newProps := make([]model.Props, 0)
			newProps = append(newProps, norms.Props...)
			for index, props := range norms.Props {
				newProps[index].Params = append(props.Params, model.KeyValuePair{
					Key:   "mappedResponse",
					Value: MappedResponse.Value,
					Type:  "static",
				})
			}
			MappedResponse, err := utils.CallFunctionFromProps(newProps, featureSessionVariables)
			if err != nil {
				fmt.Println("--------> Error while parsing mappedResponseCompletionNorms <----------------")
				return nil, err
			}
			return MappedResponse, nil
		}
	}
	return MappedResponse.Value, nil
}
