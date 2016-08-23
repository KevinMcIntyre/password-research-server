--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.4
-- Dumped by pg_dump version 9.5.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
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


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: collections; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE collections (
    id bigint NOT NULL,
    label character varying(60) NOT NULL,
    creation_date timestamp without time zone
);


ALTER TABLE collections OWNER TO postgres;

--
-- Name: collections_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE collections_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE collections_id_seq OWNER TO postgres;

--
-- Name: collections_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE collections_id_seq OWNED BY collections.id;


--
-- Name: collections_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE collections_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE collections_seq OWNER TO postgres;

--
-- Name: image_trial_stage_results; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE image_trial_stage_results (
    id bigint NOT NULL,
    trial_id integer NOT NULL,
    stage_id integer NOT NULL,
    selected_image_id integer NOT NULL,
    correct_image_id integer,
    start_time timestamp without time zone,
    end_time timestamp without time zone
);


ALTER TABLE image_trial_stage_results OWNER TO postgres;

--
-- Name: image_trial_stage_results_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE image_trial_stage_results_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE image_trial_stage_results_id_seq OWNER TO postgres;

--
-- Name: image_trial_stage_results_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE image_trial_stage_results_id_seq OWNED BY image_trial_stage_results.id;


--
-- Name: image_trial_stage_results_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE image_trial_stage_results_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE image_trial_stage_results_seq OWNER TO postgres;

--
-- Name: image_trials; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE image_trials (
    id bigint NOT NULL,
    subject_id integer NOT NULL,
    test_config_id integer NOT NULL,
    passed_auth boolean,
    start_time timestamp without time zone,
    end_time timestamp without time zone,
    notes text
);


ALTER TABLE image_trials OWNER TO postgres;

--
-- Name: image_trials_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE image_trials_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE image_trials_id_seq OWNER TO postgres;

--
-- Name: image_trials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE image_trials_id_seq OWNED BY image_trials.id;


--
-- Name: image_trials_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE image_trials_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE image_trials_seq OWNER TO postgres;

--
-- Name: random_stage_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE random_stage_images (
    id bigint NOT NULL,
    image bytea,
    image_type character varying(20) NOT NULL,
    alias character varying(40),
    test_config_id integer NOT NULL,
    stage_number integer NOT NULL,
    row_number integer NOT NULL,
    column_number integer NOT NULL,
    creation_date timestamp without time zone,
    replacement_alias character varying(255)
);


ALTER TABLE random_stage_images OWNER TO postgres;

--
-- Name: random_stage_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE random_stage_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE random_stage_images_id_seq OWNER TO postgres;

--
-- Name: random_stage_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE random_stage_images_id_seq OWNED BY random_stage_images.id;


--
-- Name: random_stage_images_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE random_stage_images_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE random_stage_images_seq OWNER TO postgres;

--
-- Name: saved_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE saved_images (
    id bigint NOT NULL,
    image bytea NOT NULL,
    image_type character varying(20) NOT NULL,
    subject_id integer,
    collection_id integer,
    alias character varying(40),
    creation_date timestamp without time zone
);


ALTER TABLE saved_images OWNER TO postgres;

--
-- Name: saved_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE saved_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE saved_images_id_seq OWNER TO postgres;

--
-- Name: saved_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE saved_images_id_seq OWNED BY saved_images.id;


--
-- Name: saved_images_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE saved_images_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE saved_images_seq OWNER TO postgres;

--
-- Name: subjects; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE subjects (
    id bigint NOT NULL,
    username character varying(30),
    email character varying(30),
    password character varying(60),
    password_entropy double precision,
    password_strength integer,
    pin_number character varying(10),
    first_name character varying(50),
    last_name character varying(50),
    birth_date date,
    creation_date timestamp without time zone,
    notes text
);


ALTER TABLE subjects OWNER TO postgres;

--
-- Name: subjects_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE subjects_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE subjects_id_seq OWNER TO postgres;

--
-- Name: subjects_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE subjects_id_seq OWNED BY subjects.id;


