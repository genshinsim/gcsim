import { model } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { useLocation } from "wouter";

export default function DBEntryTags({
  tags,
}: {
  tags: model.DBTag[] | undefined | null;
}) {
  const { t: translate } = useTranslation();
  const t = (key: string) => translate(key) as React.ReactChild; // idk why this is needed
  const [location, setLocation] = useLocation();

  return (
    <div className={"flex flex-row h-full overflow-hidden max-w-xl"}>
      {tags?.map((tag) => (
        <div
          className="hover:opacity-50 cursor-pointer bg-slate-700 text-xs font-semibold rounded-full px-2 py-1 mr-2 mt-1 whitespace-nowrap "
          key={tag}
          onClick={() => {
            setLocation(`/tag/${tag}`);
          }}
        >
          {
            // https://www.typescriptlang.org/docs/handbook/enums.html search d.ts
            // model.DBTag[tag]
            t(`db.${tag}`)
          }
        </div>
      ))}
    </div>
  );
}
