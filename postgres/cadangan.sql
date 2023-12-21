CREATE TABLE cadangan_types (
    id serial PRIMARY KEY,
    name VARCHAR(12)
);
insert into cadangan_types(name) values ('Aduan');
insert into cadangan_types(name) values ('Cadangan');
insert into cadangan_types(name) values ('Lain-lain');

CREATE TABLE cadangan (
    id serial PRIMARY KEY,
    cadangan_types_id INTEGER REFERENCES cadangan_types(id)
    cadangan_text VARCHAR(1024) IS NOT NULL,
    tindakan_text VARCHAR(1024),
    is_open BOOLEAN DEFAULT false,
    score smallint,
    create_date int8 NULL
);
