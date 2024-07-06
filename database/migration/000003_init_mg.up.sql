CREATE TABLE transactions (
    transaction_id INT NOT NULL AUTO_INCREMENT,
    user_id INT,
    book_id INT,
    status ENUM('checkout_requested', 'checkout_accepted', 'checkout_rejected', 'checkin_rejected', 'checkin_requested', 'returned') NOT NULL,
    checkout_time DATETIME,
    checkin_time DATETIME,
    PRIMARY KEY (transaction_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);