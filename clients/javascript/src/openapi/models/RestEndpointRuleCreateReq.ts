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

export class RestEndpointRuleCreateReq {
    'conditionExpression': string;
    'conditionSource': string;
    'exclusionary'?: boolean;
    'name': string;
    'priority'?: number;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "conditionExpression",
            "baseName": "condition_expression",
            "type": "string",
            "format": ""
        },
        {
            "name": "conditionSource",
            "baseName": "condition_source",
            "type": "string",
            "format": ""
        },
        {
            "name": "exclusionary",
            "baseName": "exclusionary",
            "type": "boolean",
            "format": ""
        },
        {
            "name": "name",
            "baseName": "name",
            "type": "string",
            "format": ""
        },
        {
            "name": "priority",
            "baseName": "priority",
            "type": "number",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return RestEndpointRuleCreateReq.attributeTypeMap;
    }

    public constructor() {
    }
}

