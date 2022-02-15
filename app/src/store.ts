import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { Action, configureStore, ThunkAction } from "@reduxjs/toolkit";
import { simSlice } from "/src/Pages/Sim/simSlice";
import {
  viewerInitialState,
  viewerSlice,
} from "./Pages/ViewerDashboard/viewerSlice";

let persistedState = {};
if (localStorage.getItem("gcsim_redux")) {
  let s = JSON.parse(localStorage.getItem("gcsim_redux")!);
  //reset some defaults
  s.edit_index = -1;
  s.ready = 0;
  s.run = {
    progress: -1,
    result: -1,
    time: -1,
  };
  persistedState = s;
}

const store = configureStore({
  reducer: {
    [simSlice.name]: simSlice.reducer,
    [viewerSlice.name]: viewerSlice.reducer,
  },
  preloadedState: persistedState,
});

store.subscribe(() => {
  let s = store.getState();
  s.viewer = viewerInitialState;
  localStorage.setItem("gcsim_redux", JSON.stringify(s));
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
