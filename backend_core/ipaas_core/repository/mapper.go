package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fermion/backend_core/ipaas_core/model"
	"fermion/backend_core/ipaas_core/utils"
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

type MapperTask struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	InputType string `json:"input_type"`
	InputKey  string `json:"input_key"`
	MapperId  string `json:"mapper_id"`
}

func ExecuteMapperForTask(mapperTasks []interface{}, response map[string]interface{}, mappedResponseCompletionNorms, featureSessionVariables []model.KeyValuePair) ([]model.KeyValuePair, error) {

	var mapperTasksResponse []model.KeyValuePair

	//get base path from session variable
	// basePath, err := utils.GetValueFromSessionVariablesKey("basePath", featureSessionVariables)
	// if err != nil {
	// 	return nil, err
	// }

	var mapperTasksArr []MapperTask
	mdata, err := json.Marshal(mapperTasks)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(mdata, &mapperTasksArr)
	if err != nil {
		return nil, err
	}

	for _, task := range mapperTasksArr {

		//prepare path for reading mapper file
		// mapperPath := basePath.(string) + task.Link
		mapperConfigTemplateFile, err := utils.ReadFeature("mappers", task.MapperId)
		if err != nil {
			return nil, err
		}

		inputPath, err := utils.FormatEndpoint(task.InputKey, featureSessionVariables)
		if err != nil {
			return nil, err
		}

		//type array
		if task.InputType == "array" {
			responseArr, err := utils.ParseJsonPathFromObject(response, inputPath)
			if err != nil {
				return nil, err
			}
			var mapperInput []interface{}
			mdata, err := json.Marshal(responseArr)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(mdata, &mapperInput)
			if err != nil {
				return nil, err
			}

			mappedResponse, err := ExecuteArrayMapper(mapperInput, mapperConfigTemplateFile, featureSessionVariables)
			if err != nil {
				return nil, err
			}

			response := model.KeyValuePair{
				Key:   task.Name,
				Value: mappedResponse,
			}

			for _, norms := range mappedResponseCompletionNorms {
				if norms.Key == "mapped_response" && norms.Value == task.Name {
					for index, props := range norms.Props {
						norms.Props[index].Params = append(props.Params, model.KeyValuePair{
							Key:   "mappedResponse",
							Value: mappedResponse,
							Type:  "static",
						})
					}
					response.Value, err = utils.CallFunctionFromProps(norms.Props, featureSessionVariables)
					if err != nil {
						fmt.Println("--------> Error while parsing mappedResponseCompletionNorms <----------------")
						return nil, err
					}
					break
				}
			}
			mapperTasksResponse = append(mapperTasksResponse, response)
			continue
		}

		//type object
		if task.InputType == "object" {
			mapperInput, err := utils.ParseJsonPathFromObject(response, inputPath)
			if err != nil {
				return nil, err
			}
			mappedResponse, err := ExecuteObjMapper(mapperInput, mapperConfigTemplateFile, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			response := model.KeyValuePair{
				Key:   task.Name,
				Value: mappedResponse,
			}
			for _, norms := range mappedResponseCompletionNorms {
				if norms.Key == "mapped_response" && norms.Value == task.Name {
					for index, props := range norms.Props {
						norms.Props[index].Params = append(props.Params, model.KeyValuePair{
							Key:   "mappedResponse",
							Value: mappedResponse,
							Type:  "static",
						})
					}
					response.Value, err = utils.CallFunctionFromProps(norms.Props, featureSessionVariables)

					if err != nil {
						fmt.Println("--------> Error while parsing mappedResponseCompletionNorms <----------------")
						return nil, err
					}
					break
				}
			}
			mapperTasksResponse = append(mapperTasksResponse, response)
		}
	}
	return mapperTasksResponse, nil
}
func MergeMultipleMapperPayload(required_mapped_output_array []interface{}, taskResponseArr []model.TaskEndpointResponse) interface{} {
	var payload interface{}
	for _, required_map := range required_mapped_output_array {
		if required_map.(map[string]interface{}) != nil && required_map.(map[string]interface{})["type"].(string) == "payload" {
			for _, taskResponse := range taskResponseArr {
				if taskResponse.Name == required_map.(map[string]interface{})["name"].(string) {
					for _, mapper := range taskResponse.MappedResponse {
						fmt.Println(mapper.Key, required_map.(map[string]interface{})["mapper_name"].(string))
						if mapper.Key == required_map.(map[string]interface{})["mapper_name"].(string) {
							// fmt.Println("===============Required Map=============================")
							// mdata, _ := json.MarshalIndent(taskResponse.MappedResponse, "", "\t")
							// fmt.Println(string(mdata))
							// fmt.Println("============================================")
							return mapper.Value
						}
					}
				}
			}
		}
	}
	return payload
}

func ExecuteArrayMapper(responseArrToMap []interface{}, mapperConfigTemplateFile map[string]interface{}, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	mappedArrData := make([]map[string]interface{}, 0)
	mapperFieldsArrToMap := mapperConfigTemplateFile["fields"].([]interface{})
	for responseObjIndex := 0; responseObjIndex < len(responseArrToMap); responseObjIndex++ {
		var mappedObjData = make(map[string]interface{})
		responseObj := responseArrToMap[responseObjIndex].(map[string]interface{})
		for mapperFieldObjIndex := 0; mapperFieldObjIndex < len(mapperFieldsArrToMap); mapperFieldObjIndex++ {
			mapperFieldObj := mapperFieldsArrToMap[mapperFieldObjIndex].(map[string]interface{})
			if mapperFieldObj["execution_type"] == "nested" {
				continue
			}
			mappedResponseObj, err := checkFieldTypeAndMap(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			if mapperFieldObj["output_type"] == "array" {
				MergeInterfacesIntoArray(mappedObjData, mappedResponseObj)
				continue
			}
			MergeInterfaces(mappedObjData, mappedResponseObj)
		}
		mappedArrData = append(mappedArrData, mappedObjData)
	}
	return mappedArrData, nil
}
func ExecuteObjMapper(responseObjToMapper interface{}, mapperConfigTemplateFile map[string]interface{}, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	responseObj := responseObjToMapper.(map[string]interface{})
	mapperFieldsArrToMap := mapperConfigTemplateFile["fields"].([]interface{})

	var mappedObjData = make(map[string]interface{})

	for mapperFieldObjIndex := 0; mapperFieldObjIndex < len(mapperFieldsArrToMap); mapperFieldObjIndex++ {
		mapperFieldObj := mapperFieldsArrToMap[mapperFieldObjIndex].(map[string]interface{})
		if mapperFieldObj["execution_type"] == "nested" {
			continue
		}
		mappedResponseObj, err := checkFieldTypeAndMap(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		if mapperFieldObj["output_type"] == "array" {
			MergeInterfacesIntoArray(mappedObjData, mappedResponseObj)
			continue
		}
		MergeInterfaces(mappedObjData, mappedResponseObj)
	}
	return mappedObjData, nil
}

func checkFieldTypeAndMap(mapperFieldObj map[string]interface{}, responseObj interface{}, mapperFieldsArrToMap []interface{}, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	if mapperFieldObj["input_type"] == "keyvalue" && mapperFieldObj["output_type"] == "keyvalue" {
		mappedResForKey, err := MapKeyValueToKeyValueTypeWithMapper(mapperFieldObj, responseObj, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForKey)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "object" && mapperFieldObj["output_type"] == "object" {
		mappedResForObject, err := MapObjectToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj["output"].(string), mappedResForObject)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "array" && mapperFieldObj["output_type"] == "array" {
		mappedResForArr, err := MapArrToArrTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj["output"].(string), mappedResForArr)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "array" && mapperFieldObj["output_type"] == "object" {
		mappedResForObject, err := MapArrToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj["output"].(string), mappedResForObject)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "object" && mapperFieldObj["output_type"] == "keyvalue" {
		mappedResForKeyValue, err := MapObjectToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForKeyValue)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "keyvalue" && mapperFieldObj["output_type"] == "object" {
		mappedResForObject, err := MapKeyValueToObjectTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForObject)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "keyvalue" && mapperFieldObj["output_type"] == "array" {
		mappedResForArray, err := MapKeyValueToArrayTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj["output"].(string), mappedResForArray)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "object" && mapperFieldObj["output_type"] == "array" {
		mappedResForArray, err := MapObjectToArrayTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj["output"].(string), mappedResForArray)
		if err != nil {
			return nil, err
		}
	}
	if mapperFieldObj["input_type"] == "interface" && mapperFieldObj["output_type"] == "keyvalue" {
		mappedResForKey, err := MapInterfaceToKeyValueTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap)
		if err != nil {
			return nil, err
		}
		mappedObjData, err = MergeInterfaces(mappedObjData, mappedResForKey)
		if err != nil {
			return nil, err
		}
	}
	// if mapperFieldObj["input_type"] == "keyvalue" && mapperFieldObj["output_type"] == "interface" {
	// 	mappedResForInterface := MapKeyvalueToInterfaceTypeWithMapper(mapperFieldObj, responseObj, mapperFieldsArrToMap)
	// 	mappedObjData = MergeObjectOrArrIntoInterfaces(mappedObjData, mapperFieldObj["output"].(string), mappedResForInterface)
	// }
	return mappedObjData, nil
}

func MapKeyValueToKeyValueTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}
	valueObj := responseObj[mapperFieldObj["input"].(string)]
	if valueObj != nil {
		mappedObjData[mapperFieldObj["output"].(string)] = responseObj[mapperFieldObj["input"].(string)]
		//----helper function ----------
		if len(mapperFieldObj["helper_function"].([]interface{})) > 0 {
			helperFunctionOutput, err := MapperHelperFunction(mapperFieldObj, mappedObjData, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			mappedObjData[mapperFieldObj["output"].(string)] = helperFunctionOutput
		}
		return mappedObjData, nil
	}
	defaultValueToMap := mapperFieldObj["default_value"]
	if defaultValueToMap != false {
		mappedObjData[mapperFieldObj["output"].(string)] = defaultValueToMap
	}
	//----helper function ----------
	if len(mapperFieldObj["helper_function"].([]interface{})) > 0 {
		helperFunctionOutput, err := MapperHelperFunction(mapperFieldObj, mappedObjData, featureSessionVariables)
		if err != nil {
			return nil, err
		}
		mappedObjData[mapperFieldObj["output"].(string)] = helperFunctionOutput
		return mappedObjData, nil
	}
	return mappedObjData, nil
}
func MapObjectToObjectTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, mapperFieldsArrToMap []interface{}, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}
	valueObj := responseObj[mapperFieldObj["input"].(string)]
	if valueObj == nil {
		valueObj = responseObj
	}
	if valueObj != nil {
		// fmt.Println(valueObj)
		mappedToBe := valueObj.(map[string]interface{})
		fieldsToMapWithinObj := mapperFieldObj["fields"].([]interface{})
		for fieldToMapIndex := 0; fieldToMapIndex < len(fieldsToMapWithinObj); fieldToMapIndex++ {
			fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[fieldToMapIndex].(string), mapperFieldsArrToMap)
			if err != nil {
				return nil, err
			}
			ObjectToBeMapped, err := checkFieldTypeAndMap(fieldObj, mappedToBe, mapperFieldsArrToMap, featureSessionVariables)
			if err != nil {
				return nil, err
			}
			if fieldObj["output_type"] == "array" {
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
func MapArrToArrTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, mapperFieldsArrToMap []interface{}, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	var mappedArrData = make([]map[string]interface{}, 0)
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}
	if responseObj[mapperFieldObj["input"].(string)] == nil {
		return mappedArrData, nil
	}
	responseArr := responseObj[mapperFieldObj["input"].(string)].([]interface{})
	for responseObjIndex := 0; responseObjIndex < len(responseArr); responseObjIndex++ {
		responseObj := responseArr[responseObjIndex]
		var mappedObjData = make(map[string]interface{})
		if responseObj != nil {
			objToMapWithinArr := mapperFieldObj["fields"].([]interface{})
			for fieldToMapIndex := 0; fieldToMapIndex < len(objToMapWithinArr); fieldToMapIndex++ {
				fieldObj, err := fetchFieldObjFromFieldsArr(objToMapWithinArr[fieldToMapIndex].(string), mapperFieldsArrToMap)
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
func MapArrToObjectTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, mapperFieldsArrToMap []interface{}, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	keyForInputWithIndex := mapperFieldObj["input"].(string)
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
		objToMapWithinArr := mapperFieldObj["fields"].([]interface{})
		for fieldToMapIndex := 0; fieldToMapIndex < len(objToMapWithinArr); fieldToMapIndex++ {
			fieldObj, err := fetchFieldObjFromFieldsArr(objToMapWithinArr[fieldToMapIndex].(string), mapperFieldsArrToMap)
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
func MapKeyValueToObjectTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, mapperFieldsArrToMap []interface{}, featureSessionVariables []model.KeyValuePair) (map[string]interface{}, error) {
	var mappedObjData = make(map[string]interface{})
	var valueObj = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}

	valueObj[mapperFieldObj["input"].(string)] = responseObj[mapperFieldObj["input"].(string)]
	if len(valueObj) == 0 {
		valueObj = responseObj
	}
	fieldsToMapWithinObj := mapperFieldObj["fields"].([]interface{})
	fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[0].(string), mapperFieldsArrToMap)
	if err != nil {
		return nil, err
	}

	mappedObjData[mapperFieldObj["output"].(string)], err = checkFieldTypeAndMap(fieldObj, valueObj, mapperFieldsArrToMap, featureSessionVariables)
	if err != nil {
		return nil, err
	}
	return mappedObjData, nil
}
func MapKeyValueToArrayTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, mapperFieldsArrToMap []interface{}, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	var mappedObjData = make([]map[string]interface{}, 0)
	var valueObj = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}

	valueObj[mapperFieldObj["input"].(string)] = responseObj[mapperFieldObj["input"].(string)]
	if len(valueObj) == 0 {
		valueObj = responseObj
	}
	fieldsToMapWithinObj := mapperFieldObj["fields"].([]interface{})
	fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[0].(string), mapperFieldsArrToMap)
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
func MapObjectToArrayTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, mapperFieldsArrToMap []interface{}, featureSessionVariables []model.KeyValuePair) ([]map[string]interface{}, error) {
	var mappedArrData = make([]map[string]interface{}, 0)
	var mappedObjData = make(map[string]interface{})
	responseObj, ok := responseInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("response interface is not a map[string]interface{}")
	}

	valueObj := responseObj[mapperFieldObj["input"].(string)]
	if valueObj == nil {
		valueObj = responseObj
	}
	if valueObj != nil {
		fieldsToMapWithinObj := mapperFieldObj["fields"].([]interface{})
		for fieldToMapIndex := 0; fieldToMapIndex < len(fieldsToMapWithinObj); fieldToMapIndex++ {
			fieldObj, err := fetchFieldObjFromFieldsArr(fieldsToMapWithinObj[fieldToMapIndex].(string), mapperFieldsArrToMap)
			if err != nil {
				return nil, err
			}
			ObjectToBeMapped, err := checkFieldTypeAndMap(fieldObj, valueObj.(map[string]interface{}), mapperFieldsArrToMap, featureSessionVariables)
			if err != nil {
				return nil, err
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
func MapInterfaceToKeyValueTypeWithMapper(mapperFieldObj map[string]interface{}, responseObj interface{}, mapperFieldsArrToMap []interface{}) (map[string]interface{}, error) {

	var mappedObjData = make(map[string]interface{})
	if responseObj != nil {
		mappedObjData[mapperFieldObj["output"].(string)] = responseObj
		return mappedObjData, nil
	}
	defaultValueToMap := mapperFieldObj["default_value"]
	if defaultValueToMap != false {
		mappedObjData[mapperFieldObj["output"].(string)] = defaultValueToMap
	}
	return mappedObjData, nil
}

// func MapKeyvalueToInterfaceTypeWithMapper(mapperFieldObj map[string]interface{}, responseInterface interface{}, mapperFieldsArrToMap []interface{}) interface{} {
// 	responseObj, error := responseInterface.(map[string]interface{})
// 	if error {
// 		valueObj := responseObj[mapperFieldObj["input"].(string)]
// 		if valueObj != nil {
// 			return valueObj
// 		}
// 		defaultValueToMap := mapperFieldObj["default_value"]
// 		if defaultValueToMap != false {
// 			return defaultValueToMap
// 		}
// 		return nil
// 	} else {
// 		return nil
// 	}
// }

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
func fetchFieldObjFromFieldsArr(fieldName string, fieldsArr []interface{}) (map[string]interface{}, error) {
	for i := 0; i < len(fieldsArr); i++ {
		fieldObj, ok := fieldsArr[i].(map[string]interface{})
		if !ok {
			return nil, errors.New(" error in fetchFieldObjFromFieldsArr")
		}
		if fieldObj["id"] == fieldName {
			return fieldObj, nil
		}
	}
	return nil, nil
}

// ------------------------------call mapper helper functiion-----------------------------------------------
func MapperHelperFunction(mapperFieldObj map[string]interface{}, mappedObjData map[string]interface{}, featureSessionVariables []model.KeyValuePair) (interface{}, error) {
	helpers := utils.ParsePropsFromConfig(mapperFieldObj["helper_function"].([]interface{}))
	sessionVariables := utils.ConvertObjectToKeyValuePair(mappedObjData)
	featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, sessionVariables)
	helperFuncOutput, err := utils.CallFunctionFromProps(helpers, featureSessionVariables)
	if err != nil {
		return nil, err
	}
	fmt.Println("--------> Function return Value from Props<----------------", helperFuncOutput)
	return helperFuncOutput, nil
}
