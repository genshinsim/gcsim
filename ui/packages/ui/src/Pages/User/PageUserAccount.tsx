import { Button, ButtonGroup, Checkbox } from "@blueprintjs/core";
import { Viewport } from "../../Components";
import { AppThunk, useAppDispatch, useAppSelector } from "../../Stores/store";
import {
  initialState,
  saveUserSettings,
  userActions,
} from "../../Stores/userSlice";
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
  const user = useAppSelector((state) => state.user);
  const dispatch = useAppDispatch();

  if (user.uid === "") {
    return <Login />;
  }

  return (
    <Viewport>
      <div className="flex flex-col ">
        <div>
          <Checkbox
            checked={user.data.settings.showTips}
            onChange={() => {
              dispatch(
                userActions.setUserSettings({
                  showTips: !user.data.settings.showTips,
                  showBuilder: user.data.settings.showBuilder,
                })
              );
            }}
          >
            Show tips
          </Checkbox>
          <Checkbox
            checked={user.data.settings.showBuilder}
            onChange={() => {
              dispatch(
                userActions.setUserSettings({
                  showTips: user.data.settings.showTips,
                  showBuilder: !user.data.settings.showBuilder,
                })
              );
            }}
          >
            Show builder
          </Checkbox>
        </div>
        <div className="flex flex-row place-content-center mt-2">
          <ButtonGroup>
            <Button
              icon="saved"
              large
              onClick={() => dispatch(saveUserSettings())}
            >
              Save Settings
            </Button>
            <Button
              icon="log-out"
              intent="danger"
              large
              onClick={() => dispatch(logout())}
            >
              Logout
            </Button>
          </ButtonGroup>
        </div>
      </div>
    </Viewport>
  );
}
