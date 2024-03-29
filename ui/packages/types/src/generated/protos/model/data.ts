/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import {
  AvatarCurveType,
  avatarCurveTypeFromJSON,
  avatarCurveTypeToJSON,
  BodyType,
  bodyTypeFromJSON,
  bodyTypeToJSON,
  Element,
  elementFromJSON,
  elementToJSON,
  MonsterCurveType,
  monsterCurveTypeFromJSON,
  monsterCurveTypeToJSON,
  QualityType,
  qualityTypeFromJSON,
  qualityTypeToJSON,
  StatType,
  statTypeFromJSON,
  statTypeToJSON,
  WeaponClass,
  weaponClassFromJSON,
  weaponClassToJSON,
  WeaponCurveType,
  weaponCurveTypeFromJSON,
  weaponCurveTypeToJSON,
  ZoneType,
  zoneTypeFromJSON,
  zoneTypeToJSON,
} from "./enums";

export interface AvatarDataMap {
  data?: { [key: string]: AvatarData } | undefined;
}

export interface AvatarDataMap_DataEntry {
  key: string;
  value?: AvatarData | undefined;
}

export interface AvatarData {
  id?: number | undefined;
  sub_id?: number | undefined;
  key?: string | undefined;
  rarity?: QualityType | undefined;
  body?: BodyType | undefined;
  region?: ZoneType | undefined;
  element?: Element | undefined;
  weapon_class?: WeaponClass | undefined;
  icon_name?: string | undefined;
  stats?: AvatarStatsData | undefined;
  skill_details?: AvatarSkillsData | undefined;
  "name_text_hash_map "?: number | undefined;
}

export interface AvatarStatsData {
  /**
   * TODO: base stat should be refactor to just an array of stats
   * there is no requirement that base stat can only be 3 stats; in fact
   * ER/cr/cd can be considered as base
   */
  base_hp?: number | undefined;
  base_atk?: number | undefined;
  base_def?: number | undefined;
  hp_curve?: AvatarCurveType | undefined;
  atk_curve?: AvatarCurveType | undefined;
  def_cruve?: AvatarCurveType | undefined;
  promo_data?: PromotionData[] | undefined;
}

export interface AvatarSkillsData {
  skill?: number | undefined;
  burst?: number | undefined;
  attack?: number | undefined;
  burst_energy_cost?: number | undefined;
  attack_scaling?: AvatarSkillExcelIndexData[] | undefined;
  skill_scaling?: AvatarSkillExcelIndexData[] | undefined;
  burst_scaling?: AvatarSkillExcelIndexData[] | undefined;
}

export interface AvatarSkillExcelIndexData {
  group_id?:
    | number
    | undefined;
  /** position in the param list */
  index?: number | undefined;
  level_data?: AvatarSkillExcelLevelData[] | undefined;
}

export interface AvatarSkillExcelLevelData {
  level?: number | undefined;
  value?: number | undefined;
}

export interface WeaponDataMap {
  data?: { [key: string]: WeaponData } | undefined;
}

export interface WeaponDataMap_DataEntry {
  key: string;
  value?: WeaponData | undefined;
}

export interface WeaponData {
  id?: number | undefined;
  key?:
    | string
    | undefined;
  /** for whatever reason weapon rarity is a number */
  rarity?: number | undefined;
  weapon_class?: WeaponClass | undefined;
  image_name?: string | undefined;
  base_stats?: WeaponStatsData | undefined;
  "name_text_hash_map "?: number | undefined;
}

export interface WeaponStatsData {
  base_props?: WeaponProp[] | undefined;
  promo_data?: PromotionData[] | undefined;
}

export interface WeaponProp {
  prop_type?: StatType | undefined;
  initial_value?: number | undefined;
  curve?: WeaponCurveType | undefined;
}

export interface ArtifactDataMap {
  data?: { [key: string]: ArtifactData } | undefined;
}

export interface ArtifactDataMap_DataEntry {
  key: string;
  value?: ArtifactData | undefined;
}

export interface ArtifactData {
  id?: number | undefined;
  text_map_id?: number | undefined;
  key?: string | undefined;
}

export interface PromotionData {
  max_level?: number | undefined;
  add_props?: PromotionAddProp[] | undefined;
}

export interface PromotionAddProp {
  prop_type?: StatType | undefined;
  value?: number | undefined;
}

export interface MonsterData {
  id?: number | undefined;
  key?: string | undefined;
  base_stats?: MonsterStatsData | undefined;
  name_text_hash_map?: number | undefined;
}

export interface MonsterStatsData {
  base_hp?: number | undefined;
  hp_curve?: MonsterCurveType | undefined;
  resist?: MonsterResistData | undefined;
  freeze_resist?: number | undefined;
  hp_drop?: MonsterHPDrop[] | undefined;
}

export interface MonsterResistData {
  fire_resist?: number | undefined;
  grass_resist?: number | undefined;
  water_resist?: number | undefined;
  electric_resist?: number | undefined;
  wind_resist?: number | undefined;
  ice_resist?: number | undefined;
  rock_resist?: number | undefined;
  physical_resist?: number | undefined;
}

export interface MonsterHPDrop {
  drop_id?: number | undefined;
  hp_percent?: number | undefined;
}

function createBaseAvatarDataMap(): AvatarDataMap {
  return { data: {} };
}

