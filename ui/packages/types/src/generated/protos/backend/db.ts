/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Struct } from "../../google/protobuf/struct";
import { DBTag, dBTagFromJSON, dBTagToJSON, SimMode, simModeFromJSON, simModeToJSON } from "../model/enums";
import { DescriptiveStats, SimulationResult } from "../model/result";
import { Character } from "../model/sim";

export interface Entry {
  /** basic info */
  _id?: string | undefined;
  create_date?:
    | number
    | undefined;
  /** key fields */
  config?: string | undefined;
  description?: string | undefined;
  submitter?:
    | string
    | undefined;
  /** db tagging data */
  accepted_tags?: DBTag[] | undefined;
  rejected_tags?:
    | DBTag[]
    | undefined;
  /** string upgrade_failed = 20 */
  is_db_valid?:
    | boolean
    | undefined;
  /** these fields are updated on every rerun */
  share_key?: string | undefined;
  last_update?: number | undefined;
  hash?: string | undefined;
  summary?: EntrySummary | undefined;
}

export interface EntrySummary {
  sim_duration?: DescriptiveStats | undefined;
  mode?:
    | SimMode
    | undefined;
  /** indexing data */
  total_damage?: DescriptiveStats | undefined;
  char_names?: string[] | undefined;
  target_count?: number | undefined;
  mean_dps_per_target?:
    | number
    | undefined;
  /** detailed results */
  team?: Character[] | undefined;
  dps_by_target?: { [key: string]: DescriptiveStats } | undefined;
}

export interface EntrySummary_DpsByTargetEntry {
  key: string;
  value?: DescriptiveStats | undefined;
}

export interface Entries {
  data?: Entry[] | undefined;
}

export interface QueryOpt {
  query?: { [key: string]: any } | undefined;
  sort?: { [key: string]: any } | undefined;
  project?: { [key: string]: any } | undefined;
  skip?: number | undefined;
  limit?: number | undefined;
}

export interface ComputeWork {
  _id?: string | undefined;
  config?: string | undefined;
  iterations?: number | undefined;
}

export interface GetRequest {
  query?: QueryOpt | undefined;
}

export interface GetResponse {
  data?: Entries | undefined;
}

export interface GetOneRequest {
  _id?: string | undefined;
}

export interface GetOneResponse {
  data?: Entry | undefined;
}

export interface GetPendingRequest {
  tag?: DBTag | undefined;
  query?: QueryOpt | undefined;
}

export interface GetPendingResponse {
  data?: Entries | undefined;
}

export interface GetBySubmitterRequest {
  submitter?: string | undefined;
  query?: QueryOpt | undefined;
}

export interface GetBySubmitterResponse {
  data?: Entries | undefined;
}

export interface GetAllRequest {
  query?: QueryOpt | undefined;
}

export interface GetAllResponse {
  data?: Entries | undefined;
}

export interface ApproveTagRequest {
  id?: string | undefined;
  tag?: DBTag | undefined;
}

export interface ApproveTagResponse {
  id?: string | undefined;
}

export interface RejectTagRequest {
  id?: string | undefined;
  tag?: DBTag | undefined;
}

export interface RejectTagResponse {
  id?: string | undefined;
}

export interface RejectTagAllUnapprovedRequest {
  tag?: DBTag | undefined;
}

export interface RejectTagAllUnapprovedResponse {
  count?: number | undefined;
}

export interface SubmitRequest {
  config?:
    | string
    | undefined;
  /** submitter discord id */
  submitter?: string | undefined;
  description?: string | undefined;
}

export interface SubmitResponse {
  _id?: string | undefined;
}

export interface DeletePendingRequest {
  _id?: string | undefined;
  sender?: string | undefined;
}

export interface DeletePendingResponse {
  _id?: string | undefined;
}

export interface GetWorkRequest {
}

export interface GetWorkResponse {
  data?: ComputeWork[] | undefined;
}

export interface RejectWorkRequest {
  id?: string | undefined;
  reason?: string | undefined;
  hash?: string | undefined;
}

export interface RejectWorkResponse {
}

export interface CompleteWorkRequest {
  _id?: string | undefined;
  result?: SimulationResult | undefined;
}

export interface CompleteWorkResponse {
  _id?: string | undefined;
}

export interface WorkStatusRequest {
}

export interface WorkStatusResponse {
  todo_count?: number | undefined;
  total_count?: number | undefined;
}

export interface ReplaceConfigRequest {
  _id?: string | undefined;
  config?: string | undefined;
  source_tag?: DBTag | undefined;
}

export interface ReplaceConfigResponse {
  _id?: string | undefined;
}

export interface ReplaceDescRequest {
  _id?: string | undefined;
  desc?: string | undefined;
  source_tag?: DBTag | undefined;
}

export interface ReplaceDescResponse {
  _id?: string | undefined;
}

function createBaseEntry(): Entry {
  return {
    _id: "",
    create_date: 0,
    config: "",
    description: "",
    submitter: "",
    accepted_tags: [],
    rejected_tags: [],
    is_db_valid: false,
    share_key: "",
    last_update: 0,
    hash: "",
    summary: undefined,
  };
}

