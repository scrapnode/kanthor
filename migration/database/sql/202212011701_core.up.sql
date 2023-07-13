-- migrating
CREATE TABLE IF NOT EXISTS workspace
(
    id         VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at BIGINT       NOT NULL DEFAULT 0,
    updated_at BIGINT       NOT NULL DEFAULT 0,
    deleted_at BIGINT       NOT NULL DEFAULT 0,

    owner_id   VARCHAR(64)  NOT NULL,
    name       VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS workspace_tier
(
    workspace_id VARCHAR(64)  NOT NULL,
    name         VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS workspace_privilege
(
    id           VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at   BIGINT       NOT NULL DEFAULT 0,
    updated_at   BIGINT       NOT NULL DEFAULT 0,
    deleted_at   BIGINT       NOT NULL DEFAULT 0,

    workspace_id VARCHAR(64)  NOT NULL,
    account_sub  VARCHAR(64)  NOT NULL,
    account_role VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS application
(
    id           VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at   BIGINT       NOT NULL DEFAULT 0,
    updated_at   BIGINT       NOT NULL DEFAULT 0,
    deleted_at   BIGINT       NOT NULL DEFAULT 0,

    workspace_id VARCHAR(64)  NOT NULL,
    name         VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS endpoint
(
    id         VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at BIGINT       NOT NULL DEFAULT 0,
    updated_at BIGINT       NOT NULL DEFAULT 0,
    deleted_at BIGINT       NOT NULL DEFAULT 0,

    app_id     VARCHAR(64)  NOT NULL,
    name       VARCHAR(256) NOT NULL,
    uri        TEXT         NOT NULL,
    method     VARCHAR(64)  NOT NULL
);

CREATE TABLE IF NOT EXISTS endpoint_rule
(
    id                   VARCHAR(64)  NOT NULL PRIMARY KEY,
    created_at           BIGINT       NOT NULL DEFAULT 0,
    updated_at           BIGINT       NOT NULL DEFAULT 0,
    deleted_at           BIGINT       NOT NULL DEFAULT 0,

    endpoint_id          VARCHAR(64)  NOT NULL,
    name                 VARCHAR(256) NOT NULL,
    condition_source     VARCHAR(256) NOT NULL,
    condition_expression TEXT         NOT NULL,
    priority             SMALLINT     NOT NULL DEFAULT 0,
    exclusionary         BOOLEAN      NOT NULL DEFAULT FALSE
);

-- seeding
INSERT INTO workspace
    (id, created_at, updated_at, deleted_at, owner_id, name)
VALUES ('ws_2IJycYlPQSXCTwVZ31vNC9taz20', 1669914060000, 0, 0, 'u_2IJycYlPREw5pnIp8IcMrxatSLo', 'default')
ON CONFLICT DO NOTHING;

INSERT INTO workspace_tier
    (workspace_id, name)
VALUES ('ws_2IJycYlPQSXCTwVZ31vNC9taz20', 'default')
ON CONFLICT DO NOTHING;

INSERT INTO application (id, created_at, updated_at, deleted_at, workspace_id, name)
VALUES ('app_2IJycYlPREw5nqMDss3TCdQhotU', 1669914060000, 0, 0, 'ws_2IJycYlPQSXCTwVZ31vNC9taz20', 'demo')
ON CONFLICT DO NOTHING;

INSERT INTO endpoint (id, created_at, updated_at, deleted_at, app_id, name, uri, method)
VALUES ('ep_2IJyV1PA8HtGES52vOhYFGJ4M40', 1669914060000, 0, 0, 'app_2IJycYlPREw5nqMDss3TCdQhotU', 'httpin-200',
        'https://httpbin.org/post', 'POST'),
       ('ep_2IJyV1PA8HtGES53zKkOzJbZCim', 1669914060000, 0, 0, 'app_2IJycYlPREw5nqMDss3TCdQhotU', 'httpin-500',
        'https://httpbin.org/status/500', 'PUT')
ON CONFLICT DO NOTHING;

INSERT INTO endpoint_rule (id, created_at, updated_at, deleted_at, endpoint_id, name, condition_source,
                           condition_expression,
                           priority, exclusionary)
VALUES ('epr_2IJycYfnerRGxL3uzcvA20mHVqK', 1669914060000, 0, 0, 'ep_2IJyV1PA8HtGES52vOhYFGJ4M40', 'httpin-200',
        'app_id',
        'equal::app_2IJycYlPREw5nqMDss3TCdQhotU', 0, false),
       ('epr_2IJyV1PA8HtGES53zKkT9BktUo4', 1669914060000, 0, 0, 'ep_2IJyV1PA8HtGES53zKkOzJbZCim', 'httpin-500',
        'app_id',
        'equal::app_2IJycYlPREw5nqMDss3TCdQhotU', 0, false)
ON CONFLICT DO NOTHING;