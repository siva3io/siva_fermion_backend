package purchase_invoice

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
	go CreatePurchanseInvoiceConsumer(h, KafkaConsumerTimeout)
	go UpdatePurchaseInvoiceConsumer(h, KafkaConsumerTimeout)
	// go CreatePurchaseInvoiceConsumerAck(context.Background(),h)
	// go UpdatePurchaseInvoiceConsumerAck(context.Background(),h)
}

func CreatePurchanseInvoiceConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_PURCHASE_INVOICE, eda.ACCOUNTING)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreatePurchaseInvoice(receiverData)
		}
	}
}

func UpdatePurchaseInvoiceConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_PURCHASE_INVOICE, eda.ACCOUNTING)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdatePurchaseInvoice(receiverData)
		}
	}
}

// func CreatePurchaseInvoiceConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_PURCHASE_INVOICE_ACK, eda.ACCOUNTING)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PURCHASE_INVOICE_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PURCHASE_INVOICE_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpdatePurchaseInvoiceConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_PURCHASE_INVOICE_ACK, eda.ACCOUNTING)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PURCHASE_INVOICE_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PURCHASE_INVOICE_ACK), receiverData["error_message"])
// 		}
// 	}
// }
