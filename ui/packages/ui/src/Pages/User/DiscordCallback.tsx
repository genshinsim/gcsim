import React from "react";
import { useNavigate } from "react-router";
import { useAppDispatch } from "../../Stores/store";
import { userActions } from "../../Stores/userSlice";
import { authProvider } from "./Login";

export function DiscordCallback() {
	const [error, setError] = React.useState("");
	const dispatch = useAppDispatch();
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
			.then((user) => {
				dispatch(userActions.mergeUser(user));
				navigate("/account");
			})
			.catch((error) => {
				setError(JSON.stringify(error));
				navigate("/account");
			});
	}, [dispatch, navigate]);

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
