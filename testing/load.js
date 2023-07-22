import grpc from 'k6/net/grpc';
import {check, sleep} from 'k6';

const client = new grpc.Client();
client.load(['/Users/tuannguyen/Projects/scrapnode/kanthor/services/dataplane/grpc/protos'], 'dataplane.proto');

export const options = {
    // Key configurations for spike in this section
    stages: [
        { duration: '2m', target: 200 }, // fast ramp-up to a high point
        { duration: '30m', target: 150 }, // fast ramp-up to a high point
        { duration: '1m', target: 0 }, // quick ramp-down to 0 users
    ],
};

export default () => {
    client.connect('localhost:8181', {
        plaintext: true,
        timeout: '5s'
    });

    const data = {
        "appId": "app_2IJycYlPREw5nqMDss3TCdQhotU",
        "type": "test.demo",
        "headers": {
            "x-kanthor-client": "k6"
        },
        "body": "{\"author\":\"Tuan Nguyen\"}",
        "metadata": {
            "x-kanthor-client": "k6"
        }
    };
    const response = client.invoke('kanthor.dataplane.v1.Message/Put', data);

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    client.close();
    sleep(1);
};
