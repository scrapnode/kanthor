BEGIN;

CREATE TABLE IF NOT EXISTS kanthor_message (
    app_id VARCHAR(64) NOT NULL,
    id VARCHAR(64) NOT NULL,
    PRIMARY KEY (app_id, id),

    timestamp BIGINT NOT NULL DEFAULT 0,
    tier VARCHAR(64) NOT NULL,
    type VARCHAR(256) NOT NULL,
    headers TEXT NOT NULL,
    body TEXT NOT NULL,
    metadata TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kanthor_request (
    app_id VARCHAR(64) NOT NULL,
    msg_id VARCHAR(64) NOT NULL,
    id VARCHAR(64) NOT NULL,
    PRIMARY KEY (app_id, msg_id, id),

    timestamp BIGINT NOT NULL DEFAULT 0,
    ep_id VARCHAR(64) NOT NULL,
    tier VARCHAR(64) NOT NULL,
    type VARCHAR(256) NOT NULL,
    metadata TEXT NOT NULL,
    method VARCHAR(64) NOT NULL,
    uri VARCHAR(2048) NOT NULL,
    headers TEXT NOT NULL,
    body TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kanthor_response (
    app_id VARCHAR(64) NOT NULL,
    msg_id VARCHAR(64) NOT NULL,
    id VARCHAR(64) NOT NULL,
    PRIMARY KEY (app_id, msg_id, id),
    
    timestamp BIGINT NOT NULL DEFAULT 0,
    ep_id VARCHAR(64) NOT NULL,
    req_id VARCHAR(64) NOT NULL,
    tier VARCHAR(64) NOT NULL,
    type VARCHAR(256) NOT NULL,
    metadata TEXT NOT NULL,
    status INT NOT NULL,
    uri VARCHAR(2048) NOT NULL,
    headers TEXT NOT NULL,
    body TEXT NOT NULL,
    error TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS kanthor_attempt (
    req_id VARCHAR(64) NOT NULL,
    id VARCHAR(64) NOT NULL,
    PRIMARY KEY (req_id, id),

    res_id VARCHAR(64) NOT NULL,
    tier VARCHAR(64) NOT NULL,
    status INT NOT NULL,
    scheduled_at BIGINT NOT NULL DEFAULT 0,
    completed_at BIGINT NOT NULL DEFAULT 0
)
COMMIT;