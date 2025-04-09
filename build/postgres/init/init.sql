CREATE TABLE dialogues (
    id UUID PRIMARY KEY,        -- **原生UUID类型**
    username VARCHAR(255) NOT NULL,
    create_time TIMESTAMPTZ DEFAULT NOW(),
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE messages (
    id UUID PRIMARY KEY,
    dialogue_id UUID NOT NULL,
    role VARCHAR(20) CHECK (role IN ('system','user','assistant')),
    content TEXT NOT NULL,
    seq_num INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
);

CREATE INDEX idx_dialogue_seq ON messages (dialogue_id, seq_num);