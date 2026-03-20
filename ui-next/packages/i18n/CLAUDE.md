# @gcsim/i18n

Internationalization package for the gcsim UI. Wraps i18next + react-i18next with pre-loaded translations for 7 languages.

## Structure

- `src/locales/` — JSON translation files (7 UI languages + generated game name files)
- `src/resources.ts` — builds the i18next resources object from JSON imports
- `src/init.ts` — `initI18n()` function to initialize i18next with react-i18next
- `src/index.ts` — barrel exports

## Supported Languages

| Code | Language |
|------|----------|
| en | English |
| zh | Chinese |
| ja | Japanese |
| ko | Korean |
| es | Spanish |
| ru | Russian |
| de | German |

## Namespaces

- `translation` — UI strings (from `<Language>.json` files)
- `game` — character/weapon/artifact names (merged from `names.generated.json` + `names.traveler.json`)

## Usage

```typescript
import { initI18n } from "@gcsim/i18n";

// Call once at app startup
await initI18n("en");

// Then use react-i18next hooks in components
import { useTranslation } from "react-i18next";
const { t } = useTranslation();
t("nav.simulator"); // UI string
t("game:character_names.albedo"); // Game name
```

## Special Locales

`specialLocales` array contains locale codes that need special handling (CJK): `["zh", "ja", "ko"]`.
