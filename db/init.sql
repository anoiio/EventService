-- Table: public.events

-- DROP TABLE public.events;

CREATE TABLE public.events
(
    id SERIAL PRIMARY KEY,
    created_timestamp timestamptz not null default current_timestamp,
    date_time date NOT NULL,
    type character varying(64) COLLATE pg_catalog."default" NOT NULL,
    transaction_id character varying(64) COLLATE pg_catalog."default" NOT NULL,
    ad_type character varying(64) COLLATE pg_catalog."default",
    time_to_click integer,
    user_id character varying(255) COLLATE pg_catalog."default"
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.events
    OWNER to postgres;
    
    
