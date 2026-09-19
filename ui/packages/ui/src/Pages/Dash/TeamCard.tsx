import type { db, model } from "@gcsim/types";

type TeamCardProps = {
	entry: db.Entry;
	className?: string;
};

export function TeamCard({ entry, className }: TeamCardProps) {
	const roster = entry.summary?.team ?? [];
	const slots = Array.from({ length: 4 }, (_, i) => {
		const char = roster[i] ?? null;
		return { key: char?.name ?? `empty-${i}`, char };
	});

	const names = roster
		.filter((c): c is model.Character => c != null)
		.map((c) => c.name);
	const title = entry.description || (names.length ? names.join(", ") : "Team");
	const submitter =
		entry.submitter === "migrated" ? "unknown" : entry.submitter;
	const dps = (entry.summary?.mean_dps_per_target ?? 0).toLocaleString(
		navigator.language,
		{
			notation: "compact",
			minimumSignificantDigits: 3,
			maximumSignificantDigits: 3,
		},
	);

	return (
		<a
			href={`https://gcsim.app/db/${entry._id}`}
			target="_blank"
			rel="noreferrer"
			className={
				"block rounded-g-lg border border-g-line bg-g-surface p-g-card shadow-g-card transition-colors hover:border-g-accent " +
				(className ?? "")
			}
		>
			<div className="flex flex-col">
				<div className="mb-3 flex gap-g-base-sm">
					{slots.map(({ key, char }) =>
						char ? (
							<img
								key={key}
								src={`/api/assets/avatar/${char.name}.png`}
								alt=""
								className="h-[38px] w-[38px] rounded-g-md border border-g-line bg-g-surface-2 object-cover object-top"
							/>
						) : (
							<div
								key={key}
								className="h-[38px] w-[38px] rounded-g-md border border-g-line bg-g-surface-2"
							/>
						),
					)}
				</div>
				<h3 className="mb-0.5 line-clamp-1 font-g-display text-g-h3 font-semibold text-g-ink">
					{title}
				</h3>
				<p className="mb-3 text-g-sm text-g-ink-mute">by {submitter}</p>
				<div className="flex items-end justify-between">
					<div>
						<div className="font-g-mono text-g-num-sm font-bold text-g-ink">
							{dps}
						</div>
						<div className="text-g-xs text-g-ink-mute">DPS / target</div>
					</div>
					<span className="inline-flex items-center rounded-g-pill border border-g-line-soft bg-g-surface-2 px-2.5 py-1 text-g-xs font-semibold text-g-ink-dim">
						mode {entry.summary?.mode ? "TTK" : "duration"}
					</span>
				</div>
			</div>
		</a>
	);
}
