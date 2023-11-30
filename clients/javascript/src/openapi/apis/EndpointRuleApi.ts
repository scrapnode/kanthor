import {v4 as uuidv4} from 'uuid'
import {BaseAPIRequestFactory, RequiredError, COLLECTION_FORMATS} from './baseapi';
import {Configuration} from '../configuration';
import {RequestContext, HttpMethod, ResponseContext, HttpFile, HttpInfo} from '../http/http';
import  FormData from "form-data";
import { URLSearchParams } from 'url';
import {ObjectSerializer} from '../models/ObjectSerializer';
import {ApiException} from './exception';
import {canConsumeForm, isCodeInRange} from '../util';
import {SecurityAuthentication} from '../auth/auth';


import { GatewayError } from '../models/GatewayError';
import { RestEndpointRuleCreateReq } from '../models/RestEndpointRuleCreateReq';
import { RestEndpointRuleCreateRes } from '../models/RestEndpointRuleCreateRes';
import { RestEndpointRuleDeleteRes } from '../models/RestEndpointRuleDeleteRes';
import { RestEndpointRuleGetRes } from '../models/RestEndpointRuleGetRes';
import { RestEndpointRuleListRes } from '../models/RestEndpointRuleListRes';
import { RestEndpointRuleUpdateReq } from '../models/RestEndpointRuleUpdateReq';
import { RestEndpointRuleUpdateRes } from '../models/RestEndpointRuleUpdateRes';

/**
 * no description
 */
