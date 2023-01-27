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
-- Data for Name: apps; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--

INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Inventory Orders', 'INVENTORY_ORDERS',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Inventory Tasks', 'INVENTORY_TASKS',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('MDM', 'MDM',1,  'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Omnichannel', 'OMNICHANNEL',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Orders', 'ORDERS',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Returns', 'RETURNS',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Shipping', 'SHIPPING',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Accounting', 'ACCOUNTING',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Inventory 360', 'INVENTORY_360',1, 'core_app', '1.0.0', true, 1, 1, '[]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
INSERT INTO public.installed_apps (name, code, category_id, category_name, current_version, is_core, created_by, updated_by, app_services) VALUES ('Core DB', 'CORE_DB',13, 'cloud', '1.0.0', true, 1, 1, '["STORAGE"]') ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;


-- INSERT INTO public.installed_apps (name, code, category_name, current_version, is_core, created_by, updated_by) VALUES ('WooCommerce', 'WOOCOMMERCE', 'webstore','1.0.0', false, 1, 1) ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;
-- INSERT INTO public.installed_apps (name, code, category_name, current_version, is_core, created_by, updated_by) VALUES ('Vamaship', 'VAMASHIP', 'shipping','1.0.0', false, 1, 1) ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;

SET session_replication_role = 'origin';

--
-- PostgreSQL database dump complete
--
