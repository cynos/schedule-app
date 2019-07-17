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

-- Index: idx_jadwals_deleted_at

-- DROP INDEX public.idx_jadwals_deleted_at;

CREATE INDEX idx_jadwals_deleted_at
    ON public.jadwals USING btree
    (deleted_at)
    TABLESPACE pg_default;
