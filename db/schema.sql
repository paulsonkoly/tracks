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
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


--
-- Name: filestatus; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.filestatus AS ENUM (
    'uploaded',
    'processed',
    'processing_failed'
);


--
-- Name: tracktype; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.tracktype AS ENUM (
    'track',
    'route'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: gpxfiles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.gpxfiles (
    id integer NOT NULL,
    filename text NOT NULL,
    filesize bigint NOT NULL,
    status public.filestatus NOT NULL,
    link text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone NOT NULL
);


--
-- Name: COLUMN gpxfiles.link; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.gpxfiles.link IS 'gpx metadata link field';


--
-- Name: gpxfiles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.gpxfiles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: gpxfiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.gpxfiles_id_seq OWNED BY public.gpxfiles.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: segments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.segments (
    id integer NOT NULL,
    track_id integer NOT NULL,
    geometry public.geometry
);


--
-- Name: segments_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.segments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: segments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.segments_id_seq OWNED BY public.segments.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    token character(43) NOT NULL,
    data bytea NOT NULL,
    expiry timestamp without time zone NOT NULL
);


--
-- Name: tracks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tracks (
    id integer NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    type public.tracktype NOT NULL,
    gpxfile_id integer NOT NULL
);


--
-- Name: tracks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tracks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tracks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tracks_id_seq OWNED BY public.tracks.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    hashed_password character varying(255) NOT NULL,
    created_at timestamp with time zone NOT NULL
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: gpxfiles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gpxfiles ALTER COLUMN id SET DEFAULT nextval('public.gpxfiles_id_seq'::regclass);


--
-- Name: segments id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.segments ALTER COLUMN id SET DEFAULT nextval('public.segments_id_seq'::regclass);


--
-- Name: tracks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks ALTER COLUMN id SET DEFAULT nextval('public.tracks_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: gpxfiles gpxfiles_filename_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gpxfiles
    ADD CONSTRAINT gpxfiles_filename_key UNIQUE (filename);


--
-- Name: gpxfiles gpxfiles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gpxfiles
    ADD CONSTRAINT gpxfiles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: segments segments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.segments
    ADD CONSTRAINT segments_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (token);


--
-- Name: tracks tracks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks
    ADD CONSTRAINT tracks_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: sessions_expiry_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX sessions_expiry_idx ON public.sessions USING btree (expiry);


--
-- Name: segments segments_track_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.segments
    ADD CONSTRAINT segments_track_id_fkey FOREIGN KEY (track_id) REFERENCES public.tracks(id) ON DELETE CASCADE;


--
-- Name: tracks tracks_gpxfile_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracks
    ADD CONSTRAINT tracks_gpxfile_id_fkey FOREIGN KEY (gpxfile_id) REFERENCES public.gpxfiles(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20241024143931'),
    ('20241025092553'),
    ('20241102084123'),
    ('20241107112239'),
    ('20241107141128');
