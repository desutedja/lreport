CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    user_level ENUM('admin', 'user') NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);