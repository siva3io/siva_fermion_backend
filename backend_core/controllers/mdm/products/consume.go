package products

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
	go CreateProductConsumer(h, KafkaConsumerTimeout)
	// go CreateProductConsumerAck(h, KafkaConsumerTimeout)

	go UpdateProductConsumer(h, KafkaConsumerTimeout)
	// go UpdateProductConsumerAck(h, KafkaConsumerTimeout)

	go CreateProductBundleConsumer(h, KafkaConsumerTimeout)
	// go CreateProductBundleConsumerAck(h, KafkaConsumerTimeout)

	go UpdateProductBundleConsumer(h, KafkaConsumerTimeout)
	// go UpdateProductBundleConsumerAck(h, KafkaConsumerTimeout)

	go CreateProductCategoryConsumer(h, KafkaConsumerTimeout)
	// go CreateProductCategoryConsumerAck(h, KafkaConsumerTimeout)

	go UpdateProductCategoryConsumer(h, KafkaConsumerTimeout)
	// go UpdateProductCategoryConsumerAck(h, KafkaConsumerTimeout)

	go CreateProductBrandConsumer(h, KafkaConsumerTimeout)
	// go CreateProductBrandConsumerAck(h, KafkaConsumerTimeout)

	go UpdateCustomerBrandConsumer(h, KafkaConsumerTimeout)
	// go UpdateCustomerBrandConsumerAck(h, KafkaConsumerTimeout)

	go CreateProductVariantConsumer(h, KafkaConsumerTimeout)
	// go CreateProductVariantConsumerAck(h, KafkaConsumerTimeout)

	go UpdateCustomerVariantConsumer(h, KafkaConsumerTimeout)
	// go UpdateCustomerVariantConsumerAck(h, KafkaConsumerTimeout)

	go UpsertProductTemplateConsumer(h, KafkaConsumerTimeout)
	// go UpsertProductTemplateConsumerAck(h, KafkaConsumerTimeout)

	go UpsertProductVariantConsumer(h, KafkaConsumerTimeout)
	// go UpsertProductVariantConsumerAck(h, KafkaConsumerTimeout)

}

func CreateProductConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_PRODUCT, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateProductDetails(receiverData)
		}
	}
}

// func CreateProductConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func UpdateProductConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateProduct(receiverData)
		}
	}
}

// func UpdateProductConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func CreateProductBrandConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_BRAND, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateProductBrand(receiverData)
		}
	}
}

// func CreateProductBrandConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_BRAND_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_BRAND_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_BRAND_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func UpdateCustomerBrandConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_BRAND, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateProductBrand(receiverData)
		}
	}
}

// func UpdateCustomerBrandConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_BRAND_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_BRAND_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_BRAND_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func CreateProductVariantConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_VARIANT, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateProductVariant(receiverData)
		}
	}
}

// func CreateProductVariantConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_VARIANT_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_VARIANT_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_VARIANT_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func UpdateCustomerVariantConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_VARIANT, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateProductVariant(receiverData)
		}
	}
}

// func UpdateCustomerVariantConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_VARIANT_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_VARIANT_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_VARIANT_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func CreateProductBundleConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_BUNDLE, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateBundles(receiverData)
		}
	}
}

// func CreateProductBundleConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_BUNDLE_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_BUNDLE_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_BUNDLE_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func UpdateProductBundleConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_BUNDLE, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateBundle(receiverData)
		}
	}
}

// func UpdateProductBundleConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_BUNDLE_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_BUNDLE_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_BUNDLE_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func CreateProductCategoryConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_CATEGORY, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateProductCategory(receiverData)
		}
	}
}

// func CreateProductCategoryConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.CREATE_PRODUCT_CATEGORY_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_CATEGORY_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.CREATE_PRODUCT_CATEGORY_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func UpdateProductCategoryConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_CATEGORY, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateProductCategory(receiverData)
		}
	}
}

// func UpdateProductCategoryConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPDATE_PRODUCT_CATEGORY_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_CATEGORY_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPDATE_PRODUCT_CATEGORY_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func UpsertProductTemplateConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPSERT_PRODUCT_TEMPLATE, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.ChannelProductUpsert(receiverData)
		}
	}
}

// func UpsertProductTemplateConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPSERT_PRODUCT_TEMPLATE_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_PRODUCT_TEMPLATE_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_PRODUCT_TEMPLATE_ACK), receiverData["error_message"])
// 		}
// 	}
// }

func UpsertProductVariantConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPSERT_PRODUCT_VARIANT, eda.MDM)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.ChannelProductVariantUpsert(receiverData)
		}
	}
}

// func UpsertProductVariantConsumerAck(h *handler, KafkaConsumerTimeout time.Duration) {
// 	receiver := eda.ReaderConfig(eda.UPSERT_PRODUCT_VARIANT_ACK, eda.MDM)
// 	defer receiver.Close()
// 	for {
// 		receiverData,_ := eda.FetchReceiverPayload(receiver, ctx)
// 		if receiverData["response"] != nil {
// 			fmt.Println("The success response data received from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_PRODUCT_VARIANT_ACK), receiverData["response"])
// 		}
// 		if receiverData["error_message"] != nil {
// 			fmt.Println("An error response is receieved from the topic : ")
// 			helpers.PrettyPrint(string(eda.UPSERT_PRODUCT_VARIANT_ACK), receiverData["error_message"])
// 		}
// 	}
// }
