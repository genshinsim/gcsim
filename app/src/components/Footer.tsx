import React from "react";
import DiscordLogo from "../discord-icon.svg";
import GithubLogo from "../github-icon.svg";

export default function Footer() {
  return (
    <div className="flex flex-row w-full justify-end gap-2 lg:mb-10 md:m-3 sm:m-1 items-end">
      <div className="mr-auto text-xs">
        gcsim is not affiliated with miHoYo. Genshin Impact, game content and
        materials are trademarks and copyrights of miHoYo.
      </div>
      <div className=" hover:bg-gray-600 p-2 rounded-md h-12">
        <a href="https://github.com/genshinsim/gsim" target="_blank">
          <img
            src={GithubLogo}
            alt="Github Logo"
            className="object-contain h-full"
          />
        </a>
      </div>
      <div className=" hover:bg-gray-600 p-2 rounded-md h-12 ">
        <a href="https://discord.gg/m7jvjdxx7q" target="_blank">
          <img
            src={DiscordLogo}
            alt="Discord Logo"
            className="object-contain h-full"
          />
        </a>
      </div>
    </div>
  );
}
