CREATE TABLE IF NOT EXISTS ie_articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    content TEXT NOT NULL,
    origin VARCHAR NOT NULL,
    author VARCHAR NOT NULL,
    cover_image VARCHAR NOT NULL,
    publish_date VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ie_articles_title_idx ON ie_articles(title);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON ie_articles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();