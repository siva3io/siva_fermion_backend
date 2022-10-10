package helpers

import (
	"errors"
	"fmt"
	"strings"

	"fermion/backend_core/db"
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
/*
	//--------------Function Call Example-------------------------------
	ID, err := helpers.GetLookupcodeId("LocAtion_TyPe", "retAIL")
	if err != nil {
		return err
	}
*/

func GetLookupcodeId(lookuptype string, lookupcode string) (uint, error) {
	var ID uint

	lookuptype = strings.ToUpper(lookuptype)
	lookupcode = strings.ToUpper(lookupcode)

	db := db.DbManager()

	//Fetch Lookuptype ID from LookupType table usong lookup type
	query := fmt.Sprintf("SELECT id FROM lookupcodes WHERE lookup_type_id = (SELECT id FROM lookuptypes WHERE lookup_type = '%v') AND lookup_code = '%v'", lookuptype, lookupcode)
	err := db.Raw(query).Scan(&ID)
	if ID == 0 || err.Error != nil {
		return 0, errors.New("lookup not found")
	}

	return ID, nil
}

func GetLookupcodeName(lookupId int) (string, error) {
	var LookupCode string

	db := db.DbManager()

	//Fetch Lookuptype Code from LookupType table usong lookup id
	query := fmt.Sprintf("SELECT lookup_code FROM lookupcodes WHERE id = '%v'", lookupId)
	db.Raw(query).Scan(&LookupCode)
	// if LookupCode ==  || err.Error != nil {
	// 	return 0, errors.New("lookup not found")
	// }

	return LookupCode, nil
}
