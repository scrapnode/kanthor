DROP INDEX IF EXISTS kanthor_wst_ws_unique;
DROP INDEX IF EXISTS kanthor_wsc_ws_ref;
DROP INDEX IF EXISTS kanthor_app_ws_ref;
DROP INDEX IF EXISTS kanthor_ep_app_ref;
DROP INDEX IF EXISTS kanthor_epr_ep_ref;

DROP TABLE IF EXISTS kanthor_endpoint_rule;
DROP TABLE IF EXISTS kanthor_endpoint;
DROP TABLE IF EXISTS kanthor_application;
DROP TABLE IF EXISTS kanthor_workspace_tier;
DROP TABLE IF EXISTS kanthor_workspace;