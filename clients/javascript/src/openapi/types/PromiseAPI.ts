import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration} from '../configuration'

import { AuthenticatorAccount } from '../models/AuthenticatorAccount';
import { AuthorizatorPermission } from '../models/AuthorizatorPermission';
import { EntitiesApplication } from '../models/EntitiesApplication';
import { EntitiesEndpoint } from '../models/EntitiesEndpoint';
import { EntitiesEndpointRule } from '../models/EntitiesEndpointRule';
import { GatewayError } from '../models/GatewayError';
import { RestAccountGetRes } from '../models/RestAccountGetRes';
import { RestApplicationCreateReq } from '../models/RestApplicationCreateReq';
import { RestApplicationCreateRes } from '../models/RestApplicationCreateRes';
import { RestApplicationDeleteRes } from '../models/RestApplicationDeleteRes';
import { RestApplicationGetRes } from '../models/RestApplicationGetRes';
import { RestApplicationListRes } from '../models/RestApplicationListRes';
import { RestApplicationUpdateReq } from '../models/RestApplicationUpdateReq';
import { RestApplicationUpdateRes } from '../models/RestApplicationUpdateRes';
import { RestEndpointCreateReq } from '../models/RestEndpointCreateReq';
import { RestEndpointCreateRes } from '../models/RestEndpointCreateRes';
import { RestEndpointDeleteRes } from '../models/RestEndpointDeleteRes';
import { RestEndpointGetRes } from '../models/RestEndpointGetRes';
import { RestEndpointListRes } from '../models/RestEndpointListRes';
import { RestEndpointRuleCreateReq } from '../models/RestEndpointRuleCreateReq';
import { RestEndpointRuleCreateRes } from '../models/RestEndpointRuleCreateRes';
import { RestEndpointRuleDeleteRes } from '../models/RestEndpointRuleDeleteRes';
import { RestEndpointRuleGetRes } from '../models/RestEndpointRuleGetRes';
import { RestEndpointRuleListRes } from '../models/RestEndpointRuleListRes';
import { RestEndpointRuleUpdateReq } from '../models/RestEndpointRuleUpdateReq';
import { RestEndpointRuleUpdateRes } from '../models/RestEndpointRuleUpdateRes';
import { RestEndpointUpdateReq } from '../models/RestEndpointUpdateReq';
import { RestEndpointUpdateRes } from '../models/RestEndpointUpdateRes';
import { RestMessagePutReq } from '../models/RestMessagePutReq';
import { RestMessagePutRes } from '../models/RestMessagePutRes';
import { ObservableAccountApi } from './ObservableAPI';

import { AccountApiRequestFactory, AccountApiResponseProcessor} from "../apis/AccountApi";
export class PromiseAccountApi {
    private api: ObservableAccountApi

    public constructor(
        configuration: Configuration,
        requestFactory?: AccountApiRequestFactory,
        responseProcessor?: AccountApiResponseProcessor
    ) {
        this.api = new ObservableAccountApi(configuration, requestFactory, responseProcessor);
    }

    /**
     */
    public accountMeGetWithHttpInfo(_options?: Configuration): Promise<HttpInfo<RestAccountGetRes>> {
        const result = this.api.accountMeGetWithHttpInfo(_options);
        return result.toPromise();
    }

    /**
     */
    public accountMeGet(_options?: Configuration): Promise<RestAccountGetRes> {
        const result = this.api.accountMeGet(_options);
        return result.toPromise();
    }


}



import { ObservableApplicationApi } from './ObservableAPI';

import { ApplicationApiRequestFactory, ApplicationApiResponseProcessor} from "../apis/ApplicationApi";
export class PromiseApplicationApi {
    private api: ObservableApplicationApi

