import { Spinner } from "@gcsim/primitives";
import type { db } from "@gcsim/types";
import { FullCard } from "./EntryCard";

export function ListView({ data }: { data: db.Entry[] }) {
	if (!data) {
		return (
			<div>
				<Spinner className="size-8" />
			</div>
		);
	}

	return (
		<div className="flex flex-col gap-g-base-lg">
			{data.map((entry) => (
				<FullCard key={entry._id} entry={entry} />
			))}
		</div>
	);
}
