CREATE TABLE room_chats
(
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    room_name varchar(200) NOT NULL,
    status boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);