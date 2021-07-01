import { configureStore, ThunkAction, Action } from "@reduxjs/toolkit";
import appReducer from "app/appSlice";
import simReducer from "features/sim/simSlice";
import debugReducer from "features/debug/debugSlice";
import importReducer from "features/import/importSlice";
import resultsReducer from "features/results/resultsSlice";

export const store = configureStore({
  reducer: {
    app: appReducer,
    sim: simReducer,
    debug: debugReducer,
    import: importReducer,
    results: resultsReducer,
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
