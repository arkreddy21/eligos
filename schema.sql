CREATE TABLE IF NOT EXISTS users
(
    id       uuid primary key,
    name     varchar(50) not null,
    email    text unique not null,
    password text        not null
);