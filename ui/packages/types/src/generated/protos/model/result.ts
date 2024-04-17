/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { SimMode, simModeFromJSON, simModeToJSON } from "./enums";
import { Character, Coord, Enemy, EnergySettings, SimulatorSettings } from "./sim";

export interface Version {
  major?: string | undefined;
  minor?: string | undefined;
}

export interface SimulationResult {
  /** required fields (should always be here regardless of schema version) */
  schema_version?: Version | undefined;
  sim_version?: string | undefined;
  modified?: boolean | undefined;
  build_date?: string | undefined;
  sample_seed?: string | undefined;
  config_file?: string | undefined;
  simulator_settings?: SimulatorSettings | undefined;
  energy_settings?: EnergySettings | undefined;
  initial_character?: string | undefined;
  character_details?: Character[] | undefined;
  target_details?: Enemy[] | undefined;
  player_position?: Coord | undefined;
  incomplete_characters?:
    | string[]
    | undefined;
  /** All data that changes per iteration goes here */
  statistics?:
    | SimulationStatistics
    | undefined;
  /** --- optional metadata fields below --- */
  mode?: SimMode | undefined;
  key_type?:
    | string
    | undefined;
  /** if set to -1 then should result in perm */
  created_date?: number | undefined;
}

export interface SimulationStatistics {
  /** metadata */
  min_seed?: string | undefined;
  max_seed?: string | undefined;
  p25_seed?: string | undefined;
  p50_seed?: string | undefined;
  p75_seed?: string | undefined;
  iterations?:
    | number
    | undefined;
  /** global overview (global/no group by) */
  duration?: OverviewStats | undefined;
  dps?: OverviewStats | undefined;
  rps?: OverviewStats | undefined;
  eps?: OverviewStats | undefined;
  hps?: OverviewStats | undefined;
  shp?: OverviewStats | undefined;
  total_damage?: DescriptiveStats | undefined;
  warnings?: Warnings | undefined;
  failed_actions?:
    | FailedActions[]
    | undefined;
  /** damage */
  element_dps?: { [key: string]: DescriptiveStats } | undefined;
  target_dps?: { [key: number]: DescriptiveStats } | undefined;
  character_dps?: DescriptiveStats[] | undefined;
  dps_by_element?: ElementStats[] | undefined;
  dps_by_target?: TargetStats[] | undefined;
  source_dps?: SourceStats[] | undefined;
  source_damage_instances?: SourceStats[] | undefined;
  damage_buckets?: BucketStats | undefined;
  cumu_damage_contrib?: CharacterBucketStats | undefined;
  cumu_damage?:
    | TargetBucketStats
    | undefined;
  /** shield */
  shields?:
    | { [key: string]: ShieldInfo }
    | undefined;
  /** field time */
  field_time?:
    | DescriptiveStats[]
    | undefined;
  /** total source energy */
  total_source_energy?:
    | SourceStats[]
    | undefined;
  /** source reactions */
  source_reactions?:
    | SourceStats[]
    | undefined;
  /** character actions */
  character_actions?:
    | SourceStats[]
    | undefined;
  /** target aura uptime */
  target_aura_uptime?:
    | SourceStats[]
    | undefined;
  /** misc statistics at the end of each sim */
  end_stats?: EndStats[] | undefined;
}

export interface SimulationStatistics_ElementDpsEntry {
  key: string;
  value?: DescriptiveStats | undefined;
}

export interface SimulationStatistics_TargetDpsEntry {
  key: number;
  value?: DescriptiveStats | undefined;
}

export interface SimulationStatistics_ShieldsEntry {
  key: string;
  value?: ShieldInfo | undefined;
}

export interface SignedSimulationStatistics {
  stats?: SimulationStatistics | undefined;
  hash?: string | undefined;
}

export interface OverviewStats {
  min?: number | undefined;
  max?: number | undefined;
  mean?: number | undefined;
  sd?: number | undefined;
  q1?: number | undefined;
  q2?: number | undefined;
  q3?: number | undefined;
  histogram?: number[] | undefined;
}

export interface DescriptiveStats {
  min?: number | undefined;
  max?: number | undefined;
  mean?: number | undefined;
  sd?: number | undefined;
}

export interface ElementStats {
  elements?: { [key: string]: DescriptiveStats } | undefined;
}

export interface ElementStats_ElementsEntry {
  key: string;
  value?: DescriptiveStats | undefined;
}

export interface TargetStats {
  targets?: { [key: number]: DescriptiveStats } | undefined;
}

export interface TargetStats_TargetsEntry {
  key: number;
  value?: DescriptiveStats | undefined;
}

export interface SourceStats {
  sources?: { [key: string]: DescriptiveStats } | undefined;
}

export interface SourceStats_SourcesEntry {
  key: string;
  value?: DescriptiveStats | undefined;
}

export interface BucketStats {
  bucket_size?: number | undefined;
  buckets?: DescriptiveStats[] | undefined;
}

export interface CharacterBucketStats {
  bucket_size?: number | undefined;
  characters?: CharacterBuckets[] | undefined;
}

export interface CharacterBuckets {
  buckets?: DescriptiveStats[] | undefined;
}

export interface TargetBucketStats {
  bucket_size?: number | undefined;
  targets?: { [key: number]: TargetBuckets } | undefined;
}

export interface TargetBucketStats_TargetsEntry {
  key: number;
  value?: TargetBuckets | undefined;
}

export interface TargetBuckets {
  overall?: TargetBucket | undefined;
  target?: TargetBucket | undefined;
}

export interface TargetBucket {
  min?: number[] | undefined;
  max?: number[] | undefined;
  q1?: number[] | undefined;
  q2?: number[] | undefined;
  q3?: number[] | undefined;
}

export interface Warnings {
  /** optional unnecessary, missing == false in ui */
  target_overlap?: boolean | undefined;
  insufficient_energy?: boolean | undefined;
  insufficient_stamina?: boolean | undefined;
  swap_cd?: boolean | undefined;
  skill_cd?: boolean | undefined;
  dash_cd?: boolean | undefined;
}

export interface FailedActions {
  insufficient_energy?: DescriptiveStats | undefined;
  insufficient_stamina?: DescriptiveStats | undefined;
  swap_cd?: DescriptiveStats | undefined;
  skill_cd?: DescriptiveStats | undefined;
  dash_cd?: DescriptiveStats | undefined;
}

export interface ShieldInfo {
  hp?: { [key: string]: DescriptiveStats } | undefined;
  uptime?: DescriptiveStats | undefined;
}

export interface ShieldInfo_HpEntry {
  key: string;
  value?: DescriptiveStats | undefined;
}

export interface EndStats {
  ending_energy?: DescriptiveStats | undefined;
}

function createBaseVersion(): Version {
  return { major: "", minor: "" };
}

