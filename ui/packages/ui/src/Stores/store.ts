import { Action, configureStore, createListenerMiddleware, ThunkAction, TypedStartListening, TypedStopListening } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { appSlice } from "./appSlice";
import { userDataSlice } from "./userDataSlice";
import { userSlice } from "./userSlice";
import { viewerSlice } from "./viewerSlice";

const listenerMiddleware = createListenerMiddleware();

const userDataKey = "redux-user-data-v0.0.1";
const userLocalSettings = "redux-user-local-settings";
const userAppDataKey = "redux-app-data";

const persistedState = JSON.parse(
  JSON.stringify({
    [userDataSlice.name]: userDataSlice.getInitialState(),
    [userSlice.name]: userSlice.getInitialState(),
    [viewerSlice.name]: viewerSlice.getInitialState(),
    [appSlice.name]: appSlice.getInitialState(),
  })
);

if (localStorage.getItem(userDataKey)) {
  const s = JSON.parse(localStorage.getItem(userDataKey) ?? "{}");
  persistedState.user_data = Object.assign(persistedState.user_data, s);
}

if (localStorage.getItem(userAppDataKey)) {
  const s = JSON.parse(localStorage.getItem(userAppDataKey) ?? "{}");
  persistedState.app = Object.assign(persistedState.app, {
    cfg: s.cfg ?? "",
    team: s.team ?? []
  });
}

if (localStorage.getItem(userLocalSettings)) {
  const s = JSON.parse(localStorage.getItem(userLocalSettings) ?? "{}");
  persistedState.user = Object.assign(persistedState.user, { settings: s });
}

export const store = configureStore({
  reducer: {
    [userDataSlice.name]: userDataSlice.reducer,
    [userSlice.name]: userSlice.reducer,
    [viewerSlice.name]: viewerSlice.reducer,
    [appSlice.name]: appSlice.reducer,
  },
  preloadedState: persistedState,
  middleware: (getDefaultMiddleware) => getDefaultMiddleware({
    serializableCheck: false,
  }).prepend(listenerMiddleware.middleware),
});

store.subscribe(() => {
  localStorage.setItem(userDataKey, JSON.stringify(store.getState().user_data));
  localStorage.setItem(userAppDataKey, JSON.stringify({
    cfg: store.getState().app.cfg,
    team: store.getState().app.team
  }));
  if (store.getState().user.settings) {
    localStorage.setItem(
      userLocalSettings,
      JSON.stringify(store.getState().user.settings)
    );
  }
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;

export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;

// Use throughout your app instead of plain `useDispatch` and `useSelector`
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

export const appStartListening =
    listenerMiddleware.startListening as TypedStartListening<RootState>;
export const appStopListening =
    listenerMiddleware.stopListening as TypedStopListening<RootState>;