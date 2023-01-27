package omnichannel_fields

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"fermion/backend_core/controllers/cores"
	scheduler "fermion/backend_core/controllers/scheduler/app"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
	omnichannel_repo "fermion/backend_core/internal/repository/omnichannel"
	ipaas_core "fermion/backend_core/ipaas_core/app"

	"fermion/backend_core/pkg/util/helpers"
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
type Service interface {
	CreateOmnichannelField(data *OmnichannelFieldRequestDto, access_template_id string, token_id string) (uint, error)
	ListOmnichannelFields(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error)
	ViewOmnichannelField(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error)
	UpdateOmnichannelField(query map[string]interface{}, data *OmnichannelFieldRequestDto, access_template_id string, token_id string) error
	DeleteOmnichannelField(query map[string]interface{}, access_template_id string, token_id string) error

	ViewAppFields(query ViewAppFieldsQueryDto, accessToken string, access_template_id string, token_id string) (interface{}, error)
	ViewAppFieldsData(query ViewAppFieldsDataQueryDto, access_template_id string, token_id string) (interface{}, error)
	UpsertOmnichannelFieldsData(dtoData *OmnichannelFieldDataRequestDto, access_template_id string, token_id string, register bool) (interface{}, error)
	GetAppSyncSettings(app_code string, access_template_id string, token_id string) (interface{}, error)
	UpsertAppSyncSettings(appCode string, data []OmnichannelSyncSettingsRequestDto, access_template_id string, token_id string) error
	GetAppDataFilter(query GetAppDataFilterQueryDto, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	omnichannelFieldRepository omnichannel_repo.OmnichannelField
	coreService                cores.Service
	ipaasFeaturesService       ipaas_core.Service
	schedulerService           scheduler.Service
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		omnichannelFieldRepository: omnichannel_repo.NewOmnichannelFieldRepo(),
		coreService:                cores.NewService(),
		ipaasFeaturesService:       ipaas_core.NewService(),
		schedulerService:           scheduler.NewService(),
	}
	return newServiceObj
}

func (s *service) CreateOmnichannelField(dtoData *OmnichannelFieldRequestDto, access_template_id string, token_id string) (uint, error) {
	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return fmt.Errorf("you dont have access for create sales order at view level")
	// }
	// if data_access == nil {
	// 	return fmt.Errorf("you dont have access for create sales order at data level")
	// }

	data := new(omnichannel.OmnichannelField)
	err := helpers.JsonMarshaller(&dtoData, data)
	if err != nil {
		return 0, nil
	}
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	id, err := s.omnichannelFieldRepository.Save(data)
	if err != nil {
		return 0, nil
	}
	return id, nil

}

func (s *service) ListOmnichannelFields(page *pagination.Paginatevalue, access_template_id string, token_id string, access_action string) (interface{}, error) {
	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return nil, fmt.Errorf("you dont have access for list sales order at view level")
	// }
	// if data_access == nil {
	// 	return nil, fmt.Errorf("you dont have access for list sales order at data level")
	// }

	data, err := s.omnichannelFieldRepository.FindAll(page)
	if err != nil {
		return nil, err
	}
	dtoData := new([]OmnichannelFieldResponseDto)
	err = helpers.JsonMarshaller(data, dtoData)
	if err != nil {
		return nil, err
	}
	return dtoData, nil

}

func (s *service) ViewOmnichannelField(query map[string]interface{}, access_template_id string, token_id string) (interface{}, error) {
	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at view level")
	// }
	// if data_access == nil {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at data level")
	// }

	data, err := s.omnichannelFieldRepository.FindOne(query)
	if err != nil {
		return nil, err
	}
	dtoData := new(OmnichannelFieldResponseDto)
	err = helpers.JsonMarshaller(data, dtoData)
	if err != nil {
		return nil, err
	}
	return dtoData, nil

}

