package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	// "os/exec"
	"encoding/json"
	"path/filepath"

	model_core "fermion/backend_core/internal/model/core"

	"gorm.io/gorm"
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
// SqlModel declaring a struct
type SqlModel struct {

	// defining struct variables
	ModelName string
	Link      string
	Format    bool
}

func SeedMetaTable(db *gorm.DB, models *string) bool {
	ReflectModels(db, models)
	return true
}

// ReflectModels Reflect the given models in IRModel.
func ReflectModels(db *gorm.DB, models *string) {
	// var data model.IRModel
	modelNames := strings.Split(*models, ",")
	if modelNames[0] == "all" {
		db.Raw("SELECT tablename FROM pg_catalog.pg_tables where schemaname='public'").Scan(&modelNames)
		//db.Raw("SELECT schemaname as schema_name,tablename as table_name FROM pg_catalog.pg_tables where schemaname='public' AND tablename IN " + fmt.Sprintf("('%s')", strings.Join(modelNames, "','"))).Scan(&data)
	}
	for _, m := range modelNames {
		data := new(model_core.IRModel)
		db.Raw("SELECT schemaname as schema_name,tablename as table_name FROM pg_catalog.pg_tables where schemaname='public' AND tablename = ?", m).Scan(&data)
		if db.Model(&model_core.IRModel{}).Where("table_name = ?", m).Updates(&data).RowsAffected == 0 {
			db.Model(&model_core.IRModel{}).Create(&data)
		}
		createdData := new(model_core.IRModel)
		db.Model(&model_core.IRModel{}).Where("table_name = ?", data.TableName).First(&createdData)
		if createdData.ID != 0 {
			ReflectModelFields(db, m, createdData.ID)
			ReflectModelConstraints(db, m, createdData.ID)
		}
	}

}

// ReflectModelFields Reflecting the fields and constraints of given models in IRModelFields and IRModelConstraints
func ReflectModelFields(db *gorm.DB, m string, id uint) {
	var fields []model_core.IRModelFields
	db.Raw("SELECT column_name,table_name,data_type,character_maximum_length,is_nullable,column_default FROM information_schema.columns WHERE table_name = ?", m).Scan(&fields)
	for i := range fields {
		fields[i].TableId = id
	}
	if db.Model(&model_core.IRModelFields{}).Where("table_name = ?", m).Updates(&fields).RowsAffected == 0 {
		db.Model(&model_core.IRModelFields{}).Create(&fields)
	}
	//db.Model(&model.IRModelFields{}).Create(&fields)

	//db.Raw("select schemaname,tablename,column_name,is_nullable,data_type,udt_name,conname,contype from pg_catalog.pg_tables inner join information_schema.columns on information_schema.columns.table_schema=pg_catalog.pg_tables.schemaname INNER JOIN pg_catalog.pg_class rel ON rel.relname = tablename INNER JOIN pg_catalog.pg_constraint con ON con.conrelid = rel.oid INNER JOIN pg_catalog.pg_namespace nsp ON nsp.oid = connamespace where schemaname='public' order by tablename").Scan(&data)
}

func ReflectModelConstraints(db *gorm.DB, m string, id uint) {
	var constraints []model_core.IRModelConstraints
	db.Raw("SELECT constraint_name,table_name,constraint_type FROM information_schema.table_constraints WHERE table_name = ?", m).Scan(&constraints)
	for i := range constraints {
		constraints[i].TableId = id
	}
	if db.Model(&model_core.IRModelConstraints{}).Where("table_name = ?", m).Updates(&constraints).RowsAffected == 0 {
		db.Model(&model_core.IRModelConstraints{}).Create(&constraints)
	}
}

