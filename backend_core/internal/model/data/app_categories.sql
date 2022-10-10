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
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('WEBSTORES', 'webstores',12) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;
INSERT INTO public.app_categories (category_code, display_name, id) VALUES ('CLOUD', 'cloud',13) ON CONFLICT (id) DO UPDATE SET display_name = EXCLUDED.display_name;

select setval( pg_get_serial_sequence('public.app_categories', 'id'), (select max(id) from public.app_categories));

SET session_replication_role = 'origin';