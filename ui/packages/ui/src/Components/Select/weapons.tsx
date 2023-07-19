import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";
import { IWeapon } from "@gcsim/types";
import i18n from "i18next";
import { valid_weapons } from "../../Data";

export const weapons: IWeapon[] = valid_weapons.map(e => e)

export const renderWeapon: ItemRenderer<IWeapon> = (weapon, { handleClick, modifiers, query }) => {
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
      text={highlightText(i18n.t("game:weapon_names." + weapon), query)}
    />
  );
};

export const filterWeapon: ItemPredicate<IWeapon> = (query, weapon, _index, exactMatch) => {
  const normalizedQuery = query.toLowerCase();
  const transWeapon = i18n
    .t("game:weapon_names." + weapon)
    .replace(" ", "")
    .toLowerCase();
  if (exactMatch) {
    return weapon === normalizedQuery;
  } else {
    return `${weapon} ${transWeapon}`.indexOf(normalizedQuery) >= 0;
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
