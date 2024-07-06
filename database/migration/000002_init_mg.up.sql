CREATE TABLE books (
    id INT NOT NULL AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(13) NOT NULL UNIQUE,
    publication_year INT NOT NULL,
    total_copies INT UNSIGNED NOT NULL,
    available_copies INT UNSIGNED NOT NULL,
    PRIMARY KEY (id)
);