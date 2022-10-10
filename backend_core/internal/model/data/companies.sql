--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1
-- Dumped by pg_dump version 14.3

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

INSERT INTO public.companies (id, name, addresses, phone, email, company_details, is_enterpise, parent_id, child_ids, type, company_defaults, plt_point_ids, total_points, constraints, schedulers, queue_services, notification_settings_id, notification_templates_ids, menu_hierarchy, is_enabled, is_active, created_by, updated_date, updated_by, created_date, organization_details, kyc_documents,file_preference_id,invoice_generation_id) VALUES (1, 'Siva', NULL, '9876543210', 'contact@siva3.io', '[]', true, NULL, NULL, NULL, '{"time_zone": {"id": 507, "lookup_code": "IST", "display_name": "Indian Standard Time"}, "date_format": {"id": 504, "lookup_code": "dd/mm/yyyy", "display_name": "dd/mm/yyyy"}, "time_format": {"id": 502, "lookup_code": "12_HOURS", "display_name": "12 hours"}, "access_template_details": {"default_user_template_id": 3}}', NULL, NULL, '[]', '[]', '[]', NULL, NULL, '[]', true, true, 1, '2022-05-30 12:25:57.481189+05:30', 1, '2022-05-30 12:25:57.481189+05:30', '{}','[]',( SELECT public.lookupcodes.id FROM public.lookuptypes, public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'FILE_SERVICES' AND public.lookupcodes.lookup_code = 'LOCAL' ),( SELECT public.lookupcodes.id FROM public.lookuptypes, public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'INVOICE_GENERATION' AND public.lookupcodes.lookup_code = 'ON_SHIPMENT_OF_ORDER' )) ON CONFLICT (name) DO UPDATE SET (addresses, phone, email) = (EXCLUDED.addresses, EXCLUDED.phone, EXCLUDED.email);
-- INSERT INTO public.companies (id, name, addresses, phone, email, company_details, is_enterpise, parent_id, child_ids, type, company_defaults, plt_point_ids, total_points, constraints, schedulers, queue_services, notification_settings_id, notification_templates_ids, menu_hierarchy, is_enabled, is_active, created_by, updated_date, updated_by, created_date, organization_details, kyc_documents,file_preference_id,invoice_generation_id) VALUES (1, 'Eunimart Ltd', NULL, '9876543210', 'contact@eunimart.com', '[]', true, NULL, NULL, NULL, '{"time_zone": {"id": 507, "lookup_code": "IST", "display_name": "Indian Standard Time"}, "date_format": {"id": 504, "lookup_code": "dd/mm/yyyy", "display_name": "dd/mm/yyyy"}, "time_format": {"id": 502, "lookup_code": "12_HOURS", "display_name": "12 hours"}, "access_template_details": {"default_user_template_id": 3}}', NULL, NULL, '[]', '[]', '[]', NULL, NULL, '[]', true, true, 1, '2022-05-30 12:25:57.481189+05:30', 1, '2022-09-28 13:42:16.815379+05:30');

--
-- Name: companies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: eunimartuser
--
select setval( pg_get_serial_sequence('public.companies', 'id'), (select max(id) from public.companies));

SET session_replication_role = 'origin';


--
-- PostgreSQL database dump complete
--

