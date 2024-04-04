/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Struct } from "../../google/protobuf/struct";
import { Character, Enemy } from "./sim";

export interface Sample {
  build_date?: string | undefined;
  sim_version?: string | undefined;
  modified?: boolean | undefined;
  config?: string | undefined;
  initial_character?: string | undefined;
  character_details?: Character[] | undefined;
  target_details?: Enemy[] | undefined;
  seed?: string | undefined;
  logs?: { [key: string]: any }[] | undefined;
}

function createBaseSample(): Sample {
  return {
    build_date: "",
    sim_version: undefined,
    modified: undefined,
    config: "",
    initial_character: "",
    character_details: [],
    target_details: [],
    seed: "",
    logs: [],
  };
}

export const Sample = {
  encode(message: Sample, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.build_date !== undefined && message.build_date !== "") {
      writer.uint32(18).string(message.build_date);
    }
    if (message.sim_version !== undefined) {
      writer.uint32(10).string(message.sim_version);
    }
    if (message.modified !== undefined) {
      writer.uint32(24).bool(message.modified);
    }
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(34).string(message.config);
    }
    if (message.initial_character !== undefined && message.initial_character !== "") {
      writer.uint32(42).string(message.initial_character);
    }
    if (message.character_details !== undefined && message.character_details.length !== 0) {
      for (const v of message.character_details) {
        Character.encode(v!, writer.uint32(50).fork()).ldelim();
      }
    }
    if (message.target_details !== undefined && message.target_details.length !== 0) {
      for (const v of message.target_details) {
        Enemy.encode(v!, writer.uint32(58).fork()).ldelim();
      }
    }
    if (message.seed !== undefined && message.seed !== "") {
      writer.uint32(66).string(message.seed);
    }
    if (message.logs !== undefined && message.logs.length !== 0) {
      for (const v of message.logs) {
        Struct.encode(Struct.wrap(v!), writer.uint32(74).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Sample {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSample();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag !== 18) {
            break;
          }

          message.build_date = reader.string();
          continue;
        case 1:
          if (tag !== 10) {
            break;
          }

          message.sim_version = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.modified = reader.bool();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.config = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.initial_character = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.character_details!.push(Character.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.target_details!.push(Enemy.decode(reader, reader.uint32()));
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.seed = reader.string();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.logs!.push(Struct.unwrap(Struct.decode(reader, reader.uint32())));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Sample {
    return {
      build_date: isSet(object.build_date) ? globalThis.String(object.build_date) : "",
      sim_version: isSet(object.sim_version) ? globalThis.String(object.sim_version) : undefined,
      modified: isSet(object.modified) ? globalThis.Boolean(object.modified) : undefined,
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      initial_character: isSet(object.initial_character) ? globalThis.String(object.initial_character) : "",
      character_details: globalThis.Array.isArray(object?.character_details)
        ? object.character_details.map((e: any) => Character.fromJSON(e))
        : [],
      target_details: globalThis.Array.isArray(object?.target_details)
        ? object.target_details.map((e: any) => Enemy.fromJSON(e))
        : [],
      seed: isSet(object.seed) ? globalThis.String(object.seed) : "",
      logs: globalThis.Array.isArray(object?.logs) ? [...object.logs] : [],
    };
  },

  toJSON(message: Sample): unknown {
    const obj: any = {};
    if (message.build_date !== undefined && message.build_date !== "") {
      obj.build_date = message.build_date;
    }
    if (message.sim_version !== undefined) {
      obj.sim_version = message.sim_version;
    }
    if (message.modified !== undefined) {
      obj.modified = message.modified;
    }
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.initial_character !== undefined && message.initial_character !== "") {
      obj.initial_character = message.initial_character;
    }
    if (message.character_details?.length) {
      obj.character_details = message.character_details.map((e) => Character.toJSON(e));
    }
    if (message.target_details?.length) {
      obj.target_details = message.target_details.map((e) => Enemy.toJSON(e));
    }
    if (message.seed !== undefined && message.seed !== "") {
      obj.seed = message.seed;
    }
    if (message.logs?.length) {
      obj.logs = message.logs;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Sample>, I>>(base?: I): Sample {
    return Sample.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Sample>, I>>(object: I): Sample {
    const message = createBaseSample();
    message.build_date = object.build_date ?? "";
    message.sim_version = object.sim_version ?? undefined;
    message.modified = object.modified ?? undefined;
    message.config = object.config ?? "";
    message.initial_character = object.initial_character ?? "";
    message.character_details = object.character_details?.map((e) => Character.fromPartial(e)) || [];
    message.target_details = object.target_details?.map((e) => Enemy.fromPartial(e)) || [];
    message.seed = object.seed ?? "";
    message.logs = object.logs?.map((e) => e) || [];
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
