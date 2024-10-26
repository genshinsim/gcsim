import {Character} from '@gcsim/types';
import React from 'react';
import {useTranslation} from 'react-i18next';
import {
  GenerateDefaultCharacters,
  Item,
  OmniSelect,
} from '../../Components/Select';
import {CharMap} from '../../Data';
import {appActions} from '../../Stores/appSlice';
import {RootState, useAppDispatch, useAppSelector} from '../../Stores/store';
import {Builder} from './Components/TeamBuilder/Builder';

function newCharFromKey(k: string): Character {
  console.log(k);
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
  const {t} = useTranslation();

  const {imported, team} = useAppSelector((state: RootState) => {
    return {
      imported: state.user_data.GOODImport,
      team: state.app.team,
    };
  });
  const [open, setOpen] = React.useState<boolean>(false);
  const dispatch = useAppDispatch();

  const handleRemove = (index: number) => {
    return () => {
      dispatch(appActions.deleteCharacter({index}));
    };
  };

  const handleAdd = (item: Item) => {
    setOpen(false);
    if (item.char_source === 'user') {
      const character: Character = JSON.parse(
        JSON.stringify(imported[item.key]),
      );
      dispatch(appActions.addCharacter({character}));
    } else {
      const character = newCharFromKey(item.key);
      dispatch(appActions.addCharacter({character}));
    }
  };

  // filter based on char_key,which is c.name from Character
  const disabled: string[] = team.map((c) => c.name);

  const items: Item[] = GenerateDefaultCharacters();

  Object.keys(imported).forEach((k) => {
    const e = imported[k];
    let label = e.enka_build_name !== undefined ? ` ${e.enka_build_name}` : '';
    label += e.date_added !== undefined ? ` (Imported on ${e.date_added})` : '';
    items.push({
      key: k,
      char_key: e.name,
      char_source: 'user',
      text: t('game:character_names.' + e.name),
      label: label,
    });
  });

  return (
    <div className="flex flex-col">
      <Builder
        team={team}
        handleAdd={() => setOpen(true)}
        handleRemove={handleRemove}
      />

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
