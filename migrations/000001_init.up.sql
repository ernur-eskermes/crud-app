create table if not exists users
(
    id        UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    username  varchar(255) not null unique,
    password  varchar(255) not null,
    is_active boolean      not null default false
);

create table if not exists book
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    title        varchar(255) not null,
    publish_date timestamp    not null default now(),
    rating       int          not null,
    author_id    UUID         NOT NULL,

    CONSTRAINT author_fk FOREIGN KEY (author_id) REFERENCES users (id)
);