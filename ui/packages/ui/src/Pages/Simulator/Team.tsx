import { Button, Card, Dialog, useHotkeys } from "@blueprintjs/core";
import React from "react";
import { CharacterEdit } from "./CharacterEdit";
import { Trans, useTranslation } from "react-i18next";
import { Builder } from "./Components/TeamBuilder/Builder";
import { OmniSelect, Item, GenerateDefaultCharacters } from "../../Components/Select";
import { CharMap } from "../../Data";
import { RootState, useAppDispatch, useAppSelector } from "../../Stores/store";
import { appActions } from "../../Stores/appSlice";
import { Character } from "@gcsim/types";

function newCharFromKey(k: string): Character {
  return {
    name: k,
    level: 80,
    max_level: 90,
    element: CharMap[k].element,
    cons: 0,
    weapon: {
      name: "dullblade",
      refine: 1,
      level: 1,
      max_level: 20,
    },
    talents: {
      attack: 6,
      skill: 6,
      burst: 6,
    },
    stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    snapshot: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    sets: {},
  };
}

export function Team() {
  const { t } = useTranslation();

  const { imported, team } = useAppSelector((state: RootState) => {
    return {
      imported: state.user_data.GOODImport,
      team: state.app.team,
    };
  });
  const [open, setOpen] = React.useState<boolean>(false);
  const [editIndex, setEditIndex] = React.useState<number>(-1);
  const dispatch = useAppDispatch();

  const hotkeys = React.useMemo(
    () => [
      {
        combo: "Esc",
        global: true,
        label: t("simple.exit_edit"),
        onKeyDown: () => {
          setEditIndex(-1);
        },
      },
    ],
    []
  );
  useHotkeys(hotkeys);

  const handleEdit = (index: number) => {
    return () => {
      if (index > -1 && index < team.length) {
        setEditIndex(index);
      }
    };
  };
  const handleRemove = (index: number) => {
    return () => {
      dispatch(appActions.deleteCharacter({ index }));
    };
  };

  const handleAdd = (item: Item) => {
    setOpen(false);
    //check if this is from GOOD
    if (item.notes) {
      const character: Character = JSON.parse(JSON.stringify(imported[item.key]));
      dispatch(appActions.addCharacter({ character }));
      return;
    }
    //else it's new
    const character = newCharFromKey(item.key);
    dispatch(appActions.addCharacter({ character }));
  };

  const handleChange = (char: Character, index: number) => {
    dispatch(appActions.editCharacter({ char, index }));
  };

  const disabled: string[] = team.map((c) => c.name);

  const items: Item[] = GenerateDefaultCharacters();

  Object.keys(imported).forEach((k) => {
    items.push({
      key: imported[k].name,
      text: t("game:character_names." + imported[k].name),
      label: t(`elements.${imported[k].element}`),
      notes: `Imported on ${imported[k].date_added}`,
    });
  });

  return (
    <div className="flex flex-col">
      <Builder
        team={team}
        handleAdd={() => setOpen(true)}
        handleEdit={handleEdit}
        handleRemove={handleRemove}
      />

      <Dialog
        isOpen={editIndex > -1}
        onClose={() => {
          setEditIndex(-1);
        }}
        style={{ width: "95%" }}
      >
        <Card className="m-2">
          <CharacterEdit
            index={editIndex}
            handleChange={handleChange}
            char={editIndex > -1 ? team[editIndex] : null}
          />
          <Button
            fill
            intent="primary"
            icon="edit"
            onClick={() => {
              setEditIndex(-1);
            }}
          >
            <Trans>simple.done</Trans>
          </Button>
        </Card>
      </Dialog>

      <OmniSelect
        isOpen={open}
        items={items}
        onClose={() => setOpen(false)}
        onSelect={handleAdd}
        disabled={disabled}
      />
    </div>
  );
}
