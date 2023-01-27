package contacts

import (
	"context"
	"time"

	"fermion/backend_core/controllers/eda"
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

// const GROUP_ID = "my-group"

func InitConsumers(h *handler, KafkaConsumerTimeout time.Duration) {
	go CreateContactConsumer(h, KafkaConsumerTimeout)
	go UpdateContactConsumer(h, KafkaConsumerTimeout)
	go UpsertContactConsumer(h, KafkaConsumerTimeout)
	// go CreateContactConsumerAck(context.Background(),h)
	// go UpdateContactConsumerAck(context.Background(),h)
	// go UpsertContactConsumerAck(context.Background(),h)
}

func CreateContactConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_CONTACTS, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateContact(receiverData)
		}
	}
}

func UpdateContactConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_CONTACTS, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateContact(receiverData)
		}
	}
}

func UpsertContactConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPSERT_CONTACT, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpsertContact(receiverData)
		}
	}
}

// func CreateContactConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_CONTACTS_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_CONTACTS_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_CONTACTS_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpdateContactConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_CONTACTS_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_CONTACTS_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_CONTACTS_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpsertContactConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPSERT_CONTACTS_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_CONTACTS_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_CONTACTS_ACK), receiverData["error_message"])
// 		}
// 	}
// }
