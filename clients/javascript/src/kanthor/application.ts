import { Configuration } from "../openapi/configuration";
import { PromiseApplicationApi } from "../openapi/types/PromiseAPI";
import {
  RestApplicationCreateReq,
  RestApplicationCreateRes,
  RestApplicationListRes,
  RestApplicationGetRes,
  RestApplicationUpdateReq,
  RestApplicationUpdateRes,
  RestApplicationDeleteRes,
} from "../openapi/models/all";
import { ListReq } from "./models";

export class Application {
  private readonly api: PromiseApplicationApi;

  public constructor(config: Configuration) {
    this.api = new PromiseApplicationApi(config);
  }

  public async create(
    req: RestApplicationCreateReq
  ): Promise<RestApplicationCreateRes> {
    return this.api.applicationPost(req);
  }

  public async list(req: ListReq): Promise<RestApplicationListRes> {
    return this.api.applicationGet(req.cursor, req.q, req.limit, req.ids);
  }

  public async get(id: string): Promise<RestApplicationGetRes> {
    return this.api.applicationAppIdGet(id);
  }

  public async update(
    id: string,
    req: RestApplicationUpdateReq
  ): Promise<RestApplicationUpdateRes> {
    return this.api.applicationAppIdPut(id, req);
  }

  public async delete(id: string): Promise<RestApplicationDeleteRes> {
    return this.api.applicationAppIdDelete(id);
  }
}
