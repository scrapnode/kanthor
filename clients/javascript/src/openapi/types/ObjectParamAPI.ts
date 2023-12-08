import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration} from '../configuration'

import { AuthenticatorAccount } from '../models/AuthenticatorAccount';
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

import { ObservableAccountApi } from "./ObservableAPI";
import { AccountApiRequestFactory, AccountApiResponseProcessor} from "../apis/AccountApi";

export interface AccountApiAccountMeGetRequest {
}

export class ObjectAccountApi {
    private api: ObservableAccountApi

    public constructor(configuration: Configuration, requestFactory?: AccountApiRequestFactory, responseProcessor?: AccountApiResponseProcessor) {
        this.api = new ObservableAccountApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param param the request object
     */
    public accountMeGetWithHttpInfo(param: AccountApiAccountMeGetRequest = {}, options?: Configuration): Promise<HttpInfo<RestAccountGetRes>> {
        return this.api.accountMeGetWithHttpInfo( options).toPromise();
    }

    /**
     * @param param the request object
     */
    public accountMeGet(param: AccountApiAccountMeGetRequest = {}, options?: Configuration): Promise<RestAccountGetRes> {
        return this.api.accountMeGet( options).toPromise();
    }

}

import { ObservableApplicationApi } from "./ObservableAPI";
import { ApplicationApiRequestFactory, ApplicationApiResponseProcessor} from "../apis/ApplicationApi";

export interface ApplicationApiApplicationAppIdDeleteRequest {
    /**
     * application id
     * @type string
     * @memberof ApplicationApiapplicationAppIdDelete
     */
    appId: string
}

export interface ApplicationApiApplicationAppIdGetRequest {
    /**
     * application id
     * @type string
     * @memberof ApplicationApiapplicationAppIdGet
     */
    appId: string
}

export interface ApplicationApiApplicationAppIdPutRequest {
    /**
     * application id
     * @type string
     * @memberof ApplicationApiapplicationAppIdPut
     */
    appId: string
    /**
     * application properties
     * @type RestApplicationUpdateReq
     * @memberof ApplicationApiapplicationAppIdPut
     */
    props: RestApplicationUpdateReq
}

export interface ApplicationApiApplicationGetRequest {
    /**
     * current query cursor
     * @type string
     * @memberof ApplicationApiapplicationGet
     */
    cursor?: string
    /**
     * search keyword
     * @type string
     * @memberof ApplicationApiapplicationGet
     */
    q?: string
    /**
     * limit returning records
     * @type number
     * @memberof ApplicationApiapplicationGet
     */
    limit?: number
    /**
     * only return records with selected ids
     * @type Array&lt;string&gt;
     * @memberof ApplicationApiapplicationGet
     */
    id?: Array<string>
}

export interface ApplicationApiApplicationPostRequest {
    /**
     * application properties
     * @type RestApplicationCreateReq
     * @memberof ApplicationApiapplicationPost
     */
    props: RestApplicationCreateReq
}

export class ObjectApplicationApi {
    private api: ObservableApplicationApi

