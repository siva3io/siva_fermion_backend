--
-- PostgreSQL database dump
--
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
-- Data for Name: ondc_details; Type: TABLE DATA; Schema: public; Owner: eunimartuser
--
INSERT INTO public.ondc_details ( id,subscriber_id,subscriber_url, signing_public_key, signing_private_key,encryption_private_key, encryption_public_key,unique_id,type, buyer_app_finder_fee_type, buyer_app_finder_fee_amount, is_collector ) VALUES ( 1,'ondc.eunimart.com','https://ondc.eunimart.com/api/v1/ondc/bpp/eunimart_bpp/','SigningPublicKey','SigningPrivateKey','EncryptionPrivateKey','EncryptionPublicKey', 'UniqueId',( SELECT public.lookupcodes.id FROM public.lookuptypes, public.lookupcodes WHERE public.lookuptypes.id = public.lookupcodes.lookup_type_id AND public.lookuptypes.lookup_type = 'TYPE_OF_COMPANY' AND public.lookupcodes.lookup_code = 'BPP' ), 'Amount', 10, true );

--
-- Name: ondc_details_id_seq; Type: SEQUENCE SET; Schema: public; Owner: eunimartuser
--
select setval( pg_get_serial_sequence('public.ondc_details', 'id'), (select max(id) from public.ondc_details));

SET session_replication_role = 'origin';


--
-- PostgreSQL database dump complete
--

