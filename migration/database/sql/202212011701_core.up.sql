CREATE TABLE IF NOT EXISTS kanthor_workspace
(
    id          VARCHAR(64)  NOT NULL PRIMARY KEY,
    modified_by VARCHAR(64)  NOT NULL,
    created_at  BIGINT       NOT NULL DEFAULT 0,
    updated_at  BIGINT       NOT NULL DEFAULT 0,

    owner_id    VARCHAR(64)  NOT NULL,
    name        VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS kanthor_workspace_tier
(
    id           VARCHAR(64)  NOT NULL PRIMARY KEY,
    modified_by  TEXT         NOT NULL,
    created_at   BIGINT       NOT NULL DEFAULT 0,
    updated_at   BIGINT       NOT NULL DEFAULT 0,

    workspace_id VARCHAR(64)  NOT NULL,
    name         VARCHAR(256) NOT NULL,

    CONSTRAINT workspace_unique UNIQUE (workspace_id),
    FOREIGN KEY (workspace_id) REFERENCES kanthor_workspace (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS kanthor_workspace_credentials
(
    id           VARCHAR(64)  NOT NULL PRIMARY KEY,
    modified_by  TEXT         NOT NULL,
    created_at   BIGINT       NOT NULL DEFAULT 0,
    updated_at   BIGINT       NOT NULL DEFAULT 0,

    workspace_id VARCHAR(64)  NOT NULL,
    hash         VARCHAR(256) NOT NULL,
    expired_at   BIGINT       NOT NULL DEFAULT 0,

    FOREIGN KEY (workspace_id) REFERENCES kanthor_workspace (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS kanthor_application
(
    id           VARCHAR(64)  NOT NULL PRIMARY KEY,
    modified_by  TEXT         NOT NULL,
    created_at   BIGINT       NOT NULL DEFAULT 0,
    updated_at   BIGINT       NOT NULL DEFAULT 0,

    workspace_id VARCHAR(64)  NOT NULL,
    name         VARCHAR(256) NOT NULL,

    FOREIGN KEY (workspace_id) REFERENCES kanthor_workspace (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS kanthor_endpoint
(
    id          VARCHAR(64)  NOT NULL PRIMARY KEY,
    modified_by VARCHAR(64)  NOT NULL,
    created_at  BIGINT       NOT NULL DEFAULT 0,
    updated_at  BIGINT       NOT NULL DEFAULT 0,

    app_id      VARCHAR(64)  NOT NULL,
    secret_key  VARCHAR(64)  NOT NULL,
    name        VARCHAR(256) NOT NULL,
    uri         TEXT         NOT NULL,
    method      VARCHAR(64)  NOT NULL,

    FOREIGN KEY (app_id) REFERENCES kanthor_application (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS kanthor_endpoint_rule
(
    id                   VARCHAR(64)  NOT NULL PRIMARY KEY,
    modified_by          VARCHAR(64)  NOT NULL,
    created_at           BIGINT       NOT NULL DEFAULT 0,
    updated_at           BIGINT       NOT NULL DEFAULT 0,

    endpoint_id          VARCHAR(64)  NOT NULL,
    name                 VARCHAR(256) NOT NULL,
    condition_source     VARCHAR(256) NOT NULL,
    condition_expression TEXT         NOT NULL,
    priority             SMALLINT     NOT NULL DEFAULT 0,
    exclusionary         BOOLEAN      NOT NULL DEFAULT FALSE,

    FOREIGN KEY (endpoint_id) REFERENCES kanthor_endpoint (id) ON DELETE CASCADE
);
