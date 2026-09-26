import type { Character } from "@gcsim/types";
import React from "react";

export type ImportSource = "enka" | "good";

type ImportedMap = { [key: string]: Character };

interface ImportedCharactersValue {
	imported: ImportedMap;
	loadImported: (data: Character[], source: ImportSource) => void;
}

const ImportedCharactersContext =
	React.createContext<ImportedCharactersValue | null>(null);

const legacyReduxImportsKey = "redux-user-data-v0.0.1";

function hydrate(): ImportedMap {
	try {
		const raw = localStorage.getItem(legacyReduxImportsKey);
		return raw ? (JSON.parse(raw).GOODImport ?? {}) : {};
	} catch {
		return {};
	}
}

export function ImportedCharactersProvider({
	children,
}: {
	children: React.ReactNode;
}) {
	const [imported, setImported] = React.useState<ImportedMap>(hydrate);

	React.useEffect(() => {
		localStorage.setItem(
			legacyReduxImportsKey,
			JSON.stringify({ GOODImport: imported }),
		);
	}, [imported]);

	const loadImported = React.useCallback(
		(data: Character[], source: ImportSource) => {
			if (data.length === 0) {
				return;
			}
			const next: ImportedMap = {};
			data.forEach((c) => {
				const key = `${c.name}-${source}-${c.enka_build_name ?? "none"}`;
				next[key] = { ...c, source };
			});
			setImported(next);
		},
		[],
	);

	const value = React.useMemo(
		() => ({ imported, loadImported }),
		[imported, loadImported],
	);

	return (
		<ImportedCharactersContext.Provider value={value}>
			{children}
		</ImportedCharactersContext.Provider>
	);
}

export function useImportedCharacters(): ImportedCharactersValue {
	const ctx = React.useContext(ImportedCharactersContext);
	if (!ctx) {
		throw new Error(
			"useImportedCharacters must be used within an ImportedCharactersProvider",
		);
	}
	return ctx;
}
