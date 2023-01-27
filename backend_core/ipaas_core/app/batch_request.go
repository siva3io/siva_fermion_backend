package ipaas_core

import (
	"fmt"
	"strconv"

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

func ExecuteBatchRequest(taskResponseArray []model.TaskEndpointResponse, currentTaskResponseObject model.TaskEndpointResponse, singleTaskObject model.Task, endpointTaskFile model.APIFile, featureSessionVariables []model.KeyValuePair) model.TaskEndpointResponse {
	var mappedResponseCompletionNorms []model.KeyValuePair
	Errors := map[string]interface{}{}

	//========================task list object from feature======================================================================
	tasksListObjectFromFeature := utils.GetValueFromSessionVariablesKey("tasksListObjectFromFeature", featureSessionVariables)

	//========================payload details-[ header-[], body-[],query_params-[], oauth-[] ]======================================================================
	payload := endpointTaskFile.Payload
	queryParams := payload.QueryParams

	//======================================end point url===========================================================================================================
	getURL := endpointTaskFile.Scheme + "://" + endpointTaskFile.Endpoint

	//==========================request_method & request_details====================================================================================================
	requestMethod := endpointTaskFile.RequestConfiguration.Method
	requestDetails := endpointTaskFile.RequestConfiguration.RequestDetails

	//=================================total_batch_size=============================================================================================================
	totalBatch := utils.GetValueFromSessionVariablesKey("total_batch", featureSessionVariables)
	totalBatchArray, ok := totalBatch.([]interface{})
	if !ok {
		Errors["message"] = "error for finding totalBatchArray"
		Errors["total_batch"] = "error in totalBatchArray"
		currentTaskResponseObject.Errors = Errors
		currentTaskResponseObject.Completed = false
		return currentTaskResponseObject
	}
	batchSize := len(totalBatchArray)

	//=================================per_batch_count==============================================================================================================
	batchCount := utils.GetValueFromSessionVariablesKey("per_batch", featureSessionVariables)
	perBatchSize, err := strconv.Atoi(batchCount.(string))
	if err != nil {
		Errors["message"] = "error in "
		Errors["batch_count_error"] = fmt.Sprintf("%v", err)
		currentTaskResponseObject.Errors = Errors
		currentTaskResponseObject.Completed = false
		return currentTaskResponseObject
	}

	//============================================ batch_loop =======================================================================================================
	for currentBatch := 0; currentBatch < batchSize; currentBatch += perBatchSize {

		var urlWithParams string
		var newQueryParams string
		var requestBody interface{}
		var headers []model.KeyValuePair

		//=================================update the batch requirements===========================================================================================
		var batchRequestDetails []model.KeyValuePair
		batchRequestDetails = append(batchRequestDetails, requestDetails...)
		featureSessionVariables, err = utils.UpdateBatchRequirements(featureSessionVariables, currentBatch, perBatchSize, batchRequestDetails)
		if err != nil {
			Errors["message"] = "error from UpdateCurrentBatchRequirements function"
			Errors["update_batch_request_requirements_error"] = fmt.Sprintf("%v", err)
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}

		//====================================format the url========================================================================================================
		formattedURL, err := utils.FormatEndpoint(getURL, featureSessionVariables)
		if err != nil {
			Errors["message"] = "error from format endpoint function"
			Errors["format_endpoint_error"] = fmt.Sprintf("%v", err)
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}

		//============================query params===================================================================================================================
		if len(queryParams) > 0 {
			var input []model.KeyValuePair
			input = append(input, queryParams...)
			params := utils.ParseKeyValuePair(input, featureSessionVariables)
			newQueryParams, err = utils.GetUrlQueryParams(params, featureSessionVariables)
			if err != nil {
				fmt.Println("--------> error in api call function <----------------")
				Errors["message"] = "error in api call function"
				Errors["api_error"] = fmt.Sprintf("%v", err)
				currentTaskResponseObject.Errors = Errors
				currentTaskResponseObject.Completed = false
				return currentTaskResponseObject
			}
		}
		urlWithParams = formattedURL + newQueryParams
		fmt.Println("=============> executing query params <============= " + singleTaskObject.TaskName)

		//============================request body===================================================================================================================
		if len(payload.Body) > 0 {
			var input []model.KeyValuePair
			input = append(input, payload.Body...)
			bodyKeyValuePair := utils.ParseKeyValuePair(input, featureSessionVariables)
			requestBody = utils.ConvertKeyValuePairToInterface(bodyKeyValuePair)
		}
		fmt.Println("=============> executing request_body <============= " + singleTaskObject.TaskName)

		//================================headers====================================================================================================================
		headers = append(headers, payload.Headers...)
		fmt.Println("=============> passing headers <============= " + singleTaskObject.TaskName)

		//================== call the MakeAPIRequest function and get the response of the API ========================================================================
		response, err := utils.MakeAPIRequest(requestMethod, urlWithParams, headers, requestBody, featureSessionVariables, endpointTaskFile)
		if err != nil {
			fmt.Println("--------> error in api call function <----------------")
			Errors["message"] = "error in api call function"
			Errors["api_error"] = fmt.Sprintf("%v", err)
			currentTaskResponseObject.Errors = Errors
			currentTaskResponseObject.Completed = false
			return currentTaskResponseObject
		}

		//=========================customize the API response=======================================================================================================
		//===========================session variables==============================================================================================================
		if len(endpointTaskFile.Response.SessionVariables) > 0 {
			taskSessionVariables, err := utils.ConvertArrayInterfaceToArrayStructWithJsonKeyParse(endpointTaskFile.Response.SessionVariables, response, featureSessionVariables)
			if err != nil {
				fmt.Println("=============>error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function <=============", err)
				Errors["message"] = "error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function"
				Errors["api_error"] = fmt.Sprintf("%v", err)
				currentTaskResponseObject.Errors = Errors
				currentTaskResponseObject.Completed = false
				return currentTaskResponseObject
			}
			taskSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, taskSessionVariables)
			currentTaskResponseObject.SessionVariables = append(currentTaskResponseObject.SessionVariables, taskSessionVariables...)
		}
		currentTaskResponseObject.EndpointResponse = response

		//==================if the task contains any mapper then execute the mapper function========================================================================
		if currentTaskResponseObject.Body != nil && len(singleTaskObject.Mapper) > 0 {
			currentTaskResponseObject.MappedResponse, err = repo.ExecuteMapperForTask(singleTaskObject.Mapper, currentTaskResponseObject.Body, mappedResponseCompletionNorms, featureSessionVariables)
			if err != nil {
				fmt.Println("--------> error in mapper function <----------------")
				Errors["message"] = "error in mapper function"
				Errors["mapper_error"] = fmt.Sprintf("%v", err)
				currentTaskResponseObject.Errors = Errors
				currentTaskResponseObject.Completed = false
				return currentTaskResponseObject
			}
		}

		//======================execute the on_success_tasks sequentially=========================================================================================
		if len(singleTaskObject.OnSuccess) > 0 {
			var onSuccesTasks []model.Task

			for _, taskName := range singleTaskObject.OnSuccess {

				//============loop for tasksListObjectFromFeature=================================================================================================
				jsonTasksListObjectFromFeature, ok := tasksListObjectFromFeature.([]model.Task)
				if !ok {
					fmt.Println("=============> error in batch_request on_success_tasks <=============")
					Errors["message"] = "error in batch_request on_success_tasks"
					currentTaskResponseObject.Errors = Errors
					currentTaskResponseObject.Completed = false
					return currentTaskResponseObject
				}
				for _, task := range jsonTasksListObjectFromFeature {
					if task.TaskName == taskName {
						onSuccesTasks = append(onSuccesTasks, task)
					}
				}
			}

			//===========================append or update feature_session_variable============================================================================================================
			featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, currentTaskResponseObject.SessionVariables)

			//==========================append the current endpoint task response================================================================================
			taskResponseArray = utils.AppendOrUpdateTaskResponseObject(taskResponseArray, currentTaskResponseObject)

			//===================recursively process the ExecuteFeatureTasks=====================================================================================
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

	//=================================task response==============================================================================================================
	currentTaskResponseObject.Completed = true
	return currentTaskResponseObject
}
