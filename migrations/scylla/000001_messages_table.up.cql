CREATE KEYSPACE messages_space WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor' : 1};

use messages_space;

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    user_id UUID,
    room_id UUID,
    content TEXT,
    time_created TIMESTAMP,
);

CREATE INDEX IF NOT EXISTS ON messages (user_id);
