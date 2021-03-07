import { Request, Response } from "./http";

export enum MiddlewareT {
  SERVICE,
  ROUTER,
  HANDLER_FN,
  MIDDLEWARE_FN,
}

export type Middleware = HandlerFn | MiddlewareFn | MiddlewareObj;

export interface NextFn {
  (err?: any): void;
  (deferToNext: "router"): void;
}

export type HandlerFn = (req: Request, res: Response) => void;

export type MiddlewareFn = (req: Request, res: Response, next: NextFn) => void;

export interface MiddlewareObj {
  handle: HandlerFn | MiddlewareFn;
  type: MiddlewareT;
}
