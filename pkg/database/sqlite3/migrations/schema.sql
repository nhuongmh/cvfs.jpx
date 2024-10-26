CREATE TABLE IF NOT EXISTS sentences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sentence TEXT NOT NULL,
    created_at datetime NOT NULL DEFAULT (datetime('now','localtime')),
    udpated_at datetime NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS cards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    front TEXT,
    back TEXT,
    properties TEXT,
    created_at datetime NOT NULL DEFAULT (datetime('now','localtime')),
    udpated_at datetime NOT NULL DEFAULT (datetime('now','localtime')),
    fsrs_id INTEGER,
    status VARCHAR(255),
    FOREIGN KEY(fsrs_id) REFERENCES fsrs(id)
);

CREATE TABLE IF NOT EXISTS fsrs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    card_id INTEGER,
    due datetime,
    stability REAL,
    difficulty REAL,
    elapsed_days INTEGER,
    scheduled_days INTEGER,
    reps INTEGER,
    lapses INTEGER,
    state VARCHAR(255),
    last_review datetime,
    FOREIGN KEY(card_id) REFERENCES cards(id)
);