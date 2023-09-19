import { check } from "k6";
import http from "k6/http";
import { b64encode } from "k6/encoding";
import {
  uuidv4,
  randomString,
} from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

const sdkapi = JSON.parse(open(`${__ENV.API_CREDS_PATH}/sdkapi.json`));

// init context: define k6 options
export const options = {
  vus: Number(__ENV.K6_VUS || 1),
  duration: __ENV.K6_DURATION || "30s",
};

export default () => {
  const authorization = b64encode(`${sdkapi.credentials.username}:${sdkapi.credentials.password}`);

  const url = `${__ENV.API_ENDPOINT}/api/application/${sdkapi.applications[0]}/message`;
  const payload = JSON.stringify({
    type: "testing.traffic.stability",
    body: { username: `u_${randomString(28)}@kanthorlabs.com` },
    headers: { "x-client": "k6.io" },
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
      "Idempotency-Key": uuidv4(),
      Authorization: `Basic ${authorization}`,
    },
  };

  const res = http.put(url, payload, params);
  check(res, {
    "is status 201": (r) => r.status === 201,
  });
};