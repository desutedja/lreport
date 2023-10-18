CREATE TABLE bonus (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    category_id int NOT NULL,
    new_member decimal(10,2),
    cb_sl decimal(10,2),
    rb_sl decimal(10,2),
    cb_ca decimal(10,2),
    roll_ca decimal(10,2),
    cb_sp decimal(10,2),
    rb_sp decimal(10,2),
    refferal decimal(10,2),
    promo decimal(10,2),
    total decimal(10,2),
    trans_date datetime,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_on TIMESTAMP,
    update_by CHAR(36)
);