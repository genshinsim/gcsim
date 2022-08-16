import { Spinner } from '@blueprintjs/core';
import axios from 'axios';
import React from 'react';
import { useAppDispatch, useAppSelector } from '~src/store';
import { DBAvatarSimDetails } from '~src/Types/database';
import { loadCharacter } from './dbSlice';

type CharacterViewProps = {
  char: string;
};

export function CharacterView({ char }: CharacterViewProps) {
  const { charSims } = useAppSelector((state) => {
    return {
      charSims: state.db.charSims,
    };
  });
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    if (!(char in charSims)) {
      dispatch(loadCharacter(char));
    }
  }, [charSims, dispatch]);

  return <div>hi</div>;
}
