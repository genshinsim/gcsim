import { ArrowRight } from "lucide-react";
import type { ReactNode } from "react";

type SectionHeadProps = {
	title: string;
	subtitle: string;
	linkLabel?: string;
	linkHref?: string;
	className?: string;
	children?: ReactNode;
};

// Shared section header used by both Dash layouts: a title + subtitle on the
// left, an optional "see all" link on the right.
export function SectionHead({
	title,
	subtitle,
	linkLabel,
	linkHref,
	className = "",
	children,
}: SectionHeadProps) {
	return (
		<div
			className={`flex items-baseline justify-between gap-3 ${className}`.trim()}
		>
			<div>
				<h2 className="mb-1 font-g-display text-g-h2 font-semibold text-g-ink">
					{title}
				</h2>
				<p className="text-g-body text-g-ink-dim">{subtitle}</p>
			</div>
			{children ??
				(linkLabel ? (
					<a
						href={linkHref}
						target="_blank"
						rel="noreferrer"
						className="inline-flex shrink-0 items-center gap-g-base-sm text-g-sm font-semibold text-g-accent hover:text-g-accent-hover"
					>
						{linkLabel}
						<ArrowRight className="size-3.5" />
					</a>
				) : null)}
		</div>
	);
}
