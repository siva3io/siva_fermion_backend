package ondc

import (
	"context"
	"time"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/pkg/util/cache"
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
	go UpsertOndcProductConsumerAck(h, KafkaConsumerTimeout)
	go UpsertOndcProviderConsumerAck(h, KafkaConsumerTimeout)
}

func UpsertOndcProductConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.BECKN_PRODUCTS_SYNC_ACK, eda.BECKN)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			cache.SetCacheVariable(receiverData["meta_data"].(map[string]interface{})["request_id"].(string), receiverData)
		}
	}
}

func UpsertOndcProviderConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
	receiver := eda.ReaderConfig(eda.BECKN_PROVIDER_SYNC_ACK, eda.BECKN)
	defer receiver.Close()
	for {
		receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
		if err != nil {
			receiver.Close()
			cancel()
			break
		}
		cache.SetCacheVariable(receiverData["meta_data"].(map[string]interface{})["request_id"].(string), receiverData)
	}
}
