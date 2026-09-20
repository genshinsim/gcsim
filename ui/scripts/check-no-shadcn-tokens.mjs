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

const ALLOWLIST_MARKER = "TODO(gauge-migration)";

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
// `(?!-g-)` lets the Gauge `-g-` namespace through (e.g. `--radius-g-btn`, a
// coined Tailwind token) while still catching every bare shadcn `--radius` etc.
const BARE_CUSTOM_PROP = new RegExp(
	`(?<![A-Za-z0-9_])--(?:${customPropTokenAlt})(?!-g-)\\b`,
	"g",
);
const BARE_UTILITY = new RegExp(
	`(?<![A-Za-z0-9_])(?:${colorUtilAlt})-(?:${colorTokenAlt})\\b`,
	"g",
);
const DEPRECATED_REF = /(?<![A-Za-z0-9_])-?-?deprecated-/g;

const CHECKS = [
	{ re: BARE_CUSTOM_PROP, kind: "bare shadcn token" },
	{ re: BARE_UTILITY, kind: "bare shadcn utility" },
	{ re: DEPRECATED_REF, kind: "deprecated- reference" },
];

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
	const src = readFileSync(file, "utf8");
	if (src.includes(ALLOWLIST_MARKER)) continue;
	const lines = src.split("\n");
	lines.forEach((line, i) => {
		for (const { re, kind } of CHECKS) {
			re.lastIndex = 0;
			let m = re.exec(line);
			while (m !== null) {
				violations.push({
					file: relative(ROOT, file),
					line: i + 1,
					match: m[0],
					kind,
				});
				m = re.exec(line);
			}
		}
	});
}

if (violations.length > 0) {
	console.error(
		`Found ${violations.length} non-Gauge token reference(s) — use the -g- Gauge namespace (or tag a not-shipped file with // ${ALLOWLIST_MARKER}: to allowlist known debt):`,
	);
	for (const v of violations)
		console.error(`  ${v.file}:${v.line}  ${v.match}  (${v.kind})`);
	process.exit(1);
}

console.log("OK: only Gauge (-g-) tokens in use; no deprecated layer remains.");
