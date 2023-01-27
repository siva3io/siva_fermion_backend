package sales_orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"fermion/backend_core/controllers/eda"
	products "fermion/backend_core/controllers/mdm/products"
	"fermion/backend_core/controllers/omnichannel/channel"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/pagination"
	shipping_model "fermion/backend_core/internal/model/shipping"
	"fermion/backend_core/pkg/util/helpers"

	"gorm.io/datatypes"
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
	go CreateSalesOrdersConsumer(h, KafkaConsumerTimeout)
	go UpdateSalesOrdersConsumer(h, KafkaConsumerTimeout)
	go UpsertSalesOrdersConsumer(h, KafkaConsumerTimeout)

	go CancelSalesOrdersConsumer(h, KafkaConsumerTimeout)
	go OndcOrderConsumer(h, KafkaConsumerTimeout)
	go UpdateSalesOrdersStatusConsumer(h, KafkaConsumerTimeout)
	go GetSalesOrderStatusConsumer(h, KafkaConsumerTimeout)
}

func CreateSalesOrdersConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CREATE_SALES_ORDER, eda.ORDERS)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.CreateSalesOrders(receiverData)
		}
	}
}
func UpdateSalesOrdersConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_SALES_ORDER, eda.ORDERS)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.UpdateSalesOrders(receiverData)
		}
	}
}
func UpsertSalesOrdersConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPSERT_SALES_ORDER, eda.ORDERS)
		for {
			receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
			if err != nil {
				receiver.Close()
				cancel()
				break
			}
			h.ChannelSalesOrderUpsert(receiverData)
		}
	}
}

func UpdateSalesOrdersStatusConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.UPDATE_SALES_ORDER_STATUS, eda.ORDERS)
		for {
			err := func() error {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("PANIC ERROR: ", r, " in UpdateSalesOrdersStatusConsumer")
					}
				}()

				receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
				if err != nil {
					receiver.Close()
					cancel()
					return err
				}
				var edaMetaData core.MetaData
				helpers.JsonMarshaller(receiverData["meta_data"], &edaMetaData)
				// fmt.Println("reuestId", edaMetaData.RequestId)

				if transactionId := receiverData["data"].(map[string]interface{})["transaction_id"]; transactionId != nil {
					edaMetaData.Query = map[string]interface{}{
						"reference_number": transactionId,
					}
					sales_order, _ := h.service.ViewSalesOrder(edaMetaData)
					var salesOrderResponse OndcOrdersResponseDTO
					helpers.JsonMarshaller(sales_order, &salesOrderResponse)

					if quote := receiverData["data"].(map[string]interface{})["quote"]; quote != nil {
						salesOrderResponse.OndcContext["message"].(map[string]interface{})["order"].(map[string]interface{})["quote"] = quote
						var updatedOrder orders.SalesOrders
						updatedOrder.OndcContext, _ = json.Marshal(salesOrderResponse.OndcContext)
						h.service.UpdateSalesOrder(edaMetaData, &updatedOrder)
					}
					return nil
				}

				edaMetaData.Query = map[string]interface{}{
					"reference_number": receiverData["data"].(map[string]interface{})["context"].(map[string]interface{})["transaction_id"],
				}

				OndcState := receiverData["data"].(map[string]interface{})["message"].(map[string]interface{})["order"].(map[string]interface{})["state"].(string)

				salesOrderStatusMapping := map[string]string{
					"Order-picked-up": "SHIPPED",
					"Order-delivered": "DELIVERED",
					// "RTO-Initiated":"",
					// "RTO-Delivered":"",
					"Cancelled": "CANCELLED",
				}

				statusId, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", salesOrderStatusMapping[OndcState])

				order := orders.SalesOrders{
					StatusId: &statusId,
				}
				if statusId != 0 {
					h.service.UpdateSalesOrder(edaMetaData, &order)
				}
				return nil
			}()

			if err != nil {
				break
			}
		}
	}
}
func CancelSalesOrdersConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.CANCEL_SALES_ORDER, eda.ORDERS)
		for {

			err := func() error {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("PANIC ERROR: ", r, "in CancelSalesOrdersConsumer")
					}
				}()

				receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
				if err != nil {
					receiver.Close()
					cancel()
					return err
				}

				var edaMetaData core.MetaData
				helpers.JsonMarshaller(receiverData["meta_data"], &edaMetaData)

				edaMetaData.Query = map[string]interface{}{
					"reference_number": receiverData["data"].(map[string]interface{})["context"].(map[string]interface{})["transaction_id"],
				}

				if edaMetaData.AdditionalFields == nil {
					edaMetaData.AdditionalFields = map[string]interface{}{}
				}
				edaMetaData.AdditionalFields["don't produce"] = true

				data, _ := receiverData["data"].(map[string]interface{})

				if ondcOrder := data["message"].(map[string]interface{})["order"]; ondcOrder != nil {
					if ondcOrderItems := ondcOrder.(map[string]interface{})["items"]; ondcOrderItems != nil {

						var salesOrder orders.SalesOrders

						oldSalesOrderInterface, _ := h.service.ViewSalesOrder(edaMetaData)
						var oldSalesOrder orders.SalesOrders
						err = helpers.JsonMarshaller(oldSalesOrderInterface, &oldSalesOrder)
						if err != nil {
							fmt.Println("-----Ondc Order Cancellation------", err)
							return nil
						}

						var oldOndcContext map[string]interface{}
						err = json.Unmarshal(oldSalesOrder.OndcContext, &oldOndcContext)
						if err != nil {
							fmt.Println("-----Ondc Order Cancellation------", err)
							return nil
						}

						oldOndcContext["message"].(map[string]interface{})["order"].(map[string]interface{})["quote"] = data["message"].(map[string]interface{})["order"].(map[string]interface{})["quote"]

						salesOrder.OndcContext, _ = json.Marshal(oldOndcContext)

						statusId, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "PARTIALLY_CANCELLED")

						salesOrder.StatusId = &statusId

						var breakups []interface{}

						if quote := ondcOrder.(map[string]interface{})["quote"]; quote != nil {
							if breakupInterface := quote.(map[string]interface{})["breakup"]; breakupInterface != nil {
								breakups = breakupInterface.([]interface{})
							}

							if price := quote.(map[string]interface{})["price"]; price != nil {
								if priceValue := price.(map[string]interface{})["value"]; priceValue != nil {
									priceValueFloat64, _ := strconv.ParseFloat(fmt.Sprint(priceValue), 32)
									salesOrder.SoPaymentDetails.TotalAmount = float32(priceValueFloat64)
								}
							}
						}

						for _, orderLine := range oldSalesOrder.SalesOrderLines {

							var newOrderLine orders.SalesOrderLines

							var breakupValue map[string]interface{}
							for _, breakup := range breakups {
								breakupValue = breakup.(map[string]interface{})
								if itemId := breakupValue["@ondc/org/item_id"]; itemId != nil {
									if titleType := breakupValue["@ondc/org/title_type"]; titleType != nil {
										if itemId == orderLine.SerialNumber && titleType == "item" {

											newOrderLine.ProductId = orderLine.ProductId

											if itemQuantity := breakupValue["@ondc/org/item_quantity"]; itemQuantity != nil {
												if count := itemQuantity.(map[string]interface{})["count"]; count != nil {
													priceValueInt, _ := strconv.Atoi(fmt.Sprint(count))
													newOrderLine.Quantity = int64(priceValueInt)
												}
											}

											if price := breakupValue["price"]; price != nil {
												if priceValue := price.(map[string]interface{})["value"]; priceValue != nil {
													priceValueFloat64, _ := strconv.ParseFloat(fmt.Sprint(priceValue), 32)
													newOrderLine.Amount = float32(priceValueFloat64)
												}
											}
											salesOrder.SalesOrderLines = append(salesOrder.SalesOrderLines, newOrderLine)
										}
									}
								}
							}
						}

						// helpers.PrettyPrint("salesOrder", salesOrder)
						err = h.service.UpdateSalesOrder(edaMetaData, &salesOrder)
						if err != nil {
							fmt.Println("-----Ondc Order Cancellation------", err)
							return nil
						}
						return nil
					}
				}

				helpers.PrettyPrint("edaMetaData", edaMetaData)

				statusId, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "CANCELLED")

				order := orders.SalesOrders{
					StatusId: &statusId,
				}
				h.service.UpdateSalesOrder(edaMetaData, &order)

				sales_order, _ := h.service.ViewSalesOrder(edaMetaData)
				var salesOrderResponse OndcOrdersResponseDTO
				err = helpers.JsonMarshaller(sales_order, &salesOrderResponse)

				if err != nil {
					fmt.Println("787878", err)
				}

				edaMetaData.Encryption = false
				requestPayload := map[string]interface{}{
					"meta_data":      edaMetaData,
					"data":           salesOrderResponse,
					"cancel_request": data,
				}

				eda.Produce(eda.CANCEL_SALES_ORDER_RESPONSE, requestPayload, false)
				return nil
			}()

			if err != nil {
				break
			}
		}
	}
}