export const AvatarDataMap = {
  encode(message: AvatarDataMap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.data || {}).forEach(([key, value]) => {
      AvatarDataMap_DataEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AvatarDataMap {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAvatarDataMap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = AvatarDataMap_DataEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.data![entry1.key] = entry1.value;
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

  fromJSON(object: any): AvatarDataMap {
    return {
      data: isObject(object.data)
        ? Object.entries(object.data).reduce<{ [key: string]: AvatarData }>((acc, [key, value]) => {
          acc[key] = AvatarData.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: AvatarDataMap): unknown {
    const obj: any = {};
    if (message.data) {
      const entries = Object.entries(message.data);
      if (entries.length > 0) {
        obj.data = {};
        entries.forEach(([k, v]) => {
          obj.data[k] = AvatarData.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AvatarDataMap>, I>>(base?: I): AvatarDataMap {
    return AvatarDataMap.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AvatarDataMap>, I>>(object: I): AvatarDataMap {
    const message = createBaseAvatarDataMap();
    message.data = Object.entries(object.data ?? {}).reduce<{ [key: string]: AvatarData }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = AvatarData.fromPartial(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseAvatarDataMap_DataEntry(): AvatarDataMap_DataEntry {
  return { key: "", value: undefined };
}

export const AvatarDataMap_DataEntry = {
  encode(message: AvatarDataMap_DataEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      AvatarData.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AvatarDataMap_DataEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAvatarDataMap_DataEntry();
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

          message.value = AvatarData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AvatarDataMap_DataEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? AvatarData.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: AvatarDataMap_DataEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = AvatarData.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AvatarDataMap_DataEntry>, I>>(base?: I): AvatarDataMap_DataEntry {
    return AvatarDataMap_DataEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AvatarDataMap_DataEntry>, I>>(object: I): AvatarDataMap_DataEntry {
    const message = createBaseAvatarDataMap_DataEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? AvatarData.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseAvatarData(): AvatarData {
  return {
    id: 0,
    sub_id: 0,
    key: "",
    rarity: 0,
    body: 0,
    region: 0,
    element: 0,
    weapon_class: 0,
    icon_name: "",
    stats: undefined,
    skill_details: undefined,
    "name_text_hash_map ": 0,
  };
}

export const AvatarData = {
  encode(message: AvatarData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== 0) {
      writer.uint32(8).int32(message.id);
    }
    if (message.sub_id !== undefined && message.sub_id !== 0) {
      writer.uint32(16).int32(message.sub_id);
    }
    if (message.key !== undefined && message.key !== "") {
      writer.uint32(26).string(message.key);
    }
    if (message.rarity !== undefined && message.rarity !== 0) {
      writer.uint32(32).int32(message.rarity);
    }
    if (message.body !== undefined && message.body !== 0) {
      writer.uint32(40).int32(message.body);
    }
    if (message.region !== undefined && message.region !== 0) {
      writer.uint32(48).int32(message.region);
    }
    if (message.element !== undefined && message.element !== 0) {
      writer.uint32(56).int32(message.element);
    }
    if (message.weapon_class !== undefined && message.weapon_class !== 0) {
      writer.uint32(64).int32(message.weapon_class);
    }
    if (message.icon_name !== undefined && message.icon_name !== "") {
      writer.uint32(74).string(message.icon_name);
    }
    if (message.stats !== undefined) {
      AvatarStatsData.encode(message.stats, writer.uint32(82).fork()).ldelim();
    }
    if (message.skill_details !== undefined) {
      AvatarSkillsData.encode(message.skill_details, writer.uint32(90).fork()).ldelim();
    }
    if (message["name_text_hash_map "] !== undefined && message["name_text_hash_map "] !== 0) {
      writer.uint32(96).int64(message["name_text_hash_map "]);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AvatarData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAvatarData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = reader.int32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.sub_id = reader.int32();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.key = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.rarity = reader.int32() as any;
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.body = reader.int32() as any;
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.region = reader.int32() as any;
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.element = reader.int32() as any;
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.weapon_class = reader.int32() as any;
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.icon_name = reader.string();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.stats = AvatarStatsData.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.skill_details = AvatarSkillsData.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag !== 96) {
            break;
          }

          message["name_text_hash_map "] = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AvatarData {
    return {
      id: isSet(object.id) ? globalThis.Number(object.id) : 0,
      sub_id: isSet(object.sub_id) ? globalThis.Number(object.sub_id) : 0,
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      rarity: isSet(object.rarity) ? qualityTypeFromJSON(object.rarity) : 0,
      body: isSet(object.body) ? bodyTypeFromJSON(object.body) : 0,
      region: isSet(object.region) ? zoneTypeFromJSON(object.region) : 0,
      element: isSet(object.element) ? elementFromJSON(object.element) : 0,
      weapon_class: isSet(object.weapon_class) ? weaponClassFromJSON(object.weapon_class) : 0,
      icon_name: isSet(object.icon_name) ? globalThis.String(object.icon_name) : "",
      stats: isSet(object.stats) ? AvatarStatsData.fromJSON(object.stats) : undefined,
      skill_details: isSet(object.skill_details) ? AvatarSkillsData.fromJSON(object.skill_details) : undefined,
      "name_text_hash_map ": isSet(object["name_text_hash_map "])
        ? globalThis.Number(object["name_text_hash_map "])
        : 0,
    };
  },

  toJSON(message: AvatarData): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.sub_id !== undefined && message.sub_id !== 0) {
      obj.sub_id = Math.round(message.sub_id);
    }
    if (message.key !== undefined && message.key !== "") {
      obj.key = message.key;
    }
    if (message.rarity !== undefined && message.rarity !== 0) {
      obj.rarity = qualityTypeToJSON(message.rarity);
    }
    if (message.body !== undefined && message.body !== 0) {
      obj.body = bodyTypeToJSON(message.body);
    }
    if (message.region !== undefined && message.region !== 0) {
      obj.region = zoneTypeToJSON(message.region);
    }
    if (message.element !== undefined && message.element !== 0) {
      obj.element = elementToJSON(message.element);
    }
    if (message.weapon_class !== undefined && message.weapon_class !== 0) {
      obj.weapon_class = weaponClassToJSON(message.weapon_class);
    }
    if (message.icon_name !== undefined && message.icon_name !== "") {
      obj.icon_name = message.icon_name;
    }
    if (message.stats !== undefined) {
      obj.stats = AvatarStatsData.toJSON(message.stats);
    }
    if (message.skill_details !== undefined) {
      obj.skill_details = AvatarSkillsData.toJSON(message.skill_details);
    }
    if (message["name_text_hash_map "] !== undefined && message["name_text_hash_map "] !== 0) {
      obj["name_text_hash_map "] = Math.round(message["name_text_hash_map "]);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AvatarData>, I>>(base?: I): AvatarData {
    return AvatarData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AvatarData>, I>>(object: I): AvatarData {
    const message = createBaseAvatarData();
    message.id = object.id ?? 0;
    message.sub_id = object.sub_id ?? 0;
    message.key = object.key ?? "";
    message.rarity = object.rarity ?? 0;
    message.body = object.body ?? 0;
    message.region = object.region ?? 0;
    message.element = object.element ?? 0;
    message.weapon_class = object.weapon_class ?? 0;
    message.icon_name = object.icon_name ?? "";
    message.stats = (object.stats !== undefined && object.stats !== null)
      ? AvatarStatsData.fromPartial(object.stats)
      : undefined;
    message.skill_details = (object.skill_details !== undefined && object.skill_details !== null)
      ? AvatarSkillsData.fromPartial(object.skill_details)
      : undefined;
    message["name_text_hash_map "] = object["name_text_hash_map "] ?? 0;
    return message;
  },
};

function createBaseAvatarStatsData(): AvatarStatsData {
  return { base_hp: 0, base_atk: 0, base_def: 0, hp_curve: 0, atk_curve: 0, def_cruve: 0, promo_data: [] };
}

export const AvatarStatsData = {
  encode(message: AvatarStatsData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.base_hp !== undefined && message.base_hp !== 0) {
      writer.uint32(9).double(message.base_hp);
    }
    if (message.base_atk !== undefined && message.base_atk !== 0) {
      writer.uint32(17).double(message.base_atk);
    }
    if (message.base_def !== undefined && message.base_def !== 0) {
      writer.uint32(25).double(message.base_def);
    }
    if (message.hp_curve !== undefined && message.hp_curve !== 0) {
      writer.uint32(32).int32(message.hp_curve);
    }
    if (message.atk_curve !== undefined && message.atk_curve !== 0) {
      writer.uint32(40).int32(message.atk_curve);
    }
    if (message.def_cruve !== undefined && message.def_cruve !== 0) {
      writer.uint32(48).int32(message.def_cruve);
    }
    if (message.promo_data !== undefined && message.promo_data.length !== 0) {
      for (const v of message.promo_data) {
        PromotionData.encode(v!, writer.uint32(58).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AvatarStatsData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAvatarStatsData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.base_hp = reader.double();
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.base_atk = reader.double();
          continue;
        case 3:
          if (tag !== 25) {
            break;
          }

          message.base_def = reader.double();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.hp_curve = reader.int32() as any;
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.atk_curve = reader.int32() as any;
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.def_cruve = reader.int32() as any;
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.promo_data!.push(PromotionData.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AvatarStatsData {
    return {
      base_hp: isSet(object.base_hp) ? globalThis.Number(object.base_hp) : 0,
      base_atk: isSet(object.base_atk) ? globalThis.Number(object.base_atk) : 0,
      base_def: isSet(object.base_def) ? globalThis.Number(object.base_def) : 0,
      hp_curve: isSet(object.hp_curve) ? avatarCurveTypeFromJSON(object.hp_curve) : 0,
      atk_curve: isSet(object.atk_curve) ? avatarCurveTypeFromJSON(object.atk_curve) : 0,
      def_cruve: isSet(object.def_cruve) ? avatarCurveTypeFromJSON(object.def_cruve) : 0,
      promo_data: globalThis.Array.isArray(object?.promo_data)
        ? object.promo_data.map((e: any) => PromotionData.fromJSON(e))
        : [],
    };
  },

  toJSON(message: AvatarStatsData): unknown {
    const obj: any = {};
    if (message.base_hp !== undefined && message.base_hp !== 0) {
      obj.base_hp = message.base_hp;
    }
    if (message.base_atk !== undefined && message.base_atk !== 0) {
      obj.base_atk = message.base_atk;
    }
    if (message.base_def !== undefined && message.base_def !== 0) {
      obj.base_def = message.base_def;
    }
    if (message.hp_curve !== undefined && message.hp_curve !== 0) {
      obj.hp_curve = avatarCurveTypeToJSON(message.hp_curve);
    }
    if (message.atk_curve !== undefined && message.atk_curve !== 0) {
      obj.atk_curve = avatarCurveTypeToJSON(message.atk_curve);
    }
    if (message.def_cruve !== undefined && message.def_cruve !== 0) {
      obj.def_cruve = avatarCurveTypeToJSON(message.def_cruve);
    }
    if (message.promo_data?.length) {
      obj.promo_data = message.promo_data.map((e) => PromotionData.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AvatarStatsData>, I>>(base?: I): AvatarStatsData {
    return AvatarStatsData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AvatarStatsData>, I>>(object: I): AvatarStatsData {
    const message = createBaseAvatarStatsData();
    message.base_hp = object.base_hp ?? 0;
    message.base_atk = object.base_atk ?? 0;
    message.base_def = object.base_def ?? 0;
    message.hp_curve = object.hp_curve ?? 0;
    message.atk_curve = object.atk_curve ?? 0;
    message.def_cruve = object.def_cruve ?? 0;
    message.promo_data = object.promo_data?.map((e) => PromotionData.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAvatarSkillsData(): AvatarSkillsData {
  return {
    skill: 0,
    burst: 0,
    attack: 0,
    burst_energy_cost: 0,
    attack_scaling: [],
    skill_scaling: [],
    burst_scaling: [],
  };
}

export const AvatarSkillsData = {
  encode(message: AvatarSkillsData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.skill !== undefined && message.skill !== 0) {
      writer.uint32(8).int32(message.skill);
    }
    if (message.burst !== undefined && message.burst !== 0) {
      writer.uint32(16).int32(message.burst);
    }
    if (message.attack !== undefined && message.attack !== 0) {
      writer.uint32(24).int32(message.attack);
    }
    if (message.burst_energy_cost !== undefined && message.burst_energy_cost !== 0) {
      writer.uint32(33).double(message.burst_energy_cost);
    }
    if (message.attack_scaling !== undefined && message.attack_scaling.length !== 0) {
      for (const v of message.attack_scaling) {
        AvatarSkillExcelIndexData.encode(v!, writer.uint32(42).fork()).ldelim();
      }
    }
    if (message.skill_scaling !== undefined && message.skill_scaling.length !== 0) {
      for (const v of message.skill_scaling) {
        AvatarSkillExcelIndexData.encode(v!, writer.uint32(50).fork()).ldelim();
      }
    }
    if (message.burst_scaling !== undefined && message.burst_scaling.length !== 0) {
      for (const v of message.burst_scaling) {
        AvatarSkillExcelIndexData.encode(v!, writer.uint32(58).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AvatarSkillsData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAvatarSkillsData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.skill = reader.int32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.burst = reader.int32();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.attack = reader.int32();
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.burst_energy_cost = reader.double();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.attack_scaling!.push(AvatarSkillExcelIndexData.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.skill_scaling!.push(AvatarSkillExcelIndexData.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.burst_scaling!.push(AvatarSkillExcelIndexData.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AvatarSkillsData {
    return {
      skill: isSet(object.skill) ? globalThis.Number(object.skill) : 0,
      burst: isSet(object.burst) ? globalThis.Number(object.burst) : 0,
      attack: isSet(object.attack) ? globalThis.Number(object.attack) : 0,
      burst_energy_cost: isSet(object.burst_energy_cost) ? globalThis.Number(object.burst_energy_cost) : 0,
      attack_scaling: globalThis.Array.isArray(object?.attack_scaling)
        ? object.attack_scaling.map((e: any) => AvatarSkillExcelIndexData.fromJSON(e))
        : [],
      skill_scaling: globalThis.Array.isArray(object?.skill_scaling)
        ? object.skill_scaling.map((e: any) => AvatarSkillExcelIndexData.fromJSON(e))
        : [],
      burst_scaling: globalThis.Array.isArray(object?.burst_scaling)
        ? object.burst_scaling.map((e: any) => AvatarSkillExcelIndexData.fromJSON(e))
        : [],
    };
  },

  toJSON(message: AvatarSkillsData): unknown {
    const obj: any = {};
    if (message.skill !== undefined && message.skill !== 0) {
      obj.skill = Math.round(message.skill);
    }
    if (message.burst !== undefined && message.burst !== 0) {
      obj.burst = Math.round(message.burst);
    }
    if (message.attack !== undefined && message.attack !== 0) {
      obj.attack = Math.round(message.attack);
    }
    if (message.burst_energy_cost !== undefined && message.burst_energy_cost !== 0) {
      obj.burst_energy_cost = message.burst_energy_cost;
    }
    if (message.attack_scaling?.length) {
      obj.attack_scaling = message.attack_scaling.map((e) => AvatarSkillExcelIndexData.toJSON(e));
    }
    if (message.skill_scaling?.length) {
      obj.skill_scaling = message.skill_scaling.map((e) => AvatarSkillExcelIndexData.toJSON(e));
    }
    if (message.burst_scaling?.length) {
      obj.burst_scaling = message.burst_scaling.map((e) => AvatarSkillExcelIndexData.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AvatarSkillsData>, I>>(base?: I): AvatarSkillsData {
    return AvatarSkillsData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AvatarSkillsData>, I>>(object: I): AvatarSkillsData {
    const message = createBaseAvatarSkillsData();
    message.skill = object.skill ?? 0;
    message.burst = object.burst ?? 0;
    message.attack = object.attack ?? 0;
    message.burst_energy_cost = object.burst_energy_cost ?? 0;
    message.attack_scaling = object.attack_scaling?.map((e) => AvatarSkillExcelIndexData.fromPartial(e)) || [];
    message.skill_scaling = object.skill_scaling?.map((e) => AvatarSkillExcelIndexData.fromPartial(e)) || [];
    message.burst_scaling = object.burst_scaling?.map((e) => AvatarSkillExcelIndexData.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAvatarSkillExcelIndexData(): AvatarSkillExcelIndexData {
  return { group_id: 0, index: 0, level_data: [] };
}

export const AvatarSkillExcelIndexData = {
  encode(message: AvatarSkillExcelIndexData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.group_id !== undefined && message.group_id !== 0) {
      writer.uint32(8).int32(message.group_id);
    }
    if (message.index !== undefined && message.index !== 0) {
      writer.uint32(16).int32(message.index);
    }
    if (message.level_data !== undefined && message.level_data.length !== 0) {
      for (const v of message.level_data) {
        AvatarSkillExcelLevelData.encode(v!, writer.uint32(26).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AvatarSkillExcelIndexData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAvatarSkillExcelIndexData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.group_id = reader.int32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.index = reader.int32();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.level_data!.push(AvatarSkillExcelLevelData.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AvatarSkillExcelIndexData {
    return {
      group_id: isSet(object.group_id) ? globalThis.Number(object.group_id) : 0,
      index: isSet(object.index) ? globalThis.Number(object.index) : 0,
      level_data: globalThis.Array.isArray(object?.level_data)
        ? object.level_data.map((e: any) => AvatarSkillExcelLevelData.fromJSON(e))
        : [],
    };
  },

  toJSON(message: AvatarSkillExcelIndexData): unknown {
    const obj: any = {};
    if (message.group_id !== undefined && message.group_id !== 0) {
      obj.group_id = Math.round(message.group_id);
    }
    if (message.index !== undefined && message.index !== 0) {
      obj.index = Math.round(message.index);
    }
    if (message.level_data?.length) {
      obj.level_data = message.level_data.map((e) => AvatarSkillExcelLevelData.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AvatarSkillExcelIndexData>, I>>(base?: I): AvatarSkillExcelIndexData {
    return AvatarSkillExcelIndexData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AvatarSkillExcelIndexData>, I>>(object: I): AvatarSkillExcelIndexData {
    const message = createBaseAvatarSkillExcelIndexData();
    message.group_id = object.group_id ?? 0;
    message.index = object.index ?? 0;
    message.level_data = object.level_data?.map((e) => AvatarSkillExcelLevelData.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAvatarSkillExcelLevelData(): AvatarSkillExcelLevelData {
  return { level: 0, value: 0 };
}

export const AvatarSkillExcelLevelData = {
  encode(message: AvatarSkillExcelLevelData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.level !== undefined && message.level !== 0) {
      writer.uint32(8).int32(message.level);
    }
    if (message.value !== undefined && message.value !== 0) {
      writer.uint32(17).double(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AvatarSkillExcelLevelData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAvatarSkillExcelLevelData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.level = reader.int32();
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.value = reader.double();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AvatarSkillExcelLevelData {
    return {
      level: isSet(object.level) ? globalThis.Number(object.level) : 0,
      value: isSet(object.value) ? globalThis.Number(object.value) : 0,
    };
  },

  toJSON(message: AvatarSkillExcelLevelData): unknown {
    const obj: any = {};
    if (message.level !== undefined && message.level !== 0) {
      obj.level = Math.round(message.level);
    }
    if (message.value !== undefined && message.value !== 0) {
      obj.value = message.value;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AvatarSkillExcelLevelData>, I>>(base?: I): AvatarSkillExcelLevelData {
    return AvatarSkillExcelLevelData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AvatarSkillExcelLevelData>, I>>(object: I): AvatarSkillExcelLevelData {
    const message = createBaseAvatarSkillExcelLevelData();
    message.level = object.level ?? 0;
    message.value = object.value ?? 0;
    return message;
  },
};

function createBaseWeaponDataMap(): WeaponDataMap {
  return { data: {} };
}

export const WeaponDataMap = {
  encode(message: WeaponDataMap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.data || {}).forEach(([key, value]) => {
      WeaponDataMap_DataEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WeaponDataMap {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWeaponDataMap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = WeaponDataMap_DataEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.data![entry1.key] = entry1.value;
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

  fromJSON(object: any): WeaponDataMap {
    return {
      data: isObject(object.data)
        ? Object.entries(object.data).reduce<{ [key: string]: WeaponData }>((acc, [key, value]) => {
          acc[key] = WeaponData.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: WeaponDataMap): unknown {
    const obj: any = {};
    if (message.data) {
      const entries = Object.entries(message.data);
      if (entries.length > 0) {
        obj.data = {};
        entries.forEach(([k, v]) => {
          obj.data[k] = WeaponData.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<WeaponDataMap>, I>>(base?: I): WeaponDataMap {
    return WeaponDataMap.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<WeaponDataMap>, I>>(object: I): WeaponDataMap {
    const message = createBaseWeaponDataMap();
    message.data = Object.entries(object.data ?? {}).reduce<{ [key: string]: WeaponData }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = WeaponData.fromPartial(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseWeaponDataMap_DataEntry(): WeaponDataMap_DataEntry {
  return { key: "", value: undefined };
}

export const WeaponDataMap_DataEntry = {
  encode(message: WeaponDataMap_DataEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      WeaponData.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WeaponDataMap_DataEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWeaponDataMap_DataEntry();
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

          message.value = WeaponData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): WeaponDataMap_DataEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? WeaponData.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: WeaponDataMap_DataEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = WeaponData.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<WeaponDataMap_DataEntry>, I>>(base?: I): WeaponDataMap_DataEntry {
    return WeaponDataMap_DataEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<WeaponDataMap_DataEntry>, I>>(object: I): WeaponDataMap_DataEntry {
    const message = createBaseWeaponDataMap_DataEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? WeaponData.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseWeaponData(): WeaponData {
  return {
    id: 0,
    key: "",
    rarity: 0,
    weapon_class: 0,
    image_name: "",
    base_stats: undefined,
    "name_text_hash_map ": 0,
  };
}

export const WeaponData = {
  encode(message: WeaponData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== 0) {
      writer.uint32(8).int32(message.id);
    }
    if (message.key !== undefined && message.key !== "") {
      writer.uint32(18).string(message.key);
    }
    if (message.rarity !== undefined && message.rarity !== 0) {
      writer.uint32(24).int32(message.rarity);
    }
    if (message.weapon_class !== undefined && message.weapon_class !== 0) {
      writer.uint32(32).int32(message.weapon_class);
    }
    if (message.image_name !== undefined && message.image_name !== "") {
      writer.uint32(42).string(message.image_name);
    }
    if (message.base_stats !== undefined) {
      WeaponStatsData.encode(message.base_stats, writer.uint32(50).fork()).ldelim();
    }
    if (message["name_text_hash_map "] !== undefined && message["name_text_hash_map "] !== 0) {
      writer.uint32(56).int64(message["name_text_hash_map "]);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WeaponData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWeaponData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = reader.int32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.key = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.rarity = reader.int32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.weapon_class = reader.int32() as any;
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.image_name = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.base_stats = WeaponStatsData.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message["name_text_hash_map "] = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): WeaponData {
    return {
      id: isSet(object.id) ? globalThis.Number(object.id) : 0,
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      rarity: isSet(object.rarity) ? globalThis.Number(object.rarity) : 0,
      weapon_class: isSet(object.weapon_class) ? weaponClassFromJSON(object.weapon_class) : 0,
      image_name: isSet(object.image_name) ? globalThis.String(object.image_name) : "",
      base_stats: isSet(object.base_stats) ? WeaponStatsData.fromJSON(object.base_stats) : undefined,
      "name_text_hash_map ": isSet(object["name_text_hash_map "])
        ? globalThis.Number(object["name_text_hash_map "])
        : 0,
    };
  },

  toJSON(message: WeaponData): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.key !== undefined && message.key !== "") {
      obj.key = message.key;
    }
    if (message.rarity !== undefined && message.rarity !== 0) {
      obj.rarity = Math.round(message.rarity);
    }
    if (message.weapon_class !== undefined && message.weapon_class !== 0) {
      obj.weapon_class = weaponClassToJSON(message.weapon_class);
    }
    if (message.image_name !== undefined && message.image_name !== "") {
      obj.image_name = message.image_name;
    }
    if (message.base_stats !== undefined) {
      obj.base_stats = WeaponStatsData.toJSON(message.base_stats);
    }
    if (message["name_text_hash_map "] !== undefined && message["name_text_hash_map "] !== 0) {
      obj["name_text_hash_map "] = Math.round(message["name_text_hash_map "]);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<WeaponData>, I>>(base?: I): WeaponData {
    return WeaponData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<WeaponData>, I>>(object: I): WeaponData {
    const message = createBaseWeaponData();
    message.id = object.id ?? 0;
    message.key = object.key ?? "";
    message.rarity = object.rarity ?? 0;
    message.weapon_class = object.weapon_class ?? 0;
    message.image_name = object.image_name ?? "";
    message.base_stats = (object.base_stats !== undefined && object.base_stats !== null)
      ? WeaponStatsData.fromPartial(object.base_stats)
      : undefined;
    message["name_text_hash_map "] = object["name_text_hash_map "] ?? 0;
    return message;
  },
};

function createBaseWeaponStatsData(): WeaponStatsData {
  return { base_props: [], promo_data: [] };
}

export const WeaponStatsData = {
  encode(message: WeaponStatsData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.base_props !== undefined && message.base_props.length !== 0) {
      for (const v of message.base_props) {
        WeaponProp.encode(v!, writer.uint32(10).fork()).ldelim();
      }
    }
    if (message.promo_data !== undefined && message.promo_data.length !== 0) {
      for (const v of message.promo_data) {
        PromotionData.encode(v!, writer.uint32(18).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WeaponStatsData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWeaponStatsData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.base_props!.push(WeaponProp.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.promo_data!.push(PromotionData.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): WeaponStatsData {
    return {
      base_props: globalThis.Array.isArray(object?.base_props)
        ? object.base_props.map((e: any) => WeaponProp.fromJSON(e))
        : [],
      promo_data: globalThis.Array.isArray(object?.promo_data)
        ? object.promo_data.map((e: any) => PromotionData.fromJSON(e))
        : [],
    };
  },

  toJSON(message: WeaponStatsData): unknown {
    const obj: any = {};
    if (message.base_props?.length) {
      obj.base_props = message.base_props.map((e) => WeaponProp.toJSON(e));
    }
    if (message.promo_data?.length) {
      obj.promo_data = message.promo_data.map((e) => PromotionData.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<WeaponStatsData>, I>>(base?: I): WeaponStatsData {
    return WeaponStatsData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<WeaponStatsData>, I>>(object: I): WeaponStatsData {
    const message = createBaseWeaponStatsData();
    message.base_props = object.base_props?.map((e) => WeaponProp.fromPartial(e)) || [];
    message.promo_data = object.promo_data?.map((e) => PromotionData.fromPartial(e)) || [];
    return message;
  },
};

function createBaseWeaponProp(): WeaponProp {
  return { prop_type: 0, initial_value: 0, curve: 0 };
}

export const WeaponProp = {
  encode(message: WeaponProp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.prop_type !== undefined && message.prop_type !== 0) {
      writer.uint32(8).int32(message.prop_type);
    }
    if (message.initial_value !== undefined && message.initial_value !== 0) {
      writer.uint32(17).double(message.initial_value);
    }
    if (message.curve !== undefined && message.curve !== 0) {
      writer.uint32(24).int32(message.curve);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WeaponProp {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWeaponProp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.prop_type = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.initial_value = reader.double();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.curve = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): WeaponProp {
    return {
      prop_type: isSet(object.prop_type) ? statTypeFromJSON(object.prop_type) : 0,
      initial_value: isSet(object.initial_value) ? globalThis.Number(object.initial_value) : 0,
      curve: isSet(object.curve) ? weaponCurveTypeFromJSON(object.curve) : 0,
    };
  },

  toJSON(message: WeaponProp): unknown {
    const obj: any = {};
    if (message.prop_type !== undefined && message.prop_type !== 0) {
      obj.prop_type = statTypeToJSON(message.prop_type);
    }
    if (message.initial_value !== undefined && message.initial_value !== 0) {
      obj.initial_value = message.initial_value;
    }
    if (message.curve !== undefined && message.curve !== 0) {
      obj.curve = weaponCurveTypeToJSON(message.curve);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<WeaponProp>, I>>(base?: I): WeaponProp {
    return WeaponProp.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<WeaponProp>, I>>(object: I): WeaponProp {
    const message = createBaseWeaponProp();
    message.prop_type = object.prop_type ?? 0;
    message.initial_value = object.initial_value ?? 0;
    message.curve = object.curve ?? 0;
    return message;
  },
};

function createBaseArtifactDataMap(): ArtifactDataMap {
  return { data: {} };
}

export const ArtifactDataMap = {
  encode(message: ArtifactDataMap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.data || {}).forEach(([key, value]) => {
      ArtifactDataMap_DataEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ArtifactDataMap {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseArtifactDataMap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = ArtifactDataMap_DataEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.data![entry1.key] = entry1.value;
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

  fromJSON(object: any): ArtifactDataMap {
    return {
      data: isObject(object.data)
        ? Object.entries(object.data).reduce<{ [key: string]: ArtifactData }>((acc, [key, value]) => {
          acc[key] = ArtifactData.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: ArtifactDataMap): unknown {
    const obj: any = {};
    if (message.data) {
      const entries = Object.entries(message.data);
      if (entries.length > 0) {
        obj.data = {};
        entries.forEach(([k, v]) => {
          obj.data[k] = ArtifactData.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ArtifactDataMap>, I>>(base?: I): ArtifactDataMap {
    return ArtifactDataMap.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ArtifactDataMap>, I>>(object: I): ArtifactDataMap {
    const message = createBaseArtifactDataMap();
    message.data = Object.entries(object.data ?? {}).reduce<{ [key: string]: ArtifactData }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = ArtifactData.fromPartial(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseArtifactDataMap_DataEntry(): ArtifactDataMap_DataEntry {
  return { key: "", value: undefined };
}

export const ArtifactDataMap_DataEntry = {
  encode(message: ArtifactDataMap_DataEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      ArtifactData.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ArtifactDataMap_DataEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseArtifactDataMap_DataEntry();
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

          message.value = ArtifactData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ArtifactDataMap_DataEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? ArtifactData.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: ArtifactDataMap_DataEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = ArtifactData.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ArtifactDataMap_DataEntry>, I>>(base?: I): ArtifactDataMap_DataEntry {
    return ArtifactDataMap_DataEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ArtifactDataMap_DataEntry>, I>>(object: I): ArtifactDataMap_DataEntry {
    const message = createBaseArtifactDataMap_DataEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? ArtifactData.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseArtifactData(): ArtifactData {
  return { id: 0, text_map_id: 0, key: "" };
}

export const ArtifactData = {
  encode(message: ArtifactData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== 0) {
      writer.uint32(8).int64(message.id);
    }
    if (message.text_map_id !== undefined && message.text_map_id !== 0) {
      writer.uint32(16).int64(message.text_map_id);
    }
    if (message.key !== undefined && message.key !== "") {
      writer.uint32(26).string(message.key);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ArtifactData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseArtifactData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = longToNumber(reader.int64() as Long);
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.text_map_id = longToNumber(reader.int64() as Long);
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.key = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ArtifactData {
    return {
      id: isSet(object.id) ? globalThis.Number(object.id) : 0,
      text_map_id: isSet(object.text_map_id) ? globalThis.Number(object.text_map_id) : 0,
      key: isSet(object.key) ? globalThis.String(object.key) : "",
    };
  },

  toJSON(message: ArtifactData): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.text_map_id !== undefined && message.text_map_id !== 0) {
      obj.text_map_id = Math.round(message.text_map_id);
    }
    if (message.key !== undefined && message.key !== "") {
      obj.key = message.key;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ArtifactData>, I>>(base?: I): ArtifactData {
    return ArtifactData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ArtifactData>, I>>(object: I): ArtifactData {
    const message = createBaseArtifactData();
    message.id = object.id ?? 0;
    message.text_map_id = object.text_map_id ?? 0;
    message.key = object.key ?? "";
    return message;
  },
};

function createBasePromotionData(): PromotionData {
  return { max_level: 0, add_props: [] };
}

export const PromotionData = {
  encode(message: PromotionData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.max_level !== undefined && message.max_level !== 0) {
      writer.uint32(8).int32(message.max_level);
    }
    if (message.add_props !== undefined && message.add_props.length !== 0) {
      for (const v of message.add_props) {
        PromotionAddProp.encode(v!, writer.uint32(18).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PromotionData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePromotionData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.max_level = reader.int32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.add_props!.push(PromotionAddProp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PromotionData {
    return {
      max_level: isSet(object.max_level) ? globalThis.Number(object.max_level) : 0,
      add_props: globalThis.Array.isArray(object?.add_props)
        ? object.add_props.map((e: any) => PromotionAddProp.fromJSON(e))
        : [],
    };
  },

  toJSON(message: PromotionData): unknown {
    const obj: any = {};
    if (message.max_level !== undefined && message.max_level !== 0) {
      obj.max_level = Math.round(message.max_level);
    }
    if (message.add_props?.length) {
      obj.add_props = message.add_props.map((e) => PromotionAddProp.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PromotionData>, I>>(base?: I): PromotionData {
    return PromotionData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PromotionData>, I>>(object: I): PromotionData {
    const message = createBasePromotionData();
    message.max_level = object.max_level ?? 0;
    message.add_props = object.add_props?.map((e) => PromotionAddProp.fromPartial(e)) || [];
    return message;
  },
};

function createBasePromotionAddProp(): PromotionAddProp {
  return { prop_type: 0, value: 0 };
}

export const PromotionAddProp = {
  encode(message: PromotionAddProp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.prop_type !== undefined && message.prop_type !== 0) {
      writer.uint32(8).int32(message.prop_type);
    }
    if (message.value !== undefined && message.value !== 0) {
      writer.uint32(17).double(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PromotionAddProp {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePromotionAddProp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.prop_type = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.value = reader.double();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PromotionAddProp {
    return {
      prop_type: isSet(object.prop_type) ? statTypeFromJSON(object.prop_type) : 0,
      value: isSet(object.value) ? globalThis.Number(object.value) : 0,
    };
  },

  toJSON(message: PromotionAddProp): unknown {
    const obj: any = {};
    if (message.prop_type !== undefined && message.prop_type !== 0) {
      obj.prop_type = statTypeToJSON(message.prop_type);
    }
    if (message.value !== undefined && message.value !== 0) {
      obj.value = message.value;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PromotionAddProp>, I>>(base?: I): PromotionAddProp {
    return PromotionAddProp.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PromotionAddProp>, I>>(object: I): PromotionAddProp {
    const message = createBasePromotionAddProp();
    message.prop_type = object.prop_type ?? 0;
    message.value = object.value ?? 0;
    return message;
  },
};

function createBaseMonsterData(): MonsterData {
  return { id: 0, key: "", base_stats: undefined, name_text_hash_map: 0 };
}

export const MonsterData = {
  encode(message: MonsterData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== undefined && message.id !== 0) {
      writer.uint32(8).int32(message.id);
    }
    if (message.key !== undefined && message.key !== "") {
      writer.uint32(18).string(message.key);
    }
    if (message.base_stats !== undefined) {
      MonsterStatsData.encode(message.base_stats, writer.uint32(26).fork()).ldelim();
    }
    if (message.name_text_hash_map !== undefined && message.name_text_hash_map !== 0) {
      writer.uint32(32).int64(message.name_text_hash_map);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonsterData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonsterData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.id = reader.int32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.key = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.base_stats = MonsterStatsData.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.name_text_hash_map = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MonsterData {
    return {
      id: isSet(object.id) ? globalThis.Number(object.id) : 0,
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      base_stats: isSet(object.base_stats) ? MonsterStatsData.fromJSON(object.base_stats) : undefined,
      name_text_hash_map: isSet(object.name_text_hash_map) ? globalThis.Number(object.name_text_hash_map) : 0,
    };
  },

  toJSON(message: MonsterData): unknown {
    const obj: any = {};
    if (message.id !== undefined && message.id !== 0) {
      obj.id = Math.round(message.id);
    }
    if (message.key !== undefined && message.key !== "") {
      obj.key = message.key;
    }
    if (message.base_stats !== undefined) {
      obj.base_stats = MonsterStatsData.toJSON(message.base_stats);
    }
    if (message.name_text_hash_map !== undefined && message.name_text_hash_map !== 0) {
      obj.name_text_hash_map = Math.round(message.name_text_hash_map);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MonsterData>, I>>(base?: I): MonsterData {
    return MonsterData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MonsterData>, I>>(object: I): MonsterData {
    const message = createBaseMonsterData();
    message.id = object.id ?? 0;
    message.key = object.key ?? "";
    message.base_stats = (object.base_stats !== undefined && object.base_stats !== null)
      ? MonsterStatsData.fromPartial(object.base_stats)
      : undefined;
    message.name_text_hash_map = object.name_text_hash_map ?? 0;
    return message;
  },
};

function createBaseMonsterStatsData(): MonsterStatsData {
  return { base_hp: 0, hp_curve: 0, resist: undefined, freeze_resist: 0, hp_drop: [] };
}

export const MonsterStatsData = {
  encode(message: MonsterStatsData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.base_hp !== undefined && message.base_hp !== 0) {
      writer.uint32(9).double(message.base_hp);
    }
    if (message.hp_curve !== undefined && message.hp_curve !== 0) {
      writer.uint32(16).int32(message.hp_curve);
    }
    if (message.resist !== undefined) {
      MonsterResistData.encode(message.resist, writer.uint32(26).fork()).ldelim();
    }
    if (message.freeze_resist !== undefined && message.freeze_resist !== 0) {
      writer.uint32(33).double(message.freeze_resist);
    }
    if (message.hp_drop !== undefined && message.hp_drop.length !== 0) {
      for (const v of message.hp_drop) {
        MonsterHPDrop.encode(v!, writer.uint32(42).fork()).ldelim();
      }
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonsterStatsData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonsterStatsData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.base_hp = reader.double();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.hp_curve = reader.int32() as any;
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.resist = MonsterResistData.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.freeze_resist = reader.double();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.hp_drop!.push(MonsterHPDrop.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MonsterStatsData {
    return {
      base_hp: isSet(object.base_hp) ? globalThis.Number(object.base_hp) : 0,
      hp_curve: isSet(object.hp_curve) ? monsterCurveTypeFromJSON(object.hp_curve) : 0,
      resist: isSet(object.resist) ? MonsterResistData.fromJSON(object.resist) : undefined,
      freeze_resist: isSet(object.freeze_resist) ? globalThis.Number(object.freeze_resist) : 0,
      hp_drop: globalThis.Array.isArray(object?.hp_drop)
        ? object.hp_drop.map((e: any) => MonsterHPDrop.fromJSON(e))
        : [],
    };
  },

  toJSON(message: MonsterStatsData): unknown {
    const obj: any = {};
    if (message.base_hp !== undefined && message.base_hp !== 0) {
      obj.base_hp = message.base_hp;
    }
    if (message.hp_curve !== undefined && message.hp_curve !== 0) {
      obj.hp_curve = monsterCurveTypeToJSON(message.hp_curve);
    }
    if (message.resist !== undefined) {
      obj.resist = MonsterResistData.toJSON(message.resist);
    }
    if (message.freeze_resist !== undefined && message.freeze_resist !== 0) {
      obj.freeze_resist = message.freeze_resist;
    }
    if (message.hp_drop?.length) {
      obj.hp_drop = message.hp_drop.map((e) => MonsterHPDrop.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MonsterStatsData>, I>>(base?: I): MonsterStatsData {
    return MonsterStatsData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MonsterStatsData>, I>>(object: I): MonsterStatsData {
    const message = createBaseMonsterStatsData();
    message.base_hp = object.base_hp ?? 0;
    message.hp_curve = object.hp_curve ?? 0;
    message.resist = (object.resist !== undefined && object.resist !== null)
      ? MonsterResistData.fromPartial(object.resist)
      : undefined;
    message.freeze_resist = object.freeze_resist ?? 0;
    message.hp_drop = object.hp_drop?.map((e) => MonsterHPDrop.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMonsterResistData(): MonsterResistData {
  return {
    fire_resist: 0,
    grass_resist: 0,
    water_resist: 0,
    electric_resist: 0,
    wind_resist: 0,
    ice_resist: 0,
    rock_resist: 0,
    physical_resist: 0,
  };
}

export const MonsterResistData = {
  encode(message: MonsterResistData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fire_resist !== undefined && message.fire_resist !== 0) {
      writer.uint32(9).double(message.fire_resist);
    }
    if (message.grass_resist !== undefined && message.grass_resist !== 0) {
      writer.uint32(17).double(message.grass_resist);
    }
    if (message.water_resist !== undefined && message.water_resist !== 0) {
      writer.uint32(25).double(message.water_resist);
    }
    if (message.electric_resist !== undefined && message.electric_resist !== 0) {
      writer.uint32(33).double(message.electric_resist);
    }
    if (message.wind_resist !== undefined && message.wind_resist !== 0) {
      writer.uint32(41).double(message.wind_resist);
    }
    if (message.ice_resist !== undefined && message.ice_resist !== 0) {
      writer.uint32(49).double(message.ice_resist);
    }
    if (message.rock_resist !== undefined && message.rock_resist !== 0) {
      writer.uint32(57).double(message.rock_resist);
    }
    if (message.physical_resist !== undefined && message.physical_resist !== 0) {
      writer.uint32(65).double(message.physical_resist);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonsterResistData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonsterResistData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.fire_resist = reader.double();
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.grass_resist = reader.double();
          continue;
        case 3:
          if (tag !== 25) {
            break;
          }

          message.water_resist = reader.double();
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.electric_resist = reader.double();
          continue;
        case 5:
          if (tag !== 41) {
            break;
          }

          message.wind_resist = reader.double();
          continue;
        case 6:
          if (tag !== 49) {
            break;
          }

          message.ice_resist = reader.double();
          continue;
        case 7:
          if (tag !== 57) {
            break;
          }

          message.rock_resist = reader.double();
          continue;
        case 8:
          if (tag !== 65) {
            break;
          }

          message.physical_resist = reader.double();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MonsterResistData {
    return {
      fire_resist: isSet(object.fire_resist) ? globalThis.Number(object.fire_resist) : 0,
      grass_resist: isSet(object.grass_resist) ? globalThis.Number(object.grass_resist) : 0,
      water_resist: isSet(object.water_resist) ? globalThis.Number(object.water_resist) : 0,
      electric_resist: isSet(object.electric_resist) ? globalThis.Number(object.electric_resist) : 0,
      wind_resist: isSet(object.wind_resist) ? globalThis.Number(object.wind_resist) : 0,
      ice_resist: isSet(object.ice_resist) ? globalThis.Number(object.ice_resist) : 0,
      rock_resist: isSet(object.rock_resist) ? globalThis.Number(object.rock_resist) : 0,
      physical_resist: isSet(object.physical_resist) ? globalThis.Number(object.physical_resist) : 0,
    };
  },

  toJSON(message: MonsterResistData): unknown {
    const obj: any = {};
    if (message.fire_resist !== undefined && message.fire_resist !== 0) {
      obj.fire_resist = message.fire_resist;
    }
    if (message.grass_resist !== undefined && message.grass_resist !== 0) {
      obj.grass_resist = message.grass_resist;
    }
    if (message.water_resist !== undefined && message.water_resist !== 0) {
      obj.water_resist = message.water_resist;
    }
    if (message.electric_resist !== undefined && message.electric_resist !== 0) {
      obj.electric_resist = message.electric_resist;
    }
    if (message.wind_resist !== undefined && message.wind_resist !== 0) {
      obj.wind_resist = message.wind_resist;
    }
    if (message.ice_resist !== undefined && message.ice_resist !== 0) {
      obj.ice_resist = message.ice_resist;
    }
    if (message.rock_resist !== undefined && message.rock_resist !== 0) {
      obj.rock_resist = message.rock_resist;
    }
    if (message.physical_resist !== undefined && message.physical_resist !== 0) {
      obj.physical_resist = message.physical_resist;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MonsterResistData>, I>>(base?: I): MonsterResistData {
    return MonsterResistData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MonsterResistData>, I>>(object: I): MonsterResistData {
    const message = createBaseMonsterResistData();
    message.fire_resist = object.fire_resist ?? 0;
    message.grass_resist = object.grass_resist ?? 0;
    message.water_resist = object.water_resist ?? 0;
    message.electric_resist = object.electric_resist ?? 0;
    message.wind_resist = object.wind_resist ?? 0;
    message.ice_resist = object.ice_resist ?? 0;
    message.rock_resist = object.rock_resist ?? 0;
    message.physical_resist = object.physical_resist ?? 0;
    return message;
  },
};

function createBaseMonsterHPDrop(): MonsterHPDrop {
  return { drop_id: 0, hp_percent: 0 };
}

export const MonsterHPDrop = {
  encode(message: MonsterHPDrop, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.drop_id !== undefined && message.drop_id !== 0) {
      writer.uint32(8).int32(message.drop_id);
    }
    if (message.hp_percent !== undefined && message.hp_percent !== 0) {
      writer.uint32(17).double(message.hp_percent);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MonsterHPDrop {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMonsterHPDrop();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.drop_id = reader.int32();
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.hp_percent = reader.double();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): MonsterHPDrop {
    return {
      drop_id: isSet(object.drop_id) ? globalThis.Number(object.drop_id) : 0,
      hp_percent: isSet(object.hp_percent) ? globalThis.Number(object.hp_percent) : 0,
    };
  },

  toJSON(message: MonsterHPDrop): unknown {
    const obj: any = {};
    if (message.drop_id !== undefined && message.drop_id !== 0) {
      obj.drop_id = Math.round(message.drop_id);
    }
    if (message.hp_percent !== undefined && message.hp_percent !== 0) {
      obj.hp_percent = message.hp_percent;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<MonsterHPDrop>, I>>(base?: I): MonsterHPDrop {
    return MonsterHPDrop.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<MonsterHPDrop>, I>>(object: I): MonsterHPDrop {
    const message = createBaseMonsterHPDrop();
    message.drop_id = object.drop_id ?? 0;
    message.hp_percent = object.hp_percent ?? 0;
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
