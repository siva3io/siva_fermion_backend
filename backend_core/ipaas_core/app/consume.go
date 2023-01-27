package ipaas_core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	eda "fermion/backend_core/controllers/eda"
	ipaas_core_model "fermion/backend_core/ipaas_core/model"
	ipaas_mapper "fermion/backend_core/ipaas_core/repository"
	cache "fermion/backend_core/pkg/util/cache"
	pkg_helpers "fermion/backend_core/pkg/util/helpers"
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
func InitConsumers(KafkaConsumerTimeout time.Duration) {
	go BosonConvertor(KafkaConsumerTimeout)
	go BosonConvertorAck(KafkaConsumerTimeout)
}

type BosonConvertorDTO struct {
	MetaData map[string]interface{} `json:"meta_data"`
	Data     interface{}            `json:"data"`
}

type BosonConvertorDataInput struct {
	MapperTemplate          ipaas_core_model.Mapper         `json:"mapper_template"`
	InputData               interface{}                     `json:"input_data"`
	InputType               string                          `json:"input_type"`
	FeatureSessionVariables []ipaas_core_model.KeyValuePair `json:"feature_session_variables"`
}

type BosonConvertorDataOuput struct {
	ErrorMessage   error       `json:"error_message"`
	MappedResponse interface{} `json:"mapped_response"`
}

func BosonConvertor(KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.BOSON_MAPPER_CONVERTOR, eda.BOSON_CONVERTOR)
		for {
			var NewBosonConvertor BosonConvertorDataInput
			var BosonConvertorResponse BosonConvertorDataOuput
			var input BosonConvertorDTO

			//===================convert the receiver messsage to BosonConvertorDTO format================================
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			err = pkg_helpers.JsonMarshaller(receiverData, &input)
			if err != nil {
				BosonConvertorResponse.ErrorMessage = err
				ProduceBosonConvertorResponse(input.MetaData, BosonConvertorResponse)
				continue
			}

			//=====================convert the  data to boson_convertor struct format==================================================================

			err = pkg_helpers.JsonMarshaller(input.Data, &NewBosonConvertor)
			if err != nil {
				BosonConvertorResponse.ErrorMessage = err
				ProduceBosonConvertorResponse(input.MetaData, BosonConvertorResponse)
				continue
			}

			//=================find the payload input_type======================================================
			NewBosonConvertor.InputType = "object"
			_, ok := NewBosonConvertor.InputData.([]interface{})
			if ok {
				NewBosonConvertor.InputType = "array"
			}

			//=================execute object mapper=======================================================
			if NewBosonConvertor.InputType == "object" {
				objectMappedReponse, err := ipaas_mapper.ExecuteObjMapper(NewBosonConvertor.InputData, NewBosonConvertor.MapperTemplate, NewBosonConvertor.FeatureSessionVariables)
				if err != nil {
					BosonConvertorResponse.ErrorMessage = err
				}
				BosonConvertorResponse.MappedResponse = objectMappedReponse
				ProduceBosonConvertorResponse(input.MetaData, BosonConvertorResponse)
				continue
			}

			//=================execute array mapper=======================================================
			if NewBosonConvertor.InputType == "array" {
				var mapperInput []interface{}

				err = pkg_helpers.JsonMarshaller(NewBosonConvertor.InputData, &mapperInput)
				if err != nil {
					BosonConvertorResponse.ErrorMessage = err
					ProduceBosonConvertorResponse(input.MetaData, BosonConvertorResponse)
					continue
				}
				arrayMappedResponse, err := ipaas_mapper.ExecuteArrayMapper(mapperInput, NewBosonConvertor.MapperTemplate, NewBosonConvertor.FeatureSessionVariables)
				if err != nil {
					BosonConvertorResponse.ErrorMessage = err
				}
				BosonConvertorResponse.MappedResponse = arrayMappedResponse
				ProduceBosonConvertorResponse(input.MetaData, BosonConvertorResponse)
				continue
			}
			BosonConvertorResponse.ErrorMessage = errors.New("input_type must be object or array")
			ProduceBosonConvertorResponse(input.MetaData, BosonConvertorResponse)
		}
	}
}
func ProduceBosonConvertorResponse(metaData map[string]interface{}, bosonConvertorResponse BosonConvertorDataOuput) {
	var response map[string]interface{}
	var output BosonConvertorDTO

	output.MetaData = metaData

	output.Data = bosonConvertorResponse

	pkg_helpers.JsonMarshaller(output, &response)

	if bosonConvertorResponse.ErrorMessage != nil {
		response["data"].(map[string]interface{})["error_message"] = bosonConvertorResponse.ErrorMessage.Error()
	}

	eda.Produce(eda.BOSON_MAPPER_CONVERTOR_ACK, response, metaData["encryption"].(bool))
}
func BosonConvertorAck(KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfigWithoutGroup(eda.BOSON_MAPPER_CONVERTOR_ACK)
		for {
			msg, err := receiver.ReadMessage(ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}

			var receiverData BosonConvertorDTO
			err = json.Unmarshal(msg.Value, &receiverData)
			if err != nil {
				fmt.Println("could not un-marshal message related to BsonConvertorConsumer")
			}

			request_id, ok := receiverData.MetaData["request_id"].(string)
			if ok {
				if request_id != "" {
					cache.SetCacheVariable(request_id, receiverData)
				}
			}
		}
	}
}
