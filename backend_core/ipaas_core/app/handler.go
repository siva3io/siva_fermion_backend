package ipaas_core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	model "fermion/backend_core/ipaas_core/model"
	repo "fermion/backend_core/ipaas_core/repository"
	"fermion/backend_core/ipaas_core/utils"
	pkg_helpers "fermion/backend_core/pkg/util/helpers"
	pkg_response "fermion/backend_core/pkg/util/response"

	"github.com/labstack/echo/v4"
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

type handler struct {
	service Service
}

var IpaasCoreHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if IpaasCoreHandler != nil {
		return IpaasCoreHandler
	}
	service := NewService()
	IpaasCoreHandler = &handler{service}
	return IpaasCoreHandler
}

/*
This function mainly represents to execute the feature.

It holds the inputs like

	`1. c :  context of the current HTTP request.
	 2. input : input_configuration of input parameters`

Each feature contains 'n' of tasks. This function execute the ExecuteFeatureTasks functions get the response of the function.
Finally, it returns the reponse of all the tasks as a single response

Function flow :

	`1. get the platform token to access the credentials of the 3rd party service of core API
	 2. by using the access token to get the feature file data, which neeeds to be executed
	 3. initialize the query_params and auth_token in featureSessionVariables to execute the tasks of the feature
	 4. add the request_body, base_path and task_list_object_feature to the feature_session variables`
*/
func ExecuteFeature(c echo.Context, input model.InputConfigForTask) ([]interface{}, error) {
	var Feature model.Feature
	var tasksListObjectFromFeature []model.Task
	var featureSessionVariables []model.KeyValuePair
	var taskResponseArray []model.TaskEndpointResponse

	var queryParams url.Values
	var requestBody interface{}
	var response []interface{}

	//======================dynamic app token for all the feature api calls=======================================================
	featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, utils.ConvertObjectToKeyValuePair(input.AppAuthOptions))

	//======================read feature=======================================================================================================
	fileData, err := utils.ReadFeature("features", input.Id)
	if err != nil {
		return nil, err
	}
	err = pkg_helpers.JsonMarshaller(fileData, &Feature)
	if err != nil {
		return nil, err
	}

	//==================initialize query_parameters to feature_session_variables======================================================================
	queryParams = c.Request().URL.Query()
	for key, value := range queryParams {
		if len(value) > 1 {
			featureSessionVariables = append(featureSessionVariables, utils.MakeKeyValuePair(key, value, "static", nil))
			continue
		}
		featureSessionVariables = append(featureSessionVariables, utils.MakeKeyValuePair(key, value[0], "static", nil))
	}

	//==================initialize request_body to feature_session_variables==========================================================================
	requestBody = utils.GetBodyData(c)
	featureSessionVariables = append(featureSessionVariables, model.KeyValuePair{Key: "requestPayload", Value: requestBody})

	//=================take tasks from feature========================================================================================================
	tasksListObjectFromFeature = Feature.Tasks

	//=================add basePath and tasksListObjectFromFeature to the featureSessionVariables=====================================================
	featureSessionVariables = append(featureSessionVariables,
		[]model.KeyValuePair{
			{Key: "basePath", Value: input.BasePath},
			{Key: "tasksListObjectFromFeature", Value: tasksListObjectFromFeature}}...)

	//==================execute the feature tasks=====================================================================================================
	taskResponseArray = ExecuteFeatureTasks(tasksListObjectFromFeature, featureSessionVariables, taskResponseArray)

	for _, taskResponse := range taskResponseArray {
		fmt.Println(taskResponse.Name)
		if taskResponse.Errors != nil {
			taskResponseObj := make(map[string]interface{}, 0)
			taskResponseObj["name"] = taskResponse.Name
			taskResponseObj["errors"] = taskResponse.Errors
			taskResponseObj["completed"] = taskResponse.Completed
			response = append(response, taskResponseObj)
			continue
		}
		response = append(response, taskResponse)
	}
	return response, nil
}

