--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1
-- Dumped by pg_dump version 14.3
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
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
SET session_replication_role = 'replica';

--
-- Data for Name: core_app_modules; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (1,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'DASHBOARD', 'Dashboard', NULL, NULL, NULL,'files/icons/dashboard_!.svg', 1);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (2,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'MDM', 'MDM', NULL, NULL, NULL,'files/icons/MDM.svg', 3);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (3,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ORDER', 'Order', NULL, NULL, NULL,'files/icons/orders_1.svg', 4);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (4,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'INVENTORY', 'Inventory', NULL, NULL, NULL,'files/icons/Inventory.svg', 5);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (5,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'INVENTORY_TASK','Inventory Task', NULL, NULL, NULL,'files/icons/Inventory_tasks.svg', 6);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (6,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ACCOUNTING', 'Accounting', NULL, NULL, NULL,'files/icons/Accounting.svg', 7);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (7,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'RETURNS_MANAGEMENT', 'Returns Management', NULL, NULL, NULL,'files/icons/Returns.svg', 8);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (8,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SHIPPING_MANAGEMENT','Shipping Management', NULL, NULL, NULL,'files/icons/local_shipping.svg', 9);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (9,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'OMNICHANNEL_HUB', 'OmniChannel Hub', NULL, NULL, NULL,'files/icons/Omnichannel.svg', 10);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (10,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ADMIN', 'Admin', NULL, NULL, NULL,'files/icons/Admin.svg', 11);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (11,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SETTINGS', 'Settings', NULL, NULL, NULL,'files/icons/MDM.svg', 12);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (12,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SCHEDULERS_LIST', 'Schedulers List', NULL, NULL, NULL,'files/icons/MDM.svg', 13);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (13,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ONDC_APPS', 'ONDC Apps', NULL, NULL, NULL,'files/icons/MDM.svg', 14);



INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (14,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'DASHBOARD', 'Dashboard', 1, 4032,'${PLATFORM_UI_URL}/dashboard','files/icons/Dashboard.svg', 1);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (15,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PRODUCTS', 'Products', 2, 4021,'${PLATFORM_UI_URL}/products','files/icons/products.svg', 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (16,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'CONTACTS', 'Contacts', 2, 4022,'${PLATFORM_UI_URL}/contacts',NULL, 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (17,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'LOCATIONS', 'Locations', 2, 4023,'${PLATFORM_UI_URL}/locations',NULL, 3);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (18,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'UOM', 'UoM', 2, 4001,'${PLATFORM_UI_URL}/uom',NULL, 4);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (19,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PRICING', 'Pricing', 2, 4003,'${PLATFORM_UI_URL}/pricing',NULL, 5);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (20,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'INVENTORY_DECENTRALIZED', 'Inventory-Decentralized', 2, 4018,'${PLATFORM_UI_URL}/inventory',NULL, 6);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (21,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'CENTRALISED_INVENTORY', 'Centralised Inventory', 2, 4030,'${PLATFORM_UI_URL}/cinventory',NULL, 7);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (22,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SALES_ORDERS', 'Sales Orders', 3, 4008,'${PLATFORM_UI_URL}/salesOrders','files/icons/Orders.svg', 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (23,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'IST', 'IST', 3  , 4009,'${PLATFORM_UI_URL}/ist',NULL, 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (24,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SCRAP_ORDERS', 'Scrap Orders', 3, 4002,'${PLATFORM_UI_URL}/scrapOrders','files/icons/Orders.svg', 3);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (25,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'DELIVERY_ORDERS', 'Delivery Orders', 3, 4019,'${PLATFORM_UI_URL}/deliveryOrders','files/icons/Orders.svg', 4);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (26,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PURCHASE_ORDERS', 'Purchase Orders', 3, 4020,'${PLATFORM_UI_URL}/purchaseOrders','files/icons/Orders.svg', 5);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (27,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ASN', 'ASN', 4, 4006,'${PLATFORM_UI_URL}/asn',NULL, 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (28,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'GRN', 'GRN', 4, 4007,'${PLATFORM_UI_URL}/grn',NULL, 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (29,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'INVENTORY_ADJUSTMENT', 'Inventory Adjustment', 4, 4018,'${PLATFORM_UI_URL}/inventoryAdjustment',NULL, 3);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (30,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PICK_LIST', 'Pick List', 5, 4018,'${PLATFORM_UI_URL}/pickList',NULL, 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (31,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'CYCLE_COUNT', 'Cycle Count', 5, 4018,'${PLATFORM_UI_URL}/cycleCount',NULL, 2);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (32,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SALES_INVOICE', 'Sales Invoice', 6, 4024,'${PLATFORM_UI_URL}/salesInvoice',NULL, 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (33,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PURCHASE_INVOICE', 'Purchase Invoice',6, 4024,'${PLATFORM_UI_URL}/purchaseInvoice',NULL, 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (34,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'DEBIT_NOTE', 'Debit Note', 6, 4010,'${PLATFORM_UI_URL}/debitNote',NULL, 3);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (35,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'CREDIT_NOTE', 'Credit Note', 6, 4010,'${PLATFORM_UI_URL}/creditNote',NULL, 4);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (36,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'CURRENCY_EXCHANGE', 'Currency & Exchange', 6, 4010,'${PLATFORM_UI_URL}/currencyExchange',NULL, 5);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (37,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PAYMENT_TERMS_RECORD_PAYMENT', 'Payment Terms & Record Payment', 6, 4010,'${PLATFORM_UI_URL}/paymentTermsRecordPayment',NULL, 6);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (38,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PURCHASE_RETURNS', 'Purchase Return', 7, 4011,'${PLATFORM_UI_URL}/purchaseReturns','files/icons/Returns.svg', 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (39,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SALES_RETURNS', 'Sales Returns', 7, 4012,'${PLATFORM_UI_URL}/salesReturns','files/icons/Returns.svg', 2);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (40,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SHIPPING_ORDERS', 'Shipping Orders', 8, 4013,'${PLATFORM_UI_URL}/shippingOrders','files/icons/Shipping.svg', 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (41,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'NDR', 'NDR', 8, 4014,'${PLATFORM_UI_URL}/ndr',NULL, 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (42,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'RTO', 'RTO', 8, 4014,'${PLATFORM_UI_URL}/rto',NULL, 3);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (43,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'WD', 'WD', 8, 4014,'${PLATFORM_UI_URL}/wd',NULL, 4);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence,external_id)VALUES (44,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'VIRTUAL_WAREHOUSE', 'Virtual Warehouse', 9, 4015,'${PLATFORM_UI_URL}/virtualWarehouse',NULL, 1,(SELECT id from public.app_categories where category_code = 'VIRTUAL_WAREHOUSE'));
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence,external_id)VALUES (45,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'RETAIL', 'Retail', 9, 4015,'${PLATFORM_UI_URL}/retail',NULL, 2,(SELECT id from public.app_categories where category_code = 'RETAIL'));
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence,external_id)VALUES (46,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'WEBSTORES', 'Webstores', 9, 4016,'${PLATFORM_UI_URL}/webstores',NULL, 3,(SELECT id from public.app_categories where category_code = 'WEBSTORE'));
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence,external_id)VALUES (47,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'MARKETPLACES', 'Marketplaces', 9, 4016,'${PLATFORM_UI_URL}/marketplaces',NULL, 4,(SELECT id from public.app_categories where category_code = 'MARKETPLACE'));
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence,external_id)VALUES (48,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'LOGISTICS_PARTNERS', 'Logistics Partners', 9, 4017,'${PLATFORM_UI_URL}/logisticsPartners',NULL, 5,(SELECT id from public.app_categories where category_code = 'SHIPPING'));
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence,external_id)VALUES (49,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'LOCAL_WAREHOUSE', 'Local Warehouse', 9, 4017,'${PLATFORM_UI_URL}/localWarehouse',NULL, 6,(SELECT id from public.app_categories where category_code = 'LOCAL_WAREHOUSE'));
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (50,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'POS', 'pos', 9, 4004, '${PLATFORM_UI_URL}/pos', NULL, 7);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (51,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'OMNI_ACCOUNTING', 'Accounting', 9, 4004, '${PLATFORM_UI_URL}/accounting', NULL, 8);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (52,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'LOOKUP_TYPES', 'Lookup Types', 10, 4027,'${PLATFORM_UI_URL}/lookup-types',NULL, 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (53,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'LOOKUP_CODES', 'Lookup Codes', 10, 4027,'${PLATFORM_UI_URL}/lookup-codes',NULL, 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (54,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ACCESS_TEMPLATES', 'Access Templates', 10, 4027,'${PLATFORM_UI_URL}/access-templates',NULL, 3);




INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (55,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'USER_PROFILE', 'User Profile', 11, 4028, '${PLATFORM_UI_URL}/profile','files/icons/user_profile.svg', 1);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (56,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ORGANISATION', 'Organisation', 11, 4028, '${PLATFORM_UI_URL}/organisation','files/icons/organisation.svg', 2);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (57,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BUSINESS_SETTINGS', 'Business Settings', 11, 4005, '${PLATFORM_UI_URL}/settings', 'files/icons/business_settings.svg', 3);



INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (58,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'SCHEDULER', 'Scheduler', 12, 4029, '${PLATFORM_UI_URL}/scheduler','files/icons/Omnichannel.svg', 1);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (59,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'CREATE_APP', 'Create Apps', 13, 4034, '${PLATFORM_UI_URL}/createApp',NULL,1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (60,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'MY_APPS', 'My Apps', 13, 4034, '${PLATFORM_UI_URL}/myApps', NULL, 2);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (61,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'WHITE_LISTING', 'White listing', 13, 4034, '${PLATFORM_UI_URL}/whitelisting' ,'files/icons/White_list.svg',3);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (62,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ROLES', 'Roles', 13, 4034,'${PLATFORM_UI_URL}/roles','files/icons/Roles.svg', 4);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (63,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'PAYMENTS', 'Payments', 13, 4034,'${PLATFORM_UI_URL}/payments', 'files/icons/Payment.svg', 5);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (64,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BAP_ORDER_STATUS', 'Order Status', NULL, NULL,'${PLATFORM_UI_URL}/baporderStatus', 'files/icons/orderStatus.svg', 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (65,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BAP_PAYOUT_DETAIL', 'Payout Detail', NULL, NULL,'${PLATFORM_UI_URL}/payoutDetail', 'files/icons/payoutDetails.svg', 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (66,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BAP_CUSTOMER_SUPPORT', 'Customer Support', NULL, NULL,'${PLATFORM_UI_URL}/bapcustomerSupport', 'files/icons/customerSupport.svg', 3);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (67,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_ORDER_STATUS', 'Order Status', NULL, NULL,'${PLATFORM_UI_URL}/bpporderStatus', 'files/icons/orderStatus.svg', 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (68,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_CUSTOMER_SUPPORT', 'Customer Support', NULL, NULL,'${PLATFORM_UI_URL}/bppcustomerSupport', 'files/icons/customerSupport.svg', 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (69,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_LSP', 'LSP', NULL, NULL,'${PLATFORM_UI_URL}/lsp', 'files/icons/LSP.svg', 3);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (70,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_PRODUCTS', 'Products', NULL, NULL,'${PLATFORM_UI_URL}/products', 'files/icons/products.svg', 4);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (71,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_SHIPPING', 'Shipping', NULL, NULL,'${PLATFORM_UI_URL}/shippingOrders', 'files/icons/Shipping.svg', 5);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (72,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_PAYMENTS', 'Payments', NULL, NULL,'${PLATFORM_UI_URL}/paymentsAdmin', 'files/icons/Payment.svg', 6);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (73,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_ORDERS', 'Orders', NULL, NULL,'${PLATFORM_UI_URL}/salesOrders', 'files/icons/Orders.svg', 7);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (74,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_PAYMENT_PARTNER_SETUP', 'Payment Partner Setup', NULL, NULL,'${PLATFORM_UI_URL}/paymentPartners', 'files/icons/paymentPartnersetup.svg', 8);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (75,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_SHIPPING_PARTNER_SETUP', 'Shipping Partner Setup', NULL, NULL,'${PLATFORM_UI_URL}/shippingPartners', 'files/icons/Shipping.svg', 9);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (76,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ONDC', 'Ondc', NULL, NULL, NULL,'files/icons/MDM.svg', 2);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (77,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ONDC_PRODUCTS', 'Products', 76, NULL,'${PLATFORM_UI_URL}/ondcProducts', NULL, 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (78,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ONDC_ORDERS', 'Orders', 76, NULL,'${PLATFORM_UI_URL}/ondcOrders', NULL, 2);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (79,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'ONDC_SUPPORT', 'Support', 76, NULL,'${PLATFORM_UI_URL}/bppSupport', NULL, 3);

INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (80,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BAP_MANAGE_USER', 'Manage User', NULL, NULL,'${PLATFORM_UI_URL}/manageUser', 'files/icons/manageUser.svg', 4);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (81,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'BPP_MANAGE_USER', 'Manage User', NULL, NULL,'${PLATFORM_UI_URL}/manageUser', 'files/icons/manageUser.svg', 10);


INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (82,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'OFFERS', 'Offers', NULL, NULL, NULL,'files/icons/MDM.svg', 15);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (83,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'OFFERS_MODULE', 'Offers Module', 82, 4021,'${PLATFORM_UI_URL}/promotions','files/icons/products.svg', 1);
INSERT INTO public.core_app_modules(id, is_enabled, is_active, created_by, updated_date, updated_by, deleted_by, created_date, company_id, app_id, deleted_at, module_code, module_name, parent_module, port_number, route_path, image_option, item_sequence)VALUES (84,true, true, 1, NULL, 1, NULL, NULL, 1, NULL, NULL, 'OFFERS_HISTORY', 'Offers History', 82, 4022,'${PLATFORM_UI_URL}/promotionHistory',NULL, 2);




SET session_replication_role = 'origin';


--
-- Name: core_app_modules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: eunimartuser
--

select setval( pg_get_serial_sequence('public.core_app_modules', 'id'), (select max(id) from public.core_app_modules));


--
-- PostgreSQL database dump complete
--