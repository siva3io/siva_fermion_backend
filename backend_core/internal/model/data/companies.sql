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
-- Data for Name: companies; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--
INSERT INTO public.companies ( id, name, phone, email, website, business_name, business_address, authorised_signatory, authorised_signatory_address, store_name, store_description, established_on, store_timings, seller_apps, ondc_details_id, is_enterpise, parent_id, child_ids, type, company_defaults, plt_point_ids, total_points, constraints, schedulers, queue_services, notification_settings_id, notification_templates_ids, menu_hierarchy, is_enabled, is_active, created_by, updated_date, updated_by, created_date, kyc_documents, file_preference_id, invoice_generation_id ) VALUES ( 1, 'Eunimart Ltd', '9876543210', 'contact@eunimart.com', 'www.eunimart.com', 'Eunimart Omnichannel Private Ltd', 'Hyderabad', 'shayak Mazumder', 'Hyderabad', 'Eunimart BPP', 'Eunimart is a seller PLatform', NULL,'[]', ARRAY['Eunimart BPP'] , 1, true, NULL, NULL, ( SELECT public.lookupcodes.id FROM public.lookuptypes, public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'TYPE_OF_COMPANY' AND public.lookupcodes.lookup_code = 'BPP' ), '{"time_zone": {"id": 507, "lookup_code": "IST", "display_name": "Indian Standard Time"}, "date_format": {"id": 504, "lookup_code": "dd/mm/yyyy", "display_name": "dd/mm/yyyy"}, "time_format": {"id": 502, "lookup_code": "12_HOURS", "display_name": "12 hours"}, "access_template_details": {"default_user_template_id": 3}}', NULL, NULL, '[]', '[]', '[]', NULL, NULL, '[]', true, true, 1, NULL, 1, NULL, '[]', ( SELECT public.lookupcodes.id FROM public.lookuptypes, public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'FILE_SERVICES' AND public.lookupcodes.lookup_code = 'AWS' ), ( SELECT public.lookupcodes.id FROM public.lookuptypes, public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'INVOICE_GENERATION' AND public.lookupcodes.lookup_code = 'ON_SHIPMENT_OF_ORDER' ) ) ON CONFLICT (name) DO UPDATE SET ( phone, email) = ( EXCLUDED.phone, EXCLUDED.email );

-- INSERT INTO public.companies (id, name, addresses, phone, email, company_details, is_enterpise, parent_id, child_ids, type, company_defaults, plt_point_ids, total_points, constraints, schedulers, queue_services, notification_settings_id, notification_templates_ids, menu_hierarchy, is_enabled, is_active, created_by, updated_date, updated_by, created_date, organization_details, kyc_documents,file_preference_id,invoice_generation_id) VALUES (1, 'Eunimart Ltd', NULL, '9876543210', 'contact@eunimart.com', '[]', true, NULL, NULL, NULL, '{"time_zone": {"id": 507, "lookup_code": "IST", "display_name": "Indian Standard Time"}, "date_format": {"id": 504, "lookup_code": "dd/mm/yyyy", "display_name": "dd/mm/yyyy"}, "time_format": {"id": 502, "lookup_code": "12_HOURS", "display_name": "12 hours"}, "access_template_details": {"default_user_template_id": 3}}', NULL, NULL, '[]', '[]', '[]', NULL, NULL, '[]', true, true, 1, NULL, 1, NULL);

--
-- Name: companies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: eunimartuser
--
select setval( pg_get_serial_sequence('public.companies', 'id'), (select max(id) from public.companies));

SET session_replication_role = 'origin';


--
-- PostgreSQL database dump complete
--

