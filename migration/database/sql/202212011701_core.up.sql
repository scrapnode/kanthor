BEGIN;

CREATE TABLE IF NOT EXISTS kanthor_workspace
(
    id          VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at  BIGINT       NOT NULL DEFAULT 0,
    updated_at  BIGINT       NOT NULL DEFAULT 0,

    owner_id    VARCHAR(64)  NOT NULL,
    name        VARCHAR(256) NOT NULL,
    tier        VARCHAR(256) NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS  kanthor_ws_owner_ref ON kanthor_workspace(owner_id DESC);

CREATE TABLE IF NOT EXISTS kanthor_workspace_credentials
(
    id           VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at   BIGINT       NOT NULL DEFAULT 0,
    updated_at   BIGINT       NOT NULL DEFAULT 0,

    ws_id VARCHAR(64)  NOT NULL,
    name         VARCHAR(256) NOT NULL,
    hash         VARCHAR(256) NOT NULL,
    expired_at   BIGINT       NOT NULL DEFAULT 0,

    FOREIGN KEY (ws_id) REFERENCES kanthor_workspace (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS kanthor_wsc_ws_ref ON kanthor_workspace_credentials(ws_id DESC);

CREATE TABLE IF NOT EXISTS kanthor_application
(
    id           VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at   BIGINT       NOT NULL DEFAULT 0,
    updated_at   BIGINT       NOT NULL DEFAULT 0,

    ws_id VARCHAR(64)  NOT NULL,
    name         VARCHAR(256) NOT NULL,

    FOREIGN KEY (ws_id) REFERENCES kanthor_workspace (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS kanthor_app_ws_ref ON kanthor_application(ws_id DESC);

CREATE TABLE IF NOT EXISTS kanthor_endpoint
(
    id          VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at  BIGINT       NOT NULL DEFAULT 0,
    updated_at  BIGINT       NOT NULL DEFAULT 0,

    app_id      VARCHAR(64)  NOT NULL,
    secret_key  VARCHAR(64)  NOT NULL,
    name        VARCHAR(256) NOT NULL,
    uri         TEXT         NOT NULL,
    method      VARCHAR(64)  NOT NULL,

    FOREIGN KEY (app_id) REFERENCES kanthor_application (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS kanthor_ep_app_ref ON kanthor_endpoint(app_id DESC);

CREATE TABLE IF NOT EXISTS kanthor_endpoint_rule
(
    id                   VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at           BIGINT       NOT NULL DEFAULT 0,
    updated_at           BIGINT       NOT NULL DEFAULT 0,

    ep_id          VARCHAR(64)  NOT NULL,
    name                 VARCHAR(256) NOT NULL,
    condition_source     VARCHAR(256) NOT NULL,
    condition_expression TEXT         NOT NULL,
    priority             SMALLINT     NOT NULL DEFAULT 0,
    exclusionary         BOOLEAN      NOT NULL DEFAULT FALSE,

    FOREIGN KEY (ep_id) REFERENCES kanthor_endpoint (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS kanthor_epr_ep_ref ON kanthor_endpoint_rule(ep_id DESC);

COMMIT;