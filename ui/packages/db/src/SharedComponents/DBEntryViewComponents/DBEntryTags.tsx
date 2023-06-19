import { model } from "@gcsim/types";
import { useState } from "react";
import { dbTag } from "./tag";

export default function DBEntryTags({
  tags,
}: {
  tags: model.DBTag[] | undefined | null;
}) {
  const t = (key: string) => key;
  return (
    <div
      className={
        "flex flex-row h-full overflow-hidden max-w-xl"
      }
    >
      {tags?.map((tag) => (
        <div
          className="bg-slate-700 text-xs font-semibold rounded-full px-2 py-1 mr-2 mt-1 whitespace-nowrap "
          key={tag}
        >
          {
            // https://www.typescriptlang.org/docs/handbook/enums.html search d.ts
            // model.DBTag[tag]
          tag
          }
        </div>
      ))}
      {tags?.map((tag) => (
        <div
          className="bg-slate-700 text-xs font-semibold rounded-full px-2 py-1 mr-2 mt-1 whitespace-nowrap "
          key={tag}
        >
           
           {tag}
        </div>
      ))}
     
    </div>
  );
}
