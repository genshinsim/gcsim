/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { SimulationResult } from "../model/result";

export interface ShareEntry {
  _id?: string | undefined;
  result?: SimulationResult | undefined;
  expires_at?: number | undefined;
  submitter?: string | undefined;
}

export interface CreateRequest {
  result?: SimulationResult | undefined;
  expires_at?: number | undefined;
  submitter?: string | undefined;
}

export interface CreateResponse {
  id?: string | undefined;
}

export interface ReadRequest {
  id?: string | undefined;
}

export interface ReadResponse {
  id?: string | undefined;
  result?: SimulationResult | undefined;
  expires_at?: number | undefined;
}

export interface UpdateRequest {
  id?: string | undefined;
  result?: SimulationResult | undefined;
  expires_at?: number | undefined;
  submitter?: string | undefined;
}

export interface UpdateResponse {
  id?: string | undefined;
}

export interface SetTTLRequest {
  id?: string | undefined;
  expires_at?: number | undefined;
}

export interface SetTTLResponse {
  id?: string | undefined;
}

export interface DeleteRequest {
  id?: string | undefined;
}

export interface DeleteResponse {
  /** TODO: add deleted data to response in future */
  id?: string | undefined;
}

export interface RandomRequest {
}

export interface RandomResponse {
  id?: string | undefined;
}

function createBaseShareEntry(): ShareEntry {
  return { _id: "", result: undefined, expires_at: 0, submitter: "" };
}

