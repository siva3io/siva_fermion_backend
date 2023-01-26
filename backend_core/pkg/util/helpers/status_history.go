package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	model_core "fermion/backend_core/internal/model/core"

	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
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
func UpdateStatusHistory(data datatypes.JSON, status_id uint) ([]byte, error) {

	var status model_core.Lookupcode

	db := db.DbManager()
	err := db.Preload(clause.Associations).Where("id", status_id).First(&status)

	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if err.Error != nil {
		return nil, err.Error
	}

	var status_history []map[string]interface{}
	var result []byte
	var DB_TZ = os.Getenv("DB_TZ")

	if data != nil {
		err := json.Unmarshal(data, &status_history)
		if err != nil {
			return result, err
		}
	}
	loc, _ := time.LoadLocation(DB_TZ)
	status_map := map[string]interface{}{
		"status": status.DisplayName, "updated_at": time.Now().In(loc),
	}

	status_history = append([]map[string]interface{}{status_map}, status_history...)
	fmt.Println("----------------->", status_history)
	result, _ = json.Marshal(status_history)

	return result, nil

}
