import { cn } from "@gcsim/primitives";
import { type db, model } from "@gcsim/types";
import type { CSSProperties } from "react";
import { useTranslation } from "react-i18next";
import { KQM_DB_URL } from "./kqm";

// KQM DB lavender mixed into the active theme's own tokens, so it holds in
// light and dark themes alike.
const KQM_TINT = "#b07ae0";
const kqmStyle = {
	"--kqm-bg": `color-mix(in oklab, var(--g-surface) 88%, ${KQM_TINT})`,
	"--kqm-bg-hover": `color-mix(in oklab, var(--g-surface) 82%, ${KQM_TINT})`,
	"--kqm-line": `color-mix(in oklab, var(--g-border) 55%, ${KQM_TINT})`,
	"--kqm-line-hover": `color-mix(in oklab, var(--g-border) 25%, ${KQM_TINT})`,
	"--kqm-ink": `color-mix(in oklab, ${KQM_TINT} 55%, var(--g-text))`,
} as CSSProperties;

type TeamCardProps = {
	entry: db.Entry;
	className?: string;
};

export function TeamCard({ entry, className }: TeamCardProps) {
	const { t } = useTranslation();

	const roster = entry.summary?.team ?? [];
	const slots = Array.from({ length: 4 }, (_, i) => {
		const char = roster[i] ?? null;
		return { key: char?.name ?? `empty-${i}`, char };
	});

	const names = roster
		.filter((c): c is model.Character => c != null)
		.map((c) => c.name);
	const title =
		entry.description ||
		(names.length ? names.join(", ") : t("dash.team_fallback"));
	const submitter =
		entry.submitter === "migrated" ? t("dash.author_unknown") : entry.submitter;
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
			href={`${KQM_DB_URL}/db/${entry._id}`}
			target="_blank"
			rel="noreferrer"
			style={kqmStyle}
			className={cn(
				"block rounded-g-lg border border-(--kqm-line) bg-(--kqm-bg) p-g-card shadow-g-card transition-colors hover:border-(--kqm-line-hover) hover:bg-(--kqm-bg-hover)",
				className,
			)}
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
				<p className="mb-3 text-g-sm text-g-ink-mute">
					{t("dash.by_author", { author: submitter })}
				</p>
				<div className="flex items-end justify-between">
					<div>
						<div className="font-g-mono text-g-num-sm font-bold text-g-ink">
							{dps}
						</div>
						<div className="text-g-xs text-g-ink-mute">
							{t("dash.dps_per_target")}
						</div>
					</div>
					<span className="inline-flex items-center rounded-g-pill border border-(--kqm-line) bg-(--kqm-bg-hover) px-2.5 py-1 text-g-xs font-semibold text-(--kqm-ink)">
						{t("dash.mode_label", {
							mode:
								entry.summary?.mode === model.SimMode.TTK_MODE
									? t("dash.mode_ttk")
									: t("dash.mode_duration"),
						})}
					</span>
				</div>
			</div>
		</a>
	);
}
