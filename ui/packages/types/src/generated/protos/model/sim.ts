/* eslint-disable */
import _m0 from "protobufjs/minimal";

export interface Character {
  name?: string | undefined;
  element?: string | undefined;
  level?: number | undefined;
  max_level?: number | undefined;
  cons?: number | undefined;
  weapon?: Weapon | undefined;
  talents?: CharacterTalents | undefined;
  sets?: { [key: string]: number } | undefined;
  stats?: number[] | undefined;
  snapshot?: number[] | undefined;
}

export interface Character_SetsEntry {
  key: string;
  value: number;
}

export interface CharacterTalents {
  attack?: number | undefined;
  skill?: number | undefined;
  burst?: number | undefined;
}

export interface Weapon {
  name?: string | undefined;
  refine?: number | undefined;
  level?: number | undefined;
  max_level?: number | undefined;
}

export interface Enemy {
  level?: number | undefined;
  hp?: number | undefined;
  resist?: { [key: string]: number } | undefined;
  position?: Coord | undefined;
  particle_drop_threshold?: number | undefined;
  particle_drop_count?: number | undefined;
  particle_element?: string | undefined;
  name?: string | undefined;
  modified?: boolean | undefined;
}

export interface Enemy_ResistEntry {
  key: string;
  value: number;
}

export interface Coord {
  x?: number | undefined;
  y?: number | undefined;
  r?: number | undefined;
}

export interface SimulatorSettings {
  duration?: number | undefined;
  damage_mode?: boolean | undefined;
  enable_hitlag?: boolean | undefined;
  def_halt?: boolean | undefined;
  number_of_workers?: number | undefined;
  iterations?: number | undefined;
  delays?: Delays | undefined;
  ignore_burst_energy?: boolean | undefined;
}

export interface Delays {
  skill?: number | undefined;
  burst?: number | undefined;
  attack?: number | undefined;
  charge?: number | undefined;
  aim?: number | undefined;
  dash?: number | undefined;
  jump?: number | undefined;
  swap?: number | undefined;
}

export interface EnergySettings {
  active?: boolean | undefined;
  once?: boolean | undefined;
  start?: number | undefined;
  end?: number | undefined;
  amount?: number | undefined;
  last_energy_drop?: number | undefined;
}

function createBaseCharacter(): Character {
  return {
    name: "",
    element: "",
    level: 0,
    max_level: 0,
    cons: 0,
    weapon: undefined,
    talents: undefined,
    sets: {},
    stats: [],
    snapshot: [],
  };
}

