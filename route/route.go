package route

import (
	_ "fermion/docs"

	"fmt"
	"net/http"
	"os"

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
	"fermion/backend_core/controllers/omnichannel/catalogue"
	"fermion/backend_core/controllers/omnichannel/channel"
	"fermion/backend_core/controllers/omnichannel/marketplace"
	"fermion/backend_core/controllers/omnichannel/virtual_warehouse"
	"fermion/backend_core/controllers/omnichannel/webstores"
	"fermion/backend_core/controllers/orders/delivery_orders"
	"fermion/backend_core/controllers/orders/internal_transfers"
	"fermion/backend_core/controllers/orders/purchase_orders"
	"fermion/backend_core/controllers/orders/sales_orders"
	"fermion/backend_core/controllers/orders/scrap_orders"
	"fermion/backend_core/controllers/payments/customers"
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
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
	currency_exchange.NewHandler().Route(g.Group("/api/v1/currency_and_exchange"))
	payment_terms_and_record_payment.NewHandler().Route(g.Group("/api/v1/payment_terms"))
	asn.NewHandler().Route(g.Group("/api/v1/asn"))
	sales_orders.NewHandler().Route(g.Group("/api/v1/sales_orders"))
	webstores.NewHandler().Route(g.Group("/api/v1/webstores"))
	marketplace.NewHandler().Route(g.Group("/api/v1/marketplace"))
	creditnote.NewHandler().Route(g.Group("/api/v1/creditnote"))
	debitnote.NewHandler().Route(g.Group("/api/v1/debitnote"))
	sales_returns.NewHandler().Route(g.Group("/api/v1/sales_returns"))
	purchase_returns.NewHandler().Route(g.Group("/api/v1/purchase_returns"))

	products_handler := products.NewHandler()
	products_handler.Route(g.Group("/api/v1/products"))
	products_handler.Init()
	// products.Init()
	contacts.NewHandler().Route(g.Group("/api/v1/contacts"))
	locations.NewHandler().Route(g.Group("/api/v1/locations"))
	pricing.NewHandler().Route(g.Group("/api/v1/pricing"))
	uom.NewHandler().Route(g.Group("/api/v1/uom"))
	basic_inventory.NewHandler().Route(g.Group("/api/v1/basic_inventory"))
	purchase_orders.NewHandler().Route(g.Group("/api/v1/purchase_orders"))
	delivery_orders.NewHandler().Route(g.Group("/api/v1/delivery_orders"))
	scrap_orders.NewHandler().Route(g.Group("/api/v1/scrap_orders"))
	internal_transfers.NewHandler().Route(g.Group("/api/v1/internal_transfers"))
	shipping_orders.NewHandler().Route(g.Group("/api/v1/shipping_orders"))
	shipping_partners.NewHandler().Route(g.Group("/api/v1/shipping_partners"))
	grn.NewHandler().Route(g.Group("/api/v1/grn"))
	pick_list.NewHandler().Route(g.Group("/api/v1/pick_list"))
	cycle_count.NewHandler().Route(g.Group("/api/v1/cycle_count"))
	vendors.NewHandler().Route(g.Group("/api/v1/vendors"))
	inventory_adjustments.NewHandler().Route(g.Group("/api/v1/inventory_adjustments"))
	shipping_orders_wd.NewHandler().Route(g.Group("/api/v1/shipping_orders_wd"))
	shipping_orders_ndr.NewHandler().Route(g.Group("/api/v1/shipping_orders_ndr"))
	shipping_orders_rto.NewHandler().Route(g.Group("/api/v1/shipping_orders_rto"))
	app_core.NewHandler().Route(g.Group("/api/v1/core"))
	sales_invoice.NewHandler().Route(g.Group("/api/v1/sales_invoice"))
	purchase_invoice.NewHandler().Route(g.Group("/api/v1/purchase_invoice"))
	ipaas.NewIpaasHandler().Route(g.Group("ipaas"))
	ipaas.Init()
	virtual_warehouse.NewHandler().Route(g.Group("/api/v1/virtual_warehouse"))
	catalogue.NewHandler().Route(g.Group("/api/v1/catalogue"))
	scheduler.NewHandler().Route(g.Group("/scheduler"))
	template.NewHandler().Route(g.Group("/api/v1/template"))
	module.NewHandler().Route(g.Group("/api/v1/module"))
	views.NewHandler().Route(g.Group("/api/v1/views"))
	channel.NewHandler().Route(g.Group("api/v1/channels"))
	accounting.NewHandler().Route(g.Group("/api/v1/accounting"))
	pos.NewHandler().Route(g.Group("/api/v1/pos"))
	customers.NewHandler().Route(g.Group("/api/v1/customers"))
	transactions.NewHandler().Route(g.Group("/api/v1/transactions"))
	wallets.NewHandler().Route(g.Group("/api/v1/wallets"))
}
