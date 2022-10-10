--
-- PostgreSQL database dump
--

select setval( pg_get_serial_sequence('public.channel_lookup_codes', 'id'), (select max(id) from public.channel_lookup_codes));
--
-- PostgreSQL database dump complete
--