import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";

export const weaponKeyToName: { [key: string]: string } = {
  akuoumaru: "Akuoumaru",
  alleyhunter: "Alley Hunter",
  amberbead: "Amber Bead",
  amenomakageuchi: "Amenoma Kageuchi",
  amosbow: "Amos' Bow",
  apprenticesnotes: "Apprentice's Notes",
  aquilafavonia: "Aquila Favonia",
  beginnersprotector: "Beginner's Protector",
  blackcliffagate: "Blackcliff Agate",
  blackclifflongsword: "Blackcliff Longsword",
  blackcliffpole: "Blackcliff Pole",
  blackcliffslasher: "Blackcliff Slasher",
  blackcliffwarbow: "Blackcliff Warbow",
  blacktassel: "Black Tassel",
  bloodtaintedgreatsword: "Bloodtainted Greatsword",
  calamityqueller: "Calamity Queller",
  cinnabarspindle: "Cinnabar Spindle",
  compoundbow: "Compound Bow",
  coolsteel: "Cool Steel",
  crescentpike: "Crescent Pike",
  darkironsword: "Dark Iron Sword",
  deathmatch: "Deathmatch",
  debateclub: "Debate Club",
  dodocotales: "Dodoco Tales",
  dragonsbane: "Dragon's Bane",
  dragonspinespear: "Dragonspine Spear",
  dullblade: "Dull Blade",
  ebonybow: "Ebony Bow",
  elegyfortheend: "Elegy for the End",
  emeraldorb: "Emerald Orb",
  engulfinglightning: "Engulfing Lightning",
  everlastingmoonglow: "Everlasting Moonglow",
  eyeofperception: "Eye of Perception",
  favoniuscodex: "Favonius Codex",
  favoniusgreatsword: "Favonius Greatsword",
  favoniuslance: "Favonius Lance",
  favoniussword: "Favonius Sword",
  favoniuswarbow: "Favonius Warbow",
  ferrousshadow: "Ferrous Shadow",
  festeringdesire: "Festering Desire",
  filletblade: "Fillet Blade",
  freedomsworn: "Freedom-Sworn",
  frostbearer: "Frostbearer",
  hakushinring: "Hakushin Ring",
  halberd: "Halberd",
  hamayumi: "Hamayumi",
  harbingerofdawn: "Harbinger of Dawn",
  huntersbow: "Hunter's Bow",
  ironpoint: "Iron Point",
  ironsting: "Iron Sting",
  katsuragikirinagamasa: "Katsuragikiri Nagamasa",
  kitaincrossspear: "Kitain Cross Spear",
  lionsroar: "Lion's Roar",
  lithicblade: "Lithic Blade",
  lithicspear: "Lithic Spear",
  lostprayertothesacredwinds: "Lost Prayer to the Sacred Winds",
  luxurioussealord: "Luxurious Sea-Lord",
  magicguide: "Magic Guide",
  mappamare: "Mappa Mare",
  memoryofdust: "Memory of Dust",
  messenger: "Messenger",
  mistsplitterreforged: "Mistsplitter Reforged",
  mitternachtswaltz: "Mitternachts Waltz",
  mouunsmoon: "Mouun's Moon",
  oldmercspal: "Old Merc's Pal",
  otherworldlystory: "Otherworldly Story",
  pocketgrimoire: "Pocket Grimoire",
  polarstar: "Polar Star",
  predator: "Predator",
  primordialjadecutter: "Primordial Jade Cutter",
  primordialjadewingedspear: "Primordial Jade Winged-Spear",
  prototypeamber: "Prototype Amber",
  prototypearchaic: "Prototype Archaic",
  prototypecrescent: "Prototype Crescent",
  prototyperancour: "Prototype Rancour",
  prototypestarglitter: "Prototype Starglitter",
  quartz: "Quartz",
  rainslasher: "Rainslasher",
  ravenbow: "Raven Bow",
  recurvebow: "Recurve Bow",
  redhornstonethresher: "Redhorn Stonethresher",
  royalbow: "Royal Bow",
  royalgreatsword: "Royal Greatsword",
  royalgrimoire: "Royal Grimoire",
  royallongsword: "Royal Longsword",
  royalspear: "Royal Spear",
  rust: "Rust",
  sacrificialbow: "Sacrificial Bow",
  sacrificialfragments: "Sacrificial Fragments",
  sacrificialgreatsword: "Sacrificial Greatsword",
  sacrificialsword: "Sacrificial Sword",
  seasonedhuntersbow: "Seasoned Hunter's Bow",
  serpentspine: "Serpent Spine",
  sharpshootersoath: "Sharpshooter's Oath",
  silversword: "Silver Sword",
  skyridergreatsword: "Skyrider Greatsword",
  skyridersword: "Skyrider Sword",
  skywardatlas: "Skyward Atlas",
  skywardblade: "Skyward Blade",
  skywardharp: "Skyward Harp",
  skywardpride: "Skyward Pride",
  skywardspine: "Skyward Spine",
  slingshot: "Slingshot",
  snowtombedstarsilver: "Snow-Tombed Starsilver",
  solarpearl: "Solar Pearl",
  songofbrokenpines: "Song of Broken Pines",
  staffofhoma: "Staff of Homa",
  summitshaper: "Summit Shaper",
  swordofdescension: "Sword of Descension",
  thealleyflash: "The Alley Flash",
  thebell: "The Bell",
  theblacksword: "The Black Sword",
  thecatch: '"The Catch"',
  theflagstaff: "The Flagstaff",
  theflute: "The Flute",
  thestringless: "The Stringless",
  theunforged: "The Unforged",
  theviridescenthunt: "The Viridescent Hunt",
  thewidsith: "The Widsith",
  thrillingtalesofdragonslayers: "Thrilling Tales of Dragon Slayers",
  thunderingpulse: "Thundering Pulse",
  travelershandysword: "Traveler's Handy Sword",
  twinnephrite: "Twin Nephrite",
  vortexvanquisher: "Vortex Vanquisher",
  wastergreatsword: "Waster Greatsword",
  wavebreakersfin: "Wavebreaker's Fin",
  whiteblind: "Whiteblind",
  whiteirongreatsword: "White Iron Greatsword",
  whitetassel: "White Tassel",
  windblumeode: "Windblume Ode",
  wineandsong: "Wine and Song",
  wolfsgravestone: "Wolf's Gravestone",
};