    public constructor(
        configuration: Configuration,
        requestFactory?: ApplicationApiRequestFactory,
        responseProcessor?: ApplicationApiResponseProcessor
    ) {
        this.api = new ObservableApplicationApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param appId application id
     */
    public applicationAppIdDeleteWithHttpInfo(appId: string, _options?: Configuration): Promise<HttpInfo<RestApplicationDeleteRes>> {
        const result = this.api.applicationAppIdDeleteWithHttpInfo(appId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     */
    public applicationAppIdDelete(appId: string, _options?: Configuration): Promise<RestApplicationDeleteRes> {
        const result = this.api.applicationAppIdDelete(appId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     */
    public applicationAppIdGetWithHttpInfo(appId: string, _options?: Configuration): Promise<HttpInfo<RestApplicationGetRes>> {
        const result = this.api.applicationAppIdGetWithHttpInfo(appId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     */
    public applicationAppIdGet(appId: string, _options?: Configuration): Promise<RestApplicationGetRes> {
        const result = this.api.applicationAppIdGet(appId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param props application properties
     */
    public applicationAppIdPutWithHttpInfo(appId: string, props: RestApplicationUpdateReq, _options?: Configuration): Promise<HttpInfo<RestApplicationUpdateRes>> {
        const result = this.api.applicationAppIdPutWithHttpInfo(appId, props, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param props application properties
     */
    public applicationAppIdPut(appId: string, props: RestApplicationUpdateReq, _options?: Configuration): Promise<RestApplicationUpdateRes> {
        const result = this.api.applicationAppIdPut(appId, props, _options);
        return result.toPromise();
    }

    /**
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationGetWithHttpInfo(cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<HttpInfo<RestApplicationListRes>> {
        const result = this.api.applicationGetWithHttpInfo(cursor, q, limit, id, _options);
        return result.toPromise();
    }

    /**
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationGet(cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<RestApplicationListRes> {
        const result = this.api.applicationGet(cursor, q, limit, id, _options);
        return result.toPromise();
    }

    /**
     * @param props application properties
     */
    public applicationPostWithHttpInfo(props: RestApplicationCreateReq, _options?: Configuration): Promise<HttpInfo<RestApplicationCreateRes>> {
        const result = this.api.applicationPostWithHttpInfo(props, _options);
        return result.toPromise();
    }

    /**
     * @param props application properties
     */
    public applicationPost(props: RestApplicationCreateReq, _options?: Configuration): Promise<RestApplicationCreateRes> {
        const result = this.api.applicationPost(props, _options);
        return result.toPromise();
    }


}



import { ObservableEndpointApi } from './ObservableAPI';

import { EndpointApiRequestFactory, EndpointApiResponseProcessor} from "../apis/EndpointApi";
export class PromiseEndpointApi {
    private api: ObservableEndpointApi

    public constructor(
        configuration: Configuration,
        requestFactory?: EndpointApiRequestFactory,
        responseProcessor?: EndpointApiResponseProcessor
    ) {
        this.api = new ObservableEndpointApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdDeleteWithHttpInfo(appId: string, epId: string, _options?: Configuration): Promise<HttpInfo<RestEndpointDeleteRes>> {
        const result = this.api.applicationAppIdEndpointEpIdDeleteWithHttpInfo(appId, epId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdDelete(appId: string, epId: string, _options?: Configuration): Promise<RestEndpointDeleteRes> {
        const result = this.api.applicationAppIdEndpointEpIdDelete(appId, epId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdGetWithHttpInfo(appId: string, epId: string, _options?: Configuration): Promise<HttpInfo<RestEndpointGetRes>> {
        const result = this.api.applicationAppIdEndpointEpIdGetWithHttpInfo(appId, epId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdGet(appId: string, epId: string, _options?: Configuration): Promise<RestEndpointGetRes> {
        const result = this.api.applicationAppIdEndpointEpIdGet(appId, epId, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointEpIdPutWithHttpInfo(appId: string, epId: string, props: RestEndpointUpdateReq, _options?: Configuration): Promise<HttpInfo<RestEndpointUpdateRes>> {
        const result = this.api.applicationAppIdEndpointEpIdPutWithHttpInfo(appId, epId, props, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointEpIdPut(appId: string, epId: string, props: RestEndpointUpdateReq, _options?: Configuration): Promise<RestEndpointUpdateRes> {
        const result = this.api.applicationAppIdEndpointEpIdPut(appId, epId, props, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationAppIdEndpointGetWithHttpInfo(appId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<HttpInfo<RestEndpointListRes>> {
        const result = this.api.applicationAppIdEndpointGetWithHttpInfo(appId, cursor, q, limit, id, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationAppIdEndpointGet(appId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<RestEndpointListRes> {
        const result = this.api.applicationAppIdEndpointGet(appId, cursor, q, limit, id, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointPostWithHttpInfo(appId: string, props: RestEndpointCreateReq, _options?: Configuration): Promise<HttpInfo<RestEndpointCreateRes>> {
        const result = this.api.applicationAppIdEndpointPostWithHttpInfo(appId, props, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointPost(appId: string, props: RestEndpointCreateReq, _options?: Configuration): Promise<RestEndpointCreateRes> {
        const result = this.api.applicationAppIdEndpointPost(appId, props, _options);
        return result.toPromise();
    }


}



import { ObservableEndpointRuleApi } from './ObservableAPI';

import { EndpointRuleApiRequestFactory, EndpointRuleApiResponseProcessor} from "../apis/EndpointRuleApi";
export class PromiseEndpointRuleApi {
    private api: ObservableEndpointRuleApi

    public constructor(
        configuration: Configuration,
        requestFactory?: EndpointRuleApiRequestFactory,
        responseProcessor?: EndpointRuleApiResponseProcessor
    ) {
        this.api = new ObservableEndpointRuleApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdDeleteWithHttpInfo(epId: string, eprId: string, _options?: Configuration): Promise<HttpInfo<RestEndpointRuleDeleteRes>> {
        const result = this.api.endpointEpIdRuleEprIdDeleteWithHttpInfo(epId, eprId, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdDelete(epId: string, eprId: string, _options?: Configuration): Promise<RestEndpointRuleDeleteRes> {
        const result = this.api.endpointEpIdRuleEprIdDelete(epId, eprId, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdGetWithHttpInfo(epId: string, eprId: string, _options?: Configuration): Promise<HttpInfo<RestEndpointRuleGetRes>> {
        const result = this.api.endpointEpIdRuleEprIdGetWithHttpInfo(epId, eprId, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdGet(epId: string, eprId: string, _options?: Configuration): Promise<RestEndpointRuleGetRes> {
        const result = this.api.endpointEpIdRuleEprIdGet(epId, eprId, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     * @param props rule properties
     */
    public endpointEpIdRuleEprIdPutWithHttpInfo(epId: string, eprId: string, props: RestEndpointRuleUpdateReq, _options?: Configuration): Promise<HttpInfo<RestEndpointRuleUpdateRes>> {
        const result = this.api.endpointEpIdRuleEprIdPutWithHttpInfo(epId, eprId, props, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     * @param props rule properties
     */
    public endpointEpIdRuleEprIdPut(epId: string, eprId: string, props: RestEndpointRuleUpdateReq, _options?: Configuration): Promise<RestEndpointRuleUpdateRes> {
        const result = this.api.endpointEpIdRuleEprIdPut(epId, eprId, props, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public endpointEpIdRuleGetWithHttpInfo(epId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<HttpInfo<RestEndpointRuleListRes>> {
        const result = this.api.endpointEpIdRuleGetWithHttpInfo(epId, cursor, q, limit, id, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public endpointEpIdRuleGet(epId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Promise<RestEndpointRuleListRes> {
        const result = this.api.endpointEpIdRuleGet(epId, cursor, q, limit, id, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param props rule properties
     */
    public endpointEpIdRulePostWithHttpInfo(epId: string, props: RestEndpointRuleCreateReq, _options?: Configuration): Promise<HttpInfo<RestEndpointRuleCreateRes>> {
        const result = this.api.endpointEpIdRulePostWithHttpInfo(epId, props, _options);
        return result.toPromise();
    }

    /**
     * @param epId endpoint id
     * @param props rule properties
     */
    public endpointEpIdRulePost(epId: string, props: RestEndpointRuleCreateReq, _options?: Configuration): Promise<RestEndpointRuleCreateRes> {
        const result = this.api.endpointEpIdRulePost(epId, props, _options);
        return result.toPromise();
    }


}



import { ObservableMessageApi } from './ObservableAPI';

import { MessageApiRequestFactory, MessageApiResponseProcessor} from "../apis/MessageApi";
export class PromiseMessageApi {
    private api: ObservableMessageApi

    public constructor(
        configuration: Configuration,
        requestFactory?: MessageApiRequestFactory,
        responseProcessor?: MessageApiResponseProcessor
    ) {
        this.api = new ObservableMessageApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param appId application id
     * @param props message properties
     */
    public applicationAppIdMessagePutWithHttpInfo(appId: string, props: RestMessagePutReq, _options?: Configuration): Promise<HttpInfo<RestMessagePutRes>> {
        const result = this.api.applicationAppIdMessagePutWithHttpInfo(appId, props, _options);
        return result.toPromise();
    }

    /**
     * @param appId application id
     * @param props message properties
     */
    public applicationAppIdMessagePut(appId: string, props: RestMessagePutReq, _options?: Configuration): Promise<RestMessagePutRes> {
        const result = this.api.applicationAppIdMessagePut(appId, props, _options);
        return result.toPromise();
    }


}



