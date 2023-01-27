package scheduler

import (
	"encoding/json"
	"fmt"
	"reflect"

	model "fermion/backend_core/controllers/scheduler/model"
	helpers "fermion/backend_core/pkg/util/helpers"

	"gorm.io/datatypes"
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
type CronFunctions struct {
}

func (c CronFunctions) FuncOne(param datatypes.JSON) model.SchedulerLog {
	var log model.SchedulerLog
	log.State = "success"
	fmt.Println("FuncOne", param)
	return log
}

func (c CronFunctions) FuncTwo(param datatypes.JSON) model.SchedulerLog {
	var log model.SchedulerLog
	log.State = "success"
	fmt.Println("FuncTwo", param)
	return log
}

func (c CronFunctions) ApiCall(param datatypes.JSON) model.SchedulerLog {

	var requestOption helpers.Request
	var log model.SchedulerLog

	err := json.Unmarshal(param, &requestOption)

	if err != nil {
		log.State = "failed"
		fmt.Println("-----Err-------", err)
		return log
	}
	response, err := helpers.MakeRequest(requestOption)

	if err != nil {
		log.State = "failed"
		fmt.Println("-----Err-------", err)
		return log
	}

	log.Notes, err = json.Marshal(response)

	if err != nil {
		log.State = "failed"
		fmt.Println("-----Err-------", err)
		return log
	}

	log.State = "success"
	return log
}

func CallFuncByName(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("method not found \"%s\"", funcName)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	out = m.Call(in)
	return
}
