#!/usr/bin/env python3
"""Discover playable community-data characters missing from internal/characters.

The manifest is deliberately independent from a hand-maintained character list.
It keeps workflow state across runs and treats empty directories as missing.
"""

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DATA = ROOT / "pipeline/community/data"
CHARS = ROOT / "internal/characters"
MANIFEST = ROOT / "character_imports/manifest.json"

TRAVELER = re.compile(r"^(aether|lumine)")
WEAPONS = {
    "WEAPON_SWORD_ONE_HAND": "Sword",
    "WEAPON_CLAYMORE": "Claymore",
    "WEAPON_POLE": "Polearm",
    "WEAPON_CATALYST": "Catalyst",
    "WEAPON_BOW": "Bow",
}


def version_key(value: str) -> tuple[int, ...]:
    return tuple(int(x) for x in re.findall(r"\d+", value))


def records() -> dict[str, dict]:
    latest: dict[str, dict] = {}
    for version_dir in sorted(DATA.iterdir(), key=lambda p: version_key(p.name)):
        for source in (version_dir / "characters").glob("*.json"):
            record = json.loads(source.read_text())
            if record.get("kind") != "character" or record.get("release_track") != "live":
                continue
            ident = record.get("identity", {})
            slug = ident.get("gcsim_slug")
            confirmed = record.get("implementation_inputs", {}).get("confirmed", {})
            identity = confirmed.get("identity", {})
            if not slug or not ident.get("id") or not identity.get("name"):
                continue
            if TRAVELER.match(slug):
                continue
            record["_version"] = version_dir.name
            record["_path"] = source.relative_to(ROOT).as_posix()
            latest[slug] = record
    return latest


def implemented(slug: str) -> bool:
    directory = CHARS / slug
    return (directory / "config.yml").is_file() and any(
        p.suffix == ".go" and not p.name.startswith("zz_") for p in directory.glob("*.go")
    )


def entry(record: dict, old: dict | None) -> dict:
    ident = record["identity"]
    identity = record["implementation_inputs"]["confirmed"]["identity"]
    missing = record["implementation_inputs"].get("unresolved", [])
    slug = ident["gcsim_slug"]
    result = {
        "id": int(str(ident["id"]).split("-")[0]),
        "name": slug,
        "display_name": identity["name"],
        "element": identity["element"].title(),
        "weapon": WEAPONS.get(identity["weapon"], identity["weapon"]),
        "data_version": record["_version"],
        "source_path": record["_path"],
        "lunaris_url": f"https://lunaris.moe/character/{str(ident['id']).split('-')[0]}",
        "nanoka_query": f"{identity['name']} {str(ident['id']).split('-')[0]}",
        "status": "pending",
        "simulation_level": "not-generated",
        "missing_local_fields": list(missing),
        "missing_core_support": detect_core_support(record),
        "references": [],
        "todos": list(missing),
    }
    if old:
        for key in ("status", "simulation_level", "references", "todos"):
            result[key] = old.get(key, result[key])
    if implemented(slug):
        result["status"] = old.get("status", "existing") if old else "existing"
        result["simulation_level"] = old.get("simulation_level", "unknown") if old else "unknown"
    return result


def detect_core_support(record: dict) -> list[str]:
    tables = record["implementation_inputs"]["confirmed"]["combat_tables"]
    text = " ".join(
        item.get("desc", "")
        for group in ("passives", "constellations", "skills")
        for item in tables.get(group, [])
    ).lower()
    terms = {
        "Lunar-Bloom reaction event and formula": ("lunar-bloom",),
        "Moonsign state/query": ("moonsign",),
        "Ascendant Gleam state": ("ascendant gleam",),
        "Verdant Dew team resource": ("verdant dew",),
        "Seed of Deceit entity": ("seed of deceit", "seeds of deceit"),
    }
    return [label for label, needles in terms.items() if any(needle in text for needle in needles)]


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="fail if the manifest would change")
    args = parser.parse_args()
    old_entries = {}
    if MANIFEST.exists():
        old_entries = {e["name"]: e for e in json.loads(MANIFEST.read_text()).get("characters", [])}
    discovered = records()
    document = {
        "schema_version": 1,
        "source": "pipeline/community/data",
        "characters": [entry(discovered[k], old_entries.get(k)) for k in sorted(discovered)],
    }
    rendered = json.dumps(document, indent=2, ensure_ascii=False) + "\n"
    previous = MANIFEST.read_text() if MANIFEST.exists() else ""
    if args.check:
        if previous != rendered:
            raise SystemExit("character import manifest is stale")
        return
    MANIFEST.parent.mkdir(parents=True, exist_ok=True)
    MANIFEST.write_text(rendered)
    missing = [e["name"] for e in document["characters"] if not implemented(e["name"])]
    print(f"discovered={len(document['characters'])} missing={len(missing)}")
    print("\n".join(missing))


if __name__ == "__main__":
    main()
