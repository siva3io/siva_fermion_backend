package consumer_init

import (
	"time"

	"fermion/backend_core/controllers/accounting/creditnote"
	"fermion/backend_core/controllers/accounting/currency_exchange"
	"fermion/backend_core/controllers/accounting/debitnote"
	"fermion/backend_core/controllers/accounting/payment_terms_and_record_payment"
	"fermion/backend_core/controllers/accounting/purchase_invoice"
	"fermion/backend_core/controllers/accounting/sales_invoice"
	"fermion/backend_core/controllers/auth"
	"fermion/backend_core/controllers/inventory_orders/asn"
	"fermion/backend_core/controllers/inventory_orders/grn"
	"fermion/backend_core/controllers/inventory_orders/inventory_adjustments"
	"fermion/backend_core/controllers/inventory_tasks/cycle_count"
	"fermion/backend_core/controllers/inventory_tasks/pick_list"
	"fermion/backend_core/controllers/mdm/basic_inventory"
	"fermion/backend_core/controllers/mdm/contacts"
	"fermion/backend_core/controllers/mdm/locations"
	"fermion/backend_core/controllers/mdm/pricing"
	"fermion/backend_core/controllers/mdm/products"
	"fermion/backend_core/controllers/mdm/uom"
	"fermion/backend_core/controllers/mdm/vendors"
	"fermion/backend_core/controllers/ondc"
	"fermion/backend_core/controllers/orders/delivery_orders"
	"fermion/backend_core/controllers/orders/internal_transfers"
	"fermion/backend_core/controllers/orders/purchase_orders"
	"fermion/backend_core/controllers/orders/sales_orders"
	"fermion/backend_core/controllers/orders/scrap_orders"
	"fermion/backend_core/controllers/payments/customers"
	"fermion/backend_core/controllers/payments/transactions"
	"fermion/backend_core/controllers/payments/wallets"
	"fermion/backend_core/controllers/rating"
	"fermion/backend_core/controllers/returns/purchase_returns"
	"fermion/backend_core/controllers/returns/sales_returns"
	"fermion/backend_core/controllers/shipping/shipping_orders"
	"fermion/backend_core/controllers/shipping/shipping_orders_ndr"
	"fermion/backend_core/controllers/shipping/shipping_orders_rto"
	"fermion/backend_core/controllers/shipping/shipping_orders_wd"
	"fermion/backend_core/controllers/shipping/shipping_partners"
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

func Init() {

	KafkaConsumerTimeout := 30 * time.Minute

	//======================== mdm =================================================
	go products.InitConsumers(products.NewHandler(), KafkaConsumerTimeout)
	go contacts.InitConsumers(contacts.NewHandler(), KafkaConsumerTimeout)
	go basic_inventory.InitConsumers(basic_inventory.NewHandler(), KafkaConsumerTimeout)
	go locations.InitConsumers(locations.NewHandler(), KafkaConsumerTimeout)
	go pricing.InitConsumers(pricing.NewHandler(), KafkaConsumerTimeout)
	go uom.InitConsumers(uom.NewHandler(), KafkaConsumerTimeout)
	go vendors.InitConsumers(vendors.NewHandler(), KafkaConsumerTimeout)

	//======================== orders =================================================
	go sales_orders.InitConsumers(sales_orders.NewHandler(), KafkaConsumerTimeout)
	go purchase_orders.InitConsumers(purchase_orders.NewHandler(), KafkaConsumerTimeout)
	go delivery_orders.InitConsumers(delivery_orders.NewHandler(), KafkaConsumerTimeout)
	go scrap_orders.InitConsumers(scrap_orders.NewHandler(), KafkaConsumerTimeout)
	go internal_transfers.InitConsumers(internal_transfers.NewHandler(), KafkaConsumerTimeout)

	//======================== returns =================================================
	go sales_returns.InitConsumers(sales_returns.NewHandler(), KafkaConsumerTimeout)
	go purchase_returns.InitConsumers(purchase_returns.NewHandler(), KafkaConsumerTimeout)

	//======================== inventory_orders =================================================
	go asn.InitConsumers(asn.NewHandler(), KafkaConsumerTimeout)
	go grn.InitConsumers(grn.NewHandler(), KafkaConsumerTimeout)
	go inventory_adjustments.InitConsumers(inventory_adjustments.NewHandler(), KafkaConsumerTimeout)

	//======================== inventory_tasks =================================================
	go pick_list.InitConsumers(pick_list.NewHandler(), KafkaConsumerTimeout)
	go cycle_count.InitConsumers(cycle_count.NewHandler(), KafkaConsumerTimeout)

	//======================== shipping =================================================
	go shipping_orders.InitConsumers(shipping_orders.NewHandler(), KafkaConsumerTimeout)
	go shipping_orders_wd.InitConsumers(shipping_orders_wd.NewHandler(), KafkaConsumerTimeout)
	go shipping_orders_ndr.InitConsumers(shipping_orders_ndr.NewHandler(), KafkaConsumerTimeout)
	go shipping_orders_rto.InitConsumers(shipping_orders_rto.NewHandler(), KafkaConsumerTimeout)
	go shipping_partners.InitConsumers(shipping_partners.NewHandler(), KafkaConsumerTimeout)

	//======================== accounting =================================================
	go creditnote.InitConsumers(creditnote.NewHandler(), KafkaConsumerTimeout)
	go debitnote.InitConsumers(debitnote.NewHandler(), KafkaConsumerTimeout)
	go purchase_invoice.InitConsumers(purchase_invoice.NewHandler(), KafkaConsumerTimeout)
	go sales_invoice.InitConsumers(sales_invoice.NewHandler(), KafkaConsumerTimeout)
	go currency_exchange.InitConsumers(currency_exchange.NewHandler(), KafkaConsumerTimeout)
	go payment_terms_and_record_payment.InitConsumers(payment_terms_and_record_payment.NewHandler(), KafkaConsumerTimeout)

	//======================== core =================================================
	go auth.InitConsumers(auth.NewHandler(), KafkaConsumerTimeout)

	//======================== payments =================================================
	go customers.InitConsumers(customers.NewHandler(), KafkaConsumerTimeout)
	go transactions.InitConsumers(transactions.NewHandler(), KafkaConsumerTimeout)
	go wallets.InitConsumers(wallets.NewHandler(), KafkaConsumerTimeout)

	go rating.InitConsumers(rating.NewHandler(), KafkaConsumerTimeout)

	go ondc.InitConsumers(ondc.NewHandler(), KafkaConsumerTimeout)
}
