CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  email varchar NOT NULL,
  nickname varchar NOT NULL,
  image varchar NOT NULL,
  country_code varchar NOT NULL,
  birthday date NOT NULL,
  created_at timestamptz DEFAULT (now())
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("nickname");