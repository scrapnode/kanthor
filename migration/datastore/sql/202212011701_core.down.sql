BEGIN;

DROP INDEX IF EXISTS kanthor_res_scan;
DROP TABLE IF EXISTS kanthor_response;

DROP INDEX IF EXISTS kanthor_req_scan;
DROP TABLE IF EXISTS kanthor_request;

DROP INDEX IF EXISTS kanthor_msg_scan;
DROP TABLE IF EXISTS kanthor_message;

COMMIT;