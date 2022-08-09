import { Button } from "@blueprintjs/core";
import { Viewport } from "~src/Components";
//@ts-ignore
import DiscordLogo from "~src/Assets/discord-icon.svg";
import { authProvider } from "~src/UserData/userSlice";

export function Login() {
  return (
    <Viewport>
      <div className="flex flex-row place-content-center mt-2">
        <Button
          large
          icon={
            <img
              src={DiscordLogo}
              alt="Discord Logo"
              className="object-contain h-[16px]"
            />
          }
          onClick={() => authProvider.login()}
        >
          Login with Discord
        </Button>
      </div>
    </Viewport>
  );
}
