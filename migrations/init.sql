CREATE TABLE IF NOT EXISTS urls (
  id SERIAL PRIMARY KEY,
  code VARCHAR(32) UNIQUE NOT NULL,
  original_url TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  expires_at TIMESTAMPTZ,
  clicks_count BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS clicks (
  id BIGSERIAL PRIMARY KEY,
  url_id INTEGER REFERENCES urls(id) ON DELETE CASCADE,
  code VARCHAR(32) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  ip VARCHAR(45),
  country VARCHAR(100),
  city VARCHAR(100),
  user_agent TEXT,
  referer TEXT
);

CREATE INDEX IF NOT EXISTS idx_clicks_url_id ON clicks(url_id);
CREATE INDEX IF NOT EXISTS idx_urls_code ON urls(code);
