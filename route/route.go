package route

import (
	_ "fermion/docs"

	"fmt"
	"net/http"
	"os"

	"fermion/backend_core/controllers/offers"
	"fermion/backend_core/controllers/ondc"
	"fermion/backend_core/controllers/rating"
	"fermion/backend_core/controllers/webapp"

	module "fermion/backend_core/controllers/access/module"
	"fermion/backend_core/controllers/access/template"
	"fermion/backend_core/controllers/access/views"
	"fermion/backend_core/controllers/accounting/accounting"
	"fermion/backend_core/controllers/accounting/creditnote"
	"fermion/backend_core/controllers/accounting/currency_exchange"
	"fermion/backend_core/controllers/accounting/debitnote"
	"fermion/backend_core/controllers/accounting/payment_terms_and_record_payment"
	"fermion/backend_core/controllers/accounting/pos"
	"fermion/backend_core/controllers/accounting/purchase_invoice"
	"fermion/backend_core/controllers/accounting/sales_invoice"
	"fermion/backend_core/controllers/auth"
	app_core "fermion/backend_core/controllers/cores"
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

	// "fermion/backend_core/controllers/notifications"
	"fermion/backend_core/controllers/omnichannel/catalogue"
	"fermion/backend_core/controllers/omnichannel/channel"
	"fermion/backend_core/controllers/omnichannel/marketplace"
	"fermion/backend_core/controllers/omnichannel/omnichannel_fields"
	"fermion/backend_core/controllers/omnichannel/virtual_warehouse"
	"fermion/backend_core/controllers/omnichannel/webstores"
	"fermion/backend_core/controllers/orders/delivery_orders"
	"fermion/backend_core/controllers/orders/internal_transfers"
	"fermion/backend_core/controllers/orders/purchase_orders"
	"fermion/backend_core/controllers/orders/sales_orders"
	"fermion/backend_core/controllers/orders/scrap_orders"
	"fermion/backend_core/controllers/payments/customers"
	paymentpartners "fermion/backend_core/controllers/payments/payment_partners"
	"fermion/backend_core/controllers/payments/transactions"
	"fermion/backend_core/controllers/payments/wallets"
	"fermion/backend_core/controllers/returns/purchase_returns"
	"fermion/backend_core/controllers/returns/sales_returns"
	scheduler "fermion/backend_core/controllers/scheduler/app"
	"fermion/backend_core/controllers/shipping/shipping_orders"
	"fermion/backend_core/controllers/shipping/shipping_orders_ndr"
	"fermion/backend_core/controllers/shipping/shipping_orders_rto"
	"fermion/backend_core/controllers/shipping/shipping_orders_wd"
	"fermion/backend_core/controllers/shipping/shipping_partners"
	ipaas "fermion/backend_core/ipaas_core/app"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
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
func Init(g *echo.Group) {
	var (
		APP     = os.Getenv("APP")
		VERSION = os.Getenv("VERSION")
	)

	// Index
	g.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Welcome to %s version %s", APP, VERSION)
		return c.String(http.StatusOK, message)
	})

	g.GET("/swagger/*", echoSwagger.WrapHandler)
	// Routes

	auth.NewHandler().Route(g.Group("/auth"))
	webapp.NewHandler().Route(g.Group("/webapp"))

	//======================== mdm =================================================
	products.NewHandler().Route(g.Group("/api/v1/products"))
	contacts.NewHandler().Route(g.Group("/api/v1/contacts"))
	basic_inventory.NewHandler().Route(g.Group("/api/v1/basic_inventory"))
	locations.NewHandler().Route(g.Group("/api/v1/locations"))
	pricing.NewHandler().Route(g.Group("/api/v1/pricing"))
	uom.NewHandler().Route(g.Group("/api/v1/uom"))
	vendors.NewHandler().Route(g.Group("/api/v1/vendors"))

	//======================== orders =================================================
	purchase_orders.NewHandler().Route(g.Group("/api/v1/purchase_orders"))
	delivery_orders.NewHandler().Route(g.Group("/api/v1/delivery_orders"))
	scrap_orders.NewHandler().Route(g.Group("/api/v1/scrap_orders"))
	internal_transfers.NewHandler().Route(g.Group("/api/v1/internal_transfers"))
	sales_orders.NewHandler().Route(g.Group("/api/v1/sales_orders"))

	//======================== returns =================================================
	sales_returns.NewHandler().Route(g.Group("/api/v1/sales_returns"))
	purchase_returns.NewHandler().Route(g.Group("/api/v1/purchase_returns"))

	//======================== inventory_orders =================================================
	asn.NewHandler().Route(g.Group("/api/v1/asn"))
	grn.NewHandler().Route(g.Group("/api/v1/grn"))
	inventory_adjustments.NewHandler().Route(g.Group("/api/v1/inventory_adjustments"))

	//======================== inventory_tasks =================================================
	pick_list.NewHandler().Route(g.Group("/api/v1/pick_list"))
	cycle_count.NewHandler().Route(g.Group("/api/v1/cycle_count"))

	//======================== shipping =================================================
	shipping_orders.NewHandler().Route(g.Group("/api/v1/shipping_orders"))
	shipping_partners.NewHandler().Route(g.Group("/api/v1/shipping_partners"))
	shipping_orders_wd.NewHandler().Route(g.Group("/api/v1/shipping_orders_wd"))
	shipping_orders_ndr.NewHandler().Route(g.Group("/api/v1/shipping_orders_ndr"))
	shipping_orders_rto.NewHandler().Route(g.Group("/api/v1/shipping_orders_rto"))

	//======================== accounting =================================================
	currency_exchange.NewHandler().Route(g.Group("/api/v1/currency_and_exchange"))
	payment_terms_and_record_payment.NewHandler().Route(g.Group("/api/v1/payment_terms"))
	creditnote.NewHandler().Route(g.Group("/api/v1/creditnote"))
	debitnote.NewHandler().Route(g.Group("/api/v1/debitnote"))
	purchase_invoice.NewHandler().Route(g.Group("/api/v1/purchase_invoice"))
	sales_invoice.NewHandler().Route(g.Group("/api/v1/sales_invoice"))
	accounting.NewHandler().Route(g.Group("/api/v1/accounting"))
	pos.NewHandler().Route(g.Group("/api/v1/pos"))

	//======================== core =================================================
	app_core.NewHandler().Route(g.Group("/api/v1/core"))

	//======================== omnichannel =================================================
	omnichannel_fields.NewHandler().Route(g.Group("/api/v1/omnichannel_fields"))
	marketplace.NewHandler().Route(g.Group("/api/v1/marketplace"))
	webstores.NewHandler().Route(g.Group("/api/v1/webstores"))
	virtual_warehouse.NewHandler().Route(g.Group("/api/v1/virtual_warehouse"))
	catalogue.NewHandler().Route(g.Group("/api/v1/catalogue"))
	channel.NewHandler().Route(g.Group("api/v1/channels"))

	//======================== payments =================================================
	customers.NewHandler().Route(g.Group("/api/v1/customers"))
	transactions.NewHandler().Route(g.Group("/api/v1/transactions"))
	wallets.NewHandler().Route(g.Group("/api/v1/wallets"))
	paymentpartners.NewHandler().Route(g.Group("api/v1/payment_partners"))

	//======================== access =================================================
	template.NewHandler().Route(g.Group("/api/v1/template"))
	module.NewHandler().Route(g.Group("/api/v1/module"))
	views.NewHandler().Route(g.Group("/api/v1/views"))

	//======================== ipaas =================================================
	ipaas.NewHandler().Route(g.Group("ipaas"))

	//======================== schedular =================================================
	scheduler.NewHandler().Route(g.Group("/scheduler"))

	//======================== ondc =================================================
	ondc.NewHandler().Route(g.Group("/api/v1/ondc"))

	rating.NewHandler().Route(g.Group("/api/v1/rating"))

	offers.NewHandler().Route(g.Group("/api/v1/offers"))
}
