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
-- Data for Name: currencies; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--

INSERT INTO public.currencies (name, currency_symbol, currency_code, is_base_currency, exchange_rate, exchange_rate_history, auto_update_er, id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at) VALUES ('INR', '₹', 'INR', true, NULL, NULL, NULL, 1, true, true, 1, '2022-05-30 12:25:57.676699+05:30', 1, '2022-05-30 12:25:57.676699+05:30', NULL, NULL) ON CONFLICT (currency_code) DO UPDATE SET (name, currency_symbol) = (EXCLUDED.name, EXCLUDED.currency_symbol);
INSERT INTO public.currencies (name, currency_symbol, currency_code, is_base_currency, exchange_rate, exchange_rate_history, auto_update_er, id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at) VALUES ('USD', '$', 'USD', true, NULL, NULL, NULL, 2, true, true, 1, '2022-05-30 12:25:57.678755+05:30', 1, '2022-05-30 12:25:57.678755+05:30', NULL, NULL) ON CONFLICT (currency_code) DO UPDATE SET (name, currency_symbol) = (EXCLUDED.name, EXCLUDED.currency_symbol);
INSERT INTO public.currencies (name, currency_symbol, currency_code, is_base_currency, exchange_rate, exchange_rate_history, auto_update_er, id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at) VALUES ('NPR', 'रू', 'NPR', true, NULL, NULL, NULL, 3, true, true, 1, '2022-05-30 12:25:57.678755+05:30', 1, '2022-05-30 12:25:57.678755+05:30', NULL, NULL) ON CONFLICT (currency_code) DO UPDATE SET (name, currency_symbol) = (EXCLUDED.name, EXCLUDED.currency_symbol);
INSERT INTO public.currencies (name, currency_symbol, currency_code, is_base_currency, exchange_rate, exchange_rate_history, auto_update_er, id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at) VALUES ('SGD', '$', 'SGD', true, NULL, NULL, NULL, 4, true, true, 1, '2022-05-30 12:25:57.678755+05:30', 1, '2022-05-30 12:25:57.678755+05:30', NULL, NULL) ON CONFLICT (currency_code) DO UPDATE SET (name, currency_symbol) = (EXCLUDED.name, EXCLUDED.currency_symbol);
INSERT INTO public.currencies (name, currency_symbol, currency_code, is_base_currency, exchange_rate, exchange_rate_history, auto_update_er, id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at) VALUES ('MYR', 'RM', 'MYR', true, NULL, NULL, NULL, 5, true, true, 1, '2022-05-30 12:25:57.678755+05:30', 1, '2022-05-30 12:25:57.678755+05:30', NULL, NULL) ON CONFLICT (currency_code) DO UPDATE SET (name, currency_symbol) = (EXCLUDED.name, EXCLUDED.currency_symbol);
INSERT INTO public.currencies (name, currency_symbol, currency_code, is_base_currency, exchange_rate, exchange_rate_history, auto_update_er, id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at) VALUES ('AED', 'د.إ', 'AED', true, NULL, NULL, NULL, 6, true, true, 1, '2022-05-30 12:25:57.678755+05:30', 1, '2022-05-30 12:25:57.678755+05:30', NULL, NULL) ON CONFLICT (currency_code) DO UPDATE SET (name, currency_symbol) = (EXCLUDED.name, EXCLUDED.currency_symbol);
INSERT INTO public.currencies (name, currency_symbol, currency_code, is_base_currency, exchange_rate, exchange_rate_history, auto_update_er, id, is_enabled, is_active, created_by, updated_date, updated_by, created_date, company_id, deleted_at) VALUES ('AUD','$', 'AUD', true, NULL, NULL, NULL, 7, true, true, 1, '2022-05-30 12:25:57.678755+05:30', 1, '2022-05-30 12:25:57.678755+05:30', NULL, NULL) ON CONFLICT (currency_code) DO UPDATE SET (name, currency_symbol) = (EXCLUDED.name, EXCLUDED.currency_symbol);

SET session_replication_role = 'origin';
--
-- Name: currencies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: eunimartuser
--

select setval( pg_get_serial_sequence('public.currencies', 'id'), (select max(id) from public.currencies));

--
-- PostgreSQL database dump complete
--

