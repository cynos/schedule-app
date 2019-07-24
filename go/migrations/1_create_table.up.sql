CREATE TABLE public.jadwals
(
    id bigserial NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    kegiatan text NULL,
    jam text NULL,
    userid integer NOT NULL,
    contact_list integer[],
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
    username varchar(50),
    password varchar(25),
    CONSTRAINT logins_pkey PRIMARY KEY (id)
);

ALTER TABLE public.logins
    OWNER to postgres;

CREATE TABLE public.contacts
(
    id bigserial NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    firstname varchar(50) NULL,
    lastname varchar(50) NUll,
    phone varchar(25) NULL,
    email varchar(50) NULL,
    userid integer NOT NULL,
    CONSTRAINT contacts_pkey PRIMARY KEY (id)
);

ALTER TABLE public.contacts
    OWNER to postgres;
