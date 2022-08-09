import { Button } from "@blueprintjs/core";
import { Viewport } from "~src/Components";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { logout } from "~src/UserData/userSlice";
import { Login } from "./Login";

export default function UserAccount() {
  const { user } = useAppSelector((state: RootState) => {
    return {
      user: state.user,
    };
  });
  const dispatch = useAppDispatch();

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
