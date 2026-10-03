import React from "react";
import { loadDraft, type SendOptions, saveDraft, sendDraft } from "./draft";
import {
	loadHero,
	loadSampleOnLoad,
	loadTheme,
	saveHero,
	saveSampleOnLoad,
	saveTheme,
} from "./prefs";

interface DraftValue {
	cfg: string;
	setCfg: (cfg: string) => void;
	send: (cfg: string, opts: SendOptions) => void;
}

interface PrefsValue {
	sampleOnLoad: boolean;
	setSampleOnLoad: (sampleOnLoad: boolean) => void;
	theme: string;
	setTheme: (theme: string) => void;
	hero: string;
	setHero: (hero: string) => void;
}

const DraftContext = React.createContext<DraftValue | null>(null);
const PrefsContext = React.createContext<PrefsValue | null>(null);

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
	const [theme, setTheme] = React.useState(() => loadTheme(localStorage));
	const [hero, setHero] = React.useState(() => loadHero(localStorage));

	React.useEffect(
		() => saveSampleOnLoad(localStorage, sampleOnLoad),
		[sampleOnLoad],
	);
	React.useEffect(() => {
		saveTheme(localStorage, theme);
		document.documentElement.dataset.theme = theme;
	}, [theme]);
	React.useEffect(() => saveHero(localStorage, hero), [hero]);

	const value = React.useMemo(
		() => ({ sampleOnLoad, setSampleOnLoad, theme, setTheme, hero, setHero }),
		[sampleOnLoad, theme, hero],
	);
	return (
		<PrefsContext.Provider value={value}>{children}</PrefsContext.Provider>
	);
}

export function AppStateProvider({ children }: { children: React.ReactNode }) {
	return (
		<PrefsProvider>
			<DraftProvider>{children}</DraftProvider>
		</PrefsProvider>
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
