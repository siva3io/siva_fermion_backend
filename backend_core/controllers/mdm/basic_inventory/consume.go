package basic_inventory

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

func InitConsumers(h *handler, KafkaConsumerTimeout time.Duration) {
	go CreateDecentralizedInventoryConsumer(h, KafkaConsumerTimeout)
	// go CreateDecentralizedInventoryConsumerAck(context.Background(),h)

	go UpdateDecentralizedInventoryConsumer(h, KafkaConsumerTimeout)
	// go UpdateDecentralizedInventoryConsumerAck(context.Background(),h)

	go UpsertDecentralizedInventoryConsumer(h, KafkaConsumerTimeout)
	// go UpsertDecentralizedInventoryConsumerAck(context.Background(),h)

	go CreateCentralizedInventoryConsumer(h, KafkaConsumerTimeout)
	// go CreateCentralizedInventoryConsumerAck(context.Background(),h)

	go UpdateCentralizedInventoryConsumer(h, KafkaConsumerTimeout)
	// go UpdateCentralizedInventoryConsumerAck(context.Background(),h)
}

func CreateDecentralizedInventoryConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_INVENTORY, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateDecentralizedInventory(receiverData)
		}
	}
}

func UpdateDecentralizedInventoryConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_INVENTORY, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateDecentralizedInventory(receiverData)
		}
	}
}

func UpsertDecentralizedInventoryConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPSERT_DECENTRALIZED_INVENTORY, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.DeCentralisedInventoryUpsert(receiverData)
		}
	}
}

// func CreateDecentralizedInventoryConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_INVENTORY_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_INVENTORY_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_INVENTORY_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpdateDecentralizedInventoryConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_INVENTORY_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_INVENTORY_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_INVENTORY_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpsertDecentralizedInventoryConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPSERT_DECENTRALIZED_INVENTORY_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_DECENTRALIZED_INVENTORY_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_DECENTRALIZED_INVENTORY_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func CreateCentralizedInventoryConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_CENTRALIZED_INVENTORY, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateCentrailizedInventory(receiverData)
		}
	}
}

func UpdateCentralizedInventoryConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_CENTRALIZED_INVENTORY, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateCentrailizedInventory(receiverData)
		}
	}
}

// func CreateCentralizedInventoryConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_CENTRALIZED_INVENTORY_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_CENTRALIZED_INVENTORY_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_CENTRALIZED_INVENTORY_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpdateCentralizedInventoryConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_CENTRALIZED_INVENTORY_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_CENTRALIZED_INVENTORY_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_CENTRALIZED_INVENTORY_ACK), receiverData["error_message"])
// 		}
// 	}
// }
