CREATE DATABASE stub;

USE stub;

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

CREATE TABLE sessions ( 
    token CHAR(43) PRIMARY KEY, 
    data BLOB NOT NULL, 
    expiry TIMESTAMP(6) NOT NULL 
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);