import { model } from "@gcsim/types";
import { useTranslation } from "react-i18next";
import { useLocation } from "wouter";

const defaultTranslations = {
  DB_TAG_INVALID: "Invalid",

  DB_TAG_GCSIM: "",

  DB_TAG_TESTING: "Internal Test",

  DB_TAG_KQM_GUIDE: "KQM Guide",

  DB_TAG_GEO_SIMPS: "Geo",

  DB_TAG_ITTO_SIMPS: "Itto Simps",
};

export default function DBEntryTags({
  tags,
}: {
  tags: model.DBTag[] | undefined | null;
}) {
  const { t: translate } = useTranslation();
  const t = (key: string) => translate(key) as string; // idk why this is needed
  const [location, setLocation] = useLocation();

  return (
    <div className={"flex flex-row h-full overflow-hidden max-w-xl "}>
      {tags
        ?.filter((tag) => (tag as unknown as string) !== "DB_TAG_GCSIM")
        .map((tag) => (
          <div
            className="hover:opacity-50 cursor-pointer bg-slate-700 text-xs font-semibold rounded-full px-2 py-1 mr-2 mt-1 whitespace-nowrap "
            key={tag}
            onClick={() => {
              setLocation(`/tag/${tag}`);
            }}
          >
            {defaultTranslations[tag]}
          </div>
        ))}
    </div>
  );
}
