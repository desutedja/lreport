CREATE TABLE login_history (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    device VARCHAR(255) NULL,
    ip_address VARCHAR(50) NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);