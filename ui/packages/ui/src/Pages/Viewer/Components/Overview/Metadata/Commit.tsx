import { Tag } from "@blueprintjs/core";
import { memo } from "react";
import { useTranslation } from "react-i18next";

type Props = {
  commit?: string;
}

export const Commit = memo(({ commit }: Props) => {
  const { t } = useTranslation();
  if (commit == null || commit == "") {
    return null;
  }

  const shortCommit = commit?.substring(0, 7);
  const url = "https://github.com/genshinsim/gcsim/commits/" + commit;

  return (
    <Tag large={true} intent="none" minimal={true}>
      <div className="flex flex-row items-center gap-2 font-mono text-xs">
        <div className="text-gray-400">{t<string>("result.metadata_commit")}</div>
        <a href={url} target="_blank" rel="noreferrer" className="font-bold">
          {shortCommit}
        </a>
      </div>
    </Tag>
  );
});