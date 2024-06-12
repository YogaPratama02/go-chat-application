CREATE TABLE direct_message_details
(
    id SERIAL PRIMARY KEY NOT NULL,
    direct_message_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (direct_message_id) REFERENCES direct_messages (id) ON DELETE CASCADE
);