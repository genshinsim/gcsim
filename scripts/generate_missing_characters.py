#!/usr/bin/env python3
"""Generate conservative character packages from the discovered local records."""

from __future__ import annotations

import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
MANIFEST = ROOT / "character_imports/manifest.json"
ELEMENTS = {"Anemo": "Anemo", "Geo": "Geo", "Electro": "Electro", "Dendro": "Dendro", "Cryo": "Cryo", "Pyro": "Pyro", "Hydro": "Hydro"}
MODEL_ELEMENTS = {"Anemo": "Wind", "Geo": "Rock", "Electro": "Electric", "Dendro": "Grass", "Cryo": "Ice", "Pyro": "Fire", "Hydro": "Water"}
MODEL_WEAPONS = {"Sword": "WEAPON_SWORD_ONE_HAND", "Claymore": "WEAPON_CLAYMORE", "Polearm": "WEAPON_POLE", "Catalyst": "WEAPON_CATALYST", "Bow": "WEAPON_BOW"}


def go_name(slug: str) -> str:
    return "".join(part.capitalize() for part in re.split(r"[^a-zA-Z0-9]", slug))


def promote(skill: dict) -> list[dict]:
    return [skill["promote"][str(i)] for i in range(15)]


def params(label: str) -> list[int]:
    return [int(v) - 1 for v in re.findall(r"\{param(\d+)", label)]


def scaling(skill: dict, match, occurrence=0) -> dict:
    labels = promote(skill)[0]["desc"]
    found = [x for x in labels if x and match(x)]
    label = found[occurrence] if len(found) > occurrence else ""
    indexes = params(label)
    values = [sum(level["param"][i] for i in indexes) for level in promote(skill)] if indexes else [0.0] * 15
    upper = label.upper()
    return {"label": label, "values": values, "def": " DEF" in upper, "hp": " MAX HP" in upper or " HP" in upper, "em": "ELEMENTAL MASTERY" in upper}


def seconds(skill: dict, needle: str, default: int) -> int:
    labels = promote(skill)[0]["desc"]
    for label in labels:
        if needle.lower() in label.lower():
            ids = params(label)
            if ids:
                return round(promote(skill)[0]["param"][ids[0]] * 60)
    return default * 60


def normal_scalings(attack: dict) -> list[dict]:
    labels = promote(attack)[0]["desc"]
    out = []
    for label in labels:
        if re.match(r"\d+-Hit DMG", label):
            indexes = params(label)
            out.append({"label": label, "values": [sum(level["param"][i] for i in indexes) for level in promote(attack)], "def": False, "hp": False, "em": False})
    return out or [scaling(attack, lambda x: "DMG" in x)]


def literal(s: dict) -> str:
    vals = ", ".join(f"{v:.8g}" for v in s["values"])
    return f'basicimport.Scaling{{Values: []float64{{{vals}}}, UseDef: {str(s["def"]).lower()}, UseHP: {str(s["hp"]).lower()}, UseEM: {str(s["em"]).lower()}}}'


