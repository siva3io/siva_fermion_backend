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
-- Data for Name: lookuptypes; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('GENDER', 'Gender') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PAYMENT_TYPE', 'Payment type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PAYMENT_TERMS', 'Payment terms') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PAYMENT_STATUS', 'Payment status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('WEIGHT_TYPE', 'Weight type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CONTACT_TYPE', 'Contact type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CONTACT_PROPERTIES', 'Contact properties') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('LOCATION_TYPE', 'Location type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('INVENTORY_TYPE', 'Inventory type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRICE_INCLUDES_TAX', 'Price Includes Tax') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('LOCATION_SPACE_TYPE', 'Location space type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('UOM_ITEM_TYPE', 'Uom item type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('UOM_CONVERSION_TYPE', 'Uom conversion type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ASN_STATUS', 'Asn Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ASN_SOURCE_DOCUMENT_TYPES', 'Asn Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('GRN_SOURCE_DOCUMENT_TYPES', 'Grn Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DELIVERY_ORDERS_SOURCE_DOCUMENT_TYPES', 'Deliver Orders Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PURCHASE_ORDERS_DOCUMENT_TYPES', 'Purchase Orders Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('IST_SOURCE_DOCUMENT_TYPES', 'Ist Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SCRAP_ORDERS_SOURCE_DOCUMENT_TYPES', 'Scrap Orders Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('INVENTORY_ADJUSTMENTS_SOURCE_DOCUMENT_TYPES', 'Inventory Adjustments Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SHIPPING_ORDERS_SOURCE_DOCUMENT_TYPES', 'Shipping Orders Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SALES_INVOICE_SOURCE_DOCUMENT_TYPES', 'Sales Inovice Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PURCHASE_ORDERS_SOURCE_DOCUMENT_TYPES', 'Purchase Orders Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PURCHASE_INVOICE_SOURCE_DOCUMENT_TYPES', 'Purchase Invoice Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SALES_RETURNS_SOURCE_DOCUMENT_TYPES', 'Sales Returns Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PURCHASE_RETURNS_SOURCE_DOCUMENT_TYPES', 'Purchase Returns Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CREDIT_NOTE_SOURCE_DOCUMENT_TYPES', 'Credit Note Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DEBIT_NOTE_SOURCE_DOCUMENT_TYPES', 'Debit Note Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PICK_LIST_SOURCE_DOCUMENT_TYPES', 'Pick List Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CYCLE_COUNT_SOURCE_DOCUMENT_TYPES', 'Cycle Count Source Document Types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ADJUSTMENT_TYPE', 'Adjustment Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ADJUSTMENT_STATUS', 'Adjustment Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ADJUSTMENT_REASON', 'Adjustment Reason') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PICK_LIST_SOURCE_DOCUMENT_TYPE', 'Pick List Source Document Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PICK_LIST_STATUS', 'Pick List Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CYCLE_COUNT_STATUS', 'Cycle Count Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CYCLE_COUNT_METHOD', 'Cycle Count Method') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('RETURN_TYPE', 'Return Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('BILLING_STATUS', 'Billing Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('INVOICE_STATUS', 'Invoice Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('RECEIPT_ROUTING_OPTIONS', 'Receipt Routing Options') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SHIPPING_MODE', 'Shipping Mode') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('RETURN_STATUS', 'Return Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('IST_REASON', 'Ist Reason') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('IST_STATUS', 'Ist Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ORDER_STATUS', 'Order Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('STANDARD_PRODUCT_TYPE', 'Standard Product Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('VARIANT_TYPE', 'variant Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_CONDITION', 'Product Condition') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PACKAGING_TYPE', 'packaging Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('INVENTORY_TRACKING', 'Inventory Tracking') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_TYPE', 'PRODUCT TYPE') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_CATEGORY_TYPE', 'Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_PROCUREMENT_TREATMENT', 'Product Procurement Treatment') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_STATUS', 'Product Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_TEMP_STATUS', 'Product Temp Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('STOCK_TREATMENT', 'Stock Treatment') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PACKAGING_MATERIAL', 'Packaging Material') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SOURCE_DOCUMENT_TYPES', 'Source Document types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('GRN_STATUS', 'Grn status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('MARKETPLACE_TYPE_OF_SOURCE', 'Marketplace Type Of Source') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('WEBSTORE_TYPE_OF_SOURCE', 'Webstore Type Of Source') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SCRAPPING_REASON', 'Scrapping Reason') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SCRAP_ORDER_STATUS', 'Scrap Order Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('RETAIL_PRICE_INCLUDES_TAX', 'Retail Price Includes Tax') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SHIPPING_ORDER_STATUS', 'Shipping order status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PACKAGE_DIRECTIONS', 'Package directions') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SHIPPING_PARTNER_TYPE', 'Shipping partner type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SHIPPING_TYPE', 'Shipping type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SALES_ORDER_STATUS', 'Sales order status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PURCHASE_ORDER_STATUS', 'Purchase order status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DEBIT_STATUS', 'Debit status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CREDIT_STATUS', 'Credit status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('RETURN_REASON', 'Return reason') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DELIVERY_ORDER_STATUS', 'Delivery order status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DELIVERY_ORDER_PAYMENT_TYPES', 'Delivery order payment types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DELIVERY_ORDER_SHIPPING_CARRIER', 'Delivery order shipping carrier') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('QUANTITY_VALUE_TYPE', 'Quantity value type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PURCHASE_INVOICE_STATUS','Purchase invoice status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRICING_OPTIONS','Pricing options') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SELECT_CHANNEL_OF_SALE','Select channel of sale') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('TERMS_OF_SALES','Terms of sales') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('USER_SHIPPING_REGISTRATION_STATUS','User Shipping Registration Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PACKAGING_METHOD','Packaging method') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CYCLE_COUNT_FREQUENCY','Cycle Count Frequency') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SALES_INVOICE_STATUS','Sales Invoice Status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DISCOUNT_TYPE','Discount Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('TAX_TYPE','Tax Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRICING_STATUS','Pricing status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('RTO_SHIPPING_ORDER_STATUS','RTO SHIPPING ORDER STATUS') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DEBIT_NOTE_STATUS','Debit note status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CREDIT_NOTE_STATUS','Credit note status') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PAYMENT_TERM_TYPE','Payment term type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PAYMENT_ACCOUNT_TYPE','Payment account type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('EXCHANGE','Exchange') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ACCESS_MANAGEMENT_MODULES','Access management modules') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CHANNEL_TYPE','Channel type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CREDIT_NOTE_REASONS','Credit note reasons') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DEBIT_NOTE_REASONS','Debit note reasons') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PURCHASE_ORDERS_DELIVERY_TO','Purchase Orders Delivery To') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('STORAGE_LOCATION_TYPE','Storage Location type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('TIME_FORMAT','Time format') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DATE_FORMAT','Date format') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('TIME_ZONE','Time zone') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('APP_SERVICES','App services') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('INVOICE_GENERATION', 'Invoice Generation') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('FILE_SERVICES', 'File Services') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('TYPE_OF_BUSINESS', 'Type of Business') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CONCURRENCY_TYPE', 'concurrency types') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('WEBSTORE_FUNCTION_TYPE', 'Webstore Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('MARKETPLACE_FUNCTION_TYPE', 'Marketplace Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('VIRTUAL_WAREHOUSE_FUNCTION_TYPE', 'Virtual Warehouse Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SHIPPING_FUNCTION_TYPE', 'Shipping Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('LOCAL_WAREHOUSE_FUNCTION_TYPE', 'Local Warehouse Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('POS_FUNCTION_TYPE', 'Pos Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('RETAIL_FUNCTION_TYPE', 'Retail Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ACCOUNTING_FUNCTION_TYPE', 'Accounting Function Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'TYPE_OF_COMPANY', 'Type of Company' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'TYPE_OF_USER', 'Type of User' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'STD_CODES', 'STD Codes' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'TYPE_OF_DOMAIN', 'Type of Domain' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SCHEDULER_FREQUENCY', 'Scheduler Frequency') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CANCEL_REASON', 'Cancel Reason ') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;

INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'NOTIFICATION_TYPE', 'Notification Type' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'NOTIFICATION_BROADCAST_SCOPE', 'Notification Broadcast Scopes' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'NOTIFICATION_EVENT', 'Notification_Event' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ( 'FEATURE_LIST', 'Feature List' ) ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_ESTIMATED_DELIVERY_TIME', 'Product estimated delivery time') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_FULFILLMENT', 'Product Fulfillment') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('TIME_WITHIN', 'Time Within') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ANSWER_TYPE', 'Answer Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SELLER_NP_TYPE', 'Seller Np Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('ORDER_CATEGORY', 'Order Category') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DELIVERY_TYPE', 'Delivery Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_FULFILLMENT_MANAGED_BY', 'Product Fulfillment Managed By') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;

INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('PRODUCT_DOMAIN', 'Product Domain') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;

INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('FOOD_TYPE', 'Food Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('DELIVERY_TYPE_PREFERENCES', 'Delivery type Preferences') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('OTP_PREFERENCES', 'OTP Preferences') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;

INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('OFFERS_DISCOUNT_TYPE', 'Offers Discunt Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('OFFERS_TERMS_AND_CONDTONS', 'Offers Terms And Conditions') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;

INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('MONTHS', 'Months') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;


INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('SHIPPING_FULFILLMENT_TYPE', 'Shipping Fulfillment Type') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;

INSERT INTO public.lookuptypes (lookup_type, display_name) VALUES ('CANCELLED_BY', 'Cancelled By') ON CONFLICT (lookup_type) DO UPDATE SET display_name = EXCLUDED.display_name;

SET session_replication_role = 'origin';
-- PostgreSQL database dump complete
--