func OndcOrderConsumer(h *handler, KafkaConsumerTimeout time.Duration) {

	PaymentTypesMapper := map[string]string{
		"ON-ORDER":         "PRE-PAID",
		"POST-FULFILLMENT": "COD",
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.ONDC_CREATE_ORDER, eda.BECKN)
		for {

			err := func() error {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("PANIC ERROR: ", r, "in CancelSalesOrdersConsumer")
					}
				}()

				receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
				if err != nil {
					receiver.Close()
					cancel()
					return err
				}
				var lsp_data map[string]interface{}
				if data := receiverData["data"].(map[string]interface{})["lsp"]; data != nil {
					lsp_data = data.(map[string]interface{})
				}

				mapper, err := helpers.ReadFile("./backend_core/controllers/orders/sales_orders/order_mapper_template.json")
				if err != nil {
					fmt.Println("-----OndcOrderConsumer------", err)
					return nil
				}

				authToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MTEsIlVzZXJuYW1lIjoic3VwZXJfdXNlckBldW5pbWFydC5jb20iLCJhY2Nlc3NfdGVtcGxhdGVfaWQiOjEsImNvbXBhbnlfaWQiOjEsImV4cCI6MTcwNzU0NTI0NCwiZmlyc3RfbmFtZSI6IlN1cGVyIiwibGFzdF9uYW1lIjoiVXNlciIsInJvbGVfaWQiOm51bGwsInVzZXJfdHlwZXMiOlt7ImlkIjo2MzQsIm5hbWUiOiJTRUxMRVIifV19.yFfpIOAKL2--8uutOdVvgGnO9Lyruz72bkM0cpFT8fU"

				data, _ := receiverData["data"].(map[string]interface{})

				response, err := helpers.MakeRequest(helpers.Request{
					Method: "POST",
					Scheme: "http",
					Host:   "localhost:3031",
					Path:   "ipaas/boson_convertor",
					Header: map[string]string{
						"Authorization": authToken,
						"Content-Type":  "application/json",
					},
					Body: map[string]interface{}{
						"data": map[string]interface{}{
							"input_data":      data,
							"mapper_template": mapper,
							"feature_session_variables": []map[string]interface{}{
								{
									"key":   "host",
									"value": "localhost:3031",
								},
								{
									"key":   "Authorization",
									"value": authToken,
								},
							},
						},
					},
				})

				helpers.PrettyPrint("mapper_reponse", response)

				if err != nil {
					fmt.Println("-----OndcOrderConsumer------", err)
					return nil
				}

				var order orders.SalesOrders
				var metaData core.MetaData

				helpers.PrettyPrint("meta_data", receiverData["meta_data"])

				err = errors.New("mapper error")
				if data := response.(map[string]interface{})["data"]; data != nil {
					if mapped_response := data.(map[string]interface{})["mapped_response"]; mapped_response != nil {
						err = nil
					}
				}
				if err != nil {
					fmt.Println("-----OndcOrderConsumer------", response)
					return nil
				}

				err = helpers.JsonMarshaller(response.(map[string]interface{})["data"].(map[string]interface{})["mapped_response"], &order)
				if err != nil {
					fmt.Println("-----OndcOrderConsumer------", err)
					return nil
				}
				err = helpers.JsonMarshaller(receiverData["meta_data"], &metaData)
				if err != nil {
					fmt.Println("-----OndcOrderConsumer------", err)
					return nil
				}

				currentDate := datatypes.Date(time.Now())
				order.SoDate = currentDate
				order.ExpectedShippingDate = currentDate
				order.OndcContext, _ = json.Marshal(data)

				var breakups []interface{}

				if message := data["message"]; message != nil {
					if messageOrder := message.(map[string]interface{})["order"]; messageOrder != nil {
						if quote := messageOrder.(map[string]interface{})["quote"]; quote != nil {
							if breakupInterface := quote.(map[string]interface{})["breakup"]; breakupInterface != nil {
								breakups = breakupInterface.([]interface{})
							}
						}
						if payment := messageOrder.(map[string]interface{})["payment"]; payment != nil {
							if paymentTypeInterface := payment.(map[string]interface{})["type"]; paymentTypeInterface != nil {
								paymentTypeId, err := helpers.GetLookupcodeId("PAYMENT_TYPE", PaymentTypesMapper[paymentTypeInterface.(string)])
								if err == nil {
									order.PaymentTypeId = &paymentTypeId
								}
							}
						}
					}
				}

				for index, orderLine := range order.SalesOrderLines {
					statusId, _ := helpers.GetLookupcodeId("SALES_ORDER_STATUS", "NEW")
					order.SalesOrderLines[index].StatusId = &statusId

					var breakupValue map[string]interface{}
					for _, breakup := range breakups {
						breakupValue = breakup.(map[string]interface{})
						fmt.Println(breakupValue["@ondc/org/item_id"], breakupValue["@ondc/org/title_type"])
						if itemId := breakupValue["@ondc/org/item_id"]; itemId != nil {
							if titleType := breakupValue["@ondc/org/title_type"]; titleType != nil {
								if itemId == orderLine.SerialNumber && titleType == "item" {
									if price := breakupValue["price"]; price != nil {
										if priceValue := price.(map[string]interface{})["value"]; priceValue != nil {
											priceValueFloat64, _ := strconv.ParseFloat(priceValue.(string), 32)
											order.SalesOrderLines[index].Amount = float32(priceValueFloat64)
										}
									}
									if item := breakupValue["item"]; item != nil {
										if price := item.(map[string]interface{})["price"]; price != nil {
											if value := price.(map[string]interface{})["value"]; value != nil {
												priceValueFloat64, _ := strconv.ParseFloat(value.(string), 32)
												order.SalesOrderLines[index].Price = float32(priceValueFloat64)
											}
										}
									}
								}
								// if itemId == orderLine.SerialNumber && titleType == "tax" {
								// 	if price := breakupValue["price"]; price != nil {
								// 		if priceValue := price.(map[string]interface{})["value"]; priceValue != nil {
								// 			order.SalesOrderLines[index].Tax = float32(priceValue.(float64))
								// 		}
								// 	}
								// }
							}
						}
					}

				}

				// helpers.PrettyPrint("metaData", metaData)
				// helpers.PrettyPrint("order", order)

				if lsp_data != nil {
					order.LogisticsNetworkOrderId = lsp_data["id"].(string)
					order.LogisticsNetworkTransactionId = lsp_data["context"].(map[string]interface{})["transaction_id"].(string)
					order.LogisticsSellerNpName = lsp_data["provider"].(map[string]interface{})["id"].(string)
				}
				err = h.service.CreateSalesOrder(metaData, &order)
				if err != nil {
					fmt.Println("-----OndcOrderConsumer------", err)
					return nil
				}

				//================create Shipment =========================
				shippingOrder := new(shipping_model.ShippingOrder)

				channel_name := order.ChannelName

				var channel_data channel.ChannelDTO
				fmt.Println("chanel name------------", channel_name)
				if channel_name != "" {
					query := map[string]interface{}{
						"name": channel_name,
					}
					channel_data, _ = h.ChannelService.GetChannelService(query)
				} else {
					channel_data.ID = 1
				}

				shippingOrder.ChannelId = &channel_data.ID
				fmt.Println("channel id-----------------------", channel_data.ID)
				if lsp_data != nil {
					// TODO : ERROR HANDLE

					metaData.Query = map[string]interface{}{
						"id": order.ID,
					}

					shippingOrder.OrderId = &order.ID
					sender_add := lsp_data["fulfillments"].([]interface{})[0].(map[string]interface{})["start"].(map[string]interface{})["location"].(map[string]interface{})["address"].(map[string]interface{})

					shippingOrder.SenderAddress, _ = json.Marshal(map[string]interface{}{
						"sender_name":   sender_add["name"],
						"email":         lsp_data["fulfillments"].([]interface{})[0].(map[string]interface{})["start"].(map[string]interface{})["contact"].(map[string]interface{})["email"],
						"mobile_number": lsp_data["fulfillments"].([]interface{})[0].(map[string]interface{})["start"].(map[string]interface{})["contact"].(map[string]interface{})["phone"],
						"address_line1": sender_add["building"],
						"address_line2": sender_add["locality"],
						"address_line3": sender_add["city"],
						"state":         sender_add["state"],
						"country":       sender_add["country"],
						"landmark":      sender_add["building"],
						"zipcode":       sender_add["area_code"],
						"city":          sender_add["city"],
					})

					fmt.Println("order---------", shippingOrder.OrderId)
					sales_order, _ := h.service.ViewSalesOrder(metaData)
					var salesOrderResponse SalesOrdersDTO
					helpers.JsonMarshaller(sales_order, &salesOrderResponse)
					// shippingOrder.productvalue = salesOrderResponse.SalesOrderLines[0].Product.ProductPricingDetails.SalesPrice

					shippingOrder.PackageDetails = salesOrderResponse.SalesOrderLines[0].Product.PackageDimensions
					cost := lsp_data["quote"].(map[string]interface{})["price"].(map[string]interface{})["value"].(string)
					cost1, _ := strconv.ParseFloat(cost, 64)
					shippingOrder.ShippingCost = cost1
					shippingOrder.ActualShippingCost = shippingOrder.ShippingCost
					// shippingOrder.EstimatedWeight = salesOrderResponse.SalesOrderLines[0].Product.ProductPricingDetails
				} else {
					shippingOrder.OrderId = &order.ID
					company_details, _ := h.cores_service.GetCompanyPreferences(map[string]interface{}{
						"id": order.CompanyId,
					})

					locationPagnation := new(pagination.Paginatevalue)
					lookup_id, _ := helpers.GetLookupcodeId("LOCATION_TYPE", "OFFICE")
					locationPagnation.Filters = fmt.Sprintf("[[\"company_id\",\"=\",%v],[\"location_type_id\",\"=\",%v]]", metaData.CompanyId, lookup_id)
					location_resp, _ := h.locations_service.GetLocationList(metaData, locationPagnation)

					var location_data = make([]map[string]interface{}, 0)
					helpers.JsonMarshaller(location_resp, &location_data)

					length := len(location_data)
					if length > 0 {
						address_data := location_data[0]["address"].(map[string]interface{})

						shippingOrder.SenderAddress, _ = json.Marshal(map[string]interface{}{
							"sender_name":   company_details.Name,
							"email":         company_details.Email,
							"mobile_number": company_details.Phone,
							"address_line1": address_data["address_line_1"],
							"address_line2": address_data["address_line_2"],
							"address_line3": address_data["address_line_3"],
							"state":         address_data["state"].(map[string]interface{})["name"],
							"country":       address_data["country"].(map[string]interface{})["name"],
							"landmark":      address_data["land_mark"],
							"zipcode":       address_data["pin_code"],
							"city":          address_data["city"],
						})
					}

					productQuery := map[string]interface{}{
						"id": order.SalesOrderLines[0].ProductId,
					}
					// helpers.PrettyPrint("orderdata",order.SalesOrderLines)
					product_details, _ := h.Product_service.GetTemplateView(productQuery, fmt.Sprint(metaData.TokenUserId), fmt.Sprint(metaData.AccessTemplateId))
					var template products.TemplateReponseDTO
					helpers.JsonMarshaller(product_details, &template)
					helpers.PrettyPrint("Package details", template.PackageDimensions)
					shippingOrder.PackageDetails = template.PackageDimensions
				}

				var billingAddress CustomerBillingorShippingAddress
				json.Unmarshal(order.CustomerBillingAddress, &billingAddress)

				var receivingAddress CustomerBillingorShippingAddress
				json.Unmarshal(order.CustomerShippingAddress, &receivingAddress)

				shippingOrder.BillingAddress, _ = json.Marshal(map[string]interface{}{
					"sender_name":   billingAddress.ContactPersonName,
					"mobile_number": billingAddress.ContactPersonNumber,
					"address_line1": billingAddress.AddressLine1,
					"address_line2": billingAddress.AddressLine2,
					"address_line3": billingAddress.AddressLine3,
					"zipcode":       billingAddress.PinCode,
					"city":          billingAddress.City,
					"state":         billingAddress.State,
					"country":       billingAddress.Country,
					"landmark":      billingAddress.Landmark,
				})
				shippingOrder.ReceiverAddress, _ = json.Marshal(map[string]interface{}{
					"receiver_name": receivingAddress.ContactPersonName,
					"mobile_number": receivingAddress.ContactPersonNumber,
					"address_line1": receivingAddress.AddressLine1,
					"address_line2": receivingAddress.AddressLine2,
					"address_line3": receivingAddress.AddressLine3,
					"zipcode":       receivingAddress.PinCode,
					"city":          receivingAddress.City,
					"state":         receivingAddress.State,
					"country":       receivingAddress.Country,
					"landmark":      receivingAddress.Landmark,
				})
				pincode, _ := strconv.Atoi(receivingAddress.PinCode)
				shippingOrder.DestinationZipcode = int32(pincode)
				pincode, _ = strconv.Atoi(billingAddress.PinCode)
				shippingOrder.OriginZipcode = int32(pincode)
				shippingPartnerId := uint(8)
				shippingOrder.ShippingPartnerId = &shippingPartnerId
				shippingOrder.PaymentStatus = "UNPAID"
				shippingOrder.ShippingPaymentTypeID = order.PaymentTypeId
				for _, sales_order_line := range order.SalesOrderLines {
					ShippingOrderLine := new(shipping_model.ShippingOrderLines)
					ShippingOrderLine.ProductId = sales_order_line.ProductId
					ShippingOrderLine.ProductTemplateId = sales_order_line.ProductTemplateId
					ShippingOrderLine.ItemQuantity = int32(sales_order_line.Quantity)
					ShippingOrderLine.UnitPrice = float64(sales_order_line.Price)
					shippingOrder.ShippingOrderLines = append(shippingOrder.ShippingOrderLines, *ShippingOrderLine)
				}
				shippingOrder.Quantity = int32(order.TotalQuantity)
				shippingOrder.OrderDate = order.CreatedDate
				fmt.Println("company id-------", order.CompanyId)
				shippingOrder.SetPickupDate = order.ReadyTOShip
				shippingOrder.OrderValue = float64(order.SoPaymentDetails.TotalAmount)
				temp := map[string]interface{}{}
				err = json.Unmarshal(shippingOrder.PackageDetails, &temp)
				helpers.PrettyPrint("package detail 2", temp)
				if err != nil {
					fmt.Println("-----Shipping Order Package Details------", err)
					return nil
				}
				temp["product_value"] = shippingOrder.OrderValue
				shippingOrder.PackageDetails, err = json.Marshal(temp)
				shippingOrder.OrderDate = order.CreatedDate
				lookup_id, _ := helpers.GetLookupcodeId("SHIPPING_TYPE", "DOMESTIC")
				shippingOrder.ShippingTypeId = &lookup_id

				err = h.ShippingOrderService.CreateShippingOrder(metaData, shippingOrder)
				if err != nil {
					return nil
				}
				return nil
			}()

			if err != nil {
				break
			}
		}
	}
}
func GetSalesOrderStatusConsumer(h *handler, KafkaConsumerTimeout time.Duration) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), KafkaConsumerTimeout)
		receiver := eda.ReaderConfig(eda.GET_SALES_ORDER, eda.BECKN)
		for {

			err := func() error {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("PANIC ERROR: ", r, "in GetSalesOrderStatusConsumer")
					}
				}()

				receiverData, err := eda.FetchReceiverPayload(receiver, ctx)
				if err != nil {
					receiver.Close()
					cancel()
					return err
				}

				var edaMetaData core.MetaData
				helpers.JsonMarshaller(receiverData["meta_data"], &edaMetaData)

				edaMetaData.Query = map[string]interface{}{
					"reference_number": receiverData["data"].(map[string]interface{})["context"].(map[string]interface{})["transaction_id"],
				}
				// Type := edaMetaData.AdditionalFields["type"]

				helpers.PrettyPrint("edaMetaData", edaMetaData)

				sales_order, _ := h.service.ViewSalesOrder(edaMetaData)
				var salesOrderResponse OndcOrdersResponseDTO
				helpers.JsonMarshaller(sales_order, &salesOrderResponse)
				edaMetaData.Encryption = false
				requestPayload := map[string]interface{}{
					"meta_data":     edaMetaData,
					"data":          salesOrderResponse,
					"order_request": receiverData["data"],
				}
				eda.Produce(eda.SEND_SALES_ORDER, requestPayload, false)
				return nil
			}()

			if err != nil {
				break
			}
		}
	}
}