def write_package(entry: dict) -> None:
    source = ROOT / entry["source_path"]
    record = json.loads(source.read_text())
    identity = record["implementation_inputs"]["confirmed"]["identity"]
    tables = record["implementation_inputs"]["confirmed"]["combat_tables"]
    attack, skill, burst = tables["skills"]
    normals = normal_scalings(attack)
    charge = scaling(attack, lambda x: x.startswith("Charged Attack DMG"))
    plunge = scaling(attack, lambda x: x.startswith("Plunge DMG"))
    low = scaling(attack, lambda x: x.startswith("Low/High Plunge DMG"))
    high = scaling(attack, lambda x: x.startswith("Low/High Plunge DMG"))
    low_ids = params(low["label"])
    if len(low_ids) >= 2:
        low["values"] = [level["param"][low_ids[0]] for level in promote(attack)]
        high["values"] = [level["param"][low_ids[1]] for level in promote(attack)]
    skill_mv = scaling(skill, lambda x: "DMG" in x and bool(params(x)))
    burst_mv = scaling(burst, lambda x: "DMG" in x and bool(params(x)))
    energy = round(seconds(burst, "Energy Cost", 1) / 60)
    slug = entry["name"]
    directory = ROOT / "internal/characters" / slug
    directory.mkdir(parents=True, exist_ok=True)
    config = f"""use: pipeline
kind: character
name: {slug}
source: community
version: {entry['data_version']}
override:
  id: {entry['id']}

# Exact local labels used by the community-data generator.
talents:
  - attack:
"""
    for n in normals:
        label = n["label"].replace("'", "''")
        config += f"      - 'attack({label})'\n"
    optional = (
        ("charge", "attack", charge, ""),
        ("collision", "attack", plunge, ""),
        ("lowPlunge", "attack", low, "|0"),
        ("highPlunge", "attack", high, "|1"),
        ("skill", "skill", skill_mv, ""),
        ("burst", "burst", burst_mv, ""),
    )
    for variable, talent, value, suffix in optional:
        if value["label"]:
            label = value["label"].replace("'", "''")
            config += f"  - {variable}: '{talent}({label}{suffix})'\n"
    (directory / "config.yml").write_text(config)
    base = f"""package {slug}

import (
\t"github.com/genshinsim/gcsim/internal/template/character/basicimport"
\t"github.com/genshinsim/gcsim/pkg/core"
\t"github.com/genshinsim/gcsim/pkg/core/info"
\t"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {{ *basicimport.Character }}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {{
\tc := &char{{Character: basicimport.New(s, w, generatedProfile)}}
\tw.Character = c
\treturn nil
}}
"""
    (directory / f"{slug}.go").write_text(base)
    for name, text in {
        "attack.go": "// Normal attacks use locally generated talent multipliers via basicimport.\n",
        "charge.go": "// Charged attacks use locally generated talent multipliers via basicimport.\n",
        "plunge.go": "// Plunge attacks use locally generated talent multipliers via basicimport.\n",
        "skill.go": "// TODO: replace conservative direct-hit timing with verified runtime behavior.\n",
        "burst.go": "// TODO: replace conservative direct-hit timing with verified runtime behavior.\n",
        "asc.go": "// TODO: connect confirmed passive descriptions after their required core events exist.\n",
        "cons.go": "// TODO: connect confirmed constellation descriptions without inventing missing core APIs.\n",
    }.items():
        (directory / name).write_text(f"package {slug}\n\n{text}")
    profile = f"""// Code generated by scripts/generate_missing_characters.py; DO NOT EDIT.

package {slug}

import (
\t"github.com/genshinsim/gcsim/internal/template/character/basicimport"
\t"github.com/genshinsim/gcsim/pkg/core"
\t"github.com/genshinsim/gcsim/pkg/core/attributes"
\t"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {{ core.RegisterCharFunc(keys.{go_name(slug)}, NewChar) }}

var generatedProfile = basicimport.Profile{{
\tName: {json.dumps(identity['name'])}, Element: attributes.{ELEMENTS[identity['element']]}, Weapon: {json.dumps(entry['weapon'])},
\tEnergy: {energy}, SkillCD: {seconds(skill, 'CD', 10)}, BurstCD: {seconds(burst, 'CD', 15)},
\tAttack: []basicimport.Scaling{{{", ".join(literal(v) for v in normals)}}},
\tCharge: {literal(charge)}, Collision: {literal(plunge)}, LowPlunge: {literal(low)}, HighPlunge: {literal(high)},
\tSkill: {literal(skill_mv)}, Burst: {literal(burst_mv)},
}}
"""
    (directory / f"zz_{slug}.dm.go").write_text(profile)


