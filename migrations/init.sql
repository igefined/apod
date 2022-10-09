CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    first_name VARCHAR(255)  NOT NULL,
    last_name  VARCHAR(255)  NOT NULL,
    email      VARCHAR(255)  NOT NULL,
    password   VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP     NOT NULL
);

INSERT INTO users(id, first_name, last_name, email, password, created_at)
VALUES (uuid_in(md5(random()::text || random()::text)::cstring), 'Igor', 'Pomazkov', 'ig.pomazkov@gmail.com',
        '$2a$10$qBLtELEVDECprsPT9gpz4uJPR1Sq22Jn/YCQnOpdFuOMiCr/1jWNa', NOW());

