package offers

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

func InitConsumers(h *handler) {
	// go CreateOffersConsumer(context.Background(), h)
	//go UpdateOffersConsumer(context.Background(), h)
	// go CreateOffersAck(context.Background(), h)
	// go UpdateOffersAck(context.Background(), h)
}

// func CreateOffersConsumer(ctx context.Context, h *handler) {
// 	receiver := eda.ReaderConfig(eda.CREATE_OFFER, eda.OFFERS)
// 	defer receiver.Close()
// 	for {
// 		receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
// 		if err != nil {
// 			continue
// 		}
// 		h.CreateOffers(receiverData)
// 	}
// }

// func UpdateOffersConsumer(ctx context.Context, h *handler) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_OFFER, eda.OFFERS)
// 	defer receiver.Close()
// 	for {
// 		receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
// 		if err != nil {
// 			continue
// 		}
// 		h.UpdateOffers(receiverData)
// 	}
// }
// func CreateOffersAck(ctx context.Context, h *handler) {
// 	receiver := eda.ReaderConfig(eda.CREATE_OFFER_ACK, eda.OFFERS)
// 	defer receiver.Close()
// 	for {
// 		receiverData, _ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_OFFER_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_OFFER_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpdateOffersAck(ctx context.Context, h *handler) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_OFFER_ACK, eda.OFFERS)
// 	defer receiver.Close()
// 	for {
// 		receiverData, _ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_OFFER_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_OFFER_ACK), receiverData["error_message"])
// 		}
// 	}
// }
