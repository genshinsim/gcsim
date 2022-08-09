import axios from "axios";
import { UserInfo } from "~src/Types/user";

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
      user_id: 0,
      user_key: "123145",
      user_name: "FakeUser#1234",
      token: "thisisafaketoken",
    };
  }

  async logout(): Promise<void> {}

  getAccountData(): void {}

  setAccountData(): void {}
}

const discordURL =
  "https://discord.com/api/oauth2/authorize?client_id=999720972875739226&redirect_uri=https%3A%2F%2Fnext.gcsim.app%2Fauth%2Fdiscord&response_type=code&scope=identify&prompt=none";

export class DiscordProvider implements AuthProvider {
  constructor() {}

  login(): void {
    window.location.href = discordURL;
  }

  async auth(code: string): Promise<UserInfo> {
    const response = await axios({
      method: "get",
      url: "/api/auth",
      headers: { "X-DISCORD-CODE": code },
    });
    console.log(response);
    if (response.data.token && response.data.user) {
      //extract user info from token
      return {
        ...response.data.user,
        token: response.data.token,
      };
    }
    throw "Unexpected error";
  }

  async logout(): Promise<void> {}
}
