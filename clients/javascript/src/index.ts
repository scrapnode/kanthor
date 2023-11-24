import { PromiseMiddleware as Middleware } from "./openapi/middleware";
import { RequestContext, ResponseContext } from "./openapi/http/http";
import { ServerConfiguration } from "./openapi/servers";
import { Configuration, createConfiguration } from "./openapi/configuration";
import { version } from "./metadata.json";

class UserAgentMiddleware implements Middleware {
  public pre(context: RequestContext): Promise<RequestContext> {
    context.setHeaderParam(
      "User-Agent",
      `kanthorlabs/kanthor/${version}/javascript`
    );
    return Promise.resolve(context);
  }

  public post(context: ResponseContext): Promise<ResponseContext> {
    return Promise.resolve(context);
  }
}

// @TODO: define cloud server endpoint

export class Kanthor {
  public readonly configuration: Configuration;

  public constructor(endpoint: string, token: string) {
    const credentials = token.split(":");
    this.configuration = createConfiguration({
      baseServer: new ServerConfiguration<any>(endpoint, {}),
      promiseMiddleware: [new UserAgentMiddleware()],
      authMethods: {
        BasicAuth: {
          username: credentials[0],
          password: credentials[1],
        },
      },
    });
  }
}
