import { Callout } from "@blueprintjs/core";
import { useTranslation } from "react-i18next";
//@ts-ignore
import DiscordLogo from "./Footer/discord-icon.svg";
import { SimResults } from "./Viewer/DataType";


export const AnnouncementBanner = ({}) => {
  const { t } = useTranslation();

  const link = (
    <a href="https://discord.gg/m7jvjdxx7q" target="_blank">
      <img
        src={DiscordLogo}
        alt="Discord Logo"
        className="inline object-contain h-[18px]"
      />
      {' '} discord
    </a>
  );

  const title = t("game:character_names.dehya") + " and " + t("game:character_names.mika") + " now available!"

  return (
    <Callout intent="primary" title={title} className="mt-4">
      <div>
        Visit the {link} for detailed release notes
      </div>
    </Callout>
  );
}

export const EarlyReleaseBanner = ({ data }: { data: SimResults }) => {
  const { t } = useTranslation();

  if (data.incomplete_chars == null || data.incomplete_chars.length == 0) {
    return null;
  }

  const link = (
    <a href="https://discord.gg/m7jvjdxx7q" target="_blank">
      <img
        src={DiscordLogo}
        alt="Discord Logo"
        className="inline object-contain h-[18px]"
      />
      {' '} gcsim discord!
    </a>
  );

  return (
    <div className="flex flex-col items-center">
      <Callout intent="warning" title="Early Release Characters" className="mb-4 max-w-2xl">
        <p>
          This simulation contains early release characters! These characters are fully implemented,
          but may not have optimal frame data aligned with in-game animations. We are actively
          collecting data to improve their implementation. If you wish to help,
          please reach out in the {link}
        </p>
        <div className="flex flex-col justify-start gap-1 text-xs font-mono text-gray-400">
          <span className="font-bold text-sm font-sans ">
            DB submissions for these characters will be disabled until frames are updated:
          </span>
          <ul className="list-disc pl-4 grid grid-cols-1 gap-x-3 justify-start">
            {data.incomplete_chars?.map((c) => (
              <div className="list-item">{t('game:character_names.' + c)}</div>
            ))}
          </ul>
        </div>
      </Callout>
    </div>
  );
}