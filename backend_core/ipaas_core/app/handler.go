package ipaas_core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"

	cache "fermion/backend_core/controllers/cache"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/repository"
	model "fermion/backend_core/ipaas_core/model"
	"fermion/backend_core/ipaas_core/utils"
	res "fermion/backend_core/pkg/util/response"

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
	service        Service
	coreRepository repository.Core
}

func NewIpaasHandler() *handler {
	return &handler{
		NewService(),
		repository.NewCore(),
	}
}

func Init() {
	CacheSync()
}

func ExecuteFeature(c echo.Context, input model.InputConfigForTask) ([]interface{}, error) {
	var tasksListObjectFromFeature []map[string]interface{}
	var featureSessionVariables []model.KeyValuePair
	var taskResponseArray []model.TaskEndpointResponse

	var queryParams url.Values
	var requestBody interface{}
	var response []interface{}

	//----------------dynamic app token for all the feature api calls---------------------------------------------------
	authOptions := *input.AppAuthOptions
	tokenData := authOptions["token"].(string)
	featureSessionVariables = append(featureSessionVariables, utils.MakeKeyValuePair("euni_access_token", tokenData, "static", nil))

	//-----------------read feature file-------------------------------------------------
	featureFileData, err := utils.ReadFeature("features", input.Id)
	if err != nil {
		return nil, err
	}

	//-----------------initialize query_parameters to feature_session_variables---------------------------------------------------
	queryParams = c.Request().URL.Query()
	for key, value := range queryParams {
		if len(value) > 1 {
			featureSessionVariables = append(featureSessionVariables, utils.MakeKeyValuePair(key, value, "static", nil))
			continue
		}
		featureSessionVariables = append(featureSessionVariables, utils.MakeKeyValuePair(key, value[0], "static", nil))
	}

	//---------------initialize request_body to feature_session_variables---------------------------------------------------------
	json.NewDecoder(c.Request().Body).Decode(&requestBody)
	featureSessionVariables = append(featureSessionVariables, model.KeyValuePair{Key: "requestPayload", Value: requestBody})

	//-------------------read FeatureFileData and BasePath -----------------------------------------------------------------------
	// featureFileData := input.FeatureFileData
	basePath := input.BasePath

	//--------------------get the task list from the feature file data------------------------------------------------------------
	tasksListStringFromFeature, err := json.Marshal(featureFileData["tasks"])
	if err != nil {
		fmt.Println("--------> marshal error for tasksListStringFromFeature <----------------")
	}
	if err := json.Unmarshal(tasksListStringFromFeature, &tasksListObjectFromFeature); err != nil {
		return nil, err
	}

	//------------------------add basePath and tasksListObjectFromFeature to the featureSessionVariables--------------------------
	featureSessionVariables = append(featureSessionVariables,
		[]model.KeyValuePair{
			{Key: "basePath", Value: basePath},
			{Key: "tasksListObjectFromFeature", Value: tasksListObjectFromFeature}}...)

	//----------------------execute the feature tasks----------------------------------------------------------------------------
	result := ExecuteFeatureTasks(tasksListObjectFromFeature, featureSessionVariables, taskResponseArray)

	for _, taskResponse := range result {
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
func ExecuteFeatureTasks(tasksListObjectFromFeature []map[string]interface{}, featureSessionVariables []model.KeyValuePair, taskResponseArray []model.TaskEndpointResponse) []model.TaskEndpointResponse {
	var tasksToSkip []string

	//-------------------execute each task from the list of tasks object from the feature------------------------------------------
	for task := 0; task < len(tasksListObjectFromFeature); task++ {
		singleTaskObject := tasksListObjectFromFeature[task]
		// Errors := map[string]interface{}{}

		//----------------- check dependency task is executed successfully or not--------------------------------------------------
		if singleTaskObject["dependents"] != nil && len(singleTaskObject["dependents"].([]interface{})) > 0 {
			dependencyStatus := true
			for _, dependencyTaskName := range singleTaskObject["dependents"].([]interface{}) {
				for _, taskRespose := range taskResponseArray {
					if taskRespose.Name == dependencyTaskName.(string) {
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
				taskResponseObject.Name = singleTaskObject["name"].(string)
				taskResponseObject.Completed = false
				taskResponseArray = append(taskResponseArray, taskResponseObject)
				continue
			}
		}

		//---------------------------if the task was already executed then skip the task----------------------------------------------------------------
		if utils.ContainString(tasksToSkip, singleTaskObject["name"].(string)) {
			var taskResponseObject model.TaskEndpointResponse
			taskResponseObject.Name = singleTaskObject["name"].(string)
			taskResponseObject.Completed = true
			taskResponseArray = append(taskResponseArray, taskResponseObject)
			continue
		}

		//---------------------------execute the task if the type is endpoint------------------------------------------------------
		if singleTaskObject["type"] == "endpoint" {
			fmt.Println("--------> singleTaskObject <----------------" + singleTaskObject["name"].(string))

			//------------------------------execute the endpoint task -------------------------------------------------------------
			taskResponseObject := ExecuteTaskForEndpoint(singleTaskObject, featureSessionVariables, taskResponseArray)
			fmt.Println("--------> ExecTaskForEndpoint successfully <----------------")
			fmt.Println("==================================================================================================================================================================")

			//-------------add or update the endpoint task sessionVariables to FeatureSession variables ---------------------------
			featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, taskResponseObject.SessionVariables)
			// taskResponseObject.Completed = true

			//-----------------------append the endpoint task  response -----------------------------------------------------------
			taskResponseArray = append(taskResponseArray, taskResponseObject)

			//------------------if task contain any skip tasks then add to the taskToSkip array------------------------------------
			if singleTaskObject["on_success_skip_tasks"] != nil {
				marshaldata, _ := json.Marshal(singleTaskObject["on_success_skip_tasks"])
				var skipTasks []string
				json.Unmarshal(marshaldata, &skipTasks)
				tasksToSkip = append(tasksToSkip, skipTasks...)
			}
		}
	}
	return taskResponseArray
}
func ExecuteTaskForEndpoint(singleTaskObject map[string]interface{}, featureSessionVariables []model.KeyValuePair, taskResponseArray []model.TaskEndpointResponse) model.TaskEndpointResponse {
	var taskResponseObject model.TaskEndpointResponse

	//------------------task name--------------------------------------------------------------------------------------------
	taskResponseObject.Name = singleTaskObject["name"].(string)

	//------------------get endpointTaskFile for the task--------------------------------------------------------------------
	// basePath, err := utils.GetValueFromSessionVariablesKey("basePath", featureSessionVariables)
	// if err != nil {
	// 	fmt.Println("--------> base path not exists <----------------")
	// 	return taskResponseObject
	// }

	endpointTaskFile, err := utils.ReadFeature("apis", singleTaskObject["task_id"].(string))
	if err != nil {
		fmt.Println("--------> Error reading endpointTaskFile <----------------")
		return taskResponseObject
	}

	//------------------initialize the sessionVariables to featureSessionVariables, if it exists-----------------------------
	if endpointTaskFile["session_variables"] != nil {
		jsonSessionVariables, ok := endpointTaskFile["session_variables"].([]interface{})
		if ok {
			sessionVariables := utils.ParseObjectsFromConfigPayload(jsonSessionVariables)
			sessionVariables = utils.ParseKeyValuePair(sessionVariables, featureSessionVariables)
			featureSessionVariables = utils.AppendOrUpdateKeyValuePair(featureSessionVariables, sessionVariables)
		}
	}

	//----------------------get request configuration for the task-----------------------------------------------------------
	requestConfiguration, ok := endpointTaskFile["request_configuration"].(map[string]interface{})
	if !ok {
		fmt.Println("--------> error in request_configuration <----------------")
	}

	//----------------------execute the task based on request type-----------------------------------------------------------
	//-----------------------------single request type configuration---------------------------------------------------------
	if requestConfiguration["type"].(string) == "single" {
		fmt.Println("--------> executing endpoint request type - SINGLE request <----------------")
		return ExecuteSingleRequest(taskResponseArray, taskResponseObject, singleTaskObject, endpointTaskFile, featureSessionVariables)

	}
	//-----------------------------pagination request type configuration------------------------------------------------------
	if requestConfiguration["type"].(string) == "paginated" {
		fmt.Println("--------> executing endpoint request type - PAGINATED request <----------------")
		return ExecutePaginationRequest(taskResponseArray, taskResponseObject, singleTaskObject, endpointTaskFile, featureSessionVariables)

	}
	//-----------------------------bulk request type configuration------------------------------------------------------
	if requestConfiguration["type"].(string) == "batch" {
		fmt.Println("--------> executing endpoint request type - BATCH request <----------------")
		return ExecuteBatchRequest(taskResponseArray, taskResponseObject, singleTaskObject, endpointTaskFile, featureSessionVariables)

	}

	return taskResponseObject
}

// need to handle error
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
			return res.RespSuccess(c, "Tracking details retrieved successfully", tracking_resp.Body)
		}
		//---------------------------------Task level error handling------------------------------------------------------------
		return res.RespValidationErr(c, "task name "+task_name+" not found inside "+module_name+" module", err)
	}
	//-----------------------------------Module level error handling-----------------------------------------------------------
	return res.RespValidationErr(c, "module name "+module_name+" not found", err)
}

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
	return res.RespSuccess(c, "Shipping partners rates retrieved successfully", list)
}

