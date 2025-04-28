import {Character} from '@gcsim/types';

const statKeys = [
  'n/a',
  'def%',
  'def',
  'hp',
  'hp%',
  'atk',
  'atk%',
  'er',
  'em',
  'cr',
  'cd',
  'heal',
  'pyro%',
  'hydro%',
  'cryo%',
  'electro%',
  'anemo%',
  'geo%',
  'dendro%',
  'phys%',
  'atkspd%',
  'dmg%',
];

export function charToCfg(char: Character): string {
  let str = '';
  // prettier-ignore
  str += `${char.name} char lvl=${char.level}/${char.max_level} cons=${char.cons} talent=${char.talents.attack},${char.talents.skill},${char.talents.burst};\n`;
  // prettier-ignore
  str += `${char.name} add weapon="${char.weapon.name}" refine=${char.weapon.refine} lvl=${char.weapon.level}/${char.weapon.max_level};\n`;

  //build sets
  for (const key in char.sets) {
    if (char.sets[key] > 0) {
      str += `${char.name} add set="${key}" count=${char.sets[key]};\n`;
    }
  }

  //add stats
  let count = 0;
  let statStr = `${char.name} add stats`;
  char.stats.forEach(([index, value], i) => {
    if (value === 0) return;
    count++;
    statStr += ` ${statKeys[index]}=${value.toPrecision()}`;

    // Add ";\n" after main stats, then "\n" after every artifact piece substats
    if (i === 4) {
      str += statStr + '; # main stats\n';
      statStr = `${char.name} add stats`;
    } else if ((count - 5) % 4 === 0 && i > 4) {
      str += statStr + ';\n';
      statStr = `${char.name} add stats`;
    }
  });

  // \n delimiting is under assumption that each artifact has 4 subs.
  // If that's not the case, add to the end.
  if (count > 0 && (count - 5) % 4 !== 0) {
    str += statStr + `;\n`;
  }

  return str;
}
