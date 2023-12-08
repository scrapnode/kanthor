import {
  Middleware,
  RequestContext,
  ResponseContext,
  ServerConfiguration,
  createConfiguration,
  AuthMethodsConfiguration,
  BasicAuthAuthentication,
} from "./openapi";
import { Account, Application } from "./kanthor";
import { version } from "./metadata.json";

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

    const [user, pass] = token.split(":");
    const authMethods: AuthMethodsConfiguration = {
      default: new BasicAuthAuthentication(user, pass),
    };

    const conf = createConfiguration({
      baseServer: new ServerConfiguration<any>(opts.endpoint, {}),
      promiseMiddleware: [new UserAgentMiddleware(), new AuthEngineMiddleware(user)],
      authMethods: authMethods,
    });

    this.account = new Account(conf);
    this.application = new Application(conf);
  }
}

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

class AuthEngineMiddleware implements Middleware {
  private readonly user: string;
  constructor(user: string) {
    this.user = user;
  }

  public pre(context: RequestContext): Promise<RequestContext> {
    if (this.user.startsWith("wsc_")) {
      context.setHeaderParam(
        "X-Authorization-Engine",
        `sdk.internal`
      );
    } else {
      context.setHeaderParam(
        "X-Authorization-Engine",
        `ask`
      );
    }
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
