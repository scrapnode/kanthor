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
import { RestEndpointCreateReq } from '../models/RestEndpointCreateReq';
import { RestEndpointCreateRes } from '../models/RestEndpointCreateRes';
import { RestEndpointDeleteRes } from '../models/RestEndpointDeleteRes';
import { RestEndpointGetRes } from '../models/RestEndpointGetRes';
import { RestEndpointListRes } from '../models/RestEndpointListRes';
import { RestEndpointUpdateReq } from '../models/RestEndpointUpdateReq';
import { RestEndpointUpdateRes } from '../models/RestEndpointUpdateRes';

/**
 * no description
 */
export class EndpointApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public async applicationAppIdEndpointEpIdDelete(appId: string, epId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'appId' is not null or undefined
        if (appId === null || appId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointEpIdDelete", "appId");
        }


        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointEpIdDelete", "epId");
        }


        // Path Params
        const localVarPath = '/application/{app_id}/endpoint/{ep_id}'
            .replace('{' + 'app_id' + '}', encodeURIComponent(String(appId)))
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)));

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
     * @param appId application id
     * @param epId endpoint id
     */
    public async applicationAppIdEndpointEpIdGet(appId: string, epId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'appId' is not null or undefined
        if (appId === null || appId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointEpIdGet", "appId");
        }


        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointEpIdGet", "epId");
        }


        // Path Params
        const localVarPath = '/application/{app_id}/endpoint/{ep_id}'
            .replace('{' + 'app_id' + '}', encodeURIComponent(String(appId)))
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)));

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
     * @param appId application id
     * @param epId endpoint id
     * @param props endpoint properties
     */
    public async applicationAppIdEndpointEpIdPut(appId: string, epId: string, props: RestEndpointUpdateReq, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'appId' is not null or undefined
        if (appId === null || appId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointEpIdPut", "appId");
        }


        // verify required parameter 'epId' is not null or undefined
        if (epId === null || epId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointEpIdPut", "epId");
        }


        // verify required parameter 'props' is not null or undefined
        if (props === null || props === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointEpIdPut", "props");
        }


        // Path Params
        const localVarPath = '/application/{app_id}/endpoint/{ep_id}'
            .replace('{' + 'app_id' + '}', encodeURIComponent(String(appId)))
            .replace('{' + 'ep_id' + '}', encodeURIComponent(String(epId)));

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
            ObjectSerializer.serialize(props, "RestEndpointUpdateReq", ""),
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
     * @param appId application id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public async applicationAppIdEndpointGet(appId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'appId' is not null or undefined
        if (appId === null || appId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointGet", "appId");
        }






        // Path Params
        const localVarPath = '/application/{app_id}/endpoint'
            .replace('{' + 'app_id' + '}', encodeURIComponent(String(appId)));

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
     * @param appId application id
     * @param props endpoint properties
     */
    public async applicationAppIdEndpointPost(appId: string, props: RestEndpointCreateReq, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'appId' is not null or undefined
        if (appId === null || appId === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointPost", "appId");
        }


        // verify required parameter 'props' is not null or undefined
        if (props === null || props === undefined) {
            throw new RequiredError("EndpointApi", "applicationAppIdEndpointPost", "props");
        }


        // Path Params
        const localVarPath = '/application/{app_id}/endpoint'
            .replace('{' + 'app_id' + '}', encodeURIComponent(String(appId)));

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
            ObjectSerializer.serialize(props, "RestEndpointCreateReq", ""),
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

export class EndpointApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to applicationAppIdEndpointEpIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async applicationAppIdEndpointEpIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointDeleteRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointDeleteRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointDeleteRes", ""
            ) as RestEndpointDeleteRes;
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
            const body: RestEndpointDeleteRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointDeleteRes", ""
            ) as RestEndpointDeleteRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to applicationAppIdEndpointEpIdGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async applicationAppIdEndpointEpIdGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointGetRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointGetRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointGetRes", ""
            ) as RestEndpointGetRes;
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
            const body: RestEndpointGetRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointGetRes", ""
            ) as RestEndpointGetRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to applicationAppIdEndpointEpIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async applicationAppIdEndpointEpIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointUpdateRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointUpdateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointUpdateRes", ""
            ) as RestEndpointUpdateRes;
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
            const body: RestEndpointUpdateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointUpdateRes", ""
            ) as RestEndpointUpdateRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to applicationAppIdEndpointGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async applicationAppIdEndpointGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointListRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: RestEndpointListRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointListRes", ""
            ) as RestEndpointListRes;
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
            const body: RestEndpointListRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointListRes", ""
            ) as RestEndpointListRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to applicationAppIdEndpointPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async applicationAppIdEndpointPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<RestEndpointCreateRes >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: RestEndpointCreateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointCreateRes", ""
            ) as RestEndpointCreateRes;
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
            const body: RestEndpointCreateRes = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "RestEndpointCreateRes", ""
            ) as RestEndpointCreateRes;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Buffer | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
