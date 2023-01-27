package model

import (
	"net/http"
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

// =====================Features==========================================================================
type Feature struct {
	Id          string `json:"_id" bson:"_id"`
	FeatureName string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	OnFailure   string `json:"on_failure" bson:"on_failure"` //--on_failure-[patch, rollback]
	Tasks       []Task `json:"tasks" bson:"tasks"`
}

// =====================Feature Tasks======================================================================
type Task struct {
	Id                   string                 `json:"task_id" bson:"task_id"`
	TaskName             string                 `json:"name" bson:"name"`
	Link                 string                 `json:"link" bson:"link"`
	Type                 string                 `json:"type" bson:"type"` //------type-[endpoint]
	PluginAuthRequired   *bool                  `json:"plugin_auth_required" bson:"plugin_auth_required"`
	Dependents           []string               `json:"dependents" bson:"dependents"`
	OnSuccess            []string               `json:"on_success" bson:"on_success"`
	OnSuccessSkipTasks   []string               `json:"on_success_skip_tasks" bson:"on_success_skip_tasks"`
	Mapper               []MapperTask           `json:"mapper" bson:"mapper"`
	RequiredMappedOutput []RequiredMappedOutput `json:"required_mapped_output" bson:"required_mapped_output"`
	Priority             int                    `json:"priority" bson:"priority"`
}
type MapperTask struct {
	MapperId   string `json:"mapper_id" bson:"mapper_id"`
	MapperName string `json:"name" bson:"name"`
	Link       string `json:"link" bson:"link"`
	InputType  string `json:"input_type" bson:"input_type"`
	InputKey   string `json:"input_key" bson:"input_key"`
}
type RequiredMappedOutput struct {
	MapperId   string `json:"mapper_id" bson:"mapper_id"`
	TaskId     string `json:"task_id" bson:"task_id"`
	TaskName   string `json:"name" bson:"name"`
	MapperName string `json:"mapper_name" bson:"mapper_name"`
	Type       string `json:"type" bson:"type"`
}

// ======================API Files======================================================================
type APIFile struct {
	ChannelCode          string               `json:"channel_code" bson:"channel_code"`
	Name                 string               `json:"name" bson:"name"`
	Description          string               `json:"description" bson:"description"`
	Endpoint             string               `json:"endpoint" bson:"endpoint"`
	Scheme               string               `json:"scheme" bson:"scheme"`
	RequestConfiguration RequestConfiguration `json:"request_configuration" bson:"request_configuration"`
	SessionVariables     []KeyValuePair       `json:"session_variables" bson:"session_variables"`
	Payload              Payload              `json:"payload" bson:"payload"`
	CompletionNorms      []Props              `json:"completion_norms" bson:"completion_norms"`
	SuccessCallBack      []Props              `json:"success_callback" bson:"success_callback"`
	Response             ResponseNorms        `json:"response" bson:"response"`
}
type RequestConfiguration struct {
	Method         string         `json:"method" bson:"method"` //--[GET, PUT, DELETE,POST,PATCH, OTHERS]
	Type           string         `json:"type" bson:"type"`     //--[single, paaginated, batch]
	RequestDetails []KeyValuePair `json:"request_details" bson:"request_details"`
}
type KeyValuePair struct {
	Key   string      `json:"key" bson:"key"`
	Value interface{} `json:"value" bson:"value"`
	Type  string      `json:"type" bson:"type"`
	Props []Props     `json:"props" bson:"props"`
}
type Props struct {
	Name   string         `json:"name" bson:"name"`
	Params []KeyValuePair `json:"params" bson:"params"`
}
type Payload struct {
	Type        string         `json:"type" bson:"type"`
	Headers     []KeyValuePair `json:"headers" bson:"headers"`
	QueryParams []KeyValuePair `json:"query_params" bson:"query_params"`
	Body        []KeyValuePair `json:"body" bson:"body"`
	Signature   []KeyValuePair `json:"signature" bson:"signature"`
}
type ResponseNorms struct {
	SessionVariables        []KeyValuePair `json:"session_variables" bson:"session_variables"`
	ResponseCompletionNorms []KeyValuePair `json:"response_completion_norms" bson:"response_completion_norms"`
}

// ======================Mapper======================================================================
type Mapper struct {
	Id          string        `json:"id" bson:"id"`
	ChannelCode string        `json:"channel_code" bson:"channel_code"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Fields      []MapperField `json:"fields" bson:"fields"`
}
type MapperField struct {
	Id             string      `json:"id" bson:"id"`
	Scope          string      `json:"scope" bson:"scope"`
	Input          string      `json:"input" bson:"input"`
	InputType      string      `json:"input_type" bson:"input_type"` //--[keyvalue, object, array,interface]
	Output         string      `json:"output" bson:"output"`
	OutputType     string      `json:"output_type" bson:"output_type"` //--[keyvalue, object, array]
	DefaultValue   interface{} `json:"default_value" bson:"default_value"`
	ExecutionType  string      `json:"execution_type" bson:"execution_type"`
	Fields         []string    `json:"fields" bson:"fields"`
	HelperFunction []Props     `json:"helper_function" bson:"helper_function"`
}

// ======================Task Endpoint Response======================================================================
type TaskEndpointResponse struct {
	Name string `json:"name"`
	EndpointResponse
	MappedResponse   []KeyValuePair         `json:"mapped_response"`
	Completed        bool                   `json:"completed"`
	SessionVariables []KeyValuePair         `json:"session_variables"`
	Errors           map[string]interface{} `json:"errors"`
}
type EndpointResponse struct {
	Header     http.Header            `json:"header"`
	Body       map[string]interface{} `json:"body"`
	Status     string                 `json:"status"`      // e.g. "200 OK"
	StatusCode int                    `json:"status_code"` // e.g. 200
	Request    *http.Request          `json:"-"`
}

// =====================Input Configuration For FeatureTasks==========================================================
type InputConfigForTask struct {
	Id             string                 `json:"id"`
	BasePath       string                 `json:"base_path"`
	AppAuthOptions map[string]interface{} `json:"app_auth_options"`
}
