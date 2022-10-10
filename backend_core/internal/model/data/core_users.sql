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
-- Data for Name: core_users; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--

INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Super User', 'Super', 'User', 'super_user@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[1]', 1, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Company Admin', 'Company', 'User', 'comapany_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[2]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('MDM Admin', 'Mdm', 'Admin', 'mdm_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[3]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Order Admin', 'Order', 'Admin', 'order_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[4]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Inventory Admin', 'Inventory', 'Admin', 'inventory_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[5]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Invnetory Task Admin', 'Invnetory Task', 'Admin', 'inventory_task_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[6]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Accouting Admin', 'Accounting', 'Admin', 'accounting_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[7]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Returns Admin', 'Returns', 'Admin', 'returns_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[8]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('Shipping Admin', 'Shipping', 'Admin', 'shipping_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[9]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);
INSERT INTO public.core_users (name, first_name, last_name, username, email, work_email, mobile_number, password, login_type, auth, fa_conf, device_ids, preferences, tutorials, gamification, total_points, menu_hierarchy, constraints, schedulers, queue_services, access_ids, company_id, profile, plt_point_id, is_enabled, is_active, updated_date, created_date, created_by, updated_by) VALUES ('OmniChannel Admin', 'OmniChannel', 'Admin', 'omni_channel_admin@eunimart.com', NULL, NULL, NULL, '', 0, '{}', '[]', '[]', '[]', '[]', '[]', 0, '{"otp_token": "666666"}', '[]', '[]', '[]', '[10]', NULL, '{}', NULL, true, true, '2022-06-23 15:27:43.884201+05:30', '2022-06-23 15:27:43.884201+05:30', NULL, NULL);

SET session_replication_role = 'origin';

--
-- Name: core_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: eunimartuser
--

select setval( pg_get_serial_sequence('public.core_users', 'id'), (select max(id) from public.core_users));

--
-- PostgreSQL database dump complete
--
