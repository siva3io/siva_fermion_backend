package ipaas_core

import (
	"fmt"

	model "fermion/backend_core/ipaas_core/model"
	repo "fermion/backend_core/ipaas_core/repository"
	"fermion/backend_core/ipaas_core/utils"
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
func ExecuteSingleRequest(taskResponseArray []model.TaskEndpointResponse, taskResponseObject model.TaskEndpointResponse, singleTaskObject map[string]interface{}, endpointTaskFile map[string]interface{}, featureSessionVariables []model.KeyValuePair) model.TaskEndpointResponse {
	var requestBody interface{}
	var queryParams string
	var urlWithParams string
	var mappedResponseCompletionNorms []model.KeyValuePair
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
	getURL, err := utils.FormatEndpoint(getURL, featureSessionVariables)
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
			params := utils.ParseObjectsFromConfigPayload(jsonQueryParams)
			params = utils.ParseKeyValuePair(params, featureSessionVariables)
			queryParams, err = utils.GetUrlQueryParams(params, featureSessionVariables)
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

	//-------if the task requires any mapper output then ftech the mapper output and append to the request_body ------------
	if singleTaskObject["required_mapped_output"] != nil {
		jsonRequiredMappedOutput, ok := singleTaskObject["required_mapped_output"].([]interface{})
		if !ok {
			fmt.Println("--------> error in single_request required_mapped_output <----------------")
			Errors["message"] = "error in single_request required_mapped_output"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		requestBody = repo.MergeMultipleMapperPayload(jsonRequiredMappedOutput, taskResponseArray)
	}
	// fmt.Println("------------<<>>", requestBody)

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
			bodyKeyValuePair := utils.ParseObjectsFromConfigPayload(jsonPayloadBody)
			bodyKeyValuePair = utils.ParseKeyValuePair(bodyKeyValuePair, featureSessionVariables)
			requestBody = utils.ConvertKeyValuePairToInterface(bodyKeyValuePair)
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
		headers = utils.ParseObjectsFromConfigPayload(jsonPayloadHeaders)
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
	response, err := utils.MakeAPIRequest(requestMethod, urlWithParams, headers, requestBody, featureSessionVariables, endpointTaskFile)
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
		//------------------------response_completion_norms--------------------------------------------------------------------------
		if jsonResponse["response_completion_norms"] != nil {
			responseCompletionNorms, ok := jsonResponse["response_completion_norms"].([]interface{})
			if !ok {
				fmt.Println("--------> error in single_request response[response_completion_norms] <----------------")
				Errors["message"] = "error in single_request response[response_completion_norms]"
				taskResponseObject.Errors = Errors
				taskResponseObject.Completed = false
				return taskResponseObject
			}
			mappedResponseCompletionNorms = utils.ParseObjectsFromConfigPayload(responseCompletionNorms)
		}
		//------------------------session variables--------------------------------------------------------------------------
		if jsonResponse["session_variables"] != nil {
			jsonResponseSessionVariables, ok := jsonResponse["session_variables"].([]interface{})
			if !ok {
				fmt.Println("--------> error in single_request response[session_variables] <----------------")
				Errors["message"] = "error in single_request response[session_variables]"
				taskResponseObject.Errors = Errors
				taskResponseObject.Completed = false
				return taskResponseObject
			}
			taskSessionVariables, err := utils.ConvertArrayInterfaceToArrayStructWithJsonKeyParse(jsonResponseSessionVariables, response, featureSessionVariables)
			if err != nil {
				fmt.Println("--------> error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function <----------------")
				Errors["message"] = "error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function"
				Errors["api_error"] = fmt.Sprintf("%v", err)
				taskResponseObject.Errors = Errors
				taskResponseObject.Completed = false
				return taskResponseObject
			}
			taskSessionVariables = utils.AppendOrUpdateKeyValuePair(taskSessionVariables, featureSessionVariables)
			taskResponseObject.SessionVariables = taskSessionVariables
		}
	}
	taskResponseObject.EndpointResponse = response
	if singleTaskObject["mapper"] != nil {
		mapper, ok := singleTaskObject["mapper"].([]interface{})
		if !ok {
			fmt.Println("--------> error in mapper <----------------")
			Errors["message"] = "error in mapper json file"
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		if taskResponseObject.Body != nil && len(mapper) > 0 {
			res, err := repo.ExecuteMapperForTask(mapper, taskResponseObject.Body, mappedResponseCompletionNorms, featureSessionVariables)
			if err != nil {
				fmt.Println("--------> error in mapper function <----------------")
				Errors["message"] = "error in mapper fuction"
				Errors["mapper_error"] = fmt.Sprintf("%v", err)
				taskResponseObject.Errors = Errors
				taskResponseObject.Completed = false
				return taskResponseObject
			}
			taskResponseObject.MappedResponse = res
		}
	}

	//------------------task response ------------------------------------------------------------------------------------------
	taskResponseObject.Completed = true

	return taskResponseObject
}
