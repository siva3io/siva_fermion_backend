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
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('CORE_APP', 'core_app',1) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('3PL', '3PL',2) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('CRM', 'CRM',3) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('POS', 'POS',4) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('SOCIAL_COMMERCE', 'social_commerce',5) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('ACCOUNTING', 'accounting',6) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('MARKETPLACE', 'marketplace',7) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('PAYMENT_PARTNER', 'payment_partner',8) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('SHIPPING', 'shipping',9) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('VIRTUAL_WAREHOUSE', 'virtual_warehouse',10) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('WAREHOUSE', 'warehouse',11) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('WEBSTORE', 'webstore',12) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('CLOUD', 'cloud',13) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('LOCAL_WAREHOUSE', 'local_warehouse',14) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('RETAIL', 'retail',15) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;

SET session_replication_role = 'origin';