import { AvatarCard } from "@gcsim/components";
import { Badge, Button, toast } from "@gcsim/primitives";
import type { db } from "@gcsim/types";
import { FaCopy, FaExternalLinkAlt } from "react-icons/fa";
import {
	author,
	created,
	dps,
	mode,
	simTime,
	tagNames,
	targetCount,
	team,
	viewerLink,
} from "../lib/entry";

export function copyConfig(entry: db.Entry) {
	const cfg = entry.config ?? "";
	navigator.clipboard.writeText(cfg).then(() =>
		toast("Copied config", {
			description: `${cfg.length} characters copied to clipboard`,
		}),
	);
}

export function StatChip({ label, value }: { label: string; value: string }) {
	return (
		<Badge className="gap-1.5 bg-g-surface-2 font-g-mono">
			<span className="text-g-xs lowercase text-g-ink-mute">{label}</span>
			<span className="text-g-xs text-g-ink">{value}</span>
		</Badge>
	);
}

export function TagBadges({ entry }: { entry: db.Entry }) {
	const tags = tagNames(entry);
	if (!tags.length) return null;
	return (
		<div className="flex flex-wrap gap-g-base-sm">
			{tags.map((t) => (
				<Badge key={t} className="bg-g-success/15 font-g-mono">
					<span className="text-g-xs text-g-success">{t}</span>
				</Badge>
			))}
		</div>
	);
}

export function CardActions({ entry }: { entry: db.Entry }) {
	return (
		<div className="flex flex-wrap gap-g-base-sm">
			<Button size="sm" variant="secondary" onClick={() => copyConfig(entry)}>
				<FaCopy size={12} /> Copy config
			</Button>
			<Button size="sm" asChild>
				<a href={viewerLink(entry)} target="_blank" rel="noreferrer">
					<FaExternalLinkAlt size={11} /> Open in viewer
				</a>
			</Button>
		</div>
	);
}

export function Team({ entry }: { entry: db.Entry }) {
	return (
		<div className="w-full max-w-[420px]">
			<AvatarCard chars={team(entry)} className="w-full" />
		</div>
	);
}

export function FullCard({ entry }: { entry: db.Entry }) {
	return (
		<div className="flex flex-col gap-g-base rounded-g-lg border border-g-line-soft bg-g-surface p-g-card md:flex-row md:items-stretch">
			<div className="w-full shrink-0 md:w-[420px]">
				<Team entry={entry} />
			</div>
			<div className="flex min-w-0 flex-1 flex-col gap-g-base">
				<div className="flex flex-wrap items-center gap-g-base-sm">
					<StatChip label="mode" value={mode(entry)} />
					<StatChip label="targets" value={String(targetCount(entry))} />
					<StatChip label="dps/target" value={dps(entry)} />
					<StatChip label="avg time" value={simTime(entry)} />
					<StatChip label="date" value={created(entry)} />
				</div>
				<TagBadges entry={entry} />
				<p className="text-g-sm text-g-ink-dim">
					<span className="font-semibold text-g-accent">{author(entry)}: </span>
					{entry.description}
				</p>
				<div className="mt-auto flex justify-end pt-g-base-sm">
					<CardActions entry={entry} />
				</div>
			</div>
		</div>
	);
}
