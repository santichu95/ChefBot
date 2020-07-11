use chefbot;

CREATE TABLE IF NOT EXISTS Currency (
    ID varchar(100) NOT NULL,
    Value int UNSIGNED,
    PRIMARY KEY (ID)
);

CREATE TABLE IF NOT EXISTS Users (
    ID varchar(100) NOT NULL,
    Username varchar(100) NOT NULL,
    Discriminator int NOT NULL,
    PRIMARY KEY (ID)
);

CREATE TABLE IF NOT EXISTS Shots (
    Name varchar(100) NOT NULL,
    PhotoUrl varchar(1000) NOT NULL,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
