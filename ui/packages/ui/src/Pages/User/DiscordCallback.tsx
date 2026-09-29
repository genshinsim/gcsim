import { useNavigate } from "@tanstack/react-router";
import React from "react";
import { useUser } from "../../stores/AppState";
import { authProvider } from "./Login";

export function DiscordCallback() {
	const [error, setError] = React.useState("");
	const { merge } = useUser();
	const navigate = useNavigate();

	React.useEffect(() => {
		const params = new URLSearchParams(window.location.search);
		const code = params.get("code");
		if (code === null) {
			setError("Invalid Discord auth code. Login failed.");
			return;
		}
		authProvider
			.auth(code)
			.then((info) => {
				merge(info);
				navigate({ to: "/account" });
			})
			.catch((error) => {
				setError(JSON.stringify(error));
				navigate({ to: "/account" });
			});
	}, [merge, navigate]);

	if (error === "") {
		return (
			<div className="flex flex-row place-content-center">
				Logging in... please wait.
			</div>
		);
	} else {
		return (
			<div className="flex flex-row place-content-center">
				Error encountered: {error}
			</div>
		);
	}
}
