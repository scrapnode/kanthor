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

export class EntitiesApplication {
    /**
    * I didn\'t find a way to disable automatic fields modify yet so, I use a tag to disable this feature here but, we should keep our entities stateless if we can
    */
    'createdAt'?: number;
    'id'?: string;
    'name'?: string;
    'updatedAt'?: number;
    'wsId'?: string;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
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
            "name": "name",
            "baseName": "name",
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
            "name": "wsId",
            "baseName": "ws_id",
            "type": "string",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return EntitiesApplication.attributeTypeMap;
    }

    public constructor() {
    }
}

