CREATE TABLE room_chat_messages
(
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL,
    room_chat_id INTEGER NOT NULL,
    message text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (room_chat_id) REFERENCES room_chats (id) ON DELETE CASCADE
);