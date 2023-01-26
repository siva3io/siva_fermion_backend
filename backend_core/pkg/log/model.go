package log

import (
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
	ID           string    `json:"id"`
	Header       string    `json:"request_header"`
	Body         string    `json:"request_body"`
	URL          string    `json:"url"`
	HttpMethod   string    `json:"http_method"`
	Email        string    `json:"email"`
	ErrorMessage string    `json:"error_message"`
	Level        string    `json:"level"`
	AppName      string    `json:"app_name"`
	Version      string    `json:"version"`
	Env          string    `json:"env"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LogActivity struct {
	ID        string    `json:"id"`
	TableName string    `json:"table_name"`
	Email     string    `json:"email"`
	Row       string    `json:"row"`
	NewData   string    `json:"new_data"`
	OldData   string    `json:"old_data"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LogLogin struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DateLogin   string    `json:"date_login"`
	DateLogout  string    `json:"date_logout"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
