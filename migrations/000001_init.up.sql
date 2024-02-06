CREATE TABLE if not exists users
(
    id serial primary key,
    username varchar(255) not null unique,
    password_hash varchar(255) not null,
    created_at timestamp not null default now()
);

CREATE TABLE if not exists records
(
    id serial primary key,
    user_id int references users(id) on delete cascade not null,
    name varchar(255) not null unique,
    site varchar(255),
    login varchar(255),
    password_hash varchar(255),
    info text,
    created_at timestamp not null default now()
);