import { Button } from "@blueprintjs/core";
import { StatusType } from "@gcsim/types";
import axios from "axios";
import React from "react";
import { Viewport } from "../../Components";
import { AppThunk, useAppDispatch, useAppSelector } from "../../Stores/store";
import { initialState, userActions } from "../../Stores/userSlice";
import { authProvider, Login } from "./Login";

//thunks
function logout(): AppThunk {
  return function (dispatch) {
    authProvider
      .logout()
      .then(() => dispatch(userActions.setUser(initialState)))
      .catch((err) => {
        //log out the user
        console.warn("Error occured logging out: ", err);
        dispatch(userActions.setUser(initialState));
      });
  };
}

export function PageUserAccount() {
  const [status, setStatus] = React.useState<StatusType>("idle");
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
