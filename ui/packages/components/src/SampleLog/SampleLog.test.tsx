import { fireEvent, render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type React from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("react-i18next", () => ({
	useTranslation: () => ({
		t: (k: string, o?: Record<string, unknown>) =>
			o == null ? k : `${k} ${JSON.stringify(o)}`,
	}),
	Trans: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

import { SampleLog } from "./SampleLog";
import { AllSampleOptions } from "./SampleOptions";
import { loadSampleFilter } from "./sampleFilter";
import { loadSampleView, saveSampleView } from "./sampleView";
import { sampleFixture as sample } from "./testdata";
import { columns } from "./timeline/layout";
import { modelFromSample } from "./timeline/model";

const scrolls: number[] = [];
const vscrolls: number[] = [];

beforeEach(() => {
	localStorage.clear();
	scrolls.length = 0;
	vscrolls.length = 0;
	Element.prototype.scrollTo = function (
		this: Element,
		opts?: ScrollToOptions | number,
	) {
		if (typeof opts === "object" && opts.top != null) {
			vscrolls.push(opts.top);
			this.scrollTop = opts.top;
			this.dispatchEvent(new Event("scroll"));
		} else if (typeof opts === "object") {
			scrolls.push(opts.left ?? 0);
			this.scrollLeft = opts.left ?? 0;
			this.dispatchEvent(new Event("scroll"));
		}
	} as typeof Element.prototype.scrollTo;
});

function renderLog(
	props: Partial<React.ComponentProps<typeof SampleLog>> = {},
) {
	return render(<SampleLog sample={sample} {...props} />);
}

const chipsOfType = (type: string) =>
	screen.queryAllByTitle(new RegExp(`^\\d+ · ${type}: `));

describe("SampleLog", () => {
	it("draws a lane per character plus the sim lane", () => {
		renderLog();
		for (const c of sample.character_details ?? []) {
			expect(
				screen.getByTitle(`game:character_names.${c.name}`),
			).toBeInTheDocument();
		}
		expect(screen.getAllByTitle("sample.sim_lane").length).toBeGreaterThan(0);
	});

	it("steps through search matches, centring the strip on each", async () => {
		renderLog();
		await userEvent.type(
			screen.getByPlaceholderText("sample.search_placeholder"),
			"Sesshou",
		);
		const next = screen.getByRole("button", { name: "sample.next_match" });
		await userEvent.click(next);
		await userEvent.click(next);
		expect(scrolls).toHaveLength(2);
		expect(screen.getByText(/^2\/\d+$/)).toBeInTheDocument();
	});

	it("steps the strip a column with the arrow keys from anywhere on the page", async () => {
		renderLog();
		await userEvent.keyboard("{ArrowRight}");
		expect(scrolls).toHaveLength(1);
		expect(scrolls[0]).toBeGreaterThan(0);
		await userEvent.keyboard("{ArrowLeft}");
		expect(scrolls).toEqual([scrolls[0], 0]);
	});

	it("pans the strip six columns for each wheel notch on the overview", () => {
		const { container } = renderLog();
		const overview = container.querySelector("canvas")?.parentElement;
		if (overview == null) {
			throw new Error("no overview");
		}
		for (let i = 0; i < 4; i++) {
			fireEvent.wheel(overview, { deltaY: 100 });
		}
		const cols = columns(modelFromSample(sample), new Set(AllSampleOptions));
		expect(scrolls).toEqual([6, 12, 18, 24].map((i) => cols[i].x));
	});

	it("leaves arrow keys to the search box", async () => {
		renderLog();
		await userEvent.click(
			screen.getByPlaceholderText("sample.search_placeholder"),
		);
		await userEvent.keyboard("{ArrowRight}{ArrowLeft}");
		expect(scrolls).toEqual([]);
	});

	it("leaves arrow keys to an open frame view", async () => {
		renderLog();
		await userEvent.click(
			screen.getAllByRole("button", { name: /^sample\.open_frame / })[0],
		);
		await userEvent.keyboard("{ArrowRight}");
		// the frame view centres the strip on its new frame, and nothing else scrolls it
		expect(
			within(screen.getByRole("dialog")).getByText(/sample\.frame_title/),
		).toBeInTheDocument();
		expect(scrolls).toHaveLength(1);
	});

	it("restarts the match count when the filter changes the matches", async () => {
		renderLog();
		await userEvent.type(
			screen.getByPlaceholderText("sample.search_placeholder"),
			"Sesshou",
		);
		const next = screen.getByRole("button", { name: "sample.next_match" });
		await userEvent.click(next);
		await userEvent.click(next);
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.log_options" }),
		);
		await userEvent.click(
			within(screen.getByRole("dialog")).getByText("viewer.clear"),
		);
		expect(screen.getByText("0/0")).toBeInTheDocument();
	});

	it("offers download only when a handler is given", async () => {
		const { unmount } = renderLog();
		expect(
			screen.queryByRole("button", { name: "viewer.download" }),
		).toBeNull();
		unmount();

		const onDownload = vi.fn();
		renderLog({ onDownload });
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.download" }),
		);
		expect(onDownload).toHaveBeenCalledExactlyOnceWith(sample);
	});

	it("offers generate beside the filter only when a handler is given", async () => {
		const { unmount } = renderLog();
		expect(
			screen.queryByRole("button", { name: "viewer.generate" }),
		).toBeNull();
		unmount();

		const onGenerate = vi.fn();
		renderLog({ onGenerate });
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.generate" }),
		);
		expect(onGenerate).toHaveBeenCalledOnce();
	});

	it("offers only generate when there is no sample", async () => {
		const onGenerate = vi.fn();
		renderLog({ sample: null, onGenerate });
		expect(screen.queryByTitle("sample.sim_lane")).toBeNull();
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.generate" }),
		);
		expect(onGenerate).toHaveBeenCalledOnce();
	});

	it("saves the filter and hides the types it leaves out", async () => {
		renderLog();
		expect(chipsOfType("action").length).toBeGreaterThan(0);
		await userEvent.click(
			screen.getByRole("button", { name: "viewer.log_options" }),
		);
		const popover = screen.getByRole("dialog");
		await userEvent.click(within(popover).getByText("action"));
		expect(loadSampleFilter(localStorage)).not.toContain("action");
		expect(chipsOfType("action")).toEqual([]);
	});

	it("leads a damage chip with the compact amount and keeps the exact one in its title", async () => {
		renderLog();
		await userEvent.type(
			screen.getByPlaceholderText("sample.search_placeholder"),
			"Sesshou",
		);
		await userEvent.click(
			screen.getByRole("button", { name: "sample.next_match" }),
		);
		const chip = screen.getByTitle(
			"122 · damage: Sesshou Sakura Tick [8,815] (crit)",
		);
		expect(chip).toHaveTextContent(/8\.82K Sesshou Sakura Tick \(crit\)$/);
	});

	it("opens an event's details from its chip", async () => {
		renderLog();
		const chip = chipsOfType("action")[0];
		await userEvent.click(chip);
		const heading = chip.title.replace(" · action: ", " · ");
		expect(
			within(screen.getByRole("dialog")).getByRole("heading", {
				name: heading,
			}),
		).toBeInTheDocument();
	});

	it("keeps a duration highlight when Escape closes the details dialog", async () => {
		renderLog();
		const durations = screen.getAllByRole("button", {
			name: /^sample\.show_duration /,
		});
		await userEvent.click(durations[0]);
		const clear = () => screen.queryAllByLabelText(/^sample\.clear_duration /);
		expect(clear()).toHaveLength(1);

		await userEvent.click(chipsOfType("action")[0]);
		await userEvent.keyboard("{Escape}");
		expect(screen.queryByRole("dialog")).toBeNull();
		expect(clear()).toHaveLength(1);

		await userEvent.keyboard("{Escape}");
		expect(clear()).toHaveLength(0);
	});
});

