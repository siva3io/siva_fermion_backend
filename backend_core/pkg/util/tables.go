package util

import (
	scheduler "fermion/backend_core/controllers/scheduler/model"
	"fermion/backend_core/internal/model/accounting"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/inventory_tasks"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/offers"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/orders"
	"fermion/backend_core/internal/model/payments"
	"fermion/backend_core/internal/model/returns"
	"fermion/backend_core/internal/model/shipping"
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
var Tables = map[string]interface{}{

	"migration_version_control": &MigrationVersionControl{},

	//------------------------core------------------------
	"company":                &core.Company{},
	"ondc_details":           &core.OndcDetails{},
	"core_users":             &core.CoreUsers{},
	"custom_solution":        &core.CustomSolution{},
	"tags":                   &core.Tags{},
	"state":                  &core.State{},
	"country":                &core.Country{},
	"currency":               &core.Currency{},
	"l10n":                   &core.L10n{},
	"installed_apps":         &core.InstalledApps{},
	"app_category":           &core.AppCategory{},
	"app_store":              &core.AppStore{},
	"attachments":            &core.Attachments{},
	"lookuptype":             &core.Lookuptype{},
	"lookupcode":             &core.Lookupcode{},
	"access":                 &core.Access{},
	"apps_edit":              &core.AppsEdit{},
	"eunimart_base_settings": &core.EunimartBaseSettings{},
	"platform_points":        &core.PlatformPoints{},
	"notification_templates": &core.NotificationTemplates{},
	"notification_settings":  &core.NotificationSettings{},
	"localization":           &core.Localization{},
	"db_schema":              &core.DBSchema{},
	"external_mapper":        &core.ExternalMapper{},
	"view_schema":            &core.ViewSchema{},
	"channel_lookup_codes":   &core.ChannelLookupCodes{},
	"notifications":          &core.Notifications{},

	//------------------------mdm------------------------
	"product_to_channel_map":             &mdm.ProductToChannelMap{},
	"uom_class":                          &mdm.UomClass{},
	"uom":                                &mdm.Uom{},
	"partner":                            &mdm.Partner{},
	"vendors":                            &mdm.Vendors{},
	"vendor_price_lists":                 &mdm.VendorPriceLists{},
	"product_brand":                      &mdm.ProductBrand{},
	"product_category":                   &mdm.ProductCategory{},
	"product_template":                   &mdm.ProductTemplate{},
	"product_variant":                    &mdm.ProductVariant{},
	"product_selected_attributes_values": &mdm.ProductSelectedAttributesValues{},
	"product_selected_attributes":        &mdm.ProductSelectedAttributes{},
	"product_base_attributes_values":     &mdm.ProductBaseAttributesValues{},
	"product_base_attributes":            &mdm.ProductBaseAttributes{},
	"product_bundles":                    &mdm.ProductBundles{},
	"bundle_line_items":                  &mdm.BundleLineItems{},
	"centralized_basic_inventory":        &mdm.CentralizedBasicInventory{},
	"decentralized_basic_inventory":      &mdm.DecentralizedBasicInventory{},
	"centralized_inventory_transactions": &mdm.CentralizedInventoryTransactions{},
	"user_mdm_fav":                       &mdm.UserMdmFav{},

	//------------------------shared_pricing_and_location------------------------
	"locations":           &shared_pricing_and_location.Locations{},
	"virtual_location":    &shared_pricing_and_location.VirtualLocation{},
	"local_warehouse":     &shared_pricing_and_location.LocalWarehouse{},
	"storage_location":    &shared_pricing_and_location.StorageLocation{},
	"storage_quantity":    &shared_pricing_and_location.StorageQuantity{},
	"retail":              &shared_pricing_and_location.Retail{},
	"office":              &shared_pricing_and_location.Office{},
	"pricing":             &shared_pricing_and_location.Pricing{},
	"sales_price_list":    &shared_pricing_and_location.SalesPriceList{},
	"sales_line_items":    &shared_pricing_and_location.SalesLineItems{},
	"purchase_price_list": &shared_pricing_and_location.PurchasePriceList{},
	"purchase_line_items": &shared_pricing_and_location.PurchaseLineItems{},
	"transfer_price_list": &shared_pricing_and_location.TransferPriceList{},
	"transfer_line_items": &shared_pricing_and_location.TransferLineItems{},

	//------------------------inventory_orders------------------------
	"asn":                        &inventory_orders.ASN{},
	"asn_lines":                  &inventory_orders.AsnLines{},
	"grn":                        &inventory_orders.GRN{},
	"grn_order_lines":            &inventory_orders.GRNOrderLines{},
	"inventory_adjustments":      &inventory_orders.InventoryAdjustments{},
	"inventory_adjustment_lines": &inventory_orders.InventoryAdjustmentLines{},
	"inventory_orders_fav":       &inventory_orders.InventoryOrdersFav{},

	//------------------------inventory_tasks------------------------
	"cycle_count":         &inventory_tasks.CycleCount{},
	"cycle_count_lines":   &inventory_tasks.CycleCountLines{},
	"pick_list":           &inventory_tasks.PickList{},
	"pick_list_lines":     &inventory_tasks.PickListLines{},
	"inventory_tasks_fav": &inventory_tasks.InventoryTasksFav{},

	//------------------------orders------------------------
	"internal_transfers":      &orders.InternalTransfers{},
	"internal_transfer_lines": &orders.InternalTransferLines{},
	"delivery_orders":         &orders.DeliveryOrders{},
	"delivery_order_lines":    &orders.DeliveryOrderLines{},
	"scrap_orders":            &orders.ScrapOrders{},
	"scrap_order_lines":       &orders.ScrapOrderLines{},
	"sales_orders":            &orders.SalesOrders{},
	"sales_order_lines":       &orders.SalesOrderLines{},
	"purchase_orders":         &orders.PurchaseOrders{},
	"purchase_order_lines":    &orders.PurchaseOrderLines{},
	"user_orders_fav":         &orders.UserOrdersFav{},

	//------------------------returns------------------------
	"sales_returns":         &returns.SalesReturns{},
	"sales_return_lines":    &returns.SalesReturnLines{},
	"purchase_returns":      &returns.PurchaseReturns{},
	"purchase_return_lines": &returns.PurchaseReturnLines{},
	"user_returns_fav":      &returns.UserReturnsFav{},

	//------------------------omnichannel------------------------
	"channel":                             &omnichannel.Channel{},
	"user_marketplace_registration":       &omnichannel.User_Marketplace_Registration{},
	"user_marketplace_link":               &omnichannel.User_Marketplace_Link{},
	"marketplace":                         &omnichannel.Marketplace{},
	"user_webstore_link":                  &omnichannel.User_Webstore_Link{},
	"webstore":                            &omnichannel.Webstore{},
	"omnichannel_fav":                     &omnichannel.OmnichannelFav{},
	"user_virtual_warehouse_registration": &omnichannel.User_Virtual_Warehouse_Registration{},
	"user_virtual_warehouse_link":         &omnichannel.User_Virtual_Warehouse_Link{},
	"virtual_warehouse":                   &omnichannel.VirtualWarehouse{},
	"catalogue_template":                  &omnichannel.CatalogueTemplate{},
	"catalogue_template_data":             &omnichannel.CatalogueTemplateData{},
	"omnichannel_field":                   &omnichannel.OmnichannelField{},
	"omnichannel_field_data":              &omnichannel.OmnichannelFieldData{},

	//------------------------shipping------------------------
	"shipping_order":                     &shipping.ShippingOrder{},
	"shipping_order_lines":               &shipping.ShippingOrderLines{},
	"rate_calculator":                    &shipping.RateCalculator{},
	"wd":                                 &shipping.WD{},
	"ndr":                                &shipping.NDR{},
	"ndr_lines":                          &shipping.NDRLines{},
	"rto":                                &shipping.RTO{},
	"shipping_partner":                   &shipping.ShippingPartner{},
	"user_shipping_partner_registration": &shipping.UserShippingPartnerRegistration{},
	"user_shipping_fav_unfav":            &shipping.UserShippingFavUnfav{},

	//------------------------accounting------------------------
	"debit_note":             &accounting.DebitNote{},
	"debit_note_line_items":  &accounting.DebitNoteLineItems{},
	"credit_note":            &accounting.CreditNote{},
	"credit_note_line_items": &accounting.CreditNoteLineItems{},
	"sales_invoice":          &accounting.SalesInvoice{},
	"sales_invoice_lines":    &accounting.SalesInvoiceLines{},
	"purchase_invoice":       &accounting.PurchaseInvoice{},
	"purchase_invoice_lines": &accounting.PurchaseInvoiceLines{},
	"accounting_fav":         &accounting.AccountingFav{},
	"payment_terms":          &accounting.PaymentTerms{},
	"payment_term_details":   &accounting.PaymentTermDetails{},
	"currency_exchange":      &accounting.CurrencyExchange{},
	"break_down_intervals":   &accounting.BreakDownIntervals{},
	"accounting":             &accounting.Accounting{},
	"user_accounting_link":   &accounting.UserAccountingLink{},
	"pos":                    &accounting.Pos{},
	"user_pos_link":          &accounting.UserPosLink{},

	//------------------------meta_table------------------------
	"ir_model":             &core.IRModel{},
	"ir_model_fields":      &core.IRModelFields{},
	"ir_model_constraints": &core.IRModelConstraints{},

	//------------------------scheduler------------------------
	"scheduler_job": &scheduler.SchedulerJob{},
	"scheduler_log": &scheduler.SchedulerLog{},

	//------------------------access------------------------
	"access_template":          &core.AccessTemplate{},
	"core_app_module":          &core.CoreAppModule{},
	"access_module_action":     &core.AccessModuleAction{},
	"views_level_access_items": &core.ViewsLevelAccessItems{},

	//------------------------payments------------------------
	"customers":    &payments.Customers{},
	"transactions": &payments.Transactions{},
	"wallets":      &payments.Wallets{},
	//-------------------ondc-------------------
	"offers":         &offers.Offers{},
	"offers_lines":   &offers.OfferProductDetails{},
	"hsn_codes_data": &mdm.HSNCodesData{},
}
