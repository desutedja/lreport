CREATE TABLE transaction (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    category_id int NOT NULL,
    rgs int NOT NULL, -- registration
    rd int NOT NULL, -- registration_deposit
    ap int NOT NULl, -- active_player
    
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_on TIMESTAMP,
    update_by CHAR(36)
);  