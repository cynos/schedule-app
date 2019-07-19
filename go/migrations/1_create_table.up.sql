CREATE TABLE public.jadwals
(
    id bigserial NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    kegiatan text NULL,
    jam text NULL,
    userid integer NULL,
    CONSTRAINT jadwals_pkey PRIMARY KEY (id)
);

ALTER TABLE public.jadwals
    OWNER to postgres;

CREATE TABLE public.logins
(
    id bigserial NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NULL,
    username varchar(50),
    password varchar(25),
    CONSTRAINT logins_pkey PRIMARY KEY (id)
);

ALTER TABLE public.logins
    OWNER to postgres;
