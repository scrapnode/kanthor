import { Configuration } from "../openapi/configuration";
import { PromiseAccountApi } from "../openapi/types/PromiseAPI";
import { RestAccountGetRes } from "../openapi/models/all";

export class Account {
  private readonly api: PromiseAccountApi;

  public constructor(config: Configuration) {
    this.api = new PromiseAccountApi(config);
  }

  public async me(): Promise<RestAccountGetRes> {
    return this.api.accountMeGet();
  }
}
