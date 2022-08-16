import { H3 } from '@blueprintjs/core';
import React from 'react';
import { useLocation } from 'wouter';
import { useAppDispatch, useAppSelector } from '~src/store';
import { DBAvatarSimDetails } from '~src/Types/database';
import { CharacterCard } from './CharacterCard';
import { loadCharacter } from './dbSlice';
import { TeamsList } from './TeamList';

type CharacterViewProps = {
  char: string;
};

type Teams = {
  [key in string]: DBAvatarSimDetails[];
};

function reduceTeams(sims: DBAvatarSimDetails[]): Teams {
  let next: Teams = {};
  sims.forEach((s) => {
    const key = s.metadata.char_names
      .map((x) => x)
      .sort()
      .join('-');
    if (!(key in next)) {
      next[key] = [];
    }
    next[key].push(s);
  });

  return next;
}

export function CharacterView({ char }: CharacterViewProps) {
  const charSims = useAppSelector((state) => state.db.charSims);
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();

  React.useEffect(() => {
    if (!(char in charSims)) {
      dispatch(loadCharacter(char));
    }
  }, [charSims, dispatch]);

  if (!(char in charSims) || charSims[char].length === 0) {
    return (
      <div className="flex flex-row place-content-center mt-2">
        Sorry, this character does not have any sims submitted :(
      </div>
    );
  }

  const teams = charSims[char].reduce<{ [key in string]: number }>(
    (next, s) => {
      const key = s.metadata.char_names
        .map((x) => x)
        .sort()
        .join('-');
      next[key]++;
      return next;
    },
    {}
  );

  const rows = Object.keys(teams)
    .sort()
    .map((e) => {
      const col = e
        .split('-')
        .map((char) => <CharacterCard char={char} key={char} />);
      return (
        <div
          key={e}
          className="flex flex-row bg-gray-700 rounded-md p-2 hover:cursor-pointer"
          onClick={() => setLocation(`/db/${char}/${e}`)}
        >
          {col}
        </div>
      );
    });

  return (
    <main className="flex flex-col h-full m-2 w-full xs:w-full sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto ">
      <H3>
        <span style={{ textTransform: 'capitalize' }}>{char}</span> Teams
      </H3>
      <div className="grid grid-cols-2 gap-3">{rows}</div>
    </main>
  );
}