export const Character = {
  encode(message: Character, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== undefined && message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.element !== undefined && message.element !== "") {
      writer.uint32(18).string(message.element);
    }
    if (message.level !== undefined && message.level !== 0) {
      writer.uint32(24).int32(message.level);
    }
    if (message.max_level !== undefined && message.max_level !== 0) {
      writer.uint32(32).int32(message.max_level);
    }
    if (message.cons !== undefined && message.cons !== 0) {
      writer.uint32(40).int32(message.cons);
    }
    if (message.weapon !== undefined) {
      Weapon.encode(message.weapon, writer.uint32(50).fork()).ldelim();
    }
    if (message.talents !== undefined) {
      CharacterTalents.encode(message.talents, writer.uint32(58).fork()).ldelim();
    }
    Object.entries(message.sets || {}).forEach(([key, value]) => {
      Character_SetsEntry.encode({ key: key as any, value }, writer.uint32(66).fork()).ldelim();
    });
    if (message.stats !== undefined && message.stats.length !== 0) {
      writer.uint32(74).fork();
      for (const v of message.stats) {
        writer.double(v);
      }
      writer.ldelim();
    }
    if (message.snapshot !== undefined && message.snapshot.length !== 0) {
      writer.uint32(82).fork();
      for (const v of message.snapshot) {
        writer.double(v);
      }
      writer.ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Character {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCharacter();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.element = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.level = reader.int32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.max_level = reader.int32();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.cons = reader.int32();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.weapon = Weapon.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.talents = CharacterTalents.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          const entry8 = Character_SetsEntry.decode(reader, reader.uint32());
          if (entry8.value !== undefined) {
            message.sets![entry8.key] = entry8.value;
          }
          continue;
        case 9:
          if (tag === 73) {
            message.stats!.push(reader.double());

            continue;
          }

          if (tag === 74) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.stats!.push(reader.double());
            }

            continue;
          }

          break;
        case 10:
          if (tag === 81) {
            message.snapshot!.push(reader.double());

            continue;
          }

          if (tag === 82) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.snapshot!.push(reader.double());
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

  fromJSON(object: any): Character {
    return {
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      element: isSet(object.element) ? globalThis.String(object.element) : "",
      level: isSet(object.level) ? globalThis.Number(object.level) : 0,
      max_level: isSet(object.max_level) ? globalThis.Number(object.max_level) : 0,
      cons: isSet(object.cons) ? globalThis.Number(object.cons) : 0,
      weapon: isSet(object.weapon) ? Weapon.fromJSON(object.weapon) : undefined,
      talents: isSet(object.talents) ? CharacterTalents.fromJSON(object.talents) : undefined,
      sets: isObject(object.sets)
        ? Object.entries(object.sets).reduce<{ [key: string]: number }>((acc, [key, value]) => {
          acc[key] = Number(value);
          return acc;
        }, {})
        : {},
      stats: globalThis.Array.isArray(object?.stats) ? object.stats.map((e: any) => globalThis.Number(e)) : [],
      snapshot: globalThis.Array.isArray(object?.snapshot) ? object.snapshot.map((e: any) => globalThis.Number(e)) : [],
    };
  },

  toJSON(message: Character): unknown {
    const obj: any = {};
    if (message.name !== undefined && message.name !== "") {
      obj.name = message.name;
    }
    if (message.element !== undefined && message.element !== "") {
      obj.element = message.element;
    }
    if (message.level !== undefined && message.level !== 0) {
      obj.level = Math.round(message.level);
    }
    if (message.max_level !== undefined && message.max_level !== 0) {
      obj.max_level = Math.round(message.max_level);
    }
    if (message.cons !== undefined && message.cons !== 0) {
      obj.cons = Math.round(message.cons);
    }
    if (message.weapon !== undefined) {
      obj.weapon = Weapon.toJSON(message.weapon);
    }
    if (message.talents !== undefined) {
      obj.talents = CharacterTalents.toJSON(message.talents);
    }
    if (message.sets) {
      const entries = Object.entries(message.sets);
      if (entries.length > 0) {
        obj.sets = {};
        entries.forEach(([k, v]) => {
          obj.sets[k] = Math.round(v);
        });
      }
    }
    if (message.stats?.length) {
      obj.stats = message.stats;
    }
    if (message.snapshot?.length) {
      obj.snapshot = message.snapshot;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Character>, I>>(base?: I): Character {
    return Character.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Character>, I>>(object: I): Character {
    const message = createBaseCharacter();
    message.name = object.name ?? "";
    message.element = object.element ?? "";
    message.level = object.level ?? 0;
    message.max_level = object.max_level ?? 0;
    message.cons = object.cons ?? 0;
    message.weapon = (object.weapon !== undefined && object.weapon !== null)
      ? Weapon.fromPartial(object.weapon)
      : undefined;
    message.talents = (object.talents !== undefined && object.talents !== null)
      ? CharacterTalents.fromPartial(object.talents)
      : undefined;
    message.sets = Object.entries(object.sets ?? {}).reduce<{ [key: string]: number }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = globalThis.Number(value);
      }
      return acc;
    }, {});
    message.stats = object.stats?.map((e) => e) || [];
    message.snapshot = object.snapshot?.map((e) => e) || [];
    return message;
  },
};

function createBaseCharacter_SetsEntry(): Character_SetsEntry {
  return { key: "", value: 0 };
}

