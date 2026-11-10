import { Label, Switch } from "@gcsim/primitives";
import React from "react";
import { useTranslation } from "react-i18next";
import ServerMode from "./ServerMode";
import WasmMode from "./WasmMode";

const serverModeKey = "use-server-mode";

const App = () => {
	const { t } = useTranslation();
	const serverModeParam = new URLSearchParams(window.location.search).get(
		serverModeKey,
	);
	const [serverMode, setServerMode] = React.useState<boolean>((): boolean => {
		if (serverModeParam !== null) return serverModeParam === "true";
		return localStorage.getItem(serverModeKey) === "true";
	});
	React.useEffect(() => {
		if (serverModeParam === null) {
			localStorage.setItem(serverModeKey, serverMode.toString());
		}
	}, [serverModeParam, serverMode]);

	const serverSwitch = (
		<div className="flex items-center gap-2">
			<Switch
				id="server-mode-switch"
				checked={serverMode}
				onCheckedChange={setServerMode}
			/>
			<Label htmlFor="server-mode-switch">
				{t(
					serverMode
						? "simple.server_mode_disable"
						: "simple.server_mode_enable",
				)}
			</Label>
		</div>
	);

	return (
		<>
			{serverMode ? (
				<ServerMode>{serverSwitch}</ServerMode>
			) : (
				<WasmMode>{serverSwitch}</WasmMode>
			)}
		</>
	);
};

export default App;
