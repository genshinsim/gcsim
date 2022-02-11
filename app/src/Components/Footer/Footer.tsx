//@ts-ignore
import DiscordLogo from "./discord-icon.svg";
//@ts-ignore
import GithubLogo from "./github-icon.svg";
//@ts-ignore
import KofiLogo from "./ko-fi-icon.svg";

export default function Footer() {
  return (
    <div className="flex flex-row flex-wrap md:flex-row-reverse w-full justify-end gap-2 items-end">
      <div className="flex flex-row justify-end ml-auto basis-full md:basis-auto">
        <div className=" hover:bg-gray-600 p-2 rounded-md h-12">
          <a href="https://ko-fi.com/srliao" target="_blank">
            <img
              src={KofiLogo}
              alt="Ko-Fi Logo"
              className="object-contain h-full"
            />
          </a>
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
      <div className="mr-auto text-xs basis-full  md:basis-auto p-1">
        gcsim is not affiliated with miHoYo. Genshin Impact, game content and
        materials are trademarks and copyrights of miHoYo.
      </div>
    </div>
  );
}
