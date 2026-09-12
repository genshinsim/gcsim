import tagData from "@gcsim/data/src/tags.json";
import type { model } from "@gcsim/types";
import { Link } from "wouter";

export default function DBEntryTags({
	tags,
}: {
	tags: model.DBTag[] | undefined | null;
}) {
	// const { t: translate } = useTranslation();
	// const t = (key: string) => translate(key) as string; // idk why this is needed

	return (
		<div className={"flex flex-row overflow-hidden"}>
			{tags
				?.filter((tag) => tag !== 1)
				.map((tag) => (
					<Link
						className="hover:opacity-50 cursor-pointer bg-slate-500 text-xs font-semibold rounded-full px-2 py-1 mr-2 mt-1 whitespace-nowrap "
						key={tag}
						href={`/tag/${tag}`}
					>
						{tagData[tag].display_name}
					</Link>
				))}
		</div>
	);
}