    public constructor(configuration: Configuration, requestFactory?: ApplicationApiRequestFactory, responseProcessor?: ApplicationApiResponseProcessor) {
        this.api = new ObservableApplicationApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param param the request object
     */
    public applicationAppIdDeleteWithHttpInfo(param: ApplicationApiApplicationAppIdDeleteRequest, options?: Configuration): Promise<HttpInfo<RestApplicationDeleteRes>> {
        return this.api.applicationAppIdDeleteWithHttpInfo(param.appId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdDelete(param: ApplicationApiApplicationAppIdDeleteRequest, options?: Configuration): Promise<RestApplicationDeleteRes> {
        return this.api.applicationAppIdDelete(param.appId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdGetWithHttpInfo(param: ApplicationApiApplicationAppIdGetRequest, options?: Configuration): Promise<HttpInfo<RestApplicationGetRes>> {
        return this.api.applicationAppIdGetWithHttpInfo(param.appId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdGet(param: ApplicationApiApplicationAppIdGetRequest, options?: Configuration): Promise<RestApplicationGetRes> {
        return this.api.applicationAppIdGet(param.appId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdPutWithHttpInfo(param: ApplicationApiApplicationAppIdPutRequest, options?: Configuration): Promise<HttpInfo<RestApplicationUpdateRes>> {
        return this.api.applicationAppIdPutWithHttpInfo(param.appId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdPut(param: ApplicationApiApplicationAppIdPutRequest, options?: Configuration): Promise<RestApplicationUpdateRes> {
        return this.api.applicationAppIdPut(param.appId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationGetWithHttpInfo(param: ApplicationApiApplicationGetRequest = {}, options?: Configuration): Promise<HttpInfo<RestApplicationListRes>> {
        return this.api.applicationGetWithHttpInfo(param.cursor, param.q, param.limit, param.id,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationGet(param: ApplicationApiApplicationGetRequest = {}, options?: Configuration): Promise<RestApplicationListRes> {
        return this.api.applicationGet(param.cursor, param.q, param.limit, param.id,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationPostWithHttpInfo(param: ApplicationApiApplicationPostRequest, options?: Configuration): Promise<HttpInfo<RestApplicationCreateRes>> {
        return this.api.applicationPostWithHttpInfo(param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationPost(param: ApplicationApiApplicationPostRequest, options?: Configuration): Promise<RestApplicationCreateRes> {
        return this.api.applicationPost(param.props,  options).toPromise();
    }

}

import { ObservableEndpointApi } from "./ObservableAPI";
import { EndpointApiRequestFactory, EndpointApiResponseProcessor} from "../apis/EndpointApi";

export interface EndpointApiApplicationAppIdEndpointEpIdDeleteRequest {
    /**
     * application id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointEpIdDelete
     */
    appId: string
    /**
     * endpoint id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointEpIdDelete
     */
    epId: string
}

export interface EndpointApiApplicationAppIdEndpointEpIdGetRequest {
    /**
     * application id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointEpIdGet
     */
    appId: string
    /**
     * endpoint id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointEpIdGet
     */
    epId: string
}

export interface EndpointApiApplicationAppIdEndpointEpIdPutRequest {
    /**
     * application id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointEpIdPut
     */
    appId: string
    /**
     * endpoint id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointEpIdPut
     */
    epId: string
    /**
     * endpoint properties
     * @type RestEndpointUpdateReq
     * @memberof EndpointApiapplicationAppIdEndpointEpIdPut
     */
    props: RestEndpointUpdateReq
}

export interface EndpointApiApplicationAppIdEndpointGetRequest {
    /**
     * application id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointGet
     */
    appId: string
    /**
     * current query cursor
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointGet
     */
    cursor?: string
    /**
     * search keyword
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointGet
     */
    q?: string
    /**
     * limit returning records
     * @type number
     * @memberof EndpointApiapplicationAppIdEndpointGet
     */
    limit?: number
    /**
     * only return records with selected ids
     * @type Array&lt;string&gt;
     * @memberof EndpointApiapplicationAppIdEndpointGet
     */
    id?: Array<string>
}

export interface EndpointApiApplicationAppIdEndpointPostRequest {
    /**
     * application id
     * @type string
     * @memberof EndpointApiapplicationAppIdEndpointPost
     */
    appId: string
    /**
     * endpoint properties
     * @type RestEndpointCreateReq
     * @memberof EndpointApiapplicationAppIdEndpointPost
     */
    props: RestEndpointCreateReq
}

export class ObjectEndpointApi {
    private api: ObservableEndpointApi

    public constructor(configuration: Configuration, requestFactory?: EndpointApiRequestFactory, responseProcessor?: EndpointApiResponseProcessor) {
        this.api = new ObservableEndpointApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointEpIdDeleteWithHttpInfo(param: EndpointApiApplicationAppIdEndpointEpIdDeleteRequest, options?: Configuration): Promise<HttpInfo<RestEndpointDeleteRes>> {
        return this.api.applicationAppIdEndpointEpIdDeleteWithHttpInfo(param.appId, param.epId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointEpIdDelete(param: EndpointApiApplicationAppIdEndpointEpIdDeleteRequest, options?: Configuration): Promise<RestEndpointDeleteRes> {
        return this.api.applicationAppIdEndpointEpIdDelete(param.appId, param.epId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointEpIdGetWithHttpInfo(param: EndpointApiApplicationAppIdEndpointEpIdGetRequest, options?: Configuration): Promise<HttpInfo<RestEndpointGetRes>> {
        return this.api.applicationAppIdEndpointEpIdGetWithHttpInfo(param.appId, param.epId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointEpIdGet(param: EndpointApiApplicationAppIdEndpointEpIdGetRequest, options?: Configuration): Promise<RestEndpointGetRes> {
        return this.api.applicationAppIdEndpointEpIdGet(param.appId, param.epId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointEpIdPutWithHttpInfo(param: EndpointApiApplicationAppIdEndpointEpIdPutRequest, options?: Configuration): Promise<HttpInfo<RestEndpointUpdateRes>> {
        return this.api.applicationAppIdEndpointEpIdPutWithHttpInfo(param.appId, param.epId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointEpIdPut(param: EndpointApiApplicationAppIdEndpointEpIdPutRequest, options?: Configuration): Promise<RestEndpointUpdateRes> {
        return this.api.applicationAppIdEndpointEpIdPut(param.appId, param.epId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointGetWithHttpInfo(param: EndpointApiApplicationAppIdEndpointGetRequest, options?: Configuration): Promise<HttpInfo<RestEndpointListRes>> {
        return this.api.applicationAppIdEndpointGetWithHttpInfo(param.appId, param.cursor, param.q, param.limit, param.id,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointGet(param: EndpointApiApplicationAppIdEndpointGetRequest, options?: Configuration): Promise<RestEndpointListRes> {
        return this.api.applicationAppIdEndpointGet(param.appId, param.cursor, param.q, param.limit, param.id,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointPostWithHttpInfo(param: EndpointApiApplicationAppIdEndpointPostRequest, options?: Configuration): Promise<HttpInfo<RestEndpointCreateRes>> {
        return this.api.applicationAppIdEndpointPostWithHttpInfo(param.appId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdEndpointPost(param: EndpointApiApplicationAppIdEndpointPostRequest, options?: Configuration): Promise<RestEndpointCreateRes> {
        return this.api.applicationAppIdEndpointPost(param.appId, param.props,  options).toPromise();
    }

}

import { ObservableEndpointRuleApi } from "./ObservableAPI";
import { EndpointRuleApiRequestFactory, EndpointRuleApiResponseProcessor} from "../apis/EndpointRuleApi";

export interface EndpointRuleApiEndpointEpIdRuleEprIdDeleteRequest {
    /**
     * endpoint id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleEprIdDelete
     */
    epId: string
    /**
     * rule id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleEprIdDelete
     */
    eprId: string
}

export interface EndpointRuleApiEndpointEpIdRuleEprIdGetRequest {
    /**
     * endpoint id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleEprIdGet
     */
    epId: string
    /**
     * rule id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleEprIdGet
     */
    eprId: string
}

export interface EndpointRuleApiEndpointEpIdRuleEprIdPutRequest {
    /**
     * endpoint id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleEprIdPut
     */
    epId: string
    /**
     * rule id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleEprIdPut
     */
    eprId: string
    /**
     * rule properties
     * @type RestEndpointRuleUpdateReq
     * @memberof EndpointRuleApiendpointEpIdRuleEprIdPut
     */
    props: RestEndpointRuleUpdateReq
}

export interface EndpointRuleApiEndpointEpIdRuleGetRequest {
    /**
     * endpoint id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleGet
     */
    epId: string
    /**
     * current query cursor
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleGet
     */
    cursor?: string
    /**
     * search keyword
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRuleGet
     */
    q?: string
    /**
     * limit returning records
     * @type number
     * @memberof EndpointRuleApiendpointEpIdRuleGet
     */
    limit?: number
    /**
     * only return records with selected ids
     * @type Array&lt;string&gt;
     * @memberof EndpointRuleApiendpointEpIdRuleGet
     */
    id?: Array<string>
}

export interface EndpointRuleApiEndpointEpIdRulePostRequest {
    /**
     * endpoint id
     * @type string
     * @memberof EndpointRuleApiendpointEpIdRulePost
     */
    epId: string
    /**
     * rule properties
     * @type RestEndpointRuleCreateReq
     * @memberof EndpointRuleApiendpointEpIdRulePost
     */
    props: RestEndpointRuleCreateReq
}

export class ObjectEndpointRuleApi {
    private api: ObservableEndpointRuleApi

    public constructor(configuration: Configuration, requestFactory?: EndpointRuleApiRequestFactory, responseProcessor?: EndpointRuleApiResponseProcessor) {
        this.api = new ObservableEndpointRuleApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleEprIdDeleteWithHttpInfo(param: EndpointRuleApiEndpointEpIdRuleEprIdDeleteRequest, options?: Configuration): Promise<HttpInfo<RestEndpointRuleDeleteRes>> {
        return this.api.endpointEpIdRuleEprIdDeleteWithHttpInfo(param.epId, param.eprId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleEprIdDelete(param: EndpointRuleApiEndpointEpIdRuleEprIdDeleteRequest, options?: Configuration): Promise<RestEndpointRuleDeleteRes> {
        return this.api.endpointEpIdRuleEprIdDelete(param.epId, param.eprId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleEprIdGetWithHttpInfo(param: EndpointRuleApiEndpointEpIdRuleEprIdGetRequest, options?: Configuration): Promise<HttpInfo<RestEndpointRuleGetRes>> {
        return this.api.endpointEpIdRuleEprIdGetWithHttpInfo(param.epId, param.eprId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleEprIdGet(param: EndpointRuleApiEndpointEpIdRuleEprIdGetRequest, options?: Configuration): Promise<RestEndpointRuleGetRes> {
        return this.api.endpointEpIdRuleEprIdGet(param.epId, param.eprId,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleEprIdPutWithHttpInfo(param: EndpointRuleApiEndpointEpIdRuleEprIdPutRequest, options?: Configuration): Promise<HttpInfo<RestEndpointRuleUpdateRes>> {
        return this.api.endpointEpIdRuleEprIdPutWithHttpInfo(param.epId, param.eprId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleEprIdPut(param: EndpointRuleApiEndpointEpIdRuleEprIdPutRequest, options?: Configuration): Promise<RestEndpointRuleUpdateRes> {
        return this.api.endpointEpIdRuleEprIdPut(param.epId, param.eprId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleGetWithHttpInfo(param: EndpointRuleApiEndpointEpIdRuleGetRequest, options?: Configuration): Promise<HttpInfo<RestEndpointRuleListRes>> {
        return this.api.endpointEpIdRuleGetWithHttpInfo(param.epId, param.cursor, param.q, param.limit, param.id,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRuleGet(param: EndpointRuleApiEndpointEpIdRuleGetRequest, options?: Configuration): Promise<RestEndpointRuleListRes> {
        return this.api.endpointEpIdRuleGet(param.epId, param.cursor, param.q, param.limit, param.id,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRulePostWithHttpInfo(param: EndpointRuleApiEndpointEpIdRulePostRequest, options?: Configuration): Promise<HttpInfo<RestEndpointRuleCreateRes>> {
        return this.api.endpointEpIdRulePostWithHttpInfo(param.epId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public endpointEpIdRulePost(param: EndpointRuleApiEndpointEpIdRulePostRequest, options?: Configuration): Promise<RestEndpointRuleCreateRes> {
        return this.api.endpointEpIdRulePost(param.epId, param.props,  options).toPromise();
    }

}

import { ObservableMessageApi } from "./ObservableAPI";
import { MessageApiRequestFactory, MessageApiResponseProcessor} from "../apis/MessageApi";

export interface MessageApiApplicationAppIdMessagePutRequest {
    /**
     * application id
     * @type string
     * @memberof MessageApiapplicationAppIdMessagePut
     */
    appId: string
    /**
     * message properties
     * @type RestMessagePutReq
     * @memberof MessageApiapplicationAppIdMessagePut
     */
    props: RestMessagePutReq
}

export class ObjectMessageApi {
    private api: ObservableMessageApi

    public constructor(configuration: Configuration, requestFactory?: MessageApiRequestFactory, responseProcessor?: MessageApiResponseProcessor) {
        this.api = new ObservableMessageApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param param the request object
     */
    public applicationAppIdMessagePutWithHttpInfo(param: MessageApiApplicationAppIdMessagePutRequest, options?: Configuration): Promise<HttpInfo<RestMessagePutRes>> {
        return this.api.applicationAppIdMessagePutWithHttpInfo(param.appId, param.props,  options).toPromise();
    }

    /**
     * @param param the request object
     */
    public applicationAppIdMessagePut(param: MessageApiApplicationAppIdMessagePutRequest, options?: Configuration): Promise<RestMessagePutRes> {
        return this.api.applicationAppIdMessagePut(param.appId, param.props,  options).toPromise();
    }

}
