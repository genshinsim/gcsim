import { Character } from "~src/types";

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
  "phys%",
  "dendro%",
  "atkspd%",
  "dmg%",
];

export function charToCfg(char: Character): string {
  let str = "";
  // prettier-ignore
  str += `${char.name} char lvl=${char.level}/${char.max_level} cons=${char.cons} talent=${char.talents.attack},${char.talents.skill},${char.talents.burst};\n`
  // prettier-ignore
  str += `${char.name} add weapon="${char.weapon.name}" refine=${char.weapon.refine} lvl=${char.weapon.level}/${char.weapon.max_level};\n`

  //build sets
  for (const key in char.sets) {
    if (char.sets[key] > 0) {
      str += `${char.name} add set="${key}" count=${char.sets[key]};\n`;
    }
  }

  //add stats
  let count = 0;
  let statStr = `${char.name} add stats`;
  char.stats.forEach((v, i) => {
    if (v === 0) return;
    count++;
    statStr += ` ${statKeys[i]}=${v.toPrecision()}`;
  });
  if (count > 0) {
    str += statStr + `;\n`;
  }

  return str;
}
