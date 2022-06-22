create table if not exists users (
    id uuid primary key not null,
    first_name varchar(255),
    last_name varchar(255)

);