export const Version = {
  encode(message: Version, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.major !== undefined && message.major !== "") {
      writer.uint32(10).string(message.major);
    }
    if (message.minor !== undefined && message.minor !== "") {
      writer.uint32(18).string(message.minor);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Version {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVersion();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.major = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.minor = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Version {
    return {
      major: isSet(object.major) ? globalThis.String(object.major) : "",
      minor: isSet(object.minor) ? globalThis.String(object.minor) : "",
    };
  },

  toJSON(message: Version): unknown {
    const obj: any = {};
    if (message.major !== undefined && message.major !== "") {
      obj.major = message.major;
    }
    if (message.minor !== undefined && message.minor !== "") {
      obj.minor = message.minor;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Version>, I>>(base?: I): Version {
    return Version.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Version>, I>>(object: I): Version {
    const message = createBaseVersion();
    message.major = object.major ?? "";
    message.minor = object.minor ?? "";
    return message;
  },
};

function createBaseSimulationResult(): SimulationResult {
  return {
    schema_version: undefined,
    sim_version: undefined,
    modified: undefined,
    build_date: "",
    sample_seed: "",
    config_file: "",
    simulator_settings: undefined,
    energy_settings: undefined,
    initial_character: "",
    character_details: [],
    target_details: [],
    player_position: undefined,
    incomplete_characters: [],
    statistics: undefined,
    mode: 0,
    key_type: "",
    created_date: 0,
  };
}

export const SimulationResult = {
  encode(message: SimulationResult, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.schema_version !== undefined) {
      Version.encode(message.schema_version, writer.uint32(10).fork()).ldelim();
    }
    if (message.sim_version !== undefined) {
      writer.uint32(18).string(message.sim_version);
    }
    if (message.modified !== undefined) {
      writer.uint32(24).bool(message.modified);
    }
    if (message.build_date !== undefined && message.build_date !== "") {
      writer.uint32(34).string(message.build_date);
    }
    if (message.sample_seed !== undefined && message.sample_seed !== "") {
      writer.uint32(42).string(message.sample_seed);
    }
    if (message.config_file !== undefined && message.config_file !== "") {
      writer.uint32(50).string(message.config_file);
    }
    if (message.simulator_settings !== undefined) {
      SimulatorSettings.encode(message.simulator_settings, writer.uint32(58).fork()).ldelim();
    }
    if (message.energy_settings !== undefined) {
      EnergySettings.encode(message.energy_settings, writer.uint32(66).fork()).ldelim();
    }
    if (message.initial_character !== undefined && message.initial_character !== "") {
      writer.uint32(74).string(message.initial_character);
    }
    if (message.character_details !== undefined && message.character_details.length !== 0) {
      for (const v of message.character_details) {
        Character.encode(v!, writer.uint32(82).fork()).ldelim();
      }
    }
    if (message.target_details !== undefined && message.target_details.length !== 0) {
      for (const v of message.target_details) {
        Enemy.encode(v!, writer.uint32(90).fork()).ldelim();
      }
    }
    if (message.player_position !== undefined) {
      Coord.encode(message.player_position, writer.uint32(130).fork()).ldelim();
    }
    if (message.incomplete_characters !== undefined && message.incomplete_characters.length !== 0) {
      for (const v of message.incomplete_characters) {
        writer.uint32(138).string(v!);
      }
    }
    if (message.statistics !== undefined) {
      SimulationStatistics.encode(message.statistics, writer.uint32(98).fork()).ldelim();
    }
    if (message.mode !== undefined && message.mode !== 0) {
      writer.uint32(104).int32(message.mode);
    }
    if (message.key_type !== undefined && message.key_type !== "") {
      writer.uint32(114).string(message.key_type);
    }
    if (message.created_date !== undefined && message.created_date !== 0) {
      writer.uint32(120).int64(message.created_date);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SimulationResult {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSimulationResult();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.schema_version = Version.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
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

          message.build_date = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.sample_seed = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.config_file = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.simulator_settings = SimulatorSettings.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.energy_settings = EnergySettings.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.initial_character = reader.string();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.character_details!.push(Character.decode(reader, reader.uint32()));
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.target_details!.push(Enemy.decode(reader, reader.uint32()));
          continue;
        case 16:
          if (tag !== 130) {
            break;
          }

          message.player_position = Coord.decode(reader, reader.uint32());
          continue;
        case 17:
          if (tag !== 138) {
            break;
          }

          message.incomplete_characters!.push(reader.string());
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.statistics = SimulationStatistics.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag !== 104) {
            break;
          }

          message.mode = reader.int32() as any;
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.key_type = reader.string();
          continue;
        case 15:
          if (tag !== 120) {
            break;
          }

          message.created_date = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SimulationResult {
    return {
      schema_version: isSet(object.schema_version) ? Version.fromJSON(object.schema_version) : undefined,
      sim_version: isSet(object.sim_version) ? globalThis.String(object.sim_version) : undefined,
      modified: isSet(object.modified) ? globalThis.Boolean(object.modified) : undefined,
      build_date: isSet(object.build_date) ? globalThis.String(object.build_date) : "",
      sample_seed: isSet(object.sample_seed) ? globalThis.String(object.sample_seed) : "",
      config_file: isSet(object.config_file) ? globalThis.String(object.config_file) : "",
      simulator_settings: isSet(object.simulator_settings)
        ? SimulatorSettings.fromJSON(object.simulator_settings)
        : undefined,
      energy_settings: isSet(object.energy_settings) ? EnergySettings.fromJSON(object.energy_settings) : undefined,
      initial_character: isSet(object.initial_character) ? globalThis.String(object.initial_character) : "",
      character_details: globalThis.Array.isArray(object?.character_details)
        ? object.character_details.map((e: any) => Character.fromJSON(e))
        : [],
      target_details: globalThis.Array.isArray(object?.target_details)
        ? object.target_details.map((e: any) => Enemy.fromJSON(e))
        : [],
      player_position: isSet(object.player_position) ? Coord.fromJSON(object.player_position) : undefined,
      incomplete_characters: globalThis.Array.isArray(object?.incomplete_characters)
        ? object.incomplete_characters.map((e: any) => globalThis.String(e))
        : [],
      statistics: isSet(object.statistics) ? SimulationStatistics.fromJSON(object.statistics) : undefined,
      mode: isSet(object.mode) ? simModeFromJSON(object.mode) : 0,
      key_type: isSet(object.key_type) ? globalThis.String(object.key_type) : "",
      created_date: isSet(object.created_date) ? globalThis.Number(object.created_date) : 0,
    };
  },

  toJSON(message: SimulationResult): unknown {
    const obj: any = {};
    if (message.schema_version !== undefined) {
      obj.schema_version = Version.toJSON(message.schema_version);
    }
    if (message.sim_version !== undefined) {
      obj.sim_version = message.sim_version;
    }
    if (message.modified !== undefined) {
      obj.modified = message.modified;
    }
    if (message.build_date !== undefined && message.build_date !== "") {
      obj.build_date = message.build_date;
    }
    if (message.sample_seed !== undefined && message.sample_seed !== "") {
      obj.sample_seed = message.sample_seed;
    }
    if (message.config_file !== undefined && message.config_file !== "") {
      obj.config_file = message.config_file;
    }
    if (message.simulator_settings !== undefined) {
      obj.simulator_settings = SimulatorSettings.toJSON(message.simulator_settings);
    }
    if (message.energy_settings !== undefined) {
      obj.energy_settings = EnergySettings.toJSON(message.energy_settings);
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
    if (message.player_position !== undefined) {
      obj.player_position = Coord.toJSON(message.player_position);
    }
    if (message.incomplete_characters?.length) {
      obj.incomplete_characters = message.incomplete_characters;
    }
    if (message.statistics !== undefined) {
      obj.statistics = SimulationStatistics.toJSON(message.statistics);
    }
    if (message.mode !== undefined && message.mode !== 0) {
      obj.mode = simModeToJSON(message.mode);
    }
    if (message.key_type !== undefined && message.key_type !== "") {
      obj.key_type = message.key_type;
    }
    if (message.created_date !== undefined && message.created_date !== 0) {
      obj.created_date = Math.round(message.created_date);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SimulationResult>, I>>(base?: I): SimulationResult {
    return SimulationResult.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SimulationResult>, I>>(object: I): SimulationResult {
    const message = createBaseSimulationResult();
    message.schema_version = (object.schema_version !== undefined && object.schema_version !== null)
      ? Version.fromPartial(object.schema_version)
      : undefined;
    message.sim_version = object.sim_version ?? undefined;
    message.modified = object.modified ?? undefined;
    message.build_date = object.build_date ?? "";
    message.sample_seed = object.sample_seed ?? "";
    message.config_file = object.config_file ?? "";
    message.simulator_settings = (object.simulator_settings !== undefined && object.simulator_settings !== null)
      ? SimulatorSettings.fromPartial(object.simulator_settings)
      : undefined;
    message.energy_settings = (object.energy_settings !== undefined && object.energy_settings !== null)
      ? EnergySettings.fromPartial(object.energy_settings)
      : undefined;
    message.initial_character = object.initial_character ?? "";
    message.character_details = object.character_details?.map((e) => Character.fromPartial(e)) || [];
    message.target_details = object.target_details?.map((e) => Enemy.fromPartial(e)) || [];
    message.player_position = (object.player_position !== undefined && object.player_position !== null)
      ? Coord.fromPartial(object.player_position)
      : undefined;
    message.incomplete_characters = object.incomplete_characters?.map((e) => e) || [];
    message.statistics = (object.statistics !== undefined && object.statistics !== null)
      ? SimulationStatistics.fromPartial(object.statistics)
      : undefined;
    message.mode = object.mode ?? 0;
    message.key_type = object.key_type ?? "";
    message.created_date = object.created_date ?? 0;
    return message;
  },
};

function createBaseSimulationStatistics(): SimulationStatistics {
  return {
    min_seed: "",
    max_seed: "",
    p25_seed: "",
    p50_seed: "",
    p75_seed: "",
    iterations: 0,
    duration: undefined,
    dps: undefined,
    rps: undefined,
    eps: undefined,
    hps: undefined,
    shp: undefined,
    total_damage: undefined,
    warnings: undefined,
    failed_actions: [],
    element_dps: {},
    target_dps: {},
    character_dps: [],
    dps_by_element: [],
    dps_by_target: [],
    source_dps: [],
    source_damage_instances: [],
    damage_buckets: undefined,
    cumu_damage_contrib: undefined,
    cumu_damage: undefined,
    shields: {},
    field_time: [],
    total_source_energy: [],
    source_reactions: [],
    character_actions: [],
    target_aura_uptime: [],
    end_stats: [],
  };
}

export const SimulationStatistics = {
  encode(message: SimulationStatistics, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.min_seed !== undefined && message.min_seed !== "") {
      writer.uint32(10).string(message.min_seed);
    }
    if (message.max_seed !== undefined && message.max_seed !== "") {
      writer.uint32(18).string(message.max_seed);
    }
    if (message.p25_seed !== undefined && message.p25_seed !== "") {
      writer.uint32(26).string(message.p25_seed);
    }
    if (message.p50_seed !== undefined && message.p50_seed !== "") {
      writer.uint32(34).string(message.p50_seed);
    }
    if (message.p75_seed !== undefined && message.p75_seed !== "") {
      writer.uint32(42).string(message.p75_seed);
    }
    if (message.iterations !== undefined && message.iterations !== 0) {
      writer.uint32(48).uint32(message.iterations);
    }
    if (message.duration !== undefined) {
      OverviewStats.encode(message.duration, writer.uint32(58).fork()).ldelim();
    }
    if (message.dps !== undefined) {
      OverviewStats.encode(message.dps, writer.uint32(66).fork()).ldelim();
    }
    if (message.rps !== undefined) {
      OverviewStats.encode(message.rps, writer.uint32(74).fork()).ldelim();
    }
    if (message.eps !== undefined) {
      OverviewStats.encode(message.eps, writer.uint32(82).fork()).ldelim();
    }
    if (message.hps !== undefined) {
      OverviewStats.encode(message.hps, writer.uint32(90).fork()).ldelim();
    }
    if (message.shp !== undefined) {
      OverviewStats.encode(message.shp, writer.uint32(98).fork()).ldelim();
    }
    if (message.total_damage !== undefined) {
      DescriptiveStats.encode(message.total_damage, writer.uint32(106).fork()).ldelim();
    }
    if (message.warnings !== undefined) {
      Warnings.encode(message.warnings, writer.uint32(114).fork()).ldelim();
    }
    if (message.failed_actions !== undefined && message.failed_actions.length !== 0) {
      for (const v of message.failed_actions) {
        FailedActions.encode(v!, writer.uint32(122).fork()).ldelim();
      }
    }
    Object.entries(message.element_dps || {}).forEach(([key, value]) => {
      SimulationStatistics_ElementDpsEntry.encode({ key: key as any, value }, writer.uint32(130).fork()).ldelim();
    });
    Object.entries(message.target_dps || {}).forEach(([key, value]) => {
      SimulationStatistics_TargetDpsEntry.encode({ key: key as any, value }, writer.uint32(138).fork()).ldelim();
    });
    if (message.character_dps !== undefined && message.character_dps.length !== 0) {
      for (const v of message.character_dps) {
        DescriptiveStats.encode(v!, writer.uint32(146).fork()).ldelim();
      }
    }
    if (message.dps_by_element !== undefined && message.dps_by_element.length !== 0) {
      for (const v of message.dps_by_element) {
        ElementStats.encode(v!, writer.uint32(154).fork()).ldelim();
      }
    }
    if (message.dps_by_target !== undefined && message.dps_by_target.length !== 0) {
      for (const v of message.dps_by_target) {
        TargetStats.encode(v!, writer.uint32(162).fork()).ldelim();
      }
    }
    if (message.source_dps !== undefined && message.source_dps.length !== 0) {
      for (const v of message.source_dps) {
        SourceStats.encode(v!, writer.uint32(194).fork()).ldelim();
      }
    }
    if (message.source_damage_instances !== undefined && message.source_damage_instances.length !== 0) {
      for (const v of message.source_damage_instances) {
        SourceStats.encode(v!, writer.uint32(242).fork()).ldelim();
      }
    }
    if (message.damage_buckets !== undefined) {
      BucketStats.encode(message.damage_buckets, writer.uint32(170).fork()).ldelim();
    }
    if (message.cumu_damage_contrib !== undefined) {
      CharacterBucketStats.encode(message.cumu_damage_contrib, writer.uint32(178).fork()).ldelim();
    }
    if (message.cumu_damage !== undefined) {
      TargetBucketStats.encode(message.cumu_damage, writer.uint32(250).fork()).ldelim();
    }
    Object.entries(message.shields || {}).forEach(([key, value]) => {
      SimulationStatistics_ShieldsEntry.encode({ key: key as any, value }, writer.uint32(186).fork()).ldelim();
    });
    if (message.field_time !== undefined && message.field_time.length !== 0) {
      for (const v of message.field_time) {
        DescriptiveStats.encode(v!, writer.uint32(202).fork()).ldelim();
      }
    }
    if (message.total_source_energy !== undefined && message.total_source_energy.length !== 0) {
      for (const v of message.total_source_energy) {
        SourceStats.encode(v!, writer.uint32(210).fork()).ldelim();
      }
    }
    if (message.source_reactions !== undefined && message.source_reactions.length !== 0) {
      for (const v of message.source_reactions) {
        SourceStats.encode(v!, writer.uint32(218).fork()).ldelim();
      }
    }
    if (message.character_actions !== undefined && message.character_actions.length !== 0) {
      for (const v of message.character_actions) {
        SourceStats.encode(v!, writer.uint32(226).fork()).ldelim();
      }
    }
    if (message.target_aura_uptime !== undefined && message.target_aura_uptime.length !== 0) {
      for (const v of message.target_aura_uptime) {
        SourceStats.encode(v!, writer.uint32(234).fork()).ldelim();
      }
    }
    if (message.end_stats !== undefined && message.end_stats.length !== 0) {
      for (const v of message.end_stats) {
        EndStats.encode(v!, writer.uint32(258).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SimulationStatistics {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSimulationStatistics();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.min_seed = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.max_seed = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.p25_seed = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.p50_seed = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.p75_seed = reader.string();
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.iterations = reader.uint32();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.duration = OverviewStats.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.dps = OverviewStats.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.rps = OverviewStats.decode(reader, reader.uint32());
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.eps = OverviewStats.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.hps = OverviewStats.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.shp = OverviewStats.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.total_damage = DescriptiveStats.decode(reader, reader.uint32());
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.warnings = Warnings.decode(reader, reader.uint32());
          continue;
        case 15:
          if (tag !== 122) {
            break;
          }

          message.failed_actions!.push(FailedActions.decode(reader, reader.uint32()));
          continue;
        case 16:
          if (tag !== 130) {
            break;
          }

          const entry16 = SimulationStatistics_ElementDpsEntry.decode(reader, reader.uint32());
          if (entry16.value !== undefined) {
            message.element_dps![entry16.key] = entry16.value;
          }
          continue;
        case 17:
          if (tag !== 138) {
            break;
          }

          const entry17 = SimulationStatistics_TargetDpsEntry.decode(reader, reader.uint32());
          if (entry17.value !== undefined) {
            message.target_dps![entry17.key] = entry17.value;
          }
          continue;
        case 18:
          if (tag !== 146) {
            break;
          }

          message.character_dps!.push(DescriptiveStats.decode(reader, reader.uint32()));
          continue;
        case 19:
          if (tag !== 154) {
            break;
          }

          message.dps_by_element!.push(ElementStats.decode(reader, reader.uint32()));
          continue;
        case 20:
          if (tag !== 162) {
            break;
          }

          message.dps_by_target!.push(TargetStats.decode(reader, reader.uint32()));
          continue;
        case 24:
          if (tag !== 194) {
            break;
          }

          message.source_dps!.push(SourceStats.decode(reader, reader.uint32()));
          continue;
        case 30:
          if (tag !== 242) {
            break;
          }

          message.source_damage_instances!.push(SourceStats.decode(reader, reader.uint32()));
          continue;
        case 21:
          if (tag !== 170) {
            break;
          }

          message.damage_buckets = BucketStats.decode(reader, reader.uint32());
          continue;
        case 22:
          if (tag !== 178) {
            break;
          }

          message.cumu_damage_contrib = CharacterBucketStats.decode(reader, reader.uint32());
          continue;
        case 31:
          if (tag !== 250) {
            break;
          }

          message.cumu_damage = TargetBucketStats.decode(reader, reader.uint32());
          continue;
        case 23:
          if (tag !== 186) {
            break;
          }

          const entry23 = SimulationStatistics_ShieldsEntry.decode(reader, reader.uint32());
          if (entry23.value !== undefined) {
            message.shields![entry23.key] = entry23.value;
          }
          continue;
        case 25:
          if (tag !== 202) {
            break;
          }

          message.field_time!.push(DescriptiveStats.decode(reader, reader.uint32()));
          continue;
        case 26:
          if (tag !== 210) {
            break;
          }

          message.total_source_energy!.push(SourceStats.decode(reader, reader.uint32()));
          continue;
        case 27:
          if (tag !== 218) {
            break;
          }

          message.source_reactions!.push(SourceStats.decode(reader, reader.uint32()));
          continue;
        case 28:
          if (tag !== 226) {
            break;
          }

          message.character_actions!.push(SourceStats.decode(reader, reader.uint32()));
          continue;
        case 29:
          if (tag !== 234) {
            break;
          }

          message.target_aura_uptime!.push(SourceStats.decode(reader, reader.uint32()));
          continue;
        case 32:
          if (tag !== 258) {
            break;
          }

          message.end_stats!.push(EndStats.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SimulationStatistics {
    return {
      min_seed: isSet(object.min_seed) ? globalThis.String(object.min_seed) : "",
      max_seed: isSet(object.max_seed) ? globalThis.String(object.max_seed) : "",
      p25_seed: isSet(object.p25_seed) ? globalThis.String(object.p25_seed) : "",
      p50_seed: isSet(object.p50_seed) ? globalThis.String(object.p50_seed) : "",
      p75_seed: isSet(object.p75_seed) ? globalThis.String(object.p75_seed) : "",
      iterations: isSet(object.iterations) ? globalThis.Number(object.iterations) : 0,
      duration: isSet(object.duration) ? OverviewStats.fromJSON(object.duration) : undefined,
      dps: isSet(object.dps) ? OverviewStats.fromJSON(object.dps) : undefined,
      rps: isSet(object.rps) ? OverviewStats.fromJSON(object.rps) : undefined,
      eps: isSet(object.eps) ? OverviewStats.fromJSON(object.eps) : undefined,
      hps: isSet(object.hps) ? OverviewStats.fromJSON(object.hps) : undefined,
      shp: isSet(object.shp) ? OverviewStats.fromJSON(object.shp) : undefined,
      total_damage: isSet(object.total_damage) ? DescriptiveStats.fromJSON(object.total_damage) : undefined,
      warnings: isSet(object.warnings) ? Warnings.fromJSON(object.warnings) : undefined,
      failed_actions: globalThis.Array.isArray(object?.failed_actions)
        ? object.failed_actions.map((e: any) => FailedActions.fromJSON(e))
        : [],
      element_dps: isObject(object.element_dps)
        ? Object.entries(object.element_dps).reduce<{ [key: string]: DescriptiveStats }>((acc, [key, value]) => {
          acc[key] = DescriptiveStats.fromJSON(value);
          return acc;
        }, {})
        : {},
      target_dps: isObject(object.target_dps)
        ? Object.entries(object.target_dps).reduce<{ [key: number]: DescriptiveStats }>((acc, [key, value]) => {
          acc[globalThis.Number(key)] = DescriptiveStats.fromJSON(value);
          return acc;
        }, {})
        : {},
      character_dps: globalThis.Array.isArray(object?.character_dps)
        ? object.character_dps.map((e: any) => DescriptiveStats.fromJSON(e))
        : [],
      dps_by_element: globalThis.Array.isArray(object?.dps_by_element)
        ? object.dps_by_element.map((e: any) => ElementStats.fromJSON(e))
        : [],
      dps_by_target: globalThis.Array.isArray(object?.dps_by_target)
        ? object.dps_by_target.map((e: any) => TargetStats.fromJSON(e))
        : [],
      source_dps: globalThis.Array.isArray(object?.source_dps)
        ? object.source_dps.map((e: any) => SourceStats.fromJSON(e))
        : [],
      source_damage_instances: globalThis.Array.isArray(object?.source_damage_instances)
        ? object.source_damage_instances.map((e: any) => SourceStats.fromJSON(e))
        : [],
      damage_buckets: isSet(object.damage_buckets) ? BucketStats.fromJSON(object.damage_buckets) : undefined,
      cumu_damage_contrib: isSet(object.cumu_damage_contrib)
        ? CharacterBucketStats.fromJSON(object.cumu_damage_contrib)
        : undefined,
      cumu_damage: isSet(object.cumu_damage) ? TargetBucketStats.fromJSON(object.cumu_damage) : undefined,
      shields: isObject(object.shields)
        ? Object.entries(object.shields).reduce<{ [key: string]: ShieldInfo }>((acc, [key, value]) => {
          acc[key] = ShieldInfo.fromJSON(value);
          return acc;
        }, {})
        : {},
      field_time: globalThis.Array.isArray(object?.field_time)
        ? object.field_time.map((e: any) => DescriptiveStats.fromJSON(e))
        : [],
      total_source_energy: globalThis.Array.isArray(object?.total_source_energy)
        ? object.total_source_energy.map((e: any) => SourceStats.fromJSON(e))
        : [],
      source_reactions: globalThis.Array.isArray(object?.source_reactions)
        ? object.source_reactions.map((e: any) => SourceStats.fromJSON(e))
        : [],
      character_actions: globalThis.Array.isArray(object?.character_actions)
        ? object.character_actions.map((e: any) => SourceStats.fromJSON(e))
        : [],
      target_aura_uptime: globalThis.Array.isArray(object?.target_aura_uptime)
        ? object.target_aura_uptime.map((e: any) => SourceStats.fromJSON(e))
        : [],
      end_stats: globalThis.Array.isArray(object?.end_stats)
        ? object.end_stats.map((e: any) => EndStats.fromJSON(e))
        : [],
    };
  },

  toJSON(message: SimulationStatistics): unknown {
    const obj: any = {};
    if (message.min_seed !== undefined && message.min_seed !== "") {
      obj.min_seed = message.min_seed;
    }
    if (message.max_seed !== undefined && message.max_seed !== "") {
      obj.max_seed = message.max_seed;
    }
    if (message.p25_seed !== undefined && message.p25_seed !== "") {
      obj.p25_seed = message.p25_seed;
    }
    if (message.p50_seed !== undefined && message.p50_seed !== "") {
      obj.p50_seed = message.p50_seed;
    }
    if (message.p75_seed !== undefined && message.p75_seed !== "") {
      obj.p75_seed = message.p75_seed;
    }
    if (message.iterations !== undefined && message.iterations !== 0) {
      obj.iterations = Math.round(message.iterations);
    }
    if (message.duration !== undefined) {
      obj.duration = OverviewStats.toJSON(message.duration);
    }
    if (message.dps !== undefined) {
      obj.dps = OverviewStats.toJSON(message.dps);
    }
    if (message.rps !== undefined) {
      obj.rps = OverviewStats.toJSON(message.rps);
    }
    if (message.eps !== undefined) {
      obj.eps = OverviewStats.toJSON(message.eps);
    }
    if (message.hps !== undefined) {
      obj.hps = OverviewStats.toJSON(message.hps);
    }
    if (message.shp !== undefined) {
      obj.shp = OverviewStats.toJSON(message.shp);
    }
    if (message.total_damage !== undefined) {
      obj.total_damage = DescriptiveStats.toJSON(message.total_damage);
    }
    if (message.warnings !== undefined) {
      obj.warnings = Warnings.toJSON(message.warnings);
    }
    if (message.failed_actions?.length) {
      obj.failed_actions = message.failed_actions.map((e) => FailedActions.toJSON(e));
    }
    if (message.element_dps) {
      const entries = Object.entries(message.element_dps);
      if (entries.length > 0) {
        obj.element_dps = {};
        entries.forEach(([k, v]) => {
          obj.element_dps[k] = DescriptiveStats.toJSON(v);
        });
      }
    }
    if (message.target_dps) {
      const entries = Object.entries(message.target_dps);
      if (entries.length > 0) {
        obj.target_dps = {};
        entries.forEach(([k, v]) => {
          obj.target_dps[k] = DescriptiveStats.toJSON(v);
        });
      }
    }
    if (message.character_dps?.length) {
      obj.character_dps = message.character_dps.map((e) => DescriptiveStats.toJSON(e));
    }
    if (message.dps_by_element?.length) {
      obj.dps_by_element = message.dps_by_element.map((e) => ElementStats.toJSON(e));
    }
    if (message.dps_by_target?.length) {
      obj.dps_by_target = message.dps_by_target.map((e) => TargetStats.toJSON(e));
    }
    if (message.source_dps?.length) {
      obj.source_dps = message.source_dps.map((e) => SourceStats.toJSON(e));
    }
    if (message.source_damage_instances?.length) {
      obj.source_damage_instances = message.source_damage_instances.map((e) => SourceStats.toJSON(e));
    }
    if (message.damage_buckets !== undefined) {
      obj.damage_buckets = BucketStats.toJSON(message.damage_buckets);
    }
    if (message.cumu_damage_contrib !== undefined) {
      obj.cumu_damage_contrib = CharacterBucketStats.toJSON(message.cumu_damage_contrib);
    }
    if (message.cumu_damage !== undefined) {
      obj.cumu_damage = TargetBucketStats.toJSON(message.cumu_damage);
    }
    if (message.shields) {
      const entries = Object.entries(message.shields);
      if (entries.length > 0) {
        obj.shields = {};
        entries.forEach(([k, v]) => {
          obj.shields[k] = ShieldInfo.toJSON(v);
        });
      }
    }
    if (message.field_time?.length) {
      obj.field_time = message.field_time.map((e) => DescriptiveStats.toJSON(e));
    }
    if (message.total_source_energy?.length) {
      obj.total_source_energy = message.total_source_energy.map((e) => SourceStats.toJSON(e));
    }
    if (message.source_reactions?.length) {
      obj.source_reactions = message.source_reactions.map((e) => SourceStats.toJSON(e));
    }
    if (message.character_actions?.length) {
      obj.character_actions = message.character_actions.map((e) => SourceStats.toJSON(e));
    }
    if (message.target_aura_uptime?.length) {
      obj.target_aura_uptime = message.target_aura_uptime.map((e) => SourceStats.toJSON(e));
    }
    if (message.end_stats?.length) {
      obj.end_stats = message.end_stats.map((e) => EndStats.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SimulationStatistics>, I>>(base?: I): SimulationStatistics {
    return SimulationStatistics.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SimulationStatistics>, I>>(object: I): SimulationStatistics {
    const message = createBaseSimulationStatistics();
    message.min_seed = object.min_seed ?? "";
    message.max_seed = object.max_seed ?? "";
    message.p25_seed = object.p25_seed ?? "";
    message.p50_seed = object.p50_seed ?? "";
    message.p75_seed = object.p75_seed ?? "";
    message.iterations = object.iterations ?? 0;
    message.duration = (object.duration !== undefined && object.duration !== null)
      ? OverviewStats.fromPartial(object.duration)
      : undefined;
    message.dps = (object.dps !== undefined && object.dps !== null) ? OverviewStats.fromPartial(object.dps) : undefined;
    message.rps = (object.rps !== undefined && object.rps !== null) ? OverviewStats.fromPartial(object.rps) : undefined;
    message.eps = (object.eps !== undefined && object.eps !== null) ? OverviewStats.fromPartial(object.eps) : undefined;
    message.hps = (object.hps !== undefined && object.hps !== null) ? OverviewStats.fromPartial(object.hps) : undefined;
    message.shp = (object.shp !== undefined && object.shp !== null) ? OverviewStats.fromPartial(object.shp) : undefined;
    message.total_damage = (object.total_damage !== undefined && object.total_damage !== null)
      ? DescriptiveStats.fromPartial(object.total_damage)
      : undefined;
    message.warnings = (object.warnings !== undefined && object.warnings !== null)
      ? Warnings.fromPartial(object.warnings)
      : undefined;
    message.failed_actions = object.failed_actions?.map((e) => FailedActions.fromPartial(e)) || [];
    message.element_dps = Object.entries(object.element_dps ?? {}).reduce<{ [key: string]: DescriptiveStats }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = DescriptiveStats.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.target_dps = Object.entries(object.target_dps ?? {}).reduce<{ [key: number]: DescriptiveStats }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[globalThis.Number(key)] = DescriptiveStats.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.character_dps = object.character_dps?.map((e) => DescriptiveStats.fromPartial(e)) || [];
    message.dps_by_element = object.dps_by_element?.map((e) => ElementStats.fromPartial(e)) || [];
    message.dps_by_target = object.dps_by_target?.map((e) => TargetStats.fromPartial(e)) || [];
    message.source_dps = object.source_dps?.map((e) => SourceStats.fromPartial(e)) || [];
    message.source_damage_instances = object.source_damage_instances?.map((e) => SourceStats.fromPartial(e)) || [];
    message.damage_buckets = (object.damage_buckets !== undefined && object.damage_buckets !== null)
      ? BucketStats.fromPartial(object.damage_buckets)
      : undefined;
    message.cumu_damage_contrib = (object.cumu_damage_contrib !== undefined && object.cumu_damage_contrib !== null)
      ? CharacterBucketStats.fromPartial(object.cumu_damage_contrib)
      : undefined;
    message.cumu_damage = (object.cumu_damage !== undefined && object.cumu_damage !== null)
      ? TargetBucketStats.fromPartial(object.cumu_damage)
      : undefined;
    message.shields = Object.entries(object.shields ?? {}).reduce<{ [key: string]: ShieldInfo }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = ShieldInfo.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.field_time = object.field_time?.map((e) => DescriptiveStats.fromPartial(e)) || [];
    message.total_source_energy = object.total_source_energy?.map((e) => SourceStats.fromPartial(e)) || [];
    message.source_reactions = object.source_reactions?.map((e) => SourceStats.fromPartial(e)) || [];
    message.character_actions = object.character_actions?.map((e) => SourceStats.fromPartial(e)) || [];
    message.target_aura_uptime = object.target_aura_uptime?.map((e) => SourceStats.fromPartial(e)) || [];
    message.end_stats = object.end_stats?.map((e) => EndStats.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSimulationStatistics_ElementDpsEntry(): SimulationStatistics_ElementDpsEntry {
  return { key: "", value: undefined };
}

export const SimulationStatistics_ElementDpsEntry = {
  encode(message: SimulationStatistics_ElementDpsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      DescriptiveStats.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SimulationStatistics_ElementDpsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSimulationStatistics_ElementDpsEntry();
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

  fromJSON(object: any): SimulationStatistics_ElementDpsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? DescriptiveStats.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: SimulationStatistics_ElementDpsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = DescriptiveStats.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SimulationStatistics_ElementDpsEntry>, I>>(
    base?: I,
  ): SimulationStatistics_ElementDpsEntry {
    return SimulationStatistics_ElementDpsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SimulationStatistics_ElementDpsEntry>, I>>(
    object: I,
  ): SimulationStatistics_ElementDpsEntry {
    const message = createBaseSimulationStatistics_ElementDpsEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? DescriptiveStats.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseSimulationStatistics_TargetDpsEntry(): SimulationStatistics_TargetDpsEntry {
  return { key: 0, value: undefined };
}

export const SimulationStatistics_TargetDpsEntry = {
  encode(message: SimulationStatistics_TargetDpsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).int32(message.key);
    }
    if (message.value !== undefined) {
      DescriptiveStats.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SimulationStatistics_TargetDpsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSimulationStatistics_TargetDpsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.key = reader.int32();
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

  fromJSON(object: any): SimulationStatistics_TargetDpsEntry {
    return {
      key: isSet(object.key) ? globalThis.Number(object.key) : 0,
      value: isSet(object.value) ? DescriptiveStats.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: SimulationStatistics_TargetDpsEntry): unknown {
    const obj: any = {};
    if (message.key !== 0) {
      obj.key = Math.round(message.key);
    }
    if (message.value !== undefined) {
      obj.value = DescriptiveStats.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SimulationStatistics_TargetDpsEntry>, I>>(
    base?: I,
  ): SimulationStatistics_TargetDpsEntry {
    return SimulationStatistics_TargetDpsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SimulationStatistics_TargetDpsEntry>, I>>(
    object: I,
  ): SimulationStatistics_TargetDpsEntry {
    const message = createBaseSimulationStatistics_TargetDpsEntry();
    message.key = object.key ?? 0;
    message.value = (object.value !== undefined && object.value !== null)
      ? DescriptiveStats.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseSimulationStatistics_ShieldsEntry(): SimulationStatistics_ShieldsEntry {
  return { key: "", value: undefined };
}

export const SimulationStatistics_ShieldsEntry = {
  encode(message: SimulationStatistics_ShieldsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      ShieldInfo.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SimulationStatistics_ShieldsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSimulationStatistics_ShieldsEntry();
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

          message.value = ShieldInfo.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SimulationStatistics_ShieldsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? ShieldInfo.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: SimulationStatistics_ShieldsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = ShieldInfo.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SimulationStatistics_ShieldsEntry>, I>>(
    base?: I,
  ): SimulationStatistics_ShieldsEntry {
    return SimulationStatistics_ShieldsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SimulationStatistics_ShieldsEntry>, I>>(
    object: I,
  ): SimulationStatistics_ShieldsEntry {
    const message = createBaseSimulationStatistics_ShieldsEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? ShieldInfo.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseSignedSimulationStatistics(): SignedSimulationStatistics {
  return { stats: undefined, hash: "" };
}

export const SignedSimulationStatistics = {
  encode(message: SignedSimulationStatistics, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stats !== undefined) {
      SimulationStatistics.encode(message.stats, writer.uint32(10).fork()).ldelim();
    }
    if (message.hash !== undefined && message.hash !== "") {
      writer.uint32(18).string(message.hash);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SignedSimulationStatistics {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSignedSimulationStatistics();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.stats = SimulationStatistics.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
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

  fromJSON(object: any): SignedSimulationStatistics {
    return {
      stats: isSet(object.stats) ? SimulationStatistics.fromJSON(object.stats) : undefined,
      hash: isSet(object.hash) ? globalThis.String(object.hash) : "",
    };
  },

  toJSON(message: SignedSimulationStatistics): unknown {
    const obj: any = {};
    if (message.stats !== undefined) {
      obj.stats = SimulationStatistics.toJSON(message.stats);
    }
    if (message.hash !== undefined && message.hash !== "") {
      obj.hash = message.hash;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SignedSimulationStatistics>, I>>(base?: I): SignedSimulationStatistics {
    return SignedSimulationStatistics.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SignedSimulationStatistics>, I>>(object: I): SignedSimulationStatistics {
    const message = createBaseSignedSimulationStatistics();
    message.stats = (object.stats !== undefined && object.stats !== null)
      ? SimulationStatistics.fromPartial(object.stats)
      : undefined;
    message.hash = object.hash ?? "";
    return message;
  },
};

function createBaseOverviewStats(): OverviewStats {
  return {
    min: undefined,
    max: undefined,
    mean: undefined,
    sd: undefined,
    q1: undefined,
    q2: undefined,
    q3: undefined,
    histogram: [],
  };
}

export const OverviewStats = {
  encode(message: OverviewStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.min !== undefined) {
      writer.uint32(9).double(message.min);
    }
    if (message.max !== undefined) {
      writer.uint32(17).double(message.max);
    }
    if (message.mean !== undefined) {
      writer.uint32(25).double(message.mean);
    }
    if (message.sd !== undefined) {
      writer.uint32(33).double(message.sd);
    }
    if (message.q1 !== undefined) {
      writer.uint32(41).double(message.q1);
    }
    if (message.q2 !== undefined) {
      writer.uint32(49).double(message.q2);
    }
    if (message.q3 !== undefined) {
      writer.uint32(57).double(message.q3);
    }
    if (message.histogram !== undefined && message.histogram.length !== 0) {
      writer.uint32(66).fork();
      for (const v of message.histogram) {
        writer.uint32(v);
      }
      writer.ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OverviewStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOverviewStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.min = reader.double();
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.max = reader.double();
          continue;
        case 3:
          if (tag !== 25) {
            break;
          }

          message.mean = reader.double();
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.sd = reader.double();
          continue;
        case 5:
          if (tag !== 41) {
            break;
          }

          message.q1 = reader.double();
          continue;
        case 6:
          if (tag !== 49) {
            break;
          }

          message.q2 = reader.double();
          continue;
        case 7:
          if (tag !== 57) {
            break;
          }

          message.q3 = reader.double();
          continue;
        case 8:
          if (tag === 64) {
            message.histogram!.push(reader.uint32());

            continue;
          }

          if (tag === 66) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.histogram!.push(reader.uint32());
            }

            continue;
          }

          break;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OverviewStats {
    return {
      min: isSet(object.min) ? globalThis.Number(object.min) : undefined,
      max: isSet(object.max) ? globalThis.Number(object.max) : undefined,
      mean: isSet(object.mean) ? globalThis.Number(object.mean) : undefined,
      sd: isSet(object.sd) ? globalThis.Number(object.sd) : undefined,
      q1: isSet(object.q1) ? globalThis.Number(object.q1) : undefined,
      q2: isSet(object.q2) ? globalThis.Number(object.q2) : undefined,
      q3: isSet(object.q3) ? globalThis.Number(object.q3) : undefined,
      histogram: globalThis.Array.isArray(object?.histogram)
        ? object.histogram.map((e: any) => globalThis.Number(e))
        : [],
    };
  },

  toJSON(message: OverviewStats): unknown {
    const obj: any = {};
    if (message.min !== undefined) {
      obj.min = message.min;
    }
    if (message.max !== undefined) {
      obj.max = message.max;
    }
    if (message.mean !== undefined) {
      obj.mean = message.mean;
    }
    if (message.sd !== undefined) {
      obj.sd = message.sd;
    }
    if (message.q1 !== undefined) {
      obj.q1 = message.q1;
    }
    if (message.q2 !== undefined) {
      obj.q2 = message.q2;
    }
    if (message.q3 !== undefined) {
      obj.q3 = message.q3;
    }
    if (message.histogram?.length) {
      obj.histogram = message.histogram.map((e) => Math.round(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<OverviewStats>, I>>(base?: I): OverviewStats {
    return OverviewStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<OverviewStats>, I>>(object: I): OverviewStats {
    const message = createBaseOverviewStats();
    message.min = object.min ?? undefined;
    message.max = object.max ?? undefined;
    message.mean = object.mean ?? undefined;
    message.sd = object.sd ?? undefined;
    message.q1 = object.q1 ?? undefined;
    message.q2 = object.q2 ?? undefined;
    message.q3 = object.q3 ?? undefined;
    message.histogram = object.histogram?.map((e) => e) || [];
    return message;
  },
};

function createBaseDescriptiveStats(): DescriptiveStats {
  return { min: undefined, max: undefined, mean: undefined, sd: undefined };
}

export const DescriptiveStats = {
  encode(message: DescriptiveStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.min !== undefined) {
      writer.uint32(9).double(message.min);
    }
    if (message.max !== undefined) {
      writer.uint32(17).double(message.max);
    }
    if (message.mean !== undefined) {
      writer.uint32(25).double(message.mean);
    }
    if (message.sd !== undefined) {
      writer.uint32(33).double(message.sd);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DescriptiveStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDescriptiveStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.min = reader.double();
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.max = reader.double();
          continue;
        case 3:
          if (tag !== 25) {
            break;
          }

          message.mean = reader.double();
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.sd = reader.double();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DescriptiveStats {
    return {
      min: isSet(object.min) ? globalThis.Number(object.min) : undefined,
      max: isSet(object.max) ? globalThis.Number(object.max) : undefined,
      mean: isSet(object.mean) ? globalThis.Number(object.mean) : undefined,
      sd: isSet(object.sd) ? globalThis.Number(object.sd) : undefined,
    };
  },

  toJSON(message: DescriptiveStats): unknown {
    const obj: any = {};
    if (message.min !== undefined) {
      obj.min = message.min;
    }
    if (message.max !== undefined) {
      obj.max = message.max;
    }
    if (message.mean !== undefined) {
      obj.mean = message.mean;
    }
    if (message.sd !== undefined) {
      obj.sd = message.sd;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DescriptiveStats>, I>>(base?: I): DescriptiveStats {
    return DescriptiveStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DescriptiveStats>, I>>(object: I): DescriptiveStats {
    const message = createBaseDescriptiveStats();
    message.min = object.min ?? undefined;
    message.max = object.max ?? undefined;
    message.mean = object.mean ?? undefined;
    message.sd = object.sd ?? undefined;
    return message;
  },
};

function createBaseElementStats(): ElementStats {
  return { elements: {} };
}

export const ElementStats = {
  encode(message: ElementStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.elements || {}).forEach(([key, value]) => {
      ElementStats_ElementsEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ElementStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseElementStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = ElementStats_ElementsEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.elements![entry1.key] = entry1.value;
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

  fromJSON(object: any): ElementStats {
    return {
      elements: isObject(object.elements)
        ? Object.entries(object.elements).reduce<{ [key: string]: DescriptiveStats }>((acc, [key, value]) => {
          acc[key] = DescriptiveStats.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: ElementStats): unknown {
    const obj: any = {};
    if (message.elements) {
      const entries = Object.entries(message.elements);
      if (entries.length > 0) {
        obj.elements = {};
        entries.forEach(([k, v]) => {
          obj.elements[k] = DescriptiveStats.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ElementStats>, I>>(base?: I): ElementStats {
    return ElementStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ElementStats>, I>>(object: I): ElementStats {
    const message = createBaseElementStats();
    message.elements = Object.entries(object.elements ?? {}).reduce<{ [key: string]: DescriptiveStats }>(
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

function createBaseElementStats_ElementsEntry(): ElementStats_ElementsEntry {
  return { key: "", value: undefined };
}

export const ElementStats_ElementsEntry = {
  encode(message: ElementStats_ElementsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      DescriptiveStats.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ElementStats_ElementsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseElementStats_ElementsEntry();
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

  fromJSON(object: any): ElementStats_ElementsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? DescriptiveStats.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: ElementStats_ElementsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = DescriptiveStats.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ElementStats_ElementsEntry>, I>>(base?: I): ElementStats_ElementsEntry {
    return ElementStats_ElementsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ElementStats_ElementsEntry>, I>>(object: I): ElementStats_ElementsEntry {
    const message = createBaseElementStats_ElementsEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? DescriptiveStats.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseTargetStats(): TargetStats {
  return { targets: {} };
}

export const TargetStats = {
  encode(message: TargetStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.targets || {}).forEach(([key, value]) => {
      TargetStats_TargetsEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = TargetStats_TargetsEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.targets![entry1.key] = entry1.value;
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

  fromJSON(object: any): TargetStats {
    return {
      targets: isObject(object.targets)
        ? Object.entries(object.targets).reduce<{ [key: number]: DescriptiveStats }>((acc, [key, value]) => {
          acc[globalThis.Number(key)] = DescriptiveStats.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: TargetStats): unknown {
    const obj: any = {};
    if (message.targets) {
      const entries = Object.entries(message.targets);
      if (entries.length > 0) {
        obj.targets = {};
        entries.forEach(([k, v]) => {
          obj.targets[k] = DescriptiveStats.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TargetStats>, I>>(base?: I): TargetStats {
    return TargetStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TargetStats>, I>>(object: I): TargetStats {
    const message = createBaseTargetStats();
    message.targets = Object.entries(object.targets ?? {}).reduce<{ [key: number]: DescriptiveStats }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[globalThis.Number(key)] = DescriptiveStats.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseTargetStats_TargetsEntry(): TargetStats_TargetsEntry {
  return { key: 0, value: undefined };
}

export const TargetStats_TargetsEntry = {
  encode(message: TargetStats_TargetsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).int32(message.key);
    }
    if (message.value !== undefined) {
      DescriptiveStats.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetStats_TargetsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetStats_TargetsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.key = reader.int32();
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

  fromJSON(object: any): TargetStats_TargetsEntry {
    return {
      key: isSet(object.key) ? globalThis.Number(object.key) : 0,
      value: isSet(object.value) ? DescriptiveStats.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: TargetStats_TargetsEntry): unknown {
    const obj: any = {};
    if (message.key !== 0) {
      obj.key = Math.round(message.key);
    }
    if (message.value !== undefined) {
      obj.value = DescriptiveStats.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TargetStats_TargetsEntry>, I>>(base?: I): TargetStats_TargetsEntry {
    return TargetStats_TargetsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TargetStats_TargetsEntry>, I>>(object: I): TargetStats_TargetsEntry {
    const message = createBaseTargetStats_TargetsEntry();
    message.key = object.key ?? 0;
    message.value = (object.value !== undefined && object.value !== null)
      ? DescriptiveStats.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseSourceStats(): SourceStats {
  return { sources: {} };
}

export const SourceStats = {
  encode(message: SourceStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.sources || {}).forEach(([key, value]) => {
      SourceStats_SourcesEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SourceStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSourceStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = SourceStats_SourcesEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.sources![entry1.key] = entry1.value;
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

  fromJSON(object: any): SourceStats {
    return {
      sources: isObject(object.sources)
        ? Object.entries(object.sources).reduce<{ [key: string]: DescriptiveStats }>((acc, [key, value]) => {
          acc[key] = DescriptiveStats.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: SourceStats): unknown {
    const obj: any = {};
    if (message.sources) {
      const entries = Object.entries(message.sources);
      if (entries.length > 0) {
        obj.sources = {};
        entries.forEach(([k, v]) => {
          obj.sources[k] = DescriptiveStats.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SourceStats>, I>>(base?: I): SourceStats {
    return SourceStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SourceStats>, I>>(object: I): SourceStats {
    const message = createBaseSourceStats();
    message.sources = Object.entries(object.sources ?? {}).reduce<{ [key: string]: DescriptiveStats }>(
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

function createBaseSourceStats_SourcesEntry(): SourceStats_SourcesEntry {
  return { key: "", value: undefined };
}

export const SourceStats_SourcesEntry = {
  encode(message: SourceStats_SourcesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      DescriptiveStats.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SourceStats_SourcesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSourceStats_SourcesEntry();
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

  fromJSON(object: any): SourceStats_SourcesEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? DescriptiveStats.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: SourceStats_SourcesEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = DescriptiveStats.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SourceStats_SourcesEntry>, I>>(base?: I): SourceStats_SourcesEntry {
    return SourceStats_SourcesEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SourceStats_SourcesEntry>, I>>(object: I): SourceStats_SourcesEntry {
    const message = createBaseSourceStats_SourcesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? DescriptiveStats.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseBucketStats(): BucketStats {
  return { bucket_size: 0, buckets: [] };
}

export const BucketStats = {
  encode(message: BucketStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bucket_size !== undefined && message.bucket_size !== 0) {
      writer.uint32(8).uint32(message.bucket_size);
    }
    if (message.buckets !== undefined && message.buckets.length !== 0) {
      for (const v of message.buckets) {
        DescriptiveStats.encode(v!, writer.uint32(18).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BucketStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBucketStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.bucket_size = reader.uint32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.buckets!.push(DescriptiveStats.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): BucketStats {
    return {
      bucket_size: isSet(object.bucket_size) ? globalThis.Number(object.bucket_size) : 0,
      buckets: globalThis.Array.isArray(object?.buckets)
        ? object.buckets.map((e: any) => DescriptiveStats.fromJSON(e))
        : [],
    };
  },

  toJSON(message: BucketStats): unknown {
    const obj: any = {};
    if (message.bucket_size !== undefined && message.bucket_size !== 0) {
      obj.bucket_size = Math.round(message.bucket_size);
    }
    if (message.buckets?.length) {
      obj.buckets = message.buckets.map((e) => DescriptiveStats.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<BucketStats>, I>>(base?: I): BucketStats {
    return BucketStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<BucketStats>, I>>(object: I): BucketStats {
    const message = createBaseBucketStats();
    message.bucket_size = object.bucket_size ?? 0;
    message.buckets = object.buckets?.map((e) => DescriptiveStats.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCharacterBucketStats(): CharacterBucketStats {
  return { bucket_size: 0, characters: [] };
}

export const CharacterBucketStats = {
  encode(message: CharacterBucketStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bucket_size !== undefined && message.bucket_size !== 0) {
      writer.uint32(8).uint32(message.bucket_size);
    }
    if (message.characters !== undefined && message.characters.length !== 0) {
      for (const v of message.characters) {
        CharacterBuckets.encode(v!, writer.uint32(18).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CharacterBucketStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCharacterBucketStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.bucket_size = reader.uint32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.characters!.push(CharacterBuckets.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CharacterBucketStats {
    return {
      bucket_size: isSet(object.bucket_size) ? globalThis.Number(object.bucket_size) : 0,
      characters: globalThis.Array.isArray(object?.characters)
        ? object.characters.map((e: any) => CharacterBuckets.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CharacterBucketStats): unknown {
    const obj: any = {};
    if (message.bucket_size !== undefined && message.bucket_size !== 0) {
      obj.bucket_size = Math.round(message.bucket_size);
    }
    if (message.characters?.length) {
      obj.characters = message.characters.map((e) => CharacterBuckets.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CharacterBucketStats>, I>>(base?: I): CharacterBucketStats {
    return CharacterBucketStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CharacterBucketStats>, I>>(object: I): CharacterBucketStats {
    const message = createBaseCharacterBucketStats();
    message.bucket_size = object.bucket_size ?? 0;
    message.characters = object.characters?.map((e) => CharacterBuckets.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCharacterBuckets(): CharacterBuckets {
  return { buckets: [] };
}

export const CharacterBuckets = {
  encode(message: CharacterBuckets, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.buckets !== undefined && message.buckets.length !== 0) {
      for (const v of message.buckets) {
        DescriptiveStats.encode(v!, writer.uint32(10).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CharacterBuckets {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCharacterBuckets();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.buckets!.push(DescriptiveStats.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CharacterBuckets {
    return {
      buckets: globalThis.Array.isArray(object?.buckets)
        ? object.buckets.map((e: any) => DescriptiveStats.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CharacterBuckets): unknown {
    const obj: any = {};
    if (message.buckets?.length) {
      obj.buckets = message.buckets.map((e) => DescriptiveStats.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CharacterBuckets>, I>>(base?: I): CharacterBuckets {
    return CharacterBuckets.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CharacterBuckets>, I>>(object: I): CharacterBuckets {
    const message = createBaseCharacterBuckets();
    message.buckets = object.buckets?.map((e) => DescriptiveStats.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTargetBucketStats(): TargetBucketStats {
  return { bucket_size: 0, targets: {} };
}

export const TargetBucketStats = {
  encode(message: TargetBucketStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bucket_size !== undefined && message.bucket_size !== 0) {
      writer.uint32(8).uint32(message.bucket_size);
    }
    Object.entries(message.targets || {}).forEach(([key, value]) => {
      TargetBucketStats_TargetsEntry.encode({ key: key as any, value }, writer.uint32(18).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetBucketStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetBucketStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.bucket_size = reader.uint32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          const entry2 = TargetBucketStats_TargetsEntry.decode(reader, reader.uint32());
          if (entry2.value !== undefined) {
            message.targets![entry2.key] = entry2.value;
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

  fromJSON(object: any): TargetBucketStats {
    return {
      bucket_size: isSet(object.bucket_size) ? globalThis.Number(object.bucket_size) : 0,
      targets: isObject(object.targets)
        ? Object.entries(object.targets).reduce<{ [key: number]: TargetBuckets }>((acc, [key, value]) => {
          acc[globalThis.Number(key)] = TargetBuckets.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: TargetBucketStats): unknown {
    const obj: any = {};
    if (message.bucket_size !== undefined && message.bucket_size !== 0) {
      obj.bucket_size = Math.round(message.bucket_size);
    }
    if (message.targets) {
      const entries = Object.entries(message.targets);
      if (entries.length > 0) {
        obj.targets = {};
        entries.forEach(([k, v]) => {
          obj.targets[k] = TargetBuckets.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TargetBucketStats>, I>>(base?: I): TargetBucketStats {
    return TargetBucketStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TargetBucketStats>, I>>(object: I): TargetBucketStats {
    const message = createBaseTargetBucketStats();
    message.bucket_size = object.bucket_size ?? 0;
    message.targets = Object.entries(object.targets ?? {}).reduce<{ [key: number]: TargetBuckets }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[globalThis.Number(key)] = TargetBuckets.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseTargetBucketStats_TargetsEntry(): TargetBucketStats_TargetsEntry {
  return { key: 0, value: undefined };
}

export const TargetBucketStats_TargetsEntry = {
  encode(message: TargetBucketStats_TargetsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).int32(message.key);
    }
    if (message.value !== undefined) {
      TargetBuckets.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetBucketStats_TargetsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetBucketStats_TargetsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.key = reader.int32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = TargetBuckets.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TargetBucketStats_TargetsEntry {
    return {
      key: isSet(object.key) ? globalThis.Number(object.key) : 0,
      value: isSet(object.value) ? TargetBuckets.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: TargetBucketStats_TargetsEntry): unknown {
    const obj: any = {};
    if (message.key !== 0) {
      obj.key = Math.round(message.key);
    }
    if (message.value !== undefined) {
      obj.value = TargetBuckets.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TargetBucketStats_TargetsEntry>, I>>(base?: I): TargetBucketStats_TargetsEntry {
    return TargetBucketStats_TargetsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TargetBucketStats_TargetsEntry>, I>>(
    object: I,
  ): TargetBucketStats_TargetsEntry {
    const message = createBaseTargetBucketStats_TargetsEntry();
    message.key = object.key ?? 0;
    message.value = (object.value !== undefined && object.value !== null)
      ? TargetBuckets.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseTargetBuckets(): TargetBuckets {
  return { overall: undefined, target: undefined };
}

export const TargetBuckets = {
  encode(message: TargetBuckets, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.overall !== undefined) {
      TargetBucket.encode(message.overall, writer.uint32(10).fork()).ldelim();
    }
    if (message.target !== undefined) {
      TargetBucket.encode(message.target, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetBuckets {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetBuckets();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.overall = TargetBucket.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.target = TargetBucket.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TargetBuckets {
    return {
      overall: isSet(object.overall) ? TargetBucket.fromJSON(object.overall) : undefined,
      target: isSet(object.target) ? TargetBucket.fromJSON(object.target) : undefined,
    };
  },

  toJSON(message: TargetBuckets): unknown {
    const obj: any = {};
    if (message.overall !== undefined) {
      obj.overall = TargetBucket.toJSON(message.overall);
    }
    if (message.target !== undefined) {
      obj.target = TargetBucket.toJSON(message.target);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TargetBuckets>, I>>(base?: I): TargetBuckets {
    return TargetBuckets.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TargetBuckets>, I>>(object: I): TargetBuckets {
    const message = createBaseTargetBuckets();
    message.overall = (object.overall !== undefined && object.overall !== null)
      ? TargetBucket.fromPartial(object.overall)
      : undefined;
    message.target = (object.target !== undefined && object.target !== null)
      ? TargetBucket.fromPartial(object.target)
      : undefined;
    return message;
  },
};

function createBaseTargetBucket(): TargetBucket {
  return { min: [], max: [], q1: [], q2: [], q3: [] };
}

export const TargetBucket = {
  encode(message: TargetBucket, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.min !== undefined && message.min.length !== 0) {
      writer.uint32(10).fork();
      for (const v of message.min) {
        writer.double(v);
      }
      writer.ldelim();
    }
    if (message.max !== undefined && message.max.length !== 0) {
      writer.uint32(18).fork();
      for (const v of message.max) {
        writer.double(v);
      }
      writer.ldelim();
    }
    if (message.q1 !== undefined && message.q1.length !== 0) {
      writer.uint32(26).fork();
      for (const v of message.q1) {
        writer.double(v);
      }
      writer.ldelim();
    }
    if (message.q2 !== undefined && message.q2.length !== 0) {
      writer.uint32(34).fork();
      for (const v of message.q2) {
        writer.double(v);
      }
      writer.ldelim();
    }
    if (message.q3 !== undefined && message.q3.length !== 0) {
      writer.uint32(42).fork();
      for (const v of message.q3) {
        writer.double(v);
      }
      writer.ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetBucket {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetBucket();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag === 9) {
            message.min!.push(reader.double());

            continue;
          }

          if (tag === 10) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.min!.push(reader.double());
            }

            continue;
          }

          break;
        case 2:
          if (tag === 17) {
            message.max!.push(reader.double());

            continue;
          }

          if (tag === 18) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.max!.push(reader.double());
            }

            continue;
          }

          break;
        case 3:
          if (tag === 25) {
            message.q1!.push(reader.double());

            continue;
          }

          if (tag === 26) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.q1!.push(reader.double());
            }

            continue;
          }

          break;
        case 4:
          if (tag === 33) {
            message.q2!.push(reader.double());

            continue;
          }

          if (tag === 34) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.q2!.push(reader.double());
            }

            continue;
          }

          break;
        case 5:
          if (tag === 41) {
            message.q3!.push(reader.double());

            continue;
          }

          if (tag === 42) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.q3!.push(reader.double());
            }

            continue;
          }

          break;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TargetBucket {
    return {
      min: globalThis.Array.isArray(object?.min) ? object.min.map((e: any) => globalThis.Number(e)) : [],
      max: globalThis.Array.isArray(object?.max) ? object.max.map((e: any) => globalThis.Number(e)) : [],
      q1: globalThis.Array.isArray(object?.q1) ? object.q1.map((e: any) => globalThis.Number(e)) : [],
      q2: globalThis.Array.isArray(object?.q2) ? object.q2.map((e: any) => globalThis.Number(e)) : [],
      q3: globalThis.Array.isArray(object?.q3) ? object.q3.map((e: any) => globalThis.Number(e)) : [],
    };
  },

  toJSON(message: TargetBucket): unknown {
    const obj: any = {};
    if (message.min?.length) {
      obj.min = message.min;
    }
    if (message.max?.length) {
      obj.max = message.max;
    }
    if (message.q1?.length) {
      obj.q1 = message.q1;
    }
    if (message.q2?.length) {
      obj.q2 = message.q2;
    }
    if (message.q3?.length) {
      obj.q3 = message.q3;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TargetBucket>, I>>(base?: I): TargetBucket {
    return TargetBucket.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TargetBucket>, I>>(object: I): TargetBucket {
    const message = createBaseTargetBucket();
    message.min = object.min?.map((e) => e) || [];
    message.max = object.max?.map((e) => e) || [];
    message.q1 = object.q1?.map((e) => e) || [];
    message.q2 = object.q2?.map((e) => e) || [];
    message.q3 = object.q3?.map((e) => e) || [];
    return message;
  },
};

function createBaseWarnings(): Warnings {
  return {
    target_overlap: false,
    insufficient_energy: false,
    insufficient_stamina: false,
    swap_cd: false,
    skill_cd: false,
    dash_cd: false,
  };
}

export const Warnings = {
  encode(message: Warnings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.target_overlap !== undefined && message.target_overlap !== false) {
      writer.uint32(8).bool(message.target_overlap);
    }
    if (message.insufficient_energy !== undefined && message.insufficient_energy !== false) {
      writer.uint32(16).bool(message.insufficient_energy);
    }
    if (message.insufficient_stamina !== undefined && message.insufficient_stamina !== false) {
      writer.uint32(24).bool(message.insufficient_stamina);
    }
    if (message.swap_cd !== undefined && message.swap_cd !== false) {
      writer.uint32(32).bool(message.swap_cd);
    }
    if (message.skill_cd !== undefined && message.skill_cd !== false) {
      writer.uint32(40).bool(message.skill_cd);
    }
    if (message.dash_cd !== undefined && message.dash_cd !== false) {
      writer.uint32(48).bool(message.dash_cd);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Warnings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWarnings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.target_overlap = reader.bool();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.insufficient_energy = reader.bool();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.insufficient_stamina = reader.bool();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.swap_cd = reader.bool();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.skill_cd = reader.bool();
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.dash_cd = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Warnings {
    return {
      target_overlap: isSet(object.target_overlap) ? globalThis.Boolean(object.target_overlap) : false,
      insufficient_energy: isSet(object.insufficient_energy) ? globalThis.Boolean(object.insufficient_energy) : false,
      insufficient_stamina: isSet(object.insufficient_stamina)
        ? globalThis.Boolean(object.insufficient_stamina)
        : false,
      swap_cd: isSet(object.swap_cd) ? globalThis.Boolean(object.swap_cd) : false,
      skill_cd: isSet(object.skill_cd) ? globalThis.Boolean(object.skill_cd) : false,
      dash_cd: isSet(object.dash_cd) ? globalThis.Boolean(object.dash_cd) : false,
    };
  },

  toJSON(message: Warnings): unknown {
    const obj: any = {};
    if (message.target_overlap !== undefined && message.target_overlap !== false) {
      obj.target_overlap = message.target_overlap;
    }
    if (message.insufficient_energy !== undefined && message.insufficient_energy !== false) {
      obj.insufficient_energy = message.insufficient_energy;
    }
    if (message.insufficient_stamina !== undefined && message.insufficient_stamina !== false) {
      obj.insufficient_stamina = message.insufficient_stamina;
    }
    if (message.swap_cd !== undefined && message.swap_cd !== false) {
      obj.swap_cd = message.swap_cd;
    }
    if (message.skill_cd !== undefined && message.skill_cd !== false) {
      obj.skill_cd = message.skill_cd;
    }
    if (message.dash_cd !== undefined && message.dash_cd !== false) {
      obj.dash_cd = message.dash_cd;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Warnings>, I>>(base?: I): Warnings {
    return Warnings.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Warnings>, I>>(object: I): Warnings {
    const message = createBaseWarnings();
    message.target_overlap = object.target_overlap ?? false;
    message.insufficient_energy = object.insufficient_energy ?? false;
    message.insufficient_stamina = object.insufficient_stamina ?? false;
    message.swap_cd = object.swap_cd ?? false;
    message.skill_cd = object.skill_cd ?? false;
    message.dash_cd = object.dash_cd ?? false;
    return message;
  },
};

function createBaseFailedActions(): FailedActions {
  return {
    insufficient_energy: undefined,
    insufficient_stamina: undefined,
    swap_cd: undefined,
    skill_cd: undefined,
    dash_cd: undefined,
  };
}

export const FailedActions = {
  encode(message: FailedActions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.insufficient_energy !== undefined) {
      DescriptiveStats.encode(message.insufficient_energy, writer.uint32(10).fork()).ldelim();
    }
    if (message.insufficient_stamina !== undefined) {
      DescriptiveStats.encode(message.insufficient_stamina, writer.uint32(18).fork()).ldelim();
    }
    if (message.swap_cd !== undefined) {
      DescriptiveStats.encode(message.swap_cd, writer.uint32(26).fork()).ldelim();
    }
    if (message.skill_cd !== undefined) {
      DescriptiveStats.encode(message.skill_cd, writer.uint32(34).fork()).ldelim();
    }
    if (message.dash_cd !== undefined) {
      DescriptiveStats.encode(message.dash_cd, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FailedActions {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFailedActions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.insufficient_energy = DescriptiveStats.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.insufficient_stamina = DescriptiveStats.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.swap_cd = DescriptiveStats.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.skill_cd = DescriptiveStats.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.dash_cd = DescriptiveStats.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FailedActions {
    return {
      insufficient_energy: isSet(object.insufficient_energy)
        ? DescriptiveStats.fromJSON(object.insufficient_energy)
        : undefined,
      insufficient_stamina: isSet(object.insufficient_stamina)
        ? DescriptiveStats.fromJSON(object.insufficient_stamina)
        : undefined,
      swap_cd: isSet(object.swap_cd) ? DescriptiveStats.fromJSON(object.swap_cd) : undefined,
      skill_cd: isSet(object.skill_cd) ? DescriptiveStats.fromJSON(object.skill_cd) : undefined,
      dash_cd: isSet(object.dash_cd) ? DescriptiveStats.fromJSON(object.dash_cd) : undefined,
    };
  },

  toJSON(message: FailedActions): unknown {
    const obj: any = {};
    if (message.insufficient_energy !== undefined) {
      obj.insufficient_energy = DescriptiveStats.toJSON(message.insufficient_energy);
    }
    if (message.insufficient_stamina !== undefined) {
      obj.insufficient_stamina = DescriptiveStats.toJSON(message.insufficient_stamina);
    }
    if (message.swap_cd !== undefined) {
      obj.swap_cd = DescriptiveStats.toJSON(message.swap_cd);
    }
    if (message.skill_cd !== undefined) {
      obj.skill_cd = DescriptiveStats.toJSON(message.skill_cd);
    }
    if (message.dash_cd !== undefined) {
      obj.dash_cd = DescriptiveStats.toJSON(message.dash_cd);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FailedActions>, I>>(base?: I): FailedActions {
    return FailedActions.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FailedActions>, I>>(object: I): FailedActions {
    const message = createBaseFailedActions();
    message.insufficient_energy = (object.insufficient_energy !== undefined && object.insufficient_energy !== null)
      ? DescriptiveStats.fromPartial(object.insufficient_energy)
      : undefined;
    message.insufficient_stamina = (object.insufficient_stamina !== undefined && object.insufficient_stamina !== null)
      ? DescriptiveStats.fromPartial(object.insufficient_stamina)
      : undefined;
    message.swap_cd = (object.swap_cd !== undefined && object.swap_cd !== null)
      ? DescriptiveStats.fromPartial(object.swap_cd)
      : undefined;
    message.skill_cd = (object.skill_cd !== undefined && object.skill_cd !== null)
      ? DescriptiveStats.fromPartial(object.skill_cd)
      : undefined;
    message.dash_cd = (object.dash_cd !== undefined && object.dash_cd !== null)
      ? DescriptiveStats.fromPartial(object.dash_cd)
      : undefined;
    return message;
  },
};

function createBaseShieldInfo(): ShieldInfo {
  return { hp: {}, uptime: undefined };
}

export const ShieldInfo = {
  encode(message: ShieldInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.hp || {}).forEach(([key, value]) => {
      ShieldInfo_HpEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    if (message.uptime !== undefined) {
      DescriptiveStats.encode(message.uptime, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ShieldInfo {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseShieldInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = ShieldInfo_HpEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.hp![entry1.key] = entry1.value;
          }
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.uptime = DescriptiveStats.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ShieldInfo {
    return {
      hp: isObject(object.hp)
        ? Object.entries(object.hp).reduce<{ [key: string]: DescriptiveStats }>((acc, [key, value]) => {
          acc[key] = DescriptiveStats.fromJSON(value);
          return acc;
        }, {})
        : {},
      uptime: isSet(object.uptime) ? DescriptiveStats.fromJSON(object.uptime) : undefined,
    };
  },

  toJSON(message: ShieldInfo): unknown {
    const obj: any = {};
    if (message.hp) {
      const entries = Object.entries(message.hp);
      if (entries.length > 0) {
        obj.hp = {};
        entries.forEach(([k, v]) => {
          obj.hp[k] = DescriptiveStats.toJSON(v);
        });
      }
    }
    if (message.uptime !== undefined) {
      obj.uptime = DescriptiveStats.toJSON(message.uptime);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ShieldInfo>, I>>(base?: I): ShieldInfo {
    return ShieldInfo.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ShieldInfo>, I>>(object: I): ShieldInfo {
    const message = createBaseShieldInfo();
    message.hp = Object.entries(object.hp ?? {}).reduce<{ [key: string]: DescriptiveStats }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = DescriptiveStats.fromPartial(value);
      }
      return acc;
    }, {});
    message.uptime = (object.uptime !== undefined && object.uptime !== null)
      ? DescriptiveStats.fromPartial(object.uptime)
      : undefined;
    return message;
  },
};

function createBaseShieldInfo_HpEntry(): ShieldInfo_HpEntry {
  return { key: "", value: undefined };
}

export const ShieldInfo_HpEntry = {
  encode(message: ShieldInfo_HpEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      DescriptiveStats.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ShieldInfo_HpEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseShieldInfo_HpEntry();
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

  fromJSON(object: any): ShieldInfo_HpEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? DescriptiveStats.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: ShieldInfo_HpEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = DescriptiveStats.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ShieldInfo_HpEntry>, I>>(base?: I): ShieldInfo_HpEntry {
    return ShieldInfo_HpEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ShieldInfo_HpEntry>, I>>(object: I): ShieldInfo_HpEntry {
    const message = createBaseShieldInfo_HpEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? DescriptiveStats.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseEndStats(): EndStats {
  return { ending_energy: undefined };
}

export const EndStats = {
  encode(message: EndStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ending_energy !== undefined) {
      DescriptiveStats.encode(message.ending_energy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EndStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEndStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.ending_energy = DescriptiveStats.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EndStats {
    return { ending_energy: isSet(object.ending_energy) ? DescriptiveStats.fromJSON(object.ending_energy) : undefined };
  },

  toJSON(message: EndStats): unknown {
    const obj: any = {};
    if (message.ending_energy !== undefined) {
      obj.ending_energy = DescriptiveStats.toJSON(message.ending_energy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EndStats>, I>>(base?: I): EndStats {
    return EndStats.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EndStats>, I>>(object: I): EndStats {
    const message = createBaseEndStats();
    message.ending_energy = (object.ending_energy !== undefined && object.ending_energy !== null)
      ? DescriptiveStats.fromPartial(object.ending_energy)
      : undefined;
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
