/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-namespace */
export namespace Aggregator {
  export enum Request {
    Initialize = "initialize",
    Add = "add",
    Flush = "flush",
    Validate = "validate",
    BuildInfo = "build_info",
  }

  export enum Response {
    Failed = "failed",
    Ready = "ready",
    Initialized = "initialized",
    Done = "done",
    Result = "result",
    Validate = "validated",
    BuildInfo = "build_info",
  }

  export interface FailedResponse {
    type: Response.Failed;
    reason: string;
  }

  export function FailedResponse(reason: string): FailedResponse {
    return { type: Response.Failed, reason: reason };
  }
  
  export interface ReadyResponse {
    type: Response.Ready;
  }
  
  export function ReadyResponse(): ReadyResponse {
    return { type: Response.Ready };
  }

  export interface InitializeRequest {
    type: Request.Initialize;
    cfg: string;
  }

  export function InitializeRequest(cfg: string): InitializeRequest {
    return { type: Request.Initialize, cfg: cfg };
  }

  export interface InitializeResponse {
    type: Response.Initialized;
    result: any;
  }

  export function InitializeResponse(result: any): InitializeResponse {
    return { type: Response.Initialized, result: result };
  }

  export interface AddRequest {
    type: Request.Add;
    result: Uint8Array;
  }

  export function AddRequest(result: Uint8Array): AddRequest {
    return { type: Request.Add, result: result };
  }

  export interface AddResponse {
    type: Response.Done;
  }

  export function AddResponse(): AddResponse {
    return { type: Response.Done };
  }

  export interface FlushRequest {
    type: Request.Flush;
    startTime: number;
  }

  export function FlushRequest(startTime: number): FlushRequest {
    return { type: Request.Flush, startTime: startTime };
  }

  export interface ResultResponse {
    type: Response.Result;
    result: any;
  }

  export function ResultResponse(result: any): ResultResponse {
    return { type: Response.Result, result: result };
  }
}

export namespace Helper {
  export enum Request {
    Validate = "validate",
    GenerateDebug = "gen_debug",
  }

  export enum Response {
    Failed = "failed",
    Validate = "validated",
    GenerateDebug = "gen_debug",
  }

  export interface FailedResponse {
    type: Response.Failed;
    reason: string;
  }

  export function FailedResponse(reason: string): FailedResponse {
    return { type: Response.Failed, reason: reason };
  }

  export interface ValidateRequest {
    type: Request.Validate;
    cfg: string;
  }

  export function ValidateRequest(cfg: string): ValidateRequest {
    return { type: Helper.Request.Validate, cfg: cfg };
  }

  export interface ValidateResponse {
    type: Response.Validate;
    cfg: any;
  }

  export function ValidateResponse(cfg: any): ValidateResponse {
    return { type: Response.Validate, cfg: cfg };
  }

  export interface GenerateDebugRequest {
    type: Request.GenerateDebug;
    cfg: string;
    seed: string;
  }

  export function GenerateDebugRequest(cfg: string, seed: string): GenerateDebugRequest {
    return { type: Helper.Request.GenerateDebug, cfg: cfg, seed: seed };
  }

  export interface GenerateDebugResponse {
    type: Response.GenerateDebug;
    debug: any[];
  }
}

export namespace SimWorker {
  export enum Request {
    Initialize = "initialize",
    Run = "run",
  }

  export enum Response {
    Failed = "failed",
    Ready = "ready",
    Initialized = "initialized",
    Done = "done",
  }

  export interface FailedResponse {
    type: Response.Failed;
    reason: string;
  }

  export function FailedResponse(reason: string): FailedResponse {
    return { type: Response.Failed, reason: reason };
  }
  
  export interface ReadyResponse {
    type: Response.Ready;
  }
  
  export function ReadyResponse(): ReadyResponse {
    return { type: Response.Ready };
  }

  export interface InitializeRequest {
    type: Request.Initialize;
    cfg: string;
  }

  export function InitializeRequest(cfg: string): InitializeRequest {
    return { type: Request.Initialize, cfg: cfg };
  }

  export interface InitializeResponse {
    type: Response.Initialized;
  }

  export function InitializeResponse(): InitializeResponse {
    return { type: Response.Initialized };
  }

  export interface RunRequest {
    type: Request.Run;
    itr: number;
  }

  export function RunRequest(itr: number): RunRequest {
    return { type: Request.Run, itr: itr };
  }

  export interface RunResponse {
    type: Response.Done;
    result: Uint8Array;
    itr: number;
  }

  export function RunResponse(result: Uint8Array, itr: number): RunResponse {
    return { type: Response.Done, result: result, itr: itr };
  }
}