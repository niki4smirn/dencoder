CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    link VARCHAR(300) UNIQUE NOT NULL
);