export const weapons: IWeapon[] = Object.keys(weaponKeyToName).map((k) => ({
  key: k,
  name: weaponKeyToName[k],
}));

export interface IWeapon {
  key: string;
  name: string;
}

export const renderWeapon: ItemRenderer<IWeapon> = (
  weapon,
  { handleClick, modifiers, query }
) => {
  if (!modifiers.matchesPredicate) {
    return null;
  }
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      label={weapon.name}
      key={weapon.key}
      onClick={handleClick}
      text={highlightText(weapon.name, query)}
    />
  );
};

export const filterWeapon: ItemPredicate<IWeapon> = (
  query,
  weapon,
  _index,
  exactMatch
) => {
  const normalizedQuery = query.toLowerCase();

  if (exactMatch) {
    return weapon.key === normalizedQuery;
  } else {
    return `${weapon.key} ${weapon.name}`.indexOf(normalizedQuery) >= 0;
  }
};

function escapeRegExpChars(text: string) {
  return text.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, "\\$1");
}

function highlightText(text: string, query: string) {
  let lastIndex = 0;
  const words = query
    .split(/\s+/)
    .filter((word) => word.length > 0)
    .map(escapeRegExpChars);
  if (words.length === 0) {
    return [text];
  }
  const regexp = new RegExp(words.join("|"), "gi");
  const tokens: React.ReactNode[] = [];
  while (true) {
    const match = regexp.exec(text);
    if (!match) {
      break;
    }
    const length = match[0].length;
    const before = text.slice(lastIndex, regexp.lastIndex - length);
    if (before.length > 0) {
      tokens.push(before);
    }
    lastIndex = regexp.lastIndex;
    tokens.push(<strong key={lastIndex}>{match[0]}</strong>);
  }
  const rest = text.slice(lastIndex);
  if (rest.length > 0) {
    tokens.push(rest);
  }
  return tokens;
}

export const weaponSelectProps = {
  itemPredicate: filterWeapon,
  itemRenderer: renderWeapon,
  items: weapons,
};