/*
This function mainly represents to execute the feature tasks.

It holds the input parameters like the following are :

	   `1. tasksListObjectFromFeature : list of tasks
		2. featureSessionVariables : holds the requires session variables to execute all the task of the feature
		3. taskResponseArray : holds the response of all the task of the feature`

# It returns the reponse of the all the tasks of the feature

Function flow :

	   `1. loop all the tasks of the feature
		2. check the dependency task is executed or not. if it is successfully executed then execute the current task
		3. next it checks the task is present in the skip task list or not. if it is not present then execute the current task
		4. then it checks the task is endpoint or not. if it is endpoint then execute the current task
		5. once the task is completed then store the required session_variables in featureSessionVariables
		6. it checks if any of the tasks need to skip after execution of the task, then it will add the tasks in skipTaskList`
*/
func ExecuteFeatureTasks(tasksListObjectFromFeature []model.Task, featureSessionVariables []model.KeyValuePair, taskResponseArray []model.TaskEndpointResponse) []model.TaskEndpointResponse {
	var tasksToSkip []string

	//=======================execute each task from the list of task_object from the feature============================================================
	for task := 0; task < len(tasksListObjectFromFeature); task++ {
		singleTaskObject := tasksListObjectFromFeature[task]

		//====================check dependency tasks were executed successfully or not===================================================================
		if len(singleTaskObject.Dependents) > 0 {
			dependencyStatus := true
			for _, dependencyTaskName := range singleTaskObject.Dependents {
				for _, taskRespose := range taskResponseArray {
					if taskRespose.Name == dependencyTaskName {
						if !taskRespose.Completed {
							dependencyStatus = false
							break
						}
					}
				}
				if !dependencyStatus {
					break
				}
			}
			if !dependencyStatus {
				var taskResponseObject model.TaskEndpointResponse
				taskResponseObject.Name = singleTaskObject.TaskName
				taskResponseObject.Completed = false
				taskResponseArray = append(taskResponseArray, taskResponseObject)
				continue
			}
		}

		//========================if the task was already executed then skip the task====================================================================
		if pkg_helpers.Contains(tasksToSkip, singleTaskObject.TaskName) {
			var taskResponseObject model.TaskEndpointResponse
			taskResponseObject.Name = singleTaskObject.TaskName
			taskResponseObject.Completed = true
			taskResponseArray = append(taskResponseArray, taskResponseObject)
			continue
		}

		//===============================execute the task if the type is endpoint========================================================================
		if singleTaskObject.Type == "endpoint" {
			fmt.Println("=============> singleTaskObject <============= " + singleTaskObject.TaskName)
			//=================================execute the endpoint task=================================================================================
			taskResponseObject := ExecuteTaskForEndpoint(singleTaskObject, featureSessionVariables, taskResponseArray)
			fmt.Println("=============> ExecTaskForEndpoint successfully <============= ")
			fmt.Println("===================================================================================================================================================================================================")

			//==================add or update the endpoint task sessionVariables to FeatureSession variables=============================================
			featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, taskResponseObject.SessionVariables)

			//==========================append the endpoint task  response===============================================================================
			taskResponseArray = append(taskResponseArray, taskResponseObject)

			//=================if task contain any skip tasks then add to the taskToSkip array===========================================================
			tasksToSkip = append(tasksToSkip, singleTaskObject.OnSuccessSkipTasks...)
		}
	}
	return taskResponseArray
}

