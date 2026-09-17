import { Label, Switch } from "@gcsim/primitives";
import React from "react";
import { useTranslation } from "react-i18next";
import ServerMode from "./ServerMode";
import WasmMode from "./WasmMode";

const serverModeKey = "use-server-mode";

const App = () => {
	const { t } = useTranslation();
	const [serverMode, setServerMode] = React.useState<boolean>((): boolean => {
		return localStorage.getItem(serverModeKey) === "true";
	});
	React.useEffect(() => {
		localStorage.setItem(serverModeKey, serverMode.toString());
	}, [serverMode]);

	const children = (
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
				<ServerMode>{children}</ServerMode>
			) : (
				<WasmMode>{children}</WasmMode>
			)}
		</>
	);
};

export default App;
