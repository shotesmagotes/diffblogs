#!/bin/bash
#Create meta sqlite db for blog metadata

echo "creating .meta.db and articles tables..."

sqlite3 meta.db << INITDB

CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY,
    title TEXT,
    type TEXT,
    filename TEXT,
    published TEXT,
    modified TEXT
);

INITDB

echo "done."
