CREATE TABLE user_preferences (
    id         BINARY(16)  NOT NULL PRIMARY KEY,
    user_id    BINARY(16)  NOT NULL UNIQUE,
    theme      VARCHAR(16) NOT NULL,
    created_at DATETIME    NOT NULL,
    updated_at DATETIME    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
