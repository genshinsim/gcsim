import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import type { Page } from "@playwright/test";

/**
 * A minimal, parseable GOOD payload (one character, Bennett) the import spec
 * pastes into the GOOD dialog. Committed as a plain JSON file so it reads like a
 * payload a user would paste (see also `sucroseConfig`).
 */
export const goodImport: string = readFileSync(
	fileURLToPath(new URL("../fixtures/good-import.json", import.meta.url)),
	"utf8",
);

/**
 * A minimal Enka.Network `/api/enka/:uid` payload the import dialog can parse
 * into exactly one character (Bennett). Only the fields `EnkaToGOOD` reads are
 * present: the avatar id resolves to a known char, the equip list carries one
 * recognised weapon (Aquila Favonia), and the prop/skill maps drive level and
 * talents. No reliquaries — artifact sets and stats stay at their empty base,
 * which is a valid import.
 */
export const ENKA_UID = "700000000";

export const enkaImportPayload = [
	{
		avatarId: 10000032, // Bennett
		skillDepotId: 3201,
		talentIdList: [],
		propMap: {
			"4001": { val: "90" }, // level
			"1002": { val: "6" }, // ascension
		},
		skillLevelMap: {
			"10321": 1, // attack
			"10322": 8, // skill
			"10323": 10, // burst
		},
		equipList: [
			{
				itemId: 11501, // Aquila Favonia
				weapon: { level: 90, promoteLevel: 6, affixMap: { "111501": 4 } },
				flat: { itemType: "ITEM_WEAPON" },
			},
		],
	},
];

/**
 * Stub the Enka import endpoint so the import dialog resolves offline and
 * deterministically. The dev server proxies `/api/enka/:uid` to production;
 * this returns {@link enkaImportPayload} for {@link ENKA_UID} instead.
 */
export async function installEnkaRoutes(page: Page): Promise<void> {
	await page.route(`**/api/enka/${ENKA_UID}`, (route) =>
		route.fulfill({
			contentType: "application/json",
			body: JSON.stringify(enkaImportPayload),
		}),
	);
}
