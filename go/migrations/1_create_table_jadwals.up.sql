-- Table: public.jadwals

-- DROP TABLE public.jadwals;

CREATE TABLE public.jadwals
(
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    kegiatan text COLLATE pg_catalog."default",
    jam text COLLATE pg_catalog."default",
    CONSTRAINT jadwals_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.jadwals
    OWNER to postgres;

-- Index: idx_jadwals_id

-- DROP INDEX public.idx_jadwals_id;

CREATE INDEX idx_jadwals_id
    ON public.jadwals USING btree
    (id)
    TABLESPACE pg_default;
