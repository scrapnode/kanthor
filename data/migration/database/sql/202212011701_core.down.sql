BEGIN;

DROP INDEX IF EXISTS kanthor_epr_ep_ref;

DROP TABLE IF EXISTS kanthor_endpoint_rule;

DROP INDEX IF EXISTS kanthor_ep_app_ref;

DROP TABLE IF EXISTS kanthor_endpoint;

DROP INDEX IF EXISTS kanthor_app_ws_ref;

DROP TABLE IF EXISTS kanthor_application;

DROP INDEX IF EXISTS kanthor_wsc_ws_ref;

DROP TABLE IF EXISTS kanthor_workspace_credentials;

DROP INDEX IF EXISTS kanthor_ws_owner;

DROP TABLE IF EXISTS kanthor_workspace;

COMMIT;