import type { Meta, StoryObj } from "@storybook/react-vite";

const PALETTES = [
	"cryo",
	"abyss-d",
	"abyss-l",
	"ember-d",
	"ember-l",
	"pulse-d",
	"pulse-l",
	"twilight",
	"qilin",
	"glacier",
	"azure",
	"aqua",
	"frostfall",
	"blossom",
	"lantern",
] as const;

type Palette = (typeof PALETTES)[number];

const labelStyle: React.CSSProperties = {
	textTransform:
		"var(--g-lbl-transform)" as React.CSSProperties["textTransform"],
	letterSpacing: "var(--g-lbl-track)",
	fontWeight: "var(--g-lbl-weight)" as React.CSSProperties["fontWeight"],
	fontFamily: "var(--g-lbl-font)",
};

const Label = ({ children }: { children: React.ReactNode }) => (
	<div className="text-g-ink-mute text-g-xs" style={labelStyle}>
		{children}
	</div>
);

const Swatch = ({ className, name }: { className: string; name: string }) => (
	<div className="flex flex-col gap-g-base-sm">
		<div
			className={`h-12 rounded-g-md border border-g-line ${className}`}
			style={{ borderWidth: "var(--g-bw)" }}
		/>
		<span className="text-g-ink-dim text-g-xs font-g-mono">{name}</span>
	</div>
);

const Foundation = ({ palette }: { palette: Palette }) => (
	<div
		data-theme={palette}
		className="bg-g-canvas text-g-ink font-g-body p-g-page flex flex-col gap-g-section rounded-g-lg"
	>
		<section className="flex flex-col gap-g-base">
			<Label>Surfaces</Label>
			<div className="grid grid-cols-2 gap-g-base sm:grid-cols-4">
				<Swatch className="bg-g-canvas" name="canvas" />
				<Swatch className="bg-g-surface" name="surface" />
				<Swatch className="bg-g-surface-2" name="surface-2" />
				<Swatch className="bg-g-surface-3" name="surface-3" />
			</div>
		</section>

		<section className="flex flex-col gap-g-base">
			<Label>Type ramp</Label>
			<div className="bg-g-surface p-g-card rounded-g-lg flex flex-col gap-g-base-sm">
				<div className="font-g-display text-g-hero">Hero 48</div>
				<div className="font-g-display text-g-h1">Heading 1 · 30</div>
				<div className="font-g-display text-g-h2">Heading 2 · 22</div>
				<div className="font-g-display text-g-h3">Heading 3 · 15</div>
				<div className="text-g-lg">Large body · 16</div>
				<div className="text-g-body">Body · 15 — the quick brown fox</div>
				<div className="text-g-sm text-g-ink-dim">Small · 13 — dimmed</div>
				<div className="text-g-xs text-g-ink-mute">
					Extra small · 11.5 — muted
				</div>
				<div className="font-g-mono text-g-num">1,234.5</div>
				<div className="font-g-mono text-g-num-sm text-g-ink-dim">42.0</div>
			</div>
		</section>

		<section className="flex flex-col gap-g-base">
			<Label>Accent &amp; control</Label>
			<div className="flex flex-wrap items-center gap-g-base">
				<button
					type="button"
					className="bg-g-accent text-g-accent-fg font-g-body h-g-ctrl px-g-btn-x rounded-g-btn shadow-g-pop"
				>
					Run
				</button>
				<button
					type="button"
					className="bg-g-accent-weak text-g-accent font-g-body h-g-ctrl px-g-btn-x rounded-g-btn border border-g-line"
					style={{ borderWidth: "var(--g-bw)" }}
				>
					Secondary
				</button>
				<span className="text-g-success text-g-sm">success</span>
				<span className="text-g-warning text-g-sm">warning</span>
				<span className="text-g-danger text-g-sm">danger</span>
			</div>
		</section>

		<section className="flex flex-col gap-g-base">
			<Label>Chart series</Label>
			<div className="grid grid-cols-3 gap-g-base sm:grid-cols-6">
				<Swatch className="bg-g-s1" name="s1" />
				<Swatch className="bg-g-s2" name="s2" />
				<Swatch className="bg-g-s3" name="s3" />
				<Swatch className="bg-g-s4" name="s4" />
				<Swatch className="bg-g-s5" name="s5" />
				<Swatch className="bg-g-s6" name="s6" />
			</div>
		</section>

		<section className="flex flex-col gap-g-base">
			<Label>Elements</Label>
			<div className="grid grid-cols-4 gap-g-base sm:grid-cols-7">
				<Swatch className="bg-g-anemo" name="anemo" />
				<Swatch className="bg-g-geo" name="geo" />
				<Swatch className="bg-g-electro" name="electro" />
				<Swatch className="bg-g-hydro" name="hydro" />
				<Swatch className="bg-g-pyro" name="pyro" />
				<Swatch className="bg-g-cryo" name="cryo" />
				<Swatch className="bg-g-dendro" name="dendro" />
			</div>
		</section>
	</div>
);

const meta: Meta<typeof Foundation> = {
	title: "Theme/Gauge Foundation",
	component: Foundation,
	args: { palette: "cryo" },
	argTypes: {
		palette: { control: "select", options: PALETTES },
	},
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Cryo: Story = {};

export const Frostfall: Story = {
	args: { palette: "frostfall" },
};
