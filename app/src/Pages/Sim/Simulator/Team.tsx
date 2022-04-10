import { Button, Card, ButtonGroup, Dialog } from "@blueprintjs/core";
import React from "react";
import {
  characterKeyToICharacter,
  CharacterSelect,
  ICharacter,
  newChar,
} from "~src/Components/Character";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { simActions } from "~src/Pages/Sim/simSlice";
import { CharacterEdit } from "./CharacterEdit";
import { Trans, useTranslation } from "react-i18next";
import { Builder } from "../Components/TeamBuilder/Builder";

export function Team() {
  useTranslation();

  const { team, edit_index, imported } = useAppSelector((state: RootState) => {
    return {
      team: state.sim.team,
      edit_index: state.sim.edit_index,
      imported: state.user_data.GOODImport,
    };
  });
  const dispatch = useAppDispatch();
  const [open, setOpen] = React.useState<boolean>(false);

  const handleEdit = (index: number) => {
    return () => {
      if (index > -1 && index < team.length) {
        dispatch(simActions.editCharacter({ index: index }));
      }
    };
  };
  const handleRemove = (index: number) => {
    return () => {
      if (index > -1 && index < team.length) {
        dispatch(simActions.deleteCharacter({ index: index }));
      }
    };
  };

  const handleAdd = (character: ICharacter) => {
    setOpen(false);
    //check if this is from GOOD
    if (character.notes) {
      dispatch(simActions.addCharacter({ character: imported[character.key] }));
      return;
    }
    //else it's new
    const c = newChar(character.key);
    dispatch(simActions.addCharacter({ character: c }));
  };

  let disabled: string[] = team.map((c) => c.name);

  const additionalChars = Object.keys(imported).map((k) => {
    let x = Object.assign({}, characterKeyToICharacter[k]);
    x.notes = `Imported on ${imported[k].date_added}`;
    return x;
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
        isOpen={edit_index > -1}
        onClose={() => {
          dispatch(simActions.editCharacter({ index: -1 }));
        }}
        style={{ width: "95%" }}
      >
        <Card className="m-2">
          <CharacterEdit />
          <Button
            fill
            intent="primary"
            icon="edit"
            onClick={() => {
              dispatch(simActions.editCharacter({ index: -1 }));
            }}
          >
            <Trans>simple.done</Trans>
          </Button>
        </Card>
      </Dialog>

      <CharacterSelect
        disabled={disabled}
        onClose={() => setOpen(false)}
        onSelect={handleAdd}
        isOpen={open}
        additionalOptions={additionalChars}
      />
    </div>
  );
}
