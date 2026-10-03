import React, { type ReactNode } from "react";

export type ExecutorKind = "wasm" | "server";

type HostValue = {
	executorSettings: ReactNode;
	executorKind: ExecutorKind;
};

const HostContext = React.createContext<HostValue | null>(null);

export const HostProvider = HostContext.Provider;

export function useHost(): HostValue {
	const value = React.useContext(HostContext);
	if (!value) {
		throw new Error("useHost must be used within a HostProvider");
	}
	return value;
}