def update_generated(entries: list[dict]) -> None:
    slugs = [e["name"] for e in entries]
    key_path = ROOT / "pkg/core/keys/character.dm.go"
    text = key_path.read_text()
    for slug in slugs:
        name = go_name(slug)
        if re.search(rf"^\s*{name}\s+", text, re.M):
            continue
        text = text.replace("\tInvalidChar", f"\t{name:<30} // {slug}\n\tInvalidChar", 1)
        text = text.replace('\t"invalidchar",', f'\t"{slug}",\n\t"invalidchar",', 1)
        marker = text.rfind("\tInvalidChar,")
        text = text[:marker] + f"\t{name},\n" + text[marker:]
    key_path.write_text(text)
    shortcut = ROOT / "pkg/shortcut/character.dm.go"
    text = shortcut.read_text()
    for slug in slugs:
        line = f'\t"{slug}": keys.{go_name(slug)},\n'
        if not re.search(rf'^\s*"{re.escape(slug)}"\s*:', text, re.M):
            marker = text.rfind("}\n")
            text = text[:marker] + line + text[marker:]
    shortcut.write_text(text)
    imports = ROOT / "pkg/simulation/imports.character.dm.go"
    text = imports.read_text()
    for slug in slugs:
        line = f'\t_ "github.com/genshinsim/gcsim/internal/characters/{slug}"\n'
        if line not in text:
            text = text.replace(")\n", line + ")\n", 1)
    imports.write_text(text)


def write_catalog(entries: list[dict]) -> None:
    blocks = []
    for entry in entries:
        record = json.loads((ROOT / entry["source_path"]).read_text())
        confirmed = record["implementation_inputs"]["confirmed"]
        identity = confirmed["identity"]
        stats = confirmed["base_stats"]
        skills = confirmed["combat_tables"]["skills"]
        energy = round(seconds(skills[2], "Energy Cost", 1) / 60)
        blocks.append(f"""\tkeys.{go_name(entry['name'])}: {{
\t\tId: {entry['id']}, Key: {json.dumps(entry['name'])},
\t\tRarity: model.QualityType_{identity['rarity']},
\t\tElement: model.ElementType_{MODEL_ELEMENTS[identity['element']]},
\t\tWeaponClass: model.WeaponType_{MODEL_WEAPONS[entry['weapon']]},
\t\tIconName: {json.dumps(identity['icon'])},
\t\tStats: &model.AvatarStatsData{{BaseHp: {stats['base_hp']}, BaseAtk: {stats['base_atk']}, BaseDef: {stats['base_def']}, ElementMastery: {stats.get('elemental_mastery', 0)}}},
\t\tSkillDetails: &model.AvatarSkillsData{{Attack: {skills[0]['id']}, Skill: {skills[1]['id']}, Burst: {skills[2]['id']}, BurstEnergyCost: {energy}}},
\t}},""")
    output = f"""// Code generated by scripts/generate_missing_characters.py; DO NOT EDIT.

package catalog

import (
\t"github.com/genshinsim/gcsim/pkg/core/keys"
\t"github.com/genshinsim/gcsim/pkg/model"
)

var CommunityCharacterMap = map[keys.Char]*model.AvatarData{{
{chr(10).join(blocks)}
}}

func init() {{
\tfor key, value := range CommunityCharacterMap {{
\t\tCharacterMap[key] = value
\t}}
}}
"""
    (ROOT / "pkg/catalog/character_community.dm.go").write_text(output)


def main() -> None:
    doc = json.loads(MANIFEST.read_text())
    entries = [e for e in doc["characters"] if e["status"] != "existing"]
    for entry in entries:
        write_package(entry)
        if entry["status"] == "pending":
            entry["status"] = "generated"
            entry["simulation_level"] = "partial-simulation" if entry["missing_core_support"] else "basic-simulation-ready"
        entry["todos"] = sorted(set(entry["todos"] + [
            "Verify frames, hitboxes, ICD, hitlag, particle behavior, and runtime state ordering",
            "Implement confirmed passives and constellations beyond the conservative direct-damage baseline",
        ]))
    update_generated(entries)
    write_catalog([e for e in doc["characters"] if e["status"] != "existing"])
    MANIFEST.write_text(json.dumps(doc, indent=2, ensure_ascii=False) + "\n")
    print(f"generated={len(entries)}")


if __name__ == "__main__":
    main()
