import React from "react";
import { useLocation } from "wouter";
import { useAppDispatch } from "../../Stores/store";
import { userActions } from "../../Stores/userSlice";
import { authProvider } from "./Login";

export function DiscordCallback() {
  const params = new URLSearchParams(window.location.search);
  const [error, setError] = React.useState("");
  const dispatch = useAppDispatch();
  const [_, setLocation] = useLocation();

  React.useEffect(() => {
    const code = params.get("code");
    if (code === null) {
      setError("Invalid Discord auth code. Login failed.");
      return;
    }
    authProvider
      .auth(code)
      .then((user) => {
        dispatch(userActions.setUser(user));
        setLocation("/account");
      })
      .catch((error) => {
        setError(JSON.stringify(error));
        setLocation("/auth/discord");
      });
  }, [params]);

  if (error === "") {
    return <div className="flex flex-row place-content-center">Logging in... please wait.</div>;
  } else {
    return <div className="flex flex-row place-content-center">Error encountered: {error}</div>;
  }
}
