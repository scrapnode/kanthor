import {
  Middleware,
  RequestContext,
  ResponseContext,
  ServerConfiguration,
  createConfiguration,
  AuthMethodsConfiguration,
  BasicAuthAuthentication
} from "./openapi";
import { Account, Application } from "./kanthor";
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

export interface Options {
  endpoint: string;
}

export interface Option {
  (opts: Options): void;
}

export function withEndpoint(endpoint: string): Option {
  return function (opts: Options) {
    opts.endpoint = endpoint;
  };
}

export class Kanthor {
  public readonly account: Account;
  public readonly application: Application;

  public constructor(token: string, ...options: Option[]) {
    const opts: Options = {
      // @TODO: define cloud server endpoint
      endpoint: "",
    };
    for (let option of options) {
      option(opts);
    }

    const credentials = token.split(":");
    const authMethods: AuthMethodsConfiguration = {
      default: new BasicAuthAuthentication(credentials[0], credentials[1]),
    };

    const conf = createConfiguration({
      baseServer: new ServerConfiguration<any>(opts.endpoint, {}),
      promiseMiddleware: [new UserAgentMiddleware()],
      authMethods: authMethods,
    });

    this.account = new Account(conf);
    this.application = new Application(conf);
  }
}