func SeedMasterData(db *gorm.DB) bool {
	MasterModelData := []byte(`
		[
			{
				"ModelName": "lookuptypes",
				"Link":"lookuptypes.sql"
			},
			{
				"ModelName": "lookupcodes",
				"Link":"lookupcodes.sql"
			},	
			{
				"ModelName": "core_users",
				"Link":"core_users.sql"
			},
			{
				"ModelName": "ondc_details",
				"Link":"ondc_details.sql"
			},
			{
				"ModelName": "companies",
				"Link":"companies.sql"
			},
			{
				"ModelName": "countries",
				"Link":"countries.sql"
			},
			{
				"ModelName": "states",
				"Link":"states.sql"
			},
			{
				"ModelName": "currencies",
				"Link":"currencies.sql"
			},
			{
				"ModelName": "installed_apps",
				"Link":"app.sql"
			},
			{
				"ModelName": "shipping_partners",
				"Link":"shipping_partners.sql"
			},
			{
				"ModelName": "app_categories",
				"Link":"app_categories.sql"
			},
			{
				"ModelName": "app_stores",
				"Link":"app_store.sql"
			},
			{
				"ModelName": "product_brands",
				"Link":"product_brands.sql"
			},
			{
				"ModelName": "core_app_modules",
				"Link":"core_app_modules.sql",
				"Format": true
			},
			{
				"ModelName": "access_templates",
				"Link":"access_templates.sql"
			},
			{
				"ModelName": "access_module_actions",
				"Link":"access_module_actions.sql"
			},
			{
				"ModelName": "user_template_mapper",
				"Link":"user_template_mapper.sql"
			},
			{
				"ModelName": "omnichannel_fields",
				"Link":"omnichannel_fields.sql"
			},
			{
				"ModelName": "hsn_codes_data",
				"Link":"hsn_codes_data.sql"
			}
		]
	`)
	fmt.Println("seeding master data into database...")
	var sqlModel []SqlModel

	err := json.Unmarshal(MasterModelData, &sqlModel)
	if err != nil {
		// if error is not nil
		// print error
		fmt.Println(err)
	}
	for i, model := range sqlModel {
		fmt.Println("seeding " + sqlModel[i].ModelName)
		path := filepath.Join("./backend_core/internal/model/data", sqlModel[i].Link)
		c, ioErr := ioutil.ReadFile(path)
		if ioErr != nil {
			// handle error.
			fmt.Println(ioErr)
		}
		sql := string(c)

		if model.Format {
			regex := regexp.MustCompile(`\$\{(.*?)\}`)
			keys := regex.FindAllStringSubmatch(sql, -1)

			for _, key := range keys {
				value := os.Getenv(key[1])
				sql = strings.Replace(sql, key[0], value, -1)
			}
		}

		db.Exec(sql)
		fmt.Println("seeded " + sqlModel[i].ModelName)
	}

	fmt.Println("master data seeded into database...")

	return false
}

func SeedTestData(db *gorm.DB) bool {
	MasterModelData := []byte(`
		[
			{
				"ModelName": "accounting"
			},
			{
				"ModelName": "inventory_orders"
			},
			{
				"ModelName": "inventory_tasks"
			},
			{
				"ModelName": "mdm"
			},
			{
				"ModelName": "omnichannel"
			},
			{
				"ModelName": "orders"
			},
			{
				"ModelName": "returns"
			},
			{
				"ModelName": "shipping"
			},
			{
				"ModelName": "rating"
			},
			{
				"ModelName": "offers"
			}
		]
	`)
	fmt.Println("seeding test data into database...")
	var sqlModel []SqlModel

	err := json.Unmarshal(MasterModelData, &sqlModel)
	if err != nil {
		// if error is not nil
		// print error
		fmt.Println(err)
	}
	for i := range sqlModel {
		fmt.Println("seeding " + sqlModel[i].ModelName)
		path := filepath.Join("./backend_core/internal/model/", sqlModel[i].ModelName, "/test_data.sql")
		c, ioErr := ioutil.ReadFile(path)
		if ioErr != nil {
			// handle error.
			fmt.Println(ioErr)
		}
		sql := string(c)
		db.Exec(sql)
		fmt.Println("seeded " + sqlModel[i].ModelName)
	}

	fmt.Println("test data seeded into database...")

	return false
}