/*
This function mainly represents to execute the api file of the particular task.

It holds the input parameters like the following are :

	   `1. singleTaskObject : the task to execute
		2. featureSessionVariables : the variable which holds the required values from previous tasks session_variables as well as initial values
		3. taskResponseArray : which holds the previous task response`

# It returns the endpoint response of the task

Function flow :

	   `1. read the api file
		2. initialize the session_variables of this task
		3. execute the endpoint request based on the request_type call the function

		request_types are : [`single`, `paginated`, `batch`, ... ]
		4. get the task response and return it`
*/
func ExecuteTaskForEndpoint(singleTaskObject model.Task, featureSessionVariables []model.KeyValuePair, taskResponseArray []model.TaskEndpointResponse) model.TaskEndpointResponse {
	var taskResponseObject model.TaskEndpointResponse
	var endpointTaskFile model.APIFile

	//=====================task name=====================================================================================================================
	taskResponseObject.Name = singleTaskObject.TaskName

	fileData, err := utils.ReadFeature("apis", singleTaskObject.Id)
	if err != nil {
		fmt.Println("=============> Error reading endpointTaskFile <============= ")
		return taskResponseObject
	}
	err = pkg_helpers.JsonMarshaller(fileData, &endpointTaskFile)
	if err != nil {
		return taskResponseObject
	}

	//====================initialize the sessionVariables to featureSessionVariables, if it exists======================================================
	if len(endpointTaskFile.SessionVariables) > 0 {
		sessionVariables := utils.ParseKeyValuePair(endpointTaskFile.SessionVariables, featureSessionVariables)
		featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, sessionVariables)
	}

	//===========================execute the task based on request type=================================================================================
	//=============================single request type configuration====================================================================================
	if endpointTaskFile.RequestConfiguration.Type == "single" {
		fmt.Println("=============> executing endpoint request type - SINGLE request <============= ")
		return ExecuteSingleRequest(taskResponseArray, taskResponseObject, singleTaskObject, endpointTaskFile, featureSessionVariables)

	}
	//=========================pagination request type configuration====================================================================================
	if endpointTaskFile.RequestConfiguration.Type == "paginated" {
		fmt.Println("=============> executing endpoint request type - PAGINATED request <============= ")
		return ExecutePaginationRequest(taskResponseArray, taskResponseObject, singleTaskObject, endpointTaskFile, featureSessionVariables)

	}
	//=============================bulk request type configuration======================================================================================
	if endpointTaskFile.RequestConfiguration.Type == "batch" {
		fmt.Println("=============> executing endpoint request type - BATCH request <============= ")
		return ExecuteBatchRequest(taskResponseArray, taskResponseObject, singleTaskObject, endpointTaskFile, featureSessionVariables)
	}
	return taskResponseObject
}

// ========================ipaas_rate_calculator=================================================================================================================================================================
func (h *handler) GenericHandler(c echo.Context) (err error) {
	//get dynamic host
	host := c.Request().Host
	//get dynamic scheme
	scheme := c.Scheme()
	module_name := c.Param("module")
	task_name := c.Param("task")
	user_token := c.Request().Header.Get("Authorization")
	//-----------------------------------Shipping module-----------------------------------------------------------------------
	// TODO: Remove the if conditions
	if module_name == "shipping" {
		//---------------------------------Tracking------------------------------------------------------------------------------
		if task_name == "tracking" {
			tracking_id := c.QueryParam("tracking_id")
			shipping_partner_id := c.QueryParam("shipping_partner_id")
			var requestBody interface{}
			var headers []model.KeyValuePair
			headers = append(headers, utils.MakeKeyValuePair("Content-Type", "application/json", "static", nil))
			headers = append(headers, utils.MakeKeyValuePair("Authorization", user_token, "static", nil))
			spo_resp, _ := utils.MakeAPIRequest("GET", scheme+"://"+host+"/api/v1/shipping_orders/"+tracking_id, headers, requestBody, nil)
			if spo_resp.Body["meta"].(map[string]interface{})["success"].(bool) {
				tracking_id = spo_resp.Body["data"].(map[string]interface{})["awb_number"].(string)
			}
			//apicall

			response, _ := utils.MakeAPIRequest("GET", scheme+"://"+host+"/api/v1/shipping_partners/"+shipping_partner_id, headers, requestBody, nil)
			shipping_partner_name := response.Body["data"].(map[string]interface{})["partner_name"].(string)
			shipping_partner_name = strings.ToLower(shipping_partner_name)
			tracking_resp, _ := utils.MakeAPIRequest("GET", scheme+"://"+host+"/integrations/"+shipping_partner_name+"/"+task_name+"/"+tracking_id, headers, requestBody, nil)
			return pkg_response.RespSuccess(c, "Tracking details retrieved successfully", tracking_resp.Body)
		}
		//---------------------------------Task level error handling------------------------------------------------------------
		return pkg_response.RespValidationErr(c, "task name "+task_name+" not found inside "+module_name+" module", err)
	}
	//-----------------------------------Module level error handling-----------------------------------------------------------
	return pkg_response.RespValidationErr(c, "module name "+module_name+" not found", err)
}

