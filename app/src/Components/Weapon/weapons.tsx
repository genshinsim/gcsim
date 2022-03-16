import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";
import i18n from "i18next";

export const weapons: IWeapon[] = [
  "akuoumaru",
  "alleyhunter",
  "amenomakageuchi",
  "amosbow",
  "apprenticesnotes",
  "aquilafavonia",
  "beginnersprotector",
  "blackcliffagate",
  "blackclifflongsword",
  "blackcliffpole",
  "blackcliffslasher",
  "blackcliffwarbow",
  "blacktassel",
  "bloodtaintedgreatsword",
  "calamityqueller",
  "cinnabarspindle",
  "compoundbow",
  "coolsteel",
  "crescentpike",
  "darkironsword",
  "deathmatch",
  "debateclub",
  "dodocotales",
  "dragonsbane",
  "dragonspinespear",
  "dullblade",
  "elegyfortheend",
  "emeraldorb",
  "engulfinglightning",
  "everlastingmoonglow",
  "eyeofperception",
  "favoniuscodex",
  "favoniusgreatsword",
  "favoniuslance",
  "favoniussword",
  "favoniuswarbow",
  "ferrousshadow",
  "festeringdesire",
  "filletblade",
  "freedomsworn",
  "frostbearer",
  "hakushinring",
  "halberd",
  "hamayumi",
  "harbingerofdawn",
  "huntersbow",
  "ironpoint",
  "ironsting",
  "kagurasverity",
  "katsuragikirinagamasa",
  "kitaincrossspear",
  "lionsroar",
  "lithicblade",
  "lithicspear",
  "lostprayertothesacredwinds",
  "luxurioussealord",
  "magicguide",
  "mappamare",
  "memoryofdust",
  "messenger",
  "mistsplitterreforged",
  "mitternachtswaltz",
  "mouunsmoon",
  "oathsworneye",
  "oldmercspal",
  "otherworldlystory",
  "pocketgrimoire",
  "polarstar",
  "predator",
  "primordialjadecutter",
  "primordialjadewingedspear",
  "prototypeamber",
  "prototypearchaic",
  "prototypecrescent",
  "prototyperancour",
  "prototypestarglitter",
  "rainslasher",
  "ravenbow",
  "recurvebow",
  "redhornstonethresher",
  "royalbow",
  "royalgreatsword",
  "royalgrimoire",
  "royallongsword",
  "royalspear",
  "rust",
  "sacrificialbow",
  "sacrificialfragments",
  "sacrificialgreatsword",
  "sacrificialsword",
  "seasonedhuntersbow",
  "serpentspine",
  "sharpshootersoath",
  "silversword",
  "skyridergreatsword",
  "skyridersword",
  "skywardatlas",
  "skywardblade",
  "skywardharp",
  "skywardpride",
  "skywardspine",
  "slingshot",
  "snowtombedstarsilver",
  "solarpearl",
  "songofbrokenpines",
  "staffofhoma",
  "summitshaper",
  "swordofdescension",
  "thealleyflash",
  "thebell",
  "theblacksword",
  "thecatch",
  "theflute",
  "thestringless",
  "theunforged",
  "theviridescenthunt",
  "thewidsith",
  "thrillingtalesofdragonslayers",
  "thunderingpulse",
  "travelershandysword",
  "twinnephrite",
  "vortexvanquisher",
  "wastergreatsword",
  "wavebreakersfin",
  "whiteblind",
  "whiteirongreatsword",
  "whitetassel",
  "windblumeode",
  "wineandsong",
  "wolfsgravestone",
];

export type IWeapon = string;

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
      label={""}
      key={weapon}
      onClick={handleClick}
      text={highlightText(i18n.t("weapon_names." + weapon), query)}
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
    return weapon === normalizedQuery;
  } else {
    return (
      `${weapon} ${i18n.t("weapon_names." + weapon)}`.indexOf(
        normalizedQuery
      ) >= 0
    );
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
