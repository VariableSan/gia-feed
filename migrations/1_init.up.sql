CREATE TABLE IF NOT EXISTS feeds (
    id        TEXT PRIMARY KEY,
    author_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    status INT NOT NULL DEFAULT 0,
    moderator_comment TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_feed_updated_at
AFTER UPDATE ON feeds
FOR EACH ROW
BEGIN
    UPDATE feeds SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