/**/
func (h *handler) RateCalculator(c echo.Context) (err error) {
	host := c.Request().Host
	scheme := c.Scheme()
	var request interface{}
	var requestBody interface{}
	var headers []model.KeyValuePair
	var map_data []interface{}
	user_token := c.Request().Header.Get("Authorization")
	headers = append(headers, utils.MakeKeyValuePair("Content-Type", "application/json", "static", nil))
	headers = append(headers, utils.MakeKeyValuePair("Authorization", user_token, "static", nil))
	//var fsv []model.KeyValuePair
	response, _ := utils.MakeAPIRequest("GET", scheme+"://"+host+"/api/v1/core/apps/installed?filters=[[\"category_name\",\"=\",\"shipping\"]]", headers, request, nil)
	list := make(map[string]interface{})
	list["express"] = make([]interface{}, 0)
	list["surface"] = make([]interface{}, 0)
	var str []string
	//c.Bind(&requestBody)
	requestBody = utils.GetBodyData(c)
	wg := sync.WaitGroup{}
	apps := response.Body["data"].([]interface{})
	for _, app := range apps {
		data := app.(map[string]interface{})["name"].(string)
		str = append(str, strings.ToLower(data))
	}
	// fmt.Println(str)
	for _, app := range str {
		wg.Add(1)
		go func(idx string) {
			response, _ := utils.MakeAPIRequest("POST", scheme+"://"+host+"/integrations/"+idx+"/"+"rate_calculator", headers, requestBody, nil)
			if response.Body["data"] != nil {
				data := response.Body["data"].(map[string]interface{})
				if data["express"] != nil {
					map_data = data["express"].([]interface{})
					for _, single_data := range map_data {
						single_map := single_data.(map[string]interface{})
						shipping_partner := single_map["name"].(string)
						single_map["name"] = shipping_partner + "-" + idx
						//delete(single_map,"shipping_partner")
						list["express"] = append(list["express"].([]interface{}), single_map)
					}
				}
				if data["surface"] != nil {
					map_data = data["surface"].([]interface{})
					for _, single_data := range map_data {
						single_map := single_data.(map[string]interface{})
						shipping_partner := single_map["name"].(string)
						single_map["name"] = shipping_partner + "-" + idx
						//delete(single_map,"shipping_partner")
						list["surface"] = append(list["surface"].([]interface{}), single_map)
					}
				}
			}
			//list[idx]=data
			wg.Done()
		}(app)
	}
	wg.Wait()
	return pkg_response.RespSuccess(c, "Shipping partners rates retrieved successfully", list)
}

// ========================ipass_integrations_handler============================================================================================================================================================
func (h *handler) GetFeature(c echo.Context) (err error) {

	id := c.Param("id")
	collection := c.Param("collection")

	feature, err := h.service.GetFeature(collection, id)

	if err != nil {
		return pkg_response.RespErr(c, err)
	}

	return pkg_response.RespSuccess(c, "success", feature)
}
func (h *handler) CreateFeature(c echo.Context) (err error) {

	requestBody := make(map[string]interface{}, 0)
	err = json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return pkg_response.RespErr(c, err)
	}

	collection := c.Param("collection")

	feature, err := h.service.CreateFeature(collection, requestBody)

	if err != nil {
		return pkg_response.RespErr(c, err)
	}

	return pkg_response.RespSuccess(c, "success", feature)
}
func (h *handler) UpdateFeature(c echo.Context) (err error) {

	requestBody := make(map[string]interface{}, 0)
	err = json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return pkg_response.RespErr(c, err)
	}

	id := c.Param("id")
	collection := c.Param("collection")

	feature, err := h.service.UpdateFeature(collection, id, requestBody)

	if err != nil {
		return pkg_response.RespErr(c, err)
	}

	return pkg_response.RespSuccess(c, "success", feature)
}
func (h *handler) DeleteFeature(c echo.Context) (err error) {

	id := c.Param("id")
	collection := c.Param("collection")

	data, err := h.service.DeleteFeature(collection, id)

	if err != nil {
		return pkg_response.RespErr(c, err)
	}

	return pkg_response.RespSuccess(c, "success", data)
}
func (h *handler) GetAppFeatures(c echo.Context) (err error) {

	query := map[string]interface{}{
		"app_code": c.Param("app_code"),
	}

	features, err := h.service.GetAppFeatures(query)
	if err != nil {
		return pkg_response.RespErr(c, err)
	}

	return pkg_response.RespSuccess(c, "success", features)
}