describe("SampleLog vertical", () => {
	const vertical = (over: Partial<ReturnType<typeof loadSampleView>> = {}) =>
		saveSampleView(localStorage, {
			...loadSampleView(localStorage, true),
			orientation: "vertical",
			...over,
		});

	const openView = async () => {
		await userEvent.click(screen.getByRole("button", { name: "sample.view" }));
		return screen.getByRole("dialog");
	};

	it("switches to vertical from the view menu and remembers it", async () => {
		renderLog();
		const menu = await openView();
		await userEvent.click(
			within(menu).getByRole("button", { name: "sample.vertical" }),
		);
		expect(loadSampleView(localStorage, true).orientation).toBe("vertical");
		expect(
			screen.getByRole("button", { name: "sample.hide_overview_short" }),
		).toBeInTheDocument();
	});

	it("caps slim rows and expands one on +N more", async () => {
		vertical({ slim: true, cap: 3 });
		renderLog();
		const more = screen.getAllByRole("button", { name: /^sample\.more / });
		expect(more.length).toBeGreaterThan(0);
		await userEvent.click(more[0]);
		expect(
			screen.getByRole("button", { name: /sample\.less/ }),
		).toBeInTheDocument();
	});

	it("shows every chip, wrapped, when slim chips are off", async () => {
		vertical({ slim: false });
		renderLog();
		expect(screen.queryAllByRole("button", { name: /^sample\.more / })).toEqual(
			[],
		);
		const menu = await openView();
		expect(within(menu).getByRole("slider")).toBeDisabled();
	});

	it("only offers the chip options to the vertical layout", async () => {
		renderLog();
		const menu = await openView();
		expect(within(menu).getByRole("switch")).toBeDisabled();
		expect(within(menu).getByRole("slider")).toBeDisabled();
	});

	it("saves the chip cap from the slider", async () => {
		vertical({ slim: true, cap: 4 });
		renderLog();
		const menu = await openView();
		fireEvent.change(within(menu).getByRole("slider"), {
			target: { value: "8" },
		});
		expect(loadSampleView(localStorage, true).cap).toBe(8);
	});

	it("scrolls down to each search match", async () => {
		vertical();
		renderLog();
		await userEvent.type(
			screen.getByPlaceholderText("sample.search_placeholder"),
			"Sesshou",
		);
		const next = screen.getByRole("button", { name: "sample.next_match" });
		await userEvent.click(next);
		await userEvent.click(next);
		expect(vscrolls).toHaveLength(2);
		expect(vscrolls[0]).toBeGreaterThan(0);
		expect(scrolls).toEqual([]);
	});

	it("collapses the overview to a rail and remembers it", async () => {
		vertical({ overview: true });
		const { container } = renderLog();
		expect(container.querySelector("canvas")).not.toBeNull();
		await userEvent.click(
			screen.getByRole("button", { name: "sample.hide_overview_short" }),
		);
		expect(container.querySelector("canvas")).toBeNull();
		expect(loadSampleView(localStorage, true).overview).toBe(false);
	});
});
