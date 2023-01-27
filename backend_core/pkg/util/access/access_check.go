package access_check

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/pagination"
	access_repo "fermion/backend_core/internal/repository/access"
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
func ValidateUserAccess(access_template_id string, module_access_action string, module_access_name string, user_id uint) (bool, interface{}) {

	id, _ := strconv.Atoi(access_template_id)

	// fmt.Println("AccesstemplateId", id)

	p := new(pagination.Paginatevalue)

	p.Filters = "[[\"display_name\",\"=\",\"" + module_access_name + "\"]]"

	modules_data, _ := access_repo.NewModule().FindByidModule(uint(id), p)

	for _, module_data := range modules_data {
		if module_data.DisplayName == module_access_name {

			var view_actions []map[string]interface{}
			if module_data.ViewActions != nil {
				json.Unmarshal(module_data.ViewActions, &view_actions)
			}
			var data_actions []map[string]interface{}
			if module_data.ViewActions != nil {
				json.Unmarshal(module_data.DataActions, &data_actions)
			}
			viewActionsAccess, _ := ViewActionsAccess(module_access_action, view_actions)

			if viewActionsAccess {
				dataActionsAccess, _ := DataActionsAccess(module_access_action, data_actions, user_id)
				return viewActionsAccess, dataActionsAccess
			}
			return viewActionsAccess, nil
		}
	}
	return false, nil
}

func ViewActionsAccess(module_access_action string, actions []map[string]interface{}) (bool, error) {
	for _, action := range actions {

		ctrl_flag := fmt.Sprint(action["ctrl_flag"])

		if (action["lookup_code"] == "ADMIN" || action["lookup_code"] == module_access_action) && ctrl_flag == "1" {
			return true, nil
		}
	}
	return false, nil
}

func DataActionsAccess(module_access_action string, actions []map[string]interface{}, id uint) (interface{}, error) {
	for _, action := range actions {

		ctrl_flag := fmt.Sprint(action["ctrl_flag"])
		//add comapny query in default to all usesrs for the level of access

		if action["lookup_code"] == "ADMIN" || action["lookup_code"] == module_access_action {
			if ctrl_flag == "1" {
				return "", nil
			} else if ctrl_flag == "2" {

				return [1][3]string{{"created_by", "=", fmt.Sprintf("%d", id)}}, nil
			} else if ctrl_flag == "3" {
				//add the return data here based on the sub users creation with the created by or parent id.
				return [1][3]string{{"created_by", "=", fmt.Sprintf("%d", id)}}, nil
			} else if ctrl_flag == "4" {
				return action["query"], nil
			}
		}
	}
	return nil, nil
}
