# @gcsim/data

## Purpose

Provides static game data (characters, tags) with typed exports for use across the gcsim UI.

## Usage

```typescript
import { tags, latestChars } from "@gcsim/data";
import type { TagInfo, TagMap, LatestCharsMap } from "@gcsim/data";

// Access tag info by numeric ID
const gcsimTag: TagInfo = tags["1"];
console.log(gcsimTag.display_name); // "gcsim"

// Access latest characters by version
const newChars: string[] = latestChars["v2.38"];
```

## How to Add Data

### Adding a new version's characters

Edit `src/latest-chars.json` and add a new version entry:

```json
{
  "v2.38": ["varesa", "luminepyro"],
  "v2.39": ["newcharacter"]
}
```

### Adding a new tag

Edit `src/tags.json` and add a new numeric key:

```json
{
  "10": {
    "display_name": "My Tag",
    "blurb": "Optional description of the tag."
  }
}
```

### Adding a new data category

1. Create the JSON file in `src/` (kebab-case name)
2. Create a typed TypeScript module in `src/` that imports the JSON and exports it with proper types
3. Re-export from `src/index.ts`
4. Add tests in `src/__tests__/data.test.ts`

## Don'ts

- Don't import from internal paths (`@gcsim/data/src/tags`) -- use the package index only
- Don't add runtime logic here -- this package is data only
- Don't manually define data that should come from protobuf types -- use `@gcsim/types` instead
