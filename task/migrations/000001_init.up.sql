CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(320) UNIQUE, 
    hash_pass VARCHAR(70)
);


CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    text VARCHAR(512),
    user_id INT, 
    FOREIGN KEY (user_id) REFERENCES users (id)
);
