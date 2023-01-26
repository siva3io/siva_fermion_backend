package core

import (
	"net/http"
	"time"
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
type LogError struct {
	Id            string      `json:"id"`
	RequestHeader http.Header `json:"request_header"`
	RequestBody   string      `json:"request_body"`
	URL           string      `json:"url"`
	HttpMethod    string      `json:"http_method"`
	FunctionName  string      `json:"function_name"`
	ErrorMessage  string      `json:"error_message"`
	Level         string      `json:"level"`
	CreatedAt     time.Time   `json:"created_at"`
}

type LogActivity struct {
	Id            string      `json:"id"`
	RequestHeader http.Header `json:"request_header"`
	RequestBody   string      `json:"request_body"`
	ResponseBody  string      `json:"response_body"`
	URL           string      `json:"url"`
	HttpMethod    string      `json:"http_method"`
	FunctionName  string      `json:"function_name"`
	CreatedAt     time.Time   `json:"created_at"`
}

type LogLogin struct {
	Id            string      `json:"id"`
	RequestHeader http.Header `json:"request_header"`
	RequestBody   string      `json:"request_body"`
	URL           string      `json:"url"`
	HttpMethod    string      `json:"http_method"`
	UserID        string      `json:"user_id"`
	CreatedAt     string      `json:"created_at"`
}
