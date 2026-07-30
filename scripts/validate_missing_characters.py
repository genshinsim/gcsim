#!/usr/bin/env python3
"""Structural review for generated missing-character packages."""

from __future__ import annotations

import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
MANIFEST = ROOT / "character_imports/manifest.json"
REQUIRED = {"attack.go", "charge.go", "plunge.go", "skill.go", "burst.go", "asc.go", "cons.go", "config.yml"}


def main() -> None:
    doc = json.loads(MANIFEST.read_text())
    failures = []
    for entry in doc["characters"]:
        if entry["status"] == "existing":
            continue
        slug = entry["name"]
        directory = ROOT / "internal/characters" / slug
        expected = REQUIRED | {f"{slug}.go", f"zz_{slug}.dm.go"}
        missing = sorted(name for name in expected if not (directory / name).is_file())
        if missing:
            failures.append(f"{slug}: missing {', '.join(missing)}")
            continue
        config = (directory / "config.yml").read_text()
        source = json.loads((ROOT / entry["source_path"]).read_text())
        labels = {
            label
            for skill in source["implementation_inputs"]["confirmed"]["combat_tables"]["skills"]
            for level in skill["promote"].values()
            for label in level["desc"]
            if label
        }
        configured = re.findall(r"(?:attack|skill|burst)\((.*?)\)", config)
        for label in configured:
            label = re.sub(r"\|[01]$", "", label).replace("''", "'")
            if label not in labels:
                failures.append(f"{slug}: config label absent locally: {label}")
        generated = (directory / f"zz_{slug}.dm.go").read_text()
        arrays = re.findall(r"Values:\s*\[\]float64\{([^}]+)\}", generated)
        if not arrays or any(len(a.split(",")) != 15 for a in arrays):
            failures.append(f"{slug}: generated multiplier array is not 15 levels")
        if f"keys.{slug.capitalize()}" not in generated:
            failures.append(f"{slug}: registration key missing")
    if failures:
        raise SystemExit("\n".join(failures))
    print(f"validated={sum(e['status'] != 'existing' for e in doc['characters'])}")


if __name__ == "__main__":
    main()