--
-- Name: subjects_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE subjects_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE subjects_seq OWNER TO postgres;

--
-- Name: test_config_stage_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE test_config_stage_images (
    id bigint NOT NULL,
    image bytea,
    image_type character varying(20) NOT NULL,
    stage_id integer NOT NULL,
    alias character varying(40),
    row_number integer NOT NULL,
    column_number integer NOT NULL,
    creation_date timestamp without time zone
);


ALTER TABLE test_config_stage_images OWNER TO postgres;

--
-- Name: test_config_stage_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE test_config_stage_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test_config_stage_images_id_seq OWNER TO postgres;

--
-- Name: test_config_stage_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE test_config_stage_images_id_seq OWNED BY test_config_stage_images.id;


--
-- Name: test_config_stage_images_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE test_config_stage_images_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test_config_stage_images_seq OWNER TO postgres;

--
-- Name: test_config_stages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE test_config_stages (
    id bigint NOT NULL,
    test_config_id integer NOT NULL,
    stage_number integer NOT NULL,
    creation_date timestamp without time zone
);


ALTER TABLE test_config_stages OWNER TO postgres;

--
-- Name: test_config_stages_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE test_config_stages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test_config_stages_id_seq OWNER TO postgres;

--
-- Name: test_config_stages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE test_config_stages_id_seq OWNED BY test_config_stages.id;


--
-- Name: test_configs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE test_configs (
    id bigint NOT NULL,
    name character varying(60),
    rows_in_matrix integer,
    cols_in_matrix integer,
    stage_count integer,
    image_may_not_be_present boolean,
    creation_date timestamp without time zone
);


ALTER TABLE test_configs OWNER TO postgres;

--
-- Name: test_configs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE test_configs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test_configs_id_seq OWNER TO postgres;

--
-- Name: test_configs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE test_configs_id_seq OWNED BY test_configs.id;


--
-- Name: test_configs_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE test_configs_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test_configs_seq OWNER TO postgres;

--
-- Name: uploaded_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE uploaded_images (
    id bigint NOT NULL,
    image bytea NOT NULL,
    image_type character varying(20) NOT NULL,
    subject_id integer,
    collection_id integer,
    alias character varying(40),
    creation_date timestamp without time zone
);


ALTER TABLE uploaded_images OWNER TO postgres;

--
-- Name: uploaded_images_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE uploaded_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE uploaded_images_id_seq OWNER TO postgres;

--
-- Name: uploaded_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE uploaded_images_id_seq OWNED BY uploaded_images.id;


--
-- Name: uploaded_images_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE uploaded_images_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE uploaded_images_seq OWNER TO postgres;

