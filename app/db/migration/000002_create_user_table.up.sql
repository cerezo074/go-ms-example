CREATE TABLE users (
  id uuid DEFAULT (uuid_generate_v4()),
  email varchar PRIMARY KEY,
  nickname varchar NOT NULL,
  password varchar NOT NULL,
  image_url varchar NOT NULL,
  country_code varchar NOT NULL,
  birthday date NOT NULL,
  created_at timestamptz DEFAULT (now())
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("nickname");