import { Logger } from "./logging";
import { Version } from "./version";
import { Middleware, MiddlewareObj, HandlerFn } from "./middleware";
import { ServiceImpl } from "./internal/service-impl";

export function fold(version?: Version): Service {
  return new ServiceImpl(version);
}

export interface Service extends MiddlewareObj {
  logger: Logger;
  name: string;
  settings: ServiceSettings;
  version: Version;

  start: () => void;
  get: (path: string, handler: HandlerFn) => void;
  head: (path: string, handler: HandlerFn) => void;
  post: (path: string, handler: HandlerFn) => void;
  put: (path: string, handler: HandlerFn) => void;
  delete: (path: string, handler: HandlerFn) => void;
  connect: (path: string, handler: HandlerFn) => void;
  options: (path: string, handler: HandlerFn) => void;
  trace: (path: string, handler: HandlerFn) => void;
  patch: (path: string, handler: HandlerFn) => void;
  all: (path: string, handler: HandlerFn) => void;
  use: (use: Middleware | string, middleware?: Middleware) => void;
}

export interface ServiceSettings {
  get: (setting: string) => any | undefined;
  set: (setting: string, value: any) => void;
  enable: (setting: string) => void;
  enabled: (setting: string) => boolean;
  disable: (setting: string) => void;
}
