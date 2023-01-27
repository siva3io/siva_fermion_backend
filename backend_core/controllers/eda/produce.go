package eda


import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"fermion/backend_core/pkg/util/helpers"

	kafka "github.com/segmentio/kafka-go"
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


func Produce(topic TopicName, input_data interface{}, optionalParams ...interface{}) {
	topic = GetTechTopicName(topic)
	ProduceFunction(topic, input_data, optionalParams...)
}

func OndcProduce(topic TopicName, input_data interface{}, optionalParams ...interface{}) {
	topic = GetOndcTopicName(topic)
	ProduceFunction(topic, input_data, optionalParams...)
}

// Produce a set of data under a specific topic with consumer group name to the kafka server
func ProduceFunction(topic TopicName, input_data interface{}, optionalParams ...interface{}) {
	// id := uuid.New()
	var data map[string]interface{}
	helpers.JsonMarshaller(input_data, &data)
	//==========================checking whether the data need to encrypt or not================
	isDataEncryption := true
	var ok bool

	if len(optionalParams) > 0 {
		isDataEncryption, ok = optionalParams[0].(bool)
		if !ok {
			fmt.Println("could not write message - error in optionalParams")
		}
	}

	//=============================writer configuration====================================
	l := log.New(os.Stdout, "kafka writer: ", 0)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{os.Getenv("TECH_KAFKA_BROKER_ADDRESS")},
		Topic:    string(topic),
		Balancer: &kafka.LeastBytes{},
		// CompressionCodec: &compress.ZstdCodec,
		Logger: l,
	})
	fmt.Println("NewWriter: ", w)
	//==================if encryption is true, then encrypt the data =================
	if isDataEncryption {
		payload := data["data"]
		payload_bytes, _ := json.Marshal(payload)
		encrypted_payload, _ := helpers.AESGCMEncrypt(string(payload_bytes))
		data["data"] = encrypted_payload
	}
	//==================add the encryption field in meta_data==========================
	if data["meta_data"] != nil {
		metaData := data["meta_data"].(map[string]interface{})
		metaData["encryption"] = isDataEncryption
		data["meta_data"] = metaData
		fmt.Println(data["meta_data"])
	} else {
		data["meta_data"] = map[string]interface{}{
			"encryption": isDataEncryption,
		}
	}
	data_bytes, _ := json.Marshal(data)
	//==================write the message================================================
	err := w.WriteMessages(context.TODO(), kafka.Message{
		Key:   []byte("Tech-Key"),
		Value: data_bytes,
	})
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
