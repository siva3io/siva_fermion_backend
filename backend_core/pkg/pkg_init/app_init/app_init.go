package app_init

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	db "fermion/backend_core/db"
	model_core "fermion/backend_core/internal/model/core"
	cache "fermion/backend_core/pkg/util/cache"
	helpers "fermion/backend_core/pkg/util/helpers"
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

func AppInit(app model_core.InstalledApps) {
	//TODO : unzip the folder for app
	SeedMasterDataForApp(app.Code)
	SetIntoCache(app)
	InitApiCall(app)
}

func SeedMasterDataForApp(app_code string) {
	type SqlModel struct {
		ModelName string `json:"model_name"`
		Link      string `json:"link"`
	}
	MasterDataForApp := []byte(`
		[
			{
				"ModelName": "Appfields and Channel LookupCodes",
				"Link":"master_data.sql"
			}
		]
	`)

	var sqlModel []SqlModel
	err := json.Unmarshal(MasterDataForApp, &sqlModel)
	if err != nil {
		fmt.Println(err)
	}

	for _, model := range sqlModel {
		fmt.Println("seeding master data for " + app_code + " ...")
		path := filepath.Join("./integrations/" + strings.ToLower(app_code) + "/data/" + model.Link)
		fileData, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		sqlData := string(fileData)

		err = db.DbManager().Exec(sqlData).Error
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("seeded " + app_code)
	}
}
func SetIntoCache(app model_core.InstalledApps) {
	status := db.RedisManager().Ping().Val()
	if status != "" {
		cache.SetCacheKey(app.Code, app)
	}
}
func InitApiCall(app model_core.InstalledApps) {
	fmt.Println("api call initiate started for " + app.Code + " ...")
	accessToken, _ := helpers.TokenBearerType(app.AccessToken)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": accessToken,
	}
	queryParams := map[string]string{
		"app_code": app.Code,
		"path":     fmt.Sprintf("integrations/%v/", strings.ToLower(app.Code)),
	}
	var body interface{}

	request := helpers.Request{
		Method: "GET",
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%v", os.Getenv("PORT")),
		Path:   fmt.Sprintf("integrations/%v/init", strings.ToLower(app.Code)),
		Header: headers,
		Params: queryParams,
		Body:   body,
	}
	helpers.MakeRequest(request)
	fmt.Println("api call initiate completed for " + app.Code)
}
