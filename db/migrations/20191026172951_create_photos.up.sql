CREATE TABLE photos(
  id serial PRIMARY KEY,
  location varchar NOT NULL,
  original boolean DEFAULT false NOT NULL,
  attrs jsonb DEFAULT '{}'::jsonb NOT NULL
);

CREATE INDEX photos_attr_idx ON photos USING gin (attrs);
