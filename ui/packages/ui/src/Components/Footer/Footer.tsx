import React from "react";
//@ts-ignore
import DiscordLogo from "./discord-icon.svg";
//@ts-ignore
import GithubLogo from "./github-icon.svg";
//@ts-ignore
import KofiLogo from "./ko-fi-icon.svg";
import { Trans, useTranslation } from "react-i18next";

export function Footer() {
  useTranslation();

  return (
    <div className="flex flex-row flex-wrap md:flex-row-reverse w-full justify-end gap-2 items-end">
      <div className="flex flex-row justify-end ml-auto basis-full md:basis-auto">
        <div className=" hover:bg-gray-600 p-2 rounded-md h-12">
          <a href="https://ko-fi.com/srliao" target="_blank" rel="noreferrer">
            <img
              src={KofiLogo}
              alt="Ko-Fi Logo"
              className="object-contain h-full"
            />
          </a>
        </div>
        <div className=" hover:bg-gray-600 p-2 rounded-md h-12">
          <a
            href="https://github.com/genshinsim/gsim"
            target="_blank"
            rel="noreferrer"
          >
            <img
              src={GithubLogo}
              alt="Github Logo"
              className="object-contain h-full"
            />
          </a>
        </div>
        <div className=" hover:bg-gray-600 p-2 rounded-md h-12 ">
          <a
            href="https://discord.gg/m7jvjdxx7q"
            target="_blank"
            rel="noreferrer"
          >
            <img
              src={DiscordLogo}
              alt="Discord Logo"
              className="object-contain h-full"
            />
          </a>
        </div>
      </div>
      <div className="mr-auto text-xs basis-full md:basis-auto p-1 self-center">
        <Trans>footer.gcsim_is_not</Trans>
      </div>
    </div>
  );
}
