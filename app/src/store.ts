import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { Action, configureStore, ThunkAction } from "@reduxjs/toolkit";
import { defaultRunStat, simSlice } from "/src/Pages/Sim/simSlice";
import { viewerSlice } from "./Pages/ViewerDashboard/viewerSlice";
import { userDataSlice } from "./Pages/Sim/userDataSlice";

export type RootState = ReturnType<typeof store.getState>;

const simStateKey = "redux-sim-v0.0.2";
const userDataKey = "redux-user-data-v0.0.1";

let persistedState = {};

if (localStorage.getItem(simStateKey)) {
  let s = JSON.parse(localStorage.getItem(simStateKey)!);
  //reset some defaults
  s.edit_index = -1;
  s.ready = 0;
  s.run = defaultRunStat;
  if (!s.adv_cfg_err) {
    s.adv_cfg_err = "";
  }
  if (!s.GOChars) {
    s.GOChars = [];
  }
  persistedState = Object.assign(persistedState, { sim: s });
  // localStorage.clear();
  console.log("loaded sim store from localStorage: ", persistedState);
}

if (localStorage.getItem(userDataKey)) {
  let s = JSON.parse(localStorage.getItem(userDataKey)!);
  persistedState = Object.assign(persistedState, { user_data: s });
}

const store = configureStore({
  reducer: {
    [simSlice.name]: simSlice.reducer,
    [viewerSlice.name]: viewerSlice.reducer,
    [userDataSlice.name]: userDataSlice.reducer,
  },
  preloadedState: persistedState,
});

store.subscribe(() => {
  localStorage.setItem(simStateKey, JSON.stringify(store.getState().sim));
  localStorage.setItem(userDataKey, JSON.stringify(store.getState().user_data));
});

export { store };

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
