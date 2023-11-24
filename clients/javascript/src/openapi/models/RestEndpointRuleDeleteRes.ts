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

export class RestEndpointRuleDeleteRes {
    /**
    * examples:  - equal::orders.paid  - regex::.*
    */
    'conditionExpression'?: string;
    /**
    * examples  - app_id  - type  - body  - metadata
    */
    'conditionSource'?: string;
    /**
    * I didn\'t find a way to disable automatic fields modify yet so, I use a tag to disable this feature here but, we should keep our entities stateless if we can
    */
    'createdAt'?: number;
    'epId'?: string;
    /**
    * the logic of not-false is true should be used here to guarantee default all rule will be on include mode
    */
    'exclusionary'?: boolean;
    'id'?: string;
    'name'?: string;
    'priority'?: number;
    'updatedAt'?: number;

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
            "name": "createdAt",
            "baseName": "created_at",
            "type": "number",
            "format": ""
        },
        {
            "name": "epId",
            "baseName": "ep_id",
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
            "name": "id",
            "baseName": "id",
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
            "name": "priority",
            "baseName": "priority",
            "type": "number",
            "format": ""
        },
        {
            "name": "updatedAt",
            "baseName": "updated_at",
            "type": "number",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return RestEndpointRuleDeleteRes.attributeTypeMap;
    }

    public constructor() {
    }
}

