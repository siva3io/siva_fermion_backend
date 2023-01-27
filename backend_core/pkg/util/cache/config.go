package cache

type KeyGroups string

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

const (
	PAGINATION              KeyGroups = "pagination"
	NOTIFICATION            KeyGroups = "notification"
	PRODUCT                 KeyGroups = "product"
	DECENTRALIZED_INVENTORY KeyGroups = "decentralized_inventory"
	CENTRALIZED_INVENTORY   KeyGroups = "centralized_inventory"
	CONTACT                 KeyGroups = "contact"
	LOCATION                KeyGroups = "location"
	PRICING                 KeyGroups = "pricing"
	UOM                     KeyGroups = "uom"
	UOM_CLASS               KeyGroups = "uom_class"
	DELIVER_ORDER           KeyGroups = "deliver_order"
	INTERNAL_TRANSFERS      KeyGroups = "internal_transfers"
	PURCHASE_ORDER          KeyGroups = "purchase_order"
	SALES_ORDER             KeyGroups = "sales_order"
	SCRAP_ORDER             KeyGroups = "scrap_order"
	ASN                     KeyGroups = "asn"
	GRN                     KeyGroups = "grn"
	INVADJ                  KeyGroups = "inventory_adjustment"
	CYCLE_COUNT             KeyGroups = "cycle_count"
	PICK_LIST               KeyGroups = "pick_list"
	ACCOUNTING_LINK         KeyGroups = "accounting_link"
	CREDIT_NOTE             KeyGroups = "credit_note"
	CURRENCY_EXCHANGE       KeyGroups = "currency_exchange"
	DEBIT_NOTE              KeyGroups = "debit_note"
	PAYMENT_TERMS           KeyGroups = "payment_terms"
	POS_LINK                KeyGroups = "pos_link"
	PURCHASE_INVOICE        KeyGroups = "purchase_invoice"
	SALES_INVOICE           KeyGroups = "sales_invoice"
	PURCHASE_RETURNS        KeyGroups = "purchase_returns"
	SALES_RETURNS           KeyGroups = "sales_returns"
	SHIPPING_ORDER          KeyGroups = "shipping_order"
	NDR                     KeyGroups = "ndr"
	RTO                     KeyGroups = "rto"
	WD                      KeyGroups = "wd"
	LOGISTIC_PARTNER        KeyGroups = "logistic_partner"
	MARKETPLACE_LINK        KeyGroups = "marketplace_link"
	VIRTUAL_WAREHOUSE       KeyGroups = "virtual_warehouse"
	WEBSTORE                KeyGroups = "webstore"
	OFFERS                  KeyGroups = "offers"
)
