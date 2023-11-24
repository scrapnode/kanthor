import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration} from '../configuration'
import { Observable, of, from } from '../rxjsStub';
import {mergeMap, map} from  '../rxjsStub';
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
import { RestApplicationListReq } from '../models/RestApplicationListReq';
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

import { AccountApiRequestFactory, AccountApiResponseProcessor} from "../apis/AccountApi";
export class ObservableAccountApi {
    private requestFactory: AccountApiRequestFactory;
    private responseProcessor: AccountApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: AccountApiRequestFactory,
        responseProcessor?: AccountApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new AccountApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new AccountApiResponseProcessor();
    }

    /**
     */
    public accountMeGetWithHttpInfo(_options?: Configuration): Observable<HttpInfo<RestAccountGetRes>> {
        const requestContextPromise = this.requestFactory.accountMeGet(_options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.accountMeGetWithHttpInfo(rsp)));
            }));
    }

    /**
     */
    public accountMeGet(_options?: Configuration): Observable<RestAccountGetRes> {
        return this.accountMeGetWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<RestAccountGetRes>) => apiResponse.data));
    }

}

import { ApplicationApiRequestFactory, ApplicationApiResponseProcessor} from "../apis/ApplicationApi";
export class ObservableApplicationApi {
    private requestFactory: ApplicationApiRequestFactory;
    private responseProcessor: ApplicationApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: ApplicationApiRequestFactory,
        responseProcessor?: ApplicationApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new ApplicationApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new ApplicationApiResponseProcessor();
    }

    /**
     * @param appId application id
     */
    public applicationAppIdDeleteWithHttpInfo(appId: string, _options?: Configuration): Observable<HttpInfo<RestApplicationDeleteRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdDelete(appId, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     */
    public applicationAppIdDelete(appId: string, _options?: Configuration): Observable<RestApplicationDeleteRes> {
        return this.applicationAppIdDeleteWithHttpInfo(appId, _options).pipe(map((apiResponse: HttpInfo<RestApplicationDeleteRes>) => apiResponse.data));
    }

    /**
     * @param appId application id
     */
    public applicationAppIdGetWithHttpInfo(appId: string, _options?: Configuration): Observable<HttpInfo<RestApplicationGetRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdGet(appId, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     */
    public applicationAppIdGet(appId: string, _options?: Configuration): Observable<RestApplicationGetRes> {
        return this.applicationAppIdGetWithHttpInfo(appId, _options).pipe(map((apiResponse: HttpInfo<RestApplicationGetRes>) => apiResponse.data));
    }

    /**
     * @param appId application id
     * @param props application properties
     */
    public applicationAppIdPutWithHttpInfo(appId: string, props: RestApplicationUpdateReq, _options?: Configuration): Observable<HttpInfo<RestApplicationUpdateRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdPut(appId, props, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     * @param props application properties
     */
    public applicationAppIdPut(appId: string, props: RestApplicationUpdateReq, _options?: Configuration): Observable<RestApplicationUpdateRes> {
        return this.applicationAppIdPutWithHttpInfo(appId, props, _options).pipe(map((apiResponse: HttpInfo<RestApplicationUpdateRes>) => apiResponse.data));
    }

    /**
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationGetWithHttpInfo(cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Observable<HttpInfo<RestApplicationListReq>> {
        const requestContextPromise = this.requestFactory.applicationGet(cursor, q, limit, id, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationGet(cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Observable<RestApplicationListReq> {
        return this.applicationGetWithHttpInfo(cursor, q, limit, id, _options).pipe(map((apiResponse: HttpInfo<RestApplicationListReq>) => apiResponse.data));
    }

    /**
     * @param props application properties
     */
    public applicationPostWithHttpInfo(props: RestApplicationCreateReq, _options?: Configuration): Observable<HttpInfo<RestApplicationCreateRes>> {
        const requestContextPromise = this.requestFactory.applicationPost(props, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param props application properties
     */
    public applicationPost(props: RestApplicationCreateReq, _options?: Configuration): Observable<RestApplicationCreateRes> {
        return this.applicationPostWithHttpInfo(props, _options).pipe(map((apiResponse: HttpInfo<RestApplicationCreateRes>) => apiResponse.data));
    }

}

import { EndpointApiRequestFactory, EndpointApiResponseProcessor} from "../apis/EndpointApi";
export class ObservableEndpointApi {
    private requestFactory: EndpointApiRequestFactory;
    private responseProcessor: EndpointApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: EndpointApiRequestFactory,
        responseProcessor?: EndpointApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new EndpointApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new EndpointApiResponseProcessor();
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdDeleteWithHttpInfo(appId: string, epId: string, _options?: Configuration): Observable<HttpInfo<RestEndpointDeleteRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdEndpointEpIdDelete(appId, epId, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdEndpointEpIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdDelete(appId: string, epId: string, _options?: Configuration): Observable<RestEndpointDeleteRes> {
        return this.applicationAppIdEndpointEpIdDeleteWithHttpInfo(appId, epId, _options).pipe(map((apiResponse: HttpInfo<RestEndpointDeleteRes>) => apiResponse.data));
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdGetWithHttpInfo(appId: string, epId: string, _options?: Configuration): Observable<HttpInfo<RestEndpointGetRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdEndpointEpIdGet(appId, epId, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdEndpointEpIdGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     */
    public applicationAppIdEndpointEpIdGet(appId: string, epId: string, _options?: Configuration): Observable<RestEndpointGetRes> {
        return this.applicationAppIdEndpointEpIdGetWithHttpInfo(appId, epId, _options).pipe(map((apiResponse: HttpInfo<RestEndpointGetRes>) => apiResponse.data));
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointEpIdPutWithHttpInfo(appId: string, epId: string, props: RestEndpointUpdateReq, _options?: Configuration): Observable<HttpInfo<RestEndpointUpdateRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdEndpointEpIdPut(appId, epId, props, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdEndpointEpIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     * @param epId endpoint id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointEpIdPut(appId: string, epId: string, props: RestEndpointUpdateReq, _options?: Configuration): Observable<RestEndpointUpdateRes> {
        return this.applicationAppIdEndpointEpIdPutWithHttpInfo(appId, epId, props, _options).pipe(map((apiResponse: HttpInfo<RestEndpointUpdateRes>) => apiResponse.data));
    }

    /**
     * @param appId application id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationAppIdEndpointGetWithHttpInfo(appId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Observable<HttpInfo<RestEndpointListRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdEndpointGet(appId, cursor, q, limit, id, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdEndpointGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public applicationAppIdEndpointGet(appId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Observable<RestEndpointListRes> {
        return this.applicationAppIdEndpointGetWithHttpInfo(appId, cursor, q, limit, id, _options).pipe(map((apiResponse: HttpInfo<RestEndpointListRes>) => apiResponse.data));
    }

    /**
     * @param appId application id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointPostWithHttpInfo(appId: string, props: RestEndpointCreateReq, _options?: Configuration): Observable<HttpInfo<RestEndpointCreateRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdEndpointPost(appId, props, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdEndpointPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     * @param props endpoint properties
     */
    public applicationAppIdEndpointPost(appId: string, props: RestEndpointCreateReq, _options?: Configuration): Observable<RestEndpointCreateRes> {
        return this.applicationAppIdEndpointPostWithHttpInfo(appId, props, _options).pipe(map((apiResponse: HttpInfo<RestEndpointCreateRes>) => apiResponse.data));
    }

}

import { EndpointRuleApiRequestFactory, EndpointRuleApiResponseProcessor} from "../apis/EndpointRuleApi";
export class ObservableEndpointRuleApi {
    private requestFactory: EndpointRuleApiRequestFactory;
    private responseProcessor: EndpointRuleApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: EndpointRuleApiRequestFactory,
        responseProcessor?: EndpointRuleApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new EndpointRuleApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new EndpointRuleApiResponseProcessor();
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdDeleteWithHttpInfo(epId: string, eprId: string, _options?: Configuration): Observable<HttpInfo<RestEndpointRuleDeleteRes>> {
        const requestContextPromise = this.requestFactory.endpointEpIdRuleEprIdDelete(epId, eprId, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.endpointEpIdRuleEprIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdDelete(epId: string, eprId: string, _options?: Configuration): Observable<RestEndpointRuleDeleteRes> {
        return this.endpointEpIdRuleEprIdDeleteWithHttpInfo(epId, eprId, _options).pipe(map((apiResponse: HttpInfo<RestEndpointRuleDeleteRes>) => apiResponse.data));
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdGetWithHttpInfo(epId: string, eprId: string, _options?: Configuration): Observable<HttpInfo<RestEndpointRuleGetRes>> {
        const requestContextPromise = this.requestFactory.endpointEpIdRuleEprIdGet(epId, eprId, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.endpointEpIdRuleEprIdGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     */
    public endpointEpIdRuleEprIdGet(epId: string, eprId: string, _options?: Configuration): Observable<RestEndpointRuleGetRes> {
        return this.endpointEpIdRuleEprIdGetWithHttpInfo(epId, eprId, _options).pipe(map((apiResponse: HttpInfo<RestEndpointRuleGetRes>) => apiResponse.data));
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     * @param props rule properties
     */
    public endpointEpIdRuleEprIdPutWithHttpInfo(epId: string, eprId: string, props: RestEndpointRuleUpdateReq, _options?: Configuration): Observable<HttpInfo<RestEndpointRuleUpdateRes>> {
        const requestContextPromise = this.requestFactory.endpointEpIdRuleEprIdPut(epId, eprId, props, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.endpointEpIdRuleEprIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param epId endpoint id
     * @param eprId rule id
     * @param props rule properties
     */
    public endpointEpIdRuleEprIdPut(epId: string, eprId: string, props: RestEndpointRuleUpdateReq, _options?: Configuration): Observable<RestEndpointRuleUpdateRes> {
        return this.endpointEpIdRuleEprIdPutWithHttpInfo(epId, eprId, props, _options).pipe(map((apiResponse: HttpInfo<RestEndpointRuleUpdateRes>) => apiResponse.data));
    }

    /**
     * @param epId endpoint id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public endpointEpIdRuleGetWithHttpInfo(epId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Observable<HttpInfo<RestEndpointRuleListRes>> {
        const requestContextPromise = this.requestFactory.endpointEpIdRuleGet(epId, cursor, q, limit, id, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.endpointEpIdRuleGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param epId endpoint id
     * @param cursor current query cursor
     * @param q search keyword
     * @param limit limit returning records
     * @param id only return records with selected ids
     */
    public endpointEpIdRuleGet(epId: string, cursor?: string, q?: string, limit?: number, id?: Array<string>, _options?: Configuration): Observable<RestEndpointRuleListRes> {
        return this.endpointEpIdRuleGetWithHttpInfo(epId, cursor, q, limit, id, _options).pipe(map((apiResponse: HttpInfo<RestEndpointRuleListRes>) => apiResponse.data));
    }

    /**
     * @param epId endpoint id
     * @param props rule properties
     */
    public endpointEpIdRulePostWithHttpInfo(epId: string, props: RestEndpointRuleCreateReq, _options?: Configuration): Observable<HttpInfo<RestEndpointRuleCreateRes>> {
        const requestContextPromise = this.requestFactory.endpointEpIdRulePost(epId, props, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.endpointEpIdRulePostWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param epId endpoint id
     * @param props rule properties
     */
    public endpointEpIdRulePost(epId: string, props: RestEndpointRuleCreateReq, _options?: Configuration): Observable<RestEndpointRuleCreateRes> {
        return this.endpointEpIdRulePostWithHttpInfo(epId, props, _options).pipe(map((apiResponse: HttpInfo<RestEndpointRuleCreateRes>) => apiResponse.data));
    }

}

import { MessageApiRequestFactory, MessageApiResponseProcessor} from "../apis/MessageApi";
export class ObservableMessageApi {
    private requestFactory: MessageApiRequestFactory;
    private responseProcessor: MessageApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: MessageApiRequestFactory,
        responseProcessor?: MessageApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new MessageApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new MessageApiResponseProcessor();
    }

    /**
     * @param appId application id
     * @param props message properties
     */
    public applicationAppIdMessagePutWithHttpInfo(appId: string, props: RestMessagePutReq, _options?: Configuration): Observable<HttpInfo<RestMessagePutRes>> {
        const requestContextPromise = this.requestFactory.applicationAppIdMessagePut(appId, props, _options);

        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (let middleware of this.configuration.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => this.configuration.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (let middleware of this.configuration.middleware) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.applicationAppIdMessagePutWithHttpInfo(rsp)));
            }));
    }

    /**
     * @param appId application id
     * @param props message properties
     */
    public applicationAppIdMessagePut(appId: string, props: RestMessagePutReq, _options?: Configuration): Observable<RestMessagePutRes> {
        return this.applicationAppIdMessagePutWithHttpInfo(appId, props, _options).pipe(map((apiResponse: HttpInfo<RestMessagePutRes>) => apiResponse.data));
    }

}
