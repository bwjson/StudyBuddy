CREATE TABLE users (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    username varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    is_active boolean DEFAULT FALSE,
    verification_token varchar(255)
);

CREATE TABLE tags (
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL UNIQUE,
    description varchar(255)
);

CREATE TABLE user_tags (
    id serial PRIMARY KEY,
    user_id integer NOT NULL,
    tag_id integer NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
    UNIQUE (user_id, tag_id)
);