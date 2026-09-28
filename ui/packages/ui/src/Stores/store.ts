import {
	type Action,
	combineReducers,
	configureStore,
	createListenerMiddleware,
	type ThunkAction,
	type TypedStartListening,
	type TypedStopListening,
} from "@reduxjs/toolkit";
import {
	type TypedUseSelectorHook,
	useDispatch,
	useSelector,
} from "react-redux";
import { appSlice } from "./appSlice";
import { userSlice } from "./userSlice";

const listenerMiddleware = createListenerMiddleware();

const rootReducer = combineReducers({
	[userSlice.name]: userSlice.reducer,
	[appSlice.name]: appSlice.reducer,
});

export type RootState = ReturnType<typeof rootReducer>;

const userLocalSettings = "redux-user-local-settings";
const userAppDataKey = "redux-app-data";

const persistedState: RootState = JSON.parse(
	JSON.stringify({
		[userSlice.name]: userSlice.getInitialState(),
		[appSlice.name]: appSlice.getInitialState(),
	}),
);

if (localStorage.getItem(userAppDataKey)) {
	const s = JSON.parse(localStorage.getItem(userAppDataKey) ?? "{}");
	persistedState.app = Object.assign(persistedState.app, {
		sampleOnLoad: s.sampleOnLoad ?? false,
		cfg: s.cfg ?? "",
	});
}

if (localStorage.getItem(userLocalSettings)) {
	const s = JSON.parse(localStorage.getItem(userLocalSettings) ?? "{}");
	persistedState.user.data = Object.assign(persistedState.user.data, {
		settings: s,
	});
}

export const store = configureStore({
	reducer: rootReducer,
	preloadedState: persistedState,
	middleware: (getDefaultMiddleware) =>
		getDefaultMiddleware({
			serializableCheck: false,
		}).prepend(listenerMiddleware.middleware),
});

store.subscribe(() => {
	localStorage.setItem(
		userAppDataKey,
		JSON.stringify({
			sampleOnLoad: store.getState().app.sampleOnLoad,
			cfg: store.getState().app.cfg,
		}),
	);

	if (store.getState().user.data.settings) {
		localStorage.setItem(
			userLocalSettings,
			JSON.stringify(store.getState().user.data.settings),
		);
	}
});

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

export const appStartListening =
	listenerMiddleware.startListening as TypedStartListening<RootState>;
export const appStopListening =
	listenerMiddleware.stopListening as TypedStopListening<RootState>;
