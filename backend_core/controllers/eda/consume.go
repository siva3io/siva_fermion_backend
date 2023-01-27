package eda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

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

// This method configures the kafka reader and returns a kafka reader object
func ReaderConfig(topic TopicName, group GroupName) *kafka.Reader {
	// LOGGING := log.New(os.Stdout, "kafka reader: ", 1)
	topic = GetTechTopicName(topic)
	return ReaderConfigFunction(topic, group)
}

func OndcReaderConfig(topic TopicName, group GroupName) *kafka.Reader {
	// LOGGING := log.New(os.Stdout, "kafka reader: ", 1)
	topic = GetOndcTopicName(topic)
	return ReaderConfigFunction(topic, group)
}
func ReaderConfigFunction(topic TopicName, group GroupName) *kafka.Reader {
	receiver := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("TECH_KAFKA_BROKER_ADDRESS")},
		Topic:   string(topic),
		GroupID: string(group),
		//Logger:                LOGGING,
		// CommitInterval:        5 * time.Second,
		OffsetOutOfRangeError: true,
		StartOffset:           kafka.LastOffset,
		MinBytes:              10e3, // 10KB
		MaxBytes:              10e6, // 10MB
		HeartbeatInterval:     3 * time.Second,
		SessionTimeout:        10 * time.Second,
	})
	return receiver
}

func FetchReceiverPayload(receiver *kafka.Reader, ctx context.Context) (map[string]interface{}, error) {
	msg, err := receiver.ReadMessage(ctx)
	if err != nil {
		fmt.Println("could not read message related to "+msg.Topic, err)
		time.Sleep(1 * time.Second)
		return nil, err
	}
	var receiverData map[string]interface{}
	err = json.Unmarshal(msg.Value, &receiverData)
	if err != nil {
		fmt.Println("could not un-marshal message related to "+msg.Topic, err)
		time.Sleep(1 * time.Second)
		return nil, err
	}
	meta_data, ok := receiverData["meta_data"].(map[string]interface{})
	if ok {
		if meta_data["encryption"].(bool) {
			decrypted_payload, _ := helpers.AESGCMDecrypt(receiverData["data"].(string))
			var payload interface{}
			json.Unmarshal([]byte(decrypted_payload), &payload)
			receiverData["data"] = payload
		}
	} else {
		time.Sleep(1 * time.Second)
		return nil, errors.New("no meta data found")
	}

	return receiverData, nil
}

func ReaderConfigWithoutGroup(topic TopicName) *kafka.Reader {
	// LOGGING := log.New(os.Stdout, "kafka reader: ", 1)
	topic = GetTechTopicName(topic)
	receiver := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("TECH_KAFKA_BROKER_ADDRESS")},
		Topic:   string(topic),
		GroupID: string("BosonConvertorAck"),
		// Logger:                LOGGING,
		// CommitInterval:        5 * time.Second,
		OffsetOutOfRangeError: true,
		MinBytes:              10e3, // 10KB
		MaxBytes:              10e6, // 10MB
		HeartbeatInterval:     3 * time.Second,
		SessionTimeout:        10 * time.Second,
	})
	return receiver
}
