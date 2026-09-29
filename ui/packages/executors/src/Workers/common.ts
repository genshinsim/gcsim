/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-namespace */
export namespace Aggregator {
	export enum Request {
		Ready = "ready",
		Initialize = "initialize",
		Add = "add",
		Flush = "flush",
	}

	export enum Response {
		Failed = "failed",
		Ready = "ready",
		Initialized = "initialized",
		Done = "done",
		Result = "result",
	}

	export interface FailedResponse {
		type: Response.Failed;
		reason: string;
	}

	export function FailedResponse(reason: string): FailedResponse {
		return { type: Response.Failed, reason: reason };
	}

	export interface ReadyRequest {
		type: Request.Ready;
		module: WebAssembly.Module;
	}

	export function ReadyRequest(module: WebAssembly.Module): ReadyRequest {
		return { type: Request.Ready, module: module };
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
	}

	export function FlushRequest(): FlushRequest {
		return { type: Request.Flush };
	}

	export interface ResultResponse {
		type: Response.Result;
		result: {
			hash: string;
			stats: any;
		};
	}

	export function ResultResponse(result: any): ResultResponse {
		return { type: Response.Result, result: result };
	}
}

export namespace Helper {
	export enum Request {
		Ready = "ready",
		Validate = "validate",
		Sample = "sample",
	}

	export enum Response {
		Failed = "failed",
		Ready = "ready",
		Validate = "validated",
		Sample = "sample",
	}

	// id is missing when loading the wasm failed
	export interface FailedResponse {
		id: number;
		type: Response.Failed;
		reason: string;
	}

	export function FailedResponse(id: number, reason: string): FailedResponse {
		return { id: id, type: Response.Failed, reason: reason };
	}

	export interface ReadyRequest {
		type: Request.Ready;
		wasm: string;
	}

	export function ReadyRequest(wasm: string): ReadyRequest {
		return { type: Request.Ready, wasm: wasm };
	}

	// the compiled module, for the executor to share with the other workers
	export interface ReadyResponse {
		type: Response.Ready;
		module: WebAssembly.Module;
	}

	export interface ValidateRequest {
		id: number;
		type: Request.Validate;
		cfg: string;
	}

	export function ValidateRequest(id: number, cfg: string): ValidateRequest {
		return { id: id, type: Helper.Request.Validate, cfg: cfg };
	}

	export interface ValidateResponse {
		id: number;
		type: Response.Validate;
		cfg: any;
	}

	export function ValidateResponse(id: number, cfg: any): ValidateResponse {
		return { id: id, type: Response.Validate, cfg: cfg };
	}

	export interface SampleRequest {
		id: number;
		type: Request.Sample;
		cfg: string;
		seed: string;
	}

	export function SampleRequest(
		id: number,
		cfg: string,
		seed: string,
	): SampleRequest {
		return { id: id, type: Helper.Request.Sample, cfg: cfg, seed: seed };
	}

	export interface SampleResponse {
		id: number;
		type: Response.Sample;
		sample: any;
	}
}

export namespace SimWorker {
	export enum Request {
		Ready = "ready",
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

	export interface ReadyRequest {
		type: Request.Ready;
		module: WebAssembly.Module;
	}

	export function ReadyRequest(module: WebAssembly.Module): ReadyRequest {
		return { type: Request.Ready, module: module };
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