func (s *service) UpdateOmnichannelField(query map[string]interface{}, dtoData *OmnichannelFieldRequestDto, access_template_id string, token_id string) error {
	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return fmt.Errorf("you dont have access for update sales order at view level")
	// }
	// if data_access == nil {
	// 	return fmt.Errorf("you dont have access for update sales order at data level")
	// }

	data := new(omnichannel.OmnichannelField)
	err := helpers.JsonMarshaller(&dtoData, data)
	if err != nil {
		return nil
	}
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = s.omnichannelFieldRepository.Update(query, data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) DeleteOmnichannelField(query map[string]interface{}, access_template_id string, si string) error {
	// user_id := helpers.ConvertStringToUint(si)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "SALES_ORDERS", *user_id)
	// if !access_module_flag {
	// 	return fmt.Errorf("you dont have access for delete sales order at view level")
	// }
	// if data_access == nil {
	// 	return fmt.Errorf("you dont have access for delete sales order at data level")
	// }

	err := s.omnichannelFieldRepository.Delete(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ViewAppFields(query ViewAppFieldsQueryDto, accessToken string, access_template_id string, token_id string) (interface{}, error) {

	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at view level")
	// }
	// if data_access == nil {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at data level")
	// }

	page := new(pagination.Paginatevalue)
	page.Filters = "[[\"code\",\"=\",\"OMNICHANNEL\"]]"
	omnichannelApp, err := s.coreService.ListInstalledApps(page, false)
	if err != nil {
		return nil, err
	}

	page = new(pagination.Paginatevalue)

	filters := fmt.Sprintf("[[\"channel_type_id\",\"=\",%v],[\"channel_function_id\",\"=\",%v],[\"app_id\",\"IN\",[%v,%v]],[\"is_hidden\",\"=\",%v]]", query.ChannelTypeId, query.ChannelFunctionId, omnichannelApp[0].ID, query.AppId, false)

	page.Sort = "[[\"section_sequence\", \"asc\"],[\"field_sequence\", \"asc\"]]"
	page.Filters = filters
	page.Per_page = -1

	data, err := s.omnichannelFieldRepository.FindAll(page)
	if err != nil {
		return nil, err
	}

	dtoData := make([]OmnichannelFieldResponseDto, 0)
	err = helpers.JsonMarshaller(data, &dtoData)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	for index, field := range dtoData {

		if field.DataType == "dynamic" {
			wg.Add(1)
			go func(index int, field OmnichannelFieldResponseDto, dtoData []OmnichannelFieldResponseDto) {
				defer wg.Done()

				apiResponse, _ := helpers.MakeRequest(helpers.Request{
					Method: "GET",
					Scheme: "http",
					Host:   "localhost:3031",
					Path:   field.DataSource,
					Header: map[string]string{
						"Authorization": accessToken,
						"Content-Type":  "application/json",
					},
					Params: map[string]string{
						"per_page": "20",
					},
				})

				if apiResponse != nil {
					responseObj := apiResponse.(map[string]interface{})
					if responseObj["data"] != nil {
						if field.DataSourceParams == nil {
							dtoData[index].AllowedValues, _ = json.Marshal(responseObj["data"])
							return
						}
						paramObj := make(map[string]interface{})
						json.Unmarshal(field.DataSourceParams, &paramObj)
						dataArr := make([]map[string]interface{}, 0)
						helpers.JsonMarshaller(responseObj["data"], &dataArr)
						allowedValuesArr := make([]interface{}, len(dataArr))
						for index, value := range dataArr {
							allowedValuesArr[index] = map[string]interface{}{
								"id":    value["id"],
								"value": value[paramObj["key"].(string)],
							}
						}
						dtoData[index].AllowedValues, _ = json.Marshal(allowedValuesArr)
					}
				}

			}(index, field, dtoData)
		}

	}

	wg.Wait()

	return dtoData, nil
}

func (s *service) UpsertOmnichannelFieldsData(dtoData *OmnichannelFieldDataRequestDto, access_template_id string, token_id string, register bool) (interface{}, error) {
	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return fmt.Errorf("you dont have access for create sales order at view level")
	// }
	// if data_access == nil {
	// 	return fmt.Errorf("you dont have access for create sales order at data level")
	// }

	data := make([]omnichannel.OmnichannelFieldData, len(dtoData.Fields))
	err := helpers.JsonMarshaller(&dtoData.Fields, &data)
	if err != nil {
		return nil, err
	}

	appData, err := helpers.GetAppFromCache(dtoData.AppCode)
	if err != nil {
		return nil, err
	}

	appId := uint(appData["id"].(float64))

	for index := range data {
		data[index].FieldAppId = &appId
	}

	err = s.omnichannelFieldRepository.UpsertOmnichannelFieldsData(&data)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"redirect": false,
		"url":      nil,
		"message":  nil,
	}

	if !register {
		return response, nil
	}

	apiResponse, err := helpers.MakeRequest(helpers.Request{
		Method: "GET",
		Scheme: "http",
		Host:   "localhost:3031",
		Path:   "integrations/" + strings.ToLower(dtoData.AppCode) + "/register",
		Header: map[string]string{
			"Authorization": appData["access_token"].(string),
		},
	})
	if err != nil {
		return nil, err
	}

	if apiResponse.(map[string]interface{})["data"] == nil {
		return apiResponse.(map[string]interface{})["meta"], nil
	}

	return apiResponse.(map[string]interface{})["data"], nil
}

func (s *service) GetAppSyncSettings(app_code string, access_template_id string, token_id string) (interface{}, error) {

	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at view level")
	// }
	// if data_access == nil {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at data level")
	// }

	query := map[string]interface{}{
		"app_code":    app_code,
		"is_syncable": true,
	}

	data, err := s.ipaasFeaturesService.GetAppFeatures(query)
	if err != nil {
		return nil, err
	}

	var dtoData []OmnichannelSyncSettingsResponseDto

	err = helpers.JsonMarshaller(data, &dtoData)
	if err != nil {
		return nil, err
	}

	for index, syncSetting := range dtoData {
		query := map[string]interface{}{
			"source_type": "features",
			"source_id":   syncSetting.ID,
		}
		schedulerJob, err := s.schedulerService.GetOneSchedulerJob(query)
		if err != nil {
			return nil, err
		}
		if schedulerJob.ID != 0 {
			dtoData[index].JobId = schedulerJob.ID
			dtoData[index].State = *schedulerJob.State
			dtoData[index].Frequency = schedulerJob.Frequency
		}
	}

	query = map[string]interface{}{
		"lookup_type": "SCHEDULER_FREQUENCY",
	}
	var page pagination.Paginatevalue
	frequencyValues, err := s.coreService.GetLookupCodes(query, &page)
	if err != nil {
		return nil, err
	}
	dtoDataWithDropDown := OmnichannelSyncSettingsResponseDtoWithDropDownDto{
		SyncSettings:      dtoData,
		FrequencyDropDown: frequencyValues,
	}

	return dtoDataWithDropDown, nil
}

func (s *service) UpsertAppSyncSettings(appCode string, data []OmnichannelSyncSettingsRequestDto, access_template_id string, token_id string) error {

	for _, syncSetting := range data {
		query := map[string]interface{}{
			"source_type": "features",
			"source_id":   syncSetting.ID,
		}
		schedulerJob, _ := s.schedulerService.GetOneSchedulerJob(query)

		user_id, _ := strconv.Atoi(token_id)

		if schedulerJob.ID != 0 {
			query := map[string]interface{}{
				"id":      schedulerJob.ID,
				"user_id": user_id,
			}

			newSchedulerJob := scheduler.SchedulerJobDto{
				Name:      syncSetting.Name,
				Frequency: syncSetting.Frequency,
				State:     &syncSetting.State,
				Params:    schedulerJob.Params,
			}

			err := s.schedulerService.UpdateSchedulerJob(query, &newSchedulerJob)
			if err != nil {
				return err
			}
			continue
		}

		if !syncSetting.State {
			continue
		}

		newSchedulerJob := scheduler.SchedulerJobDto{
			Name:       syncSetting.Name,
			Function:   "ApiCall",
			State:      &syncSetting.State,
			Params:     []byte{},
			Frequency:  syncSetting.Frequency,
			SourceType: "features",
			SourceId:   syncSetting.ID,
		}

		requestOption := helpers.Request{
			Method: "GET",
			Scheme: "http",
			Host:   "localhost:3031",
			Path:   "/integrations/" + strings.ToLower(appCode) + "/features/" + syncSetting.ID,
		}
		newSchedulerJob.Params, _ = json.Marshal(requestOption)

		err := s.schedulerService.AddSchedulerJob(newSchedulerJob)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) ViewAppFieldsData(query ViewAppFieldsDataQueryDto, access_template_id string, token_id string) (interface{}, error) {

	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at view level")
	// }
	// if data_access == nil {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at data level")
	// }

	newQuery := map[string]interface{}{
		"omnichannel_fields.channel_function_id": query.ChannelFunctionId,
		"omnichannel_field_data.field_app_id":    query.FieldAppId,
		"omnichannel_fields.is_hidden":           false,
	}

	data, err := s.omnichannelFieldRepository.GetOmnichannelFieldsData(newQuery)
	if err != nil {
		return nil, err
	}

	dtoData := new([]OmnichannelFieldDataResponseDto)
	err = helpers.JsonMarshaller(data, dtoData)
	if err != nil {
		return nil, err
	}

	return dtoData, nil
}

func (s *service) GetAppDataFilter(query GetAppDataFilterQueryDto, access_template_id string, token_id string) (interface{}, error) {

	// token_user_id := helpers.ConvertStringToUint(token_id)
	// access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	// if !access_module_flag {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at view level")
	// }
	// if data_access == nil {
	// 	return nil, fmt.Errorf("you dont have access for view sales order at data level")
	// }
	page := new(pagination.Paginatevalue)
	page.Filters = "[[\"code\",\"=\",\"OMNICHANNEL\"]]"
	omnichannelApp, err := s.coreService.ListInstalledApps(page, false)
	if err != nil {
		return nil, err
	}

	page = new(pagination.Paginatevalue)
	page.Filters = fmt.Sprintf("[[\"name\",\"IN\",%v],[\"app_id\",\"IN\",[%v,%v]]]", query.Fields, query.AppId, omnichannelApp[0].ID)
	page.Per_page = -1

	data, err := s.omnichannelFieldRepository.FindAll(page)
	if err != nil {
		return nil, err
	}

	dtoData := make([]AppDataFilterResponseDto, 0)
	err = helpers.JsonMarshaller(data, &dtoData)
	if err != nil {
		return nil, err
	}

	objResponse := make(map[string]interface{})

	for _, field := range data.([]omnichannel.OmnichannelField) {
		page.Filters = fmt.Sprintf("[[\"omnichannel_field_id\",\"=\",%v],[\"field_app_id\",\"=\",%v]]", field.ID, query.AppId)
		page.Per_page = -1

		data, err := s.omnichannelFieldRepository.FindAllFieldData(page)
		if err != nil {
			return nil, err
		}
		modelData := data.([]omnichannel.OmnichannelFieldData)
		var fieldValue interface{}
		if len(modelData) > 0 {
			fieldValue = data.([]omnichannel.OmnichannelFieldData)[0].Data
		}
		objResponse[field.Name] = map[string]interface{}{
			"id":    field.ID,
			"value": fieldValue,
		}
	}

	return objResponse, nil
}
