import { Button } from "@blueprintjs/core";
import axios from "axios";
import React from "react";
import { Viewport } from "../../Components";
import { useAppDispatch, useAppSelector } from "../../Stores/store";
import { statusType } from "../../Types";
import { logout } from "../../Stores/userSlice";
import { Login } from "./Login";

export function PageUserAccount() {
  const [status, setStatus] = React.useState<statusType>("idle");
  const [errMsg, setErrMsg] = React.useState<string>("");

  const user = useAppSelector((state) => state.user);
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    if (status === "idle" && user.token && user.token !== "") {
      axios
        .get(`/api/${user.user_id}/sims`)
        .then((resp) => {
          console.log(resp.data);
          setStatus("done");
        })
        .catch((err) => {
          setStatus("error");
          setErrMsg(`Error encountered loading sims for user: ${err}`);
        });
    }
  }, [status, dispatch, user.token]);

  if (user.token === "" || user.token === undefined) {
    return <Login />;
  }

  return (
    <Viewport>
      <div className="flex flex-row place-content-center mt-2">
        <Button icon="log-out" large onClick={() => dispatch(logout())}>
          Logout
        </Button>
      </div>
    </Viewport>
  );
}
