import { TeamTile } from "@gcsim/components";
import { Badge, Button, toast } from "@gcsim/primitives";
import type { db } from "@gcsim/types";
import {
	FaCheck,
	FaExchangeAlt,
	FaExternalLinkAlt,
	FaTimes,
} from "react-icons/fa";
import {
	author,
	created,
	dpsFull,
	mode,
	simTime,
	tagNames,
	team,
	viewerLink,
} from "./lib/entry";

function copy(cmd: string, label: string) {
	navigator.clipboard
		.writeText(cmd)
		.then(() => toast(label, { description: cmd }));
}

export const commands = {
	approve: (id: string) => copy(`/approve id:${id}`, "Copied approve command"),
	reject: (id: string) => copy(`/reject id:${id}`, "Copied reject command"),
	replace: (dupId: string, shareKey: string) =>
		copy(
			`/replace id:${dupId} link:https://gcsim.app/sh/${shareKey}`,
			"Copied replace command",
		),
};

export function ViewerLink({ entry }: { entry: db.Entry }) {
	return (
		<Button size="sm" variant="secondary" asChild>
			<a href={viewerLink(entry)} target="_blank" rel="noreferrer">
				<FaExternalLinkAlt size={11} /> Result viewer
			</a>
		</Button>
	);
}

export function DecisionButtons({ entry }: { entry: db.Entry }) {
	const id = entry._id ?? "";
	return (
		<>
			<Button variant="destructive" onClick={() => commands.reject(id)}>
				<FaTimes size={12} /> Copy reject
			</Button>
			<Button
				className="bg-g-success text-g-accent-fg hover:bg-g-success/90"
				onClick={() => commands.approve(id)}
			>
				<FaCheck size={12} /> Copy approve
			</Button>
			<ViewerLink entry={entry} />
		</>
	);
}

export function MetaChips({ entry }: { entry: db.Entry }) {
	return (
		<div className="flex flex-wrap gap-g-base-sm font-g-mono text-g-xs text-g-ink-dim">
			<Badge className="bg-g-surface-2">mode {mode(entry)}</Badge>
			<Badge className="bg-g-surface-2">{simTime(entry)}</Badge>
			<Badge className="bg-g-surface-2">{created(entry)}</Badge>
			{tagNames(entry).map((t) => (
				<Badge key={t} className="bg-g-success/15 text-g-success">
					{t}
				</Badge>
			))}
		</div>
	);
}

export function MiniTeam({ entry }: { entry: db.Entry }) {
	return (
		<div className="flex -space-x-1">
			{team(entry).map((c, i) => (
				<img
					key={c?.name ?? i}
					src={
						c
							? `/api/assets/avatar/${c.name}.png`
							: "/api/assets/misc/default.png"
					}
					alt={c?.name ?? ""}
					className="size-8 rounded-g-sm border border-g-line-soft bg-g-surface-2 object-cover"
				/>
			))}
		</div>
	);
}

export function DuplicateRow({
	entry,
	main,
}: {
	entry: db.Entry;
	main: db.Entry;
}) {
	return (
		<div className="flex flex-col gap-g-base rounded-g-md border border-g-line-soft bg-g-surface p-g-card md:flex-row md:items-center">
			<div className="flex min-w-0 flex-1 items-center gap-g-base">
				<MiniTeam entry={entry} />
				<div className="min-w-0 flex-1">
					<div className="flex items-center gap-g-base">
						<span className="font-g-mono text-g-sm font-semibold text-g-ink">
							{dpsFull(entry)}
						</span>
						<span className="text-g-xs text-g-ink-mute">
							DPS · {simTime(entry)}
						</span>
					</div>
					<p className="line-clamp-1 text-g-xs text-g-ink-dim">
						by {author(entry)} · {entry.description}
					</p>
				</div>
			</div>
			<div className="flex shrink-0 gap-g-base-sm [&>*]:flex-1 md:[&>*]:flex-none">
				<Button
					size="sm"
					variant="outline"
					onClick={() =>
						commands.replace(entry._id ?? "", main.share_key ?? "")
					}
				>
					<FaExchangeAlt size={11} /> Replace
				</Button>
				<ViewerLink entry={entry} />
			</div>
		</div>
	);
}

export function ReviewHeadline({ entry }: { entry: db.Entry }) {
	return (
		<div className="flex flex-col gap-g-base-sm">
			<div className="w-full max-w-[420px]">
				<TeamTile chars={team(entry)} className="w-full" />
			</div>
			<MetaChips entry={entry} />
			<p className="text-g-sm text-g-ink-dim">
				<span className="font-semibold text-g-accent">{author(entry)}: </span>
				{entry.description}
			</p>
		</div>
	);
}
