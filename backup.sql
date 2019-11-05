--
-- PostgreSQL database dump
--

-- Dumped from database version 10.10
-- Dumped by pg_dump version 10.10

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

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: sex_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.sex_type AS ENUM (
    'male',
    'female'
);


ALTER TYPE public.sex_type OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: message; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.message (
    id bigint NOT NULL,
    sender character varying(255) NOT NULL,
    "time" timestamp with time zone NOT NULL,
    body text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.message OWNER TO postgres;

--
-- Name: message_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.message_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.message_id_seq OWNER TO postgres;

--
-- Name: message_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.message_id_seq OWNED BY public.message.id;


--
-- Name: receiver; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.receiver (
    id bigint NOT NULL,
    mailid bigint,
    email character varying(255) NOT NULL
);


ALTER TABLE public.receiver OWNER TO postgres;

--
-- Name: receiver_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.receiver_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.receiver_id_seq OWNER TO postgres;

--
-- Name: receiver_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.receiver_id_seq OWNED BY public.receiver.id;


--
-- Name: session; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.session (
    login character varying(64) NOT NULL,
    token uuid NOT NULL
);


ALTER TABLE public.session OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    login character varying(64) NOT NULL,
    password bytea NOT NULL,
    sault bytea NOT NULL,
    avatar character varying(255) DEFAULT 'default.png'::character varying,
    firstname character varying(255),
    secondname character varying(255),
    sex public.sex_type,
    birthdate date
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: message id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message ALTER COLUMN id SET DEFAULT nextval('public.message_id_seq'::regclass);


--
-- Name: receiver id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receiver ALTER COLUMN id SET DEFAULT nextval('public.receiver_id_seq'::regclass);


--
-- Data for Name: message; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.message (id, sender, "time", body) FROM stdin;
1	admin@mail.ru	2019-12-21 12:31:00+03	<div>Hello</div>
2	admin@mail.ru	2019-12-21 12:31:00+03	<div>Hello</div>
3	admin@mail.ru	2019-12-21 12:31:00+03	<div>Hello</div>
4	slave@mail.ru	2019-12-21 12:31:00+03	<div>Hello</div>
5	ivan	0001-01-01 02:30:17+02:30:17	
6	aa@aa.aa	1900-01-01 00:00:00+02:30:17	Hell
12	ivanov.vanya.111@mail.ru	0001-01-01 02:30:17+02:30:17	\n<HTML><BODY>Hello<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
13	ivanov.vanya.111@mail.ru	0001-01-01 02:30:17+02:30:17	\n<HTML><BODY>Hello<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
15	ivanov.vanya.111@mail.ru	0001-01-01 02:30:17+02:30:17	\n<HTML><BODY>Hello<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
16	aa@aa.aa	1900-01-01 00:00:00+02:30:17	Hell
17	a@a.a	1000-01-01 00:00:00+02:30:17	H
23	ivanov.vanya.111@mail.ru	0001-01-01 02:30:17+02:30:17	\n<HTML><BODY>Hello<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
24	ivanov.vanya.111@mail.ru	0001-01-01 02:30:17+02:30:17	\n<HTML><BODY>Hello<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
25	ivanov.vanya.111@mail.ru	2019-11-04 13:46:46+03	\n<HTML><BODY>Hello<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
26	saf	2006-01-02 12:34:12+03	body
27	ivanov.vanya.111@mail.ru	2019-11-04 14:07:15+03	\n<HTML><BODY>Hello<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
28	andreykochnov@yandex.ru	2019-11-05 10:43:50+03	<div>Bye</div><div> </div><div>------------------------------------------</div><div>С уважением,</div><div>Андрей К.</div><div> </div>\n\n\n
29	andreykochnov@yandex.ru	2019-11-05 10:46:50+03	<div>Bye</div><div> </div><div>------------------------------------------</div><div>С уважением,</div><div>Андрей К.</div><div> </div>\n\n\n
30	ivanov.vanya.111@mail.ru	2019-11-05 17:19:44+03	\n<HTML><BODY>Hello, world<br><br><br>-- <br>Иван Кочубей</BODY></HTML>\n
\.


--
-- Data for Name: receiver; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.receiver (id, mailid, email) FROM stdin;
1	1	ivan@nlmail.ddns.net
2	4	ivanov@nlmail.ddns.net
3	23	aaa
4	24	aaa@nlmail.ddns.net
5	25	aaa@nlmail.ddns.net
6	27	aaa@nlmail.ddns.net
7	28	aaa@nlmail.ddns.net
8	29	aaa@nlmail.ddns.net
9	30	aaa@nlmail.ddns.net
\.


--
-- Data for Name: session; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.session (login, token) FROM stdin;
testhash	08817441-ffbd-11e9-b094-98fa9b864510
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (login, password, sault, avatar, firstname, secondname, sex, birthdate) FROM stdin;
ivanovo	\\x3132333435	\\x		Ivan	Ivanov	male	1984-01-20
eladminoderussiafuckall	\\x6c69636b6d79617373	\\x		El Admino	De Rusia	female	2020-12-12
matsu	\\x31323334	\\x		Ma	Tsu	female	2019-03-31
aaa	\\x3132333435	\\x		AA	aa	male	1988-01-01
ian	\\x3132333435	\\x		I	I	male	1000-01-01
testhash	\\xf8dbddbc1325e451564e65f17c755db9de959ce6402064ef1f94ee5ea53c135b	\\x7465737468617368d41d8cd98f00b204e9800998ecf8427e		MyName	Ivanov	male	1432-12-12
\.


--
-- Name: message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.message_id_seq', 30, true);


--
-- Name: receiver_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.receiver_id_seq', 9, true);


--
-- Name: message message_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_pkey PRIMARY KEY (id);


--
-- Name: receiver receiver_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receiver
    ADD CONSTRAINT receiver_pkey PRIMARY KEY (id);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (login);


--
-- Name: receiver receiver_mailid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.receiver
    ADD CONSTRAINT receiver_mailid_fkey FOREIGN KEY (mailid) REFERENCES public.message(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

