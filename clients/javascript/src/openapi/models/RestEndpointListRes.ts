/**
 * Kanthor SDK API
 * SDK API
 *
 * OpenAPI spec version: 1.0
 * Contact: support@kanthorlabs.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { EntitiesEndpoint } from '../models/EntitiesEndpoint';
import { HttpFile } from '../http/http';

export class RestEndpointListRes {
    'cursor'?: string;
    'data'?: Array<EntitiesEndpoint>;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "cursor",
            "baseName": "cursor",
            "type": "string",
            "format": ""
        },
        {
            "name": "data",
            "baseName": "data",
            "type": "Array<EntitiesEndpoint>",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return RestEndpointListRes.attributeTypeMap;
    }

    public constructor() {
    }
}

