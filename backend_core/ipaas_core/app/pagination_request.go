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
func ExecutePaginationRequest(taskResponseArray []model.TaskEndpointResponse, currentTaskResponseObject model.TaskEndpointResponse, singleTaskObject map[string]interface{}, endpointTaskFile map[string]interface{}, featureSessionVariables []model.KeyValuePair) model.TaskEndpointResponse {
	var requestBody interface{}
	var params []model.KeyValuePair
	var response model.EndpointResponse
	var requestMethod string
	var mappedResponseCompletionNorms []model.KeyValuePair
	var headers []model.KeyValuePair
	var requestDetails []model.KeyValuePair
	Errors := map[string]interface{}{}

	tasksListObjectFromFeature, err := utils.GetValueFromSessionVariablesKey("tasksListObjectFromFeature", featureSessionVariables)
	if err != nil {
		Errors["message"] = "error from GetValueFromSessionVariablesKey function"
		Errors["task_list_object_from_feature"] = fmt.Sprintf("%v", err)
		currentTaskResponseObject.Errors = Errors
		currentTaskResponseObject.Completed = false
		return currentTaskResponseObject
	}

	//---------------------payload details-[ header-[], body-[],query_params-[], oauth-[] ]--------------------------------------
	payload, ok := endpointTaskFile["payload"].(map[string]interface{})
	if !ok {
		fmt.Println("--------> error in paginated_request payload <----------------")
		Errors["message"] = "error in paginated_request payload"
		currentTaskResponseObject.Errors = Errors
		currentTaskResponseObject.Completed = false
		return currentTaskResponseObject
	}

	//-----------------------------end point url---------------------------------------------------------------------------------
	getURL := endpointTaskFile["scheme"].(string) + "://" + endpointTaskFile["endpoint"].(string)
	getURL, err = utils.FormatEndpoint(getURL, featureSessionVariables)
	if err != nil {
		Errors["message"] = "error from format endpoint function"
		Errors["format_endpoint_error"] = fmt.Sprintf("%v", err)
		currentTaskResponseObject.Errors = Errors
		currentTaskResponseObject.Completed = false
		return currentTaskResponseObject
	}
	//-----------------------------headers---------------------------------------------------------------------------------------
	if payload["headers"] != nil {
		jsonPayloadHeaders, ok := payload["headers"].([]interface{})
		if !ok {
			fmt.Println("--------> error in paginated_request payload[headers] <----------------")
			Errors["message"] = "error in paginated_request payload[headers]"
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}
		headers = utils.ParseObjectsFromConfigPayload(jsonPayloadHeaders)
	}
	fmt.Println("--------> passing headers <----------------" + currentTaskResponseObject.Name)

	//------------------------------request method-----------------------------------------------------------------------------
	if endpointTaskFile["request_configuration"] != nil {
		jsonRequestConfiguration, ok := endpointTaskFile["request_configuration"].(map[string]interface{})
		if !ok {
			fmt.Println("--------> error in paginated_request request_configuration <----------------")
			Errors["message"] = "error in paginated_request request_configuration"
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}
		requestMethod, ok = jsonRequestConfiguration["method"].(string)
		if !ok {
			fmt.Println("--------> error in paginated_request request_configuration[method] <----------------")
			Errors["message"] = "error in paginated_request request_configuration[method]"
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}
		jsonRequestDetails, ok := jsonRequestConfiguration["request_details"].([]interface{})
		if !ok {
			fmt.Println("--------> error in batch_request request_configuration[request_details] <----------------")
			Errors["message"] = "error in batch_request request_configuration[request_details]"
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}
		requestDetails = utils.ParseObjectsFromConfigPayload(jsonRequestDetails)
	}
	//-----------------------------request body----------------------------------------------------------------------------------
	if payload["body"] != nil {
		jsonPayloadBody, ok := payload["body"].([]interface{})
		if !ok {
			fmt.Println("--------> error in single_request payload[body] <----------------")
			Errors["message"] = "error in single_request payload[body]"
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}
		if len(jsonPayloadBody) > 0 {
			bodyKeyValuePair := utils.ParseObjectsFromConfigPayload(jsonPayloadBody)
			bodyKeyValuePair = utils.ParseKeyValuePair(bodyKeyValuePair, featureSessionVariables)
			requestBody = utils.ConvertKeyValuePairToInterface(bodyKeyValuePair)
		}
	}

	for {
		var urlWithParams string
		var queryParams string
		var newHeaders []model.KeyValuePair

		if len(requestDetails) > 0 {
			requestDetailsRespose := utils.ParseKeyValuePair(requestDetails, featureSessionVariables)
			featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, requestDetailsRespose)
		}

		//-----------------------------query params-------------------------------------------------------------------------------
		if payload["query_params"] != nil {
			jsonQueryParams, ok := payload["query_params"].([]interface{})
			if !ok {
				fmt.Println("--------> error in single_request query_params <----------------")
				Errors["message"] = "error in single_request query_params"
				currentTaskResponseObject.Errors = Errors
				currentTaskResponseObject.Completed = false
				return currentTaskResponseObject
			}
			if len(jsonQueryParams) > 0 {
				params = utils.ParseObjectsFromConfigPayload(jsonQueryParams)
				params = utils.ParseKeyValuePair(params, featureSessionVariables)
				queryParams, err = utils.GetUrlQueryParams(params, featureSessionVariables)
				if err != nil {
					fmt.Println("--------> error in api call function <----------------")
					Errors["message"] = "error in api call function"
					Errors["api_error"] = fmt.Sprintf("%v", err)
					currentTaskResponseObject.Errors = Errors
					currentTaskResponseObject.Completed = false
					return currentTaskResponseObject
				}
			}
		}
		urlWithParams = getURL + queryParams
		fmt.Println("--------> api endpoint <----------------" + urlWithParams)

		//-----------------------------headers-------------------------------------------------------------------------------
		newHeaders = append(newHeaders, headers...)

		//-------------call the MakeAPIRequest function and get the response of the API-------------------------------------------
		response, err = utils.MakeAPIRequest(requestMethod, urlWithParams, newHeaders, requestBody, featureSessionVariables, endpointTaskFile)
		if err != nil {
			fmt.Println("--------> error in api call function <----------------")
			Errors["message"] = "error in api call function"
			Errors["api_error"] = fmt.Sprintf("%v", err)
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}

		//-----------------customize the API response------------------------------------------------
		if endpointTaskFile["response"] != nil {
			jsonResponse, ok := endpointTaskFile["response"].(map[string]interface{})
			if !ok {
				fmt.Println("--------> error in paginated_request response <----------------")
				Errors["message"] = "error in paginated_request response"
				currentTaskResponseObject.Errors = Errors
				currentTaskResponseObject.Completed = false
				return currentTaskResponseObject
			}

			//------------------------response_completion_norms--------------------------------------------------------------------------
			if jsonResponse["response_completion_norms"] != nil {
				responseCompletionNorms, ok := jsonResponse["response_completion_norms"].([]interface{})
				if !ok {
					fmt.Println("--------> error in single_request response[response_completion_norms] <----------------")
					Errors["message"] = "error in single_request response[response_completion_norms]"
					currentTaskResponseObject.Errors = Errors
					currentTaskResponseObject.Completed = false
					return currentTaskResponseObject
				}
				mappedResponseCompletionNorms = utils.ParseObjectsFromConfigPayload(responseCompletionNorms)
			}
			//------------------------1.session variables--------------------------------------------------------------------------
			if jsonResponse["session_variables"] != nil {
				jsonResponseSessionVariables, ok := jsonResponse["session_variables"].([]interface{})
				if !ok {
					fmt.Println("--------> error in paginated_request response[session_variables] <----------------")
					Errors["message"] = "error in paginated_request response[session_variables]"
					currentTaskResponseObject.Errors = Errors
					currentTaskResponseObject.Completed = false
					return currentTaskResponseObject
				}
				taskSessionVariables, err := utils.ConvertArrayInterfaceToArrayStructWithJsonKeyParse(jsonResponseSessionVariables, response, featureSessionVariables)
				if err != nil {
					fmt.Println("--------> error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function <----------------", err)
					Errors["message"] = "error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function"
					Errors["api_error"] = fmt.Sprintf("%v", err)
					currentTaskResponseObject.Errors = Errors
					currentTaskResponseObject.Completed = false
					return currentTaskResponseObject
				}
				taskSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, taskSessionVariables)
				currentTaskResponseObject.SessionVariables = taskSessionVariables
			}
		}
		currentTaskResponseObject.EndpointResponse = response
		//-------------add or update the endpoint task sessionVariables to FeatureSession variables ---------------------------
		featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, currentTaskResponseObject.SessionVariables)

		//need to check mapper is exists or not ?
		//------------if the task contains any mapper then execute the mapper function ----------------------------------------
		if singleTaskObject["mapper"] != nil {
			mapper, ok := singleTaskObject["mapper"].([]interface{})
			if !ok {
				fmt.Println("--------> error in mapper <----------------")
				Errors["message"] = "error in mapper json file"
				currentTaskResponseObject.Errors = Errors
				currentTaskResponseObject.Completed = false
			}
			if currentTaskResponseObject.Body != nil && len(mapper) > 0 {
				currentTaskResponseObject.MappedResponse, err = repo.ExecuteMapperForTask(mapper, currentTaskResponseObject.Body, mappedResponseCompletionNorms, featureSessionVariables)
				if err != nil {
					fmt.Println("--------> error in mapper function <----------------")
					Errors["message"] = "error in mapper function"
					Errors["mapper_error"] = fmt.Sprintf("%v", err)
					currentTaskResponseObject.Errors = Errors
					currentTaskResponseObject.Completed = false
					return currentTaskResponseObject
				}
			}
		}

		//-------------------------------completion norms to check for next API call----------------------------------------------
		completionNorms := utils.ParsePropsFromConfig(endpointTaskFile["completion_norms"].([]interface{}))
		for index, props := range completionNorms {
			completionNorms[index].Params = append(props.Params, model.KeyValuePair{
				Key:   "response",
				Value: response,
				Type:  "static",
			})
		}
		completionNormsStatus, err := utils.CallFunctionFromProps(completionNorms, featureSessionVariables)
		if err != nil {
			fmt.Println("--------> error in completion norms function <----------------")
			Errors["message"] = "error in completion norms function"
			Errors["completion_norms_error"] = fmt.Sprintf("%v", err)
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}
		fmt.Println("--------> Function return Value from Props - completionNormsStatus<----------------", completionNormsStatus)
		if completionNormsStatus.(bool) {
			break
		}

		//-------------------------update the params for next API call------------------------------------------------------------
		// featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, params)

		//------------------execute the on_success_tasks sequentially-------------------------------------------------------------
		if singleTaskObject["on_success"] != nil {
			var onSuccesTasks []map[string]interface{}

			//-------loop for on_success_tasks -----------------------------------------------------------------------------------
			jsonOnSuccessTasks, ok := singleTaskObject["on_success"].([]interface{})
			if !ok {
				fmt.Println("--------> error in paginated_request on_success_tasks <----------------")
				Errors["message"] = "error in paginated_request on_success_tasks"
				currentTaskResponseObject.Errors = Errors
				currentTaskResponseObject.Completed = false
				return currentTaskResponseObject
			}
			for _, taskName := range jsonOnSuccessTasks {
				//----loop for tasksListObjectFromFeature-------------------------------------------------------------------------
				jsonTasksListObjectFromFeature, ok := tasksListObjectFromFeature.([]map[string]interface{})
				if !ok {
					fmt.Println("--------> error in paginated_request on_success_tasks <----------------")
					Errors["message"] = "error in paginated_request on_success_tasks"
					currentTaskResponseObject.Errors = Errors
					currentTaskResponseObject.Completed = false
					return currentTaskResponseObject
				}
				for _, task := range jsonTasksListObjectFromFeature {
					if task["name"].(string) == taskName.(string) {
						onSuccesTasks = append(onSuccesTasks, task)
					}
				}
			}

			//-----------------------append the current endpoint task response -----------------------------------------------------------
			taskResponseArray = utils.AppendOrUpdateTaskResponseObject(taskResponseArray, currentTaskResponseObject)

			//===================recursively process the ExecuteFeatureTasks======================================================
			ExecuteFeatureTasksResponse := ExecuteFeatureTasks(onSuccesTasks, featureSessionVariables, taskResponseArray)
			ErrorsArray := make([]map[string]interface{}, 0)
			for _, taskResponse := range ExecuteFeatureTasksResponse {
				if taskResponse.Errors != nil {
					NewError := make(map[string]interface{}, 0)
					NewError["name"] = taskResponse.Name
					NewError["errors"] = taskResponse.Errors
					NewError["completed"] = taskResponse.Completed
					ErrorsArray = append(ErrorsArray, NewError)
				}
			}
			if len(ErrorsArray) > 0 {
				Errors["message"] = "on success task failures"
				Errors["on_success_tasks_errors"] = ErrorsArray
				currentTaskResponseObject.Errors = Errors
			}
		}
	}

	//------------------task response ------------------------------------------------------------------------------------------
	currentTaskResponseObject.Completed = true
	return currentTaskResponseObject
}
