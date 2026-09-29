import type { UserInfo, UserSettings } from "@gcsim/types";
import React from "react";
import { loadDraft, type SendOptions, saveDraft, sendDraft } from "./draft";
import { loadSampleOnLoad, saveSampleOnLoad } from "./prefs";
import { defaultUser, loadUser, mergeUser, saveUserSettings } from "./user";

interface DraftValue {
	cfg: string;
	setCfg: (cfg: string) => void;
	send: (cfg: string, opts: SendOptions) => void;
}

interface PrefsValue {
	sampleOnLoad: boolean;
	setSampleOnLoad: (sampleOnLoad: boolean) => void;
	settingsOpen: boolean;
	setSettingsOpen: (open: boolean) => void;
}

interface UserValue {
	user: UserInfo;
	merge: (user: UserInfo) => void;
	setSettings: (settings: UserSettings) => void;
	reset: () => void;
}

const DraftContext = React.createContext<DraftValue | null>(null);
const PrefsContext = React.createContext<PrefsValue | null>(null);
const UserContext = React.createContext<UserValue | null>(null);

function DraftProvider({ children }: { children: React.ReactNode }) {
	const [cfg, setCfg] = React.useState(() => loadDraft(localStorage));

	React.useEffect(() => saveDraft(localStorage, cfg), [cfg]);

	const send = React.useCallback(
		(incoming: string, opts: SendOptions) =>
			setCfg((current) => sendDraft(current, incoming, opts)),
		[],
	);
	const value = React.useMemo(() => ({ cfg, setCfg, send }), [cfg, send]);
	return (
		<DraftContext.Provider value={value}>{children}</DraftContext.Provider>
	);
}

function PrefsProvider({ children }: { children: React.ReactNode }) {
	const [sampleOnLoad, setSampleOnLoad] = React.useState(() =>
		loadSampleOnLoad(localStorage),
	);
	const [settingsOpen, setSettingsOpen] = React.useState(false);

	React.useEffect(
		() => saveSampleOnLoad(localStorage, sampleOnLoad),
		[sampleOnLoad],
	);

	const value = React.useMemo(
		() => ({ sampleOnLoad, setSampleOnLoad, settingsOpen, setSettingsOpen }),
		[sampleOnLoad, settingsOpen],
	);
	return (
		<PrefsContext.Provider value={value}>{children}</PrefsContext.Provider>
	);
}

function UserProvider({ children }: { children: React.ReactNode }) {
	const [user, setUser] = React.useState(() => loadUser(localStorage));

	React.useEffect(() => saveUserSettings(localStorage, user), [user]);

	const actions = React.useMemo(
		() => ({
			merge: (incoming: UserInfo) =>
				setUser((prev) => mergeUser(prev, incoming)),
			setSettings: (settings: UserSettings) =>
				setUser((prev) => ({ ...prev, data: { ...prev.data, settings } })),
			reset: () => setUser(defaultUser),
		}),
		[],
	);
	const value = React.useMemo(() => ({ user, ...actions }), [user, actions]);
	return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
}

export function AppStateProvider({ children }: { children: React.ReactNode }) {
	return (
		<UserProvider>
			<PrefsProvider>
				<DraftProvider>{children}</DraftProvider>
			</PrefsProvider>
		</UserProvider>
	);
}

function useRequired<T>(context: React.Context<T | null>, name: string): T {
	const value = React.useContext(context);
	if (!value) {
		throw new Error(`${name} must be used within an AppStateProvider`);
	}
	return value;
}

export const useDraft = () => useRequired(DraftContext, "useDraft");
export const usePrefs = () => useRequired(PrefsContext, "usePrefs");
export const useUser = () => useRequired(UserContext, "useUser");
