#!/usr/bin/env node
import { readdirSync, readFileSync } from "node:fs";
import { join, relative } from "node:path";
import { fileURLToPath } from "node:url";

const ROOT = join(fileURLToPath(new URL(".", import.meta.url)), "..");
const PACKAGES = join(ROOT, "packages");

const SKIP_DIRS = new Set([
	"node_modules",
	"dist",
	"storybook-static",
	".turbo",
	"playwright-report",
	"test-results",
]);
const SKIP_FILES = new Set(["stats.html"]);
const EXTS = [".ts", ".tsx", ".css", ".html"];

const COLOR_TOKENS = [
	"background",
	"foreground",
	"card",
	"popover",
	"primary",
	"secondary",
	"muted",
	"accent",
	"destructive",
	"success",
	"warning",
	"border",
	"input",
	"ring",
];
const CUSTOM_PROP_TOKENS = [...COLOR_TOKENS, "radius"];
const COLOR_UTILS = [
	"bg",
	"text",
	"border",
	"ring-offset",
	"ring",
	"outline",
	"fill",
	"stroke",
];

const colorTokenAlt = COLOR_TOKENS.join("|");
const customPropTokenAlt = CUSTOM_PROP_TOKENS.join("|");
const colorUtilAlt = COLOR_UTILS.join("|");
const BARE_CUSTOM_PROP = new RegExp(
	`(?<![A-Za-z0-9_])--(?:${customPropTokenAlt})\\b`,
	"g",
);
const BARE_UTILITY = new RegExp(
	`(?<![A-Za-z0-9_])(?:${colorUtilAlt})-(?:${colorTokenAlt})\\b`,
	"g",
);

function* walk(dir) {
	for (const entry of readdirSync(dir, { withFileTypes: true })) {
		if (entry.isDirectory()) {
			if (!SKIP_DIRS.has(entry.name)) yield* walk(join(dir, entry.name));
		} else if (
			!SKIP_FILES.has(entry.name) &&
			EXTS.some((ext) => entry.name.endsWith(ext))
		) {
			yield join(dir, entry.name);
		}
	}
}

const violations = [];
for (const file of walk(PACKAGES)) {
	const lines = readFileSync(file, "utf8").split("\n");
	lines.forEach((line, i) => {
		for (const re of [BARE_CUSTOM_PROP, BARE_UTILITY]) {
			re.lastIndex = 0;
			let m = re.exec(line);
			while (m !== null) {
				violations.push({
					file: relative(ROOT, file),
					line: i + 1,
					match: m[0],
				});
				m = re.exec(line);
			}
		}
	});
}

if (violations.length > 0) {
	console.error(
		`Found ${violations.length} bare shadcn token(s)/utilities — they must use the --deprecated- form:`,
	);
	for (const v of violations)
		console.error(`  ${v.file}:${v.line}  ${v.match}`);
	process.exit(1);
}

console.log("OK: no bare shadcn tokens or utilities found.");
