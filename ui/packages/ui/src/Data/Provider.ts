import { UserInfo } from "@gcsim/types";
import axios from "axios";
import { initialState } from "../Stores/userSlice";

export interface AuthProvider {
  login(): void;
  auth(code: string): Promise<UserInfo>;
  logout(): Promise<void>;
}

export class MockProvider implements AuthProvider {
  constructor() {}

  login(): void {
    //fake navigating back to the site from discord
    window.location.href = "/auth/discord?code=fake";
  }

  async auth(_: string): Promise<UserInfo> {
    console.log("hello im logging in");
    return {
      uid: "1560962267213",
      name: "FakeUser#1234",
      role: 0,
      permalinks: [],
      data: {
        settings: initialState.data.settings,
      },
    };
  }

  async logout(): Promise<void> {}

  getAccountData(): void {}

  setAccountData(): void {}
}

const callbackURL =
  window.location.protocol + "//" + window.location.host + "/auth/discord";
const discordURL =
  "https://discord.com/api/oauth2/authorize?client_id=1040701711783829566&redirect_uri=" +
  encodeURIComponent(callbackURL) +
  "&response_type=code&scope=identify&prompt=none";

export class DiscordProvider implements AuthProvider {
  private started: boolean;
  constructor() {
    this.started = false;
  }

  login(): void {
    window.location.href = discordURL;
  }

  async auth(code: string): Promise<UserInfo> {
    if (this.started) {
      throw "Auth already started";
    }
    this.started = true;
    const response = await axios({
      method: "get",
      url: "/api/login",
      headers: {
        "X-DISCORD-CODE": code,
        "x-discord-redirect": callbackURL, //TODO: check prod vs dev
      },
    });
    this.started = false;
    console.log(response);
    if (response.data) {
      return response.data;
    }
    throw "Unexpected error";
  }

  async logout(): Promise<void> {}
}