func CacheSync() {
	p := new(pagination.Paginatevalue)
	p.Per_page = 1000
	installed_apps, _ := NewIpaasHandler().coreRepository.ListInstalledApps(p)
	var applist []string
	for _, app := range installed_apps {
		app_code := app.Code
		app_code = strings.ReplaceAll(app_code, " ", "")
		cache.SetCacheVariable(app_code+"AUTHORIZATION", app.AccessToken)
		applist = append(applist, app_code)
	}
	fmt.Println("Cache sync is successfull for the following installed applications")
	fmt.Println(applist)
}

func (h *handler) GetFeature(c echo.Context) (err error) {

	id := c.Param("id")
	collection := c.Param("collection")

	feature, err := h.service.GetFeature(collection, id)

	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "success", feature)
}

func (h *handler) CreateFeature(c echo.Context) (err error) {

	requestBody := make(map[string]interface{}, 0)
	err = json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return res.RespErr(c, err)
	}

	collection := c.Param("collection")

	feature, err := h.service.CreateFeature(collection, requestBody)

	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "success", feature)
}

func (h *handler) UpdateFeature(c echo.Context) (err error) {

	requestBody := make(map[string]interface{}, 0)
	err = json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	collection := c.Param("collection")

	feature, err := h.service.UpdateFeature(collection, id, requestBody)

	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "success", feature)
}

func (h *handler) DeleteFeature(c echo.Context) (err error) {

	id := c.Param("id")
	collection := c.Param("collection")

	data, err := h.service.DeleteFeature(collection, id)

	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "success", data)
}
