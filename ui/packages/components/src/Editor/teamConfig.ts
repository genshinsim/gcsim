import type { model } from "@gcsim/types";

// TEMPORARY: backs the TeamView add/remove crutch; delete when TeamView is view-only.

const statKeys = [
	"n/a",
	"def%",
	"def",
	"hp",
	"hp%",
	"atk",
	"atk%",
	"er",
	"em",
	"cr",
	"cd",
	"heal",
	"pyro%",
	"hydro%",
	"cryo%",
	"electro%",
	"anemo%",
	"geo%",
	"dendro%",
	"phys%",
	"atkspd%",
	"dmg%",
];

export const charLinesRegEx =
	/^(\w+) (?:char|add) (?:lvl|weapon|set|stats).+$(?:\r\n|\r|\n)?/gm;

export function charToCfg(char: model.Character): string {
	const name = char.name ?? "";
	const talents = char.talents ?? {};
	const weapon = char.weapon ?? {};
	const sets = char.sets ?? {};

	let str = "";
	str += `${name} char lvl=${char.level}/${char.max_level} cons=${char.cons} talent=${talents.attack},${talents.skill},${talents.burst};\n`;
	str += `${name} add weapon="${weapon.name}" refine=${weapon.refine} lvl=${weapon.level}/${weapon.max_level};\n`;

	for (const key in sets) {
		if (sets[key] > 0) {
			str += `${name} add set="${key}" count=${sets[key]};\n`;
		}
	}

	let count = 0;
	let statStr = `${name} add stats`;
	(char.stats ?? []).forEach((v, i) => {
		if (v === 0) return;
		count++;
		statStr += ` ${statKeys[i]}=${v.toPrecision()}`;
	});
	if (count > 0) {
		str += `${statStr};\n`;
	}

	return str;
}

export function cfgFromTeam(team: model.Character[], cfg: string): string {
	let next = "";
	team.forEach((c) => {
		next += `${charToCfg(c)}\n`;
	});

	let out = cfg.replace(charLinesRegEx, "");
	out = next + out;
	out = out.replace(/(\r\n|\r|\n){2,}/g, "$1\n");

	return out;
}