export const Character_SetsEntry = {
  encode(message: Character_SetsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== 0) {
      writer.uint32(16).int32(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Character_SetsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCharacter_SetsEntry();
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
          if (tag !== 16) {
            break;
          }

          message.value = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Character_SetsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? globalThis.Number(object.value) : 0,
    };
  },

  toJSON(message: Character_SetsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== 0) {
      obj.value = Math.round(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Character_SetsEntry>, I>>(base?: I): Character_SetsEntry {
    return Character_SetsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Character_SetsEntry>, I>>(object: I): Character_SetsEntry {
    const message = createBaseCharacter_SetsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? 0;
    return message;
  },
};

function createBaseCharacterTalents(): CharacterTalents {
  return { attack: 0, skill: 0, burst: 0 };
}

export const CharacterTalents = {
  encode(message: CharacterTalents, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.attack !== undefined && message.attack !== 0) {
      writer.uint32(8).int32(message.attack);
    }
    if (message.skill !== undefined && message.skill !== 0) {
      writer.uint32(16).int32(message.skill);
    }
    if (message.burst !== undefined && message.burst !== 0) {
      writer.uint32(24).int32(message.burst);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CharacterTalents {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCharacterTalents();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.attack = reader.int32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.skill = reader.int32();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.burst = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CharacterTalents {
    return {
      attack: isSet(object.attack) ? globalThis.Number(object.attack) : 0,
      skill: isSet(object.skill) ? globalThis.Number(object.skill) : 0,
      burst: isSet(object.burst) ? globalThis.Number(object.burst) : 0,
    };
  },

  toJSON(message: CharacterTalents): unknown {
    const obj: any = {};
    if (message.attack !== undefined && message.attack !== 0) {
      obj.attack = Math.round(message.attack);
    }
    if (message.skill !== undefined && message.skill !== 0) {
      obj.skill = Math.round(message.skill);
    }
    if (message.burst !== undefined && message.burst !== 0) {
      obj.burst = Math.round(message.burst);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CharacterTalents>, I>>(base?: I): CharacterTalents {
    return CharacterTalents.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CharacterTalents>, I>>(object: I): CharacterTalents {
    const message = createBaseCharacterTalents();
    message.attack = object.attack ?? 0;
    message.skill = object.skill ?? 0;
    message.burst = object.burst ?? 0;
    return message;
  },
};

function createBaseWeapon(): Weapon {
  return { name: "", refine: 0, level: 0, max_level: 0 };
}

export const Weapon = {
  encode(message: Weapon, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== undefined && message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.refine !== undefined && message.refine !== 0) {
      writer.uint32(16).int32(message.refine);
    }
    if (message.level !== undefined && message.level !== 0) {
      writer.uint32(24).int32(message.level);
    }
    if (message.max_level !== undefined && message.max_level !== 0) {
      writer.uint32(32).int32(message.max_level);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Weapon {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWeapon();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.refine = reader.int32();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.level = reader.int32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.max_level = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Weapon {
    return {
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      refine: isSet(object.refine) ? globalThis.Number(object.refine) : 0,
      level: isSet(object.level) ? globalThis.Number(object.level) : 0,
      max_level: isSet(object.max_level) ? globalThis.Number(object.max_level) : 0,
    };
  },

  toJSON(message: Weapon): unknown {
    const obj: any = {};
    if (message.name !== undefined && message.name !== "") {
      obj.name = message.name;
    }
    if (message.refine !== undefined && message.refine !== 0) {
      obj.refine = Math.round(message.refine);
    }
    if (message.level !== undefined && message.level !== 0) {
      obj.level = Math.round(message.level);
    }
    if (message.max_level !== undefined && message.max_level !== 0) {
      obj.max_level = Math.round(message.max_level);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Weapon>, I>>(base?: I): Weapon {
    return Weapon.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Weapon>, I>>(object: I): Weapon {
    const message = createBaseWeapon();
    message.name = object.name ?? "";
    message.refine = object.refine ?? 0;
    message.level = object.level ?? 0;
    message.max_level = object.max_level ?? 0;
    return message;
  },
};

function createBaseEnemy(): Enemy {
  return {
    level: 0,
    hp: 0,
    resist: {},
    position: undefined,
    particle_drop_threshold: 0,
    particle_drop_count: 0,
    particle_element: "",
    name: "",
    modified: false,
  };
}

export const Enemy = {
  encode(message: Enemy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.level !== undefined && message.level !== 0) {
      writer.uint32(8).int32(message.level);
    }
    if (message.hp !== undefined && message.hp !== 0) {
      writer.uint32(17).double(message.hp);
    }
    Object.entries(message.resist || {}).forEach(([key, value]) => {
      Enemy_ResistEntry.encode({ key: key as any, value }, writer.uint32(26).fork()).ldelim();
    });
    if (message.position !== undefined) {
      Coord.encode(message.position, writer.uint32(34).fork()).ldelim();
    }
    if (message.particle_drop_threshold !== undefined && message.particle_drop_threshold !== 0) {
      writer.uint32(41).double(message.particle_drop_threshold);
    }
    if (message.particle_drop_count !== undefined && message.particle_drop_count !== 0) {
      writer.uint32(49).double(message.particle_drop_count);
    }
    if (message.particle_element !== undefined && message.particle_element !== "") {
      writer.uint32(58).string(message.particle_element);
    }
    if (message.name !== undefined && message.name !== "") {
      writer.uint32(66).string(message.name);
    }
    if (message.modified !== undefined && message.modified !== false) {
      writer.uint32(72).bool(message.modified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Enemy {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnemy();
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

          message.hp = reader.double();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          const entry3 = Enemy_ResistEntry.decode(reader, reader.uint32());
          if (entry3.value !== undefined) {
            message.resist![entry3.key] = entry3.value;
          }
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.position = Coord.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 41) {
            break;
          }

          message.particle_drop_threshold = reader.double();
          continue;
        case 6:
          if (tag !== 49) {
            break;
          }

          message.particle_drop_count = reader.double();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.particle_element = reader.string();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.name = reader.string();
          continue;
        case 9:
          if (tag !== 72) {
            break;
          }

          message.modified = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Enemy {
    return {
      level: isSet(object.level) ? globalThis.Number(object.level) : 0,
      hp: isSet(object.hp) ? globalThis.Number(object.hp) : 0,
      resist: isObject(object.resist)
        ? Object.entries(object.resist).reduce<{ [key: string]: number }>((acc, [key, value]) => {
          acc[key] = Number(value);
          return acc;
        }, {})
        : {},
      position: isSet(object.position) ? Coord.fromJSON(object.position) : undefined,
      particle_drop_threshold: isSet(object.particle_drop_threshold)
        ? globalThis.Number(object.particle_drop_threshold)
        : 0,
      particle_drop_count: isSet(object.particle_drop_count) ? globalThis.Number(object.particle_drop_count) : 0,
      particle_element: isSet(object.particle_element) ? globalThis.String(object.particle_element) : "",
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      modified: isSet(object.modified) ? globalThis.Boolean(object.modified) : false,
    };
  },

  toJSON(message: Enemy): unknown {
    const obj: any = {};
    if (message.level !== undefined && message.level !== 0) {
      obj.level = Math.round(message.level);
    }
    if (message.hp !== undefined && message.hp !== 0) {
      obj.hp = message.hp;
    }
    if (message.resist) {
      const entries = Object.entries(message.resist);
      if (entries.length > 0) {
        obj.resist = {};
        entries.forEach(([k, v]) => {
          obj.resist[k] = v;
        });
      }
    }
    if (message.position !== undefined) {
      obj.position = Coord.toJSON(message.position);
    }
    if (message.particle_drop_threshold !== undefined && message.particle_drop_threshold !== 0) {
      obj.particle_drop_threshold = message.particle_drop_threshold;
    }
    if (message.particle_drop_count !== undefined && message.particle_drop_count !== 0) {
      obj.particle_drop_count = message.particle_drop_count;
    }
    if (message.particle_element !== undefined && message.particle_element !== "") {
      obj.particle_element = message.particle_element;
    }
    if (message.name !== undefined && message.name !== "") {
      obj.name = message.name;
    }
    if (message.modified !== undefined && message.modified !== false) {
      obj.modified = message.modified;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Enemy>, I>>(base?: I): Enemy {
    return Enemy.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Enemy>, I>>(object: I): Enemy {
    const message = createBaseEnemy();
    message.level = object.level ?? 0;
    message.hp = object.hp ?? 0;
    message.resist = Object.entries(object.resist ?? {}).reduce<{ [key: string]: number }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = globalThis.Number(value);
      }
      return acc;
    }, {});
    message.position = (object.position !== undefined && object.position !== null)
      ? Coord.fromPartial(object.position)
      : undefined;
    message.particle_drop_threshold = object.particle_drop_threshold ?? 0;
    message.particle_drop_count = object.particle_drop_count ?? 0;
    message.particle_element = object.particle_element ?? "";
    message.name = object.name ?? "";
    message.modified = object.modified ?? false;
    return message;
  },
};

function createBaseEnemy_ResistEntry(): Enemy_ResistEntry {
  return { key: "", value: 0 };
}

export const Enemy_ResistEntry = {
  encode(message: Enemy_ResistEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== 0) {
      writer.uint32(17).double(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Enemy_ResistEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnemy_ResistEntry();
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

  fromJSON(object: any): Enemy_ResistEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? globalThis.Number(object.value) : 0,
    };
  },

  toJSON(message: Enemy_ResistEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== 0) {
      obj.value = message.value;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Enemy_ResistEntry>, I>>(base?: I): Enemy_ResistEntry {
    return Enemy_ResistEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Enemy_ResistEntry>, I>>(object: I): Enemy_ResistEntry {
    const message = createBaseEnemy_ResistEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? 0;
    return message;
  },
};

function createBaseCoord(): Coord {
  return { x: 0, y: 0, r: 0 };
}

export const Coord = {
  encode(message: Coord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.x !== undefined && message.x !== 0) {
      writer.uint32(9).double(message.x);
    }
    if (message.y !== undefined && message.y !== 0) {
      writer.uint32(17).double(message.y);
    }
    if (message.r !== undefined && message.r !== 0) {
      writer.uint32(25).double(message.r);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Coord {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCoord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.x = reader.double();
          continue;
        case 2:
          if (tag !== 17) {
            break;
          }

          message.y = reader.double();
          continue;
        case 3:
          if (tag !== 25) {
            break;
          }

          message.r = reader.double();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Coord {
    return {
      x: isSet(object.x) ? globalThis.Number(object.x) : 0,
      y: isSet(object.y) ? globalThis.Number(object.y) : 0,
      r: isSet(object.r) ? globalThis.Number(object.r) : 0,
    };
  },

  toJSON(message: Coord): unknown {
    const obj: any = {};
    if (message.x !== undefined && message.x !== 0) {
      obj.x = message.x;
    }
    if (message.y !== undefined && message.y !== 0) {
      obj.y = message.y;
    }
    if (message.r !== undefined && message.r !== 0) {
      obj.r = message.r;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Coord>, I>>(base?: I): Coord {
    return Coord.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Coord>, I>>(object: I): Coord {
    const message = createBaseCoord();
    message.x = object.x ?? 0;
    message.y = object.y ?? 0;
    message.r = object.r ?? 0;
    return message;
  },
};

function createBaseSimulatorSettings(): SimulatorSettings {
  return {
    duration: 0,
    damage_mode: false,
    enable_hitlag: false,
    def_halt: false,
    number_of_workers: 0,
    iterations: 0,
    delays: undefined,
    ignore_burst_energy: false,
  };
}

export const SimulatorSettings = {
  encode(message: SimulatorSettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.duration !== undefined && message.duration !== 0) {
      writer.uint32(9).double(message.duration);
    }
    if (message.damage_mode !== undefined && message.damage_mode !== false) {
      writer.uint32(16).bool(message.damage_mode);
    }
    if (message.enable_hitlag !== undefined && message.enable_hitlag !== false) {
      writer.uint32(24).bool(message.enable_hitlag);
    }
    if (message.def_halt !== undefined && message.def_halt !== false) {
      writer.uint32(32).bool(message.def_halt);
    }
    if (message.number_of_workers !== undefined && message.number_of_workers !== 0) {
      writer.uint32(40).uint32(message.number_of_workers);
    }
    if (message.iterations !== undefined && message.iterations !== 0) {
      writer.uint32(48).uint32(message.iterations);
    }
    if (message.delays !== undefined) {
      Delays.encode(message.delays, writer.uint32(58).fork()).ldelim();
    }
    if (message.ignore_burst_energy !== undefined && message.ignore_burst_energy !== false) {
      writer.uint32(64).bool(message.ignore_burst_energy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SimulatorSettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSimulatorSettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 9) {
            break;
          }

          message.duration = reader.double();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.damage_mode = reader.bool();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.enable_hitlag = reader.bool();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.def_halt = reader.bool();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.number_of_workers = reader.uint32();
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

          message.delays = Delays.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.ignore_burst_energy = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SimulatorSettings {
    return {
      duration: isSet(object.duration) ? globalThis.Number(object.duration) : 0,
      damage_mode: isSet(object.damage_mode) ? globalThis.Boolean(object.damage_mode) : false,
      enable_hitlag: isSet(object.enable_hitlag) ? globalThis.Boolean(object.enable_hitlag) : false,
      def_halt: isSet(object.def_halt) ? globalThis.Boolean(object.def_halt) : false,
      number_of_workers: isSet(object.number_of_workers) ? globalThis.Number(object.number_of_workers) : 0,
      iterations: isSet(object.iterations) ? globalThis.Number(object.iterations) : 0,
      delays: isSet(object.delays) ? Delays.fromJSON(object.delays) : undefined,
      ignore_burst_energy: isSet(object.ignore_burst_energy) ? globalThis.Boolean(object.ignore_burst_energy) : false,
    };
  },

  toJSON(message: SimulatorSettings): unknown {
    const obj: any = {};
    if (message.duration !== undefined && message.duration !== 0) {
      obj.duration = message.duration;
    }
    if (message.damage_mode !== undefined && message.damage_mode !== false) {
      obj.damage_mode = message.damage_mode;
    }
    if (message.enable_hitlag !== undefined && message.enable_hitlag !== false) {
      obj.enable_hitlag = message.enable_hitlag;
    }
    if (message.def_halt !== undefined && message.def_halt !== false) {
      obj.def_halt = message.def_halt;
    }
    if (message.number_of_workers !== undefined && message.number_of_workers !== 0) {
      obj.number_of_workers = Math.round(message.number_of_workers);
    }
    if (message.iterations !== undefined && message.iterations !== 0) {
      obj.iterations = Math.round(message.iterations);
    }
    if (message.delays !== undefined) {
      obj.delays = Delays.toJSON(message.delays);
    }
    if (message.ignore_burst_energy !== undefined && message.ignore_burst_energy !== false) {
      obj.ignore_burst_energy = message.ignore_burst_energy;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SimulatorSettings>, I>>(base?: I): SimulatorSettings {
    return SimulatorSettings.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SimulatorSettings>, I>>(object: I): SimulatorSettings {
    const message = createBaseSimulatorSettings();
    message.duration = object.duration ?? 0;
    message.damage_mode = object.damage_mode ?? false;
    message.enable_hitlag = object.enable_hitlag ?? false;
    message.def_halt = object.def_halt ?? false;
    message.number_of_workers = object.number_of_workers ?? 0;
    message.iterations = object.iterations ?? 0;
    message.delays = (object.delays !== undefined && object.delays !== null)
      ? Delays.fromPartial(object.delays)
      : undefined;
    message.ignore_burst_energy = object.ignore_burst_energy ?? false;
    return message;
  },
};

function createBaseDelays(): Delays {
  return { skill: 0, burst: 0, attack: 0, charge: 0, aim: 0, dash: 0, jump: 0, swap: 0 };
}

export const Delays = {
  encode(message: Delays, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.skill !== undefined && message.skill !== 0) {
      writer.uint32(8).int32(message.skill);
    }
    if (message.burst !== undefined && message.burst !== 0) {
      writer.uint32(16).int32(message.burst);
    }
    if (message.attack !== undefined && message.attack !== 0) {
      writer.uint32(24).int32(message.attack);
    }
    if (message.charge !== undefined && message.charge !== 0) {
      writer.uint32(32).int32(message.charge);
    }
    if (message.aim !== undefined && message.aim !== 0) {
      writer.uint32(40).int32(message.aim);
    }
    if (message.dash !== undefined && message.dash !== 0) {
      writer.uint32(48).int32(message.dash);
    }
    if (message.jump !== undefined && message.jump !== 0) {
      writer.uint32(56).int32(message.jump);
    }
    if (message.swap !== undefined && message.swap !== 0) {
      writer.uint32(64).int32(message.swap);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Delays {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDelays();
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
          if (tag !== 32) {
            break;
          }

          message.charge = reader.int32();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.aim = reader.int32();
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.dash = reader.int32();
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.jump = reader.int32();
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.swap = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Delays {
    return {
      skill: isSet(object.skill) ? globalThis.Number(object.skill) : 0,
      burst: isSet(object.burst) ? globalThis.Number(object.burst) : 0,
      attack: isSet(object.attack) ? globalThis.Number(object.attack) : 0,
      charge: isSet(object.charge) ? globalThis.Number(object.charge) : 0,
      aim: isSet(object.aim) ? globalThis.Number(object.aim) : 0,
      dash: isSet(object.dash) ? globalThis.Number(object.dash) : 0,
      jump: isSet(object.jump) ? globalThis.Number(object.jump) : 0,
      swap: isSet(object.swap) ? globalThis.Number(object.swap) : 0,
    };
  },

  toJSON(message: Delays): unknown {
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
    if (message.charge !== undefined && message.charge !== 0) {
      obj.charge = Math.round(message.charge);
    }
    if (message.aim !== undefined && message.aim !== 0) {
      obj.aim = Math.round(message.aim);
    }
    if (message.dash !== undefined && message.dash !== 0) {
      obj.dash = Math.round(message.dash);
    }
    if (message.jump !== undefined && message.jump !== 0) {
      obj.jump = Math.round(message.jump);
    }
    if (message.swap !== undefined && message.swap !== 0) {
      obj.swap = Math.round(message.swap);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Delays>, I>>(base?: I): Delays {
    return Delays.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Delays>, I>>(object: I): Delays {
    const message = createBaseDelays();
    message.skill = object.skill ?? 0;
    message.burst = object.burst ?? 0;
    message.attack = object.attack ?? 0;
    message.charge = object.charge ?? 0;
    message.aim = object.aim ?? 0;
    message.dash = object.dash ?? 0;
    message.jump = object.jump ?? 0;
    message.swap = object.swap ?? 0;
    return message;
  },
};

function createBaseEnergySettings(): EnergySettings {
  return { active: false, once: false, start: 0, end: 0, amount: 0, last_energy_drop: 0 };
}

export const EnergySettings = {
  encode(message: EnergySettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.active !== undefined && message.active !== false) {
      writer.uint32(8).bool(message.active);
    }
    if (message.once !== undefined && message.once !== false) {
      writer.uint32(16).bool(message.once);
    }
    if (message.start !== undefined && message.start !== 0) {
      writer.uint32(24).int32(message.start);
    }
    if (message.end !== undefined && message.end !== 0) {
      writer.uint32(32).int32(message.end);
    }
    if (message.amount !== undefined && message.amount !== 0) {
      writer.uint32(40).int32(message.amount);
    }
    if (message.last_energy_drop !== undefined && message.last_energy_drop !== 0) {
      writer.uint32(48).int32(message.last_energy_drop);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EnergySettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnergySettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.active = reader.bool();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.once = reader.bool();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.start = reader.int32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.end = reader.int32();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.amount = reader.int32();
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.last_energy_drop = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EnergySettings {
    return {
      active: isSet(object.active) ? globalThis.Boolean(object.active) : false,
      once: isSet(object.once) ? globalThis.Boolean(object.once) : false,
      start: isSet(object.start) ? globalThis.Number(object.start) : 0,
      end: isSet(object.end) ? globalThis.Number(object.end) : 0,
      amount: isSet(object.amount) ? globalThis.Number(object.amount) : 0,
      last_energy_drop: isSet(object.last_energy_drop) ? globalThis.Number(object.last_energy_drop) : 0,
    };
  },

  toJSON(message: EnergySettings): unknown {
    const obj: any = {};
    if (message.active !== undefined && message.active !== false) {
      obj.active = message.active;
    }
    if (message.once !== undefined && message.once !== false) {
      obj.once = message.once;
    }
    if (message.start !== undefined && message.start !== 0) {
      obj.start = Math.round(message.start);
    }
    if (message.end !== undefined && message.end !== 0) {
      obj.end = Math.round(message.end);
    }
    if (message.amount !== undefined && message.amount !== 0) {
      obj.amount = Math.round(message.amount);
    }
    if (message.last_energy_drop !== undefined && message.last_energy_drop !== 0) {
      obj.last_energy_drop = Math.round(message.last_energy_drop);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EnergySettings>, I>>(base?: I): EnergySettings {
    return EnergySettings.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EnergySettings>, I>>(object: I): EnergySettings {
    const message = createBaseEnergySettings();
    message.active = object.active ?? false;
    message.once = object.once ?? false;
    message.start = object.start ?? 0;
    message.end = object.end ?? 0;
    message.amount = object.amount ?? 0;
    message.last_energy_drop = object.last_energy_drop ?? 0;
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

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
