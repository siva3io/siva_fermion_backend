package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"fermion/backend_core/internal/model/core"

	"gorm.io/gorm"
)

const PLATFORM_VERSION = "1.0.0.1"
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
type MigrationVersionControl struct {
	core.Model
	Version         string `json:"version" gorm:"unique"`
	Description     string `json:"description"`
	ModelsAdded     int    `json:"models_added"`
	ModelsUpdated   int    `json:"models_updated"`
	ModelsDeleted   int    `json:"models_deleted"`
	ColumnsAdded    int    `json:"columns_added"`
	ColumnsUpdated  int    `json:"columns_updated"`
	ColumnsDeleted  int    `json:"columns_deleted"`
	SeedDataAdded   int    `json:"seed_data_added"`
	SeedDataUpdated int    `json:"seed_data_updated"`
	SeedDataDeleted int    `json:"seed_data_deleted"`
}

func Migrate(db *gorm.DB) {

	upgradeVersionString := strings.Split(PLATFORM_VERSION, ".")[3]

	var migrationVersionControl MigrationVersionControl
	err := db.Order("id desc").First(&migrationVersionControl).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("MIGRATION ERROR:", err)
		return
	}

	upgradeVersion, _ := strconv.Atoi(upgradeVersionString)
	currentVersion, _ := strconv.Atoi(migrationVersionControl.Version)

	fmt.Println("Current DB Version: ", currentVersion)

	if upgradeVersion <= currentVersion {
		fmt.Println("Already in latest DB Version: ", currentVersion)
		return
	}

	for version := currentVersion + 1; version <= upgradeVersion; version++ {

		fileBytes, err := ioutil.ReadFile("./migrations/" + fmt.Sprint(version) + "/index.json")
		if err != nil {
			fmt.Println("MIGRATION ERROR:", err)
			return
		}

		var index map[string]interface{}
		err = json.Unmarshal(fileBytes, &index)
		if err != nil {
			fmt.Println("MIGRATION ERROR:", err)
			return
		}

		newModels := index["auto_migrate_models"].([]interface{})
		for _, modelName := range newModels {
			model := Tables[modelName.(string)]
			err = db.AutoMigrate(model)
			if err != nil {
				fmt.Println("MIGRATION ERROR:", err)
				return
			}
		}

		sqlFiles := index["sql_files"].([]interface{})
		for _, fileName := range sqlFiles {
			fileBytes, err = ioutil.ReadFile("./migrations/" + fmt.Sprint(version) + "/" + fileName.(string))
			if err != nil {
				fmt.Println("MIGRATION ERROR:", err)
				return
			}

			sql := string(fileBytes)

			//to handle WILDCARDS
			regex := regexp.MustCompile(`\$\{(.*?)\}`)
			keys := regex.FindAllStringSubmatch(sql, -1)
			for _, key := range keys {
				value := os.Getenv(key[1])
				sql = strings.Replace(sql, key[0], value, -1)
			}

			// fmt.Println(sql)

			// migrationSqlCommand := fmt.Sprintf("INSERT INTO public.migration_version_controls (version, description) VALUES ('%v', '%v');", version, index["description"])

			// sql = sql + "\n" + migrationSqlCommand

			tx := db.Exec(sql)
			if tx.Error != nil {
				fmt.Println("MIGRATION ERROR:", tx.Error)
				return
			}
		}
	}

	fmt.Println("Updated DB Version: ", upgradeVersion)
}
