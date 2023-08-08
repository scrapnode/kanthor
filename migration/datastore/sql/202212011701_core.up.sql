CREATE TABLE IF NOT EXISTS kanthor_message
(
    id        VARCHAR(64)  NOT NULL PRIMARY KEY,
    timestamp BIGINT       NOT NULL DEFAULT 0,
    bucket    VARCHAR(64)  NOT NULL,

    tier      VARCHAR(64)  NOT NULL,
    app_id    VARCHAR(64)  NOT NULL,
    type      VARCHAR(256) NOT NULL,
    headers   TEXT         NOT NULL,
    body      TEXT         NOT NULL,
    metadata  TEXT         NOT NULL
);
CREATE INDEX kanthor_msg_scan ON kanthor_message (bucket DESC, tier DESC, id DESC);

CREATE TABLE IF NOT EXISTS kanthor_request
(
    id        VARCHAR(64)   NOT NULL PRIMARY KEY,
    timestamp BIGINT        NOT NULL DEFAULT 0,
    bucket    VARCHAR(64)   NOT NULL,

    tier      VARCHAR(64)   NOT NULL,
    app_id    VARCHAR(64)   NOT NULL,
    type      VARCHAR(256)  NOT NULL,
    metadata  TEXT          NOT NULL,

    method    VARCHAR(64)   NOT NULL,
    uri       VARCHAR(2048) NOT NULL,
    headers   TEXT          NOT NULL,
    body      TEXT          NOT NULL
);
CREATE INDEX kanthor_req_scan ON kanthor_request (bucket DESC, tier DESC, id DESC);

CREATE TABLE IF NOT EXISTS kanthor_response
(
    id        VARCHAR(64)   NOT NULL PRIMARY KEY,
    timestamp BIGINT        NOT NULL DEFAULT 0,
    bucket    VARCHAR(64)   NOT NULL,

    tier      VARCHAR(64)   NOT NULL,
    app_id    VARCHAR(64)   NOT NULL,
    type      VARCHAR(256)  NOT NULL,
    metadata  TEXT          NOT NULL,

    status    INT   NOT NULL,
    uri       VARCHAR(2048) NOT NULL,
    headers   TEXT          NOT NULL,
    body      TEXT          NOT NULL,

    error TEXT NOT NULL
);
CREATE INDEX kanthor_res_scan ON kanthor_response (bucket DESC, tier DESC, id DESC);