export const ShareEntry = {
  encode(message: ShareEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    if (message.result !== undefined) {
      SimulationResult.encode(message.result, writer.uint32(18).fork()).ldelim();
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      writer.uint32(24).uint64(message.expires_at);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(34).string(message.submitter);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ShareEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseShareEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message._id = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.result = SimulationResult.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.expires_at = longToNumber(reader.uint64() as Long);
          continue;
        case 4:
          if (tag !== 34) {
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

  fromJSON(object: any): ShareEntry {
    return {
      _id: isSet(object._id) ? globalThis.String(object._id) : "",
      result: isSet(object.result) ? SimulationResult.fromJSON(object.result) : undefined,
      expires_at: isSet(object.expires_at) ? globalThis.Number(object.expires_at) : 0,
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
    };
  },

  toJSON(message: ShareEntry): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    if (message.result !== undefined) {
      obj.result = SimulationResult.toJSON(message.result);
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      obj.expires_at = Math.round(message.expires_at);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ShareEntry>, I>>(base?: I): ShareEntry {
    return ShareEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ShareEntry>, I>>(object: I): ShareEntry {
    const message = createBaseShareEntry();
    message._id = object._id ?? "";
    message.result = (object.result !== undefined && object.result !== null)
      ? SimulationResult.fromPartial(object.result)
      : undefined;
    message.expires_at = object.expires_at ?? 0;
    message.submitter = object.submitter ?? "";
    return message;
  },
};

function createBaseCreateRequest(): CreateRequest {
  return { result: undefined, expires_at: 0, submitter: "" };
}

export const CreateRequest = {
  encode(message: CreateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.result !== undefined) {
      SimulationResult.encode(message.result, writer.uint32(10).fork()).ldelim();
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      writer.uint32(16).uint64(message.expires_at);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(26).string(message.submitter);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.result = SimulationResult.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.expires_at = longToNumber(reader.uint64() as Long);
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

  fromJSON(object: any): CreateRequest {
    return {
      result: isSet(object.result) ? SimulationResult.fromJSON(object.result) : undefined,
      expires_at: isSet(object.expires_at) ? globalThis.Number(object.expires_at) : 0,
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
    };
  },

  toJSON(message: CreateRequest): unknown {
    const obj: any = {};
    if (message.result !== undefined) {
      obj.result = SimulationResult.toJSON(message.result);
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      obj.expires_at = Math.round(message.expires_at);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateRequest>, I>>(base?: I): CreateRequest {
    return CreateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateRequest>, I>>(object: I): CreateRequest {
    const message = createBaseCreateRequest();
    message.result = (object.result !== undefined && object.result !== null)
      ? SimulationResult.fromPartial(object.result)
      : undefined;
    message.expires_at = object.expires_at ?? 0;
    message.submitter = object.submitter ?? "";
    return message;
  },
};

function createBaseCreateResponse(): CreateResponse {
  return { id: "" };
}

export const CreateResponse = {
  encode(message: CreateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateResponse {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: CreateResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateResponse>, I>>(base?: I): CreateResponse {
    return CreateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateResponse>, I>>(object: I): CreateResponse {
    const message = createBaseCreateResponse();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseReadRequest(): ReadRequest {
  return { id: "" };
}

export const ReadRequest = {
  encode(message: ReadRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReadRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReadRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReadRequest {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: ReadRequest): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ReadRequest>, I>>(base?: I): ReadRequest {
    return ReadRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ReadRequest>, I>>(object: I): ReadRequest {
    const message = createBaseReadRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseReadResponse(): ReadResponse {
  return { id: "", result: undefined, expires_at: 0 };
}

export const ReadResponse = {
  encode(message: ReadResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.result !== undefined) {
      SimulationResult.encode(message.result, writer.uint32(18).fork()).ldelim();
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      writer.uint32(24).uint64(message.expires_at);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReadResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReadResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.result = SimulationResult.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.expires_at = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReadResponse {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      result: isSet(object.result) ? SimulationResult.fromJSON(object.result) : undefined,
      expires_at: isSet(object.expires_at) ? globalThis.Number(object.expires_at) : 0,
    };
  },

  toJSON(message: ReadResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    if (message.result !== undefined) {
      obj.result = SimulationResult.toJSON(message.result);
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      obj.expires_at = Math.round(message.expires_at);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ReadResponse>, I>>(base?: I): ReadResponse {
    return ReadResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ReadResponse>, I>>(object: I): ReadResponse {
    const message = createBaseReadResponse();
    message.id = object.id ?? "";
    message.result = (object.result !== undefined && object.result !== null)
      ? SimulationResult.fromPartial(object.result)
      : undefined;
    message.expires_at = object.expires_at ?? 0;
    return message;
  },
};

function createBaseUpdateRequest(): UpdateRequest {
  return { id: "", result: undefined, expires_at: 0, submitter: "" };
}

export const UpdateRequest = {
  encode(message: UpdateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.result !== undefined) {
      SimulationResult.encode(message.result, writer.uint32(18).fork()).ldelim();
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      writer.uint32(24).uint64(message.expires_at);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(34).string(message.submitter);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.result = SimulationResult.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.expires_at = longToNumber(reader.uint64() as Long);
          continue;
        case 4:
          if (tag !== 34) {
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

  fromJSON(object: any): UpdateRequest {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      result: isSet(object.result) ? SimulationResult.fromJSON(object.result) : undefined,
      expires_at: isSet(object.expires_at) ? globalThis.Number(object.expires_at) : 0,
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
    };
  },

  toJSON(message: UpdateRequest): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    if (message.result !== undefined) {
      obj.result = SimulationResult.toJSON(message.result);
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      obj.expires_at = Math.round(message.expires_at);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRequest>, I>>(base?: I): UpdateRequest {
    return UpdateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRequest>, I>>(object: I): UpdateRequest {
    const message = createBaseUpdateRequest();
    message.id = object.id ?? "";
    message.result = (object.result !== undefined && object.result !== null)
      ? SimulationResult.fromPartial(object.result)
      : undefined;
    message.expires_at = object.expires_at ?? 0;
    message.submitter = object.submitter ?? "";
    return message;
  },
};

function createBaseUpdateResponse(): UpdateResponse {
  return { id: "" };
}

export const UpdateResponse = {
  encode(message: UpdateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateResponse {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: UpdateResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateResponse>, I>>(base?: I): UpdateResponse {
    return UpdateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateResponse>, I>>(object: I): UpdateResponse {
    const message = createBaseUpdateResponse();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseSetTTLRequest(): SetTTLRequest {
  return { id: "", expires_at: 0 };
}

export const SetTTLRequest = {
  encode(message: SetTTLRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      writer.uint32(24).uint64(message.expires_at);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetTTLRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetTTLRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.expires_at = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetTTLRequest {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      expires_at: isSet(object.expires_at) ? globalThis.Number(object.expires_at) : 0,
    };
  },

  toJSON(message: SetTTLRequest): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    if (message.expires_at !== undefined && message.expires_at !== 0) {
      obj.expires_at = Math.round(message.expires_at);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetTTLRequest>, I>>(base?: I): SetTTLRequest {
    return SetTTLRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetTTLRequest>, I>>(object: I): SetTTLRequest {
    const message = createBaseSetTTLRequest();
    message.id = object.id ?? "";
    message.expires_at = object.expires_at ?? 0;
    return message;
  },
};

function createBaseSetTTLResponse(): SetTTLResponse {
  return { id: "" };
}

export const SetTTLResponse = {
  encode(message: SetTTLResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetTTLResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetTTLResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetTTLResponse {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: SetTTLResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetTTLResponse>, I>>(base?: I): SetTTLResponse {
    return SetTTLResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetTTLResponse>, I>>(object: I): SetTTLResponse {
    const message = createBaseSetTTLResponse();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseDeleteRequest(): DeleteRequest {
  return { id: "" };
}

export const DeleteRequest = {
  encode(message: DeleteRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteRequest {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: DeleteRequest): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteRequest>, I>>(base?: I): DeleteRequest {
    return DeleteRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteRequest>, I>>(object: I): DeleteRequest {
    const message = createBaseDeleteRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseDeleteResponse(): DeleteResponse {
  return { id: "" };
}

export const DeleteResponse = {
  encode(message: DeleteResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteResponse {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: DeleteResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteResponse>, I>>(base?: I): DeleteResponse {
    return DeleteResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteResponse>, I>>(object: I): DeleteResponse {
    const message = createBaseDeleteResponse();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseRandomRequest(): RandomRequest {
  return {};
}

export const RandomRequest = {
  encode(_: RandomRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RandomRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRandomRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): RandomRequest {
    return {};
  },

  toJSON(_: RandomRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<RandomRequest>, I>>(base?: I): RandomRequest {
    return RandomRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RandomRequest>, I>>(_: I): RandomRequest {
    const message = createBaseRandomRequest();
    return message;
  },
};

function createBaseRandomResponse(): RandomResponse {
  return { id: "" };
}

export const RandomResponse = {
  encode(message: RandomResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RandomResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRandomResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RandomResponse {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: RandomResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RandomResponse>, I>>(base?: I): RandomResponse {
    return RandomResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RandomResponse>, I>>(object: I): RandomResponse {
    const message = createBaseRandomResponse();
    message.id = object.id ?? "";
    return message;
  },
};

export interface ShareStore {
  Create(request: CreateRequest): Promise<CreateResponse>;
  Read(request: ReadRequest): Promise<ReadResponse>;
  Update(request: UpdateRequest): Promise<UpdateResponse>;
  SetTTL(request: SetTTLRequest): Promise<SetTTLResponse>;
  Delete(request: DeleteRequest): Promise<DeleteResponse>;
  Random(request: RandomRequest): Promise<RandomResponse>;
}

export const ShareStoreServiceName = "share.ShareStore";
export class ShareStoreClientImpl implements ShareStore {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || ShareStoreServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Read = this.Read.bind(this);
    this.Update = this.Update.bind(this);
    this.SetTTL = this.SetTTL.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Random = this.Random.bind(this);
  }
  Create(request: CreateRequest): Promise<CreateResponse> {
    const data = CreateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreateResponse.decode(_m0.Reader.create(data)));
  }

  Read(request: ReadRequest): Promise<ReadResponse> {
    const data = ReadRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Read", data);
    return promise.then((data) => ReadResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdateRequest): Promise<UpdateResponse> {
    const data = UpdateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateResponse.decode(_m0.Reader.create(data)));
  }

  SetTTL(request: SetTTLRequest): Promise<SetTTLResponse> {
    const data = SetTTLRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "SetTTL", data);
    return promise.then((data) => SetTTLResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteRequest): Promise<DeleteResponse> {
    const data = DeleteRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteResponse.decode(_m0.Reader.create(data)));
  }

  Random(request: RandomRequest): Promise<RandomResponse> {
    const data = RandomRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Random", data);
    return promise.then((data) => RandomResponse.decode(_m0.Reader.create(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(globalThis.Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
