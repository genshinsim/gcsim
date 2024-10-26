import {Character} from '@gcsim/types';
import {createSlice, PayloadAction} from '@reduxjs/toolkit';

export interface UserData {
  GOODImport: {[key: string]: Character};
}

const initialState: UserData = {
  GOODImport: {},
};

type ImportSource = 'enka' | 'good';

export const userDataSlice = createSlice({
  name: 'user_data',
  initialState: initialState,
  reducers: {
    loadFromGOOD: (
      state,
      action: PayloadAction<{data: Character[]; source: ImportSource}>,
    ) => {
      // if there are characters, do something
      if (action.payload.data.length > 0) {
        //make it
        state.GOODImport = {};
        action.payload.data.forEach((c) => {
          c.source = action.payload.source;
          //unique key should be name + source + optional enka build name
          const key = `${c.name}-${c.source}-${c.enka_build_name ?? 'none'}`;
          state.GOODImport[key] = c;
        });
      }
      return state;
    },
  },
});

export const userDataActions = userDataSlice.actions;

export type UserDataSlice = {
  [userDataSlice.name]: ReturnType<(typeof userDataSlice)['reducer']>;
};
