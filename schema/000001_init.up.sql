CREATE TABLE users (
   id serial PRIMARY KEY,
   name varchar(255) NOT NULL,
   username varchar(255) NOT NULL UNIQUE,
   password_hash varchar(255) NOT NULL
);

CREATE TABLE tag_list (
   id serial PRIMARY KEY,
   user_id integer NOT NULL UNIQUE,
   FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE tag (
  id serial PRIMARY KEY,
  tag_list_id integer NOT NULL,
  title varchar(255) NOT NULL,
  description varchar(255),
  FOREIGN KEY (tag_list_id) REFERENCES tag_list (id) ON DELETE CASCADE
);
