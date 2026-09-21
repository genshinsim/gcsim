import tagData from "@gcsim/data/src/tags.json";
import type { db, model } from "@gcsim/types";

const TAGS = tagData as Record<string, { display_name: string }>;

/** The implicit "gcsim" tag carried by every entry; never surfaced as a chip. */
const GCSIM_TAG_ID = 1;

export function team(entry: db.Entry): (model.Character | null)[] {
	const t: (model.Character | null)[] = [...(entry.summary?.team ?? [])];
	while (t.length < 4) t.push(null);
	return t.slice(0, 4);
}

export function dpsFull(entry: db.Entry): string {
	const v = entry.summary?.mean_dps_per_target ?? 0;
	return Math.round(v).toLocaleString();
}

export function mode(entry: db.Entry): string {
	return entry.summary?.mode ? "TTK" : "Duration";
}

export function simTime(entry: db.Entry): string {
	const m = entry.summary?.sim_duration?.mean;
	return m ? `${m.toPrecision(3)}s` : "—";
}

export function created(entry: db.Entry): string {
	const d = entry.create_date;
	if (!d) return "unknown";
	return new Date((d as unknown as number) * 1000).toLocaleDateString();
}

export function author(entry: db.Entry): string {
	return entry.submitter === "migrated"
		? "Unknown author"
		: (entry.submitter ?? "unknown");
}

export function tagNames(entry: db.Entry): string[] {
	return (entry.accepted_tags ?? [])
		.filter((t) => t !== GCSIM_TAG_ID)
		.map((t) => TAGS[String(t)]?.display_name)
		.filter(Boolean) as string[];
}

export function viewerLink(entry: db.Entry): string {
	return `https://gcsim.app/db/${entry._id}`;
}
