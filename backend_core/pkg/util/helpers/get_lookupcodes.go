package helpers

import (
	"encoding/json"

	"fermion/backend_core/db"
	model_core "fermion/backend_core/internal/model/core"

	"gorm.io/datatypes"
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
func GetLookupCodesFromArrOfObjectIdsToArrOfObjects(input datatypes.JSON) datatypes.JSON {

	db := db.DbManager()
	lookupCodes := make([]model_core.Lookupcode, 0)

	inputArray := make([]map[string]uint, 0)
	inputIds := make([]uint, 0)

	json.Unmarshal(input, &inputArray)
	if len(inputArray) == 0 {
		return nil
	}
	for _, val := range inputArray {
		inputIds = append(inputIds, val["id"])
	}

	db.Model(model_core.Lookupcode{}).Find(&lookupCodes, inputIds)

	output, _ := json.Marshal(lookupCodes)

	return output
}