// ==============================Boson convertor==================================================================================================================================================================
func (h *handler) BosonConvertor(c echo.Context) (err error) {
	// isDataEncryption := true
	// encryption := c.QueryParam("encryption")
	// if encryption != "" {
	// 	isDataEncryption, _ = strconv.ParseBool(encryption)
	// }

	// requestId := uuid.New().String()

	var response BosonConvertorDTO
	var data BosonConvertorDataOuput

	requestPayload := map[string]interface{}{}
	err = c.Bind(&requestPayload)
	if err != nil {
		response.Data = BosonConvertorDataOuput{
			ErrorMessage: err,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	if requestPayload["data"] == nil {
		response.Data = BosonConvertorDataOuput{
			ErrorMessage: errors.New("data not found in request payload"),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	var bosonConvertorDataInput BosonConvertorDataInput
	err = pkg_helpers.JsonMarshaller(requestPayload["data"], &bosonConvertorDataInput)
	if err != nil {
		response.Data = BosonConvertorDataOuput{
			ErrorMessage: err,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	if metaData := requestPayload["meta_data"]; metaData != nil {
		response.MetaData = metaData.(map[string]interface{})
	}

	// if requestPayload["meta_data"] != nil {
	// 	metaData, ok := requestPayload["meta_data"].(map[string]interface{})
	// 	if !ok {
	// 		pkg_response.RespErr(c, errors.New("oops! Request payload meta_data error"))
	// 	}
	// 	if metaData["request_id"] == nil {
	// 		metaData["request_id"] = requestId
	// 	}
	// 	metaData["encryption"] = isDataEncryption
	// 	requestPayload["meta_data"] = metaData
	// } else {
	// 	requestPayload["meta_data"] = map[string]interface{}{
	// 		"request_id": requestId,
	// 		"encryption": isDataEncryption,
	// 	}
	// }
	// eda.Produce(eda.BOSON_MAPPER_CONVERTOR, requestPayload, isDataEncryption)

	// edaResponse := map[string]interface{}{}
	// for {
	// 	output := cache.GetCacheVariable(requestId)
	// 	if output != nil {
	// 		var resp map[string]interface{}
	// 		err := pkg_helpers.JsonMarshaller(output, &resp)
	// 		if err != nil {
	// 			return pkg_response.RespErr(c, err)
	// 		}
	// 		metaData := resp["meta_data"].(map[string]interface{})
	// 		if metaData["encryption"].(bool) {
	// 			decrypted_payload, _ := pkg_helpers.AESGCMDecrypt(resp["data"].(string))
	// 			var decryptedResponseData interface{}
	// 			err := json.Unmarshal([]byte(decrypted_payload), &decryptedResponseData)
	// 			if err != nil {
	// 				pkg_response.RespErr(c, err)
	// 			}
	// 			resp["data"] = decryptedResponseData
	// 		}
	// 		err = pkg_helpers.JsonMarshaller(resp, &edaResponse)
	// 		if err != nil {
	// 			pkg_response.RespErr(c, err)
	// 		}
	// 		break
	// 	}
	// }

	mapperInputType := "object"
	_, ok := bosonConvertorDataInput.InputData.([]interface{})
	if ok {
		mapperInputType = "array"
	}

	// var data BosonConvertorDataOuput

	//=================execute object mapper=======================================================
	if mapperInputType == "object" {
		data.MappedResponse, data.ErrorMessage = repo.ExecuteObjMapper(bosonConvertorDataInput.InputData, bosonConvertorDataInput.MapperTemplate, bosonConvertorDataInput.FeatureSessionVariables)
		if data.ErrorMessage != nil {
			response.Data = data
			return c.JSON(http.StatusBadRequest, response)
		}
	}
	//=================execute array mapper=======================================================
	if mapperInputType == "array" {
		var mapperInput []interface{}

		err = pkg_helpers.JsonMarshaller(bosonConvertorDataInput.InputData, &mapperInput)
		if err != nil {
			response.Data = BosonConvertorDataOuput{
				ErrorMessage: err,
			}
			return c.JSON(http.StatusBadRequest, response)
		}
		data.MappedResponse, data.ErrorMessage = repo.ExecuteArrayMapper(mapperInput, bosonConvertorDataInput.MapperTemplate, bosonConvertorDataInput.FeatureSessionVariables)
		if data.ErrorMessage != nil {
			response.Data = data
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	response.Data = data
	return c.JSON(http.StatusOK, response)
}
