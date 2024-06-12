CREATE TABLE users
(
    id SERIAL PRIMARY KEY NOT NULL,
    name varchar(200) NOT NULL,
    email varchar(200) NOT NULL UNIQUE,
    password varchar(200) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone
);