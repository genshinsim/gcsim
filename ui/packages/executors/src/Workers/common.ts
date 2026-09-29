/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-namespace */
// Requests to the aggregator and sim workers (other than Ready) carry the id of the run they
// belong to, and the responses echo it, so the executor can drop responses from a cancelled run.
//
// A Failed response with fatal set means the worker's wasm instance can't be used anymore: it
// failed to load, or its Go program exited or threw during a call (a fatal error such as running
// out of memory). The executor then terminates the worker and starts a new one when needed.
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
		run: number;
		reason: string;
		fatal?: boolean;
	}

	export function FailedResponse(run: number, reason: string): FailedResponse {
		return { type: Response.Failed, run: run, reason: reason };
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
		run: number;
		cfg: string;
	}

	export function InitializeRequest(
		run: number,
		cfg: string,
	): InitializeRequest {
		return { type: Request.Initialize, run: run, cfg: cfg };
	}

	export interface InitializeResponse {
		type: Response.Initialized;
		run: number;
		result: any;
	}

	export function InitializeResponse(
		run: number,
		result: any,
	): InitializeResponse {
		return { type: Response.Initialized, run: run, result: result };
	}

	export interface AddRequest {
		type: Request.Add;
		run: number;
		result: Uint8Array;
	}

	export function AddRequest(run: number, result: Uint8Array): AddRequest {
		return { type: Request.Add, run: run, result: result };
	}

	export interface AddResponse {
		type: Response.Done;
		run: number;
	}

	export function AddResponse(run: number): AddResponse {
		return { type: Response.Done, run: run };
	}

	export interface FlushRequest {
		type: Request.Flush;
		run: number;
		// the last flush of the run; echoed in the response
		final: boolean;
	}

	export function FlushRequest(run: number, final: boolean): FlushRequest {
		return { type: Request.Flush, run: run, final: final };
	}

	export interface ResultResponse {
		type: Response.Result;
		run: number;
		final: boolean;
		// time the aggregator spent on the flush, in ms
		ms: number;
		result: {
			hash: string;
			stats: any;
		};
	}

	export function ResultResponse(
		run: number,
		final: boolean,
		ms: number,
		result: any,
	): ResultResponse {
		return {
			type: Response.Result,
			run: run,
			final: final,
			ms: ms,
			result: result,
		};
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
		fatal?: boolean;
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
		run: number;
		reason: string;
		fatal?: boolean;
	}

	export function FailedResponse(run: number, reason: string): FailedResponse {
		return { type: Response.Failed, run: run, reason: reason };
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
		run: number;
		cfg: string;
	}

	export function InitializeRequest(
		run: number,
		cfg: string,
	): InitializeRequest {
		return { type: Request.Initialize, run: run, cfg: cfg };
	}

	export interface InitializeResponse {
		type: Response.Initialized;
		run: number;
	}

	export function InitializeResponse(run: number): InitializeResponse {
		return { type: Response.Initialized, run: run };
	}

	export interface RunRequest {
		type: Request.Run;
		run: number;
		itr: number;
	}

	export function RunRequest(run: number, itr: number): RunRequest {
		return { type: Request.Run, run: run, itr: itr };
	}

	export interface RunResponse {
		type: Response.Done;
		run: number;
		result: Uint8Array;
		itr: number;
	}

	export function RunResponse(
		run: number,
		result: Uint8Array,
		itr: number,
	): RunResponse {
		return { type: Response.Done, run: run, result: result, itr: itr };
	}
}
