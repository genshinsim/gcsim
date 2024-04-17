/* eslint-disable */
import _m0 from "protobufjs/minimal";

export interface ComputeFailedEvent {
  dbId?: string | undefined;
  config?: string | undefined;
  submitter?: string | undefined;
  reason?: string | undefined;
}

export interface ComputeCompletedEvent {
  dbId?: string | undefined;
  shareId?: string | undefined;
}

export interface SubmissionDeleteEvent {
  dbId?: string | undefined;
  config?: string | undefined;
  submitter?: string | undefined;
}

export interface EntryReplaceEvent {
  dbId?: string | undefined;
  config?: string | undefined;
  oldConfig?: string | undefined;
}

export interface DescReplaceEvent {
  dbId?: string | undefined;
  desc?: string | undefined;
  oldDesc?: string | undefined;
}

function createBaseComputeFailedEvent(): ComputeFailedEvent {
  return { dbId: "", config: "", submitter: "", reason: "" };
}

export const ComputeFailedEvent = {
  encode(message: ComputeFailedEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.dbId !== undefined && message.dbId !== "") {
      writer.uint32(10).string(message.dbId);
    }
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(18).string(message.config);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(26).string(message.submitter);
    }
    if (message.reason !== undefined && message.reason !== "") {
      writer.uint32(34).string(message.reason);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ComputeFailedEvent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseComputeFailedEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.dbId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.config = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.submitter = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.reason = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ComputeFailedEvent {
    return {
      dbId: isSet(object.dbId) ? globalThis.String(object.dbId) : "",
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
      reason: isSet(object.reason) ? globalThis.String(object.reason) : "",
    };
  },

  toJSON(message: ComputeFailedEvent): unknown {
    const obj: any = {};
    if (message.dbId !== undefined && message.dbId !== "") {
      obj.dbId = message.dbId;
    }
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    if (message.reason !== undefined && message.reason !== "") {
      obj.reason = message.reason;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ComputeFailedEvent>, I>>(base?: I): ComputeFailedEvent {
    return ComputeFailedEvent.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ComputeFailedEvent>, I>>(object: I): ComputeFailedEvent {
    const message = createBaseComputeFailedEvent();
    message.dbId = object.dbId ?? "";
    message.config = object.config ?? "";
    message.submitter = object.submitter ?? "";
    message.reason = object.reason ?? "";
    return message;
  },
};

function createBaseComputeCompletedEvent(): ComputeCompletedEvent {
  return { dbId: "", shareId: "" };
}

export const ComputeCompletedEvent = {
  encode(message: ComputeCompletedEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.dbId !== undefined && message.dbId !== "") {
      writer.uint32(10).string(message.dbId);
    }
    if (message.shareId !== undefined && message.shareId !== "") {
      writer.uint32(18).string(message.shareId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ComputeCompletedEvent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseComputeCompletedEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.dbId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.shareId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ComputeCompletedEvent {
    return {
      dbId: isSet(object.dbId) ? globalThis.String(object.dbId) : "",
      shareId: isSet(object.shareId) ? globalThis.String(object.shareId) : "",
    };
  },

  toJSON(message: ComputeCompletedEvent): unknown {
    const obj: any = {};
    if (message.dbId !== undefined && message.dbId !== "") {
      obj.dbId = message.dbId;
    }
    if (message.shareId !== undefined && message.shareId !== "") {
      obj.shareId = message.shareId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ComputeCompletedEvent>, I>>(base?: I): ComputeCompletedEvent {
    return ComputeCompletedEvent.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ComputeCompletedEvent>, I>>(object: I): ComputeCompletedEvent {
    const message = createBaseComputeCompletedEvent();
    message.dbId = object.dbId ?? "";
    message.shareId = object.shareId ?? "";
    return message;
  },
};

function createBaseSubmissionDeleteEvent(): SubmissionDeleteEvent {
  return { dbId: "", config: "", submitter: "" };
}

export const SubmissionDeleteEvent = {
  encode(message: SubmissionDeleteEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.dbId !== undefined && message.dbId !== "") {
      writer.uint32(10).string(message.dbId);
    }
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(18).string(message.config);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(26).string(message.submitter);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubmissionDeleteEvent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubmissionDeleteEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.dbId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.config = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.submitter = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubmissionDeleteEvent {
    return {
      dbId: isSet(object.dbId) ? globalThis.String(object.dbId) : "",
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
    };
  },

  toJSON(message: SubmissionDeleteEvent): unknown {
    const obj: any = {};
    if (message.dbId !== undefined && message.dbId !== "") {
      obj.dbId = message.dbId;
    }
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubmissionDeleteEvent>, I>>(base?: I): SubmissionDeleteEvent {
    return SubmissionDeleteEvent.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SubmissionDeleteEvent>, I>>(object: I): SubmissionDeleteEvent {
    const message = createBaseSubmissionDeleteEvent();
    message.dbId = object.dbId ?? "";
    message.config = object.config ?? "";
    message.submitter = object.submitter ?? "";
    return message;
  },
};

function createBaseEntryReplaceEvent(): EntryReplaceEvent {
  return { dbId: "", config: "", oldConfig: "" };
}

export const EntryReplaceEvent = {
  encode(message: EntryReplaceEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.dbId !== undefined && message.dbId !== "") {
      writer.uint32(10).string(message.dbId);
    }
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(18).string(message.config);
    }
    if (message.oldConfig !== undefined && message.oldConfig !== "") {
      writer.uint32(26).string(message.oldConfig);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EntryReplaceEvent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntryReplaceEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.dbId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.config = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.oldConfig = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EntryReplaceEvent {
    return {
      dbId: isSet(object.dbId) ? globalThis.String(object.dbId) : "",
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      oldConfig: isSet(object.oldConfig) ? globalThis.String(object.oldConfig) : "",
    };
  },

  toJSON(message: EntryReplaceEvent): unknown {
    const obj: any = {};
    if (message.dbId !== undefined && message.dbId !== "") {
      obj.dbId = message.dbId;
    }
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.oldConfig !== undefined && message.oldConfig !== "") {
      obj.oldConfig = message.oldConfig;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EntryReplaceEvent>, I>>(base?: I): EntryReplaceEvent {
    return EntryReplaceEvent.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EntryReplaceEvent>, I>>(object: I): EntryReplaceEvent {
    const message = createBaseEntryReplaceEvent();
    message.dbId = object.dbId ?? "";
    message.config = object.config ?? "";
    message.oldConfig = object.oldConfig ?? "";
    return message;
  },
};

function createBaseDescReplaceEvent(): DescReplaceEvent {
  return { dbId: "", desc: "", oldDesc: "" };
}

export const DescReplaceEvent = {
  encode(message: DescReplaceEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.dbId !== undefined && message.dbId !== "") {
      writer.uint32(10).string(message.dbId);
    }
    if (message.desc !== undefined && message.desc !== "") {
      writer.uint32(18).string(message.desc);
    }
    if (message.oldDesc !== undefined && message.oldDesc !== "") {
      writer.uint32(26).string(message.oldDesc);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DescReplaceEvent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDescReplaceEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.dbId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.desc = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.oldDesc = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DescReplaceEvent {
    return {
      dbId: isSet(object.dbId) ? globalThis.String(object.dbId) : "",
      desc: isSet(object.desc) ? globalThis.String(object.desc) : "",
      oldDesc: isSet(object.oldDesc) ? globalThis.String(object.oldDesc) : "",
    };
  },

  toJSON(message: DescReplaceEvent): unknown {
    const obj: any = {};
    if (message.dbId !== undefined && message.dbId !== "") {
      obj.dbId = message.dbId;
    }
    if (message.desc !== undefined && message.desc !== "") {
      obj.desc = message.desc;
    }
    if (message.oldDesc !== undefined && message.oldDesc !== "") {
      obj.oldDesc = message.oldDesc;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DescReplaceEvent>, I>>(base?: I): DescReplaceEvent {
    return DescReplaceEvent.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DescReplaceEvent>, I>>(object: I): DescReplaceEvent {
    const message = createBaseDescReplaceEvent();
    message.dbId = object.dbId ?? "";
    message.desc = object.desc ?? "";
    message.oldDesc = object.oldDesc ?? "";
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