export const Entry = {
  encode(message: Entry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    if (message.create_date !== undefined && message.create_date !== 0) {
      writer.uint32(16).uint64(message.create_date);
    }
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(26).string(message.config);
    }
    if (message.description !== undefined && message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(42).string(message.submitter);
    }
    if (message.accepted_tags !== undefined && message.accepted_tags.length !== 0) {
      writer.uint32(50).fork();
      for (const v of message.accepted_tags) {
        writer.int32(v);
      }
      writer.ldelim();
    }
    if (message.rejected_tags !== undefined && message.rejected_tags.length !== 0) {
      writer.uint32(58).fork();
      for (const v of message.rejected_tags) {
        writer.int32(v);
      }
      writer.ldelim();
    }
    if (message.is_db_valid !== undefined && message.is_db_valid !== false) {
      writer.uint32(64).bool(message.is_db_valid);
    }
    if (message.share_key !== undefined && message.share_key !== "") {
      writer.uint32(74).string(message.share_key);
    }
    if (message.last_update !== undefined && message.last_update !== 0) {
      writer.uint32(80).uint64(message.last_update);
    }
    if (message.hash !== undefined && message.hash !== "") {
      writer.uint32(90).string(message.hash);
    }
    if (message.summary !== undefined) {
      EntrySummary.encode(message.summary, writer.uint32(98).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Entry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntry();
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
          if (tag !== 16) {
            break;
          }

          message.create_date = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.config = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.description = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.submitter = reader.string();
          continue;
        case 6:
          if (tag === 48) {
            message.accepted_tags!.push(reader.int32() as any);

            continue;
          }

          if (tag === 50) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.accepted_tags!.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 7:
          if (tag === 56) {
            message.rejected_tags!.push(reader.int32() as any);

            continue;
          }

          if (tag === 58) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.rejected_tags!.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.is_db_valid = reader.bool();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.share_key = reader.string();
          continue;
        case 10:
          if (tag !== 80) {
            break;
          }

          message.last_update = longToNumber(reader.uint64() as Long);
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.hash = reader.string();
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.summary = EntrySummary.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Entry {
    return {
      _id: isSet(object._id) ? globalThis.String(object._id) : "",
      create_date: isSet(object.create_date) ? globalThis.Number(object.create_date) : 0,
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      description: isSet(object.description) ? globalThis.String(object.description) : "",
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
      accepted_tags: globalThis.Array.isArray(object?.accepted_tags)
        ? object.accepted_tags.map((e: any) => dBTagFromJSON(e))
        : [],
      rejected_tags: globalThis.Array.isArray(object?.rejected_tags)
        ? object.rejected_tags.map((e: any) => dBTagFromJSON(e))
        : [],
      is_db_valid: isSet(object.is_db_valid) ? globalThis.Boolean(object.is_db_valid) : false,
      share_key: isSet(object.share_key) ? globalThis.String(object.share_key) : "",
      last_update: isSet(object.last_update) ? globalThis.Number(object.last_update) : 0,
      hash: isSet(object.hash) ? globalThis.String(object.hash) : "",
      summary: isSet(object.summary) ? EntrySummary.fromJSON(object.summary) : undefined,
    };
  },

  toJSON(message: Entry): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    if (message.create_date !== undefined && message.create_date !== 0) {
      obj.create_date = Math.round(message.create_date);
    }
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.description !== undefined && message.description !== "") {
      obj.description = message.description;
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    if (message.accepted_tags?.length) {
      obj.accepted_tags = message.accepted_tags.map((e) => dBTagToJSON(e));
    }
    if (message.rejected_tags?.length) {
      obj.rejected_tags = message.rejected_tags.map((e) => dBTagToJSON(e));
    }
    if (message.is_db_valid !== undefined && message.is_db_valid !== false) {
      obj.is_db_valid = message.is_db_valid;
    }
    if (message.share_key !== undefined && message.share_key !== "") {
      obj.share_key = message.share_key;
    }
    if (message.last_update !== undefined && message.last_update !== 0) {
      obj.last_update = Math.round(message.last_update);
    }
    if (message.hash !== undefined && message.hash !== "") {
      obj.hash = message.hash;
    }
    if (message.summary !== undefined) {
      obj.summary = EntrySummary.toJSON(message.summary);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Entry>, I>>(base?: I): Entry {
    return Entry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Entry>, I>>(object: I): Entry {
    const message = createBaseEntry();
    message._id = object._id ?? "";
    message.create_date = object.create_date ?? 0;
    message.config = object.config ?? "";
    message.description = object.description ?? "";
    message.submitter = object.submitter ?? "";
    message.accepted_tags = object.accepted_tags?.map((e) => e) || [];
    message.rejected_tags = object.rejected_tags?.map((e) => e) || [];
    message.is_db_valid = object.is_db_valid ?? false;
    message.share_key = object.share_key ?? "";
    message.last_update = object.last_update ?? 0;
    message.hash = object.hash ?? "";
    message.summary = (object.summary !== undefined && object.summary !== null)
      ? EntrySummary.fromPartial(object.summary)
      : undefined;
    return message;
  },
};

function createBaseEntrySummary(): EntrySummary {
  return {
    sim_duration: undefined,
    mode: 0,
    total_damage: undefined,
    char_names: [],
    target_count: 0,
    mean_dps_per_target: 0,
    team: [],
    dps_by_target: {},
  };
}

export const EntrySummary = {
  encode(message: EntrySummary, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sim_duration !== undefined) {
      DescriptiveStats.encode(message.sim_duration, writer.uint32(10).fork()).ldelim();
    }
    if (message.mode !== undefined && message.mode !== 0) {
      writer.uint32(16).int32(message.mode);
    }
    if (message.total_damage !== undefined) {
      DescriptiveStats.encode(message.total_damage, writer.uint32(26).fork()).ldelim();
    }
    if (message.char_names !== undefined && message.char_names.length !== 0) {
      for (const v of message.char_names) {
        writer.uint32(34).string(v!);
      }
    }
    if (message.target_count !== undefined && message.target_count !== 0) {
      writer.uint32(40).int32(message.target_count);
    }
    if (message.mean_dps_per_target !== undefined && message.mean_dps_per_target !== 0) {
      writer.uint32(49).double(message.mean_dps_per_target);
    }
    if (message.team !== undefined && message.team.length !== 0) {
      for (const v of message.team) {
        Character.encode(v!, writer.uint32(58).fork()).ldelim();
      }
    }
    Object.entries(message.dps_by_target || {}).forEach(([key, value]) => {
      EntrySummary_DpsByTargetEntry.encode({ key: key as any, value }, writer.uint32(66).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EntrySummary {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntrySummary();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.sim_duration = DescriptiveStats.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.mode = reader.int32() as any;
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.total_damage = DescriptiveStats.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.char_names!.push(reader.string());
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.target_count = reader.int32();
          continue;
        case 6:
          if (tag !== 49) {
            break;
          }

          message.mean_dps_per_target = reader.double();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.team!.push(Character.decode(reader, reader.uint32()));
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          const entry8 = EntrySummary_DpsByTargetEntry.decode(reader, reader.uint32());
          if (entry8.value !== undefined) {
            message.dps_by_target![entry8.key] = entry8.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EntrySummary {
    return {
      sim_duration: isSet(object.sim_duration) ? DescriptiveStats.fromJSON(object.sim_duration) : undefined,
      mode: isSet(object.mode) ? simModeFromJSON(object.mode) : 0,
      total_damage: isSet(object.total_damage) ? DescriptiveStats.fromJSON(object.total_damage) : undefined,
      char_names: globalThis.Array.isArray(object?.char_names)
        ? object.char_names.map((e: any) => globalThis.String(e))
        : [],
      target_count: isSet(object.target_count) ? globalThis.Number(object.target_count) : 0,
      mean_dps_per_target: isSet(object.mean_dps_per_target) ? globalThis.Number(object.mean_dps_per_target) : 0,
      team: globalThis.Array.isArray(object?.team) ? object.team.map((e: any) => Character.fromJSON(e)) : [],
      dps_by_target: isObject(object.dps_by_target)
        ? Object.entries(object.dps_by_target).reduce<{ [key: string]: DescriptiveStats }>((acc, [key, value]) => {
          acc[key] = DescriptiveStats.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: EntrySummary): unknown {
    const obj: any = {};
    if (message.sim_duration !== undefined) {
      obj.sim_duration = DescriptiveStats.toJSON(message.sim_duration);
    }
    if (message.mode !== undefined && message.mode !== 0) {
      obj.mode = simModeToJSON(message.mode);
    }
    if (message.total_damage !== undefined) {
      obj.total_damage = DescriptiveStats.toJSON(message.total_damage);
    }
    if (message.char_names?.length) {
      obj.char_names = message.char_names;
    }
    if (message.target_count !== undefined && message.target_count !== 0) {
      obj.target_count = Math.round(message.target_count);
    }
    if (message.mean_dps_per_target !== undefined && message.mean_dps_per_target !== 0) {
      obj.mean_dps_per_target = message.mean_dps_per_target;
    }
    if (message.team?.length) {
      obj.team = message.team.map((e) => Character.toJSON(e));
    }
    if (message.dps_by_target) {
      const entries = Object.entries(message.dps_by_target);
      if (entries.length > 0) {
        obj.dps_by_target = {};
        entries.forEach(([k, v]) => {
          obj.dps_by_target[k] = DescriptiveStats.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EntrySummary>, I>>(base?: I): EntrySummary {
    return EntrySummary.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EntrySummary>, I>>(object: I): EntrySummary {
    const message = createBaseEntrySummary();
    message.sim_duration = (object.sim_duration !== undefined && object.sim_duration !== null)
      ? DescriptiveStats.fromPartial(object.sim_duration)
      : undefined;
    message.mode = object.mode ?? 0;
    message.total_damage = (object.total_damage !== undefined && object.total_damage !== null)
      ? DescriptiveStats.fromPartial(object.total_damage)
      : undefined;
    message.char_names = object.char_names?.map((e) => e) || [];
    message.target_count = object.target_count ?? 0;
    message.mean_dps_per_target = object.mean_dps_per_target ?? 0;
    message.team = object.team?.map((e) => Character.fromPartial(e)) || [];
    message.dps_by_target = Object.entries(object.dps_by_target ?? {}).reduce<{ [key: string]: DescriptiveStats }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = DescriptiveStats.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseEntrySummary_DpsByTargetEntry(): EntrySummary_DpsByTargetEntry {
  return { key: "", value: undefined };
}

export const EntrySummary_DpsByTargetEntry = {
  encode(message: EntrySummary_DpsByTargetEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      DescriptiveStats.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EntrySummary_DpsByTargetEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntrySummary_DpsByTargetEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = DescriptiveStats.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EntrySummary_DpsByTargetEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? DescriptiveStats.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: EntrySummary_DpsByTargetEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = DescriptiveStats.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EntrySummary_DpsByTargetEntry>, I>>(base?: I): EntrySummary_DpsByTargetEntry {
    return EntrySummary_DpsByTargetEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EntrySummary_DpsByTargetEntry>, I>>(
    object: I,
  ): EntrySummary_DpsByTargetEntry {
    const message = createBaseEntrySummary_DpsByTargetEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? DescriptiveStats.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseEntries(): Entries {
  return { data: [] };
}

export const Entries = {
  encode(message: Entries, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== undefined && message.data.length !== 0) {
      for (const v of message.data) {
        Entry.encode(v!, writer.uint32(10).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Entries {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntries();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data!.push(Entry.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Entries {
    return { data: globalThis.Array.isArray(object?.data) ? object.data.map((e: any) => Entry.fromJSON(e)) : [] };
  },

  toJSON(message: Entries): unknown {
    const obj: any = {};
    if (message.data?.length) {
      obj.data = message.data.map((e) => Entry.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Entries>, I>>(base?: I): Entries {
    return Entries.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Entries>, I>>(object: I): Entries {
    const message = createBaseEntries();
    message.data = object.data?.map((e) => Entry.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQueryOpt(): QueryOpt {
  return { query: undefined, sort: undefined, project: undefined, skip: 0, limit: 0 };
}

export const QueryOpt = {
  encode(message: QueryOpt, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      Struct.encode(Struct.wrap(message.query), writer.uint32(10).fork()).ldelim();
    }
    if (message.sort !== undefined) {
      Struct.encode(Struct.wrap(message.sort), writer.uint32(18).fork()).ldelim();
    }
    if (message.project !== undefined) {
      Struct.encode(Struct.wrap(message.project), writer.uint32(26).fork()).ldelim();
    }
    if (message.skip !== undefined && message.skip !== 0) {
      writer.uint32(32).int64(message.skip);
    }
    if (message.limit !== undefined && message.limit !== 0) {
      writer.uint32(40).int64(message.limit);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryOpt {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryOpt();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.query = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.sort = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.project = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.skip = longToNumber(reader.int64() as Long);
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.limit = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): QueryOpt {
    return {
      query: isObject(object.query) ? object.query : undefined,
      sort: isObject(object.sort) ? object.sort : undefined,
      project: isObject(object.project) ? object.project : undefined,
      skip: isSet(object.skip) ? globalThis.Number(object.skip) : 0,
      limit: isSet(object.limit) ? globalThis.Number(object.limit) : 0,
    };
  },

  toJSON(message: QueryOpt): unknown {
    const obj: any = {};
    if (message.query !== undefined) {
      obj.query = message.query;
    }
    if (message.sort !== undefined) {
      obj.sort = message.sort;
    }
    if (message.project !== undefined) {
      obj.project = message.project;
    }
    if (message.skip !== undefined && message.skip !== 0) {
      obj.skip = Math.round(message.skip);
    }
    if (message.limit !== undefined && message.limit !== 0) {
      obj.limit = Math.round(message.limit);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<QueryOpt>, I>>(base?: I): QueryOpt {
    return QueryOpt.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<QueryOpt>, I>>(object: I): QueryOpt {
    const message = createBaseQueryOpt();
    message.query = object.query ?? undefined;
    message.sort = object.sort ?? undefined;
    message.project = object.project ?? undefined;
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseComputeWork(): ComputeWork {
  return { _id: "", config: "", iterations: 0 };
}

export const ComputeWork = {
  encode(message: ComputeWork, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(18).string(message.config);
    }
    if (message.iterations !== undefined && message.iterations !== 0) {
      writer.uint32(24).int32(message.iterations);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ComputeWork {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseComputeWork();
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

          message.config = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.iterations = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ComputeWork {
    return {
      _id: isSet(object._id) ? globalThis.String(object._id) : "",
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      iterations: isSet(object.iterations) ? globalThis.Number(object.iterations) : 0,
    };
  },

  toJSON(message: ComputeWork): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.iterations !== undefined && message.iterations !== 0) {
      obj.iterations = Math.round(message.iterations);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ComputeWork>, I>>(base?: I): ComputeWork {
    return ComputeWork.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ComputeWork>, I>>(object: I): ComputeWork {
    const message = createBaseComputeWork();
    message._id = object._id ?? "";
    message.config = object.config ?? "";
    message.iterations = object.iterations ?? 0;
    return message;
  },
};

function createBaseGetRequest(): GetRequest {
  return { query: undefined };
}

export const GetRequest = {
  encode(message: GetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      QueryOpt.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.query = QueryOpt.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRequest {
    return { query: isSet(object.query) ? QueryOpt.fromJSON(object.query) : undefined };
  },

  toJSON(message: GetRequest): unknown {
    const obj: any = {};
    if (message.query !== undefined) {
      obj.query = QueryOpt.toJSON(message.query);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRequest>, I>>(base?: I): GetRequest {
    return GetRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRequest>, I>>(object: I): GetRequest {
    const message = createBaseGetRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? QueryOpt.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseGetResponse(): GetResponse {
  return { data: undefined };
}

export const GetResponse = {
  encode(message: GetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== undefined) {
      Entries.encode(message.data, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = Entries.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetResponse {
    return { data: isSet(object.data) ? Entries.fromJSON(object.data) : undefined };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    if (message.data !== undefined) {
      obj.data = Entries.toJSON(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetResponse>, I>>(base?: I): GetResponse {
    return GetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(object: I): GetResponse {
    const message = createBaseGetResponse();
    message.data = (object.data !== undefined && object.data !== null) ? Entries.fromPartial(object.data) : undefined;
    return message;
  },
};

function createBaseGetOneRequest(): GetOneRequest {
  return { _id: "" };
}

export const GetOneRequest = {
  encode(message: GetOneRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOneRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOneRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message._id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetOneRequest {
    return { _id: isSet(object._id) ? globalThis.String(object._id) : "" };
  },

  toJSON(message: GetOneRequest): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetOneRequest>, I>>(base?: I): GetOneRequest {
    return GetOneRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetOneRequest>, I>>(object: I): GetOneRequest {
    const message = createBaseGetOneRequest();
    message._id = object._id ?? "";
    return message;
  },
};

function createBaseGetOneResponse(): GetOneResponse {
  return { data: undefined };
}

export const GetOneResponse = {
  encode(message: GetOneResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== undefined) {
      Entry.encode(message.data, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOneResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOneResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = Entry.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetOneResponse {
    return { data: isSet(object.data) ? Entry.fromJSON(object.data) : undefined };
  },

  toJSON(message: GetOneResponse): unknown {
    const obj: any = {};
    if (message.data !== undefined) {
      obj.data = Entry.toJSON(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetOneResponse>, I>>(base?: I): GetOneResponse {
    return GetOneResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetOneResponse>, I>>(object: I): GetOneResponse {
    const message = createBaseGetOneResponse();
    message.data = (object.data !== undefined && object.data !== null) ? Entry.fromPartial(object.data) : undefined;
    return message;
  },
};

function createBaseGetPendingRequest(): GetPendingRequest {
  return { tag: 0, query: undefined };
}

export const GetPendingRequest = {
  encode(message: GetPendingRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tag !== undefined && message.tag !== 0) {
      writer.uint32(8).int32(message.tag);
    }
    if (message.query !== undefined) {
      QueryOpt.encode(message.query, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPendingRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPendingRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.tag = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.query = QueryOpt.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetPendingRequest {
    return {
      tag: isSet(object.tag) ? dBTagFromJSON(object.tag) : 0,
      query: isSet(object.query) ? QueryOpt.fromJSON(object.query) : undefined,
    };
  },

  toJSON(message: GetPendingRequest): unknown {
    const obj: any = {};
    if (message.tag !== undefined && message.tag !== 0) {
      obj.tag = dBTagToJSON(message.tag);
    }
    if (message.query !== undefined) {
      obj.query = QueryOpt.toJSON(message.query);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetPendingRequest>, I>>(base?: I): GetPendingRequest {
    return GetPendingRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetPendingRequest>, I>>(object: I): GetPendingRequest {
    const message = createBaseGetPendingRequest();
    message.tag = object.tag ?? 0;
    message.query = (object.query !== undefined && object.query !== null)
      ? QueryOpt.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseGetPendingResponse(): GetPendingResponse {
  return { data: undefined };
}

export const GetPendingResponse = {
  encode(message: GetPendingResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== undefined) {
      Entries.encode(message.data, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPendingResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPendingResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = Entries.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetPendingResponse {
    return { data: isSet(object.data) ? Entries.fromJSON(object.data) : undefined };
  },

  toJSON(message: GetPendingResponse): unknown {
    const obj: any = {};
    if (message.data !== undefined) {
      obj.data = Entries.toJSON(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetPendingResponse>, I>>(base?: I): GetPendingResponse {
    return GetPendingResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetPendingResponse>, I>>(object: I): GetPendingResponse {
    const message = createBaseGetPendingResponse();
    message.data = (object.data !== undefined && object.data !== null) ? Entries.fromPartial(object.data) : undefined;
    return message;
  },
};

function createBaseGetBySubmitterRequest(): GetBySubmitterRequest {
  return { submitter: "", query: undefined };
}

export const GetBySubmitterRequest = {
  encode(message: GetBySubmitterRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(10).string(message.submitter);
    }
    if (message.query !== undefined) {
      QueryOpt.encode(message.query, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBySubmitterRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBySubmitterRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.submitter = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.query = QueryOpt.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetBySubmitterRequest {
    return {
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
      query: isSet(object.query) ? QueryOpt.fromJSON(object.query) : undefined,
    };
  },

  toJSON(message: GetBySubmitterRequest): unknown {
    const obj: any = {};
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    if (message.query !== undefined) {
      obj.query = QueryOpt.toJSON(message.query);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetBySubmitterRequest>, I>>(base?: I): GetBySubmitterRequest {
    return GetBySubmitterRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetBySubmitterRequest>, I>>(object: I): GetBySubmitterRequest {
    const message = createBaseGetBySubmitterRequest();
    message.submitter = object.submitter ?? "";
    message.query = (object.query !== undefined && object.query !== null)
      ? QueryOpt.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseGetBySubmitterResponse(): GetBySubmitterResponse {
  return { data: undefined };
}

export const GetBySubmitterResponse = {
  encode(message: GetBySubmitterResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== undefined) {
      Entries.encode(message.data, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBySubmitterResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBySubmitterResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = Entries.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetBySubmitterResponse {
    return { data: isSet(object.data) ? Entries.fromJSON(object.data) : undefined };
  },

  toJSON(message: GetBySubmitterResponse): unknown {
    const obj: any = {};
    if (message.data !== undefined) {
      obj.data = Entries.toJSON(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetBySubmitterResponse>, I>>(base?: I): GetBySubmitterResponse {
    return GetBySubmitterResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetBySubmitterResponse>, I>>(object: I): GetBySubmitterResponse {
    const message = createBaseGetBySubmitterResponse();
    message.data = (object.data !== undefined && object.data !== null) ? Entries.fromPartial(object.data) : undefined;
    return message;
  },
};

function createBaseGetAllRequest(): GetAllRequest {
  return { query: undefined };
}

export const GetAllRequest = {
  encode(message: GetAllRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      QueryOpt.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAllRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.query = QueryOpt.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAllRequest {
    return { query: isSet(object.query) ? QueryOpt.fromJSON(object.query) : undefined };
  },

  toJSON(message: GetAllRequest): unknown {
    const obj: any = {};
    if (message.query !== undefined) {
      obj.query = QueryOpt.toJSON(message.query);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAllRequest>, I>>(base?: I): GetAllRequest {
    return GetAllRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAllRequest>, I>>(object: I): GetAllRequest {
    const message = createBaseGetAllRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? QueryOpt.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseGetAllResponse(): GetAllResponse {
  return { data: undefined };
}

export const GetAllResponse = {
  encode(message: GetAllResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== undefined) {
      Entries.encode(message.data, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAllResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = Entries.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAllResponse {
    return { data: isSet(object.data) ? Entries.fromJSON(object.data) : undefined };
  },

  toJSON(message: GetAllResponse): unknown {
    const obj: any = {};
    if (message.data !== undefined) {
      obj.data = Entries.toJSON(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAllResponse>, I>>(base?: I): GetAllResponse {
    return GetAllResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAllResponse>, I>>(object: I): GetAllResponse {
    const message = createBaseGetAllResponse();
    message.data = (object.data !== undefined && object.data !== null) ? Entries.fromPartial(object.data) : undefined;
    return message;
  },
};

function createBaseApproveTagRequest(): ApproveTagRequest {
  return { id: "", tag: 0 };
}

export const ApproveTagRequest = {
  encode(message: ApproveTagRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.tag !== undefined && message.tag !== 0) {
      writer.uint32(16).int32(message.tag);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApproveTagRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApproveTagRequest();
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
          if (tag !== 16) {
            break;
          }

          message.tag = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ApproveTagRequest {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      tag: isSet(object.tag) ? dBTagFromJSON(object.tag) : 0,
    };
  },

  toJSON(message: ApproveTagRequest): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    if (message.tag !== undefined && message.tag !== 0) {
      obj.tag = dBTagToJSON(message.tag);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ApproveTagRequest>, I>>(base?: I): ApproveTagRequest {
    return ApproveTagRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ApproveTagRequest>, I>>(object: I): ApproveTagRequest {
    const message = createBaseApproveTagRequest();
    message.id = object.id ?? "";
    message.tag = object.tag ?? 0;
    return message;
  },
};

function createBaseApproveTagResponse(): ApproveTagResponse {
  return { id: "" };
}

export const ApproveTagResponse = {
  encode(message: ApproveTagResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApproveTagResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApproveTagResponse();
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

  fromJSON(object: any): ApproveTagResponse {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: ApproveTagResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ApproveTagResponse>, I>>(base?: I): ApproveTagResponse {
    return ApproveTagResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ApproveTagResponse>, I>>(object: I): ApproveTagResponse {
    const message = createBaseApproveTagResponse();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseRejectTagRequest(): RejectTagRequest {
  return { id: "", tag: 0 };
}

export const RejectTagRequest = {
  encode(message: RejectTagRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.tag !== undefined && message.tag !== 0) {
      writer.uint32(16).int32(message.tag);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RejectTagRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRejectTagRequest();
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
          if (tag !== 16) {
            break;
          }

          message.tag = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RejectTagRequest {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      tag: isSet(object.tag) ? dBTagFromJSON(object.tag) : 0,
    };
  },

  toJSON(message: RejectTagRequest): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    if (message.tag !== undefined && message.tag !== 0) {
      obj.tag = dBTagToJSON(message.tag);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RejectTagRequest>, I>>(base?: I): RejectTagRequest {
    return RejectTagRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RejectTagRequest>, I>>(object: I): RejectTagRequest {
    const message = createBaseRejectTagRequest();
    message.id = object.id ?? "";
    message.tag = object.tag ?? 0;
    return message;
  },
};

function createBaseRejectTagResponse(): RejectTagResponse {
  return { id: "" };
}

export const RejectTagResponse = {
  encode(message: RejectTagResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RejectTagResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRejectTagResponse();
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

  fromJSON(object: any): RejectTagResponse {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: RejectTagResponse): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RejectTagResponse>, I>>(base?: I): RejectTagResponse {
    return RejectTagResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RejectTagResponse>, I>>(object: I): RejectTagResponse {
    const message = createBaseRejectTagResponse();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseRejectTagAllUnapprovedRequest(): RejectTagAllUnapprovedRequest {
  return { tag: 0 };
}

export const RejectTagAllUnapprovedRequest = {
  encode(message: RejectTagAllUnapprovedRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tag !== undefined && message.tag !== 0) {
      writer.uint32(8).int32(message.tag);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RejectTagAllUnapprovedRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRejectTagAllUnapprovedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.tag = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RejectTagAllUnapprovedRequest {
    return { tag: isSet(object.tag) ? dBTagFromJSON(object.tag) : 0 };
  },

  toJSON(message: RejectTagAllUnapprovedRequest): unknown {
    const obj: any = {};
    if (message.tag !== undefined && message.tag !== 0) {
      obj.tag = dBTagToJSON(message.tag);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RejectTagAllUnapprovedRequest>, I>>(base?: I): RejectTagAllUnapprovedRequest {
    return RejectTagAllUnapprovedRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RejectTagAllUnapprovedRequest>, I>>(
    object: I,
  ): RejectTagAllUnapprovedRequest {
    const message = createBaseRejectTagAllUnapprovedRequest();
    message.tag = object.tag ?? 0;
    return message;
  },
};

function createBaseRejectTagAllUnapprovedResponse(): RejectTagAllUnapprovedResponse {
  return { count: 0 };
}

export const RejectTagAllUnapprovedResponse = {
  encode(message: RejectTagAllUnapprovedResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== undefined && message.count !== 0) {
      writer.uint32(8).int64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RejectTagAllUnapprovedResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRejectTagAllUnapprovedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.count = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RejectTagAllUnapprovedResponse {
    return { count: isSet(object.count) ? globalThis.Number(object.count) : 0 };
  },

  toJSON(message: RejectTagAllUnapprovedResponse): unknown {
    const obj: any = {};
    if (message.count !== undefined && message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RejectTagAllUnapprovedResponse>, I>>(base?: I): RejectTagAllUnapprovedResponse {
    return RejectTagAllUnapprovedResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RejectTagAllUnapprovedResponse>, I>>(
    object: I,
  ): RejectTagAllUnapprovedResponse {
    const message = createBaseRejectTagAllUnapprovedResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseSubmitRequest(): SubmitRequest {
  return { config: "", submitter: "", description: "" };
}

export const SubmitRequest = {
  encode(message: SubmitRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(10).string(message.config);
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      writer.uint32(18).string(message.submitter);
    }
    if (message.description !== undefined && message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubmitRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubmitRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.config = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.submitter = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.description = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubmitRequest {
    return {
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      submitter: isSet(object.submitter) ? globalThis.String(object.submitter) : "",
      description: isSet(object.description) ? globalThis.String(object.description) : "",
    };
  },

  toJSON(message: SubmitRequest): unknown {
    const obj: any = {};
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.submitter !== undefined && message.submitter !== "") {
      obj.submitter = message.submitter;
    }
    if (message.description !== undefined && message.description !== "") {
      obj.description = message.description;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubmitRequest>, I>>(base?: I): SubmitRequest {
    return SubmitRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SubmitRequest>, I>>(object: I): SubmitRequest {
    const message = createBaseSubmitRequest();
    message.config = object.config ?? "";
    message.submitter = object.submitter ?? "";
    message.description = object.description ?? "";
    return message;
  },
};

function createBaseSubmitResponse(): SubmitResponse {
  return { _id: "" };
}

export const SubmitResponse = {
  encode(message: SubmitResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubmitResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubmitResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message._id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubmitResponse {
    return { _id: isSet(object._id) ? globalThis.String(object._id) : "" };
  },

  toJSON(message: SubmitResponse): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubmitResponse>, I>>(base?: I): SubmitResponse {
    return SubmitResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SubmitResponse>, I>>(object: I): SubmitResponse {
    const message = createBaseSubmitResponse();
    message._id = object._id ?? "";
    return message;
  },
};

function createBaseDeletePendingRequest(): DeletePendingRequest {
  return { _id: "", sender: "" };
}

export const DeletePendingRequest = {
  encode(message: DeletePendingRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    if (message.sender !== undefined && message.sender !== "") {
      writer.uint32(18).string(message.sender);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeletePendingRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeletePendingRequest();
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

          message.sender = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeletePendingRequest {
    return {
      _id: isSet(object._id) ? globalThis.String(object._id) : "",
      sender: isSet(object.sender) ? globalThis.String(object.sender) : "",
    };
  },

  toJSON(message: DeletePendingRequest): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    if (message.sender !== undefined && message.sender !== "") {
      obj.sender = message.sender;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeletePendingRequest>, I>>(base?: I): DeletePendingRequest {
    return DeletePendingRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeletePendingRequest>, I>>(object: I): DeletePendingRequest {
    const message = createBaseDeletePendingRequest();
    message._id = object._id ?? "";
    message.sender = object.sender ?? "";
    return message;
  },
};

function createBaseDeletePendingResponse(): DeletePendingResponse {
  return { _id: "" };
}

export const DeletePendingResponse = {
  encode(message: DeletePendingResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeletePendingResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeletePendingResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message._id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeletePendingResponse {
    return { _id: isSet(object._id) ? globalThis.String(object._id) : "" };
  },

  toJSON(message: DeletePendingResponse): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeletePendingResponse>, I>>(base?: I): DeletePendingResponse {
    return DeletePendingResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeletePendingResponse>, I>>(object: I): DeletePendingResponse {
    const message = createBaseDeletePendingResponse();
    message._id = object._id ?? "";
    return message;
  },
};

function createBaseGetWorkRequest(): GetWorkRequest {
  return {};
}

export const GetWorkRequest = {
  encode(_: GetWorkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetWorkRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetWorkRequest();
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

  fromJSON(_: any): GetWorkRequest {
    return {};
  },

  toJSON(_: GetWorkRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<GetWorkRequest>, I>>(base?: I): GetWorkRequest {
    return GetWorkRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetWorkRequest>, I>>(_: I): GetWorkRequest {
    const message = createBaseGetWorkRequest();
    return message;
  },
};

function createBaseGetWorkResponse(): GetWorkResponse {
  return { data: [] };
}

export const GetWorkResponse = {
  encode(message: GetWorkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== undefined && message.data.length !== 0) {
      for (const v of message.data) {
        ComputeWork.encode(v!, writer.uint32(10).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetWorkResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetWorkResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data!.push(ComputeWork.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetWorkResponse {
    return { data: globalThis.Array.isArray(object?.data) ? object.data.map((e: any) => ComputeWork.fromJSON(e)) : [] };
  },

  toJSON(message: GetWorkResponse): unknown {
    const obj: any = {};
    if (message.data?.length) {
      obj.data = message.data.map((e) => ComputeWork.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetWorkResponse>, I>>(base?: I): GetWorkResponse {
    return GetWorkResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetWorkResponse>, I>>(object: I): GetWorkResponse {
    const message = createBaseGetWorkResponse();
    message.data = object.data?.map((e) => ComputeWork.fromPartial(e)) || [];
    return message;
  },
};

function createBaseRejectWorkRequest(): RejectWorkRequest {
  return { id: "", reason: "", hash: "" };
}

export const RejectWorkRequest = {
  encode(message: RejectWorkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.reason !== undefined && message.reason !== "") {
      writer.uint32(18).string(message.reason);
    }
    if (message.hash !== undefined && message.hash !== "") {
      writer.uint32(26).string(message.hash);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RejectWorkRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRejectWorkRequest();
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

          message.reason = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.hash = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RejectWorkRequest {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      reason: isSet(object.reason) ? globalThis.String(object.reason) : "",
      hash: isSet(object.hash) ? globalThis.String(object.hash) : "",
    };
  },

  toJSON(message: RejectWorkRequest): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== "") {
      obj.id = message.id;
    }
    if (message.reason !== undefined && message.reason !== "") {
      obj.reason = message.reason;
    }
    if (message.hash !== undefined && message.hash !== "") {
      obj.hash = message.hash;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RejectWorkRequest>, I>>(base?: I): RejectWorkRequest {
    return RejectWorkRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RejectWorkRequest>, I>>(object: I): RejectWorkRequest {
    const message = createBaseRejectWorkRequest();
    message.id = object.id ?? "";
    message.reason = object.reason ?? "";
    message.hash = object.hash ?? "";
    return message;
  },
};

function createBaseRejectWorkResponse(): RejectWorkResponse {
  return {};
}

export const RejectWorkResponse = {
  encode(_: RejectWorkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RejectWorkResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRejectWorkResponse();
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

  fromJSON(_: any): RejectWorkResponse {
    return {};
  },

  toJSON(_: RejectWorkResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<RejectWorkResponse>, I>>(base?: I): RejectWorkResponse {
    return RejectWorkResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RejectWorkResponse>, I>>(_: I): RejectWorkResponse {
    const message = createBaseRejectWorkResponse();
    return message;
  },
};

function createBaseCompleteWorkRequest(): CompleteWorkRequest {
  return { _id: "", result: undefined };
}

export const CompleteWorkRequest = {
  encode(message: CompleteWorkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    if (message.result !== undefined) {
      SimulationResult.encode(message.result, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CompleteWorkRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCompleteWorkRequest();
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
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CompleteWorkRequest {
    return {
      _id: isSet(object._id) ? globalThis.String(object._id) : "",
      result: isSet(object.result) ? SimulationResult.fromJSON(object.result) : undefined,
    };
  },

  toJSON(message: CompleteWorkRequest): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    if (message.result !== undefined) {
      obj.result = SimulationResult.toJSON(message.result);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CompleteWorkRequest>, I>>(base?: I): CompleteWorkRequest {
    return CompleteWorkRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CompleteWorkRequest>, I>>(object: I): CompleteWorkRequest {
    const message = createBaseCompleteWorkRequest();
    message._id = object._id ?? "";
    message.result = (object.result !== undefined && object.result !== null)
      ? SimulationResult.fromPartial(object.result)
      : undefined;
    return message;
  },
};

function createBaseCompleteWorkResponse(): CompleteWorkResponse {
  return { _id: "" };
}

export const CompleteWorkResponse = {
  encode(message: CompleteWorkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CompleteWorkResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCompleteWorkResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message._id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CompleteWorkResponse {
    return { _id: isSet(object._id) ? globalThis.String(object._id) : "" };
  },

  toJSON(message: CompleteWorkResponse): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CompleteWorkResponse>, I>>(base?: I): CompleteWorkResponse {
    return CompleteWorkResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CompleteWorkResponse>, I>>(object: I): CompleteWorkResponse {
    const message = createBaseCompleteWorkResponse();
    message._id = object._id ?? "";
    return message;
  },
};

function createBaseWorkStatusRequest(): WorkStatusRequest {
  return {};
}

export const WorkStatusRequest = {
  encode(_: WorkStatusRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WorkStatusRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWorkStatusRequest();
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

  fromJSON(_: any): WorkStatusRequest {
    return {};
  },

  toJSON(_: WorkStatusRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<WorkStatusRequest>, I>>(base?: I): WorkStatusRequest {
    return WorkStatusRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<WorkStatusRequest>, I>>(_: I): WorkStatusRequest {
    const message = createBaseWorkStatusRequest();
    return message;
  },
};

function createBaseWorkStatusResponse(): WorkStatusResponse {
  return { todo_count: 0, total_count: 0 };
}

export const WorkStatusResponse = {
  encode(message: WorkStatusResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.todo_count !== undefined && message.todo_count !== 0) {
      writer.uint32(8).int32(message.todo_count);
    }
    if (message.total_count !== undefined && message.total_count !== 0) {
      writer.uint32(16).int32(message.total_count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WorkStatusResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWorkStatusResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.todo_count = reader.int32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.total_count = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): WorkStatusResponse {
    return {
      todo_count: isSet(object.todo_count) ? globalThis.Number(object.todo_count) : 0,
      total_count: isSet(object.total_count) ? globalThis.Number(object.total_count) : 0,
    };
  },

  toJSON(message: WorkStatusResponse): unknown {
    const obj: any = {};
    if (message.todo_count !== undefined && message.todo_count !== 0) {
      obj.todo_count = Math.round(message.todo_count);
    }
    if (message.total_count !== undefined && message.total_count !== 0) {
      obj.total_count = Math.round(message.total_count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<WorkStatusResponse>, I>>(base?: I): WorkStatusResponse {
    return WorkStatusResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<WorkStatusResponse>, I>>(object: I): WorkStatusResponse {
    const message = createBaseWorkStatusResponse();
    message.todo_count = object.todo_count ?? 0;
    message.total_count = object.total_count ?? 0;
    return message;
  },
};

function createBaseReplaceConfigRequest(): ReplaceConfigRequest {
  return { _id: "", config: "", source_tag: 0 };
}

export const ReplaceConfigRequest = {
  encode(message: ReplaceConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    if (message.config !== undefined && message.config !== "") {
      writer.uint32(18).string(message.config);
    }
    if (message.source_tag !== undefined && message.source_tag !== 0) {
      writer.uint32(24).int32(message.source_tag);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplaceConfigRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplaceConfigRequest();
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

          message.config = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.source_tag = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplaceConfigRequest {
    return {
      _id: isSet(object._id) ? globalThis.String(object._id) : "",
      config: isSet(object.config) ? globalThis.String(object.config) : "",
      source_tag: isSet(object.source_tag) ? dBTagFromJSON(object.source_tag) : 0,
    };
  },

  toJSON(message: ReplaceConfigRequest): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    if (message.config !== undefined && message.config !== "") {
      obj.config = message.config;
    }
    if (message.source_tag !== undefined && message.source_tag !== 0) {
      obj.source_tag = dBTagToJSON(message.source_tag);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplaceConfigRequest>, I>>(base?: I): ReplaceConfigRequest {
    return ReplaceConfigRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ReplaceConfigRequest>, I>>(object: I): ReplaceConfigRequest {
    const message = createBaseReplaceConfigRequest();
    message._id = object._id ?? "";
    message.config = object.config ?? "";
    message.source_tag = object.source_tag ?? 0;
    return message;
  },
};

function createBaseReplaceConfigResponse(): ReplaceConfigResponse {
  return { _id: "" };
}

export const ReplaceConfigResponse = {
  encode(message: ReplaceConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplaceConfigResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplaceConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message._id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplaceConfigResponse {
    return { _id: isSet(object._id) ? globalThis.String(object._id) : "" };
  },

  toJSON(message: ReplaceConfigResponse): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplaceConfigResponse>, I>>(base?: I): ReplaceConfigResponse {
    return ReplaceConfigResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ReplaceConfigResponse>, I>>(object: I): ReplaceConfigResponse {
    const message = createBaseReplaceConfigResponse();
    message._id = object._id ?? "";
    return message;
  },
};

function createBaseReplaceDescRequest(): ReplaceDescRequest {
  return { _id: "", desc: "", source_tag: 0 };
}

export const ReplaceDescRequest = {
  encode(message: ReplaceDescRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    if (message.desc !== undefined && message.desc !== "") {
      writer.uint32(18).string(message.desc);
    }
    if (message.source_tag !== undefined && message.source_tag !== 0) {
      writer.uint32(24).int32(message.source_tag);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplaceDescRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplaceDescRequest();
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

          message.desc = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.source_tag = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplaceDescRequest {
    return {
      _id: isSet(object._id) ? globalThis.String(object._id) : "",
      desc: isSet(object.desc) ? globalThis.String(object.desc) : "",
      source_tag: isSet(object.source_tag) ? dBTagFromJSON(object.source_tag) : 0,
    };
  },

  toJSON(message: ReplaceDescRequest): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    if (message.desc !== undefined && message.desc !== "") {
      obj.desc = message.desc;
    }
    if (message.source_tag !== undefined && message.source_tag !== 0) {
      obj.source_tag = dBTagToJSON(message.source_tag);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplaceDescRequest>, I>>(base?: I): ReplaceDescRequest {
    return ReplaceDescRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ReplaceDescRequest>, I>>(object: I): ReplaceDescRequest {
    const message = createBaseReplaceDescRequest();
    message._id = object._id ?? "";
    message.desc = object.desc ?? "";
    message.source_tag = object.source_tag ?? 0;
    return message;
  },
};

function createBaseReplaceDescResponse(): ReplaceDescResponse {
  return { _id: "" };
}

export const ReplaceDescResponse = {
  encode(message: ReplaceDescResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message._id !== undefined && message._id !== "") {
      writer.uint32(10).string(message._id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplaceDescResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplaceDescResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message._id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplaceDescResponse {
    return { _id: isSet(object._id) ? globalThis.String(object._id) : "" };
  },

  toJSON(message: ReplaceDescResponse): unknown {
    const obj: any = {};
    if (message._id !== undefined && message._id !== "") {
      obj._id = message._id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplaceDescResponse>, I>>(base?: I): ReplaceDescResponse {
    return ReplaceDescResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ReplaceDescResponse>, I>>(object: I): ReplaceDescResponse {
    const message = createBaseReplaceDescResponse();
    message._id = object._id ?? "";
    return message;
  },
};

export interface DBStore {
  /** generic get for pulling from approved db */
  Get(request: GetRequest): Promise<GetResponse>;
  GetAll(request: GetAllRequest): Promise<GetAllResponse>;
  GetOne(request: GetOneRequest): Promise<GetOneResponse>;
  GetPending(request: GetPendingRequest): Promise<GetPendingResponse>;
  GetBySubmitter(request: GetBySubmitterRequest): Promise<GetBySubmitterResponse>;
  /** tagging */
  ApproveTag(request: ApproveTagRequest): Promise<ApproveTagResponse>;
  RejectTag(request: RejectTagRequest): Promise<RejectTagResponse>;
  RejectTagAllUnapproved(request: RejectTagAllUnapprovedRequest): Promise<RejectTagAllUnapprovedResponse>;
  /** submissions */
  Submit(request: SubmitRequest): Promise<SubmitResponse>;
  DeletePending(request: DeletePendingRequest): Promise<DeletePendingResponse>;
  /** work related */
  GetWork(request: GetWorkRequest): Promise<GetWorkResponse>;
  CompleteWork(request: CompleteWorkRequest): Promise<CompleteWorkResponse>;
  RejectWork(request: RejectWorkRequest): Promise<RejectWorkResponse>;
  WorkStatus(request: WorkStatusRequest): Promise<WorkStatusResponse>;
  /** admin endpoint */
  ReplaceConfig(request: ReplaceConfigRequest): Promise<ReplaceConfigResponse>;
  ReplaceDesc(request: ReplaceDescRequest): Promise<ReplaceDescResponse>;
}

export const DBStoreServiceName = "db.DBStore";
export class DBStoreClientImpl implements DBStore {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || DBStoreServiceName;
    this.rpc = rpc;
    this.Get = this.Get.bind(this);
    this.GetAll = this.GetAll.bind(this);
    this.GetOne = this.GetOne.bind(this);
    this.GetPending = this.GetPending.bind(this);
    this.GetBySubmitter = this.GetBySubmitter.bind(this);
    this.ApproveTag = this.ApproveTag.bind(this);
    this.RejectTag = this.RejectTag.bind(this);
    this.RejectTagAllUnapproved = this.RejectTagAllUnapproved.bind(this);
    this.Submit = this.Submit.bind(this);
    this.DeletePending = this.DeletePending.bind(this);
    this.GetWork = this.GetWork.bind(this);
    this.CompleteWork = this.CompleteWork.bind(this);
    this.RejectWork = this.RejectWork.bind(this);
    this.WorkStatus = this.WorkStatus.bind(this);
    this.ReplaceConfig = this.ReplaceConfig.bind(this);
    this.ReplaceDesc = this.ReplaceDesc.bind(this);
  }
  Get(request: GetRequest): Promise<GetResponse> {
    const data = GetRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetResponse.decode(_m0.Reader.create(data)));
  }

  GetAll(request: GetAllRequest): Promise<GetAllResponse> {
    const data = GetAllRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetAll", data);
    return promise.then((data) => GetAllResponse.decode(_m0.Reader.create(data)));
  }

  GetOne(request: GetOneRequest): Promise<GetOneResponse> {
    const data = GetOneRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetOne", data);
    return promise.then((data) => GetOneResponse.decode(_m0.Reader.create(data)));
  }

  GetPending(request: GetPendingRequest): Promise<GetPendingResponse> {
    const data = GetPendingRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetPending", data);
    return promise.then((data) => GetPendingResponse.decode(_m0.Reader.create(data)));
  }

  GetBySubmitter(request: GetBySubmitterRequest): Promise<GetBySubmitterResponse> {
    const data = GetBySubmitterRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetBySubmitter", data);
    return promise.then((data) => GetBySubmitterResponse.decode(_m0.Reader.create(data)));
  }

  ApproveTag(request: ApproveTagRequest): Promise<ApproveTagResponse> {
    const data = ApproveTagRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "ApproveTag", data);
    return promise.then((data) => ApproveTagResponse.decode(_m0.Reader.create(data)));
  }

  RejectTag(request: RejectTagRequest): Promise<RejectTagResponse> {
    const data = RejectTagRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RejectTag", data);
    return promise.then((data) => RejectTagResponse.decode(_m0.Reader.create(data)));
  }

  RejectTagAllUnapproved(request: RejectTagAllUnapprovedRequest): Promise<RejectTagAllUnapprovedResponse> {
    const data = RejectTagAllUnapprovedRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RejectTagAllUnapproved", data);
    return promise.then((data) => RejectTagAllUnapprovedResponse.decode(_m0.Reader.create(data)));
  }

  Submit(request: SubmitRequest): Promise<SubmitResponse> {
    const data = SubmitRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Submit", data);
    return promise.then((data) => SubmitResponse.decode(_m0.Reader.create(data)));
  }

  DeletePending(request: DeletePendingRequest): Promise<DeletePendingResponse> {
    const data = DeletePendingRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "DeletePending", data);
    return promise.then((data) => DeletePendingResponse.decode(_m0.Reader.create(data)));
  }

  GetWork(request: GetWorkRequest): Promise<GetWorkResponse> {
    const data = GetWorkRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetWork", data);
    return promise.then((data) => GetWorkResponse.decode(_m0.Reader.create(data)));
  }

  CompleteWork(request: CompleteWorkRequest): Promise<CompleteWorkResponse> {
    const data = CompleteWorkRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CompleteWork", data);
    return promise.then((data) => CompleteWorkResponse.decode(_m0.Reader.create(data)));
  }

  RejectWork(request: RejectWorkRequest): Promise<RejectWorkResponse> {
    const data = RejectWorkRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RejectWork", data);
    return promise.then((data) => RejectWorkResponse.decode(_m0.Reader.create(data)));
  }

  WorkStatus(request: WorkStatusRequest): Promise<WorkStatusResponse> {
    const data = WorkStatusRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "WorkStatus", data);
    return promise.then((data) => WorkStatusResponse.decode(_m0.Reader.create(data)));
  }

  ReplaceConfig(request: ReplaceConfigRequest): Promise<ReplaceConfigResponse> {
    const data = ReplaceConfigRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "ReplaceConfig", data);
    return promise.then((data) => ReplaceConfigResponse.decode(_m0.Reader.create(data)));
  }

  ReplaceDesc(request: ReplaceDescRequest): Promise<ReplaceDescResponse> {
    const data = ReplaceDescRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "ReplaceDesc", data);
    return promise.then((data) => ReplaceDescResponse.decode(_m0.Reader.create(data)));
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

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
