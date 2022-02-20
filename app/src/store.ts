import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { Action, configureStore, ThunkAction } from "@reduxjs/toolkit";
import { defaultRunStat, simSlice } from "/src/Pages/Sim/simSlice";
import { viewerSlice } from "./Pages/ViewerDashboard/viewerSlice";

const storageKey = "redux-sim-v0.0.2";

let persistedState = {};
if (localStorage.getItem(storageKey)) {
  let s = JSON.parse(localStorage.getItem(storageKey)!);
  //reset some defaults
  s.edit_index = -1;
  s.ready = 0;
  s.run = defaultRunStat;
  persistedState = { sim: s };
  localStorage.clear();
  console.log("loaded sim store from localStorage: ", persistedState);
}

const store = configureStore({
  reducer: {
    [simSlice.name]: simSlice.reducer,
    [viewerSlice.name]: viewerSlice.reducer,
  },
  preloadedState: persistedState,
});

store.subscribe(() => {
  localStorage.setItem(storageKey, JSON.stringify(store.getState().sim));
});

export { store };

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;

// Use throughout your app instead of plain `useDispatch` and `useSelector`
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
