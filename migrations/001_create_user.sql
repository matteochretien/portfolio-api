-- Write your migrate up statements here
CREATE TABLE users
(
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name  VARCHAR(255) NOT NULL,
    last_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    password    VARCHAR(255) NOT NULL,
    permissions varchar(255)[]   DEFAULT '{}',
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----
DROP TABLE IF EXISTS users;