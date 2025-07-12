CREATE TABLE articles (
    id         UUID PRIMARY KEY,
    author_id  UUID REFERENCES authors(id),
    title      TEXT NOT NULL,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    tsv        tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(body, '')), 'B')
    ) STORED
);

CREATE INDEX idx_articles_created_at ON articles(created_at DESC);
CREATE INDEX idx_articles_tsv ON articles USING GIN(tsv);
