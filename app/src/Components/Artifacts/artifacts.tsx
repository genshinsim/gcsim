import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";
import i18n from "i18next";

export const items = [
  "archaicpetra",
  "berserker",
  "blizzardstrayer",
  "bloodstainedchivalry",
  "braveheart",
  "crimsonwitchofflames",
  "defenderswill",
  "echoesofanoffering",
  "emblemofseveredfate",
  "gambler",
  "gladiatorsfinale",
  "heartofdepth",
  "huskofopulentdreams",
  "instructor",
  "lavawalker",
  "maidenbeloved",
  "martialartist",
  "noblesseoblige",
  "oceanhuedclam",
  "paleflame",
  "prayersfordestiny",
  "prayersforillumination",
  "prayersforwisdom",
  "prayerstospringtime",
  "resolutionofsojourner",
  "retracingbolide",
  "scholar",
  "shimenawasreminiscence",
  "tenacityofthemillelith",
  "theexile",
  "thunderingfury",
  "thundersoother",
  "tinymiracle",
  "vermillionhereafter",
  "viridescentvenerer",
  "wandererstroupe",
];

export type IArtifact = string;

export const render: ItemRenderer<IArtifact> = (
  item,
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
      key={item}
      onClick={handleClick}
      text={highlightText(i18n.t("artifact_names." + item), query)}
    />
  );
};

export const filter: ItemPredicate<IArtifact> = (
  query,
  item,
  _index,
  exactMatch
) => {
  const normalizedQuery = query.toLowerCase();
  const transItem = i18n
    .t("artifact_names." + item)
    .replace(" ", "")
    .toLowerCase();
  if (exactMatch) {
    return item === normalizedQuery;
  } else {
    return `${item} ${transItem}`.indexOf(normalizedQuery) >= 0;
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

export const artifactSelectProps = {
  itemPredicate: filter,
  itemRenderer: render,
  items: items,
};
