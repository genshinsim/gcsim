import { configureStore, ThunkAction, Action } from "@reduxjs/toolkit";
import { dbSlice } from "PageDatabase/dbSlice";

export const store = configureStore({
  reducer: {
    [dbSlice.name]: dbSlice.reducer,
  },
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
