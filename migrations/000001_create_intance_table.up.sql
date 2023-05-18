CREATE TABLE IF NOT EXISTS instances (
    id bigserial NOT NULL PRIMARY KEY,
    instance_id text NOT NULL,
    instance_name text NOT NULL,
    name text NOT NULL,
    module text NOT NULL,
    running boolean DEFAULT false,
    suspended boolean DEFAULT false,
    cpu_usage smallint NOT NULL,
    cpu_max smallint NOT NULL,
    cpu_percent smallint NOT NULL,
    memory_usage bigint NOT NULL,
    memory_max bigint NOT NULL,
    memory_percent smallint NOT NULL,
    users_active smallint NOT NULL,
    users_max smallint NOT NULL
);