--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY collections ALTER COLUMN id SET DEFAULT nextval('collections_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trial_stage_results ALTER COLUMN id SET DEFAULT nextval('image_trial_stage_results_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trials ALTER COLUMN id SET DEFAULT nextval('image_trials_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY random_stage_images ALTER COLUMN id SET DEFAULT nextval('random_stage_images_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY saved_images ALTER COLUMN id SET DEFAULT nextval('saved_images_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY subjects ALTER COLUMN id SET DEFAULT nextval('subjects_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_config_stage_images ALTER COLUMN id SET DEFAULT nextval('test_config_stage_images_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_config_stages ALTER COLUMN id SET DEFAULT nextval('test_config_stages_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_configs ALTER COLUMN id SET DEFAULT nextval('test_configs_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY uploaded_images ALTER COLUMN id SET DEFAULT nextval('uploaded_images_id_seq'::regclass);


--
-- Data for Name: collections; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY collections (id, label, creation_date) FROM stdin;
\.


--
-- Name: collections_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('collections_id_seq', 2, true);


--
-- Name: collections_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('collections_seq', 1, false);


--
-- Data for Name: image_trial_stage_results; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY image_trial_stage_results (id, trial_id, stage_id, selected_image_id, correct_image_id, start_time, end_time) FROM stdin;
\.


--
-- Name: image_trial_stage_results_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('image_trial_stage_results_id_seq', 1, false);


--
-- Name: image_trial_stage_results_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('image_trial_stage_results_seq', 1, false);


--
-- Data for Name: image_trials; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY image_trials (id, subject_id, test_config_id, passed_auth, start_time, end_time, notes) FROM stdin;
\.


--
-- Name: image_trials_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('image_trials_id_seq', 1, false);


--
-- Name: image_trials_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('image_trials_seq', 1, false);


--
-- Data for Name: random_stage_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY random_stage_images (id, image, image_type, alias, test_config_id, stage_number, row_number, column_number, creation_date, replacement_alias) FROM stdin;
\.


--
-- Name: random_stage_images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('random_stage_images_id_seq', 13239, true);


--
-- Name: random_stage_images_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('random_stage_images_seq', 1, false);


--
-- Data for Name: saved_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY saved_images (id, image, image_type, subject_id, collection_id, alias, creation_date) FROM stdin;
\.


--
-- Name: saved_images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('saved_images_id_seq', 28, true);


--
-- Name: saved_images_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('saved_images_seq', 1, false);


--
-- Data for Name: subjects; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY subjects (id, username, email, password, password_entropy, password_strength, pin_number, first_name, last_name, birth_date, creation_date, notes) FROM stdin;
\.


--
-- Name: subjects_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('subjects_id_seq', 3, true);


--
-- Name: subjects_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('subjects_seq', 1, false);


--
-- Data for Name: test_config_stage_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY test_config_stage_images (id, image, image_type, stage_id, alias, row_number, column_number, creation_date) FROM stdin;
\.


--
-- Name: test_config_stage_images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('test_config_stage_images_id_seq', 2374, true);


--
-- Name: test_config_stage_images_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('test_config_stage_images_seq', 1, false);


--
-- Data for Name: test_config_stages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY test_config_stages (id, test_config_id, stage_number, creation_date) FROM stdin;
\.


--
-- Name: test_config_stages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('test_config_stages_id_seq', 238, true);


--
-- Data for Name: test_configs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY test_configs (id, name, rows_in_matrix, cols_in_matrix, stage_count, image_may_not_be_present, creation_date) FROM stdin;
\.


--
-- Name: test_configs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('test_configs_id_seq', 107, true);


--
-- Name: test_configs_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('test_configs_seq', 1, false);


--
-- Data for Name: uploaded_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY uploaded_images (id, image, image_type, subject_id, collection_id, alias, creation_date) FROM stdin;
\.


--
-- Name: uploaded_images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('uploaded_images_id_seq', 82, true);


--
-- Name: uploaded_images_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('uploaded_images_seq', 1, false);


--
-- Name: collections_label_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY collections
    ADD CONSTRAINT collections_label_key UNIQUE (label);


--
-- Name: pk_collections; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY collections
    ADD CONSTRAINT pk_collections PRIMARY KEY (id);


--
-- Name: pk_custom_stages; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_config_stages
    ADD CONSTRAINT pk_custom_stages PRIMARY KEY (id);


--
-- Name: pk_image_trial_stage_results; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trial_stage_results
    ADD CONSTRAINT pk_image_trial_stage_results PRIMARY KEY (id);


--
-- Name: pk_image_trials; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trials
    ADD CONSTRAINT pk_image_trials PRIMARY KEY (id);


--
-- Name: pk_random_stage_images; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY random_stage_images
    ADD CONSTRAINT pk_random_stage_images PRIMARY KEY (id);


--
-- Name: pk_saved_images; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY saved_images
    ADD CONSTRAINT pk_saved_images PRIMARY KEY (id);


--
-- Name: pk_subjects; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY subjects
    ADD CONSTRAINT pk_subjects PRIMARY KEY (id);


--
-- Name: pk_test_config; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_configs
    ADD CONSTRAINT pk_test_config PRIMARY KEY (id);


--
-- Name: pk_test_config_stage_images; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_config_stage_images
    ADD CONSTRAINT pk_test_config_stage_images PRIMARY KEY (id);


--
-- Name: pk_uploaded_images; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY uploaded_images
    ADD CONSTRAINT pk_uploaded_images PRIMARY KEY (id);


--
-- Name: subjects_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY subjects
    ADD CONSTRAINT subjects_email_key UNIQUE (email);


--
-- Name: test_configs_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_configs
    ADD CONSTRAINT test_configs_name_key UNIQUE (name);


--
-- Name: collection_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX collection_id ON collections USING btree (id);


--
-- Name: image_trial_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX image_trial_id ON image_trials USING btree (id);


--
-- Name: image_trial_result_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX image_trial_result_id ON image_trial_stage_results USING btree (id);


--
-- Name: random_stage_image_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX random_stage_image_id ON random_stage_images USING btree (id);


--
-- Name: saved_image_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX saved_image_id ON saved_images USING btree (id);


--
-- Name: subject_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX subject_id ON subjects USING btree (id);


--
-- Name: test_config_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX test_config_id ON test_configs USING btree (id);


--
-- Name: test_config_stage_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX test_config_stage_id ON test_config_stages USING btree (id);


--
-- Name: test_config_stage_image_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX test_config_stage_image_id ON test_config_stage_images USING btree (id);


--
-- Name: uploaded_image_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX uploaded_image_id ON uploaded_images USING btree (id);


--
-- Name: fk_image_trial_stage_results_correct_image_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trial_stage_results
    ADD CONSTRAINT fk_image_trial_stage_results_correct_image_id_01 FOREIGN KEY (correct_image_id) REFERENCES uploaded_images(id);


--
-- Name: fk_image_trial_stage_results_selected_image_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trial_stage_results
    ADD CONSTRAINT fk_image_trial_stage_results_selected_image_id_01 FOREIGN KEY (selected_image_id) REFERENCES test_config_stage_images(id);


--
-- Name: fk_image_trial_stage_results_stage_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trial_stage_results
    ADD CONSTRAINT fk_image_trial_stage_results_stage_id_01 FOREIGN KEY (stage_id) REFERENCES test_config_stages(id);


--
-- Name: fk_image_trial_stage_results_trial_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trial_stage_results
    ADD CONSTRAINT fk_image_trial_stage_results_trial_id_01 FOREIGN KEY (trial_id) REFERENCES image_trials(id);


--
-- Name: fk_image_trials_subject_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trials
    ADD CONSTRAINT fk_image_trials_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects(id);


--
-- Name: fk_image_trials_test_config_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY image_trials
    ADD CONSTRAINT fk_image_trials_test_config_id_01 FOREIGN KEY (test_config_id) REFERENCES test_configs(id);


--
-- Name: fk_random_stage_images_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY random_stage_images
    ADD CONSTRAINT fk_random_stage_images_01 FOREIGN KEY (test_config_id) REFERENCES test_configs(id);


--
-- Name: fk_saved_images_collection_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY saved_images
    ADD CONSTRAINT fk_saved_images_collection_id_01 FOREIGN KEY (collection_id) REFERENCES collections(id);


--
-- Name: fk_saved_images_subject_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY saved_images
    ADD CONSTRAINT fk_saved_images_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects(id);


--
-- Name: fk_test_config_stage_images_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_config_stage_images
    ADD CONSTRAINT fk_test_config_stage_images_01 FOREIGN KEY (stage_id) REFERENCES test_config_stages(id);


--
-- Name: fk_test_config_stages_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY test_config_stages
    ADD CONSTRAINT fk_test_config_stages_01 FOREIGN KEY (test_config_id) REFERENCES test_configs(id);


--
-- Name: fk_uploaded_images_collection_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY uploaded_images
    ADD CONSTRAINT fk_uploaded_images_collection_id_01 FOREIGN KEY (collection_id) REFERENCES collections(id);


--
-- Name: fk_uploaded_images_subject_id_01; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY uploaded_images
    ADD CONSTRAINT fk_uploaded_images_subject_id_01 FOREIGN KEY (subject_id) REFERENCES subjects(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

