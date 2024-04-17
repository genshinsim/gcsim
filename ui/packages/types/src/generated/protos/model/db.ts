/* eslint-disable */
import _m0 from "protobufjs/minimal";

export interface DBStatus {
  db_total_count?: number | undefined;
  compute_count?: number | undefined;
}

function createBaseDBStatus(): DBStatus {
  return { db_total_count: 0, compute_count: 0 };
}

export const DBStatus = {
  encode(message: DBStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.db_total_count !== undefined && message.db_total_count !== 0) {
      writer.uint32(8).int32(message.db_total_count);
    }
    if (message.compute_count !== undefined && message.compute_count !== 0) {
      writer.uint32(16).int32(message.compute_count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DBStatus {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDBStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.db_total_count = reader.int32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.compute_count = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DBStatus {
    return {
      db_total_count: isSet(object.db_total_count) ? globalThis.Number(object.db_total_count) : 0,
      compute_count: isSet(object.compute_count) ? globalThis.Number(object.compute_count) : 0,
    };
  },

  toJSON(message: DBStatus): unknown {
    const obj: any = {};
    if (message.db_total_count !== undefined && message.db_total_count !== 0) {
      obj.db_total_count = Math.round(message.db_total_count);
    }
    if (message.compute_count !== undefined && message.compute_count !== 0) {
      obj.compute_count = Math.round(message.compute_count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DBStatus>, I>>(base?: I): DBStatus {
    return DBStatus.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DBStatus>, I>>(object: I): DBStatus {
    const message = createBaseDBStatus();
    message.db_total_count = object.db_total_count ?? 0;
    message.compute_count = object.compute_count ?? 0;
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
