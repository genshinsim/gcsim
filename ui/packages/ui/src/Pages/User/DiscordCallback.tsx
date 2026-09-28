import React from "react";
import { useNavigate } from "react-router";
import { user } from "../../Stores";
import { authProvider } from "./Login";

export function DiscordCallback() {
	const [error, setError] = React.useState("");
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
				user.merge(info);
				navigate("/account");
			})
			.catch((error) => {
				setError(JSON.stringify(error));
				navigate("/account");
			});
	}, [navigate]);

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
