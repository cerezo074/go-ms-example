CREATE TABLE users (
  id uuid DEFAULT (uuid_generate_v4()),
  email varchar PRIMARY KEY,
  nickname varchar NOT NULL,
  password varchar NOT NULL,
  image_url varchar NOT NULL,
  country_code varchar NOT NULL,
  birthday date NOT NULL,
  created_at timestamptz DEFAULT (now()),
  updated_at timestamptz DEFAULT (now())
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("nickname");

CREATE OR REPLACE FUNCTION set_updated_at_on_users()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_on_users();