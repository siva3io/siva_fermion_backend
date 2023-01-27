package vendors

import (
	"context"
	"time"

	"fermion/backend_core/controllers/eda"
)

func InitConsumers(h *handler, KafkaConsumerTimeout time.Duration) {
	go CreateVendorConsumer(h, KafkaConsumerTimeout)
	go UpdateVendorConsumer(h, KafkaConsumerTimeout)
	go UpsertVendorConsumer(h, KafkaConsumerTimeout)
	// go CreateVendorConsumerAck(context.Background(),h)
	// go UpdateVendorConsumerAck(context.Background(),h)
	// go UpsertVendorConsumerAck(h, KafkaConsumerTimeout)

}

func CreateVendorConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_VENDOR, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateVendor(receiverData)
		}
	}
}

func UpdateVendorConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_VENDOR, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateVendor(receiverData)
		}
	}
}
func UpsertVendorConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPSERT_VENDOR, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpsertVendor(receiverData)
		}
	}
}

// func CreateVendorConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_VENDOR_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_VENDOR_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_VENDOR_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpdateVendorConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_VENDOR_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_VENDOR_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_VENDOR_ACK), receiverData["error_message"])
// 		}
// 	}
// }

// func UpsertVendorConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPSERT_VENDORS_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_VENDORS_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_VENDORS_ACK), receiverData["error_message"])
// 		}
// 	}
// }
