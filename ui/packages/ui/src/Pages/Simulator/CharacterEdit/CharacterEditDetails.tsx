import { Button } from "@blueprintjs/core";
import { ascLvlMax, ascLvlMin, ascToMaxLvl, maxLvlToAsc } from "../Components/util";
import { NumberInput } from "../../../Components/NumberInput";
import React from "react";
import { Trans, useTranslation } from "react-i18next";
import { GenerateDefaultCharacters, Item, OmniSelect } from "../../../Components/Select";
import { RootState, useAppSelector } from "../../../Stores/store";
import { CharMap } from "../../../Data";
import { Character } from "@gcsim/types";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};

export function CharacterEditDetails({ char, onChange }: Props) {
  const { t } = useTranslation();

  const { team } = useAppSelector((state: RootState) => {
    return {
      team: state.app.team,
    };
  });
  const [open, setOpen] = React.useState<boolean>(false);

  const handleChangeCharacter = (w: Item) => {
    setOpen(false);
    //do nothing if this char already exists
    for (let i = 0; i < team.length; i++) {
      if (team[i].name === w.key) {
        return;
      }
    }
    const next = JSON.parse(JSON.stringify(char));
    next.name = w.key;
    next.element = CharMap[next.name].element;
    onChange(next);
  };

  const handleChangeTalent = (key: "attack" | "skill" | "burst") => {
    return (val: number) => {
      const next = JSON.parse(JSON.stringify(char));
      next.talents[key] = val;
      onChange(next);
    };
  };

  const handleChangeAsc = (val: number) => {
    if (val < 0 || val > 6) {
      return;
    }
    const m = ascToMaxLvl(val);
    let l = char.level;
    const asc = maxLvlToAsc(m);
    if (l > m) {
      l = m;
    } else if (l < ascLvlMin(asc)) {
      l = ascLvlMin(asc);
    }
    const next = JSON.parse(JSON.stringify(char));
    next.max_level = m;
    next.level = l;
    onChange(next);
  };

  const handleChangeLvl = (val: number) => {
    if (val <= 0 || val > 90) {
      return;
    }
    const next = JSON.parse(JSON.stringify(char));
    next.level = val;
    onChange(next);
  };

  const handleChangeCons = (val: number) => {
    if (val < 0 || val > 6) {
      return;
    }
    const next = JSON.parse(JSON.stringify(char));
    next.cons = val;
    onChange(next);
  };

  const asc = maxLvlToAsc(char.max_level);

  const disabled = team.map((c) => c.name);

  const items: Item[] = GenerateDefaultCharacters();

  return (
    <div className="flex flex-row gap-2 flex-wrap">
      <div className="flex flex-col place-items-center gap-1 basis-full hd:basis-36">
        <img
          src={"/api/assets/avatar/" + char.name + ".png"}
          alt={char.name}
          className="w-28"
        />
        <Button icon="swap-horizontal" fill onClick={() => setOpen(true)}>
          <Trans>characteredit.change</Trans>
        </Button>
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 hd:basis-0 flex flex-col gap-y-2">
        <NumberInput
          label={t("characteredit.ascension")}
          onChange={handleChangeAsc}
          value={asc}
          min={0}
          max={6}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.level")}
          onChange={handleChangeLvl}
          value={char.level}
          min={ascLvlMin(asc)}
          max={ascLvlMax(asc)}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.cons")}
          onChange={handleChangeCons}
          value={char.cons}
          integerOnly
        />
      </div>
      <div className="bg-gray-600 rounded-md basis-full flex-grow p-2 hd:basis-0 flex flex-col gap-y-2">
        <NumberInput
          label={t("characteredit.attack")}
          onChange={handleChangeTalent("attack")}
          min={1}
          max={10}
          value={char.talents.attack}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.skill")}
          onChange={handleChangeTalent("skill")}
          min={1}
          max={10}
          value={char.talents.skill}
          integerOnly
        />
        <NumberInput
          label={t("characteredit.burst")}
          onChange={handleChangeTalent("burst")}
          min={1}
          max={10}
          value={char.talents.burst}
          integerOnly
        />
      </div>
      <OmniSelect
        isOpen={open}
        items={items}
        onClose={() => setOpen(false)}
        onSelect={handleChangeCharacter}
        disabled={disabled}
      />
    </div>
  );
}
