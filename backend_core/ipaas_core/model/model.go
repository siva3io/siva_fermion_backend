package model

import (
	"net/http"
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
type InputConfigForTask struct {
	AppAuthOptions  *map[string]interface{} `json:"app_auth_options"`
	TaskName        string                  `json:"task_name"`
	BasePath        string                  `json:"base_path"`
	Id              string                  `json:"id"`
	FeatureFileData map[string]interface{}  `json:"feature_file_data"`
}

type KeyValuePair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	Props []Props     `json:"props"`
}

type Props struct {
	Name   string         `json:"name"`
	Params []KeyValuePair `json:"params"`
}

type EndpointResponse struct {
	Header     http.Header            `json:"header"`
	Body       map[string]interface{} `json:"body"`
	Status     string                 `json:"status"`      // e.g. "200 OK"
	StatusCode int                    `json:"status_code"` // e.g. 200
	Request    *http.Request          `json:"-"`
}

type TaskEndpointResponse struct {
	Name             string `json:"name"`
	EndpointResponse `json:"endpoint_response"`
	MappedResponse   []KeyValuePair         `json:"mapped_response"`
	Completed        bool                   `json:"completed"`
	SessionVariables []KeyValuePair         `json:"session_variables"`
	Errors           map[string]interface{} `json:"errors"`
}
