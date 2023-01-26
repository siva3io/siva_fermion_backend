--
-- PostgreSQL database dump
--

-- Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
-- All rights reserved.
-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU Lesser General Public License v3.0 as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.
-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU Lesser General Public License v3.0 for more details.
-- You should have received a copy of the GNU Lesser General Public License v3.0
-- along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.

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
-- Data for Name: states; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--

-- Country: India

INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (1, 1);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (2, 2);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (3, 3);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (4, 4);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (5, 5);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (6, 6);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (7, 7);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (8, 8);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (9, 9);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (10, 10);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (11, 1);
INSERT INTO public.user_template_mapper( core_users_id, access_template_id) VALUES (11, 2);




SET session_replication_role = 'origin';


--
-- Name: user_template_mapper_id; Type: SEQUENCE SET; Schema: public; Owner: eunimartuser
--

-- select setval( pg_get_serial_sequence('public.user_template_mapper', 'id'), (select max(id) from public.user_template_mapper));


--
-- PostgreSQL database dump complete
--