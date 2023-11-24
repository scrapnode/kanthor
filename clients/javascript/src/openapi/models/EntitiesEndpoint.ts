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

import { HttpFile } from '../http/http';

export class EntitiesEndpoint {
    'appId'?: string;
    /**
    * I didn\'t find a way to disable automatic fields modify yet so, I use a tag to disable this feature here but, we should keep our entities stateless if we can
    */
    'createdAt'?: number;
    'id'?: string;
    /**
    * HTTP: POST/PUT/PATCH
    */
    'method'?: string;
    'name'?: string;
    'secretKey'?: string;
    'updatedAt'?: number;
    /**
    * format: scheme \":\" [\"//\" authority] path [\"?\" query] [\"#\" fragment] HTTP: https:://httpbin.org/post?app=kanthor.webhook gRPC: grpc:://app.kanthorlabs.com
    */
    'uri'?: string;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "appId",
            "baseName": "app_id",
            "type": "string",
            "format": ""
        },
        {
            "name": "createdAt",
            "baseName": "created_at",
            "type": "number",
            "format": ""
        },
        {
            "name": "id",
            "baseName": "id",
            "type": "string",
            "format": ""
        },
        {
            "name": "method",
            "baseName": "method",
            "type": "string",
            "format": ""
        },
        {
            "name": "name",
            "baseName": "name",
            "type": "string",
            "format": ""
        },
        {
            "name": "secretKey",
            "baseName": "secret_key",
            "type": "string",
            "format": ""
        },
        {
            "name": "updatedAt",
            "baseName": "updated_at",
            "type": "number",
            "format": ""
        },
        {
            "name": "uri",
            "baseName": "uri",
            "type": "string",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return EntitiesEndpoint.attributeTypeMap;
    }

    public constructor() {
    }
}

