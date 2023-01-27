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

func ExecuteSingleRequest(taskResponseArray []model.TaskEndpointResponse, taskResponseObject model.TaskEndpointResponse, singleTaskObject model.Task, endpointTaskFile model.APIFile, featureSessionVariables []model.KeyValuePair) model.TaskEndpointResponse {
	var requestBody interface{}
	var queryParams string
	var urlWithParams string
	var mappedResponseCompletionNorms []model.KeyValuePair
	var requestMethod string
	Errors := map[string]interface{}{}

	//========================payload details-[ header-[], body-[],query_params-[], signature-[] ]======================================================================================
	payload := endpointTaskFile.Payload

	//================================end point url=====================================================================================================================================
	getURL := endpointTaskFile.Scheme + "://" + endpointTaskFile.Endpoint
	getURL, err := utils.FormatEndpoint(getURL, featureSessionVariables)
	if err != nil {
		Errors["message"] = "error from format endpoint function"
		Errors["format_endpoint_error"] = fmt.Sprintf("%v", err)
		taskResponseObject.Errors = Errors
		taskResponseObject.Completed = false
		return taskResponseObject
	}

	//=================================query params=====================================================================================================================================
	if len(payload.QueryParams) > 0 {
		params := utils.ParseKeyValuePair(payload.QueryParams, featureSessionVariables)
		queryParams, err = utils.GetUrlQueryParams(params, featureSessionVariables)
		if err != nil {
			fmt.Println("=============> error in GetUrlQueryParams <============= ")
			Errors["message"] = "error in GetUrlQueryParams"
			Errors["query_params"] = fmt.Sprintf("%v", err)
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
	}
	urlWithParams = getURL + queryParams
	fmt.Println("=============> executing query params <============= " + singleTaskObject.TaskName)

	//===================if the task requires any mapper output then ftech the mapper output and append to the request_body============================================================
	if len(singleTaskObject.RequiredMappedOutput) > 0 {
		requestBody = repo.MergeMultipleMapperPayload(singleTaskObject.RequiredMappedOutput, taskResponseArray)
	}

	//==============================request body=======================================================================================================================================
	if len(payload.Body) > 0 {
		bodyKeyValuePair := utils.ParseKeyValuePair(payload.Body, featureSessionVariables)
		requestBody = utils.ConvertKeyValuePairToInterface(bodyKeyValuePair)
	}
	fmt.Println("=============> executing request_body <============= " + singleTaskObject.TaskName)

	//============================headers and request_method===========================================================================================================================
	headers := payload.Headers
	requestMethod = endpointTaskFile.RequestConfiguration.Method
	fmt.Println("=============> passing headers <============= " + singleTaskObject.TaskName)

	//==============call the MakeAPIRequest function and get the response of the API call==============================================================================================
	response, err := utils.MakeAPIRequest(requestMethod, urlWithParams, headers, requestBody, featureSessionVariables, endpointTaskFile)
	if err != nil {
		fmt.Println("=============> error in api call function <============= ")
		Errors["message"] = "error in api call function"
		Errors["api_error"] = fmt.Sprintf("%v", err)
		taskResponseObject.Errors = Errors
		taskResponseObject.Completed = false
		return taskResponseObject
	}

	//========================customize the API response===============================================================================================================================
	//===========================session variables=====================================================================================================================================
	if len(endpointTaskFile.Response.SessionVariables) > 0 {
		taskSessionVariables, err := utils.ConvertArrayInterfaceToArrayStructWithJsonKeyParse(endpointTaskFile.Response.SessionVariables, response, featureSessionVariables)
		if err != nil {
			fmt.Println("=============> error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function <=============")
			Errors["message"] = "error in ConvertArrayInterfaceToArrayStructWithJsonKeyParse function"
			Errors["response_session_variable _error"] = fmt.Sprintf("%v", err)
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		taskSessionVariables = utils.AppendOrUpdateKeyValuePair(taskSessionVariables, featureSessionVariables)
		taskResponseObject.SessionVariables = taskSessionVariables
	}

	taskResponseObject.EndpointResponse = response

	//=============================if mapper exists then execute the mapper===========================================================================================================
	if taskResponseObject.Body != nil && len(singleTaskObject.Mapper) > 0 {
		res, err := repo.ExecuteMapperForTask(singleTaskObject.Mapper, taskResponseObject.Body, mappedResponseCompletionNorms, featureSessionVariables)
		if err != nil {
			fmt.Println("=============> error in mapper function <=============")
			Errors["message"] = "error in mapper fuction"
			Errors["mapper_error"] = fmt.Sprintf("%v", err)
			taskResponseObject.Errors = Errors
			taskResponseObject.Completed = false
			return taskResponseObject
		}
		taskResponseObject.MappedResponse = res
	}

	//==============================================task response====================================================================================================================
	taskResponseObject.Completed = true

	return taskResponseObject
}
