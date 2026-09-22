import type { Character } from "@gcsim/types";

// TEMPORARY. These team<->config serializers back the TeamView add/remove crutch
// that only exists so GOOD/Enka imports have somewhere to land. Once TeamView
// becomes view-only and imports get a first-class path, delete this file.

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

export function charToCfg(char: Character): string {
	let str = "";
	str += `${char.name} char lvl=${char.level}/${char.max_level} cons=${char.cons} talent=${char.talents.attack},${char.talents.skill},${char.talents.burst};\n`;
	str += `${char.name} add weapon="${char.weapon.name}" refine=${char.weapon.refine} lvl=${char.weapon.level}/${char.weapon.max_level};\n`;

	for (const key in char.sets) {
		if (char.sets[key] > 0) {
			str += `${char.name} add set="${key}" count=${char.sets[key]};\n`;
		}
	}

	let count = 0;
	let statStr = `${char.name} add stats`;
	char.stats.forEach((v, i) => {
		if (v === 0) return;
		count++;
		statStr += ` ${statKeys[i]}=${v.toPrecision()}`;
	});
	if (count > 0) {
		str += `${statStr};\n`;
	}

	return str;
}

export function cfgFromTeam(team: Character[], cfg: string): string {
	let next = "";
	team.forEach((c) => {
		next += `${charToCfg(c)}\n`;
	});

	let out = cfg.replace(charLinesRegEx, "");
	out = next + out;
	out = out.replace(/(\r\n|\r|\n){2,}/g, "$1\n");

	return out;
}
