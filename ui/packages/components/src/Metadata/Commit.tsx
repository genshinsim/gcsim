import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Badge } from "../common/ui/badge";
import { cn } from "../lib/utils";

type Props = {
  commit?: string;
  className?: string;
};

export const Commit = memo(({ commit, className }: Props) => {
  const { t } = useTranslation();
  if (commit == null || commit == "") {
    return null;
  }

  const shortCommit = commit?.substring(0, 7);
  const url = "https://github.com/genshinsim/gcsim/commits/" + commit;

  const cc = cn("text-sm font-mono", className);

  return (
    <Badge className={cc}>
      <div className="flex flex-row items-center gap-2">
        <div className="text-gray-400">
          {t<string>("result.metadata_commit")}
        </div>
        <a href={url} target="_blank" rel="noreferrer" className="font-bold">
          {shortCommit}
        </a>
      </div>
    </Badge>
  );
});