export class EndpointRuleApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public async endpointEpIdRuleEprIdDelete(epId: string, eprId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleEprIdDelete", "epId");
        }


        // verify required parameter 'eprId' is not null or undefined
        if (eprId === null || eprId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleEprIdDelete", "eprId");
        }


        // Path Params
        const localVarPath = '/endpoint/{ep_id}/rule/{epr_id}'
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)))
            .replace('{' + 'epr_id' + '}', encodeURIComponent(String(eprId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.DELETE);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")
        requestContext.setHeaderParam("Idempotency-Key", uuidv4()) 



        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["BasicAuth"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public async endpointEpIdRuleEprIdGet(epId: string, eprId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleEprIdGet", "epId");
        }


        // verify required parameter 'eprId' is not null or undefined
        if (eprId === null || eprId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleEprIdGet", "eprId");
        }


        // Path Params
        const localVarPath = '/endpoint/{ep_id}/rule/{epr_id}'
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)))
            .replace('{' + 'epr_id' + '}', encodeURIComponent(String(eprId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")
        requestContext.setHeaderParam("Idempotency-Key", uuidv4()) 



        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["BasicAuth"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     * @param props rule properties
     */
    public async endpointEpIdRuleEprIdPut(epId: string, eprId: string, props: RestEndpointRuleUpdateReq, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleEprIdPut", "epId");
        }


        // verify required parameter 'eprId' is not null or undefined
        if (eprId === null || eprId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleEprIdPut", "eprId");
        }


        // verify required parameter 'props' is not null or undefined
        if (props === null || props === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleEprIdPut", "props");
        }


        // Path Params
        const localVarPath = '/endpoint/{ep_id}/rule/{epr_id}'
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)))
            .replace('{' + 'epr_id' + '}', encodeURIComponent(String(eprId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")
        requestContext.setHeaderParam("Idempotency-Key", uuidv4()) 



        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(props, "RestEndpointRuleUpdateReq", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["BasicAuth"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * @param epId endpoint id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public async endpointEpIdRuleGet(epId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRuleGet", "epId");
        }






        // Path Params
        const localVarPath = '/endpoint/{ep_id}/rule'
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")
        requestContext.setHeaderParam("Idempotency-Key", uuidv4()) 


        // Query Params
        if (cursor !== undefined) {
            requestContext.setQueryParam("_cursor", ObjectSerializer.serialize(cursor, "string", ""));
        }

        // Query Params
        if (q !== undefined) {
            requestContext.setQueryParam("_q", ObjectSerializer.serialize(q, "string", ""));
        }

        // Query Params
        if (limit !== undefined) {
            requestContext.setQueryParam("_limit", ObjectSerializer.serialize(limit, "number", ""));
        }

        // Query Params
        if (id !== undefined) {
            requestContext.setQueryParam("_id", ObjectSerializer.serialize(id, "Array<string>", ""));
        }


        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["BasicAuth"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * @param epId endpoint id
     * @param props rule properties
     */
    public async endpointEpIdRulePost(epId: string, props: RestEndpointRuleCreateReq, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRulePost", "epId");
        }


        // verify required parameter 'props' is not null or undefined
        if (props === null || props === undefined) {
            throw new RequiredError("EndpointRuleApi", "endpointEpIdRulePost", "props");
        }


        // Path Params
        const localVarPath = '/endpoint/{ep_id}/rule'
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")
        requestContext.setHeaderParam("Idempotency-Key", uuidv4()) 



        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(props, "RestEndpointRuleCreateReq", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        let authMethod: SecurityAuthentication | undefined;
        // Apply auth methods
        authMethod = _config.authMethods["BasicAuth"]
        if (authMethod?.applySecurityAuthentication) {
            await authMethod?.applySecurityAuthentication(requestContext);
        }
        
        const defaultAuth: SecurityAuthentication | undefined = _options?.authMethods?.default || this.configuration?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

}

export class EndpointRuleApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to endpointEpIdRuleEprIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async endpointEpIdRuleEprIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointRuleDeleteRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointRuleDeleteRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleDeleteRes", ""
            ) as RestEndpointRuleDeleteRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("0", response.httpStatusCode)) {
            const body: GatewayError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GatewayError", ""
            ) as GatewayError;
            throw new ApiException<GatewayError>(response.httpStatusCode, "", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: RestEndpointRuleDeleteRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleDeleteRes", ""
            ) as RestEndpointRuleDeleteRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to endpointEpIdRuleEprIdGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async endpointEpIdRuleEprIdGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointRuleGetRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointRuleGetRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleGetRes", ""
            ) as RestEndpointRuleGetRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("0", response.httpStatusCode)) {
            const body: GatewayError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GatewayError", ""
            ) as GatewayError;
            throw new ApiException<GatewayError>(response.httpStatusCode, "", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: RestEndpointRuleGetRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleGetRes", ""
            ) as RestEndpointRuleGetRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to endpointEpIdRuleEprIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async endpointEpIdRuleEprIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointRuleUpdateRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointRuleUpdateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleUpdateRes", ""
            ) as RestEndpointRuleUpdateRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("0", response.httpStatusCode)) {
            const body: GatewayError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GatewayError", ""
            ) as GatewayError;
            throw new ApiException<GatewayError>(response.httpStatusCode, "", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: RestEndpointRuleUpdateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleUpdateRes", ""
            ) as RestEndpointRuleUpdateRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to endpointEpIdRuleGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async endpointEpIdRuleGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointRuleListRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointRuleListRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleListRes", ""
            ) as RestEndpointRuleListRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("0", response.httpStatusCode)) {
            const body: GatewayError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GatewayError", ""
            ) as GatewayError;
            throw new ApiException<GatewayError>(response.httpStatusCode, "", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: RestEndpointRuleListRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleListRes", ""
            ) as RestEndpointRuleListRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to endpointEpIdRulePost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async endpointEpIdRulePostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointRuleCreateRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: RestEndpointRuleCreateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleCreateRes", ""
            ) as RestEndpointRuleCreateRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("0", response.httpStatusCode)) {
            const body: GatewayError = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GatewayError", ""
            ) as GatewayError;
            throw new ApiException<GatewayError>(response.httpStatusCode, "", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: RestEndpointRuleCreateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointRuleCreateRes", ""
            ) as RestEndpointRuleCreateRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
