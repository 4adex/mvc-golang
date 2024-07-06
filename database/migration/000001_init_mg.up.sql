CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role ENUM('client', 'admin') NOT NULL,
    request_status ENUM('pending', 'accepted', 'rejected', 'not_requested') DEFAULT 'not_requested',
    salt VARCHAR(255),
    PRIMARY KEY (id)
);