import { Button, Card, Dialog } from '@blueprintjs/core';
import React from 'react';
import { RootState, useAppDispatch, useAppSelector } from '~src/store';
import { simActions } from '~src/Pages/Sim/simSlice';
import { CharacterEdit } from './CharacterEdit';
import { Trans, useTranslation } from 'react-i18next';
import { Builder } from '../Components/TeamBuilder/Builder';
import {
  OmniSelect,
  Item,
  GenerateDefaultCharacters,
} from '~src/Components/Select';
import { Character } from '~src/Types/sim';
import { CharMap, TransformTravelerKeyToName, TravelerCheck } from '~src/Data';

function newCharFromKey(k: string): Character {
  return {
    name: k,
    level: 80,
    max_level: 90,
    element: CharMap[k].element,
    cons: 0,
    weapon: {
      name: 'dullblade',
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
    snapshot: [
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    sets: {},
  };
}

export function Team() {
  let { t } = useTranslation();

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

  const handleAdd = (item: Item) => {
    setOpen(false);
    //check if this is from GOOD
    if (item.notes) {
      dispatch(simActions.addCharacter({ character: imported[item.key] }));
      return;
    }
    //else it's new
    const c = newCharFromKey(item.key);
    dispatch(simActions.addCharacter({ character: c }));
  };

  let disabled: string[] = team.map((c) => c.name);

  let items: Item[] = GenerateDefaultCharacters();

  Object.keys(imported).forEach((k) => {
    items.push({
      key: imported[k].name,
      text: t('game:character_names.' + imported[k].name),
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
        isOpen={edit_index > -1}
        onClose={() => {
          dispatch(simActions.editCharacter({ index: -1 }));
        }}
        style={{ width: '95%' }}
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
