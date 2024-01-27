CREATE TABLE IF NOT EXISTS users
(
    id       uuid primary key,
    name     varchar(50) not null,
    email    text unique not null,
    password text        not null
);

CREATE TABLE IF NOT EXISTS spaces
(
    id   uuid primary key,
    name varchar(50) not null
);

CREATE TABLE IF NOT EXISTS userspaces
(
    userid  uuid not null references users (id),
    spaceid uuid not null references spaces (id)
);
